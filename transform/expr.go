package transform

import (
	"errors"
	"github.com/qntfy/kazaam/registry"
	"go/ast"
	"go/constant"
	"go/parser"
	"go/token"
	"strconv"
)

// support for basic expression evaluation support.
// NOTE: expressions must evaluate to a bool value, i.e. true or false
// supports
//  ! - unary not operator
// && - binary logical and
// || - binary logical or
// == - binary equality
// != - binary equality negated
// <  - binary less than
// >  - binary greater than
// <= - binary less than or equal to
// >= - binary greater than or equal to
// ( ) - parentheses grouping
// true - boolean true constant
// false - boolean false constant
// <number constants> - number constants
// <string constants> - string constants  ( "" )
// <func calls> -  converter function call,    f( jsonPath, stringArgs ) where, f is a converter id, jsonPath is a json path, and args are arguments to the converter
// <json path resolution> - a json path.. does not support path parameters, i.e. ? conditionals and converter extensions
type BasicExpr struct {
	fullExpression string
	data           []byte
	astExpr        ast.Expr
}

// json data for evaluating json paths, and the expression string to evaluate to a boolean true/false constant value
// err will be returned if the expression can't be parsed
func NewBasicExpr(data []byte, exprStr string) (expr *BasicExpr, err error) {
	astExpr, err := parser.ParseExpr(exprStr)
	if err != nil {
		return
	}

	expr = new(BasicExpr)

	expr.data = data
	expr.fullExpression = exprStr
	expr.astExpr = astExpr

	return
}

// returns true|false, otherwise error if the expression can not be evaluated
func (expr *BasicExpr) Eval() (val bool, err error) {
	var evaluation constant.Value
	evaluation, err = expr.evalExpr(expr.astExpr)
	if err != nil {
		return
	}
	val = constant.BoolVal(evaluation)
	return
}

func (expr *BasicExpr) evalBinaryExpr(exp *ast.BinaryExpr) (val constant.Value, err error) {
	var left, right constant.Value

	left, err = expr.evalExpr(exp.X)
	if err != nil {
		return
	}

	right, err = expr.evalExpr(exp.Y)
	if err != nil {
		return
	}

	switch exp.Op {
	// logical operators
	case token.LOR:
		fallthrough
	case token.LAND:

		l := left.Kind()
		r := right.Kind()

		if l != constant.Bool || r != constant.Bool {
			err = errors.New("logical operators require bool values")
			return
		}

		val = constant.BinaryOp(left, exp.Op, right)
		return

		// comparison operators
	case token.GTR:
		fallthrough
	case token.LSS:
		fallthrough
	case token.EQL:
		fallthrough
	case token.NEQ:
		fallthrough
	case token.GEQ:
		fallthrough
	case token.LEQ:

		l := left.Kind()
		r := right.Kind()

		if l != r && !(expr.isNumKind(l) && expr.isNumKind(r)) {
			err = errors.New("comparison operators require types to be the same")
			return
		}

		val = constant.MakeBool(constant.Compare(left, exp.Op, right))
		return
	}

	err = errors.New("unsupported operator")

	return
}

func (expr *BasicExpr) isNumKind(k constant.Kind) bool {
	return k == constant.Int || k == constant.Float
}

func (expr *BasicExpr) evalExpr(exp ast.Expr) (val constant.Value, err error) {
	switch exp := exp.(type) {
	case *ast.UnaryExpr:
		if exp.Op == token.NOT {
			val, err = expr.evalExpr(exp.X)
			if err != nil {
				return
			}
			val = constant.MakeBool(!(constant.BoolVal(val)))
			return
		}
	case *ast.ParenExpr:
		val, err = expr.evalExpr(exp.X)
		return
	case *ast.BinaryExpr:
		val, err = expr.evalBinaryExpr(exp)
		return
	case *ast.Ident:
		if exp.Name == "true" || exp.Name == "false" {
			val = constant.MakeBool(exp.Name == "true")
			return
		} else if exp.Name == "null" || exp.Name == "nil" {
			val =  constant.MakeUnknown()
			return
		} else {
			// assumed to be a json path variable -- without selector syntax, e.g top level prop
			val, err = expr.evalJsonPath(exp.Name)
			return
		}
	case *ast.SelectorExpr:
		pos := exp.Pos()
		end := exp.End()
		path := expr.fullExpression[ pos-1 : end-1 ]
		val, err = expr.evalJsonPath(path)
		return
	case *ast.BasicLit:
		switch exp.Kind {
		case token.STRING:
			val = constant.MakeFromLiteral(exp.Value, exp.Kind, 0)
			return
		case token.INT:
			val = constant.MakeFromLiteral(exp.Value, exp.Kind, 0)
			return
		case token.FLOAT:
			val = constant.MakeFromLiteral(exp.Value, exp.Kind, 0)
			return
		}
	case *ast.CallExpr:
		val, err = expr.evalConverterFunc(exp.Fun.(*ast.Ident).Name, exp.Args)
		return
	}

	err = errors.New("unsupported expression syntax")

	return
}

func (expr *BasicExpr) evalJsonPath(path string) (val constant.Value, err error) {
	var jsonPathSimpleValue *JSONValue

	jsonPathSimpleValue, err = GetJsonPathValue(expr.data, path)
	if err != nil {
		return
	}

	switch jsonPathSimpleValue.valueType {
	case JSONNull:
		val = constant.MakeUnknown()
	case JSONString:
		val = constant.MakeString(jsonPathSimpleValue.GetStringValue())
	case JSONInt:
		val = constant.MakeInt64(jsonPathSimpleValue.GetIntValue())
	case JSONFloat:
		val = constant.MakeFloat64(jsonPathSimpleValue.GetFloatValue())
	case JSONBool:
		val = constant.MakeBool(jsonPathSimpleValue.GetBoolValue())
	default:
		err = errors.New("unsupported type")
	}

	return
}

func (expr *BasicExpr) evalConverterFunc(name string, args []ast.Expr) (val constant.Value, err error) {
	var jsonPath, converterArgs string

	if len(args) > 0 {
		argValues := make([]constant.Value, 0, len(args))
		for _, a := range args {
			v, e := expr.evalExpr(a)
			if e != nil {
				err = e
				return
			}
			argValues = append(argValues, v)
		}
		if argValues[0].Kind() == constant.String {
			jsonPath = constant.StringVal(argValues[0])
		} else {
			err = errors.New("expected json path as string")
			return
		}
		if len(args) > 1 {
			if argValues[1].Kind() == constant.String {
				converterArgs = constant.StringVal(argValues[1])
			} else {
				err = errors.New("expected converter arguments as string")
				return
			}
		}
	} else {
		err = errors.New("expected path string, and optional arguments as string")
		return
	}

	conv := registry.GetConverter(name)
	if conv != nil {

		var jsonPathBytes, converterArgsBytes []byte

		// get json path (string)'s value as simple value
		jsonPathValue, e := GetJsonPathValue(expr.data, jsonPath)
		if e != nil {
			err = e
			return
		}

		// get the bytes for this value
		jsonPathBytes = jsonPathValue.GetData()

		if len(converterArgs) > 0 {
			converterArgsBytes = []byte(strconv.Quote(converterArgs))
		}

		// run them through the converter, and get the new value's bytes
		var newValue []byte
		newValue, err = conv.Convert(expr.data, jsonPathBytes, converterArgsBytes)
		if err != nil {
			return
		}

		// convert to simple value
		var simpleValue *JSONValue
		simpleValue, err = NewJSONValue(newValue)
		if err != nil {
			return
		}

		switch simpleValue.valueType {
		case JSONNull:
			val = constant.MakeUnknown()
		case JSONString:
			val = constant.MakeString(simpleValue.GetStringValue())
		case JSONInt:
			val = constant.MakeInt64(simpleValue.GetIntValue())
		case JSONFloat:
			val = constant.MakeFloat64(simpleValue.GetFloatValue())
		case JSONBool:
			val = constant.MakeBool(simpleValue.GetBoolValue())
		default:
			err = errors.New("unsupported type")
		}

	} else {
		val = constant.MakeUnknown()
		err = errors.New("missing converter")
	}

	return
}
