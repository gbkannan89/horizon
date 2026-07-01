import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';

final debtProvider = FutureProvider.autoDispose<Map<String, dynamic>>((ref) async {
  final r = await ref.read(apiClientProvider).get('/planning/debt-payoff');
  return (r.data as Map<String, dynamic>)['data'] as Map<String, dynamic>? ?? {};
});

class DebtPayoffPage extends ConsumerWidget {
  const DebtPayoffPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(debtProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Debt Payoff')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (d) {
          final debts = d['debts'] as List? ?? [];
          return ListView(padding: const EdgeInsets.all(AppSpacing.md), children: [
            Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: Column(children: [
              Text('₹${d['total_debt']}', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold, color: Colors.red)),
              Text('Total Debt · DTI ${d['debt_to_income_ratio']}%', style: theme.textTheme.bodySmall),
              Text('Debt-free by ${d['debt_free_date']}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.primary)),
            ]))),
            ...debts.map((debt) => Card(margin: const EdgeInsets.only(bottom: AppSpacing.sm), child: ListTile(
              title: Text(debt['name'] ?? ''),
              subtitle: Text('₹${debt['principal']} @ ${debt['interest_rate']}%'),
              trailing: Text('₹${debt['monthly_emi']}/mo'),
            ))),
          ]);
        },
      ),
    );
  }
}
