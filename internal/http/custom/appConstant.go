package custom

import "net/http"

type ResponseStatus int
type Headers int
type General int

const (
	Success ResponseStatus = iota + 1
	DataNotFound
	BadRequest
	Unauthorized
	Forbidden
	UnprocessableEntity
	InternalServerError
	Conflict
	NoContent
)

func (r ResponseStatus) GetResponseStatus() bool {
	responseStatus := [...]int{http.StatusOK, http.StatusNotFound, http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusUnprocessableEntity, http.StatusInternalServerError, http.StatusConflict, http.StatusNoContent}[r-1]
	return responseStatus != http.StatusOK
}

func (r ResponseStatus) GetResponseMessage() string {
	return [...]string{"Success", "Data Not Found", "Bad Request", "Unauthorized", "Internal Server Error", "Conflict", "No Content"}[r-1]
}
