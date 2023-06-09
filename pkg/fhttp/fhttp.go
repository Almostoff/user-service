package fhttp

import (
	"UsersService/config"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	"strconv"
	"time"
)

type Client struct {
	cfg *config.Config
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg: cfg,
	}
}

var timeout = time.Millisecond * 15000
var requestSerialID string
var scenarioName string
var requestTimeLimit = make(map[string]int)
var requestDelay = make(map[string]time.Duration)
var lastRequestTime = make(map[string]time.Time)

func init() {
	for k, v := range requestTimeLimit {
		if v == 0 {
			panic("Ты че ебан? Какие запрос 0 раз в секунду?!")
		}
		requestDelay[k] = time.Duration(1000/v) * time.Millisecond
	}
	scenarioName = "fullOK"
	requestSerialID = uuid.New().String()
}

func (h *Client) validateRequestTime(url string) {
	if _, ok := requestTimeLimit[url]; !ok {
		return
	}
	delay := time.Until(lastRequestTime[url].Add(requestDelay[url])) // через сколько нужно отправить следующий запрос
	if delay > 0 {
		time.Sleep(delay)
	}
}

func (h *Client) Request(method string, url string, body []byte, queryParams map[string]string, headers map[string]string) (responseBody []byte, statusCode int, err error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	var queryCount int
	for k, v := range queryParams {
		switch queryCount {
		case 0:
			url += fmt.Sprintf("?%s=%s", k, v)
		default:
			url += fmt.Sprintf("&%s=%s", k, v)
		}
		queryCount++
	}
	req.Header.Set("requestSerialID", requestSerialID)
	req.Header.Set("scenarioName", scenarioName)

	req.SetBody(body)
	req.SetRequestURI(url)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	client := &fasthttp.Client{}

	if err = client.DoTimeout(req, res, timeout); err != nil {
		return
	}

	statusCode = res.StatusCode()

	responseBody = make([]byte, len(res.Body()))
	copy(responseBody, res.Body())
	req.SetConnectionClose()
	return
}

func (h *Client) RequestProxy(method string, url string, proxy string, body []byte, queryParams map[string]string, headers map[string]string) (responseBody []byte, statusCode int, err error) {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	var queryCount int
	req.SetConnectionClose()
	for k, v := range queryParams {
		switch queryCount {
		case 0:
			url += fmt.Sprintf("?%s=%s", k, v)
		default:
			url += fmt.Sprintf("&%s=%s", k, v)
		}
		queryCount++
	}
	req.Header.Set("requestSerialID", requestSerialID)
	req.Header.Set("scenarioName", scenarioName)

	req.SetBody(body)
	req.SetRequestURI(url)
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	client := &fasthttp.Client{
		Dial: fasthttpproxy.FasthttpHTTPDialer(proxy),
	}

	if err = client.DoTimeout(req, res, timeout); err != nil {
		return
	}

	statusCode = res.StatusCode()

	responseBody = make([]byte, len(res.Body()))
	copy(responseBody, res.Body())
	return
}

func (h *Client) AuthHeaders(body, privateKey, publicKey string) map[string]string {
	headers := make(map[string]string)
	mac := hmac.New(sha512.New, []byte(privateKey))
	timeStamp := time.Now().Unix()
	timeStampString := strconv.Itoa(int(timeStamp))
	mac.Write([]byte(timeStampString + body))
	signature := hex.EncodeToString(mac.Sum(nil))
	headers["Content-Type"] = "application/json"
	headers["ApiPublic"] = publicKey
	headers["Signature"] = signature
	headers["TimeStamp"] = timeStampString
	return headers
}
