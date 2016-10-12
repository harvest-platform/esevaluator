package elastic

// EqualityTerm checks if a value is equal or not equal to another value
type EqualityTerm struct {
	field string
	value interface{}
	sign  string
}

// Translate translates an equality term into an ES "term" term
func (e *EqualityTerm) Translate() map[string]interface{} {
	term := map[string]interface{}{
		"term": map[string]interface{}{
			e.field: e.value,
		},
	}
	if e.sign == "-" {
		boolTerm := Negate(term)
		return boolTerm
	}
	return term
}

// DefinitionTerm checks if a value is defined for a given field
type DefinitionTerm struct {
	field   string
	defined bool
}

// Translate translates a definition term into an ES "exists" term
func (d *DefinitionTerm) Translate() map[string]interface{} {
	existsTerm := map[string]interface{}{
		"exists": map[string]interface{}{
			"field": d.field,
		},
	}
	if !d.defined {
		boolTerm := Negate(existsTerm)
		return boolTerm
	}
	return existsTerm
}
