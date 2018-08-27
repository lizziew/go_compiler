package compiler

// Differentiate between different scopes for symbols
type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

// Stores Name, Scope, and Index for a given symbol
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// Associate string identifiers with scope and unique number
type SymbolTable struct {
	Outer          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
}

func BuildSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func BuildInnerSymbolTable(outer *SymbolTable) *SymbolTable {
	inner := BuildSymbolTable()
	inner.Outer = outer
	return inner
}

// Create and store a symbol from an identifier
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}

	// Set scope
	if s.Outer == nil {
		// Is outer symbol table
		symbol.Scope = GlobalScope
	} else {
		// Is inner symbol table
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
	return symbol
}

// Retrieve a symbol for an identifier
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	// Check outer environment if exists
	if !ok && s.Outer != nil {
		obj, ok := s.Outer.Resolve(name)
		return obj, ok
	} else {
		return obj, ok
	}
}
