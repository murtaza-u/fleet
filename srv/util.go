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
// foo.example.com -> foo
// bar.foo.example.com -> bar
// blah.bar.foo.example.com -> blah
func subdomain(host string) string {
	host = strings.TrimSpace(host)
	parts := strings.Split(host, ".")
	if len(parts) == 2 {
		return ""
	}
	return parts[0]
}
