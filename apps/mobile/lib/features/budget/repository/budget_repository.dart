import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/budget_models.dart';

final budgetRepositoryProvider = Provider<BudgetRepository>((ref) {
  return BudgetRepository(apiClient: ref.read(apiClientProvider));
});

class BudgetRepository {
  final ApiClient apiClient;
  BudgetRepository({required this.apiClient});

  Future<BudgetListResponse> getBudgets() async {
    final resp = await apiClient.get('/budgets');
    return BudgetListResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<BudgetDetailResponse> getBudget(String budgetId) async {
    final resp = await apiClient.get('/budgets/$budgetId');
    return BudgetDetailResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<BudgetDetailResponse> getBudgetVsActual(String budgetId) async {
    final resp = await apiClient.get('/budgets/$budgetId/vs-actual');
    return BudgetDetailResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<BudgetResultResponse> createBudget(CreateBudgetRequest request) async {
    final resp = await apiClient.post('/budgets', data: request.toJson());
    return BudgetResultResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<BudgetResultResponse> activateBudget(String budgetId) async {
    final resp = await apiClient.post('/budgets/$budgetId/activate');
    return BudgetResultResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<BudgetResultResponse> pauseBudget(String budgetId) async {
    final resp = await apiClient.post('/budgets/$budgetId/pause');
    return BudgetResultResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<BudgetResultResponse> resumeBudget(String budgetId) async {
    final resp = await apiClient.post('/budgets/$budgetId/resume');
    return BudgetResultResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<BudgetResultResponse> completeBudget(String budgetId) async {
    final resp = await apiClient.post('/budgets/$budgetId/complete');
    return BudgetResultResponse.fromJson(resp.data as Map<String, dynamic>);
  }

  Future<BudgetResultResponse> archiveBudget(String budgetId) async {
    final resp = await apiClient.post('/budgets/$budgetId/archive');
    return BudgetResultResponse.fromJson(resp.data as Map<String, dynamic>);
  }
}
