import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
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
    if (state.loading) return const Center(child: CircularProgressIndicator());
    if (state.error != null) {
      return Center(
      child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
        Icon(Icons.error_outline, size: 64, color: theme.colorScheme.error),
        const SizedBox(height: 16), const Text('Could not load goal details'),
        const SizedBox(height: 16),
        FilledButton.icon(onPressed: () => ref.read(goalDetailProvider(widget.goalId).notifier).load(), icon: const Icon(Icons.refresh), label: const Text('Try Again')),
      ]),
    );
    }
    if (state.dashboard == null) return const Center(child: Text('No data'));

    return RefreshIndicator(
      onRefresh: () => ref.read(goalDetailProvider(widget.goalId).notifier).load(),
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          GoalHeaderCard(goal: state.dashboard!),
          const SizedBox(height: 16),
          Text('Dashboard', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
          const SizedBox(height: 8),
          ...state.dashboard!.cards.map((c) => Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: GoalDashboardCard(
              title: c.title, summary: c.summary, value: c.formattedValue,
              icon: _cardIcon(c.cardType), color: _cardColor(c.cardType),
            ),
          )),
          if (state.projection != null) ...[
            const SizedBox(height: 16),
            Text('Projection', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            Card(child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _row('On Track', state.projection!.onTrack ? 'Yes' : 'No'),
                  _row('Projected Date', state.projection!.projectedDate ?? '--'),
                  _row('Monthly', '₹${state.projection!.monthlyContribution.toInt()}'),
                ],
              ),
            )),
          ],
          if (state.recs != null && state.recs!.recs.isNotEmpty) ...[
            const SizedBox(height: 16),
            Text('Recommendations', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            ...state.recs!.recs.map((r) => Card(child: ListTile(
              leading: const Icon(Icons.lightbulb, color: Colors.amber),
              title: Text(r.title), subtitle: Text(r.summary),
              trailing: Text(r.impact, style: TextStyle(color: r.impact == 'High' ? Colors.green : Colors.orange, fontWeight: FontWeight.w600)),
            ))),
          ],
          if (state.opts != null && state.opts!.opts.isNotEmpty) ...[
            const SizedBox(height: 16),
            Text('Optimization', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            ...state.opts!.opts.map((o) => Card(child: ListTile(
              leading: const Icon(Icons.auto_graph, color: Colors.cyan),
              title: Text(o.strategy), subtitle: Text(o.summary),
              trailing: Text(o.score.toStringAsFixed(1), style: const TextStyle(fontWeight: FontWeight.bold, color: Colors.cyan)),
            ))),
          ],
          if (state.timeline != null && state.timeline!.events.isNotEmpty) ...[
            const SizedBox(height: 16),
            Text('Timeline', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            ...state.timeline!.events.take(5).map((e) => Card(child: ListTile(
              leading: Icon(Icons.circle, size: 12, color: theme.colorScheme.primary),
              title: Text(e.title, style: theme.textTheme.bodyMedium),
              subtitle: Text(e.timestamp, style: theme.textTheme.bodySmall),
            ))),
          ],
          if (state.milestones != null && state.milestones!.milestones.isNotEmpty) ...[
            const SizedBox(height: 16),
            Text('Milestones', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
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

  Widget _row(String label, String value) => Padding(
    padding: const EdgeInsets.symmetric(vertical: 4),
    child: Row(children: [SizedBox(width: 120, child: Text(label, style: TextStyle(color: Theme.of(context).colorScheme.onSurfaceVariant))), Text(value)]),
  );

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
