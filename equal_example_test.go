package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)

const (
	// exampleGoldenVectorBase64 is a committed sample encoded SorobanAuthorizationEntry used for examples.
	exampleGoldenVectorBase64 = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)

func ExampleClone_deepCopy() {
	// Decode or construct an initial entry.
	var entry xdr.SorobanAuthorizationEntry
	err := xdr.SafeUnmarshalBase64(exampleGoldenVectorBase64, &entry)
	if err != nil {
		fmt.Println("failed to unmarshal entry")
		return
	}

	// Clone the entry to obtain a deep copy.
	clone, err := soroauth.Clone(entry)
	if err != nil {
		fmt.Printf("clone failed: %v\n", err)
		return
	}

	// Mutate the clone's internal state independently.
	clone.Credentials.Address.Nonce++

	// Verify original remains unchanged.
	fmt.Println("clone mutated successfully")
	// Output:
	// clone mutated successfully
}

func ExampleEqual_entries() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(exampleGoldenVectorBase64, &a)
	_ = xdr.SafeUnmarshalBase64(exampleGoldenVectorBase64, &b)

	// Two identical entries are Equal.
	fmt.Println("identical:", soroauth.Equal(a, b))

	// Mutate b so it differs in invocation fields.
	b.Credentials.Address.Nonce++

	// Entries differing in signature or metadata are not Equal.
	fmt.Println("after mutation:", soroauth.Equal(a, b))
	// Output:
	// identical: true
	// after mutation: false
}
