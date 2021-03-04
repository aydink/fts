package fts

type Token struct {
	start, end, position uint32
	value                string
}

type Tokenizer interface {
	Tokenize(string) []Token
}

type TokenFilterer interface {
	Filter([]Token) []Token
}

type Analyzer interface {
	Analyze(string) []Token
}

type SimpleAnalyzer struct {
	tokenizer    Tokenizer
	tokenFilters []TokenFilterer
}

func NewSimpleAnalyzer(t Tokenizer) *SimpleAnalyzer {
	return &SimpleAnalyzer{t, make([]TokenFilterer, 0)}
}

func (sa *SimpleAnalyzer) AddTokenFilter(f TokenFilterer) {
	sa.tokenFilters = append(sa.tokenFilters, f)
}

func (sa *SimpleAnalyzer) Analyze(s string) []Token {
	t := sa.tokenizer.Tokenize(s)
	for _, tf := range sa.tokenFilters {
		t = tf.Filter(t)
	}
	return t
}

var p = map[byte]bool{'\'': true, ',': true, '.': true, ':': true, ';': true, '!': true, '?': true, '(': true, ')': true, '"': true, ' ': true, '\t': true, '\n': true, '\r': true, '|': true, '\\': true, '/': true}

type SimpleTokenizer struct{}

func NewSimpleTokenizer() SimpleTokenizer {
	return SimpleTokenizer{}
}

func (tk SimpleTokenizer) Tokenize(s string) []Token {
	var posToken uint32 = 0

	i := 0
	tokens := []Token{}
	token := Token{start: 0, end: 0, position: 0, value: ""}

	for i < len(s) {

		for (i < len(s)) && p[s[i]] {
			i++
		}

		token.start = uint32(i)

		for (i < len(s)) && !p[s[i]] {
			i++
		}
		token.end = uint32(i)
		token.value = s[token.start:token.end]

		// handle zero length tokens
		if token.start != token.end {
			token.position = posToken
			posToken++
			tokens = append(tokens, token)
		}
	}

	return tokens
}
