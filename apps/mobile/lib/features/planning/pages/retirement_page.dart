import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';

final retirementProvider = FutureProvider.autoDispose<Map<String, dynamic>>((ref) async {
  final r = await ref.read(apiClientProvider).get('/planning/retirement');
  return (r.data as Map<String, dynamic>)['data'] as Map<String, dynamic>? ?? {};
});

class RetirementPage extends ConsumerWidget {
  const RetirementPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(retirementProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Retirement Planning')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (d) => ListView(padding: const EdgeInsets.all(AppSpacing.md), children: [
          Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: Column(children: [
            Text('₹${d['current_corpus']}', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold, color: theme.colorScheme.primary)),
            Text('Current Corpus', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: AppSpacing.sm),
            LinearProgressIndicator(value: ((d['progress_pct'] as num?)?.toDouble() ?? 0) / 100),
            Text('${d['progress_pct']}% of ₹${d['target_corpus']} target', style: theme.textTheme.bodySmall),
          ]))),
          Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: _row(theme, 'Monthly Contribution', '₹${d['monthly_contribution']}'))),
          Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: _row(theme, 'Recommended', '₹${d['recommended_contribution']}'))),
          Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: _row(theme, 'Projected Retirement Age', '${d['projected_retirement_age']}'))),
        ]),
      ),
    );
  }
  Widget _row(ThemeData t, String label, String value) => Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
    Text(label, style: t.textTheme.bodyMedium), Text(value, style: t.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
  ]);
}
