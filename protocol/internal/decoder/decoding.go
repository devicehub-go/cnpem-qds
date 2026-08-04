package decoder

import (
	"strconv"
	"strings"
)

type DecodeOptions struct {
	Encoded   string
	ChunkSize int
	LeastStep float64
}

var DecodingMap = map[rune]string{
	'G': "00", 'H': "01", 'I': "02", 'J': "03",
	'K': "04", 'L': "05", 'M': "06", 'N': "07",
	'O': "08", 'P': "09", 'Q': "0A", 'R': "0B",
	'S': "0C", 'T': "0D", 'U': "0E", 'V': "0F",
	'W': "10", 'X': "20", 'Y': "30", 'Z': "40",
	'a': "50", 'b': "60", 'c': "70", 'd': "80",
	'e': "90", 'f': "A0", 'g': "B0", 'h': "C0",
	'i': "D0", 'j': "E0", 'k': "F0", 'l': "F1",
	'm': "F2", 'n': "F3", 'o': "F4", 'p': "F5",
	'q': "F6", 'r': "F7", 's': "F8", 't': "F9",
	'u': "FA", 'v': "FB", 'w': "FC", 'x': "FD",
	'y': "FE", 'z': "FF", '+': "1F", ';': "2F",
	'=': "3F", '<': "4F", '@': "5F", '#': "6F",
	'$': "7F", '%': "8F", '^': "9F", '&': "AF",
	'*': "BF", '?': "CF", '>': "DF", '-': "EF",
}

// Decompress a message according to decoding map
func decompress(compressed string) string {
	var builder strings.Builder
	builder.Grow(len(compressed) * 2)

	for _, symbol := range compressed {
		if decoded, ok := DecodingMap[symbol]; ok {
			builder.WriteString(decoded)
		} else {
			builder.WriteRune(symbol)
		}
	}
	return builder.String()
}

// 18-bit two's complement conversion using bitwise shifts
func hexToSample(hexStr string, leastStep float64) float64 {
	if hexStr == "" {
		return 0.0
	}
	val, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		return 0.0
	}
	signedVal := int64(val<<(64-18)) >> (64 - 18)
	return float64(signedVal) * leastStep
}

// Converts a QDS string channel content to float array
func Decode(options DecodeOptions) []float64 {
	var samples []float64

	if options.ChunkSize <= 0 {
		options.ChunkSize = 5
	}
	if options.LeastStep <= 0 {
		options.LeastStep = 1.0
	}

	decompressed := decompress(options.Encoded)
	for i := 0; i < len(decompressed); i += options.ChunkSize {
		end := min(i+options.ChunkSize, len(decompressed))
		value := hexToSample(decompressed[i:end], options.LeastStep)
		samples = append(samples, value)
	}

	return samples
}

// Rounds a number to next multiple of
func RoundMultiple(number int, multiple int) int {
	return ((number + multiple - 1) / multiple) * multiple
}
