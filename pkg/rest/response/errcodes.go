package response

type ErrCode uint

const (
	// ---- Part of technical errors on our side (can not unmarshal request etc) ----

	InvalidRequestStructure ErrCode = 10001
	InternalErrorCode       ErrCode = 10002

	// ---- Part of user errors. Missed argument, invalid type, etc ----

	WrongValue       ErrCode = 11000
	MissedValue      ErrCode = 11001
	UnsupportedValue ErrCode = 11002
	WrongFormat      ErrCode = 11003
	MaxLengthText    ErrCode = 1104
)
