package version

// These variables should be initialized by the linker. They should not be initialized in code.
var (
	// GitCommit is the commit that was compiled, it is filled in by the compiler.
	GitCommit string
	GitTag    string
)
