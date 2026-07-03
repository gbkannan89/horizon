import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/features/recurring/models/recurring_models.dart';
import 'package:horizon_mobile/features/recurring/repository/recurring_repository.dart';

final recurringListProvider = FutureProvider.autoDispose<List<RecurringTransaction>>((ref) async {
  final repo = ref.watch(recurringRepositoryProvider);
  final resp = await repo.getRecurringTransactions();
  return resp.data ?? [];
});

final recurringDetailProvider = FutureProvider.family.autoDispose<RecurringTransaction?, String>((ref, id) async {
  final repo = ref.watch(recurringRepositoryProvider);
  final resp = await repo.getRecurringTransaction(id);
  return resp.data;
});
