package handlers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"happy-api/internal/config"
	"happy-api/internal/http/mw"
	"happy-api/internal/models"
	"happy-api/internal/payments"
	"happy-api/internal/payments/robokassa"
	"happy-api/internal/repo"
	"happy-api/internal/security"
	"happy-api/internal/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	cfg  config.Config
	repo *repo.Repo
	pay  payments.Provider
	robo *robokassa.Client
}

func (h *Handler) SetRobokassa(c *robokassa.Client) { h.robo = c }

func New(cfg config.Config, r *repo.Repo, pay payments.Provider) *Handler {
	return &Handler{cfg: cfg, repo: r, pay: pay}
}

func userPublic(u models.User) gin.H {
	return gin.H{
		"id":       u.ID.Hex(),
		"email":    u.Email,
		"username": u.Username,
		"role":     u.Role,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var in struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_JSON"})
		return
	}

	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.Username = strings.TrimSpace(in.Username)

	if !utils.IsValidEmail(in.Email) || !utils.IsValidUsername(in.Username) || len(in.Password) < 8 {
		c.JSON(400, gin.H{"ok": false, "error": "INVALID_INPUT"})
		return
	}

	pwHash, err := security.HashPassword(in.Password)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "HASH_FAIL"})
		return
	}

	recoveryPlain := security.RandomCode(10)
	recHash, _ := security.HashPassword(recoveryPlain)

	u, err := h.repo.CreateUser(c, in.Email, in.Username, pwHash, recHash)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(409, gin.H{"ok": false, "error": "ALREADY_EXISTS"})
			return
		}
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}

	token, err := security.SignJWT(h.cfg.JWTSecret, u.ID.Hex(), string(u.Role), 30*24*time.Hour)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "JWT_FAIL"})
		return
	}
	h.setAuthCookie(c, token, 7*24*time.Hour)

	c.JSON(200, gin.H{
		"ok":           true,
		"user":         userPublic(u),
		"recoveryCode": recoveryPlain, // показываем один раз сразу после регистрации
	})
}

func (h *Handler) Login(c *gin.Context) {
	var in struct {
		LoginOrEmail string `json:"loginOrEmail"`
		Password     string `json:"password"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_JSON"})
		return
	}
	in.LoginOrEmail = strings.ToLower(strings.TrimSpace(in.LoginOrEmail))

	u, err := h.repo.FindUserByLogin(c, in.LoginOrEmail)
	if err != nil {
		c.JSON(401, gin.H{"ok": false, "error": "INVALID_CREDENTIALS"})
		return
	}
	if !security.CheckPassword(u.PasswordHash, in.Password) {
		c.JSON(401, gin.H{"ok": false, "error": "INVALID_CREDENTIALS"})
		return
	}
	token, err := security.SignJWT(h.cfg.JWTSecret, u.ID.Hex(), string(u.Role), 30*24*time.Hour)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "JWT_FAIL"})
		return
	}
	h.setAuthCookie(c, token, 7*24*time.Hour)
	c.JSON(200, gin.H{"ok": true, "user": userPublic(u)})
}

func (h *Handler) Logout(c *gin.Context) {
	h.clearAuthCookie(c)
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) Me(c *gin.Context) {
	uid := mw.UserID(c)
	if uid == primitive.NilObjectID {
		c.JSON(401, gin.H{"ok": false, "error": "UNAUTHORIZED"})
		return
	}

	u, err := h.repo.FindUserByID(c, uid)
	if err != nil {
		c.JSON(200, gin.H{"ok": true, "user": nil})
		return
	}

	c.JSON(200, gin.H{"ok": true, "user": userPublic(u)})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var in struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_JSON"})
		return
	}
	if len(in.NewPassword) < 8 {
		c.JSON(400, gin.H{"ok": false, "error": "WEAK_PASSWORD"})
		return
	}

	uid := mw.UserID(c)
	u, err := h.repo.FindUserByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(401, gin.H{"ok": false, "error": "UNAUTHORIZED"})
		return
	}
	if !security.CheckPassword(u.PasswordHash, in.OldPassword) {
		c.JSON(400, gin.H{"ok": false, "error": "WRONG_OLD_PASSWORD"})
		return
	}
	pwHash, _ := security.HashPassword(in.NewPassword)
	if err := h.repo.UpdateUserPasswordAndRecovery(c, uid, pwHash, u.RecoveryCodeHash, u.RecoveryCodeUsedAt); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) RecoverReset(c *gin.Context) {
	var in struct {
		LoginOrEmail string `json:"loginOrEmail"`
		RecoveryCode string `json:"recoveryCode"`
		NewPassword  string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_JSON"})
		return
	}
	in.LoginOrEmail = strings.ToLower(strings.TrimSpace(in.LoginOrEmail))
	in.RecoveryCode = strings.TrimSpace(in.RecoveryCode)
	if len(in.NewPassword) < 8 || in.RecoveryCode == "" {
		c.JSON(400, gin.H{"ok": false, "error": "INVALID_INPUT"})
		return
	}

	u, err := h.repo.FindUserByLogin(c, in.LoginOrEmail)
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "INVALID_RECOVERY"})
		return
	}

	// одноразовый код
	if u.RecoveryCodeUsedAt != nil {
		c.JSON(400, gin.H{"ok": false, "error": "RECOVERY_ALREADY_USED"})
		return
	}
	if !security.CheckPassword(u.RecoveryCodeHash, in.RecoveryCode) {
		c.JSON(400, gin.H{"ok": false, "error": "INVALID_RECOVERY"})
		return
	}

	pwHash, _ := security.HashPassword(in.NewPassword)
	newRecovery := security.RandomCode(10)
	newRecHash, _ := security.HashPassword(newRecovery)
	now := time.Now()
	if err := h.repo.UpdateUserPasswordAndRecovery(c, u.ID, pwHash, newRecHash, &now); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}

	token, _ := security.SignJWT(h.cfg.JWTSecret, u.ID.Hex(), string(u.Role), 30*24*time.Hour)
	h.setAuthCookie(c, token, 7*24*time.Hour)

	c.JSON(200, gin.H{"ok": true, "newRecoveryCode": newRecovery})
}

func (h *Handler) RecoverRegenerate(c *gin.Context) {
	uid := mw.UserID(c)
	u, err := h.repo.FindUserByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(401, gin.H{"ok": false, "error": "UNAUTHORIZED"})
		return
	}
	newRecovery := security.RandomCode(10)
	newRecHash, _ := security.HashPassword(newRecovery)
	if err := h.repo.UpdateUserPasswordAndRecovery(c, u.ID, u.PasswordHash, newRecHash, nil); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "recoveryCode": newRecovery})
}

func (h *Handler) ListGreetings(c *gin.Context) {
	uid := mw.UserID(c)
	items, err := h.repo.ListGreetingsByOwner(c, uid)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "items": items})
}

func (h *Handler) CreateGreeting(c *gin.Context) {
	var in struct {
		Type      models.GreetingType `json:"type"`
		Subdomain string              `json:"subdomain"`
		Title     string              `json:"title"`
		Body      string              `json:"body"`
		Theme     models.Theme        `json:"theme"`
		Gift      models.Gift         `json:"gift"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_JSON"})
		return
	}
	if in.Type != models.GreetingCode && in.Type != models.GreetingSubdomain {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_TYPE"})
		return
	}
	if in.Type == models.GreetingSubdomain {
		in.Subdomain = utils.NormalizeSubdomain(in.Subdomain)
		if in.Subdomain != "" && !utils.IsValidSubdomainLabel(in.Subdomain) {
			c.JSON(400, gin.H{"ok": false, "error": "BAD_SUBDOMAIN"})
			return
		}
	}

	uid := mw.UserID(c)
	g, err := h.repo.CreateGreeting(c, uid, models.Greeting{
		Type:      in.Type,
		Subdomain: in.Subdomain,
		Title:     strings.TrimSpace(in.Title),
		Body:      strings.TrimSpace(in.Body),
		Theme:     in.Theme,
		Gift:      in.Gift,
	})
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(409, gin.H{"ok": false, "error": "DUPLICATE"})
			return
		}
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "item": g})
}

func (h *Handler) GetGreeting(c *gin.Context) {
	uid := mw.UserID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}
	g, err := h.repo.GetGreetingOwned(c, id, uid)
	if err != nil {
		c.JSON(404, gin.H{"ok": false, "error": "NOT_FOUND"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "item": g})
}

func (h *Handler) UpdateGreeting(c *gin.Context) {
	uid := mw.UserID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}

	var in struct {
		Subdomain *string       `json:"subdomain"`
		Title     *string       `json:"title"`
		Body      *string       `json:"body"`
		Theme     *models.Theme `json:"theme"`
		Gift      *models.Gift  `json:"gift"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_JSON"})
		return
	}

	patch := bson.M{}
	if in.Subdomain != nil {
		sub := utils.NormalizeSubdomain(*in.Subdomain)
		if sub != "" && !utils.IsValidSubdomainLabel(sub) {
			c.JSON(400, gin.H{"ok": false, "error": "BAD_SUBDOMAIN"})
			return
		}
		patch["subdomain"] = sub
	}
	if in.Title != nil {
		patch["title"] = strings.TrimSpace(*in.Title)
	}
	if in.Body != nil {
		patch["body"] = strings.TrimSpace(*in.Body)
	}
	if in.Theme != nil {
		patch["theme"] = *in.Theme
	}
	if in.Gift != nil {
		patch["gift"] = *in.Gift
	}

	if len(patch) == 0 {
		c.JSON(200, gin.H{"ok": true})
		return
	}

	if err := h.repo.UpdateGreetingOwned(c, id, uid, patch); err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.JSON(409, gin.H{"ok": false, "error": "DUPLICATE"})
			return
		}
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) UploadPhotos(c *gin.Context) {
	uid := mw.UserID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}
	// ensure ownership
	g, err := h.repo.GetGreetingOwned(c, id, uid)
	if err != nil {
		c.JSON(404, gin.H{"ok": false, "error": "NOT_FOUND"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil || form.File == nil {
		c.JSON(400, gin.H{"ok": false, "error": "NO_FILES"})
		return
	}
	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(400, gin.H{"ok": false, "error": "NO_FILES"})
		return
	}
	dir := filepath.Join(h.cfg.UploadDir, g.ID.Hex())
	if err := os.MkdirAll(dir, 0o755); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "UPLOAD_FAIL"})
		return
	}

	var saved []string
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
			continue
		}
		name := security.RandomCode(6) + ext
		dst := filepath.Join(dir, name)
		if err := c.SaveUploadedFile(f, dst); err == nil {
			saved = append(saved, name)
		}
	}
	if len(saved) == 0 {
		c.JSON(400, gin.H{"ok": false, "error": "NO_VALID_IMAGES"})
		return
	}

	newPhotos := append(g.Photos, saved...)
	if err := h.repo.UpdateGreetingOwned(c, id, uid, bson.M{"photos": newPhotos}); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}

	c.JSON(200, gin.H{"ok": true, "photos": newPhotos})
}

func (h *Handler) ServeMedia(c *gin.Context) {
	gid, err := primitive.ObjectIDFromHex(c.Param("gid"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	fn := filepath.Base(c.Param("fn"))
	if fn == "." || fn == "/" || strings.Contains(fn, "..") {
		c.AbortWithStatus(400)
		return
	}

	g, err := h.repo.GetGreetingByID(c, gid)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	// allowed if owner OR public paid+not expired
	uidAny, ok := c.Get(mw.CtxUserID)
	ownerOK := ok && uidAny.(primitive.ObjectID) == g.OwnerID

	publicOK := false
	if g.Paid {
		if g.ExpiresAt == nil || g.ExpiresAt.After(time.Now()) {
			publicOK = true
		}
	}

	if !ownerOK && !publicOK {
		c.AbortWithStatus(403)
		return
	}

	path := filepath.Join(h.cfg.UploadDir, g.ID.Hex(), fn)
	c.File(path)
}

func (h *Handler) DeleteGreeting(c *gin.Context) {
	uid := mw.UserID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}
	g, err := h.repo.GetGreetingOwned(c, id, uid)
	if err != nil {
		c.JSON(404, gin.H{"ok": false, "error": "NOT_FOUND"})
		return
	}
	_ = os.RemoveAll(filepath.Join(h.cfg.UploadDir, g.ID.Hex()))
	if err := h.repo.DeleteGreeting(c, g.ID); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) PublishGreeting(c *gin.Context) {
	uid := mw.UserID(c)
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}
	g, err := h.repo.GetGreetingOwned(c, id, uid)
	if err != nil {
		c.JSON(404, gin.H{"ok": false, "error": "NOT_FOUND"})
		return
	}

	if g.Type == models.GreetingSubdomain {
		sub := utils.NormalizeSubdomain(g.Subdomain)
		if sub == "" || !utils.IsValidSubdomainLabel(sub) {
			c.JSON(400, gin.H{"ok": false, "error": "SUBDOMAIN_REQUIRED"})
			return
		}
	}

	amount := h.cfg.PriceCodeRUB
	if g.Type == models.GreetingSubdomain {
		amount = h.cfg.PriceSubRUB
	}

	order, err := h.repo.CreateOrder(c, uid, g.ID, g.Type, amount, h.pay.Name())
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "ORDER_CREATE_FAIL"})
		return
	}

	checkoutURL, providerRef, err := h.pay.CreateCheckout(c, order)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "PAYMENT_CREATE_FAIL"})
		return
	}
	_ = h.repo.SetOrderProviderRef(c, order.ID, providerRef)

	c.JSON(200, gin.H{
		"ok":          true,
		"order":       order,
		"checkoutUrl": checkoutURL,
	})
}

func (h *Handler) ListOrders(c *gin.Context) {
	uid := mw.UserID(c)
	items, err := h.repo.ListOrdersByOwner(c, uid)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "items": items})
}

func (h *Handler) GetOrder(c *gin.Context) {
	uid := mw.UserID(c)
	oid, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}
	o, err := h.repo.GetOrderOwned(c, oid, uid)
	if err != nil {
		c.JSON(404, gin.H{"ok": false, "error": "NOT_FOUND"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "order": o})
}

func (h *Handler) PublicByCode(c *gin.Context) {
	code := strings.TrimSpace(c.Param("code"))
	g, err := h.repo.FindPublicByCode(c, code)
	if err != nil {
		c.JSON(404, gin.H{"ok": false, "error": "NOT_FOUND"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "item": g})
}

func (h *Handler) PublicBySubdomain(c *gin.Context) {
	sub := utils.NormalizeSubdomain(c.Param("sub"))
	g, err := h.repo.FindPublicBySubdomain(c, sub)
	if err != nil {
		c.JSON(404, gin.H{"ok": false, "error": "NOT_FOUND"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "item": g})
}

func (h *Handler) PaymentsWebhook(c *gin.Context) {
	// This endpoint is prepared for real payment providers.
	// With stub provider it does nothing.
	body, _ := c.GetRawData()
	headers := map[string]string{}
	for k, v := range c.Request.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	ref, paid, err := h.pay.HandleWebhook(c, body, headers)
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "WEBHOOK_INVALID"})
		return
	}
	if !paid || ref == "" {
		c.JSON(200, gin.H{"ok": true})
		return
	}
	// find order by providerRef and mark as paid
	o, err := h.repo.FindOrderByProviderRef(context.Background(), ref)
	if err != nil {
		c.JSON(200, gin.H{"ok": true})
		return
	}
	_, _ = h.activateOrderPaid(context.Background(), o.ID)
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) AdminStats(c *gin.Context) {
	countPaid, sumPaid, err := h.repo.AdminSalesStats(c)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "paidOrders": countPaid, "paidSumRUB": sumPaid})
}

func (h *Handler) AdminUsers(c *gin.Context) {
	users, err := h.repo.AdminListUsers(c, 200)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	out := make([]gin.H, 0, len(users))
	for _, u := range users {
		out = append(out, userPublic(u))
	}
	c.JSON(200, gin.H{"ok": true, "items": out})
}

func (h *Handler) AdminDeleteUser(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}
	if err := h.repo.AdminDeleteUser(c, id); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

func (h *Handler) AdminOrders(c *gin.Context) {
	orders, err := h.repo.AdminListOrders(c, 300)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "items": orders})
}

func (h *Handler) AdminMarkPaid(c *gin.Context) {
	orderID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}
	o, err := h.activateOrderPaid(c, orderID)
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"ok": true, "order": o})
}

func (h *Handler) activateOrderPaid(ctx context.Context, orderID primitive.ObjectID) (models.Order, error) {
	o, err := h.repo.MarkOrderPaid(ctx, orderID)
	if err != nil {
		return models.Order{}, fmt.Errorf("ORDER_MARK_FAIL")
	}
	g, err := h.repo.GetGreetingByID(ctx, o.GreetingID)
	if err != nil {
		return o, fmt.Errorf("GREETING_NOT_FOUND")
	}

	var dns *models.DNSRecord
	if g.Type == models.GreetingSubdomain {
		sub := utils.NormalizeSubdomain(g.Subdomain)
		if sub == "" || !utils.IsValidSubdomainLabel(sub) {
			return o, fmt.Errorf("SUBDOMAIN_REQUIRED")
		}

	}

	if err := h.repo.MarkGreetingPaidAndPublish(ctx, g.ID, h.cfg.PublishTTL, dns); err != nil {
		return o, fmt.Errorf("PUBLISH_FAIL")
	}
	return o, nil
}

func (h *Handler) AdminGreetings(c *gin.Context) {
	items, err := h.repo.AdminListGreetings(c, 300)
	if err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true, "items": items})
}

func (h *Handler) AdminDeleteGreeting(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_ID"})
		return
	}
	g, err := h.repo.GetGreetingByID(c, id)
	if err != nil {
		c.JSON(404, gin.H{"ok": false, "error": "NOT_FOUND"})
		return
	}

	// try remove DNS

	_ = os.RemoveAll(filepath.Join(h.cfg.UploadDir, g.ID.Hex()))
	if err := h.repo.DeleteGreeting(c, g.ID); err != nil {
		c.JSON(500, gin.H{"ok": false, "error": "DB_FAIL"})
		return
	}
	c.JSON(200, gin.H{"ok": true})
}

// Media отдаёт файлы загрузок: /api/media/:gid/:filename
func (h *Handler) Media(c *gin.Context) {
	gid := c.Param("gid")
	filename := c.Param("filename")
	if gid == "" || filename == "" {
		c.JSON(400, gin.H{"ok": false, "error": "BAD_REQUEST"})
		return
	}

	// Безопасно собираем путь: /data/uploads/<greetingId>/<filename>
	full := filepath.Join(h.cfg.UploadDir, gid, filepath.Clean("/"+filename))
	// filepath.Clean("/"+filename) гарантирует что не будет ../ с выходом из каталога

	c.File(full)
}
