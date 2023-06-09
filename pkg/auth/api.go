package auth

import (
	"UsersService/config"
	"UsersService/internal/admins"
	"UsersService/internal/cConstants"
	"UsersService/internal/cErrors"
	"UsersService/internal/users"
	"UsersService/pkg/secure"
	"UsersService/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type Client struct {
	cfg        *config.Config
	httpClient *resty.Client
}

func GetClient(params *GetClientParams) ServiceAuth {
	return &Client{
		cfg: params.Config,
		httpClient: resty.New().OnBeforeRequest(SignatureMiddleware(params)).
			EnableTrace().SetDebug(true).SetBaseURL("http://localhost:11012"),
	}
}

func (c *Client) SendCodeToEmail(params *SendCodeToEmailParams) *ResponseSuccess {
	var responseModel ResponseSuccess
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.SendCodeToEmail)
	if err != nil {
		return &ResponseSuccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &ResponseSuccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseSuccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &responseModel
}

func (c *Client) CheckCodeFromEmail(params *VerCodeParams) *ResponseSuccess {
	var responseModel ResponseSuccess
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.CheckCodeForEmail)
	if err != nil {
		return &ResponseSuccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &ResponseSuccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseSuccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &responseModel
}

func SignatureMiddleware(connData *GetClientParams) resty.RequestMiddleware {
	return func(c *resty.Client, request *resty.Request) error {
		timestamp := utils.GetEuropeTime().Unix()
		jsonBody, _ := json.Marshal(request.Body)
		body := createRequestBody(strconv.FormatInt(timestamp, 10), string(jsonBody))
		request.SetHeader("Content-Type", "application/json")
		request.SetHeader("ApiPublic", *connData.Public)
		request.SetHeader("Signature", secure.CalcSignature(*connData.Private, body))
		request.SetHeader("TimeStamp", strconv.FormatInt(timestamp, 10))

		return nil
	}
}

func createRequestBody(timestamp, jsonBody string) string {
	return timestamp + jsonBody
}

func (c *Client) ConfirmKycInit(params *KycInitParams) *ResponseKycConfirmInit {
	var responseModel ResponseKycConfirmInit
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.KycConfInit)
	if err != nil {
		return &ResponseKycConfirmInit{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &ResponseKycConfirmInit{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseKycConfirmInit{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &responseModel
}

func (c *Client) ConfirmKyc(params *KycParams) *ResponseKycConfirm {
	var responseModel ResponseKycConfirm
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.KycConf)
	if err != nil {
		return &ResponseKycConfirm{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &ResponseKycConfirm{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseKycConfirm{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &responseModel
}

func (c *Client) ChangeNickname(params *users.UpdateUserNickNameParams) *users.ResponseUpdateUserNickName {
	var responseModel users.ResponseUpdateUserNickName
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ChangeNickname)
	if err != nil {
		return &users.ResponseUpdateUserNickName{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &users.ResponseUpdateUserNickName{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &users.ResponseUpdateUserNickName{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &responseModel
}

func (c *Client) ChangeTg(params *ChangeTgParams) *ResponseChangeTg {
	var responseModel ResponseChangeTg
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ClientChangeTg)
	if err != nil {
		return &ResponseChangeTg{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
			Data: &ResponseSuccessModel{Success: false},
		}
	}
	if response == nil {
		return &ResponseChangeTg{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
			Data: &ResponseSuccessModel{Success: false},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseChangeTg{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
			Data: &ResponseSuccessModel{Success: false},
		}
	}

	return &responseModel
}

func (c *Client) ClientChangePassword(params *users.ChangeClientPasswordParams) (*users.ResponseClientChangePassword, error) {
	var responseModel users.ResponseClientChangePassword
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ClientChangePassword)
	if err != nil {
		fmt.Println(err)
		return &users.ResponseClientChangePassword{}, err
	}
	if response == nil {
		return &users.ResponseClientChangePassword{}, errors.New("response is nil somehow")
	}
	//fmt.Println(response)
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &users.ResponseClientChangePassword{}, fmt.Errorf("status code {%d}", statusCode)
	}

	return &responseModel, nil
}

func (c *Client) CreateDndClient(params *users.SignUpSAParams) (*users.ResponseClientSignUp, error) {
	var responseModel users.ResponseClientSignUp
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.CreateDndClient)
	if err != nil {
		fmt.Println(err)
		return &users.ResponseClientSignUp{}, err
	}
	if response == nil {
		return &users.ResponseClientSignUp{}, errors.New("response is nil somehow")
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &users.ResponseClientSignUp{}, fmt.Errorf("status code {%d}", statusCode)
	}

	return &responseModel, nil
}

func (c *Client) GetClientPrivate(params *users.ClientID) *ResponseGetClientPrivate {
	var responseModel ResponseGetClientPrivate
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.GetClientPrivate)
	if err != nil {
		return &ResponseGetClientPrivate{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
			Data: &ResponseGetClientPrivateModel{},
		}
	}
	if response == nil {
		return &ResponseGetClientPrivate{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
			Data: &ResponseGetClientPrivateModel{},
		}
	}
	fmt.Println(response)
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseGetClientPrivate{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d statuc ccode", statusCode),
			},
			Data: &ResponseGetClientPrivateModel{},
		}
	}

	return &responseModel
}

func (c *Client) GetUserSessions(params *GetUserSessions) *ResponseGetUserSessions {
	var responseModel ResponseGetUserSessions
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.GetUserSessions)
	if err != nil {
		fmt.Println(err)
		return &ResponseGetUserSessions{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &ResponseGetUserSessions{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "nil response",
			},
		}
	}
	//fmt.Println(response)
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseGetUserSessions{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: statusCode,
				Message:      fmt.Sprintf("Response with status code %d", statusCode),
			},
		}
	}

	return &responseModel
}

func (c *Client) ClientSignUp(params *users.SignUpSAParams) (*users.ResponseClientSignUp, error) {
	var responseModel users.ResponseClientSignUp
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ClientSignUp)
	if err != nil {
		fmt.Println(err)
		return &users.ResponseClientSignUp{}, err
	}
	if response == nil {
		return &users.ResponseClientSignUp{}, errors.New("response is nil somehow")
	}
	//fmt.Println(response)
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &users.ResponseClientSignUp{}, fmt.Errorf("status code {%d}", statusCode)
	}

	return &responseModel, nil

}

func (c *Client) ClientSignIn(params *users.SignInSAParams) (*users.ResponseClientSignIn, error) {
	var responseModel users.ResponseClientSignIn
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ClientSignIn)
	if err != nil {
		fmt.Println(err)
		return &users.ResponseClientSignIn{}, err
	}
	if response == nil {
		return &users.ResponseClientSignIn{}, errors.New("response is nil somehow")
	}
	//fmt.Println(response)
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &users.ResponseClientSignIn{}, fmt.Errorf("status code {%d}", statusCode)
	}

	return &responseModel, nil
}

func (c *Client) ClientSignInTg(params *users.ClientSignInTGParams) (*users.ResponseClientSignIn, error) {
	var responseModel users.ResponseClientSignIn
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ClientSignInTg)
	if err != nil {
		fmt.Println(err)
		return &users.ResponseClientSignIn{}, err
	}
	if response == nil {
		return &users.ResponseClientSignIn{}, errors.New("response is nil somehow")
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &users.ResponseClientSignIn{}, fmt.Errorf("status code {%d}", statusCode)
	}

	return &responseModel, nil
}

func (c *Client) RefreshAccess(params *RefreshAccessParams) *ResponseRefreshAccess {
	var responseModel ResponseRefreshAccess
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.RefreshAccess)
	if err != nil {
		fmt.Println(err)
		return &ResponseRefreshAccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &ResponseRefreshAccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseRefreshAccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &responseModel
}

func (c *Client) Logout(params *LogoutParams) *ResponseLogout {
	var responseModel ResponseLogout
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.Logout)
	if err != nil {
		fmt.Println(err)
		return &ResponseLogout{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &ResponseLogout{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseLogout{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &responseModel
}

func (c *Client) ValidateAccess(params *ValidateAccessTokenParams) *ResponseValidateAccess {
	var responseModel ResponseValidateAccess
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ValidateAccess)
	if err != nil {
		fmt.Println(err)
		return &ResponseValidateAccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      err.Error(),
			},
		}
	}
	if response == nil {
		return &ResponseValidateAccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      "response is nil somehow",
			},
		}
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseValidateAccess{
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.StatusInternalServerError,
				StandartCode: cErrors.StatusInternalServerError,
				Message:      fmt.Sprintf("response with %d status code", statusCode),
			},
		}
	}

	return &responseModel
}

func (c *Client) ConfirmMailReq(params *admins.ConfirmMailReqParams) error {
	var responseModel users.ResponseClientChangePassword
	fmt.Println("EMAIL", params.ClientUUID, params.LanguageIso)
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ConfirmMailReq)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if response == nil {
		return errors.New("response is nil somehow")
	}
	fmt.Println(response)
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return fmt.Errorf("status code {%d}", statusCode)
	}

	return nil
}

func (c *Client) GetClient(params *users.ClientID) (*ResponseUserFromAuth, error) {
	var responseModel ResponseUserFromAuth
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ClientGetClient)
	if err != nil {
		fmt.Println(err)
		return &ResponseUserFromAuth{}, err
	}
	if response == nil {
		return &ResponseUserFromAuth{}, errors.New("response is nil somehow")
	}
	fmt.Println(response)
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseUserFromAuth{}, fmt.Errorf("status code {%d}", statusCode)
	}

	return &responseModel, nil
}

func (c *Client) GetClientVerification(params *users.ClientID) (*ResponseGetClientVerification, error) {
	var responseModel ResponseGetClientVerification
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.GetClientVerification)
	if err != nil {
		fmt.Println(err)
		return &ResponseGetClientVerification{}, err
	}
	if response == nil {
		return &ResponseGetClientVerification{}, errors.New("response is nil somehow")
	}
	fmt.Println(response)
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &ResponseGetClientVerification{}, fmt.Errorf("status code {%d}", statusCode)
	}

	return &responseModel, nil
}

func (c *Client) AdminSignIn(params *admins.AdminSignInParamsAS) (*admins.ResponseAdminSignIn, error) {
	fmt.Println("here")
	var responseModel admins.ResponseAdminSignIn
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.AdminSignIn)
	if err != nil {
		fmt.Println(err.Error())
		return &admins.ResponseAdminSignIn{}, err
	}
	if response == nil {
		return &admins.ResponseAdminSignIn{}, errors.New("response is nil somehow")
	}
	statusCode := int64(response.StatusCode())
	if statusCode != 200 {
		return &admins.ResponseAdminSignIn{}, fmt.Errorf("status code {%d}", statusCode)
	}

	return &responseModel, nil
}

func (c *Client) ConfirmEmailByHash(params *users.ConfirmEmailByHashParams) *ResponseConfirmEmailByHash {
	var responseModel ResponseConfirmEmailByHash
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ConfirmEmailByHash)
	if err != nil {
		fmt.Println(response)
		fmt.Println(err.Error())
	}
	return &responseModel
}

func (c Client) RecoveryInit(params *RecoveryInitParams) *ResponseRecovery {
	var responseModel ResponseRecovery
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.AuthRecoveryInit)
	if err != nil {
		return &ResponseRecovery{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErr,
				StandartCode: 500,
			},
		}
	}
	if response == nil {
		return &ResponseRecovery{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErrNil,
				StandartCode: 500,
			},
		}
	}
	if !response.IsSuccess() {
		return &ResponseRecovery{
			Data:  responseModel.Data,
			Error: &cErrors.ResponseErrorModel{},
		}
	}

	return &responseModel
}

func (c Client) RecoveryConfirm(params *RecoveryConfirmParams) *ResponseRecoveryConfirm {
	var responseModel ResponseRecoveryConfirm
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.AuthRecoveryPasswordByEmail)
	if err != nil {
		return &ResponseRecoveryConfirm{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErr,
				StandartCode: 500,
			},
		}
	}
	if response == nil {
		return &ResponseRecoveryConfirm{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErrNil,
				StandartCode: 500,
			},
		}
	}
	if !response.IsSuccess() {
		return &ResponseRecoveryConfirm{
			Data:  responseModel.Data,
			Error: &cErrors.ResponseErrorModel{},
		}
	}

	return &responseModel
}

func (c *Client) ConfirmEmail(params *users.ConfirmEmailAuthParams) *users.ResponseConfirmEmail {
	var responseModel users.ResponseConfirmEmail
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.ConfirmEmailReq)
	if err != nil {
		return &users.ResponseConfirmEmail{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErr,
				StandartCode: 500,
			},
		}
	}
	if response == nil {
		return &users.ResponseConfirmEmail{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErrNil,
				StandartCode: 500,
			},
		}
	}
	if !response.IsSuccess() {
		return &users.ResponseConfirmEmail{
			Data:  responseModel.Data,
			Error: &cErrors.ResponseErrorModel{},
		}
	}

	return &responseModel
}

func (c *Client) AddTotp(params *AddTotpParams) *ResponseAddTotp {
	var responseModel ResponseAddTotp
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.AddTotp)
	if err != nil {
		return &ResponseAddTotp{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErr,
				StandartCode: 500,
			},
		}
	}
	if response == nil {
		return &ResponseAddTotp{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErrNil,
				StandartCode: 500,
			},
		}
	}
	if !response.IsSuccess() {
		return &ResponseAddTotp{
			Data:  responseModel.Data,
			Error: &cErrors.ResponseErrorModel{},
		}
	}

	return &responseModel
}

func (c *Client) VerifyTotp(params *VerifyTotpParams) *ResponseVerifyTotp {
	var responseModel ResponseVerifyTotp
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.VerTotp)
	if err != nil {
		return &ResponseVerifyTotp{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErr,
				StandartCode: 500,
			},
		}
	}
	if response == nil {
		return &ResponseVerifyTotp{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErrNil,
				StandartCode: 500,
			},
		}
	}
	if !response.IsSuccess() {
		return &ResponseVerifyTotp{
			Data:  responseModel.Data,
			Error: &cErrors.ResponseErrorModel{},
		}
	}

	return &responseModel
}

func (c *Client) VerifyTotpInit(params *VerifyTotpParams) *ResponseVerifyTotp {
	var responseModel ResponseVerifyTotp
	response, err := c.httpClient.R().SetResult(&responseModel).SetBody(params).Post(cConstants.VerTotpInit)
	if err != nil {
		return &ResponseVerifyTotp{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErr,
				StandartCode: 500,
			},
		}
	}
	if response == nil {
		return &ResponseVerifyTotp{
			Data: responseModel.Data,
			Error: &cErrors.ResponseErrorModel{
				InternalCode: cErrors.RatingServiceErrNil,
				StandartCode: 500,
			},
		}
	}
	if !response.IsSuccess() {
		return &ResponseVerifyTotp{
			Data:  responseModel.Data,
			Error: &cErrors.ResponseErrorModel{},
		}
	}

	return &responseModel
}
