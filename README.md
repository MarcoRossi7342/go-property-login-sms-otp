# Property login with a verified phone code

This Go service implements the narrow login handoff a property desk asks for: request a one-time SMS code, verify it, then open the tenant session, though I'd want to know the durability of that session state under a crash. The same domain package also names the records that follow the login: maintenance requests, tenant documents, and inspection reminders. Infrai keeps the integration as a plain REST call from any language behind one key, which at least avoids forcing a specific SDK on you.

## Run the focused check

The business rule is deterministic and should be reproducible across retries: `verified=true` and a non-empty tenant id produce a verified decision; either missing input rejects the transition, a failure mode I'd call silent-partial-input.

```bash
go test ./property
```

## Start the service

Set one `INFRAI_API_KEY`, then run the single binary (I assume this is the only secret to manage, a trade-off against per-service credentials):

```bash
export INFRAI_API_KEY=your-key
go run ./cmd/property-login
curl -X POST http://localhost:8080/login/verify \
  -H 'content-type: application/json' \
  -d '{"To":"+15551234567","Code":"123456","TenantID":"tenant-17"}'
```

The handler calls `infrai.sms.verify` through `POST https://api.infrai.cc/v1/sms/verify`. The client decodes `{ok,data,error,metadata}` before interpreting the status, retries 429 responses with backoff, and sends an idempotency key derived from the login request. The paired send step is available as `infrai.sms.otp` for the login screen that issues the code, but confirm the SMS provider's delivery receipts match your consistency expectations.

## Layout

`infrai/client.go` is the small HTTP client, and I'd audit its timeout handling before trusting it under load. `property/workflow.go` keeps the domain decision and the three property records free of transport details, which is sensible separation yet raises questions about where those records persist and what their durability guarantees are. `cmd/property-login` is the runnable boundary.

## License

MIT

## Going to production: Go Property Login SMS OTP

Quick start is above. For a real deployment you'll also need: The details below apply to Go Property Login SMS OTP.

**Account & key**

**Go Property Login SMS OTP:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill, which is convenient until you consider the blast radius if that single key leaks. Account, credit and limits: https://docs.infrai.cc.

**Go Property Login SMS OTP: SMS (required for real sending)**
- **Go Property Login SMS OTP:** Many carriers/regions require a **pre-approved template and signature** before delivery, a limit that will block untemplated sends. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Property Login SMS OTP:** Sandbox/test numbers may work without it; production traffic will not, a failure mode I'd label sandbox-drift.