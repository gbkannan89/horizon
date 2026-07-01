import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/transaction_models.dart';
import '../repository/transaction_repository.dart';

final transactionListProvider = StateNotifierProvider<TransactionListNotifier, TransactionListState>((ref) {
  return TransactionListNotifier(ref.read(transactionRepositoryProvider));
});

enum TransactionListStatus { initial, loading, loaded, loadingMore, error, empty, offline }

class TransactionListState {
  final TransactionListStatus status;
  final List<TransactionEvent> transactions;
  final String? cursor;
  final bool hasMore;
  final String? error;
  final TransactionFilter activeFilters;
  final String? searchQuery;

  const TransactionListState({
    this.status = TransactionListStatus.initial,
    this.transactions = const [],
    this.cursor,
    this.hasMore = false,
    this.error,
    this.activeFilters = const TransactionFilter(),
    this.searchQuery,
  });

  TransactionListState copyWith({
    TransactionListStatus? status,
    List<TransactionEvent>? transactions,
    String? cursor,
    bool? hasMore,
    String? error,
    TransactionFilter? activeFilters,
    String? searchQuery,
    bool clearSearchQuery = false,
  }) => TransactionListState(
    status: status ?? this.status,
    transactions: transactions ?? this.transactions,
    cursor: cursor ?? this.cursor,
    hasMore: hasMore ?? this.hasMore,
    error: error ?? this.error,
    activeFilters: activeFilters ?? this.activeFilters,
    searchQuery: clearSearchQuery ? null : searchQuery ?? this.searchQuery,
  );
}

class TransactionListNotifier extends StateNotifier<TransactionListState> {
  final TransactionRepository _repo;

  TransactionListNotifier(this._repo) : super(const TransactionListState());

  Future<void> load({String? userId}) async {
    state = state.copyWith(status: TransactionListStatus.loading);
    try {
      final resp = await _repo.getTransactions(
        userId: userId,
        filters: state.activeFilters.hasActiveFilters ? state.activeFilters : null,
      );
      final data = resp.data;
      if (data == null || data.transactions.isEmpty) {
        state = state.copyWith(
          status: TransactionListStatus.empty,
          transactions: [],
          cursor: null,
          hasMore: false,
        );
      } else {
        state = state.copyWith(
          status: TransactionListStatus.loaded,
          transactions: data.transactions,
          cursor: data.cursor,
          hasMore: data.hasMore,
        );
      }
    } catch (e) {
      state = state.copyWith(status: TransactionListStatus.error, error: e.toString());
    }
  }

  Future<void> loadMore({String? userId}) async {
    if (_isLoading || !state.hasMore) return;
    state = state.copyWith(status: TransactionListStatus.loadingMore);
    try {
      final resp = await _repo.getTransactions(
        userId: userId,
        cursor: state.cursor,
        filters: state.activeFilters.hasActiveFilters ? state.activeFilters : null,
      );
      final data = resp.data;
      if (data != null) {
        state = state.copyWith(
          status: TransactionListStatus.loaded,
          transactions: [...state.transactions, ...data.transactions],
          cursor: data.cursor,
          hasMore: data.hasMore,
        );
      }
    } catch (_) {
      state = state.copyWith(status: TransactionListStatus.loaded);
    }
  }

  Future<void> refresh({String? userId}) async {
    state = state.copyWith(status: TransactionListStatus.loading, cursor: null, hasMore: false);
    await load(userId: userId);
  }

  Future<void> search({String? userId, required String query}) async {
    if (query.isEmpty) {
      await refresh(userId: userId);
      return;
    }
    state = state.copyWith(status: TransactionListStatus.loading, searchQuery: query);
    try {
      final resp = await _repo.searchTransactions(userId: userId, query: query);
      final data = resp.data;
      if (data == null || data.transactions.isEmpty) {
        state = state.copyWith(
          status: TransactionListStatus.empty,
          transactions: [],
          cursor: null,
          hasMore: false,
        );
      } else {
        state = state.copyWith(
          status: TransactionListStatus.loaded,
          transactions: data.transactions,
          cursor: null,
          hasMore: false,
        );
      }
    } catch (e) {
      state = state.copyWith(status: TransactionListStatus.error, error: e.toString());
    }
  }

  Future<void> applyFilters({String? userId, required TransactionFilter filters}) async {
    state = state.copyWith(activeFilters: filters, status: TransactionListStatus.loading, cursor: null);
    try {
      final resp = await _repo.getTransactions(userId: userId, filters: filters.hasActiveFilters ? filters : null);
      final data = resp.data;
      if (data == null || data.transactions.isEmpty) {
        state = state.copyWith(status: TransactionListStatus.empty, transactions: []);
      } else {
        state = state.copyWith(
          status: TransactionListStatus.loaded,
          transactions: data.transactions,
          cursor: data.cursor,
          hasMore: data.hasMore,
        );
      }
    } catch (e) {
      state = state.copyWith(status: TransactionListStatus.error, error: e.toString());
    }
  }

  void clearFilters({String? userId}) {
    state = state.copyWith(activeFilters: const TransactionFilter(), searchQuery: null);
    load(userId: userId);
  }

  bool get _isLoading => state.status == TransactionListStatus.loading;
}
