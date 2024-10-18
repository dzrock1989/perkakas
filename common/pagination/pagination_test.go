package pagination

import "testing"

func TestIsSortSave(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
		given    string
	}{
		{"empty", false, ""},
		{"SQLi", false, "name; drop table users;"},
		{"save 1", true, "name"},
		{"save 2", true, "Name"},
		{"save 3", true, "az asc"},
		{"save asc", true, "name asc"},
		{"save desc", true, "name desc"},
		{"underscore save", true, "full_name asc"},
		{"underscore not save", false, "full_name; asc"},
		{"not save asc typo", false, "name ascc"},
		{"not save desc typo", false, "name descs"},
		{"not save desc", false, "name a desc"},
		{"not save asc", false, "name a asc"},
		{"unknown char", false, "name; asc"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual := IsSortSave(tt.given)
			if actual != tt.expected {
				t.Errorf("(%s): expected %v, actual %v", tt.given, tt.expected, actual)
			}
		})
	}
}
