package model

type CookingItem struct {
	Code     string `gorm:"column:code;primaryKey"`
	Name     string `gorm:"column:name"`
	UserCode string `gorm:"column:user_code"`
}

func (CookingItem) TableName() string {
	return "cooking_items"
}
