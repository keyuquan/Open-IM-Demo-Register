package mysql_model

type Register struct {
	PhoneNumber string `gorm:"column:phone_number"`
	Password    string `gorm:"column:password"`
}
