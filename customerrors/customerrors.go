package customerrors

import "errors"

var (
	ErrOnlyAdminAccess = errors.New("only admin can access this")
	ErrSessionIsClosed = errors.New("session is closed")
	ErrDuplicateSession = errors.New("a session with the same name is already open")
)
