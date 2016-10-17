package elastic

// Negate wraps a term in a boolean "not" condition
func Negate(t Term) Term {
	b := Term{
		"bool": Term{
			"not": t,
		},
	}
	return b
}

// Nest wraps a term in a "nested" statement
func Nest(t Term) Term {
	n := Term{
		"nested": t,
	}
	return n
}

// Query wraps a term in a "query" statement
func Query(t Term) Term {
	q := Term{
		"query": t,
	}
	return q
}

// Filter wraps a term in a "filter" statement
func Filter(t Term) Term {
	f := Term{
		"filter": t,
	}
	return f
}

// Term is an alias for map[string]interface{}
type Term map[string]interface{}

// BooleanStatement blueprints a Must or Should term
type BooleanStatement interface {
	AddParam(Term)
	Encode() Term
}

// MustTerm represents an ES boolean "must" statement
type MustTerm struct {
	params []Term
}

// AddParam adds a parameter to a MustTerm
func (m *MustTerm) AddParam(t Term) {
	m.params = append(m.params, t)
}

// Encode translates a MustTerm into an ES boolean "must" statement
func (m MustTerm) Encode() Term {
	b := Term{
		"bool": Term{
			"must": m.params,
		},
	}
	return b
}

// ShouldTerm represents an ES boolean "should" statement
type ShouldTerm struct {
	params []Term
}

// AddParam adds a parameter to a ShouldTerm
func (s *ShouldTerm) AddParam(t Term) {
	s.params = append(s.params, t)
}

// Encode translates a ShouldTerm into an ES boolean "should" statement
func (s ShouldTerm) Encode() Term {
	b := Term{
		"bool": Term{
			"should":               s.params,
			"minimum_should_match": 1,
		},
	}
	return b
}

// EqualityTerm checks if a value is equal or not equal to another value
type EqualityTerm struct {
	field  string
	value  interface{}
	negate bool
}

// Translate translates an equality term into an ES "term" term
func (e *EqualityTerm) Translate() Term {
	t := Term{
		"term": Term{
			e.field: e.value,
		},
	}
	if e.negate {
		return Negate(t)
	}
	return t
}

// DefinitionTerm checks if a value is defined for a given field
type DefinitionTerm struct {
	field   string
	defined bool
}

// Translate translates a definition term into an ES "exists" term
func (d *DefinitionTerm) Translate() Term {
	t := Term{
		"exists": Term{
			"field": d.field,
		},
	}
	if !d.defined {
		return Negate(t)
	}
	return t
}

// OneTerm checks if a field values matches one of a given set of values
type OneTerm struct {
	field  string
	values []string
}

// Translate translates a OneTerm into an ES "constant_score" statement
func (o *OneTerm) Translate() Term {
	t := Term{
		"constant_score": Term{
			"filter": Term{
				"terms": o.values,
			},
		},
	}

	return t
}

// MatchTerm uses the free text ES match statement
type MatchTerm struct {
	field string
	value string
}

// Translate translates a MatchTerm into an ES "match" statement
func (c *MatchTerm) Translate() Term {
	t := Term{
		"match": Term{
			c.field: c.value,
		},
	}
	return t
}

// QueryTerm uses the free text ES query statement
type QueryTerm struct {
	field string
	value string
}

// Translate translates a QueryTerm into an ES "simple_query_string" statement
func (q *QueryTerm) Translate() Term {
	t := Term{
		"simple_query_string": Term{
			"default_field": q.field,
			"query":         q.value,
		},
	}
	return t
}

// EmptinessTerm checks if an array field is empty
type EmptinessTerm struct {
	field string
	empty bool
}

// Translate translates an Emptiness term into a corresponsing ES "term" statement
func (e *EmptinessTerm) Translate() Term {
	t := Term{
		"constant_score": Term{
			"filter": Term{
				"field":     e.field,
				"existence": true,
			},
		},
	}
	if !e.empty {
		return Negate(t)
	}
	return t
}

// MemberTerm checks if a field value is an element of an array
type MemberTerm struct {
	field string
	value interface{}
}

// Translate translates an Emptiness term into a corresponding ES "terms" statement
func (m *MemberTerm) Translate() Term {
	t := Term{
		"term": Term{
			m.field: m.value,
		},
	}
	return t
}

// SubsetTerm checks if an array is a subset of an array field
type SubsetTerm struct {
	field string
	set   []interface{}
}

// Translate translates a SubsetTerm into a corresponding ES "terms" statement
func (s *SubsetTerm) Translate() Term {
	l := len(s.set)

	t := Term{
		"terms": Term{
			s.field:                s.set,
			"minimum_should_match": l,
		},
	}

	return t
}

var dateFormat = "dd/MM/yyyy||yyyy"

// GreaterThanTerm checks if a field value is greater than a given value
type GreaterThanTerm struct {
	field     string
	value     interface{}
	inclusive bool
}

// Translate translates a GreaterThanTerm into a corresponding range ES term
func (g *GreaterThanTerm) Translate() Term {
	t := Term{}

	var key string
	if g.inclusive {
		key = "gte"
	} else {
		key = "gt"
	}
	t[key] = g.value
	if _, ok := g.value.(string); ok {
		t["format"] = dateFormat
	}
	r := Term{
		"range": Term{
			g.field: t,
		},
	}

	return r

}

// LessThanTerm checks if a field value is less than a given value
type LessThanTerm struct {
	field     string
	value     interface{}
	inclusive bool
}

// Translate translates a LessThanTerm into a corresponding range ES term
func (l *LessThanTerm) Translate() Term {
	t := Term{}

	var key string
	if l.inclusive {
		key = "lte"
	} else {
		key = "lt"
	}
	t[key] = l.value
	if _, ok := l.value.(string); ok {
		t["format"] = dateFormat
	}
	r := Term{
		"range": Term{
			l.field: t,
		},
	}

	return r

}

// RangeTerm checks if a field value is within a given range
type RangeTerm struct {
	field string
	lt    interface{}
	gt    interface{}
}

// Translate translates a RangeTerm into a corresponding range ES term
func (rt *RangeTerm) Translate() Term {
	t := Term{}

	t["lte"] = rt.lt
	t["gte"] = rt.gt
	if _, ok := rt.lt.(string); ok {
		t["format"] = dateFormat
	}

	r := Term{
		"range": Term{
			rt.field: t,
		},
	}
	return r
}
