import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/dashboard_models.dart';

final dashboardRepositoryProvider = Provider<DashboardRepository>((ref) {
  return DashboardRepository(apiClient: ref.read(apiClientProvider));
});

class DashboardRepository {
  final ApiClient apiClient;

  DashboardRepository({required this.apiClient});

  Future<DashboardResponse> getDashboard({String? userId, String? mode}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    if (mode != null) params['mode'] = mode;
    final response = await apiClient.get('/dashboard', queryParameters: params);
    return DashboardResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<SummaryResponse> getSummary({String? userId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final response = await apiClient.get('/dashboard/summary', queryParameters: params);
    return SummaryResponse.fromJson(response.data as Map<String, dynamic>);
  }
}
