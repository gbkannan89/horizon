class GoalsListResponse {
  final bool success;
  final GoalSummaryList? data;
  GoalsListResponse({required this.success, this.data});
  factory GoalsListResponse.fromJson(Map<String, dynamic> json) => GoalsListResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? GoalSummaryList.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class GoalSummaryList {
  final List<GoalSummary> goals;
  final int total;
  GoalSummaryList({required this.goals, required this.total});
  factory GoalSummaryList.fromJson(Map<String, dynamic> json) => GoalSummaryList(
    goals: (json['goals'] as List?)?.map((e) => GoalSummary.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    total: (json['total'] as num?)?.toInt() ?? 0,
  );
}

class GoalSummary {
  final String goalId;
  final String name;
  final int progress;
  final String status;
  final String importance;
  final int priority;
  final bool hasRecommendation;

  GoalSummary({
    required this.goalId, required this.name, required this.progress,
    required this.status, required this.importance, required this.priority,
    required this.hasRecommendation,
  });

  factory GoalSummary.fromJson(Map<String, dynamic> json) => GoalSummary(
    goalId: json['goal_id'] as String? ?? '',
    name: json['name'] as String? ?? '',
    progress: (json['progress'] as num?)?.toInt() ?? 0,
    status: json['status'] as String? ?? '',
    importance: json['importance'] as String? ?? '',
    priority: (json['priority'] as num?)?.toInt() ?? 0,
    hasRecommendation: json['has_recommendation'] as bool? ?? false,
  );
}

// Goal Dashboard response uses embedded Go structs from the engine
class GoalDashboardResponse {
  final bool success;
  final GoalDashboardData? data;
  GoalDashboardResponse({required this.success, this.data});
  factory GoalDashboardResponse.fromJson(Map<String, dynamic> json) => GoalDashboardResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? GoalDashboardData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class GoalDashboardData {
  final String goalId;
  final String goalName;
  final int progress;
  final String status;
  final String importance;
  final double targetAmount;
  final double currentValue;
  final double fundingGap;
  final String? targetDate;
  final List<GoalCardData> cards;

  GoalDashboardData({
    required this.goalId, required this.goalName, required this.progress,
    required this.status, required this.importance,
    required this.targetAmount, required this.currentValue,
    required this.fundingGap, this.targetDate, required this.cards,
  });

  factory GoalDashboardData.fromJson(Map<String, dynamic> json) => GoalDashboardData(
    goalId: json['goal_id'] as String? ?? '',
    goalName: json['goal_name'] as String? ?? '',
    progress: (json['progress'] as num?)?.toInt() ?? 0,
    status: json['status'] as String? ?? '',
    importance: json['importance'] as String? ?? '',
    targetAmount: (json['target_amount'] as num?)?.toDouble() ?? 0,
    currentValue: (json['current_value'] as num?)?.toDouble() ?? 0,
    fundingGap: (json['funding_gap'] as num?)?.toDouble() ?? 0,
    targetDate: json['target_date'] as String?,
    cards: (json['cards'] as List?)?.map((e) => GoalCardData.fromJson(e as Map<String, dynamic>)).toList() ?? [],
  );
}

class GoalCardData {
  final String cardId;
  final String cardType;
  final String title;
  final String summary;
  final dynamic data;
  final int priority;

  GoalCardData({required this.cardId, required this.cardType, required this.title, required this.summary, this.data, required this.priority});

  factory GoalCardData.fromJson(Map<String, dynamic> json) => GoalCardData(
    cardId: json['card_id'] as String? ?? '',
    cardType: json['card_type'] as String? ?? '',
    title: json['title'] as String? ?? '',
    summary: json['summary'] as String? ?? '',
    data: json['data'],
    priority: (json['priority'] as num?)?.toInt() ?? 0,
  );

  String get formattedValue {
    if (data == null) return '';
    if (data is num) return _fmt((data as num).toDouble());
    if (data is Map) {
      final m = data as Map;
      if (m['pct'] != null) return '${m['pct']}%';
      if (m['gap'] != null) return _fmt((m['gap'] as num).toDouble());
      if (m['remaining'] != null) return _fmt((m['remaining'] as num).toDouble());
      if (m['projected'] != null) return _fmt((m['projected'] as num).toDouble());
    }
    return summary;
  }

  static String _fmt(double v) {
    if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr';
    if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L';
    return '₹${v.toInt()}';
  }
}

// Individual response types
class GoalProgressResponse {
  final bool success;
  final GoalProgressData? data;
  GoalProgressResponse({required this.success, this.data});
  factory GoalProgressResponse.fromJson(Map<String, dynamic> json) => GoalProgressResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? GoalProgressData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class GoalProgressData {
  final String goalId; final String goalName; final int progressPct;
  final double currentValue; final double targetAmount; final double remainingAmount;
  final String status; final double monthlyContribution; final int monthsToTarget;
  GoalProgressData({required this.goalId, required this.goalName, required this.progressPct,
    required this.currentValue, required this.targetAmount, required this.remainingAmount,
    required this.status, required this.monthlyContribution, required this.monthsToTarget});
  factory GoalProgressData.fromJson(Map<String, dynamic> json) => GoalProgressData(
    goalId: json['goal_id'] as String? ?? '', goalName: json['goal_name'] as String? ?? '',
    progressPct: (json['progress_pct'] as num?)?.toInt() ?? 0,
    currentValue: (json['current_value'] as num?)?.toDouble() ?? 0,
    targetAmount: (json['target_amount'] as num?)?.toDouble() ?? 0,
    remainingAmount: (json['remaining_amount'] as num?)?.toDouble() ?? 0,
    status: json['status'] as String? ?? '',
    monthlyContribution: (json['monthly_contribution'] as num?)?.toDouble() ?? 0,
    monthsToTarget: (json['months_to_target'] as num?)?.toInt() ?? 0,
  );
}

class GoalProjectionResponse {
  final bool success; final GoalProjectionData? data;
  GoalProjectionResponse({required this.success, this.data});
  factory GoalProjectionResponse.fromJson(Map<String, dynamic> json) => GoalProjectionResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? GoalProjectionData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class GoalProjectionData {
  final String goalId; final String? projectedDate; final double projectedValue;
  final bool onTrack; final double monthlyContribution;
  GoalProjectionData({required this.goalId, this.projectedDate, required this.projectedValue, required this.onTrack, required this.monthlyContribution});
  factory GoalProjectionData.fromJson(Map<String, dynamic> json) => GoalProjectionData(
    goalId: json['goal_id'] as String? ?? '',
    projectedDate: json['projected_date'] as String?,
    projectedValue: (json['projected_value'] as num?)?.toDouble() ?? 0,
    onTrack: json['on_track'] as bool? ?? false,
    monthlyContribution: (json['monthly_contribution'] as num?)?.toDouble() ?? 0,
  );
}

class GoalRecsResponse {
  final bool success; final GoalRecsData? data;
  GoalRecsResponse({required this.success, this.data});
  factory GoalRecsResponse.fromJson(Map<String, dynamic> json) => GoalRecsResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? GoalRecsData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class GoalRecsData {
  final String goalId; final List<RecItem> recs; final int totalCount;
  GoalRecsData({required this.goalId, required this.recs, required this.totalCount});
  factory GoalRecsData.fromJson(Map<String, dynamic> json) => GoalRecsData(
    goalId: json['goal_id'] as String? ?? '',
    recs: (json['recommendations'] as List?)?.map((e) => RecItem.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    totalCount: (json['total_count'] as num?)?.toInt() ?? 0,
  );
}

class RecItem {
  final String recId; final String title; final String summary; final int priority; final String impact;
  RecItem({required this.recId, required this.title, required this.summary, required this.priority, required this.impact});
  factory RecItem.fromJson(Map<String, dynamic> json) => RecItem(
    recId: json['rec_id'] as String? ?? '', title: json['title'] as String? ?? '',
    summary: json['summary'] as String? ?? '', priority: (json['priority'] as num?)?.toInt() ?? 0,
    impact: json['impact'] as String? ?? '',
  );
}

class GoalOptsResponse {
  final bool success; final GoalOptsData? data;
  GoalOptsResponse({required this.success, this.data});
  factory GoalOptsResponse.fromJson(Map<String, dynamic> json) => GoalOptsResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? GoalOptsData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class GoalOptsData {
  final String goalId; final List<OptItem> opts; final int totalCount;
  GoalOptsData({required this.goalId, required this.opts, required this.totalCount});
  factory GoalOptsData.fromJson(Map<String, dynamic> json) => GoalOptsData(
    goalId: json['goal_id'] as String? ?? '',
    opts: (json['optimizations'] as List?)?.map((e) => OptItem.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    totalCount: (json['total_count'] as num?)?.toInt() ?? 0,
  );
}

class OptItem {
  final String optId; final String strategy; final double score; final String summary;
  OptItem({required this.optId, required this.strategy, required this.score, required this.summary});
  factory OptItem.fromJson(Map<String, dynamic> json) => OptItem(
    optId: json['opt_id'] as String? ?? '', strategy: json['strategy'] as String? ?? '',
    score: (json['score'] as num?)?.toDouble() ?? 0, summary: json['summary'] as String? ?? '',
  );
}

class GoalTimelineResponse {
  final bool success; final GoalTimelineData? data;
  GoalTimelineResponse({required this.success, this.data});
  factory GoalTimelineResponse.fromJson(Map<String, dynamic> json) => GoalTimelineResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? GoalTimelineData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class GoalTimelineData {
  final String goalId; final List<GoalEventItem> events; final int totalCount;
  GoalTimelineData({required this.goalId, required this.events, required this.totalCount});
  factory GoalTimelineData.fromJson(Map<String, dynamic> json) => GoalTimelineData(
    goalId: json['goal_id'] as String? ?? '',
    events: (json['events'] as List?)?.map((e) => GoalEventItem.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    totalCount: (json['total_count'] as num?)?.toInt() ?? 0,
  );
}

class GoalEventItem {
  final String eventId; final String eventType; final String title; final String timestamp;
  GoalEventItem({required this.eventId, required this.eventType, required this.title, required this.timestamp});
  factory GoalEventItem.fromJson(Map<String, dynamic> json) => GoalEventItem(
    eventId: json['event_id'] as String? ?? '', eventType: json['event_type'] as String? ?? '',
    title: json['title'] as String? ?? '', timestamp: json['timestamp'] as String? ?? '',
  );
}

class GoalMilestoneResponse {
  final bool success; final GoalMilestoneData? data;
  GoalMilestoneResponse({required this.success, this.data});
  factory GoalMilestoneResponse.fromJson(Map<String, dynamic> json) => GoalMilestoneResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? GoalMilestoneData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class GoalMilestoneData {
  final String goalId; final List<MilestoneItem> milestones; final int totalCount;
  GoalMilestoneData({required this.goalId, required this.milestones, required this.totalCount});
  factory GoalMilestoneData.fromJson(Map<String, dynamic> json) => GoalMilestoneData(
    goalId: json['goal_id'] as String? ?? '',
    milestones: (json['milestones'] as List?)?.map((e) => MilestoneItem.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    totalCount: (json['total_count'] as num?)?.toInt() ?? 0,
  );
}

class MilestoneItem {
  final String milestoneId; final String title; final int progressPct; final bool reached;
  MilestoneItem({required this.milestoneId, required this.title, required this.progressPct, required this.reached});
  factory MilestoneItem.fromJson(Map<String, dynamic> json) => MilestoneItem(
    milestoneId: json['milestone_id'] as String? ?? '', title: json['title'] as String? ?? '',
    progressPct: (json['progress_pct'] as num?)?.toInt() ?? 0,
    reached: json['reached'] as bool? ?? false,
  );
}
