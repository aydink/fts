package fts

import (
	"strings"
	"unicode"
)

// TurkishLowercaseFilter lowercases all tokens respecting Tukish specila lowercase rules like "İ"->"i", "I"->"ı"
type turkishLowercaseFilter struct{}

func NewTurkishLowercaseFilter() TokenFilterer {
	filter := turkishLowercaseFilter{}
	return filter
}

func (tf turkishLowercaseFilter) Filter(tokens []Token) []Token {
	for i := range tokens {
		tokens[i].value = strings.ToLowerSpecial(unicode.TurkishCase, tokens[i].value)
	}
	return tokens
}

type turkishAccentFilter struct{}

func NewTurkishAccentFilter() TokenFilterer {
	filter := turkishAccentFilter{}
	return filter
}

func (tf turkishAccentFilter) Filter(tokens []Token) []Token {
	replacer := strings.NewReplacer("â", "a", "î", "i", "û", "u", "Â", "A", "Î", "İ", "Û", "U")
	for i := range tokens {
		tokens[i].value = replacer.Replace(tokens[i].value)
	}
	return tokens
}
