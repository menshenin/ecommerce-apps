// Package httphandlers HTTP-Хендлеры
package httphandlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"route256.ozon.ru/project/cart/internal/model"
	"route256.ozon.ru/project/cart/internal/service/cart"
)

// ErrInvalidPathParam Ошибка невалидного параметра
type ErrInvalidPathParam struct {
	Param string
}

func (e ErrInvalidPathParam) Error() string {
	return fmt.Sprintf("invalid path param: %s", e.Param)
}

// Handler Хендлер
type Handler func(writer http.ResponseWriter, request *http.Request) error

func (h Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := h(writer, request)
	if err == nil {
		return
	}
	status := http.StatusInternalServerError
	var t validator.ValidationErrors
	if errors.As(err, &t) {
		status = http.StatusBadRequest
	}
	http.Error(writer, err.Error(), status)
}

// AddItem Добавление товара
func AddItem(cartService *cart.Service) Handler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		userID, err := strconv.Atoi(request.PathValue("user_id"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil
		}
		skuID, err := strconv.Atoi(request.PathValue("sku_id"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil
		}
		var body struct {
			Count uint16 `json:"count"`
		}
		decoder := json.NewDecoder(request.Body)
		err = decoder.Decode(&body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil
		}

		err = cartService.AddItem(request.Context(), cart.AddItemRequest{
			UserID: int64(userID),
			SkuID:  int64(skuID),
			Count:  body.Count,
		})
		if errors.Is(err, model.ErrInsufficientStocks) {
			http.Error(writer, err.Error(), http.StatusPreconditionFailed)
			return nil
		}
		return err
	}
}

// DeleteItem Удаление товара
func DeleteItem(cartService *cart.Service) Handler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		userID, err := strconv.Atoi(request.PathValue("user_id"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil
		}
		skuID, err := strconv.Atoi(request.PathValue("sku_id"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil
		}
		err = cartService.DeleteItem(request.Context(), cart.DeleteItemRequest{
			UserID: int64(userID),
			SkuID:  int64(skuID),
		})
		if err == nil {
			writer.WriteHeader(http.StatusNoContent)
		}
		return err
	}
}

// Clear Очистка корзины
func Clear(cartService *cart.Service) Handler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		userID, err := strconv.Atoi(request.PathValue("user_id"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil
		}
		err = cartService.DeleteAllItems(request.Context(), cart.DeleteAllItemsRequest{
			UserID: int64(userID),
		})
		if errors.Is(err, model.ErrCartNotFound) {
			writer.WriteHeader(http.StatusNoContent)
		}
		return err
	}
}

// ListItem Товар в списке
type ListItem struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

// ListItemsResponse Список товаров
type ListItemsResponse struct {
	Items      []ListItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

// ListItems Вывод списка товаров
func ListItems(cartService *cart.Service) Handler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		userID, err := strconv.Atoi(request.PathValue("user_id"))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil
		}
		items, err := cartService.GetItems(request.Context(), cart.GetItemsRequest{
			UserID: int64(userID),
		})

		if err != nil {
			if errors.Is(err, model.ErrCartNotFound) {
				writer.WriteHeader(http.StatusNotFound)
				return nil
			}
		}
		if len(items.Items) == 0 {
			writer.WriteHeader(http.StatusNotFound)
			return nil
		}
		resp := ListItemsResponse{
			Items:      make([]ListItem, len(items.Items)),
			TotalPrice: items.TotalPrice,
		}

		for i, item := range items.Items {
			resp.Items[i] = ListItem{
				SkuID: int64(item.CartItem.SKU),
				Name:  item.Name,
				Count: item.Count,
				Price: item.Price,
			}
		}
		encoder := json.NewEncoder(writer)
		return encoder.Encode(resp)
	}
}

// Checkout Оформление заказа
func Checkout(cartService *cart.Service) Handler {
	return func(writer http.ResponseWriter, request *http.Request) error {
		var body struct {
			UserID int64 `json:"user_id"`
		}
		decoder := json.NewDecoder(request.Body)
		err := decoder.Decode(&body)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return nil
		}
		resp, err := cartService.Checkout(request.Context(), cart.CheckoutRequest{UserID: body.UserID})
		if err != nil {
			return err
		}
		result := struct {
			OrderID int64 `json:"order_id"`
		}{
			OrderID: int64(resp.OrderID),
		}

		encoder := json.NewEncoder(writer)
		return encoder.Encode(&result)
	}
}
