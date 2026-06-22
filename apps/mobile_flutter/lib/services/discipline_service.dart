import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';
import '../models/wishlist_item.dart';
import '../models/discipline_aggregates.dart';
import 'package:flutter/foundation.dart';
import 'api_service.dart';

class DisciplineService {
  static String get baseUrl => ApiService.baseUrl;

  Future<Map<String, String>> _getHeaders() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('jwt_token');
    return {
      'Content-Type': 'application/json',
      if (token != null) 'Authorization': 'Bearer $token',
    };
  }

  Future<List<WishlistItem>> getWishlist() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/discipline/wishlist'),
      headers: await _getHeaders(),
    );
    if (response.statusCode == 200) {
      final List<dynamic> data = jsonDecode(response.body);
      return data.map((json) => WishlistItem.fromJson(json)).toList();
    } else {
      throw Exception('Failed to load wishlist: ${response.statusCode}');
    }
  }

  Future<WishlistItem> addWishlistItem(String name, double amount, int lockDurationDays) async {
    final response = await http.post(
      Uri.parse('$baseUrl/api/discipline/wishlist'),
      headers: await _getHeaders(),
      body: jsonEncode({
        'name': name,
        'amount': amount,
        'lock_duration_days': lockDurationDays,
      }),
    );
    if (response.statusCode == 201) {
      return WishlistItem.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to add wishlist item: ${response.statusCode}');
    }
  }

  Future<WishlistItem> updateWishlistItemStatus(int id, String status) async {
    final response = await http.put(
      Uri.parse('$baseUrl/api/discipline/wishlist/$id'),
      headers: await _getHeaders(),
      body: jsonEncode({'status': status}),
    );
    if (response.statusCode == 200) {
      return WishlistItem.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to update wishlist item status: ${response.statusCode}');
    }
  }

  Future<FinancialGuardrails> getGuardrails() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/discipline/guardrails'),
      headers: await _getHeaders(),
    );
    if (response.statusCode == 200) {
      return FinancialGuardrails.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to load guardrails: ${response.statusCode}');
    }
  }

  Future<DebtRepaymentStrategy> getDebtStrategy() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/discipline/debt-strategy'),
      headers: await _getHeaders(),
    );
    if (response.statusCode == 200) {
      return DebtRepaymentStrategy.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to load debt strategy: ${response.statusCode}');
    }
  }

  Future<ZeroBasedBudget> getZeroBasedBudget() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/discipline/zero-based-budget'),
      headers: await _getHeaders(),
    );
    if (response.statusCode == 200) {
      return ZeroBasedBudget.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to load zero based budget: ${response.statusCode}');
    }
  }

  Future<PayYourselfFirst> getPayYourselfFirst() async {
    final response = await http.get(
      Uri.parse('$baseUrl/api/discipline/pay-yourself-first'),
      headers: await _getHeaders(),
    );
    if (response.statusCode == 200) {
      return PayYourselfFirst.fromJson(jsonDecode(response.body));
    } else {
      throw Exception('Failed to load pay yourself first data: ${response.statusCode}');
    }
  }
}
