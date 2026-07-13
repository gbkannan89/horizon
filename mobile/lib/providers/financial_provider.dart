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

  // Phase 5: Family & Profile Tracking
  List<LocalFamilyMember> familyMembers = [];
  double recurringFamilyCosts = 0.0;
  Map<String, dynamic>? recurringCostsBreakdown;

  // Phase 6: Budget Revamp
  LocalBudgetPlan? budgetPlan;
  LocalBudgetComparison? budgetComparison;

  // Phase 2: Income & Salary Deep Dive
  List<LocalSalaryDetail> salaryDetails = [];
  LocalSalaryGrowth? salaryGrowth;

  // Phase 3: Vehicle Enhanced Tracking
  List<LocalVehicleService> vehicleServices = [];
  List<LocalVehicleFuel> vehicleFuels = [];
  LocalVehicleLoan? vehicleLoan;

  // Phase 4: Electronics / Device Tracking
  List<LocalElectronic> electronics = [];
  List<LocalElectronicService> electronicServices = [];
  LocalElectronicEmi? electronicEmi;

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
  List<dynamic> investmentReturns = [];
  Map<String, dynamic>? taxEstimate;
  Map<String, dynamic>? debtDashboard;
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
      final data = await _apiService.get(
        '/api/v1/dashboard/overview?month=$selectedMonth&year=$selectedYear',
      );
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
        assets = data
            .map(
              (a) => LocalAsset(
                id: a['id'].toString(),
                name: a['name'],
                amount: (a['amount'] ?? 0).toDouble(),
                interestRate: (a['interest_rate'] ?? 0).toDouble(),
                isLiability: a['is_liability'] ?? false,
                generatesIncome: a['generates_income'] ?? false,
                purchasePrice: a['purchase_price']?.toDouble(),
                purchaseDate: a['purchase_date'],
              ),
            )
            .toList();
      }
    } catch (e) {
      print('Assets load error: $e');
    }
  }

  Future<void> _reloadPortfolioSummary() async {
    try {
      portfolioSummary = await _apiService.get(
        '/api/analytics/portfolio-summary',
      );
    } catch (e) {
      print('Portfolio summary load error: $e');
    }
  }

  Future<void> _reloadVehicles() async {
    try {
      final data = await _apiService.get('/api/assets/vehicles');
      if (data is List) {
        vehicles = data
            .map((v) => LocalVehicle.fromJson(v as Map<String, dynamic>))
            .toList();
      }
    } catch (e) {
      print('Vehicles load error: $e');
    }
  }

  Future<void> reloadVehicleServiceRecords(String vehicleId) async {
    try {
      final intId = int.parse(vehicleId);
      final data = await _apiService.get('/api/assets/vehicles/$intId/service');
      if (data is List) {
        vehicleServices = data
            .map((x) => LocalVehicleService.fromJson(x as Map<String, dynamic>))
            .toList();
      }
      notifyListeners();
    } catch (e) {
      print('Service records load error: $e');
    }
  }

  Future<void> reloadVehicleFuelRecords(String vehicleId) async {
    try {
      final intId = int.parse(vehicleId);
      final data = await _apiService.get('/api/assets/vehicles/$intId/fuel');
      if (data is List) {
        vehicleFuels = data
            .map((x) => LocalVehicleFuel.fromJson(x as Map<String, dynamic>))
            .toList();
      }
      notifyListeners();
    } catch (e) {
      print('Fuel records load error: $e');
    }
  }

  Future<void> reloadVehicleLoan(String vehicleId) async {
    try {
      final intId = int.parse(vehicleId);
      final data = await _apiService.get('/api/assets/vehicles/$intId/loan');
      if (data != null) {
        vehicleLoan = LocalVehicleLoan.fromJson(data);
      } else {
        vehicleLoan = null;
      }
      notifyListeners();
    } catch (e) {
      print('Vehicle loan load error: $e');
    }
  }

  Future<void> reloadElectronics() async {
    try {
      final data = await _apiService.get('/api/electronics');
      if (data is List) {
        electronics = data
            .map((e) => LocalElectronic.fromJson(e as Map<String, dynamic>))
            .toList();
      }
      notifyListeners();
    } catch (e) {
      print('Electronics load error: $e');
    }
  }

  Future<void> reloadElectronicServices(String deviceId) async {
    try {
      final intId = int.parse(deviceId);
      final data = await _apiService.get('/api/electronics/$intId/service');
      if (data is List) {
        electronicServices = data
            .map(
              (x) => LocalElectronicService.fromJson(x as Map<String, dynamic>),
            )
            .toList();
      }
      notifyListeners();
    } catch (e) {
      print('Electronic service load error: $e');
    }
  }

  Future<void> reloadElectronicEmi(String deviceId) async {
    try {
      final intId = int.parse(deviceId);
      final data = await _apiService.get('/api/electronics/$intId/emi');
      if (data != null) {
        electronicEmi = LocalElectronicEmi.fromJson(data);
      } else {
        electronicEmi = null;
      }
      notifyListeners();
    } catch (e) {
      print('Electronic EMI load error: $e');
    }
  }

  Future<void> _reloadBills() async {
    try {
      final data = await _apiService.get('/api/bills');
      if (data is List) {
        bills = data
            .map(
              (b) => LocalBill(
                id: b['id'].toString(),
                name: b['name'],
                amount: (b['amount'] ?? 0).toDouble(),
                frequency: b['frequency'] ?? 'monthly',
                dueDate: DateTime.now(),
                isEmi: b['is_emi'] ?? false,
                emiTotalMonths: b['emi_total_months'] ?? 0,
                emiMonthsPaid: b['emi_months_paid'] ?? 0,
              ),
            )
            .toList();
      }
    } catch (e) {
      print('Bills load error: $e');
    }
  }

  Future<void> _reloadLiabilities() async {
    try {
      final data = await _apiService.get('/api/liabilities');
      if (data is List) {
        liabilities = data
            .map(
              (l) => LocalLiability(
                id: l['id'].toString(),
                name: l['name'],
                amount: (l['outstanding'] ?? 0).toDouble(),
                interestRate: (l['interest_rate'] ?? 0).toDouble(),
              ),
            )
            .toList();
      }
    } catch (e) {
      print('Liabilities load error: $e');
    }
  }

  Future<void> _reloadIncomes() async {
    try {
      final data = await _apiService.get('/api/incomes');
      if (data is List) {
        incomes = data
            .map(
              (i) => LocalIncome(
                id: i['id'].toString(),
                label: i['label'] ?? _typeToLabel(i['type'] ?? 'salary'),
                type: i['type'] ?? 'salary',
                amount: (i['amount'] ?? 0).toDouble(),
                frequency: i['frequency'] ?? 'monthly',
                companyName: i['company_name'],
                notes: i['notes'],
              ),
            )
            .toList();
      }
    } catch (e) {
      print('Incomes load error: $e');
    }
  }

  Map<String, dynamic>? _householdSummary;

  Future<void> _reloadHousehold() async {
    try {
      final hhData = await _apiService.get('/api/household/summary');
      _householdSummary = hhData;
      if (hhData['contributingMembers'] != null) {
        final members = hhData['contributingMembers'] as List;
        householdMembers = members
            .map(
              (m) => LocalHouseholdMember(
                id: m['id'].toString(),
                name: m['name'],
                monthlyIncome: (m['monthlyIncome'] ?? 0).toDouble(),
                contribution: (m['contributionToHousehold'] ?? 0).toDouble(),
                relationship: m['relationship'] ?? 'Member',
              ),
            )
            .toList();
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
        wishlist = wData
            .map((w) => WishlistItem.fromJson(w as Map<String, dynamic>))
            .toList();
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

  Future<void> reloadInvestmentReturns() async {
    try {
      final data = await _apiService.get('/api/analytics/returns');
      if (data is List) {
        investmentReturns = data;
        notifyListeners();
      }
    } catch (e) {
      print('Investment returns load error: $e');
    }
  }

  Future<void> reloadTaxEstimate() async {
    try {
      final data = await _apiService.get('/api/advisor/tax-estimate');
      taxEstimate = data;
      notifyListeners();
    } catch (e) {
      print('Tax estimate load error: $e');
    }
  }

  Future<void> reloadDebtDashboard() async {
    try {
      final data = await _apiService.get('/api/advisor/debt-dashboard');
      debtDashboard = data;
      notifyListeners();
    } catch (e) {
      print('Debt dashboard load error: $e');
    }
  }

  // ── PHASE 1: GOLD ASSETS ──────────────────────────────────────────────────
  Future<void> _reloadGoldAssets() async {
    try {
      final data = await _apiService.get('/api/assets/gold');
      if (data is List) {
        goldAssets = data
            .map((g) => LocalGoldAsset.fromJson(g as Map<String, dynamic>))
            .toList();
      }
    } catch (e) {
      print('Gold assets load error: $e');
    }
  }

  // ── PHASE 1: STOCK HOLDINGS ────────────────────────────────────────────────
  Future<void> _reloadStockHoldings() async {
    try {
      final data = await _apiService.get('/api/assets/stocks');
      if (data is List) {
        stockHoldings = data
            .map((s) => LocalStockHolding.fromJson(s as Map<String, dynamic>))
            .toList();
      }
    } catch (e) {
      print('Stock holdings load error: $e');
    }
  }

  // ── PHASE 1: PF ASSETS ─────────────────────────────────────────────────────
  Future<void> _reloadPFAssets() async {
    try {
      final data = await _apiService.get('/api/assets/pf');
      if (data is List) {
        pfAssets = data
            .map((p) => LocalPFAsset.fromJson(p as Map<String, dynamic>))
            .toList();
      }
    } catch (e) {
      print('PF assets load error: $e');
    }
  }

  // ── PHASE 1: LENDING RECORDS ───────────────────────────────────────────────
  Future<void> _reloadLendingRecords() async {
    try {
      final data = await _apiService.get('/api/lending');
      if (data is List) {
        lendingRecords = data
            .map((l) => LocalLendingRecord.fromJson(l as Map<String, dynamic>))
            .toList();
      }
    } catch (e) {
      print('Lending records load error: $e');
    }
  }

  Future<void> _reloadLendingOverview() async {
    try {
      lendingOverview = await _apiService.get('/api/lending/overview');
    } catch (e) {
      print('Lending overview load error: $e');
    }
  }

  Future<void> _reloadFamilyMembers() async {
    try {
      final data = await _apiService.get('/api/profile/family');
      if (data is List) {
        familyMembers = data
            .map((f) => LocalFamilyMember.fromJson(f as Map<String, dynamic>))
            .toList();
      }
    } catch (e) {
      print('Family members load error: $e');
    }
  }

  Future<void> _reloadRecurringFamilyCosts() async {
    try {
      final data = await _apiService.get('/api/profile/family/recurring-costs');
      recurringFamilyCosts = (data['total_monthly'] ?? 0).toDouble();
      recurringCostsBreakdown = data;
    } catch (e) {
      print('Recurring costs load error: $e');
    }
  }

  Future<void> _reloadBudgetPlan() async {
    try {
      final data = await _apiService.get(
        '/api/budget/plan?month=$selectedMonth&year=$selectedYear',
      );
      budgetPlan = LocalBudgetPlan.fromJson(data);
    } catch (e) {
      print('Budget plan load error: $e');
    }
  }

  Future<void> _reloadBudgetComparison() async {
    try {
      final data = await _apiService.get(
        '/api/budget/compare?month=$selectedMonth&year=$selectedYear',
      );
      budgetComparison = LocalBudgetComparison.fromJson(data);
    } catch (e) {
      print('Budget comparison load error: $e');
    }
  }

  Future<void> reloadSalaryDetails(String incomeId) async {
    try {
      final intId = int.parse(incomeId);
      final data = await _apiService.get('/api/incomes/$intId/salary');
      if (data is List) {
        salaryDetails = data
            .map((x) => LocalSalaryDetail.fromJson(x as Map<String, dynamic>))
            .toList();
      }
      notifyListeners();
    } catch (e) {
      print('Salary details load error: $e');
    }
  }

  Future<void> loadSalaryGrowth() async {
    try {
      final data = await _apiService.get('/api/incomes/salary-growth');
      salaryGrowth = LocalSalaryGrowth.fromJson(data);
      notifyListeners();
    } catch (e) {
      print('Salary growth load error: $e');
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
      _reloadGoldAssets(),
      _reloadStockHoldings(),
      _reloadPFAssets(),
      _reloadFamilyMembers(),
      _reloadRecurringFamilyCosts(),
      _reloadBudgetPlan(),
      _reloadBudgetComparison(),
      loadSalaryGrowth(),
      reloadElectronics(),
    ]);
    _reloadLendingRecords();
    _reloadLendingOverview();

    loadBudgetBreakdown();
    loadNudges();
    reloadInvestmentReturns();
    reloadTaxEstimate();
    reloadDebtDashboard();

    _isLoading = false;
    notifyListeners();
  }

  static String _typeToLabel(String type) {
    switch (type) {
      case 'salary':
        return 'Salary';
      case 'business':
        return 'Business';
      case 'passive':
        return 'Passive';
      default:
        return type[0].toUpperCase() + type.substring(1);
    }
  }

  // ── EXPENSE ────────────────────────────────────────────────────────────────
  Future<void> addExpense(
    String name,
    double amount,
    String category,
    String bucket,
  ) async {
    final created = await _apiService.post('/api/v1/dashboard/expense', {
      'name': name,
      'amount': amount,
      'category': category,
      'bucket': bucket,
      'date': DateTime.now().toIso8601String().split('T').first,
      'icon': 'receipt',
    });
    if (created != null && created['id'] != null) {
      recentExpenses.insert(0, created);
    }
    notifyListeners();
    _reloadDashboard();
  }

  // ── INCOME ─────────────────────────────────────────────────────────────────
  Future<void> addIncome(
    String label,
    String type,
    double amount,
    String frequency, {
    String? companyName,
    String? notes,
  }) async {
    final created = await _apiService.post('/api/incomes', {
      'label': label,
      'type': type,
      'amount': amount,
      'frequency': frequency,
      'company_name': companyName,
      'notes': notes,
    });
    incomes.add(
      LocalIncome(
        id: created['id'].toString(),
        label: created['label'] ?? label,
        type: created['type'] ?? type,
        amount: (created['amount'] ?? amount).toDouble(),
        frequency: created['frequency'] ?? frequency,
        companyName: created['company_name'] ?? companyName,
        notes: created['notes'] ?? notes,
      ),
    );
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
  Future<void> addRecurringBill(
    String name,
    double amount,
    String category,
    String bucket, {
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
    bills.add(
      LocalBill(
        id: created['id'].toString(),
        name: created['name'] ?? name,
        amount: (created['amount'] ?? amount).toDouble(),
        frequency: created['frequency'] ?? frequency,
        dueDate: DateTime.now(),
        isEmi: created['is_emi'] ?? isEmi,
        emiTotalMonths: created['emi_total_months'] ?? emiTotalMonths ?? 0,
        emiMonthsPaid: created['emi_months_paid'] ?? 0,
      ),
    );
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
      pendingRecurringSuggestions = List<dynamic>.from(
        res['analysis']['recurring_detections'] ?? [],
      );
    }
    notifyListeners();
    return res['inserted'] ?? 0;
  }

  Future<void> confirmRecurringSuggestion(
    Map<String, dynamic> suggestion, {
    int dueDay = 1,
  }) async {
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
      pendingRecurringSuggestions.removeWhere(
        (s) => s['name'] == suggestion['name'],
      );
      notifyListeners();
      await _reloadDashboard();
    } catch (e) {
      print('Confirm recurring error: $e');
      rethrow;
    }
  }

  Future<void> loadSpendingPatterns() async {
    try {
      spendingPatterns = await _apiService.get(
        '/api/analytics/spending-patterns',
      );
    } catch (e) {
      print('Load spending patterns error: $e');
    }
    notifyListeners();
  }

  Future<void> detectSubscriptions() async {
    try {
      subscriptionCandidates = await _apiService.get(
        '/api/analytics/subscriptions/detect',
      );
    } catch (e) {
      print('Detect subscriptions error: $e');
    }
    notifyListeners();
  }

  Future<void> detectLapsedSubscriptions() async {
    try {
      lapsedSubscriptionDetections = await _apiService.get(
        '/api/analytics/subscriptions/lapsed',
      );
    } catch (e) {
      print('Detect lapsed subscriptions error: $e');
    }
    notifyListeners();
  }

  Future<void> runFullAnalysis() async {
    try {
      lastAnalysisResult = await _apiService.post(
        '/api/analytics/run-analysis',
        {},
      );
      pendingRecurringSuggestions = List<dynamic>.from(
        lastAnalysisResult?['recurring_detections'] ?? [],
      );
      await loadAllAdvisorData();
    } catch (e) {
      print('Run full analysis error: $e');
    }
    notifyListeners();
  }

  // ── ASSETS ─────────────────────────────────────────────────────────────────
  Future<void> addAsset(
    String name,
    String type,
    double amount, {
    double interestRate = 0.0,
    bool isLiability = false,
    bool isEmergency = false,
    bool generatesIncome = false,
    String? incomeFrequency,
    double? purchasePrice,
    String? purchaseDate,
    int? yearsOfDeposit,
  }) async {
    final created = await _apiService.post('/api/assets', {
      'name': name,
      'type': type,
      'amount': amount,
      'interest_rate': interestRate,
      'is_liability': isLiability,
      'is_emergency': isEmergency,
      'generates_income': generatesIncome,
      'income_frequency': incomeFrequency,
      'purchase_price': purchasePrice,
      'purchase_date': purchaseDate,
      'years_of_deposit': yearsOfDeposit,
    });
    assets.add(
      LocalAsset(
        id: created['id'].toString(),
        name: created['name'] ?? name,
        amount: (created['amount'] ?? amount).toDouble(),
        interestRate: (created['interest_rate'] ?? interestRate).toDouble(),
        isLiability: created['is_liability'] ?? isLiability,
        generatesIncome: created['generates_income'] ?? generatesIncome,
        purchasePrice: created['purchase_price']?.toDouble(),
        purchaseDate: created['purchase_date'],
      ),
    );
    notifyListeners();
    _reloadDashboard();
  }

  Future<void> addVehicle(
    String makeModel,
    double purchaseCost, {
    DateTime? insuranceRenewalDate,
    int? modelYear,
    int? purchaseYear,
    String? fuelType,
    double? mileageKmpl,
    double? insuranceIdv,
    double? insuranceRenewalAmount,
    String? registrationNumber,
  }) async {
    final created = await _apiService.post('/api/assets/vehicles', {
      'make_model': makeModel,
      'purchase_cost': purchaseCost,
      'insurance_renewal_date': insuranceRenewalDate
          ?.toIso8601String()
          .split('T')
          .first,
      'model_year': modelYear,
      'purchase_year': purchaseYear,
      'fuel_type': fuelType,
      'mileage_kmpl': mileageKmpl,
      'insurance_idv': insuranceIdv,
      'insurance_renewal_amount': insuranceRenewalAmount,
      'registration_number': registrationNumber,
    });
    vehicles.add(LocalVehicle.fromJson(created as Map<String, dynamic>));
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

  Future<void> updateAsset(int id, Map<String, dynamic> body) async {
    final updated = await _apiService.put('/api/assets/$id', body);
    final idx = assets.indexWhere((a) => a.id == id.toString());
    if (idx >= 0) {
      assets[idx] = LocalAsset(
        id: updated['id'].toString(),
        name: updated['name'] ?? assets[idx].name,
        amount: (updated['amount'] ?? assets[idx].amount).toDouble(),
        interestRate: (updated['interest_rate'] ?? assets[idx].interestRate)
            .toDouble(),
        isLiability: updated['is_liability'] ?? false,
        generatesIncome: updated['generates_income'] ?? false,
        purchasePrice: updated['purchase_price']?.toDouble(),
        purchaseDate: updated['purchase_date'],
      );
    }
    notifyListeners();
    _reloadDashboard();
  }

  // ── LIABILITIES ────────────────────────────────────────────────────────────
  Future<void> addLiability(
    String name,
    String type,
    double outstanding,
    double emi,
    double interestRate,
  ) async {
    final created = await _apiService.post('/api/liabilities', {
      'name': name,
      'type': type,
      'outstanding': outstanding,
      'emi': emi,
      'interest_rate': interestRate,
    });
    liabilities.add(
      LocalLiability(
        id: created['id'].toString(),
        name: created['name'] ?? name,
        amount: (created['outstanding'] ?? outstanding).toDouble(),
        interestRate: (created['interest_rate'] ?? interestRate).toDouble(),
      ),
    );
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
  Future<void> addContributingMember(
    String name,
    double monthlyIncome,
    double contribution,
    String relationship,
  ) async {
    final created = await _apiService
        .post('/api/household/contributing-members', {
          'name': name,
          'monthly_income': monthlyIncome,
          'contribution_to_household': contribution,
          'relationship': relationship,
        });
    householdMembers.add(
      LocalHouseholdMember(
        id: created['id'].toString(),
        name: created['name'] ?? name,
        monthlyIncome: (created['monthlyIncome'] ?? monthlyIncome).toDouble(),
        contribution: (created['contributionToHousehold'] ?? contribution)
            .toDouble(),
        relationship: created['relationship'] ?? relationship,
      ),
    );
    notifyListeners();
    _reloadDashboard();
  }

  Future<Map<String, dynamic>> inviteFamilyMember(
    String name,
    String email,
    String relationship,
  ) async {
    final result = await _apiService.post('/api/household/invite-by-email', {
      'name': name,
      'email': email,
      'relationship': relationship,
    });
    await _reloadHousehold();
    return result;
  }

  Future<Map<String, dynamic>> joinHousehold(String inviteCode) async {
    final result = await _apiService.post('/api/household/join', {
      'invite_code': inviteCode,
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
      budgetBreakdown = await _apiService.get(
        '/api/analytics/budget-breakdown',
      );
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

  Future<Map<String, dynamic>> lifecycleSimulate(
    Map<String, dynamic> params,
  ) async {
    final result = await _apiService.post(
      '/api/advisor/lifecycle-simulate',
      params,
    );
    return result;
  }

  Future<Map<String, dynamic>> runSimulation(
    Map<String, dynamic> params,
  ) async {
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
    await _apiService.put('/api/advisor/subscription-insights/$billId', {
      'status': status,
    });
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

  Future<void> updateCollection(
    int id, {
    String? label,
    String? description,
    String? status,
  }) async {
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
      goldAssets = data
          .map((g) => LocalGoldAsset.fromJson(g as Map<String, dynamic>))
          .toList();
    }
    notifyListeners();
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 1: STOCK HOLDINGS CRUD
  // ═══════════════════════════════════════════════════════════════════════════
  Future<Map<String, dynamic>> addStockHolding(
    Map<String, dynamic> body,
  ) async {
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
      stockHoldings = data
          .map((s) => LocalStockHolding.fromJson(s as Map<String, dynamic>))
          .toList();
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
    if (idx >= 0) {
      pfAssets[idx] = LocalPFAsset.fromJson(updated);
    } else {
      pfAssets.add(LocalPFAsset.fromJson(updated));
    }
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
  Future<Map<String, dynamic>> addLendingRecord(
    Map<String, dynamic> body,
  ) async {
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

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 5: FAMILY CRUD OPERATIONS
  // ═══════════════════════════════════════════════════════════════════════════
  Future<void> addFamilyMember(Map<String, dynamic> body) async {
    await _apiService.post('/api/profile/family', body);
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> updateFamilyMember(int id, Map<String, dynamic> body) async {
    await _apiService.put('/api/profile/family/$id', body);
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> deleteFamilyMember(int id) async {
    await _apiService.delete('/api/profile/family/$id');
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  // Schooling
  Future<void> addSchooling(int memberId, Map<String, dynamic> body) async {
    await _apiService.post('/api/profile/family/$memberId/schooling', body);
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> updateSchooling(
    int memberId,
    int schoolingId,
    Map<String, dynamic> body,
  ) async {
    await _apiService.put(
      '/api/profile/family/$memberId/schooling/$schoolingId',
      body,
    );
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> deleteSchooling(int memberId, int schoolingId) async {
    await _apiService.delete(
      '/api/profile/family/$memberId/schooling/$schoolingId',
    );
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<dynamic> getSchoolingPayments(int memberId, int schoolingId) async {
    try {
      return await _apiService.get(
        '/api/profile/family/$memberId/schooling/$schoolingId/payments',
      );
    } catch (e) {
      return [];
    }
  }

  Future<void> deleteSchoolingPayment(
    int memberId,
    int schoolingId,
    int paymentId,
  ) async {
    await _apiService.delete(
      '/api/profile/family/$memberId/schooling/$schoolingId/payments/$paymentId',
    );
    await _reloadFamilyMembers();
    notifyListeners();
  }

  Future<void> recordSchoolingPayment(
    int memberId,
    int schoolingId,
    Map<String, dynamic> body,
  ) async {
    await _apiService.post(
      '/api/profile/family/$memberId/schooling/$schoolingId/payments',
      body,
    );
    await _reloadFamilyMembers();
    notifyListeners();
  }

  // Checkups
  Future<void> addCheckup(int memberId, Map<String, dynamic> body) async {
    await _apiService.post('/api/profile/family/$memberId/checkups', body);
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> updateCheckup(
    int memberId,
    int checkupId,
    Map<String, dynamic> body,
  ) async {
    await _apiService.put(
      '/api/profile/family/$memberId/checkups/$checkupId',
      body,
    );
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> deleteCheckup(int memberId, int checkupId) async {
    await _apiService.delete(
      '/api/profile/family/$memberId/checkups/$checkupId',
    );
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  // Medicines
  Future<void> addMedicine(int memberId, Map<String, dynamic> body) async {
    await _apiService.post('/api/profile/family/$memberId/medicines', body);
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> updateMedicine(
    int memberId,
    int medicineId,
    Map<String, dynamic> body,
  ) async {
    await _apiService.put(
      '/api/profile/family/$memberId/medicines/$medicineId',
      body,
    );
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> deleteMedicine(int memberId, int medicineId) async {
    await _apiService.delete(
      '/api/profile/family/$memberId/medicines/$medicineId',
    );
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  // Vaccinations
  Future<void> addVaccination(int memberId, Map<String, dynamic> body) async {
    await _apiService.post('/api/profile/family/$memberId/vaccinations', body);
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> updateVaccination(
    int memberId,
    int vaccinationId,
    Map<String, dynamic> body,
  ) async {
    await _apiService.put(
      '/api/profile/family/$memberId/vaccinations/$vaccinationId',
      body,
    );
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  Future<void> deleteVaccination(int memberId, int vaccinationId) async {
    await _apiService.delete(
      '/api/profile/family/$memberId/vaccinations/$vaccinationId',
    );
    await _reloadFamilyMembers();
    await _reloadRecurringFamilyCosts();
    notifyListeners();
  }

  // Earnings
  Future<void> addEarnings(int memberId, Map<String, dynamic> body) async {
    await _apiService.post('/api/profile/family/$memberId/earnings', body);
    await _reloadFamilyMembers();
    notifyListeners();
  }

  Future<void> updateEarnings(
    int memberId,
    int earningsId,
    Map<String, dynamic> body,
  ) async {
    await _apiService.put(
      '/api/profile/family/$memberId/earnings/$earningsId',
      body,
    );
    await _reloadFamilyMembers();
    notifyListeners();
  }

  Future<void> deleteEarnings(int memberId, int earningsId) async {
    await _apiService.delete(
      '/api/profile/family/$memberId/earnings/$earningsId',
    );
    await _reloadFamilyMembers();
    notifyListeners();
  }

  // Insurance Link
  Future<void> linkInsurance(int memberId, Map<String, dynamic> body) async {
    await _apiService.post('/api/profile/family/$memberId/insurance', body);
    await _reloadFamilyMembers();
    notifyListeners();
  }

  Future<void> unlinkInsurance(int memberId, int insuranceId) async {
    await _apiService.delete(
      '/api/profile/family/$memberId/insurance/$insuranceId',
    );
    await _reloadFamilyMembers();
    notifyListeners();
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 6: BUDGET PLAN ACTIONS
  // ═══════════════════════════════════════════════════════════════════════════
  Future<void> autoPopulateBudget() async {
    await _apiService.post(
      '/api/budget/plan/auto-populate?month=$selectedMonth&year=$selectedYear',
      {},
    );
    await _reloadBudgetPlan();
    await _reloadBudgetComparison();
    notifyListeners();
  }

  Future<void> addBudgetItem(Map<String, dynamic> body) async {
    await _apiService.post(
      '/api/budget/items?month=$selectedMonth&year=$selectedYear',
      body,
    );
    await _reloadBudgetPlan();
    await _reloadBudgetComparison();
    notifyListeners();
  }

  Future<void> updateBudgetItem(int id, Map<String, dynamic> body) async {
    await _apiService.put('/api/budget/items/$id', body);
    await _reloadBudgetPlan();
    await _reloadBudgetComparison();
    notifyListeners();
  }

  Future<void> deleteBudgetItem(int id) async {
    await _apiService.delete('/api/budget/items/$id');
    await _reloadBudgetPlan();
    await _reloadBudgetComparison();
    notifyListeners();
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 2: SALARY DETAIL ACTIONS
  // ═══════════════════════════════════════════════════════════════════════════
  Future<void> addSalaryDetail(
    String incomeId,
    Map<String, dynamic> body,
  ) async {
    final intId = int.parse(incomeId);
    await _apiService.post('/api/incomes/$intId/salary', body);
    await reloadSalaryDetails(incomeId);
    await loadSalaryGrowth();
    await _reloadIncomes(); // reload income list as amount might have synced
    notifyListeners();
  }

  Future<void> updateSalaryDetail(
    String incomeId,
    int salaryId,
    Map<String, dynamic> body,
  ) async {
    final intId = int.parse(incomeId);
    await _apiService.put('/api/incomes/$intId/salary/$salaryId', body);
    await reloadSalaryDetails(incomeId);
    await loadSalaryGrowth();
    await _reloadIncomes();
    notifyListeners();
  }

  Future<void> deleteSalaryDetail(String incomeId, int salaryId) async {
    final intId = int.parse(incomeId);
    await _apiService.delete('/api/incomes/$intId/salary/$salaryId');
    await reloadSalaryDetails(incomeId);
    await loadSalaryGrowth();
    await _reloadIncomes();
    notifyListeners();
  }

  Future<void> setCurrentSalaryDetail(String incomeId, int salaryId) async {
    final intId = int.parse(incomeId);
    await _apiService.post('/api/incomes/$intId/salary/$salaryId/current', {});
    await reloadSalaryDetails(incomeId);
    await loadSalaryGrowth();
    await _reloadIncomes();
    notifyListeners();
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 3: VEHICLE LOGGING ACTIONS
  // ═══════════════════════════════════════════════════════════════════════════
  Future<void> addServiceRecord(
    String vehicleId,
    Map<String, dynamic> body,
  ) async {
    final intId = int.parse(vehicleId);
    await _apiService.post('/api/assets/vehicles/$intId/service', body);
    await reloadVehicleServiceRecords(vehicleId);
    await _reloadVehicles();
    notifyListeners();
  }

  Future<void> deleteServiceRecord(String vehicleId, int serviceId) async {
    final intId = int.parse(vehicleId);
    await _apiService.delete('/api/assets/vehicles/$intId/service/$serviceId');
    await reloadVehicleServiceRecords(vehicleId);
    await _reloadVehicles();
    notifyListeners();
  }

  Future<void> addFuelRecord(
    String vehicleId,
    Map<String, dynamic> body,
  ) async {
    final intId = int.parse(vehicleId);
    await _apiService.post('/api/assets/vehicles/$intId/fuel', body);
    await reloadVehicleFuelRecords(vehicleId);
    await _reloadVehicles();
    notifyListeners();
  }

  Future<void> deleteFuelRecord(String vehicleId, int fuelId) async {
    final intId = int.parse(vehicleId);
    await _apiService.delete('/api/assets/vehicles/$intId/fuel/$fuelId');
    await reloadVehicleFuelRecords(vehicleId);
    await _reloadVehicles();
    notifyListeners();
  }

  Future<void> addOrUpdateVehicleLoan(
    String vehicleId,
    Map<String, dynamic> body, {
    bool isUpdate = false,
  }) async {
    final intId = int.parse(vehicleId);
    if (isUpdate) {
      await _apiService.put('/api/assets/vehicles/$intId/loan', body);
    } else {
      await _apiService.post('/api/assets/vehicles/$intId/loan', body);
    }
    await reloadVehicleLoan(vehicleId);
    notifyListeners();
  }

  Future<void> deleteVehicleLoan(String vehicleId) async {
    final intId = int.parse(vehicleId);
    await _apiService.delete('/api/assets/vehicles/$intId/loan');
    await reloadVehicleLoan(vehicleId);
    notifyListeners();
  }

  // ═══════════════════════════════════════════════════════════════════════════
  // PHASE 4: ELECTRONICS ACTIONS
  // ═══════════════════════════════════════════════════════════════════════════
  Future<void> addElectronic(Map<String, dynamic> body) async {
    await _apiService.post('/api/electronics', body);
    await reloadElectronics();
    _reloadDashboard();
  }

  Future<void> updateElectronic(String id, Map<String, dynamic> body) async {
    final intId = int.parse(id);
    await _apiService.put('/api/electronics/$intId', body);
    await reloadElectronics();
    _reloadDashboard();
  }

  Future<void> deleteElectronic(String id) async {
    final intId = int.parse(id);
    await _apiService.delete('/api/electronics/$intId');
    await reloadElectronics();
    _reloadDashboard();
  }

  Future<void> addElectronicService(
    String deviceId,
    Map<String, dynamic> body,
  ) async {
    final intId = int.parse(deviceId);
    await _apiService.post('/api/electronics/$intId/service', body);
    await reloadElectronicServices(deviceId);
    await reloadElectronics();
    notifyListeners();
  }

  Future<void> deleteElectronicService(String deviceId, int serviceId) async {
    final intId = int.parse(deviceId);
    await _apiService.delete('/api/electronics/$intId/service/$serviceId');
    await reloadElectronicServices(deviceId);
    await reloadElectronics();
    notifyListeners();
  }

  Future<void> addOrUpdateElectronicEmi(
    String deviceId,
    Map<String, dynamic> body, {
    bool isUpdate = false,
  }) async {
    final intId = int.parse(deviceId);
    if (isUpdate) {
      await _apiService.put('/api/electronics/$intId/emi', body);
    } else {
      await _apiService.post('/api/electronics/$intId/emi', body);
    }
    await reloadElectronicEmi(deviceId);
    await reloadElectronics();
    notifyListeners();
  }

  Future<void> deleteElectronicEmi(String deviceId) async {
    final intId = int.parse(deviceId);
    await _apiService.delete('/api/electronics/$intId/emi');
    await reloadElectronicEmi(deviceId);
    await reloadElectronics();
    notifyListeners();
  }
}
