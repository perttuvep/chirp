package auth_test

import (
	"net/http"
	"testing"

	"github.com/perttuvep/chirp/internal/auth"
)

func TestHashPass(t *testing.T) {
	tests := []struct {
		name    string
		pw      string
		wantErr bool
	}{
		{
			name:    "valid password",
			pw:      "password123",
			wantErr: false,
		},
		{
			name:    "empty password",
			pw:      "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.HashPass(tt.pw)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("HashPass() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			if gotErr == nil {
				if len(got) == 0 {
					t.Errorf("HashPass() returned empty hash for pw: %v", tt.pw)
				}
				// Indirectly test by checking if CheckPwHash works with it.
				checkErr := auth.CheckPwHash(tt.pw, got)
				if checkErr != nil {
					t.Errorf("HashPass() generated a hash that CheckPwHash() rejected: %v", checkErr)
				}
			}
		})
	}
}

func TestCheckPwHash(t *testing.T) {
	password := "securePassword"
	hash, err := auth.HashPass(password)
	if err != nil {
		t.Fatalf("Failed to generate hash: %v", err)
	}

	tests := []struct {
		name    string
		pw      string
		hash    string
		wantErr bool
	}{
		{
			name:    "valid password and hash",
			pw:      password,
			hash:    hash,
			wantErr: false,
		},
		{
			name:    "invalid password",
			pw:      "wrongPassword",
			hash:    hash,
			wantErr: true,
		},
		{
			name:    "invalid hash",
			pw:      password,
			hash:    "invalidHash",
			wantErr: true,
		},
		{
			name:    "empty password",
			pw:      "",
			hash:    hash,
			wantErr: true,
		},
		{
			name:    "empty hash",
			pw:      password,
			hash:    "",
			wantErr: true,
		},
		{
			name:    "empty password and hash",
			pw:      "",
			hash:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := auth.CheckPwHash(tt.pw, tt.hash)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("CheckPwHash() error = %v, wantErr %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		header  http.Header
		want    string
		wantErr bool
	}{
		{
			name:    "Valid Bearer Token",
			header:  http.Header{"Authorization": []string{"Bearer mytoken123"}},
			want:    "mytoken123",
			wantErr: false,
		},
		{
			name:    "Missing Authorization Header",
			header:  http.Header{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Invalid Prefix",
			header:  http.Header{"Authorization": []string{"Basic mytoken123"}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Short Header",
			header:  http.Header{"Authorization": []string{"Bearer"}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Bearer with extra spaces",
			header:  http.Header{"Authorization": []string{"Bearer    mytoken"}},
			want:    "   mytoken",
			wantErr: false,
		},
		{
			name:    "Bearer with tabs",
			header:  http.Header{"Authorization": []string{"Bearer\tmytoken"}},
			want:    "\tmytoken",
			wantErr: false,
		},
		{
			name:    "Bearer with leading and trailing space",
			header:  http.Header{"Authorization": []string{" Bearer mytoken "}},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Bearer with leading space",
			header:  http.Header{"Authorization": []string{" Bearer mytoken"}},
			want:    "",
			wantErr: true,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := auth.GetBearerToken(tt.header)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetBearerToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetBearerToken() succeeded unexpectedly")
			}
			if got != tt.want {
				t.Errorf("GetBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
