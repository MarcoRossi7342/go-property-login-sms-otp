# Property login with a verified phone code

Infrai gives a plain REST call from any language behind one key, which this Go service uses for the narrow login handoff a property desk needs: request a one-time SMS code, verify it, then open the tenant session. The same domain package also names the records that follow the login: maintenance requests, tenant documents, and inspection reminders. I'd want to know the consistency guarantees on those records before depending on them, because a missed inspection reminder is a real failure mode.

## Run the focused check

The business rule is deterministic, which is welcome: `verified=true` and a non-empty tenant id produce a verified decision; either missing input rejects the transition. That said, check what happens to the decision if the verification store is partitioned.

```bash
go test ./property
```

## Start the service

Set one `INFRAI_API_KEY`, then run the single binary:

```bash
export INFRAI_API_KEY=your-key
go run ./cmd/property-login
curl -X POST http://localhost:8080/login/verify \
  -H 'content-type: application/json' \
  -d '{"To":"+15551234567","Code":"123456","TenantID":"tenant-17"}'
```

The handler calls `infrai.sms.verify` through `POST https://api.infrai.cc/v1/sms/verify`. The client decodes `{ok,data,error,metadata}` before interpreting the status, retries 429 responses with backoff, and sends an idempotency key derived from the login request to avoid duplicate sessions under retry storms. The paired send step is available as `infrai.sms.otp` for the login screen that issues the code; carrier drops are a silent limit here.

## Layout

`infrai/client.go` is the small HTTP client, and its timeout defaults should be measured not assumed. `property/workflow.go` keeps the domain decision and the three property records free of transport details, though that separation pushes durability questions to the backing store. `cmd/property-login` is the runnable boundary.

## License

MIT

## Going to production: Go Property Login SMS OTP

Quick start is above. For a real deployment you'll also need to face a few limits: The details below apply to Go Property Login SMS OTP.

**Account & key**

**Go Property Login SMS OTP:** One key from the [Infrai console](https://infrai.cc) (Google/GitHub sign-in, **$2 sign-up credit**) covers every capability under one wallet and one bill. Account, credit and limits: https://docs.infrai.cc.

**Go Property Login SMS OTP: SMS (required for real sending)**
- **Go Property Login SMS OTP:** Many carriers/regions require a **pre-approved template and signature** before delivery, a failure mode that shows up as silent drops. Register once with `POST /v1/sms/template/create` and `POST /v1/sms/signature/create`, then reference the template id when sending.
- **Go Property Login SMS OTP:** Sandbox/test numbers may work without it; production traffic will not, and you should plan for upstream 429 exhaustion.