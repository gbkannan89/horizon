class PlanningDashboardResponse {
  final bool success; final PlanningDashboard? data;
  PlanningDashboardResponse({required this.success, this.data});
  factory PlanningDashboardResponse.fromJson(Map<String, dynamic> json) => PlanningDashboardResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? PlanningDashboard.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class PlanningDashboard {
  final int healthScore; final String riskLevel; final int scenarioCount;
  PlanningDashboard({this.healthScore = 0, this.riskLevel = 'moderate', this.scenarioCount = 0});
  factory PlanningDashboard.fromJson(Map<String, dynamic> json) => PlanningDashboard(
    healthScore: (json['health_score'] as num?)?.toInt() ?? 0,
    riskLevel: json['risk_level'] as String? ?? 'moderate',
    scenarioCount: (json['scenario_count'] as num?)?.toInt() ?? 0,
  );
}

class ProjectionData {
  final List<ProjectionPoint> points;
  ProjectionData({this.points = const []});
  factory ProjectionData.fromJson(Map<String, dynamic> json) {
    final raw = json['projections'] as List? ?? json['data'] as List? ?? [];
    return ProjectionData(points: raw.map((e) => ProjectionPoint.fromJson(e as Map<String, dynamic>)).toList());
  }
}

class ProjectionPoint {
  final String period; final double income; final double expenses; final double netWorth;
  ProjectionPoint({required this.period, this.income = 0, this.expenses = 0, this.netWorth = 0});
  factory ProjectionPoint.fromJson(Map<String, dynamic> json) => ProjectionPoint(
    period: json['period'] as String? ?? json['year'] as String? ?? '',
    income: (json['income'] as num?)?.toDouble() ?? 0,
    expenses: (json['expenses'] as num?)?.toDouble() ?? 0,
    netWorth: (json['net_worth'] as num?)?.toDouble() ?? 0,
  );
}

class Scenario {
  final String id; final String name; final String? description; final String type;
  Scenario({required this.id, required this.name, this.description, this.type = ''});
  factory Scenario.fromJson(Map<String, dynamic> json) => Scenario(
    id: json['id'] as String? ?? json['scenario_id'] as String? ?? '',
    name: json['name'] as String? ?? json['title'] as String? ?? '',
    description: json['description'] as String?,
    type: json['type'] as String? ?? '',
  );
}
