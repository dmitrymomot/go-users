package users

type (
	// Condition struct
	Condition struct {
		query  string
		params []interface{}
	}

	// ConditionFunc func
	ConditionFunc func() Condition
)

// Query string
func (c Condition) Query() string {
	return c.query
}

// Params for query string
func (c Condition) Params() []interface{} {
	return c.params
}

// Disabled func add condition to select items only with disabled=v
func Disabled(v bool) ConditionFunc {
	return func() Condition {
		return Condition{
			query:  "disabled=?",
			params: []interface{}{v},
		}
	}
}

// Confirmed func add condition to select items only with confirmed=v
func Confirmed(v bool) ConditionFunc {
	return func() Condition {
		return Condition{
			query:  "confirmed=?",
			params: []interface{}{v},
		}
	}
}

// OrderBy func add ordering to selected list
func OrderBy(order ...Order) ConditionFunc {
	return func() Condition {
		return Condition{
			query:  "ORDER BY",
			params: []interface{}{order},
		}
	}
}
