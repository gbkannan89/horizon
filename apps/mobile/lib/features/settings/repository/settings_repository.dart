import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/settings_models.dart';

final settingsRepositoryProvider = Provider<SettingsRepository>((ref) {
  return SettingsRepository(apiClient: ref.read(apiClientProvider));
});

class SettingsRepository {
  final ApiClient apiClient;
  SettingsRepository({required this.apiClient});

  Future<UserProfile> getProfile({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/users/${userId ?? 'me'}', queryParameters: p);
    final data = r.data is Map ? (r.data as Map)['data'] as Map? ?? r.data as Map : {};
    return UserProfile.fromJson(data as Map<String, dynamic>);
  }

  Future<UserProfile> updateProfile({String? userId, required Map<String, dynamic> data}) async {
    final r = await apiClient.put('/users/${userId ?? 'me'}', data: data);
    final resp = r.data is Map ? (r.data as Map)['data'] as Map? ?? r.data as Map : {};
    return UserProfile.fromJson(resp as Map<String, dynamic>);
  }

  Future<UserPreferences> getPreferences({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/users/preferences', queryParameters: p);
    final data = r.data is Map ? (r.data as Map)['data'] as Map? ?? r.data as Map : {};
    return UserPreferences.fromJson(data as Map<String, dynamic>);
  }

  Future<UserPreferences> updatePreferences({String? userId, required UserPreferences prefs}) async {
    final r = await apiClient.put('/users/preferences', data: prefs.toJson());
    final data = r.data is Map ? (r.data as Map)['data'] as Map? ?? r.data as Map : {};
    return UserPreferences.fromJson(data as Map<String, dynamic>);
  }

  Future<PrivacySettings> getPrivacy({String? userId}) async {
    final p = <String, dynamic>{}; if (userId != null) p['user_id'] = userId;
    final r = await apiClient.get('/users/privacy', queryParameters: p);
    final data = r.data is Map ? (r.data as Map)['data'] as Map? ?? r.data as Map : {};
    return PrivacySettings.fromJson(data as Map<String, dynamic>);
  }

  Future<PrivacySettings> updatePrivacy({String? userId, required PrivacySettings privacy}) async {
    final r = await apiClient.put('/users/privacy', data: privacy.toJson());
    final data = r.data is Map ? (r.data as Map)['data'] as Map? ?? r.data as Map : {};
    return PrivacySettings.fromJson(data as Map<String, dynamic>);
  }
}
