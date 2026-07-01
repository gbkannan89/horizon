import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/goal_models.dart';
import '../repository/goal_repository.dart';
import '../widgets/goal_widgets.dart';

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

    return Scaffold(
      appBar: AppBar(title: Text(state.dashboard?.goalName ?? 'Goal Details')),
      body: _buildBody(theme, state),
    );
  }

  Widget _buildBody(ThemeData theme, GoalDetailState state) {
    if (state.loading) return const LoadingView();
    if (state.error != null) return ErrorView(
      message: 'Could not load goal details',
      onRetry: () => ref.read(goalDetailProvider(widget.goalId).notifier).load(),
    );
    if (state.dashboard == null) return const EmptyState(icon: Icons.search_off, title: 'No data');

    return RefreshIndicator(
      onRefresh: () => ref.read(goalDetailProvider(widget.goalId).notifier).load(),
      child: ListView(
        padding: const EdgeInsets.all(AppTheme.spacingLg),
        children: [
          GoalHeaderCard(goal: state.dashboard!),
          const SizedBox(height: AppTheme.spacingLg),
          SectionHeader(title: 'Dashboard'),
          const SizedBox(height: AppTheme.spacingSm),
          ...state.dashboard!.cards.map((c) => Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: IconCard(
              title: c.title, subtitle: c.summary, trailing: c.formattedValue,
              icon: _cardIcon(c.cardType), color: _cardColor(c.cardType),
            ),
          )),
          if (state.projection != null) ...[
            const SizedBox(height: AppTheme.spacingLg),
            SectionHeader(title: 'Projection'),
            const SizedBox(height: AppTheme.spacingSm),
            AppCard(child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                InfoRow(label: 'On Track', value: state.projection!.onTrack ? 'Yes' : 'No'),
                InfoRow(label: 'Projected Date', value: state.projection!.projectedDate ?? '--'),
                InfoRow(label: 'Monthly', value: '₹${state.projection!.monthlyContribution.toInt()}'),
              ],
            )),
          ],
          if (state.recs != null && state.recs!.recs.isNotEmpty) ...[
            const SizedBox(height: AppTheme.spacingLg),
            SectionHeader(title: 'Recommendations'),
            const SizedBox(height: AppTheme.spacingSm),
            ...state.recs!.recs.map((r) => Card(child: ListTile(
              leading: Icon(Icons.lightbulb, color: Colors.amber),
              title: Text(r.title), subtitle: Text(r.summary),
              trailing: Text(r.impact, style: TextStyle(color: r.impact == 'High' ? Colors.green : Colors.orange, fontWeight: FontWeight.w600)),
            ))),
          ],
          if (state.opts != null && state.opts!.opts.isNotEmpty) ...[
            const SizedBox(height: AppTheme.spacingLg),
            SectionHeader(title: 'Optimization'),
            const SizedBox(height: AppTheme.spacingSm),
            ...state.opts!.opts.map((o) => Card(child: ListTile(
              leading: Icon(Icons.auto_graph, color: Colors.cyan),
              title: Text(o.strategy), subtitle: Text(o.summary),
              trailing: Text('${o.score.toStringAsFixed(1)}', style: TextStyle(fontWeight: FontWeight.bold, color: Colors.cyan)),
            ))),
          ],
          if (state.timeline != null && state.timeline!.events.isNotEmpty) ...[
            const SizedBox(height: AppTheme.spacingLg),
            SectionHeader(title: 'Timeline'),
            const SizedBox(height: AppTheme.spacingSm),
            ...state.timeline!.events.take(5).map((e) => Card(child: ListTile(
              leading: Icon(Icons.circle, size: 12, color: theme.colorScheme.primary),
              title: Text(e.title, style: theme.textTheme.bodyMedium),
              subtitle: Text(e.timestamp, style: theme.textTheme.bodySmall),
            ))),
          ],
          if (state.milestones != null && state.milestones!.milestones.isNotEmpty) ...[
            const SizedBox(height: AppTheme.spacingLg),
            SectionHeader(title: 'Milestones'),
            const SizedBox(height: AppTheme.spacingSm),
            ...state.milestones!.milestones.map((m) => Card(child: ListTile(
              leading: Icon(m.reached ? Icons.check_circle : Icons.radio_button_unchecked, color: m.reached ? Colors.green : Colors.grey),
              title: Text(m.title),
              trailing: Text('${m.progressPct}%'),
            ))),
          ],
        ],
      ),
    );
  }

  IconData _cardIcon(String type) {
    switch (type) {
      case 'Overview': return Icons.info_outline;
      case 'Progress': return Icons.pie_chart;
      case 'Funding': return Icons.account_balance_wallet;
      case 'Projection': return Icons.trending_up;
      case 'Recommendation': return Icons.lightbulb;
      case 'Optimization': return Icons.auto_graph;
      case 'Timeline': return Icons.history;
      case 'Milestone': return Icons.flag;
      default: return Icons.circle;
    }
  }

  Color _cardColor(String type) {
    switch (type) {
      case 'Overview': return Colors.blue;
      case 'Progress': return Colors.green;
      case 'Funding': return Colors.teal;
      case 'Projection': return Colors.cyan;
      case 'Recommendation': return Colors.amber;
      case 'Optimization': return Colors.purple;
      case 'Timeline': return Colors.indigo;
      case 'Milestone': return Colors.orange;
      default: return Colors.grey;
    }
  }
}
