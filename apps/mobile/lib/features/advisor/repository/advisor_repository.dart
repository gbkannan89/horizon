import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/advisor_models.dart';

final advisorRepositoryProvider = Provider<AdvisorRepository>((ref) {
  return AdvisorRepository(apiClient: ref.read(apiClientProvider));
});

class AdvisorRepository {
  final ApiClient apiClient;

  AdvisorRepository({required this.apiClient});

  Future<ChatResponse> sendMessage({String? sessionId, required String message, Map<String, dynamic>? context}) async {
    final body = <String, dynamic>{'message': message};
    if (sessionId != null) body['session_id'] = sessionId;
    if (context != null) body['context'] = context;
    final response = await apiClient.post('/ai/chat', data: body);
    return ChatResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<AdvisorContext> getContext() async {
    final response = await apiClient.get('/ai/context');
    final json = response.data as Map<String, dynamic>;
    final data = json['data'] as Map<String, dynamic>? ?? json;
    return AdvisorContext.fromJson(data);
  }

  Future<List<SuggestedPrompt>> getPrompts() async {
    final response = await apiClient.get('/ai/prompts');
    final json = response.data as Map<String, dynamic>;
    final data = json['data'] as List? ?? [];
    return data.map((e) => SuggestedPrompt(
      id: (e['id'] ?? '').toString(),
      title: (e['title'] ?? '').toString(),
      prompt: (e['prompt'] ?? e['title'] ?? '').toString(),
      category: (e['category'] ?? 'general').toString(),
    )).toList();
  }

  Future<List<ProviderInfo>> getProviders() async {
    final response = await apiClient.get('/ai/providers');
    final json = response.data as Map<String, dynamic>;
    final data = json['data'] as List? ?? [];
    return data.map((e) => ProviderInfo.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<void> switchProvider(String provider) async {
    await apiClient.post('/ai/provider/switch', data: {'provider': provider});
  }
}
