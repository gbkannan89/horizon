import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../repository/automation_repository.dart';
import '../state/automation_state.dart';

class AutomationDetailPage extends ConsumerWidget {
  final String ruleId;
  const AutomationDetailPage({super.key, required this.ruleId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncRule = ref.watch(ruleDetailProvider(ruleId));
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Rule Details'),
        actions: [
          IconButton(
            icon: const Icon(Icons.delete, color: Colors.red),
            onPressed: () async {
              final confirm = await showDialog<bool>(
                context: context,
                builder: (ctx) => AlertDialog(
                  title: const Text('Delete Rule?'),
                  content: const Text('This cannot be undone.'),
                  actions: [
                    TextButton(onPressed: () => ctx.pop(false), child: const Text('Cancel')),
                    TextButton(onPressed: () => ctx.pop(true), child: const Text('Delete', style: TextStyle(color: Colors.red))),
                  ],
                ),
              );
              if (confirm == true) {
                await ref.read(automationRepositoryProvider).deleteRule(ruleId);
                ref.invalidate(rulesProvider);
                if (context.mounted) context.pop();
              }
            },
          ),
        ],
      ),
      body: asyncRule.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (rule) {
          if (rule == null) return const Center(child: Text('Not found'));
          return ListView(
            padding: const EdgeInsets.all(AppSpacing.md),
            children: [
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Expanded(child: Text(rule.name, style: theme.textTheme.titleLarge)),
                          Icon(rule.enabled ? Icons.toggle_on : Icons.toggle_off_outlined,
                              color: rule.enabled ? Colors.green : Colors.grey, size: 32),
                        ],
                      ),
                      if (rule.description.isNotEmpty) ...[
                        const SizedBox(height: AppSpacing.sm),
                        Text(rule.description, style: theme.textTheme.bodyMedium),
                      ],
                      const SizedBox(height: AppSpacing.md),
                      _infoRow(theme, 'Category', rule.category),
                      _infoRow(theme, 'Priority', '${rule.priority}'),
                      _infoRow(theme, 'Enabled', rule.enabled ? 'Yes' : 'No'),
                    ],
                  ),
                ),
              ),
              if (rule.conditions.isNotEmpty) ...[
                const SizedBox(height: AppSpacing.md),
                const SharedSectionHeader(title: 'Conditions'),
                const SizedBox(height: AppSpacing.sm),
                ...rule.conditions.map((c) => Card(
                  margin: const EdgeInsets.only(bottom: AppSpacing.sm),
                  child: Padding(
                    padding: const EdgeInsets.all(AppSpacing.md),
                    child: Text('${c.field} ${c.operator} ${c.value}'),
                  ),
                )),
              ],
              if (rule.actions.isNotEmpty) ...[
                const SizedBox(height: AppSpacing.md),
                const SharedSectionHeader(title: 'Actions'),
                const SizedBox(height: AppSpacing.sm),
                ...rule.actions.map((a) => Card(
                  margin: const EdgeInsets.only(bottom: AppSpacing.sm),
                  child: Padding(
                    padding: const EdgeInsets.all(AppSpacing.md),
                    child: Text('Type: ${a.type}', style: theme.textTheme.bodyMedium),
                  ),
                )),
              ],
            ],
          );
        },
      ),
    );
  }

  Widget _infoRow(ThemeData t, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(label, style: t.textTheme.bodySmall?.copyWith(color: t.colorScheme.onSurfaceVariant)),
          Text(value, style: t.textTheme.bodyMedium),
        ],
      ),
    );
  }
}
