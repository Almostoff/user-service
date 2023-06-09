package admins

import (
	"UsersService/internal/cErrors"
	"UsersService/internal/users"
	"time"
)

type Admin struct {
	ID        int64  `json:"id" db:"id"`
	IsBlocked bool   `json:"is_blocked" db:"is_blocked"`
	Nickname  string `json:"nickname" db:"nickname"`
	LastEntry string `json:"last_entry" db:"last_entry"`
	Name      string `json:"name"`
	ProfileID string `json:"profile_id" db:"profile_id"`
	Totp      string `json:"totp_token"`
	Ip        string `json:"ip" db:"ip"`
}

// ----------------------------------------------------------------------------------------------------

type GetAdminByNicknameParams struct {
	Nickname string `json:"nickname" db:"nickname"`
}

type GetAllBlockUsersParams struct {
	Page        int64  `json:"page"`
	Limit       int64  `json:"limit"`
	SearchField string `json:"search_field"`
	IsDnd       string `json:"is_dnd"`
	IsBlocked   string `json:"is_blocked"`
}

type Search struct {
	Page        int64     `json:"page"`
	Limit       int64     `json:"limit"`
	FromDate    time.Time `json:"date_from,omitempty"`
	ToDate      time.Time `json:"date_to,omitempty"`
	SearchField string    `json:"search_field"`
}

type Count struct {
	Count int64 `db:"total_count"`
}

type BlockUserParams struct {
	ClientID     int64     `json:"client_id"`
	BlockedUntil time.Time `json:"blocked_until"`
}

type ChangeBlockParams struct {
	Block        bool
	ClientID     int64
	BlockedUntil time.Time
}

type UnBlockUserParams struct {
	ClientID int64 `json:"client_id"`
}

type AdminSignInParams struct {
	Nickname       string `json:"nickname" db:"nickname"`
	Password       string `json:"password" db:"password"`
	SecondPassword string `json:"second_password" db:"second_password"`
}

type AdminSignInParamsAS struct {
	AdminID        int64  `json:"admin_id" db:"admin_id"`
	Password       string `json:"password" db:"password"`
	SecondPassword string `json:"second_password" db:"second_password"`
}
type GetAdminRoleParams struct {
	ClientID int64 `json:"client_id" db:"client_id"`
}

type RoleParams struct {
	Role string `json:"role" db:"role"`
}

type ChangeDndNicknameParams struct {
	ClientID    int64  `json:"client_id" db:"client_id"`
	NewNickname string `json:"new_nickname"`
}

type UpdateAdminLastEntryParams struct {
	ClientID int64 `json:"client_id" db:"client_id"`
}
type GetAdminByAccessTokenParams struct {
	Access string `json:"Access"`
}
type GetAdminByEmailParams struct {
	Email string `json:"email" db:"email"`
}
type IsAdminBlockedParams struct {
	ClientID int64 `json:"client_id" db:"client_id"`
}

type GetAdminByIDParams struct {
	ClientID int64 `json:"client_id" db:"client_id"`
}

type ConfirmMailReqParams struct {
	ClientUUID  string `json:"client_uuid" db:"client_uuid"`
	LanguageIso string `json:"language_iso"`
}

type GetAdminIDByNicknameParams struct {
	Nickname string `json:"nickname" db:"nickname"`
}

// ----------------------------------------------------------------------------------------------------

type ResponseSuccessModel struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ResponseGetAllClientKycModel struct {
	KysList *[]users.AuthKyc
}

type ResponseAdminSignInModel struct {
	RefreshToken string `json:"Refresh"`
	AccessToken  string `json:"Access"`
}
type ResponseGetAdminRoleModel struct {
	Roles []string `json:"roles"`
}

type ResponseGetAdminIDByNicknameModel struct {
	ClientID int64 `json:"client_id" db:"client_id"`
}

// ----------------------------------------------------------------------------------------------------

type ResponseUpdateAdminLastEntry struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseConfirmMailReq struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseGetAdminByAccessToken struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *Admin                      `json:"data"`
}

type ResponseAdminSignIn struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseAdminSignInModel   `json:"data"`
}

type ResponseGetAdminRole struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseGetAdminRoleModel  `json:"data"`
}

type ResponseGetAdminIDByNickname struct {
	Error *cErrors.ResponseErrorModel        `json:"error"`
	Data  *ResponseGetAdminIDByNicknameModel `json:"data"`
}

type ResponseGetAllBlockUsers struct {
	Error *cErrors.ResponseErrorModel    `json:"error"`
	Data  *ResponseGetAllBlockUsersModel `json:"data"`
}

type ResponseGetAllBlockUsersModel struct {
	Users *[]users.User `json:"users"`
	Total int64         `json:"total"`
	Pages int64         `json:"pages"`
}

type ResponseBlockUser struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseUnBlockUser struct {
	Error *cErrors.ResponseErrorModel `json:"error"`
	Data  *ResponseSuccessModel       `json:"data"`
}

type ResponseGetAllClientKyc struct {
	Error *cErrors.ResponseErrorModel   `json:"error"`
	Data  *ResponseGetAllClientKycModel `json:"data"`
}
