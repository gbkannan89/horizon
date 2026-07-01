import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../storage/secure_storage_service.dart';

final tokenManagerProvider = Provider<TokenManager>((ref) {
  final storage = ref.read(secureStorageProvider);
  return TokenManager(storage: storage);
});

class TokenManager {
  final SecureStorageService storage;
  String? _accessToken;
  String? _refreshToken;

  TokenManager({required this.storage});

  String? get accessToken => _accessToken;
  String? get refreshToken => _refreshToken;

  Future<void> init() async {
    _accessToken = await storage.read('access_token');
    _refreshToken = await storage.read('refresh_token');
  }

  Future<void> setTokens({required String accessToken, String? refreshToken}) async {
    _accessToken = accessToken;
    _refreshToken = refreshToken;
    await storage.write('access_token', accessToken);
    if (refreshToken != null) await storage.write('refresh_token', refreshToken);
  }

  Future<void> clear() async {
    _accessToken = null;
    _refreshToken = null;
    await storage.delete('access_token');
    await storage.delete('refresh_token');
  }

  bool get isAuthenticated => _accessToken != null;
}
