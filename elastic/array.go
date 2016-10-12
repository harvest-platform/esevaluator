package elastic

// EmptinessTerm checks if an array field is empty
type EmptinessTerm struct {
	field string
	empty bool
}

// Translate translates an Emptiness term into a corresponsing ES "term" statement
func (e *EmptinessTerm) Translate() map[string]interface{} {
	constantTerm := map[string]interface{}{
		"constant_score": map[string]interface{}{
			"filter": map[string]interface{}{
				"field":     e.field,
				"existence": true,
			},
		},
	}
	if !e.empty {
		boolTerm := Negate(constantTerm)
		return boolTerm
	}
	return constantTerm
}

// MemberTerm checks if a field value is an element of an array
type MemberTerm struct {
	field string
	value interface{}
}

// Translate translates an Emptiness term into a corresponding ES "terms" statement
func (m *MemberTerm) Translate() map[string]interface{} {
	term := map[string]interface{}{
		"term": map[string]interface{}{
			m.field: m.value,
		},
	}
	return term
}

// SubsetTerm checks if an array is a subset of an array field
type SubsetTerm struct {
	field string
	set   []interface{}
}

// Translate translates a SubsetTerm into a corresponding ES "terms" statement
func (s *SubsetTerm) Translate() map[string]interface{} {
	setSize := len(s.set)

	term := map[string]interface{}{
		"terms": map[string]interface{}{
			s.field:                s.set,
			"minimum_should_match": setSize,
		},
	}

	return term
}
