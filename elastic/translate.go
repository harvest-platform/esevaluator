package elastic

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Leaf is a terminal node of a query
// All terminal query terms implement the Translate() method
type Leaf interface {
	Translate() map[string]interface{}
}

func encodeParam(concept string, param map[string]interface{}) map[string]interface{} {
	_ = "breakpoint"
	operator := param["operator"].(string)
	var conceptID string
	if operator == "set" {
		conceptID = "_id"
	} else {
		conceptID = param["id"].(string)
	}
	var conceptPath string
	if concept != "" && operator != "set" {
		conceptPath = concept + "." + conceptID
	} else {
		conceptPath = conceptID
	}
	value := param["value"]
	var paramLeaf Leaf
	switch operator {
	case "set":
		paramLeaf = &MemberTerm{conceptPath, value}
	case "eq":
		paramLeaf = &EqualityTerm{conceptPath, value, "+"}
	case "-eq":
		paramLeaf = &EqualityTerm{conceptPath, value, "-"}
	case "undefined":
		paramLeaf = &DefinitionTerm{conceptPath, false}
	case "defined":
		paramLeaf = &DefinitionTerm{conceptPath, true}
	case "one":
		sliceValue := value.([]string)
		paramLeaf = &OneTerm{conceptPath, sliceValue}
	case "match":
		stringValue := value.(string)
		paramLeaf = &MatchTerm{conceptPath, stringValue}
	case "query":
		stringValue := value.(string)
		paramLeaf = &QueryTerm{conceptPath, stringValue}
	case "gt":
		paramLeaf = &GreaterThanTerm{conceptPath, value, false}
	case "gte":
		paramLeaf = &GreaterThanTerm{conceptPath, value, true}
	case "lt":
		paramLeaf = &LessThanTerm{conceptPath, value, false}
	case "lte":
		paramLeaf = &LessThanTerm{conceptPath, value, true}
	case "range":
		sliceValue := value.([]interface{})
		gt := sliceValue[0]
		lt := sliceValue[1]
		paramLeaf = &RangeTerm{conceptPath, lt, gt}
	case "empty":
		paramLeaf = &EmptinessTerm{conceptPath, true}
	case "nonempty":
		paramLeaf = &EmptinessTerm{conceptPath, false}
	case "member":
		paramLeaf = &MemberTerm{conceptPath, value}
	case "subset":
		sliceValue := value.([]interface{})
		paramLeaf = &SubsetTerm{conceptPath, sliceValue}
	default:
		log.Panic(fmt.Sprintf("Operator %s not found", operator))
	}
	return paramLeaf.Translate()
}

// BooleanStatement blueprints a Must or Should term
type BooleanStatement interface {
	AddParam(map[string]interface{})
	Encode() map[string]interface{}
}

// MustTerm represents an ES boolean "must" statement
type MustTerm struct {
	params []map[string]interface{}
}

// AddParam adds a parameter to a MustTerm
func (m *MustTerm) AddParam(p map[string]interface{}) {
	m.params = append(m.params, p)
}

// Encode translates a MustTerm into an ES boolean "must" statement
func (m MustTerm) Encode() map[string]interface{} {
	boolTerm := map[string]interface{}{
		"bool": map[string]interface{}{
			"must": m.params,
		},
	}
	return boolTerm
}

// ShouldTerm represents an ES boolean "should" statement
type ShouldTerm struct {
	params []map[string]interface{}
}

// AddParam adds a parameter to a ShouldTerm
func (s *ShouldTerm) AddParam(p map[string]interface{}) {
	s.params = append(s.params, p)
}

// Encode translates a ShouldTerm into an ES boolean "should" statement
func (s ShouldTerm) Encode() map[string]interface{} {
	boolTerm := map[string]interface{}{
		"bool": map[string]interface{}{
			"should":               s.params,
			"minimum_should_match": 1,
		},
	}
	return boolTerm
}

func encodeConcept(t map[string]interface{}) map[string]interface{} {
	concept := t["concept"].(string)
	params := t["params"].([]interface{})
	var must MustTerm
	for _, param := range params {
		paramMap := param.(map[string]interface{})
		encodedParam := encodeParam(concept, paramMap)
		must.AddParam(encodedParam)
	}
	filterStatement := map[string]interface{}{
		"filter": must.Encode(),
	}
	if strings.Contains(concept, ".") {
		return Nest(filterStatement)
	}
	return filterStatement
}

func encodeBranch(t map[string]interface{}) map[string]interface{} {
	operator := t["operator"].(string)
	terms := t["terms"].([]interface{})
	var booleanStatement BooleanStatement
	if operator == "or" {
		booleanStatement = &ShouldTerm{}
	} else if operator == "and" {
		booleanStatement = &MustTerm{}
	} else {
		log.Panic(fmt.Sprintf("Unknown boolean operator: %s", operator))
	}
	for _, term := range terms {
		termMap := term.(map[string]interface{})
		termType := termMap["type"].(string)
		if termType == "concept" {
			booleanStatement.AddParam(encodeConcept(termMap))
		} else if termType == "branch" {
			booleanStatement.AddParam(encodeBranch(termMap))
		}
	}
	return booleanStatement.Encode()
}

// EncodeQuery takes a query byte array and encodes it to a valid ES query
func EncodeQuery(q []byte) map[string]interface{} {
	var queryMap map[string]interface{}
	json.Unmarshal(q, &queryMap)
	term := queryMap["term"].(map[string]interface{})
	termType := term["type"].(string)
	// Some queries may contain just one concept
	if termType == "concept" {
		encodedConcept := encodeConcept(term)
		return Filter(encodedConcept)
	} else if termType == "branch" {
		encodedBranch := encodeBranch(term)
		return Filter(encodedBranch)
	} else {
		log.Panic(fmt.Sprintf("Unknown term type: %s", termType))
	}
	return nil
}
