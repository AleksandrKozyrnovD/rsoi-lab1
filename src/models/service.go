package models

type Person struct {
	Id      int32  `json:"id" gorm:"primaryKey;autoIncrement" binding:"required"`
	Name    string `json:"name" gorm:"not null" binding:"required"`
	Age     int32  `json:"age"`
	Address string `json:"address"`
	Work    string `json:"work"`
}
