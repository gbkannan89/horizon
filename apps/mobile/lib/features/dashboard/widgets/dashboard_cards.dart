import 'package:flutter/material.dart';
import 'package:horizon_mobile/features/dashboard/models/dashboard_models.dart';

class HealthScoreCard extends StatelessWidget {
  final WidgetModel widget;
  final int score;
  final String grade;

  const HealthScoreCard({super.key, required this.widget, this.score = 0, this.grade = ''});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = score >= 60 ? Colors.green : (score >= 40 ? Colors.orange : Colors.red);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.favorite, color: color, size: 20),
              const SizedBox(width: 8),
              Text('Health Score', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Row(
              crossAxisAlignment: CrossAxisAlignment.end,
              children: [
                Text('$score', style: theme.textTheme.displaySmall?.copyWith(fontWeight: FontWeight.bold, color: color)),
                const SizedBox(width: 8),
                Padding(
                  padding: const EdgeInsets.only(bottom: 4),
                  child: Text('/100', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ),
                const Spacer(),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: color.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(grade, style: TextStyle(color: color, fontWeight: FontWeight.w600, fontSize: 12)),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class NetWorthCard extends StatelessWidget {
  final WidgetModel widget;
  final int netWorth;
  const NetWorthCard({super.key, required this.widget, this.netWorth = 0});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.account_balance_wallet, color: theme.colorScheme.primary, size: 20),
              const SizedBox(width: 8),
              Text('Net Worth', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Text(WidgetModel.formatMoney(netWorth), style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 4),
            Text('Total financial position', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ],
        ),
      ),
    );
  }
}

class CashFlowCard extends StatelessWidget {
  final int income;
  final int expenses;
  final WidgetModel? widget;
  const CashFlowCard({super.key, this.income = 0, this.expenses = 0, this.widget});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final surplus = income - expenses;
    final color = surplus >= 0 ? Colors.green : Colors.red;
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.swap_horiz, color: theme.colorScheme.primary, size: 20),
              const SizedBox(width: 8),
              Text('Cash Flow', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Text(_flowLine('Income', income), style: theme.textTheme.bodyMedium),
            Text(_flowLine('Expenses', expenses), style: theme.textTheme.bodyMedium),
            const Divider(height: 16),
            Row(
              children: [
                Text('Net: ', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                Text(WidgetModel.formatMoney(surplus), style: TextStyle(fontWeight: FontWeight.bold, color: color)),
              ],
            ),
          ],
        ),
      ),
    );
  }

  String _flowLine(String label, int amount) => '$label: ${WidgetModel.formatMoney(amount)}';
}

class GoalProgressCard extends StatelessWidget {
  final int onTrack;
  final int total;
  final WidgetModel? widget;
  const GoalProgressCard({super.key, this.onTrack = 0, this.total = 0, this.widget});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final pct = total > 0 ? onTrack / total : 0.0;
    final color = pct >= 0.8 ? Colors.green : (pct >= 0.5 ? Colors.orange : Colors.red);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.flag, color: color, size: 20),
              const SizedBox(width: 8),
              Text('Goal Progress', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Text('$onTrack / $total on track', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(value: pct, backgroundColor: color.withValues(alpha: 0.1), color: color, minHeight: 8),
            ),
          ],
        ),
      ),
    );
  }
}

class PortfolioCard extends StatelessWidget {
  final int value;
  final WidgetModel? widget;
  const PortfolioCard({super.key, this.value = 0, this.widget});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.pie_chart, color: theme.colorScheme.primary, size: 20),
              const SizedBox(width: 8),
              Text('Portfolio', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Text(WidgetModel.formatMoney(value), style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 4),
            Text('Total portfolio value', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ],
        ),
      ),
    );
  }
}

class RiskScoreCard extends StatelessWidget {
  final int score;
  final String level;
  final WidgetModel? widget;
  const RiskScoreCard({super.key, this.score = 0, this.level = '', this.widget});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = score >= 80 ? Colors.red : (score >= 60 ? Colors.orange : Colors.green);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.shield, color: color, size: 20),
              const SizedBox(width: 8),
              Text('Risk', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Row(
              children: [
                Text('$score', style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold, color: color)),
                const SizedBox(width: 8),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(8)),
                  child: Text(level, style: TextStyle(color: color, fontSize: 12, fontWeight: FontWeight.w600)),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class RecommendationCard extends StatelessWidget {
  final WidgetModel widget;
  const RecommendationCard({super.key, required this.widget});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              const Icon(Icons.lightbulb, color: Colors.amber, size: 20),
              const SizedBox(width: 8),
              Text('Recommendation', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Text(widget.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
            if (widget.numericValue.isNotEmpty) ...[
              const SizedBox(height: 4),
              Text(widget.numericValue, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            ],
          ],
        ),
      ),
    );
  }
}

class AccountsCard extends StatelessWidget {
  final int count;
  final int balance;
  final WidgetModel? widget;
  const AccountsCard({super.key, this.count = 0, this.balance = 0, this.widget});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.account_balance, color: theme.colorScheme.primary, size: 20),
              const SizedBox(width: 8),
              Text('Accounts', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Text('$count accounts', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            Text(WidgetModel.formatMoney(balance), style: theme.textTheme.bodyMedium),
          ],
        ),
      ),
    );
  }
}

class DebtSummaryCard extends StatelessWidget {
  final int debt;
  final WidgetModel? widget;
  const DebtSummaryCard({super.key, this.debt = 0, this.widget});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    if (debt <= 0) return const SizedBox.shrink();
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              const Icon(Icons.credit_card, color: Colors.purple, size: 20),
              const SizedBox(width: 8),
              Text('Total Debt', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: 12),
            Text(WidgetModel.formatMoney(debt), style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold, color: Colors.purple)),
          ],
        ),
      ),
    );
  }
}

class MilestoneCard extends StatelessWidget {
  final WidgetModel widget;
  const MilestoneCard({super.key, required this.widget});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Row(children: [
          const Icon(Icons.flag, color: Colors.green, size: 20),
          const SizedBox(width: 12),
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(widget.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
              if (widget.numericValue.isNotEmpty) Text(widget.numericValue, style: theme.textTheme.bodySmall),
            ],
          )),
          Icon(Icons.chevron_right, color: theme.colorScheme.onSurfaceVariant),
        ]),
      ),
    );
  }
}
