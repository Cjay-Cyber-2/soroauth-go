# Contributing a new signer

This guide explains how to add a new signer adapter to `soroauth`: the required module layout, the contract, the conformance suite, the test expectations, and the documentation it must carry.

There are a dozen signers across several modules, and the pattern is learned by reading existing implementations such as `adapters/gcpkms`, `adapters/awskms`, and `adapters/walletsdk`. The parts that are easy to get wrong — typed-nil clients, context handling, never logging a secret — are covered here.

---

## 1. Module Layout and `go.mod` Placement

Signer adapters that pull in heavy cloud or SDK dependencies (such as AWS KMS, Google Cloud KMS, or the wallet SDK) are structured as **independent Go modules** in their own subdirectories under `adapters/` (e.g., `adapters/awskms`, `adapters/gcpkms`, `adapters/walletsdk`).

### Why separate modules?
Keeping each adapter in its own module ensures that `go get github.com/soroauth/soroauth-go` does not pull heavy cloud SDKs (like `cloud.google.com/go/kms` or `aws-sdk-go-v2`) into every application that signs a Soroban entry. The root module stays lean.

### Rules for a new adapter module:
1. **Directory:** Place your adapter in `adapters/<name>`.
2. **`go.mod`:** Initialize a separate Go module with path `github.com/soroauth/soroauth-go/adapters/<name>`. Use the same Go version as the rest of the repository or the minimum version required by your dependencies.
3. **CI Job:** Add a workflow job in `.github/workflows/ci-adapters.yml` to run `gofmt`, `go vet ./...`, and `go test -race ./...` inside your adapter directory.
4. **No secret logging:** Never log, print, or expose private keys, secrets, or sensitive request payloads.

---

## 2. The `Signer` Contract and Conformance Suite

Every signer implements the core `soroauth.Signer` interface defined in `signer.go`:

```go
type Signer interface {
    Address() string
    Sign(ctx context.Context, preimage xdr.HashIdPreimage, payload [32]byte) (xdr.ScVal, error)
}
```

### Contract Rules:
- **`Address()`**: Returns the classic Stellar account address (`G…`) that this signer authorizes.
- **`Sign()`**: 
  - Takes a context, the `xdr.HashIdPreimage` (for context/inspection, though classic single-key adapters sign the 32-byte hash/payload directly), and the exact 32-byte `payload` (`SHA-256(XDR(HashIdPreimage))`).
  - Returns an `xdr.ScVal` containing the signature (e.g., produced via `soroauth.Ed25519SignatureScVal`).
  - **Must verify locally**: Before returning the signed `ScVal`, the signer must verify the produced signature against the public key derived from its address. If the signature does not verify or has the wrong length, return `soroauth.ErrSignatureMismatch`. Do not emit bad signatures that will only fail on-chain after transaction fees are paid.

### Conformance Testing:
Every adapter must include unit tests verifying:
- Rejection of empty or invalid addresses / key identifiers.
- Context cancellation and timeout handling before calling external services.
- Successful signing resulting in a valid, verified classic account signature `ScVal`.
- Rejection of malformed, wrong-length, or non-verifying signatures.
- Propagation of client/network errors.

---

## 3. Recurring Traps and Existing Adapter Fixes

### Trap 1: The Typed-Nil Client Bug
When a caller declares `var client *kms.KeyManagementClient` and passes it to an adapter constructor, Go passes an interface containing a concrete pointer type with a nil value (`client != nil` evaluates to `true`). 

- **The Problem:** If the constructor only checks `if client == nil`, it lets the typed-nil client through. Later, when `Sign()` is called, it panics with a nil pointer dereference.
- **The Fix (`adapters/gcpkms` and `adapters/awskms`):** Use reflection via `reflect.ValueOf` to inspect whether the interface holds a nil pointer:
  ```go
  func isNilClient(client any) bool {
      if client == nil {
          return true
      }
      v := reflect.ValueOf(client)
      switch v.Kind() {
      case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
          return v.IsNil()
      default:
          return false
      }
  }
  ```
  Refuse typed-nil clients at construction (`NewSigner`) and return `fmt.Errorf(..., soroauth.ErrMissingSigner)` so the caller catches it immediately.

### Trap 2: Context Handling
Always check `ctx.Err()` at the very beginning of `Sign()` before invoking external KMS or hardware calls:
```go
if err := ctx.Err(); err != nil {
    return xdr.ScVal{}, fmt.Errorf("adapter: sign: %w", err)
}
```

### Trap 3: Never Logging Secrets
Ensure logging or error formatting never includes raw private keys, seed phrases, or sensitive message payloads. Test doubles (fakes) should be used in unit tests so signing operations run entirely offline without real hardware or cloud credentials.

---

## 4. Documentation Requirements in `docs/signers.md`

When you add a new signer adapter, you must update `docs/signers.md`:
1. **Add your signer to the summary table:** Include the module path/name, the threat model it addresses, evidence and caveats, and recommendations.
2. **State the evidence honestly:** Be clear about whether the signer has live e2e testnet evidence in `e2e/RESULTS.md` or unit-tested adapter fakes. An untested hardware or cloud claim is not acceptable.
