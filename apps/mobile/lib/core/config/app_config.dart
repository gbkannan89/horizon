import 'environments/environment.dart';

export 'environments/environment.dart';

class AppConfig {
  final String apiBaseUrl;
  final int connectTimeout;
  final int receiveTimeout;
  final bool useHttps;
  final String appVersion;
  final bool enableLogging;
  final bool enableOfflineCache;

  AppConfig._({
    required this.apiBaseUrl,
    required this.connectTimeout,
    required this.receiveTimeout,
    required this.useHttps,
    required this.appVersion,
    required this.enableLogging,
    required this.enableOfflineCache,
  });

  static AppConfig _instance = AppConfig._(
    apiBaseUrl: 'http://localhost:8081/api/v1',
    connectTimeout: 15000,
    receiveTimeout: 30000,
    useHttps: false,
    appVersion: '0.1.0',
    enableLogging: true,
    enableOfflineCache: true,
  );

  static AppConfig get instance => _instance;

  static void init(Environment env) {
    _instance = AppConfig._(
      apiBaseUrl: env.apiBaseUrl,
      connectTimeout: 15000,
      receiveTimeout: 30000,
      useHttps: env.apiBaseUrl.startsWith('https'),
      appVersion: '0.1.0',
      enableLogging: env.enableLogging,
      enableOfflineCache: true,
    );
  }
}
