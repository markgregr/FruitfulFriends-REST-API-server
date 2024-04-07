package response

func ResolveError(err error) Error {
	/*switch {
	case errors.Is(err, auth.ErrDuplicateEmail):
		return NewDuplicateEntryError()
	case errors.Is(err, auth.ErrInvalidCredentials):
		return NewInvalidCredentialsError()
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NewNotFoundError()
	case errors.Is(err, auth.ErrInActiveProfile):
		return NewPermissionDeniedError()
	}*/

	return NewInternalError()
}
