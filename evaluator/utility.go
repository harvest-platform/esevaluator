package evaluator

// Negate wraps a term in a boolean "not" condition
func Negate(term map[string]interface{}) map[string]interface{} {
	var notTerm map[string]interface{}
	notTerm = make(map[string]interface{})

	var boolTerm map[string]interface{}
	boolTerm = make(map[string]interface{})

	notTerm["not"] = term
	boolTerm["bool"] = notTerm

	return boolTerm
}

// Nest wraps a term in a "nested" condition
func Nest(term map[string]interface{}) map[string]interface{} {
	var nestTerm map[string]interface{}
	nestTerm = make(map[string]interface{})

	nestTerm["nested"] = term
	return nestTerm
}
