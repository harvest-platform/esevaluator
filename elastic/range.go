package elastic

var dateFormat = "dd/MM/yyyy||yyyy"

// GreaterThanTerm checks if a field value is greater than a given value
type GreaterThanTerm struct {
	field     string
	value     interface{}
	inclusive bool
}

// Translate translates a GreaterThanTerm into a corresponding range ES term
func (g *GreaterThanTerm) Translate() map[string]interface{} {
	paramsTerm := make(map[string]interface{})

	var inclusionKey string
	switch g.inclusive {
	case true:
		inclusionKey = "gte"
	case false:
		inclusionKey = "gt"
	}

	switch g.value.(type) {
	case string:
		paramsTerm[inclusionKey] = g.value
		paramsTerm["format"] = dateFormat
	default:
		paramsTerm[inclusionKey] = g.value
	}

	rangeTerm := map[string]interface{}{
		"range": map[string]interface{}{
			g.field: paramsTerm,
		},
	}

	return rangeTerm

}

// LessThanTerm checks if a field value is less than a given value
type LessThanTerm struct {
	field     string
	value     interface{}
	inclusive bool
}

// Translate translates a LessThanTerm into a corresponding range ES term
func (l *LessThanTerm) Translate() map[string]interface{} {
	paramsTerm := make(map[string]interface{})

	var inclusionKey string
	switch l.inclusive {
	case true:
		inclusionKey = "lte"
	case false:
		inclusionKey = "lt"
	}

	switch l.value.(type) {
	case string:
		paramsTerm[inclusionKey] = l.value
		paramsTerm["format"] = dateFormat
	default:
		paramsTerm[inclusionKey] = l.value
	}

	rangeTerm := map[string]interface{}{
		"range": map[string]interface{}{
			l.field: paramsTerm,
		},
	}

	return rangeTerm

}

// RangeTerm checks if a field value is within a given range
type RangeTerm struct {
	field string
	lt    interface{}
	gt    interface{}
}

// Translate translates a RangeTerm into a corresponding range ES term
func (r *RangeTerm) Translate() map[string]interface{} {
	paramsTerm := make(map[string]interface{})

	paramsTerm["lte"] = r.lt
	paramsTerm["gte"] = r.gt
	switch r.lt.(type) {
	case string:
		paramsTerm["format"] = dateFormat
	}

	rangeTerm := map[string]interface{}{
		"range": map[string]interface{}{
			r.field: paramsTerm,
		},
	}
	return rangeTerm
}
