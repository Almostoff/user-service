package repository

import (
	"UsersService/internal/cErrors"
	"UsersService/internal/users"
	"UsersService/pkg/secure"
	"UsersService/pkg/utils"
	"fmt"
	"log"
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

func NewPostgresRepository(db *sqlx.DB, shield *secure.Shield) users.Repository {
	return &postgresRepository{db: db, shield: shield}
}

func (p postgresRepository) GetAuthKyc(params *users.ClientID) (*users.AuthKyc, *cErrors.ResponseErrorModel) {
	var data []users.AuthKyc
	err := p.db.Select(&data, `SELECT * FROM user_db.public.kyc_auth
	WHERE client_id=$1 `, params.ClientID)
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
	return &data[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetActiveNotice(params *users.GetActiveNoticeParams) (*[]users.Notice, *cErrors.ResponseErrorModel) {
	var data []users.NoticeWithoutInfo
	err := p.db.Select(&data, `SELECT * FROM user_db.public.notice 
	WHERE is_read=false AND type=$1 AND client_id=$2 AND extract(epoch from (now() + '3 hour' - create_time)) < 180`,
		params.TypeNotice, params.ClientID)
	if err != nil {
		return nil, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	if len(data) == 0 {
		return &[]users.Notice{}, &cErrors.ResponseErrorModel{}
	}
	var newOrder []users.NewOrderNotice
	var notice []users.Notice
	for _, n := range data {
		err = p.db.Select(&newOrder, `SELECT * FROM user_db.public.notice_new_order 
	WHERE id=$1`,
			n.Id)
		if err != nil {
			return nil, &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			}
		}
		if len(newOrder) == 0 {
			return &[]users.Notice{}, &cErrors.ResponseErrorModel{}
		}
		notice = append(notice, users.Notice{
			ClientID:   n.ClientID,
			InternalID: n.InternalID,
			IsRead:     n.IsRead,
			CreateTime: n.CreateTime,
			Info:       &newOrder[0],
		})
	}

	return &notice, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) AddAuthKyc(params *users.AuthKycParams) *cErrors.ResponseErrorModel {
	_, err := p.db.Exec(`INSERT INTO user_db.public.kyc_auth (client_id, auth_kyc, create_time) VALUES ($1, $2, $3)`,
		params.ClientID, params.AuthToken, utils.GetEuropeTime())
	if err != nil {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	return &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetNicknameChanges(params *users.UpdateUserNickNameParams) (*[]users.ChangeNickname, *cErrors.ResponseErrorModel) {
	var data []users.ChangeNickname
	err := p.db.Select(&data, `SELECT * FROM user_db.public.nickname_history WHERE id=$1`, params.ClientID)
	if err != nil {
		return &data, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}

	return &data, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetUserNicknameByID(params *users.GetUserNicknameByIDParams) (string, *cErrors.ResponseErrorModel) {
	var data []string
	err := p.db.Select(&data, `SELECT nickname FROM user_db.public.users WHERE id=$1`, params.ClientID)
	if err != nil {
		return "", &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	if len(data) == 0 {
		return "", &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusBadRequest,
			StandartCode: cErrors.StatusBadRequest,
			Message:      "no such user",
		}
	}

	return data[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetRegistration(params *users.GetRegistrationParams) (int64, *cErrors.ResponseErrorModel) {
	var t time.Time
	if params.ToDate == t || params.FromDate == t {
		params.FromDate, params.ToDate = utils.GetTodayAndYesterday()
	}
	q := `SELECT count(*) FROM user_db.public.users  WHERE is_dnd=$1 AND registration_date<$2 AND registration_date>$3`
	var count []users.Count
	err := p.db.Select(&count, q, params.Dnd, params.ToDate, params.FromDate)
	if err != nil {
		return 0, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	return count[0].Count, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) ReadNotice(params *users.ReadNoticeParams) *cErrors.ResponseErrorModel {
	res, err := p.db.NamedExec(`UPDATE user_db.public.notice SET is_read=true
	WHERE client_id=:client_id AND internal_id=:internal_id`, params)
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

func (p postgresRepository) AddNotice(params *users.ADDNotice) *cErrors.ResponseErrorModel {
	var result int64
	err := p.db.QueryRow(`INSERT INTO notice (internal_id, client_id, type, is_read, create_time) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id`, params.InternalID, params.ClientID, params.Type, false, params.CreateTime).Scan(&result)
	if err != nil {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}

	_, err = p.db.Exec(`INSERT INTO notice_new_order (id, amount_to, amount_to_token, amount_from,
                                             amount_from_token, nickname) 
	VALUES ($1, $2, $3, $4, $5, $6)`,
		result, params.AmountTo, params.AmountToToken, params.AmountFrom, params.AmountFromToken,
		params.ContrParty)
	if err != nil {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	return &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) UpdateNicknameChanges(params *users.ChangeNickname) *cErrors.ResponseErrorModel {
	_, err := p.db.Exec(`INSERT INTO user_db.public.nickname_history (client_id, old_nickname, change_time) VALUES ($1, $2, $3)`,
		params.ClientID, params.OldNickname, utils.GetEuropeTime())
	if err != nil {
		return &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	return &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetLanguage(params *users.IsValidLanguageParams) (*users.Language, *cErrors.ResponseErrorModel) {
	var data []users.Language
	err := p.db.Select(&data, `SELECT * FROM user_db.public.languages WHERE code=$1`,
		params.LanguageIso)
	if err != nil {
		return nil, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusInternalServerError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	if len(data) == 0 {
		return nil, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.StatusBadRequest,
			StandartCode: cErrors.StatusBadRequest,
			Message:      "no such language",
		}
	}

	return &data[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) UpdateUserLastActivity(params *users.UpdateUserLastActivityParams) (bool, *cErrors.ResponseErrorModel) {
	res, err := p.db.Exec("UPDATE user_db.public.users SET last_activity=$1 WHERE id=$2", utils.GetEuropeTime(), params.ClientID)
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if insertedRow == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return true, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) UpdateLastLogin(params *users.UpdateLastLoginParams) (bool, *cErrors.ResponseErrorModel) {
	res, err := p.db.Exec("UPDATE users SET last_entry=$1, ip=$3 WHERE email=$2",
		utils.GetEuropeTime(), params.Email, params.Ip)
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if insertedRow == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return true, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetUserByNickNameWithID(params *users.GetUserByNicknameWithID) (*users.User, *cErrors.ResponseErrorModel) {
	var user []users.User
	err := p.db.Select(&user, `SELECT *  FROM user_db.public.users WHERE nickname=$1`,
		params.Nickname)
	if err != nil {
		return &user[0], &cErrors.ResponseErrorModel{
			InternalCode: cErrors.GetUserByNickNameWithIDRepoPGError,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      fmt.Sprintf(err.Error()),
		}
	}
	if len(user) == 0 {
		return &user[0], &cErrors.ResponseErrorModel{
			InternalCode: cErrors.GetUserByNickNameWithIDRepoLenZero,
			StandartCode: cErrors.StatusBadRequest,
			Message:      "no such user",
		}
	}
	return &user[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) CreateClient(params *users.CreateClientParamsRepo) (int64, *cErrors.ResponseErrorModel) {
	if params.Language == "" {
		params.Language = "ru"
	}
	var lastInsertId int64
	err := p.db.QueryRow(`INSERT INTO user_db.public.users 
    (nickname, email, is_blocked, language, last_entry, registration_date, avatar, is_dnd) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		params.NickName, params.Email, false, params.Language, params.TimeNow, params.TimeNow, params.Avatar, params.IsDnD).Scan(&lastInsertId)
	if err != nil {
		return 0, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.CreateClient_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if lastInsertId == 0 {
		return lastInsertId, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.CreateClient_RW_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	return lastInsertId, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetUserByEmail(params *users.GetUserByEmailParams) (*users.User, *cErrors.ResponseErrorModel) {
	var data []users.User
	err := p.db.Select(&data, "SELECT * FROM users WHERE email=$1", params.Email)
	if err != nil {
		return &users.User{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByEmail_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			}
	}

	if len(data) == 0 {
		return &users.User{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByEmail_Repo_No_Such_User,
				StandartCode: cErrors.StatusBadRequest,
				Message:      BadRequest,
			}
	}
	return &data[0], &cErrors.ResponseErrorModel{}
}
func (p postgresRepository) GetUserByNickName(params *users.GetUserByNickNameParams) (*users.User, *cErrors.ResponseErrorModel) {
	var data []users.User
	err := p.db.Select(&data, "SELECT * FROM user_db.public.users WHERE nickname=$1", params.Nickname)
	if err != nil {
		fmt.Println(err)
		return &users.User{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByEmail_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			}
	}

	if len(data) == 0 {
		return &users.User{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByNickname_Repo_No_Such_User,
				StandartCode: cErrors.StatusBadRequest,
				Message:      BadRequest,
			}
	}
	//fmt.Println(data[0])
	return &data[0], &cErrors.ResponseErrorModel{}

}

func (p postgresRepository) GetUserByID(params *users.GetUserByIDParams) (*users.User, *cErrors.ResponseErrorModel) {
	var data []users.User
	err := p.db.Select(&data, "SELECT * FROM user_db.public.users WHERE id=$1", params.ClientID)
	if err != nil {
		return &users.User{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByID_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			}
	}
	if len(data) == 0 {
		return &users.User{},
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByID_Repo_No_Such_User,
				StandartCode: cErrors.StatusBadRequest,
				Message:      BadRequest,
			}
	}
	return &data[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) IsUserBlocked(params *users.IsUserBlockedParams) (bool, *cErrors.ResponseErrorModel) {
	var data []bool
	err := p.db.Select(&data, "SELECT is_blocked FROM users WHERE id=$1", params.ClientID)
	if err != nil {
		return false,
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_IsUserBlocked_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			}
	}
	return data[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) ChangeDefaultUserLanguage(params *users.ChangeDefaultUserLanguageParams) (bool, *cErrors.ResponseErrorModel) {
	res, err := p.db.Exec("UPDATE user_db.public.users SET language=$1 WHERE id=$2", params.Ticker, params.ClientID)
	if err != nil {
		fmt.Println(err)
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)

		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if insertedRow == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_ChangeDefaultUserLanguage_Repo_PG_Error,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return true, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetClientIP(params *users.GetUserIpParams) ([]string, *cErrors.ResponseErrorModel) {
	var data []string
	err := p.db.Select(&data, "SELECT ip FROM user_db.public.users WHERE id=$1", params.ClientID)
	if err != nil {
		return data, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_GetClientIP_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	return data, nil
}

func (p postgresRepository) UpdateUserLastEntry(params *users.UpdateUserLastEntryParams) (bool, *cErrors.ResponseErrorModel) {
	var data []string
	err := p.db.Select(&data, "UPDATE user_db.public.users SET last_entry=$1 WHERE id=$2", utils.GetEuropeTime(), params.ClientID)
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_UpdateUserLastEntry_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if len(data) == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.Clients_UpdateUserLastEntry_NoSuchUser,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	return true, nil
}

func (p postgresRepository) UpdateUserAvatar(params *users.UpdateUserAvatarParams) (bool, *cErrors.ResponseErrorModel) {
	res, err := p.db.Exec("UPDATE user_db.public.users SET avatar=$1 WHERE id=$2",
		params.NewAvatar, params.ClientID)
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserAvatar_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserAvatar_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      err.Error(),
		}
	}
	if insertedRow == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserAvatar_Repo_PG_Error,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return true, nil
}

func (p postgresRepository) UpdateUserNickName(params *users.UpdateUserNickNameParams) (bool, *cErrors.ResponseErrorModel) {
	res, err := p.db.Exec("UPDATE user_db.public.users SET nickname=$1 WHERE id=$2",
		params.NewNickName, params.ClientID)
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserNickName_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserNickName_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if insertedRow == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserNickName_Repo_PG_Error,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return true, nil
}

func (p postgresRepository) UpdateUserBio(params *users.UpdateUserBioParams) (bool, *cErrors.ResponseErrorModel) {
	res, err := p.db.Exec("UPDATE users SET bio=$1 WHERE id=$2",
		params.Bio, params.ClientID)
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserBio_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	insertedRow, err := res.RowsAffected()
	if err != nil {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserBio_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if insertedRow == 0 {
		return false, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.UpdateUserBio_Repo_PG_Error,
			StandartCode: cErrors.StatusBadRequest,
			Message:      BadRequest,
		}
	}
	return true, nil
}

func (p postgresRepository) GetUserIDByNickName(params *users.GetUserIDByNickNameParams) (int64, *cErrors.ResponseErrorModel) {
	var data []int64
	err := p.db.Select(&data, "SELECT id FROM user_db.public.users WHERE nickname=$1", params.Nickname)
	//fmt.Println(data)
	if err != nil {
		fmt.Println(err)
		return 0,
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByEmail_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			}
	}

	if len(data) == 0 {
		return 0,
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByEmail_Repo_No_Such_User,
				StandartCode: cErrors.StatusBadRequest,
				Message:      BadRequest,
			}
	}
	return data[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) CreateClientUuid(params *users.CreateClientUidParamsRepo) (string, *cErrors.ResponseErrorModel) {
	var lastInsertUuid string
	err := p.db.QueryRow(`INSERT INTO auth_db.public.client_user 
    							(user_id) 
								VALUES ($1) RETURNING client_uuid`,
		params.ClientID).Scan(&lastInsertUuid)
	if err != nil {
		return "", &cErrors.ResponseErrorModel{
			InternalCode: cErrors.CreateClient_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if lastInsertUuid == "" {
		fmt.Println(err)
		return lastInsertUuid, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.CreateClient_RW_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	return lastInsertUuid, &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) GetUserIdByUuid(uuid string) (int64, *cErrors.ResponseErrorModel) {
	var data []int64
	err := p.db.Select(&data, "SELECT user_id FROM client_user WHERE client_uuid=$1", uuid)
	if err != nil { // public_client_user это связывающая таблица
		log.Println("zdez?", err)
		return 0,
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByID_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			}
	}
	if len(data) == 0 {
		return 0,
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByID_Repo_No_Such_User,
				StandartCode: cErrors.StatusBadRequest,
				Message:      "not found uuid",
			}
	}
	return data[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) ClientUUID(params *users.ClientUuidByIDParams) (string, *cErrors.ResponseErrorModel) {
	var data []string
	err := p.db.Select(&data, "SELECT * FROM auth_db.public.client_user WHERE user_id=$1", params.UserId)
	if err != nil { // public_client_user это связывающая таблица
		return "",
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByID_Repo_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			}
	}
	if len(data) == 0 {
		return "",
			&cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByID_Repo_No_Such_User,
				StandartCode: cErrors.StatusBadRequest,
				Message:      "not found uuid",
			}
	}
	return data[0], &cErrors.ResponseErrorModel{}
}

func (p postgresRepository) AddClientUser(params *users.AddClientUserParamsRepo) (string, *cErrors.ResponseErrorModel) {
	var lastInsertUuid string
	err := p.db.QueryRow(`INSERT INTO public.client_user 
    							(user_id, client_uuid) 
								VALUES ($1, $2) RETURNING client_uuid`,
		params.UserId, params.ClientUuid).Scan(&lastInsertUuid)
	if err != nil {
		return "", &cErrors.ResponseErrorModel{
			InternalCode: cErrors.AddClientUser_RW_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	if lastInsertUuid == "" {
		return lastInsertUuid, &cErrors.ResponseErrorModel{
			InternalCode: cErrors.AddClientUser_RW_Repo_PG_Error,
			StandartCode: cErrors.StatusInternalServerError,
			Message:      Internal,
		}
	}
	return lastInsertUuid, &cErrors.ResponseErrorModel{}
}
