package rating

import (
	"UsersService/config"
	"UsersService/internal/cErrors"
	"UsersService/internal/users"
	"time"
)

type GetClientParams struct {
	Config  *config.Config
	Public  *string
	Private *string
	BaseUrl *string
}

type GetClientSRRatingParams struct {
	ClientID int64 `json:"client_id" db:"client_id"`
}

type GetAllReviewsParams struct {
	ClientID int64  `json:"client_id" db:"client_id"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
	Type     string `json:"type"`
}

type UpdateCommentAdminParams struct {
	HistoryId   int64    `json:"history_id" db:"history_id"`
	AdminId     int64    `json:"admin_id" db:"admin_id"`
	AdminRole   []string `json:"admin_role"`
	IsSatisfied bool     `json:"is_satisfied" db:"is_satisfied"`
}

type UpdateCommentClientParams struct {
	ClientID    int64  `json:"client_id" db:"client_id"`
	Rate        bool   `json:"rate" db:"rate"`
	SwapID      int64  `json:"swap_id" db:"swap_id"`
	CreatedDate string `json:"created_date" `
	Text        string `json:"text" db:"text"`
}

type AddCommentParams struct {
	ClientID    int64     `json:"client_id" db:"client_id"`
	ReviewerID  int64     `json:"client_reviewer_id"`
	Rate        bool      `json:"rate" db:"rate"`
	InternalID  string    `json:"internal_id" db:"internal_id"`
	CreatedDate time.Time `json:"created_date" `
	Text        string    `json:"text" db:"text"`
}

type TransactionsVolume struct {
	Token string `json:"token"`
	Value int64  `json:"value"`
}

type IsCommentExistParams struct {
	ClientID   int64  `json:"client_id" db:"client_id"`
	InternalID string `json:"internal_id" db:"internal_id"`
}

type ResponseGetClientSRRatingModel struct {
	Rate              float64 `json:"rate"`
	FeedbacksAmount   int64   `json:"feedbacksAmount"`
	FeedbacksPositive int64   `json:"feedbacksPositive"`
	FeedbacksNegative int64   `json:"feedbacksNegative"`
}

type ResponseSuccessModel struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ResponseGetClientSRRating struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *users.UserStatistic        `json:"data"`
}

type ResponseUpdateComment struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseUpdateCommentClient struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseAddComment struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

//
//type ResponseGetAllReviews struct {
//	Error *cErrors.ResponseErrorModel `json:"error"`
//	Data  *Comments                   `json:"data"`
//}

type Comments struct {
	Data  *[]users.Comment            `json:"data"`
	Error *cErrors.ResponseErrorModel `json:"error"`
}

type ResponseGetClientRatingForOrders struct {
}

type Statistic struct {
	PercentDoneOrders  float64    `json:"percent_done_orders"`
	Orders             int64      `json:"orders" db:"orders"`
	AveragePMTimeTaker float64    `json:"average_pm_time_taker" db:"average_pm_time"`
	AveragePMTimeMaker float64    `json:"average_pm_time_maker" db:"average_pm_time"`
	FirstPM            *time.Time `json:"first_pm" db:"first_pm"`
	Amount             float64    `json:"amount" db:"amount"`
	Amount30           float64    `json:"amount_30" db:"amount_30"`
	Orders30           int64      `json:"orders_30" db:"orders_30"`
	Buy                int64      `json:"buy"`
	Sell               int64      `json:"sell"`
}

type ResponseGetClientStatistics struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *Statistic                  `json:"data"`
}

type ResponseGetClientStatisticsModel struct {
	*Statistic
}

type ResponseGetAllReviews struct {
	Error *cErrors.ResponseErrorModel  `json:"error"`
	Data  *ResponseGetAllReviewsModels `json:"data"`
}

type ResponseIsCommentExist struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseGetAllReviewsModels struct {
	Total    int64            `json:"total"`
	Pages    int64            `json:"pages"`
	Comments *[]users.Comment `json:"comments"`
}

type Comment struct {
	ClientID         int64  `json:"client_id" db:"client_id"`
	ClientReviewerId int64  `json:"client_reviewer_id" db:"client_reviewer_id"`
	Rate             bool   `json:"rate" db:"rate"`
	SwapID           int64  `json:"swap_id" db:"swap_id"`
	CreatedDate      string `json:"created_date" `
	Text             string `json:"text" db:"text"`
}
