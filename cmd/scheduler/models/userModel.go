package models

type User struct {
	ID                 uint `gorm:"primaryKey"`
	Email              string
	Password           string
	TaskDefintionModel []TaskDefintionModel `gorm:"foreignKey:UserID"`
}
