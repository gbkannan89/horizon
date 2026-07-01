import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/insight_models.dart';

final insightRepositoryProvider = Provider<InsightRepository>((ref) {
  return InsightRepository(apiClient: ref.read(apiClientProvider));
});

class InsightRepository {
  final ApiClient apiClient;
  InsightRepository({required this.apiClient});

  Future<InsightDashboardResponse> getSummary() async {
    final r = await apiClient.get('/insights/summary');
    return InsightDashboardResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<InsightListResponse> getInsights() async {
    final r = await apiClient.get('/insights');
    return InsightListResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<InsightListResponse> getOpportunities() async {
    final r = await apiClient.get('/insights/opportunities');
    return InsightListResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<InsightListResponse> getWarnings() async {
    final r = await apiClient.get('/insights/warnings');
    return InsightListResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<InsightListResponse> getAchievements() async {
    final r = await apiClient.get('/insights/achievements');
    return InsightListResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<InsightListResponse> getForecast() async {
    final r = await apiClient.get('/insights/forecast');
    return InsightListResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<TrendData> getTrends() async {
    final r = await apiClient.get('/insights/trends');
    final json = r.data as Map<String, dynamic>;
    final data = json['data'] as Map<String, dynamic>? ?? json;
    return TrendData.fromJson(data);
  }

  Future<InsightListResponse> search(String query) async {
    final r = await apiClient.get('/insights/search', queryParameters: {'q': query});
    return InsightListResponse.fromJson(r.data as Map<String, dynamic>);
  }
}
