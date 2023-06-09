package rating

type ServiceRating interface {
	GetClientSRRating(params *GetClientSRRatingParams) *ResponseGetClientSRRating
	UpdateCommentAdmin(params *UpdateCommentAdminParams) *ResponseUpdateComment
	UpdateCommentClient(params *UpdateCommentClientParams) *ResponseUpdateCommentClient
	AddComment(params *AddCommentParams) *ResponseAddComment
	GetAllReviews(params *GetAllReviewsParams) *ResponseGetAllReviews
	GetClientRatingForOrders(params *GetClientSRRatingParams) *ResponseGetClientStatistics

	IsCommentExist(params *IsCommentExistParams) *ResponseIsCommentExist
}
