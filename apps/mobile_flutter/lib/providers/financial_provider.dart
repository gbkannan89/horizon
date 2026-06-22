import 'package:flutter/material.dart';
import '../models/financial_models.dart';
import '../services/api_service.dart';

class FinancialProvider extends ChangeNotifier {
  final ApiService _apiService = ApiService();
  bool _isLoading = false;

  List<LocalIncome> incomes = [];
  List<LocalAsset> assets = [];
  List<LocalLiability> liabilities = [];

  List<LocalBill> bills = [];
  List<LocalHouseholdMember> householdMembers = [];

  // Dashboard Fields
  double netWorth = 0;
  int finScoreVal = 0;
  double totalIncomeAgg = 0;
  double totalSpent = 0;
  double totalLeft = 0;
  double needsSpent = 0;
  double wantsSpent = 0;
  double savingsSpent = 0;
  double needsBudget = 1;
  double wantsBudget = 1;
  double savingsBudget = 1;

  List<dynamic> goals = [];
  List<dynamic> recentExpenses = [];

  bool get isLoading => _isLoading;

  Future<void> loadAllData() async {
    _isLoading = true;
    notifyListeners();

    try {
      final dashboardData = await _apiService.get('/api/v1/dashboard/overview');
      netWorth = (dashboardData['net_worth'] ?? 0).toDouble();
      finScoreVal = dashboardData['fin_score'] ?? 0;
      totalIncomeAgg = (dashboardData['total_income'] ?? 0).toDouble();
      totalSpent = (dashboardData['total_spent'] ?? 0).toDouble();
      totalLeft = (dashboardData['total_left'] ?? 0).toDouble();
      needsSpent = (dashboardData['needs_spent'] ?? 0).toDouble();
      wantsSpent = (dashboardData['wants_spent'] ?? 0).toDouble();
      savingsSpent = (dashboardData['savings_spent'] ?? 0).toDouble();
      needsBudget = (dashboardData['needs_budget'] ?? 1).toDouble();
      wantsBudget = (dashboardData['wants_budget'] ?? 1).toDouble();
      savingsBudget = (dashboardData['savings_budget'] ?? 1).toDouble();
      goals = dashboardData['goals'] ?? [];
      recentExpenses = dashboardData['recent_expenses'] ?? [];
    } catch (e) {
      print('Dashboard load error: $e');
    }

    try {
      final assetsData = await _apiService.get('/api/assets');
      if (assetsData is List) {
        assets = assetsData.map((a) => LocalAsset(
          id: a['id'].toString(), name: a['name'], amount: (a['amount'] ?? 0).toDouble()
        )).toList();
      }
    } catch (e) { print('Assets load error: $e'); }

    try {
      final liabilitiesData = await _apiService.get('/api/liabilities');
      if (liabilitiesData is List) {
        liabilities = liabilitiesData.map((l) => LocalLiability(
          id: l['id'].toString(), name: l['name'],
          amount: (l['outstanding'] ?? 0).toDouble(),
          interestRate: (l['interest_rate'] ?? 0).toDouble()
        )).toList();
      }
    } catch (e) { print('Liabilities load error: $e'); }

    try {
      final incomesData = await _apiService.get('/api/incomes');
      if (incomesData is List) {
        incomes = incomesData.map((i) => LocalIncome(
          id: i['id'].toString(),
          label: i['label'] ?? _typeToLabel(i['type'] ?? 'salary'),
          type: i['type'] ?? 'salary',
          amount: (i['amount'] ?? 0).toDouble(),
          frequency: i['frequency'] ?? 'monthly',
        )).toList();
      }
    } catch (e) { print('Incomes load error: $e'); }

    try {
      final hhData = await _apiService.get('/api/household/summary');
      if (hhData['contributingMembers'] != null) {
        final members = hhData['contributingMembers'] as List;
        householdMembers = members.map((m) => LocalHouseholdMember(
          id: m['id'].toString(),
          name: m['name'],
          monthlyIncome: (m['monthlyIncome'] ?? 0).toDouble(),
          contribution: (m['contributionToHousehold'] ?? 0).toDouble(),
          relationship: m['relationship'] ?? 'Member',
        )).toList();
      }
    } catch (e) {
      print('Household load error: $e');
      householdMembers = [];
    }

    _isLoading = false;
    notifyListeners();
  }

  static String _typeToLabel(String type) {
    switch (type) {
      case 'salary':   return 'Salary';
      case 'business': return 'Business';
      case 'passive':  return 'Passive';
      default: return type[0].toUpperCase() + type.substring(1);
    }
  }

  // ── EXPENSE ────────────────────────────────────────────────────────────────
  Future<void> addExpense(String name, double amount, String category, String bucket) async {
    await _apiService.post('/api/v1/dashboard/expense', {
      'name': name, 'amount': amount, 'category': category, 'bucket': bucket,
      'date': DateTime.now().toIso8601String().split('T').first, 'icon': 'receipt'
    });
    await loadAllData();
  }

  // ── INCOME ─────────────────────────────────────────────────────────────────
  Future<void> addIncome(String label, String type, double amount, String frequency) async {
    await _apiService.post('/api/incomes', {
      'label': label, 'type': type, 'amount': amount, 'frequency': frequency,
    });
    await loadAllData();
  }

  Future<void> deleteIncome(String id) async {
    await _apiService.delete('/api/incomes/$id');
    await loadAllData();
  }

  // ── ASSETS ─────────────────────────────────────────────────────────────────
  Future<void> addAsset(String name, String type, double amount) async {
    await _apiService.post('/api/assets', {
      'name': name, 'type': type, 'amount': amount,
    });
    await loadAllData();
  }

  Future<void> deleteAsset(String id) async {
    await _apiService.delete('/api/assets/$id');
    await loadAllData();
  }

  // ── LIABILITIES ────────────────────────────────────────────────────────────
  Future<void> addLiability(String name, String type, double outstanding, double emi, double interestRate) async {
    await _apiService.post('/api/liabilities', {
      'name': name, 'type': type, 'outstanding': outstanding,
      'emi': emi, 'interest_rate': interestRate,
    });
    await loadAllData();
  }

  Future<void> deleteLiability(String id) async {
    await _apiService.delete('/api/liabilities/$id');
    await loadAllData();
  }

  // ── HOUSEHOLD MEMBERS ──────────────────────────────────────────────────────
  Future<void> addContributingMember(String name, double monthlyIncome, double contribution, String relationship) async {
    await _apiService.post('/api/household/contributing-members', {
      'name': name, 'monthly_income': monthlyIncome,
      'contribution_to_household': contribution, 'relationship': relationship
    });
    await loadAllData();
  }
}
