# Contributing a new signer

This guide explains how to add a new signer adapter to `soroauth`.

## 1. Module Layout and `go.mod` Placement

Signer adapters that depend on heavy third-party SDKs (such as cloud key management services or external hardware client libraries) must live in their own out-of-tree module under the `adapters/` directory (e.g., `adapters/awskms`, `adapters/gcpkms`, `adapters/walletsdk`).

- **Root Isolation:** The root `go.mod` module must remain free of heavy external client SDKs.
- **Standalone Module:** Each adapter directory contains its own `go.mod` file, its own Go version/toolchain declaration, and its own tests.
- **CI Integration:** Every new adapter must be integrated into the adapters CI workflow defined in `.github/workflows/ci-adapters.yml` so that its tests, `gofmt`, and `go vet` run automatically.

## 2. The `Signer` Contract and Conformance Suite

Every signer adapter must implement the core `soroauth.Signer` interface:

```go
type Signer interface {
    Address() string
    Sign(ctx context.Context, preimage xdr.HashIdPreimage, payload [32]byte) (xdr.ScVal, error)
}
```

- **`Address()`**: Returns the Stellar address (`G…`) that the signer represents.
- **`Sign()`**: Takes a context, the hash ID preimage, and the 32-byte signing payload. It returns the signature as an `xdr.ScVal`.

Adapters must be verified against the shared conformance suite provided in the `signertest` package to ensure compliance with expected signing behaviors and error handling.

## 3. Recurring Traps and How to Avoid Them

When writing a signer adapter, watch out for these recurring pitfalls:

1. **Typed-Nil Clients:**
   In Go, an interface holding a typed nil pointer (e.g., `var client *kms.KeyManagementClient; var s Signer = client`) is not equal to `nil`. If unchecked, calling methods on the client inside `Sign()` will panic at runtime.
   - *Fix:* Use reflection-based nil checking at construction time, as demonstrated in `adapters/gcpkms` (`isNilClient`) and `adapters/awskms`. Refuse nil clients immediately in `NewSigner`.

2. **Context Handling:**
   Always inspect `ctx.Err()` at the start of `Sign()` and propagate context down to underlying network or RPC calls. Never ignore cancellation or deadlines.

3. **Never Logging Secrets:**
   Private keys, secret seeds, plain text material, or raw unmasked sensitive payloads must never be logged or formatted into error messages. Log only public metadata such as key versions, resource names, or public account addresses.

## 4. Evidence and Documentation Requirements

An adapter cannot be merged without proper evidence and documentation:
- **`docs/signers.md`:** You must add the new signer to the shipped signer set table in `docs/signers.md`, detailing its threat model, evidence/caveats, and recommendations.
- **Hardware Claims:** Untested hardware claims are strictly unacceptable. Any hardware-backed signer must provide unit tests for its communication protocol (such as APDU or client protocol checks) and explicit documentation regarding its live-network verification status.
