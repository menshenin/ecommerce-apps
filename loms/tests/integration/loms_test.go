package integration

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/loms/migrations"
	orderSuite "route256.ozon.ru/project/loms/tests/integration/repository/order"
	stockSuite "route256.ozon.ru/project/loms/tests/integration/repository/stock"
)

const testPostgresDsn = "postgresql://test:test@localhost:55629/postgres"

type MainSuite struct {
	suite.Suite
	stdDB *sql.DB
}

func (s *MainSuite) SetupSuite() {
	conf, err := pgx.ParseConfig(testPostgresDsn)
	s.NoError(err)
	goose.SetBaseFS(migrations.Migrations)
	s.stdDB = stdlib.OpenDB(*conf)
	s.NoError(goose.SetDialect("postgres"))
	s.NoError(goose.Up(s.stdDB, "."))
}

func (s *MainSuite) TearDownSuite() {
	s.NoError(goose.Down(s.stdDB, "."))
}

func (s *MainSuite) TestRepositories() {
	conn, err := pgx.Connect(context.Background(), testPostgresDsn)
	s.NoError(err)
	suite.Run(s.T(), orderSuite.New(conn))
	suite.Run(s.T(), stockSuite.New(conn))
}

func TestIntegration(t *testing.T) {
	suite.Run(t, new(MainSuite))
}
