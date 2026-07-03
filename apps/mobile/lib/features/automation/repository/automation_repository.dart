import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/automation_models.dart';

final automationRepositoryProvider = Provider<AutomationRepository>((ref) {
  return AutomationRepository(apiClient: ref.read(apiClientProvider));
});

class AutomationRepository {
  final ApiClient apiClient;
  AutomationRepository({required this.apiClient});

  Future<RuleListResponse> getRules() async {
    final resp = await apiClient.get('/rules');
    return RuleListResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<RuleDetailResponse> getRule(String id) async {
    final resp = await apiClient.get('/rules/$id');
    return RuleDetailResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<void> createRule(Map<String, dynamic> data) async {
    await apiClient.post('/rules', data: data);
  }

  Future<void> deleteRule(String id) async {
    await apiClient.delete('/rules/$id');
  }
}
