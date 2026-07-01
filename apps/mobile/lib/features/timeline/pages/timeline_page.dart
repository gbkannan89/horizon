import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../state/timeline_state.dart';
import '../widgets/timeline_widgets.dart';

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

    return Scaffold(
      appBar: AppBar(
        title: _searching
            ? HorizonSearchBar(
                hintText: 'Search timeline...',
                onSearch: (q) { if (q.isNotEmpty) ref.read(timelineStateProvider.notifier).search(query: q); },
                onClose: () { setState(() { _searching = false; _searchCtrl.clear(); }); ref.read(timelineStateProvider.notifier).refresh(); },
              )
            : const Text('Timeline'),
        actions: [
          IconButton(
            icon: const Icon(Icons.search),
            onPressed: () => setState(() => _searching = !_searching),
          ),
          IconButton(icon: const Icon(Icons.filter_list), onPressed: _showFilterSheet),
        ],
      ),
      body: _buildBody(context, theme, state),
    );
  }

  Widget _buildBody(BuildContext context, ThemeData theme, TimelineState state) {
    switch (state.status) {
      case TimelineStatus.initial:
      case TimelineStatus.loading:
        return const SkeletonList(itemCount: 8, itemHeight: 72);
      case TimelineStatus.loadingMore:
        return _buildListWithLoading(context, theme, state);
      case TimelineStatus.loaded:
        return _buildList(context, theme, state);
      case TimelineStatus.empty:
        return EmptyState(
          icon: Icons.history,
          title: state.activeFilters.q != null ? 'No results found' : 'No timeline events',
          subtitle: state.activeFilters.q != null ? 'Try a different search term' : 'Events will appear here as they happen',
        );
      case TimelineStatus.error:
        return ErrorView(
          message: 'Could not load timeline',
          onRetry: () => ref.read(timelineStateProvider.notifier).refresh(),
        );
      case TimelineStatus.offline:
        return const OfflineView();
    }
  }

  Widget _buildList(BuildContext context, ThemeData theme, TimelineState state) {
    if (state.items.isEmpty) return EmptyState(
      icon: Icons.history,
      title: state.activeFilters.q != null ? 'No results found' : 'No timeline events',
      subtitle: state.activeFilters.q != null ? 'Try a different search term' : 'Events will appear here as they happen',
    );
    return RefreshIndicator(
      onRefresh: () => ref.read(timelineStateProvider.notifier).refresh(),
      child: ListView.builder(
        controller: _scrollCtrl,
        padding: const EdgeInsets.all(AppTheme.spacingLg),
        itemCount: state.items.length,
        itemBuilder: (context, index) {
          final item = state.items[index];
          return Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: TimelineCard(
              item: item,
              onTap: () => context.push('/timeline/${item.timelineId}', extra: item),
            ),
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
          Positioned(bottom: 0, left: 0, right: 0, child: Center(child: Padding(padding: const EdgeInsets.all(AppTheme.spacingLg), child: CircularProgressIndicator(strokeWidth: 2, color: theme.colorScheme.primary)))),
      ],
    );
  }
}
