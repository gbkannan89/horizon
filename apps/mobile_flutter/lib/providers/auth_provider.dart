import 'package:flutter/material.dart';
import '../services/api_service.dart';

class AuthProvider extends ChangeNotifier {
  final ApiService _apiService = ApiService();
  bool _isAuthenticated = false;
  bool _isLoading = false;

  bool get isAuthenticated => _isAuthenticated;
  bool get isLoading => _isLoading;

  Map<String, dynamic>? _user;
  Map<String, dynamic>? get user => _user;

  Future<bool> login(String username, String password) async {
    _isLoading = true;
    notifyListeners();

    final success = await _apiService.login(username, password);
    _isAuthenticated = success;
    
    if (success) await fetchUser();

    _isLoading = false;
    notifyListeners();
    return success;
  }
  Future<bool> register(String name, String email, String password) async {
    _isLoading = true;
    notifyListeners();

    final success = await _apiService.register(name, email, password);
    _isAuthenticated = success;
    
    if (success) await fetchUser();

    _isLoading = false;
    notifyListeners();
    return success;
  }

  Future<void> fetchUser() async {
    try {
      final data = await _apiService.get('/api/auth/me');
      _user = data;
      notifyListeners();
    } catch (e) {
      print('Failed to fetch user: $e');
    }
  }

  Future<void> logout() async {
    await _apiService.logout();
    _isAuthenticated = false;
    _user = null;
    notifyListeners();
  }
}
