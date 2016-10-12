package elastic

// Negate wraps a term in a boolean "not" condition
func Negate(term map[string]interface{}) map[string]interface{} {
	boolTerm := map[string]interface{}{
		"bool": map[string]interface{}{
			"not": term,
		},
	}
	return boolTerm
}

// Nest wraps a term in a "nested" statement
func Nest(term map[string]interface{}) map[string]interface{} {
	nestTerm := map[string]interface{}{
		"nested": term,
	}
	return nestTerm
}

// Query wraps a term in a "query" statement
func Query(term map[string]interface{}) map[string]interface{} {
	queryTerm := map[string]interface{}{
		"query": term,
	}
	return queryTerm
}

// Filter wraps a term in a "filter" statement
func Filter(term map[string]interface{}) map[string]interface{} {
	filterTerm := map[string]interface{}{
		"filter": term,
	}
	return filterTerm
}
