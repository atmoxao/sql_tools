package global

import "github.com/jmoiron/sqlx"

var Conn *sqlx.DB

var Conf Config

type Config struct {
	Debug bool
	System
}

type System struct {
	EnableMysql bool
	EnableRedis bool
}
