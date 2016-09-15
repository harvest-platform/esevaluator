package evaluator

type EmptinessTerm struct {
	field string
	empty string
}

// Translate translates an Emptiness term into a corresponsing ES "term" query
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

type MemberTerm struct {
	field string
	value string
}

// Translate translates an Emptiness term into a corresponsing ES "terms" query
func (m *MemberTerm) Translate() {
	term := map[string]interface{}{
		"term": map[string]string{
			m.field: m.value,
		},
	}
	return term
}

type SubsetTerm struct {
	field string
	set   []interface{}
}

func (s *SubsetTerm) Translate() {
	setSize = len(s.set)

	term := map[string]interface{}{
		"terms": map[string]interface{}{
			s.field:                s.set,
			"minimum_should_match": setSize,
		},
	}

	return term
}
