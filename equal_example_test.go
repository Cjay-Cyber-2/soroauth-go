package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)



func ExampleEqual() {
	var orig xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &orig)
	clone, _ := soroauth.Clone(orig)

	fmt.Println(soroauth.Equal(orig, clone))
	// Output:
	// true
}
func ExampleEqual_and_Clone() {
	var orig xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &orig)
	clone, _ := soroauth.Clone(orig)

	fmt.Println(soroauth.Equal(orig, clone))
	// Output:
	// true
}

func ExampleEqual_entries() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &b)

	fmt.Println("identical:", soroauth.Equal(a, b))

	b.Credentials.Address.Nonce++

	fmt.Println("after mutation:", soroauth.Equal(a, b))
	// Output:
	// identical: true
	// after mutation: false
}

func ExampleEqual_clone() {
	var orig xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &orig)
	cloned, _ := soroauth.Clone(orig)
	fmt.Println(soroauth.Equal(orig, cloned))

	var different xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &different)
	different.Credentials.Type = xdr.SorobanCredentialsTypeSorobanCredentialsAddress
	fmt.Println(soroauth.Equal(orig, different))
	// Output:
	// true
	// false
}

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

func ExampleClone_deepCopy() {
	var entry xdr.SorobanAuthorizationEntry
	err := xdr.SafeUnmarshalBase64(equalExampleVector, &entry)
	if err != nil {
		fmt.Println("failed to unmarshal entry")
		return
	}

	clone, err := soroauth.Clone(entry)
	if err != nil {
		fmt.Printf("clone failed: %v\n", err)
		return
	}

	if clone.Credentials.AddressV2 != nil {
		clone.Credentials.AddressV2.Nonce++
	}

	fmt.Println("clone mutated successfully")
	// Output:
	// clone mutated successfully
}

func ExampleClone() {
	var orig xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(equalExampleVector, &orig)

	clone, err := soroauth.Clone(orig)
	if err != nil {
		panic(err)
	}

	if clone.Credentials.AddressV2 != nil {
		clone.Credentials.AddressV2.Nonce = 999
	}

	fmt.Println("original equals clone after clone mutation:", soroauth.Equal(orig, clone))
	// Output:
	// original equals clone after clone mutation: false
}

const (
	equalExampleVector = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)
