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

func AreKeySetsMatching(existingKeys []string, expectedKeys map[string]bool) bool {
	// Make sure all required keys exist.
	for k, expected := range expectedKeys {
		// If this key is expected, look for it in existingKeys
		if expected {
			found := false
			for _, existing := range existingKeys {
				if k == existing {
					found = true
					break
				}
			}
			// If we didn't find it, then we already know keysets aren't matching,
			// so we return false.
			if !found {
				return false
			}
		}
	}

	// Make sure we don't have any extra, unexpected keys.
	for _, existing := range existingKeys {
		// If there is a mapping whose key is `existing`, then regardless of whether
		// it is required or not, it's valid, so we allow it. Only if the key is
		// invalid (i.e. not required and not even an optional field) do we return
		// false.
		if _, ok := expectedKeys[existing]; !ok { // here we are checking if the key `existing` is in the map `expectedKeys`
			return false
		}
	}

	// If we managed to reach this point, then the keysets are matching.
	return true
}
