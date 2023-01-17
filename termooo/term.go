package termooo

import (
	"regexp"
	"strings"
)

type letterStatus struct {
	letter string
	color  string
}

type word struct {
	guess  string
	status map[int]letterStatus
}

func (w *word) applyGreen(target string) map[rune]int {
	tFreq := frequency(target)

	for index, letter := range w.guess {
		if letter == rune(target[index]) {

			tFreq[letter]--
			w.status[index] = letterStatus{color: "Green", letter: string(letter)}
		}
	}

	return tFreq
}

func (w *word) applyYellow(targetFreq map[rune]int, target string) {
	for index, letter := range w.guess {
		if w.status[index].color != "Green" &&
			strings.Contains(target, string(letter)) &&
			targetFreq[letter] != 0 {

			targetFreq[letter]--
			w.status[index] = letterStatus{color: "Yellow", letter: string(letter)}
		}
	}
}

func (w *word) applyGray(target string) {
	for index, letter := range w.guess {
		if w.status[index].color != "Green" && w.status[index].color != "Yellow" {
			w.status[index] = letterStatus{color: "Gray", letter: string(letter)}
		}
	}
}

func (w *word) applyAll(target string) {
	freq := w.applyGreen(target)
	w.applyYellow(freq, target)
	w.applyGray(target)

}

func frequency(target string) map[rune]int {
	freq := make(map[rune]int)
	for _, value := range target {
		freq[value] = freq[value] + 1
	}

	return freq
}

func isAlpha(s string) bool {
	return regexp.MustCompile(`^[A-Za-zÀ-ÿ]+$`).MatchString(s)
}

func generateTarget() (string, error) {
	lines, err := readFile("words.txt")
	if err != nil {
		return "", err
	}

	element := selectRandomFromSlice(lines)

	return strings.ToLower(element), nil
}
