package models

type Users struct {
	Model
	Email string `db:"email" json:"email"`
}

func (m Users) TableName() string {
	return m.Model.TableName("users")
}
