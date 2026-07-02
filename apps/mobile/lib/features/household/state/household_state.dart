import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/features/budget/models/budget_models.dart';
import 'package:horizon_mobile/features/goals/models/goal_models.dart';
import 'package:horizon_mobile/features/household/models/household_models.dart';
import 'package:horizon_mobile/features/household/repository/household_repository.dart';

final householdRepositoryProvider = Provider<HouseholdRepository>((ref) {
  final apiClient = ref.watch(apiClientProvider);
  return HouseholdRepository(apiClient.dio);
});

final householdsProvider = FutureProvider<List<HouseholdView>>((ref) async {
  final repository = ref.watch(householdRepositoryProvider);
  return repository.getHouseholds();
});

final householdDetailProvider = FutureProvider.family<HouseholdDetailView, String>((ref, id) async {
  final repository = ref.watch(householdRepositoryProvider);
  return repository.getHouseholdDetail(id);
});

final householdSummaryProvider = FutureProvider.family<HouseholdFinancialSummary, String>((ref, id) async {
  final repository = ref.watch(householdRepositoryProvider);
  return repository.getHouseholdSummary(id);
});

final householdGoalsProvider = FutureProvider.family<List<GoalSummary>, String>((ref, householdId) async {
  final repository = ref.watch(householdRepositoryProvider);
  return repository.getHouseholdGoals(householdId);
});

final householdBudgetsProvider = FutureProvider.family<List<BudgetModel>, String>((ref, householdId) async {
  final repository = ref.watch(householdRepositoryProvider);
  return repository.getHouseholdBudgets(householdId);
});
