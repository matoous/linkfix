// Code generated by "stringer -type=Severity"; DO NOT EDIT.

package severity

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Error-0]
	_ = x[Warn-1]
	_ = x[Info-2]
}

const _Severity_name = "ErrorWarnInfo"

var _Severity_index = [...]uint8{0, 5, 9, 13}

func (i Severity) String() string {
	if i >= Severity(len(_Severity_index)-1) {
		return "Severity(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Severity_name[_Severity_index[i]:_Severity_index[i+1]]
}
