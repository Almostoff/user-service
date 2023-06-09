package iConnection

type UseCase interface {
	GetInnerConnection(params *GetInnerConnectionParams) (*InnerConnection, error)
	Validate(params *ValidateParams) (*bool, error)
}
