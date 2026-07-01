import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../state/notification_state.dart';

class NotificationsPage extends ConsumerStatefulWidget {
  const NotificationsPage({super.key});
  @override
  ConsumerState<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends ConsumerState<NotificationsPage> {
  @override
  void initState() { super.initState(); Future.microtask(() => ref.read(notificationProvider.notifier).load()); }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(notificationProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Notifications')),
      body: _buildBody(theme, state),
    );
  }

  Widget _buildBody(ThemeData theme, NotificationState state) {
    switch (state.status) {
      case NotificationStatus.initial:
      case NotificationStatus.loading:
        return const SharedLoadingView();
      case NotificationStatus.error:
        return SharedErrorView(message: state.error, onRetry: () => ref.read(notificationProvider.notifier).load());
      case NotificationStatus.empty:
        return SharedEmptyView(icon: Icons.notifications_none, title: 'No notifications', subtitle: 'You\'re all caught up');
      case NotificationStatus.loaded:
        return RefreshIndicator(
          onRefresh: () => ref.read(notificationProvider.notifier).load(),
          child: ListView(
            padding: const EdgeInsets.all(AppSpacing.md),
            children: [
              _buildTabs(theme, state),
              const SizedBox(height: AppSpacing.sm),
              ...state.notifications.map((n) => Dismissible(
                key: Key(n.id),
                direction: DismissDirection.endToStart,
                background: Container(alignment: Alignment.centerRight, padding: const EdgeInsets.only(right: AppSpacing.md), color: theme.colorScheme.error, child: const Icon(Icons.archive, color: Colors.white)),
                onDismissed: (_) => ref.read(notificationProvider.notifier).archive(n.id),
                child: Card(
                  clipBehavior: Clip.antiAlias,
                  margin: const EdgeInsets.only(bottom: AppSpacing.sm),
                  child: InkWell(
                    onTap: () => ref.read(notificationProvider.notifier).markRead(n.id),
                    child: Padding(
                      padding: const EdgeInsets.all(AppSpacing.md),
                      child: Row(children: [
                        if (!n.isRead) Container(width: 8, height: 8, decoration: BoxDecoration(color: theme.colorScheme.primary, shape: BoxShape.circle)),
                        if (!n.isRead) const SizedBox(width: AppSpacing.sm),
                        Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                          Text(n.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: n.isRead ? FontWeight.normal : FontWeight.w600)),
                          if (n.body != null) Text(n.body!, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 2, overflow: TextOverflow.ellipsis),
                        ])),
                      ]),
                    ),
                  ),
                ),
              )),
            ],
          ),
        );
    }
  }

  Widget _buildTabs(ThemeData theme, NotificationState state) {
    return Row(children: [
      _tab(theme, 'All', state.activeTab == 'all', () => ref.read(notificationProvider.notifier).load()),
      const SizedBox(width: AppSpacing.sm),
      _tab(theme, 'Unread', state.activeTab == 'unread', () => ref.read(notificationProvider.notifier).loadUnread()),
    ]);
  }

  Widget _tab(ThemeData theme, String label, bool selected, VoidCallback onTap) {
    return ChoiceChip(label: Text(label), selected: selected, onSelected: (_) => onTap(), visualDensity: VisualDensity.compact);
  }
}
