package services

import (
	"github.com/atmoxao/sql-tools/pkg/global"
	"github.com/jmoiron/sqlx"
)

type MysqlServices struct {
	Db *sqlx.DB
}

func NewMysqlServices() MysqlServices {
	return MysqlServices{Db: global.Conn}
}
