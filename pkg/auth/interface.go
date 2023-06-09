package auth

import (
	"UsersService/internal/admins"
	"UsersService/internal/users"
)

type ServiceAuth interface {
	Logout(params *LogoutParams) *ResponseLogout
	ClientSignUp(params *users.SignUpSAParams) (*users.ResponseClientSignUp, error)
	ClientSignIn(params *users.SignInSAParams) (*users.ResponseClientSignIn, error)
	ClientSignInTg(params *users.ClientSignInTGParams) (*users.ResponseClientSignIn, error)

	CreateDndClient(params *users.SignUpSAParams) (*users.ResponseClientSignUp, error)

	GetClient(params *users.ClientID) (*ResponseUserFromAuth, error)
	GetClientVerification(params *users.ClientID) (*ResponseGetClientVerification, error)
	GetClientPrivate(params *users.ClientID) *ResponseGetClientPrivate
	GetUserSessions(params *GetUserSessions) *ResponseGetUserSessions

	AdminSignIn(params *admins.AdminSignInParamsAS) (*admins.ResponseAdminSignIn, error)
	ConfirmMailReq(params *admins.ConfirmMailReqParams) error
	ConfirmEmailByHash(params *users.ConfirmEmailByHashParams) *ResponseConfirmEmailByHash

	RecoveryInit(params *RecoveryInitParams) *ResponseRecovery
	RecoveryConfirm(params *RecoveryConfirmParams) *ResponseRecoveryConfirm
	ConfirmEmail(params *users.ConfirmEmailAuthParams) *users.ResponseConfirmEmail

	AddTotp(params *AddTotpParams) *ResponseAddTotp
	VerifyTotp(params *VerifyTotpParams) *ResponseVerifyTotp
	VerifyTotpInit(params *VerifyTotpParams) *ResponseVerifyTotp
	ClientChangePassword(params *users.ChangeClientPasswordParams) (*users.ResponseClientChangePassword, error)
	ChangeTg(params *ChangeTgParams) *ResponseChangeTg
	ChangeNickname(params *users.UpdateUserNickNameParams) *users.ResponseUpdateUserNickName

	RefreshAccess(params *RefreshAccessParams) *ResponseRefreshAccess
	ValidateAccess(params *ValidateAccessTokenParams) *ResponseValidateAccess

	ConfirmKycInit(params *KycInitParams) *ResponseKycConfirmInit
	ConfirmKyc(params *KycParams) *ResponseKycConfirm

	SendCodeToEmail(params *SendCodeToEmailParams) *ResponseSuccess
	CheckCodeFromEmail(params *VerCodeParams) *ResponseSuccess
}
