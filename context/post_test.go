package context

import "testing"

func Test_splitFilnameMatcher(t *testing.T) {
	var tests = []struct {
		pattern string
		match   bool
	}{
		{"2044-02-04-My_post", true},
		{"20424-02-04-My_post", false},
		{"2044-02-4-My_post", false},
		{"s2044-02-04-My_post", false},
		{"20d44-02-04-My_post", false},
		{"2044-02-4", false},
		{"2044-02-04 asdf", true},
	}

	for _, f := range tests {
		_, _, err := splitFilname(f.pattern)
		if err != nil && f.match == true {
			t.Errorf("date: %s should not match\n", f.pattern)
		}
	}
}

func Test_splitFilnameTitelizer(t *testing.T) {
	var tests = []struct {
		pattern string
		title   string
	}{
		{"2044-02-04-My_post.md", "My_post"},
		{"2044-02-04 my_post.markdown", "My_post"},
		{"2044-02-04my_post     .mkdn", "My_post"},
		{"2044-02-04        my_post-sub.mkd", "My_post-Sub"},
		{"2044-02-04    -my_post-sub-.md", "My_post-Sub"},
		{"2044-02-04 --        my_post-sub.mkdown", "My_post-Sub"},
	}

	for _, f := range tests {
		_, title, err := splitFilname(f.pattern)
		if err != nil || f.title != title {
			t.Errorf("correct: %s got %s\n", f.title, title)
		}
	}
}
