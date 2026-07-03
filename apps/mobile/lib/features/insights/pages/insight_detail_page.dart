import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/insight_models.dart';
import '../repository/insight_repository.dart';
import '../widgets/insight_widgets.dart';

final insightDetailProvider = FutureProvider.autoDispose.family<InsightItem?, String>((ref, id) async {
  final repo = ref.read(insightRepositoryProvider);
  final resp = await repo.getInsights();
  return resp.data?.insights.where((i) => i.id == id).firstOrNull;
});

class InsightDetailPage extends ConsumerWidget {
  final String insightId;
  const InsightDetailPage({super.key, required this.insightId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(insightDetailProvider(insightId));
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Insight')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (item) {
          if (item == null) return const SharedEmptyView(icon: Icons.search_off, title: 'Insight not found');
          return ListView(
            padding: const EdgeInsets.all(AppSpacing.md),
            children: [
              InsightCard(insight: item),
              const SizedBox(height: AppSpacing.md),
              if (item.description != null) Card(
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                    Text('Details', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                    const SizedBox(height: AppSpacing.sm),
                    Text(item.description!, style: theme.textTheme.bodyMedium),
                  ]),
                ),
              ),
              if (item.impact != null || item.confidence != null) Card(
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: Column(children: [
                    if (item.impact != null) _row(theme, 'Impact', '${item.impact!.toStringAsFixed(0)}%'),
                    if (item.confidence != null) _row(theme, 'Confidence', '${item.confidence!.toStringAsFixed(0)}%'),
                  ]),
                ),
              ),
            ],
          );
        },
      ),
    );
  }

  Widget _row(ThemeData theme, String label, String value) => Padding(
    padding: const EdgeInsets.symmetric(vertical: AppSpacing.xs),
    child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
      Text(label, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
      Text(value, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500)),
    ]),
  );
}
