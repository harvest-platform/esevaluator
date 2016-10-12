package elastic

// OneTerm checks if a field values matches one of a given set of values
type OneTerm struct {
	field  string
	values []string
}

// Translate translates a OneTerm into an ES "constant_score" statement
func (o *OneTerm) Translate() map[string]interface{} {
	constantTerm := map[string]interface{}{
		"constant_score": map[string]interface{}{
			"filter": map[string]interface{}{
				"terms": o.values,
			},
		},
	}

	return constantTerm
}

// MatchTerm uses the free text ES match statement
type MatchTerm struct {
	field string
	value string
}

// Translate translates a MatchTerm into an ES "match" statement
func (c *MatchTerm) Translate() map[string]interface{} {
	matchTerm := map[string]interface{}{
		"match": map[string]interface{}{
			c.field: c.value,
		},
	}
	return matchTerm
}

// QueryTerm uses the free text ES query statement
type QueryTerm struct {
	field string
	value string
}

// Translate translates a QueryTerm into an ES "simple_query_string" statement
func (q *QueryTerm) Translate() map[string]interface{} {
	queryTerm := map[string]interface{}{
		"simple_query_string": map[string]interface{}{
			"default_field": q.field,
			"query":         q.value,
		},
	}
	return queryTerm
}
