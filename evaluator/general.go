package evaluator

import (
	"encoding/json"
	"log"
)

type EqualityTerm struct {
	field string
	value interface{}
	sign  string
}

// Translate translates an equality term into an ES "term" term
func (e *EqualityTerm) Translate() string {
	if e.sign == "-" {
		var fieldExpression map[string]interface{}
		fieldExpression = make(map[string]interface{})
		fieldExpression[e.field] = e.value

		var term map[string]interface{}
		term = make(map[string]interface{})

		term["term"] = fieldExpression

		var boolTerm map[string]interface{}
		boolTerm = Negate(term)

		termSlice, err := json.Marshal(boolTerm)
		if err != nil {
			log.Panic("...")
		}
		termString := string(termSlice[:len(termSlice)])
		return termString
	} else {

		var fieldExpression map[string]interface{}
		fieldExpression = make(map[string]interface{})
		fieldExpression[e.field] = e.value

		var term map[string]interface{}
		term = make(map[string]interface{})

		termSlice, err := json.Marshal(term)
		if err != nil {
			log.Panic("...")
		}
		termString := string(termSlice[:len(termSlice)])
		return termString
	}
}

type DefinitionTerm struct {
	field   string
	defined bool
}

// Translate translates a definition term into an ES "exists" term
func (d *DefinitionTerm) Translate() string {
	if d.defined == true {
		var fieldTerm map[string]interface{}
		fieldTerm = make(map[string]interface{})
		fieldTerm["field"] = d.field

		var existsTerm map[string]interface{}
		existsTerm = make(map[string]interface{})
		existsTerm["exists"] = fieldTerm

		termSlice, err := json.Marshal(existsTerm)
		if err != nil {
			log.Panic("...")
		}
		termString := string(termSlice[:len(termSlice)])
		return termString
	} else {
		var fieldTerm map[string]interface{}
		fieldTerm = make(map[string]interface{})
		fieldTerm["field"] = d.field

		var existsTerm map[string]interface{}
		existsTerm = make(map[string]interface{})
		existsTerm["exists"] = fieldTerm

		var boolTerm map[string]interface{}
		boolTerm = Negate(existsTerm)

		termSlice, err := json.Marshal(boolTerm)
		if err != nil {
			log.Panic("...")
		}
		termString := string(termSlice[:len(termSlice)])
		return termString
	}
}
