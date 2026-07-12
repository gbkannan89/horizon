class LocalIncome {
  final String id;
  final String label;  // user-friendly name, e.g. "My Salary", "Partner Income"
  final String type;   // backend type: 'salary', 'business', 'passive'
  final double amount;
  final String frequency;
  final String? companyName;
  final String? notes;

  LocalIncome({
    required this.id,
    required this.label,
    required this.type,
    required this.amount,
    this.frequency = 'monthly',
    this.companyName,
    this.notes,
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
  final int? modelYear;
  final int? purchaseYear;
  final String? fuelType;
  final double? mileageKmpl;
  final double fuelCostTotal;
  final double kmDriven;
  final double? insuranceIdv;
  final double? insuranceRenewalAmount;
  final String? registrationNumber;
  // Computed fields
  final double costPerKm;
  final double suggestedIdv;
  final double totalServiceCost;
  final Map<String, dynamic>? nextService;

  LocalVehicle({
    required this.id,
    required this.makeModel,
    required this.purchaseCost,
    this.insuranceRenewalDate,
    this.modelYear,
    this.purchaseYear,
    this.fuelType,
    this.mileageKmpl,
    this.fuelCostTotal = 0.0,
    this.kmDriven = 0.0,
    this.insuranceIdv,
    this.insuranceRenewalAmount,
    this.registrationNumber,
    this.costPerKm = 0.0,
    this.suggestedIdv = 0.0,
    this.totalServiceCost = 0.0,
    this.nextService,
  });

  factory LocalVehicle.fromJson(Map<String, dynamic> json) {
    return LocalVehicle(
      id: json['id'].toString(),
      makeModel: json['make_model'] ?? '',
      purchaseCost: (json['purchase_cost'] ?? 0).toDouble(),
      insuranceRenewalDate: json['insurance_renewal_date'] != null
          ? DateTime.tryParse(json['insurance_renewal_date'])
          : null,
      modelYear: json['model_year'],
      purchaseYear: json['purchase_year'],
      fuelType: json['fuel_type'],
      mileageKmpl: json['mileage_kmpl'] != null ? (json['mileage_kmpl'] as num).toDouble() : null,
      fuelCostTotal: (json['fuel_cost_total'] ?? 0).toDouble(),
      kmDriven: (json['km_driven'] ?? 0).toDouble(),
      insuranceIdv: json['insurance_idv'] != null ? (json['insurance_idv'] as num).toDouble() : null,
      insuranceRenewalAmount: json['insurance_renewal_amount'] != null
          ? (json['insurance_renewal_amount'] as num).toDouble()
          : null,
      registrationNumber: json['registration_number'],
      costPerKm: (json['cost_per_km'] ?? 0).toDouble(),
      suggestedIdv: (json['suggested_idv'] ?? 0).toDouble(),
      totalServiceCost: (json['total_service_cost'] ?? 0).toDouble(),
      nextService: json['next_service'],
    );
  }
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

class LocalFamilyMember {
  final int id;
  final int householdId;
  final String name;
  final String? dob;
  final String? bloodGroup;
  final String relationship;
  final String? avatarColor;
  final bool isSelf;
  final bool isActive;
  final List<LocalSchooling> schooling;
  final List<LocalCheckup> checkups;
  final List<LocalMedicine> medicines;
  final List<LocalVaccination> vaccinations;
  final List<LocalEarnings> earnings;
  final List<LocalInsuranceLink> insuranceLinks;

  LocalFamilyMember({
    required this.id,
    required this.householdId,
    required this.name,
    this.dob,
    this.bloodGroup,
    required this.relationship,
    this.avatarColor,
    required this.isSelf,
    required this.isActive,
    required this.schooling,
    required this.checkups,
    required this.medicines,
    required this.vaccinations,
    required this.earnings,
    required this.insuranceLinks,
  });

  factory LocalFamilyMember.fromJson(Map<String, dynamic> json) {
    return LocalFamilyMember(
      id: json['id'],
      householdId: json['household_id'],
      name: json['name'] ?? '',
      dob: json['dob'],
      bloodGroup: json['blood_group'],
      relationship: json['relationship'] ?? 'other',
      avatarColor: json['avatar_color'],
      isSelf: json['is_self'] ?? false,
      isActive: json['is_active'] ?? true,
      schooling: (json['schooling'] as List?)
              ?.map((e) => LocalSchooling.fromJson(e))
              .toList() ??
          [],
      checkups: (json['checkups'] as List?)
              ?.map((e) => LocalCheckup.fromJson(e))
              .toList() ??
          [],
      medicines: (json['medicines'] as List?)
              ?.map((e) => LocalMedicine.fromJson(e))
              .toList() ??
          [],
      vaccinations: (json['vaccinations'] as List?)
              ?.map((e) => LocalVaccination.fromJson(e))
              .toList() ??
          [],
      earnings: json['earnings'] is Map
              ? [LocalEarnings.fromJson(json['earnings'] as Map<String, dynamic>)]
              : (json['earnings'] as List?)
                      ?.map((e) => LocalEarnings.fromJson(e))
                      .toList() ??
                  [],
      insuranceLinks: (json['insurance_links'] as List?)
              ?.map((e) => LocalInsuranceLink.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class LocalSchooling {
  final int id;
  final int memberId;
  final String institutionName;
  final double feeAmount;
  final String feeFrequency;
  final String? lastPaidDate;
  final String? nextDueDate;
  final String? notes;
  final double monthlyEquivalent;

  LocalSchooling({
    required this.id,
    required this.memberId,
    required this.institutionName,
    required this.feeAmount,
    required this.feeFrequency,
    this.lastPaidDate,
    this.nextDueDate,
    this.notes,
    required this.monthlyEquivalent,
  });

  factory LocalSchooling.fromJson(Map<String, dynamic> json) {
    return LocalSchooling(
      id: json['id'],
      memberId: json['member_id'],
      institutionName: json['institution_name'] ?? '',
      feeAmount: (json['fee_amount'] ?? 0).toDouble(),
      feeFrequency: json['fee_frequency'] ?? 'monthly',
      lastPaidDate: json['last_paid_date'],
      nextDueDate: json['next_due_date'],
      notes: json['notes'],
      monthlyEquivalent: (json['monthly_equivalent'] ?? 0).toDouble(),
    );
  }
}

class LocalSchoolingPayment {
  final int id;
  final int schoolingId;
  final double amount;
  final String paidDate;
  final String? receiptRef;
  final String? notes;

  LocalSchoolingPayment({
    required this.id,
    required this.schoolingId,
    required this.amount,
    required this.paidDate,
    this.receiptRef,
    this.notes,
  });

  factory LocalSchoolingPayment.fromJson(Map<String, dynamic> json) {
    return LocalSchoolingPayment(
      id: json['id'],
      schoolingId: json['schooling_id'],
      amount: (json['amount'] ?? 0).toDouble(),
      paidDate: json['paid_date'] ?? '',
      receiptRef: json['receipt_ref'],
      notes: json['notes'],
    );
  }
}

class LocalCheckup {
  final int id;
  final int memberId;
  final String? checkupType;
  final String frequency;
  final double recurringCost;
  final String? lastCheckupDate;
  final String? nextDueDate;
  final String? notes;
  final double monthlyEquivalent;

  LocalCheckup({
    required this.id,
    required this.memberId,
    this.checkupType,
    required this.frequency,
    required this.recurringCost,
    this.lastCheckupDate,
    this.nextDueDate,
    this.notes,
    required this.monthlyEquivalent,
  });

  factory LocalCheckup.fromJson(Map<String, dynamic> json) {
    return LocalCheckup(
      id: json['id'],
      memberId: json['member_id'],
      checkupType: json['checkup_type'],
      frequency: json['frequency'] ?? 'yearly',
      recurringCost: (json['recurring_cost'] ?? 0).toDouble(),
      lastCheckupDate: json['last_checkup_date'],
      nextDueDate: json['next_due_date'],
      notes: json['notes'],
      monthlyEquivalent: (json['monthly_equivalent'] ?? 0).toDouble(),
    );
  }
}

class LocalMedicine {
  final int id;
  final int memberId;
  final String medicineName;
  final double monthlyCost;
  final String? purpose;
  final bool isRegular;
  final String? prescribedBy;

  LocalMedicine({
    required this.id,
    required this.memberId,
    required this.medicineName,
    required this.monthlyCost,
    this.purpose,
    required this.isRegular,
    this.prescribedBy,
  });

  factory LocalMedicine.fromJson(Map<String, dynamic> json) {
    return LocalMedicine(
      id: json['id'],
      memberId: json['member_id'],
      medicineName: json['medicine_name'] ?? '',
      monthlyCost: (json['monthly_cost'] ?? 0).toDouble(),
      purpose: json['purpose'],
      isRegular: json['is_regular'] ?? true,
      prescribedBy: json['prescribed_by'],
    );
  }
}

class LocalVaccination {
  final int id;
  final int memberId;
  final String vaccineName;
  final String frequency;
  final double recurringCost;
  final String? lastVaccinationDate;
  final String? nextDueDate;
  final String? notes;
  final double monthlyEquivalent;

  LocalVaccination({
    required this.id,
    required this.memberId,
    required this.vaccineName,
    required this.frequency,
    required this.recurringCost,
    this.lastVaccinationDate,
    this.nextDueDate,
    this.notes,
    required this.monthlyEquivalent,
  });

  factory LocalVaccination.fromJson(Map<String, dynamic> json) {
    return LocalVaccination(
      id: json['id'],
      memberId: json['member_id'],
      vaccineName: json['vaccine_name'] ?? '',
      frequency: json['frequency'] ?? 'one_time',
      recurringCost: (json['recurring_cost'] ?? 0).toDouble(),
      lastVaccinationDate: json['last_vaccination_date'],
      nextDueDate: json['next_due_date'],
      notes: json['notes'],
      monthlyEquivalent: (json['monthly_equivalent'] ?? 0).toDouble(),
    );
  }
}

class LocalEarnings {
  final int id;
  final int memberId;
  final double monthlyIncome;
  final double contributionToHousehold;
  final String? occupation;

  LocalEarnings({
    required this.id,
    required this.memberId,
    required this.monthlyIncome,
    required this.contributionToHousehold,
    this.occupation,
  });

  factory LocalEarnings.fromJson(Map<String, dynamic> json) {
    return LocalEarnings(
      id: json['id'],
      memberId: json['member_id'],
      monthlyIncome: (json['monthly_income'] ?? 0).toDouble(),
      contributionToHousehold: (json['contribution_to_household'] ?? 0).toDouble(),
      occupation: json['occupation'],
    );
  }
}

class LocalInsuranceLink {
  final int id;
  final int memberId;
  final int insuranceId;
  final String relationship;
  final String? provider;
  final String? policyName;
  final String? type;

  LocalInsuranceLink({
    required this.id,
    required this.memberId,
    required this.insuranceId,
    required this.relationship,
    this.provider,
    this.policyName,
    this.type,
  });

  factory LocalInsuranceLink.fromJson(Map<String, dynamic> json) {
    return LocalInsuranceLink(
      id: json['id'],
      memberId: json['member_id'],
      insuranceId: json['insurance_id'],
      relationship: json['relationship'] ?? 'self',
      provider: json['provider'],
      policyName: json['policy_name'],
      type: json['type'],
    );
  }
}

class LocalBudgetItem {
  final int id;
  final int budgetPlanId;
  final String category;
  final String label;
  final double amount;
  final String frequency;
  final String bucket;
  final String source;
  final int? sourceId;
  final String? sourceLabel;
  final bool isCommitted;
  final int sortOrder;

  LocalBudgetItem({
    required this.id,
    required this.budgetPlanId,
    required this.category,
    required this.label,
    required this.amount,
    required this.frequency,
    required this.bucket,
    required this.source,
    this.sourceId,
    this.sourceLabel,
    required this.isCommitted,
    required this.sortOrder,
  });

  factory LocalBudgetItem.fromJson(Map<String, dynamic> json) {
    return LocalBudgetItem(
      id: json['id'],
      budgetPlanId: json['budget_plan_id'],
      category: json['category'] ?? 'other',
      label: json['label'] ?? '',
      amount: (json['amount'] ?? 0).toDouble(),
      frequency: json['frequency'] ?? 'monthly',
      bucket: json['bucket'] ?? 'Needs',
      source: json['source'] ?? 'manual',
      sourceId: json['source_id'],
      sourceLabel: json['source_label'],
      isCommitted: json['is_committed'] ?? true,
      sortOrder: json['sort_order'] ?? 0,
    );
  }

  bool get isAuto => source != 'manual';
}

class LocalBudgetPlan {
  final int id;
  final int userId;
  final int month;
  final int year;
  final double totalBudgeted;
  final String? notes;
  final List<LocalBudgetItem> items;

  LocalBudgetPlan({
    required this.id,
    required this.userId,
    required this.month,
    required this.year,
    required this.totalBudgeted,
    this.notes,
    required this.items,
  });

  factory LocalBudgetPlan.fromJson(Map<String, dynamic> json) {
    return LocalBudgetPlan(
      id: json['id'],
      userId: json['user_id'],
      month: json['month'],
      year: json['year'],
      totalBudgeted: (json['total_budgeted'] ?? 0).toDouble(),
      notes: json['notes'],
      items: (json['items'] as List?)
              ?.map((e) => LocalBudgetItem.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class LocalBudgetSummary {
  final double total;
  final double needs;
  final double wants;
  final double savings;

  LocalBudgetSummary({
    required this.total,
    required this.needs,
    required this.wants,
    required this.savings,
  });

  factory LocalBudgetSummary.fromJson(Map<String, dynamic> json) {
    return LocalBudgetSummary(
      total: (json['total'] ?? 0).toDouble(),
      needs: (json['needs'] ?? 0).toDouble(),
      wants: (json['wants'] ?? 0).toDouble(),
      savings: (json['savings'] ?? 0).toDouble(),
    );
  }
}

class LocalCompareItem {
  final String label;
  final double budgeted;
  final double actual;
  final double remaining;
  final String bucket;

  LocalCompareItem({
    required this.label,
    required this.budgeted,
    required this.actual,
    required this.remaining,
    required this.bucket,
  });

  factory LocalCompareItem.fromJson(Map<String, dynamic> json) {
    return LocalCompareItem(
      label: json['label'] ?? '',
      budgeted: (json['budgeted'] ?? 0).toDouble(),
      actual: (json['actual'] ?? 0).toDouble(),
      remaining: (json['remaining'] ?? 0).toDouble(),
      bucket: json['bucket'] ?? 'Needs',
    );
  }
}

class LocalBudgetComparison {
  final LocalBudgetSummary budgeted;
  final LocalBudgetSummary actual;
  final LocalBudgetSummary remaining;
  final List<LocalCompareItem> items;

  LocalBudgetComparison({
    required this.budgeted,
    required this.actual,
    required this.remaining,
    required this.items,
  });

  factory LocalBudgetComparison.fromJson(Map<String, dynamic> json) {
    return LocalBudgetComparison(
      budgeted: LocalBudgetSummary.fromJson(json['budgeted'] ?? {}),
      actual: LocalBudgetSummary.fromJson(json['actual'] ?? {}),
      remaining: LocalBudgetSummary.fromJson(json['remaining'] ?? {}),
      items: (json['items'] as List?)
              ?.map((e) => LocalCompareItem.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class LocalSalaryDetail {
  final int id;
  final int incomeId;
  final String? companyName;
  final int fromYear;
  final int? toYear;
  final bool isCurrent;
  final double? fixedPay;
  final double? basicPay;
  final double? hra;
  final double? lta;
  final double? pfEmployee;
  final double? pfEmployer;
  final double? specialAllowance;
  final double? mealCard;
  final double? variablePayPercentage;
  final double? variablePayAmount;
  final double? grossAnnual;
  final double? monthlyInHand;

  LocalSalaryDetail({
    required this.id,
    required this.incomeId,
    this.companyName,
    required this.fromYear,
    this.toYear,
    required this.isCurrent,
    this.fixedPay,
    this.basicPay,
    this.hra,
    this.lta,
    this.pfEmployee,
    this.pfEmployer,
    this.specialAllowance,
    this.mealCard,
    this.variablePayPercentage,
    this.variablePayAmount,
    this.grossAnnual,
    this.monthlyInHand,
  });

  factory LocalSalaryDetail.fromJson(Map<String, dynamic> json) {
    return LocalSalaryDetail(
      id: json['id'],
      incomeId: json['income_id'],
      companyName: json['company_name'],
      fromYear: json['from_year'],
      toYear: json['to_year'],
      isCurrent: json['is_current'] ?? false,
      fixedPay: json['fixed_pay'] != null ? (json['fixed_pay'] as num).toDouble() : null,
      basicPay: json['basic_pay'] != null ? (json['basic_pay'] as num).toDouble() : null,
      hra: json['hra'] != null ? (json['hra'] as num).toDouble() : null,
      lta: json['lta'] != null ? (json['lta'] as num).toDouble() : null,
      pfEmployee: json['pf_employee'] != null ? (json['pf_employee'] as num).toDouble() : null,
      pfEmployer: json['pf_employer'] != null ? (json['pf_employer'] as num).toDouble() : null,
      specialAllowance: json['special_allowance'] != null ? (json['special_allowance'] as num).toDouble() : null,
      mealCard: json['meal_card'] != null ? (json['meal_card'] as num).toDouble() : null,
      variablePayPercentage: json['variable_pay_percentage'] != null ? (json['variable_pay_percentage'] as num).toDouble() : null,
      variablePayAmount: json['variable_pay_amount'] != null ? (json['variable_pay_amount'] as num).toDouble() : null,
      grossAnnual: json['gross_annual'] != null ? (json['gross_annual'] as num).toDouble() : null,
      monthlyInHand: json['monthly_in_hand'] != null ? (json['monthly_in_hand'] as num).toDouble() : null,
    );
  }
}

class LocalSalaryGrowthItem {
  final int year;
  final double grossAnnual;
  final String? company;
  final double? growthPct;

  LocalSalaryGrowthItem({
    required this.year,
    required this.grossAnnual,
    this.company,
    this.growthPct,
  });

  factory LocalSalaryGrowthItem.fromJson(Map<String, dynamic> json) {
    return LocalSalaryGrowthItem(
      year: json['year'],
      grossAnnual: (json['gross_annual'] ?? 0).toDouble(),
      company: json['company'],
      growthPct: json['growth_pct'] != null ? (json['growth_pct'] as num).toDouble() : null,
    );
  }
}

class LocalSalaryGrowth {
  final List<LocalSalaryGrowthItem> growthData;
  final double avgGrowthPct;
  final double totalGrowthPct;
  final int timePeriodYears;

  LocalSalaryGrowth({
    required this.growthData,
    required this.avgGrowthPct,
    required this.totalGrowthPct,
    required this.timePeriodYears,
  });

  factory LocalSalaryGrowth.fromJson(Map<String, dynamic> json) {
    return LocalSalaryGrowth(
      growthData: (json['growth_data'] as List?)
              ?.map((e) => LocalSalaryGrowthItem.fromJson(e))
              .toList() ??
          [],
      avgGrowthPct: (json['avg_growth_pct'] ?? 0).toDouble(),
      totalGrowthPct: (json['total_growth_pct'] ?? 0).toDouble(),
      timePeriodYears: json['time_period_years'] ?? 0,
    );
  }
}

class LocalVehicleService {
  final int id;
  final int vehicleId;
  final String serviceDate;
  final String serviceType;
  final String? description;
  final double cost;
  final double? kmAtService;
  final String? serviceCenter;
  final double? nextServiceKm;

  LocalVehicleService({
    required this.id,
    required this.vehicleId,
    required this.serviceDate,
    required this.serviceType,
    this.description,
    required this.cost,
    this.kmAtService,
    this.serviceCenter,
    this.nextServiceKm,
  });

  factory LocalVehicleService.fromJson(Map<String, dynamic> json) {
    return LocalVehicleService(
      id: json['id'],
      vehicleId: json['vehicle_id'],
      serviceDate: json['service_date'] ?? '',
      serviceType: json['service_type'] ?? 'regular',
      description: json['description'],
      cost: (json['cost'] ?? 0).toDouble(),
      kmAtService: json['km_at_service'] != null ? (json['km_at_service'] as num).toDouble() : null,
      serviceCenter: json['service_center'],
      nextServiceKm: json['next_service_km'] != null ? (json['next_service_km'] as num).toDouble() : null,
    );
  }
}

class LocalVehicleFuel {
  final int id;
  final int vehicleId;
  final String fillDate;
  final double amount;
  final double liters;
  final double? kmAtFill;
  final double? pricePerLiter;
  final bool isFullTank;

  LocalVehicleFuel({
    required this.id,
    required this.vehicleId,
    required this.fillDate,
    required this.amount,
    required this.liters,
    this.kmAtFill,
    this.pricePerLiter,
    required this.isFullTank,
  });

  factory LocalVehicleFuel.fromJson(Map<String, dynamic> json) {
    return LocalVehicleFuel(
      id: json['id'],
      vehicleId: json['vehicle_id'],
      fillDate: json['fill_date'] ?? '',
      amount: (json['amount'] ?? 0).toDouble(),
      liters: (json['liters'] ?? 0).toDouble(),
      kmAtFill: json['km_at_fill'] != null ? (json['km_at_fill'] as num).toDouble() : null,
      pricePerLiter: json['price_per_liter'] != null ? (json['price_per_liter'] as num).toDouble() : null,
      isFullTank: json['is_full_tank'] ?? true,
    );
  }
}

class LocalVehicleLoan {
  final int id;
  final int vehicleId;
  final String? bankName;
  final double loanAmount;
  final double interestRate;
  final int tenureMonths;
  final double emi;
  final String startDate;
  final int emiPaid;
  final double outstanding;

  LocalVehicleLoan({
    required this.id,
    required this.vehicleId,
    this.bankName,
    required this.loanAmount,
    required this.interestRate,
    required this.tenureMonths,
    required this.emi,
    required this.startDate,
    required this.emiPaid,
    required this.outstanding,
  });

  factory LocalVehicleLoan.fromJson(Map<String, dynamic> json) {
    return LocalVehicleLoan(
      id: json['id'],
      vehicleId: json['vehicle_id'],
      bankName: json['bank_name'],
      loanAmount: (json['loan_amount'] ?? 0).toDouble(),
      interestRate: (json['interest_rate'] ?? 0).toDouble(),
      tenureMonths: json['tenure_months'] ?? 12,
      emi: (json['emi'] ?? 0).toDouble(),
      startDate: json['start_date'] ?? '',
      emiPaid: json['emi_paid'] ?? 0,
      outstanding: (json['outstanding'] ?? 0).toDouble(),
    );
  }
}

class LocalElectronicService {
  final int id;
  final int electronicId;
  final String serviceDate;
  final String serviceType;
  final String? description;
  final double cost;
  final String? serviceCenter;

  LocalElectronicService({
    required this.id,
    required this.electronicId,
    required this.serviceDate,
    required this.serviceType,
    this.description,
    required this.cost,
    this.serviceCenter,
  });

  factory LocalElectronicService.fromJson(Map<String, dynamic> json) {
    return LocalElectronicService(
      id: json['id'],
      electronicId: json['electronic_id'],
      serviceDate: json['service_date'] ?? '',
      serviceType: json['service_type'] ?? 'repair',
      description: json['description'],
      cost: (json['cost'] ?? 0).toDouble(),
      serviceCenter: json['service_center'],
    );
  }
}

class LocalElectronicEmi {
  final int id;
  final int electronicId;
  final String bankName;
  final double emiAmount;
  final double interestRate;
  final int totalMonths;
  final int monthsPaid;
  final String startDate;
  final bool startImmediately;

  LocalElectronicEmi({
    required this.id,
    required this.electronicId,
    required this.bankName,
    required this.emiAmount,
    required this.interestRate,
    required this.totalMonths,
    required this.monthsPaid,
    required this.startDate,
    required this.startImmediately,
  });

  factory LocalElectronicEmi.fromJson(Map<String, dynamic> json) {
    return LocalElectronicEmi(
      id: json['id'],
      electronicId: json['electronic_id'],
      bankName: json['bank_name'] ?? '',
      emiAmount: (json['emi_amount'] ?? 0).toDouble(),
      interestRate: (json['interest_rate'] ?? 0).toDouble(),
      totalMonths: json['total_months'] ?? 1,
      monthsPaid: json['months_paid'] ?? 0,
      startDate: json['start_date'] ?? '',
      startImmediately: json['start_immediately'] ?? true,
    );
  }
}

class LocalElectronic {
  final String id;
  final String name;
  final String category;
  final String? brand;
  final String? model;
  final DateTime purchaseDate;
  final double purchaseAmount;
  final int warrantyYears;
  final DateTime? warrantyExpiryDate;
  final int expectedLifeYears;
  final double currentValue;
  final String? notes;
  final String warrantyStatus;
  final int warrantyDaysRemaining;
  final double totalServiceCost;
  final LocalElectronicEmi? emi;

  LocalElectronic({
    required this.id,
    required this.name,
    required this.category,
    this.brand,
    this.model,
    required this.purchaseDate,
    required this.purchaseAmount,
    required this.warrantyYears,
    this.warrantyExpiryDate,
    required this.expectedLifeYears,
    required this.currentValue,
    this.notes,
    required this.warrantyStatus,
    required this.warrantyDaysRemaining,
    required this.totalServiceCost,
    this.emi,
  });

  factory LocalElectronic.fromJson(Map<String, dynamic> json) {
    return LocalElectronic(
      id: json['id'].toString(),
      name: json['name'] ?? '',
      category: json['category'] ?? 'mobile',
      brand: json['brand'],
      model: json['model'],
      purchaseDate: DateTime.parse(json['purchase_date']),
      purchaseAmount: (json['purchase_amount'] ?? 0).toDouble(),
      warrantyYears: json['warranty_years'] ?? 1,
      warrantyExpiryDate: json['warranty_expiry_date'] != null
          ? DateTime.parse(json['warranty_expiry_date'])
          : null,
      expectedLifeYears: json['expected_life_years'] ?? 3,
      currentValue: (json['current_value'] ?? 0).toDouble(),
      notes: json['notes'],
      warrantyStatus: json['warranty_status'] ?? 'expired',
      warrantyDaysRemaining: json['warranty_days_remaining'] ?? 0,
      totalServiceCost: (json['total_service_cost'] ?? 0).toDouble(),
      emi: json['emi'] != null ? LocalElectronicEmi.fromJson(json['emi']) : null,
    );
  }
}
