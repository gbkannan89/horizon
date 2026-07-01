import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/search_models.dart';
import '../repository/search_repository.dart';

final searchStateProvider = StateNotifierProvider<SearchStateNotifier, SearchState>((ref) {
  return SearchStateNotifier(ref.read(searchRepositoryProvider));
});

enum SearchStatus { initial, loading, loaded, empty, error }

class SearchState {
  final SearchStatus status;
  final List<SearchResult> results;
  final List<RecentSearch> recentSearches;
  final SearchFilters filters;
  final String query;
  final String? error;

  const SearchState({
    this.status = SearchStatus.initial,
    this.results = const [],
    this.recentSearches = const [],
    this.filters = const SearchFilters(),
    this.query = '',
    this.error,
  });

  SearchState copyWith({
    SearchStatus? status, List<SearchResult>? results, List<RecentSearch>? recentSearches,
    SearchFilters? filters, String? query, String? error,
  }) => SearchState(
    status: status ?? this.status,
    results: results ?? this.results,
    recentSearches: recentSearches ?? this.recentSearches,
    filters: filters ?? this.filters,
    query: query ?? this.query,
    error: error ?? this.error,
  );
}

class SearchStateNotifier extends StateNotifier<SearchState> {
  final SearchRepository _repo;

  SearchStateNotifier(this._repo) : super(SearchState());

  Future<void> search(String query) async {
    if (query.trim().isEmpty) {
      state = state.copyWith(status: SearchStatus.initial, results: [], query: '');
      return;
    }
    state = state.copyWith(status: SearchStatus.loading, query: query, error: null);
    try {
      final results = await _repo.search(query: query, modules: state.filters.modules.isNotEmpty ? state.filters.modules : null);
      if (results.isEmpty) {
        state = state.copyWith(status: SearchStatus.empty, results: []);
      } else {
        state = state.copyWith(status: SearchStatus.loaded, results: results);
        _saveRecentSearch(query);
      }
    } catch (e) {
      state = state.copyWith(status: SearchStatus.error, error: e.toString());
    }
  }

  void applyFilters(SearchFilters filters) {
    state = state.copyWith(filters: filters);
    if (state.query.isNotEmpty) search(state.query);
  }

  void clearFilters() {
    state = state.copyWith(filters: const SearchFilters());
    if (state.query.isNotEmpty) search(state.query);
  }

  void clearHistory() {
    state = state.copyWith(recentSearches: []);
  }

  void _saveRecentSearch(String query) {
    final existing = List<RecentSearch>.from(state.recentSearches);
    existing.removeWhere((r) => r.query == query);
    existing.insert(0, RecentSearch(query: query, timestamp: DateTime.now()));
    if (existing.length > 10) existing.removeLast();
    state = state.copyWith(recentSearches: existing);
  }
}
