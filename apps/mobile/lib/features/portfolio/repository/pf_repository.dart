import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/pf_models.dart';

final pfRepositoryProvider = Provider<PfRepository>((ref) {
  return PfRepository(apiClient: ref.read(apiClientProvider));
});

class PfRepository {
  final ApiClient apiClient;
  PfRepository({required this.apiClient});

  Future<PfDashboardResponse> getDashboard({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/experience', queryParameters: p);
    return PfDashboardResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<AllocationResponse> getAllocation({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/allocation', queryParameters: p);
    return AllocationResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<PerformanceResponse> getPerformance({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/performance', queryParameters: p);
    return PerformanceResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<RiskResponse> getRisk({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/risk', queryParameters: p);
    return RiskResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<ProjectionResponse> getProjection({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/projection', queryParameters: p);
    return ProjectionResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<CardViewResponse> getRecommendations({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/recommendations', queryParameters: p);
    return CardViewResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<CardViewResponse> getOptimization({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/optimization', queryParameters: p);
    return CardViewResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<SimResponse> getSimulations({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/simulations', queryParameters: p);
    return SimResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<CardViewResponse> getTimeline({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/portfolio/timeline', queryParameters: p);
    return CardViewResponse.fromJson(r.data as Map<String, dynamic>);
  }
}
