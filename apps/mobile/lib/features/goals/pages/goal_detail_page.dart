import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_animate/flutter_animate.dart';
import '../models/goal_models.dart';
import '../repository/goal_repository.dart';
import '../widgets/goal_widgets.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';

final goalDetailProvider = StateNotifierProvider.family<GoalDetailNotifier, GoalDetailState, String>((ref, goalId) {
  return GoalDetailNotifier(ref.read(goalRepositoryProvider), goalId);
});

class GoalDetailState {
  final bool loading; final GoalDashboardData? dashboard; final GoalProgressData? progress;
  final GoalProjectionData? projection; final GoalRecsData? recs; final GoalOptsData? opts;
  final GoalTimelineData? timeline; final GoalMilestoneData? milestones; final String? error;
  GoalDetailState({this.loading = false, this.dashboard, this.progress, this.projection, this.recs, this.opts, this.timeline, this.milestones, this.error});
}

class GoalDetailNotifier extends StateNotifier<GoalDetailState> {
  final GoalRepository _repo; final String _goalId;
  GoalDetailNotifier(this._repo, this._goalId) : super(GoalDetailState());

  Future<void> load({String? userId}) async {
    state = GoalDetailState(loading: true);
    try {
      final results = await Future.wait([
        _repo.getDashboard(userId: userId, goalId: _goalId),
        _repo.getProgress(userId: userId, goalId: _goalId),
        _repo.getProjection(userId: userId, goalId: _goalId),
        _repo.getRecommendations(userId: userId, goalId: _goalId),
        _repo.getOptimization(userId: userId, goalId: _goalId),
        _repo.getTimeline(userId: userId, goalId: _goalId),
        _repo.getMilestones(userId: userId, goalId: _goalId),
      ]);
      state = GoalDetailState(
        dashboard: (results[0] as GoalDashboardResponse).data,
        progress: (results[1] as GoalProgressResponse).data,
        projection: (results[2] as GoalProjectionResponse).data,
        recs: (results[3] as GoalRecsResponse).data,
        opts: (results[4] as GoalOptsResponse).data,
        timeline: (results[5] as GoalTimelineResponse).data,
        milestones: (results[6] as GoalMilestoneResponse).data,
      );
    } catch (e) {
      state = GoalDetailState(error: e.toString());
    }
  }
}

class GoalDetailPage extends ConsumerStatefulWidget {
  final String goalId;
  const GoalDetailPage({super.key, required this.goalId});
  @override
  ConsumerState<GoalDetailPage> createState() => _GoalDetailPageState();
}

class _GoalDetailPageState extends ConsumerState<GoalDetailPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(goalDetailProvider(widget.goalId).notifier).load());
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(goalDetailProvider(widget.goalId));
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        title: Text(state.dashboard?.goalName ?? 'Goal Details', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: isDark 
                ? [AppColors.darkBackground, AppColors.navy900]
                : [AppColors.lightBackground, AppColors.teal50.withOpacity(0.5)],
          ),
        ),
        child: _buildBody(theme, state),
      ),
    );
  }

  Widget _buildBody(ThemeData theme, GoalDetailState state) {
    if (state.loading) {
      return Center(
        child: CircularProgressIndicator(color: AppColors.teal500).animate().fade(),
      );
    }
    if (state.error != null) {
      return Center(
        child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
          Icon(Icons.error_outline_rounded, size: 64, color: AppColors.red500),
          const SizedBox(height: 16), const Text('Could not load goal details'),
          const SizedBox(height: 16),
          FilledButton.icon(
            style: FilledButton.styleFrom(backgroundColor: AppColors.teal500),
            onPressed: () => ref.read(goalDetailProvider(widget.goalId).notifier).load(), 
            icon: const Icon(Icons.refresh_rounded), 
            label: const Text('Try Again')
          ),
        ]),
      );
    }
    if (state.dashboard == null) return const Center(child: Text('No data'));

    final children = <Widget>[
      GoalHeaderCard(goal: state.dashboard!),
      const SizedBox(height: 16),
      Text('Dashboard', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
      const SizedBox(height: 8),
      ...state.dashboard!.cards.map((c) => Padding(
        padding: const EdgeInsets.only(bottom: AppSpacing.md),
        child: GoalDashboardCard(
          title: c.title, summary: c.summary, value: c.formattedValue,
          icon: _cardIcon(c.cardType), color: _cardColor(c.cardType),
        ),
      )),
      if (state.projection != null) ...[
        const SizedBox(height: 16),
        Text('Projection', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        GlassCard(child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _row('On Track', state.projection!.onTrack ? 'Yes' : 'No', context),
            _row('Projected Date', state.projection!.projectedDate ?? '--', context),
            _row('Monthly', '₹${state.projection!.monthlyContribution.toInt()}', context),
          ],
        )),
      ],
      if (state.recs != null && state.recs!.recs.isNotEmpty) ...[
        const SizedBox(height: 16),
        Text('Recommendations', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        ...state.recs!.recs.map((r) => Padding(
          padding: const EdgeInsets.only(bottom: AppSpacing.sm),
          child: GlassCard(child: ListTile(
            contentPadding: EdgeInsets.zero,
            leading: const Icon(Icons.lightbulb_rounded, color: AppColors.amber500),
            title: Text(r.title, style: const TextStyle(fontWeight: FontWeight.w600)), 
            subtitle: Text(r.summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            trailing: Text(r.impact, style: TextStyle(color: r.impact == 'High' ? AppColors.teal500 : AppColors.amber500, fontWeight: FontWeight.bold)),
          )),
        )),
      ],
      if (state.opts != null && state.opts!.opts.isNotEmpty) ...[
        const SizedBox(height: 16),
        Text('Optimization', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        ...state.opts!.opts.map((o) => Padding(
          padding: const EdgeInsets.only(bottom: AppSpacing.sm),
          child: GlassCard(child: ListTile(
            contentPadding: EdgeInsets.zero,
            leading: const Icon(Icons.auto_graph_rounded, color: Colors.cyan),
            title: Text(o.strategy, style: const TextStyle(fontWeight: FontWeight.w600)), 
            subtitle: Text(o.summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            trailing: Text(o.score.toStringAsFixed(1), style: const TextStyle(fontWeight: FontWeight.bold, color: Colors.cyan)),
          )),
        )),
      ],
      if (state.timeline != null && state.timeline!.events.isNotEmpty) ...[
        const SizedBox(height: 16),
        Text('Timeline', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        ...state.timeline!.events.take(5).map((e) => Padding(
          padding: const EdgeInsets.only(bottom: AppSpacing.sm),
          child: GlassCard(child: ListTile(
            contentPadding: EdgeInsets.zero,
            leading: Icon(Icons.circle, size: 12, color: theme.colorScheme.primary),
            title: Text(e.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
            subtitle: Text(e.timestamp, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          )),
        )),
      ],
      if (state.milestones != null && state.milestones!.milestones.isNotEmpty) ...[
        const SizedBox(height: 16),
        Text('Milestones', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        ...state.milestones!.milestones.map((m) => Padding(
          padding: const EdgeInsets.only(bottom: AppSpacing.sm),
          child: GlassCard(child: ListTile(
            contentPadding: EdgeInsets.zero,
            leading: Icon(m.reached ? Icons.check_circle_rounded : Icons.radio_button_unchecked_rounded, color: m.reached ? AppColors.teal500 : AppColors.slate500),
            title: Text(m.title, style: const TextStyle(fontWeight: FontWeight.w600)),
            trailing: Text('${m.progressPct}%', style: const TextStyle(fontWeight: FontWeight.bold)),
          )),
        )),
      ],
      const SizedBox(height: 100), // padding for navbar
    ];

    return RefreshIndicator(
      color: AppColors.teal500,
      onRefresh: () => ref.read(goalDetailProvider(widget.goalId).notifier).load(),
      child: ListView.builder(
        padding: const EdgeInsets.all(AppSpacing.md),
        itemCount: children.length,
        itemBuilder: (context, index) {
          return children[index]
            .animate()
            .fade(duration: 400.ms, delay: (20 * index).ms)
            .slideY(begin: 0.1, end: 0, duration: 400.ms, curve: Curves.easeOutQuad, delay: (20 * index).ms);
        },
      ),
    );
  }

  Widget _row(String label, String value, BuildContext context) => Padding(
    padding: const EdgeInsets.symmetric(vertical: 6),
    child: Row(children: [
      SizedBox(width: 140, child: Text(label, style: TextStyle(color: Theme.of(context).colorScheme.onSurfaceVariant, fontWeight: FontWeight.w500))), 
      Text(value, style: const TextStyle(fontWeight: FontWeight.bold))
    ]),
  );

  IconData _cardIcon(String type) {
    switch (type) {
      case 'Overview': return Icons.info_outline_rounded;
      case 'Progress': return Icons.pie_chart_rounded;
      case 'Funding': return Icons.account_balance_wallet_rounded;
      case 'Projection': return Icons.trending_up_rounded;
      case 'Recommendation': return Icons.lightbulb_rounded;
      case 'Optimization': return Icons.auto_graph_rounded;
      case 'Timeline': return Icons.history_rounded;
      case 'Milestone': return Icons.flag_rounded;
      default: return Icons.circle;
    }
  }

  Color _cardColor(String type) {
    switch (type) {
      case 'Overview': return Colors.blue;
      case 'Progress': return AppColors.teal500;
      case 'Funding': return AppColors.teal500;
      case 'Projection': return Colors.cyan;
      case 'Recommendation': return AppColors.amber500;
      case 'Optimization': return Colors.purple;
      case 'Timeline': return Colors.indigo;
      case 'Milestone': return AppColors.amber500;
      default: return AppColors.slate500;
    }
  }
}
