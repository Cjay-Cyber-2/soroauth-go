package soroauth_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/soroauth/soroauth-go"
)

// ExampleNonceTracker shows how to record a nonce and refuse a reuse.
//
// Note that NewInMemoryNonceTracker is per-process only and shares nothing
// across process restarts or different instances unless backed by durable storage.
func ExampleNonceTracker() {
	tracker := soroauth.NewInMemoryNonceTracker()
	ctx := context.Background()
	const address = "GBEXAMPLE"

	// Record the first use of nonce 100.
	if err := tracker.Reserve(ctx, address, 100); err != nil {
		fmt.Println("error:", err)
		return
	}

	// A second use of the same nonce for the same address is refused.
	err := tracker.Reserve(ctx, address, 100)
	fmt.Println("reuse refused:", errors.Is(err, soroauth.ErrNonceAlreadyReserved))

	// Output:
	// reuse refused: true
}
