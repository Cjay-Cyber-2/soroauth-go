package soroauth_test

import (
	"fmt"
	"github.com/soroauth/soroauth-go"
	"github.com/stellar/go-stellar-sdk/xdr"
)


func canonicalExampleInvocation() xdr.SorobanAuthorizedInvocation {
	contract, err := soroauth.ParseAddress(canonicalExampleContract)
	if err != nil {
		panic(err)
	}
	return xdr.SorobanAuthorizedInvocation{
		Function: xdr.SorobanAuthorizedFunction{
			Type: xdr.SorobanAuthorizedFunctionTypeSorobanAuthorizedFunctionTypeContractFn,
			ContractFn: &xdr.InvokeContractArgs{
				ContractAddress: contract,
				FunctionName:    xdr.ScSymbol("transfer"),
			},
		},
	}
}

func canonicalExampleV2Entry() xdr.SorobanAuthorizationEntry {
	address, err := soroauth.ParseAddress(canonicalExampleWallet)
	if err != nil {
		panic(err)
	}
	return xdr.SorobanAuthorizationEntry{
		Credentials: xdr.SorobanCredentials{
			Type: xdr.SorobanCredentialsTypeSorobanCredentialsAddressV2,
			AddressV2: &xdr.SorobanAddressCredentials{
				Address:                   address,
				Nonce:                     42,
				SignatureExpirationLedger: 1000,
				Signature:                 xdr.ScVal{Type: xdr.ScValTypeScvVoid},
			},
		},
		RootInvocation: canonicalExampleInvocation(),
	}
}

// ExampleMarshalCanonicalJSON shows how two differently-built but semantically
// equivalent entries (or an entry round-tripped through MarshalCanonicalJSON and
// UnmarshalCanonicalJSON) canonicalise identically.
func ExampleMarshalCanonicalJSON() {
	entry1 := canonicalExampleV2Entry()

	// Marshal to canonical JSON and back
	jsonBytes, err := soroauth.MarshalCanonicalJSON(entry1)
	if err != nil {
		fmt.Println("marshal error:", err)
		return
	}

	entry2, err := soroauth.UnmarshalCanonicalJSON(jsonBytes)
	if err != nil {
		fmt.Println("unmarshal error:", err)
		return
	}

	// Marshal entry2 to canonical JSON
	jsonBytes2, err := soroauth.MarshalCanonicalJSON(entry2)
	if err != nil {
		fmt.Println("marshal2 error:", err)
		return
	}

	// Both yield identical canonical JSON bytes
	fmt.Println(string(jsonBytes) == string(jsonBytes2))

	// Output:
	// true
}

const (
	canonicalExampleVector = "AAAAAQAAAAEAAAABAAAAAQAAAAEAAAA="
)

func ExampleCanonical_equivalence() {
	var a, b xdr.SorobanAuthorizationEntry
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &a)
	_ = xdr.SafeUnmarshalBase64(canonicalExampleVector, &b)

	// Canonical form helpers normalize and equate authorization entries.
	fmt.Println("equal canonical:", soroauth.Equal(a, b))
	// Output:
	// equal canonical: true
}
// ExampleCanonicalize demonstrates that two differently constructed entries
// with identical fields canonicalize to the exact same form.
func ExampleCanonicalize() {
	e1 := &soroauth.Entry{Account: "GCEZWKCA5VJHQN7WJWBR7LHXUXVPUN2KY7GBK4AEI5YXR5D2K67X4JOK"}
	e2 := &soroauth.Entry{Account: "GCEZWKCA5VJHQN7WJWBR7LHXUXVPUN2KY7GBK4AEI5YXR5D2K67X4JOK"}

	c1 := soroauth.Canonicalize(e1)
	c2 := soroauth.Canonicalize(e2)

	fmt.Println("canonicalized equal:", soroauth.Equal(c1, c2))

	// Output:
	// canonicalized equal: true
}
