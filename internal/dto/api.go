package dto

type ApiResponse[T any] struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Length  int    `json:"length"`
}

type ApiResponseWithPaginate[T any, P any] struct {
	Error        bool   `json:"error"`
	Message      string `json:"message"`
	Data         T      `json:"data"`
	PaginateData P      `json:"paginate"`
	Length       int    `json:"length"`
}

// for swagger documentation
type APIErrorResponse struct {
	Error   bool   `json:"error"   example:"true"`
	Message string `json:"message" example:"Invalid credentials"`
	Data    string `json:"data"    example:"null"`
	Length  int    `json:"length"  example:"0"`
}

type APIObjectResponse struct {
	Error   bool   `json:"error"   example:"false"`
	Message string `json:"message" example:"Success"`
	Data    any    `json:"data"`
	Length  int    `json:"length"  example:"171"`
}
