class HouseholdView {
  final String householdId;
  final String name;
  final String householdType;
  final String headOfHouseholdId;
  final int memberCount;
  final String status;
  final String currency;
  final String country;
  final int totalAssets;
  final int totalLiabilities;
  final int totalNetWorth;
  final String health;
  final List<String> tags;
  final String createdAt;
  final String updatedAt;

  HouseholdView({
    required this.householdId,
    required this.name,
    required this.householdType,
    required this.headOfHouseholdId,
    required this.memberCount,
    required this.status,
    required this.currency,
    required this.country,
    required this.totalAssets,
    required this.totalLiabilities,
    required this.totalNetWorth,
    required this.health,
    this.tags = const [],
    required this.createdAt,
    required this.updatedAt,
  });

  factory HouseholdView.fromJson(Map<String, dynamic> json) {
    return HouseholdView(
      householdId: json['household_id'] ?? '',
      name: json['name'] ?? '',
      householdType: json['household_type'] ?? '',
      headOfHouseholdId: json['head_of_household_id'] ?? '',
      memberCount: json['member_count'] ?? 0,
      status: json['status'] ?? '',
      currency: json['currency'] ?? '',
      country: json['country'] ?? '',
      totalAssets: json['total_assets'] ?? 0,
      totalLiabilities: json['total_liabilities'] ?? 0,
      totalNetWorth: json['total_net_worth'] ?? 0,
      health: json['health'] ?? '',
      tags: List<String>.from(json['tags'] ?? []),
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class HouseholdMemberView {
  final String userId;
  final String role;
  final String addedAt;
  final String inviteStatus;

  HouseholdMemberView({
    required this.userId,
    required this.role,
    required this.addedAt,
    required this.inviteStatus,
  });

  factory HouseholdMemberView.fromJson(Map<String, dynamic> json) {
    return HouseholdMemberView(
      userId: json['user_id'] ?? '',
      role: json['role'] ?? '',
      addedAt: json['added_at'] ?? '',
      inviteStatus: json['invite_status'] ?? '',
    );
  }
}

class LinkedAccountView {
  final String accountId;
  final String addedBy;
  final String addedAt;

  LinkedAccountView({
    required this.accountId,
    required this.addedBy,
    required this.addedAt,
  });

  factory LinkedAccountView.fromJson(Map<String, dynamic> json) {
    return LinkedAccountView(
      accountId: json['account_id'] ?? '',
      addedBy: json['added_by'] ?? '',
      addedAt: json['added_at'] ?? '',
    );
  }
}

class LinkedGoalView {
  final String goalId;
  final String addedBy;
  final String addedAt;

  LinkedGoalView({
    required this.goalId,
    required this.addedBy,
    required this.addedAt,
  });

  factory LinkedGoalView.fromJson(Map<String, dynamic> json) {
    return LinkedGoalView(
      goalId: json['goal_id'] ?? '',
      addedBy: json['added_by'] ?? '',
      addedAt: json['added_at'] ?? '',
    );
  }
}

class LinkedBudgetView {
  final String budgetId;
  final String addedBy;
  final String addedAt;

  LinkedBudgetView({
    required this.budgetId,
    required this.addedBy,
    required this.addedAt,
  });

  factory LinkedBudgetView.fromJson(Map<String, dynamic> json) {
    return LinkedBudgetView(
      budgetId: json['budget_id'] ?? '',
      addedBy: json['added_by'] ?? '',
      addedAt: json['added_at'] ?? '',
    );
  }
}

class GoalContributionView {
  final String goalId;
  final String userId;
  final int amount;
  final String date;

  GoalContributionView({
    required this.goalId,
    required this.userId,
    required this.amount,
    required this.date,
  });

  factory GoalContributionView.fromJson(Map<String, dynamic> json) {
    return GoalContributionView(
      goalId: json['goal_id'] ?? '',
      userId: json['user_id'] ?? '',
      amount: json['amount'] ?? 0,
      date: json['date'] ?? '',
    );
  }
}

class HouseholdDetailView {
  final String householdId;
  final String name;
  final String householdType;
  final String headOfHouseholdId;
  final int memberCount;
  final String status;
  final String currency;
  final String country;
  final int totalAssets;
  final int totalLiabilities;
  final int totalNetWorth;
  final String health;
  final List<String> tags;
  final String createdAt;
  final String updatedAt;
  final List<HouseholdMemberView> members;
  final String? notes;
  final List<LinkedAccountView> linkedAccounts;
  final List<LinkedGoalView> linkedGoals;
  final List<LinkedBudgetView> linkedBudgets;
  final List<GoalContributionView> goalContributions;

  HouseholdDetailView({
    required this.householdId,
    required this.name,
    required this.householdType,
    required this.headOfHouseholdId,
    required this.memberCount,
    required this.status,
    required this.currency,
    required this.country,
    required this.totalAssets,
    required this.totalLiabilities,
    required this.totalNetWorth,
    required this.health,
    this.tags = const [],
    required this.createdAt,
    required this.updatedAt,
    this.members = const [],
    this.notes,
    this.linkedAccounts = const [],
    this.linkedGoals = const [],
    this.linkedBudgets = const [],
    this.goalContributions = const [],
  });

  factory HouseholdDetailView.fromJson(Map<String, dynamic> json) {
    return HouseholdDetailView(
      householdId: json['household_id'] ?? '',
      name: json['name'] ?? '',
      householdType: json['household_type'] ?? '',
      headOfHouseholdId: json['head_of_household_id'] ?? '',
      memberCount: json['member_count'] ?? 0,
      status: json['status'] ?? '',
      currency: json['currency'] ?? '',
      country: json['country'] ?? '',
      totalAssets: json['total_assets'] ?? 0,
      totalLiabilities: json['total_liabilities'] ?? 0,
      totalNetWorth: json['total_net_worth'] ?? 0,
      health: json['health'] ?? '',
      tags: List<String>.from(json['tags'] ?? []),
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      members: (json['members'] as List?)?.map((e) => HouseholdMemberView.fromJson(e)).toList() ?? [],
      notes: json['notes'],
      linkedAccounts: (json['linked_accounts'] as List?)?.map((e) => LinkedAccountView.fromJson(e)).toList() ?? [],
      linkedGoals: (json['linked_goals'] as List?)?.map((e) => LinkedGoalView.fromJson(e)).toList() ?? [],
      linkedBudgets: (json['linked_budgets'] as List?)?.map((e) => LinkedBudgetView.fromJson(e)).toList() ?? [],
      goalContributions: (json['goal_contributions'] as List?)?.map((e) => GoalContributionView.fromJson(e)).toList() ?? [],
    );
  }
}

class HouseholdFinancialSummary {
  final String householdId;
  final int totalAssets;
  final int totalLiabilities;
  final int totalNetWorth;
  final String currency;
  final int linkedAccounts;
  final int linkedGoals;
  final int linkedBudgets;
  final int totalContributions;

  HouseholdFinancialSummary({
    required this.householdId,
    required this.totalAssets,
    required this.totalLiabilities,
    required this.totalNetWorth,
    required this.currency,
    required this.linkedAccounts,
    required this.linkedGoals,
    required this.linkedBudgets,
    required this.totalContributions,
  });

  factory HouseholdFinancialSummary.fromJson(Map<String, dynamic> json) {
    return HouseholdFinancialSummary(
      householdId: json['household_id'] ?? '',
      totalAssets: json['total_assets'] ?? 0,
      totalLiabilities: json['total_liabilities'] ?? 0,
      totalNetWorth: json['total_net_worth'] ?? 0,
      currency: json['currency'] ?? '',
      linkedAccounts: json['linked_accounts'] ?? 0,
      linkedGoals: json['linked_goals'] ?? 0,
      linkedBudgets: json['linked_budgets'] ?? 0,
      totalContributions: json['total_contributions'] ?? 0,
    );
  }
}
