package controllers

func errInternal(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}

func errRequest(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
	}
}
