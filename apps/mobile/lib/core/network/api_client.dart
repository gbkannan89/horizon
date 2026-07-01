import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../config/app_config.dart';
import '../auth/token_manager.dart';
import '../connectivity/connectivity_monitor.dart';

final apiClientProvider = Provider<ApiClient>((ref) {
  final tokenManager = ref.read(tokenManagerProvider);
  final connectivityMonitor = ref.read(connectivityMonitorProvider);
  return ApiClient(tokenManager: tokenManager, connectivityMonitor: connectivityMonitor);
});

class ApiClient {
  late final Dio dio;
  final TokenManager tokenManager;
  final ConnectivityMonitor connectivityMonitor;

  ApiClient({required this.tokenManager, required this.connectivityMonitor}) {
    final config = AppConfig.instance;
    dio = Dio(BaseOptions(
      baseUrl: config.apiBaseUrl,
      connectTimeout: Duration(milliseconds: config.connectTimeout),
      receiveTimeout: Duration(milliseconds: config.receiveTimeout),
      headers: {'Content-Type': 'application/json', 'Accept': 'application/json'},
    ));

    dio.interceptors.addAll([
      _AuthInterceptor(tokenManager),
      _CorrelationInterceptor(),
      _LogInterceptor(config.enableLogging),
      _RetryInterceptor(dio, tokenManager, connectivityMonitor),
    ]);
  }

  Future<Response<T>> get<T>(String path, {Map<String, dynamic>? queryParameters, Options? options}) =>
      dio.get<T>(path, queryParameters: queryParameters, options: options);

  Future<Response<T>> post<T>(String path, {dynamic data, Map<String, dynamic>? queryParameters, Options? options}) =>
      dio.post<T>(path, data: data, queryParameters: queryParameters, options: options);

  Future<Response<T>> put<T>(String path, {dynamic data, Map<String, dynamic>? queryParameters, Options? options}) =>
      dio.put<T>(path, data: data, queryParameters: queryParameters, options: options);

  Future<Response<T>> patch<T>(String path, {dynamic data, Map<String, dynamic>? queryParameters, Options? options}) =>
      dio.patch<T>(path, data: data, queryParameters: queryParameters, options: options);

  Future<Response<T>> delete<T>(String path, {Map<String, dynamic>? queryParameters, Options? options}) =>
      dio.delete<T>(path, queryParameters: queryParameters, options: options);
}

class _AuthInterceptor extends Interceptor {
  final TokenManager tokenManager;
  _AuthInterceptor(this.tokenManager);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    final token = tokenManager.accessToken;
    if (token != null) {
      options.headers['Authorization'] = 'Bearer $token';
    }
    handler.next(options);
  }
}

class _CorrelationInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    options.headers['X-Correlation-ID'] = '${DateTime.now().millisecondsSinceEpoch}-${_randomString(8)}';
    handler.next(options);
  }

  String _randomString(int length) {
    const chars = 'abcdefghijklmnopqrstuvwxyz0123456789';
    return List.generate(length, (_) => chars[DateTime.now().microsecondsSinceEpoch % chars.length]).join();
  }
}

class _LogInterceptor extends Interceptor {
  final bool enabled;
  _LogInterceptor(this.enabled);

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    if (enabled) print('[API] ${options.method} ${options.path}');
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    if (enabled) print('[API] ${response.statusCode} ${response.requestOptions.path}');
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (enabled) print('[API] ERROR ${err.response?.statusCode} ${err.requestOptions.path}: ${err.message}');
    handler.next(err);
  }
}

class _RetryInterceptor extends Interceptor {
  final Dio dio;
  final TokenManager tokenManager;
  final ConnectivityMonitor connectivityMonitor;
  int _retryCount = 0;
  static const _maxRetries = 1;

  _RetryInterceptor(this.dio, this.tokenManager, this.connectivityMonitor);

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response?.statusCode == 401 && _retryCount < _maxRetries) {
      _retryCount++;
      final refreshTokenValue = tokenManager.refreshToken;
      if (refreshTokenValue == null) {
        await tokenManager.clear();
        handler.next(err);
        return;
      }

      try {
        final refreshDio = Dio(BaseOptions(baseUrl: dio.options.baseUrl));
        final response = await refreshDio.post(
          '/auth/refresh',
          data: {'refresh_token': refreshTokenValue},
          options: Options(
            headers: {'Content-Type': 'application/json', 'Accept': 'application/json'},
          ),
        );

        if (response.statusCode == 200) {
          final data = response.data as Map<String, dynamic>;
          await tokenManager.setTokens(
            accessToken: data['access_token'] as String,
            refreshToken: data['refresh_token'] as String?,
          );

          final opts = err.requestOptions;
          opts.headers['Authorization'] = 'Bearer ${tokenManager.accessToken}';
          final retryResponse = await dio.fetch(opts);
          _retryCount = 0;
          handler.resolve(retryResponse);
          return;
        }
      } catch (_) {
        // Refresh failed — network or server error
      }

      await tokenManager.clear();
    }
    handler.next(err);
  }
}
