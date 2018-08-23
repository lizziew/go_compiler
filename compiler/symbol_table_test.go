package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": Symbol{"a", GlobalScope, 0},
		"b": Symbol{"b", GlobalScope, 1},
	}

	global := BuildSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Fatalf("a is wrong (define)")
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Fatalf("b is wrong (define)")
	}
}

func TestResolveGlobal(t *testing.T) {
	global := BuildSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		{"a", GlobalScope, 0},
		{"b", GlobalScope, 1},
	}

	for _, e := range expected {
		result, ok := global.Resolve(e.Name)
		if !ok {
			t.Fatalf("not resolvable")
		}

		if result != e {
			t.Fatalf("resolve error")
		}
	}
}
