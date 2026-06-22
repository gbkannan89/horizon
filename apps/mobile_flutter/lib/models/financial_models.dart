class LocalIncome {
  final String id;
  final String label;  // user-friendly name, e.g. "My Salary", "Partner Income"
  final String type;   // backend type: 'salary', 'business', 'passive'
  final double amount;
  final String frequency;

  LocalIncome({
    required this.id,
    required this.label,
    required this.type,
    required this.amount,
    this.frequency = 'monthly',
  });

  // Backwards-compat getter used in older code
  String get name => label;
}

class LocalAsset {
  final String id;
  final String name;
  final double amount;

  LocalAsset({required this.id, required this.name, required this.amount});
}

class LocalLiability {
  final String id;
  final String name;
  final double amount;
  final double interestRate;

  LocalLiability({required this.id, required this.name, required this.amount, this.interestRate = 0.0});
}

class LocalBill {
  final String id;
  final String name;
  final double amount;
  final DateTime dueDate;

  LocalBill({required this.id, required this.name, required this.amount, required this.dueDate});
}

class LocalHouseholdMember {
  final String id;
  final String name;
  final double monthlyIncome;
  final double contribution;
  final String relationship;

  LocalHouseholdMember({
    required this.id,
    required this.name,
    required this.monthlyIncome,
    required this.contribution,
    required this.relationship,
  });
}
