package utils

// Rounds a number to next multiple of
func RoundMultiple(number int, multiple int) int {
	return ((number + multiple - 1) / multiple) * multiple
}
