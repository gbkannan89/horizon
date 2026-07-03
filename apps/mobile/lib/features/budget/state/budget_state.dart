import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/budget_models.dart';
import '../repository/budget_repository.dart';

final budgetsProvider = FutureProvider.autoDispose<List<BudgetModel>>((ref) async {
  final repo = ref.watch(budgetRepositoryProvider);
  final resp = await repo.getBudgets();
  return resp.data?.budgets ?? [];
});

final budgetDetailProvider = FutureProvider.family.autoDispose<BudgetModel?, String>((ref, id) async {
  final repo = ref.watch(budgetRepositoryProvider);
  final resp = await repo.getBudgetVsActual(id);
  return resp.data;
});

final budgetTransitionProvider = FutureProvider.family.autoDispose<BudgetResultResponse?, TransitionParams>((ref, params) async {
  final repo = ref.watch(budgetRepositoryProvider);
  switch (params.action) {
    case BudgetAction.activate:
      return repo.activateBudget(params.budgetId);
    case BudgetAction.pause:
      return repo.pauseBudget(params.budgetId);
    case BudgetAction.resume:
      return repo.resumeBudget(params.budgetId);
    case BudgetAction.complete:
      return repo.completeBudget(params.budgetId);
    case BudgetAction.archive:
      return repo.archiveBudget(params.budgetId);
  }
});

class TransitionParams {
  final String budgetId;
  final BudgetAction action;
  const TransitionParams({required this.budgetId, required this.action});

  @override
  bool operator ==(Object other) =>
      identical(this, other) || other is TransitionParams && budgetId == other.budgetId && action == other.action;

  @override
  int get hashCode => budgetId.hashCode ^ action.hashCode;
}

enum BudgetAction { activate, pause, resume, complete, archive }
