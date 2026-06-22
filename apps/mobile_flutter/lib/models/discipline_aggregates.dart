class FinancialGuardrails {
  final double emergencyFundRatio;
  final double housingCostRatio;
  final String housingStatus;
  final String runwayStatus;

  FinancialGuardrails({
    required this.emergencyFundRatio,
    required this.housingCostRatio,
    required this.housingStatus,
    required this.runwayStatus,
  });

  factory FinancialGuardrails.fromJson(Map<String, dynamic> json) {
    return FinancialGuardrails(
      emergencyFundRatio: (json['emergency_fund_ratio'] as num).toDouble(),
      housingCostRatio: (json['housing_cost_ratio'] as num).toDouble(),
      housingStatus: json['housing_status'],
      runwayStatus: json['runway_status'],
    );
  }
}

class DebtRepaymentStrategy {
  final int snowballMonths;
  final int avalancheMonths;
  final double snowballInterest;
  final double avalancheInterest;
  final String recommendedStrategy;

  DebtRepaymentStrategy({
    required this.snowballMonths,
    required this.avalancheMonths,
    required this.snowballInterest,
    required this.avalancheInterest,
    required this.recommendedStrategy,
  });

  factory DebtRepaymentStrategy.fromJson(Map<String, dynamic> json) {
    return DebtRepaymentStrategy(
      snowballMonths: json['snowball_months_to_freedom'],
      avalancheMonths: json['avalanche_months_to_freedom'],
      snowballInterest: (json['snowball_total_interest'] as num).toDouble(),
      avalancheInterest: (json['avalanche_total_interest'] as num).toDouble(),
      recommendedStrategy: json['recommended_strategy'],
    );
  }
}

class ZeroBasedBudget {
  final double totalIncome;
  final double totalAllocated;
  final double unallocated;
  final String status;

  ZeroBasedBudget({
    required this.totalIncome,
    required this.totalAllocated,
    required this.unallocated,
    required this.status,
  });

  factory ZeroBasedBudget.fromJson(Map<String, dynamic> json) {
    return ZeroBasedBudget(
      totalIncome: (json['total_income'] as num).toDouble(),
      totalAllocated: (json['total_allocated'] as num).toDouble(),
      unallocated: (json['unallocated'] as num).toDouble(),
      status: json['status'],
    );
  }
}

class PayYourselfFirst {
  final double targetPercentage;
  final double targetAmount;
  final double actualSavings;
  final String status;

  PayYourselfFirst({
    required this.targetPercentage,
    required this.targetAmount,
    required this.actualSavings,
    required this.status,
  });

  factory PayYourselfFirst.fromJson(Map<String, dynamic> json) {
    return PayYourselfFirst(
      targetPercentage: (json['target_percentage'] as num).toDouble(),
      targetAmount: (json['target_amount'] as num).toDouble(),
      actualSavings: (json['actual_savings'] as num).toDouble(),
      status: json['status'],
    );
  }
}
