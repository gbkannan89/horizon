import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/providers/auth_state.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import 'package:horizon_mobile/features/auth/repositories/auth_repository.dart';
import '../models/settings_models.dart';
import '../repository/settings_repository.dart';

final privacyProvider = FutureProvider.autoDispose<PrivacySettings>((ref) async {
  final repo = ref.read(settingsRepositoryProvider);
  return await repo.getPrivacy();
});

class PrivacyPage extends ConsumerWidget {
  const PrivacyPage({super.key});

  Future<void> _toggle(WidgetRef ref, PrivacySettings current, PrivacySettings Function(PrivacySettings) update) async {
    final updated = update(current);
    try {
      await ref.read(settingsRepositoryProvider).updatePrivacy(privacy: updated);
      ref.invalidate(privacyProvider);
    } catch (_) {}
  }

  Future<void> _deleteAccount(BuildContext context, WidgetRef ref) async {
    final confirmed = await showDialog<bool>(context: context, builder: (ctx) => AlertDialog(
      title: const Text('Delete Account'),
      content: const Text('This action cannot be undone. All your financial data will be permanently deleted.'),
      actions: [
        TextButton(onPressed: () => Navigator.pop(ctx, false), child: const Text('Cancel')),
        FilledButton(onPressed: () => Navigator.pop(ctx, true), style: FilledButton.styleFrom(backgroundColor: Colors.red), child: const Text('Delete Permanently')),
      ],
    ));
    if (confirmed != true) return;
    try {
      await ref.read(authRepositoryProvider).logout();
      ref.read(authStateProvider.notifier).unauthenticated();
      if (context.mounted) context.go('/login');
    } catch (_) {}
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(privacyProvider);
    return Scaffold(
      appBar: AppBar(title: const Text('Privacy & Data')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (p) => ListView(
          padding: const EdgeInsets.all(16),
          children: [
            Card(child: Column(children: [
              SwitchListTile(title: const Text('Data Sharing'), subtitle: const Text('Share data to improve your experience'), value: p.dataSharingEnabled, onChanged: (v) => _toggle(ref, p, (c) => PrivacySettings(dataSharingEnabled: v, analyticsEnabled: c.analyticsEnabled, personalizeEnabled: c.personalizeEnabled, thirdPartySharing: c.thirdPartySharing, marketingEnabled: c.marketingEnabled))),
              SwitchListTile(title: const Text('Analytics'), subtitle: const Text('Help us improve with usage data'), value: p.analyticsEnabled, onChanged: (v) => _toggle(ref, p, (c) => PrivacySettings(dataSharingEnabled: c.dataSharingEnabled, analyticsEnabled: v, personalizeEnabled: c.personalizeEnabled, thirdPartySharing: c.thirdPartySharing, marketingEnabled: c.marketingEnabled))),
              SwitchListTile(title: const Text('Personalization'), subtitle: const Text('Personalized recommendations'), value: p.personalizeEnabled, onChanged: (v) => _toggle(ref, p, (c) => PrivacySettings(dataSharingEnabled: c.dataSharingEnabled, analyticsEnabled: c.analyticsEnabled, personalizeEnabled: v, thirdPartySharing: c.thirdPartySharing, marketingEnabled: c.marketingEnabled))),
              SwitchListTile(title: const Text('Third-party Sharing'), subtitle: const Text('Share with trusted partners'), value: p.thirdPartySharing, onChanged: (v) => _toggle(ref, p, (c) => PrivacySettings(dataSharingEnabled: c.dataSharingEnabled, analyticsEnabled: c.analyticsEnabled, personalizeEnabled: c.personalizeEnabled, thirdPartySharing: v, marketingEnabled: c.marketingEnabled))),
              SwitchListTile(title: const Text('Marketing'), subtitle: const Text('Receive marketing communications'), value: p.marketingEnabled, onChanged: (v) => _toggle(ref, p, (c) => PrivacySettings(dataSharingEnabled: c.dataSharingEnabled, analyticsEnabled: c.analyticsEnabled, personalizeEnabled: c.personalizeEnabled, thirdPartySharing: c.thirdPartySharing, marketingEnabled: v))),
            ])),
            const SizedBox(height: 16),
            OutlinedButton.icon(onPressed: () => ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Data export coming soon'))), icon: const Icon(Icons.download), label: const Text('Export My Data')),
            const SizedBox(height: 8),
            OutlinedButton.icon(onPressed: () => _deleteAccount(context, ref), icon: const Icon(Icons.delete_forever, color: Colors.red), label: const Text('Delete Account', style: TextStyle(color: Colors.red))),
          ],
        ),
      ),
    );
  }
}
