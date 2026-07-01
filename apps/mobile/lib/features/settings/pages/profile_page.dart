import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
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
        loading: () => const LoadingView(),
        error: (e, _) => Center(child: Text('Could not load profile: $e')),
        data: (state) {
          if (state.error != null) return Center(child: Text(state.error!));
          final p = state.profile;
          return ListView(
            padding: const EdgeInsets.all(AppTheme.spacingLg),
            children: [
              Center(
                child: Column(children: [
                  CircleAvatar(radius: 48, backgroundColor: theme.colorScheme.primaryContainer,
                    child: Text(p?.name?[0]?.toUpperCase() ?? p?.email[0]?.toUpperCase() ?? '?',
                      style: TextStyle(fontSize: 32, color: theme.colorScheme.primary))),
                  const SizedBox(height: AppTheme.spacingMd),
                  Text(p?.name ?? 'User', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
                  Text(p?.email ?? '', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ]),
              ),
              const SizedBox(height: AppTheme.spacingXl),
              AppCard(child: Column(children: [
                InfoRow(label: 'Phone', value: p?.phone ?? '--', labelWidth: 140),
                InfoRow(label: 'Country', value: p?.country ?? '--', labelWidth: 140),
                InfoRow(label: 'Currency', value: p?.baseCurrency ?? '--', labelWidth: 140),
                InfoRow(label: 'Timezone', value: p?.timezone ?? '--', labelWidth: 140),
                InfoRow(label: 'Language', value: p?.language ?? '--', labelWidth: 140),
                InfoRow(label: 'Member since', value: p?.createdAt ?? '--', labelWidth: 140),
              ])),
            ],
          );
        },
      ),
    );
  }
}
