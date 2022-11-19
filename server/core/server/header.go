package server

import (
	"fmt"
	"net/http"
)

const (
	HeaderXPagination = "X-Pagination"
)

func XPaginationHeader(collectionLen int, pageSize int) func(http.ResponseWriter) error {
	hasMore := collectionLen > pageSize

	return func(w http.ResponseWriter) error {
		w.Header().Set(HeaderXPagination, XPaginationHeaderValue(hasMore))
		return nil
	}
}

func XPaginationHeaderValue(hasMore bool) string {
	return fmt.Sprintf(`{"has_more":%t}`, hasMore)
}
