// Package nodes defines PostgreSQL parse tree node types.
// These types are designed to be 100% compatible with PostgreSQL's internal representation.
package nodes

// Node is the interface implemented by all parse tree nodes.
type Node interface {
	// Tag returns the NodeTag for this node type.
	Tag() NodeTag
}

// ParseLoc represents a location in the source text (byte offset).
// -1 means "unknown".
type ParseLoc int

// Oid represents a PostgreSQL object identifier.
type Oid uint32

// Constants for special OID values.
const (
	InvalidOid Oid = 0
)

// List represents a PostgreSQL List.
// In PostgreSQL, List is a generic linked list that can hold any node type.
type List struct {
	// In Go, we use a slice instead of a linked list.
	// The element type depends on context.
	Items []Node
}

func (l *List) Tag() NodeTag { return T_List }

// Len returns the number of items in the list.
func (l *List) Len() int {
	if l == nil {
		return 0
	}
	return len(l.Items)
}

// IntList represents a list of integers.
type IntList struct {
	Items []int
}

func (l *IntList) Tag() NodeTag { return T_IntList }

// OidList represents a list of OIDs.
type OidList struct {
	Items []Oid
}

func (l *OidList) Tag() NodeTag { return T_OidList }

// String represents a PostgreSQL String value node.
type String struct {
	Str string
}

func (s *String) Tag() NodeTag { return T_String }

// Integer represents a PostgreSQL Integer value node.
type Integer struct {
	Ival int64
}

func (i *Integer) Tag() NodeTag { return T_Integer }

// Float represents a PostgreSQL Float value node.
// Note: PostgreSQL stores floats as strings to preserve precision.
type Float struct {
	Fval string
}

func (f *Float) Tag() NodeTag { return T_Float }

// Boolean represents a PostgreSQL Boolean value node.
type Boolean struct {
	Boolval bool
}

func (b *Boolean) Tag() NodeTag { return T_Boolean }

// BitString represents a PostgreSQL bit string value node.
type BitString struct {
	Bsval string
}

func (b *BitString) Tag() NodeTag { return T_BitString }
