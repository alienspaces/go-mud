package model

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
)

const (
	PageSize   = "page_size"
	PageNumber = "page_number"
	SortColumn = "sort_column"
)

const (
	DefaultPageSize              = "10"
	DefaultPageNumber            = "1"
	DefaultOrderDescendingColumn = "created_at"
)

func ResolveQueryParams(l logger.Logger, queryParams map[string]any) (*coresql.Options, int, error) {
	params := map[string]any{}
	for k, v := range queryParams {
		params[k] = v
	}

	params, opts, pageSize, err := resolvePagination(params, l)
	if err != nil {
		l.Warn("failed to resolve pagination params", err)
		return nil, 0, err
	}

	params, opts, err = resolveSortOrder(params, opts)
	if err != nil {
		l.Warn("failed to resolve sort_column params", err)
		return nil, 0, err
	}

	for col, val := range params {
		valueStr, ok := val.([]string)
		if !ok {
			return nil, 0, fmt.Errorf("query param value with key >%s< not []string type >%#v<, but >%s<", col, valueStr, reflect.TypeOf(valueStr).Kind())
		}

		a := val
		if len(valueStr) == 1 {
			a = valueStr[0]
		}

		opts.Params = append(opts.Params, coresql.Param{
			Col: col,
			Val: a,
		})
	}

	return opts, pageSize, err
}

// resolvePagination mutates queryParams
func resolvePagination(queryParams map[string]any, l logger.Logger) (map[string]any, *coresql.Options, int, error) {
	queryParams, opts, adjustedPageSize, err := resolveLimit(queryParams, &coresql.Options{})
	pageSize := adjustedPageSize - 1
	if err != nil {
		l.Warn(fmt.Sprintf("failed to resolve limit >%v<", err))
		return queryParams, opts, pageSize, err
	}

	queryParams, opts, err = resolveOffset(queryParams, opts, pageSize)
	if err != nil {
		l.Warn(fmt.Sprintf("failed to resolve offset >%v<", err))
		return queryParams, opts, pageSize, err
	}

	return queryParams, opts, pageSize, err
}

func resolveLimit(params map[string]any, opts *coresql.Options) (map[string]any, *coresql.Options, int, error) {
	params, pageSize, err := extractIntQueryParam(params, PageSize, DefaultPageSize)
	if err != nil {
		return params, opts, 0, err
	}
	if pageSize < 1 {
		return params, opts, 0, coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< is less than 1 >%d<", PageSize, pageSize))
	}

	opts.Limit = pageSize + 1
	return params, opts, pageSize + 1, nil
}

func resolveOffset(params map[string]any, opts *coresql.Options, pageSize int) (map[string]any, *coresql.Options, error) {
	params, pageNumber, err := extractIntQueryParam(params, PageNumber, DefaultPageNumber)
	if err != nil {
		return params, opts, err
	}
	if pageNumber < 1 {
		return params, opts, coreerror.NewQueryParamError(fmt.Sprintf("Query parameter >%s< is less than 1 >%d<", PageNumber, pageNumber))
	}

	opts.Offset = (pageNumber - 1) * pageSize
	return params, opts, nil
}

// extractIntQueryParam extracts the value associated with the key and removes the key, mutating the params map.
// The params map value is expected to be a string slice.
func extractIntQueryParam(params map[string]any, key string, defaultValue string) (map[string]any, int, error) {
	params, valueStr, err := extractQueryParam(params, key)
	if err != nil {
		return params, 0, err
	}
	if valueStr == nil {
		valueStr = []string{defaultValue}
	}

	if len(valueStr) != 1 {
		return params, 0, coreerror.NewQueryParamError(fmt.Sprintf("query parameter >%s< should be a single value but is >%+v<", key, valueStr))
	}

	valueInt, err := strconv.Atoi(valueStr[0])
	if err != nil {
		return params, 0, coreerror.NewQueryParamError(fmt.Sprintf("query parameter >%s< has an invalid value >%+v<", key, valueStr))
	}

	return params, valueInt, nil
}

func resolveSortOrder(params map[string]any, opts *coresql.Options) (map[string]any, *coresql.Options, error) {
	params, sortColumns, err := extractQueryParam(params, SortColumn)
	if err != nil {
		return params, opts, err
	}
	if sortColumns == nil {
		opts.OrderBy = []coresql.OrderBy{
			{
				Col:       DefaultOrderDescendingColumn,
				Direction: coresql.OrderDirectionDESC,
			},
		}

		return params, opts, nil
	}

	for _, col := range sortColumns {
		dir := coresql.OrderDirectionASC
		if strings.HasPrefix(col, "-") {
			col = strings.TrimPrefix(col, "-")
			dir = coresql.OrderDirectionDESC
		}

		opts.OrderBy = append(opts.OrderBy, coresql.OrderBy{
			Col:       col,
			Direction: dir,
		})
	}

	return params, opts, nil
}

func extractQueryParam(params map[string]any, key string) (map[string]any, []string, error) {
	value, ok := params[key]
	if !ok {
		return params, nil, nil
	}

	// Go pkg for query params stores query param values as []string type
	valueStr, ok := value.([]string)
	if !ok {
		return params, nil, fmt.Errorf("query param value with key >%s< not []string type >%#v<, but >%s<", key, valueStr, reflect.TypeOf(valueStr).Kind())
	}

	delete(params, key)
	return params, valueStr, nil
}
