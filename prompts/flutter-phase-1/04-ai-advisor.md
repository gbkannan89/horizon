# Horizon — Milestone 4: AI Advisor

**Phase 3 of 9** | **Estimated: 2 sessions** | **Status: Pending**

---

## CONTEXT REMINDER

**Architecture is frozen. Specifications are frozen. Technology stack is frozen. Do NOT redesign the architecture.**

### Prior Work Completed

- **M1 Foundation**: Auth, Settings, Networking, Error handling stabilized.
- **M2 Transactions**: Full transactions module (list, detail, create, archive, search, filters, pagination).
- **M3 Global Search**: Global search across modules + shared widget library (`SharedLoadingView`, `SharedErrorView`, `SharedEmptyView`, `SharedOfflineView`, `SharedCard`, `SharedSectionHeader`, `SharedStatusBadge`, `SharedSearchBar`, `SharedFilterChip`) + theme constants (`AppSpacing`, `AppRadius`, `AppElevation`, `AppIconSize`).

### Design System References

- `horizon-spec/products/UX_PHILOSOPHY.md`: Teaching moments, explanation patterns, non-judgmental language, confidence-building, progressive disclosure, financial storytelling.
- `horizon-spec/products/DESIGN_SYSTEM.md`: AI Boundary section — what AI may/may not do with design. Cards philosophy, information hierarchy, confidence indicators.
- `horizon-spec/products/NAVIGATION_SYSTEM.md` Section 7 (Advisor Journey Flow): Dashboard -> Advisor (daily briefing) -> View Recommendation -> Ask Question -> View Plan.
- `horizon-spec/engineering/FRONTEND_ENGINEERING_GUIDE.md`: Deterministic Rendering, Presentation Without Business Logic.
- `horizon-spec/products/PROJECT_CONTEXT.md`: AI Principles section — AI assists, does not decide. Every AI output is explainable. AI never presents as human.

### Current App State
- `/advisor` route exists with a PLACEHOLDER page (7-line stub: `Center(child: Text('Advisor'))`)
- Theme: Material 3 teal (`#00897B`)
- Shared widgets available in `lib/shared/widgets/`

### Flutter App Location

`C:\Kannan\Horizon\horizon\apps\mobile\`

---

## STANDARD OPERATING PROCEDURE

1. **Backend Audit** — Verify every required endpoint exists.
2. **Implementation Plan** — List every file to create/modify.
3. **Implementation** — Build the complete feature. NO PLACEHOLDERS. NO TODOs.
4. **Build Verification** — `flutter analyze` — zero errors.
5. **Functional Verification** — Checklist below.
6. **Final Report**.
7. **Message to Solution Architect**.

---

## Step 1: Backend Audit

### Required Endpoints

Check these exist in the codebase at `C:\Kannan\Horizon\horizon\services\ai\` and `C:\Kannan\Horizon\horizon\services\experiences\advisor\`:

| Method | Path | Purpose | Expected Finding |
|--------|------|---------|------------------|
| `POST` | `/api/v1/ai/chat` | Send message, get AI reply | ✅ EXISTS — blocking HTTP, no streaming |
| `GET` | `/api/v1/ai/context` | Financial context for user | ✅ EXISTS |
| `GET` | `/api/v1/ai/prompts` | Suggested prompts | ✅ EXISTS |
| `GET` | `/api/v1/ai/providers` | Available AI providers | ✅ EXISTS |
| `GET` | `/api/v1/ai/providers/active` | Currently active provider | ✅ EXISTS |
| `POST` | `/api/v1/ai/provider/switch` | Switch provider | ✅ EXISTS |
| `GET` | `/api/v1/advisor/context` | AdvisorWorkspace (financial summary) | ✅ EXISTS |
| `GET` | `/api/v1/ai/sessions` | Session history | ❌ DOES NOT EXIST — in-memory only |

For each endpoint, record:
- Full path + method
- Request body/query params shape
- Response JSON shape
- File path

**Critical findings to document:**
- NO streaming endpoint. Chat is synchronous HTTP request-response.
- NO conversation list/history/delete endpoint. Sessions are in-memory and lost on server restart.
- The AI service (`/api/v1/ai/chat`) provides LLM chat. The Advisor experience (`/api/v1/advisor/*`) provides card-based deterministic data.
- flutter_markdown package will be needed for rendering AI responses.

### If ANY endpoint is missing

STOP and report. Do NOT fabricate client code for endpoints that don't exist. If chat endpoint is missing, the milestone cannot proceed.

---

## Step 2: Implementation Plan

Structure:
```
features/advisor/
├── models/advisor_models.dart
├── repository/advisor_repository.dart
├── state/advisor_state.dart
├── pages/advisor_page.dart       (REPLACE placeholder)
└── widgets/advisor_widgets.dart
```

Packages to add to `pubspec.yaml`:
- `flutter_markdown: ^0.7.0` (for rendering AI markdown responses)

---

## Step 3: Implementation

### Models (`features/advisor/models/advisor_models.dart`)

- `ChatMessage`: `String id`, `String role` (user/assistant), `String content`, `DateTime timestamp`, `bool isLoading`
- `ChatSession`: `String sessionId`, `List<ChatMessage> messages`, `DateTime createdAt`, `DateTime? updatedAt`
- `SuggestedPrompt`: `String id`, `String title`, `String prompt`, `String category`
- `AdvisorContext`: Map the backend `AdvisorContext` shape (financial summary for context awareness)
- `ProviderInfo`: `String name`, `List<String> capabilities`, `bool isActive`
- `ChatResponse`: `String reply`, `double? confidence`, `String? provider`, `String? sessionId`

### Repository (`features/advisor/repository/advisor_repository.dart`)

- `sendMessage({String? sessionId, required String message, Map<String, dynamic>? context})` — `POST /api/v1/ai/chat`
- `getContext()` — `GET /api/v1/ai/context`
- `getPrompts()` — `GET /api/v1/ai/prompts`
- `getProviders()` — `GET /api/v1/ai/providers`
- `switchProvider(String provider)` — `POST /api/v1/ai/provider/switch`

### State (`features/advisor/state/advisor_state.dart`)

- `AdvisorStatus` enum: `initial, loading, ready, thinking, error`
- `AdvisorState` with `copyWith`: `status`, `messages[]`, `error`, `context`, `suggestedPrompts[]`, `sessionId`, `providers[]`, `activeProvider`
- `AdvisorNotifier`:
  - `init()` — load context + prompts + providers on first launch
  - `sendMessage(String message)` — add user message, call API, add AI response, handle errors
  - `clearConversation()` — reset messages, start new session
  - `refreshContext()` — reload financial context
  - `switchProvider(String name)` — call API, update active provider

### Page (`features/advisor/pages/advisor_page.dart`)

REPLACE the existing 7-line placeholder with a full chat interface.

**Layout:**
- `Scaffold` with `AppBar` title "Advisor"
- AppBar action: clear conversation button
- Body: `Column` with:
  - Top: Context summary card (compact: net worth, health score quick stats)
  - Middle (Expanded): `ListView.builder` of chat messages
  - Bottom: Input bar with `TextField` + send `IconButton`

**States:**
- `initial`: Show loading spinner while context loads
- `loading`: Shimmer while initial data loads
- `ready` with no messages: Show welcome text + suggested prompts as horizontal scrollable `Wrap` of chips
- `ready` with messages: Show chat bubbles
- `thinking`: Show typing indicator dots animation while waiting for AI response
- `error`: Show error with retry option on the failed message

**Chat Messages:**
- User messages: Right-aligned, primary color background, white text
- AI messages: Left-aligned, surface color background, with small avatar/icon
- AI messages rendered with `flutter_markdown`'s `MarkdownBody` (support bold, italic, lists, code blocks)
- Loading message: Animated dots while waiting
- Scroll-to-bottom when new messages arrive (use `ScrollController.animateTo`)

**Input Bar:**
- `TextField` with hint "Ask your financial advisor..."
- Send `IconButton` (disabled while thinking)
- Disable input and button while waiting for response

### Widgets (`features/advisor/widgets/advisor_widgets.dart`)

- `ChatBubble` — renders user or AI message bubble with markdown support
- `SuggestedPromptChip` — tappable chip for suggested prompts
- `TypingIndicator` — animated dots ("thinking...")
- `ContextSummaryCard` — compact financial summary card

### Navigation

- The existing `/advisor` route in `lib/app/router.dart` already points to `AdvisorPage`. Just update the import path (remove placeholder, point to real page).

---

## Step 4: Build Verification

```powershell
cd C:\Kannan\Horizon\horizon\apps\mobile
flutter pub get
flutter analyze
```

Zero errors required.

---

## VERIFICATION CHECKLIST

- [ ] `flutter analyze` — zero errors
- [ ] flutter_markdown added to pubspec.yaml and `flutter pub get` succeeds
- [ ] Advisor page shows context summary + suggested prompts on first load
- [ ] Tapping a suggested prompt sends it as a message
- [ ] Typing a message and pressing send calls `POST /api/v1/ai/chat`
- [ ] User message appears right-aligned
- [ ] AI response renders as markdown (bold, italic, lists render correctly)
- [ ] Loading indicator shown while AI is responding (typing dots)
- [ ] Multiple messages maintain conversation history (scrollable)
- [ ] Scroll-to-bottom works on new messages
- [ ] Clear conversation resets chat
- [ ] Error state on network failure
- [ ] Retry on failed message works
- [ ] Context summary card shows financial data
- [ ] Provider switch works (if multiple providers available)
- [ ] App builds without errors

---

## DELIVERABLE TEMPLATE

```markdown
## 1. API Inventory Used
...
## 2. Streaming Limitation Documentation
...
## 3. Files Created
...
## 4. Files Modified
...
## 5. Packages Added
...
## 6. Remaining Backend Gaps
...
## 7. Verification Checklist
...
## 8. Message to Solution Architect
```
