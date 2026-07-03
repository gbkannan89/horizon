import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../state/timeline_state.dart';
import '../widgets/timeline_widgets.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';

class TimelinePage extends ConsumerStatefulWidget {
  const TimelinePage({super.key});
  @override
  ConsumerState<TimelinePage> createState() => _TimelinePageState();
}

class _TimelinePageState extends ConsumerState<TimelinePage> {
  final _searchCtrl = TextEditingController();
  final _scrollCtrl = ScrollController();
  bool _searching = false;

  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(timelineStateProvider.notifier).load());
    _scrollCtrl.addListener(_onScroll);
  }

  @override
  void dispose() {
    _searchCtrl.dispose();
    _scrollCtrl.dispose();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollCtrl.position.pixels >= _scrollCtrl.position.maxScrollExtent - 200) {
      ref.read(timelineStateProvider.notifier).loadMore();
    }
  }

  void _showFilterSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (_) => FilterSheet(
        current: ref.read(timelineStateProvider).activeFilters,
        onApply: (filters) => ref.read(timelineStateProvider.notifier).applyFilters(filters: filters),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(timelineStateProvider);
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        title: _searching
            ? TextField(
                controller: _searchCtrl,
                autofocus: true,
                style: theme.textTheme.titleMedium,
                decoration: InputDecoration(
                  hintText: 'Search timeline...',
                  hintStyle: theme.textTheme.titleMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                  border: InputBorder.none,
                  suffixIcon: IconButton(
                    icon: const Icon(Icons.close_rounded),
                    onPressed: () {
                      setState(() { _searching = false; _searchCtrl.clear(); });
                      ref.read(timelineStateProvider.notifier).refresh();
                    },
                  ),
                ),
                onSubmitted: (q) {
                  if (q.isNotEmpty) ref.read(timelineStateProvider.notifier).search(query: q);
                },
              )
            : Text('Timeline', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        actions: [
          IconButton(
            icon: Icon(_searching ? Icons.search_rounded : Icons.search_rounded),
            onPressed: () => setState(() => _searching = !_searching),
          ),
          IconButton(icon: const Icon(Icons.filter_list_rounded), onPressed: _showFilterSheet),
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
        child: _buildBody(context, theme, state),
      ),
    );
  }

  Widget _buildBody(BuildContext context, ThemeData theme, TimelineState state) {
    switch (state.status) {
      case TimelineStatus.initial:
      case TimelineStatus.loading:
        return _buildSkeleton();
      case TimelineStatus.loadingMore:
        return _buildListWithLoading(context, theme, state);
      case TimelineStatus.loaded:
        return _buildList(context, theme, state);
      case TimelineStatus.empty:
        return _buildEmpty(theme, state);
      case TimelineStatus.error:
        return _buildError(theme, state);
      case TimelineStatus.offline:
        return _buildOffline(theme);
    }
  }

  Widget _buildSkeleton() {
    return const SharedSkeletonList(itemCount: 8, itemHeight: 88);
  }

  Widget _buildList(BuildContext context, ThemeData theme, TimelineState state) {
    if (state.items.isEmpty) return _buildEmpty(theme, state);
    return RefreshIndicator(
      color: AppColors.teal500,
      onRefresh: () => ref.read(timelineStateProvider.notifier).refresh(),
      child: ListView.builder(
        controller: _scrollCtrl,
        padding: const EdgeInsets.all(AppSpacing.md),
        itemCount: state.items.length + 1, // +1 for navbar padding
        itemBuilder: (context, index) {
          if (index == state.items.length) return const SizedBox(height: 100);
          final item = state.items[index];
          return Padding(
            padding: const EdgeInsets.only(bottom: AppSpacing.sm),
            child: TimelineCard(
              item: item,
              onTap: () => context.push('/timeline/${item.timelineId}', extra: item),
            ).animate().fade(duration: 400.ms, delay: (20 * index).ms).slideY(begin: 0.1, end: 0, duration: 400.ms, curve: Curves.easeOutQuad, delay: (20 * index).ms),
          );
        },
      ),
    );
  }

  Widget _buildListWithLoading(BuildContext context, ThemeData theme, TimelineState state) {
    return Stack(
      children: [
        _buildList(context, theme, state),
        if (state.status == TimelineStatus.loadingMore)
          Positioned(bottom: 24, left: 0, right: 0, child: Center(child: Padding(padding: const EdgeInsets.all(16), child: CircularProgressIndicator(strokeWidth: 2, color: AppColors.teal500)))),
      ],
    );
  }

  Widget _buildEmpty(ThemeData theme, TimelineState state) {
    return SharedEmptyView(
      icon: Icons.history_rounded,
      title: state.activeFilters.q != null ? 'No results found' : 'No timeline events',
      subtitle: state.activeFilters.q != null ? 'Try a different search term' : 'Events will appear here as they happen',
    );
  }

  Widget _buildError(ThemeData theme, TimelineState state) {
    return SharedErrorView(
      message: state.error,
      onRetry: () => ref.read(timelineStateProvider.notifier).refresh(),
    );
  }

  Widget _buildOffline(ThemeData theme) {
    return const SharedOfflineView();
  }
}
