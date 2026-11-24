package soroauth_test

import (
	"context"
	"fmt"

	"github.com/soroauth/soroauth-go"
)

const (
	nonceTrackerExampleKey = "GAB7BCQ72A23Z3BXTJVRW66JXYW25J6Y3D7V3X6X44J6H2H2H2H2"
)

func ExampleNewMemoryNonceTracker() {
	tracker := soroauth.NewMemoryNonceTracker()
	var nonce uint64 = 100

	firstUse := tracker.CheckAndRecord(nonce)
	secondUse := tracker.CheckAndRecord(nonce)

	fmt.Println("first use allowed:", firstUse)
	fmt.Println("second use allowed:", secondUse)
	// Output:
	// first use allowed: true
	// second use allowed: false
}

func ExampleNewInMemoryNonceTracker() {
	ctx := context.Background()
	tracker := soroauth.NewInMemoryNonceTracker()
	var nonce int64 = 100

	err1 := tracker.Reserve(ctx, nonceTrackerExampleKey, nonce)
	err2 := tracker.Reserve(ctx, nonceTrackerExampleKey, nonce)

	fmt.Println("first reserve error:", err1)
	fmt.Println("second reserve error:", err2 != nil)
	// Output:
	// first reserve error: <nil>
	// second reserve error: true
}

func ExampleNonceTracker() {
	ctx := context.Background()
	tracker := soroauth.NewInMemoryNonceTracker()
	var nonce int64 = 42

	fmt.Println(tracker.Reserve(ctx, nonceTrackerExampleKey, nonce))
	fmt.Println(tracker.Reserve(ctx, nonceTrackerExampleKey, nonce) != nil)
	// Output:
	// <nil>
	// true
}
