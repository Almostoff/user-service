package usecase

import (
	"UsersService/internal/cConstants"
	"UsersService/internal/iConnection"
	"UsersService/pkg/secure"
	"errors"
	"fmt"
	"strings"
)

type iConnectionUseCase struct {
	repo iConnection.Repository
}

func NewIConnectionUsecase(repo iConnection.Repository) iConnection.UseCase {
	return &iConnectionUseCase{repo: repo}
}

func (u *iConnectionUseCase) Validate(params *iConnection.ValidateParams) (*bool, error) {

	connection, err := u.repo.GetServiceByPublic(&iConnection.GetServiceByPublicParams{Public: &params.Public})
	if err != nil {
		return &cConstants.False, err
	}

	body := createRequestBody(params.Timestamp, strings.TrimRight(params.Message, "\n"))
	hash := secure.CalcSignature(*connection.Private, body)
	if hash != params.Signature {
		return &cConstants.False, fmt.Errorf("hash {%s} != signature {%s}, %s, %s", hash, params.Signature, *connection.Private, body)
	}

	return &cConstants.True, nil
}

func createRequestBody(timestamp, jsonBody string) string {
	return timestamp + jsonBody
}

func (u *iConnectionUseCase) GetInnerConnection(params *iConnection.GetInnerConnectionParams) (*iConnection.InnerConnection, error) {

	connection, err := u.repo.GetInnerConnection(&iConnection.GetInnerConnectionParams{Name: params.Name})
	if err != nil {
		return &iConnection.InnerConnection{}, err
	}
	if connection == nil {
		return &iConnection.InnerConnection{}, errors.New("connection is nil somehow")
	}
	return connection, nil
}
