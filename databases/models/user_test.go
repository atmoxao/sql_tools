package models

import (
	"fmt"
	"github.com/atmoxao/sql-tools/pkg/global"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
	log "github.com/sirupsen/logrus"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dsn := "root:123456@tcp(vm202002:3306)/demo"
	var err error
	global.Conn, err = sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return
	}

	logger := log.New()
	logger.Level = log.InfoLevel            // miminum level
	logger.Formatter = &log.TextFormatter{} // logrus automatically add time field
	global.Conn.DB = sqldblogger.OpenDriver(
		dsn,
		global.Conn.Driver(),
		logrusadapter.New(logger),
		// optional config...
	)

	global.Conn.SetMaxOpenConns(200)
	global.Conn.SetMaxIdleConns(10)

	err = global.Conn.Ping()

	if err != nil {
		fmt.Printf("connect server failed, err:%v\n", err)
		return
	}
}

func TestUserInsert(t *testing.T) {
	var users []*Users
	u1 := &Users{
		Email: "1",
	}
	u2 := &Users{
		Email: "2",
	}

	users = append(users, u1, u2)
	newUser := new(Users)
	var userStruct = sqlbuilder.NewStruct(newUser)
	sb := userStruct.InsertInto(newUser.TableName(), sqlbuilder.Flatten(users)...)

	sql, args := sb.Build()

	global.Conn.Exec(sql, args...)

	log.Info(sql)
	log.Info(args)
}

func TestUsersSelect(t *testing.T) {
	var users []Users

	newUser := new(Users)
	var userStruct = sqlbuilder.NewStruct(newUser)
	sb := userStruct.SelectFrom(newUser.TableName())

	sql, args := sb.Build()

	err := global.Conn.Select(&users, sql, args...)
	if err != nil {
		log.Fatalln(err)
	}

	log.Info(sql)
	log.Info(args)
	log.Infof("users: %+v", users)
}

func TestUserDelete(t *testing.T) {
	var users = &Users{
		Model: Model{
			Id: 3,
		},
		Email: "at@admin.com",
	}

	var userStruct = sqlbuilder.NewStruct(users)
	sb := userStruct.DeleteFrom(users.TableName())
	sb.Where(sb.E("id", users.Id))
	sql, args := sb.Build()

	res, err := global.Conn.Exec(sql, args...)
	if err != nil {
		log.Fatalln(err)
	}

	affect, _ := res.RowsAffected()

	log.Infof("affect:%d", affect)
}

func TestUserUpdate(t *testing.T) {
	//var users = &Users{
	//	Model: Model{
	//		Id: 5,
	//	},
	//	Email: "at@admin.com",
	//}
	//
	//var userStruct = sqlbuilder.NewStruct(users)
	//sb := userStruct.Update(users.TableName(), users)
	//sb.Where(sb.GE("id", users.Id))
	//sql, args := sb.Build()
	var user = new(Users)
	ub := sqlbuilder.NewUpdateBuilder()
	ub.Update(user.TableName())
	ub.Set(ub.Assign("email", "at@admin.com"))
	ub.Where(ub.GE("id", 5))

	sql, args := ub.Build()

	res, err := global.Conn.Exec(sql, args...)
	if err != nil {
		log.Fatalln(err)
	}

	affect, _ := res.RowsAffected()

	log.Infof("affect:%d", affect)
}
