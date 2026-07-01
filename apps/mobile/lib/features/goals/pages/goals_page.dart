import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
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
            ? HorizonSearchBar(
                hintText: 'Search goals...',
                onSearch: (_) => setState(() {}),
                onClose: () => setState(() { _searching = false; _searchCtrl.clear(); }),
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
    if (state.loading) return const LoadingView();
    if (state.error != null) return ErrorView(
      message: 'Could not load goals',
      onRetry: () => ref.read(goalsListProvider.notifier).load(),
    );
    if (state.goals.isEmpty) return EmptyState(
      icon: Icons.flag_outlined,
      title: 'No goals yet',
      subtitle: 'Create your first financial goal',
    );

    final filtered = _searchCtrl.text.isEmpty ? state.goals : state.goals.where((g) =>
      g.name.toLowerCase().contains(_searchCtrl.text.toLowerCase())).toList();

    return RefreshIndicator(
      onRefresh: () => ref.read(goalsListProvider.notifier).load(),
      child: ListView.builder(
        padding: const EdgeInsets.all(AppTheme.spacingLg),
        itemCount: filtered.length,
        itemBuilder: (_, i) => Padding(
          padding: const EdgeInsets.only(bottom: AppTheme.spacingMd),
          child: GoalListCard(goal: filtered[i], onTap: () => context.push('/goals/${filtered[i].goalId}')),
        ),
      ),
    );
  }
}
