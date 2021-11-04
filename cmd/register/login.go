package register

import (
	"Company-Chat-Register/common/config"
	"Company-Chat-Register/common/db/mysql_model"
	my_err "Company-Chat-Register/common/err"
	"Company-Chat-Register/common/log"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type ParamsLogin struct {
	PhoneNumber string `json:"phoneNumber" binding:"required,min=11,max=11"`
	Password    string `json:"password"`
	Platform    int32  `json:"platform"`
}

func Login(c *gin.Context) {

	log.InfoByKv("Login api is statrting...", "")

	params := ParamsLogin{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": my_err.FormattingError, "errMsg": err.Error()})
		return
	}

	log.InfoByKv("api Login get params", params.PhoneNumber)

	queryParams := mysql_model.Register{
		PhoneNumber: params.PhoneNumber,
		Password:    params.Password,
	}

	canLogin := mysql_model.Login(&queryParams)
	if canLogin == 1 {
		log.ErrorByKv("Incorrect phone number password", params.PhoneNumber, "err", "Mobile phone number is not registered")
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.LogicalError, "errMsg": "Mobile phone number is not registered"})
		return
	}
	if canLogin == 2 {
		log.ErrorByKv("Incorrect phone number password", params.PhoneNumber, "err", "Incorrect password")
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.LogicalError, "errMsg": "Incorrect password"})
		return
	}

	resp, err := OpenIMToken(params.PhoneNumber, params.Platform)
	if err != nil {
		log.ErrorByKv("get token by phone number err", params.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.HttpError, "errMsg": err.Error()})
		return
	}
	response, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.ErrorByKv("Failed to read file", params.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.IoErrot, "errMsg": err.Error()})
		return
	}
	imRep := IMRegisterResp{}
	err = json.Unmarshal(response, &imRep)
	if err != nil {
		log.ErrorByKv("json parsing failed", params.PhoneNumber, "err", err.Error())
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.FormattingError, "errMsg": err.Error()})
		return
	}

	if imRep.ErrCode != 0 {
		log.ErrorByKv("openIM Login request failed", params.PhoneNumber, "err")
		c.JSON(http.StatusOK, gin.H{"errCode": my_err.HttpError, "errMsg": imRep.ErrMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"errCode": my_err.NoError, "errMsg": "", "data": imRep.Data})
	return

}

func OpenIMToken(phoneNumber string, platform int32) (*http.Response, error) {
	url := fmt.Sprintf("http://%s:10000/auth/user_token", config.Config.ImServer.IP)

	client := &http.Client{}
	params := make(map[string]interface{})

	params["secret"] = "tuoyun"
	params["platform"] = platform
	params["uid"] = phoneNumber
	con, err := json.Marshal(params)
	if err != nil {
		log.ErrorByKv("json parsing failed", phoneNumber, "err", err.Error())
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(con)))
	if err != nil {
		log.ErrorByKv("request error", "/auth/user_token", "err", err.Error())
		return nil, err
	}

	resp, err := client.Do(req)
	return resp, err
}
