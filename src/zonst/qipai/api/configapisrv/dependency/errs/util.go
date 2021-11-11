package errs

func GetString(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}

	vi, exist := m[key]
	if !exist {
		return ""
	}

	v, convert := vi.(string)
	if !convert {
		return ""
	}
	return v
}

func StringDefault(src string, defaultV string) string {
	if src == "" {
		return defaultV
	}
	return src
}
