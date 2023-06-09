package iConnection

type Repository interface {
	GetInnerConnection(params *GetInnerConnectionParams) (*InnerConnection, error)
	GetServiceByPublic(params *GetServiceByPublicParams) (*InnerConnection, error)
}
