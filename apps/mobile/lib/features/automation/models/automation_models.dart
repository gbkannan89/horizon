class RuleCondition {
  final String field;
  final String operator;
  final String value;

  RuleCondition({required this.field, required this.operator, required this.value});

  factory RuleCondition.fromJson(Map<String, dynamic> json) => RuleCondition(
    field: json['field'] as String? ?? '',
    operator: json['operator'] as String? ?? '',
    value: json['value'] as String? ?? '',
  );

  Map<String, dynamic> toJson() => {'field': field, 'operator': operator, 'value': value};
}

class RuleAction {
  final String type;
  final Map<String, dynamic> params;

  RuleAction({required this.type, this.params = const {}});

  factory RuleAction.fromJson(Map<String, dynamic> json) => RuleAction(
    type: json['type'] as String? ?? '',
    params: json['params'] as Map<String, dynamic>? ?? {},
  );

  Map<String, dynamic> toJson() => {'type': type, 'params': params};
}

class RuleModel {
  final String ruleId;
  final String name;
  final String description;
  final String category;
  final int priority;
  final bool enabled;
  final List<RuleCondition> conditions;
  final List<RuleAction> actions;
  final String createdAt;
  final String updatedAt;

  RuleModel({
    required this.ruleId,
    required this.name,
    this.description = '',
    required this.category,
    required this.priority,
    required this.enabled,
    required this.conditions,
    required this.actions,
    required this.createdAt,
    required this.updatedAt,
  });

  factory RuleModel.fromJson(Map<String, dynamic> json) => RuleModel(
    ruleId: json['rule_id'] as String? ?? json['id'] as String? ?? '',
    name: json['name'] as String? ?? '',
    description: json['description'] as String? ?? '',
    category: json['category'] as String? ?? '',
    priority: (json['priority'] as num?)?.toInt() ?? 0,
    enabled: json['enabled'] as bool? ?? false,
    conditions: (json['conditions'] as List?)?.map((e) => RuleCondition.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    actions: (json['actions'] as List?)?.map((e) => RuleAction.fromJson(e as Map<String, dynamic>)).toList() ?? [],
    createdAt: json['created_at'] as String? ?? '',
    updatedAt: json['updated_at'] as String? ?? '',
  );

  Map<String, dynamic> toJson() => {
    'name': name,
    'description': description,
    'category': category,
    'priority': priority,
    'enabled': enabled,
    'conditions': conditions.map((e) => e.toJson()).toList(),
    'actions': actions.map((e) => e.toJson()).toList(),
  };
}

class RuleListResponse {
  final bool success;
  final List<RuleModel>? data;
  RuleListResponse({required this.success, this.data});
  factory RuleListResponse.fromJson(Map<String, dynamic> json) => RuleListResponse(
    success: json['success'] as bool? ?? false,
    data: (json['data'] as List?)?.map((e) => RuleModel.fromJson(e as Map<String, dynamic>)).toList(),
  );
}

class RuleDetailResponse {
  final bool success;
  final RuleModel? data;
  RuleDetailResponse({required this.success, this.data});
  factory RuleDetailResponse.fromJson(Map<String, dynamic> json) => RuleDetailResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? RuleModel.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}
