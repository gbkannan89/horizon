import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/goal_models.dart';

final goalRepositoryProvider = Provider<GoalRepository>((ref) {
  return GoalRepository(apiClient: ref.read(apiClientProvider));
});

class GoalRepository {
  final ApiClient apiClient;
  GoalRepository({required this.apiClient});

  Future<GoalsListResponse> getGoals({String? userId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final resp = await apiClient.get('/goals/experience', queryParameters: params);
    return GoalsListResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<GoalDashboardResponse> getDashboard({String? userId, required String goalId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final resp = await apiClient.get('/goals/$goalId/dashboard', queryParameters: params);
    return GoalDashboardResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<GoalProgressResponse> getProgress({String? userId, required String goalId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final resp = await apiClient.get('/goals/$goalId/progress', queryParameters: params);
    return GoalProgressResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<GoalProjectionResponse> getProjection({String? userId, required String goalId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final resp = await apiClient.get('/goals/$goalId/projection', queryParameters: params);
    return GoalProjectionResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<GoalRecsResponse> getRecommendations({String? userId, required String goalId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final resp = await apiClient.get('/goals/$goalId/recommendations', queryParameters: params);
    return GoalRecsResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<GoalOptsResponse> getOptimization({String? userId, required String goalId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final resp = await apiClient.get('/goals/$goalId/optimization', queryParameters: params);
    return GoalOptsResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<GoalTimelineResponse> getTimeline({String? userId, required String goalId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final resp = await apiClient.get('/goals/$goalId/timeline', queryParameters: params);
    return GoalTimelineResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<GoalMilestoneResponse> getMilestones({String? userId, required String goalId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final resp = await apiClient.get('/goals/$goalId/milestones', queryParameters: params);
    return GoalMilestoneResponse.fromJson(resp.data as Map<String, dynamic>);
  }
}
