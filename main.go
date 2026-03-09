package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"
)

const (
	lowerChars = "abcdefghijklmnopqrstuvwxyz"
	upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars = "0123456789"
	allSymbols = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	c := flag.Int("c", 12, "number of strings to generate")
	n := flag.Int("n", 16, "length of each string")
	a := flag.Bool("a", false, "include lowercase letters (a-z)")
	A := flag.Bool("A", false, "include uppercase letters (A-Z)")
	zero := flag.Bool("0", false, "include digits (0-9)")
	s := flag.Bool("s", false, "include all symbols")
	symbols := flag.String("symbols", "", "include specific symbols (e.g. '!@#$')")
	symbolMax := flag.Int("symbol-max", 20, "max percentage of symbols in output")
	flag.Parse()

	if *c <= 0 {
		return fmt.Errorf("count must be positive")
	}
	if *n <= 0 {
		return fmt.Errorf("length must be positive")
	}

	alphaPool, symbolPool := buildPools(*a, *A, *zero, *s, *symbols)

	for i := 0; i < *c; i++ {
		result, err := generate(alphaPool, symbolPool, *n, *symbolMax)
		if err != nil {
			return err
		}
		fmt.Println(result)
	}
	return nil
}

func buildPools(l, u, d, s bool, symbols string) (string, string) {
	noCharFlags := !l && !u && !d && !s && symbols == ""

	var alphaPool string
	if l || noCharFlags {
		alphaPool += lowerChars
	}
	if u || noCharFlags {
		alphaPool += upperChars
	}
	if d || noCharFlags {
		alphaPool += digitChars
	}

	var symbolPool string
	if s {
		symbolPool = allSymbols
	} else if symbols != "" {
		seen := make(map[byte]bool)
		for i := 0; i < len(symbols); i++ {
			ch := symbols[i]
			if strings.IndexByte(allSymbols, ch) != -1 && !seen[ch] {
				symbolPool += string(ch)
				seen[ch] = true
			}
		}
	}

	return alphaPool, symbolPool
}

func generate(alphaPool, symbolPool string, length int, symbolMaxPct int) (string, error) {
	if alphaPool == "" && symbolPool == "" {
		return "", nil
	}

	maxSymbols := length
	if alphaPool != "" && symbolPool != "" {
		maxSymbols = length * symbolMaxPct / 100
	}

	fullPool := alphaPool + symbolPool
	result := make([]byte, length)
	symbolCount := 0

	for i := range result {
		pool := fullPool
		if symbolCount >= maxSymbols && alphaPool != "" {
			pool = alphaPool
		}

		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(pool))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		ch := pool[idx.Int64()]
		result[i] = ch
		if strings.IndexByte(symbolPool, ch) != -1 {
			symbolCount++
		}
	}
	return string(result), nil
}
