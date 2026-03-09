package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
)

const (
	lowerChars  = "abcdefghijklmnopqrstuvwxyz"
	upperChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digitChars  = "0123456789"
	symbolChars = "!@#$%^&*()-_=+[]{}|;:,.<>?"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	n := flag.Int("n", 16, "length of password")
	l := flag.Bool("l", false, "include lowercase letters")
	u := flag.Bool("u", false, "include uppercase letters")
	d := flag.Bool("d", false, "include digits")
	s := flag.Bool("s", false, "include symbols")
	flag.Parse()

	if *n <= 0 {
		return fmt.Errorf("length must be positive")
	}

	charset := buildCharset(*l, *u, *d, *s)

	password, err := generate(charset, *n)
	if err != nil {
		return err
	}

	fmt.Println(password)
	return nil
}

func buildCharset(l, u, d, s bool) string {
	// フラグ未指定なら全部入り
	if !l && !u && !d && !s {
		return lowerChars + upperChars + digitChars + symbolChars
	}

	var charset string
	if l {
		charset += lowerChars
	}
	if u {
		charset += upperChars
	}
	if d {
		charset += digitChars
	}
	if s {
		charset += symbolChars
	}
	return charset
}

func generate(charset string, length int) (string, error) {
	result := make([]byte, length)
	for i := range result {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		result[i] = charset[idx.Int64()]
	}
	return string(result), nil
}
