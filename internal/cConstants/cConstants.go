package cConstants

const (
	IotaCounter           int64  = 20_000
	MaxBodyLimit          int64  = 1024
	MaxBodyLimitForAvatar int64  = 2097152
	MAxSizeStringStruct   string = "500"
	MinSizeStringStruct   string = "5"
)

var (
	AuthService   = "auth_service"
	SsoService    = "sso_service"
	User          = "user"
	RatingService = "rating_service"
)

const (
	ClientSignUp                   string = "/clients/sign_up"
	ClientSignIn                   string = "/clients/sign_in"
	ClientSignInTg                 string = "/clients/sign_in_tg"
	ClientChangePassword           string = "/clients/change/password"
	ClientChangeTg                 string = "/clients/change/tg"
	ChangeNickname                 string = "/clients/change/nickname"
	ClientGetClient                string = "/clients/get_client"
	GetClientVerification          string = "/clients/get_auth_level"
	RatingClientRatingFull         string = "/rating/get_clients_rating_full"
	RatingClientRating             string = "/rating/get_clients_rating"
	RatingUpdateCommentAdmin       string = "/rating/update_comment"
	RatingUpdateCommentClient      string = "/rating/update_comment_request"
	RatingAddComment               string = "/rating/add_comment"
	AdminSignIn                    string = "/admin/sign_in"
	ConfirmMailReq                 string = "/clients/req_to_confirm_email"
	GetAllReviews                  string = "/rating/get/all_reviews"
	RatingClientStatisticForOrders string = "/rating/get_clients_statistics"
	ConfirmEmailByHash             string = "/clients/confirm/verify/email"
	ConfirmEmailReq                string = "/clients/req_to_confirm_email"
	AuthRecoveryInit               string = "clients/recovery/"
	AuthRecoveryConfirm            string = "user/recover/password/:hash"
	AuthRecoveryPasswordByEmail    string = "/clients/recovery/confirm_by_email"
	AddTotp                        string = "/clients/confirm/add/totp"
	VerTotp                        string = "/clients/confirm/verify/totp"
	VerTotpInit                    string = "/clients/confirm/verify/totp_init"
	RefreshAccess                  string = "/clients/token/update/access"
	ValidateAccess                 string = "/clients/token/validate/access"
	Logout                         string = "/clients/logout"
	GetUserSessions                string = "/clients/get/sessions"
	IsCommentExist                 string = "/rating/comment/exist"
	CreateDndClient                string = "/clients/sign_up_dnd"
	GetClientPrivate               string = "/clients/get/private"
	KycConf                        string = "/clients/confirm/verify/kyc"
	KycConfInit                    string = "/clients/confirm/verify/kyc_init"

	SendCodeToEmail   string = "/clients/confirm/verify/code_init"
	CheckCodeForEmail string = "/clients/confirm/verify/code"

	SsoClientSignUp            string = "/sso/sign_up"
	SsoClientSignIn            string = "/sso/sign_in"
	SsoClientVerification      string = "/sso/my/get/verify_level"
	SsoGetClientPrivate        string = "/sso/my/get/private_info"
	SsoConfirmMailReq          string = "/sso/my/confirm/add/email"
	SsoVerTotpInit             string = "/sso/my/confirm/verify/totp_init"
	SsoSendCodeToEmail         string = "/sso/my/confirm/add/code"
	SsoConfirmEmailReq         string = "/sso/my/confirm/add/email"
	SsoKycConfInit             string = "/sso/my/confirm/add/kyc"
	SsoVerTotp                 string = "/sso/my/confirm/verify/totp"
	SsoRefreshAccess           string = "/sso/token/update/access"
	SsoValidateAccess          string = "/sso/token/validate/access"
	SsoRecoveryPasswordByEmail string = "/sso/recovery/by_email"
	SsoClientChangePassword    string = "/clients/change/password"
)

var (
	False bool = false
	True  bool = true
)

const (
	CheckIfExistsQuery string = "SELECT * FROM wallet WHERE is_active = true AND user_id = $1 AND network = $2 AND token = $3"
	CreateAddressQuery string = "INSERT INTO wallet (time_create, is_active, address, user_id, network, token) VALUES ($1, $2, $3, $4, $5, $6)"
	GetWalletQuery     string = "SELECT * FROM wallet WHERE is_active = true AND user_id = $1 AND network = $2 AND token = $3"

	GetInnerConnectionQuery string = "SELECT * FROM inner_connection WHERE name = $1"
	GetServiceByPublicQuery string = "SELECT * FROM inner_connection WHERE public = $1"
)
