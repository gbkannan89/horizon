import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import 'package:flutter/foundation.dart';

// ─────────────────────────────────────────────────────────────────────────────
// BASE URL CONFIGURATION
// Change the value that matches your current setup and hot-restart.
// ─────────────────────────────────────────────────────────────────────────────
class _Config {
  // ① Android Emulator (local dev — emulator talks to host via 10.0.2.2)
  static const String emulator = 'http://10.0.2.2:8000';

  // ② NUC Server on local Wi-Fi (real device / after deploying to NUC)
  //    Find NUC IP with: ip addr show   →  usually 192.168.x.x
  static const String nucLocal = 'http://192.168.0.100:8000';

  // ③ Cloudflare tunnel / public domain (after Cloudflare setup)
  static const String production = 'https://horizon.thaari.in';

  // ── ACTIVE ENVIRONMENT ── switch this line only ───────────────────────────
  static const String active = nucLocal; // ← change to production when ready
}

class ApiService {
  static String get baseUrl {
    if (kIsWeb) return 'http://127.0.0.1:8000';
    return _Config.active;
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
        return true;
      }
      print('Login failed: ${response.statusCode} ${response.body}');
      return false;
    } catch (e) {
      print('Login error: $e');
      return false;
    }
  }
  Future<bool> register(String name, String email, String password) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/api/auth/register'),
        headers: {'Content-Type': 'application/json'},
        body: json.encode({
          'name': name,
          'email': email,
          'password': password,
          'user_type': 'primary',
          'risk_profile': 'moderate'
        }),
      );
      if (response.statusCode == 201) {
        final data = json.decode(response.body);
        final prefs = await SharedPreferences.getInstance();
        await prefs.setString('jwt_token', data['access_token']);
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
  }

  // ── GET ───────────────────────────────────────────────────────────────────
  Future<dynamic> get(String endpoint) async {
    try {
      final response = await http.get(
        Uri.parse('$baseUrl$endpoint'),
        headers: await _getHeaders(),
      );
      if (response.statusCode == 200) return json.decode(response.body);
      throw Exception('GET $endpoint → ${response.statusCode}: ${response.body}');
    } catch (e) {
      print('ApiService.get error: $e');
      rethrow;
    }
  }

  // ── POST ──────────────────────────────────────────────────────────────────
  Future<dynamic> post(String endpoint, Map<String, dynamic> body) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl$endpoint'),
        headers: await _getHeaders(),
        body: json.encode(body),
      );
      if (response.statusCode >= 200 && response.statusCode < 300) {
        return json.decode(response.body);
      }
      throw Exception('POST $endpoint → ${response.statusCode}: ${response.body}');
    } catch (e) {
      print('ApiService.post error: $e');
      rethrow;
    }
  }

  // ── DELETE ────────────────────────────────────────────────────────────────
  Future<void> delete(String endpoint) async {
    try {
      final response = await http.delete(
        Uri.parse('$baseUrl$endpoint'),
        headers: await _getHeaders(),
      );
      if (response.statusCode < 200 || response.statusCode >= 300) {
        throw Exception('DELETE $endpoint → ${response.statusCode}: ${response.body}');
      }
    } catch (e) {
      print('ApiService.delete error: $e');
      rethrow;
    }
  }

  // ── PUT ───────────────────────────────────────────────────────────────────
  Future<dynamic> put(String endpoint, Map<String, dynamic> body) async {
    try {
      final response = await http.put(
        Uri.parse('$baseUrl$endpoint'),
        headers: await _getHeaders(),
        body: json.encode(body),
      );
      if (response.statusCode >= 200 && response.statusCode < 300) {
        return json.decode(response.body);
      }
      throw Exception('PUT $endpoint → ${response.statusCode}: ${response.body}');
    } catch (e) {
      print('ApiService.put error: $e');
      rethrow;
    }
  }
}
