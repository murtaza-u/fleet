package cli

import "strings"

// extractDomain extracts the domain name from the host, removing the
// subdomain.
//
// Example:
// example.com -> example.com
// localhost -> localhost
// foo.localhost -> localhost
// foo.localhost:8080 -> localhost
// foo.example.com -> example.com
// bar.foo.example.com -> foo.example.com
// blah.bar.foo.example.com -> bar.foo.example.com
func extractDomain(host string) string {
	host = strings.TrimSpace(host)
	host, _, _ = strings.Cut(host, ":")
	parts := strings.Split(host, ".")
	if len(parts) == 1 {
		return host
	}
	if len(parts) == 2 && !strings.HasPrefix(parts[1], "localhost") {
		return host
	}
	return strings.Join(parts[1:], ".")
}
