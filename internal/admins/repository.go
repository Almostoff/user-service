package admins

import (
	"UsersService/internal/cErrors"
	"UsersService/internal/model"
	"UsersService/internal/users"
)

type Repository interface {
	UpdateAdminLastEntry(params *UpdateAdminLastEntryParams) (bool, *cErrors.ResponseErrorModel)
	GetAdminByNickname(params *GetAdminByNicknameParams) (*Admin, *cErrors.ResponseErrorModel)
	GetAdminByID(params *GetAdminByIDParams) (*Admin, *cErrors.ResponseErrorModel)
	IsAdminBlocked(params *IsAdminBlockedParams) (bool, *cErrors.ResponseErrorModel)
	GetAdminRoles(params *GetAdminRoleParams) ([]string, *cErrors.ResponseErrorModel)
	GetAdminIDByNickname(params *GetAdminIDByNicknameParams) (int64, *cErrors.ResponseErrorModel)
	ChangeBlock(params *ChangeBlockParams) (bool, *cErrors.ResponseErrorModel)
	GetAllBlockUsers(params *GetAllBlockUsersParams) (*ResponseGetAllBlockUsersModel, *cErrors.ResponseErrorModel)

	GetAuthKycForAdmin(params *Search) (*[]users.AuthKyc, *cErrors.ResponseErrorModel)
	DeleteKycToken(params *model.ClientID) *cErrors.ResponseErrorModel

	ChangeDndNickname(params *ChangeDndNicknameParams) *cErrors.ResponseErrorModel
}
