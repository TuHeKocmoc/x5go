package service_test

import (
	"errors"
	"testing"

	"example.com/layered-orders/internal/model"
	"example.com/layered-orders/internal/service"
)

type fakeRepository struct {
	calls   []model.Order
	saveErr error
}

func (r *fakeRepository) Save(order model.Order) error {
	r.calls = append(r.calls, order)
	return r.saveErr
}

func (r *fakeRepository) ByID(int) (model.Order, error) {
	return model.Order{}, errors.New("неожиданный вызов ByID")
}

func TestCreateUsesRepository(t *testing.T) {
	repo := &fakeRepository{}
	orders := service.NewService(repo)
	for index, customer := range []string{"Анна", "  Борис  "} {
		want := model.Order{ID: index + 1, Customer: customer, Status: "new"}
		got, err := orders.Create(customer)
		if err != nil || got != want {
			t.Fatalf("Create(%q): получено %+v, %v; ожидалось %+v", customer, got, err, want)
		}
		if len(repo.calls) != index+1 || repo.calls[index] != want {
			t.Fatalf("Репозиторию передан неверный заказ: %+v", repo.calls)
		}
	}
}

func TestCreateRejectsEmptyCustomer(t *testing.T) {
	repo := &fakeRepository{}
	orders := service.NewService(repo)
	for _, customer := range []string{"", " \t\n"} {
		got, err := orders.Create(customer)
		if !errors.Is(err, service.ErrEmptyCustomer) || got != (model.Order{}) {
			t.Fatalf("Create(%q): ожидалась ErrEmptyCustomer, получено %+v, %v", customer, got, err)
		}
	}
	if len(repo.calls) != 0 {
		t.Fatalf("Некорректные заказы не должны сохраняться: %+v", repo.calls)
	}
	got, err := orders.Create("Анна")
	if err != nil || got.ID != 1 {
		t.Fatalf("Некорректный ввод не должен занимать ID: %+v, %v", got, err)
	}
}

func TestCreateWrapsSaveErrorAndRetriesID(t *testing.T) {
	cause := errors.New("хранилище недоступно")
	repo := &fakeRepository{saveErr: cause}
	orders := service.NewService(repo)
	got, err := orders.Create("Анна")
	if !errors.Is(err, cause) || err == cause || got != (model.Order{}) {
		t.Fatalf("Ожидалась обёрнутая ошибка сохранения, получено %+v, %v", got, err)
	}
	if len(repo.calls) != 1 || repo.calls[0].ID != 1 {
		t.Fatalf("Неверная первая попытка сохранения: %+v", repo.calls)
	}

	repo.saveErr = nil
	got, err = orders.Create("Борис")
	want := model.Order{ID: 1, Customer: "Борис", Status: "new"}
	if err != nil || got != want || len(repo.calls) != 2 || repo.calls[1] != want {
		t.Fatalf("После ошибки нужно повторно использовать ID 1: %+v, %v; вызовы %+v", got, err, repo.calls)
	}
}
