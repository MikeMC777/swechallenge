package repo

import "testing"

// Test table-driven para parseMoney. Cubre valores con símbolo $, comas,
// vacíos/espacios y casos límite.
func TestParseMoney(t *testing.T) {
	cases := []struct {
		in   string
		want *float64
	}{
		{"$525.00", f(525.00)},
		{"$1,234.50", f(1234.50)},
		{"$0.00", f(0)},
		{"", nil},
		{" ", nil},
		{"$", nil},
		{"$1,234", f(1234)},
		{"123.45", f(123.45)}, // por si llegan sin símbolo
	}

	for _, tc := range cases {
		got, err := parseMoney(tc.in)
		if err != nil {
			t.Fatalf("parseMoney(%q) error: %v", tc.in, err)
		}
		if (got == nil) != (tc.want == nil) {
			t.Fatalf("parseMoney(%q) nil-mismatch: got=%v want=%v", tc.in, got, tc.want)
		}
		if got != nil && *got != *tc.want {
			t.Fatalf("parseMoney(%q) = %v; want %v", tc.in, *got, *tc.want)
		}
	}
}

func TestParseMoney_Signed(t *testing.T) {
	cases := []struct {
		in   string
		want *float64
	}{
		{"+$1,000.00", f(1000)},
		{"-$12.34", f(-12.34)},
	}
	for _, tc := range cases {
		got, err := parseMoney(tc.in)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
		if (got == nil) != (tc.want == nil) {
			t.Fatalf("nil mismatch for %q", tc.in)
		}
		if got != nil && *got != *tc.want {
			t.Fatalf("want %v got %v", *tc.want, *got)
		}
	}
}

func TestParseMoney_Invalid(t *testing.T) {
	// Caso realmente inválido: separador mal puesto
	if v, err := parseMoney("$12,34.56"); err == nil {
		t.Fatalf("expected error, got %v", v)
	}
}

func f(v float64) *float64 { return &v }
