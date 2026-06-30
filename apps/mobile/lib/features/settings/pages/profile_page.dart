import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../repository/settings_repository.dart';

final profileProvider = FutureProvider<UserProfileState>((ref) async {
  try {
    final repo = ref.read(settingsRepositoryProvider);
    final profile = await repo.getProfile();
    return UserProfileState(profile: profile);
  } catch (e) {
    return UserProfileState(error: e.toString());
  }
});

class UserProfileState {
  final dynamic profile; final String? error;
  UserProfileState({this.profile, this.error});
  bool get loading => profile == null && error == null;
}

class ProfilePage extends ConsumerWidget {
  const ProfilePage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(profileProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Profile')),
      body: async.when(
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, _) => Center(child: Text('Could not load profile: $e')),
        data: (state) {
          if (state.error != null) return Center(child: Text(state.error!));
          final p = state.profile;
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              // Avatar
              Center(
                child: Column(children: [
                  CircleAvatar(radius: 48, backgroundColor: theme.colorScheme.primaryContainer,
                    child: Text(p?.name?[0]?.toUpperCase() ?? p?.email[0]?.toUpperCase() ?? '?',
                      style: TextStyle(fontSize: 32, color: theme.colorScheme.primary))),
                  const SizedBox(height: 12),
                  Text(p?.name ?? 'User', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
                  Text(p?.email ?? '', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ]),
              ),
              const SizedBox(height: 24),
              Card(child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(children: [
                  _row(theme, 'Phone', p?.phone ?? '--'),
                  _row(theme, 'Country', p?.country ?? '--'),
                  _row(theme, 'Currency', p?.baseCurrency ?? '--'),
                  _row(theme, 'Timezone', p?.timezone ?? '--'),
                  _row(theme, 'Language', p?.language ?? '--'),
                  _row(theme, 'Member since', p?.createdAt ?? '--'),
                ]),
              )),
            ],
          );
        },
      ),
    );
  }

  Widget _row(ThemeData t, String l, String v) => Padding(
    padding: const EdgeInsets.symmetric(vertical: 6),
    child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
      Text(l, style: t.textTheme.bodyMedium?.copyWith(color: t.colorScheme.onSurfaceVariant)),
      Text(v, style: t.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500)),
    ]),
  );
}
