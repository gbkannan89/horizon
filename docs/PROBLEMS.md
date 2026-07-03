# Horizon (FinArtha) Problem Statements & Solutions

This document outlines the core problems Horizon (FinArtha) is built to solve. It serves as the guiding "why" behind the platform's architecture and capabilities, ensuring every feature maps to a real, felt problem for the target user.

## The Unifying Problem Statement
> A busy, financially-inexperienced salaried person (and their family) has all the data and enough income to build real wealth - but no honest, always-available guide that turns that data into "what should I do now," protects them from their own emotions and blind spots, and keeps them going after inevitable mistakes.

Horizon exists to close the gap between "who you intend to be with money" and "who you actually are" - and to make wealth and peace almost inevitable given what you already earn.

---

## 1. Core Problems

### 1.1. Finance apps only look backward.
* **The Problem**: Every tool (Excel, tracking apps) is a mirror - it shows where money went. Nobody tells you *what to do next*. You need a windshield, not a rear-view mirror.
* **The Solution**: 
  * **Projection Engine**: Runs 30-year horizon simulations mapped across multiple paths.
  * **Optimization & Recommendation Engines**: Analyzes your data to produce forward-looking, actionable next steps instead of just parsing historical spend.

### 1.2. Good financial advice is locked behind advisors you can't afford.
* **The Problem**: The knowledge that builds wealth is mostly free and well-established - but a busy, non-expert person has no one to apply it to *their* specific numbers, consistently.
* **The Solution**: 
  * **AI Provider Framework & Advisor Experience**: Securely aggregates 14 data dimensions (net worth, goals, risk metrics) and runs them through structured models to provide hyper-personalized, context-aware financial advice 24/7.

### 1.3. "Am I okay?" has no honest answer.
* **The Problem**: You can't quickly tell if you're financially healthy or fragile - or why. There's no simple, explained verdict on your current situation.
* **The Solution**: 
  * **Health Score Engine**: Assesses financial health across 13 dimensions to provide a clear, explained metric score.
  * **Risk Engine**: Continuously monitors 18 indicators and subjects finances to 5 distinct stress tests.

### 1.4. Paralysis on investing + future regret.
* **The Problem**: "The market is moving, I know I should act, but I don't know where/how, so I do nothing" - and later regret the years of inaction.
* **The Solution**: 
  * **Optimization Engine**: Provides concrete target portfolio allocations and ranks investment candidates to reduce decision paralysis.
  * **Simulation Engine**: Enables side-by-side scenario comparisons to instantly visualize the long-term cost of inaction versus action.

### 1.5. Manual tracking dies in two weeks.
* **The Problem**: A workaholic won't maintain a spreadsheet. Any solution that demands constant manual effort fails.
* **The Solution**: 
  * **Automation Engines**: Auto-categorization, smart alerts, and scheduled reports handle background data crunching.
  * **Import Service**: Automates parsing for bank statements (CSVs, Excel, PDFs) to minimize manual entry friction.

---

## 2. Behavioral Problems

### 2.1. One slip makes people quit.
* **The Problem**: People don't get poor from a single bad month - they get poor from *abandoning the plan* after it. There's no tool that absorbs a mistake and gently recovers you.
* **The Solution**: 
  * **Dynamic Recalibration via Projections**: Dynamically recalibrates your path following a bad month, illustrating that a temporary setback does not derail a 30-year plan, instead of showing a failed static budget.

### 2.2. Optimism bias hides the real trajectory.
* **The Problem**: You promise "I'll pay ₹5k extra and close the loan in 5 years" - but reality drifts to 8.5 years and lakhs more interest, silently. Nobody shows you the path you're *actually* on vs. the one you promised.
* **The Solution**: 
  * **Goal Engine**: Constantly tracks target goals against actual savings rates and debt paydown schedules, explicitly flagging drift and projecting reality.

### 2.3. Windfalls dissolve invisibly.
* **The Problem**: A ₹20k bonus lands with no plan, gets absorbed into random spending, and disappears. Undirected money defaults to consumption.
* **The Solution**: 
  * **Rules & Automation Engines**: Detect transactions matching windfall characteristics and trigger smart alerts or automated allocation suggestions *before* the money is spent.

### 2.4. Lifestyle creep eats every raise.
* **The Problem**: Income grows, spending quietly grows to match, savings rate stays flat - you earn more but never get ahead.
* **The Solution**: 
  * **Recommendation Engine**: Tracks income increases against savings rates and alerts users to step up their automatic investment contributions proportionately to fight lifestyle creep.

---

## 3. Family & Life Problems

### 3.1. Money causes guilt and marital friction.
* **The Problem**: Enjoyment (dinners, trips, gifts) feels guilty; spending decisions become "spender vs brake" arguments between partners. There's no neutral referee.
* **The Solution**: 
  * **Household Domain**: Provides a neutral, data-backed planner where decisions are visualized mathematically, removing emotional bias.

### 3.2. Two earners, no shared picture.
* **The Problem**: Separate incomes but joint commitments (home loan, rent) - no fair, transparent, consolidated "are we okay?" view with clear contribution splits.
* **The Solution**: 
  * **Household Collaboration Module**: Supports joint budgets, shared goals, and access control profiles so partners have a single consolidated dashboard.

### 3.3. The family is financially blind if the earner is gone.
* **The Problem**: No single clear place holds the loans, insurance, nominees, and the plan - and no way for the family to *feel* that it was all planned for them.
* **The Solution**: 
  * **Centralized Dashboard & Permissions**: Aggregates all assets, accounts, liabilities, and plans into a single, secure environment accessible by designated household members via role-based permissions.

---

## 4. Blind-spot Problems (Things People Ignore)

### 4.1. Inflation silently makes "saving" a loss.
* **The Problem**: Money not growing >= 6% is losing value.
* **The Solution**: 
  * **Projections & Health Engines**: Factor inflation rates into all long-term projections and penalize excessive, underperforming cash holdings via the Health Score.

### 4.2. The earner is unprotected.
* **The Problem**: Investing is cheered while term/health insurance is missing (the #1 middle-class wealth-destroyer).
* **The Solution**: 
  * **Risk Engine**: Explicitly tracks insurance adequacy as one of its 18 core indicators, warning the user if vital safety nets are missing.

### 4.3. Big future costs arrive as a shock.
* **The Problem**: Kids' education inflation arrives unexpectedly, and there is no defined "freedom number" (finish line) to aim at.
* **The Solution**: 
  * **Goal Engine**: Allows defining milestones for high-inflation events (like education), calculating a dynamic funding pathway and projecting an ultimate "freedom number."
