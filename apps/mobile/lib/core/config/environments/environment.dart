class Environment {
  final String name;
  final String apiBaseUrl;
  final bool enableLogging;

  const Environment._({
    required this.name,
    required this.apiBaseUrl,
    required this.enableLogging,
  });

  static const development = Environment._(
    name: 'development',
    apiBaseUrl: 'http://10.0.2.2:8081/api/v1',
    enableLogging: true,
  );

  static const staging = Environment._(
    name: 'staging',
    apiBaseUrl: 'https://staging.horizon.app/api/v1',
    enableLogging: true,
  );

  static const production = Environment._(
    name: 'production',
    apiBaseUrl: 'https://api.horizon.app/api/v1',
    enableLogging: false,
  );

  static Environment current = development;
}
