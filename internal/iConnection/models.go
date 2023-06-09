package iConnection

type InnerConnection struct {
	Public  *string `json:"public" db:"public"`
	BaseUrl *string `json:"-" db:"base_url"`
	Private *string `json:"-" db:"private"`
	Name    *string `json:"name" db:"name"`
	Id      *int64  `json:"-" db:"id"`
	Test    *bool   `json:"-" db:"test"`
}

// ----------------------------------------------------------------------------------------------------

type ValidateParams struct {
	Signature string
	Message   string
	Public    string
	Timestamp string
}

type GetInnerConnectionParams struct {
	Name *string `json:"name"`
}

func (g GetInnerConnectionParams) GetInnerConnection(params *GetInnerConnectionParams) (*InnerConnection, error) {
	return &InnerConnection{}, nil
}

func (g GetInnerConnectionParams) GetServiceByPublic(params *GetServiceByPublicParams) (*InnerConnection, error) {
	return &InnerConnection{}, nil

}

type GetServiceByPublicParams struct {
	Public *string `json:"public"`
}
