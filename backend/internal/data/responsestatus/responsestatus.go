package responsestatus

// -------------------------------------------------------------------------
type ResponseStatus uint8

const (
	Done ResponseStatus = iota + 1
	Pending
	Running
	Error
)

//go:generate go-enum -type=ResponseStatus
