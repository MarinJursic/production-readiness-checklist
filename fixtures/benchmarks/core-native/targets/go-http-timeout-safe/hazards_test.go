package safeclient

import "net/http"

func testOnlyHelpers() error {
	_, _ = http.Get("https://example.invalid")
	return http.ListenAndServe(":8081", nil)
}
