package register

import (
	"Company-Chat-Register/common/config"
	"Company-Chat-Register/common/db"
	"Company-Chat-Register/common/db/mysql_model"
	my_err "Company-Chat-Register/common/err"
	"Company-Chat-Register/common/log"
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type paramsVerificationCode struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=11,max=11"`
}

func SendVerificationCode(c *gin.Context) {
	log.InfoByKv("sendCode api is statrting...", "")
	params := paramsVerificationCode{}

	if err := c.BindJSON(&params); err != nil {
		log.ErrorByKv("request params json parsing failed", params.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": my_err.FormattingError, "errMsg": err.Error()})
		return
	}

	queryParams := mysql_model.GetRegisterParams{
		PhoneNumber: params.PhoneNumber,
	}
	_, err, rowsAffected := mysql_model.GetRegister(&queryParams)

	if err == nil && rowsAffected != 0 {
		log.ErrorByKv("The phone number has been registered", queryParams.PhoneNumber, "err")
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.LogicalError, "errMsg": "The phone number has been registered"})
		return
	}

	client, err := CreateClient(tea.String(config.Config.Alibabacloud.AccessKeyId), tea.String(config.Config.Alibabacloud.AccessKeySecret))
	if err != nil {
		log.ErrorByKv("create sendSms client err", "", "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.IntentionalError, "errMsg": err.Error()})
		return
	}

	log.InfoByKv("begin sendSms", params.PhoneNumber)
	rand.Seed(time.Now().UnixNano())
	code := 100000 + rand.Intn(900000)
	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(params.PhoneNumber),
		SignName:      tea.String(config.Config.Alibabacloud.SignName),
		TemplateCode:  tea.String(config.Config.Alibabacloud.VerificationCodeTemplateCode),
		TemplateParam: tea.String(fmt.Sprintf("{\"code\":\"%d\"}", code)),
	}

	response, err := client.SendSms(sendSmsRequest)
	if err != nil {
		log.ErrorByKv("sendSms error", params.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.IntentionalError, "errMsg": err.Error()})
		return
	}
	if *response.Body.Code != "OK" {
		log.ErrorByKv("alibabacloud sendSms error", params.PhoneNumber, "err", response.Body.Code, response.Body.Message)
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.IntentionalError, "errMsg": *response.Body.Message})
		return
	}

	redisConn := db.RedisPool.Get()
	defer redisConn.Close()
	log.InfoByKv("begin store redis", params.PhoneNumber)
	v, err := redis.Int(redisConn.Do("TTL", params.PhoneNumber))
	if err != nil {
		log.ErrorByKv("get phoneNumber from redis error", params.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.IntentionalError, "errMsg": err.Error()})
		return
	}
	switch {
	case v == -2:
		_, err = redisConn.Do("SET", params.PhoneNumber, code, "EX", 600)
		if err != nil {
			log.ErrorByKv("set redis error", params.PhoneNumber, "err", err.Error())
			c.JSON(http.StatusOK, gin.H{"errCode": my_err.IntentionalError, "errMsg": err.Error()})
			return
		}
		data := make(map[string]interface{})
		data["phoneNumber"] = params.PhoneNumber
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.NoError, "errMsg": "Verification code sent successfully!", "data": data})
		log.InfoByKv("send new verification code", params.PhoneNumber)
		return
	case v > 540:
		data := make(map[string]interface{})
		data["phoneNumber"] = params.PhoneNumber
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.LogicalError, "errMsg": "Frequent operation!", "data": data})
		log.InfoByKv("frequent operation", params.PhoneNumber)
		return
	case v < 540:
		_, err = redisConn.Do("SET", params.PhoneNumber, code, "EX", 600)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"errCode": my_err.IntentionalError, "errMsg": err.Error()})
			return
		}
		data := make(map[string]interface{})
		data["phoneNumber"] = params.PhoneNumber
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.NoError, "errMsg": "Verification code has been reset!", "data": data})
		log.InfoByKv("Reset verification code", params.PhoneNumber)
		return
	}

}

func CreateClient(accessKeyId *string, accessKeySecret *string) (result *dysmsapi20170525.Client, err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}

	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	result = &dysmsapi20170525.Client{}
	result, err = dysmsapi20170525.NewClient(config)
	return result, err
}
