package termooo

import (
	"bufio"
	"math/rand"
	"os"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func normalizeWord(word string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, word)

	return result
}

func readFile(filepath string) ([]string, error) {
	var lines []string

	f, err := os.Open(filepath)
	if err != nil {
		return lines, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, normalizeWord(scanner.Text()))
	}

	return lines, nil
}

func selectRandomFromSlice(slice []string) string {
	rand.Seed(time.Now().Unix())
	randomIndex := rand.Intn(len(slice))

	return slice[randomIndex]
}
