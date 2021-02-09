package utils

import "encoding/base64"

// Base64Decode takes base64-encoded string and returns its original value
func Base64Decode(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}

	decoded := string(data)
	return decoded, nil
}
