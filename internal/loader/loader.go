package loader

import (
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"sparky-back/internal/models"
)

const DsnTemplate = "postgres://%s:%s@%s:%d/%s?sslmode=disable"

func New(host string, port int, user, password, dbName string) *bun.DB {
	dsn := fmt.Sprintf(DsnTemplate, user, password, host, port, dbName)
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(pgdb, pgdialect.New())
	db.RegisterModel((*models.Message)(nil), (*models.Reaction)(nil), (*models.User)(nil))
	return db
}
