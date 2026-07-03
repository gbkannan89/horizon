
class DashboardResponse {
  final bool success;
  final DashboardData? data;
  final Map<String, dynamic>? metadata;

  DashboardResponse({required this.success, this.data, this.metadata});

  factory DashboardResponse.fromJson(Map<String, dynamic> json) => DashboardResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? DashboardData.fromJson(json['data'] as Map<String, dynamic>) : null,
    metadata: json['metadata'] as Map<String, dynamic>?,
  );
}

class DashboardData {
  final String dashboardId;
  final String mode;
  final String state;
  final List<WidgetModel> tier1;
  final List<WidgetModel> tier2;
  final List<WidgetModel> tier3;
  final WidgetModel? criticalAlert;
  final String lastRefreshTime;

  DashboardData({
    required this.dashboardId, required this.mode, required this.state,
    required this.tier1, required this.tier2, required this.tier3,
    this.criticalAlert, required this.lastRefreshTime,
  });

  factory DashboardData.fromJson(Map<String, dynamic> json) => DashboardData(
    dashboardId: json['dashboard_id'] as String? ?? '',
    mode: json['mode'] as String? ?? '',
    state: json['state'] as String? ?? '',
    tier1: (json['tier1'] as List?)?.map((e) => WidgetModel.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    tier2: (json['tier2'] as List?)?.map((e) => WidgetModel.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    tier3: (json['tier3'] as List?)?.map((e) => WidgetModel.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    criticalAlert: json['critical_alert'] != null ? WidgetModel.fromJson(json['critical_alert'] as Map<String, dynamic>) : null,
    lastRefreshTime: json['last_refresh_time'] as String? ?? '',
  );
}

class WidgetModel {
  final String widgetId;
  final String widgetType;
  final String title;
  final int tier;
  final int priority;
  final dynamic data;
  final String confidence;
  final bool visible;

  WidgetModel({
    required this.widgetId, required this.widgetType, required this.title,
    required this.tier, required this.priority, this.data,
    required this.confidence, required this.visible,
  });

  factory WidgetModel.fromJson(Map<String, dynamic> json) => WidgetModel(
    widgetId: json['widget_id'] as String? ?? '',
    widgetType: json['widget_type'] as String? ?? '',
    title: json['title'] as String? ?? '',
    tier: json['tier'] as int? ?? 3,
    priority: json['priority'] as int? ?? 0,
    data: json['data'],
    confidence: json['confidence'] as String? ?? 'Medium',
    visible: json['visible'] as bool? ?? true,
  );

  String get numericValue {
    if (data is num) return formatMoney((data as num).toInt());
    if (data is String) return data as String;
    return title;
  }

  static String formatMoney(int amount) {
    if (amount >= 10000000) return '₹${(amount / 10000000).toStringAsFixed(2)}Cr';
    if (amount >= 100000) return '₹${(amount / 100000).toStringAsFixed(2)}L';
    return '₹$amount';
  }
}

class SummaryResponse {
  final bool success;
  final SummaryData? data;

  SummaryResponse({required this.success, this.data});

  factory SummaryResponse.fromJson(Map<String, dynamic> json) => SummaryResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? SummaryData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class SummaryData {
  final int netWorth;
  final int healthScore;
  final String healthGrade;
  final int riskScore;
  final String riskLevel;
  final int goalsOnTrack;
  final int totalGoals;
  final int cashBalance;
  final int portfolioValue;
  final int totalDebt;

  SummaryData({
    required this.netWorth, required this.healthScore, required this.healthGrade,
    required this.riskScore, required this.riskLevel,
    required this.goalsOnTrack, required this.totalGoals,
    required this.cashBalance, required this.portfolioValue, required this.totalDebt,
  });

  factory SummaryData.fromJson(Map<String, dynamic> json) => SummaryData(
    netWorth: (json['net_worth'] as num?)?.toInt() ?? 0,
    healthScore: (json['health_score'] as num?)?.toInt() ?? 0,
    healthGrade: json['health_grade'] as String? ?? '',
    riskScore: (json['risk_score'] as num?)?.toInt() ?? 0,
    riskLevel: json['risk_level'] as String? ?? '',
    goalsOnTrack: (json['goals_on_track'] as num?)?.toInt() ?? 0,
    totalGoals: (json['total_goals'] as num?)?.toInt() ?? 0,
    cashBalance: (json['cash_balance'] as num?)?.toInt() ?? 0,
    portfolioValue: (json['portfolio_value'] as num?)?.toInt() ?? 0,
    totalDebt: (json['total_debt'] as num?)?.toInt() ?? 0,
  );
}
