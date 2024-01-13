package srv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	var prev string
	for i := 0; i < 100; i++ {
		curr := newID()
		if prev == curr {
			t.Fatal("generated ids are not unique")
		}
		t.Logf("[%d] %s", i, curr)
	}
}

func TestSubdomain(t *testing.T) {
	cases := map[string]string{
		"example.com":         "",
		"foo.localhost":       "foo",
		"foo.localhost:8080":  "foo",
		"foo.example.com":     "foo",
		"foo.bar.example.com": "foo",
		"foo-bar.example.com": "foo-bar",
	}
	for host, sub := range cases {
		assert.Equal(t, sub, subdomain(host))
	}
}
