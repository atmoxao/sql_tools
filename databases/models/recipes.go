package models

type Recipes struct {
	Model
	UserId       int64  `db:"user_id" json:"user_id"`
	Level        int64  `db:"level" json:"level"`
	IsVegetarian int64  `db:"is_vegetarian" json:"is_vegetarian"`
	PrepTime     string `db:"prep_time" json:"prep_time"`
	Title        string `db:"title" json:"title"`
}

func (m Recipes) TableName() string {
	return m.Model.TableName("recipes")
}
