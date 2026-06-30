import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/acct_models.dart';

final acctRepositoryProvider = Provider<AcctRepository>((ref) {
  return AcctRepository(apiClient: ref.read(apiClientProvider));
});

class AcctRepository {
  final ApiClient apiClient;
  AcctRepository({required this.apiClient});

  Future<AcctDashboardResponse> getDashboard({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/accounts/experience', queryParameters: p);
    return AcctDashboardResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<AcctListResponse> getAccounts({String? userId, String? cursor, int limit = 25, String? type, String? status}) async {
    final p = <String, dynamic>{'limit': limit}; if (userId != null) p['user_id'] = userId;
    if (cursor != null) p['cursor'] = cursor; if (type != null) p['type'] = type; if (status != null) p['status'] = status;
    final r = await apiClient.get('/accounts/summary', queryParameters: p);
    return AcctListResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<AcctDetailResponse> getDetail({String? userId, required String accountId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/accounts/$accountId', queryParameters: p);
    return AcctDetailResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<BalanceSummaryResponse> getBalances({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/accounts/balances', queryParameters: p);
    return BalanceSummaryResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<CashFlowResponse> getCashFlow({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/accounts/cashflow', queryParameters: p);
    return CashFlowResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<AcctHealthResponse> getHealth({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/accounts/health', queryParameters: p);
    return AcctHealthResponse.fromJson(r.data as Map<String, dynamic>);
  }

  Future<AcctListResponse> getRecommendations({String? userId}) async { return AcctListResponse(success: true); }
  Future<AcctListResponse> getProjections({String? userId}) async { return AcctListResponse(success: true); }
  Future<AcctListResponse> getTimeline({String? userId}) async { return AcctListResponse(success: true); }
}
