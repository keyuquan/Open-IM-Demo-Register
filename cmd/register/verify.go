package register

import (
	"Company-Chat-Register/common/config"
	"Company-Chat-Register/common/db"
	my_err "Company-Chat-Register/common/err"
	"Company-Chat-Register/common/log"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

type paramsCertification struct {
	PhoneNumber      string `json:"phoneNumber" binding:"required,min=11,max=11"`
	VerificationCode string `json:"verificationCode"`
}

func Verify(c *gin.Context) {
	log.InfoByKv("Verify api is statrting...", "")
	params := paramsCertification{}

	if err := c.BindJSON(&params); err != nil {
		log.ErrorByKv("request params json parsing failed", "", "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": my_err.FormattingError, "errMsg": err.Error()})
		return
	}

	if params.VerificationCode == config.Config.ImServer.SuperCode {
		log.InfoByKv("Super Code Verified successfully", params.PhoneNumber)
		data := make(map[string]interface{})
		data["phoneNumber"] = params.PhoneNumber
		data["verificationCode"] = params.VerificationCode
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.NoError, "errMsg": "Verified successfully!", "data": data})
		return
	}

	log.InfoByKv("begin get form redis", params.PhoneNumber)
	redisConn := db.RedisPool.Get()
	defer redisConn.Close()
	v, err := redis.String(redisConn.Do("GET", params.PhoneNumber))
	log.InfoByKv("redis phone number and verificating Code", params.PhoneNumber, v)
	if err != nil {
		log.ErrorByKv("Verification code expired", params.PhoneNumber, "err", err.Error())
		data := make(map[string]interface{})
		data["phoneNumber"] = params.PhoneNumber
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.LogicalError, "errMsg": "Verification code expired!", "data": data})
		return
	}
	if params.VerificationCode == v {
		log.InfoByKv("Verified successfully", params.PhoneNumber)
		data := make(map[string]interface{})
		data["phoneNumber"] = params.PhoneNumber
		data["verificationCode"] = params.VerificationCode
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.NoError, "errMsg": "Verified successfully!", "data": data})
		return
	} else {
		log.InfoByKv("Verification code error", params.PhoneNumber, params.VerificationCode)
		data := make(map[string]interface{})
		data["phoneNumber"] = params.PhoneNumber
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.LogicalError, "errMsg": "Verification code error!", "data": data})
	}

}
