import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';

final efProvider = FutureProvider.autoDispose<Map<String, dynamic>>((ref) async {
  final r = await ref.read(apiClientProvider).get('/planning/emergency-fund');
  return (r.data as Map<String, dynamic>)['data'] as Map<String, dynamic>? ?? {};
});

class EmergencyFundPage extends ConsumerWidget {
  const EmergencyFundPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(efProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Emergency Fund')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (d) => ListView(padding: const EdgeInsets.all(AppSpacing.md), children: [
          Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: Column(children: [
            Text('₹${d['current_savings']}', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold, color: theme.colorScheme.primary)),
            Text('of ₹${d['target_amount']} target', style: theme.textTheme.bodySmall),
            LinearProgressIndicator(value: ((d['progress_pct'] as num?)?.toDouble() ?? 0) / 100),
            Text('${d['months_covered']} of ${d['target_months']} months covered', style: theme.textTheme.bodySmall),
          ]))),
          Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: _row(theme, 'Monthly Expenses', '₹${d['monthly_expenses']}'))),
          Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: _row(theme, 'Status', d['status'] ?? ''))),
        ]),
      ),
    );
  }
  Widget _row(ThemeData t, String label, String value) => Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
    Text(label, style: t.textTheme.bodyMedium), Text(value, style: t.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
  ]);
}
