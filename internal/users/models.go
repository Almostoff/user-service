package users

import (
	"UsersService/internal/cErrors"
	"time"
)

type User struct {
	ClientID         int64   `json:"client_id" db:"id"`
	Nickname         string  `json:"nickname" db:"nickname"`
	Email            string  `json:"email" db:"email"`
	Avatar           string  `json:"avatar" db:"avatar"`
	Bio              *string `json:"bio"`
	IsBlocked        bool    `json:"is_blocked" db:"is_blocked"`
	BlockedUntil     string  `json:"blockedUntil" db:"blocked_until"`
	LastEntry        string  `json:"last_entry" db:"last_entry"`
	LastActivity     string  `json:"last_activity" db:"last_activity"`
	RegistrationDate string  `json:"registrationDate" db:"registration_date"`
	Language         string  `json:"language" db:"language"`
	IsDnD            bool    `json:"isDnD" db:"is_dnd"`
	Ip               string  `json:"ip" db:"ip"`
	KYC              bool    `json:"isVerified" db:"-"`
	Merchant         bool    `json:"-" db:"merchant"`
}

type UserAuth struct {
	Email            string `json:"email" db:"email"`
	Phone            string `json:"phone" db:"phone"`
	IsBlocked        bool   `json:"is_blocked" db:"is_blocked"`
	Language         string `json:"language" db:"language"`
	RegistrationDate string `json:"registration_date" db:"registration_date"`
	PhoneConfirm     bool   `json:"phone_confirm" db:"phone_confirm"`
	EmailConfirm     bool   `json:"email_confirm" db:"email_confirm"`
	KycConfirm       bool   `json:"kyc_confirm" db:"kyc_confirm"`
	TG               string `json:"tg" db:"tg"`
}

type ChangeNickname struct {
	ID          int64     `json:"-" db:"id"`
	ClientID    int64     `db:"client_id"`
	OldNickname string    `db:"old_nickname"`
	ChangeTime  time.Time `db:"change_time"`
}

type Notice struct {
	Id         int64  `json:"-" db:"id"`
	ClientID   int64  `json:"client_id" db:"client_id"`
	InternalID string `json:"internal_id" db:"internal_id"`
	Type       string `json:"type" db:"type"`
	IsRead     bool   `json:"is_read" db:"is_read"`
	CreateTime string `json:"create_time" db:"create_time"`
	Info       *NewOrderNotice
}

type KycHistory struct {
	CurrentStatus string    `json:"status" db:"status"`
	CreateTime    time.Time `json:"create_time" db:"create_time"`
	ClientUuid    string    `json:"client_uuid" db:"client_uuid"`
}

type ADDNotice struct {
	Id              int64   `json:"-" db:"id"`
	ClientID        int64   `json:"client_id" db:"client_id"`
	InternalID      string  `json:"internal_id" db:"internal_id"`
	Type            string  `json:"type" db:"type"`
	IsRead          bool    `json:"is_read" db:"is_read"`
	CreateTime      string  `json:"create_time" db:"create_time"`
	AmountTo        float64 `json:"amount_to" db:"amount_to"`
	AmountToToken   string  `json:"amount_to_token" db:"amount_to_token"`
	AmountFrom      float64 `json:"amount_from" db:"amount_from"`
	AmountFromToken string  `json:"amount_from_token" db:"amount_from_token"`
	ContrParty      string  `json:"nickname" db:"nickname"`
}

type NoticeWithoutInfo struct {
	Id         int64  `json:"-" db:"id"`
	ClientID   int64  `json:"client_id" db:"client_id"`
	InternalID string `json:"internal_id" db:"internal_id"`
	Type       string `json:"type" db:"type"`
	IsRead     bool   `json:"is_read" db:"is_read"`
	CreateTime string `json:"create_time" db:"create_time"`
}

type NewOrderNotice struct {
	Id              int64   `json:"-" db:"id"`
	AmountTo        float64 `json:"amount_to" db:"amount_to"`
	AmountToToken   string  `json:"amount_to_token" db:"amount_to_token"`
	AmountFrom      float64 `json:"amount_from" db:"amount_from"`
	AmountFromToken string  `json:"amount_from_token" db:"amount_from_token"`
	ContrParty      string  `json:"nickname" db:"nickname"`
}

type Language struct {
	Id        int64  `json:"-" db:"id"`
	Code      string `db:"code"`
	Language  string `db:"language"`
	Available bool   `db:"available"`
}

type UserStatistic struct {
	Rate              float64 `json:"rate"`
	FeedbacksAmount   int64   `json:"feedbacksAmount"`
	FeedbacksPositive int64   `json:"feedbacksPositive"`
	FeedbacksNegative int64   `json:"feedbacksNegative"`
}

type RecoveryParams struct {
	Email string `json:"email" db:"email"`
}
type RecoveryConfirmParams struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

type ConfirmEmailParams struct {
	Access string `json:"access"`
}

type IsValidLanguageParams struct {
	LanguageIso string `json:"language_iso" db:"code"`
}

type ConfirmEmailAuthParams struct {
	ClientUuid  string `json:"client_uuid" db:"client_uuid"`
	LanguageIso string `json:"language_iso"`
	Email       string `json:"email"`
}

type ClientID struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type ClientUUID struct {
	ClientUUID string `json:"client_uuid" db:"client_uuid"`
}

type VerificationModel struct {
	Email bool `json:"email"`
	Phone bool `json:"phone"`
	KYC   bool `json:"kyc"`
}

type transactionsVolume struct {
	Token string `json:"token"`
	Value int64  `json:"value"`
}

type ResponseVerification struct {
	PhoneConfirm bool `json:"phone_confirm" db:"phone_confirm"`
	EmailConfirm bool `json:"email_confirm" db:"email_confirm"`
	KycConfirm   bool `json:"kyc_confirm" db:"kyc_confirm"`
	TotpConfirm  bool `json:"totp_confirm" db:"totp_confirm"`
}

type Comment struct {
	ClientID         int64     `json:"client_id" db:"id"`
	ClientReviewerId int64     `json:"client_reviewer_id" db:"client_reviewer_id"`
	Rate             bool      `json:"rate" db:"rate"`
	InternalID       string    `json:"internal_id" db:"internal_id"`
	CreatedDate      time.Time `json:"create_date" `
	Text             string    `json:"text" db:"text"`
}

type CommentResponse struct {
	ReviewerNickname string    `json:"reviewerNickname"`
	ReviewerAvatar   string    `json:"reviewer_avatar"`
	Rate             bool      `json:"rate"`
	InternalID       string    `json:"internal_id"`
	CreatedDate      time.Time `json:"create_date" `
	Text             string    `json:"text"`
}

type FullUserInfo struct {
	NickName               string  `json:"nickname"`
	Avatar                 string  `json:"avatar"`
	Bio                    string  `json:"bio"`
	Email                  string  `json:"email"`
	RegistrationDate       string  `json:"registrationDate"`
	LastVisit              string  `json:"lastVisit" `
	LastActivity           string  `json:"last_activity"`
	Rating                 float64 `json:"rating"`
	BlockedUntil           string  `json:"blockedUntil"`
	IsBlocked              bool    `json:"is_blocked"`
	Language               string  `json:"language"`
	Ip                     string  `json:",omitempty"`
	UserStatistic          *UserStatistic
	UserStatisticForOrders *Statistic
	Verification           *ResponseVerification `json:"verification"`
}

type FullMeInfo struct {
	NickName         string `json:"nickname"`
	Avatar           string `json:"avatar"`
	Bio              string `json:"bio"`
	Email            string `json:"email"`
	Tg               string `json:"tg"`
	RegistrationDate string `json:"registrationDate"`
	LastVisit        string `json:"lastVisit" `
	LastActivity     string `json:"last_activity"`
	//Rating           float64 `json:"rating"`
	BlockedUntil string `json:"blockedUntil"`
	IsBlocked    bool   `json:"is_blocked"`
	Language     string `json:"language"`
	Ip           string `json:",omitempty"`
	//UserStatistic          *UserStatistic
	//UserStatisticForOrders *Statistic
	//Verification           *ResponseVerification `json:"verification"`
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

type FullUserInfoWithID struct {
	ClientID         int64   `json:"client_id"`
	NickName         string  `json:"nickname"`
	Avatar           string  `json:"avatar"`
	Bio              string  `json:"bio"`
	Email            string  `json:"email"`
	RegistrationDate string  `json:"registrationDate"`
	LastVisit        string  `json:"lastVisit" `
	Rating           float64 `json:"rating"`
	BlockedUntil     string  `json:"blockedUntil"`
	IsBlocked        bool    `json:"is_blocked"`
	Language         string  `json:"language"`
	UserStatistic    *UserStatistic
	Verification     *ResponseVerification `json:"verification"`
	AllReview        *[]CommentResponse    `json:"allReviews"`
}

// ----------------------------------------------------------------------------------------------------

type UpdateUserLastEntryParams struct {
	ClientID int64 `json:"client_id" db:"id"`
}
type UpdateLastLoginParams struct {
	Email string
	Ip    string
}

type UpdateUserLastActivityParams struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type GetRegistrationParams struct {
	FromDate time.Time `json:"date_from"`
	ToDate   time.Time `json:"date_to"`
	Dnd      bool      `json:"dnd"`
}

type GetRegistrationStringParams struct {
	FromDate string `json:"date_from"`
	ToDate   string `json:"date_to"`
	Dnd      bool   `json:"dnd"`
}

type Count struct {
	Count int64 `db:"count"`
}

type GetUserIDByNickNameParams struct {
	Nickname string `json:"nickname" db:"nickname"`
}

type AddCommentParams struct {
}

type LogOutParams struct {
	ClientUUID string `json:"client_uuid" db:"client_uuid"`
	UA         string `json:"ua"`
}

type SignInParams struct {
	Email    string `json:"email"`
	UA       string `json:"ua" db:"ua"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
}

type ClientSignInTGParams struct {
	TgUserName string `json:"tg_user_name"`
	TgUserId   int64  `json:"tg_user_id"`
}

type SignInSAParams struct {
	Email    string `json:"email" db:"email"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	UA       string `json:"ua"`
	Ip       string
	Language string `json:"language"`
}

type SignUpSAParams struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone" db:"phone"`
	TG       string `json:"tg" db:"tg"`
	IsDnd    bool   `json:"is_dnd" db:"is_dnd"`
	UA       string `json:"ua"`
}

type SignUpVirginParams struct {
	Phone      string `json:"phone" db:"phone"`
	NickName   string `json:"nickname" db:"nickname"`
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Password   string `json:"password" db:"password"`
	Email      string `json:"email" db:"email"`
	KeyWord    string `json:"key_word" db:"key_word"`
	Hash       string `json:"hash" db:"hash"`
	TG         string `json:"tg" db:"tg"`
	UA         string `json:"UA"`
	IsDnd      bool   `json:"is_dnd" db:"is_dnd"`
}
type SignUpParams struct {
	Language string `json:"language,omitempty" db:"language"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone"`
	IP       string `json:"ip"`
	UA       string `json:"UA"`
	Password string `json:"password"`
}

type UaParams struct {
	UA string `json:"UA"`
}

type ConfirmEmailByHashParams struct {
	Hash string `json:"hash"`
}
type GetUserByNickNameParams struct {
	Nickname string `json:"nickname" db:"nickname"`
	Page     int64  `json:"page"`
	Limit    int64  `json:"limit"`
	Type     string `json:"type"`
}

type GetUserByNicknameWithID struct {
	Nickname string `json:"nickname" db:"nickname"`
}

type GetAllReviewsParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
	Page       int64  `json:"page"`
	Limit      int64  `json:"limit"`
}

type ResponseGetAuthKycModel struct {
	Auth *AuthKyc
}

type AuthKyc struct {
	ClientUuid string    `json:"client_uuid" db:"client_uuid"`
	CreateTime time.Time `json:"create_time" db:"create_time"`
	Auth       string    `json:"auth" db:"auth_kyc"`
	IsActive   bool      `json:"-" db:"is_active"`
}

type GetUserNicknameByIDParams struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type GetUserByEmailParams struct {
	Email string `json:"email" db:"email"`
}

type ValidateParams struct {
	Signature string
	Message   string
	Public    string
}
type GetUserByIDParams struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type ClientUuidByIDParams struct {
	UserId int64 `json:"user_id" db:"id"`
}

type IsUserBlockedParams struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type GetClientChangeNicknameHistoryParams struct {
	ClientUuid string `json:"client_uuid" db:"client_uuid"`
}

type ResponseReplaceIdByNickname struct {
	Comments *[]CommentResponse
}

type ChangePasswordParams struct {
	Access           string `json:"Access"`
	OldPassword      string `json:"old_password"`
	NewPassword      string `json:"new_password"`
	NewPasswordAgain string `json:"new_password_again"`
	Totp             string `json:"totp_token" size:"6"`
}

type ChangeClientPasswordParams struct {
	ClientUuid       string `json:"client_uuid" db:"client_uuid"`
	OldPassword      string `json:"old_password"`
	NewPassword      string `json:"new_password"`
	NewPasswordAgain string `json:"new_password_again"`
	Totp             string `json:"totp_token" size:"6"`
}

type Stat struct {
	FeedbackPositive       int64   `json:"feedbackPositive"`
	FeedbackNegative       int64   `json:"feedbackNegative"`
	TradesCompleted        float64 `json:"tradesCompleted"`
	TradesCompletedPercent float64 `json:"tradesCompletedPercent"`
}

type GetActiveNoticeParams struct {
	ClientID   int64  `db:"client_id"`
	TypeNotice string `json:"type" db:"type"`
}

type ReadNoticeParams struct {
	ClientID   int64  `db:"client_id"`
	InternalID string `json:"internal_id" db:"internal_id"`
}

type ChangeDefaultUserLanguageParams struct {
	ClientID int64  `json:"client_id" db:"id"`
	Ticker   string `json:"ticker" db:"ticker"`
}

type GetUserIpParams struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type GetUserByAccessTokenParams struct {
	Token string `json:"Access"`
}

type GetUserReviewsByAccessTokenParams struct {
	Access string `json:"Access"`
}

type GetClientRatingParams struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type GetUserIDByAccessTokenParams struct {
	Access string `json:"Access"`
}

type CreateClientParamsRepo struct {
	NickName string    `json:"nick_name" db:"nick_name"`
	Email    string    `json:"email" db:"email"`
	Avatar   string    `json:"avatar" db:"avatar"`
	TimeNow  time.Time `json:"timeNow"`
	Language string    `json:"language" db:"language"`
	IsDnD    bool      `json:"isDnD" db:"is_dnd"`
}

type CreateClientParams struct {
	Language string `json:"language" db:"language"`
	Email    string `json:"email" db:"email"`
	IsDnD    bool   `json:"is_dnd" db:"is_dnd"`
}

type UpdateUserAvatarParams struct {
	ClientID  int64 `json:"client_id"`
	NewAvatar string
}

type UpdateUserNickNameParams struct {
	ClientID    int64  `json:"client_id" db:"id"`
	NewNickName string `json:"new_nickname" db:"nickname" size:"20" min:"4"`
	Refresh     string `json:"refresh"`
	UA          string `json:"ua" db:"ua"`
}

type ResponseGetClientChangeNicknameHistoryModel struct {
}

type UpdateUserBioParams struct {
	ClientID int64  `json:"client_id" db:"id"`
	Bio      string `json:"new_bio" db:"bio"`
}

type CreateClientUidParamsRepo struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type AddClientUserParamsRepo struct {
	UserId     int64  `db:"user_id"`
	ClientUuid string `db:"client_uuid"`
}

// ----------------------------------------------------------------------------------------------------

type ResponseChangeNicknameModel struct {
	Access string `json:"Access"`
}

type HashKycConfirm struct {
	Hash       string `json:"hash"`
	SuccessUrl string `json:"successUrl"`
	Unverified string `json:"unverifiedUrl"`
	CallBack   string `json:"callbackUrl"`
}

type HashKycConfirmBody struct {
	ClientID string `json:"clientId"`
	CallBack string `json:"callbackUrl"`
}

type ResponseSuccessModel struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type AuthKycParams struct {
	AuthToken string `json:"authToken"`
	ClientID  int64  `json:"client_id"`
}

type ResponseGetUserNicknameByIDModel struct {
	Nickname string `json:"nickname"`
}

type ResponseGetUserIDByNickNameModel struct {
	ClientID int64 `json:"client_id" db:"id"`
}

type ResponseGetClientIpModel struct {
	IpList []string `json:"ip_list"`
}

type ResponseIsUserBlockedModel struct {
	IsBlocked    bool   `json:"is_blocked"`
	BlockedUntil string `json:"blocked_until"`
}
type ResponseClientSignUpModel struct {
	AccessToken  string `json:"Access"`
	RefreshToken string `json:"Refresh"`
}
type ResponseClientSignInModel struct {
	RefreshToken string `json:"Refresh"`
	AccessToken  string `json:"Access"`
}

type ResponseGetUserByIDModel struct {
	*User
	*Stat
}
type ResponseGetRegistrationModel struct {
	Count int64 `json:"amount_registration"`
}

type ResponseClientChangePasswordModel struct {
	ResponseSuccessModel
}
type ResponseClientChangeDefaultLanguagesModel struct {
}

type ResponseAddCommentModel struct {
}

type ResponseGetClientRatingModel struct {
	Rate                float64               `json:"rate"`
	FeedbacksAmount     int64                 `json:"feedbacksAmount"`
	FeedbacksPositive   int64                 `json:"feedbacksPositive"`
	FeedbacksNegative   int64                 `json:"feedbacksNegative"`
	FirstTradeDate      string                `json:"firstTradeDate"`
	TradeCompleted      int64                 `json:"tradeCompleted"`
	AverageEscrowTime   int64                 `json:"averageEscrowTime"`
	AverageResponseTime int64                 `json:"averageResponseTime"`
	TransactionsVolume  *[]transactionsVolume `json:"transactionsVolume"`
}

type ResponseCreateClientModel struct {
	ClientID int64  `json:"client_id" db:"id"`
	NickName string `json:"nickname" db:"nickname"`
}

// ----------------------------------------------------------------------------------------------------

type GetClientIpListResponse struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseGetClientIpModel   `json:"data"`
}

type ResponseGetClientByNickName struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *FullUserInfo               `json:"data"`
}

type ResponseGetClientByEmail struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *User                       `json:"data"`
}

type ResponseGetUserByAccessToken struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *FullMeInfo                 `json:"data"`
}

type ResponseGetUserByAccessTokenWithID struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *FullUserInfoWithID         `json:"data"`
}

type ResponseGetUserByID struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseGetUserByIDModel   `json:"data"`
}

type ResponseUpdateUserLastEntry struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseIsUserBlocked struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseIsUserBlockedModel `json:"data"`
}

type ResponseClientSignUp struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseClientSignUpModel  `json:"data"`
}

type ResponseClientSignIn struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseClientSignInModel  `json:"data"`
}

type ResponseClientChangePassword struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseClientChangeDefaultLanguage struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseLogOut struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseGetClientRating struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *UserStatistic              `json:"data"`
}

type ResponseGetUserIDByAccessToken struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ClientID                   `json:"data"`
}

type ResponseGetUserUUIDByAccessToken struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ClientUUID                 `json:"data"`
}

type ResponseValidate struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseAddComment struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseAddCommentModel    `json:"data"`
}

type ResponseCreateClient struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseCreateClientModel  `json:"data"`
}

type ResponseGetUserByNicknameWithID struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *User                       `json:"data"`
}

type ResponseUpdateUserAvatar struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseConfirmEmailByHash struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseGetUserNicknameByID struct {
	Error *cErrors.ResponseErrorModel       `json:"error"`
	Data  *ResponseGetUserNicknameByIDModel `json:"data"`
}

type ResponseUpdateUserNickName struct {
	Error *cErrors.ResponseErrorModel  `json:"error"`
	Data  *ResponseChangeNicknameModel `json:"data"`
}

type ResponseGetUserIDByNickName struct {
	Error *cErrors.ResponseErrorModel       `json:"error"`
	Data  *ResponseGetUserIDByNickNameModel `json:"data"`
}

type ResponseUpdateUserBio struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseGetUserReviewsByAccessToken struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *[]CommentResponse          `json:"data"`
}

type ResponseGetUserReviewsByNickname struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *[]CommentResponse          `json:"data"`
}
type ResponseGetAllReviews struct {
	Error *cErrors.ResponseErrorModel  `json:"error"`
	Data  *ResponseGetAllReviewsModels `json:"data"`
}

type ResponseRecovery struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseRecoveryConfirm struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseConfirmEmail struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseUpdateUserLastActivity struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseReadNotice struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseGetActiveNotice struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *[]Notice                   `json:"data"`
}

type ResponseAddNotice struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseConfirmKycModel struct {
	AuthToken string `json:"authToken"`
}

type ResponseConfirmKyc struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseConfirmKycModel    `json:"data"`
}

type ResponseGetClientChangeNicknameHistory struct {
	Error *cErrors.ResponseErrorModel                  `json:"error"`
	Data  *ResponseGetClientChangeNicknameHistoryModel `json:"data"`
}

type ResponseGetAllReviewsModels struct {
	Total    int64              `json:"total"`
	Pages    int64              `json:"pages"`
	Comments *[]CommentResponse `json:"comments"`
}

type ResponseGetRegistration struct {
	Error *cErrors.ResponseErrorModel   `json:"error"`
	Data  *ResponseGetRegistrationModel `json:"data"`
}

type ResponseGetAuthKyc struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseGetAuthKycModel    `json:"data"`
}

type KYCResponse struct {
	ClientId    string `json:"clientId"`
	ScanRef     string `json:"scanRef"`
	ExternalRef string `json:"externalRef"`
	Platform    string `json:"platform"`
	StartTime   int    `json:"startTime"`
	FinishTime  int    `json:"finishTime"`
	Status      struct {
		Overall          string        `json:"overall"`
		SuspicionReasons []interface{} `json:"suspicionReasons"`
		MismatchTags     []interface{} `json:"mismatchTags"`
		AutoDocument     string        `json:"autoDocument"`
		AutoFace         string        `json:"autoFace"`
		ManualDocument   string        `json:"manualDocument"`
		ManualFace       string        `json:"manualFace"`
	} `json:"status"`
	Data struct {
		SelectedCountry     string      `json:"selectedCountry"`
		DocFirstName        string      `json:"docFirstName"`
		DocLastName         string      `json:"docLastName"`
		DocNumber           string      `json:"docNumber"`
		DocPersonalCode     string      `json:"docPersonalCode"`
		DocExpiry           string      `json:"docExpiry"`
		DocDob              string      `json:"docDob"`
		DocType             string      `json:"docType"`
		DocSex              string      `json:"docSex"`
		DocNationality      string      `json:"docNationality"`
		DocIssuingCountry   string      `json:"docIssuingCountry"`
		ManuallyDataChanged bool        `json:"manuallyDataChanged"`
		OrgFirstName        string      `json:"orgFirstName"`
		OrgLastName         string      `json:"orgLastName"`
		OrgNationality      string      `json:"orgNationality"`
		OrgBirthPlace       string      `json:"orgBirthPlace"`
		OrgAuthority        interface{} `json:"orgAuthority"`
		OrgAddress          interface{} `json:"orgAddress"`
	} `json:"data"`
	FileUrls struct {
		FRONT string `json:"FRONT"`
		BACK  string `json:"BACK"`
		FACE  string `json:"FACE"`
	} `json:"fileUrls"`
	AML []struct {
		Status struct {
			ServiceSuspected bool   `json:"serviceSuspected"`
			CheckSuccessful  bool   `json:"checkSuccessful"`
			ServiceFound     bool   `json:"serviceFound"`
			ServiceUsed      bool   `json:"serviceUsed"`
			OverallStatus    string `json:"overallStatus"`
		} `json:"status"`
		Data             []interface{} `json:"data"`
		ServiceName      string        `json:"serviceName"`
		ServiceGroupType string        `json:"serviceGroupType"`
		Uid              string        `json:"uid"`
		ErrorMessage     interface{}   `json:"errorMessage"`
	} `json:"AML"`
	LID []struct {
		Status struct {
			ServiceSuspected bool   `json:"serviceSuspected"`
			CheckSuccessful  bool   `json:"checkSuccessful"`
			ServiceFound     bool   `json:"serviceFound"`
			ServiceUsed      bool   `json:"serviceUsed"`
			OverallStatus    string `json:"overallStatus"`
		} `json:"status"`
		Data             []interface{} `json:"data"`
		ServiceName      string        `json:"serviceName"`
		ServiceGroupType string        `json:"serviceGroupType"`
		Uid              string        `json:"uid"`
		ErrorMessage     interface{}   `json:"errorMessage"`
	} `json:"LID"`
}

// ----------------------------------------------------------------------------------------------------
