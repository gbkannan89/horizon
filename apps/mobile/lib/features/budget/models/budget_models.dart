class BudgetCategoryModel {
  final String id;
  final String category;
  final String subcategory;
  final double budgetedAmount;
  final double spentAmount;
  final double remainingAmount;
  final bool rollover;
  final double spentPct;

  BudgetCategoryModel({
    required this.id,
    required this.category,
    this.subcategory = '',
    required this.budgetedAmount,
    required this.spentAmount,
    required this.remainingAmount,
    required this.rollover,
    required this.spentPct,
  });

  factory BudgetCategoryModel.fromJson(Map<String, dynamic> json) => BudgetCategoryModel(
    id: json['id'] as String? ?? '',
    category: json['category'] as String? ?? '',
    subcategory: json['subcategory'] as String? ?? '',
    budgetedAmount: (json['budgeted_amount'] as num?)?.toDouble() ?? 0.0,
    spentAmount: (json['spent_amount'] as num?)?.toDouble() ?? 0.0,
    remainingAmount: (json['remaining_amount'] as num?)?.toDouble() ?? 0.0,
    rollover: json['rollover'] as bool? ?? false,
    spentPct: (json['spent_pct'] as num?)?.toDouble() ?? 0.0,
  );

  Map<String, dynamic> toJson() => {
    'id': id,
    'category': category,
    'subcategory': subcategory,
    'budgeted_amount': budgetedAmount,
    'spent_amount': spentAmount,
    'remaining_amount': remainingAmount,
    'rollover': rollover,
    'spent_pct': spentPct,
  };
}

class BudgetModel {
  final String budgetId;
  final String name;
  final String period;
  final String startDate;
  final String endDate;
  final String status;
  final double totalBudgeted;
  final double totalSpent;
  final double totalRemaining;
  final String currency;
  final List<BudgetCategoryModel> categories;
  final List<String> tags;
  final String createdAt;
  final String updatedAt;

  BudgetModel({
    required this.budgetId,
    required this.name,
    required this.period,
    required this.startDate,
    required this.endDate,
    required this.status,
    required this.totalBudgeted,
    required this.totalSpent,
    required this.totalRemaining,
    required this.currency,
    required this.categories,
    required this.tags,
    required this.createdAt,
    required this.updatedAt,
  });

  factory BudgetModel.fromJson(Map<String, dynamic> json) => BudgetModel(
    budgetId: json['budget_id'] as String? ?? '',
    name: json['name'] as String? ?? '',
    period: json['period'] as String? ?? 'Monthly',
    startDate: json['start_date'] as String? ?? '',
    endDate: json['end_date'] as String? ?? '',
    status: json['status'] as String? ?? 'Draft',
    totalBudgeted: (json['total_budgeted'] as num?)?.toDouble() ?? 0.0,
    totalSpent: (json['total_spent'] as num?)?.toDouble() ?? 0.0,
    totalRemaining: (json['total_remaining'] as num?)?.toDouble() ?? 0.0,
    currency: json['currency'] as String? ?? 'INR',
    categories: (json['categories'] as List?)
            ?.map((e) => BudgetCategoryModel.fromJson(e as Map<String, dynamic>))
            .toList() ??
        [],
    tags: (json['tags'] as List?)?.map((e) => e.toString()).toList() ?? [],
    createdAt: json['created_at'] as String? ?? '',
    updatedAt: json['updated_at'] as String? ?? '',
  );

  Map<String, dynamic> toJson() => {
    'budget_id': budgetId,
    'name': name,
    'period': period,
    'start_date': startDate,
    'end_date': endDate,
    'status': status,
    'total_budgeted': totalBudgeted,
    'total_spent': totalSpent,
    'total_remaining': totalRemaining,
    'currency': currency,
    'categories': categories.map((e) => e.toJson()).toList(),
    'tags': tags,
    'created_at': createdAt,
    'updated_at': updatedAt,
  };

  double get spentPercentage {
    if (totalBudgeted == 0) return 0.0;
    return (totalSpent / totalBudgeted) * 100.0;
  }
}

class BudgetListResponse {
  final bool success;
  final BudgetListData? data;

  BudgetListResponse({required this.success, this.data});

  factory BudgetListResponse.fromJson(Map<String, dynamic> json) => BudgetListResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? BudgetListData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class BudgetListData {
  final List<BudgetModel> budgets;
  final String? nextCursor;
  final bool hasMore;

  BudgetListData({
    required this.budgets,
    this.nextCursor,
    this.hasMore = false,
  });

  factory BudgetListData.fromJson(Map<String, dynamic> json) => BudgetListData(
    budgets: (json['budgets'] as List?)
            ?.map((e) => BudgetModel.fromJson(e as Map<String, dynamic>))
            .toList() ??
        [],
    nextCursor: json['next_cursor'] as String?,
    hasMore: json['has_more'] as bool? ?? false,
  );
}

class BudgetDetailResponse {
  final bool success;
  final BudgetModel? data;

  BudgetDetailResponse({required this.success, this.data});

  factory BudgetDetailResponse.fromJson(Map<String, dynamic> json) => BudgetDetailResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? BudgetModel.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class CreateCategoryRequest {
  final String category;
  final String subcategory;
  final double budgetedAmount;
  final bool rollover;

  CreateCategoryRequest({
    required this.category,
    this.subcategory = '',
    required this.budgetedAmount,
    required this.rollover,
  });

  Map<String, dynamic> toJson() => {
    'category': category,
    if (subcategory.isNotEmpty) 'subcategory': subcategory,
    'budgeted_amount': budgetedAmount.toInt(), // Go API expects int64
    'rollover': rollover,
  };
}

class CreateBudgetRequest {
  final String name;
  final String period;
  final String startDate;
  final String endDate;
  final String currency;
  final List<CreateCategoryRequest> categories;
  final List<String> tags;

  CreateBudgetRequest({
    required this.name,
    required this.period,
    required this.startDate,
    required this.endDate,
    required this.currency,
    required this.categories,
    required this.tags,
  });

  Map<String, dynamic> toJson() => {
    'name': name,
    'period': period,
    'start_date': startDate,
    'end_date': endDate,
    'currency': currency,
    'categories': categories.map((e) => e.toJson()).toList(),
    'tags': tags,
  };
}

class BudgetResultResponse {
  final bool success;
  final BudgetResultData? data;

  BudgetResultResponse({required this.success, this.data});

  factory BudgetResultResponse.fromJson(Map<String, dynamic> json) => BudgetResultResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? BudgetResultData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class BudgetResultData {
  final String budgetId;
  final String status;
  final bool success;

  BudgetResultData({
    required this.budgetId,
    required this.status,
    required this.success,
  });

  factory BudgetResultData.fromJson(Map<String, dynamic> json) => BudgetResultData(
    budgetId: json['budget_id'] as String? ?? '',
    status: json['status'] as String? ?? '',
    success: json['success'] as bool? ?? false,
  );
}
