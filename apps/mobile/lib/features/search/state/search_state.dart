import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/search_models.dart';
import '../repository/search_repository.dart';

final searchProvider = StateNotifierProvider<SearchNotifier, SearchState>((ref) {
  return SearchNotifier(ref.read(searchRepositoryProvider));
});

enum SearchStatus { idle, searching, loaded, error }

class SearchState {
  final SearchStatus status;
  final List<SearchResultItem> results;
  final List<String> recentSearches;
  final String? error;
  final String? query;

  const SearchState({
    this.status = SearchStatus.idle,
    this.results = const [],
    this.recentSearches = const [],
    this.error,
    this.query,
  });

  SearchState copyWith({
    SearchStatus? status,
    List<SearchResultItem>? results,
    List<String>? recentSearches,
    String? error,
    String? query,
    bool clearQuery = false,
    bool clearError = false,
  }) => SearchState(
    status: status ?? this.status,
    results: results ?? this.results,
    recentSearches: recentSearches ?? this.recentSearches,
    error: clearError ? null : error ?? this.error,
    query: clearQuery ? null : query ?? this.query,
  );
}

class SearchNotifier extends StateNotifier<SearchState> {
  final SearchRepository _repo;

  SearchNotifier(this._repo) : super(const SearchState());

  Future<void> search(String query, {String? moduleFilter}) async {
    if (query.isEmpty) {
      state = state.copyWith(status: SearchStatus.idle, results: [], clearQuery: true);
      return;
    }
    state = state.copyWith(status: SearchStatus.searching, query: query);
    try {
      final results = await _repo.searchAll(query);
      _saveRecent(query);
      state = state.copyWith(
        status: SearchStatus.loaded,
        results: results,
      );
    } catch (e) {
      state = state.copyWith(status: SearchStatus.error, error: e.toString());
    }
  }

  void clearSearch() {
    state = const SearchState();
  }

  void _saveRecent(String query) {
    final recent = List<String>.from(state.recentSearches);
    recent.remove(query);
    recent.insert(0, query);
    if (recent.length > 10) recent.removeLast();
    state = state.copyWith(recentSearches: recent);
  }

  void clearHistory() {
    state = state.copyWith(recentSearches: []);
  }
}
