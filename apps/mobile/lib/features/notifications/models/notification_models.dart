class NotificationItem {
  final String id;
  final String title;
  final String? body;
  final String category;
  final String priority;
  final bool isRead;
  final String createdAt;

  NotificationItem({
    required this.id, required this.title, this.body,
    this.category = '', this.priority = 'P3', this.isRead = false,
    this.createdAt = '',
  });

  factory NotificationItem.fromJson(Map<String, dynamic> json) => NotificationItem(
    id: json['id'] as String? ?? json['notification_id'] as String? ?? '',
    title: json['title'] as String? ?? '',
    body: json['body'] as String? ?? json['message'] as String?,
    category: json['category'] as String? ?? '',
    priority: json['priority'] as String? ?? 'P3',
    isRead: json['is_read'] as bool? ?? json['read'] as bool? ?? false,
    createdAt: json['created_at'] as String? ?? '',
  );
}

class NotificationPreferences {
  final Map<String, NotificationCategoryPref> categories;

  NotificationPreferences({this.categories = const {}});

  factory NotificationPreferences.fromJson(Map<String, dynamic> json) {
    final prefs = json['preferences'] as List? ?? json['data'] as List? ?? [];
    final map = <String, NotificationCategoryPref>{};
    for (final p in prefs) {
      final pm = p as Map<String, dynamic>;
      final cat = NotificationCategoryPref.fromJson(pm);
      map[cat.category] = cat;
    }
    return NotificationPreferences(categories: map);
  }

  Map<String, dynamic> toJson() => {
    'preferences': categories.values.map((c) => c.toJson()).toList(),
  };
}

class NotificationCategoryPref {
  final String category;
  final bool inApp;
  final bool push;
  final bool email;
  final String digest;
  final bool enabled;

  NotificationCategoryPref({
    required this.category, this.inApp = true, this.push = true,
    this.email = false, this.digest = 'daily', this.enabled = true,
  });

  factory NotificationCategoryPref.fromJson(Map<String, dynamic> json) => NotificationCategoryPref(
    category: json['category'] as String? ?? '',
    inApp: json['in_app'] as bool? ?? true,
    push: json['push'] as bool? ?? true,
    email: json['email'] as bool? ?? false,
    digest: json['digest'] as String? ?? 'daily',
    enabled: json['enabled'] as bool? ?? true,
  );

  Map<String, dynamic> toJson() => {
    'category': category, 'in_app': inApp, 'push': push,
    'email': email, 'digest': digest, 'enabled': enabled,
  };
}
