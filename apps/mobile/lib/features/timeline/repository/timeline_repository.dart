import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import '../models/timeline_models.dart';

final timelineRepositoryProvider = Provider<TimelineRepository>((ref) {
  return TimelineRepository(apiClient: ref.read(apiClientProvider));
});

class TimelineRepository {
  final ApiClient apiClient;

  TimelineRepository({required this.apiClient});

  Future<TimelineResponse> getTimeline({
    String? userId, String? view, String? cursor, int limit = 50, FilterParams? filters,
  }) async {
    final params = <String, dynamic>{'limit': limit};
    if (userId != null) params['user_id'] = userId;
    if (view != null) params['view'] = view;
    if (cursor != null) params['cursor'] = cursor;
    if (filters != null) params.addAll(filters.toQuery());
    final response = await apiClient.get('/timeline', queryParameters: params);
    return TimelineResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<TimelineResponse> getRecent({String? userId, int limit = 10}) async {
    final params = <String, dynamic>{'limit': limit};
    if (userId != null) params['user_id'] = userId;
    final response = await apiClient.get('/timeline/recent', queryParameters: params);
    return TimelineResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<TimelineResponse> getFiltered({String? userId, required FilterParams filters, int limit = 50, String? cursor}) async {
    final params = <String, dynamic>{'limit': limit, ...filters.toQuery()};
    if (userId != null) params['user_id'] = userId;
    if (cursor != null) params['cursor'] = cursor;
    final response = await apiClient.get('/timeline/filter', queryParameters: params);
    return TimelineResponse.fromJson(response.data as Map<String, dynamic>);
  }

  Future<List<TimelineItem>> search({String? userId, required String query}) async {
    final params = <String, dynamic>{'q': query};
    if (userId != null) params['user_id'] = userId;
    final response = await apiClient.get('/timeline/search', queryParameters: params);
    final json = response.data as Map<String, dynamic>;
    final items = (json['data'] as Map<String, dynamic>?)?['items'] as List? ?? [];
    return items.map((e) => TimelineItem.fromJson(e as Map<String, dynamic>)).toList();
  }

  Future<TimelineItem?> getById({String? userId, required String id}) async {
    final params = <String, dynamic>{};
    if (userId != null) params['user_id'] = userId;
    final response = await apiClient.get('/timeline/$id', queryParameters: params);
    final json = response.data as Map<String, dynamic>;
    final data = (json['data'] as Map<String, dynamic>?);
    if (data == null) return null;
    return TimelineItem.fromJson(data);
  }
}
