import 'dart:io';
import 'package:flutter/material.dart';
import '../models/financial_models.dart';
import '../models/wishlist_item.dart';
import '../services/api_service.dart';

class FinancialProvider extends ChangeNotifier {
  final ApiService _apiService = ApiService();
  bool _isLoading = false;

  List<LocalIncome> incomes = [];
  List<LocalAsset> assets = [];
  List<LocalLiability> liabilities = [];

  List<LocalBill> bills = [];
  List<LocalHouseholdMember> householdMembers = [];
  List<LocalVehicle> vehicles = [];
  List<dynamic> insurances = [];

  // Dashboard Fields
  double netWorth = 0;
  double netWorthChange = 0;
  String netWorthChangePeriod = '';
  int finScoreVal = 0;
  List<dynamic> scoreHistory = [];
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
  List<dynamic> upcomingBills = [];
  List<WishlistItem> wishlist = [];

  // Analytics
  List<dynamic> spendingTrend = [];
  Map<String, dynamic>? budgetBreakdown;
  Map<String, dynamic>? portfolioSummary;

  int selectedMonth = DateTime.now().month;
  int selectedYear = DateTime.now().year;

  bool get isLoading => _isLoading;

  void nextMonth() {
    if (selectedMonth == 12) {
      selectedMonth = 1;
      selectedYear++;
    } else {
      selectedMonth++;
    }
    loadAllData();
  }

  void previousMonth() {
    if (selectedMonth == 1) {
      selectedMonth = 12;
      selectedYear--;
    } else {
      selectedMonth--;
    }
    loadAllData();
  }

  Future<void> loadAllData() async {
    _isLoading = true;
    notifyListeners();

    try {
      final dashboardData = await _apiService.get('/api/v1/dashboard/overview?month=$selectedMonth&year=$selectedYear');
      netWorth = (dashboardData['net_worth'] ?? 0).toDouble();
      netWorthChange = (dashboardData['net_worth_change'] ?? 0).toDouble();
      netWorthChangePeriod = dashboardData['net_worth_change_period'] ?? '';
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
      upcomingBills = dashboardData['upcoming_bills'] ?? [];
      spendingTrend = dashboardData['spending_trend'] ?? [];
    } catch (e) {
      print('Dashboard load error: $e');
    }

    try {
      final assetsData = await _apiService.get('/api/assets');
      if (assetsData is List) {
        assets = assetsData.map((a) => LocalAsset(
          id: a['id'].toString(), 
          name: a['name'], 
          amount: (a['amount'] ?? 0).toDouble(),
          interestRate: (a['interest_rate'] ?? 0).toDouble(),
          isLiability: a['is_liability'] ?? false,
          generatesIncome: a['generates_income'] ?? false,
          purchasePrice: a['purchase_price']?.toDouble(),
          purchaseDate: a['purchase_date'],
        )).toList();
      }
    } catch (e) { print('Assets load error: $e'); }

    try {
      portfolioSummary = await _apiService.get('/api/analytics/portfolio-summary');
    } catch (e) { print('Portfolio summary load error: $e'); }

    try {
      final vehiclesData = await _apiService.get('/api/assets/vehicles');
      if (vehiclesData is List) {
        vehicles = vehiclesData.map((v) => LocalVehicle(
          id: v['id'].toString(),
          makeModel: v['make_model'],
          purchaseCost: (v['purchase_cost'] ?? 0).toDouble(),
          insuranceRenewalDate: v['insurance_renewal_date'] != null ? DateTime.parse(v['insurance_renewal_date']) : null,
        )).toList();
      }
    } catch (e) { print('Vehicles load error: $e'); }

    try {
      final billsData = await _apiService.get('/api/bills');
      if (billsData is List) {
        bills = billsData.map((b) => LocalBill(
          id: b['id'].toString(),
          name: b['name'],
          amount: (b['amount'] ?? 0).toDouble(),
          frequency: b['frequency'] ?? 'monthly',
          dueDate: DateTime.now(), // Fallback since due_day is an int, ideally construct full date
          isEmi: b['is_emi'] ?? false,
          emiTotalMonths: b['emi_total_months'] ?? 0,
          emiMonthsPaid: b['emi_months_paid'] ?? 0,
        )).toList();
      }
    } catch (e) { print('Bills load error: $e'); }

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

    try {
      final wData = await _apiService.get('/api/discipline/wishlist');
      if (wData is List) {
        wishlist = wData.map((w) => WishlistItem.fromJson(w as Map<String, dynamic>)).toList();
      }
    } catch (e) {
      print('Bills load error: $e');
    }

    try {
      final insurancesData = await _apiService.get('/api/insurance');
      if (insurancesData is List) {
        insurances = insurancesData;
      }
    } catch (e) {
      print('Insurance load error: $e');
    }

    try {
      final shData = await _apiService.get('/api/score/history');
      if (shData is List) {
        scoreHistory = shData.reversed.toList(); // Oldest first for charts
      }
    } catch (e) {
      print('Score history load error: $e');
    }

    // Load background data
    loadBudgetBreakdown();
    loadNudges();

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

  // ── RECURRING BILLS ────────────────────────────────────────────────────────
  Future<void> addRecurringBill(String name, double amount, String category, String bucket, {
    String frequency = 'monthly',
    bool isEmi = false,
    int? emiTotalMonths,
    DateTime? startDate,
  }) async {
    await _apiService.post('/api/bills', {
      'name': name, 
      'amount': amount, 
      'category': category, 
      'bucket': bucket,
      'frequency': frequency,
      'is_emi': isEmi,
      'emi_total_months': emiTotalMonths,
      'start_date': startDate?.toIso8601String().split('T').first,
    });
    await loadAllData();
  }

  Future<void> deleteRecurringBill(String id) async {
    await _apiService.delete('/api/bills/$id');
    await loadAllData();
  }

  // ── INSURANCE ─────────────────────────────────────────────────────────────
  Future<void> addInsurance(Map<String, dynamic> data) async {
    await _apiService.post('/insurance/', data);
    await loadAllData();
  }

  Future<void> deleteInsurance(String id) async {
    await _apiService.delete('/insurance/$id');
    await loadAllData();
  }

  // ── TRANSACTIONS ───────────────────────────────────────────────────────────
  Future<void> uploadStatement(String path) async {
    await _apiService.uploadFile('/api/transactions/upload', 'file', path);
    await loadAllData();
  }

  // ── ASSETS ─────────────────────────────────────────────────────────────────
  Future<void> addAsset(String name, String type, double amount, {
    double interestRate = 0.0,
    bool isLiability = false,
    bool generatesIncome = false,
    String? incomeFrequency,
    double? purchasePrice,
    String? purchaseDate,
  }) async {
    await _apiService.post('/api/assets', {
      'name': name, 'type': type, 'amount': amount,
      'interest_rate': interestRate,
      'is_liability': isLiability,
      'generates_income': generatesIncome,
      'income_frequency': incomeFrequency,
      'purchase_price': purchasePrice,
      'purchase_date': purchaseDate,
    });
    await loadAllData();
  }

  Future<void> addVehicle(String makeModel, double purchaseCost, {DateTime? insuranceRenewalDate}) async {
    await _apiService.post('/api/assets/vehicles', {
      'make_model': makeModel,
      'purchase_cost': purchaseCost,
      'insurance_renewal_date': insuranceRenewalDate?.toIso8601String().split('T').first,
    });
    await loadAllData();
  }

  Future<void> deleteVehicle(String id) async {
    await _apiService.delete('/api/assets/vehicles/$id');
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

  // ── REPORT ────────────────────────────────────────────────────────────────
  Future<String> downloadReport() async {
    final bytes = await _apiService.downloadFile('/api/report/summary-pdf');
    if (bytes.isEmpty) throw Exception('Empty response');
    final dir = Directory.systemTemp;
    final file = File('${dir.path}/horizon_report.pdf');
    await file.writeAsBytes(bytes);
    return file.path;
  }

  // ── ANALYTICS ─────────────────────────────────────────────────────────────
  Future<void> loadBudgetBreakdown() async {
    try {
      budgetBreakdown = await _apiService.get('/api/analytics/budget-breakdown');
    } catch (e) {
      print('Budget breakdown load error: $e');
    }
    notifyListeners();
  }

  // ── ADVISOR ────────────────────────────────────────────────────────────────
  List<dynamic> nudges = [];
  List<dynamic> simulations = [];
  Map<String, dynamic>? debtOptimizer;
  List<dynamic> subscriptionInsights = [];

  Future<void> loadNudges() async {
    try {
      final data = await _apiService.get('/api/advisor/nudges');
      if (data is List) nudges = data;
    } catch (e) {
      print('Nudges load error: $e');
    }
    notifyListeners();
  }

  Future<void> dismissNudge(int id) async {
    await _apiService.put('/api/advisor/nudges/$id/dismiss', {});
    await loadNudges();
  }

  Future<Map<String, dynamic>> lifecycleSimulate(Map<String, dynamic> params) async {
    final result = await _apiService.post('/api/advisor/lifecycle-simulate', params);
    return result;
  }

  Future<Map<String, dynamic>> runSimulation(Map<String, dynamic> params) async {
    final result = await _apiService.post('/api/advisor/simulate', params);
    await loadSimulations();
    return result;
  }

  Future<void> loadSimulations() async {
    try {
      final data = await _apiService.get('/api/advisor/simulations');
      if (data is List) simulations = data;
    } catch (e) {
      print('Simulations load error: $e');
    }
    notifyListeners();
  }

  Future<void> loadDebtOptimizer() async {
    try {
      debtOptimizer = await _apiService.get('/api/advisor/debt-optimizer');
    } catch (e) {
      print('Debt optimizer load error: $e');
    }
    notifyListeners();
  }

  Future<void> loadSubscriptionInsights() async {
    try {
      final data = await _apiService.get('/api/advisor/subscription-insights');
      if (data is List) subscriptionInsights = data;
    } catch (e) {
      print('Subscription insights load error: $e');
    }
    notifyListeners();
  }

  Future<void> flagSubscription(int billId, String status) async {
    await _apiService.put('/api/advisor/subscription-insights/$billId', {'status': status});
    await loadSubscriptionInsights();
  }

  Future<void> loadAllAdvisorData() async {
    try {
      await Future.wait([
        loadNudges(),
        loadSimulations(),
        loadDebtOptimizer(),
        loadSubscriptionInsights(),
      ]);
    } catch (e) {
      print('Advisor data load error: $e');
    }
  }
}
