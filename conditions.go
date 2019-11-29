package users

import (
	"strings"

	"github.com/bradfitz/slice"
)

// Predefined query section type
const (
	QsWhere QuerySection = iota + 1
	QsOrderBy
	QsLimit
	QsOffset
)

type (
	// Condition struct
	Condition interface {
		Query() string
		Params() []interface{}
		Type() QuerySection
	}

	// QuerySection type
	QuerySection int

	condition struct {
		query  string
		params []interface{}
		t      QuerySection
	}
)

// Query string
func (c condition) Query() string {
	return c.query
}

// Params for query string
func (c condition) Params() []interface{} {
	return c.params
}

// Type of query section string
func (c condition) Type() QuerySection {
	return c.t
}

// Disabled func add condition to select items only with disabled=v
func Disabled(v bool) Condition {
	return condition{
		query:  "disabled=?",
		params: []interface{}{v},
		t:      QsWhere,
	}
}

// Confirmed func add condition to select items only with confirmed=v
func Confirmed(v bool) Condition {
	return condition{
		query:  "confirmed=?",
		params: []interface{}{v},
		t:      QsWhere,
	}
}

// OrderBy func add ordering to selected list
func OrderBy(order ...Order) Condition {
	q := "ORDER BY"
	for k, o := range order {
		if k > 0 {
			q = q + ", " + o.String()
		} else {
			q = q + " " + o.String()
		}
	}
	return condition{
		query: q,
		t:     QsOrderBy,
	}
}

// Limit func add limit to select query
func Limit(v int) Condition {
	if v <= 0 {
		v = 100
	}
	return condition{
		query:  "LIMIT ?",
		params: []interface{}{v},
		t:      QsLimit,
	}
}

// Offset func add offset to select query
func Offset(v int) Condition {
	if v < 0 {
		v = 0
	}
	return condition{
		query:  "OFFSET ?",
		params: []interface{}{v},
		t:      QsOffset,
	}
}

// ConditionsToQuery represents conditions slice to query string and parameters slice
func ConditionsToQuery(cs ...Condition) (q string, params []interface{}) {
	if len(cs) == 0 {
		return
	}

	slice.Sort(cs[:], func(i, j int) bool {
		return cs[i].Type() < cs[j].Type()
	})

	var where, limit, offset, orderBy string
	for _, c := range cs {
		switch c.Type() {
		case QsWhere:
			where += " AND " + c.Query()
		case QsLimit:
			limit = c.Query()
		case QsOffset:
			offset = c.Query()
		case QsOrderBy:
			orderBy = c.Query()
		}
		params = append(params, c.Params()...)
	}

	if where != "" {
		q = " WHERE " + strings.TrimPrefix(where, " AND ")
	}
	if orderBy != "" {
		q += " " + orderBy
	}
	if limit != "" {
		q += " " + limit
	}
	if offset != "" {
		q += " " + offset
	}

	return q, params
}
