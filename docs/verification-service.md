# Verification service

`server/` is a small HTTP service around `VerifyEntry`, for wallets and custody
systems that want a pre-submission check without embedding Go. It answers one
question per entry — do the signatures on it verify over the payload it commits
to? — with the same verdicts and the same JSON shape as
`soroauth verify --json`.

## The contract

The HTTP surface — every endpoint, status code, error shape and the exact JSON
schema of a report — is [`server/openapi.yaml`](../server/openapi.yaml). That
document is the source of truth; this page explains the deployment around it
rather than restating it.

The spec is not a description that can drift. The conformance check in
[`server/openapi/`](../server/openapi/README.md) starts the real service, drives
it as a client would, validates every response against the document, and
exercises the client generated from it. A handler change the spec does not
describe fails that check. In particular the verdict enum keeps `cannot_check`
as a distinct value, so a custom account's signature is reported as uncheckable
rather than flattened to a boolean.

## Configuration

Every bound has a default an operator does not have to set. Each is read from
the environment at startup, and a malformed value refuses to start rather than
falling back silently.

| Variable                  | Default           | Meaning                                                                                                       |
| ------------------------- | ----------------- | ------------------------------------------------------------------------------------------------------------- |
| `ADDR`                    | `127.0.0.1:8080`  | API listen address. A non-loopback host requires `AUTH_TOKEN` (see below).                                    |
| `METRICS_ADDR`            | `127.0.0.1:9090`  | Prometheus listen address. `off`, `none` or `-` disables the metrics server entirely.                         |
| `AUTH_TOKEN`              | unset             | Bearer token `/verify` requires. Unset means no authentication, which is only allowed on a loopback bind.     |
| `READ_TIMEOUT`            | `5s`              | Whole-request read deadline.                                                                                  |
| `READ_HEADER_TIMEOUT`     | `5s`              | Header read deadline, the classic slowloris bound.                                                            |
| `WRITE_TIMEOUT`           | `10s`             | Response write deadline.                                                                                      |
| `IDLE_TIMEOUT`            | `120s`            | Keep-alive idle timeout.                                                                                      |
| `MAX_BODY_BYTES`          | `1048576` (1 MiB) | Request body cap. An entry that decodes within the library's own depth limit fits well under it.              |
| `MAX_CONCURRENT_REQUESTS` | `64`              | In-flight `/verify` handlers. A request over the bound is refused with `503` and `Retry-After`, never queued. |

The graceful-shutdown deadline is fixed at 10s and is not configurable; a signal
drains the API before the metrics listener, so a final scrape still describes a
live process.

## Threat model

The service is hardened for untrusted exposure, but "hardened" is not the same
as "public internet, unexamined". What is now safe and what is still not:

**Safe by construction.**

- **No secret can reach it.** Verification needs public data only — the entry,
the network passphrase, an optional expiration assertion — so there is no
field that accepts a seed, nothing is stored, and no key is ever involved. The
service verifies; it never signs.
- **Bounded input.** The body is capped by `MAX_BODY_BYTES`, and the library
refuses entries past its own decode-depth and size limits, so a hostile entry
cannot drive unbounded decoding.
- **Bounded concurrency.** At most `MAX_CONCURRENT_REQUESTS` verifications run
at once; the rest are refused with `503` rather than piling up and slowing
every request until the timeout.
- **Bounded time.** Read, header, write and idle deadlines are all set, so a
stalled or slow request cannot hold a connection open indefinitely.
- **No bytes in the logs.** One structured line is written per request, and it
carries the method, the matched route, the status, the duration, the response
size and the peer address — never the request body, the entry, the
`Authorization` header or the token. A test asserts this.
- **Fail closed on the bind.** A non-loopback address with no `AUTH_TOKEN`
refuses to start rather than accepting entries from anyone who can reach it. A
warning nobody reads is how a service ends up on a public interface with no
authentication.

**Still the operator's concern.**

- **Transport security.** The service speaks plain HTTP. Terminate TLS in front
of it; a bearer token sent over an unencrypted hop is a bearer token someone
else can read.
- **Rate limiting and abuse control.** Concurrency is bounded, but there is no
request-rate limiter or per-client quota. Put one in front if the service is
reachable by untrusted clients.
- **Token issuance, rotation and storage.** `AUTH_TOKEN` is a shared secret the
service compares in constant time; how it is minted, rotated and delivered to
callers belongs to the deployment.
- **Launch arguments and container inspection.** An environment variable is
visible to anything that can inspect the process or container
(`/proc/<pid>/environ`, `docker inspect`). Treat `AUTH_TOKEN` with the care
any process environment deserves.

There is no authentication on `/healthz`, so a load balancer's check needs no
token. It reports liveness and the version only.

If signing is ever added to this service, authentication stops being optional:
the token check and its fail-closed bind already assume that, and the metrics
below already carry a signer-latency signal for when it arrives.

## Operational Runbook

### 1. Configuration Reference

Every option is read from the environment at startup. A malformed value causes
the process to abort immediately rather than falling back silently.

| Variable                  | Default           | Meaning & Effect                                                                                     | Failure Mode on Malformed Value                                                                   |
| ------------------------- | ----------------- | ---------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| `ADDR`                    | `127.0.0.1:8080`  | TCP host and port for the API server. Binding to a non-loopback interface requires `AUTH_TOKEN` to be set. | Process fails to bind (EADDRINUSE or permission denied) or refuses to start if non-loopback lacks token. |
| `METRICS_ADDR`            | `127.0.0.1:9090`  | TCP host and port for Prometheus metrics scraping. Set to `off`, `none`, or `-` to disable the metrics listener. | Startup aborts if port is invalid or already bound.                                               |
| `AUTH_TOKEN`              | unset             | Bearer token required for all `/verify` requests. Unset means no authentication, restricted strictly to loopback addresses. | Process fails to start if `ADDR` is non-loopback and `AUTH_TOKEN` is missing.                     |
| `READ_TIMEOUT`            | `5s`              | Maximum duration for reading the entire request, including the body.                                 | HTTP 408 Request Timeout / connection dropped.                                                    |
| `READ_HEADER_TIMEOUT`     | `5s`              | Maximum duration for reading request headers (mitigates Slowloris attacks).                          | Connection terminated abruptly.                                                                 |
| `WRITE_TIMEOUT`           | `10s`             | Maximum duration before timing out writes of the HTTP response.                                      | Server cancels write, client receives truncated or reset connection.                              |
| `IDLE_TIMEOUT`            | `120s`            | Maximum time to wait for the next request when keep-alives are enabled.                              | Keep-alive connection closed.                                                                     |
| `MAX_BODY_BYTES`          | `1048576` (1 MiB) | Strict upper bound on request body size. Larger bodies are rejected before XDR decoding.             | HTTP 413 Payload Too Large.                                                                       |
| `MAX_CONCURRENT_REQUESTS` | `64`              | Maximum in-flight `/verify` evaluations allowed simultaneously.                                      | HTTP 503 Service Unavailable with `Retry-After` header.                                           |

### 2. Metrics & Alerting Signals

Operators should configure Prometheus and Alertmanager to watch the following
key signals:

- **Rising Refusal Rate (`soroauth_verify_refusals_total`):**
  - *Indicator:* A sharp or sustained increase in refusals categorized by
    reason (`invalid_body`, `decode_error`, `depth_exceeded`,
    `unsupported_arm`).
  - *Meaning:* Clients are submitting malformed payloads, out-of-spec
    authorization entries, or unsupported credential types. If specific clients
    spike in refusals, coordinate with client teams to validate their XDR
    serialization and SDK versions.

- **HTTP 503 Overload (`soroauth_http_requests_total{status="503"}`):**
  - *Indicator:* Non-zero or rising count of 503 responses.
  - *Meaning:* The concurrency ceiling (`MAX_CONCURRENT_REQUESTS = 64`) is
    saturated. The service is rejecting excess traffic rather than queuing.
  - *Action:* Scale horizontally behind a load balancer or evaluate if downstream
    clients are retrying aggressively.

- **Elevated Verification Latency (`soroauth_http_request_duration_seconds`):**
  - *Indicator:* P95 or P99 latency exceeding 50ms for the `/verify` route.
  - *Meaning:* Although stateless verification is computationally bounded,
    excessive nesting depth in authorization delegate trees or deep ScVal
    structures can inflate CPU usage. Inspect request size payloads.

- **Signer Latency / Activity (`soroauth_signer_duration_seconds`):**
  - *Indicator:* Spikes in signing phase duration.
  - *Meaning:* While the standalone verification service does not sign, custom
    deployments or integrations embedding signing hooks may experience HSM or
    external KMS latency. Watch for queueing or rate limits at the upstream
    signer.

### 3. Failure Response Playbook

#### Scenario A: Wedged Signer / Unresponsive Integration Hook

- *Symptom:* Requests hang or saturate `MAX_CONCURRENT_REQUESTS`, leading to 503
  spikes and timeouts.
- *Diagnosis:* Check container resource utilization (CPU/Memory). Verify if any
  custom signer hooks or downstream KMS wrappers are blocking indefinitely.
- *Response Steps:*
  1. If container is non-responsive, invoke graceful shutdown (or SIGTERM). The
     process drains active requests up to the 10s deadline.
  2. If Kubernetes orchestrates deployment, restart the pod via `kubectl rollout
     restart deployment/soroauth-verifier`.
  3. Investigate external signer connectivity and IAM credential validity if KMS
     integration is active.

#### Scenario B: Unreachable RPC or Network Partition

- *Symptom:* Verification service itself remains live (`/healthz` returns 200
  OK), but upstream verification dependencies or telemetry sinks fail.
- *Diagnosis:* Check container egress rules, DNS resolution, and Stellar RPC
  endpoint reachability (`curl -X POST ...`).
- *Response Steps:*
  1. Note that the verification service is designed to be **stateless and
     self-contained**: it does *not* make live RPC calls during standard
     `VerifyEntry` execution (as verification requires public data only: entry,
     passphrase, optional expiration).
  2. If the service is integrated into a larger signing or simulation pipeline
     that *does* rely on RPC, verify RPC provider status.
  3. If DNS/routing failures affect metrics emission (`METRICS_ADDR`), verify
     Prometheus scraper connectivity.

#### Scenario C: Memory Pressure / OOM Killer Termination

- *Symptom:* Container restarts unexpectedly with exit code 137.
- *Diagnosis:* Inspect kernel logs (`dmesg -T | grep -i oom`) and Prometheus
  memory usage metrics.
- *Response Steps:*
  1. Extremely large or deeply nested hostile XDR structures can temporarily
     spike heap allocations during recursive decoding.
  2. Verify that `MAX_BODY_BYTES` is strictly enforced at 1 MiB.
  3. Adjust container memory limits upward (recommended minimum: 256 MiB
     headroom) if legitimate multi-delegate trees cause memory spikes.

### 4. Deployment Assumptions & Threat Model Summary

- **Stateless & Scalable:** The service maintains no session state, database,
  or local cache. It can be scaled horizontally behind any standard Layer 4 or
  Layer 7 load balancer.
- **Network Topology:** Deploy behind a TLS-terminating reverse proxy (e.g.,
  Nginx, Envoy, or cloud load balancer). The service itself speaks plain HTTP.
- **Resource Boundaries:** CPU-bound verification ensures predictable resource
  profiles. Memory bounds are enforced via strict body sizing and XDR depth
  limits.
- **Security Isolation:** Operates with zero secrets or private keys in
  memory. Even complete container compromise yields no signing credentials.

## Metrics

The service exposes Prometheus metrics on `METRICS_ADDR`, **separate from the
API port** so it can stay on loopback or a private network and does not have to
be reachable wherever `/verify` is. `off`/`none`/`-` disables the listener.

| Metric                                   | Type      | Labels                       | Meaning                                               |
| ---------------------------------------- | --------- | ---------------------------- | ----------------------------------------------------- |
| `soroauth_http_requests_total`           | counter   | `method`, `route`, `status`  | HTTP requests handled.                                |
| `soroauth_http_request_duration_seconds` | histogram | `route`                      | Request latency, default buckets.                     |
| `soroauth_verify_verdicts_total`         | counter   | `credential_type`, `verdict` | One per credential node, by arm and verdict.          |
| `soroauth_verify_refusals_total`         | counter   | `reason`                     | Verification and transport refusals, by reason.       |
| `soroauth_signer_duration_seconds`       | histogram | `phase`, `credential_type`   | Signing-phase latency, fed by `soroauth.Hook` events. |

**Every label set is closed.** A value outside the set is clamped to a fallback
rather than creating a new series, so no label can carry an address, a nonce or
any other unbounded string — the usual way a metrics change takes down a
Prometheus. A test walks every gathered series and asserts each label name and
value is in the set the code declares.

- `method`: the standard methods; anything else is `other`.
- `route`: `/verify`, `/healthz`, or `unmatched` — the path is logged as the
  matched route template, so an arbitrary path cannot become a label or a
  log-injection vector.
- `status`: `200`, `400`, `401`, `404`, `405`, `413`, `500`, `503`; anything
  else is `other`.
- `credential_type`: `source_account`, `address`, `address_v2`,
  `address_with_delegates`, or `unknown`.
- `verdict`: `verified`, `unsigned`, `invalid`, `cannot_check`, or `unknown`.
- `reason`: `invalid`, `cannot_check`, `unsigned`, `expiration_mismatch`,
  `unauthorized`, `method_not_allowed`, `too_large`, `overloaded`, `internal`,
  `bad_request`, `not_found`, or `other`.
- `phase`: `preimage`, `sign`, `post_sign`, `write`, or `unknown`.

`cannot_check` is a distinct verdict and a distinct refusal reason: a custom
account whose signature only `__check_auth` can decide is visible as such, not
counted as a pass.

The signer histogram is the bridge to [`hooks.go`](../hooks.go): its only writer
is a `soroauth.Hook` adapter, so metric values come from the one event stream
rather than from instrumenting a signing path separately. The service verifies
rather than signs today, so the histogram is empty until signing is added; the
dashboard's panel says so.

A reference Grafana dashboard is committed at
[`docs/dashboards/soroauth-verification.json`](dashboards/soroauth-verification.json),
and [`docs/dashboards/example-metrics.txt`](dashboards/example-metrics.txt) is
an export from a real run of the endpoint, so the dashboard queries are known to
have data behind them. Import the JSON and point its `DS_PROMETHEUS` input at
the Prometheus scraping `METRICS_ADDR`.

## Running it

```sh
go run ./server                                   # loopback :8080, metrics on 127.0.0.1:9090
ADDR=0.0.0.0:8080 AUTH_TOKEN=$(openssl rand -hex 32) go run ./server   # exposed, authenticated
```

`make server` builds `bin/soroauth-server` with the release version stamped in
(reported by `/healthz` and `soroauth-server version`). The published container
image is CLI-only; serving this over a network is a deployment choice the image
does not make for you.

## Coverage

`server/main_test.go` serves the real routes over `httptest` and covers a signed
entry (verified), a partially signed delegates tree (one verified node, the rest
unsigned), malformed requests (missing fields, non-base64, base64-but-not-XDR,
non-JSON, wrong method), a tampered expiration (invalid), a contradicted
`valid_until_ledger` assertion, and an envelope input with its one-based
positions. The verdicts are asserted against the library's `Verdict*` constants
rather than string literals, so the tests fail if the service ever stops
speaking the CLI's verdicts.

`server/hardening_test.go` covers the bounds and the fail-closed bind, and
asserts the request log carries no entry bytes or secret.
`server/metrics_test.go` asserts the metric label sets are closed and that the
committed dashboard and example scrape reference the metrics the service
actually emits.
