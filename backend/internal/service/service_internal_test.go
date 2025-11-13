package service

import "testing"

func TestServiceSanitizeEmail(t *testing.T) {
	svc := &Service{}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "lowercase and trim",
			input: "  USER@Example.COM  ",
			want:  "user@example.com",
		},
		{
			name:  "empty",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		if got := svc.sanitizeEmail(tt.input); got != tt.want {
			t.Fatalf("%s: sanitizeEmail(%q) = %q, want %q", tt.name, tt.input, got, tt.want)
		}
	}
}

func TestServiceClampPagination(t *testing.T) {
	svc := &Service{}

	tests := []struct {
		name             string
		page, perPage    int
		wantPage, wantPP int
	}{
		{
			name:     "defaults when values are zero",
			page:     0,
			perPage:  0,
			wantPage: 1,
			wantPP:   20,
		},
		{
			name:     "clamp negative values",
			page:     -5,
			perPage:  -10,
			wantPage: 1,
			wantPP:   20,
		},
		{
			name:     "preserve valid values",
			page:     2,
			perPage:  50,
			wantPage: 2,
			wantPP:   50,
		},
		{
			name:     "cap per page at 100",
			page:     1,
			perPage:  200,
			wantPage: 1,
			wantPP:   100,
		},
	}

	for _, tt := range tests {
		gotPage, gotPP := svc.clampPagination(tt.page, tt.perPage)
		if gotPage != tt.wantPage || gotPP != tt.wantPP {
			t.Fatalf("%s: clampPagination(%d, %d) = (%d, %d), want (%d, %d)", tt.name, tt.page, tt.perPage, gotPage, gotPP, tt.wantPage, tt.wantPP)
		}
	}
}
