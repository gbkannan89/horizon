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

  // ═══════════════════════════════════════════════════════════════════════════
  // TARGETED RELOAD METHODS — no loading state, just data refresh
  // ═══════════════════════════════════════════════════════════════════════════

  Future<void> _reloadDashboard() async {
    try {
      final data = await _apiService.get('/api/v1/dashboard/overview?month=$selectedMonth&year=$selectedYear');
      netWorth = (data['net_worth'] ?? 0).toDouble();
      netWorthChange = (data['net_worth_change'] ?? 0).toDouble();
      netWorthChangePeriod = data['net_worth_change_period'] ?? '';
      finScoreVal = data['fin_score'] ?? 0;
      totalIncomeAgg = (data['total_income'] ?? 0).toDouble();
      totalSpent = (data['total_spent'] ?? 0).toDouble();
      totalLeft = (data['total_left'] ?? 0).toDouble();
      needsSpent = (data['needs_spent'] ?? 0).toDouble();
      wantsSpent = (data['wants_spent'] ?? 0).toDouble();
      savingsSpent = (data['savings_spent'] ?? 0).toDouble();
      needsBudget = (data['needs_budget'] ?? 1).toDouble();
      wantsBudget = (data['wants_budget'] ?? 1).toDouble();
      savingsBudget = (data['savings_budget'] ?? 1).toDouble();
      goals = data['goals'] ?? [];
      recentExpenses = data['recent_expenses'] ?? [];
      upcomingBills = data['upcoming_bills'] ?? [];
      spendingTrend = data['spending_trend'] ?? [];
      notifyListeners();
    } catch (e) {
      print('Dashboard load error: $e');
    }
  }

  Future<void> _reloadAssets() async {
    try {
      final data = await _apiService.get('/api/assets');
      if (data is List) {
        assets = data.map((a) => LocalAsset(
          id: a['id'].toString(), name: a['name'],
          amount: (a['amount'] ?? 0).toDouble(),
          interestRate: (a['interest_rate'] ?? 0).toDouble(),
          isLiability: a['is_liability'] ?? false,
          generatesIncome: a['generates_income'] ?? false,
          purchasePrice: a['purchase_price']?.toDouble(),
          purchaseDate: a['purchase_date'],
        )).toList();
      }
    } catch (e) { print('Assets load error: $e'); }
  }

  Future<void> _reloadPortfolioSummary() async {
    try {
      portfolioSummary = await _apiService.get('/api/analytics/portfolio-summary');
    } catch (e) { print('Portfolio summary load error: $e'); }
  }

  Future<void> _reloadVehicles() async {
    try {
      final data = await _apiService.get('/api/assets/vehicles');
      if (data is List) {
        vehicles = data.map((v) => LocalVehicle(
          id: v['id'].toString(),
          makeModel: v['make_model'],
          purchaseCost: (v['purchase_cost'] ?? 0).toDouble(),
          insuranceRenewalDate: v['insurance_renewal_date'] != null ? DateTime.parse(v['insurance_renewal_date']) : null,
        )).toList();
      }
    } catch (e) { print('Vehicles load error: $e'); }
  }

  Future<void> _reloadBills() async {
    try {
      final data = await _apiService.get('/api/bills');
      if (data is List) {
        bills = data.map((b) => LocalBill(
          id: b['id'].toString(),
          name: b['name'],
          amount: (b['amount'] ?? 0).toDouble(),
          frequency: b['frequency'] ?? 'monthly',
          dueDate: DateTime.now(),
          isEmi: b['is_emi'] ?? false,
          emiTotalMonths: b['emi_total_months'] ?? 0,
          emiMonthsPaid: b['emi_months_paid'] ?? 0,
        )).toList();
      }
    } catch (e) { print('Bills load error: $e'); }
  }

  Future<void> _reloadLiabilities() async {
    try {
      final data = await _apiService.get('/api/liabilities');
      if (data is List) {
        liabilities = data.map((l) => LocalLiability(
          id: l['id'].toString(), name: l['name'],
          amount: (l['outstanding'] ?? 0).toDouble(),
          interestRate: (l['interest_rate'] ?? 0).toDouble()
        )).toList();
      }
    } catch (e) { print('Liabilities load error: $e'); }
  }

  Future<void> _reloadIncomes() async {
    try {
      final data = await _apiService.get('/api/incomes');
      if (data is List) {
        incomes = data.map((i) => LocalIncome(
          id: i['id'].toString(),
          label: i['label'] ?? _typeToLabel(i['type'] ?? 'salary'),
          type: i['type'] ?? 'salary',
          amount: (i['amount'] ?? 0).toDouble(),
          frequency: i['frequency'] ?? 'monthly',
        )).toList();
      }
    } catch (e) { print('Incomes load error: $e'); }
  }

  Future<void> _reloadHousehold() async {
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
      } else {
        householdMembers = [];
      }
    } catch (e) {
      print('Household load error: $e');
      householdMembers = [];
    }
  }

  Future<void> _reloadWishlist() async {
    try {
      final wData = await _apiService.get('/api/discipline/wishlist');
      if (wData is List) {
        wishlist = wData.map((w) => WishlistItem.fromJson(w as Map<String, dynamic>)).toList();
      }
    } catch (e) {
      print('Wishlist load error: $e');
    }
  }

  Future<void> _reloadInsurances() async {
    try {
      final data = await _apiService.get('/api/insurance');
      if (data is List) {
        insurances = data;
      }
    } catch (e) {
      print('Insurance load error: $e');
    }
  }

  Future<void> _reloadScoreHistory() async {
    try {
      final shData = await _apiService.get('/api/score/history');
      if (shData is List) {
        scoreHistory = shData.reversed.toList();
      }
    } catch (e) {
      print('Score history load error: $e');
    }
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // FULL DATA LOAD — for initial load, month navigation, and pull-to-refresh
  // ═══════════════════════════════════════════════════════════════════════════

  Future<void> loadAllData() async {
    _isLoading = true;
    notifyListeners();

    await Future.wait([
      _reloadDashboard(),
      _reloadAssets(),
      _reloadPortfolioSummary(),
      _reloadVehicles(),
      _reloadBills(),
      _reloadLiabilities(),
      _reloadIncomes(),
      _reloadHousehold(),
      _reloadWishlist(),
      _reloadInsurances(),
      _reloadScoreHistory(),
    ]);

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
    final created = await _apiService.post('/api/v1/dashboard/expense', {
      'name': name, 'amount': amount, 'category': category, 'bucket': bucket,
      'date': DateTime.now().toIso8601String().split('T').first, 'icon': 'receipt'
    });
    if (created != null && created['id'] != null) {
      recentExpenses.insert(0, created);
    }
    notifyListeners();
    _reloadDashboard();
  }

  // ── INCOME ─────────────────────────────────────────────────────────────────
  Future<void> addIncome(String label, String type, double amount, String frequency) async {
    final created = await _apiService.post('/api/incomes', {
      'label': label, 'type': type, 'amount': amount, 'frequency': frequency,
    });
    incomes.add(LocalIncome(
      id: created['id'].toString(),
      label: created['label'] ?? label,
      type: created['type'] ?? type,
      amount: (created['amount'] ?? amount).toDouble(),
      frequency: created['frequency'] ?? frequency,
    ));
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> deleteIncome(String id) async {
    final index = incomes.indexWhere((i) => i.id == id);
    if (index == -1) return;
    final removed = incomes.removeAt(index);
    notifyListeners();
    try {
      await _apiService.delete('/api/incomes/$id');
      _reloadDashboard();
    } catch (_) {
      incomes.insert(index, removed);
      notifyListeners();
      rethrow;
    }
  }

  // ── RECURRING BILLS ────────────────────────────────────────────────────────
  Future<void> addRecurringBill(String name, double amount, String category, String bucket, {
    String frequency = 'monthly',
    bool isEmi = false,
    int? emiTotalMonths,
    DateTime? startDate,
  }) async {
    final created = await _apiService.post('/api/bills', {
      'name': name, 
      'amount': amount, 
      'category': category, 
      'bucket': bucket,
      'frequency': frequency,
      'is_emi': isEmi,
      'emi_total_months': emiTotalMonths,
      'start_date': startDate?.toIso8601String().split('T').first,
    });
    bills.add(LocalBill(
      id: created['id'].toString(),
      name: created['name'] ?? name,
      amount: (created['amount'] ?? amount).toDouble(),
      frequency: created['frequency'] ?? frequency,
      dueDate: DateTime.now(),
      isEmi: created['is_emi'] ?? isEmi,
      emiTotalMonths: created['emi_total_months'] ?? emiTotalMonths ?? 0,
      emiMonthsPaid: created['emi_months_paid'] ?? 0,
    ));
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> deleteRecurringBill(String id) async {
    final index = bills.indexWhere((b) => b.id == id);
    if (index == -1) return;
    final removed = bills.removeAt(index);
    notifyListeners();
    try {
      await _apiService.delete('/api/bills/$id');
      _reloadDashboard();
    } catch (_) {
      bills.insert(index, removed);
      notifyListeners();
      rethrow;
    }
  }

  // ── INSURANCE ─────────────────────────────────────────────────────────────
  Future<void> addInsurance(Map<String, dynamic> data) async {
    final created = await _apiService.post('/insurance/', data);
    insurances.add(created);
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> deleteInsurance(String id) async {
    final index = insurances.indexWhere((ins) => ins['id'].toString() == id);
    if (index == -1) return;
    final removed = insurances.removeAt(index);
    notifyListeners();
    try {
      await _apiService.delete('/insurance/$id');
      _reloadDashboard();
    } catch (_) {
      insurances.insert(index, removed);
      notifyListeners();
      rethrow;
    }
  }

  // ── TRANSACTIONS ───────────────────────────────────────────────────────────
  Future<int> uploadStatement(String path) async {
    final res = await _apiService.uploadFile('/api/transactions/upload', 'file', path);
    await loadAllData();
    return res['inserted'] ?? 0;
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
    final created = await _apiService.post('/api/assets', {
      'name': name, 'type': type, 'amount': amount,
      'interest_rate': interestRate,
      'is_liability': isLiability,
      'generates_income': generatesIncome,
      'income_frequency': incomeFrequency,
      'purchase_price': purchasePrice,
      'purchase_date': purchaseDate,
    });
    assets.add(LocalAsset(
      id: created['id'].toString(),
      name: created['name'] ?? name,
      amount: (created['amount'] ?? amount).toDouble(),
      interestRate: (created['interest_rate'] ?? interestRate).toDouble(),
      isLiability: created['is_liability'] ?? isLiability,
      generatesIncome: created['generates_income'] ?? generatesIncome,
      purchasePrice: created['purchase_price']?.toDouble(),
      purchaseDate: created['purchase_date'],
    ));
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> addVehicle(String makeModel, double purchaseCost, {DateTime? insuranceRenewalDate}) async {
    final created = await _apiService.post('/api/assets/vehicles', {
      'make_model': makeModel,
      'purchase_cost': purchaseCost,
      'insurance_renewal_date': insuranceRenewalDate?.toIso8601String().split('T').first,
    });
    vehicles.add(LocalVehicle(
      id: created['id'].toString(),
      makeModel: created['make_model'] ?? makeModel,
      purchaseCost: (created['purchase_cost'] ?? purchaseCost).toDouble(),
      insuranceRenewalDate: created['insurance_renewal_date'] != null ? DateTime.parse(created['insurance_renewal_date']) : null,
    ));
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> deleteVehicle(String id) async {
    final index = vehicles.indexWhere((v) => v.id == id);
    if (index == -1) return;
    final removed = vehicles.removeAt(index);
    notifyListeners();
    try {
      await _apiService.delete('/api/assets/vehicles/$id');
      _reloadDashboard();
    } catch (_) {
      vehicles.insert(index, removed);
      notifyListeners();
      rethrow;
    }
  }

  Future<void> deleteAsset(String id) async {
    final index = assets.indexWhere((a) => a.id == id);
    if (index == -1) return;
    final removed = assets.removeAt(index);
    notifyListeners();
    try {
      await _apiService.delete('/api/assets/$id');
      _reloadDashboard();
    } catch (_) {
      assets.insert(index, removed);
      notifyListeners();
      rethrow;
    }
  }

  // ── LIABILITIES ────────────────────────────────────────────────────────────
  Future<void> addLiability(String name, String type, double outstanding, double emi, double interestRate) async {
    final created = await _apiService.post('/api/liabilities', {
      'name': name, 'type': type, 'outstanding': outstanding,
      'emi': emi, 'interest_rate': interestRate,
    });
    liabilities.add(LocalLiability(
      id: created['id'].toString(),
      name: created['name'] ?? name,
      amount: (created['outstanding'] ?? outstanding).toDouble(),
      interestRate: (created['interest_rate'] ?? interestRate).toDouble(),
    ));
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> deleteLiability(String id) async {
    final index = liabilities.indexWhere((l) => l.id == id);
    if (index == -1) return;
    final removed = liabilities.removeAt(index);
    notifyListeners();
    try {
      await _apiService.delete('/api/liabilities/$id');
      _reloadDashboard();
    } catch (_) {
      liabilities.insert(index, removed);
      notifyListeners();
      rethrow;
    }
  }

  // ── HOUSEHOLD MEMBERS ──────────────────────────────────────────────────────
  Future<void> addContributingMember(String name, double monthlyIncome, double contribution, String relationship) async {
    final created = await _apiService.post('/api/household/contributing-members', {
      'name': name, 'monthly_income': monthlyIncome,
      'contribution_to_household': contribution, 'relationship': relationship
    });
    householdMembers.add(LocalHouseholdMember(
      id: created['id'].toString(),
      name: created['name'] ?? name,
      monthlyIncome: (created['monthlyIncome'] ?? monthlyIncome).toDouble(),
      contribution: (created['contributionToHousehold'] ?? contribution).toDouble(),
      relationship: created['relationship'] ?? relationship,
    ));
    notifyListeners();
    _reloadDashboard();
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
