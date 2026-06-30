import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/dashboard_models.dart';
import '../repository/dashboard_repository.dart';

final dashboardStateProvider = StateNotifierProvider<DashboardStateNotifier, DashboardState>((ref) {
  final repo = ref.read(dashboardRepositoryProvider);
  return DashboardStateNotifier(repo);
});

enum DashboardLoadStatus { initial, loading, loaded, error, offline }

class DashboardState {
  final DashboardLoadStatus status;
  final DashboardData? dashboard;
  final SummaryData? summary;
  final String? error;
  final bool isOffline;

  const DashboardState({
    this.status = DashboardLoadStatus.initial,
    this.dashboard,
    this.summary,
    this.error,
    this.isOffline = false,
  });

  DashboardState copyWith({
    DashboardLoadStatus? status, DashboardData? dashboard, SummaryData? summary,
    String? error, bool? isOffline,
  }) => DashboardState(
    status: status ?? this.status,
    dashboard: dashboard ?? this.dashboard,
    summary: summary ?? this.summary,
    error: error ?? this.error,
    isOffline: isOffline ?? this.isOffline,
  );
}

class DashboardStateNotifier extends StateNotifier<DashboardState> {
  final DashboardRepository _repo;

  DashboardStateNotifier(this._repo) : super(const DashboardState());

  Future<void> load({String? userId}) async {
    state = state.copyWith(status: DashboardLoadStatus.loading, error: null);
    try {
      final results = await Future.wait([
        _repo.getDashboard(userId: userId),
        _repo.getSummary(userId: userId),
      ]);
      state = state.copyWith(
        status: DashboardLoadStatus.loaded,
        dashboard: (results[0] as DashboardResponse).data,
        summary: (results[1] as SummaryResponse).data,
      );
    } catch (e) {
      state = state.copyWith(
        status: DashboardLoadStatus.error,
        error: e.toString(),
      );
    }
  }

  Future<void> refresh() async {
    await load();
  }

  String get healthLabel {
    final s = state.summary;
    if (s == null) return '--';
    final change = _healthChange;
    if (change > 0) return '${s.healthScore} ($s.healthGrade) ↑$change';
    if (change < 0) return '${s.healthScore} ($s.healthGrade) ↓${-change}';
    return '${s.healthScore} ($s.healthGrade)';
  }

  int get _healthChange {
    // Placeholder — real change tracking from engine
    return 0;
  }
}
