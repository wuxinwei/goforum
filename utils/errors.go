package utils

func ErrInternal(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}

func ErrRequest(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}
