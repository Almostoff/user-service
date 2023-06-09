package cdnService

import (
	"UsersService/config"
	"UsersService/internal/iConnection"
	"time"
)

type TxInfo struct {
	TrackerId     string      `json:"TrackerId"`
	ServiceClient string      `json:"ServiceClient"`
	Token         string      `json:"Token"`
	Type          string      `json:"Type"`
	Status        string      `json:"Status"`
	Amount        float64     `json:"Amount"`
	Commission    float64     `json:"Commission"`
	CreateTime    time.Time   `json:"CreateTime"`
	UpdateTime    time.Time   `json:"UpdateTime"`
	CompleteTime  time.Time   `json:"CompleteTime"`
	Extra         interface{} `json:"Extra"`
}

type GetClientParams struct {
	Config  *config.Config
	IConn   iConnection.UseCase
	Public  string
	Private string
	BaseUrl string
}

type TransactionDescriptionStatus struct {
	Description string `json:"description,omitempty"`
	Status      string `json:"status"`
}

// ------------------------------------------------------------------------------------------------

type CreateModel struct {
	ServiceClient string  `json:"service_client"`
	Commission    float64 `json:"Commission"`
	Amount        float64 `json:"amount"`
	Token         string  `json:"token"`
}

type ApproveTransactionParams struct {
	InternalId string `json:"internal_id"`
}

type CancelTransactionParams struct {
	InternalId string `json:"internal_id"`
}

type CreateWithdrawParams CreateModel

type CreateInvoiceParams CreateModel

type GetBalanceParams struct {
	ServiceClient string `json:"service_client"`
	Token         string `json:"token"`
}

// ------------------------------------------------------------------------------------------------

type BalanceModel struct {
	Token     string  `json:"token"`
	Freeze    float64 `json:"freeze"`
	Available float64 `json:"available"`
}

type ApproveTransactionResponse struct {
	TransactionDescriptionStatus
}

type CancelTransactionResponse struct {
	TransactionDescriptionStatus
}

type TxInfoG struct {
	TrackerId     string      `json:"tracker_id"`
	ServiceClient string      `json:"service_client"`
	Token         string      `json:"token"`
	Type          string      `json:"type"`
	Status        string      `json:"status"`
	Amount        float64     `json:"amount"`
	Commission    float64     `json:"commission"`
	CreateTime    time.Time   `json:"create_time"`
	UpdateTime    time.Time   `json:"update_time"`
	CompleteTime  time.Time   `json:"complete_time"`
	Extra         interface{} `json:"extra"`
}

type CreateWithdrawResponse struct {
	TransactionDescriptionStatus
	Result *TxInfoG `json:"result"`
}

type CreateInvoiceResponse struct {
	TransactionDescriptionStatus
	Result *TxInfoG `json:"result"`
}

type GetBalanceResponse struct {
	TransactionDescriptionStatus
	Result *BalanceResponse `json:"result"`
}

type BalanceResponse struct {
	Balance *[]BalanceModel `json:"balance"`
}

type TransferParams struct {
	Sender     string  `json:"service_sender"`
	Receiver   string  `json:"service_receiver"`
	Amount     float64 `json:"amount"`
	Commission float64 `json:"commission"`
	Extra      string  `json:"extra"`
}

type TransferResponse struct {
	TransactionDescriptionStatus
}

type SaveImageParams struct {
	ImageBase64 string `json:"image_base_64"`
}

type SaveImageCDNParams struct {
	Data string `json:"data"`
}

// -----------------------------------------------------------------------------------------------------------------

type CreateInternalLiquidityParams struct {
	ServiceClientId string  `json:"service_client"`
	Commission      float64 `json:"commission"`
	Amount          float64 `json:"amount"`
	Token           string  `json:"token"`
	Virtual         bool    `json:"virtual"`
}

type CloseInternalLiquidityParams struct {
	TrackerId string `json:"tracker_id"`
}

type WithdrawInternalLiquidityParams struct {
	Commission float64 `json:"commission"`
	Amount     float64 `json:"amount"`
	TrackerId  string  `json:"-"`
}

type CloseInternalLiquidityResponse struct {
	TransactionDescriptionStatus
}

type CreateInternalLiquidityResponse struct {
	TransactionDescriptionStatus
	Result InternalLiquidityModel `json:"result"`
}

type InternalLiquidityModel struct {
	TrackerId string `json:"tracker_id"`
}

type WithdrawInternalLiquidityResponse struct {
	TransactionDescriptionStatus
	Result InternalLiquidityModel `json:"result"`
}

type SaveImageResponse struct {
	Status string `json:"status"`
	File   string `json:"file"`
}
