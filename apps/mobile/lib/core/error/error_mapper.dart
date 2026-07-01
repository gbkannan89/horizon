import 'package:dio/dio.dart';

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
    if (error is DioException) {
      final statusCode = error.response?.statusCode;

      if (statusCode == 401) {
        final msg = _extractMessage(error) ?? 'Session expired. Please sign in again.';
        return AppError.auth(msg);
      }
      if (statusCode == 403) {
        return AppError.auth('You do not have permission to perform this action.');
      }
      if (statusCode == 404) {
        return AppError.server('Resource not found.', statusCode: statusCode);
      }
      if (statusCode != null && statusCode >= 500) {
        final msg = _extractMessage(error) ?? 'Something went wrong on our end. Please try again.';
        return AppError.server(msg, statusCode: statusCode);
      }
      if (statusCode != null && statusCode >= 400) {
        final msg = _extractMessage(error) ?? 'Invalid request. Please check your input.';
        return AppError.validation(msg);
      }

      switch (error.type) {
        case DioExceptionType.connectionTimeout:
        case DioExceptionType.sendTimeout:
        case DioExceptionType.receiveTimeout:
          return AppError.network('Connection timed out. Check that the server is running.');
        case DioExceptionType.connectionError:
          return AppError.network('Cannot connect to server. Check your internet connection.');
        case DioExceptionType.badResponse:
          return AppError.server('Unexpected server response.', statusCode: statusCode);
        default:
          return AppError.network('A network error occurred. Please try again.');
      }
    }

    if (error is Map && error['error'] is Map) {
      final e = error['error'] as Map;
      final code = e['code'] as String? ?? 'UNKNOWN';
      final message = e['message'] as String? ?? 'An unexpected error occurred';
      return AppError(code: code, message: message);
    }

    return AppError.unknown(error.toString());
  }

  static String? _extractMessage(DioException error) {
    try {
      final data = error.response?.data;
      if (data is Map && data['error'] is Map) {
        return (data['error'] as Map)['message'] as String?;
      }
    } catch (_) {}
    return null;
  }
}
