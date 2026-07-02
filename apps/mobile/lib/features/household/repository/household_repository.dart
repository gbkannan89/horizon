import 'package:dio/dio.dart';
import 'package:horizon_mobile/features/budget/models/budget_models.dart';
import 'package:horizon_mobile/features/household/models/household_models.dart';
import 'package:horizon_mobile/features/goals/models/goal_models.dart';

class HouseholdRepository {
  final Dio _dio;

  HouseholdRepository(this._dio);

  Future<List<HouseholdView>> getHouseholds({int limit = 25, String? cursor}) async {
    final response = await _dio.get('/api/v1/households', queryParameters: {
      'limit': limit,
      if (cursor != null) 'cursor': cursor,
    });

    if (response.data['success'] == true) {
      final jsonItems = response.data['data']['households'] as List? ?? [];
      return jsonItems.map((e) => HouseholdView.fromJson(e)).toList();
    }
    throw Exception(response.data['error']?['message'] ?? 'Failed to load households');
  }

  Future<HouseholdDetailView> getHouseholdDetail(String id) async {
    final response = await _dio.get('/api/v1/households/$id');
    if (response.data['success'] == true) {
      return HouseholdDetailView.fromJson(response.data['data']);
    }
    throw Exception(response.data['error']?['message'] ?? 'Failed to load household details');
  }

  Future<HouseholdFinancialSummary> getHouseholdSummary(String id) async {
    final response = await _dio.get('/api/v1/households/$id/summary');
    if (response.data['success'] == true) {
      return HouseholdFinancialSummary.fromJson(response.data['data']);
    }
    throw Exception(response.data['error']?['message'] ?? 'Failed to load household summary');
  }

  Future<void> createHousehold(String name, String type, String currency, String country) async {
    final response = await _dio.post('/api/v1/households', data: {
      'name': name,
      'household_type': type,
      'currency': currency,
      'country': country,
      'head_of_household_id': 'mobile_user',
    });
    if (response.data['success'] != true) {
      throw Exception(response.data['error']?['message'] ?? 'Failed to create household');
    }
  }

  Future<List<GoalSummary>> getHouseholdGoals(String householdId) async {
    final response = await _dio.get('/api/v1/households/$householdId/goals');
    final data = response.data;
    final goals = (data['goals'] as List?)?.map((e) => GoalSummary.fromJson(e)).toList() ?? [];
    return goals;
  }

  Future<List<BudgetModel>> getHouseholdBudgets(String householdId) async {
    final response = await _dio.get('/api/v1/households/$householdId/budgets');
    final data = response.data;
    final budgets = (data['data']['budgets'] as List?)?.map((e) => BudgetModel.fromJson(e)).toList() ?? [];
    return budgets;
  }

  Future<void> inviteMember(String householdId, String userId, String role) async {
    final response = await _dio.post('/api/v1/households/$householdId/members', data: {
      'user_id': userId,
      'role': role,
    });
    if (response.data['success'] != true) {
      throw Exception(response.data['error']?['message'] ?? 'Failed to invite member');
    }
  }

  Future<void> acceptInvite(String householdId, String userId) async {
    final response = await _dio.post('/api/v1/households/$householdId/members/$userId/accept');
    if (response.data['success'] != true) {
      throw Exception(response.data['error']?['message'] ?? 'Failed to accept invite');
    }
  }

  Future<void> removeMember(String householdId, String userId) async {
    final response = await _dio.delete('/api/v1/households/$householdId/members/$userId');
    if (response.data['success'] != true) {
      throw Exception(response.data['error']?['message'] ?? 'Failed to remove member');
    }
  }

  Future<void> updateMemberRole(String householdId, String userId, String role) async {
    final response = await _dio.put('/api/v1/households/$householdId/members/$userId/role', data: {
      'role': role,
    });
    if (response.data['success'] != true) {
      throw Exception(response.data['error']?['message'] ?? 'Failed to update member role');
    }
  }

  Future<void> dissolveHousehold(String householdId) async {
    final response = await _dio.post('/api/v1/households/$householdId/dissolve');
    if (response.data['success'] != true) {
      throw Exception(response.data['error']?['message'] ?? 'Failed to dissolve household');
    }
  }
}
