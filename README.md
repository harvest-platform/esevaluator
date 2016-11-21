# ES-Evaluator
This project implements a Harvest-Next query evaluator targeting Elasticsearch.

The evaluator acts as a validator for Harvest-Next queries and as an interpreter targeting Elasticsearch 2+.  Much like Harvest 2.x, the Harvest-Next query language is json-based and composed of concepts and individual fields.  A query will consist of at least one concept term, which is composed of a set of field parameters.  A field parameter is an atomic term which contains an identifier for the field (if applicable), an operator, and a value.  

A query may also contain one or more branch terms, which allow for boolean logic in concept querying, consist of concept terms and/or other branch terms.

### Example
Input query:
```json
{
  "type": "query",
  "term": {
    "type": "branch",
    "operator": "and",
    "terms": [
      {
        "type": "concept",
        "concept": "audiogram.puretone.f500",
        "params": [
          {
            "id": "response",
            "operator": "range",
            "value": [20,30]
          },
          {
            "operator": "set",
            "value": [1,2,3,4,5]
          },
          {
            "id": "masking_used",
            "operator": "-eq",
            "value": false
          },
          {
            "id": "no_response",
            "operator": "eq",
            "value": false
          },
          {
            "id": "vibrotactial_response",
            "operator": "eq",
            "value": false
          }
        ]
      },
      {
        "type": "branch",
        "operator": "or",
        "terms": [
          {
            "type": "concept",
            "concept": "audiogram.puretone.pta",
            "params": [
              {
                "id": "ear",
                "operator": "eq",
                "value": "Left"
              },
              {
                "id": "value",
                "operator": "range",
                "value": [30,40]
              }
            ]
          },
          {
            "type": "concept",
            "concept": "audiogram.puretone.pta",
            "params": [
              {
                "id": "ear",
                "operator": "eq",
                "value": "Right"
              },
              {
                "id": "value",
                "operator": "gt",
                "value": 30
              }
            ]
          }
        ]
      }
    ]
  }
}
```

Output Elasticsearch query:
```json
{
  "filter": {
    "bool": {
      "must": [{
        "nested": {
          "filter": {
            "bool": {
              "must": [{
                "range": {
                  "audiogram.puretone.f500.response": {
                    "gte": 20,
                    "lte": 30
                  }
                }
              }, {
                "term": {
                  "_id": [1, 2, 3, 4, 5]
                }
              }, {
                "bool": {
                  "not": {
                    "term": {
                      "audiogram.puretone.f500.masking_used": false
                    }
                  }
                }
              }, {
                "term": {
                  "audiogram.puretone.f500.no_response": false
                }
              }, {
                "term": {
                  "audiogram.puretone.f500.vibrotactial_response": false
                }
              }]
            }
          }
        }
      }, {
        "bool": {
          "minimum_should_match": 1,
          "should": [{
            "nested": {
              "filter": {
                "bool": {
                  "must": [{
                    "term": {
                      "audiogram.puretone.pta.ear": "Left"
                    }
                  }, {
                    "range": {
                      "audiogram.puretone.pta.value": {
                        "gte": 30,
                        "lte": 40
                      }
                    }
                  }]
                }
              }
            }
          }, {
            "nested": {
              "filter": {
                "bool": {
                  "must": [{
                    "term": {
                      "audiogram.puretone.pta.ear": "Right"
                    }
                  }, {
                    "range": {
                      "audiogram.puretone.pta.value": {
                        "gt": 30
                      }
                    }
                  }]
                }
              }
            }
          }]
        }
      }]
    }
  },
  "source": "_id"
}
```




## Supported Field Terms/Operators:

#### set
Returns records whose identifiers match a set of values.
```json
{
  "operator": "set",
  "values": [1,2,3,4,5]
}
```

#### eq
Equality term.  Accepts numerical, boolean, and string values.
```json
{
  "id": "response",
  "operator": "eq",
  "value": 24
}
```

#### -eq
Non-equality term.  Accepts numerical, boolean, and string values.
```json
{
  "id": "response",
  "operator": "-eq",
  "value": 24
}
```

#### defined
Term that checks whether a value for a given field is defined.
```json
{
  "id": "race",
  "operator": "defined"
}
```

#### undefined
Term that checks whether a value for a given field is undefined.
```json
{
  "id": "race",
  "operator": "undefined"
}
```

#### one
Term that, given a list of values, checks if the value of a field matches at least one of the values.  Accepts numerical, boolean, and string values.
```json
{
  "id": "race",
  "operator": "one",
  "values": ["white", "asian"]
}
```

#### match
Implements a ES free-text "match" on a text or string field.
```json
{
  "id": "note",
  "operator": "match",
  "value": "cancer oncology"
}
```

#### query
Implements a free-text ES query search of a text or string field.
```json
{
  "id": "note",
  "operator": "query",
  "value": "+(cancer | oncology) neoblastoma"
}
```

#### gt
Queries for values greater (exclusive) than the supplied value.  Accepts numeric and date values with a "dd/MM/yyyy||yyyy" format.
```json
{
  "id": "dob",
  "operator": "gt",
  "value": "1998"
}
```
#### gte
Queries for values greater (inclusive) than the supplied value.  Accepts numeric and date values with a "dd/MM/yyyy||yyyy" format.
```json
{
  "id": "dob",
  "operator": "gte",
  "value": "1998"
}
```

#### lt
Queries for values less (exclusive) than the supplied value.  Accepts numeric and date values with a "dd/MM/yyyy||yyyy" format.
```json
{
  "id": "dob",
  "operator": "lt",
  "value": "1998"
}
```

#### lte
Queries for values less (inclusive) than the supplied value.  Accepts numeric and date values with a "dd/MM/yyyy||yyyy" format.
```json
{
  "id": "dob",
  "operator": "lte",
  "value": "1998"
}
```

#### range
Queries for values (inclusive) in a given range.  Accepts numeric and date values with a "dd/MM/yyyy||yyyy" format.
```json
{
  "id": "dob",
  "operator": "range",
  "value": ["1998", "2003"]
}
```

#### empty
Checks whether a given array field is has an empty value.
```json
{
  "id": "dob",
  "operator": "empty"
}
```

#### nonempty
Checks whether a given array field is has an non-empty value.
```json
{
  "id": "dob",
  "operator": "nonempty"
}
```

#### member
Checks whether a value is a member of an array field. Accepts numerical, boolean, and string values.
```json
{
  "id": "transducers",
  "operator": "member",
  "value": "ER-3A Inserts"
}
```

#### subset
Checks whether an array is a subset of an array field. Accepts numerical, boolean, and string values.
```json
{
  "id": "transducers",
  "operator": "subset",
  "value": ["ER-3A Inserts", "TDH-39 Phones"]
}
```
