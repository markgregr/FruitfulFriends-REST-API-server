package response

type ErrCode uint

const (
	InvalidRequestStructure ErrCode = 10001
	MissedValue             ErrCode = 11001
	EmptyField              ErrCode = 11002
)
