class TransactionListResponse {
  final bool success;
  final TransactionListData? data;
  final Map<String, dynamic>? metadata;

  TransactionListResponse({required this.success, this.data, this.metadata});

  factory TransactionListResponse.fromJson(Map<String, dynamic> json) => TransactionListResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? TransactionListData.fromJson(json['data'] as Map<String, dynamic>) : null,
    metadata: json['metadata'] as Map<String, dynamic>?,
  );
}

class TransactionListData {
  final List<TransactionEvent> transactions;
  final String? cursor;
  final bool hasMore;
  final int total;

  TransactionListData({
    required this.transactions,
    this.cursor,
    this.hasMore = false,
    this.total = 0,
  });

  factory TransactionListData.fromJson(Map<String, dynamic> json) => TransactionListData(
    transactions: (json['transactions'] as List?)?.map((e) => TransactionEvent.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    cursor: json['cursor'] as String?,
    hasMore: json['has_more'] as bool? ?? false,
    total: (json['total'] as num?)?.toInt() ?? 0,
  );
}

class TransactionDetailResponse {
  final bool success;
  final TransactionEvent? data;
  final Map<String, dynamic>? metadata;

  TransactionDetailResponse({required this.success, this.data, this.metadata});

  factory TransactionDetailResponse.fromJson(Map<String, dynamic> json) => TransactionDetailResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? TransactionEvent.fromJson(json['data'] as Map<String, dynamic>) : null,
    metadata: json['metadata'] as Map<String, dynamic>?,
  );
}

class TransactionSummaryResponse {
  final bool success;
  final TransactionSummary? data;
  final Map<String, dynamic>? metadata;

  TransactionSummaryResponse({required this.success, this.data, this.metadata});

  factory TransactionSummaryResponse.fromJson(Map<String, dynamic> json) => TransactionSummaryResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? TransactionSummary.fromJson(json['data'] as Map<String, dynamic>) : null,
    metadata: json['metadata'] as Map<String, dynamic>?,
  );
}

class TransactionSummary {
  final double periodIncome;
  final double periodExpenses;
  final double netFlow;
  final int incomeCount;
  final int expenseCount;
  final int totalCount;

  TransactionSummary({
    this.periodIncome = 0,
    this.periodExpenses = 0,
    this.netFlow = 0,
    this.incomeCount = 0,
    this.expenseCount = 0,
    this.totalCount = 0,
  });

  factory TransactionSummary.fromJson(Map<String, dynamic> json) => TransactionSummary(
    periodIncome: (json['period_income'] as num?)?.toDouble() ?? 0,
    periodExpenses: (json['period_expenses'] as num?)?.toDouble() ?? 0,
    netFlow: (json['net_flow'] as num?)?.toDouble() ?? 0,
    incomeCount: (json['income_count'] as num?)?.toInt() ?? 0,
    expenseCount: (json['expense_count'] as num?)?.toInt() ?? 0,
    totalCount: (json['total_count'] as num?)?.toInt() ?? 0,
  );
}

class TransactionEvent {
  final String eventId;
  final String userId;
  final String eventType;
  final double amount;
  final String currency;
  final String eventDate;
  final String? effectiveDate;
  final String? description;
  final String state;
  final String? origin;
  final String? confidence;
  final String? createdBy;
  final String? source;
  final String? destination;
  final String? reference;
  final String? category;
  final String? notes;
  final List<String>? tags;
  final String? importedFrom;
  final String? reversalOfEventId;
  final String? correlationId;
  final String? orderIndex;
  final String createdAt;
  final String? updatedAt;

  TransactionEvent({
    required this.eventId,
    this.userId = '',
    this.eventType = '',
    this.amount = 0,
    this.currency = 'INR',
    this.eventDate = '',
    this.effectiveDate,
    this.description,
    this.state = 'posted',
    this.origin,
    this.confidence,
    this.createdBy,
    this.source,
    this.destination,
    this.reference,
    this.category,
    this.notes,
    this.tags,
    this.importedFrom,
    this.reversalOfEventId,
    this.correlationId,
    this.orderIndex,
    this.createdAt = '',
    this.updatedAt,
  });

  factory TransactionEvent.fromJson(Map<String, dynamic> json) => TransactionEvent(
    eventId: json['event_id'] as String? ?? json['id'] as String? ?? '',
    userId: json['user_id'] as String? ?? '',
    eventType: json['event_type'] as String? ?? json['type'] as String? ?? 'other',
    amount: (json['amount'] as num?)?.toDouble() ?? 0,
    currency: json['currency'] as String? ?? 'INR',
    eventDate: json['event_date'] as String? ?? '',
    effectiveDate: json['effective_date'] as String?,
    description: json['description'] as String?,
    state: json['state'] as String? ?? 'posted',
    origin: json['origin'] as String?,
    confidence: json['confidence'] as String?,
    createdBy: json['created_by'] as String?,
    source: json['source'] as String?,
    destination: json['destination'] as String?,
    reference: json['reference'] as String?,
    category: json['category'] as String?,
    notes: json['notes'] as String?,
    tags: json['tags'] is List ? (json['tags'] as List).map((e) => e.toString()).toList() : null,
    importedFrom: json['imported_from'] as String?,
    reversalOfEventId: json['reversal_of_event_id'] as String?,
    correlationId: json['correlation_id'] as String?,
    orderIndex: json['order_index'] as String?,
    createdAt: json['created_at'] as String? ?? '',
    updatedAt: json['updated_at'] as String?,
  );

  bool get isIncome => amount >= 0;
  bool get isExpense => amount < 0;
  String get displayAmount => isIncome ? '+$amount' : '$amount';

  String get formattedDate {
    try {
      final dt = DateTime.parse(eventDate);
      final now = DateTime.now();
      final diff = now.difference(dt);
      if (diff.inMinutes < 1) return 'Just now';
      if (diff.inHours < 1) return '${diff.inMinutes}m ago';
      if (diff.inDays < 1) return '${diff.inHours}h ago';
      if (diff.inDays == 1) return 'Yesterday';
      if (diff.inDays < 7) return '${diff.inDays}d ago';
      return '${dt.day}/${dt.month}/${dt.year}';
    } catch (_) {
      return eventDate;
    }
  }

  String get formattedAmount {
    final abs = amount.abs();
    if (abs >= 10000000) return '${(abs / 10000000).toStringAsFixed(1)}Cr';
    if (abs >= 100000) return '${(abs / 100000).toStringAsFixed(1)}L';
    if (abs >= 1000) return '${(abs / 1000).toStringAsFixed(1)}K';
    return abs.toStringAsFixed(0);
  }
}

class TransactionFilter {
  final DateTime? startDate;
  final DateTime? endDate;
  final String? accountId;
  final String? category;
  final String? type;
  final double? minAmount;
  final double? maxAmount;
  final String? merchant;

  const TransactionFilter({
    this.startDate,
    this.endDate,
    this.accountId,
    this.category,
    this.type,
    this.minAmount,
    this.maxAmount,
    this.merchant,
  });

  TransactionFilter copyWith({
    DateTime? startDate,
    DateTime? endDate,
    String? accountId,
    String? category,
    String? type,
    double? minAmount,
    double? maxAmount,
    String? merchant,
    bool clearStartDate = false,
    bool clearEndDate = false,
    bool clearAccountId = false,
    bool clearCategory = false,
    bool clearType = false,
    bool clearMinAmount = false,
    bool clearMaxAmount = false,
    bool clearMerchant = false,
  }) => TransactionFilter(
    startDate: clearStartDate ? null : startDate ?? this.startDate,
    endDate: clearEndDate ? null : endDate ?? this.endDate,
    accountId: clearAccountId ? null : accountId ?? this.accountId,
    category: clearCategory ? null : category ?? this.category,
    type: clearType ? null : type ?? this.type,
    minAmount: clearMinAmount ? null : minAmount ?? this.minAmount,
    maxAmount: clearMaxAmount ? null : maxAmount ?? this.maxAmount,
    merchant: clearMerchant ? null : merchant ?? this.merchant,
  );

  Map<String, dynamic> toQuery() {
    final map = <String, dynamic>{};
    if (startDate != null) map['start_date'] = startDate!.toIso8601String().split('T')[0];
    if (endDate != null) map['end_date'] = endDate!.toIso8601String().split('T')[0];
    if (accountId != null) map['account_id'] = accountId;
    if (category != null) map['category'] = category;
    if (type != null) map['type'] = type;
    if (minAmount != null) map['min_amount'] = minAmount;
    if (maxAmount != null) map['max_amount'] = maxAmount;
    if (merchant != null) map['merchant'] = merchant;
    return map;
  }

  bool get hasActiveFilters =>
    startDate != null || endDate != null || accountId != null ||
    category != null || type != null || minAmount != null ||
    maxAmount != null || merchant != null;

  int get activeFilterCount {
    int count = 0;
    if (startDate != null) count++;
    if (endDate != null) count++;
    if (accountId != null) count++;
    if (category != null) count++;
    if (type != null) count++;
    if (minAmount != null || maxAmount != null) count++;
    if (merchant != null) count++;
    return count;
  }
}

enum TransactionSortField { date, amount, category }
enum TransactionSortOrder { ascending, descending }
