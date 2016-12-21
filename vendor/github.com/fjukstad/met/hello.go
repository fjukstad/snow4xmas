package met

func SecureHello() (string, error) {
	endpoint := baseUrl + "/tests/secureHello"
	body, err := get(endpoint)
	if err != nil {
		return "", err
	}

	return string(body), nil

}
