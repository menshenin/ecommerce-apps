// Package cart Сервис корзины
package cart

import (
	"context"

	"github.com/go-playground/validator/v10"
	"route256.ozon.ru/project/cart/internal/model"
)

// CartsRepository Репозиторий корзины
//
//go:generate minimock -i CartsRepository -s "_mock.go" -o ./mocks/
type CartsRepository interface {
	GetByUserID(ctx context.Context, userID model.UserID) (*model.Cart, error)
	Create(ctx context.Context, userID model.UserID) (*model.Cart, error)
}

// ItemRepository Репоизторий товаров
//
//go:generate minimock -i ItemRepository -s "_mock.go" -o ./mocks/
type ItemRepository interface {
	GetItemsBySKU(ctx context.Context, sku ...model.SKU) (map[model.SKU]*model.Item, error)
}

// LomsClient Клиент для работы с сервисом lOMS
//
//go:generate minimock -i LomsClient -s "_mock.go" -o ./mocks/
type LomsClient interface {
	Checkout(ctx context.Context, cart *model.Cart) (model.OrderID, error)
	AvailableCount(ctx context.Context, sku model.SKU) (int32, error)
}

// Service Сервис для работы с корзиной
type Service struct {
	cartRepo   CartsRepository
	itemRepo   ItemRepository
	lomsClient LomsClient
	validate   *validator.Validate
}

// New Конструктор
func New(cartRepo CartsRepository, itemRepo ItemRepository, lomsClient LomsClient) (*Service, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidationCtx("has_product", productValidator(itemRepo))
	if err != nil {
		return nil, err
	}
	return &Service{
		cartRepo:   cartRepo,
		itemRepo:   itemRepo,
		lomsClient: lomsClient,
		validate:   validate,
	}, nil
}

// productValidator Валидатор, проверяющий наличие товара по его SKU
func productValidator(repository ItemRepository) func(ctx context.Context, fl validator.FieldLevel) bool {
	return func(ctx context.Context, fl validator.FieldLevel) bool {
		if fl.Field().CanInt() {
			sku := model.SKU(fl.Field().Int())
			items, err := repository.GetItemsBySKU(ctx, sku)
			return err == nil && items != nil && items[sku] != nil
		}
		return false
	}
}
