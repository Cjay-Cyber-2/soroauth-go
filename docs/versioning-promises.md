# Versioning and Compatibility Promises

`soroauth` maintains distinct stability and compatibility promises across its various surfaces. Because this library signs transactions that move money, stability rules are strict and explicitly bounded per surface.

## 1. Go API Compatibility
- **Governance:** Semantic Versioning (SemVersion 2.0.0) governs the public Go package API.
- **Pre-1.0 (`0.y.z`):** Minor releases (`0.y`) may introduce breaking changes to exported functions, types, or interfaces. Patch releases (`0.y.z`) are strictly bug fixes and backward-compatible additions.
- **Post-1.0 (`x.y.z` where `x >= 1`):** Major releases are required for any breaking change to exported signatures. Deprecations follow the library's deprecation and support policy, providing at least one minor release cycle of warning before removal.

## 2. Wire Output Compatibility
- **Governance:** Stellar Core Protocols and Protocol-specific CAPs (e.g., CAP-46, CAP-71, CAP-85).
- **Promise:** Absolute stability. The byte representation emitted by `soroauth` (such as `HashIdPreimage` hashes and signed `SorobanAuthorizationEntry` blobs) is determined entirely by the underlying network protocol and `go-stellar-sdk`. `soroauth` does not alter or customize wire output; changes are strictly tied to upstream protocol upgrades.

## 3. CLI Flags and Output
- **Governance:** Best-effort stability with strict backward compatibility for standard scripting flags.
- **Promise:** Core flags (`--entry`, `--network-passphrase`, `--valid-until`) are preserved across minor versions. However, diagnostic output formatting, auxiliary subcommands, and non-essential CLI text are not guaranteed against change in minor releases.

## 4. JSON Shapes
- **Governance:** Explicit schema versioning.
- **Promise:** JSON outputs from inspection and simulation commands are structured for durability. Adding new optional fields is permitted in minor releases, but removing or renaming existing fields requires a major version bump or a explicit versioned schema header.

## 5. Vector Schema & Golden Vectors
- **Governance:** Cryptographic reproducibility and upstream SDK pin.
- **Promise:** Golden vectors (`testdata/vectors/`) are locked against the pinned `go-stellar-sdk` version (see `docs/sdk-support.md`). Any change to emitted bytes requires citing the governing CAP and regenerating vectors via `make vectors`.

## Reconciliation
This document operates in full alignment with `docs/sdk-support.md` and the repository's deprecation policy, ensuring zero contradiction between SDK version pinning and library API stability.
