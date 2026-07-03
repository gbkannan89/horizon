class PfDashboardResponse {
  final bool success; final PfDashboardData? data;
  PfDashboardResponse({required this.success, this.data});
  factory PfDashboardResponse.fromJson(Map<String, dynamic> json) => PfDashboardResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? PfDashboardData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class PfDashboardData {
  final int portfolioValue; final double totalReturn; final double totalReturnPct;
  final int riskScore; final String riskLevel; final List<PfCard> cards;
  PfDashboardData({required this.portfolioValue, required this.totalReturn, required this.totalReturnPct, required this.riskScore, required this.riskLevel, required this.cards});
  factory PfDashboardData.fromJson(Map<String, dynamic> json) => PfDashboardData(
    portfolioValue: (json['portfolio_value'] as num?)?.toInt() ?? 0,
    totalReturn: (json['total_return'] as num?)?.toDouble() ?? 0,
    totalReturnPct: (json['total_return_pct'] as num?)?.toDouble() ?? 0,
    riskScore: (json['risk_score'] as num?)?.toInt() ?? 0,
    riskLevel: json['risk_level'] as String? ?? '',
    cards: (json['cards'] as List?)?.map((e) => PfCard.fromJson(e as Map<String, dynamic>)).toList() ?? [],
  );
}

class PfCard {
  final String cardId; final String cardType; final String title; final String summary; final dynamic data; final int priority;
  PfCard({required this.cardId, required this.cardType, required this.title, required this.summary, this.data, required this.priority});
  factory PfCard.fromJson(Map<String, dynamic> json) => PfCard(
    cardId: json['card_id'] as String? ?? '', cardType: json['card_type'] as String? ?? '',
    title: json['title'] as String? ?? '', summary: json['summary'] as String? ?? '',
    data: json['data'], priority: (json['priority'] as num?)?.toInt() ?? 0,
  );
}

class AllocationResponse {
  final bool success; final AllocationData? data;
  AllocationResponse({required this.success, this.data});
  factory AllocationResponse.fromJson(Map<String, dynamic> json) => AllocationResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? AllocationData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class AllocationData {
  final List<AllocEntry> allocations; final int totalValue; final double diversification;
  AllocationData({required this.allocations, required this.totalValue, required this.diversification});
  factory AllocationData.fromJson(Map<String, dynamic> json) => AllocationData(
    allocations: (json['allocations'] as List?)?.map((e) => AllocEntry.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    totalValue: (json['total_value'] as num?)?.toInt() ?? 0,
    diversification: (json['diversification_score'] as num?)?.toDouble() ?? 0,
  );
}

class AllocEntry {
  final String label; final int value; final double percent; final double target; final double drift;
  AllocEntry({required this.label, required this.value, required this.percent, required this.target, required this.drift});
  factory AllocEntry.fromJson(Map<String, dynamic> json) => AllocEntry(
    label: json['label'] as String? ?? '',
    value: (json['value'] as num?)?.toInt() ?? 0,
    percent: (json['percent'] as num?)?.toDouble() ?? 0,
    target: (json['target'] as num?)?.toDouble() ?? 0,
    drift: (json['drift'] as num?)?.toDouble() ?? 0,
  );
}

class PerformanceResponse {
  final bool success; final PerformanceData? data;
  PerformanceResponse({required this.success, this.data});
  factory PerformanceResponse.fromJson(Map<String, dynamic> json) => PerformanceResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? PerformanceData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class PerformanceData {
  final double periodReturn; final double periodReturnPct; final double benchmarkReturn;
  final int unrealizedGL; final int realizedGL; final String period;
  PerformanceData({required this.periodReturn, required this.periodReturnPct, required this.benchmarkReturn, required this.unrealizedGL, required this.realizedGL, required this.period});
  factory PerformanceData.fromJson(Map<String, dynamic> json) => PerformanceData(
    periodReturn: (json['period_return'] as num?)?.toDouble() ?? 0,
    periodReturnPct: (json['period_return_pct'] as num?)?.toDouble() ?? 0,
    benchmarkReturn: (json['benchmark_return'] as num?)?.toDouble() ?? 0,
    unrealizedGL: (json['unrealized_gain_loss'] as num?)?.toInt() ?? 0,
    realizedGL: (json['realized_gain_loss'] as num?)?.toInt() ?? 0,
    period: json['period'] as String? ?? '1M',
  );
}

class RiskResponse {
  final bool success; final RiskData? data;
  RiskResponse({required this.success, this.data});
  factory RiskResponse.fromJson(Map<String, dynamic> json) => RiskResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? RiskData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class RiskData {
  final int riskScore; final String riskLevel; final double var_; final double sharpeRatio; final double volatility; final double maxDrawdown;
  RiskData({required this.riskScore, required this.riskLevel, required this.var_, required this.sharpeRatio, required this.volatility, required this.maxDrawdown});
  factory RiskData.fromJson(Map<String, dynamic> json) => RiskData(
    riskScore: (json['risk_score'] as num?)?.toInt() ?? 0,
    riskLevel: json['risk_level'] as String? ?? '',
    var_: (json['var'] as num?)?.toDouble() ?? 0,
    sharpeRatio: (json['sharpe_ratio'] as num?)?.toDouble() ?? 0,
    volatility: (json['volatility'] as num?)?.toDouble() ?? 0,
    maxDrawdown: (json['max_drawdown'] as num?)?.toDouble() ?? 0,
  );
}

class ProjectionResponse {
  final bool success; final ProjectionData? data;
  ProjectionResponse({required this.success, this.data});
  factory ProjectionResponse.fromJson(Map<String, dynamic> json) => ProjectionResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? ProjectionData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class ProjectionData {
  final double projectedValue; final String confidence; final int horizonYears; final double annualReturn;
  ProjectionData({required this.projectedValue, required this.confidence, required this.horizonYears, required this.annualReturn});
  factory ProjectionData.fromJson(Map<String, dynamic> json) => ProjectionData(
    projectedValue: (json['projected_value'] as num?)?.toDouble() ?? 0,
    confidence: json['confidence'] as String? ?? '',
    horizonYears: (json['horizon_years'] as num?)?.toInt() ?? 0,
    annualReturn: (json['annual_return'] as num?)?.toDouble() ?? 0,
  );
}

class SimResponse {
  final bool success; final SimData? data;
  SimResponse({required this.success, this.data});
  factory SimResponse.fromJson(Map<String, dynamic> json) => SimResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? SimData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class SimData {
  final List<SimEntry> simulations; final int totalCount;
  SimData({required this.simulations, required this.totalCount});
  factory SimData.fromJson(Map<String, dynamic> json) => SimData(
    simulations: (json['simulations'] as List?)?.map((e) => SimEntry.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    totalCount: (json['total_count'] as num?)?.toInt() ?? 0,
  );
}

class SimEntry {
  final String simId; final String scenario; final String outcome;
  SimEntry({required this.simId, required this.scenario, required this.outcome});
  factory SimEntry.fromJson(Map<String, dynamic> json) => SimEntry(
    simId: json['sim_id'] as String? ?? '', scenario: json['scenario'] as String? ?? '',
    outcome: json['outcome'] as String? ?? '',
  );
}

class CardViewResponse {
  final bool success; final CardViewData? data;
  CardViewResponse({required this.success, this.data});
  factory CardViewResponse.fromJson(Map<String, dynamic> json) => CardViewResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? CardViewData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class CardViewData {
  final List<CardViewItem> items; final int count;
  CardViewData({required this.items, required this.count});
  factory CardViewData.fromJson(Map<String, dynamic> json) => CardViewData(
    items: (json['items'] as List?)?.map((e) => CardViewItem.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    count: (json['count'] as num?)?.toInt() ?? 0,
  );
}

class CardViewItem {
  final String id; final String title; final String value;
  CardViewItem({required this.id, required this.title, required this.value});
  factory CardViewItem.fromJson(Map<String, dynamic> json) => CardViewItem(
    id: json['id'] as String? ?? '', title: json['title'] as String? ?? '',
    value: json['value'] as String? ?? '',
  );
}

String _fmt(int v) {
  if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr';
  if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L';
  return '₹$v';
}
