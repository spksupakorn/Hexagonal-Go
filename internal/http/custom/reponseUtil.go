package custom

import (
	"dungeons-dragon-service/internal/dto"
	"reflect"
)

func Null() interface{} {
	return nil
}

func getLength(data interface{}) int {
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		return val.Len()
	}
	return 0
}

func BuildResponse_[T any](status bool, message string, data T) dto.ApiResponse[T] {
	return dto.ApiResponse[T]{
		Error:   status,
		Message: message,
		Data:    data,
		Length:  getLength(data),
	}
}

func BuildResponse[T any](responseStatus ResponseStatus, data T) dto.ApiResponse[T] {
	return BuildResponse_(responseStatus.GetResponseStatus(), responseStatus.GetResponseMessage(), data)
}

func BuildResponseWithPaginate_[T any, P any](status bool, message string, data T, paginateData P) dto.ApiResponseWithPaginate[T, P] {
	return dto.ApiResponseWithPaginate[T, P]{
		Error:        status,
		Message:      message,
		Data:         data,
		PaginateData: paginateData,
		Length:       getLength(data),
	}
}

func BuildResponseWithPaginate[T any, P any](responseStatus ResponseStatus, data T, paginateData P) dto.ApiResponseWithPaginate[T, P] {
	return BuildResponseWithPaginate_(responseStatus.GetResponseStatus(), responseStatus.GetResponseMessage(), data, paginateData)
}
