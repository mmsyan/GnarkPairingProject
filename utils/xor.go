package utils

func Xor(a, b []byte) []byte {
	var length int
	if len(a) > len(b) {
		length = len(b)
	} else {
		length = len(a)
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = a[i] ^ b[i]
	}
	return result
}
