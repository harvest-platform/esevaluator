package elastic

import (
	"errors"
	"fmt"
	"strings"
)

// Leaf is a terminal node of a query
// All terminal query terms implement the Translate() method
type leaf interface {
	Translate() Term
}

// encodeParam encodes a parameter from a query to its corresponding leaf type
func encodeParam(concept string, param map[string]interface{}) (map[string]interface{}, error) {
	o := param["operator"].(string)
	var cid string
	if o == "set" {
		cid = "_id"
	} else {
		cid = param["id"].(string)
	}
	var conceptPath string
	if concept != "" && o != "set" {
		conceptPath = concept + "." + cid
	} else {
		conceptPath = cid
	}
	v := param["value"]
	var l leaf
	switch o {
	case "set":
		l = &MemberTerm{conceptPath, v}
	case "eq":
		l = &EqualityTerm{conceptPath, v, false}
	case "-eq":
		l = &EqualityTerm{conceptPath, v, true}
	case "undefined":
		l = &DefinitionTerm{conceptPath, false}
	case "defined":
		l = &DefinitionTerm{conceptPath, true}
	case "one":
		s := v.([]string)
		l = &OneTerm{conceptPath, s}
	case "match":
		s := v.(string)
		l = &MatchTerm{conceptPath, s}
	case "query":
		s := v.(string)
		l = &QueryTerm{conceptPath, s}
	case "gt":
		l = &GreaterThanTerm{conceptPath, v, false}
	case "gte":
		l = &GreaterThanTerm{conceptPath, v, true}
	case "lt":
		l = &LessThanTerm{conceptPath, v, false}
	case "lte":
		l = &LessThanTerm{conceptPath, v, true}
	case "range":
		s := v.([]interface{})
		gt := s[0]
		lt := s[1]
		l = &RangeTerm{conceptPath, lt, gt}
	case "empty":
		l = &EmptinessTerm{conceptPath, true}
	case "nonempty":
		l = &EmptinessTerm{conceptPath, false}
	case "member":
		l = &MemberTerm{conceptPath, v}
	case "subset":
		s := v.([]interface{})
		l = &SubsetTerm{conceptPath, s}
	default:
		return nil, errors.New((fmt.Sprintf("Operator %s not found", o)))
	}
	return l.Translate(), nil
}

func encodeConcept(t Term) (Term, error) {
	concept, ok := t["concept"].(string)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unable to encode concept value to string: %v", t["concept"]))
	}
	params, ok := t["params"].([]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("Unable to encode concept params to interface slice: %v", t["params"]))
	}
	var must MustTerm
	for _, p := range params {
		paramMap, ok := p.(map[string]interface{})
		if !ok {
			return nil, errors.New(fmt.Sprintf("Unable to encode param to map: %v", p))
		}
		paramTerm := Term(paramMap)
		encodedParam, err := encodeParam(concept, paramTerm)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to encode parameter: %v", paramTerm))
		}
		must.AddParam(encodedParam)
	}

	filterStatement := Term{
		"filter": must.Encode(),
	}

	if strings.Contains(concept, ".") {
		return Nest(filterStatement), nil
	}

	return filterStatement, nil
}

func encodeBranch(t Term) (Term, error) {
	op, ok := t["operator"].(string)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Cannot encode branch operator to string: %v", t["operator"]))
	}
	terms, ok := t["terms"].([]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("Cannot encode branch terms to interface slice: %v", t["terms"]))
	}
	var b BooleanStatement
	if op == "or" {
		b = &ShouldTerm{}
	} else if op == "and" {
		b = &MustTerm{}
	} else {
		return nil, errors.New(fmt.Sprintf("Unrecognized branching operator: %s", op))
	}

	for _, term := range terms {
		tmap, ok := term.(map[string]interface{})
		if !ok {
			return nil, errors.New(fmt.Sprintf("Cannot encode term to map: %v", term))
		}
		tterm := Term(tmap)
		ttype := tterm["type"].(string)

		if ttype == "concept" {
			c, err := encodeConcept(tterm)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Unable to encode concept: %v", c))
			}
			b.AddParam(c)
		} else if ttype == "branch" {
			c, err := encodeBranch(tterm)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Unable to encode branch: %v", c))
			}
			b.AddParam(c)
		}
	}

	return b.Encode(), nil
}

// Translate encodes a query into its ES equivalent
func Translate(query Term) (Term, error) {
	mapterm, ok := query["term"].(map[string]interface{})
	if !ok {
		return nil, errors.New(fmt.Sprintf("Cannot parse query term: %v", mapterm))
	}
	term := Term(mapterm)
	ttype, ok := term["type"].(string)
	if !ok {
		return nil, errors.New(fmt.Sprintf("Cannot parse term type to string, invalid json: %v", term["type"]))
	}
	switch ttype {
	case "concept":
		t, err := encodeConcept(term)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to encode concept term %v", term))
		}
		return Filter(t), nil

	case "branch":
		t, err := encodeBranch(term)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Unable to encode branch term %v", term))
		}
		return Filter(t), nil
	}

	return nil, errors.New(fmt.Sprintf("Unknown term type: %s", ttype))
}
