package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strconv"
	"strings"
	"time"
)

func NumUidGenerator() string {
	timeNum := strconv.FormatInt(time.Now().UnixNano()/100, 10)
	//当前id是秒的时间戳+纳秒时间戳最后两位
	id := timeNum[:10] + timeNum[len(timeNum)-2:]
	return id
}

// EncryptPassword 密码加密
func EncryptPassword(password string) string {
	generatePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(generatePassword)
}

// CheckPasswordIsRight 检查密码是否正确
func CheckPasswordIsRight(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateUUID() string {
	newUUID, _ := uuid.NewUUID()
	return strings.ReplaceAll(newUUID.String(), "-", "")
}
