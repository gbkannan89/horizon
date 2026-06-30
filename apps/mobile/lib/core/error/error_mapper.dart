class AppError {
  final String code;
  final String message;
  final int? statusCode;
  final dynamic originalError;

  const AppError({
    required this.code,
    required this.message,
    this.statusCode,
    this.originalError,
  });

  factory AppError.network(String message) => AppError(code: 'NETWORK_ERROR', message: message);
  factory AppError.auth(String message) => AppError(code: 'AUTH_ERROR', message: message);
  factory AppError.server(String message, {int? statusCode}) => AppError(code: 'SERVER_ERROR', message: message, statusCode: statusCode);
  factory AppError.validation(String message) => AppError(code: 'VALIDATION_ERROR', message: message);
  factory AppError.unknown(String message) => AppError(code: 'UNKNOWN_ERROR', message: message);

  @override
  String toString() => 'AppError($code): $message';
}

class ErrorMapper {
  static AppError fromDioException(dynamic error) {
    if (error is Map && error['error'] is Map) {
      final e = error['error'] as Map;
      final code = e['code'] as String? ?? 'UNKNOWN';
      final message = e['message'] as String? ?? 'An unexpected error occurred';
      return AppError(code: code, message: message);
    }
    return AppError.unknown(error.toString());
  }
}
