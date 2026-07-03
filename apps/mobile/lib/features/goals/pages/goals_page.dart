import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/goal_models.dart';
import '../repository/goal_repository.dart';
import '../widgets/goal_widgets.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';

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
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
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
            : Text('Goals', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
        actions: [
          IconButton(icon: Icon(_searching ? Icons.search_off : Icons.search), onPressed: () => setState(() => _searching = !_searching)),
        ],
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

  Widget _buildBody(ThemeData theme, GoalsListState state) {
    if (state.loading) {
      return Center(
        child: CircularProgressIndicator(color: AppColors.teal500).animate().fade(),
      );
    }
    if (state.error != null) {
      return SharedErrorView(
        message: state.error,
        onRetry: () => ref.read(goalsListProvider.notifier).load(),
      );
    }
    if (state.goals.isEmpty) {
      return const SharedEmptyView(
        icon: Icons.flag_outlined,
        title: 'No goals yet',
        subtitle: 'Create your first financial goal',
      );
    }

    final filtered = _searchCtrl.text.isEmpty ? state.goals : state.goals.where((g) =>
      g.name.toLowerCase().contains(_searchCtrl.text.toLowerCase())).toList();

    return RefreshIndicator(
      color: AppColors.teal500,
      onRefresh: () => ref.read(goalsListProvider.notifier).load(),
      child: ListView.builder(
        padding: const EdgeInsets.all(AppSpacing.md),
        itemCount: filtered.length,
        itemBuilder: (_, i) => Padding(
          padding: const EdgeInsets.only(bottom: AppSpacing.md),
          child: GoalListCard(
            goal: filtered[i], 
            onTap: () => context.push('/goals/${filtered[i].goalId}')
          ).animate()
           .fade(duration: 400.ms, delay: (50 * i).ms)
           .slideY(begin: 0.1, end: 0, duration: 400.ms, curve: Curves.easeOutQuad, delay: (50 * i).ms),
        ),
      ),
    );
  }
}
