import 'package:flutter/material.dart';

class InsightListResponse {
  final bool success;
  final InsightListData? data;
  InsightListResponse({required this.success, this.data});
  factory InsightListResponse.fromJson(Map<String, dynamic> json) => InsightListResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? InsightListData.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class InsightListData {
  final List<InsightItem> insights;
  final int total;
  InsightListData({this.insights = const [], this.total = 0});
  factory InsightListData.fromJson(Map<String, dynamic> json) => InsightListData(
    insights: (json['insights'] as List?)?.map((e) => InsightItem.fromJson(e as Map<String, dynamic>)).toList() ?? (json['data'] as List?)?.map((e) => InsightItem.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    total: (json['total'] as num?)?.toInt() ?? (json['count'] as num?)?.toInt() ?? 0,
  );
}

class InsightDashboardResponse {
  final bool success;
  final InsightDashboard? data;
  InsightDashboardResponse({required this.success, this.data});
  factory InsightDashboardResponse.fromJson(Map<String, dynamic> json) => InsightDashboardResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? InsightDashboard.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class InsightDashboard {
  final int total;
  final int criticalCount;
  final int newCount;
  InsightDashboard({this.total = 0, this.criticalCount = 0, this.newCount = 0});
  factory InsightDashboard.fromJson(Map<String, dynamic> json) => InsightDashboard(
    total: (json['total'] as num?)?.toInt() ?? (json['total_count'] as num?)?.toInt() ?? 0,
    criticalCount: (json['critical_count'] as num?)?.toInt() ?? (json['critical'] as num?)?.toInt() ?? 0,
    newCount: (json['new_count'] as num?)?.toInt() ?? 0,
  );
}

class InsightItem {
  final String id;
  final String category;
  final String title;
  final String summary;
  final String? description;
  final String priority;
  final double? impact;
  final double? confidence;
  final String? createdAt;

  InsightItem({
    required this.id,
    this.category = '',
    required this.title,
    this.summary = '',
    this.description,
    this.priority = 'P3',
    this.impact,
    this.confidence,
    this.createdAt,
  });

  factory InsightItem.fromJson(Map<String, dynamic> json) => InsightItem(
    id: json['id'] as String? ?? json['insight_id'] as String? ?? '',
    category: json['category'] as String? ?? '',
    title: json['title'] as String? ?? '',
    summary: json['summary'] as String? ?? '',
    description: json['description'] as String?,
    priority: json['priority'] as String? ?? 'P3',
    impact: (json['impact'] as num?)?.toDouble(),
    confidence: (json['confidence'] as num?)?.toDouble(),
    createdAt: json['created_at'] as String?,
  );

  Color get priorityColor {
    switch (priority) {
      case 'P1': return Colors.red;
      case 'P2': return Colors.orange;
      case 'P3': return Colors.amber;
      case 'P4': return Colors.grey;
      case 'P5': return Colors.blue;
      default: return Colors.grey;
    }
  }

  IconData get categoryIcon {
    switch (category.toLowerCase()) {
      case 'spending': return Icons.shopping_cart;
      case 'income': return Icons.trending_up;
      case 'goal': return Icons.flag;
      case 'portfolio': return Icons.pie_chart;
      case 'risk': return Icons.warning;
      case 'opportunity': return Icons.lightbulb;
      case 'warning': return Icons.error_outline;
      case 'achievement': return Icons.emoji_events;
      case 'forecast': return Icons.query_stats;
      case 'trend': return Icons.show_chart;
      default: return Icons.insights;
    }
  }
}

class TrendData {
  final List<double> values;
  final List<String> labels;
  final double changePercent;
  TrendData({this.values = const [], this.labels = const [], this.changePercent = 0});
  factory TrendData.fromJson(Map<String, dynamic> json) => TrendData(
    values: (json['values'] as List?)?.map((e) => (e as num).toDouble()).toList() ?? [],
    labels: (json['labels'] as List?)?.map((e) => e.toString()).toList() ?? [],
    changePercent: (json['change_percent'] as num?)?.toDouble() ?? 0,
  );
}
