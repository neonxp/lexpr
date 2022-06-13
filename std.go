package lexpr

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Operator struct {
	handler   func(ts *TokenStack) error
	priority  int
	leftAssoc bool
}

var Operators = map[string]Operator{
	".": {
		handler: func(ts *TokenStack) error {
			t2 := ts.Pop()
			t1 := ts.Pop()
			switch t2.typ {
			case str, word:
				m := map[string]json.RawMessage{}
				if err := json.Unmarshal([]byte(t1.value), &m); err != nil {
					return fmt.Errorf("invalid json %s err: %s", t1.value, err.Error())
				}
				val, ok := m[t2.value]
				if !ok {
					return fmt.Errorf("invalid json key %s key: %s", t1.value, t2.value)
				}
				ts.Push(Token{
					typ:   str,
					value: strings.Trim(string(val), `"`),
				})
			case number:
				m := []json.RawMessage{}
				if err := json.Unmarshal([]byte(t1.value), &m); err != nil {
					return fmt.Errorf("invalid json %s err: %s", t1.value, err.Error())
				}
				if len(m) <= t2.ivalue {
					return fmt.Errorf("invalid json key %s key: %s", t1.value, t2.value)
				}
				val := m[t2.ivalue]
				ts.Push(Token{
					typ:   str,
					value: strings.Trim(string(val), `"`),
				})
			default:
				return fmt.Errorf("invalid json key: %+v", t2)
			}
			return nil
		},
		priority:  140,
		leftAssoc: false,
	},
	// Math operators
	"**": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			ts.Push(Token{
				typ:    number,
				ivalue: int(math.Pow(float64(t1.ivalue), float64(t2.ivalue))),
			})
			return nil
		},
		priority:  130,
		leftAssoc: true,
	},
	"*": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			ts.Push(Token{
				typ:    number,
				ivalue: t1.ivalue * t2.ivalue,
			})
			return nil
		},
		priority:  120,
		leftAssoc: false,
	},
	"/": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			ts.Push(Token{
				typ:    number,
				ivalue: t1.ivalue / t2.ivalue,
			})
			return nil
		},
		priority:  120,
		leftAssoc: false,
	},
	"%": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			ts.Push(Token{
				typ:    number,
				ivalue: t1.ivalue % t2.ivalue,
			})
			return nil
		},
		priority:  120,
		leftAssoc: false,
	},
	"+": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			ts.Push(Token{
				typ:    number,
				ivalue: t1.ivalue + t2.ivalue,
			})
			return nil
		},
		priority:  110,
		leftAssoc: false,
	},
	"-": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			ts.Push(Token{
				typ:    number,
				ivalue: t1.ivalue - t2.ivalue,
			})
			return nil
		},
		priority:  110,
		leftAssoc: false,
	},

	// Logic operators
	"!": {
		handler: func(ts *TokenStack) error {
			t := ts.Pop()
			switch ts.Pop().typ {
			case number:
				t.ivalue = ^t.ivalue
				ts.Push(t)
			default:
				return fmt.Errorf("Argument must be number, got %+v", t)
			}
			return nil
		},
		priority:  50,
		leftAssoc: false,
	},
	">": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			r := 0
			if t2.ivalue > t1.ivalue {
				r = 1
			}
			ts.Push(Token{
				typ:    number,
				ivalue: r,
			})
			return nil
		},
		priority:  20,
		leftAssoc: false,
	},
	">=": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			r := 0
			if t2.ivalue >= t1.ivalue {
				r = 1
			}
			ts.Push(Token{
				typ:    number,
				ivalue: r,
			})
			return nil
		},
		priority:  20,
		leftAssoc: false,
	},
	"<": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			r := 0
			if t2.ivalue < t1.ivalue {
				r = 1
			}
			ts.Push(Token{
				typ:    number,
				ivalue: r,
			})
			return nil
		},
		priority:  20,
		leftAssoc: false,
	},
	"<=": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			r := 0
			if t2.ivalue <= t1.ivalue {
				r = 1
			}
			ts.Push(Token{
				typ:    number,
				ivalue: r,
			})
			return nil
		},
		priority:  20,
		leftAssoc: false,
	},
	"==": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			r := 0
			if t1.typ == number && t2.typ == number && t1.ivalue == t2.ivalue {
				r = 1
			} else if t1.value == t2.value {
				r = 1
			}
			ts.Push(Token{
				typ:    number,
				ivalue: r,
			})
			return nil
		},
		priority:  20,
		leftAssoc: false,
	},
	"!=": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			r := 0
			if t1.typ == number && t2.typ == number && t1.ivalue != t2.ivalue {
				r = 1
			} else if t1.value != t2.value {
				r = 1
			}
			ts.Push(Token{
				typ:    number,
				ivalue: r,
			})
			return nil
		},
		priority:  20,
		leftAssoc: false,
	},
	"&&": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			b1 := true
			b2 := true
			if t1.ivalue == 0 {
				b1 = false
			}
			if t2.ivalue == 0 {
				b2 = false
			}
			r := 0
			if b1 && b2 {
				r = 1
			}
			ts.Push(Token{
				typ:    number,
				ivalue: r,
			})
			return nil
		},
		priority:  10,
		leftAssoc: false,
	},
	"||": {
		handler: func(ts *TokenStack) error {
			t1 := ts.Pop()
			t2 := ts.Pop()
			if t1.typ != number || t2.typ != number {
				return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
			}
			b1 := true
			b2 := true
			if t1.ivalue == 0 {
				b1 = false
			}
			if t2.ivalue == 0 {
				b2 = false
			}
			r := 0
			if b1 || b2 {
				r = 1
			}
			ts.Push(Token{
				typ:    number,
				ivalue: r,
			})
			return nil
		},
		priority:  0,
		leftAssoc: false,
	},
}

var Functions = map[string]func(ts *TokenStack) error{
	"max": func(ts *TokenStack) error {
		t1 := ts.Pop()
		t2 := ts.Pop()
		if t1.typ != number || t2.typ != number {
			return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
		}
		max := t1.ivalue
		if t2.ivalue > max {
			max = t2.ivalue
		}
		ts.Push(Token{
			typ:    number,
			ivalue: max,
		})
		return nil
	},
	"min": func(ts *TokenStack) error {
		t1 := ts.Pop()
		t2 := ts.Pop()
		if t1.typ != number || t2.typ != number {
			return fmt.Errorf("Both arguments must be number, got op1 = %+v, op2 = %+v", t1, t2)
		}
		min := t1.ivalue
		if t2.ivalue < min {
			min = t2.ivalue
		}
		ts.Push(Token{
			typ:    number,
			ivalue: min,
		})
		return nil
	},
	"len": func(ts *TokenStack) error {
		t := ts.Pop()
		ts.Push(Token{
			typ:    number,
			ivalue: len(t.value),
		})
		return nil
	},
	"atoi": func(ts *TokenStack) error {
		t := ts.Pop()
		if t.typ != str && t.typ != word {
			return fmt.Errorf("atoi requires string argument, got %+v", t)
		}
		n, err := strconv.Atoi(t.value)
		if err != nil {
			return err
		}
		ts.Push(Token{
			typ:    number,
			ivalue: n,
		})
		return nil
	},
	"itoa": func(ts *TokenStack) error {
		t := ts.Pop()
		if t.typ != number {
			return fmt.Errorf("itoa requires number argument, got %+v", t)
		}
		s := strconv.Itoa(t.ivalue)
		ts.Push(Token{
			typ:   str,
			value: s,
		})
		return nil
	},
}
