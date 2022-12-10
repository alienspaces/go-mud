package schema

// Request -
type Request struct {
	Pagination RequestPagination `json:"pagination"`
}

// RequestPagination -
type RequestPagination struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

// Response -
type Response struct {
	// Populated when response code is not 2XX
	Error *ResponseError `json:"error,omitempty"`
	// Populated when response is a collection
	Pagination *ResponsePagination `json:"pagination,omitempty"`
}

// ResponseError -
type ResponseError struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// ResponsePagination -
type ResponsePagination struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	PageCount  int `json:"page_count"`
}
