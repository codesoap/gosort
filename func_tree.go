package gosort

import "sort"

type FuncTree struct {
	Left  *FuncTree
	value map[string][]string
	order []string // The order of the keys in Value
	Right *FuncTree
}

// Order returns the order of the keys in value.
func (t *FuncTree) Order() []string {
	return t.order
}

// Value returns the value.
func (t *FuncTree) Value() map[string][]string {
	return t.value
}

// SetValue sets the value and order. The order is alphabetic.
// Setting and using the order ensures deterministic results.
func (t *FuncTree) SetValue(value map[string][]string) {
	t.value = value
	var order = []string{}
	for caller := range value {
		order = append(order, caller)
	}
	sort.Strings(order)
	t.order = order
}

// Print is used to print the whole tree for debugging.
func (t *FuncTree) Print(indent int) {
	prefix := ""
	for i := 0; i < indent; i++ {
		prefix += "\t"
	}
	if t.Value() != nil {
		print(prefix, "Value: ")
		for _, f := range t.Order() {
			print(f, " ")
		}
		println()
	}
	if t.Left != nil {
		print(prefix, "Left:\n")
		t.Left.Print(indent + 1)
	}
	if t.Right != nil {
		print(prefix, "Right:\n")
		t.Right.Print(indent + 1)
	}
}
