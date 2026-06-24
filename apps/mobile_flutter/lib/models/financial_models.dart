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
  final double interestRate;
  final bool isLiability;
  final bool generatesIncome;
  final double? purchasePrice;
  final String? purchaseDate;

  LocalAsset({
    required this.id, 
    required this.name, 
    required this.amount,
    this.interestRate = 0.0,
    this.isLiability = false,
    this.generatesIncome = false,
    this.purchasePrice,
    this.purchaseDate,
  });
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
  final String frequency;
  final DateTime dueDate;
  final bool isEmi;
  final int emiTotalMonths;
  final int emiMonthsPaid;

  LocalBill({
    required this.id, 
    required this.name, 
    required this.amount, 
    this.frequency = 'monthly',
    required this.dueDate,
    this.isEmi = false,
    this.emiTotalMonths = 0,
    this.emiMonthsPaid = 0,
  });
}

class LocalVehicle {
  final String id;
  final String makeModel;
  final double purchaseCost;
  final DateTime? insuranceRenewalDate;

  LocalVehicle({
    required this.id,
    required this.makeModel,
    required this.purchaseCost,
    this.insuranceRenewalDate,
  });
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
