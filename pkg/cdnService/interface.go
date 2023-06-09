package cdnService

type UseCase interface {
	SaveImage(params *SaveImageParams) (*SaveImageResponse, error)
}
