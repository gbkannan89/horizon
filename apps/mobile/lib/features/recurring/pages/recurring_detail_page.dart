import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/recurring_models.dart';
import '../repository/recurring_repository.dart';
import '../state/recurring_state.dart';

class RecurringDetailPage extends ConsumerWidget {
  final String id;
  const RecurringDetailPage({super.key, required this.id});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncData = ref.watch(recurringDetailProvider(id));
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Recurring Transaction')),
      body: asyncData.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (item) {
          if (item == null) return const Center(child: Text('Not found'));
          return ListView(
            padding: const EdgeInsets.all(AppSpacing.md),
            children: [
              _headerCard(context, ref, item),
              const SizedBox(height: AppSpacing.md),
              _detailsCard(theme, item),
              const SizedBox(height: AppSpacing.md),
              _scheduleCard(theme, item),
              const SizedBox(height: AppSpacing.md),
              _actionButtons(context, ref, item),
              const SizedBox(height: AppSpacing.xl),
            ],
          );
        },
      ),
    );
  }

  Widget _headerCard(BuildContext context, WidgetRef ref, RecurringTransaction item) {
    final theme = Theme.of(context);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(item.name, style: theme.textTheme.headlineSmall),
            if (item.description.isNotEmpty) ...[
              const SizedBox(height: AppSpacing.sm),
              Text(item.description, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            ],
            const SizedBox(height: AppSpacing.md),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                _stat(theme, 'Amount', '₹${item.amount.toStringAsFixed(0)}', Colors.blue),
                _stat(theme, 'Frequency', item.frequency, Colors.teal),
                _stat(theme, 'Status', item.status, _statusColor(item.status)),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _detailsCard(ThemeData theme, RecurringTransaction item) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Details', style: theme.textTheme.titleMedium),
            const Divider(),
            _detailRow(theme, 'Currency', item.currency),
            _detailRow(theme, 'Interval', 'Every ${item.interval} ${item.frequency.toLowerCase()}'),
            if (item.startDate.isNotEmpty) _detailRow(theme, 'Start Date', item.startDate),
            if (item.endDate != null) _detailRow(theme, 'End Date', item.endDate!),
          ],
        ),
      ),
    );
  }

  Widget _scheduleCard(ThemeData theme, RecurringTransaction item) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Schedule', style: theme.textTheme.titleMedium),
            const Divider(),
            if (item.nextOccurrence != null) ...[
              Row(children: [
                const Icon(Icons.schedule, size: 18, color: Colors.green),
                const SizedBox(width: AppSpacing.sm),
                Text('Next: ${item.nextOccurrence}', style: theme.textTheme.bodyMedium),
              ]),
            ] else ...[
              Row(children: [
                const Icon(Icons.check_circle, size: 18, color: Colors.grey),
                const SizedBox(width: AppSpacing.sm),
                Text('No upcoming occurrences', style: theme.textTheme.bodyMedium),
              ]),
            ],
          ],
        ),
      ),
    );
  }

  Widget _actionButtons(BuildContext context, WidgetRef ref, RecurringTransaction item) {
    final repo = ref.read(recurringRepositoryProvider);
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (item.status == 'Active') ...[
          _actionBtn(context, 'Pause', Icons.pause, Colors.orange, () async {
            await repo.pauseRecurring(item.id);
            ref.invalidate(recurringDetailProvider(item.id));
            ref.invalidate(recurringListProvider);
            if (context.mounted) context.pop();
          }),
          _actionBtn(context, 'Cancel', Icons.cancel, Colors.red, () async {
            await repo.cancelRecurring(item.id);
            ref.invalidate(recurringDetailProvider(item.id));
            ref.invalidate(recurringListProvider);
            if (context.mounted) context.pop();
          }),
        ],
        if (item.status == 'Paused')
          _actionBtn(context, 'Activate', Icons.play_arrow, Colors.green, () async {
            await repo.activateRecurring(item.id);
            ref.invalidate(recurringDetailProvider(item.id));
            ref.invalidate(recurringListProvider);
            if (context.mounted) context.pop();
          }),
        if (item.status == 'Cancelled' || item.status == 'Completed')
          _actionBtn(context, 'Archive', Icons.archive, Colors.grey, () async {
            await repo.archiveRecurring(item.id);
            ref.invalidate(recurringDetailProvider(item.id));
            ref.invalidate(recurringListProvider);
            if (context.mounted) context.pop();
          }),
      ],
    );
  }

  Widget _actionBtn(BuildContext context, String label, IconData icon, Color color, VoidCallback onTap) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: OutlinedButton.icon(
        onPressed: onTap,
        icon: Icon(icon, color: color),
        label: Text(label, style: TextStyle(color: color)),
        style: OutlinedButton.styleFrom(side: BorderSide(color: color)),
      ),
    );
  }

  Widget _stat(ThemeData t, String label, String value, Color color) {
    return Column(children: [
      Text(value, style: t.textTheme.titleMedium?.copyWith(color: color, fontWeight: FontWeight.bold)),
      Text(label, style: t.textTheme.labelSmall?.copyWith(color: t.colorScheme.onSurfaceVariant)),
    ]);
  }

  Widget _detailRow(ThemeData t, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: t.textTheme.bodyMedium?.copyWith(color: t.colorScheme.onSurfaceVariant)),
          Text(value, style: t.textTheme.bodyMedium),
        ],
      ),
    );
  }

  Color _statusColor(String status) {
    switch (status) {
      case 'Active': return Colors.green;
      case 'Paused': return Colors.orange;
      case 'Cancelled': return Colors.red;
      case 'Completed': return Colors.blue;
      case 'Archived': return Colors.grey;
      default: return Colors.grey;
    }
  }
}
