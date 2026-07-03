class AcctDashboardResponse {
  final bool success; final AcctDashboardData? data;
  AcctDashboardResponse({required this.success, this.data});
  factory AcctDashboardResponse.fromJson(Map<String, dynamic> json) => AcctDashboardResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? AcctDashboardData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class AcctDashboardData {
  final int totalBalance; final int totalAccounts; final Map<String, int> accountsByType;
  final List<AcctDashboardCard> cards;
  AcctDashboardData({required this.totalBalance, required this.totalAccounts, required this.accountsByType, required this.cards});
  factory AcctDashboardData.fromJson(Map<String, dynamic> json) => AcctDashboardData(
    totalBalance: (json['total_balance'] as num?)?.toInt() ?? 0,
    totalAccounts: (json['total_accounts'] as num?)?.toInt() ?? 0,
    accountsByType: (json['accounts_by_type'] as Map<String, dynamic>?)?.map((k, v) => MapEntry(k, (v as num).toInt())) ?? {},
    cards: (json['cards'] as List?)?.map((e) => AcctDashboardCard.fromJson(e as Map<String, dynamic>)).toList() ?? [],
  );
}

class AcctDashboardCard {
  final String cardId; final String cardType; final String title; final String summary; final dynamic data; final int priority;
  AcctDashboardCard({required this.cardId, required this.cardType, required this.title, required this.summary, this.data, required this.priority});
  factory AcctDashboardCard.fromJson(Map<String, dynamic> json) => AcctDashboardCard(
    cardId: json['card_id'] as String? ?? '', cardType: json['card_type'] as String? ?? '',
    title: json['title'] as String? ?? '', summary: json['summary'] as String? ?? '',
    data: json['data'], priority: (json['priority'] as num?)?.toInt() ?? 0,
  );
}

class AcctListResponse {
  final bool success; final AcctListData? data;
  AcctListResponse({required this.success, this.data});
  factory AcctListResponse.fromJson(Map<String, dynamic> json) => AcctListResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? AcctListData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class AcctListData {
  final List<AcctCardData> accounts; final int total; final String? cursor; final bool hasMore; final AcctSummaryData? summary;
  AcctListData({required this.accounts, required this.total, this.cursor, required this.hasMore, this.summary});
  factory AcctListData.fromJson(Map<String, dynamic> json) => AcctListData(
    accounts: (json['accounts'] as List?)?.map((e) => AcctCardData.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    total: (json['total'] as num?)?.toInt() ?? 0,
    cursor: json['cursor'] as String?, hasMore: json['has_more'] as bool? ?? false,
    summary: json['summary'] != null ? AcctSummaryData.fromJson(json['summary'] as Map<String, dynamic>) : null,
  );
}

class AcctCardData {
  final String accountId; final String accountName; final String accountType;
  final String status; final String currency; final int currentBalance; final int availableBalance;
  final String? institutionName; final bool hasImport;
  AcctCardData({required this.accountId, required this.accountName, required this.accountType,
    required this.status, required this.currency, required this.currentBalance,
    required this.availableBalance, this.institutionName, required this.hasImport});
  factory AcctCardData.fromJson(Map<String, dynamic> json) => AcctCardData(
    accountId: json['account_id'] as String? ?? '', accountName: json['account_name'] as String? ?? '',
    accountType: json['account_type'] as String? ?? '', status: json['status'] as String? ?? '',
    currency: json['currency'] as String? ?? '',
    currentBalance: (json['current_balance'] as num?)?.toInt() ?? 0,
    availableBalance: (json['available_balance'] as num?)?.toInt() ?? 0,
    institutionName: json['institution_name'] as String?, hasImport: json['has_import'] as bool? ?? false,
  );
}

class AcctSummaryData {
  final int totalBalance; final Map<String, int> countByType; final Map<String, int> countByStatus;
  AcctSummaryData({required this.totalBalance, required this.countByType, required this.countByStatus});
  factory AcctSummaryData.fromJson(Map<String, dynamic> json) => AcctSummaryData(
    totalBalance: (json['total_balance'] as num?)?.toInt() ?? 0,
    countByType: (json['count_by_type'] as Map<String, dynamic>?)?.map((k, v) => MapEntry(k, (v as num).toInt())) ?? {},
    countByStatus: (json['count_by_status'] as Map<String, dynamic>?)?.map((k, v) => MapEntry(k, (v as num).toInt())) ?? {},
  );
}

class AcctDetailResponse {
  final bool success; final AcctDetailData? data;
  AcctDetailResponse({required this.success, this.data});
  factory AcctDetailResponse.fromJson(Map<String, dynamic> json) => AcctDetailResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? AcctDetailData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class AcctDetailData {
  final String accountId; final String accountName; final String accountType;
  final String status; final String currency; final String? institutionName;
  final String? importStatus; final String? lastSyncTime;
  final String openedDate; final String? closedDate; final String liquidityProfile;
  final String accountHealth; final int? creditLimit; final double? interestRate;
  final String visibility; final List<String> tags; final String notes;
  final String createdAt; final List<BalanceEntry> balances;
  final List<TransPreview> transactions;

  AcctDetailData({required this.accountId, required this.accountName, required this.accountType,
    required this.status, required this.currency, this.institutionName,
    this.importStatus, this.lastSyncTime, required this.openedDate,
    this.closedDate, required this.liquidityProfile, required this.accountHealth,
    this.creditLimit, this.interestRate, required this.visibility,
    required this.tags, required this.notes, required this.createdAt,
    required this.balances, required this.transactions});

  factory AcctDetailData.fromJson(Map<String, dynamic> json) => AcctDetailData(
    accountId: json['account_id'] as String? ?? '', accountName: json['account_name'] as String? ?? '',
    accountType: json['account_type'] as String? ?? '', status: json['status'] as String? ?? '',
    currency: json['currency'] as String? ?? '',
    institutionName: json['institution_name'] as String?,
    importStatus: json['import_status'] as String?, lastSyncTime: json['last_sync_time'] as String?,
    openedDate: json['opened_date'] as String? ?? '', closedDate: json['closed_date'] as String?,
    liquidityProfile: json['liquidity_profile'] as String? ?? '',
    accountHealth: json['account_health'] as String? ?? '',
    creditLimit: (json['credit_limit'] as num?)?.toInt(),
    interestRate: (json['interest_rate'] as num?)?.toDouble(),
    visibility: json['visibility'] as String? ?? '',
    tags: (json['tags'] as List?)?.map((e) => e as String).toList() ?? [],
    notes: json['notes'] as String? ?? '',
    createdAt: json['created_at'] as String? ?? '',
    balances: (json['balances'] as List?)?.map((e) => BalanceEntry.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    transactions: (json['transactions'] as List?)?.map((e) => TransPreview.fromJson(e as Map<String, dynamic>)).toList() ?? [],
  );
}

class BalanceEntry {
  final String type; final int value; final String label; final String description;
  BalanceEntry({required this.type, required this.value, required this.label, required this.description});
  factory BalanceEntry.fromJson(Map<String, dynamic> json) => BalanceEntry(
    type: json['type'] as String? ?? '', value: (json['value'] as num?)?.toInt() ?? 0,
    label: json['label'] as String? ?? '', description: json['description'] as String? ?? '',
  );
}

class TransPreview {
  final String eventId; final String eventType; final int amount; final String currency;
  final String description; final String category; final String eventDate; final int runningBalance;
  TransPreview({required this.eventId, required this.eventType, required this.amount,
    required this.currency, required this.description, required this.category,
    required this.eventDate, required this.runningBalance});
  factory TransPreview.fromJson(Map<String, dynamic> json) => TransPreview(
    eventId: json['event_id'] as String? ?? '', eventType: json['event_type'] as String? ?? '',
    amount: (json['amount'] as num?)?.toInt() ?? 0, currency: json['currency'] as String? ?? '',
    description: json['description'] as String? ?? '', category: json['category'] as String? ?? '',
    eventDate: json['event_date'] as String? ?? '',
    runningBalance: (json['running_balance'] as num?)?.toInt() ?? 0,
  );
}

class BalanceSummaryResponse {
  final bool success; final BalanceSummaryData? data;
  BalanceSummaryResponse({required this.success, this.data});
  factory BalanceSummaryResponse.fromJson(Map<String, dynamic> json) => BalanceSummaryResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? BalanceSummaryData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class BalanceSummaryData {
  final int totalBalance; final int totalAvailable; final int totalSpendable;
  final double creditUtilizationPct; final List<AcctBalanceItem> accountBalances;
  BalanceSummaryData({required this.totalBalance, required this.totalAvailable, required this.totalSpendable,
    required this.creditUtilizationPct, required this.accountBalances});
  factory BalanceSummaryData.fromJson(Map<String, dynamic> json) => BalanceSummaryData(
    totalBalance: (json['total_balance'] as num?)?.toInt() ?? 0,
    totalAvailable: (json['total_available'] as num?)?.toInt() ?? 0,
    totalSpendable: (json['total_spendable'] as num?)?.toInt() ?? 0,
    creditUtilizationPct: (json['credit_utilization_pct'] as num?)?.toDouble() ?? 0,
    accountBalances: (json['account_balances'] as List?)?.map((e) => AcctBalanceItem.fromJson(e as Map<String, dynamic>)).toList() ?? [],
  );
}

class AcctBalanceItem {
  final String accountId; final String name; final String type; final int balance; final int available;
  AcctBalanceItem({required this.accountId, required this.name, required this.type, required this.balance, required this.available});
  factory AcctBalanceItem.fromJson(Map<String, dynamic> json) => AcctBalanceItem(
    accountId: json['account_id'] as String? ?? '', name: json['name'] as String? ?? '',
    type: json['type'] as String? ?? '', balance: (json['balance'] as num?)?.toInt() ?? 0,
    available: (json['available'] as num?)?.toInt() ?? 0,
  );
}

class CashFlowResponse {
  final bool success; final CashFlowData? data;
  CashFlowResponse({required this.success, this.data});
  factory CashFlowResponse.fromJson(Map<String, dynamic> json) => CashFlowResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? CashFlowData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class CashFlowData {
  final int periodInflow; final int periodOutflow; final int netFlow; final int projectedFlow;
  CashFlowData({required this.periodInflow, required this.periodOutflow, required this.netFlow, required this.projectedFlow});
  factory CashFlowData.fromJson(Map<String, dynamic> json) => CashFlowData(
    periodInflow: (json['period_inflow'] as num?)?.toInt() ?? 0,
    periodOutflow: (json['period_outflow'] as num?)?.toInt() ?? 0,
    netFlow: (json['net_flow'] as num?)?.toInt() ?? 0,
    projectedFlow: (json['projected_flow'] as num?)?.toInt() ?? 0,
  );
}

class AcctHealthResponse {
  final bool success; final AcctHealthData? data;
  AcctHealthResponse({required this.success, this.data});
  factory AcctHealthResponse.fromJson(Map<String, dynamic> json) => AcctHealthResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? AcctHealthData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class AcctHealthData {
  final Map<String, int> accountsByHealth; final double healthyPct;
  AcctHealthData({required this.accountsByHealth, required this.healthyPct});
  factory AcctHealthData.fromJson(Map<String, dynamic> json) => AcctHealthData(
    accountsByHealth: (json['accounts_by_health'] as Map<String, dynamic>?)?.map((k, v) => MapEntry(k, (v as num).toInt())) ?? {},
    healthyPct: (json['healthy_pct'] as num?)?.toDouble() ?? 0,
  );
}

class HouseholdAcctsResponse {
  final bool success; final HouseholdAcctsData? data;
  HouseholdAcctsResponse({required this.success, this.data});
  factory HouseholdAcctsResponse.fromJson(Map<String, dynamic> json) => HouseholdAcctsResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? HouseholdAcctsData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class HouseholdAcctsData {
  final List<HouseholdAccountView> accounts; final String? nextCursor; final bool hasMore;
  HouseholdAcctsData({required this.accounts, this.nextCursor, required this.hasMore});
  factory HouseholdAcctsData.fromJson(Map<String, dynamic> json) => HouseholdAcctsData(
    accounts: (json['accounts'] as List?)?.map((e) => HouseholdAccountView.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    nextCursor: json['next_cursor'] as String?, hasMore: json['has_more'] as bool? ?? false,
  );
}

class HouseholdAccountView {
  final String accountId; final String accountName; final String accountType;
  final String classification; final String currency; final String status;
  final String ownerId; final String visibility; final String liquidity;
  final String health; final String createdAt;
  HouseholdAccountView({required this.accountId, required this.accountName,
    required this.accountType, required this.classification, required this.currency,
    required this.status, required this.ownerId, required this.visibility,
    required this.liquidity, required this.health, required this.createdAt});
  factory HouseholdAccountView.fromJson(Map<String, dynamic> json) => HouseholdAccountView(
    accountId: json['account_id'] as String? ?? '',
    accountName: json['account_name'] as String? ?? '',
    accountType: json['account_type'] as String? ?? '',
    classification: json['classification'] as String? ?? '',
    currency: json['currency'] as String? ?? '',
    status: json['status'] as String? ?? '',
    ownerId: json['owner_id'] as String? ?? '',
    visibility: json['visibility'] as String? ?? '',
    liquidity: json['liquidity_profile'] as String? ?? '',
    health: json['account_health'] as String? ?? '',
    createdAt: json['created_at'] as String? ?? '',
  );
}


