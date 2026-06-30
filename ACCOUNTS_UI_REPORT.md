# ACCOUNTS UI REPORT

**Phase:** 8.7
**Status:** Complete
**Build:** `flutter analyze`: 0 errors | `flutter build apk --debug`: SUCCESS

---

## Screen Architecture

```
┌──────────────────────────────────────────────────────────────────────┐
│                     Accounts Dashboard Screen                         │
├──────────────────────────────────────────────────────────────────────┤
│  AppBar: "Accounts"                                                   │
├──────────────────────────────────────────────────────────────────────┤
│  Balance Summary Card (Total / Available / Spendable)                │
├──────────────────────────────────────────────────────────────────────┤
│  Dashboard Cards (8 from API)                                        │
├──────────────────────────────────────────────────────────────────────┤
│  Cash Flow Section (Inflow / Outflow / Net / Projected)              │
├──────────────────────────────────────────────────────────────────────┤
│  Account Health Section (Healthy % + breakdown)                      │
├──────────────────────────────────────────────────────────────────────┤
│  Accounts List (filtered by type chips)                              │
│  ┌──────────────────────────────────────────────────────────────────┐│
│  │ [Savings] [Checking] [Credit] [Investment]                      ││
│  ├──────────────────────────────────────────────────────────────────┤│
│  │ AcctCard                                                        ││
│  │  [Icon] Account Name                    [Type Badge]             ││
│  │         ₹5,00,000  •  INR               [Institution]           ││
│  └──────────────────────────────────────────────────────────────────┘│
├──────────────────────────────────────────────────────────────────────┤
│  Account Detail Screen (push on tap)                                 │
└──────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────┐
│                      Account Detail Screen                            │
├──────────────────────────────────────────────────────────────────────┤
│  AppBar: Account Name                                                │
├──────────────────────────────────────────────────────────────────────┤
│  Header: Icon + Name + Type + Currency                               │
│         Status, Health, Opened, Interest Rate, Credit Limit          │
├──────────────────────────────────────────────────────────────────────┤
│  Balances Section (7 balance types)                                  │
│  ┌──────────────────────────────────────────────────────────────────┐│
│  │ Current         ₹2,50,000    Net sum of all posted transactions ││
│  │ Available       ₹2,00,000    Current + uncleared - reserved     ││
│  │ Cleared         ₹2,40,000    Confirmed and posted transactions  ││
│  │ Uncleared       ₹10,000      Pending transactions               ││
│  │ Reserved        ₹30,000      Funds committed but not allocated  ││
│  │ Allocated       ₹20,000      Assigned to goals                  ││
│  │ Spendable       ₹1,80,000    Available - allocated              ││
│  └──────────────────────────────────────────────────────────────────┘│
├──────────────────────────────────────────────────────────────────────┤
│  Recent Transactions (up to 10)                                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

## Widget Hierarchy

| Widget | Purpose | Data Source |
|---|---|---|
| `AcctCard` | Account card with icon, name, type, balance | `AcctCardData` |
| `AcctBalanceCard` | Total/Available/Spendable summary | `BalanceSummaryData` |
| `AcctCardWidget` | Generic dashboard card | `AcctDashboardCard` |

## API Integration (7 endpoints)

| Endpoint | Section |
|---|---|
| `GET /accounts/experience` | Dashboard |
| `GET /accounts/summary` | Account list |
| `GET /accounts/balances` | Balance summary |
| `GET /accounts/cashflow` | Cash flow |
| `GET /accounts/health` | Account health |
| `GET /accounts/{id}` | Account detail |

All 5 dashboard endpoints loaded in parallel. Detail screen fetches `/accounts/{id}`.

## Account Type Mapping

| Type | Icon | Color |
|---|---|---|
| Savings | `savings` | green |
| Checking | `account_balance` | blue |
| Credit | `credit_card` | purple |
| Investment | `trending_up` | teal |
| Loan | `receipt` | red |

## Filtering

- Filter chips for account type (All, Savings, Checking, Credit, Investment)
- Client-side filter applied against loaded accounts

## Detail Screen

- Header with all account metadata
- 7 balance types with label + description + value
- Up to 10 recent transactions with amount, category, running balance

## Build Status

| Check | Status |
|---|---|
| `flutter analyze` | ✅ **0 errors** |
| `flutter build apk --debug` | ✅ **SUCCESS** |
| Dashboard | ✅ Balance, cards, cash flow, health |
| Account list | ✅ Cards with type icons + filter chips |
| Account detail | ✅ Header, 7 balance types, 10 transactions |
| Parallel loading | ✅ 5 endpoints via Future.wait |

## File Inventory

| File | Lines | Purpose |
|---|---|---|
| `features/accounts/models/acct_models.dart` | 220 | 10 response model classes |
| `features/accounts/repository/acct_repository.dart` | 53 | 10 endpoint methods |
| `features/accounts/widgets/acct_widgets.dart` | 105 | AcctCard, AcctBalanceCard, AcctCardWidget |
| `features/accounts/pages/accounts_page.dart` | 210 | Full accounts dashboard + list |
| `features/accounts/pages/account_detail_page.dart` | 130 | Account detail with balances + transactions |

## Ready for Phase 8.8 — Settings/Profile UI

```
Flutter analyze: 0 errors
APK build:       SUCCESS
Endpoints:       7 integrated
Parallel loads:  5 endpoints
Screens:         2 (list + detail)
Card types:      8 mapped
Balance types:   7 displayed
```
