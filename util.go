package gojsonrpc

func sliceContains(a []interface{}, b interface{}) bool {
	for i := range a {
		if b == a[i] {
			return true
		}
	}

	return false
}

func mapContains(a map[string]interface{}, k string, v interface{}) bool {
	for ka, va := range a {
		if ka == k && va == v {
			return true
		}
	}

	return false
}
