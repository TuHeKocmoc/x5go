package repository_test

import (
	"errors"
	"testing"

	"example.com/layered-orders/internal/model"
	"example.com/layered-orders/internal/repository"
)

func TestMemorySaveAndByID(t *testing.T) {
	repo := repository.NewMemory()
	want := model.Order{ID: 42, Customer: "Анна", Status: "new"}
	order := want
	if err := repo.Save(order); err != nil {
		t.Fatal(err)
	}
	order.Customer = "Борис"
	got, err := repo.ByID(42)
	if err != nil || got != want {
		t.Fatalf("Сохранённый заказ: %+v, %v; ожидалось %+v", got, err, want)
	}
	got.Status = "paid"
	again, err := repo.ByID(42)
	if err != nil || again != want {
		t.Fatalf("Изменение полученной копии затронуло хранилище: %+v, %v", again, err)
	}
}

func TestMemoryNotFound(t *testing.T) {
	got, err := repository.NewMemory().ByID(42)
	if !errors.Is(err, repository.ErrNotFound) || got != (model.Order{}) {
		t.Fatalf("Ожидалась ErrNotFound, получено %+v, %v", got, err)
	}
}
