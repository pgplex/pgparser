// Package nodes provides nodeToString functionality for serializing parse tree nodes.
// This is a Go implementation of PostgreSQL's outfuncs.c.
package nodes

import (
	"fmt"
	"strconv"
	"strings"
)

// NodeToString converts a Node to its string representation.
// This matches PostgreSQL's nodeToString() function format.
func NodeToString(node Node) string {
	if node == nil {
		return "<>"
	}
	var sb strings.Builder
	writeNode(&sb, node)
	return sb.String()
}

// writeNode writes a node to the string builder in PostgreSQL format.
func writeNode(sb *strings.Builder, node Node) {
	if node == nil {
		sb.WriteString("<>")
		return
	}

	switch n := node.(type) {
	case *List:
		writeList(sb, n)
	case *Integer:
		writeInteger(sb, n)
	case *Float:
		writeFloat(sb, n)
	case *Boolean:
		writeBoolean(sb, n)
	case *String:
		writeString(sb, n)
	case *BitString:
		writeBitString(sb, n)
	case *SelectStmt:
		writeSelectStmt(sb, n)
	case *InsertStmt:
		writeInsertStmt(sb, n)
	case *UpdateStmt:
		writeUpdateStmt(sb, n)
	case *DeleteStmt:
		writeDeleteStmt(sb, n)
	case *RangeVar:
		writeRangeVar(sb, n)
	case *ColumnRef:
		writeColumnRef(sb, n)
	case *A_Const:
		writeAConst(sb, n)
	case *A_Expr:
		writeAExpr(sb, n)
	case *TypeCast:
		writeTypeCast(sb, n)
	case *TypeName:
		writeTypeName(sb, n)
	case *FuncCall:
		writeFuncCall(sb, n)
	case *ResTarget:
		writeResTarget(sb, n)
	case *SortBy:
		writeSortBy(sb, n)
	case *Alias:
		writeAlias(sb, n)
	case *A_Star:
		writeAStar(sb, n)
	case *BoolExpr:
		writeBoolExpr(sb, n)
	case *NullTest:
		writeNullTest(sb, n)
	case *SubLink:
		writeSubLink(sb, n)
	case *JoinExpr:
		writeJoinExpr(sb, n)
	case *CreateStmt:
		writeCreateStmt(sb, n)
	case *ColumnDef:
		writeColumnDef(sb, n)
	case *Constraint:
		writeConstraint(sb, n)
	case *IndexStmt:
		writeIndexStmt(sb, n)
	case *ViewStmt:
		writeViewStmt(sb, n)
	case *RawStmt:
		writeRawStmt(sb, n)
	default:
		// Generic fallback for unhandled node types
		sb.WriteString("{")
		sb.WriteString(NodeTagName(node.Tag()))
		sb.WriteString("}")
	}
}

// Helper functions for writing specific node types

func writeList(sb *strings.Builder, n *List) {
	sb.WriteString("(")
	for i, item := range n.Items {
		if i > 0 {
			sb.WriteString(" ")
		}
		writeNode(sb, item)
	}
	sb.WriteString(")")
}

func writeInteger(sb *strings.Builder, n *Integer) {
	sb.WriteString(strconv.FormatInt(n.Ival, 10))
}

func writeFloat(sb *strings.Builder, n *Float) {
	sb.WriteString(n.Fval)
}

func writeBoolean(sb *strings.Builder, n *Boolean) {
	if n.Boolval {
		sb.WriteString("true")
	} else {
		sb.WriteString("false")
	}
}

func writeString(sb *strings.Builder, n *String) {
	sb.WriteString("\"")
	sb.WriteString(escapeString(n.Str))
	sb.WriteString("\"")
}

func writeBitString(sb *strings.Builder, n *BitString) {
	sb.WriteString("b\"")
	sb.WriteString(n.Bsval)
	sb.WriteString("\"")
}

func writeSelectStmt(sb *strings.Builder, n *SelectStmt) {
	sb.WriteString("{SELECTSTMT")
	if n.DistinctClause != nil {
		sb.WriteString(" :distinctClause ")
		writeNode(sb, n.DistinctClause)
	}
	if n.TargetList != nil {
		sb.WriteString(" :targetList ")
		writeNode(sb, n.TargetList)
	}
	if n.FromClause != nil {
		sb.WriteString(" :fromClause ")
		writeNode(sb, n.FromClause)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :whereClause ")
		writeNode(sb, n.WhereClause)
	}
	if n.GroupClause != nil {
		sb.WriteString(" :groupClause ")
		writeNode(sb, n.GroupClause)
	}
	if n.HavingClause != nil {
		sb.WriteString(" :havingClause ")
		writeNode(sb, n.HavingClause)
	}
	if n.SortClause != nil {
		sb.WriteString(" :sortClause ")
		writeNode(sb, n.SortClause)
	}
	if n.LimitOffset != nil {
		sb.WriteString(" :limitOffset ")
		writeNode(sb, n.LimitOffset)
	}
	if n.LimitCount != nil {
		sb.WriteString(" :limitCount ")
		writeNode(sb, n.LimitCount)
	}
	if n.Op != SETOP_NONE {
		sb.WriteString(fmt.Sprintf(" :op %d", n.Op))
	}
	if n.All {
		sb.WriteString(" :all true")
	}
	if n.Larg != nil {
		sb.WriteString(" :larg ")
		writeNode(sb, n.Larg)
	}
	if n.Rarg != nil {
		sb.WriteString(" :rarg ")
		writeNode(sb, n.Rarg)
	}
	sb.WriteString("}")
}

func writeInsertStmt(sb *strings.Builder, n *InsertStmt) {
	sb.WriteString("{INSERTSTMT")
	if n.Relation != nil {
		sb.WriteString(" :relation ")
		writeNode(sb, n.Relation)
	}
	if n.Cols != nil {
		sb.WriteString(" :cols ")
		writeNode(sb, n.Cols)
	}
	if n.SelectStmt != nil {
		sb.WriteString(" :selectStmt ")
		writeNode(sb, n.SelectStmt)
	}
	if n.ReturningList != nil {
		sb.WriteString(" :returningList ")
		writeNode(sb, n.ReturningList)
	}
	sb.WriteString("}")
}

func writeUpdateStmt(sb *strings.Builder, n *UpdateStmt) {
	sb.WriteString("{UPDATESTMT")
	if n.Relation != nil {
		sb.WriteString(" :relation ")
		writeNode(sb, n.Relation)
	}
	if n.TargetList != nil {
		sb.WriteString(" :targetList ")
		writeNode(sb, n.TargetList)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :whereClause ")
		writeNode(sb, n.WhereClause)
	}
	if n.FromClause != nil {
		sb.WriteString(" :fromClause ")
		writeNode(sb, n.FromClause)
	}
	if n.ReturningList != nil {
		sb.WriteString(" :returningList ")
		writeNode(sb, n.ReturningList)
	}
	sb.WriteString("}")
}

func writeDeleteStmt(sb *strings.Builder, n *DeleteStmt) {
	sb.WriteString("{DELETESTMT")
	if n.Relation != nil {
		sb.WriteString(" :relation ")
		writeNode(sb, n.Relation)
	}
	if n.UsingClause != nil {
		sb.WriteString(" :usingClause ")
		writeNode(sb, n.UsingClause)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :whereClause ")
		writeNode(sb, n.WhereClause)
	}
	if n.ReturningList != nil {
		sb.WriteString(" :returningList ")
		writeNode(sb, n.ReturningList)
	}
	sb.WriteString("}")
}

func writeRangeVar(sb *strings.Builder, n *RangeVar) {
	sb.WriteString("{RANGEVAR")
	if n.Catalogname != "" {
		sb.WriteString(" :catalogname ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Catalogname))
		sb.WriteString("\"")
	}
	if n.Schemaname != "" {
		sb.WriteString(" :schemaname ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Schemaname))
		sb.WriteString("\"")
	}
	if n.Relname != "" {
		sb.WriteString(" :relname ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Relname))
		sb.WriteString("\"")
	}
	sb.WriteString(fmt.Sprintf(" :inh %t", n.Inh))
	sb.WriteString(fmt.Sprintf(" :relpersistence %c", n.Relpersistence))
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeColumnRef(sb *strings.Builder, n *ColumnRef) {
	sb.WriteString("{COLUMNREF")
	if n.Fields != nil {
		sb.WriteString(" :fields ")
		writeNode(sb, n.Fields)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeAConst(sb *strings.Builder, n *A_Const) {
	sb.WriteString("{A_CONST")
	if n.Isnull {
		sb.WriteString(" :isnull true")
	} else if n.Val != nil {
		sb.WriteString(" :val ")
		writeNode(sb, n.Val)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeAExpr(sb *strings.Builder, n *A_Expr) {
	sb.WriteString("{A_EXPR")
	sb.WriteString(fmt.Sprintf(" :kind %d", n.Kind))
	if n.Name != nil {
		sb.WriteString(" :name ")
		writeNode(sb, n.Name)
	}
	if n.Lexpr != nil {
		sb.WriteString(" :lexpr ")
		writeNode(sb, n.Lexpr)
	}
	if n.Rexpr != nil {
		sb.WriteString(" :rexpr ")
		writeNode(sb, n.Rexpr)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeTypeCast(sb *strings.Builder, n *TypeCast) {
	sb.WriteString("{TYPECAST")
	if n.Arg != nil {
		sb.WriteString(" :arg ")
		writeNode(sb, n.Arg)
	}
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeTypeName(sb *strings.Builder, n *TypeName) {
	sb.WriteString("{TYPENAME")
	if n.Names != nil {
		sb.WriteString(" :names ")
		writeNode(sb, n.Names)
	}
	sb.WriteString(fmt.Sprintf(" :typeOid %d", n.TypeOid))
	sb.WriteString(fmt.Sprintf(" :setof %t", n.Setof))
	sb.WriteString(fmt.Sprintf(" :pct_type %t", n.PctType))
	if n.Typmods != nil {
		sb.WriteString(" :typmods ")
		writeNode(sb, n.Typmods)
	}
	sb.WriteString(fmt.Sprintf(" :typemod %d", n.Typemod))
	if n.ArrayBounds != nil {
		sb.WriteString(" :arrayBounds ")
		writeNode(sb, n.ArrayBounds)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeFuncCall(sb *strings.Builder, n *FuncCall) {
	sb.WriteString("{FUNCCALL")
	if n.Funcname != nil {
		sb.WriteString(" :funcname ")
		writeNode(sb, n.Funcname)
	}
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	if n.AggOrder != nil {
		sb.WriteString(" :agg_order ")
		writeNode(sb, n.AggOrder)
	}
	if n.AggFilter != nil {
		sb.WriteString(" :agg_filter ")
		writeNode(sb, n.AggFilter)
	}
	if n.Over != nil {
		sb.WriteString(" :over ")
		writeNode(sb, n.Over)
	}
	sb.WriteString(fmt.Sprintf(" :agg_within_group %t", n.AggWithinGroup))
	sb.WriteString(fmt.Sprintf(" :agg_star %t", n.AggStar))
	sb.WriteString(fmt.Sprintf(" :agg_distinct %t", n.AggDistinct))
	sb.WriteString(fmt.Sprintf(" :func_variadic %t", n.FuncVariadic))
	sb.WriteString(fmt.Sprintf(" :funcformat %d", n.FuncFormat))
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeResTarget(sb *strings.Builder, n *ResTarget) {
	sb.WriteString("{RESTARGET")
	if n.Name != "" {
		sb.WriteString(" :name ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Name))
		sb.WriteString("\"")
	}
	if n.Indirection != nil {
		sb.WriteString(" :indirection ")
		writeNode(sb, n.Indirection)
	}
	if n.Val != nil {
		sb.WriteString(" :val ")
		writeNode(sb, n.Val)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeSortBy(sb *strings.Builder, n *SortBy) {
	sb.WriteString("{SORTBY")
	if n.Node != nil {
		sb.WriteString(" :node ")
		writeNode(sb, n.Node)
	}
	sb.WriteString(fmt.Sprintf(" :sortby_dir %d", n.SortbyDir))
	sb.WriteString(fmt.Sprintf(" :sortby_nulls %d", n.SortbyNulls))
	if n.UseOp != nil {
		sb.WriteString(" :useOp ")
		writeNode(sb, n.UseOp)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeAlias(sb *strings.Builder, n *Alias) {
	sb.WriteString("{ALIAS")
	if n.Aliasname != "" {
		sb.WriteString(" :aliasname ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Aliasname))
		sb.WriteString("\"")
	}
	if n.Colnames != nil {
		sb.WriteString(" :colnames ")
		writeNode(sb, n.Colnames)
	}
	sb.WriteString("}")
}

func writeAStar(sb *strings.Builder, n *A_Star) {
	sb.WriteString("{A_STAR}")
}

func writeBoolExpr(sb *strings.Builder, n *BoolExpr) {
	sb.WriteString("{BOOLEXPR")
	sb.WriteString(fmt.Sprintf(" :boolop %d", n.Boolop))
	if n.Args != nil {
		sb.WriteString(" :args ")
		writeNode(sb, n.Args)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeNullTest(sb *strings.Builder, n *NullTest) {
	sb.WriteString("{NULLTEST")
	if n.Arg != nil {
		sb.WriteString(" :arg ")
		writeNode(sb, n.Arg)
	}
	sb.WriteString(fmt.Sprintf(" :nulltesttype %d", n.Nulltesttype))
	sb.WriteString(fmt.Sprintf(" :argisrow %t", n.Argisrow))
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeSubLink(sb *strings.Builder, n *SubLink) {
	sb.WriteString("{SUBLINK")
	sb.WriteString(fmt.Sprintf(" :subLinkType %d", n.SubLinkType))
	sb.WriteString(fmt.Sprintf(" :subLinkId %d", n.SubLinkId))
	if n.Testexpr != nil {
		sb.WriteString(" :testexpr ")
		writeNode(sb, n.Testexpr)
	}
	if n.OperName != nil {
		sb.WriteString(" :operName ")
		writeNode(sb, n.OperName)
	}
	if n.Subselect != nil {
		sb.WriteString(" :subselect ")
		writeNode(sb, n.Subselect)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeJoinExpr(sb *strings.Builder, n *JoinExpr) {
	sb.WriteString("{JOINEXPR")
	sb.WriteString(fmt.Sprintf(" :jointype %d", n.Jointype))
	sb.WriteString(fmt.Sprintf(" :isNatural %t", n.IsNatural))
	if n.Larg != nil {
		sb.WriteString(" :larg ")
		writeNode(sb, n.Larg)
	}
	if n.Rarg != nil {
		sb.WriteString(" :rarg ")
		writeNode(sb, n.Rarg)
	}
	if n.UsingClause != nil {
		sb.WriteString(" :usingClause ")
		writeNode(sb, n.UsingClause)
	}
	if n.Quals != nil {
		sb.WriteString(" :quals ")
		writeNode(sb, n.Quals)
	}
	if n.Alias != nil {
		sb.WriteString(" :alias ")
		writeNode(sb, n.Alias)
	}
	sb.WriteString(fmt.Sprintf(" :rtindex %d", n.Rtindex))
	sb.WriteString("}")
}

func writeCreateStmt(sb *strings.Builder, n *CreateStmt) {
	sb.WriteString("{CREATESTMT")
	if n.Relation != nil {
		sb.WriteString(" :relation ")
		writeNode(sb, n.Relation)
	}
	if n.TableElts != nil {
		sb.WriteString(" :tableElts ")
		writeNode(sb, n.TableElts)
	}
	if n.InhRelations != nil {
		sb.WriteString(" :inhRelations ")
		writeNode(sb, n.InhRelations)
	}
	if n.Constraints != nil {
		sb.WriteString(" :constraints ")
		writeNode(sb, n.Constraints)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :oncommit %d", n.OnCommit))
	if n.Tablespacename != "" {
		sb.WriteString(" :tablespacename ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Tablespacename))
		sb.WriteString("\"")
	}
	sb.WriteString(fmt.Sprintf(" :if_not_exists %t", n.IfNotExists))
	sb.WriteString("}")
}

func writeColumnDef(sb *strings.Builder, n *ColumnDef) {
	sb.WriteString("{COLUMNDEF")
	if n.Colname != "" {
		sb.WriteString(" :colname ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Colname))
		sb.WriteString("\"")
	}
	if n.TypeName != nil {
		sb.WriteString(" :typeName ")
		writeNode(sb, n.TypeName)
	}
	sb.WriteString(fmt.Sprintf(" :compression %q", n.Compression))
	sb.WriteString(fmt.Sprintf(" :inhcount %d", n.Inhcount))
	sb.WriteString(fmt.Sprintf(" :is_local %t", n.IsLocal))
	sb.WriteString(fmt.Sprintf(" :is_not_null %t", n.IsNotNull))
	sb.WriteString(fmt.Sprintf(" :is_from_type %t", n.IsFromType))
	sb.WriteString(fmt.Sprintf(" :storage %c", n.Storage))
	if n.RawDefault != nil {
		sb.WriteString(" :raw_default ")
		writeNode(sb, n.RawDefault)
	}
	if n.CookedDefault != nil {
		sb.WriteString(" :cooked_default ")
		writeNode(sb, n.CookedDefault)
	}
	if n.CollClause != nil {
		sb.WriteString(" :collClause ")
		writeNode(sb, n.CollClause)
	}
	if n.Constraints != nil {
		sb.WriteString(" :constraints ")
		writeNode(sb, n.Constraints)
	}
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString("}")
}

func writeConstraint(sb *strings.Builder, n *Constraint) {
	sb.WriteString("{CONSTRAINT")
	sb.WriteString(fmt.Sprintf(" :contype %d", n.Contype))
	if n.Conname != "" {
		sb.WriteString(" :conname ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Conname))
		sb.WriteString("\"")
	}
	sb.WriteString(fmt.Sprintf(" :deferrable %t", n.Deferrable))
	sb.WriteString(fmt.Sprintf(" :initdeferred %t", n.Initdeferred))
	sb.WriteString(fmt.Sprintf(" :location %d", n.Location))
	sb.WriteString(fmt.Sprintf(" :is_no_inherit %t", n.IsNoInherit))
	if n.RawExpr != nil {
		sb.WriteString(" :raw_expr ")
		writeNode(sb, n.RawExpr)
	}
	if n.Keys != nil {
		sb.WriteString(" :keys ")
		writeNode(sb, n.Keys)
	}
	if n.Including != nil {
		sb.WriteString(" :including ")
		writeNode(sb, n.Including)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	if n.Indexname != "" {
		sb.WriteString(" :indexname ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Indexname))
		sb.WriteString("\"")
	}
	if n.Indexspace != "" {
		sb.WriteString(" :indexspace ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Indexspace))
		sb.WriteString("\"")
	}
	sb.WriteString("}")
}

func writeIndexStmt(sb *strings.Builder, n *IndexStmt) {
	sb.WriteString("{INDEXSTMT")
	if n.Idxname != "" {
		sb.WriteString(" :idxname ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.Idxname))
		sb.WriteString("\"")
	}
	if n.Relation != nil {
		sb.WriteString(" :relation ")
		writeNode(sb, n.Relation)
	}
	if n.AccessMethod != "" {
		sb.WriteString(" :accessMethod ")
		sb.WriteString("\"")
		sb.WriteString(escapeString(n.AccessMethod))
		sb.WriteString("\"")
	}
	if n.IndexParams != nil {
		sb.WriteString(" :indexParams ")
		writeNode(sb, n.IndexParams)
	}
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	if n.WhereClause != nil {
		sb.WriteString(" :whereClause ")
		writeNode(sb, n.WhereClause)
	}
	sb.WriteString(fmt.Sprintf(" :unique %t", n.Unique))
	sb.WriteString(fmt.Sprintf(" :primary %t", n.Primary))
	sb.WriteString(fmt.Sprintf(" :isconstraint %t", n.Isconstraint))
	sb.WriteString(fmt.Sprintf(" :deferrable %t", n.Deferrable))
	sb.WriteString(fmt.Sprintf(" :initdeferred %t", n.Initdeferred))
	sb.WriteString(fmt.Sprintf(" :concurrent %t", n.Concurrent))
	sb.WriteString(fmt.Sprintf(" :if_not_exists %t", n.IfNotExists))
	sb.WriteString("}")
}

func writeViewStmt(sb *strings.Builder, n *ViewStmt) {
	sb.WriteString("{VIEWSTMT")
	if n.View != nil {
		sb.WriteString(" :view ")
		writeNode(sb, n.View)
	}
	if n.Aliases != nil {
		sb.WriteString(" :aliases ")
		writeNode(sb, n.Aliases)
	}
	if n.Query != nil {
		sb.WriteString(" :query ")
		writeNode(sb, n.Query)
	}
	sb.WriteString(fmt.Sprintf(" :replace %t", n.Replace))
	if n.Options != nil {
		sb.WriteString(" :options ")
		writeNode(sb, n.Options)
	}
	sb.WriteString(fmt.Sprintf(" :withCheckOption %d", n.WithCheckOption))
	sb.WriteString("}")
}

func writeRawStmt(sb *strings.Builder, n *RawStmt) {
	sb.WriteString("{RAWSTMT")
	if n.Stmt != nil {
		sb.WriteString(" :stmt ")
		writeNode(sb, n.Stmt)
	}
	sb.WriteString(fmt.Sprintf(" :stmt_location %d", n.StmtLocation))
	sb.WriteString(fmt.Sprintf(" :stmt_len %d", n.StmtLen))
	sb.WriteString("}")
}

// escapeString escapes special characters in a string for output.
func escapeString(s string) string {
	var sb strings.Builder
	for _, c := range s {
		switch c {
		case '\\':
			sb.WriteString("\\\\")
		case '"':
			sb.WriteString("\\\"")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		default:
			sb.WriteRune(c)
		}
	}
	return sb.String()
}
