package handler

import "example.com/layered-orders/internal/model"

type OrderService interface {
	Create(customer string) (model.Order, error)
}

type Handler struct {
	orders OrderService
}

type CreateRequest struct {
	Customer string
}

func New(orders OrderService) *Handler {
	return &Handler{orders: orders}
}

func (h *Handler) Create(request CreateRequest) (model.Order, error) {
	return h.orders.Create(request.Customer)
}
