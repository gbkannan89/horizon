import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../models/goal_models.dart';
import '../repository/goal_repository.dart';
import '../widgets/goal_widgets.dart';

final goalsListProvider = StateNotifierProvider<GoalsListNotifier, GoalsListState>((ref) {
  return GoalsListNotifier(ref.read(goalRepositoryProvider));
});

class GoalsListState {
  final bool loading; final List<GoalSummary> goals; final String? error;
  GoalsListState({this.loading = false, this.goals = const [], this.error});
}

class GoalsListNotifier extends StateNotifier<GoalsListState> {
  final GoalRepository _repo;
  GoalsListNotifier(this._repo) : super(GoalsListState());

  Future<void> load({String? userId}) async {
    state = GoalsListState(loading: true);
    try {
      final resp = await _repo.getGoals(userId: userId);
      final list = resp.data;
      state = GoalsListState(goals: list?.goals ?? []);
    } catch (e) {
      state = GoalsListState(error: e.toString());
    }
  }
}

class GoalsPage extends ConsumerStatefulWidget {
  const GoalsPage({super.key});
  @override
  ConsumerState<GoalsPage> createState() => _GoalsPageState();
}

class _GoalsPageState extends ConsumerState<GoalsPage> {
  final _searchCtrl = TextEditingController();
  bool _searching = false;

  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(goalsListProvider.notifier).load());
  }

  @override
  void dispose() { _searchCtrl.dispose(); super.dispose(); }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(goalsListProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: _searching
            ? TextField(
                controller: _searchCtrl,
                autofocus: true,
                decoration: InputDecoration(hintText: 'Search goals...', border: InputBorder.none,
                  suffixIcon: IconButton(icon: const Icon(Icons.close), onPressed: () {
                    setState(() { _searching = false; _searchCtrl.clear(); });
                  }),
                ),
              )
            : const Text('Goals'),
        actions: [
          IconButton(icon: Icon(_searching ? Icons.search_off : Icons.search), onPressed: () => setState(() => _searching = !_searching)),
        ],
      ),
      body: _buildBody(theme, state),
    );
  }

  Widget _buildBody(ThemeData theme, GoalsListState state) {
    if (state.loading) return const Center(child: CircularProgressIndicator());
    if (state.error != null) return Center(
      child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
        Icon(Icons.cloud_off, size: 64, color: theme.colorScheme.error),
        const SizedBox(height: 16), Text('Could not load goals', style: theme.textTheme.titleMedium),
        const SizedBox(height: 16),
        FilledButton.icon(onPressed: () => ref.read(goalsListProvider.notifier).load(), icon: const Icon(Icons.refresh), label: const Text('Try Again')),
      ]),
    );
    if (state.goals.isEmpty) return Center(
      child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
        Icon(Icons.flag_outlined, size: 64, color: theme.colorScheme.onSurfaceVariant.withValues(alpha: 0.4)),
        const SizedBox(height: 16), Text('No goals yet', style: theme.textTheme.titleMedium),
        const SizedBox(height: 8), Text('Create your first financial goal', style: theme.textTheme.bodySmall),
      ]),
    );

    final filtered = _searchCtrl.text.isEmpty ? state.goals : state.goals.where((g) =>
      g.name.toLowerCase().contains(_searchCtrl.text.toLowerCase())).toList();

    return RefreshIndicator(
      onRefresh: () => ref.read(goalsListProvider.notifier).load(),
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: filtered.length,
        itemBuilder: (_, i) => Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: GoalListCard(goal: filtered[i], onTap: () => context.push('/goals/${filtered[i].goalId}')),
        ),
      ),
    );
  }
}
