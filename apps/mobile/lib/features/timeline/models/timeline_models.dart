class TimelineResponse {
  final bool success;
  final TimelineOutput? data;
  final Map<String, dynamic>? metadata;

  TimelineResponse({required this.success, this.data, this.metadata});

  factory TimelineResponse.fromJson(Map<String, dynamic> json) => TimelineResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? TimelineOutput.fromJson(json['data'] as Map<String, dynamic>) : null,
    metadata: json['metadata'] as Map<String, dynamic>?,
  );
}

class TimelineOutput {
  final String view;
  final TimelineFeed feed;
  final Map<String, dynamic>? filters;
  final int total;

  TimelineOutput({required this.view, required this.feed, this.filters, required this.total});

  factory TimelineOutput.fromJson(Map<String, dynamic> json) => TimelineOutput(
    view: json['view'] as String? ?? 'Today',
    feed: TimelineFeed.fromJson(json['feed'] as Map<String, dynamic>? ?? {}),
    filters: json['filters'] as Map<String, dynamic>?,
    total: (json['total'] as num?)?.toInt() ?? 0,
  );
}

class TimelineFeed {
  final List<TimelineItem> items;
  final String? cursor;
  final bool hasMore;
  final int total;

  TimelineFeed({required this.items, this.cursor, required this.hasMore, required this.total});

  factory TimelineFeed.fromJson(Map<String, dynamic> json) => TimelineFeed(
    items: (json['items'] as List?)?.map((e) => TimelineItem.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    cursor: json['cursor'] as String?,
    hasMore: json['has_more'] as bool? ?? false,
    total: (json['total'] as num?)?.toInt() ?? 0,
  );
}

class TimelineItem {
  final String timelineId;
  final String timestamp;
  final String eventType;
  final String category;
  final String title;
  final String summary;
  final String description;
  final String severity;
  final String? relatedEntity;
  final String? relatedAgg;
  final String? relatedGoal;
  final String? relatedAccount;
  final String? relatedAsset;
  final int amount;
  final Map<String, dynamic>? metadata;
  final String icon;
  final String color;

  TimelineItem({
    required this.timelineId, required this.timestamp, required this.eventType,
    required this.category, required this.title, required this.summary,
    this.description = '', required this.severity, this.relatedEntity,
    this.relatedAgg, this.relatedGoal, this.relatedAccount, this.relatedAsset,
    this.amount = 0, this.metadata, this.icon = 'circle', this.color = 'grey',
  });

  factory TimelineItem.fromJson(Map<String, dynamic> json) => TimelineItem(
    timelineId: json['timeline_id'] as String? ?? '',
    timestamp: json['timestamp'] as String? ?? '',
    eventType: json['event_type'] as String? ?? '',
    category: json['category'] as String? ?? '',
    title: json['title'] as String? ?? '',
    summary: json['summary'] as String? ?? '',
    description: json['description'] as String? ?? '',
    severity: json['severity'] as String? ?? 'info',
    relatedEntity: json['related_entity'] as String?,
    relatedAgg: json['related_aggregate'] as String?,
    relatedGoal: json['related_goal'] as String?,
    relatedAccount: json['related_account'] as String?,
    relatedAsset: json['related_asset'] as String?,
    amount: (json['amount'] as num?)?.toInt() ?? 0,
    metadata: json['metadata'] as Map<String, dynamic>?,
    icon: json['icon'] as String? ?? 'circle',
    color: json['color'] as String? ?? 'grey',
  );

  String get formattedDate {
    try {
      final dt = DateTime.parse(timestamp);
      final now = DateTime.now();
      final diff = now.difference(dt);
      if (diff.inMinutes < 1) return 'Just now';
      if (diff.inHours < 1) return '${diff.inMinutes}m ago';
      if (diff.inDays < 1) return '${diff.inHours}h ago';
      if (diff.inDays == 1) return 'Yesterday';
      if (diff.inDays < 7) return '${diff.inDays}d ago';
      return '${dt.day}/${dt.month}/${dt.year}';
    } catch (_) {
      return timestamp;
    }
  }
}

class FilterParams {
  final List<String>? categories;
  final String? startDate;
  final String? endDate;
  final String? severity;
  final String? goal;
  final String? account;
  final String? asset;
  final String? liability;
  final String? portfolio;
  final String? q;

  const FilterParams({
    this.categories, this.startDate, this.endDate, this.severity,
    this.goal, this.account, this.asset, this.liability, this.portfolio, this.q,
  });

  Map<String, dynamic> toQuery() {
    final map = <String, dynamic>{};
    if (categories != null && categories!.isNotEmpty) map['categories'] = categories!.join(',');
    if (startDate != null) map['start_date'] = startDate;
    if (endDate != null) map['end_date'] = endDate;
    if (severity != null) map['severity'] = severity;
    if (goal != null) map['goal'] = goal;
    if (account != null) map['account'] = account;
    if (asset != null) map['asset'] = asset;
    if (liability != null) map['liability'] = liability;
    if (portfolio != null) map['portfolio'] = portfolio;
    if (q != null) map['q'] = q;
    return map;
  }
}
