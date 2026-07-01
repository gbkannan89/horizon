import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/search_models.dart';
import '../state/search_state.dart';

class SearchPage extends ConsumerStatefulWidget {
  const SearchPage({super.key});
  @override
  ConsumerState<SearchPage> createState() => _SearchPageState();
}

class _SearchPageState extends ConsumerState<SearchPage> {
  final _searchCtrl = TextEditingController();
  Timer? _debounce;

  @override
  void dispose() {
    _searchCtrl.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onSearchChanged(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 300), () {
      ref.read(searchProvider.notifier).search(value.trim());
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(searchProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _searchCtrl,
          autofocus: true,
          decoration: InputDecoration(
            hintText: 'Search accounts, transactions, goals...',
            border: InputBorder.none,
            prefixIcon: const Icon(Icons.search),
            suffixIcon: _searchCtrl.text.isNotEmpty
                ? IconButton(
                    icon: const Icon(Icons.close),
                    onPressed: () {
                      _searchCtrl.clear();
                      ref.read(searchProvider.notifier).clearSearch();
                    },
                  )
                : null,
          ),
          onChanged: _onSearchChanged,
        ),
      ),
      body: _buildBody(context, theme, state),
    );
  }

  Widget _buildBody(BuildContext context, ThemeData theme, SearchState state) {
    switch (state.status) {
      case SearchStatus.idle:
        return _buildIdle(theme, state);
      case SearchStatus.searching:
        return const SharedLoadingView();
      case SearchStatus.loaded:
        return _buildResults(theme, state);
      case SearchStatus.error:
        return SharedErrorView(
          message: state.error,
          onRetry: () => ref.read(searchProvider.notifier).search(_searchCtrl.text.trim()),
        );
    }
  }

  Widget _buildIdle(ThemeData theme, SearchState state) {
    if (state.recentSearches.isEmpty) {
      return SharedEmptyView(
        icon: Icons.search,
        title: 'Search your finances',
        subtitle: 'Find accounts, transactions, goals, and more',
      );
    }
    return ListView(
      padding: const EdgeInsets.all(AppSpacing.md),
      children: [
        SharedSectionHeader(
          title: 'Recent Searches',
          trailing: TextButton(
            onPressed: () => ref.read(searchProvider.notifier).clearHistory(),
            child: const Text('Clear'),
          ),
        ),
        ...state.recentSearches.map((q) => ListTile(
          leading: const Icon(Icons.history),
          title: Text(q),
          trailing: const Icon(Icons.chevron_right, size: 18),
          onTap: () {
            _searchCtrl.text = q;
            _onSearchChanged(q);
          },
        )),
      ],
    );
  }

  Widget _buildResults(ThemeData theme, SearchState state) {
    if (state.results.isEmpty) {
      return SharedEmptyView(
        icon: Icons.search_off,
        title: 'No results found',
        subtitle: 'Try a different search term',
      );
    }

    final grouped = <String, List<SearchResultItem>>{};
    for (final item in state.results) {
      grouped.putIfAbsent(item.module, () => []).add(item);
    }

    return ListView(
      padding: const EdgeInsets.all(AppSpacing.md),
      children: grouped.entries.map((entry) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SharedSectionHeader(title: _moduleLabel(entry.key)),
            const SizedBox(height: AppSpacing.sm),
            ...entry.value.map((item) => Padding(
              padding: const EdgeInsets.only(bottom: AppSpacing.sm),
              child: SharedCard(
                onTap: () => context.push(item.route),
                child: Row(
                  children: [
                    Container(
                      width: 40, height: 40,
                      decoration: BoxDecoration(
                        color: item.color.withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(AppRadius.sm),
                      ),
                      child: Icon(item.icon, size: AppIconSize.md, color: item.color),
                    ),
                    const SizedBox(width: AppSpacing.sm),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(item.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500), maxLines: 1, overflow: TextOverflow.ellipsis),
                          Text(item.subtitle, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 1, overflow: TextOverflow.ellipsis),
                        ],
                      ),
                    ),
                    if (item.amount != null)
                      Text(item.amount!, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                  ],
                ),
              ),
            )),
            const SizedBox(height: AppSpacing.sm),
          ],
        );
      }).toList(),
    );
  }

  String _moduleLabel(String module) {
    switch (module) {
      case 'transaction': return 'Transactions';
      case 'timeline': return 'Timeline';
      case 'account': return 'Accounts';
      case 'goal': return 'Goals';
      case 'portfolio': return 'Portfolio';
      default: return module;
    }
  }
}
