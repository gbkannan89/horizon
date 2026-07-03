import 'package:flutter/material.dart';
import 'package:horizon_mobile/features/dashboard/models/dashboard_models.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';
import 'package:horizon_mobile/core/ui_kit/animated_stat.dart';

class HealthScoreCard extends StatelessWidget {
  final WidgetModel widget;
  final int score;
  final String grade;

  const HealthScoreCard({super.key, required this.widget, this.score = 0, this.grade = ''});

  @override
  Widget build(BuildContext context) {
    final color = score >= 60 ? AppColors.teal500 : (score >= 40 ? AppColors.amber500 : AppColors.red500);
    final theme = Theme.of(context);
    
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.favorite, color: color, size: 20),
            const SizedBox(width: 8),
            Text('Health Score', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          Row(
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              AnimatedStatValue(
                value: score.toDouble(),
                format: '#',
                style: theme.textTheme.displaySmall?.copyWith(fontWeight: FontWeight.bold, color: color),
              ),
              const SizedBox(width: 4),
              Padding(
                padding: const EdgeInsets.only(bottom: 6),
                child: Text('/100', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ),
              const Spacer(),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                decoration: BoxDecoration(
                  color: color.withOpacity(0.15),
                  borderRadius: BorderRadius.circular(AppRadius.sm),
                ),
                child: Text(grade, style: TextStyle(color: color, fontWeight: FontWeight.bold, fontSize: 13)),
              ),
            ],
          ),
        ],
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
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.account_balance_wallet, color: AppColors.teal500, size: 20),
            const SizedBox(width: 8),
            Text('Net Worth', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          AnimatedStatValue(
            value: netWorth.toDouble(),
            style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.w800, color: theme.colorScheme.onSurface),
          ),
          const SizedBox(height: 4),
          Text('Total financial position', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        ],
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
    final color = surplus >= 0 ? AppColors.teal500 : AppColors.red500;
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.swap_horiz, color: AppColors.amber500, size: 20),
            const SizedBox(width: 8),
            Text('Cash Flow', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          Text(_flowLine('Income', income), style: theme.textTheme.bodyMedium),
          Text(_flowLine('Expenses', expenses), style: theme.textTheme.bodyMedium),
          const Divider(height: 16, color: AppColors.slate200),
          Row(
            children: [
              Text('Net: ', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
              AnimatedStatValue(
                value: surplus.toDouble(),
                style: TextStyle(fontWeight: FontWeight.bold, color: color),
              ),
            ],
          ),
        ],
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
    final color = pct >= 0.8 ? AppColors.teal500 : (pct >= 0.5 ? AppColors.amber500 : AppColors.red500);
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.flag_rounded, color: color, size: 20),
            const SizedBox(width: 8),
            Text('Goal Progress', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          Text('$onTrack / $total on track', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w700)),
          const SizedBox(height: 12),
          ClipRRect(
            borderRadius: BorderRadius.circular(AppRadius.sm),
            child: LinearProgressIndicator(
              value: pct, 
              backgroundColor: color.withOpacity(0.15), 
              color: color, 
              minHeight: 6
            ),
          ),
        ],
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
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.pie_chart_rounded, color: AppColors.teal500, size: 20),
            const SizedBox(width: 8),
            Text('Portfolio', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          AnimatedStatValue(
            value: value.toDouble(),
            style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 4),
          Text('Total portfolio value', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        ],
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
    final color = score >= 80 ? AppColors.red500 : (score >= 60 ? AppColors.amber500 : AppColors.teal500);
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.shield_rounded, color: color, size: 20),
            const SizedBox(width: 8),
            Text('Risk', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          Row(
            children: [
              AnimatedStatValue(
                value: score.toDouble(),
                format: '#',
                style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold, color: color),
              ),
              const SizedBox(width: 8),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.sm)),
                child: Text(level, style: TextStyle(color: color, fontSize: 12, fontWeight: FontWeight.bold)),
              ),
            ],
          ),
        ],
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
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            const Icon(Icons.lightbulb_rounded, color: AppColors.amber500, size: 20),
            const SizedBox(width: 8),
            Text('Recommendation', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          Text(widget.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
          if (widget.numericValue.isNotEmpty) ...[
            const SizedBox(height: 4),
            Text(widget.numericValue, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ],
        ],
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
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.account_balance_rounded, color: AppColors.teal500, size: 20),
            const SizedBox(width: 8),
            Text('Accounts', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          Text('$count accounts', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          AnimatedStatValue(
            value: balance.toDouble(),
            style: theme.textTheme.bodyMedium,
          ),
        ],
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
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            const Icon(Icons.credit_card_rounded, color: AppColors.red500, size: 20),
            const SizedBox(width: 8),
            Text('Total Debt', style: theme.textTheme.titleSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ]),
          const SizedBox(height: 12),
          AnimatedStatValue(
            value: debt.toDouble(),
            style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold, color: AppColors.red500),
          ),
        ],
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
    return GlassCard(
      child: Row(children: [
        const Icon(Icons.emoji_events_rounded, color: AppColors.amber500, size: 20),
        const SizedBox(width: 12),
        Expanded(child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(widget.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
            if (widget.numericValue.isNotEmpty) Text(widget.numericValue, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ],
        )),
        Icon(Icons.chevron_right_rounded, color: theme.colorScheme.onSurfaceVariant),
      ]),
    );
  }
}

