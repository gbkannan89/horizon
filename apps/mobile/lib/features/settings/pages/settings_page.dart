import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/providers/app_state.dart';
import 'package:horizon_mobile/shared/providers/auth_state.dart';
import 'package:horizon_mobile/features/auth/repositories/auth_repository.dart';
import 'package:horizon_mobile/features/settings/repository/settings_repository.dart';

class SettingsPage extends ConsumerWidget {
  const SettingsPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Settings')),
      body: ListView(
        children: [
          _section(theme, 'Profile', [
            ListTile(leading: const Icon(Icons.person_outline), title: const Text('Profile'), trailing: const Icon(Icons.chevron_right, size: 18), onTap: () => context.push('/settings/profile')),
            ListTile(leading: const Icon(Icons.edit_outlined), title: const Text('Edit Profile'), trailing: const Icon(Icons.chevron_right, size: 18), onTap: () => context.push('/settings/profile/edit')),
          ]),
          _section(theme, 'Preferences', [
            ListTile(leading: const Icon(Icons.tune), title: const Text('Preferences'), trailing: const Icon(Icons.chevron_right, size: 18), onTap: () => context.push('/settings/preferences')),
            ListTile(leading: const Icon(Icons.palette_outlined), title: const Text('Theme'), onTap: () => _showThemePicker(context, ref)),
          ]),
          _section(theme, 'Privacy & Security', [
            ListTile(leading: const Icon(Icons.lock_outline), title: const Text('Privacy'), trailing: const Icon(Icons.chevron_right, size: 18), onTap: () => context.push('/settings/privacy')),
            ListTile(leading: const Icon(Icons.security), title: const Text('Security'), trailing: const Icon(Icons.chevron_right, size: 18), onTap: () => context.push('/settings/security')),
          ]),
          _section(theme, 'About', [
            ListTile(leading: const Icon(Icons.info_outline), title: const Text('About'), trailing: const Icon(Icons.chevron_right, size: 18), onTap: () => context.push('/settings/about')),
            ListTile(
              leading: Icon(Icons.logout, color: theme.colorScheme.error),
              title: Text('Sign Out', style: TextStyle(color: theme.colorScheme.error)),
              onTap: () async {
                await ref.read(authRepositoryProvider).logout();
                ref.read(authStateProvider.notifier).unauthenticated();
                if (context.mounted) context.go('/login');
              },
            ),
          ]),
        ],
      ),
    );
  }

  Widget _section(ThemeData theme, String title, List<Widget> children) => Column(
    crossAxisAlignment: CrossAxisAlignment.start,
    children: [
      Padding(
        padding: const EdgeInsets.fromLTRB(16, 24, 16, 8),
        child: Text(title, style: TextStyle(fontWeight: FontWeight.w600, color: theme.colorScheme.primary, fontSize: 13)),
      ),
      ...children,
    ],
  );

  void _showThemePicker(BuildContext context, WidgetRef ref) {
    showModalBottomSheet(
      context: context,
      builder: (_) => SafeArea(
        child: Column(mainAxisSize: MainAxisSize.min, children: [
          ListTile(title: const Text('Light'), leading: const Icon(Icons.light_mode), onTap: () { Navigator.pop(context); _setTheme(ref, 'light'); }),
          ListTile(title: const Text('Dark'), leading: const Icon(Icons.dark_mode), onTap: () { Navigator.pop(context); _setTheme(ref, 'dark'); }),
          ListTile(title: const Text('System'), leading: const Icon(Icons.settings_brightness), onTap: () { Navigator.pop(context); _setTheme(ref, 'system'); }),
        ]),
      ),
    );
  }

  void _setTheme(WidgetRef ref, String theme) async {
    if (theme == 'dark') ref.read(themeModeProvider.notifier).state = ThemeMode.dark;
    else if (theme == 'light') ref.read(themeModeProvider.notifier).state = ThemeMode.light;
    else ref.read(themeModeProvider.notifier).state = ThemeMode.system;
    try {
      final repo = ref.read(settingsRepositoryProvider);
      final current = await repo.getPreferences();
      current.theme = theme;
      await repo.updatePreferences(prefs: current);
    } catch (_) {}
  }
}
