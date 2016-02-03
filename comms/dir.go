package comms

// Specifies what kind of connection to make.
type Dir int

const (
	// Default value - this is not valid
	None Dir = iota
	// Inbound half of connection only.
	Input
	// Outbound half of connection only.
	Output
	// Full duplex
	Duplex = Input | Output
)
