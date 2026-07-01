import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/search_models.dart';
import '../state/search_state.dart';
import '../widgets/search_widgets.dart';

class SearchPage extends ConsumerStatefulWidget {
  final String initialQuery;

  const SearchPage({super.key, this.initialQuery = ''});

  @override
  ConsumerState<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends ConsumerState<SearchPage> {
  final _searchCtrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    if (widget.initialQuery.isNotEmpty) {
      _searchCtrl.text = widget.initialQuery;
      Future.microtask(() => ref.read(searchStateProvider.notifier).search(widget.initialQuery));
    }
  }

  @override
  void dispose() {
    _searchCtrl.dispose();
    super.dispose();
  }

  void _onSearch(String query) {
    ref.read(searchStateProvider.notifier).search(query);
  }

  void _showFilterSheet() {
    final current = ref.read(searchStateProvider).filters;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (_) => SearchFilterSheet(
        current: current,
        onApply: (filters) => ref.read(searchStateProvider.notifier).applyFilters(filters),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(searchStateProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _searchCtrl,
          autofocus: true,
          decoration: InputDecoration(
            hintText: 'Search accounts, goals, transactions...',
            border: InputBorder.none,
            prefixIcon: const Icon(Icons.search),
            suffixIcon: _searchCtrl.text.isNotEmpty
                ? IconButton(icon: const Icon(Icons.close), onPressed: () { _searchCtrl.clear(); _onSearch(''); })
                : null,
          ),
          onChanged: _onSearch,
        ),
        actions: [
          IconButton(icon: const Icon(Icons.filter_list), onPressed: _showFilterSheet),
        ],
      ),
      body: _buildBody(context, theme, state),
    );
  }

  Widget _buildBody(BuildContext context, ThemeData theme, SearchState state) {
    switch (state.status) {
      case SearchStatus.initial:
        return _buildInitial(context, theme, state);
      case SearchStatus.loading:
        return const LoadingView();
      case SearchStatus.loaded:
        return _buildResults(context, theme, state);
      case SearchStatus.empty:
        return EmptyState(
          icon: Icons.search_off,
          title: 'No results found',
          subtitle: "Try a different search term or check your filters",
        );
      case SearchStatus.error:
        return ErrorView(
          message: 'Search failed',
          detail: state.error,
          onRetry: () => _onSearch(state.query),
        );
    }
  }

  Widget _buildInitial(BuildContext context, ThemeData theme, SearchState state) {
    if (state.recentSearches.isEmpty) {
      return EmptyState(
        icon: Icons.search,
        title: 'Search Horizon',
        subtitle: 'Find accounts, goals, transactions, portfolio items, and timeline events',
      );
    }
    return ListView(
      padding: const EdgeInsets.all(16),
      children: [
        SectionHeader(
          title: 'Recent Searches',
          trailing: TextButton(
            onPressed: () => ref.read(searchStateProvider.notifier).clearHistory(),
            child: const Text('Clear'),
          ),
        ),
        ...state.recentSearches.map((r) => RecentSearchTile(
          recent: r,
          onTap: () {
            _searchCtrl.text = r.query;
            _onSearch(r.query);
          },
          onDelete: () {},
        )),
      ],
    );
  }

  Widget _buildResults(BuildContext context, ThemeData theme, SearchState state) {
    final grouped = <SearchModule, List<SearchResult>>{};
    for (final r in state.results) {
      grouped.putIfAbsent(r.module, () => []).add(r);
    }
    if (state.filters.hasActiveFilters && state.results.isNotEmpty) {
      return Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          _buildActiveFilterBar(context, theme, state),
          Expanded(child: _buildGroupedResults(context, theme, grouped)),
        ],
      );
    }
    return _buildGroupedResults(context, theme, grouped);
  }

  Widget _buildActiveFilterBar(BuildContext context, ThemeData theme, SearchState state) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        children: [
          Icon(Icons.filter_alt, size: 16, color: theme.colorScheme.primary),
          const SizedBox(width: 8),
          Text('${state.results.length} results', style: theme.textTheme.bodySmall),
          const Spacer(),
          TextButton(
            onPressed: () => ref.read(searchStateProvider.notifier).clearFilters(),
            child: const Text('Clear Filters'),
          ),
        ],
      ),
    );
  }

  Widget _buildGroupedResults(BuildContext context, ThemeData theme, Map<SearchModule, List<SearchResult>> grouped) {
    return ListView(
      padding: const EdgeInsets.all(16),
      children: grouped.entries.map((entry) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SectionHeader(title: '${entry.key.label} (${entry.value.length})'),
            ...entry.value.map((r) => Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: SearchResultCard(result: r),
            )),
          ],
        );
      }).toList(),
    );
  }
}
