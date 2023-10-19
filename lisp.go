package main

import "fmt"

type LispValue interface {
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
		return nil
	}
	return car(cdr(v))
}

func caddr(v LispValue) LispValue {
	if v, ok := v.([]LispValue); ok {
		if len(v) >= 3 {
			return v[2]
		}
		return nil
	}
	return car(cdr(cdr(v)))
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
var define_sym = lisp.sym("define")
var gt_sym = lisp.init_sym(">", _gt)
var println_sym = lisp.init_sym("println", lisp_println)
var quote_sym = lisp.init_sym("quote", quote_macro)

func quote_macro(scope *LispScope, code LispValue) LispValue {
	return cadr(code)
}

func lisp_println(values []LispValue) {
	for _, v := range values {
		fmt.Printf("%v", v)
	}
	println("")
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
	values      map[int]LispValue
}

func (l *LispScope) set_value(name Symbol, value LispValue) {
	if l == nil {
		lisp.set_value(name, value)
		return
	}
	if _, ok := l.values[name.id]; ok {
		l.values[name.id] = value
	}
	l.parentScope.set_value(name, value)
}

func (l *LispScope) get_value(name Symbol) LispValue {
	if l == nil {
		return lisp.get_value(name)
	}
	if v, ok := l.values[name.id]; ok {
		return v
	}
	return l.parentScope.get_value(name)
}

type Symbol struct {
	id int
}

type Number interface {
	~int
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
	var sym, ok = cadr(code).(Symbol)
	if !ok {
		return nil
	}
	var value = eval(scope, caddr(code))
	scope.set_value(sym, value)
	return nil
}

func eval_progn(scope *LispScope, code LispValue) LispValue {
	var result LispValue = nil

	simple, isSlice := code.([]LispValue)
	if isSlice {
		for _, v := range simple {
			result = eval(scope, v)
		}
		return result
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
			return nil
		}
		first := arr[0]
		if s, ok := first.(Symbol); ok {
			f := scope.get_value(s)
			switch f2 := f.(type) {
			case func(*LispScope, LispValue) LispValue:
				return f2(scope, code)
			case func(LispValue) LispValue:
				return f2(eval(scope, arr[1]))
			case func(LispValue, LispValue) LispValue:
				return f2(eval(scope, arr[1]), eval(scope, arr[2]))
			case func(LispValue, LispValue, LispValue) LispValue:
				return f2(eval(scope, arr[1]), eval(scope, arr[2]), eval(scope, arr[3]))
			case func([]LispValue) LispValue:
				args := make([]LispValue, len(arr)-1)
				for i := range args {
					args[i] = eval(scope, arr[1+i])
				}
				return f2(args)
			case func([]LispValue):
				args := make([]LispValue, len(arr)-1)
				for i := range args {
					args[i] = eval(scope, arr[1+i])
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

		panic(fmt.Sprintf("%v", code))
	}
	if s, ok := code.(Symbol); ok {
		return scope.get_value(s)
	}

	return code
}
