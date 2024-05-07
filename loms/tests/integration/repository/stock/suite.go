// Package stock тесты репозитория складов
package stock

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/loms/internal/model"
	"route256.ozon.ru/project/loms/internal/repository/dbstock"
)

// Suite Тест-сьют
type Suite struct {
	suite.Suite
	repo *dbstock.Repository
	conn *pgx.Conn
}

// New Конструктор
func New(conn *pgx.Conn) *Suite {
	return &Suite{
		repo: dbstock.New(conn, conn),
		conn: conn,
	}
}

// SetupSuite Подготовка к тесту
func (s *Suite) SetupSuite() {
	_, err := s.conn.Exec(context.Background(), "TRUNCATE stock_item")
	s.NoError(err)
}

// TearDownSuite Очистка данных
func (s *Suite) TearDownSuite() {
	s.SetupSuite()
}

// TestStocks Тест складов
func (s *Suite) TestStocks() {
	testItems := []model.StockItem{
		{
			SKU:        1,
			TotalCount: 10,
			Reserved:   5,
		},
		{
			SKU:        2,
			TotalCount: 20,
			Reserved:   11,
		},
	}
	stock, err := s.repo.GetCurrentStock()
	s.NoError(err)

	s.Run("test load", func() {
		err := stock.Load(context.Background(), testItems)
		s.NoError(err)
	})

	s.Run("get available count", func() {
		count, err := stock.AvailableCount(context.Background(), model.SKU(2))
		s.NoError(err)
		s.EqualValues(9, count)
	})

	s.Run("reserve", func() {
		err := stock.Reserve(context.Background(), []model.OrderItem{
			{
				SKU:   1,
				Count: 2,
			},
		})
		s.NoError(err)
		count, err := stock.AvailableCount(context.Background(), model.SKU(1))
		s.NoError(err)
		s.EqualValues(3, count)
	})

	s.Run("cancel reserve", func() {
		err := stock.CancelReserve(context.Background(), []model.OrderItem{
			{
				SKU:   1,
				Count: 2,
			},
		})
		s.NoError(err)
		count, err := stock.AvailableCount(context.Background(), model.SKU(1))
		s.NoError(err)
		s.EqualValues(5, count)
	})

	s.Run("write off", func() {
		err := stock.WriteOff(context.Background(), []model.OrderItem{
			{
				SKU:   1,
				Count: 2,
			},
		})
		s.NoError(err)
		row := s.conn.QueryRow(
			context.Background(),
			`SELECT total_count, reserved FROM stock_item WHERE sku = $1`, 1)
		var count, reserved int
		s.NoError(row.Scan(&count, &reserved))
		s.Equal(count, 8)
		s.Equal(reserved, 3)
	})
}
