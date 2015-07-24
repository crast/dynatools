/*
package safeupdate generates safe UpdateItem query expressions from dynago.

PutItem is available and makes putting a record with a dynago.Document easy.
In theory you can use PutItem for updates if you're doing a full update of a
record, but this can be dangerous if you consider the case of having
potentially multiple versions of an application deployed simultaneously, some
of which know about more fields than the other.

UpdateItem is available, but using it for a full update of a record with a
number of attributes is really boilerplatey, especially dealing with reserved
words in expression attributes.

So the SafeUpdate package will take a Document and make an UpdateItem with the
appropriate reserved words handled as expression attribute names, and the
appropriate value expression, compressed.
*/
package safeupdate

import (
	"fmt"
	"sort"
	"strings"

	"gopkg.in/underarmour/dynago.v1"
)

type Update struct {
	Expression string          // Generated update expression
	Key        dynago.Document // Key extracted from attributes
	EANames    Params          // renamed reserved words
	EAValues   Params          // expression attibute values
}

// Apply the attributes of this update, returning an UpdateItem ready to execute.
func (u Update) Apply(client *dynago.Client, tableName string) *dynago.UpdateItem {
	return client.UpdateItem(tableName, u.Key).
		UpdateExpression(u.Expression, u.EANames, u.EAValues)
}

/*
Build a descriptor for a safe UpdateItem query on the specified document.

keyNames should be a slice with either one or two elements depending
on whether the document has a hash key or a hash-range key.

doc should be a document more or less the same as you would provide to
something like PutItem.

Given a document with keys like Id,Name,Age,Address it generates an expression
something like:
	SET Address=:a, Age=:b, #c=:c
(Name is a reserved word in dynamo so it's aliased as #c) and the appropriate
parameter set for the expression attribute names and values.

As a special case, due to special considerations in DynamoDB, if a value
is set to the empty string or an empty set, then the expression will have
REMOVE expressions for those properties (Dynamo can't store empty strings
nor can it store empty sets)
*/
func Build(keyNames []string, doc dynago.Document) *Update {
	var key = [2]string{keyNames[0], ""}
	if len(keyNames) == 2 {
		key[1] = keyNames[1]
	}
	schema, derivedKey := prepare(key, doc)
	u := &Update{
		Expression: schema.expression,
		Key:        derivedKey,
		EANames:    schema.eaNames,
		EAValues:   make(Params, len(schema.mapped)),
	}

	for i, attr := range schema.mapped {
		u.EAValues[i] = dynago.Param{attr.paramName, doc[attr.key]}
	}
	return u
}

// prepare splits out the keys from the document and returns the appropriate schema and keys.
func prepare(keyNames [2]string, doc dynago.Document) (*schema, dynago.Document) {
	derivedKey := make(dynago.Document, 2)
	var exprAttributes = make([]string, 0, len(doc))
	var emptyAttributes []string
	for k, v := range doc {
		if k == keyNames[0] || k == keyNames[1] {
			derivedKey[k] = v
		} else if isEmpty(v) {
			emptyAttributes = append(emptyAttributes, k)
		} else {
			exprAttributes = append(exprAttributes, k)
		}
	}
	// We sort the attributes for idempotency, but in the future we might also use this
	// feature so that we can cache our schema that we build.
	sort.Strings(exprAttributes)
	sort.Strings(emptyAttributes)
	return buildSchema(keyNames, exprAttributes, emptyAttributes), derivedKey
}

func buildSchema(keyNames [2]string, exprAttributes, emptyAttributes []string) *schema {
	nameGen := variableGen{[]rune{'#'}}
	valueGen := variableGen{[]rune{':'}}
	mapped := make([]attributeDesc, len(exprAttributes))
	var eaNames Params
	makeName := func(s string) (eaName string) {
		if reservedWords[strings.ToUpper(s)] {
			eaName = nameGen.next()
			eaNames = append(eaNames, dynago.Param{eaName, s})
		}
		return
	}
	var toSet []string
	for i, s := range exprAttributes {
		d := attributeDesc{
			key:       s,
			paramName: valueGen.next(),
		}
		setKey := s
		if eaName := makeName(s); eaName != "" {
			d.eaName = eaName
			setKey = eaName
		} else {
			nameGen.next() // just provides some handy symmetry.
		}
		toSet = append(toSet, fmt.Sprintf("%s=%s", setKey, d.paramName))
		mapped[i] = d
	}
	expression := fmt.Sprintf("SET %s", strings.Join(toSet, ", "))

	if len(emptyAttributes) > 0 {
		for i, s := range emptyAttributes {
			name := makeName(s)
			if name != "" {
				emptyAttributes[i] = name
			}
		}
		expression += fmt.Sprintf(" REMOVE %s", strings.Join(emptyAttributes, ", "))
	}
	return &schema{
		expression: expression,
		key:        keyNames,
		mapped:     mapped,
		eaNames:    eaNames,
	}
}

type schema struct {
	expression string
	key        [2]string
	mapped     []attributeDesc
	eaNames    Params
}

type attributeDesc struct {
	key       string
	eaName    string
	paramName string
}

type variableGen struct {
	current []rune
}

// This relies on the ascii a-z being a core part of runes
func (v *variableGen) next() string {
	if len(v.current) == 1 {
		v.current = append(v.current, 'a')
	} else {
		v.current = pushUp(v.current, len(v.current)-1)
	}
	return string(v.current)
}

func pushUp(b []rune, lastIndex int) []rune {
	if b[lastIndex] == 'z' {
		b[lastIndex] = 'a'
		if lastIndex == 1 {
			return append(b, 'a')
		} else {
			return pushUp(b, lastIndex-1)
		}
	} else {
		b[lastIndex] = b[lastIndex] + 1
	}
	return b
}

// Implement the dynago Params interface to use without double converting
type Params []dynago.Param

func (p Params) AsParams() []dynago.Param {
	return p
}

func isEmpty(val interface{}) bool {
	// Simplistic version of the empty check
	switch val := val.(type) {
	case string:
		return val == ""
	case dynago.StringSet:
		return len(val) == 0
	case dynago.NumberSet:
		return len(val) == 0
	case dynago.BinarySet:
		return len(val) == 0
	default:
		return false
	}
}
