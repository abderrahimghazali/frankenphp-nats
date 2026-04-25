package nats_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/caddyserver/caddy/v2/caddytest"
	_ "github.com/dunglas/frankenphp/caddy"

	_ "github.com/abderrahimghazali/frankenphp-nats"
)

// TestFrankenPHPNats boots a FrankenPHP test server, hits testdata/index.php,
// and verifies an end-to-end pub/sub round trip against a NATS server running
// on localhost:4222 (provided by CI as a service container, or by the user
// locally via `docker run -p 4222:4222 nats:latest`).
func TestFrankenPHPNats(t *testing.T) {
	tester := caddytest.NewTester(t)
	tester.InitServer(`
		{
			skip_install_trust
			admin localhost:2999
			http_port 9080
			https_port 9443

			frankenphp
		}

		localhost:9080 {
			root testdata/
			php_server
		}
		`, "caddyfile")

	_, body := tester.AssertGetResponse(
		"http://localhost:9080/index.php",
		http.StatusOK,
		"",
	)

	for _, want := range []string{
		"connected: yes",
		"subject: frankenphp.test",
		"data: hello from frankenphp",
		"header x-trace: abc123",
		"in_msgs >= 1: yes",
		"out_msgs >= 1: yes",
		"closed",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("response missing %q\nfull body:\n%s", want, body)
		}
	}
}
