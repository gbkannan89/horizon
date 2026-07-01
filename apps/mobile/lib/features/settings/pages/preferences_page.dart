import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/settings_models.dart';
import '../repository/settings_repository.dart';

final prefsProvider = FutureProvider.autoDispose<UserPreferences>((ref) async {
  final repo = ref.read(settingsRepositoryProvider);
  return await repo.getPreferences();
});

class PreferencesPage extends ConsumerStatefulWidget {
  const PreferencesPage({super.key});
  @override
  ConsumerState<PreferencesPage> createState() => _PreferencesPageState();
}

class _PreferencesPageState extends ConsumerState<PreferencesPage> {
  Future<void> _updatePrefs(UserPreferences prefs, {String? currency, String? language, bool? timelineCompact, bool? pushNotifications, bool? emailNotifications, bool? smsNotifications, bool? digestEnabled}) async {
    final updated = UserPreferences(
      baseCurrency: currency ?? prefs.baseCurrency,
      language: language ?? prefs.language,
      theme: prefs.theme,
      emailNotifications: emailNotifications ?? prefs.emailNotifications,
      pushNotifications: pushNotifications ?? prefs.pushNotifications,
      smsNotifications: smsNotifications ?? prefs.smsNotifications,
      digestEnabled: digestEnabled ?? prefs.digestEnabled,
      digestFrequency: prefs.digestFrequency,
      dashboardDefaultView: prefs.dashboardDefaultView,
      timelineCompact: timelineCompact ?? prefs.timelineCompact,
    );
    try {
      final repo = ref.read(settingsRepositoryProvider);
      await repo.updatePreferences(prefs: updated);
      ref.invalidate(prefsProvider);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Failed to update: $e'), behavior: SnackBarBehavior.floating));
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final async = ref.watch(prefsProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Preferences')),
      body: async.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Error: $e')),
        data: (prefs) => ListView(
          padding: const EdgeInsets.all(16),
          children: [
            _section(theme, 'Display', [
              ListTile(title: const Text('Currency'), subtitle: Text(prefs.baseCurrency), trailing: const Icon(Icons.chevron_right, size: 18)),
              ListTile(title: const Text('Language'), subtitle: Text(prefs.language.toUpperCase()), trailing: const Icon(Icons.chevron_right, size: 18)),
              SwitchListTile(title: const Text('Compact Timeline'), value: prefs.timelineCompact, onChanged: (v) => _updatePrefs(prefs, timelineCompact: v)),
            ]),
            _section(theme, 'Notifications', [
              SwitchListTile(title: const Text('Push Notifications'), value: prefs.pushNotifications, onChanged: (v) => _updatePrefs(prefs, pushNotifications: v)),
              SwitchListTile(title: const Text('Email Notifications'), value: prefs.emailNotifications, onChanged: (v) => _updatePrefs(prefs, emailNotifications: v)),
              SwitchListTile(title: const Text('SMS Notifications'), value: prefs.smsNotifications, onChanged: (v) => _updatePrefs(prefs, smsNotifications: v)),
              SwitchListTile(title: const Text('Weekly Digest'), value: prefs.digestEnabled, onChanged: (v) => _updatePrefs(prefs, digestEnabled: v)),
            ]),
          ],
        ),
      ),
    );
  }

  Widget _section(ThemeData t, String title, List<Widget> children) => Column(
    crossAxisAlignment: CrossAxisAlignment.start,
    children: [
      Padding(padding: const EdgeInsets.only(bottom: 8), child: Text(title, style: TextStyle(color: t.colorScheme.primary, fontWeight: FontWeight.w600))),
      Card(child: Column(children: children)),
      const SizedBox(height: 16),
    ],
  );
}
