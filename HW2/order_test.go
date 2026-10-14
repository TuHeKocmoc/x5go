package main

import "testing"

func TestNewOrder(t *testing.T) {
	for _, id := range []int{0, -1} {
		order, err := NewOrder(id)
		if err == nil || order != nil {
			t.Fatalf("NewOrder(%d): ожидались nil и ошибка", id)
		}
	}

	order, err := NewOrder(1)
	if err != nil || order == nil {
		t.Fatalf("NewOrder(1): %v", err)
	}
	if order.ID != 1 || order.Status != "new" || len(order.Items) != 0 {
		t.Fatalf("Неверный новый заказ: %+v", order)
	}
	if Total(*order) != 0 {
		t.Fatal("Сумма пустого заказа должна быть 0")
	}
}

func TestAddItemAndTotal(t *testing.T) {
	order, err := NewOrder(1)
	if err != nil {
		t.Fatal(err)
	}
	if err := AddItem(order, "Кофе", 12000); err != nil {
		t.Fatal(err)
	}
	if err := AddItem(order, "Чай", 7500); err != nil {
		t.Fatal(err)
	}
	if len(order.Items) != 2 || Total(*order) != 19500 {
		t.Fatalf("Ожидались две позиции и сумма 19500: %+v", order)
	}
	if order.Items[0].Name != "Кофе" || order.Items[1].Price != 7500 {
		t.Fatalf("Данные позиций не сохранились: %+v", order.Items)
	}
	if err := AddItem(order, "Подарок", 0); err != nil {
		t.Fatalf("Нулевая цена допустима: %v", err)
	}
	if len(order.Items) != 3 || Total(*order) != 19500 {
		t.Fatal("Позиция с нулевой ценой должна добавляться без изменения суммы")
	}
}

func TestAddItemInvalidInput(t *testing.T) {
	if err := AddItem(nil, "Кофе", 12000); err == nil {
		t.Fatal("Ожидалась ошибка для nil")
	}

	order := Order{ID: 1, Status: "new", Items: []Item{{Name: "Чай", Price: 7500}}}
	for _, item := range []Item{
		{Name: "", Price: 100},
		{Name: " \t ", Price: 100},
		{Name: "Кофе", Price: -1},
	} {
		if err := AddItem(&order, item.Name, item.Price); err == nil {
			t.Fatalf("Ожидалась ошибка для %+v", item)
		}
		if order.ID != 1 || order.Status != "new" || len(order.Items) != 1 ||
			order.Items[0] != (Item{Name: "Чай", Price: 7500}) {
			t.Fatalf("Некорректный ввод изменил заказ: %+v", order)
		}
	}
}

func TestCopyOrder(t *testing.T) {
	order := Order{ID: 1, Status: "new", Items: make([]Item, 1, 4)}
	order.Items[0] = Item{Name: "Чай", Price: 7500}
	orderCopy := CopyOrder(order)
	orderCopy.ID = 2
	orderCopy.Status = "paid"
	orderCopy.Items[0].Name = "Кофе"
	orderCopy.Items[0].Price = 12000
	if err := AddItem(&orderCopy, "Подарок", 0); err != nil {
		t.Fatal(err)
	}
	if order.ID != 1 || order.Status != "new" || len(order.Items) != 1 ||
		order.Items[0] != (Item{Name: "Чай", Price: 7500}) {
		t.Fatalf("Изменение копии затронуло оригинал: %+v", order)
	}
	if len(orderCopy.Items) != 2 || Total(orderCopy) != 12000 {
		t.Fatalf("Неверное содержимое копии: %+v", orderCopy)
	}

	emptyCopy := CopyOrder(Order{ID: 3})
	if emptyCopy.ID != 3 || emptyCopy.Items != nil || Total(emptyCopy) != 0 {
		t.Fatalf("Неверная копия пустого заказа: %+v", emptyCopy)
	}
}
