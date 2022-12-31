package app

var bigLetters = map[rune]string{
	'a': "Ａ",
	'b': "Ｂ",
	'c': "Ｃ",
	'd': "Ｄ",
	'e': "Ｅ",
	'f': "Ｆ",
	'g': "Ｇ",
	'h': "Ｈ",
	'i': "Ｉ",
	'j': "Ｊ",
	'k': "Ｋ",
	'l': "Ｌ",
	'm': "Ｍ",
	'n': "Ｎ",
	'o': "Ｏ",
	'p': "Ｐ",
	'q': "Ｑ",
	'r': "Ｒ",
	's': "Ｓ",
	't': "Ｔ",
	'u': "Ｕ",
	'v': "Ｖ",
	'w': "Ｗ",
	'x': "Ｘ",
	'y': "Ｙ",
	'z': "Ｚ",
}

func turnBig(input string) (result string) {
	for _, char := range input {
		if value, ok := bigLetters[char]; ok {
			result += value
		} else {
			result += string(char)
		}

	}
	return result
}
