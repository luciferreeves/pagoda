package validators

const MaxJade uint64 = 99999

func IsValidJade(amount uint64) bool {
	return amount <= MaxJade
}