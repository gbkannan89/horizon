import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../../notifications/models/notification_models.dart';
import '../../notifications/repository/notification_repository.dart';

final notifPrefsProvider = FutureProvider.autoDispose<NotificationPreferences>((ref) async {
  return ref.read(notificationRepositoryProvider).getPreferences();
});

class NotificationPreferencesPage extends ConsumerStatefulWidget {
  const NotificationPreferencesPage({super.key});
  @override
  ConsumerState<NotificationPreferencesPage> createState() => _NotificationPreferencesPageState();
}

class _NotificationPreferencesPageState extends ConsumerState<NotificationPreferencesPage> {
  Future<void> _toggle(NotificationCategoryPref pref, {bool? inApp, bool? push, bool? email, bool? enabled}) async {
    final updated = NotificationCategoryPref(
      category: pref.category,
      inApp: inApp ?? pref.inApp,
      push: push ?? pref.push,
      email: email ?? pref.email,
      digest: pref.digest,
      enabled: enabled ?? pref.enabled,
    );
    try {
      final repo = ref.read(notificationRepositoryProvider);
      final current = await repo.getPreferences();
      final newMap = Map<String, NotificationCategoryPref>.from(current.categories);
      newMap[pref.category] = updated;
      await repo.updatePreferences(NotificationPreferences(categories: newMap));
      ref.invalidate(notifPrefsProvider);
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Failed: $e')));
    }
  }

  @override
  Widget build(BuildContext context) {
    final async = ref.watch(notifPrefsProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Notification Preferences')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (prefs) {
          final cats = prefs.categories.values.toList();
          if (cats.isEmpty) return const SharedEmptyView(icon: Icons.notifications_off, title: 'No preferences available');
          return ListView(
            padding: const EdgeInsets.all(AppSpacing.md),
            children: cats.map((cat) => Card(
              margin: const EdgeInsets.only(bottom: AppSpacing.sm),
              child: Padding(
                padding: const EdgeInsets.all(AppSpacing.sm),
                child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Text(cat.category, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                  const SizedBox(height: AppSpacing.xs),
                  SwitchListTile(title: const Text('In-App'), value: cat.inApp, onChanged: (v) => _toggle(cat, inApp: v), dense: true),
                  SwitchListTile(title: const Text('Push'), value: cat.push, onChanged: (v) => _toggle(cat, push: v), dense: true),
                  SwitchListTile(title: const Text('Email'), value: cat.email, onChanged: (v) => _toggle(cat, email: v), dense: true),
                ]),
              ),
            )).toList(),
          );
        },
      ),
    );
  }
}
