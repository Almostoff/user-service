package cdnService

import (
	"UsersService/config"
	"UsersService/internal/cConstants"
	"UsersService/internal/iConnection"
	"UsersService/pkg/secure"
	"UsersService/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/http2"
)

type Client struct {
	cfg        *config.Config
	httpClient *resty.Client
	iConnUC    iConnection.UseCase
}

func GetClient(params *GetClientParams) UseCase {

	http2.ConfigureTransport(http.DefaultTransport.(*http.Transport))
	customTransport := http.DefaultTransport.(*http.Transport)
	return &Client{
		cfg: params.Config,
		httpClient: resty.New().OnBeforeRequest(SignatureMiddleware(params)).
			EnableTrace().SetDebug(true).SetBaseURL(params.BaseUrl).SetTransport(customTransport),
		iConnUC: params.IConn,
	}
}

func SignatureMiddleware(connData *GetClientParams) resty.RequestMiddleware {
	return func(c *resty.Client, request *resty.Request) error {

		var (
			body      string
			timestamp int64
		)

		timestamp = utils.GetEuropeTime().Unix()
		if request.Method == "GET" {
			params := strings.Split(request.URL, "?")[1]
			body = calculateGetBody(strconv.FormatInt(timestamp, 10), params)
		} else if request.Method == "POST" {
			jsonBody, _ := json.Marshal(request.Body)
			body = createRequestBody(strconv.FormatInt(timestamp, 10), string(jsonBody))
		} else {
			jsonBody := ""
			body = createRequestBody(strconv.FormatInt(timestamp, 10), string(jsonBody))
		}

		request.SetHeader("Content-Type", "application/json")
		request.SetHeader("ApiPublic", connData.Public)
		request.SetHeader("Signature", secure.CalcSignature(connData.Private, body))
		request.SetHeader("TimeStamp", strconv.FormatInt(timestamp, 10))

		return nil
	}
}

func calculateGetBody(timestamp, query string) string {
	return timestamp + query
}

func createRequestBody(timestamp, jsonBody string) string {
	return timestamp + jsonBody
}

func (c Client) SaveImage(params *SaveImageParams) (*SaveImageResponse, error) {
	var responseModel SaveImageResponse

	response, err := c.httpClient.R().SetBody(&SaveImageCDNParams{Data: params.ImageBase64}).SetResult(&responseModel).Post(cConstants.Save)
	if err = c.handleResponse(response, err); err != nil {
		return &SaveImageResponse{}, err
	}
	if responseModel.Status != "SUCCESS" {
		return &SaveImageResponse{}, err
	}
	s := strings.Split(responseModel.File, "/")
	responseModel.File = cConstants.CdnUrl + s[len(s)-1]
	return &responseModel, nil
}

func (c *Client) handleResponse(response *resty.Response, err error) error {

	if err != nil {
		return err
	}
	if response == nil {
		return errors.New("response is nil somehow")
	}

	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return fmt.Errorf("status code {%d}", statusCode)
	}

	return nil
}
