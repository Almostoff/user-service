package repository

import (
	"UsersService/internal/admins"
	"UsersService/internal/cErrors"
	"UsersService/internal/model"
	"UsersService/internal/users"
	"UsersService/pkg/secure"
	"UsersService/pkg/utils"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	shield *secure.Shield
	db     *sqlx.DB
}

const (
	Internal   string = "Users: Internal Server Error"
	BadRequest string = "Users: Bad Request"
	Success    string = "Users: Success"
)

func NewPostgresRepository(db *sqlx.DB, shield *secure.Shield) admins.Repository {
	return &postgresRepository{db: db, shield: shield}
}

func (p postgresRepository) ChangeDndNickname(params *admins.ChangeDndNicknameParams) *cErrors.ResponseErrorModel {
	res, err := p.db.Exec(`UPDATE user_db.public.users set nickname=$1 WHERE id=$2 
                                              AND is_dnd=true AND NOT check_nickname_exists($1)`,
		params.NewNickname, params.ClientID)
	if err != nil {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusBadRequest,
			StandartCode: cErrors.StatusBadRequest,
			Message:      "already exist nickname",
		}
	}
	if insertedRow == 0 {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusBadRequest,
			StandartCode: cErrors.StatusBadRequest,
		}
	}
	return &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) DeleteKycToken(params *model.ClientID) *cErrors.ResponseErrorModel {
	res, err := p.db.Exec("UPDATE user_db.public.kyc_auth set is_active=false WHERE client_id=$1", params.ClientID)
	if err != nil {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusBadRequest,
			StandartCode: cErrors.StatusBadRequest,
		}
	}
	if insertedRow == 0 {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusBadRequest,
			StandartCode: cErrors.StatusBadRequest,
		}
	}
	return &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetAuthKycForAdmin(params *admins.Search) (*[]users.AuthKyc, *cErrors.ResponseErrorModel) {
	var (
		data []users.AuthKyc
		q    string = `SELECT *, u.nickname, u.email
						FROM user_db.public.kyc_auth
						INNER JOIN user_db.public.users u
						on kyc_auth.client_id = u.id
						WHERE is_active = true`
	)
	if params.SearchField != "" {
		q += ` AND (client_id || nickname || email) like '%` + params.SearchField + `%'`
	}
	var t time.Time
	if params.FromDate != t || params.ToDate != t {
		if params.FromDate != t {
			q += fmt.Sprintf(` AND create_time >%s`, params.FromDate)
		}
		if params.ToDate != t {
			q += fmt.Sprintf(` AND create_time <%s`, params.FromDate)
		}
	}

	err := p.db.Select(&data, q)
	if err != nil {
		return nil, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	if len(data) == 0 {
		return nil, &cErrors.ResponseErrorModel{}
	}
	return &data, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetAllBlockUsers(params *admins.GetAllBlockUsersParams) (*admins.ResponseGetAllBlockUsersModel, *cErrors.ResponseErrorModel) {
	var _users []users.User
	q := `SELECT * FROM user_db.public.users`
	q += ` WHERE 1=1`
	if params.SearchField != "" {
		q += ` AND (nickname || email || bio || id) like '%` + params.SearchField + `%'`
	}
	if params.IsBlocked != "" {
		var b bool
		if params.IsBlocked == "blocked" {
			b = true
		}
		if params.SearchField != "" {
			q += fmt.Sprintf(" AND is_blocked=%t", b)
		} else {
			q += fmt.Sprintf(" AND is_blocked=%t", b)
		}
	}
	if params.IsDnd != "" {
		var b bool
		if params.IsDnd == "true" {
			b = true
		}
		q += fmt.Sprintf(" AND is_dnd=%t", b)
	}
	offset := q + ` OFFSET $1 LIMIT $2`
	fmt.Println(offset)
	err := p.db.Select(&_users, offset,
		utils.CalculateOffset(params.Limit, params.Page), params.Limit)
	if err != nil {
		return &admins.ResponseGetAllBlockUsersModel{}, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_GetAdminRoles_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	var count []admins.Count
	err = p.db.Select(&count,
		fmt.Sprintf("SELECT COUNT(*) as total_count FROM (%s) as sub", q))
	return &admins.ResponseGetAllBlockUsersModel{
		Users: &_users,
		Total: count[0].Count,
		Pages: utils.CalculatePages(params.Limit, count[0].Count),
	}, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetAdminIDByNickname(params *admins.GetAdminIDByNicknameParams) (int64, *cErrors.ResponseErrorModel) {
	var id []int64
	err := p.db.Select(&id, "SELECT id FROM user_db.public.admins WHERE nickname=$1", params.Nickname)
	if err != nil {
		return 0, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_GetAdminRoles_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	if len(id) == 0 {
		return 0, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_GetAdminRoles_NoSuchAdmin,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	fmt.Println(id[0])

	return id[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetAdminRoles(params *admins.GetAdminRoleParams) ([]string, *cErrors.ResponseErrorModel) {
	var roles []string
	err := p.db.Select(&roles, "SELECT name FROM user_db.public.roles WHERE id=(SELECT id    FROM user_db.public.admin_roles          WHERE admin_id=$1)", params.ClientID)
	if err != nil {
		return roles, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_GetAdminRoles_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	fmt.Println(roles)
	if len(roles) == 0 {
		return nil, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_GetAdminRoles_NoSuchAdmin,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}

	return roles, nil
}

func (p postgresRepository) IsAdminBlocked(params *admins.IsAdminBlockedParams) (bool, *cErrors.ResponseErrorModel) {
	var data []bool
	err := p.db.Select(&data, "SELECT is_blocked FROM admins.admins WHERE id=$1", params.ClientID)
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_IsAdminBlocked_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if len(data) == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_IsAdminBlocked_NoSuchAdmin,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return data[0], nil
}

func (p postgresRepository) GetAdminByNickname(params *admins.GetAdminByNicknameParams) (*admins.Admin, *cErrors.ResponseErrorModel) {
	var data []admins.Admin
	err := p.db.Select(&data, "SELECT * FROM user_db.public.admins WHERE nickname=$1", params.Nickname)
	if err != nil {
		fmt.Println(err)
		return &admins.Admin{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_GetAdminByEmail_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			}
	}
	if len(data) == 0 {
		return &admins.Admin{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_IsAdminBlocked_NoSuchUser,
				StandartCode: cErrors.StatusBadRequest,
				Message:      BadRequest,
			}
	}
	return &data[0], nil
}

func (p postgresRepository) GetAdminByID(params *admins.GetAdminByIDParams) (*admins.Admin, *cErrors.ResponseErrorModel) {
	var data []admins.Admin
	err := p.db.Select(&data, "SELECT * FROM admins.admins WHERE id=$1", params.ClientID)
	if err != nil {
		return &admins.Admin{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_GetAdminByID_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			}
	}
	if len(data) == 0 {
		return nil, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_GetAdminByID_NoSuchUser,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return &data[0], nil
}

func (p postgresRepository) UpdateAdminLastEntry(params *admins.UpdateAdminLastEntryParams) (bool, *cErrors.ResponseErrorModel) {
	var data []string
	t := time.Now().Add(time.Hour * 3)
	err := p.db.Select(&data, "UPDATE admins.admins SET last_entry=$1 WHERE id=$2", t, params.ClientID)
	if err != nil {
		return false,
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Admins_GetAdminByID_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			}
	}
	if len(data) == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Admins_GetAdminByID_NoSuchUser,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return true, nil
}

func (p postgresRepository) ChangeBlock(params *admins.ChangeBlockParams) (bool, *cErrors.ResponseErrorModel) {

	res, err := p.db.Exec("UPDATE user_db.public.users SET is_blocked=$1, blocked_until=$2 WHERE id=$3",
		params.Block, params.BlockedUntil, params.ClientID)
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if insertedRow == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return true, &cErrors.ResponseErrorModel{}

}

func (p postgresRepository) searchFields(params *admins.ChangeBlockParams) (bool, *cErrors.ResponseErrorModel) {
	panic("need")
}
