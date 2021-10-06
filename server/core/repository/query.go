package repository

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrOperatorNotSupported -
	ErrOperatorNotSupported = errors.New("operator not supported")

	// OperatorIsNull -
	OperatorIsNull = "__is_null"
	// OperatorIsNotNull -
	OperatorIsNotNull = "__is_not_null"
	// OperatorOrderByAscending -
	OperatorOrderByAscending = "__order_by_asc"
	// OperatorOrderByDescending -
	OperatorOrderByDescending = "__order_by_desc"
	// OperatorLimit -
	OperatorLimit = "__limit"
	// OperatorOffset -
	OperatorOffset = "__offset"

	// OperatorBetween -
	OperatorBetween = "between"
	// OperatorLike -
	OperatorLike = "like"
	// OperatorNotEqualTo -
	OperatorNotEqualTo = "!="
	// OperatorLessThanEqualTo -
	OperatorLessThanEqualTo = "<="
	// OperatorLessThan -
	OperatorLessThan = "<"
	// OperatorGreaterThanEqualTo -
	OperatorGreaterThanEqualTo = ">="
	// OperatorGreaterThan -
	OperatorGreaterThan = ">"
)

// sqlFromParamsAndOperator -
func (r *Repository) sqlFromParamsAndOperator(
	initialSQL string,
	params map[string]interface{},
	operators map[string]string) (string, map[string]interface{}, error) {

	sqlStmt := initialSQL

	// Copy the params into a new map
	queryParams := map[string]interface{}{}
	for k, v := range params {
		queryParams[k] = v
	}

	for param, val := range params {

		operator, found := operators[param]
		if !found {
			sqlStmt += fmt.Sprintf("AND %s", param)

			switch val.(type) {
			case []string:
				if len(val.([]string)) == 0 {
					return "", nil, fmt.Errorf("value for param >%s< is empty", param)
				}
				// in
				sqlStmt += " IN ("
				for idx, paramVal := range val.([]string) {
					paramName := fmt.Sprintf("%s%d", param, idx)
					sqlStmt += fmt.Sprintf(":%s", paramName)
					queryParams[paramName] = paramVal
					if idx < len(val.([]string))-1 {
						sqlStmt += ", "
					}
				}
				sqlStmt += ")\n"

				// delete original param
				delete(queryParams, param)

			case []int64:
				if len(val.([]int64)) == 0 {
					return "", nil, fmt.Errorf("value for param >%s< is empty", param)
				}
				// in
				sqlStmt += " IN ("
				for idx, paramVal := range val.([]int64) {
					paramName := fmt.Sprintf("%s%d", param, idx)
					sqlStmt += fmt.Sprintf(":%s", paramName)
					queryParams[paramName] = paramVal
					if idx < len(val.([]int64))-1 {
						sqlStmt += ", "
					}
				}
				sqlStmt += ")\n"

				// delete original param
				delete(queryParams, param)

			default:
				// equals
				sqlStmt += fmt.Sprintf(" = :%s\n", param)
			}
			continue
		}

		switch operator {
		case OperatorBetween:
			valStr, ok := val.(string)
			if !ok {
				return "", nil, fmt.Errorf("value for param >%s< is not a string", param)
			}
			split := strings.Split(valStr, ",")
			if len(split) != 2 {
				return "", nil, fmt.Errorf("Param >%s< should have 2 values separated by a comma", param)
			}

			firstParamName := param + "_1"
			secondParamName := param + "_2"

			sqlStmt += fmt.Sprintf("AND %s >= :%s\n", param, firstParamName)
			sqlStmt += fmt.Sprintf("AND %s <= :%s\n", param, secondParamName)

			// delete the old param from the queryParams.
			delete(queryParams, param)

			// add the new params to queryParams.
			queryParams[firstParamName] = split[0]
			queryParams[secondParamName] = split[1]

		case OperatorNotEqualTo, OperatorLessThanEqualTo, OperatorGreaterThanEqualTo, OperatorLessThan, OperatorGreaterThan, OperatorLike:
			sqlStmt += fmt.Sprintf("AND %s %s :%s\n", param, operator, param)
		default:
			return "", nil, ErrOperatorNotSupported
		}

	}

	// IS NULL
	if val, ok := operators[OperatorIsNull]; ok {
		vals := strings.Split(val, ",")
		for _, val := range vals {
			sqlStmt += fmt.Sprintf("AND %s IS NULL\n", val)
		}
	}

	// IS NOT NULL
	if val, ok := operators[OperatorIsNotNull]; ok {
		vals := strings.Split(val, ",")
		for _, val := range vals {
			sqlStmt += fmt.Sprintf("AND %s IS NOT NULL\n", val)
		}
	}

	// ORDER BY
	if val, ok := operators[OperatorOrderByAscending]; ok {
		sqlStmt += fmt.Sprintf("ORDER BY %s ASC\n", val)
	}
	if val, ok := operators[OperatorOrderByDescending]; ok {
		sqlStmt += fmt.Sprintf("ORDER BY %s DESC\n", val)
	}

	// LIMIT
	if val, ok := operators[OperatorLimit]; ok {
		sqlStmt += fmt.Sprintf("LIMIT %s\n", val)
	}

	// OFFSET
	if val, ok := operators[OperatorOffset]; ok {
		sqlStmt += fmt.Sprintf("OFFSET %s\n", val)
	}

	return sqlStmt, queryParams, nil
}
