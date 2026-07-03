import 'package:flutter/material.dart';
import '../models/transaction_models.dart';

class TransactionCard extends StatelessWidget {
  final TransactionEvent transaction;
  final VoidCallback? onTap;

  const TransactionCard({super.key, required this.transaction, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isExpense = transaction.isExpense;
    final amountColor = isExpense ? theme.colorScheme.error : theme.colorScheme.primary;

    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
          child: Row(
            children: [
              Container(
                width: 40, height: 40,
                decoration: BoxDecoration(
                  color: isExpense
                      ? theme.colorScheme.errorContainer
                      : theme.colorScheme.primaryContainer,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(
                  _iconForType(transaction.eventType),
                  size: 20,
                  color: isExpense ? theme.colorScheme.error : theme.colorScheme.primary,
                ),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Flexible(
                          child: Text(
                            transaction.description ?? transaction.eventType,
                            style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500),
                            maxLines: 1,
                            overflow: TextOverflow.ellipsis,
                          ),
                        ),
                        if (transaction.tags?.contains('recurring') == true)
                          Padding(
                            padding: const EdgeInsets.only(left: 6),
                            child: Icon(Icons.repeat, size: 14, color: theme.colorScheme.primary),
                          ),
                      ],
                    ),
                    const SizedBox(height: 2),
                    Text(
                      '${transaction.category ?? transaction.eventType} · ${transaction.formattedDate}',
                      style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                      maxLines: 1,
                      overflow: TextOverflow.ellipsis,
                    ),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Text(
                    '${isExpense ? '-' : '+'}${_formatAmount(transaction.amount.abs())}',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      fontWeight: FontWeight.w600,
                      color: amountColor,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    transaction.currency,
                    style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  IconData _iconForType(String eventType) {
    switch (eventType.toLowerCase()) {
      case 'salary': case 'income': return Icons.trending_up;
      case 'purchase': case 'expense': return Icons.shopping_cart_outlined;
      case 'bill': case 'utility': return Icons.receipt_long_outlined;
      case 'transfer': return Icons.swap_horiz;
      case 'investment': return Icons.trending_up;
      case 'withdrawal': return Icons.phone_iphone;
      default: return Icons.receipt_outlined;
    }
  }

  String _formatAmount(double amount) {
    if (amount >= 10000000) return '${(amount / 10000000).toStringAsFixed(1)}Cr';
    if (amount >= 100000) return '${(amount / 100000).toStringAsFixed(1)}L';
    if (amount >= 1000) return '${(amount / 1000).toStringAsFixed(1)}K';
    return amount.toStringAsFixed(0);
  }
}

class TransactionSummaryCard extends StatelessWidget {
  final TransactionSummary summary;

  const TransactionSummaryCard({super.key, required this.summary});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('This Month', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _summaryItem(
                    theme, 'Income', summary.periodIncome,
                    theme.colorScheme.primary, Icons.arrow_downward,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: _summaryItem(
                    theme, 'Expenses', summary.periodExpenses,
                    theme.colorScheme.error, Icons.arrow_upward,
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Divider(color: theme.colorScheme.outlineVariant),
            const SizedBox(height: 12),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('Net Flow', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500)),
                Text(
                  '${summary.netFlow >= 0 ? '+' : ''}${summary.netFlow.toStringAsFixed(0)}',
                  style: theme.textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: summary.netFlow >= 0 ? theme.colorScheme.primary : theme.colorScheme.error,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _summaryItem(ThemeData t, String label, double amount, Color color, IconData icon) {
    return Row(
      children: [
        Container(
          padding: const EdgeInsets.all(6),
          decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(8)),
          child: Icon(icon, size: 16, color: color),
        ),
        const SizedBox(width: 8),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(label, style: t.textTheme.bodySmall?.copyWith(color: t.colorScheme.onSurfaceVariant)),
              Text(amount.toStringAsFixed(0), style: t.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600, color: color)),
            ],
          ),
        ),
      ],
    );
  }
}

class TransactionStatusBadge extends StatelessWidget {
  final String state;

  const TransactionStatusBadge({super.key, required this.state});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final (Color color, String label) = switch (state.toLowerCase()) {
      'posted' => (theme.colorScheme.primary, 'Posted'),
      'pending' => (Colors.orange, 'Pending'),
      'draft' => (theme.colorScheme.onSurfaceVariant, 'Draft'),
      'confirmed' => (Colors.green, 'Confirmed'),
      'cancelled' => (theme.colorScheme.error, 'Cancelled'),
      'reversed' => (Colors.red, 'Reversed'),
      'archived' => (theme.colorScheme.onSurfaceVariant, 'Archived'),
      _ => (theme.colorScheme.onSurfaceVariant, state),
    };

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Text(label, style: theme.textTheme.labelSmall?.copyWith(color: color, fontWeight: FontWeight.w500)),
    );
  }
}
