package models




type Email struct {
	Email    string  `gorm:"type:varchar(100);unique;not null"`
}