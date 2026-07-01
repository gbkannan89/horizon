import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/settings_models.dart';
import '../repository/settings_repository.dart';

final privacyProvider = FutureProvider.autoDispose<PrivacySettings>((ref) async {
  final repo = ref.read(settingsRepositoryProvider);
  return await repo.getPrivacy();
});

class PrivacyPage extends ConsumerWidget {
  const PrivacyPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(privacyProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Privacy & Data')),
      body: async.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Error: $e')),
        data: (p) => ListView(
          padding: const EdgeInsets.all(16),
          children: [
            Card(child: Column(children: [
              SwitchListTile(title: const Text('Data Sharing'), subtitle: const Text('Share data to improve your experience'), value: p.dataSharingEnabled, onChanged: (_) {}),
              SwitchListTile(title: const Text('Analytics'), subtitle: const Text('Help us improve with usage data'), value: p.analyticsEnabled, onChanged: (_) {}),
              SwitchListTile(title: const Text('Personalization'), subtitle: const Text('Personalized recommendations'), value: p.personalizeEnabled, onChanged: (_) {}),
              SwitchListTile(title: const Text('Third-party Sharing'), subtitle: const Text('Share with trusted partners'), value: p.thirdPartySharing, onChanged: (_) {}),
              SwitchListTile(title: const Text('Marketing'), subtitle: const Text('Receive marketing communications'), value: p.marketingEnabled, onChanged: (_) {}),
            ])),
            const SizedBox(height: 16),
            OutlinedButton.icon(onPressed: () {}, icon: const Icon(Icons.download), label: const Text('Export My Data')),
            const SizedBox(height: 8),
            OutlinedButton.icon(onPressed: () => _confirmDelete(context), icon: const Icon(Icons.delete_forever, color: Colors.red), label: Text('Delete Account', style: TextStyle(color: Colors.red))),
          ],
        ),
      ),
    );
  }

  void _confirmDelete(BuildContext context) {
    showDialog(context: context, builder: (_) => AlertDialog(
      title: const Text('Delete Account'),
      content: const Text('This action cannot be undone. All your financial data will be permanently deleted.'),
      actions: [
        TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
        FilledButton(onPressed: () => Navigator.pop(context), style: FilledButton.styleFrom(backgroundColor: Colors.red), child: const Text('Delete Permanently')),
      ],
    ));
  }
}
