package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/everFinance/everpay/account"
	"github.com/gin-gonic/gin"
)

type err struct {
	Code int
	Msg  string
}

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Error  string `json:"error"`
}

var (
	Err_request_params = err{4001, "err_request_params"}
	Err_mist_response  = err{4010, "err_mist_response"}
)

func ErrorResponse(c *gin.Context, errorCode int, err, msg string) {
	var respCode int
	// 客服端 error
	if errorCode == Err_request_params.Code {
		respCode = http.StatusBadRequest
	} else {
		// 服务器端 error
		respCode = http.StatusInternalServerError
	}
	c.JSON(respCode, Response{
		Status: errorCode,
		Error:  err,
		Msg:    msg,
	})

}

func SuccessResponse(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, result)
}

func main() {
	// 创建gin实例
	r := gin.Default()

	// 配置路由
	r.GET("/riskScore", riskScore)

	// 运行服务器
	r.Run()
}

func riskScore(c *gin.Context) {
	api_key := ""
	address := c.DefaultQuery("address", "")
	accType, accid, err := account.IDCheck(address)
	if err != nil || accType != account.AccountTypeEVM {
		ErrorResponse(c, Err_request_params.Code, Err_request_params.Msg, "address incorrect")
		return
	}
	mistApiURL := fmt.Sprintf("https://openapi.misttrack.io/v1/risk_score?coin=ETH&address=%v&api_key=%v", accid, api_key)
	resp, err := http.Get(mistApiURL)
	if err != nil {
		ErrorResponse(c, Err_mist_response.Code, Err_mist_response.Msg, err.Error())
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorResponse(c, Err_mist_response.Code, Err_mist_response.Msg, err.Error())
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}
