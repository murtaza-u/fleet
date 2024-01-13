package srv

import (
	"fmt"
	"strings"
	"time"
)

// newID generates a time-based unique ID.
//
// Example: 1705072430052621932
func newID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// subdomain extracts the subdomain from the host.
//
// Example:
// example.com -> ""
// foo.localhost -> foo
// foo.localhost:8080 -> foo
// foo.example.com -> foo
// bar.foo.example.com -> bar
// blah.bar.foo.example.com -> blah
func subdomain(host string) string {
	host = strings.TrimSpace(host)
	host, _, _ = strings.Cut(host, ":")
	parts := strings.Split(host, ".")
	if len(parts) == 2 && parts[1] != "localhost" {
		return ""
	}
	return parts[0]
}
