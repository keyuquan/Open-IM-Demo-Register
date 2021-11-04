package mysql_model

import (
	"Company-Chat-Register/common/db"
	_ "github.com/jinzhu/gorm"
)

type GetRegisterParams struct {
	PhoneNumber string `json:"phoneNumber"`
}

func GetRegister(params *GetRegisterParams) (Register, error, int64) {
	var r Register
	result := db.DB.
		Where(&Register{PhoneNumber: params.PhoneNumber}).
		Find(&r)
	return r, result.Error, result.RowsAffected
}

type SetPasswordParams struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Uid         string `json:"uid"`
}

func SetPassword(params *SetPasswordParams) (Register, error) {
	r := Register{
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
	}

	result := db.DB.Create(&r)

	return r, result.Error
}

func Login(params *Register) int64 {
	var r Register
	result := db.DB.
		Where(&Register{PhoneNumber: params.PhoneNumber}).
		Find(&r)
	if result.Error != nil && result.RowsAffected == 0 {
		return 1
	}
	if r.Password != params.Password {
		return 2
	}
	return 0
}
