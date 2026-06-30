import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../network/api_client.dart';
import '../../core/auth/token_manager.dart';

final authServiceProvider = Provider<AuthService>((ref) {
  final apiClient = ref.read(apiClientProvider);
  final tokenManager = ref.read(tokenManagerProvider);
  return AuthService(apiClient: apiClient, tokenManager: tokenManager);
});

class AuthService {
  final ApiClient apiClient;
  final TokenManager tokenManager;

  AuthService({required this.apiClient, required this.tokenManager});

  Future<AuthResult> login(String email, String password) async {
    try {
      final response = await apiClient.post('/auth/login', data: {
        'email': email,
        'password': password,
      });
      final data = response.data as Map<String, dynamic>;
      await tokenManager.setTokens(
        accessToken: data['access_token'] as String,
        refreshToken: data['refresh_token'] as String?,
      );
      return AuthResult.success(userData: data['user']);
    } on Exception catch (e) {
      return AuthResult.failure(message: e.toString());
    }
  }

  Future<void> logout() async {
    await tokenManager.clear();
  }

  Future<bool> isAuthenticated() async {
    await tokenManager.init();
    return tokenManager.isAuthenticated;
  }
}

class AuthResult {
  final bool isSuccess;
  final String? message;
  final Map<String, dynamic>? userData;

  AuthResult._({required this.isSuccess, this.message, this.userData});

  factory AuthResult.success({Map<String, dynamic>? userData}) =>
      AuthResult._(isSuccess: true, userData: userData);

  factory AuthResult.failure({String? message}) =>
      AuthResult._(isSuccess: false, message: message);
}
