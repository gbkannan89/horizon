# Phase 2 — Milestone 6: AI Streaming

**Context Reminder:** Architecture frozen. Specs frozen. Stack frozen. Monorepo. No Redis/K8s/Grafana/testing.

**Milestones completed:** M1-M5 (Docker, Auth/Profile, Transaction CRUD, Planning, Data Management).

**Current state:** AI chat is blocking HTTP (POST /ai/chat sends request, waits for full response, returns). No streaming.

## Objective

Add SSE (Server-Sent Events) streaming to the AI chat endpoint so responses stream token-by-token.

## Backend Audit

Read:
- `services/ai/internal/api/handlers.go` — current PostChat handler (blocking)
- `services/ai/internal/provider/ollama/ollama.go` — check if Ollama client supports streaming
- `services/ai/internal/provider/provider.go` — check provider interface for Stream or Streaming field
- `services/ai/internal/provider/mock.go` — check mock provider

Key questions:
1. Does the AI provider interface support streaming? Look for `Stream` or `ChatStream` method.
2. Does the Ollama provider support streaming? (Ollama API supports SSE natively)
3. Does the mock provider support streaming?

## Implementation

### Backend

1. Create a new endpoint `POST /api/v1/ai/chat/stream` that returns `text/event-stream` content type
2. Use the existing provider's streaming capability (if available)
3. Send SSE events: `data: {"token": "Hello"}\n\n` for each token
4. Send `data: [DONE]\n\n` at the end
5. If no streaming provider available, fall back to blocking response wrapped in SSE format

### Flutter

- Update `advisor_repository.dart` to add `sendMessageStream()` method
- Use `http.Client` with SSE parsing to receive tokens progressively
- Update `advisor_page.dart` to show streaming tokens as they arrive (append to current message content)
- Keep existing `sendMessage()` as fallback

## Files to Modify

| File | Change |
|------|--------|
| `services/ai/internal/api/handlers.go` | Add streaming endpoint |
| `services/ai/register/register.go` | Register streaming route |
| `lib/features/advisor/repository/advisor_repository.dart` | Add streaming method |
| `lib/features/advisor/state/advisor_state.dart` | Support progressive token appending |
| `lib/features/advisor/pages/advisor_page.dart` | Use streaming when available |

## Verification

- [ ] `go build ./cmd/server/...` passes
- [ ] `POST /ai/chat/stream` returns SSE events
- [ ] `flutter analyze` — zero errors
- [ ] Chat messages appear token-by-token (not all at once)
- [ ] Non-streaming fallback works when streaming unavailable

## Deliverables

SSE streaming endpoint, Flutter streaming support.
