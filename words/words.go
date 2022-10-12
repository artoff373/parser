package words

import (
	"math"
	"strings"
	"unicode"
	"unicode/utf8"
)

const porog = 0.62

type Branch struct {
	letter rune
	end    bool
	posle  []*Branch
}

// Построение дерева
func Tree(keys []string) *Branch {
	node := new(Branch)
	node.letter, _ = utf8.DecodeRune([]byte("."))
	node.end = false
	for i := range keys {
		p := node
		buildATree([]byte(keys[i]), p)
	}
	return node
}

// Построение ветвей дерева
func buildATree(word []byte, p *Branch) {
	r, width := utf8.DecodeRune(word[0:])
	flag := true
	ostatok := word[width:]
	newLetter := new(Branch)
	for i := range p.posle {
		if r == p.posle[i].letter {
			flag = false
			newLetter = p.posle[i]
		}
	}
	if flag {
		p.posle = append(p.posle, newLetter)
		newLetter.letter = r
		newLetter.end = false
	}
	if len(ostatok) > 0 {
		buildATree(ostatok, newLetter)
	} else {
		newLetter.end = true
	}
}

// Обход по дереву
func throughTree(p *Branch, word []byte, deep int) (int, bool) {
	r, width := utf8.DecodeRune(word[0:])
	ostatok := word[width:]
	end := p.end
	for i := range p.posle {
		if p.posle[i].letter == r {
			deep++
			deep, end = throughTree(p.posle[i], ostatok, deep)
		}
	}
	return deep, end
}

// Поиск ключей в дереве
func SearchKeys(keys *Branch, dictionary []string) float64 {
	deep := 0
	relev := 0.0
	for i := range dictionary {
		deep, end := throughTree(keys, []byte(dictionary[i]), deep)
		wordLen := utf8.RuneCountInString(dictionary[i])
		if deep != 0 {
			perc := float64(deep) / float64(wordLen)
			switch {
			case perc == 1 && !end:
				continue
			case perc > porog || (math.Abs(float64(wordLen-deep)) < 3 && wordLen > 3):
				relev += perc
			}
		}
	}
	return relev
}

// Функция отсеивания пробелов и знаков пунктуации в строке
func FindWords(data []byte, atEOF bool) (token string) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) && !unicode.IsPunct(r) {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) || unicode.IsPunct(r) {
			return string(data[start:i])
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return string(data[start:])
	}
	// Request more data.
	return ""
}

// Сортировка строки с отсеиванием повторяющихся слов и слов меньше 3 знаков длиной
func SortUniq(massiv []string) []string {
	if len(massiv) < 2 {
		if len(massiv) == 2 && strings.Compare(massiv[0], massiv[1]) == -1 {
			massiv[0], massiv[1] = massiv[1], massiv[0]
		}
		return massiv
	}
	opora := FindWords([]byte(massiv[len(massiv)/2]), true)
	var left, right []string
	left, right = nil, nil
	for i := range massiv {
		massiv[i] = FindWords([]byte(massiv[i]), true)
		if (utf8.RuneCount([]byte(massiv[i]))) < 3 {
			continue
		}
		switch strings.Compare(massiv[i], opora) {
		case -1:
			left = append(left, massiv[i])
		case 0:
			continue
		case 1:
			right = append(right, massiv[i])
		}
	}
	var slice []string
	if (utf8.RuneCount([]byte(opora))) < 3 {
		slice = append(SortUniq(left), SortUniq(right)...)
	} else {
		slice = append(append(SortUniq(left), opora), SortUniq(right)...)
	}

	return slice
}
