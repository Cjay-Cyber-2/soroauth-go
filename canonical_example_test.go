package soroauth_test

import (
	"fmt"

	"github.com/stellar/go-stellar-sdk/xdr"

	"github.com/soroauth/soroauth-go"
)

// Fixed, derived from labels so the output is reproducible. Public test values
// holding nothing, named for this file so each example file stays independent.
const (
	canonicalExampleWallet   = "GASI5PH5KDJFVFU4XGLLRGRNA426GVNZQCLA6PENWMAJRJ4JC3GC4CQ5"
	canonicalExampleContract = "CDHAN6JJTQCD6QDPW3PJPV2NFV3QDQY5HVXEUHCAGH3GZRU2V2Q5V6RA"
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
