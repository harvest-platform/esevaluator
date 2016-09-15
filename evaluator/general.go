package evaluator

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
	if !defined {
		boolTerm := Negate(existsTerm)
		return boolTerm
	}
	return existsTerm
}
