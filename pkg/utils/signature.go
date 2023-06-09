package utils

import (
	"UsersService/config"
	"UsersService/internal/iConnection"
	"UsersService/pkg/secure"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type GetClientParams struct {
	Config  *config.Config
	Public  *string
	Private *string
	BaseUrl *string
}

type GetClientCdnParams struct {
	Config  *config.Config
	IConn   iConnection.UseCase
	Public  string
	Private string
	BaseUrl string
}

func SignatureMiddleware(connData *GetClientParams) resty.RequestMiddleware {
	return func(c *resty.Client, request *resty.Request) error {
		timestamp := GetEuropeTime().Unix()
		jsonBody, _ := json.Marshal(request.Body)
		body := createRequestBody(strconv.FormatInt(timestamp, 10), string(jsonBody))

		request.SetHeader("Content-Type", "application/json")
		request.SetHeader("ApiPublic", *connData.Public)
		request.SetHeader("Signature", secure.CalcSignature(*connData.Private, body))
		request.SetHeader("TimeStamp", strconv.FormatInt(timestamp, 10))

		return nil
	}
}

func createRequestBody(timestamp, jsonBody string) string {
	return timestamp + jsonBody
}
