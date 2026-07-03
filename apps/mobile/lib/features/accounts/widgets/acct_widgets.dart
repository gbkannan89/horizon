import 'package:flutter/material.dart';
import '../models/acct_models.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';
import 'package:horizon_mobile/core/ui_kit/animated_stat.dart';

String _fmt(int v) { if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr'; if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L'; return '₹$v'; }

Color _typeColor(String type) {
  switch (type) {
    case 'Savings': return AppColors.teal500;
    case 'Checking': return Colors.blue;
    case 'Credit': return Colors.purple;
    case 'Investment': return AppColors.teal500;
    case 'Loan': return AppColors.red500;
    case 'Insurance': return AppColors.amber500;
    default: return AppColors.slate500;
  }
}

Color _healthColor(String h) {
  switch (h) { 
    case 'healthy': 
    case 'good': return AppColors.teal500; 
    case 'warning': return AppColors.amber500; 
    case 'critical': return AppColors.red500; 
    default: return AppColors.slate500; 
  }
}

class AcctCard extends StatelessWidget {
  final AcctCardData acct; final VoidCallback? onTap;
  const AcctCard({super.key, required this.acct, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = _typeColor(acct.accountType);
    return GlassCard(
      onTap: onTap,
      child: Row(children: [
        Container(
          width: 44, height: 44,
          decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.md)),
          child: Icon(_typeIcon(acct.accountType), color: color, size: 24),
        ),
        const SizedBox(width: AppSpacing.md),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(children: [
                Expanded(child: Text(acct.accountName, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold), maxLines: 1, overflow: TextOverflow.ellipsis)),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                  decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.xs)),
                  child: Text(acct.accountType, style: TextStyle(color: color, fontSize: 10, fontWeight: FontWeight.bold)),
                ),
              ]),
              const SizedBox(height: 6),
              Row(children: [
                Expanded(child: Text('${AcctFmt(acct.currentBalance)}  •  ${acct.currency}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant, fontWeight: FontWeight.w600))),
                if (acct.institutionName != null)
                  Text(acct.institutionName!, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ]),
            ],
          ),
        ),
        const SizedBox(width: AppSpacing.sm),
        Icon(Icons.chevron_right_rounded, color: theme.colorScheme.onSurfaceVariant, size: 20),
      ]),
    );
  }

  IconData _typeIcon(String type) {
    switch (type) {
      case 'Savings': return Icons.savings_rounded;
      case 'Checking': return Icons.account_balance_rounded;
      case 'Credit': return Icons.credit_card_rounded;
      case 'Investment': return Icons.trending_up_rounded;
      case 'Loan': return Icons.receipt_rounded;
      default: return Icons.account_balance_rounded;
    }
  }
}

class AcctBalanceCard extends StatelessWidget {
  final int total; final int available; final int spendable;
  const AcctBalanceCard({super.key, required this.total, required this.available, required this.spendable});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Icon(Icons.account_balance_wallet_rounded, color: AppColors.teal500, size: 20),
              const SizedBox(width: 8),
              Text('Balance Summary', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
            ]
          ),
          const SizedBox(height: 16),
          Row(children: [
            _stat('Total', total.toDouble(), theme), 
            _stat('Available', available.toDouble(), theme), 
            _stat('Spendable', spendable.toDouble(), theme),
          ]),
        ],
      ),
    );
  }

  Widget _stat(String label, double value, ThemeData t) => Expanded(child: Column(children: [
    AnimatedStatValue(value: value, style: t.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w900)), 
    Text(label, style: t.textTheme.bodySmall?.copyWith(color: t.colorScheme.onSurfaceVariant)),
  ]));
}

class AcctCardWidget extends StatelessWidget {
  final String title; final String summary; final IconData icon; final Color color;
  const AcctCardWidget({super.key, required this.title, required this.summary, required this.icon, required this.color});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return GlassCard(
      child: Row(children: [
        Container(
          width: 44, height: 44, 
          decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.md)),
          child: Icon(icon, color: color, size: 24)
        ),
        const SizedBox(width: AppSpacing.md),
        Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Text(title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
          Text(summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        ])),
      ]),
    );
  }
}

String AcctFmt(int v) {
  if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr';
  if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L';
  return '₹$v';
}
