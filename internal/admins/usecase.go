package admins

import (
	"UsersService/internal/model"
)

type AdminCase interface {
	UpdateAdminLastEntry(params *UpdateAdminLastEntryParams) *ResponseUpdateAdminLastEntry
	GetAdminByAccessToken(params *GetAdminByAccessTokenParams) *ResponseGetAdminByAccessToken
	AdminSignIn(params *AdminSignInParams) *ResponseAdminSignIn
	GetAdminRoles(params *GetAdminRoleParams) *ResponseGetAdminRole
	GetAdminIDByNickname(params *GetAdminByAccessTokenParams) *ResponseGetAdminIDByNickname

	BlockUser(params *BlockUserParams) *ResponseBlockUser
	UnBlockUser(params *UnBlockUserParams) *ResponseUnBlockUser
	GetAllBlockUsers(params *GetAllBlockUsersParams) *ResponseGetAllBlockUsers

	GetAllClientKyc(params *Search) *ResponseGetAllClientKyc
	RefreshKyc(params *model.ClientID) *model.ResponseStandard

	ChangeDndNickname(params *ChangeDndNicknameParams) *model.ResponseStandard
}
