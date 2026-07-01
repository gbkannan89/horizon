import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../models/transaction_models.dart';
import '../repository/transaction_repository.dart';
import '../widgets/transaction_widgets.dart';

final transactionDetailProvider = FutureProvider.autoDispose.family<TransactionEvent?, String>((ref, id) async {
  final repo = ref.read(transactionRepositoryProvider);
  final resp = await repo.getTransaction(id: id);
  return resp.data;
});

class TransactionDetailPage extends ConsumerWidget {
  final String transactionId;

  const TransactionDetailPage({super.key, required this.transactionId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(transactionDetailProvider(transactionId));
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Transaction'),
        actions: [
          PopupMenuButton<String>(
            onSelected: (value) async {
              if (value == 'archive') {
                try {
                  final repo = ref.read(transactionRepositoryProvider);
                  await repo.archiveTransaction(id: transactionId);
                  if (context.mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Transaction archived'), behavior: SnackBarBehavior.floating),
                    );
                    context.pop();
                  }
                } catch (e) {
                  if (context.mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(content: Text('Failed to archive: $e'), behavior: SnackBarBehavior.floating),
                    );
                  }
                }
              }
            },
            itemBuilder: (_) => [
              const PopupMenuItem(value: 'archive', child: ListTile(leading: Icon(Icons.archive), title: Text('Archive'))),
            ],
          ),
        ],
      ),
      body: async.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(Icons.error_outline, size: 48, color: theme.colorScheme.error),
              const SizedBox(height: 16),
              Text('Could not load transaction', style: theme.textTheme.titleMedium),
              const SizedBox(height: 8),
              Text(e.toString(), style: theme.textTheme.bodySmall, textAlign: TextAlign.center),
            ],
          ),
        ),
        data: (tx) {
          if (tx == null) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.search_off, size: 48, color: theme.colorScheme.onSurfaceVariant),
                  const SizedBox(height: 16),
                  Text('Transaction not found', style: theme.textTheme.titleMedium),
                ],
              ),
            );
          }
          return _buildDetail(context, theme, tx);
        },
      ),
    );
  }

  Widget _buildDetail(BuildContext context, ThemeData theme, TransactionEvent tx) {
    final isExpense = tx.isExpense;
    final amountColor = isExpense ? theme.colorScheme.error : theme.colorScheme.primary;

    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        Card(
          child: Padding(
            padding: const EdgeInsets.all(24),
            child: Column(
              children: [
                Container(
                  width: 56, height: 56,
                  decoration: BoxDecoration(
                    color: isExpense ? theme.colorScheme.errorContainer : theme.colorScheme.primaryContainer,
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Icon(
                    isExpense ? Icons.arrow_upward : Icons.arrow_downward,
                    size: 28,
                    color: amountColor,
                  ),
                ),
                const SizedBox(height: 16),
                Text(
                  '${isExpense ? '-' : '+'}${tx.amount.toStringAsFixed(2)}',
                  style: theme.textTheme.headlineMedium?.copyWith(
                    fontWeight: FontWeight.bold,
                    color: amountColor,
                  ),
                ),
                const SizedBox(height: 4),
                Text(tx.currency, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                if (tx.description != null && tx.description!.isNotEmpty) ...[
                  const SizedBox(height: 12),
                  Text(tx.description!, style: theme.textTheme.titleMedium, textAlign: TextAlign.center),
                ],
                const SizedBox(height: 12),
                TransactionStatusBadge(state: tx.state),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                _detailRow(theme, 'Category', tx.category ?? tx.eventType),
                const Divider(height: 24),
                _detailRow(theme, 'Date', tx.formattedDate),
                const Divider(height: 24),
                _detailRow(theme, 'Type', tx.eventType),
                if (tx.source != null) ...[
                  const Divider(height: 24),
                  _detailRow(theme, 'Source', tx.source!),
                ],
                if (tx.destination != null) ...[
                  const Divider(height: 24),
                  _detailRow(theme, 'Destination', tx.destination!),
                ],
                if (tx.reference != null) ...[
                  const Divider(height: 24),
                  _detailRow(theme, 'Reference', tx.reference!),
                ],
                if (tx.origin != null) ...[
                  const Divider(height: 24),
                  _detailRow(theme, 'Origin', tx.origin!),
                ],
                if (tx.confidence != null) ...[
                  const Divider(height: 24),
                  _detailRow(theme, 'Confidence', tx.confidence!),
                ],
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Notes', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 8),
                Text(
                  tx.notes ?? 'No notes added',
                  style: theme.textTheme.bodyMedium?.copyWith(
                    color: tx.notes == null ? theme.colorScheme.onSurfaceVariant : null,
                    fontStyle: tx.notes == null ? FontStyle.italic : FontStyle.normal,
                  ),
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Attachments', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 8),
                Row(
                  children: [
                    Icon(Icons.attach_file, size: 16, color: theme.colorScheme.onSurfaceVariant),
                    const SizedBox(width: 8),
                    Text(
                      'Not yet supported',
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: theme.colorScheme.onSurfaceVariant,
                        fontStyle: FontStyle.italic,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 16),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Tags', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 8),
                if (tx.tags != null && tx.tags!.isNotEmpty)
                  Wrap(
                    spacing: 8, runSpacing: 4,
                    children: tx.tags!.map((tag) => Chip(
                      label: Text(tag, style: theme.textTheme.labelSmall),
                      visualDensity: VisualDensity.compact,
                      materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                    )).toList(),
                  )
                else
                  Text(
                    'No tags',
                    style: theme.textTheme.bodyMedium?.copyWith(
                      color: theme.colorScheme.onSurfaceVariant,
                      fontStyle: FontStyle.italic,
                    ),
                  ),
              ],
            ),
          ),
        ),
      ],
    );
  }

  Widget _detailRow(ThemeData theme, String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        Flexible(
          child: Text(
            value,
            style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500),
            textAlign: TextAlign.right,
          ),
        ),
      ],
    );
  }
}
