package robokassa

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"happy-api/internal/config"
	"happy-api/internal/models"
)

type Provider struct {
	MerchantLogin string
	Pass1         string
	Pass2         string
	IsTest        bool
	PayURL        string
}

func New(cfg config.Config) *Provider {
	return &Provider{
		MerchantLogin: cfg.RoboMerchantLogin,
		Pass1:         cfg.RoboPass1,
		Pass2:         cfg.RoboPass2,
		IsTest:        cfg.RoboIsTest,
		PayURL:        cfg.RoboPayURL,
	}
}

func (p *Provider) Name() string { return "robokassa" }

// Важно: InvId можно НЕ задавать, а связку делаем через Shp_orderId.
// Тогда providerRef = order.ID.Hex() (это мы же и положим в Shp_orderId)
func (p *Provider) CreateCheckout(_ context.Context, order models.Order) (string, string, error) {
	if p.MerchantLogin == "" || p.Pass1 == "" {
		return "", "", errors.New("robokassa: missing creds")
	}

	outSum := fmt.Sprintf("%.2f", float64(order.AmountRUB))

	q := url.Values{}
	q.Set("MerchantLogin", p.MerchantLogin)
	q.Set("OutSum", outSum)
	q.Set("Description", fmt.Sprintf("Order %s", order.ID.Hex()))
	if p.IsTest {
		q.Set("IsTest", "1")
	}

	// наш orderId уедет как дополнительный параметр
	q.Set("Shp_orderId", order.ID.Hex())

	sign := md5hex(signaturePay(p.MerchantLogin, outSum, "", p.Pass1, collectShp(q)))
	q.Set("SignatureValue", sign)

	payURL := p.PayURL
	if payURL == "" {
		payURL = "https://auth.robokassa.ru/Merchant/Index.aspx"
	}

	return payURL + "?" + q.Encode(), order.ID.Hex(), nil
}

// body ожидаем JSON map[string]string (его соберём в webhook-хендлере из Form+Query)
func (p *Provider) HandleWebhook(_ context.Context, body []byte, _ map[string]string) (string, bool, error) {
	if p.Pass2 == "" {
		return "", false, errors.New("robokassa: missing pass2")
	}

	var m map[string]string
	if err := json.Unmarshal(body, &m); err != nil {
		return "", false, err
	}

	// поддержим разные регистры ключей
	get := func(k string) string {
		if v, ok := m[k]; ok {
			return v
		}
		for kk, vv := range m {
			if strings.EqualFold(kk, k) {
				return vv
			}
		}
		return ""
	}

	outSum := get("OutSum")
	invId := get("InvId")
	sig := get("SignatureValue")
	orderID := get("Shp_orderId")

	if outSum == "" || sig == "" || orderID == "" {
		return "", false, errors.New("robokassa: missing required params")
	}

	// соберём Shp_* из входящих данных
	shp := make([]string, 0, 4)
	for k, v := range m {
		if strings.HasPrefix(strings.ToLower(k), "shp_") {
			// robokassa требует точный формат "Shp_xxx=yyy" в подписи
			// регистр ключа обычно "Shp_..." — оставим как пришло:
			shp = append(shp, fmt.Sprintf("%s=%s", k, v))
		}
	}
	sort.Strings(shp)

	expected := md5hex(signatureResult(outSum, invId, p.Pass2, shp))
	if !strings.EqualFold(expected, sig) {
		return "", false, errors.New("robokassa: bad signature")
	}

	// ResultURL приходит только при успешной оплате — считаем paid=true
	return orderID, true, nil
}

func signaturePay(merchantLogin, outSum, invId, pass1 string, shp []string) string {
	// Формат из доков: MerchantLogin:OutSum:InvId:Password#1(:Shp_...=...) :contentReference[oaicite:1]{index=1}
	base := fmt.Sprintf("%s:%s:%s:%s", merchantLogin, outSum, invId, pass1)
	for _, kv := range shp {
		base += ":" + kv
	}
	return base
}

func signatureResult(outSum, invId, pass2 string, shp []string) string {
	// Стандартно для ResultURL: OutSum:InvId:Password#2(:Shp_...=...)
	base := fmt.Sprintf("%s:%s:%s", outSum, invId, pass2)
	for _, kv := range shp {
		base += ":" + kv
	}
	return base
}

func collectShp(q url.Values) []string {
	var shp []string
	for k, vs := range q {
		if strings.HasPrefix(strings.ToLower(k), "shp_") && len(vs) > 0 {
			shp = append(shp, fmt.Sprintf("%s=%s", k, vs[0]))
		}
	}
	sort.Strings(shp)
	return shp
}
