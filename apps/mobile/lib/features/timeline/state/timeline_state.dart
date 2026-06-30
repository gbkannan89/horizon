import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/timeline_models.dart';
import '../repository/timeline_repository.dart';

final timelineStateProvider = StateNotifierProvider<TimelineStateNotifier, TimelineState>((ref) {
  return TimelineStateNotifier(ref.read(timelineRepositoryProvider));
});

enum TimelineStatus { initial, loading, loaded, loadingMore, error, empty, offline }

class TimelineState {
  final TimelineStatus status;
  final List<TimelineItem> items;
  final String? cursor;
  final bool hasMore;
  final String? error;
  final FilterParams activeFilters;
  final String view;

  const TimelineState({
    this.status = TimelineStatus.initial, this.items = const [], this.cursor,
    this.hasMore = false, this.error, this.activeFilters = const FilterParams(), this.view = 'ThisMonth',
  });

  TimelineState copyWith({
    TimelineStatus? status, List<TimelineItem>? items, String? cursor, bool? hasMore,
    String? error, FilterParams? activeFilters, String? view,
  }) => TimelineState(
    status: status ?? this.status, items: items ?? this.items, cursor: cursor ?? this.cursor,
    hasMore: hasMore ?? this.hasMore, error: error ?? this.error,
    activeFilters: activeFilters ?? this.activeFilters, view: view ?? this.view,
  );
}

class TimelineStateNotifier extends StateNotifier<TimelineState> {
  final TimelineRepository _repo;

  TimelineStateNotifier(this._repo) : super(const TimelineState());

  Future<void> load({String? userId}) async {
    state = state.copyWith(status: TimelineStatus.loading);
    try {
      final resp = await _repo.getTimeline(userId: userId, view: state.view);
      final feed = resp.data?.feed;
      if (feed == null || feed.items.isEmpty) {
        state = state.copyWith(status: TimelineStatus.empty, items: [], cursor: null, hasMore: false);
      } else {
        state = state.copyWith(status: TimelineStatus.loaded, items: feed.items, cursor: feed.cursor, hasMore: feed.hasMore);
      }
    } catch (e) {
      state = state.copyWith(status: TimelineStatus.error, error: e.toString());
    }
  }

  Future<void> loadMore({String? userId}) async {
    if (_isLoading || !state.hasMore) return;
    state = state.copyWith(status: TimelineStatus.loadingMore);
    try {
      final resp = await _repo.getTimeline(userId: userId, view: state.view, cursor: state.cursor);
      final feed = resp.data?.feed;
      if (feed != null) {
        state = state.copyWith(
          status: TimelineStatus.loaded,
          items: [...state.items, ...feed.items],
          cursor: feed.cursor, hasMore: feed.hasMore,
        );
      }
    } catch (_) {
      state = state.copyWith(status: TimelineStatus.loaded);
    }
  }

  Future<void> refresh({String? userId}) async {
    state = state.copyWith(status: TimelineStatus.loading, cursor: null, hasMore: false);
    await load(userId: userId);
  }

  Future<void> applyFilters({String? userId, required FilterParams filters}) async {
    state = state.copyWith(activeFilters: filters, status: TimelineStatus.loading, cursor: null);
    try {
      final resp = await _repo.getFiltered(userId: userId, filters: filters);
      final feed = resp.data?.feed;
      if (feed == null || feed.items.isEmpty) {
        state = state.copyWith(status: TimelineStatus.empty, items: []);
      } else {
        state = state.copyWith(status: TimelineStatus.loaded, items: feed.items, cursor: feed.cursor, hasMore: feed.hasMore);
      }
    } catch (e) {
      state = state.copyWith(status: TimelineStatus.error, error: e.toString());
    }
  }

  Future<void> search({String? userId, required String query}) async {
    state = state.copyWith(status: TimelineStatus.loading);
    try {
      final items = await _repo.search(userId: userId, query: query);
      state = state.copyWith(
        status: items.isEmpty ? TimelineStatus.empty : TimelineStatus.loaded,
        items: items, cursor: null, hasMore: false,
      );
    } catch (e) {
      state = state.copyWith(status: TimelineStatus.error, error: e.toString());
    }
  }

  bool get _isLoading => state.status == TimelineStatus.loading;
}
