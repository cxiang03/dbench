package dbench

import (
	"context"
	"database/sql"

	_ "github.com/go-sql-driver/mysql" // mysql
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/mysqldialect"
)

type Client struct {
	*bun.DB
}

func NewDB(dsn string) *Client {
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqlDB, mysqldialect.New())
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(0)
	db.SetConnMaxIdleTime(0)
	return &Client{DB: db}
}

func (c *Client) Insert(ctx context.Context, records []*Record) error {
	_, err := c.NewInsert().Model(&records).Exec(ctx)
	return err
}
