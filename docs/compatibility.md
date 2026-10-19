# Versioning and Compatibility Promises

`soroauth` maintains distinct stability and versioning guarantees across
its various exposed surfaces. Because different components serve different
layers of the ecosystem (protocol wire output vs. Go library API vs.
operational CLI flags), consumer expectations differ per surface.

## 1. Go Library API (`soroauth` package and sub-packages)

- **Version Scheme:** Semantic Versioning (SemVersion `vMAJOR.MINOR.PATCH`).
- **Pre-1.0 Expectations:** While below `v1.0.0`, minor releases (`0.y.z`)
  may introduce breaking changes to exported function signatures, types,
  or internal abstractions with appropriate notice in `CHANGELOG.md`.
  Consumers should pin exact minor/patch versions.
- **Minor Release Stability (Post-1.0):** Once `v1.0.0` is released, backward
  compatibility for exported Go APIs will be maintained across minor and
  patch releases. Breaking changes will be reserved for major version bumps
  (`v(N+1).0.0`), preceded by a deprecation cycle as outlined in the
  deprecation policy (`SUPPORT.md` / `CHANGELOG.md`).

## 2. Wire Output and Protocol Signatures

- **Version Scheme:** Governed by the Stellar Network Protocol and Soroban
  Core Application Proposals (CAPs, such as CAP-46 and CAP-71).
- **Stability:** **Fixed by protocol.** The byte sequences and XDR structures
  emitted for preimage hashing and signature verification are immutable for
  a given protocol version. `soroauth` does not alter wire output
  unilaterally.
- **Compatibility:** Changes to wire formats are driven exclusively by
  network protocol upgrades and upstream `go-stellar-sdk` XDR definitions.
  Refer to `docs/sdk-support.md` for SDK upgrade cadences and verification
  practices.

## 3. CLI Flags and Command-Line Interface

- **Version Scheme:** Semantic Versioning corresponding to the library
  version.
- **Stability:** CLI commands, subcommand names, and flags are intended to
  remain stable within major versions. However, minor releases may introduce
  new flags, optional arguments, or diagnostic output format additions.
  Destructive CLI flag removals or renaming will be avoided or deprecated
  with warning messages.

## 4. JSON Shapes and Structured Outputs

- **Version Scheme:** Semantic Versioning.
- **Stability:** JSON outputs (such as inspection payloads and diagnostic
  reports) are structured for machine consumption. Addition of new optional
  JSON keys may occur in minor releases, but existing fields and structural
  types will not be removed or retyped without a major version change or
  explicit deprecation window.

## 5. Vector Schema and Golden Test Vectors

- **Version Scheme:** Tied to reference SDK and protocol releases under
  `testdata/vectors/`.
- **Stability:** Golden test vectors are append-only or replaced strictly
  when required by protocol CAPs. Vectors are never hand-edited and serve
  as byte-level regression invariants.

## Reconciliation with SDK Support and Deprecation Policy

- **SDK Pinning:** As detailed in `docs/sdk-support.md`,
  `github.com/stellar/go-stellar-sdk` is treated as a correctness dependency
  rather than a routine version bump. Upgrades are deliberate and accompanied
  by golden vector verification (`make vectors-check`).
- **Deprecations:** Any deprecated function or CLI flag follows a clear sunset
  timeline documented in `CHANGELOG.md`, ensuring zero contradiction between
  SDK support policies, deprecations, and these stability guarantees.
