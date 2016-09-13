package evaluator

import (
	"encoding/json"
	"log"
)

type EmptinessTerm struct {
	field string
	empty string
}

// Translate translates an Emptiness term into a corresponsing ES "term" query
func (e *EmptinessTerm) Translate() string {
	var constantTerm map[string]interface{}
	constantTerm = make(map[string]interface{})

	var filterTerm map[string]interface{}
	filterTerm = make(map[string]interface{})

	var missingTerm map[string]interface{}
	missingTerm = make(map[string]interface{})

	missingTerm["field"] = e.field
	missingTerm["existence"] = true

	filterTerm["filter"] = missingTerm

	constantTerm["constant_score"] = filterTerm

	if e.empty == false {
		var boolTerm map[string]interface{}
		boolTerm = Negate(constantTerm)

		termSlice, err := json.Marshal(boolTern)
		if err != nil {
			log.Panic("...")
		}
		termString := string(termSlice[:len(termSlice)])
		return termString
	}

	termSlice, err := json.Marshal(constantTerm)
	if err != nil {
		log.Panic("...")
	}
	termString := string(termSlice[:len(termSlice)])
	return termString

}

type MemberTerm struct {
	field string
	value string
}

// Translate translates an Emptiness term into a corresponsing ES "terms" query
func (m *MemberTerm) Translate() string {
	var term map[string]interface{}
	term = make(map[string]interface{})

	var fieldTerm map[string]interface{}
	fieldTerm = make(map[string]interface{})

	fieldTerm[m.field] = m.value
	term["term"] = fieldTerm

	termSlice, err := json.Marshal(term)
	if err != nil {
		log.Panic("...")
	}
	termString := string(termSlice[:len(termSlice)])
	return termString
}

type SubsetTerm struct {
	field string
	set   []interface{}
}

func (s *SubsetTerm) Translate() {
	setSize = len(s.set)

	var term map[string]interface{}
	term = make(map[string]interface{})

	var fieldTerm map[string]interface{}
	fieldTerm = make(map[string]interface{})

	fieldTerm[s.field] = s.set
	fieldTerm["minimum_should_match"] = setSize

	term["terms"] = fieldTerm

	termSlice, err := json.Marshal(term)
	if err != nil {
		log.Panic("...")
	}
	termString := string(termSlice[:len(termSlice)])
	return termString

}
