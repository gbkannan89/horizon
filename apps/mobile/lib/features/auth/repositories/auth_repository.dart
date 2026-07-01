import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/core/auth/token_manager.dart';

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  return AuthRepository(
    apiClient: ref.read(apiClientProvider),
    tokenManager: ref.read(tokenManagerProvider),
  );
});

class AuthRepository {
  final ApiClient apiClient;
  final TokenManager tokenManager;

  AuthRepository({required this.apiClient, required this.tokenManager});

  Future<LoginResult> login(String email, String password, {bool rememberMe = false}) async {
    try {
      final response = await apiClient.post('/auth/login', data: {
        'email': email, 'password': password,
      });
      final data = response.data as Map<String, dynamic>;
      await tokenManager.setTokens(
        accessToken: data['access_token'] as String,
        refreshToken: data['refresh_token'] as String?,
      );
      final user = data['user'] as Map<String, dynamic>?;
      return LoginResult.success(
        accessToken: data['access_token'] as String,
        user: user != null ? UserModel.fromJson(user) : null,
      );
    } on Exception catch (e) {
      return LoginResult.failure(message: _extractError(e));
    }
  }

  Future<ForgotPasswordResult> forgotPassword(String email) async {
    try {
      await apiClient.post('/auth/forgot-password', data: {'email': email});
      return ForgotPasswordResult.success();
    } on Exception catch (e) {
      return ForgotPasswordResult.failure(message: _extractError(e));
    }
  }

  Future<ResetPasswordResult> resetPassword(String token, String password) async {
    try {
      await apiClient.post('/auth/reset-password', data: {
        'token': token, 'password': password,
      });
      return ResetPasswordResult.success();
    } on Exception catch (e) {
      return ResetPasswordResult.failure(message: _extractError(e));
    }
  }

  Future<void> logout() async {
    try { await apiClient.post('/auth/logout'); } catch (_) {}
    await tokenManager.clear();
  }

  Future<bool> isAuthenticated() async {
    await tokenManager.init();
    return tokenManager.isAuthenticated;
  }

  String _extractError(dynamic error) {
    if (error is DioException) {
      try {
        final data = error.response?.data;
        if (data is Map && data['error'] is Map) {
          return (data['error'] as Map)['message'] as String? ?? 'An error occurred';
        }
      } catch (_) {}
      if (error.type == DioExceptionType.connectionTimeout || error.type == DioExceptionType.receiveTimeout) {
        return 'Connection timed out. Check that the backend is running.';
      }
      if (error.type == DioExceptionType.connectionError) {
        return 'Cannot connect to backend. Is the server running on port 8081?';
      }
    }
    return 'Connection error. Please try again.';
  }
}

class LoginResult {
  final bool isSuccess;
  final String? message;
  final String? accessToken;
  final UserModel? user;

  LoginResult._({required this.isSuccess, this.message, this.accessToken, this.user});

  factory LoginResult.success({String? accessToken, UserModel? user}) =>
      LoginResult._(isSuccess: true, accessToken: accessToken, user: user);

  factory LoginResult.failure({String? message}) =>
      LoginResult._(isSuccess: false, message: message);
}

class ForgotPasswordResult {
  final bool isSuccess;
  final String? message;
  ForgotPasswordResult._({required this.isSuccess, this.message});
  factory ForgotPasswordResult.success() => ForgotPasswordResult._(isSuccess: true);
  factory ForgotPasswordResult.failure({String? message}) => ForgotPasswordResult._(isSuccess: false, message: message);
}

class ResetPasswordResult {
  final bool isSuccess;
  final String? message;
  ResetPasswordResult._({required this.isSuccess, this.message});
  factory ResetPasswordResult.success() => ResetPasswordResult._(isSuccess: true);
  factory ResetPasswordResult.failure({String? message}) => ResetPasswordResult._(isSuccess: false, message: message);
}

class UserModel {
  final String id;
  final String email;
  final String? name;

  const UserModel({required this.id, required this.email, this.name});

  factory UserModel.fromJson(Map<String, dynamic> json) => UserModel(
    id: json['id'] as String? ?? '',
    email: json['email'] as String? ?? '',
    name: json['name'] as String?,
  );
}
