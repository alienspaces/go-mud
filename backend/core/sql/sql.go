package sql

import (
	"fmt"
	"strings"

	"gitlab.com/alienspaces/go-mud/backend/core/collection/set"
)

type Operator string

const (
	OpBetween            Operator = "BETWEEN"
	OpILike              Operator = "ILIKE"
	OpLike               Operator = "LIKE"
	OpEqualTo            Operator = "="
	OpNotEqualTo         Operator = "!="
	OpLessThanEqualTo    Operator = "<="
	OpLessThan           Operator = "<"
	OpGreaterThanEqualTo Operator = ">="
	OpGreaterThan        Operator = ">"
	OpAny                Operator = "ANY"
	OpNotIn              Operator = "NOT IN"
	OpIn                 Operator = "IN"
	OpContains           Operator = "@>"
	OpContainedBy        Operator = "<@"
	OpOverlap            Operator = "&&"
	OpIsNull             Operator = "IS NULL"
	OpIsNotNull          Operator = "IS NOT NULL"
)

type Options struct {
	Params  []Param
	OrderBy []OrderBy
	Limit   int
	Offset  int
	*Lock
}

// Param defines operators and operands for an SQL query (e.g., customer_countries @> array[:customer_countries]).
//
// Col defines the column that the operator applies to (e.g., customer_countries). Col should always be specified.
//
// Op defines the operator (e.g., @>). Op should not be specified for the following operators: IN, =, @>, ANY. Doing so
// short-circuits the core repository processing of Options, and requires Array instead of Val to be specified for IN and @> operators.
//
// Val defines the operand for the operator (e.g., an int, []int, string or []string).
// Val should only ever be empty for operators that have no operands (e.g., IS NULL).
//
// ValB defines a second operand for operators requiring two operands (e.g., BETWEEN).
//
// Array should not be manually specified. If Val is a slice, Array is automatically populated and Val is emptied by the core repository for the following operators:
// IN, &&, @>, <@ and any other named array operands.
type Param struct {
	Col string
	Op  Operator

	// Operator Param operands
	Val   any
	ValB  any
	Array []any
}

var ArrayOps = set.New[Operator](OpIn, OpNotIn, OpOverlap, OpContains, OpContainedBy)

type OrderBy struct {
	Col       string
	Direction OrderDirection
}

type OrderDirection string

var (
	OrderDirectionASC  OrderDirection = "ASC"
	OrderDirectionDESC OrderDirection = "DESC"
)

// Lock with LockStrength LockStrengthUpdate should be specified when:
//
// (1) updating a key field (e.g., primary key, composite key); or
//
// (2) updating a unique field that is used as a FK.
//
// LockOption should generally not need to be specified in the API and model layers.
// The default lock conflict behaviour is to wait for the lock to be released.
//
// LockOptionSkipLocked may be useful for pulling from job queues.
type Lock struct {
	Strength LockStrength
	Option   LockOption
}

// LockStrength details can be found here: https://www.postgresql.org/docs/current/explicit-locking.html
type LockStrength string

var (
	// LockStrengthUpdate is the default locking option and should be used when:
	//
	// (1) updating a key field (e.g., primary key, composite key); or
	//
	// (2) updating a unique field that is used as a FK.
	LockStrengthUpdate LockStrength = "UPDATE"
	// LockStrengthNoKeyUpdate may be used when the update does not delete or update any key fields or any unique fields
	// that are used as a foreign key.
	//
	// Further details can be found here: https://www.postgresql.org/docs/current/explicit-locking.html
	LockStrengthNoKeyUpdate LockStrength = "NO KEY UPDATE"
	LockStrengthShare       LockStrength = "SHARE"
	LockStrengthKeyShare    LockStrength = "KEY SHARE"
)

// LockOption details can be found here: https://www.postgresql.org/docs/current/sql-select.html#:~:text=To%20prevent%20the%20operation%20from,be%20immediately%20locked%20are%20skipped.
type LockOption string

var (
	LockOptionSkipLocked LockOption = "SKIP LOCKED"
	LockOptionNoWait     LockOption = "NOWAIT"
)

var (
	// ForUpdate should be used by default in the handler layer.
	ForUpdate = &Lock{
		Strength: LockStrengthUpdate,
	}
	// ForUpdateSkipLocked may be used by any job queue implementations.
	ForUpdateSkipLocked = &Lock{
		Strength: LockStrengthUpdate,
		Option:   LockOptionSkipLocked,
	}
	// ForUpdateNoWait should be used by default in the model layer.
	ForUpdateNoWait = &Lock{
		Strength: LockStrengthUpdate,
		Option:   LockOptionNoWait,
	}
)

func From(initialSQL string, opts *Options) (string, map[string]any, error) {
	if opts == nil {
		return initialSQL, nil, nil
	}

	sql := initialSQL
	queryArgs := map[string]any{}

	for _, param := range opts.Params {
		op := param.Op

		if ArrayOps.Contains(op) && len(param.Array) == 0 {
			return "", nil, fmt.Errorf("missing param Array for op >%s< sql >%s<", op, sql)
		}

		if op != OpIsNull && op != OpIsNotNull {
			if param.ValB != nil {
				queryArgs[param.Col+"0"] = param.Val
				queryArgs[param.Col+"1"] = param.ValB
			} else if param.Val != nil {
				queryArgs[param.Col] = param.Val
			} else if len(param.Array) > 0 {
				for i, a := range param.Array {
					col := fmt.Sprintf("%s%d", param.Col, i)
					queryArgs[col] = a
				}
			} else {
				return "", nil, fmt.Errorf("missing param A for op >%s< sql >%s<", op, sql)
			}
		}

		opClause := "AND "
		switch op {
		case OpEqualTo, OpNotEqualTo, OpLessThanEqualTo, OpGreaterThanEqualTo, OpLessThan, OpGreaterThan, OpLike, OpILike:
			opClause += fmt.Sprintf("%s %s :%s", param.Col, op, param.Col)
		case OpBetween:
			opClause += fmt.Sprintf("%s %s :%s0 AND :%s1", param.Col, op, param.Col, param.Col)
			if param.ValB == "" {
				return "", nil, fmt.Errorf("missing param B for op >%s< column name >%s< sql >%s<", op, param.Col, sql)
			}
		case OpIn, OpNotIn:
			namedParams := toNamedArrayParams(param.Array, param.Col)
			opClause += fmt.Sprintf("%s %s (%s)", param.Col, op, namedParams)
		case OpContains, OpContainedBy, OpOverlap:
			namedParams := toNamedArrayParams(param.Array, param.Col)
			opClause += fmt.Sprintf("%s %s array[%s]", param.Col, op, namedParams)
		case OpAny:
			opClause += fmt.Sprintf(":%s = %s(%s)", param.Col, op, param.Col)
		case OpIsNull, OpIsNotNull:
			opClause += fmt.Sprintf("%s %s", param.Col, op)
		default:
			return "", nil, fmt.Errorf("unknown op >%s< for >%s<", op, sql)
		}

		sql += fmt.Sprintf("%s\n", opClause)
	}

	if len(opts.OrderBy) > 0 {
		orderBy := "ORDER BY"
		for _, o := range opts.OrderBy {
			orderBy += fmt.Sprintf("\n%s %s,", o.Col, o.Direction)
		}
		orderBy = strings.TrimSuffix(orderBy, ",")
		sql += fmt.Sprintf("%s\n", orderBy)
	}

	if opts.Limit != 0 {
		sql += fmt.Sprintf("LIMIT %d\n", opts.Limit)
	}

	if opts.Offset != 0 {
		sql += fmt.Sprintf("OFFSET %d\n", opts.Offset)
	}

	if opts.Lock != nil && opts.Lock.Strength != "" {
		sql += fmt.Sprintf("FOR %s %s\n", opts.Lock.Strength, opts.Lock.Option)
	}

	return sql, queryArgs, nil
}

func toNamedArrayParams(array []any, columnName string) string {
	var namedParams string
	for i := range array {
		namedParams += fmt.Sprintf(":%s%d,", columnName, i)
	}

	return strings.TrimSuffix(namedParams, ",")
}
