package utils

const (
	IDR = "IDR"
	USD = "USD"
)

func IsCurrencySupported(curr string) bool {
	switch curr {
	case IDR, USD:
		return true
	}

	return false
}
