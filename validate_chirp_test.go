package main

import "testing"

func Test_profanity(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		s    string
		want string
	}{
		{"kerfuffle", "this kerfuFFles", "this kerfuFFles"},
		{"sharbert", "sharbert this", "**** this"},
		{"fornax", "sharbert this FORNAX", "**** this ****"},
		{"multiple words", "kerfuffle kerfuffle", "**** ****"},
		{"mixed case", "KeRfUfFlE", "****"},
		{"with punctuation", "sharbert!", "sharbert!"}, // shouldn't match
		{"middle of string", "hello kerfuffle world", "hello **** world"},
		{"all three words", "fornax sharbert kerfuffle", "**** **** ****"},

		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := profanity(tt.s)
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("profanity() = %v, want %v", got, tt.want)
			}
		})
	}
}
