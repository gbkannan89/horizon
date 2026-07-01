import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
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
            ? TextField(
                controller: _searchCtrl,
                autofocus: true,
                decoration: InputDecoration(
                  hintText: 'Search timeline...',
                  border: InputBorder.none,
                  suffixIcon: IconButton(
                    icon: const Icon(Icons.close),
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
            : const Text('Timeline'),
        actions: [
          IconButton(
            icon: Icon(_searching ? Icons.search : Icons.search),
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
    return const SharedSkeletonList(itemCount: 8, itemHeight: 72);
  }

  Widget _buildList(BuildContext context, ThemeData theme, TimelineState state) {
    if (state.items.isEmpty) return _buildEmpty(theme, state);
    return RefreshIndicator(
      onRefresh: () => ref.read(timelineStateProvider.notifier).refresh(),
      child: ListView.builder(
        controller: _scrollCtrl,
        padding: const EdgeInsets.all(16),
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
          Positioned(bottom: 0, left: 0, right: 0, child: Center(child: Padding(padding: const EdgeInsets.all(16), child: CircularProgressIndicator(strokeWidth: 2, color: theme.colorScheme.primary)))),
      ],
    );
  }

  Widget _buildEmpty(ThemeData theme, TimelineState state) {
    return SharedEmptyView(
      icon: Icons.history,
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
