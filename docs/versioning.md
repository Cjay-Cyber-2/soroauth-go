# Versioning and Compatibility Promises

This document defines the stability and compatibility promises soroauth makes for each of its public surfaces. Because soroauth bridges Go application logic, cryptographic signing preimages, command-line usage, JSON metadata representations, and protocol wire formats, different surfaces carry different stability guarantees.

## Summary of Promises by Surface

| Surface                                            | Governing Scheme                      | Compatibility Promise                                                                                                             |
| -------------------------------------------------- | ------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| **Go API** (exported functions, interfaces, types) | Semantic Versioning (vX.Y.Z)          | Breaking changes only in major versions (v1.0.0+); pre-1.0 (`v0.y.z`) permits breaking changes in minor releases (`y`).             |
| **Wire Output** (`SorobanAuthorizationEntry` XDR)  | Stellar Protocol / CAP Specifications | Fixed by protocol requirements and upstream Stellar Core / SDK versions. Not controlled by soroauth.                              |
| **CLI Flags & Arguments**                          | Semantic Versioning for CLI binary    | Backward-compatible within major versions; flags are preserved or deprecated with warnings before removal.                        |
| **JSON Shapes & Metadata**                         | Semantic Versioning                   | JSON metadata structures emitted by inspection or diagnostics preserve existing keys; additive fields may appear in minor versions. |
| **Vector Schema & Corpus**                         | Protocol & Format versioning          | Golden vectors (`testdata/vectors`) follow protocol standards; schema changes follow CAP adoption.                                |

---

## 1. Go API Compatibility

soroauth adheres to Semantic Versioning (SemVer) 2.0.0 for its exported Go packages:

- **Pre-1.0 state (`v0.y.z`):** While the library is pre-1.0, **minor releases (`y`) may introduce breaking API changes**. Consumers should pin exact versions or minor ranges if depending on pre-1.0 APIs.
- **Post-1.0 state (`v1.0.0+`):** Once a `1.0.0` release is published, public exported functions, interfaces (such as `Signer`), and types will not have breaking changes in patch (`z`) or minor (`y`) releases. Breaking changes will be reserved exclusively for major (`x`) version bumps.

## 2. Wire Output and Cryptographic Preimages

The bytes emitted by `soroauth` (such as `HashIdPreimage` hashes and signed `SorobanAuthorizationEntry` XDR structures) are strictly governed by the Stellar Core protocol and Stellar Core Architecture Proposals (CAPs), specifically CAP-46-11 and CAP-71.

- **No arbitrary changes:** `soroauth` cannot and does not alter wire output arbitrarily. Every byte emitted is required to match reference outputs produced by `@stellar/stellar-sdk` and verified against golden vectors (`testdata/vectors`).
- **Upstream dependency:** The underlying XDR types and serialization are provided by `github.com/stellar/go-stellar-sdk` and `go-xdr` (see `docs/sdk-support.md`).

## 3. CLI Flags and Exit Codes

The `soroauth` command-line interface (`cmd/soroauth`) provides deterministic CLI commands (`sign`, `inspect`, `verify`, `batch`, etc.):

- Existing command names, flag names, and expected exit codes are preserved across minor releases.
- Deprecated flags will emit standard warnings on stderr for at least one minor release cycle before being retired.

## 4. JSON Shapes and Inspection Output

Inspection metadata, JSON formatting from CLI commands, and diagnostic outputs maintain backward compatibility for existing fields.

- New fields may be added in minor or patch releases.
- Existing keys and types will not be removed or renamed without a major version increment.

## 5. Reconcile with SDK Support and Deprecation Policy

This policy is fully aligned with `docs/sdk-support.md` and the library deprecation policy:

- SDK bumps are deliberate correctness inputs managed via `docs/sdk-support.md` rather than automated dependency drift.
- Deprecations in Go types follow standard deprecation comments and a transition window before removal.
