package handlers

import (
	"encoding/json"
	"net/http"

	"happy-api/internal/http/mw"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *Handler) RoboInitPayment(c *gin.Context) {
	if h.pay == nil || h.pay.Name() != "robokassa" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "PAYMENT_PROVIDER_NOT_ROBOKASSA"})
		return
	}

	var in struct {
		OrderID string `json:"orderId"`
	}
	if err := c.ShouldBindJSON(&in); err != nil || in.OrderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "BAD_REQUEST"})
		return
	}

	oid, err := primitive.ObjectIDFromHex(in.OrderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "error": "BAD_ORDER_ID"})
		return
	}

	uid := mw.UserID(c)
	if uid == primitive.NilObjectID {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "error": "UNAUTHORIZED"})
		return
	}

	order, err := h.repo.GetOrderOwned(c, oid, uid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"ok": false, "error": "ORDER_NOT_FOUND"})
		return
	}

	payURL, providerRef, err := h.pay.CreateCheckout(c, order)
	if err != nil || payURL == "" {
		c.JSON(http.StatusBadGateway, gin.H{"ok": false, "error": "ROBO_INIT_FAILED"})
		return
	}

	// на всякий — сохраним providerRef
	_ = h.repo.SetOrderProviderRef(c, order.ID, providerRef)

	c.JSON(http.StatusOK, gin.H{"ok": true, "payUrl": payURL})
}

// ResultURL должен получить ответ "OK<InvId>", иначе робокасса будет ретраить.
func (h *Handler) PaymentWebhook(c *gin.Context) {
	if h.pay == nil || h.pay.Name() != "robokassa" {
		c.String(http.StatusBadRequest, "bad signature")
		return
	}

	_ = c.Request.ParseForm()

	// соберём ВСЕ параметры (query + form) в map
	params := map[string]string{}
	for k, vs := range c.Request.Form {
		if len(vs) > 0 {
			params[k] = vs[0]
		}
	}
	invId := params["InvId"]

	body, _ := json.Marshal(params)

	headers := map[string]string{}
	for k, vs := range c.Request.Header {
		if len(vs) > 0 {
			headers[k] = vs[0]
		}
	}

	providerRef, paid, err := h.pay.HandleWebhook(c, body, headers)
	if err != nil {
		c.String(http.StatusBadRequest, "bad signature")
		return
	}
	if !paid {
		c.String(http.StatusOK, "OK"+invId)
		return
	}

	order, err := h.repo.FindOrderByProviderRef(c, providerRef)
	if err == nil {
		_, _ = h.activateOrderPaid(c, order.ID)
	}

	c.String(http.StatusOK, "OK"+invId)
}

func (h *Handler) RoboResultWebhook(c *gin.Context) {
	h.PaymentWebhook(c)
}
