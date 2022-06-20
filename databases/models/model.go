package models

import (
	"fmt"
	"time"
)

type LocalTime struct {
	time.Time
}

type Model struct {
	Id uint `db:"id" json:"id"`
	//CreatedAt LocalTime `db:"created_at" json:"createdAt"`
	//UpdatedAt LocalTime `db:"updated_at" json:"updatedAt"`
	//DeletedAt LocalTime `db:"deleted_at" json:"deletedAt"`
}

// TableName 表名设置
func (Model) TableName(name string) string {
	// 添加表前缀
	return fmt.Sprintf("%s%s", "", name)
}
