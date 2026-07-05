import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
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
  Future<bool> register(String name, String email, String password, {String userType = 'salaried', String riskProfile = 'moderate', String phone = ''}) async {
    _isLoading = true;
    notifyListeners();

    final success = await _apiService.register(name, email, password, userType: userType, riskProfile: riskProfile, phone: phone);
    _isAuthenticated = success;
    
    if (success) await fetchUser();

    _isLoading = false;
    notifyListeners();
    return success;
  }

  Future<void> checkAuthStatus() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('jwt_token');
    if (token != null && token.isNotEmpty) {
      _isAuthenticated = true;
      await fetchUser();
    }
    notifyListeners();
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
    try {
      await _apiService.logout();
    } catch (_) {}
    _isAuthenticated = false;
    _user = null;
    notifyListeners();
  }

  Future<void> deleteAccount(String password) async {
    await _apiService.deleteAccount(password);
    _isAuthenticated = false;
    _user = null;
    notifyListeners();
  }
}
