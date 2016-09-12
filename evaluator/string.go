package elastic

import (
	"encoding/json"
	"log"
)

type OneTerm struct {
	field  string
	values []string
}

func (o *OneTerm) Translate() string {
	var constantTerm map[string]interface{}
	constantTerm = make(map[string]interface{})

	var filterTerm map[string]interface{}
	filterTerm = make(map[string]interface{})

	var term map[string]interface{}
	term = make(map[string]interface{})

	term["terms"] = o.values
	filterTerm["filter"] = term
	constantTerm["constant_score"] = filterTerm

	termSlice, err := json.Marshal(constantTerm)
	if err != nil {
		log.Panic("...")
	}

	termString := string(termSlice[:len(termSlice)])
	return termString
}

type MatchTerm struct {
	field string
	value string
}

func (c *MatchTerm) Translate() string {
	var matchTerm map[string]interface{}
	matchTerm = make(map[string]interface{})

	var messageTerm map[string]interface{}
	messageTerm = make(map[string]interface{})

	messageTerm[c.field] = c.value
	matchTerm["match"] = messageTerm

	termSlice, err := json.Marshal(matchTerm)
	if err != nil {
		log.Panic("...")
	}

	termString := string(termSlice[:len(termSlice)])
	return termString
}

type QueryTerm struct {
	field string
	value string
}

func (q *QueryTerm) Translate() string {
	var queryTerm map[string]interface{}
	queryTerm = make(map[string]interface{})

	var paramTerm map[string]string
	paramTerm = make(map[string]string)

	paramTerm["default_field"] = q.field
	paramTerm["query"] = q.value

	queryTerm["simple_query_string"] = paramTerm

	termSlice, err := json.Marshal(queryTerm)
	if err != nil {
		log.Panic("...")
	}

	termString := string(termSlice[:len(termSlice)])
	return termString

}
