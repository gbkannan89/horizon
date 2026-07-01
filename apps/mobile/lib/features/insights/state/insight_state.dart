import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/insight_models.dart';
import '../repository/insight_repository.dart';

final insightProvider = StateNotifierProvider<InsightNotifier, InsightState>((ref) {
  return InsightNotifier(ref.read(insightRepositoryProvider));
});

enum InsightStatus { initial, loading, loaded, error, empty }

class InsightState {
  final InsightStatus status;
  final InsightDashboard? summary;
  final List<InsightItem> insights;
  final List<InsightItem> opportunities;
  final List<InsightItem> warnings;
  final List<InsightItem> achievements;
  final TrendData? trends;
  final String activeTab;
  final String? searchQuery;
  final String? error;

  const InsightState({
    this.status = InsightStatus.initial,
    this.summary,
    this.insights = const [],
    this.opportunities = const [],
    this.warnings = const [],
    this.achievements = const [],
    this.trends,
    this.activeTab = 'all',
    this.searchQuery,
    this.error,
  });

  InsightState copyWith({
    InsightStatus? status, InsightDashboard? summary, List<InsightItem>? insights,
    List<InsightItem>? opportunities, List<InsightItem>? warnings,
    List<InsightItem>? achievements, TrendData? trends, String? activeTab,
    String? searchQuery, String? error, bool clearError = false,
  }) => InsightState(
    status: status ?? this.status, summary: summary ?? this.summary,
    insights: insights ?? this.insights,
    opportunities: opportunities ?? this.opportunities,
    warnings: warnings ?? this.warnings,
    achievements: achievements ?? this.achievements,
    trends: trends ?? this.trends, activeTab: activeTab ?? this.activeTab,
    searchQuery: searchQuery ?? this.searchQuery,
    error: clearError ? null : error ?? this.error,
  );
}

class InsightNotifier extends StateNotifier<InsightState> {
  final InsightRepository _repo;
  InsightNotifier(this._repo) : super(const InsightState());

  Future<void> init() async {
    state = state.copyWith(status: InsightStatus.loading);
    try {
      final results = await Future.wait([_repo.getSummary(), _repo.getInsights()]);
      state = state.copyWith(
        status: InsightStatus.loaded,
        summary: (results[0] as InsightDashboardResponse).data,
        insights: (results[1] as InsightListResponse).data?.insights ?? [],
      );
    } catch (e) {
      state = state.copyWith(status: InsightStatus.error, error: e.toString());
    }
  }

  Future<void> loadTab(String tab) async {
    if (tab == state.activeTab) return;
    state = state.copyWith(activeTab: tab, status: InsightStatus.loading);
    try {
      switch (tab) {
        case 'opportunities':
          final r = await _repo.getOpportunities();
          state = state.copyWith(status: InsightStatus.loaded, opportunities: r.data?.insights ?? []);
        case 'warnings':
          final r = await _repo.getWarnings();
          state = state.copyWith(status: InsightStatus.loaded, warnings: r.data?.insights ?? []);
        case 'achievements':
          final r = await _repo.getAchievements();
          state = state.copyWith(status: InsightStatus.loaded, achievements: r.data?.insights ?? []);
        case 'forecast':
          final r = await _repo.getForecast();
          state = state.copyWith(status: InsightStatus.loaded, achievements: r.data?.insights ?? []);
        case 'trends':
          final t = await _repo.getTrends();
          state = state.copyWith(status: InsightStatus.loaded, trends: t);
        default:
          final r = await _repo.getInsights();
          state = state.copyWith(status: InsightStatus.loaded, insights: r.data?.insights ?? []);
      }
    } catch (e) {
      state = state.copyWith(status: InsightStatus.error, error: e.toString());
    }
  }

  Future<void> search(String query) async {
    if (query.isEmpty) { state = state.copyWith(searchQuery: null); await init(); return; }
    state = state.copyWith(status: InsightStatus.loading, searchQuery: query);
    try {
      final r = await _repo.search(query);
      state = state.copyWith(
        status: r.data?.insights.isNotEmpty == true ? InsightStatus.loaded : InsightStatus.empty,
        insights: r.data?.insights ?? [],
      );
    } catch (e) {
      state = state.copyWith(status: InsightStatus.error, error: e.toString());
    }
  }

  Future<void> refresh() async {
    state = state.copyWith(status: InsightStatus.loading, searchQuery: null);
    await init();
  }
}
