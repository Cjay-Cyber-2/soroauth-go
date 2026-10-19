# Versioning and Compatibility Promises

soroauth publishes multiple distinct surfaces, and the stability guarantees and versioning schemes differ per surface. This document defines what a consumer may rely on across each interface.

## Summary Matrix

| Surface | Governing Version Scheme | Stability Level | Can Change in Minor Release? |
| ------- | ------------------------ | --------------- | ------------------------------ |
| **Go API** | Semantic Versioning (SemVer) | Pre-1.0 (v0.x) unstable API; 1.0+ stable | Yes (for v0.x) / No (for v1.0+) |
| **Wire Output (Preimages & Signatures)** | Stellar Protocol / CAP-governed | Protocol-bound (absolute) | Yes, when protocol or CAPs require |
| **CLI Flags & Commands** | Semantic Versioning (SemVer) | Unstable / Best-effort | Yes |
| **JSON Shapes & RPC Responses** | Semantic Versioning (SemVer) | Experimental / Evolving | Yes |
| **Vector Schema** | Internal testdata versioning | Schema-versioned | Yes |

---

## 1. Go API Compatibility

- **Scheme:** Semantic Versioning (`vMAJOR.MINOR.PATCH`).
- **Pre-1.0 Expectation:** While at version `0.y.z`, minor releases (`0.y`) may introduce breaking changes to exported functions, types, and interfaces. Consumers should pin their dependency exact version or use appropriate version selectors.
- **Post-1.0 Expectation:** Once `v1.0.0` is released, major version bumps are required for breaking changes in public Go APIs.
- **Support Policy:** Reconciled with [docs/sdk-support.md](sdk-support.md) and [SUPPORT.md](../SUPPORT.md). Upgrades to `go-stellar-sdk` dependencies are deliberate and documented in changelogs.

## 2. Wire Output (Preimages & Signatures)

- **Scheme:** Stellar Network Protocol versions & Core Architecture Proposals (CAPs).
- **Promise:** The cryptographic bytes produced by soroauth (such as `HashIdPreimage`, payloads, and signed `SorobanAuthorizationEntry` XDR structures) are strictly governed by the Stellar protocol and specific CAPs (e.g., CAP-46, CAP-71). 
- **Non-negotiable:** soroauth does not define wire compatibility; it faithfully implements the protocol specifications. Any change to emitted bytes requires citing the governing CAP and updating golden vectors as described in [CONTRIBUTING.md](../CONTRIBUTING.md).

## 3. CLI Flags and Commands

- **Scheme:** Semantic Versioning.
- **Promise:** CLI flags (`--entry`, `--passphrase`, etc.) and command outputs are utility layers. While backward compatibility is prioritized where practical, minor or patch releases may adjust CLI options, error messages, or exit codes as the tool evolves.

## 4. JSON Shapes and Structured Outputs

- **Scheme:** Semantic Versioning.
- **Promise:** JSON representations produced by diagnostic or inspection commands are considered evolving unless explicitly stated otherwise in documentation.

## 5. Vector Schema

- **Scheme:** Internal testdata versioning under `testdata/vectors/`.
- **Promise:** Golden test vectors are bound to reference SDK implementations (such as `@stellar/stellar-sdk`). They are never edited by hand and are regenerated via `make vectors` whenever protocol updates dictate.
