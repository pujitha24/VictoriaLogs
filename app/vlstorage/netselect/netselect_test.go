package netselect

import (
	"errors"
	"testing"
)

func TestGetFirstError_AllowPartialResponse(t *testing.T) {
	unavailable := &errUnavailableBackend{err: errors.New("connection refused")}
	misconfigured := errors.New("unexpected status code 400")

	f := func(errs []error, allowPartialResponse, wantErr bool) {
		t.Helper()

		err := getFirstError(errs, allowPartialResponse)
		if (err != nil) != wantErr {
			t.Fatalf("unexpected error state; got %v; wantErr=%v", err, wantErr)
		}
	}

	// all the nodes are healthy
	f([]error{nil, nil}, true, false)
	f([]error{nil, nil}, false, false)

	// unavailable nodes are skipped if there are healthy nodes
	f([]error{nil, unavailable}, true, false)
	f([]error{unavailable, nil}, true, false)

	// all the nodes are unavailable
	f([]error{unavailable, unavailable}, true, true)

	// a node returning an error must be reported irrespective of its position
	f([]error{misconfigured, nil}, true, true)
	f([]error{nil, misconfigured}, true, true)
	f([]error{nil, unavailable, misconfigured}, true, true)

	// partial responses are disallowed
	f([]error{nil, unavailable}, false, true)
}
