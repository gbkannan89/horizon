import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../state/transaction_state.dart';
import '../widgets/transaction_widgets.dart';
import '../widgets/transaction_filters.dart';

class TransactionsPage extends ConsumerStatefulWidget {
  const TransactionsPage({super.key});
  @override
  ConsumerState<TransactionsPage> createState() => _TransactionsPageState();
}

class _TransactionsPageState extends ConsumerState<TransactionsPage> {
  final _searchCtrl = TextEditingController();
  final _scrollCtrl = ScrollController();
  Timer? _debounceTimer;
  bool _searching = false;

  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(transactionListProvider.notifier).load());
    _scrollCtrl.addListener(_onScroll);
  }

  @override
  void dispose() {
    _searchCtrl.dispose();
    _scrollCtrl.dispose();
    _debounceTimer?.cancel();
    super.dispose();
  }

  void _onScroll() {
    if (_scrollCtrl.position.pixels >= _scrollCtrl.position.maxScrollExtent - 200) {
      ref.read(transactionListProvider.notifier).loadMore();
    }
  }

  void _onSearchChanged(String query) {
    _debounceTimer?.cancel();
    _debounceTimer = Timer(const Duration(milliseconds: 300), () {
      ref.read(transactionListProvider.notifier).search(query: query);
    });
  }

  void _showFilterSheet() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (_) => TransactionFilterSheet(
        current: ref.read(transactionListProvider).activeFilters,
        onApply: (filters) => ref.read(transactionListProvider.notifier).applyFilters(filters: filters),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(transactionListProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: _searching
            ? SharedSearchBar(
                controller: _searchCtrl,
                hintText: 'Search transactions...',
                onChanged: _onSearchChanged,
                onCancel: () {
                  setState(() { _searching = false; _searchCtrl.clear(); });
                  ref.read(transactionListProvider.notifier).refresh();
                },
              )
            : const Text('Transactions'),
        actions: [
          if (state.activeFilters.hasActiveFilters)
            Badge(
              label: Text('${state.activeFilters.activeFilterCount}'),
              child: IconButton(icon: const Icon(Icons.filter_list), onPressed: _showFilterSheet),
            )
          else
            IconButton(icon: const Icon(Icons.filter_list), onPressed: _showFilterSheet),
          IconButton(
            icon: Icon(_searching ? Icons.search_off : Icons.search),
            onPressed: () => setState(() => _searching = !_searching),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () async {
          final result = await context.push<bool>('/transactions/add');
          if (result == true && mounted) {
            ref.read(transactionListProvider.notifier).refresh();
          }
        },
        child: const Icon(Icons.add),
      ),
      body: _buildBody(context, theme, state),
    );
  }

  Widget _buildBody(BuildContext context, ThemeData theme, TransactionListState state) {
    switch (state.status) {
      case TransactionListStatus.initial:
      case TransactionListStatus.loading:
        return _buildSkeleton();
      case TransactionListStatus.loadingMore:
        return _buildListWithLoading(context, theme, state);
      case TransactionListStatus.loaded:
        return _buildList(context, theme, state);
      case TransactionListStatus.empty:
        return _buildEmpty(theme, state);
      case TransactionListStatus.error:
        return _buildError(theme, state);
      case TransactionListStatus.offline:
        return _buildOffline(theme);
    }
  }

  Widget _buildSkeleton() {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: 8,
      itemBuilder: (_, __) => Padding(
        padding: const EdgeInsets.only(bottom: 8),
        child: Card(
          child: SizedBox(
            height: 72,
            child: Center(
              child: CircularProgressIndicator(strokeWidth: 2, color: Colors.grey.withValues(alpha: 0.3)),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildList(BuildContext context, ThemeData theme, TransactionListState state) {
    if (state.transactions.isEmpty) return _buildEmpty(theme, state);
    return RefreshIndicator(
      onRefresh: () => ref.read(transactionListProvider.notifier).refresh(),
      child: ListView.builder(
        controller: _scrollCtrl,
        padding: const EdgeInsets.fromLTRB(16, 16, 16, 80),
        itemCount: state.transactions.length + (state.summary != null ? 1 : 0),
        itemBuilder: (context, index) {
          if (state.summary != null && index == 0) {
            return Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: TransactionSummaryCard(summary: state.summary!),
            );
          }
          final txIdx = state.summary != null ? index - 1 : index;
          final tx = state.transactions[txIdx];
          return Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: TransactionCard(
              transaction: tx,
              onTap: () => context.push('/transactions/${tx.eventId}'),
            ),
          );
        },
      ),
    );
  }

  Widget _buildListWithLoading(BuildContext context, ThemeData theme, TransactionListState state) {
    return Stack(
      children: [
        _buildList(context, theme, state),
        if (state.status == TransactionListStatus.loadingMore)
          Positioned(
            bottom: 0, left: 0, right: 0,
            child: Center(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: CircularProgressIndicator(strokeWidth: 2, color: theme.colorScheme.primary),
              ),
            ),
          ),
      ],
    );
  }

  Widget _buildEmpty(ThemeData theme, TransactionListState state) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              state.searchQuery != null ? Icons.search_off : Icons.receipt_long_outlined,
              size: 64,
              color: theme.colorScheme.onSurfaceVariant.withValues(alpha: 0.4),
            ),
            const SizedBox(height: 16),
            Text(
              state.searchQuery != null ? 'No transactions found' : 'No transactions yet',
              style: theme.textTheme.titleMedium,
            ),
            const SizedBox(height: 8),
            Text(
              state.searchQuery != null
                  ? 'Try a different search term or filters'
                  : 'Add your first transaction to get started',
              textAlign: TextAlign.center,
              style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant),
            ),
            if (state.searchQuery != null || state.activeFilters.hasActiveFilters) ...[
              const SizedBox(height: 16),
              TextButton.icon(
                onPressed: () => ref.read(transactionListProvider.notifier).clearFilters(),
                icon: const Icon(Icons.clear_all),
                label: const Text('Clear filters'),
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildError(ThemeData theme, TransactionListState state) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.cloud_off, size: 64, color: theme.colorScheme.error),
          const SizedBox(height: 16),
          Text('Could not load transactions', style: theme.textTheme.titleMedium),
          const SizedBox(height: 8),
          Text(state.error ?? '', textAlign: TextAlign.center, style: theme.textTheme.bodySmall),
          const SizedBox(height: 24),
          FilledButton.icon(
            onPressed: () => ref.read(transactionListProvider.notifier).refresh(),
            icon: const Icon(Icons.refresh),
            label: const Text('Try Again'),
          ),
        ],
      ),
    );
  }

  Widget _buildOffline(ThemeData theme) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.wifi_off, size: 64, color: theme.colorScheme.onSurfaceVariant),
          const SizedBox(height: 16),
          Text('You\'re offline', style: theme.textTheme.titleMedium),
          const SizedBox(height: 8),
          Text('Connect to the internet to view transactions', textAlign: TextAlign.center, style: theme.textTheme.bodySmall),
        ],
      ),
    );
  }
}
