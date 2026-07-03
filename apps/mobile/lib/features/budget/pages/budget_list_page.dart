import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../state/budget_state.dart';

class BudgetListPage extends ConsumerWidget {
  const BudgetListPage({super.key});
  
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncBudgets = ref.watch(budgetsProvider);
    final theme = Theme.of(context);
    
    return Scaffold(
      appBar: AppBar(
        title: const Text('Budgets'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => context.push('/planning/budget/add'),
          ),
        ],
      ),
      body: asyncBudgets.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (budgets) {
          if (budgets.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.account_balance_wallet_outlined, size: 48, color: Colors.grey),
                  const SizedBox(height: AppSpacing.md),
                  const Text('No budgets found'),
                  const SizedBox(height: AppSpacing.md),
                  ElevatedButton(
                    onPressed: () => context.push('/planning/budget/add'),
                    child: const Text('Create Budget'),
                  )
                ],
              ),
            );
          }
          return ListView.builder(
            padding: const EdgeInsets.all(AppSpacing.md),
            itemCount: budgets.length,
            itemBuilder: (context, index) {
              final b = budgets[index];
              return Card(
                margin: const EdgeInsets.only(bottom: AppSpacing.sm),
                child: ListTile(
                  title: Text(b.name, style: theme.textTheme.titleMedium),
                  subtitle: Text('${b.period} • ${b.status}'),
                  trailing: Text('₹${b.totalBudgeted}', style: theme.textTheme.titleSmall),
                  onTap: () => context.push('/planning/budget/${b.budgetId}'),
                ),
              );
            },
          );
        },
      ),
    );
  }
}
