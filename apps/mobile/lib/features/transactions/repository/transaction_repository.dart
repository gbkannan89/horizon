import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/transaction_models.dart';

final transactionRepositoryProvider = Provider<TransactionRepository>((ref) {
  return TransactionRepository(apiClient: ref.read(apiClientProvider));
});

class TransactionRepository {
  final ApiClient apiClient;

  TransactionRepository({required this.apiClient});

  Future<TransactionListResponse> getTransactions({
    String? userId,
    String? cursor,
    int limit = 50,
    TransactionFilter? filters,
  }) async {
    final params = <String, dynamic>{'limit': limit};
    if (userId != null) params['user_id'] = userId;
    if (cursor != null) params['cursor'] = cursor;
    if (filters != null) params.addAll(filters.toQuery());
    final response = await apiClient.get('/transactions', queryParameters: params);
    return TransactionListResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<TransactionDetailResponse> getTransaction({String? userId, required String id}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final response = await apiClient.get('/transactions/$id', queryParameters: params);
    return TransactionDetailResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<TransactionListResponse> searchTransactions({
    String? userId,
    required String query,
    int limit = 50,
  }) async {
    final params = <String, dynamic>{'q': query, 'limit': limit};
    if (userId != null) params['user_id'] = userId;
    final response = await apiClient.get('/transactions/search', queryParameters: params);
    return TransactionListResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<TransactionSummaryResponse> getSummary({String? userId}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final response = await apiClient.get('/transactions/summary', queryParameters: params);
    return TransactionSummaryResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<Map<String, dynamic>> createTransaction({String? userId, required Map<String, dynamic> data}) async {
    final body = <String, dynamic>{...data};
    if (userId != null) body['user_id'] = userId;
    final response = await apiClient.post('/events', data: body);
    return (response.data as Map<String, dynamic>?) ?? {};
  }

  Future<Map<String, dynamic>> archiveTransaction({String? userId, required String id}) async {
    final body = <String, dynamic>{};
    if (userId != null) body['user_id'] = userId;
    final response = await apiClient.post('/events/$id/archive', data: body);
    return (response.data as Map<String, dynamic>?) ?? {};
  }
}
