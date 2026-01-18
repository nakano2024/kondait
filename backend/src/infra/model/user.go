package model

type User struct {
	Code string `gorm:"column:code;primaryKey"`
	Sub  string `gorm:"column:sub;"`
}
