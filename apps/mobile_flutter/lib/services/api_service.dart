import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import '../main.dart';
import '../utils/ui_utils.dart';
import '../screens/login_screen.dart';

// ─────────────────────────────────────────────────────────────────────────────
// BASE URL CONFIGURATION
// Change the value that matches your current setup and hot-restart.
// ─────────────────────────────────────────────────────────────────────────────
class _Config {
  // ① Android Emulator (local dev — emulator talks to host via 10.0.2.2)
  static const String emulator = 'http://10.0.2.2:8000';

  // ② NUC Server on local Wi-Fi (real device / after deploying to NUC)
  static const String nucLocal = 'http://192.168.0.100:8000';

  // ③ Cloudflare tunnel / public domain (after Cloudflare setup)
  static const String production = 'https://horizon.thaari.in';
}

class ApiService {
  static String get baseUrl {
    if (kIsWeb) return 'http://localhost:8000';
    if (kReleaseMode) {
      return _Config.production;
    } else {
      return _Config.emulator; // Using 10.0.2.2 for Android Emulator locally
    }
  }

  Future<String?> _getToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('jwt_token');
  }

  Future<Map<String, String>> _getHeaders() async {
    final token = await _getToken();
    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }

  // ── AUTH ──────────────────────────────────────────────────────────────────
  Future<bool> login(String username, String password) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/api/auth/login'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({'email': username, 'password': password}),
      );
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        final prefs = await SharedPreferences.getInstance();
        await prefs.setString('jwt_token', data['access_token']);
        if (data['refresh_token'] != null) {
          await prefs.setString('refresh_token', data['refresh_token']);
        }
        return true;
      }
      print('Login failed: ${response.statusCode} ${response.body}');
      return false;
    } catch (e) {
      print('Login error: $e');
      return false;
    }
  }
  Future<bool> register(String name, String email, String password, {String userType = 'salaried', String riskProfile = 'moderate', String phone = ''}) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/api/auth/register'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'name': name,
          'email': email,
          'password': password,
          'user_type': userType,
          'risk_profile': riskProfile,
          'phone': phone.isEmpty ? null : phone,
        }),
      );
      if (response.statusCode == 201) {
        final data = json.decode(response.body);
        final prefs = await SharedPreferences.getInstance();
        await prefs.setString('jwt_token', data['access_token']);
        if (data['refresh_token'] != null) {
          await prefs.setString('refresh_token', data['refresh_token']);
        }
        return true;
      }
      print('Register failed: ${response.statusCode} ${response.body}');
      return false;
    } catch (e) {
      print('Register error: $e');
      return false;
    }
  }

  Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('jwt_token');
    await prefs.remove('refresh_token');
  }

  void _handleError(http.Response response) {
    if (response.statusCode == 401) {
      logout();
      if (navigatorKey.currentContext != null) {
        Navigator.of(navigatorKey.currentContext!).pushAndRemoveUntil(
          MaterialPageRoute(builder: (_) => const LoginScreen()),
          (route) => false,
        );
        UiUtils.showSnack(navigatorKey.currentContext!, 'Session expired. Please log in again.', isError: true);
      }
    } else {
      if (navigatorKey.currentContext != null) {
        UiUtils.showSnack(navigatorKey.currentContext!, 'Error: ${response.statusCode} - ${response.body}', isError: true);
      }
    }
    throw Exception('API Error: ${response.statusCode} - ${response.body}');
  }

  // ── REFRESH TOKEN ──────────────────────────────────────────────────────────
  Future<bool> _refreshToken() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final refreshToken = prefs.getString('refresh_token');
      if (refreshToken == null) return false;

      final response = await http.post(
        Uri.parse('$baseUrl/api/auth/refresh'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({'refresh_token': refreshToken}),
      );

      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        await prefs.setString('jwt_token', data['access_token']);
        if (data['refresh_token'] != null) {
          await prefs.setString('refresh_token', data['refresh_token']);
        }
        return true;
      }
      return false;
    } catch (e) {
      return false;
    }
  }

  // ── GET ───────────────────────────────────────────────────────────────────
  Future<dynamic> get(String endpoint) async {
    try {
      var response = await http.get(
        Uri.parse('$baseUrl$endpoint'),
        headers: await _getHeaders(),
      );
      
      if (response.statusCode == 401) {
        final refreshed = await _refreshToken();
        if (refreshed) {
          response = await http.get(
            Uri.parse('$baseUrl$endpoint'),
            headers: await _getHeaders(),
          );
        }
      }

      if (response.statusCode == 200) return json.decode(response.body);
      _handleError(response);
    } catch (e) {
      print('ApiService.get error: $e');
      rethrow;
    }
  }

  // ── POST ──────────────────────────────────────────────────────────────────
  Future<dynamic> post(String endpoint, Map<String, dynamic> body) async {
    try {
      var response = await http.post(
        Uri.parse('$baseUrl$endpoint'),
        headers: await _getHeaders(),
        body: json.encode(body),
      );
      
      if (response.statusCode == 401) {
        final refreshed = await _refreshToken();
        if (refreshed) {
          response = await http.post(
            Uri.parse('$baseUrl$endpoint'),
            headers: await _getHeaders(),
            body: json.encode(body),
          );
        }
      }

      if (response.statusCode >= 200 && response.statusCode < 300) {
        return json.decode(response.body);
      }
      _handleError(response);
    } catch (e) {
      print('ApiService.post error: $e');
      rethrow;
    }
  }

  // ── DELETE ────────────────────────────────────────────────────────────────
  Future<void> delete(String endpoint) async {
    try {
      var response = await http.delete(
        Uri.parse('$baseUrl$endpoint'),
        headers: await _getHeaders(),
      );
      
      if (response.statusCode == 401) {
        final refreshed = await _refreshToken();
        if (refreshed) {
          response = await http.delete(
            Uri.parse('$baseUrl$endpoint'),
            headers: await _getHeaders(),
          );
        }
      }

      if (response.statusCode < 200 || response.statusCode >= 300) {
        _handleError(response);
      }
    } catch (e) {
      print('ApiService.delete error: $e');
      rethrow;
    }
  }

  // ── PUT ───────────────────────────────────────────────────────────────────
  Future<dynamic> put(String endpoint, Map<String, dynamic> body) async {
    try {
      var response = await http.put(
        Uri.parse('$baseUrl$endpoint'),
        headers: await _getHeaders(),
        body: json.encode(body),
      );
      
      if (response.statusCode == 401) {
        final refreshed = await _refreshToken();
        if (refreshed) {
          response = await http.put(
            Uri.parse('$baseUrl$endpoint'),
            headers: await _getHeaders(),
            body: json.encode(body),
          );
        }
      }

      if (response.statusCode >= 200 && response.statusCode < 300) {
        return json.decode(response.body);
      }
      _handleError(response);
    } catch (e) {
      print('ApiService.put error: $e');
      rethrow;
    }
  }

  // ── MULTIPART UPLOAD ──────────────────────────────────────────────────────
  Future<dynamic> uploadFile(String endpoint, String fileField, String filePath) async {
    try {
      final request = http.MultipartRequest('POST', Uri.parse('$baseUrl$endpoint'));
      
      // Add Headers
      final token = await _getToken();
      if (token != null) {
        request.headers['Authorization'] = 'Bearer $token';
      }

      // Add File
      request.files.add(await http.MultipartFile.fromPath(fileField, filePath));

      final streamedResponse = await request.send();
      final response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode >= 200 && response.statusCode < 300) {
        return json.decode(response.body);
      }
      _handleError(response);
    } catch (e) {
      print('ApiService.uploadFile error: $e');
      rethrow;
    }
  }

  // ── FILE DOWNLOAD ──────────────────────────────────────────────────────────
  Future<List<int>> downloadFile(String endpoint) async {
    try {
      final uri = Uri.parse('$baseUrl$endpoint');
      final request = http.Request('GET', uri);
      final token = await _getToken();
      if (token != null) {
        request.headers['Authorization'] = 'Bearer $token';
      }
      final streamedResponse = await request.send();
      final response = await http.Response.fromStream(streamedResponse);
      if (response.statusCode >= 200 && response.statusCode < 300) {
        return response.bodyBytes;
      }
      _handleError(response);
      return [];
    } catch (e) {
      print('ApiService.downloadFile error: $e');
      rethrow;
    }
  }
}
