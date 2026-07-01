import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';

final budgetProvider = FutureProvider.autoDispose<Map<String, dynamic>>((ref) async {
  final r = await ref.read(apiClientProvider).get('/planning/budget');
  final json = r.data as Map<String, dynamic>;
  return json['data'] as Map<String, dynamic>? ?? {};
});

class BudgetPage extends ConsumerWidget {
  const BudgetPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(budgetProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Budget Planning')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (data) {
          final categories = data['categories'] as List? ?? [];
          return ListView(
            padding: const EdgeInsets.all(AppSpacing.md),
            children: [
              Card(child: Padding(
                padding: const EdgeInsets.all(AppSpacing.md),
                child: Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: [
                  _stat(theme, 'Budget', '₹${data['total_budget']}', theme.colorScheme.primary),
                  _stat(theme, 'Spent', '₹${data['total_spent']}', data['surplus'] >= 0 ? Colors.green : Colors.red),
                  _stat(theme, 'Left', '₹${data['surplus']}', data['surplus'] >= 0 ? Colors.green : Colors.red),
                ]),
              )),
              const SizedBox(height: AppSpacing.md),
              ...categories.map((c) => Card(
                margin: const EdgeInsets.only(bottom: AppSpacing.sm),
                child: ListTile(
                  title: Text(c['category'] ?? ''),
                  trailing: Text('₹${c['actual']} / ₹${c['planned']}'),
                  subtitle: LinearProgressIndicator(value: ((c['actual'] as num?)?.toDouble() ?? 0) / ((c['planned'] as num)?.toDouble() ?? 1)),
                ),
              )),
            ],
          );
        },
      ),
    );
  }
  Widget _stat(ThemeData t, String label, String value, Color color) => Column(children: [
    Text(value, style: t.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold, color: color)),
    Text(label, style: t.textTheme.labelSmall?.copyWith(color: t.colorScheme.onSurfaceVariant)),
  ]);
}
