package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDomain(t *testing.T) {
	cases := map[string]string{
		"":                         "",
		"example.com":              "example.com",
		"localhost":                "localhost",
		"foo.localhost":            "localhost",
		"foo.localhost:8080":       "localhost",
		"foo.example.com":          "example.com",
		"bar.foo.example.com":      "foo.example.com",
		"blah.bar.foo.example.com": "bar.foo.example.com",
	}
	for host, sub := range cases {
		assert.Equal(t, sub, extractDomain(host))
	}
}
