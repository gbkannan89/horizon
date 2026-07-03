import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../../advisor/models/advisor_models.dart';
import '../../advisor/repository/advisor_repository.dart';

final aiSettingsProvider = FutureProvider.autoDispose<List<ProviderInfo>>((ref) async {
  return ref.read(advisorRepositoryProvider).getProviders();
});

class AiSettingsPage extends ConsumerWidget {
  const AiSettingsPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(aiSettingsProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('AI Settings')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (providers) => ListView(
          padding: const EdgeInsets.all(AppSpacing.md),
          children: [
            const SharedSectionHeader(title: 'AI Provider'),
            ...providers.map((p) => Card(
              margin: const EdgeInsets.only(bottom: AppSpacing.sm),
              child: ListTile(
                leading: Icon(p.isActive ? Icons.check_circle : Icons.radio_button_unchecked, color: p.isActive ? theme.colorScheme.primary : null),
                title: Text(p.name),
                subtitle: Text(p.capabilities.join(', ')),
                trailing: p.isActive ? const Chip(label: Text('Active'), visualDensity: VisualDensity.compact) : null,
                onTap: p.isActive ? null : () async {
                  try {
                    await ref.read(advisorRepositoryProvider).switchProvider(p.name);
                    ref.invalidate(aiSettingsProvider);
                  } catch (e) {
                    if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Failed: $e')));
                  }
                },
              ),
            )),
          ],
        ),
      ),
    );
  }
}
