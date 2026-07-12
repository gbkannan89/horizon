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

  // Phase 1: Gold Assets
  List<LocalGoldAsset> goldAssets = [];
  // Phase 1: Stock Holdings
  List<LocalStockHolding> stockHoldings = [];
  // Phase 1: PF Assets
  List<LocalPFAsset> pfAssets = [];
  // Phase 1: Lending Records
  List<LocalLendingRecord> lendingRecords = [];
  Map<String, dynamic>? lendingOverview;

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
  Map<String, dynamic>? spendingPatterns;
  List<dynamic> subscriptionCandidates = [];
  List<dynamic> lapsedSubscriptionDetections = [];

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

  Map<String, dynamic>? _householdSummary;

  Future<void> _reloadHousehold() async {
    try {
      final hhData = await _apiService.get('/api/household/summary');
      _householdSummary = hhData;
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

  Future<Map<String, dynamic>> loadHouseholdSummary() async {
    try {
      final data = await _apiService.get('/api/household/summary');
      _householdSummary = data;
      return data;
    } catch (e) {
      return {};
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

  // ── PHASE 1: GOLD ASSETS ──────────────────────────────────────────────────
  Future<void> _reloadGoldAssets() async {
    try {
      final data = await _apiService.get('/api/assets/gold');
      if (data is List) {
        goldAssets = data.map((g) => LocalGoldAsset.fromJson(g as Map<String, dynamic>)).toList();
      }
    } catch (e) { print('Gold assets load error: $e'); }
  }

  // ── PHASE 1: STOCK HOLDINGS ────────────────────────────────────────────────
  Future<void> _reloadStockHoldings() async {
    try {
      final data = await _apiService.get('/api/assets/stocks');
      if (data is List) {
        stockHoldings = data.map((s) => LocalStockHolding.fromJson(s as Map<String, dynamic>)).toList();
      }
    } catch (e) { print('Stock holdings load error: $e'); }
  }

  // ── PHASE 1: PF ASSETS ─────────────────────────────────────────────────────
  Future<void> _reloadPFAssets() async {
    try {
      final data = await _apiService.get('/api/assets/pf');
      if (data is List) {
        pfAssets = data.map((p) => LocalPFAsset.fromJson(p as Map<String, dynamic>)).toList();
      }
    } catch (e) { print('PF assets load error: $e'); }
  }

  // ── PHASE 1: LENDING RECORDS ───────────────────────────────────────────────
  Future<void> _reloadLendingRecords() async {
    try {
      final data = await _apiService.get('/api/lending');
      if (data is List) {
        lendingRecords = data.map((l) => LocalLendingRecord.fromJson(l as Map<String, dynamic>)).toList();
      }
    } catch (e) { print('Lending records load error: $e'); }
  }

  Future<void> _reloadLendingOverview() async {
    try {
      lendingOverview = await _apiService.get('/api/lending/overview');
    } catch (e) { print('Lending overview load error: $e'); }
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
      _reloadGoldAssets(),
      _reloadStockHoldings(),
      _reloadPFAssets(),
    ]);
    _reloadLendingRecords();
    _reloadLendingOverview();

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

  // ── INSIGHTS / ANALYSIS ────────────────────────────────────────────────────
  Map<String, dynamic>? lastAnalysisResult;
  List<dynamic> pendingRecurringSuggestions = [];

  Future<int> uploadStatement(String path, {bool runAnalysis = false}) async {
    String endpoint = '/api/transactions/upload';
    if (runAnalysis) {
      endpoint += '?run_analysis=true';
    }
    final res = await _apiService.uploadFile(endpoint, 'file', path);
    await loadAllData();

    if (runAnalysis && res['analysis'] != null) {
      lastAnalysisResult = res['analysis'];
      pendingRecurringSuggestions = List<dynamic>.from(res['analysis']['recurring_detections'] ?? []);
    }
    notifyListeners();
    return res['inserted'] ?? 0;
  }

  Future<void> confirmRecurringSuggestion(Map<String, dynamic> suggestion, {int dueDay = 1}) async {
    try {
      await _apiService.post('/api/analytics/recurring-detection/confirm', {
        'name': suggestion['name'],
        'amount': suggestion['amount'],
        'frequency': suggestion['frequency'],
        'category': suggestion['category'],
        'bucket': suggestion['bucket'],
        'due_day': dueDay,
        'is_subscription': suggestion['is_subscription'] ?? false,
      });
      pendingRecurringSuggestions.removeWhere((s) => s['name'] == suggestion['name']);
      notifyListeners();
      await _reloadDashboard();
    } catch (e) {
      print('Confirm recurring error: $e');
      rethrow;
    }
  }

  Future<void> loadSpendingPatterns() async {
    try {
      spendingPatterns = await _apiService.get('/api/analytics/spending-patterns');
    } catch (e) {
      print('Load spending patterns error: $e');
    }
    notifyListeners();
  }

  Future<void> detectSubscriptions() async {
    try {
      subscriptionCandidates = await _apiService.get('/api/analytics/subscriptions/detect');
    } catch (e) {
      print('Detect subscriptions error: $e');
    }
    notifyListeners();
  }

  Future<void> detectLapsedSubscriptions() async {
    try {
      lapsedSubscriptionDetections = await _apiService.get('/api/analytics/subscriptions/lapsed');
    } catch (e) {
      print('Detect lapsed subscriptions error: $e');
    }
    notifyListeners();
  }

  Future<void> runFullAnalysis() async {
    try {
      lastAnalysisResult = await _apiService.post('/api/analytics/run-analysis', {});
      pendingRecurringSuggestions = List<dynamic>.from(lastAnalysisResult?['recurring_detections'] ?? []);
      await loadAllAdvisorData();
    } catch (e) {
      print('Run full analysis error: $e');
    }
    notifyListeners();
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

  Future<Map<String, dynamic>> inviteFamilyMember(String name, String email, String relationship) async {
    final result = await _apiService.post('/api/household/invite-by-email', {
      'name': name, 'email': email, 'relationship': relationship,
    });
    await _reloadHousehold();
    return result;
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

  // ── COLLECTIONS ─────────────────────────────────────────────────────────────
  List<dynamic> collections = [];
  bool isLoadingCollections = false;
  String? collectionFilter;

  List<dynamic> get filteredCollections {
    if (collectionFilter == null) return collections;
    return collections.where((c) => c['status'] == collectionFilter).toList();
  }

  void setCollectionFilter(String? filter) {
    collectionFilter = filter;
    notifyListeners();
  }

  Future<void> loadCollections() async {
    isLoadingCollections = true;
    notifyListeners();
    try {
      final data = await _apiService.get('/api/collections');
      if (data is List) collections = data;
    } catch (e) {
      print('Load collections error: $e');
    }
    isLoadingCollections = false;
    notifyListeners();
  }

  Future<Map<String, dynamic>?> getCollectionDetail(int id) async {
    try {
      final data = await _apiService.get('/api/collections/$id');
      final index = collections.indexWhere((c) => c['id'] == id);
      if (index >= 0) collections[index] = data;
      notifyListeners();
      return data;
    } catch (e) {
      print('Get collection detail error: $e');
      return null;
    }
  }

  Future<void> createCollection(Map<String, dynamic> data) async {
    try {
      final created = await _apiService.post('/api/collections', data);
      collections.insert(0, created);
      notifyListeners();
    } catch (e) {
      print('Create collection error: $e');
      rethrow;
    }
  }

  Future<void> updateCollection(int id, {String? label, String? description, String? status}) async {
    try {
      final body = <String, dynamic>{};
      if (label != null) body['label'] = label;
      if (description != null) body['description'] = description;
      if (status != null) body['status'] = status;
      final updated = await _apiService.put('/api/collections/$id', body);
      final index = collections.indexWhere((c) => c['id'] == id);
      if (index >= 0) collections[index] = updated;
      notifyListeners();
    } catch (e) {
      print('Update collection error: $e');
      rethrow;
    }
  }

  Future<void> deleteCollection(int id) async {
    try {
      await _apiService.delete('/api/collections/$id');
      collections.removeWhere((c) => c['id'] == id);
      notifyListeners();
    } catch (e) {
      print('Delete collection error: $e');
      rethrow;
    }
  }

  Future<void> recordPayment(int memberId, double amount) async {
    try {
      final today = DateTime.now().toIso8601String().split('T').first;
      await _apiService.put('/api/collections/members/$memberId', {
        'paid_amount': amount,
        'paid_date': today,
      });
    } catch (e) {
      print('Record payment error: $e');
      rethrow;
    }
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 1: GOLD ASSETS CRUD
  // ═══════════════════════════════════════════════════════════════════════════
  Future<Map<String, dynamic>> addGoldAsset(Map<String, dynamic> body) async {
    final created = await _apiService.post('/api/assets/gold', body);
    goldAssets.add(LocalGoldAsset.fromJson(created));
    notifyListeners();
    _reloadDashboard();
    return created;
  }

  Future<void> updateGoldAsset(int id, Map<String, dynamic> body) async {
    final updated = await _apiService.put('/api/assets/gold/$id', body);
    final idx = goldAssets.indexWhere((g) => g.id == id);
    if (idx >= 0) goldAssets[idx] = LocalGoldAsset.fromJson(updated);
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> deleteGoldAsset(int id) async {
    final idx = goldAssets.indexWhere((g) => g.id == id);
    if (idx == -1) return;
    final removed = goldAssets.removeAt(idx);
    notifyListeners();
    try {
      await _apiService.delete('/api/assets/gold/$id');
      _reloadDashboard();
    } catch (_) {
      goldAssets.insert(idx, removed);
      notifyListeners();
      rethrow;
    }
  }

  Future<void> refreshGoldValues() async {
    final data = await _apiService.post('/api/assets/gold/refresh', {});
    if (data is List) {
      goldAssets = data.map((g) => LocalGoldAsset.fromJson(g as Map<String, dynamic>)).toList();
    }
    notifyListeners();
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 1: STOCK HOLDINGS CRUD
  // ═══════════════════════════════════════════════════════════════════════════
  Future<Map<String, dynamic>> addStockHolding(Map<String, dynamic> body) async {
    final created = await _apiService.post('/api/assets/stocks', body);
    stockHoldings.add(LocalStockHolding.fromJson(created));
    notifyListeners();
    _reloadDashboard();
    return created;
  }

  Future<void> updateStockHolding(int id, Map<String, dynamic> body) async {
    final updated = await _apiService.put('/api/assets/stocks/$id', body);
    final idx = stockHoldings.indexWhere((s) => s.id == id);
    if (idx >= 0) stockHoldings[idx] = LocalStockHolding.fromJson(updated);
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> deleteStockHolding(int id) async {
    final idx = stockHoldings.indexWhere((s) => s.id == id);
    if (idx == -1) return;
    final removed = stockHoldings.removeAt(idx);
    notifyListeners();
    try {
      await _apiService.delete('/api/assets/stocks/$id');
      _reloadDashboard();
    } catch (_) {
      stockHoldings.insert(idx, removed);
      notifyListeners();
      rethrow;
    }
  }

  Future<void> refreshStockValues() async {
    final data = await _apiService.post('/api/assets/stocks/refresh', {});
    if (data is List) {
      stockHoldings = data.map((s) => LocalStockHolding.fromJson(s as Map<String, dynamic>)).toList();
    }
    notifyListeners();
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 1: PF ASSETS CRUD
  // ═══════════════════════════════════════════════════════════════════════════
  Future<Map<String, dynamic>> addPFAsset(Map<String, dynamic> body) async {
    final created = await _apiService.post('/api/assets/pf', body);
    pfAssets.add(LocalPFAsset.fromJson(created));
    notifyListeners();
    _reloadDashboard();
    return created;
  }

  Future<void> updatePFAsset(Map<String, dynamic> body) async {
    final updated = await _apiService.put('/api/assets/pf', body);
    final idx = pfAssets.indexWhere((p) => p.id == updated['id']);
    if (idx >= 0) pfAssets[idx] = LocalPFAsset.fromJson(updated);
    else pfAssets.add(LocalPFAsset.fromJson(updated));
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> deletePFAsset() async {
    if (pfAssets.isEmpty) return;
    pfAssets.clear();
    notifyListeners();
    try {
      await _apiService.delete('/api/assets/pf');
      _reloadDashboard();
    } catch (_) {
      _reloadPFAssets();
      rethrow;
    }
  }

  Future<Map<String, dynamic>?> getPFProjection() async {
    try {
      return await _apiService.get('/api/assets/pf/projection');
    } catch (e) {
      print('PF projection error: $e');
      return null;
    }
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 1: LENDING RECORDS CRUD
  // ═══════════════════════════════════════════════════════════════════════════
  Future<Map<String, dynamic>> addLendingRecord(Map<String, dynamic> body) async {
    final created = await _apiService.post('/api/lending', body);
    lendingRecords.add(LocalLendingRecord.fromJson(created));
    notifyListeners();
    _reloadLendingOverview();
    _reloadDashboard();
    return created;
  }

  Future<void> updateLendingRecord(int id, Map<String, dynamic> body) async {
    final updated = await _apiService.put('/api/lending/$id', body);
    final idx = lendingRecords.indexWhere((l) => l.id == id);
    if (idx >= 0) lendingRecords[idx] = LocalLendingRecord.fromJson(updated);
    notifyListeners();
    _reloadLendingOverview();
    _reloadDashboard();
  }

  Future<void> deleteLendingRecord(int id) async {
    final idx = lendingRecords.indexWhere((l) => l.id == id);
    if (idx == -1) return;
    final removed = lendingRecords.removeAt(idx);
    notifyListeners();
    try {
      await _apiService.delete('/api/lending/$id');
      _reloadLendingOverview();
      _reloadDashboard();
    } catch (_) {
      lendingRecords.insert(idx, removed);
      notifyListeners();
      rethrow;
    }
  }

  Map<String, dynamic>? get lendingSummary => lendingOverview;
}
