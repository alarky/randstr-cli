package main

import (
	"math"
	"strings"
	"testing"
)

// helper: generate multiple strings
func generateN(t *testing.T, alphaPool, symbolPool string, count, length, symbolMaxPct int) []string {
	t.Helper()
	results := make([]string, count)
	for i := 0; i < count; i++ {
		s, err := generate(alphaPool, symbolPool, length, symbolMaxPct)
		if err != nil {
			t.Fatalf("generate failed: %v", err)
		}
		results[i] = s
	}
	return results
}

func TestDefaultState(t *testing.T) {
	t.Run("generates strings of specified length", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		results := generateN(t, alpha, sym, 12, 16, 20)
		if len(results) != 12 {
			t.Fatalf("expected 12 strings, got %d", len(results))
		}
		for _, s := range results {
			if len(s) != 16 {
				t.Errorf("expected length 16, got %d: %s", len(s), s)
			}
		}
	})

	t.Run("default results contain only alphanumeric chars", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		results := generateN(t, alpha, sym, 12, 16, 20)
		for _, s := range results {
			for _, c := range s {
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
					t.Errorf("unexpected char %c in %s", c, s)
				}
			}
		}
	})

	t.Run("default pools include lowercase uppercase and digits", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		if !strings.Contains(alpha, lowerChars) {
			t.Error("default should include lowercase")
		}
		if !strings.Contains(alpha, upperChars) {
			t.Error("default should include uppercase")
		}
		if !strings.Contains(alpha, digitChars) {
			t.Error("default should include digits")
		}
		if sym != "" {
			t.Error("default should not include symbols")
		}
	})
}

func TestCountAndLength(t *testing.T) {
	t.Run("changing count updates number of results", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		results := generateN(t, alpha, sym, 5, 16, 20)
		if len(results) != 5 {
			t.Fatalf("expected 5 strings, got %d", len(results))
		}
	})

	t.Run("changing length updates string length", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		results := generateN(t, alpha, sym, 3, 8, 20)
		for _, s := range results {
			if len(s) != 8 {
				t.Errorf("expected length 8, got %d", len(s))
			}
		}
	})
}

func TestCharacterSets(t *testing.T) {
	t.Run("lowercase only", func(t *testing.T) {
		alpha, sym := buildPools(true, false, false, false, "")
		results := generateN(t, alpha, sym, 12, 16, 20)
		for _, s := range results {
			for _, c := range s {
				if c < 'a' || c > 'z' {
					t.Errorf("expected only lowercase, got %c in %s", c, s)
				}
			}
		}
	})

	t.Run("uppercase only", func(t *testing.T) {
		alpha, sym := buildPools(false, true, false, false, "")
		results := generateN(t, alpha, sym, 12, 16, 20)
		for _, s := range results {
			for _, c := range s {
				if c < 'A' || c > 'Z' {
					t.Errorf("expected only uppercase, got %c in %s", c, s)
				}
			}
		}
	})

	t.Run("digits only", func(t *testing.T) {
		alpha, sym := buildPools(false, false, true, false, "")
		results := generateN(t, alpha, sym, 12, 16, 20)
		for _, s := range results {
			for _, c := range s {
				if c < '0' || c > '9' {
					t.Errorf("expected only digits, got %c in %s", c, s)
				}
			}
		}
	})

	t.Run("uppercase and digits", func(t *testing.T) {
		alpha, sym := buildPools(false, true, true, false, "")
		results := generateN(t, alpha, sym, 12, 16, 20)
		for _, s := range results {
			for _, c := range s {
				if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
					t.Errorf("expected only uppercase+digits, got %c in %s", c, s)
				}
			}
		}
	})

	t.Run("all unchecked produces empty strings", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		// Simulate: all char flags off but no symbols either
		// With the default behavior, no flags = all alphanumeric.
		// To truly get "all unchecked", we need at least one flag set to trigger non-default mode
		// but none of a/A/0/s/symbols set. This is the edge case where
		// someone passes -a=false -A=false -0=false explicitly.
		// In practice, buildPools with all false and empty symbols returns default (all alphanumeric).
		// The "all unchecked" case from web doesn't directly map since CLI defaults differ.
		// Instead, test with empty pools directly.
		s, err := generate("", "", 16, 20)
		if err != nil {
			t.Fatal(err)
		}
		if s != "" {
			t.Errorf("expected empty string for empty pools, got %s", s)
		}
		_ = alpha
		_ = sym
	})
}

func TestSymbols(t *testing.T) {
	t.Run("all symbols flag includes all 32 symbols", func(t *testing.T) {
		_, sym := buildPools(false, false, false, true, "")
		if sym != allSymbols {
			t.Errorf("expected all symbols, got %s", sym)
		}
		if len(allSymbols) != 32 {
			t.Errorf("expected 32 symbols, got %d", len(allSymbols))
		}
	})

	t.Run("individual symbols selection", func(t *testing.T) {
		alpha, sym := buildPools(true, false, false, false, "!@#$")
		results := generateN(t, alpha, sym, 20, 100, 100)
		for _, s := range results {
			for _, c := range s {
				if !((c >= 'a' && c <= 'z') || c == '!' || c == '@' || c == '#' || c == '$') {
					t.Errorf("unexpected char %c in %s", c, s)
				}
			}
		}
	})

	t.Run("individual symbols deduplication", func(t *testing.T) {
		_, sym := buildPools(false, false, false, false, "!!@@##")
		if sym != "!@#" {
			t.Errorf("expected deduplicated '!@#', got %s", sym)
		}
	})

	t.Run("invalid symbols are ignored", func(t *testing.T) {
		_, sym := buildPools(false, false, false, false, "abc!@")
		// a, b, c are not in allSymbols, only ! and @ should remain
		if sym != "!@" {
			t.Errorf("expected '!@', got %s", sym)
		}
	})

	t.Run("symbol max limits symbol ratio", func(t *testing.T) {
		alpha, sym := buildPools(true, false, false, true, "")
		results := generateN(t, alpha, sym, 20, 100, 10)
		for _, s := range results {
			symbolCount := 0
			for _, c := range s {
				if strings.ContainsRune(allSymbols, c) {
					symbolCount++
				}
			}
			if symbolCount > 10 { // 100 * 10% = 10
				t.Errorf("symbol count %d exceeds max 10 in %s", symbolCount, s)
			}
		}
	})

	t.Run("symbol max with only symbols has no limit", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, true, "")
		// When only symbols, alpha pool is empty (since -s is a symbol flag, not a char flag)
		// buildPools with s=true and no l/u/d flags: noCharFlags is false because s is true
		// so alphaPool is empty, symbolPool = allSymbols
		results := generateN(t, alpha, sym, 5, 16, 10)
		for _, s := range results {
			if len(s) != 16 {
				t.Errorf("expected length 16, got %d", len(s))
			}
			// All chars should be symbols
			for _, c := range s {
				if !strings.ContainsRune(allSymbols, c) {
					t.Errorf("expected only symbols, got %c", c)
				}
			}
		}
	})

	t.Run("symbol max with only alphanumeric has no limit", func(t *testing.T) {
		alpha, sym := buildPools(true, true, true, false, "")
		results := generateN(t, alpha, sym, 5, 16, 10)
		for _, s := range results {
			for _, c := range s {
				if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
					t.Errorf("expected only alphanumeric, got %c", c)
				}
			}
		}
	})
}

func TestBuildPools(t *testing.T) {
	t.Run("s flag overrides individual symbols", func(t *testing.T) {
		_, sym := buildPools(false, false, false, true, "!@#")
		if sym != allSymbols {
			t.Errorf("expected all symbols when -s is set, got %s", sym)
		}
	})

	t.Run("lowercase and symbols combined", func(t *testing.T) {
		alpha, sym := buildPools(true, false, false, false, "!@")
		if alpha != lowerChars {
			t.Errorf("expected lowercase chars, got %s", alpha)
		}
		if sym != "!@" {
			t.Errorf("expected '!@', got %s", sym)
		}
	})
}

// Chi-squared test helper
func chiSquared(observed map[byte]int, expected float64) float64 {
	sum := 0.0
	for _, count := range observed {
		diff := float64(count) - expected
		sum += diff * diff / expected
	}
	return sum
}

func TestRandomnessDistribution(t *testing.T) {
	t.Run("alphanumeric characters are uniformly distributed", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		results := generateN(t, alpha, sym, 500, 200, 20)

		freq := make(map[byte]int)
		for _, s := range results {
			for i := 0; i < len(s); i++ {
				freq[s[i]]++
			}
		}

		totalChars := 500 * 200
		poolSize := 62 // 26 + 26 + 10
		expectedPerChar := float64(totalChars) / float64(poolSize)

		if len(freq) != poolSize {
			t.Errorf("expected %d distinct chars, got %d", poolSize, len(freq))
		}

		// df=61, critical value at p=0.001 is ~99.6
		chi2 := chiSquared(freq, expectedPerChar)
		if chi2 >= 100 {
			t.Errorf("chi-squared %f exceeds critical value 100 (distribution not uniform)", chi2)
		}
	})

	t.Run("lowercase-only characters are uniformly distributed", func(t *testing.T) {
		alpha, sym := buildPools(true, false, false, false, "")
		results := generateN(t, alpha, sym, 500, 200, 20)

		freq := make(map[byte]int)
		for _, s := range results {
			for i := 0; i < len(s); i++ {
				freq[s[i]]++
			}
		}

		totalChars := 500 * 200
		poolSize := 26
		expectedPerChar := float64(totalChars) / float64(poolSize)

		if len(freq) != poolSize {
			t.Errorf("expected %d distinct chars, got %d", poolSize, len(freq))
		}

		// df=25, critical value at p=0.001 is ~52.6
		chi2 := chiSquared(freq, expectedPerChar)
		if chi2 >= 53 {
			t.Errorf("chi-squared %f exceeds critical value 53", chi2)
		}
	})

	t.Run("with symbols enabled distribution is reasonable", func(t *testing.T) {
		alpha, sym := buildPools(true, true, true, true, "")
		results := generateN(t, alpha, sym, 500, 200, 100) // symbolMax=100 to disable limiting

		freq := make(map[byte]int)
		for _, s := range results {
			for i := 0; i < len(s); i++ {
				freq[s[i]]++
			}
		}

		totalChars := 500 * 200
		actualPoolSize := len(freq)
		expectedPerChar := float64(totalChars) / float64(actualPoolSize)

		// 62 alphanumeric + 32 symbols = 94
		if actualPoolSize < 90 {
			t.Errorf("expected at least 90 distinct chars, got %d", actualPoolSize)
		}

		// df~94, critical value at p=0.001 is ~135
		chi2 := chiSquared(freq, expectedPerChar)
		if chi2 >= 135 {
			t.Errorf("chi-squared %f exceeds critical value 135", chi2)
		}
	})

	t.Run("generated strings are unique", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		results := generateN(t, alpha, sym, 10000, 32, 20)

		seen := make(map[string]bool)
		for _, s := range results {
			if seen[s] {
				t.Errorf("duplicate string found: %s", s)
			}
			seen[s] = true
		}
	})
}

func TestSymbolMaxPercentages(t *testing.T) {
	cases := []struct {
		name      string
		length    int
		maxPct    int
		maxSymbol int
	}{
		{"20% of 100", 100, 20, 20},
		{"50% of 100", 100, 50, 50},
		{"10% of 50", 50, 10, 5},
		{"0% means no symbols", 100, 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			alpha, sym := buildPools(true, false, false, true, "")
			results := generateN(t, alpha, sym, 50, tc.length, tc.maxPct)
			for _, s := range results {
				symbolCount := 0
				for _, c := range s {
					if strings.ContainsRune(allSymbols, c) {
						symbolCount++
					}
				}
				if symbolCount > tc.maxSymbol {
					t.Errorf("symbol count %d exceeds max %d in %s", symbolCount, tc.maxSymbol, s)
				}
			}
		})
	}
}

func TestGenerateEdgeCases(t *testing.T) {
	t.Run("length 1", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		s, err := generate(alpha, sym, 1, 20)
		if err != nil {
			t.Fatal(err)
		}
		if len(s) != 1 {
			t.Errorf("expected length 1, got %d", len(s))
		}
	})

	t.Run("very long string", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "")
		s, err := generate(alpha, sym, 1000, 20)
		if err != nil {
			t.Fatal(err)
		}
		if len(s) != 1000 {
			t.Errorf("expected length 1000, got %d", len(s))
		}
	})

	t.Run("single symbol in pool", func(t *testing.T) {
		alpha, sym := buildPools(false, false, false, false, "!")
		// noCharFlags is false because symbols != "", so alphaPool is empty
		// Actually: noCharFlags = !l && !u && !d && !s && symbols == ""
		// symbols is "!", so noCharFlags = false → alphaPool is empty
		results := generateN(t, alpha, sym, 5, 10, 20)
		for _, s := range results {
			for _, c := range s {
				if c != '!' {
					t.Errorf("expected only '!', got %c", c)
				}
			}
		}
	})
}

// Ensure math import is used (for potential future use)
var _ = math.Sqrt
