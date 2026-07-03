import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../state/automation_state.dart';

class AutomationListPage extends ConsumerWidget {
  const AutomationListPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final asyncRules = ref.watch(rulesProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Automation Rules'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () => context.push('/automation/add'),
          ),
        ],
      ),
      body: asyncRules.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (rules) {
          if (rules.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.auto_awesome, size: 48, color: Colors.grey.shade400),
                  const SizedBox(height: 16),
                  const Text('No automation rules yet'),
                  const SizedBox(height: 8),
                  ElevatedButton(
                    onPressed: () => context.push('/automation/add'),
                    child: const Text('Create Rule'),
                  ),
                ],
              ),
            );
          }
          return RefreshIndicator(
            onRefresh: () async => ref.invalidate(rulesProvider),
            child: ListView.builder(
              padding: const EdgeInsets.all(AppSpacing.md),
              itemCount: rules.length,
              itemBuilder: (_, i) {
                final rule = rules[i];
                return Card(
                  margin: const EdgeInsets.only(bottom: AppSpacing.sm),
                  child: ListTile(
                    title: Text(rule.name, style: theme.textTheme.titleSmall),
                    subtitle: Text('${rule.category} • Priority: ${rule.priority}'),
                    trailing: Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Icon(rule.enabled ? Icons.toggle_on : Icons.toggle_off_outlined,
                            color: rule.enabled ? Colors.green : Colors.grey, size: 28),
                        const SizedBox(width: 4),
                        const Icon(Icons.chevron_right),
                      ],
                    ),
                    onTap: () => context.push('/automation/${rule.ruleId}'),
                  ),
                );
              },
            ),
          );
        },
      ),
    );
  }
}
