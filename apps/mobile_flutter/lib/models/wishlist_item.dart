class WishlistItem {
  final int id;
  final int userId;
  final String name;
  final double amount;
  final DateTime addedDate;
  final DateTime unlockDate;
  final String status;

  WishlistItem({
    required this.id,
    required this.userId,
    required this.name,
    required this.amount,
    required this.addedDate,
    required this.unlockDate,
    required this.status,
  });

  factory WishlistItem.fromJson(Map<String, dynamic> json) {
    return WishlistItem(
      id: json['id'],
      userId: json['user_id'],
      name: json['name'],
      amount: (json['amount'] as num).toDouble(),
      addedDate: DateTime.parse(json['added_date']),
      unlockDate: DateTime.parse(json['unlock_date']),
      status: json['status'],
    );
  }

  bool get isUnlocked => DateTime.now().isAfter(unlockDate);
}
