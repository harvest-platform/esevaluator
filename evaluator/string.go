package elastic

type OneTerm struct {
	field  string
	values []string
}

func (o *OneTerm) Translate() string {
	constantTerm := map[string]interface{}{
		"constant_score": map[string]interface{}{
			"filter": map[string]interface{}{
				"terms": o.values,
			},
		},
	}

	return constantTerm
}

type MatchTerm struct {
	field string
	value string
}

func (c *MatchTerm) Translate() string {
	matchTerm := map[string]interface{}{
		"match": map[string]interface{}{
			c.field: c.value,
		},
	}
	return matchTerm
}

type QueryTerm struct {
	field string
	value string
}

func (q *QueryTerm) Translate() string {
	queryTerm := map[string]interface{}{
		"simple_query_string": map[string]interface{}{
			"default_field": q.field,
			"query":         q.value,
		},
	}
	return queryTerm
}
