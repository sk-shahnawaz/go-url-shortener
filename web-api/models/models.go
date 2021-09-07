package models

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type ExtendedContext struct {
	echo.Context
	Db *pgxpool.Pool
}
