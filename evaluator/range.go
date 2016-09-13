package elastic

import (
	"encoding/json"
	"log"
)

var dateFormat string = "dd/MM/yyyy||yyyy"

type GreaterThanTerm struct {
	field     string
	value     interface{}
	inclusive bool
}

// Translate translates a GreaterThanTerm into a corresponding range ES term
func (g *GreaterThanTerm) Translate() string {
	var rangeTerm map[string]interface{}
	rangeTerm = make(map[string]interface{})

	var fieldTerm map[string]interface{}
	fieldTerm = make(map[string]interface{})

	var paramsTerm map[string]interface{}
	paramsTerm = make(map[string]interface{})

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

	fieldTerm[g.field] = paramsTerm
	rangeTerm["range"] = fieldTerm

	termSlice, err := json.Marshal(rangeTerm)
	if err != nil {
		log.Panic("...")
	}

	termString := string(termSlice[:len(termSlice)])
	return termString

}

type LessThanTerm struct {
	field     string
	value     interface{}
	inclusive bool
}

// Translate translates a GreaterThanTerm into a corresponding range ES term
func (l *LessThanTerm) Translate() string {
	var rangeTerm map[string]interface{}
	rangeTerm = make(map[string]interface{})

	var fieldTerm map[string]interface{}
	fieldTerm = make(map[string]interface{})

	var paramsTerm map[string]interface{}
	paramsTerm = make(map[string]interface{})

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

	fieldTerm[l.field] = paramsTerm
	rangeTerm["range"] = fieldTerm

	termSlice, err := json.Marshal(rangeTerm)
	if err != nil {
		log.Panic("...")
	}

	termString := string(termSlice[:len(termSlice)])
	return termString
}

type RangeTerm struct {
	field string
	lt    interface{}
	gt    interface{}
}

// Translate translates a GreaterThanTerm into a corresponding range ES term
func (r *RangeTerm) Translate() string {
	var rangeTerm map[string]interface{}
	rangeTerm = make(map[string]interface{})

	var fieldTerm map[string]interface{}
	fieldTerm = make(map[string]interface{})

	var paramsTerm map[string]interface{}
	paramsTerm = make(map[string]interface{})

	paramsTerm["lte"] = r.lt
	paramsTerm["gte"] = r.gt
	switch r.lt.(type) {
	case string:
		paramsTerm["format"] = dateFormat
	}

	fieldTerm[r.field] = paramsTerm
	rangeTerm["range"] = fieldTerm

	termSlice, err := json.Marshal(rangeTerm)
	if err != nil {
		log.Panic("...")
	}

	termString := string(termSlice[:len(termSlice)])
	return termString
}
