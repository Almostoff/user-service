package auth

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

type UserFromAuth struct {
	Email            string `json:"email" db:"email"`
	Phone            string `json:"phone" db:"phone"`
	Password         string `json:"password" db:"phone"`
	RegistrationDate string `json:"registration_date" db:"registration_date"`
	Hash             string `json:"hash" db:"hash"`
	AuthLevelID      int64  `json:"auth_level_id" db:"auth_level_id"`
	TG               string `json:"tg" db:"tg"`
}

type UserAgent struct {
	ID         int64     `json:"-" db:"id"`
	ClientID   int64     `json:"-" db:"client_id"`
	UA         string    `json:"ua" db:"ua"`
	SignInDate time.Time `json:"sign_in_date" db:"sign_in_date"`
	LogoutDate time.Time `json:"-" db:"logout_date"`
	IP         string    `json:"-" db:"ip"`
	Logout     string    `json:"-" db:"logout"`
}

type ResponseVerification struct {
	PhoneConfirm bool `json:"phone_confirm" db:"phone_confirm"`
	EmailConfirm bool `json:"email_confirm" db:"email_confirm"`
	KycConfirm   bool `json:"kyc_confirm" db:"kyc_confirm"`
	TotpConfirm  bool `json:"totp_confirm" db:"totp_confirm"`
}

type Access struct {
	Access string `json:"Access"`
}

type ResponseAuthClientContent struct {
	RefreshToken string `json:"Refresh"`
	AccessToken  string `json:"Access"`
}

type ChangeTgParams struct {
	ClientID int64  `json:"client_id" db:"client_id"`
	NewTg    string `json:"tg_user_name" db:"tg"`
	TgID     string `json:"tg_user_id" db:"tg_id"`
}

type RecoveryInitParams struct {
	ClientID    int64  `json:"client_id" db:"client_id"`
	LanguageIso string `json:"language_iso"`
}

type AddTotpParams struct {
	ClientID int64 `json:"client_id" db:"client_id"`
}

type VerifyTotpParams struct {
	ClientID int64  `json:"client_id" db:"client_id"`
	Token    string `json:"totp_token"`
}

type LogoutParams struct {
	ClientID int64  `json:"client_id" db:"client_id"`
	UA       string `json:"ua"`
}

type RefreshAccessParams struct {
	ClientID int64  `json:"client_id" db:"client_id"`
	Refresh  string `json:"Refresh"`
	UA       string `json:"ua"`
}

type ValidateAccessTokenParams struct {
	ClientID int64  `json:"client_id" db:"client_id"`
	Access   string `json:"Access"`
	UA       string `json:"ua"`
}

type RecoveryConfirmParams struct {
	Password    string `json:"password"`
	Hash        string `json:"hash"`
	LanguageIso string `json:"language_iso"`
}
type ResponseSuccessModel struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ResponseSuccessKycModel struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Hash    string `json:"hash"`
}

type ResponseGetClientPrivateModel struct {
	Tg string `json:"tg"`
}

type KycParams struct {
	ClientID int64  `json:"clients_id" db:"clients_id"`
	Hash     string `json:"hash"`
	Success  bool
}

type SendCodeToEmailParams struct {
	ClientID    int64  `json:"client_id" db:"client_id"`
	Type        string `json:"type"`
	LanguageIso string `json:"language_iso"`
}

type VerCodeParams struct {
	CodeInput string `json:"code_input"`
	ClientID  int64  `json:"client_id" db:"client_id"`
	Type      string `json:"type,omitempty"`
	//Time      time.Time `json:"withdraw_time"`
	//UA        string    `json:"withdraw_ua"`
	//Ip        string    `json:"withdraw_ip"`
}

type KycInitParams struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type ConfirmEmailParams struct {
	ClientID    int64  `json:"client_id" db:"client_id"`
	LanguageIso string `json:"language_iso"`
	Email       string `json:"email"`
}

type GetUserSessions struct {
	ClientID int64 `json:"client_id" db:"client_id"`
}

type GetClientByIDParams struct {
	ID int64 `json:"id" db:"id"`
}

type ResponseGetClientByID struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseVerification       `json:"data"`
}

type ResponseUserFromAuth struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *UserFromAuth               `json:"data"`
}

type ResponseGetClientVerification struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *users.ResponseVerification `json:"data"`
}

type ResponseClientSignUpDnd struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *users.ResponseSuccessModel `json:"data"`
}

type ResponseRecovery struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseRecoveryConfirm struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseConfirmEmailByHash struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseAddTotp struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseAddTotpModel       `json:"data"`
}

type ResponseAddTotpModel struct {
	AccountName string `json:"account_name"`
	Secret      string `json:"secret"`
	Link        string `json:"qr"`
	File        string `json:"file"`
}

type ResponseVerifyTotp struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseSuccess struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseChangeTg struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseRefreshAccess struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *Access                     `json:"data"`
}

type ResponseValidateAccess struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseLogout struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseGetUserSessions struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *[]UserAgent                `json:"data"`
}

type ResponseGetClientPrivate struct {
	Error *cErrors.ResponseErrorModel    `json:"error"`
	Data  *ResponseGetClientPrivateModel `json:"data"`
}

type ResponseKycConfirmInit struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessKycModel    `json:"data"`
}

type ResponseKycConfirm struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}
