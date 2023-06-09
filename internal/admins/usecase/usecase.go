package usecase

import (
	"UsersService/internal/admins"
	"UsersService/internal/cErrors"
	"UsersService/internal/model"
	"UsersService/internal/users"
	"UsersService/pkg/auth"
	"UsersService/pkg/logger"
	"UsersService/pkg/secure"
	"UsersService/pkg/utils"
	"fmt"
	"github.com/golang-jwt/jwt"
	"sync"
)

type AdminUseCase struct {
	logger *logger.ApiLogger
	repo   admins.Repository
	shield *secure.Shield
	authSR auth.ServiceAuth
}

const (
	Internal   string = "Users: Internal Server Error"
	BadRequest string = "Users: Bad Request"
	Success    string = "Users: Success"
)

func NewAdminUseCase(logger *logger.ApiLogger, repo admins.Repository, shield *secure.Shield, authSR auth.ServiceAuth) admins.AdminCase {
	return &AdminUseCase{logger: logger, repo: repo, shield: shield, authSR: authSR}
}

func (a AdminUseCase) ChangeDndNickname(params *admins.ChangeDndNicknameParams) *model.ResponseStandard {
	cErr := a.repo.ChangeDndNickname(params)
	var ok bool
	if cErr.InternalCode == 0 {
		ok = true
	}
	return &model.ResponseStandard{
		Error:   cErr,
		Success: &model.ResponseSuccessModel{Success: ok},
	}
}

func (a AdminUseCase) GetAllClientKyc(params *admins.Search) *admins.ResponseGetAllClientKyc {
	data, err := a.repo.GetAuthKycForAdmin(params)
	return &admins.ResponseGetAllClientKyc{
		Error: err,
		Data:  &admins.ResponseGetAllClientKycModel{KysList: data},
	}
}

func (a AdminUseCase) RefreshKyc(params *model.ClientID) *model.ResponseStandard {
	cErr := a.repo.DeleteKycToken(params)
	var ok bool
	if cErr.InternalCode == 0 {
		ok = true
	}
	return &model.ResponseStandard{
		Error:   cErr,
		Success: &model.ResponseSuccessModel{Success: ok},
	}
}

func (a AdminUseCase) GetAllBlockUsers(params *admins.GetAllBlockUsersParams) *admins.ResponseGetAllBlockUsers {
	_users, cErr := a.repo.GetAllBlockUsers(params)
	us := *_users.Users
	var wg sync.WaitGroup
	wg.Add(len(us))
	for i, v := range us {
		go func(i int, v users.User) {
			defer wg.Done()
			res, err := a.authSR.GetClientVerification(&users.ClientID{ClientID: v.ClientID})
			if err != nil {
				us[i].KYC = false
			} else {
				us[i].KYC = res.Data.KycConfirm
			}
		}(i, v)
	}
	wg.Wait()
	return &admins.ResponseGetAllBlockUsers{
		Data:  _users,
		Error: cErr,
	}
}

func (a AdminUseCase) GetAdminIDByNickname(params *admins.GetAdminByAccessTokenParams) *admins.ResponseGetAdminIDByNickname {
	nickname := a.decodeToken(params.Access)
	id, err := a.repo.GetAdminIDByNickname(&admins.GetAdminIDByNicknameParams{Nickname: nickname})
	return &admins.ResponseGetAdminIDByNickname{
		Data:  &admins.ResponseGetAdminIDByNicknameModel{ClientID: id},
		Error: err,
	}
}

func (a AdminUseCase) GetAdminByAccessToken(params *admins.GetAdminByAccessTokenParams) *admins.ResponseGetAdminByAccessToken {
	email := a.decodeToken(params.Access)

	if email == "" {
		return &admins.ResponseGetAdminByAccessToken{
			Data: &admins.Admin{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_GetAdminByAccessToken_FailedDecode,
				StandartCode: cErrors.StatusConflict,
				Message:      BadRequest,
			},
		}
	}
	//var paramsForRepo = &admins.GetAdminByEmailParams{Email: email}
	admin, err := a.repo.GetAdminByID(&admins.GetAdminByIDParams{})
	if err != nil {
		return &admins.ResponseGetAdminByAccessToken{
			Data: &admins.Admin{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_GetAdinByEmail_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	return &admins.ResponseGetAdminByAccessToken{
		Data:  admin,
		Error: err,
	}
}

func (a AdminUseCase) GetAdminRoles(params *admins.GetAdminRoleParams) *admins.ResponseGetAdminRole {
	roles, err := a.repo.GetAdminRoles(params)
	if err != nil {
		fmt.Println(err)
		return &admins.ResponseGetAdminRole{
			Data: &admins.ResponseGetAdminRoleModel{Roles: roles},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_GetAdminRoles_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	if len(roles) == 0 {
		return &admins.ResponseGetAdminRole{
			Data: &admins.ResponseGetAdminRoleModel{Roles: roles},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_GetAdminRoles_NoSuchAdmin,
				StandartCode: cErrors.StatusBadRequest,
				Message:      BadRequest,
			},
		}
	}
	return &admins.ResponseGetAdminRole{
		Data:  &admins.ResponseGetAdminRoleModel{Roles: roles},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (a AdminUseCase) AdminSignIn(params *admins.AdminSignInParams) *admins.ResponseAdminSignIn {
	admin, errr := a.repo.GetAdminByNickname(&admins.GetAdminByNicknameParams{Nickname: params.Nickname})
	if errr != nil {
		fmt.Println(errr)
		return &admins.ResponseAdminSignIn{
			Data:  &admins.ResponseAdminSignInModel{},
			Error: errr,
		}
	}
	response, err := a.authSR.AdminSignIn(&admins.AdminSignInParamsAS{
		AdminID:        admin.ID,
		Password:       params.Password,
		SecondPassword: params.SecondPassword,
	})
	if err != nil {
		return &admins.ResponseAdminSignIn{
			Data: &admins.ResponseAdminSignInModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.AdminsAuthServiceBadReq,
				StandartCode: cErrors.StatusInternalServerError,
			},
		}
	}
	return &admins.ResponseAdminSignIn{
		Data:  &admins.ResponseAdminSignInModel{RefreshToken: response.Data.RefreshToken, AccessToken: response.Data.AccessToken},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (a AdminUseCase) UpdateAdminLastEntry(params *admins.UpdateAdminLastEntryParams) *admins.ResponseUpdateAdminLastEntry {
	ok, err := a.repo.UpdateAdminLastEntry(params)
	if err != nil {
		return &admins.ResponseUpdateAdminLastEntry{
			Data: &admins.ResponseSuccessModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_UpdateAdminLastEntry_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	return &admins.ResponseUpdateAdminLastEntry{
		Data:  &admins.ResponseSuccessModel{Success: ok},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (a AdminUseCase) decodeToken(token string) string {
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	var nickname string
	nickname = fmt.Sprint(claims["nickname"])
	fmt.Println(nickname)
	return nickname
}

func (a AdminUseCase) BlockUser(params *admins.BlockUserParams) *admins.ResponseBlockUser {
	res, err := a.repo.ChangeBlock(&admins.ChangeBlockParams{
		Block:        true,
		ClientID:     params.ClientID,
		BlockedUntil: params.BlockedUntil,
	})
	return &admins.ResponseBlockUser{
		Error: err,
		Data:  &admins.ResponseSuccessModel{Success: res},
	}
}

func (a AdminUseCase) UnBlockUser(params *admins.UnBlockUserParams) *admins.ResponseUnBlockUser {
	res, err := a.repo.ChangeBlock(&admins.ChangeBlockParams{
		Block:        false,
		ClientID:     params.ClientID,
		BlockedUntil: utils.GetEuropeTime(),
	})
	return &admins.ResponseUnBlockUser{
		Error: err,
		Data:  &admins.ResponseSuccessModel{Success: res},
	}
}
