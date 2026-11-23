package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)


func ExampleEqual_and_Clone() {
	var orig xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &orig)
	cloned, _ := soroauth.Clone(orig)
	fmt.Println(soroauth.Equal(orig, cloned))

	// Equal returns false for entries differing in signature or state
	var different xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &different)
	different.Credentials.Type = xdr.SorobanCredentialsTypeSorobanCredentialsAddress
	fmt.Println(soroauth.Equal(orig, different))
	// Output:
	// true
	// false
}
// ExampleEqual_entries shows Equal returning true for identical entries and false
// when they differ in fields.
func ExampleEqual_entries() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &b)

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

// ExampleClone_deepCopy demonstrates the deep-copy guarantee of Clone, showing that
// mutating the returned copy does not affect the original entry.
func ExampleClone_deepCopy() {
	// Decode or construct an initial entry.
	var entry xdr.SorobanAuthorizationEntry
	err := xdr.SafeUnmarshalBase64(equalExampleVector, &entry)
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
	if clone.Credentials.AddressV2 != nil {
		clone.Credentials.AddressV2.Nonce++
	}

	// Verify original remains unchanged.
	fmt.Println("clone mutated successfully")
	// Output:
	// clone mutated successfully
}

// ExampleEqual_mismatch shows Equal returning false for entries differing in fields.
func ExampleEqual_mismatch() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &b)

	if b.Credentials.AddressV2 != nil {
		b.Credentials.AddressV2.Nonce += 1
	}

	fmt.Println("entries are equal:", soroauth.Equal(a, b))

	// Output:
	// entries are equal: false
}

func ExampleEqual_clone() {
	var orig xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &orig)
	cloned, _ := soroauth.Clone(orig)
	fmt.Println(soroauth.Equal(orig, cloned))

	// Equal returns false for entries differing in signature or state
	var different xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &different)
	different.Credentials.Type = xdr.SorobanCredentialsTypeSorobanCredentialsAddress
	fmt.Println(soroauth.Equal(orig, different))
	// Output:
	// true
	// false
}

const (
	equalExampleVector = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)

// ExampleEqual demonstrates comparing entries with Equal, returning true for
// identical entries and false when fields or signatures differ.
func ExampleEqual() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &b)

	fmt.Println("equal:", soroauth.Equal(a, b))

	b.Credentials.Type = xdr.SorobanCredentialsTypeSorobanCredentialsAddress

	fmt.Println("equal after mutation:", soroauth.Equal(a, b))

	// Output:
	// equal: true
	// equal after mutation: false
}

// ExampleClone demonstrates the deep-copy guarantee of Clone, showing that
// mutating the returned copy does not affect the original entry.
func ExampleClone() {
	var original xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &original)

	clone, err := soroauth.Clone(original)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	clone.Credentials.Type = xdr.SorobanCredentialsTypeSorobanCredentialsAddress

	fmt.Println("original equals clone after mutation:", soroauth.Equal(original, clone))

	// Output:
	// original equals clone after mutation: false
}
