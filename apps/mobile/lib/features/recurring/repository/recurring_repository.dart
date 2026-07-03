import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/features/recurring/models/recurring_models.dart';

final recurringRepositoryProvider = Provider<RecurringRepository>((ref) {
  return RecurringRepository(apiClient: ref.read(apiClientProvider));
});

class RecurringRepository {
  final ApiClient apiClient;

  RecurringRepository({required this.apiClient});

  Future<RecurringListResponse> getRecurringTransactions() async {
    final resp = await apiClient.get('/recurring');
    return RecurringListResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<RecurringDetailResponse> getRecurringTransaction(String id) async {
    final resp = await apiClient.get('/recurring/$id');
    return RecurringDetailResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<void> createRecurring(Map<String, dynamic> data) async {
    await apiClient.post('/recurring', data: data);
  }

  Future<RecurringActionResponse> activateRecurring(String id) async {
    final resp = await apiClient.post('/recurring/$id/activate');
    return RecurringActionResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<RecurringActionResponse> pauseRecurring(String id) async {
    final resp = await apiClient.post('/recurring/$id/pause');
    return RecurringActionResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<RecurringActionResponse> cancelRecurring(String id) async {
    final resp = await apiClient.post('/recurring/$id/cancel');
    return RecurringActionResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<RecurringActionResponse> archiveRecurring(String id) async {
    final resp = await apiClient.post('/recurring/$id/archive');
    return RecurringActionResponse.fromJson(resp.data as Map<String, dynamic>);
  }
}
