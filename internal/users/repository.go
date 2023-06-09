package users

import "UsersService/internal/cErrors"

type Repository interface {
	GetUserByNickName(params *GetUserByNickNameParams) (*User, *cErrors.ResponseErrorModel)
	GetUserByEmail(params *GetUserByEmailParams) (*User, *cErrors.ResponseErrorModel)
	GetUserByID(params *GetUserByIDParams) (*User, *cErrors.ResponseErrorModel)
	GetUserNicknameByID(params *GetUserNicknameByIDParams) (string, *cErrors.ResponseErrorModel)

	IsUserBlocked(params *IsUserBlockedParams) (bool, *cErrors.ResponseErrorModel)
	ChangeDefaultUserLanguage(params *ChangeDefaultUserLanguageParams) (bool, *cErrors.ResponseErrorModel)
	GetClientIP(params *GetUserIpParams) ([]string, *cErrors.ResponseErrorModel)
	UpdateUserLastEntry(params *UpdateUserLastEntryParams) (bool, *cErrors.ResponseErrorModel)
	CreateClient(params *CreateClientParamsRepo) (int64, *cErrors.ResponseErrorModel)
	UpdateUserAvatar(params *UpdateUserAvatarParams) (bool, *cErrors.ResponseErrorModel)
	UpdateUserNickName(params *UpdateUserNickNameParams) (bool, *cErrors.ResponseErrorModel)
	UpdateUserBio(params *UpdateUserBioParams) (bool, *cErrors.ResponseErrorModel)
	GetUserIDByNickName(params *GetUserIDByNickNameParams) (int64, *cErrors.ResponseErrorModel)
	GetUserByNickNameWithID(params *GetUserByNicknameWithID) (*User, *cErrors.ResponseErrorModel)
	GetLanguage(params *IsValidLanguageParams) (*Language, *cErrors.ResponseErrorModel)
	GetNicknameChanges(params *UpdateUserNickNameParams) (*[]ChangeNickname, *cErrors.ResponseErrorModel)

	UpdateLastLogin(params *UpdateLastLoginParams) (bool, *cErrors.ResponseErrorModel)
	UpdateUserLastActivity(params *UpdateUserLastActivityParams) (bool, *cErrors.ResponseErrorModel)
	UpdateNicknameChanges(params *ChangeNickname) *cErrors.ResponseErrorModel

	GetActiveNotice(params *GetActiveNoticeParams) (*[]Notice, *cErrors.ResponseErrorModel)
	ReadNotice(params *ReadNoticeParams) *cErrors.ResponseErrorModel
	AddNotice(params *ADDNotice) *cErrors.ResponseErrorModel

	GetRegistration(params *GetRegistrationParams) (int64, *cErrors.ResponseErrorModel)

	AddAuthKyc(params *AuthKycParams) *cErrors.ResponseErrorModel
	GetAuthKyc(params *ClientID) (*AuthKyc, *cErrors.ResponseErrorModel)
	//AddKycAttempt(params *ClientID) *cErrors.ResponseErrorModel
	//GetKycAttempt(params *ClientID) *KycHistory

	CreateClientUuid(params *CreateClientUidParamsRepo) (string, *cErrors.ResponseErrorModel)
	ClientUUID(params *ClientUuidByIDParams) (string, *cErrors.ResponseErrorModel)
	AddClientUser(params *AddClientUserParamsRepo) (string, *cErrors.ResponseErrorModel)
}
