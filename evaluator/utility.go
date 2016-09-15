package evaluator

// Negate wraps a term in a boolean "not" condition
func Negate(term map[string]interface{}) map[string]interface{} {
	boolTerm := map[string]interface{}{
		"bool": map[string]interface{}{
			"not": term,
		},
	}
	return boolTerm
}

// Nest wraps a term in a "nested" condition
func Nest(term map[string]interface{}) map[string]interface{} {
	nestTerm := map[string]interface{}{
		"nested": term,
	}
	return nestTerm
}
