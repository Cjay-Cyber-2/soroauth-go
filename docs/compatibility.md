# Versioning and compatibility promises

This document defines the stability and compatibility promises soroauth makes for
its distinct consumer-facing surfaces. Because different surfaces serve
different roles in the Stellar ecosystem — ranging from fixed consensus wire
protocols to evolving developer tooling — the stability guarantees differ per
surface.

## 1. Summary Matrix

| Surface                | Governing Version Scheme            | Stability Promise                                                                                                | Minor Release Policy                                                                      |
| ---------------------- | ----------------------------------- | ---------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| **Go API**             | Semantic Versioning (SemVer v2.0.0) | Stable API for `v1.0+`; pre-1.0 breaking changes may occur in minor releases.                                    | May introduce breaking API changes prior to `v1.0.0`; strict SemVer after `v1.0.0`.        |
| **Wire Output**        | Protocol-governed (CAPs)            | 100% deterministic, byte-for-byte reproducibility against official Stellar SDK and reference implementations.      | Unaffected by soroauth minor versions; changes strictly follow Stellar Core Protocol upgrades and CAPs. |
| **CLI Flags & Output** | Best-effort stability               | Flags are backward-compatible; stdout/json text output may evolve with new diagnostic details.                   | New flags may be added; existing flags are deprecated per the deprecation policy before removal. |
| **JSON Shapes**        | Backward-compatible expansion       | JSON serialization of inspection, verification, and batch reports preserves existing fields.                     | New fields may be appended; existing fields are never renamed or removed without a major version bump. |
| **Vector Schema**      | Versioned fixtures (`v0.1.0+`)      | Golden vectors under `testdata/vectors/` are immutable regression fixtures.                                      | Existing vectors are never modified by hand; new vectors are added alongside new CAPs or protocol versions. |

---

## 2. Detailed Surface Policies

### Go API (`github.com/soroauth/soroauth-go`)

- **Scheme:** Semantic Versioning (`vMAJOR.MINOR.PATCH`).
- **Pre-1.0 Expectations:** While operating under `v0.y.z`, minor releases (`y`)
  may introduce breaking changes to exported function signatures, interfaces
  (such as `Signer`), and package structures. Consumers should pin exact
  versions or vendor dependencies.
- **Post-1.0 Expectations:** Once `v1.0.0` is tagged, public APIs follow strict
  semantic versioning. Deprecations follow the library's deprecation policy.

### Wire Output & Authorization Preimages

- **Scheme:** Stellar Consensus Protocol & Core Architecture Proposals (CAPs).
- **Stability:** The bytes emitted by soroauth when constructing preimages and
  signing authorization entries are entirely fixed by the Soroban host and
  underlying XDR definitions (`go-stellar-sdk`).
- **Compatibility:** soroauth does not independently define wire formats. Any
  change to emitted bytes requires citing the specific CAP (e.g., CAP-46, CAP-71)
  and updating golden reference vectors via `make vectors`.

### CLI Flags and Commands (`cmd/soroauth`)

- **Scheme:** Backward-compatible command-line stability.
- **Stability:** Existing flags (`--entry`, `--passphrase`, etc.) and
  subcommands remain supported. Flag deprecations follow a multi-release warning
  cycle before removal.
- **Output Stability:** CLI text output is designed for operator legibility and
  may receive supplementary formatting or informational fields in minor
  versions.

### JSON Shapes and Inspection Metadata

- **Scheme:** Additive schema evolution.
- **Stability:** JSON outputs (such as inspection reports, verification
  verdicts, and batch summaries) guarantee that existing keys and types remain
  present and unchanged. New fields may be introduced in minor releases.

### Vector Schema & Test Fixtures

- **Scheme:** Immutable testdata versioning.
- **Stability:** Golden vectors located under `testdata/vectors/` and
  `testdata/differential/` are fixed assets. They are never edited by hand; they
  are regenerated through official reference SDK runners (`make vectors`).

---

## 3. Alignment with SDK Support and Deprecation Policy

- **SDK Alignment:** As outlined in [`docs/sdk-support.md`](sdk-support.md),
  `go-stellar-sdk` version pins are correctness inputs rather than routine
  dependencies. Upgrading the SDK is a deliberate action tied to protocol
  feature requirements.
- **Deprecation Policy:** Deprecated functions, fields, or CLI flags emit clear
  warnings and observe a deprecation window before any removal in accordance
  with standard maintenance practices.
