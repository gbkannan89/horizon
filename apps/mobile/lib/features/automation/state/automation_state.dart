import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/automation_models.dart';
import '../repository/automation_repository.dart';

final rulesProvider = FutureProvider.autoDispose<List<RuleModel>>((ref) async {
  final repo = ref.watch(automationRepositoryProvider);
  final resp = await repo.getRules();
  return resp.data ?? [];
});

final ruleDetailProvider = FutureProvider.family.autoDispose<RuleModel?, String>((ref, id) async {
  final repo = ref.watch(automationRepositoryProvider);
  final resp = await repo.getRule(id);
  return resp.data;
});
