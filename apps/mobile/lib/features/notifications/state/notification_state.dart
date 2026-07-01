import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/notification_models.dart';
import '../repository/notification_repository.dart';

final notificationProvider = StateNotifierProvider<NotificationNotifier, NotificationState>((ref) {
  return NotificationNotifier(ref.read(notificationRepositoryProvider));
});

enum NotificationStatus { initial, loading, loaded, error, empty }

class NotificationState {
  final NotificationStatus status;
  final List<NotificationItem> notifications;
  final String activeTab;
  final String? error;

  const NotificationState({
    this.status = NotificationStatus.initial,
    this.notifications = const [],
    this.activeTab = 'all',
    this.error,
  });

  NotificationState copyWith({
    NotificationStatus? status, List<NotificationItem>? notifications,
    String? activeTab, String? error, bool clearError = false,
  }) => NotificationState(
    status: status ?? this.status, notifications: notifications ?? this.notifications,
    activeTab: activeTab ?? this.activeTab, error: clearError ? null : error ?? this.error,
  );
}

class NotificationNotifier extends StateNotifier<NotificationState> {
  final NotificationRepository _repo;
  NotificationNotifier(this._repo) : super(const NotificationState());

  Future<void> load() async {
    state = state.copyWith(status: NotificationStatus.loading, activeTab: 'all');
    try {
      final items = await _repo.getNotifications();
      if (items.isEmpty) { state = state.copyWith(status: NotificationStatus.empty); return; }
      state = state.copyWith(status: NotificationStatus.loaded, notifications: items);
    } catch (e) { state = state.copyWith(status: NotificationStatus.error, error: e.toString()); }
  }

  Future<void> loadUnread() async {
    state = state.copyWith(status: NotificationStatus.loading, activeTab: 'unread');
    try {
      final items = await _repo.getUnread();
      if (items.isEmpty) { state = state.copyWith(status: NotificationStatus.empty); return; }
      state = state.copyWith(status: NotificationStatus.loaded, notifications: items);
    } catch (e) { state = state.copyWith(status: NotificationStatus.error, error: e.toString()); }
  }

  Future<void> markRead(String id) async {
    try { await _repo.markRead(id); } catch (_) {}
    final updated = state.notifications.map((n) => n.id == id ? NotificationItem(id: n.id, title: n.title, body: n.body, category: n.category, priority: n.priority, isRead: true, createdAt: n.createdAt) : n).toList();
    state = state.copyWith(notifications: updated);
  }

  Future<void> archive(String id) async {
    try { await _repo.archive(id); } catch (_) {}
    state = state.copyWith(notifications: state.notifications.where((n) => n.id != id).toList());
    if (state.notifications.isEmpty) state = state.copyWith(status: NotificationStatus.empty);
  }

  Future<void> search(String query) async {
    if (query.isEmpty) { await load(); return; }
    state = state.copyWith(status: NotificationStatus.loading);
    try {
      final items = await _repo.search(query);
      if (items.isEmpty) { state = state.copyWith(status: NotificationStatus.empty); return; }
      state = state.copyWith(status: NotificationStatus.loaded, notifications: items);
    } catch (e) { state = state.copyWith(status: NotificationStatus.error, error: e.toString()); }
  }
}
