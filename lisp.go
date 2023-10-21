package main

import "fmt"

type LispValue interface {
}

func quote(args ...LispValue) []LispValue {
	return args
}

func length(v LispValue) int {
	switch str := v.(type) {
	case []LispValue:
		return len(str)
	}
	return 0
}

func car(v LispValue) LispValue {
	switch str := v.(type) {
	case []LispValue:
		if len(str) == 0 {
			return nil
		}
		return str[0]
	}
	return nil
}

func cdr(v LispValue) LispValue {
	switch str := v.(type) {
	case []LispValue:
		if len(str) == 0 {
			return nil
		}
		return str[1:]
	}
	return nil
}

func cadr(v LispValue) LispValue {
	if v, ok := v.([]LispValue); ok {
		if len(v) >= 2 {
			return v[1]
		}
	}
	return nil
}

func caddr(v LispValue) LispValue {
	if v, ok := v.([]LispValue); ok {
		if len(v) >= 3 {
			return v[2]
		}
	}
	return nil
}

type Condition interface {
	Error() string
}

type LispCondition struct {
	error LispValue
}

func (c LispCondition) Error() string {
	return fmt.Sprintf("%v", c.error)
}

type LispGlobalScope struct {
	sym_lookup map[string]int
	sym_names  []string
	values     []LispValue
}

func (l *LispGlobalScope) sym(name string) Symbol {
	if id, ok := l.sym_lookup[name]; ok {
		return Symbol{id: id}
	}
	i := len(l.sym_names)
	l.sym_names = append(l.sym_names, name)
	if l.sym_lookup == nil {
		l.sym_lookup = make(map[string]int)
	}
	l.sym_lookup[name] = i
	return Symbol{id: i}
}

func (l *LispGlobalScope) init_sym(name string, value LispValue) Symbol {
	sym := l.sym(name)
	l.set_value(sym, value)
	return sym
}

func (l *LispGlobalScope) get_value(name Symbol) LispValue {
	if len(l.values) <= name.id {
		return nil
	}
	return l.values[name.id]
}

func (l *LispGlobalScope) set_value(name Symbol, value LispValue) {
	if len(l.values) <= name.id {
		prev := l.values
		l.values = make([]LispValue, max(128, name.id*3/2))
		if prev != nil {
			copy(l.values, prev)
		}
	}

	l.values[name.id] = value
}

var lisp LispGlobalScope = LispGlobalScope{}

var plus_sym = lisp.init_sym("+", add)
var loop_sym = lisp.init_sym("loop", eval_loop)
var set_sym = lisp.init_sym("set", eval_set)
var gt_sym = lisp.init_sym(">", _gt)
var println_sym = lisp.init_sym("println", lisp_println)
var raise_sym = lisp.init_sym("raise", lisp_raise)
var quote_sym = lisp.init_sym("quote", quote_macro)
var error_handler_sym = lisp.init_sym("error-handler", lisp_error_handler)
var let_sym = lisp.init_sym("let", lisp_let)
var eval_sym = lisp.init_sym("eval", eval)

func quote_macro(scope *LispScope, code LispValue) LispValue {
	return cadr(code)
}

func lisp_println(values []LispValue) {
	for _, v := range values {
		fmt.Printf("%v", v)
	}
	println("")
}

func lisp_let(scope *LispScope, code LispValue) LispValue {
	arr, ok := code.([]LispValue)
	if !ok || len(arr) <= 2 {
		return lisp_raise("...")
	}

	args, ok := arr[1].([]LispValue)
	if !ok {
		return lisp_raise("...")
	}
	names := make([]int, len(args))
	values := make([]LispValue, len(args))
	for i := range args {
		arg := args[i]
		arg2, ok := arg.([]LispValue)
		if !ok {
			return lisp_raise("...")
		}
		sym, ok := arg2[0].(Symbol)
		if !ok {
			return lisp_raise("...")
		}
		value := arg2[1]
		names[i] = sym.id
		values[i] = value
	}
	scope2 := LispScope{
		parentScope: scope,
		values:      values,
		valuesName:  names}

	return eval_progn(&scope2, arr[2:])

}

func lisp_raise(errorCode LispValue) LispValue {
	return LispCondition{error: errorCode}
}

func lisp_error_handler(scope *LispScope, code LispValue) LispValue {
	arr, ok := code.([]LispValue)
	if !ok {
		return lisp_raise("invalid syntax for error-handler")
	}
	code1 := arr[0]
	code2 := arr[1]
	result := eval(scope, code1)
	if _, ok := result.(LispCondition); ok {
		handler_arg, ok := car(car(code2)).(Symbol)
		if !ok {
			return lisp_raise("error-handler must have a symbol as second argument.")
		}
		subScope := LispScope{parentScope: scope, values: []LispValue{result}, valuesName: []int{handler_arg.id}}
		return eval(&subScope, cadr(code2))
	}
	return result
}

func _gt(a LispValue, b LispValue) LispValue {

	a2, ok1 := a.(int)
	b2, ok2 := b.(int)
	if ok1 && ok2 {
		if a2 > b2 {
			return true
		}
	}
	return nil
}

type LispScope struct {
	parentScope *LispScope
	values      []LispValue
	valuesName  []int
}

type Symbol struct {
	id int
}

func (l *LispScope) set_value(name Symbol, value LispValue) {
	if l == nil {
		lisp.set_value(name, value)
		return
	}
	for i := range l.values {
		if l.valuesName[i] == name.id {
			l.values[i] = value
			return
		}
	}

	l.parentScope.set_value(name, value)
}

func (l *LispScope) get_value(name Symbol) LispValue {
	if l == nil {
		return lisp.get_value(name)
	}
	for i := range l.values {
		if l.valuesName[i] == name.id {
			return l.values[i]
		}
	}

	return l.parentScope.get_value(name)
}

func add(a LispValue, b LispValue) LispValue {
	switch v := a.(type) {
	case int:
		switch v2 := b.(type) {
		case int:
			return v + v2
		case float64:
			return float64(v) + v2
		}
	case float64:
		switch v2 := b.(type) {
		case int:
			return v + float64(v2)
		case float64:
			return v + v2
		}
	}
	return nil
}

func eval_set(scope *LispScope, code LispValue) LispValue {

	var code2, ok = (code).([]LispValue)
	if !ok || len(code2) != 3 {
		return nil
	}
	sym, ok := code2[1].(Symbol)
	if !ok {
		lisp_raise("value must be a ...")
		return nil
	}
	var value = eval(scope, code2[2])
	scope.set_value(sym, value)
	return nil
}

func eval_progn(scope *LispScope, code LispValue) LispValue {
	var result LispValue = nil

	simple, isSlice := code.([]LispValue)
	if isSlice {
		for i := range simple {
			result = eval(scope, simple[i])
		}
	}

	return result
}

func eval_loop(scope *LispScope, code LispValue) LispValue {
	condition := cadr(code)
	body := cdr(cdr(code))
	var result LispValue = nil
	for {
		r := eval(scope, condition)
		if r == nil {
			break
		}
		result = eval_progn(scope, body)
	}
	return result
}

func eval(scope *LispScope, code LispValue) LispValue {

	if arr, ok := code.([]LispValue); ok {
		if len(arr) == 0 {
			return LispCondition{error: "err"}
		}

		first := arr[0]
		if s, ok := first.(Symbol); ok {
			f := scope.get_value(s)
			switch f2 := f.(type) {
			case func(*LispScope, LispValue) LispValue:
				return f2(scope, code)
			case func(LispValue) LispValue:
				arg := eval(scope, arr[1])
				if _, ok := arg.(LispCondition); ok {
					return arg
				}
				return f2(arg)
			case func(LispValue, LispValue) LispValue:
				a := eval(scope, arr[1])
				b := eval(scope, arr[2])
				if _, ok := a.(LispCondition); ok {
					return a
				}
				if _, ok := b.(LispCondition); ok {
					return b
				}
				return f2(a, b)
			case func(LispValue, LispValue, LispValue) LispValue:
				a := eval(scope, arr[1])
				b := eval(scope, arr[2])
				c := eval(scope, arr[3])
				if _, ok := a.(LispCondition); ok {
					return a
				}
				if _, ok := b.(LispCondition); ok {
					return b
				}
				if _, ok := c.(LispCondition); ok {
					return c
				}
				return f2(a, b, c)
			case func([]LispValue) LispValue:
				args := make([]LispValue, len(arr)-1)
				for i := range args {
					args[i] = eval(scope, arr[1+i])
					if _, ok := args[i].(LispCondition); ok {
						return args[i]
					}
				}
				return f2(args)
			case func([]LispValue):
				args := make([]LispValue, len(arr)-1)
				for i := range args {
					args[i] = eval(scope, arr[1+i])
					if _, ok := args[i].(LispCondition); ok {
						return args[i]
					}
				}
				f2(args)
				return nil

			case func(*LispScope, LispValue):
				{
					f2(scope, code)
					return nil
				}
			}
		}
		return lisp_raise("Invalid code")
	}
	if s, ok := code.(Symbol); ok {
		return scope.get_value(s)
	}

	return code
}
