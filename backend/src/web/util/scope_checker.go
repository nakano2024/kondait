package util

func HasAnyScope(requireds []string, granteds []string) bool {
	if len(requireds) == 0 {
		return false
	}

	if len(granteds) == 0 {
		return false
	}

	reqMap := make(map[string]struct{}, len(requireds))
	for _, req := range requireds {
		reqMap[req] = struct{}{}
	}

	for _, grant := range granteds {
		if _, found := reqMap[grant]; found {
			return true
		}
	}

	return false
}
