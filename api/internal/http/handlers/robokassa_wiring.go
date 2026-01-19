package handlers

import "happy-api/internal/payments/robokassa"

func (h *Handler) SetRobo(c *robokassa.Client) {
	h.robo = c
}
