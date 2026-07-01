import 'package:flutter/material.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/acct_models.dart';

Color _typeColor(String type) {
  switch (type) {
    case 'Savings': return Colors.green;
    case 'Checking': return Colors.blue;
    case 'Credit': return Colors.purple;
    case 'Investment': return Colors.teal;
    case 'Loan': return Colors.red;
    case 'Insurance': return Colors.orange;
    default: return Colors.grey;
  }
}

class AcctCard extends StatelessWidget {
  final AcctCardData acct; final VoidCallback? onTap;
  const AcctCard({super.key, required this.acct, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = _typeColor(acct.accountType);
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(14),
          child: Row(children: [
            Container(
              width: 40, height: 40,
              decoration: BoxDecoration(color: color.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(10)),
              child: Icon(_typeIcon(acct.accountType), color: color, size: 20),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(children: [
                    Expanded(child: Text(acct.accountName, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600), maxLines: 1, overflow: TextOverflow.ellipsis)),
                    StatusChip(label: acct.accountType, color: color, fontSize: 10),
                  ]),
                  const SizedBox(height: 4),
                  Row(children: [
                    Expanded(child: Text('${formatMoney(acct.currentBalance)}  •  ${acct.currency}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant))),
                    if (acct.institutionName != null)
                      Text(acct.institutionName!, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  ]),
                ],
              ),
            ),
            const SizedBox(width: 4),
            Icon(Icons.chevron_right, color: theme.colorScheme.onSurfaceVariant, size: 20),
          ]),
        ),
      ),
    );
  }

  IconData _typeIcon(String type) {
    switch (type) {
      case 'Savings': return Icons.savings;
      case 'Checking': return Icons.account_balance;
      case 'Credit': return Icons.credit_card;
      case 'Investment': return Icons.trending_up;
      case 'Loan': return Icons.receipt;
      default: return Icons.account_balance;
    }
  }
}

class AcctBalanceCard extends StatelessWidget {
  final int total; final int available; final int spendable;
  const AcctBalanceCard({super.key, required this.total, required this.available, required this.spendable});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: AppTheme.cardPadding,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Balance Summary', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: AppTheme.spacingMd),
            Row(children: [
              _stat('Total', formatMoney(total), theme),
              _stat('Available', formatMoney(available), theme),
              _stat('Spendable', formatMoney(spendable), theme),
            ]),
          ],
        ),
      ),
    );
  }

  Widget _stat(String label, String value, ThemeData t) => Expanded(child: Column(children: [
    Text(value, style: t.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
    Text(label, style: t.textTheme.bodySmall),
  ]));
}

class AcctCardWidget extends StatelessWidget {
  final String title; final String summary; final IconData icon; final Color color;
  const AcctCardWidget({super.key, required this.title, required this.summary, required this.icon, required this.color});

  @override
  Widget build(BuildContext context) {
    return IconCard(title: title, subtitle: summary, icon: icon, color: color);
  }
}
