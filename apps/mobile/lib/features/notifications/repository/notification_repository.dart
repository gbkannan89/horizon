import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/notification_models.dart';

final notificationRepositoryProvider = Provider<NotificationRepository>((ref) {
  return NotificationRepository(apiClient: ref.read(apiClientProvider));
});

class NotificationRepository {
  final ApiClient apiClient;
  NotificationRepository({required this.apiClient});

  Future<List<NotificationItem>> getNotifications() async {
    final r = await apiClient.get('/notifications');
    return _parseList(r.data);
  }

  Future<List<NotificationItem>> getUnread() async {
    final r = await apiClient.get('/notifications/unread');
    return _parseList(r.data);
  }

  Future<List<NotificationItem>> getHistory() async {
    final r = await apiClient.get('/notifications/history');
    return _parseList(r.data);
  }

  Future<void> markRead(String id) async {
    await apiClient.post('/notifications/$id/read');
  }

  Future<void> archive(String id) async {
    await apiClient.post('/notifications/$id/archive');
  }

  Future<NotificationPreferences> getPreferences() async {
    final r = await apiClient.get('/notifications/preferences');
    final json = r.data as Map<String, dynamic>;
    final data = json['data'] as Map<String, dynamic>? ?? json;
    return NotificationPreferences.fromJson(data);
  }

  Future<void> updatePreferences(NotificationPreferences prefs) async {
    await apiClient.put('/notifications/preferences', data: prefs.toJson());
  }

  List<NotificationItem> _parseList(dynamic json) {
    final map = json is Map<String, dynamic> ? json : <String, dynamic>{};
    final data = map['data'] as Map<String, dynamic>? ?? map;
    final items = data['notifications'] as List? ?? data['items'] as List? ?? [];
    return items.map((e) => NotificationItem.fromJson(e as Map<String, dynamic>)).toList();
  }
}
