package handler_test

import (
	"errors"
	"testing"

	"example.com/layered-orders/internal/handler"
	"example.com/layered-orders/internal/model"
)

type stubOrderService struct {
	customers []string
	order     model.Order
	err       error
}

func (s *stubOrderService) Create(customer string) (model.Order, error) {
	s.customers = append(s.customers, customer)
	return s.order, s.err
}

func TestCreateDelegatesToService(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
	}{
		{name: "success"},
		{name: "error", err: errors.New("ошибка сервиса")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := model.Order{ID: 42, Customer: "Анна", Status: "new"}
			orders := &stubOrderService{order: want, err: tc.err}
			h := handler.New(orders)
			request := handler.CreateRequest{Customer: "  Анна  "}
			got, err := h.Create(request)
			if got != want || err != tc.err {
				t.Fatalf("Результат сервиса изменён: %+v, %v", got, err)
			}
			if len(orders.customers) != 1 || orders.customers[0] != request.Customer {
				t.Fatalf("Неверные вызовы сервиса: %q", orders.customers)
			}
		})
	}
}
