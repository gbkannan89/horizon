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

class LocalGoldAsset {
  final int id;
  final int carat;
  final double grams;
  final double purchasePricePerGram;
  final String purchaseDate;
  final double? currentValue;
  final double? totalPurchaseCost;
  final double? returnPct;
  final String? notes;

  LocalGoldAsset({
    required this.id,
    required this.carat,
    required this.grams,
    required this.purchasePricePerGram,
    required this.purchaseDate,
    this.currentValue,
    this.totalPurchaseCost,
    this.returnPct,
    this.notes,
  });

  factory LocalGoldAsset.fromJson(Map<String, dynamic> json) {
    return LocalGoldAsset(
      id: json['id'],
      carat: json['carat'],
      grams: (json['grams'] ?? 0).toDouble(),
      purchasePricePerGram: (json['purchase_price_per_gram'] ?? 0).toDouble(),
      purchaseDate: json['purchase_date'] ?? '',
      currentValue: (json['current_value'] as num?)?.toDouble(),
      totalPurchaseCost: (json['total_purchase_cost'] as num?)?.toDouble(),
      returnPct: (json['return_pct'] as num?)?.toDouble(),
      notes: json['notes'],
    );
  }

  String get caratLabel => '${carat}K';
  double? get gramsValue => grams;
}

class LocalStockHolding {
  final int id;
  final String ticker;
  final String exchange;
  final String? name;
  final int quantity;
  final double avgPurchasePrice;
  final String purchaseDate;
  final double totalInvested;
  final double? currentValue;
  final double? returnPct;
  final double? lastCurrentPrice;

  LocalStockHolding({
    required this.id,
    required this.ticker,
    this.exchange = 'NSE',
    this.name,
    required this.quantity,
    required this.avgPurchasePrice,
    required this.purchaseDate,
    required this.totalInvested,
    this.currentValue,
    this.returnPct,
    this.lastCurrentPrice,
  });

  factory LocalStockHolding.fromJson(Map<String, dynamic> json) {
    return LocalStockHolding(
      id: json['id'],
      ticker: json['ticker'],
      exchange: json['exchange'] ?? 'NSE',
      name: json['name'],
      quantity: json['quantity'] ?? 0,
      avgPurchasePrice: (json['avg_purchase_price'] ?? 0).toDouble(),
      purchaseDate: json['purchase_date'] ?? '',
      totalInvested: (json['total_invested'] ?? 0).toDouble(),
      currentValue: (json['current_value'] as num?)?.toDouble(),
      returnPct: (json['return_pct'] as num?)?.toDouble(),
      lastCurrentPrice: (json['last_current_price'] as num?)?.toDouble(),
    );
  }
}

class LocalPFAsset {
  final int id;
  final String startDate;
  final int retirementAge;
  final double currentBalance;
  final double monthlyContribution;
  final double? employerContribution;
  final double interestRate;
  final double? projectedCorpus;
  final int? yearsUntilRetirement;
  final String? notes;

  LocalPFAsset({
    required this.id,
    required this.startDate,
    required this.retirementAge,
    required this.currentBalance,
    required this.monthlyContribution,
    this.employerContribution,
    required this.interestRate,
    this.projectedCorpus,
    this.yearsUntilRetirement,
    this.notes,
  });

  factory LocalPFAsset.fromJson(Map<String, dynamic> json) {
    return LocalPFAsset(
      id: json['id'],
      startDate: json['start_date'] ?? '',
      retirementAge: json['retirement_age'] ?? 58,
      currentBalance: (json['current_balance'] ?? 0).toDouble(),
      monthlyContribution: (json['monthly_contribution'] ?? 0).toDouble(),
      employerContribution: (json['employer_contribution'] as num?)?.toDouble(),
      interestRate: (json['interest_rate'] ?? 8.25).toDouble(),
      projectedCorpus: (json['projected_corpus'] as num?)?.toDouble(),
      yearsUntilRetirement: json['years_until_retirement'],
      notes: json['notes'],
    );
  }
}

class LocalLendingRecord {
  final int id;
  final String direction;
  final String personName;
  final double amount;
  final String dateGiven;
  final String? promisedReturnDate;
  final String? actualReturnDate;
  final double returnedAmount;
  final double remaining;
  final String status;
  final String? notes;

  LocalLendingRecord({
    required this.id,
    required this.direction,
    required this.personName,
    required this.amount,
    required this.dateGiven,
    this.promisedReturnDate,
    this.actualReturnDate,
    this.returnedAmount = 0,
    this.remaining = 0,
    this.status = 'open',
    this.notes,
  });

  factory LocalLendingRecord.fromJson(Map<String, dynamic> json) {
    return LocalLendingRecord(
      id: json['id'],
      direction: json['direction'] ?? 'lent',
      personName: json['person_name'] ?? '',
      amount: (json['amount'] ?? 0).toDouble(),
      dateGiven: json['date_given'] ?? '',
      promisedReturnDate: json['promised_return_date'],
      actualReturnDate: json['actual_return_date'],
      returnedAmount: (json['returned_amount'] ?? 0).toDouble(),
      remaining: (json['remaining'] ?? 0).toDouble(),
      status: json['status'] ?? 'open',
      notes: json['notes'],
    );
  }

  bool get isOverdue => status == 'overdue';
  bool get isClosed => status == 'closed';
  bool get isPartial => status == 'partial';
  bool get isLent => direction == 'lent';
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
