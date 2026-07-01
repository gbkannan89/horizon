import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/planning_models.dart';

final planningRepositoryProvider = Provider<PlanningRepository>((ref) {
  return PlanningRepository(apiClient: ref.read(apiClientProvider));
});

class PlanningRepository {
  final ApiClient apiClient;
  PlanningRepository({required this.apiClient});

  Future<PlanningDashboardResponse> getDashboard() async {
    final r = await apiClient.get('/planning/dashboard');
    return PlanningDashboardResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<ProjectionData> getProjections() async {
    final r = await apiClient.get('/planning/projections');
    final json = r.data as Map<String, dynamic>;
    final data = json['data'] as Map<String, dynamic>? ?? json;
    return ProjectionData.fromJson(data);
  }

  Future<List<Scenario>> getScenarios() async {
    final r = await apiClient.get('/planning/scenarios');
    final json = r.data as Map<String, dynamic>;
    final data = json['data'] as Map<String, dynamic>? ?? json;
    final list = data['scenarios'] as List? ?? data['data'] as List? ?? [];
    return list.map((e) => Scenario.fromJson(e as Map<String, dynamic>)).toList();
  }
}
