package robokassa

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Client struct {
	MerchantLogin string
	Pass1         string
	Pass2         string
	PayURL        string
	IsTest        bool
}

func (c *Client) IsConfigured() bool {
	return c.MerchantLogin != "" && c.Pass1 != "" && c.Pass2 != "" && c.PayURL != ""
}

// InvId делаем стабильным int из orderId (последние 8 hex -> uint32)
func InvIDFromOrderHex(orderHex string) int {
	if len(orderHex) < 8 {
		return 1
	}
	last8 := orderHex[len(orderHex)-8:]
	n, err := strconv.ParseUint(last8, 16, 32)
	if err != nil || n == 0 {
		return 1
	}
	return int(n)
}

func (c *Client) PaymentURL(outSum string, invID int, desc string, shp map[string]string) string {
	q := url.Values{}
	q.Set("MerchantLogin", c.MerchantLogin)
	q.Set("OutSum", outSum)
	q.Set("InvId", strconv.Itoa(invID))
	q.Set("Description", desc)

	// Shp_*
	keys := make([]string, 0, len(shp))
	for k := range shp {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		q.Set(k, shp[k])
	}

	sign := c.signPay(outSum, invID, shp)
	q.Set("SignatureValue", sign)

	if c.IsTest {
		q.Set("IsTest", "1")
	}

	return strings.TrimRight(c.PayURL, "?") + "?" + q.Encode()
}

func (c *Client) VerifyResult(outSum string, invID int, signature string, shp map[string]string) bool {
	my := c.signResult(outSum, invID, shp)
	return strings.EqualFold(my, signature)
}

func (c *Client) signPay(outSum string, invID int, shp map[string]string) string {
	// MerchantLogin:OutSum:InvId:Pass1[:Shp_x=...]
	base := fmt.Sprintf("%s:%s:%d:%s", c.MerchantLogin, outSum, invID, c.Pass1)
	return md5hex(base + shpSuffix(shp))
}

func (c *Client) signResult(outSum string, invID int, shp map[string]string) string {
	// OutSum:InvId:Pass2[:Shp_x=...]
	base := fmt.Sprintf("%s:%d:%s", outSum, invID, c.Pass2)
	return md5hex(base + shpSuffix(shp))
}

func shpSuffix(shp map[string]string) string {
	if len(shp) == 0 {
		return ""
	}
	keys := make([]string, 0, len(shp))
	for k := range shp {
		keys = append(keys, k)
	}
	sort.Strings(keys) // важно: алфавитно :contentReference[oaicite:3]{index=3}

	var b strings.Builder
	for _, k := range keys {
		b.WriteString(":")
		b.WriteString(k)
		b.WriteString("=")
		b.WriteString(shp[k])
	}
	return b.String()
}

func md5hex(s string) string {
	sum := md5.Sum([]byte(s))
	return strings.ToUpper(hex.EncodeToString(sum[:]))
}
