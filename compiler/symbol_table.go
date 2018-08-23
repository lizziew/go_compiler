package compiler

// Differentiate between different scopes for symbols
type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

// Stores Name, Scope, and Index for a given symbol
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// Associate string identifiers with scope and unique number
type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func BuildSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

// Create and store a symbol from an identifier
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{name, GlobalScope, s.numDefinitions}

	s.store[name] = symbol
	s.numDefinitions++

	return symbol
}

// Retrieve a symbol for an identifier
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	return obj, ok
}
