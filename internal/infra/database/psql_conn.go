package database

import (
	"fmt"
	"time"

	// pgx driver for postgres
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/internal/infra/env"
	"github.com/projectsprintdev-mikroserpis01/gogomanager-api/pkg/log"
)

func NewPgsqlConn() *sqlx.DB {
	dataSourceName := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable ",
		env.AppEnv.DBHost,
		env.AppEnv.DBPort,
		env.AppEnv.DBUser,
		env.AppEnv.DBPass,
		env.AppEnv.DBName,
	)

	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		log.Panic(log.LogInfo{
			"error": err.Error(),
		}, "[DB][NewPgsqlConn] failed to connect to database")
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
