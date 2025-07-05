package utils

import "unicode"

func YankFromStrings(lines []string, startX, startY, endX, endY int) []string {
	var selectedText []string

	if (endY < startY) || (endY == startY && endX < startX) {
		startX, endX = endX, startX
		startY, endY = endY, startY
	}

	for y := startY; y <= endY && y < len(lines); y++ {
		line := lines[y]

		if y == startY && y == endY {
			selectedText = append(selectedText, line[startX:endX])
		} else if y == startY {
			selectedText = append(selectedText, line[startX:])
		} else if y == endY {
			selectedText = append(selectedText, line[:endX])
		} else {
			selectedText = append(selectedText, line)
		}
	}

	return selectedText
}

func IsWordChar(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_'
}

func MoveToNextWord(s string, pos int) int {
	runes := []rune(s)
	n := len(runes)

	for pos < n && IsWordChar(runes[pos]) {
		pos++
	}

	for pos < n && !IsWordChar(runes[pos]) {
		pos++
	}

	return pos
}

func MoveToPrevWord(s string, pos int) int {
	runes := []rune(s)

	for pos > 0 && IsWordChar(runes[pos-1]) {
		pos--
	}

	for pos > 0 && !IsWordChar(runes[pos-1]) {
		pos--
	}

	return pos
}

func MoveToNextRune(s string, pos int, r rune) int {
	runes := []rune(s)
	n := len(runes)

	for pos < n && runes[pos] == r {
		pos++
	}

	for pos < n && runes[pos] != r {
		pos++
	}

	return pos
}

func MoveToPrevRune(s string, pos int, r rune) int {
	runes := []rune(s)

	for pos > 0 && runes[pos-1] == r {
		pos--
	}

	for pos > 0 && runes[pos-1] != r {
		pos--
	}

	return pos
}
