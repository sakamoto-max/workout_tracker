package customerrors

import "errors"

var (
	ErrOnlyAdminAccess = errors.New("only admin can access this")
	ErrSessionIsClosed = errors.New("session is closed")
	ErrDuplicateSession = errors.New("a session with the same name is already open")
	ErrTokenIsInvalid = errors.New("token is invalid")
	ErrEmailDoesntExist = errors.New("please sign up first")
	ErrExerciseAlrExists = errors.New("exercise already exists")
	ErrPlanDoesNotExist = errors.New("plan does not exist")
	ErrExerciseDoesNotExist = errors.New("exercise does not exist in plan")
)

