package usecase

import (
	"UsersService/config"
	"UsersService/internal/cErrors"
	"UsersService/internal/users"
	"UsersService/pkg/auth"
	"UsersService/pkg/logger"
	"UsersService/pkg/rating"
	"UsersService/pkg/secure"
	"UsersService/pkg/sso"
	"UsersService/pkg/utils"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt"
	"log"
	"strings"
	"sync"
)

const (
	Internal string = "Users: Internal Server Error"
)

type UsersUsecase struct {
	logger   *logger.ApiLogger
	repo     users.Repository
	shield   *secure.Shield
	authSR   auth.ServiceAuth
	ssoSR    sso.ServiceSso
	ratingSR rating.ServiceRating
	kyc      config.Kyc
}

func NewUsersUsecase(logger *logger.ApiLogger, repo users.Repository, shield *secure.Shield, authSR auth.ServiceAuth,
	ssoSR sso.ServiceSso, ratingSR rating.ServiceRating, kyc config.Kyc) users.UseCase {
	return &UsersUsecase{
		logger:   logger,
		repo:     repo,
		shield:   shield,
		authSR:   authSR,
		ssoSR:    ssoSR,
		ratingSR: ratingSR,
		kyc:      kyc,
	}
}

func (u *UsersUsecase) AddAuthKyc(params *users.AuthKycParams) {
	u.repo.AddAuthKyc(params)
}

func (u *UsersUsecase) GetAuthKyc(params *users.ClientID) *users.ResponseGetAuthKyc {
	auth, cErr := u.repo.GetAuthKyc(params)
	return &users.ResponseGetAuthKyc{
		Error: cErr,
		Data: &users.ResponseGetAuthKycModel{
			Auth: auth,
		},
	}
}

func (u *UsersUsecase) GetUserNicknameByID(params *users.GetUserNicknameByIDParams) *users.ResponseGetUserNicknameByID {
	log.Println("в юзкейсе пытаюсь взять никнейм юзера")
	nickname, cErr := u.repo.GetUserNicknameByID(params)
	return &users.ResponseGetUserNicknameByID{
		Data:  &users.ResponseGetUserNicknameByIDModel{Nickname: nickname},
		Error: cErr,
	}
}

func (u *UsersUsecase) GetRegistration(params *users.GetRegistrationParams) *users.ResponseGetRegistration {
	count, cErr := u.repo.GetRegistration(params)
	return &users.ResponseGetRegistration{
		Error: cErr,
		Data:  &users.ResponseGetRegistrationModel{Count: count},
	}
}

func (u *UsersUsecase) GetActiveNotice(params *users.GetActiveNoticeParams) *users.ResponseGetActiveNotice {
	notices, cErr := u.repo.GetActiveNotice(params)
	return &users.ResponseGetActiveNotice{
		Error: cErr,
		Data:  notices,
	}
}

func (u *UsersUsecase) ReadNotice(params *users.ReadNoticeParams) *users.ResponseReadNotice {
	cError := u.repo.ReadNotice(params)
	var suc bool
	if cError.InternalCode == 0 {
		suc = true
	}
	return &users.ResponseReadNotice{
		Error: cError,
		Data:  &users.ResponseSuccessModel{Success: suc},
	}
}

func (u *UsersUsecase) AddNotice(params *users.ADDNotice) *users.ResponseAddNotice {
	cErr := u.repo.AddNotice(params)
	var suc bool
	if cErr.InternalCode == 0 {
		suc = true
	}
	return &users.ResponseAddNotice{
		Error: cErr,
		Data:  &users.ResponseSuccessModel{Success: suc},
	}
}

func (u *UsersUsecase) UpdateUserLastActivity(params *users.UpdateUserLastActivityParams) *users.ResponseUpdateUserLastActivity {
	ok, cErr := u.repo.UpdateUserLastActivity(params)
	return &users.ResponseUpdateUserLastActivity{
		Error: cErr,
		Data:  &users.ResponseSuccessModel{Success: ok},
	}
}

func (u *UsersUsecase) Recovery(params *users.RecoveryParams) *users.ResponseRecovery {
	user, err := u.repo.GetUserByEmail((*users.GetUserByEmailParams)(params))
	if err.InternalCode != 0 {
		return &users.ResponseRecovery{
			Data:  &users.ResponseSuccessModel{Success: false},
			Error: err,
		}
	}
	clientUuid, err := u.repo.ClientUUID(&users.ClientUuidByIDParams{
		UserId: user.ClientID,
	})
	if err != nil {
		return &users.ResponseRecovery{
			Data: &users.ResponseSuccessModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Message,
			},
		}
	}
	ssoRes := u.ssoSR.RecoveryInit(&sso.RecoveryInitParams{ClientUuid: clientUuid, LanguageIso: user.Language})
	return &users.ResponseRecovery{
		Data:  &users.ResponseSuccessModel{Success: ssoRes.Data.Success},
		Error: ssoRes.Error,
	}
}

func (u *UsersUsecase) RecoveryConfirm(params *users.RecoveryConfirmParams) *users.ResponseRecoveryConfirm {
	ssoRes := u.ssoSR.RecoveryConfirm(&sso.RecoveryConfirmParams{
		Password: params.Password,
		Hash:     params.Hash,
	})
	return &users.ResponseRecoveryConfirm{
		Data:  &users.ResponseSuccessModel{Success: ssoRes.Data.Success},
		Error: ssoRes.Error,
	}
}

func (u *UsersUsecase) UpdateLastLogin(params *users.UpdateLastLoginParams) {
	_, err := u.repo.UpdateLastLogin(params)
	if err.InternalCode != 0 {
		u.logger.Errorf("Error update last entry: %s", err.Message)
	}
}

func (u *UsersUsecase) GetUserByEmail(params *users.GetUserByEmailParams) *users.ResponseGetClientByEmail {
	user, err := u.repo.GetUserByEmail(params)
	return &users.ResponseGetClientByEmail{
		Data:  user,
		Error: err,
	}
}

func (u *UsersUsecase) GetUserByNicknameWithID(params *users.GetUserByNicknameWithID) *users.ResponseGetUserByNicknameWithID {
	user, err := u.repo.GetUserByNickNameWithID(params)
	return &users.ResponseGetUserByNicknameWithID{
		Data:  user,
		Error: err,
	}
}

func (u *UsersUsecase) GetUserReviewsByNickname(params *users.GetUserByNickNameParams) *users.ResponseGetAllReviews {
	var CommentRes users.ResponseGetAllReviewsModels
	user := u.GetUserIDByNickName(&users.GetUserIDByNickNameParams{Nickname: params.Nickname})
	if user.Error.InternalCode != 0 {
		return &users.ResponseGetAllReviews{
			Error: user.Error,
			Data:  &CommentRes,
		}
	}
	userAllReviews := u.ratingSR.GetAllReviews(&rating.GetAllReviewsParams{
		ClientID: user.Data.ClientID, Limit: params.Limit, Page: params.Page, Type: params.Type})
	if userAllReviews.Error.InternalCode != 0 {
		return &users.ResponseGetAllReviews{
			Data:  &CommentRes,
			Error: userAllReviews.Error,
		}
	}
	if userAllReviews.Data.Comments != nil {
		CommentResNickname, err := u.replaceIdByNicknameForComments(userAllReviews.Data.Comments)
		if err.InternalCode != 0 {
			return &users.ResponseGetAllReviews{
				Data:  &users.ResponseGetAllReviewsModels{},
				Error: err,
			}
		}
		return &users.ResponseGetAllReviews{
			Error: user.Error,
			Data: &users.ResponseGetAllReviewsModels{
				Total:    userAllReviews.Data.Total,
				Pages:    userAllReviews.Data.Pages,
				Comments: &CommentResNickname,
			},
		}
	}
	fmt.Println(&CommentRes)
	return &users.ResponseGetAllReviews{
		Error: user.Error,
		Data:  &CommentRes,
	}
}

func (u *UsersUsecase) ConfirmKyc(params *users.HashKycConfirm) *users.ResponseConfirmKyc {
	var responseModel users.ResponseConfirmKycModel
	http := resty.New().EnableTrace().SetDebug(true)
	var body users.HashKycConfirmBody
	body.CallBack = "https://apiv1.exnode.ru/proxy/5/kyc/confirm"
	body.ClientID = params.Hash
	response, err := http.R().SetBasicAuth(u.kyc.Username, u.kyc.Pass).SetBody(body).SetResult(&responseModel).Post(u.kyc.Url)
	if err != nil {
		return &users.ResponseConfirmKyc{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &users.ResponseConfirmKyc{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode > 399 {
		return &users.ResponseConfirmKyc{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &users.ResponseConfirmKyc{
		Error: &cErrors.ResponseErrorModel{},
		Data:  &responseModel,
	}
}

func (u *UsersUsecase) GetUserIDByAccessToken(params *users.GetUserIDByAccessTokenParams) *users.ResponseGetUserIDByAccessToken {
	//email := u.decodeToken(params.Access)
	email := u.decodeToken(params.Access)
	if email == "" {
		return &users.ResponseGetUserIDByAccessToken{
			Data: nil,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByAccessToken_FailedDecode,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	user, err := u.repo.GetUserByEmail(&users.GetUserByEmailParams{Email: email})
	if err.InternalCode != 0 {
		return &users.ResponseGetUserIDByAccessToken{
			Data: nil,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByAccessToken_FailedDecode,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	return &users.ResponseGetUserIDByAccessToken{
		Data:  &users.ClientID{ClientID: user.ClientID},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) GetClientUUIDByAccessToken(params *users.GetUserIDByAccessTokenParams) *users.ResponseGetUserUUIDByAccessToken {
	//email := u.decodeToken(params.Access)
	//if email == "" {
	//	return &users.ResponseGetUserUUIDByAccessToken{
	//		Data: nil,
	//		Error: &cErrors.ResponseErrorModel{
	//			InternalCode: cErrors.Clients_GetUserByAccessToken_FailedDecode,
	//			StandartCode: cErrors.StatusInternalServerError,
	//			Message:      Internal,
	//		},
	//	}
	//}
	uuid := u.decodeUuidToken(params.Access)
	if uuid == "" {
		return &users.ResponseGetUserUUIDByAccessToken{
			Data: nil,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByAccessToken_FailedDecode,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}

	log.Println("Check UUID from token", uuid)
	//user, err := u.repo.GetUserByEmail(&users.GetUserByEmailParams{Email: email})
	//if err.InternalCode != 0 {
	//	return &users.ResponseGetUserUUIDByAccessToken{
	//		Data: nil,
	//		Error: &cErrors.ResponseErrorModel{
	//			InternalCode: cErrors.Clients_GetUserByAccessToken_FailedDecode,
	//			StandartCode: cErrors.StatusInternalServerError,
	//			Message:      Internal,
	//		},
	//	}
	//}
	//clientUuid, err := u.repo.ClientUUID(&users.ClientUuidByIDParams{
	//	UserId: user.ClientID,
	//})
	//return &users.ResponseGetUserUUIDByAccessToken{
	//	Data:  &users.ClientUUID{ClientUUID: clientUuid},
	//	Error: &cErrors.ResponseErrorModel{},
	//}
	return &users.ResponseGetUserUUIDByAccessToken{
		Data:  &users.ClientUUID{ClientUUID: uuid},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) GetClientIP(params *users.GetUserIpParams) *users.GetClientIpListResponse {
	var ipList []string
	ipList, err := u.repo.GetClientIP(params)
	if err.InternalCode != 0 {
		return &users.GetClientIpListResponse{
			Data: &users.ResponseGetClientIpModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetClientIP_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	if len(ipList) == 0 {
		return &users.GetClientIpListResponse{
			Data: &users.ResponseGetClientIpModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetClientIP_NoSuchClient,
				StandartCode: cErrors.StatusBadRequest,
				Message:      "no client with such id or ip list is empty",
			},
		}
	}
	return &users.GetClientIpListResponse{
		Data:  &users.ResponseGetClientIpModel{IpList: ipList},
		Error: err,
	}
}

func (u *UsersUsecase) GetUserByAccessToken(params *users.GetUserByAccessTokenParams) *users.ResponseGetUserByAccessToken {
	var (
		//wg          sync.WaitGroup
		cErr *cErrors.ResponseErrorModel
		//userAuth    *sso.ResponseGetClientVerification
		//userPrivate *sso.ResponseGetClientPrivate
		//userRate    *rating.ResponseGetClientSRRating
		//userRate2   *rating.ResponseGetClientStatistics
	)
	//wg.Add(4)
	//nickname := u.decodeToken(params.Token) // переписываю на mail
	//if nickname == "" {
	//	return &users.ResponseGetUserByAccessToken{
	//		Data: nil,
	//		Error: &cErrors.ResponseErrorModel{
	//			InternalCode: cErrors.Clients_GetUserByAccessToken_FailedDecode,
	//			StandartCode: cErrors.StatusInternalServerError,
	//			Message:      Internal,
	//		},
	//	}
	//}

	email := u.decodeToken(params.Token)
	if email == "" {
		return &users.ResponseGetUserByAccessToken{
			Data: nil,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByAccessToken_FailedDecode,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	//userUser, cErr := u.repo.GetUserByNickName(&users.GetUserByNickNameParams{Nickname: nickname})
	//if cErr.InternalCode != 0 {
	//	return &users.ResponseGetUserByAccessToken{
	//		Data:  nil,
	//		Error: cErr,
	//	}
	//}
	userUser, cErr := u.repo.GetUserByEmail(&users.GetUserByEmailParams{Email: email})
	if cErr.InternalCode != 0 {
		return &users.ResponseGetUserByAccessToken{
			Data:  nil,
			Error: cErr,
		}
	}

	//go func() {
	//	userAuth, _ = u.ssoSR.GetClientVerification(&users.ClientID{ClientID: userUser.ClientID})
	//	wg.Done()
	//}()
	//
	//go func() {
	//	userPrivate = u.ssoSR.GetClientPrivate(&users.ClientID{ClientID: userUser.ClientID})
	//	wg.Done()
	//}()

	//go func() {
	//	userRate = u.ratingSR.GetClientSRRating(&rating.GetClientSRRatingParams{ClientID: userUser.ClientID})
	//	wg.Done()
	//}()
	//
	//go func() {
	//	userRate2 = u.ratingSR.GetClientRatingForOrders(&rating.GetClientSRRatingParams{ClientID: userUser.ClientID})
	//	wg.Done()
	//}()

	//wg.Wait()
	//if userRate2.Error.InternalCode != 0 {
	//	return &users.ResponseGetUserByAccessToken{
	//		Data:  nil,
	//		Error: userRate2.Error,
	//	}
	//}
	//if userRate.Error.InternalCode != 0 {
	//	return &users.ResponseGetUserByAccessToken{
	//		Data:  nil,
	//		Error: userRate.Error,
	//	}
	//}
	//if userAuth.Error.InternalCode != 0 {
	//	return &users.ResponseGetUserByAccessToken{
	//		Data:  nil,
	//		Error: userAuth.Error,
	//	}
	//}

	return &users.ResponseGetUserByAccessToken{
		Data: &users.FullMeInfo{
			NickName:         userUser.Nickname,
			Avatar:           userUser.Avatar,
			Bio:              "",
			RegistrationDate: userUser.RegistrationDate,
			LastVisit:        userUser.LastEntry,
			LastActivity:     userUser.LastActivity,
			Email:            u.emailEncode(userUser.Email),
			//Tg:               userPrivate.Data.Tg,
			Tg:           "",
			BlockedUntil: userUser.BlockedUntil,
			Language:     userUser.Language,
			IsBlocked:    userUser.IsBlocked,
			Ip:           userUser.Ip,
			//UserStatistic:          userRate.Data,
			//UserStatisticForOrders: (*users.Statistic)(userRate2.Data),
			//Verification:           userAuth.Data,
		},
		Error: cErr,
	}
}

func (u *UsersUsecase) emailEncode(email string) string {
	s := strings.Split(email, "@")
	var newString string
	newString = s[0][0:1] + strings.Repeat("*", len(s[0]))
	return newString + "@" + s[1]
}

func (u *UsersUsecase) decodeToken(token string) string {
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	var email string
	email = fmt.Sprint(claims["email"])
	return email
	//claims := jwt.MapClaims{}
	//_, _ = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
	//	return nil, nil
	//})
	//
	//data := make(map[string]interface{})
	//for key, value := range claims {
	//	fmt.Println(value, key)
	//	data[key] = value
	//}
	//
	//email := fmt.Sprint(data["email"])
	//
	//return email
}

func (u *UsersUsecase) decodeUuidToken(token string) string {
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	var uuid string
	uuid = fmt.Sprint(claims["uuid"])
	return uuid
}

func (u *UsersUsecase) GetUserByNickName(params *users.GetUserByNickNameParams) *users.ResponseGetClientByNickName {
	var wg sync.WaitGroup
	wg.Add(3)
	userUser, err := u.repo.GetUserByNickName(params)
	if err.InternalCode != 0 {
		return &users.ResponseGetClientByNickName{
			Data:  nil,
			Error: err,
		}
	}
	var userAuth *sso.ResponseGetClientVerification
	go func() {
		userAuth, _ = u.ssoSR.GetClientVerification(&users.ClientID{ClientID: userUser.ClientID})
		wg.Done()
	}()

	var userRate *rating.ResponseGetClientSRRating
	go func() {
		userRate = u.ratingSR.GetClientSRRating(&rating.GetClientSRRatingParams{ClientID: userUser.ClientID})
		wg.Done()
	}()

	var userRate2 *rating.ResponseGetClientStatistics
	go func() {
		userRate2 = u.ratingSR.GetClientRatingForOrders(&rating.GetClientSRRatingParams{ClientID: userUser.ClientID})
		if userRate2.Error.InternalCode != 0 {
			userRate2.Data = &rating.Statistic{}
		}
		if userRate2.Data == nil {
			userRate2.Data = &rating.Statistic{}
		}

		wg.Done()
	}()
	wg.Wait()
	if userAuth.Error.InternalCode != 0 {
		return &users.ResponseGetClientByNickName{
			Data:  nil,
			Error: err,
		}
	}
	if userRate.Error.InternalCode != 0 {
		return &users.ResponseGetClientByNickName{
			Data:  nil,
			Error: err,
		}
	}
	return &users.ResponseGetClientByNickName{
		Data: &users.FullUserInfo{
			NickName:               userUser.Nickname,
			Avatar:                 userUser.Avatar,
			RegistrationDate:       userUser.RegistrationDate,
			LastVisit:              userUser.LastEntry,
			LastActivity:           userUser.LastActivity,
			BlockedUntil:           userUser.BlockedUntil,
			Language:               userUser.Language,
			IsBlocked:              userUser.IsBlocked,
			UserStatistic:          userRate.Data,
			UserStatisticForOrders: (*users.Statistic)(userRate2.Data),
			Verification:           userAuth.Data,
		},
		Error: err,
	}
}

func (u *UsersUsecase) GetUserByAccessTokenWithID(params *users.GetUserByAccessTokenParams) *users.ResponseGetUserByAccessTokenWithID {
	email := u.decodeToken(params.Token)
	if email == "" {
		return &users.ResponseGetUserByAccessTokenWithID{
			Data: nil,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_GetUserByAccessToken_FailedDecode,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	userUser, err := u.repo.GetUserByEmail(&users.GetUserByEmailParams{Email: email})
	if err.InternalCode != 0 {
		return &users.ResponseGetUserByAccessTokenWithID{
			Data:  nil,
			Error: err,
		}
	}
	userSso, _ := u.ssoSR.GetClientVerification(&users.ClientID{ClientID: userUser.ClientID})
	if userSso.Error.InternalCode != 0 {
		return &users.ResponseGetUserByAccessTokenWithID{
			Data:  nil,
			Error: err,
		}
	}
	fmt.Println("2")
	userRate := u.ratingSR.GetClientSRRating(&rating.GetClientSRRatingParams{ClientID: userUser.ClientID})
	if userRate.Error.InternalCode != 0 {
		return &users.ResponseGetUserByAccessTokenWithID{
			Data:  nil,
			Error: err,
		}
	}

	//userAllReviews := u.ratingSR.GetAllReviews(&rating.GetAllReviewsParams{ClientID: userUser.ID})
	//if userAllReviews.Error.InternalCode != 0 {
	//	return &users.ResponseGetUserByAccessTokenWithID{
	//		Data:  nil,
	//		Error: err,
	//	}
	//}
	//var CommentRes []users.CommentResponse
	//if userAllReviews.Data != nil {
	//	CommentRes, err = u.replaceIdByNicknameForComments(userAllReviews.Data.Comments)
	//	if err.InternalCode != 0 {
	//		return &users.ResponseGetUserByAccessTokenWithID{
	//			Data:  nil,
	//			Error: err,
	//		}
	//	}
	//}
	return &users.ResponseGetUserByAccessTokenWithID{
		Data: &users.FullUserInfoWithID{
			ClientID:         userUser.ClientID,
			NickName:         userUser.Nickname,
			Email:            userUser.Email,
			Avatar:           userUser.Avatar,
			RegistrationDate: userUser.RegistrationDate,
			LastVisit:        userUser.LastEntry,
			BlockedUntil:     userUser.BlockedUntil,
			Language:         userUser.Language,
			IsBlocked:        userUser.IsBlocked,
			UserStatistic:    userRate.Data,
			Verification:     userSso.Data,
			//AllReview:        &CommentRes,
		},
		Error: err,
	}
}

func (u *UsersUsecase) GetUserByID(params *users.GetUserByIDParams) *users.ResponseGetUserByID {
	log.Println("сюда захожу вообще?")
	user, err := u.repo.GetUserByID(params)
	if err.InternalCode != 0 {
		return &users.ResponseGetUserByID{
			Data:  nil,
			Error: err,
		}
	}
	//var wg sync.WaitGroup
	//var rate1 *rating.ResponseGetClientSRRating
	//var rate2 *rating.ResponseGetClientStatistics
	//var _auth *sso.ResponseGetClientVerification
	//wg.Add(3)
	//go func() {
	//	defer wg.Done()
	//	rate1 = u.ratingSR.GetClientSRRating(&rating.GetClientSRRatingParams{ClientID: params.ClientID})
	//}()
	//go func() {
	//	defer wg.Done()
	//	rate2 = u.ratingSR.GetClientRatingForOrders(&rating.GetClientSRRatingParams{ClientID: params.ClientID})
	//
	//}()
	//go func() {
	//	defer wg.Done()
	//	_auth, _ = u.ssoSR.GetClientVerification(&users.ClientID{ClientID: params.ClientID})
	//}()
	//wg.Wait()
	//if rate1.Error.InternalCode != 0 {
	//	rate1 = &rating.ResponseGetClientSRRating{}
	//}
	//if rate2.Error.InternalCode != 0 {
	//	rate2 = &rating.ResponseGetClientStatistics{}
	//}
	//if _auth.Error.InternalCode != 0 {
	//	rate2 = &rating.ResponseGetClientStatistics{}
	//}
	//user.KYC = _auth.Data.KycConfirm
	log.Println("check err: ", err)
	return &users.ResponseGetUserByID{
		Data: &users.ResponseGetUserByIDModel{
			User: user,
			Stat: &users.Stat{
				//FeedbackPositive:       rate1.Data.FeedbacksPositive,
				//FeedbackNegative:       rate1.Data.FeedbacksNegative,
				//TradesCompletedPercent: rate2.Data.PercentDoneOrders,
				//TradesCompleted:        float64(rate2.Data.Orders),
				FeedbackPositive:       1,
				FeedbackNegative:       1,
				TradesCompletedPercent: 1,
				TradesCompleted:        float64(1),
			},
		},
		Error: err,
	}
}

func (u *UsersUsecase) IsUserBlocked(params *users.IsUserBlockedParams) *users.ResponseIsUserBlocked {
	ok, err := u.repo.IsUserBlocked(params)
	if err.InternalCode != 0 {
		return &users.ResponseIsUserBlocked{
			Data: &users.ResponseIsUserBlockedModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_IsUserBlocked_PG_Error,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      Internal,
			},
		}
	}
	return &users.ResponseIsUserBlocked{
		Data:  &users.ResponseIsUserBlockedModel{IsBlocked: ok},
		Error: err,
	}
}

func (u *UsersUsecase) ClientSignUp(params *users.SignUpParams) *users.ResponseClientSignUp {
	log.Println("зашел в clientSignUp")
	email := strings.ToLower(params.Email)
	_, err := u.repo.GetUserByEmail(&users.GetUserByEmailParams{Email: email})
	if err.InternalCode == 0 {
		return &users.ResponseClientSignUp{
			Data: nil,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.Clients_ClientSignUp_AuthError,
				StandartCode: cErrors.StatusBadRequest,
				Message:      "Client already exist",
			},
		}
	}
	cErr := utils.ValidPassword(params.Password)
	if cErr.InternalCode != 0 {
		return &users.ResponseClientSignUp{
			Data:  nil,
			Error: &cErr,
		}
	}

	//userUuid, err := u.repo.CreateClientUuid(&users.CreateClientUidParamsRepo{
	//	ClientID: userID.Data.ClientID,
	//})

	//response, errSR := u.ssoSR.ClientSignUp(&users.SignUpSAParams{
	//	Phone:    params.Phone,
	//	IsDnd:    false,
	//	Email:    email,
	//	Password: params.Password,
	//	UA:       params.UA,
	//})

	response, errSR := u.ssoSR.ClientSignUp(&users.SignUpParams{
		Phone:    params.Phone,
		IP:       "testIp",
		Email:    email,
		Password: params.Password,
		UA:       params.UA,
	})
	log.Println("читаю ответ из ссо", response)
	if errSR != nil {
		return &users.ResponseClientSignUp{
			Data: &users.ResponseClientSignUpModel{},
			Error: &cErrors.ResponseErrorModel{
				Message: errSR.Error(),
			},
		}
	}
	if params.Language == "" {
		params.Language = "ru"
	}

	log.Println("check access token:", response.Data.AccessToken)

	uuid := u.GetClientUUIDByAccessToken(&users.GetUserIDByAccessTokenParams{
		Access: response.Data.AccessToken,
	})
	log.Println("перед записью в бд?")
	userID := u.CreateClient(&users.CreateClientParams{IsDnD: false, Language: params.Language, Email: email})
	if userID.Error.InternalCode != 0 {
		return &users.ResponseClientSignUp{
			Error: userID.Error,
			Data:  nil,
		}
	}
	log.Println("перед записью в связующую таблицу?")
	_, err = u.repo.AddClientUser(&users.AddClientUserParamsRepo{
		UserId:     userID.Data.ClientID,
		ClientUuid: uuid.Data.ClientUUID,
	})
	if err.InternalCode != 0 {
		return &users.ResponseClientSignUp{
			Error: err,
			Data:  nil,
		}
	}

	//go u.ssoSR.ConfirmMailReq(&admins.ConfirmMailReqParams{
	//	ClientUUID:  uuid.Data.ClientUUID,
	//	LanguageIso: params.Language,
	//})
	//if errSR != nil {
	//	return &users.ResponseClientSignUp{
	//		Data:  &users.ResponseClientSignUpModel{},
	//		Error: &cErrors.ResponseErrorModel{},
	//	}
	//}
	fmt.Println("3", response)
	return &users.ResponseClientSignUp{
		Data: &users.ResponseClientSignUpModel{
			RefreshToken: response.Data.RefreshToken,
			AccessToken:  response.Data.AccessToken,
		},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) ClientChangePassword(params *users.ChangeClientPasswordParams) *users.ResponseClientChangePassword {
	if params.NewPasswordAgain != params.NewPassword {
		return &users.ResponseClientChangePassword{
			Data: &users.ResponseSuccessModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.PasswordsNotEqual,
				StandartCode: cErrors.StatusBadRequest,
				Message:      "passwords must be equal",
			},
		}
	}
	cErr := utils.ValidPassword(params.NewPassword)
	if cErr.InternalCode != 0 {
		return &users.ResponseClientChangePassword{
			Data:  &users.ResponseSuccessModel{},
			Error: &cErr,
		}
	}

	response, err := u.ssoSR.ClientChangePassword(params)
	if err != nil {
		return &users.ResponseClientChangePassword{
			Data:  nil,
			Error: response.Error,
		}
	}
	return &users.ResponseClientChangePassword{
		Data:  &users.ResponseSuccessModel{Success: response.Data.Success},
		Error: response.Error,
	}
}

func (u *UsersUsecase) ClientSignIn(params *users.SignInSAParams) *users.ResponseClientSignIn {
	//userId := u.GetUserByEmail(&users.GetUserByEmailParams{
	//	Email: params.Email,
	//})
	//uI, err := strconv.Atoi(userId.Data.Email) // строки не было
	//clientUuid, err := u.repo.ClientUUID(&users.ClientUuidByIDParams{
	//	UserId: int64(uI), // userId был 07.06 17:21
	//})
	//if err != nil {
	//	return &users.ResponseClientSignIn{
	//		Data: &users.ResponseClientSignInModel{},
	//		Error: &cErrors.ResponseErrorModel{
	//			InternalCode: cErrors.StatusInternalServerError,
	//			StandartCode: cErrors.StatusInternalServerError,
	//			Message:      "error HERE",
	//		},
	//	}
	//}

	response, cErr := u.ssoSR.ClientSignIn(&sso.SignInSAParams{
		ClientUUID: params.Email,
		Password:   params.Password,
		Ip:         params.Ip,
		UA:         params.UA,
	})
	if cErr != nil {
		return &users.ResponseClientSignIn{
			Data: &users.ResponseClientSignInModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      cErr.Error(),
			},
		}
	}
	return &users.ResponseClientSignIn{
		Data: &users.ResponseClientSignInModel{
			RefreshToken: response.Data.RefreshToken,
			AccessToken:  response.Data.AccessToken,
		},
		Error: response.Error,
	}
}

func (u *UsersUsecase) ClientSignInTG(params *users.ClientSignInTGParams) *users.ResponseClientSignIn {
	response, err := u.authSR.ClientSignInTg(params)
	if err != nil {
		return &users.ResponseClientSignIn{
			Data: &users.ResponseClientSignInModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	return &users.ResponseClientSignIn{
		Data:  &users.ResponseClientSignInModel{RefreshToken: response.Data.RefreshToken, AccessToken: response.Data.AccessToken},
		Error: response.Error,
	}
}

func (u *UsersUsecase) ClientChangeDefaultLanguage(params *users.ChangeDefaultUserLanguageParams) *users.ResponseClientChangeDefaultLanguage {
	lang, err := u.repo.GetLanguage(&users.IsValidLanguageParams{LanguageIso: params.Ticker})
	if err.InternalCode != 0 {
		return &users.ResponseClientChangeDefaultLanguage{
			Data:  &users.ResponseSuccessModel{Success: false},
			Error: err,
		}
	}
	if !lang.Available {
		return &users.ResponseClientChangeDefaultLanguage{
			Data: &users.ResponseSuccessModel{
				Success: false,
				Message: "lang not available",
			},
			Error: err,
		}
	}
	ok, err := u.repo.ChangeDefaultUserLanguage(params)
	return &users.ResponseClientChangeDefaultLanguage{
		Data:  &users.ResponseSuccessModel{Success: ok},
		Error: err,
	}
}

func (u *UsersUsecase) UpdateUserLastEntry(params *users.UpdateUserLastEntryParams) *users.ResponseUpdateUserLastEntry {
	ok, err := u.repo.UpdateUserLastEntry(params)
	if err != nil {
		return &users.ResponseUpdateUserLastEntry{
			Data:  &users.ResponseSuccessModel{Success: false},
			Error: err,
		}
	}

	return &users.ResponseUpdateUserLastEntry{
		Data:  &users.ResponseSuccessModel{Success: ok},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) GetClientRating(params *users.GetClientRatingParams) *users.ResponseGetClientRating {
	response := u.ratingSR.GetClientSRRating(&rating.GetClientSRRatingParams{ClientID: params.ClientID})
	if response.Error != nil {
		return &users.ResponseGetClientRating{
			Data:  nil,
			Error: response.Error,
		}
	}
	return &users.ResponseGetClientRating{
		Data:  response.Data,
		Error: response.Error,
	}
}

func (u *UsersUsecase) Validate(params *users.ValidateParams) *users.ResponseValidate {
	return &users.ResponseValidate{
		Data:  &users.ResponseSuccessModel{Success: true},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) CreateClient(params *users.CreateClientParams) *users.ResponseCreateClient {
	nickname := utils.GenerateNickName()
	//for {
	//	_, err := u.repo.GetUserByNickName(&users.GetUserByNickNameParams{Nickname: nickname})
	//	if err.InternalCode == 4040 {
	//		break
	//	} else {
	//		nickname = utils.GenerateNickName()
	//	}
	//	log.Println("застрял в цикле?")
	//}
	log.Println("Прошел цикл")
	avatar := utils.GenerateAvatar(nickname)
	time := utils.GetEuropeTime()
	virginUser := &users.CreateClientParamsRepo{
		NickName: nickname,
		Email:    params.Email,
		Avatar:   avatar,
		TimeNow:  time,
		IsDnD:    false,
		Language: params.Language,
	}

	userId, err := u.repo.CreateClient(virginUser)
	if err.InternalCode != 0 {
		return &users.ResponseCreateClient{
			Data:  nil,
			Error: err,
		}
	}
	log.Println("добавил юзера в бд?")
	return &users.ResponseCreateClient{
		Data:  &users.ResponseCreateClientModel{ClientID: userId, NickName: nickname},
		Error: err,
	}

}

func (u *UsersUsecase) CreateClientDND(params *users.CreateClientParams) *users.ResponseCreateClient {
	nickname := utils.GenerateNickName()
	for {
		_, err := u.repo.GetUserByNickName(&users.GetUserByNickNameParams{Nickname: nickname})
		if err.InternalCode == 4040 {
			break
		} else {
			nickname = utils.GenerateNickName()
		}

	}
	//avatar := utils.GenerateAvatar(nickname)
	//time := utils.GetEuropeTime()
	//virginUser := &users.CreateClientParamsRepo{
	//	NickName: nickname,
	//	Email:    nickname + "@gmail.com",
	//	Avatar:   avatar,
	//	TimeNow:  time,
	//	IsDnD:    true,
	//	Language: "ru",
	//}

	//userId, err := u.repo.CreateClient(virginUser)
	//if err.InternalCode != 0 {
	//	return &users.ResponseCreateClient{
	//		Data:  nil,
	//		Error: err,
	//	}
	//}

	//userUuid, err := u.repo.CreateClientUuid(&users.CreateClientUidParamsRepo{
	//	ClientID: userId,
	//})

	virginUser1 := &users.SignUpSAParams{
		Email: nickname + "@gmail.com",
		IsDnd: true,
	}
	tokens, ssoErr := u.ssoSR.CreateDndClient(virginUser1)
	if ssoErr != nil {
		return &users.ResponseCreateClient{
			Data: nil,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      ssoErr.Error(),
			},
		}
	}

	userId := u.GetUserIDByAccessToken(&users.GetUserIDByAccessTokenParams{
		Access: tokens.Data.AccessToken,
	})

	return &users.ResponseCreateClient{
		Data: &users.ResponseCreateClientModel{ClientID: userId.Data.ClientID, NickName: nickname},
	}
}

func (u *UsersUsecase) UpdateUserAvatar(params *users.UpdateUserAvatarParams) *users.ResponseUpdateUserAvatar {
	ok, err := u.repo.UpdateUserAvatar(params)
	if err != nil {
		return &users.ResponseUpdateUserAvatar{
			Data:  &users.ResponseSuccessModel{Success: false},
			Error: err,
		}
	}

	return &users.ResponseUpdateUserAvatar{
		Data:  &users.ResponseSuccessModel{Success: ok},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) UpdateUserNickName(params *users.UpdateUserNickNameParams) *users.ResponseUpdateUserNickName {
	_, cErr := u.repo.GetUserByNickName(&users.GetUserByNickNameParams{Nickname: params.NewNickName})
	if cErr.StandartCode != 400 {
		return &users.ResponseUpdateUserNickName{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusBadRequest,
				StandartCode: cErrors.StatusBadRequest,
			},
		}
	}
	changes, _ := u.repo.GetNicknameChanges(params)
	if len(*changes) != 0 {
		return &users.ResponseUpdateUserNickName{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusBadRequest,
				StandartCode: cErrors.StatusBadRequest,
				Message:      "nickname change limit",
			},
		}
	}
	user, _ := u.repo.GetUserByID(&users.GetUserByIDParams{ClientID: params.ClientID})

	_, err := u.repo.UpdateUserNickName(params)
	if err != nil {
		return &users.ResponseUpdateUserNickName{
			Error: err,
		}
	}
	res := u.authSR.ChangeNickname(params)
	if res.Error.InternalCode != 0 {
		return &users.ResponseUpdateUserNickName{
			Error: err,
		}
	}

	cErr = u.repo.UpdateNicknameChanges(&users.ChangeNickname{ClientID: params.ClientID, OldNickname: user.Nickname})
	if res.Error.InternalCode != 0 {
		return &users.ResponseUpdateUserNickName{
			Error: cErr,
		}
	}

	return &users.ResponseUpdateUserNickName{
		Data:  &users.ResponseChangeNicknameModel{Access: res.Data.Access},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) UpdateUserBio(params *users.UpdateUserBioParams) *users.ResponseUpdateUserBio {
	ok, err := u.repo.UpdateUserBio(params)
	if err != nil {
		return &users.ResponseUpdateUserBio{
			Data:  &users.ResponseSuccessModel{Success: false},
			Error: err,
		}
	}

	return &users.ResponseUpdateUserBio{
		Data:  &users.ResponseSuccessModel{Success: ok},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) replaceIdByNicknameForComments(params *[]users.Comment) ([]users.CommentResponse, *cErrors.ResponseErrorModel) {
	var resModel []users.CommentResponse
	for _, volume := range *params {
		nickname, err := u.repo.GetUserByID(&users.GetUserByIDParams{ClientID: volume.ClientReviewerId})
		if err.InternalCode != 0 {
			nick := utils.GenerateNickName()
			nickname = &users.User{
				Nickname: nick,
				Avatar:   utils.GenerateAvatar(nick),
			}
		}
		resModel = append(resModel, users.CommentResponse{
			ReviewerNickname: nickname.Nickname,
			ReviewerAvatar:   nickname.Avatar,
			Rate:             volume.Rate,
			InternalID:       volume.InternalID,
			CreatedDate:      volume.CreatedDate,
			Text:             volume.Text,
		})
	}

	return resModel, &cErrors.ResponseErrorModel{}
}

func (u *UsersUsecase) GetUserIDByNickName(params *users.GetUserIDByNickNameParams) *users.ResponseGetUserIDByNickName {
	res, err := u.repo.GetUserIDByNickName(params)
	if err.InternalCode != 0 {
		return &users.ResponseGetUserIDByNickName{
			Data:  nil,
			Error: &cErrors.ResponseErrorModel{},
		}
	}
	return &users.ResponseGetUserIDByNickName{
		Data:  &users.ResponseGetUserIDByNickNameModel{ClientID: res},
		Error: &cErrors.ResponseErrorModel{},
	}
}

func (u *UsersUsecase) ConfirmEmailByHash(params *users.ConfirmEmailByHashParams) *users.ResponseConfirmEmailByHash {
	data := u.authSR.ConfirmEmailByHash(params)
	return &users.ResponseConfirmEmailByHash{
		Data:  (*users.ResponseSuccessModel)(data.Data),
		Error: data.Error,
	}
}

func (u *UsersUsecase) ConfirmEmail(params *users.ConfirmEmailParams) *users.ResponseConfirmEmail {
	var responseModel users.ResponseConfirmEmail
	user := u.GetUserByAccessTokenWithID(&users.GetUserByAccessTokenParams{Token: params.Access})
	if user.Error.InternalCode != 0 {
		return &users.ResponseConfirmEmail{
			Error: user.Error,
			Data:  responseModel.Data,
		}
	}

	userUuid, err := u.repo.ClientUUID(&users.ClientUuidByIDParams{
		UserId: user.Data.ClientID,
	})
	if err != nil {
		return &users.ResponseConfirmEmail{
			Data: &users.ResponseSuccessModel{},
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Message,
			},
		}
	}

	data := u.ssoSR.ConfirmEmail(&users.ConfirmEmailAuthParams{
		Email:       user.Data.Email,
		LanguageIso: user.Data.Language,
		ClientUuid:  userUuid,
	})
	return &users.ResponseConfirmEmail{
		Data:  data.Data,
		Error: data.Error,
	}
}

func (u *UsersUsecase) GetClientVerification(params *users.GetClientRatingParams) {

}
