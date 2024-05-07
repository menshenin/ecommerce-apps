// Package order тесты репозитория заказов
package order

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/repository/dborder"
)

// Suite Тест-сьют
type Suite struct {
	suite.Suite
	repo *dborder.Repository
	conn *pgx.Conn
}

// New Конструктор
func New(conn *pgx.Conn) *Suite {
	return &Suite{
		repo: dborder.New(conn, conn),
		conn: conn,
	}
}

// SetupSuite SetupSuite Подготовка к тесту
func (s *Suite) SetupSuite() {
	_, err := s.conn.Exec(context.Background(), `TRUNCATE "order" CASCADE`)
	s.NoError(err)
}

// TearDownSuite Очистка данных после теста
func (s *Suite) TearDownSuite() {
	s.SetupSuite()
}

// TestOrder Тесты заказа
func (s *Suite) TestOrder() {
	testItems := []model.OrderItem{
		{
			SKU:   10,
			Count: 3,
		},
		{
			SKU:   11,
			Count: 4,
		},
	}
	var testID model.OrderID
	s.Run("create order", func() {
		order, err := s.repo.Create(context.Background(), model.UserID(123), testItems)
		s.NoError(err)
		s.NotNil(order)
		s.Equal(order.UserID, model.UserID(123))
		testID = order.ID
	})

	s.Run("check created order", func() {
		order, err := s.repo.GetByID(context.Background(), testID)
		s.NoError(err)
		s.Equal(order.UserID, model.UserID(123))

		status, err := order.GetStatus(context.Background())
		s.NoError(err)
		s.Equal(model.OrderStatusNew, status)

		items, err := order.GetItems(context.Background())
		s.NoError(err)
		s.Equal(testItems, items)
	})

	s.Run("check update status", func() {
		order, err := s.repo.GetByID(context.Background(), testID)
		s.NoError(err)
		s.Equal(order.UserID, model.UserID(123))

		err = order.UpdateStatus(context.Background(), model.OrderStatusPayed)
		s.NoError(err)

		status, err := order.GetStatus(context.Background())
		s.NoError(err)
		s.Equal(model.OrderStatusPayed, status)
	})
}
