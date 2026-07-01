import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';

final invProvider = FutureProvider.autoDispose<Map<String, dynamic>>((ref) async {
  final r = await ref.read(apiClientProvider).get('/planning/investment');
  return (r.data as Map<String, dynamic>)['data'] as Map<String, dynamic>? ?? {};
});

class InvestmentPage extends ConsumerWidget {
  const InvestmentPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(invProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Investment Planning')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (d) {
          final allocs = d['allocations'] as List? ?? [];
          return ListView(padding: const EdgeInsets.all(AppSpacing.md), children: [
            Card(child: Padding(padding: const EdgeInsets.all(AppSpacing.md), child: Column(children: [
              Text('₹${d['portfolio_value']}', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
              Text('Return: ${d['total_return']}% · ${d['risk_level']} Risk', style: theme.textTheme.bodySmall),
            ]))),
            ...allocs.map((a) => Card(margin: const EdgeInsets.only(bottom: AppSpacing.sm), child: ListTile(
              title: Text(a['type'] ?? ''),
              trailing: Text('${a['current']}%'),
              subtitle: LinearProgressIndicator(value: ((a['current'] as num?)?.toDouble() ?? 0) / 100),
            ))),
          ]);
        },
      ),
    );
  }
}
