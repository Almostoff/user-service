package repository

import (
	"UsersService/internal/cConstants"
	"UsersService/internal/iConnection"
	"UsersService/pkg/secure"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	shield *secure.Shield
	db     *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB, shield *secure.Shield) iConnection.Repository {
	return &postgresRepository{db: db, shield: shield}
}

func (repo *postgresRepository) decryptInnerConnection(innerConnection *iConnection.InnerConnection) {
	innerConnection.Private = repo.shield.DecryptMessageLink(*innerConnection.Private)
	innerConnection.BaseUrl = repo.shield.DecryptMessageLink(*innerConnection.BaseUrl)
}

func (repo *postgresRepository) GetInnerConnection(params *iConnection.GetInnerConnectionParams) (*iConnection.InnerConnection, error) {
	var data []iConnection.InnerConnection
	var d iConnection.InnerConnection

	if err := repo.db.Select(&data, cConstants.GetInnerConnectionQuery, params.Name); err != nil {
		return &iConnection.InnerConnection{}, err
	}

	//if len(data) == 0 {
	//	return &iConnection.InnerConnection{}, fmt.Errorf("connection with name {%s} not found", *params.Name)
	//}

	//repo.decryptInnerConnection(&data[0])
	//return &data[0], nil
	return &d, nil
}

func (repo *postgresRepository) GetServiceByPublic(params *iConnection.GetServiceByPublicParams) (*iConnection.InnerConnection, error) {

	var data []iConnection.InnerConnection

	if err := repo.db.Select(&data, cConstants.GetServiceByPublicQuery, params.Public); err != nil {
		return &iConnection.InnerConnection{}, err
	}

	if len(data) == 0 {
		return &iConnection.InnerConnection{}, fmt.Errorf("connection with public {%s} not found", *params.Public)
	}

	repo.decryptInnerConnection(&data[0])
	return &data[0], nil
}
