package cErrors

type ResponseErrorModel struct {
	InternalCode int64  `json:"internal_error_code,omitempty"`
	StandartCode int64  `json:"error_code,omitempty"`
	Message      string `json:"message,omitempty"`
}

func (r ResponseErrorModel) Error() string {
	//TODO implement me
	panic("implement me")
}

//func (r ResponseErrorModel) Error() string {
//	//TODO implement me
//	panic("implement me")
//}

const (
	StatusContinue                      = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols            = 101 // RFC 7231, 6.2.2
	StatusProcessing                    = 102 // RFC 2518, 10.1
	StatusEarlyHints                    = 103 // RFC 8297
	StatusOK                            = 200 // RFC 7231, 6.3.1
	StatusCreated                       = 201 // RFC 7231, 6.3.2
	StatusAccepted                      = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInformation   = 203 // RFC 7231, 6.3.4
	StatusNoContent                     = 204 // RFC 7231, 6.3.5
	StatusResetContent                  = 205 // RFC 7231, 6.3.6
	StatusPartialContent                = 206 // RFC 7233, 4.1
	StatusMultiStatus                   = 207 // RFC 4918, 11.1
	StatusAlreadyReported               = 208 // RFC 5842, 7.1
	StatusIMUsed                        = 226 // RFC 3229, 10.4.1
	StatusMultipleChoices               = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently              = 301 // RFC 7231, 6.4.2
	StatusFound                         = 302 // RFC 7231, 6.4.3
	StatusSeeOther                      = 303 // RFC 7231, 6.4.4
	StatusNotModified                   = 304 // RFC 7232, 4.1
	StatusUseProxy                      = 305 // RFC 7231, 6.4.5
	StatusTemporaryRedirect             = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect             = 308 // RFC 7538, 3
	StatusBadRequest                    = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                  = 401 // RFC 7235, 3.1
	StatusPaymentRequired               = 402 // RFC 7231, 6.5.2
	StatusForbidden                     = 403 // RFC 7231, 6.5.3
	StatusNotFound                      = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed              = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                 = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired             = 407 // RFC 7235, 3.2
	StatusRequestTimeout                = 408 // RFC 7231, 6.5.7
	StatusConflict                      = 409 // RFC 7231, 6.5.8
	StatusGone                          = 410 // RFC 7231, 6.5.9
	StatusLengthRequired                = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed            = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge         = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong             = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType          = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable  = 416 // RFC 7233, 4.4
	StatusExpectationFailed             = 417 // RFC 7231, 6.5.14
	StatusTeapot                        = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest            = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity           = 422 // RFC 4918, 11.2
	StatusLocked                        = 423 // RFC 4918, 11.3
	StatusFailedDependency              = 424 // RFC 4918, 11.4
	StatusTooEarly                      = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired               = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired          = 428 // RFC 6585, 3
	StatusTooManyRequests               = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge   = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons    = 451 // RFC 7725, 3
	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)

const (
	Clients_GetClientIP_NoSuchClient = iota + 4000
	Clients_GetClientIP_PG_Error
	Clients_GetUserByAccessToken_FailedDecode
	Clients_GetUserByEmail_PG_Error
	Clients_GetUserByID_PG_Error
	Clients_IsUserBlocked_PG_Error
	Clients_ClientSignUp_PG_Error
	Clients_GetUserByEmail_Repo_PG_Error
	Clients_GetUserByEmail_Repo_No_Such_User
	Clients_GetUserByID_Repo_No_Such_User
	Clients_GetUserByID_Repo_PG_Error
	Clients_IsUserBlocked_Repo_PG_Error
	Clients_ChangeDefaultUserLanguage_Repo_PG_Error
	Clients_GetClientIP_Repo_PG_Error
	Clients_UpdateUserLastEntry_Repo_PG_Error
	Clients_UpdateUserLastEntry_NoSuchUser
	Clients_ClientSignUp_AuthError
	Clients_ClientChangeDefaultLanguage_PGError
	Clients_UpdateUserLastEntry_PGError
	Admins_GetAdminRoles_Repo_PG_Error
	Admins_IsAdminBlocked_Repo_PG_Error
	Admins_GetAdminByEmail_Repo_PG_Error
	Admins_IsAdminBlocked_NoSuchUser
	Admins_GetAdminByID_Repo_PG_Error
	Admins_GetAdminByID_NoSuchUser
	Admins_GetAdminByAccessToken_FailedDecode
	Admins_GetAdinByEmail_PG_Error
	Admins_GetAdminRoles_PG_Error
	Admins_GetAdminRoles_NoSuchAdmin
	Admins_AdminSignIn_PG_Error
	Admins_UpdateAdminLastEntry_PG_Error
	Admins_IsAdminBlocked_NoSuchAdmin
	CreateClient_Repo_PG_Error
	CreateClient_RW_Repo_PG_Error
	AddClientUser_RW_Repo_PG_Error
	CreateClient_RW_ZeroRowInserted
	UpdateUserAvatar_Repo_PG_Error
	UpdateUserNickName_Repo_PG_Error
	UpdateUserBio_Repo_PG_Error
	RatingServiceErr
	RatingServiceErrNil
	Clients_GetUserByNickname_Repo_No_Such_User
	AdminsAuthServiceBadReq
	GetUserByNickNameWithIDRepoPGError
	GetUserByNickNameWithIDRepoLenZero
	PasswordsNotEqual
)
