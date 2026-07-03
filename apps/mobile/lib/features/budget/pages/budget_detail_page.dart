import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../state/budget_state.dart';
import '../widgets/budget_widgets.dart';

class BudgetDetailPage extends ConsumerWidget {
  final String budgetId;
  const BudgetDetailPage({super.key, required this.budgetId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncData = ref.watch(budgetDetailProvider(budgetId));
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Budget Detail')),
      body: asyncData.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (budget) {
          if (budget == null) return const Center(child: Text('Not found'));

          return ListView(
            padding: const EdgeInsets.all(AppSpacing.md),
            children: [
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Expanded(child: Text(budget.name, style: theme.textTheme.headlineSmall)),
                          BudgetStatusBadge(status: budget.status),
                        ],
                      ),
                      const SizedBox(height: AppSpacing.sm),
                      Text('${budget.period} • ${budget.startDate} to ${budget.endDate}', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                      const SizedBox(height: AppSpacing.md),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          _stat(theme, 'Budgeted', '₹${budget.totalBudgeted.toStringAsFixed(0)}', Colors.blue),
                          _stat(theme, 'Spent', '₹${budget.totalSpent.toStringAsFixed(0)}', budget.totalSpent > budget.totalBudgeted ? Colors.red : Colors.green),
                          _stat(theme, 'Remaining', '₹${budget.totalRemaining.toStringAsFixed(0)}', budget.totalRemaining >= 0 ? Colors.green : Colors.red),
                        ],
                      ),
                      const SizedBox(height: AppSpacing.md),
                      BudgetProgressBar(
                        spent: budget.totalSpent,
                        budgeted: budget.totalBudgeted,
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: AppSpacing.md),
              const SharedSectionHeader(title: 'Budget vs Actual'),
              const SizedBox(height: AppSpacing.sm),
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: VsActualChart(categories: budget.categories),
                ),
              ),
              const SizedBox(height: AppSpacing.md),
              SharedSectionHeader(title: 'Categories (${budget.categories.length})'),
              const SizedBox(height: AppSpacing.sm),
              ...budget.categories.map((c) => Card(
                margin: const EdgeInsets.only(bottom: AppSpacing.sm),
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Expanded(child: Text(c.category, style: theme.textTheme.titleSmall)),
                          Text('₹${c.spentAmount.toStringAsFixed(0)} / ₹${c.budgetedAmount.toStringAsFixed(0)}', style: theme.textTheme.bodySmall),
                        ],
                      ),
                      const SizedBox(height: AppSpacing.sm),
                      BudgetProgressBar(spent: c.spentAmount, budgeted: c.budgetedAmount),
                      const SizedBox(height: 4),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text('Remaining: ₹${c.remainingAmount.toStringAsFixed(0)}', style: theme.textTheme.labelSmall),
                          Text('${c.spentPct.toStringAsFixed(1)}%', style: theme.textTheme.labelSmall),
                        ],
                      ),
                    ],
                  ),
                ),
              )),
              if (budget.categories.isNotEmpty) const SizedBox(height: AppSpacing.md),
              _actionButtons(context, ref, budget),
              const SizedBox(height: AppSpacing.xl),
            ],
          );
        },
      ),
    );
  }

  Widget _actionButtons(BuildContext context, WidgetRef ref, dynamic budget) {
    final status = budget.status as String;
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        if (status == 'Draft')
          _actionButton(context, ref, 'Activate Budget', Icons.check_circle, Colors.green, BudgetAction.activate),
        if (status == 'Active') ...[
          _actionButton(context, ref, 'Pause Budget', Icons.pause_circle, Colors.orange, BudgetAction.pause),
          _actionButton(context, ref, 'Complete Budget', Icons.task_alt, Colors.blue, BudgetAction.complete),
        ],
        if (status == 'Paused')
          _actionButton(context, ref, 'Resume Budget', Icons.play_circle, Colors.green, BudgetAction.resume),
        if (status == 'Completed' || status == 'Paused' || status == 'Draft')
          _actionButton(context, ref, 'Archive Budget', Icons.archive, Colors.grey, BudgetAction.archive),
      ],
    );
  }

  Widget _actionButton(BuildContext context, WidgetRef ref, String label, IconData icon, Color color, BudgetAction action) {
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: OutlinedButton.icon(
        onPressed: () async {
          await ref.read(budgetTransitionProvider(TransitionParams(budgetId: budgetId, action: action)).future);
          ref.invalidate(budgetDetailProvider(budgetId));
          ref.invalidate(budgetsProvider);
          if (context.mounted) {
            ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('$label completed')));
            context.pop();
          }
        },
        icon: Icon(icon, color: color),
        label: Text(label, style: TextStyle(color: color)),
        style: OutlinedButton.styleFrom(side: BorderSide(color: color)),
      ),
    );
  }

  Widget _stat(ThemeData t, String label, String value, Color color) {
    return Column(
      children: [
        Text(value, style: t.textTheme.titleMedium?.copyWith(color: color, fontWeight: FontWeight.bold)),
        Text(label, style: t.textTheme.labelSmall?.copyWith(color: t.colorScheme.onSurfaceVariant)),
      ],
    );
  }
}
