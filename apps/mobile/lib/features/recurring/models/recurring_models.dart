class RecurringTransaction {
  final String id;
  final String name;
  final String description;
  final double amount;
  final String currency;
  final String frequency;
  final int interval;
  final String startDate;
  final String? endDate;
  final String? nextOccurrence;
  final String status;

  RecurringTransaction({
    required this.id,
    required this.name,
    this.description = '',
    required this.amount,
    required this.currency,
    required this.frequency,
    required this.interval,
    required this.startDate,
    this.endDate,
    this.nextOccurrence,
    required this.status,
  });

  factory RecurringTransaction.fromJson(Map<String, dynamic> json) {
    return RecurringTransaction(
      id: json['id'] as String? ?? '',
      name: json['name'] as String? ?? '',
      description: json['description'] as String? ?? '',
      amount: (json['amount'] as num?)?.toDouble() ?? 0.0,
      currency: json['currency'] as String? ?? 'INR',
      frequency: json['frequency'] as String? ?? 'Monthly',
      interval: (json['interval'] as num?)?.toInt() ?? 1,
      startDate: json['start_date'] as String? ?? '',
      endDate: json['end_date'] as String?,
      nextOccurrence: json['next_occurrence'] as String?,
      status: json['status'] as String? ?? 'Active',
    );
  }

  Map<String, dynamic> toJson() => {
    'id': id,
    'name': name,
    'description': description,
    'amount': amount.toInt(),
    'currency': currency,
    'frequency': frequency,
    'interval': interval,
    'start_date': startDate,
    if (endDate != null) 'end_date': endDate,
    'status': status,
  };
}

class RecurringListResponse {
  final bool success;
  final List<RecurringTransaction>? data;
  RecurringListResponse({required this.success, this.data});
  factory RecurringListResponse.fromJson(Map<String, dynamic> json) => RecurringListResponse(
    success: json['success'] as bool? ?? false,
    data: (json['data'] as List?)?.map((e) => RecurringTransaction.fromJson(e as Map<String, dynamic>)).toList(),
  );
}

class RecurringDetailResponse {
  final bool success;
  final RecurringTransaction? data;
  RecurringDetailResponse({required this.success, this.data});
  factory RecurringDetailResponse.fromJson(Map<String, dynamic> json) => RecurringDetailResponse(
    success: json['success'] as bool? ?? false,
    data: json['data'] != null ? RecurringTransaction.fromJson(json['data'] as Map<String, dynamic>) : null,
  );
}

class RecurringActionResponse {
  final String? status;
  RecurringActionResponse({this.status});
  factory RecurringActionResponse.fromJson(Map<String, dynamic> json) => RecurringActionResponse(
    status: json['status'] as String?,
  );
}
