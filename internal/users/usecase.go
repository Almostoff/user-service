package users

type UseCase interface {
	CreateClient(params *CreateClientParams) *ResponseCreateClient
	CreateClientDND(params *CreateClientParams) *ResponseCreateClient
	ClientSignIn(params *SignInSAParams) *ResponseClientSignIn
	ClientSignUp(params *SignUpParams) *ResponseClientSignUp
	ClientSignInTG(params *ClientSignInTGParams) *ResponseClientSignIn

	GetUserByNickName(params *GetUserByNickNameParams) *ResponseGetClientByNickName
	GetUserByEmail(params *GetUserByEmailParams) *ResponseGetClientByEmail
	GetUserByID(params *GetUserByIDParams) *ResponseGetUserByID
	GetUserByAccessToken(params *GetUserByAccessTokenParams) *ResponseGetUserByAccessToken

	GetUserByAccessTokenWithID(params *GetUserByAccessTokenParams) *ResponseGetUserByAccessTokenWithID
	GetUserByNicknameWithID(params *GetUserByNicknameWithID) *ResponseGetUserByNicknameWithID

	GetUserIDByAccessToken(params *GetUserIDByAccessTokenParams) *ResponseGetUserIDByAccessToken
	GetClientUUIDByAccessToken(params *GetUserIDByAccessTokenParams) *ResponseGetUserUUIDByAccessToken
	GetUserIDByNickName(params *GetUserIDByNickNameParams) *ResponseGetUserIDByNickName
	GetUserNicknameByID(params *GetUserNicknameByIDParams) *ResponseGetUserNicknameByID

	ClientChangePassword(params *ChangeClientPasswordParams) *ResponseClientChangePassword
	ClientChangeDefaultLanguage(params *ChangeDefaultUserLanguageParams) *ResponseClientChangeDefaultLanguage

	UpdateLastLogin(params *UpdateLastLoginParams)
	UpdateUserBio(params *UpdateUserBioParams) *ResponseUpdateUserBio
	UpdateUserLastActivity(params *UpdateUserLastActivityParams) *ResponseUpdateUserLastActivity
	UpdateUserNickName(params *UpdateUserNickNameParams) *ResponseUpdateUserNickName
	UpdateUserLastEntry(params *UpdateUserLastEntryParams) *ResponseUpdateUserLastEntry
	UpdateUserAvatar(params *UpdateUserAvatarParams) *ResponseUpdateUserAvatar

	GetUserReviewsByNickname(params *GetUserByNickNameParams) *ResponseGetAllReviews
	GetClientRating(params *GetClientRatingParams) *ResponseGetClientRating
	GetClientIP(params *GetUserIpParams) *GetClientIpListResponse

	IsUserBlocked(params *IsUserBlockedParams) *ResponseIsUserBlocked
	Validate(params *ValidateParams) *ResponseValidate

	Recovery(params *RecoveryParams) *ResponseRecovery
	RecoveryConfirm(params *RecoveryConfirmParams) *ResponseRecoveryConfirm

	ConfirmEmail(params *ConfirmEmailParams) *ResponseConfirmEmail
	ConfirmEmailByHash(params *ConfirmEmailByHashParams) *ResponseConfirmEmailByHash

	GetActiveNotice(params *GetActiveNoticeParams) *ResponseGetActiveNotice
	ReadNotice(params *ReadNoticeParams) *ResponseReadNotice
	AddNotice(params *ADDNotice) *ResponseAddNotice

	GetRegistration(params *GetRegistrationParams) *ResponseGetRegistration

	ConfirmKyc(params *HashKycConfirm) *ResponseConfirmKyc
	AddAuthKyc(params *AuthKycParams)
	GetAuthKyc(params *ClientID) *ResponseGetAuthKyc
	//AddKycAttempt(params *KycHistory) *cErrors.ResponseErrorModel
	//GetKycAttempt(params *ClientID) *KycHistory
}
