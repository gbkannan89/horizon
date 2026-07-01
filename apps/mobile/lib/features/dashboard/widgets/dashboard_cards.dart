import 'package:flutter/material.dart';
import 'package:horizon_mobile/features/dashboard/models/dashboard_models.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';

class HealthScoreCard extends StatelessWidget {
  final WidgetModel widget;
  final int score;
  final String grade;

  const HealthScoreCard({super.key, required this.widget, this.score = 0, this.grade = ''});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = AppTheme.healthColor(score);
    return Card(
      child: Padding(
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.favorite, color: color, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Health Score', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
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
                StatusChip(label: grade, color: color, fontSize: 12),
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
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.account_balance_wallet, color: theme.colorScheme.primary, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Net Worth', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
            Text(formatMoney(netWorth), style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: AppTheme.spacingXs),
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
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.swap_horiz, color: theme.colorScheme.primary, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Cash Flow', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
            Text('Income: ${formatMoney(income)}', style: theme.textTheme.bodyMedium),
            Text('Expenses: ${formatMoney(expenses)}', style: theme.textTheme.bodyMedium),
            const Divider(height: 16),
            Row(
              children: [
                Text('Net: ', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                Text(formatMoney(surplus), style: TextStyle(fontWeight: FontWeight.bold, color: color)),
              ],
            ),
          ],
        ),
      ),
    );
  }
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
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.flag, color: color, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Goal Progress', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
            Text('$onTrack / $total on track', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: AppTheme.spacingSm),
            ClipRRect(
              borderRadius: BorderRadius.circular(AppTheme.radiusSm),
              child: LinearProgressIndicator(value: pct, backgroundColor: color.withOpacity(0.1), color: color, minHeight: 8),
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
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.pie_chart, color: theme.colorScheme.primary, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Portfolio', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
            Text(formatMoney(value), style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: AppTheme.spacingXs),
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
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.shield, color: color, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Risk', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
            Row(
              children: [
                Text('$score', style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold, color: color)),
                const SizedBox(width: 8),
                StatusChip(label: level, color: color, fontSize: 12),
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
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.lightbulb, color: Colors.amber, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Recommendation', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
            Text(widget.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
            if (widget.numericValue.isNotEmpty) ...[
              const SizedBox(height: AppTheme.spacingXs),
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
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.account_balance, color: theme.colorScheme.primary, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Accounts', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
            Text('$count accounts', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            Text(formatMoney(balance), style: theme.textTheme.bodyMedium),
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
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.credit_card, color: Colors.purple, size: AppTheme.iconMd),
              const SizedBox(width: 8),
              Text('Total Debt', style: theme.textTheme.titleSmall),
            ]),
            const SizedBox(height: AppTheme.spacingMd),
            Text(formatMoney(debt), style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold, color: Colors.purple)),
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
        padding: AppTheme.cardPadding,
        child: Row(children: [
          Icon(Icons.flag, color: Colors.green, size: AppTheme.iconMd),
          const SizedBox(width: AppTheme.spacingMd),
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
