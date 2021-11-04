package register

import (
	"Company-Chat-Register/common/config"
	"Company-Chat-Register/common/db"
	"Company-Chat-Register/common/db/mysql_model"
	my_err "Company-Chat-Register/common/err"
	"Company-Chat-Register/common/log"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type ParamsSetPassword struct {
	PhoneNumber      string `json:"phoneNumber" binding:"required,min=11,max=11"`
	Password         string `json:"password"`
	VerificationCode string `json:"verificationCode"`
}

type Data struct {
	ExpiredTime int64  `json:"expiredTime"`
	Token       string `json:"token"`
	Uid         string `json:"uid"`
}

type IMRegisterResp struct {
	Data    Data   `json:"data"`
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

func SetPassword(c *gin.Context) {
	log.InfoByKv("setPassword api is statrting...", "")
	params := ParamsSetPassword{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": my_err.FormattingError, "errMsg": err.Error()})
		return
	}

	log.InfoByKv("begin store redis", params.PhoneNumber)
	redisConn := db.RedisPool.Get()
	defer redisConn.Close()
	v, err := redis.String(redisConn.Do("GET", params.PhoneNumber))

	if params.VerificationCode == config.Config.ImServer.SuperCode {
		goto openIMRegisterTab
	}

	fmt.Println("Get Redis:", v, err)
	if err != nil {
		log.ErrorByKv("password Verification code expired", params.PhoneNumber, "err", err.Error())
		data := make(map[string]interface{})
		data["phoneNumber"] = params.PhoneNumber
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.LogicalError, "errMsg": "Verification expired!", "data": data})
		return
	}
	if v != params.VerificationCode {
		log.InfoByKv("password Verification code error", params.PhoneNumber, params.VerificationCode)
		data := make(map[string]interface{})
		data["PhoneNumber"] = params.PhoneNumber
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.LogicalError, "errMsg": "Verification code error!", "data": data})
		return
	}

openIMRegisterTab:
	log.InfoByKv("openIM register begin", params.PhoneNumber)
	resp, err := OpenIMRegister(params.PhoneNumber)

	log.InfoByKv("openIM register resp", params.PhoneNumber, resp, err)
	if err != nil {
		log.ErrorByKv("request openIM register error", params.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.HttpError, "errMsg": err.Error()})
		return
	}
	response, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.IoErrot, "errMsg": err.Error()})
		return
	}
	imrep := IMRegisterResp{}
	err = json.Unmarshal(response, &imrep)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.FormattingError, "errMsg": err.Error()})
		return
	}
	if imrep.ErrCode != 0 {
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.HttpError, "errMsg": imrep.ErrMsg})
		return
	}

	queryParams := mysql_model.SetPasswordParams{
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
	}

	log.InfoByKv("begin store mysql", params.PhoneNumber, params.Password)
	_, err = mysql_model.SetPassword(&queryParams)
	if err != nil {
		log.ErrorByKv("set phone number password error", params.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.DatabaseError, "errMsg": err.Error()})
		return
	}

	log.InfoByKv("end setPassword", params.PhoneNumber)
	c.JSON(http.StatusOK, gin.H{"errCode": my_err.NoError, "errMsg": "", "data": imrep.Data})
	return
}

func OpenIMRegister(phoneNumber string) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:10000/auth/user_register", config.Config.ImServer.IP)
	fmt.Println("1:", config.Config.ImServer.Secret)

	client := &http.Client{}

	params := make(map[string]interface{})

	params["secret"] = config.Config.ImServer.Secret
	params["platform"] = 7
	params["uid"] = phoneNumber
	params["name"] = phoneNumber
	params["icon"] = ""
	params["gender"] = 0
	params["mobile"] = phoneNumber
	params["birth"] = ""
	params["email"] = ""
	params["ex"] = ""
	con, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	log.InfoByKv("openIM register params", phoneNumber, params)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(con)))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)

	return resp, err
}
