package model

import (
	"fmt"
	"strconv"

	coreerror "gitlab.com/alienspaces/go-mud/server/core/error"
	"gitlab.com/alienspaces/go-mud/server/core/repository"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
)

const (
	PageSize   = "page_size"
	PageNumber = "page_number"
)

const (
	DefaultPageSize              = "10"
	DefaultPageNumber            = "1"
	DefaultOrderDescendingColumn = "created_at"
)

func ResolvePaginationParams(queryParams map[string]interface{}, l logger.Logger) (map[string]interface{}, map[string]string, int, error) {
	params := map[string]interface{}{}
	for k, v := range queryParams {
		params[k] = v
	}

	operators := map[string]string{}

	params, operators, adjustedPageSize, err := resolveLimit(params, operators)
	pageSize := adjustedPageSize - 1
	if err != nil {
		l.Warn(fmt.Sprintf("failed to resolve limit >%v<", err))
		return params, operators, pageSize, err
	}

	params, operators, err = resolveOffset(params, operators, pageSize)
	if err != nil {
		l.Warn(fmt.Sprintf("failed to resolve offset >%v<", err))
		return params, operators, pageSize, err
	}

	params, operators = resolveSortOrder(params, operators)

	return params, operators, pageSize, err
}

func resolveLimit(params map[string]interface{}, operators map[string]string) (map[string]interface{}, map[string]string, int, error) {
	params, pageSize, err := extractParam(params, PageSize, DefaultPageSize)
	if err != nil {
		return params, operators, 0, err
	}
	if pageSize < 1 {
		return params, operators, 0, coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< is less than 1 >%d<", PageSize, pageSize))
	}

	limit := pageSize + 1
	operators[repository.OperatorLimit] = strconv.Itoa(limit)

	return params, operators, limit, nil
}

func resolveOffset(params map[string]interface{}, operators map[string]string, pageSize int) (map[string]interface{}, map[string]string, error) {
	params, pageNumber, err := extractParam(params, PageNumber, DefaultPageNumber)
	if err != nil {
		return params, operators, err
	}
	if pageNumber < 1 {
		return params, operators, coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< is less than 1 >%d<", PageNumber, pageNumber))
	}

	offset := (pageNumber - 1) * pageSize
	operators[repository.OperatorOffset] = strconv.Itoa(offset)
	return params, operators, nil
}

// TODO sorting collection endpoint will be resolved in future story
func resolveSortOrder(params map[string]interface{}, operators map[string]string) (map[string]interface{}, map[string]string) {
	operators[repository.OperatorOrderByDescending] = DefaultOrderDescendingColumn
	return params, operators
}

// extractParam extracts the value associated with the key and removes the key, mutating the params map.
// The map value is expected to be a string slice.
func extractParam(params map[string]interface{}, key string, defaultValue string) (map[string]interface{}, int, error) {
	value, ok := params[key]
	if !ok {
		value = []string{defaultValue}
	}

	valueStr, ok := value.([]string)
	if !ok {
		return params, 0, coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< has an invalid value >%v<", key, value))
	}

	if len(valueStr) != 1 {
		return params, 0, coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< should be a single value but is >%+v<", key, valueStr))
	}

	valueInt, err := strconv.Atoi(valueStr[0])
	if err != nil {
		return params, 0, coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< has an invalid value >%v<", key, value))
	}

	delete(params, key)
	return params, valueInt, nil
}
