import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/planning_models.dart';
import '../repository/planning_repository.dart';

final planningProvider = FutureProvider.autoDispose<PlanningDashboard?>((ref) async {
  try {
    final repo = ref.read(planningRepositoryProvider);
    final r = await repo.getDashboard();
    return r.data;
  } catch (_) { return null; }
});

class PlanningPage extends ConsumerWidget {
  const PlanningPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(planningProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Planning')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (dash) => ListView(
          padding: const EdgeInsets.all(AppSpacing.md),
          children: [
            SharedSectionHeader(title: 'Financial Planning'),
            const SizedBox(height: AppSpacing.sm),
            _section(theme, 'Projections', Icons.query_stats, Colors.cyan, 'View your financial projections', () => context.push('/planning/projections')),
            _section(theme, 'Scenarios', Icons.compare_arrows, Colors.indigo, 'Create and compare scenarios', () => context.push('/planning/scenarios')),
            const SizedBox(height: AppSpacing.md),
            SharedSectionHeader(title: 'Coming Soon'),
            const SizedBox(height: AppSpacing.sm),
            _comingSoon(theme, 'Budget Planning', Icons.account_balance_wallet, 'Set and track budgets'),
            _comingSoon(theme, 'Retirement Planning', Icons.beach_access, 'Plan for retirement'),
            _comingSoon(theme, 'Emergency Fund', Icons.shield_outlined, 'Build your safety net'),
            _comingSoon(theme, 'Debt Payoff', Icons.credit_score, 'Create a debt payoff strategy'),
            _comingSoon(theme, 'Investment Planning', Icons.trending_up, 'Optimize your investments'),
          ],
        ),
      ),
    );
  }

  Widget _section(ThemeData t, String title, IconData icon, Color color, String subtitle, VoidCallback onTap) {
    return Card(
      clipBehavior: Clip.antiAlias,
      margin: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.md),
          child: Row(children: [
            Container(
              width: 48, height: 48,
              decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(AppRadius.md)),
              child: Icon(icon, color: color),
            ),
            const SizedBox(width: AppSpacing.md),
            Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(title, style: t.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w600)),
              Text(subtitle, style: t.textTheme.bodySmall?.copyWith(color: t.colorScheme.onSurfaceVariant)),
            ])),
            const Icon(Icons.chevron_right, size: 18),
          ]),
        ),
      ),
    );
  }

  Widget _comingSoon(ThemeData t, String title, IconData icon, String subtitle) {
    return Card(
      margin: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: Opacity(
        opacity: 0.5,
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.md),
          child: Row(children: [
            Container(
              width: 48, height: 48,
              decoration: BoxDecoration(color: t.colorScheme.surfaceContainerHighest, borderRadius: BorderRadius.circular(AppRadius.md)),
              child: Icon(icon, color: t.colorScheme.onSurfaceVariant),
            ),
            const SizedBox(width: AppSpacing.md),
            Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Text(title, style: t.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w600)),
              Text(subtitle, style: t.textTheme.bodySmall?.copyWith(color: t.colorScheme.onSurfaceVariant)),
            ])),
            Chip(label: Text('Soon', style: t.textTheme.labelSmall), visualDensity: VisualDensity.compact),
          ]),
        ),
      ),
    );
  }
}
