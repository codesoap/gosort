package gosort

// BranchOutByTopology moves functions to its Left and Right, so that
// functions in the Left node are not called by the functions in the
// Right node.
func (t *FuncTree) BranchOutByTopology() {
	var allCallees = map[string]bool{}
	for _, callees := range t.Value() {
		for _, callee := range callees {
			allCallees[callee] = true
		}
	}
	tmpLeftValue := make(map[string][]string)
	tmpRightValue := make(map[string][]string)
	for _, caller := range t.Order() {
		callees := t.Value()[caller]
		if _, ok := allCallees[caller]; !ok {
			// The function is not called by functions in t.Value.
			tmpLeftValue[caller] = callees
		} else {
			// The function is called by functions in t.Value.
			tmpRightValue[caller] = callees
		}
	}
	if len(tmpLeftValue) > 0 && len(tmpRightValue) > 0 {
		t.SetValue(nil)
		t.Left = &FuncTree{}
		t.Left.SetValue(tmpLeftValue)
		t.Right = &FuncTree{}
		t.Right.SetValue(tmpRightValue)
		t.Right.BranchOutByTopology()
	}
}

// BranchOutByCallOrder splits leaves, so that functions that are
// called before others are left of functions that are called later.
// previousCallees is a list of callees, that were called before, in the
// order they fcalleesr.
func (t *FuncTree) BranchOutByCallOrder(previousCallees []string) []string {
	if t.Left != nil {
		previousCallees = t.Left.BranchOutByCallOrder(previousCallees)
	}
	if t.Value() != nil {
		previousCallees = t.appendNewCallees(previousCallees)
		tmpLeftValue, tmpRightValue := t.putUncalledFunctionsLeft(previousCallees)
		if len(tmpLeftValue) == 0 {
			tmpLeftValue, tmpRightValue = t.putFirstCalledFunctionLeft(previousCallees)
		}
		if len(tmpRightValue) > 0 {
			t.Left = &FuncTree{}
			t.Left.SetValue(tmpLeftValue)
			t.Right = &FuncTree{}
			t.Right.SetValue(tmpRightValue)
			t.SetValue(nil)
		}
	}
	if t.Right != nil {
		t.Right.BranchOutByCallOrder(previousCallees)
	}
	return previousCallees
}

func (t *FuncTree) appendNewCallees(previousCallees []string) []string {
	for _, caller := range t.Order() {
		for _, callee := range t.Value()[caller] {
			if !contains(previousCallees, callee) {
				previousCallees = append(previousCallees, callee)
			}
		}
	}
	return previousCallees
}

func (t *FuncTree) putUncalledFunctionsLeft(previousCallees []string) (map[string][]string, map[string][]string) {
	tmpLeftValue := make(map[string][]string)
	tmpRightValue := make(map[string][]string)
	for _, caller := range t.Order() {
		if contains(previousCallees, caller) {
			tmpRightValue[caller] = t.Value()[caller]
		} else {
			tmpLeftValue[caller] = t.Value()[caller]
		}
	}
	return tmpLeftValue, tmpRightValue
}

func (t *FuncTree) putFirstCalledFunctionLeft(previousCallees []string) (map[string][]string, map[string][]string) {
	tmpLeftValue := make(map[string][]string)
	tmpRightValue := make(map[string][]string)
	leftFilled := false
	for _, callee := range previousCallees {
		if callees, ok := t.Value()[callee]; ok {
			if !leftFilled {
				tmpLeftValue[callee] = callees
				leftFilled = true
			} else {
				tmpRightValue[callee] = callees
			}
		}
	}
	return tmpLeftValue, tmpRightValue
}
