import 'package:flutter/material.dart';

class ChatMessage {
  final String id;
  final String role;
  final String content;
  final DateTime timestamp;
  final bool isLoading;

  ChatMessage({
    required this.id,
    required this.role,
    required this.content,
    required this.timestamp,
    this.isLoading = false,
  });

  bool get isUser => role == 'user';
  bool get isAssistant => role == 'assistant';
}

class SuggestedPrompt {
  final String id;
  final String title;
  final String prompt;
  final String category;
  final IconData icon;

  SuggestedPrompt({
    required this.id,
    required this.title,
    required this.prompt,
    this.category = 'general',
    this.icon = Icons.lightbulb_outline,
  });
}

class AdvisorContext {
  final double netWorth;
  final double monthlyIncome;
  final double monthlyExpenses;
  final int healthScore;
  final int goalCount;
  final int activeRecommendations;
  final String? riskLevel;

  AdvisorContext({
    this.netWorth = 0,
    this.monthlyIncome = 0,
    this.monthlyExpenses = 0,
    this.healthScore = 0,
    this.goalCount = 0,
    this.activeRecommendations = 0,
    this.riskLevel,
  });

  factory AdvisorContext.fromJson(Map<String, dynamic> json) => AdvisorContext(
    netWorth: (json['net_worth'] as num?)?.toDouble() ?? 0,
    monthlyIncome: (json['monthly_income'] as num?)?.toDouble() ?? 0,
    monthlyExpenses: (json['monthly_expenses'] as num?)?.toDouble() ?? 0,
    healthScore: (json['health_score'] as num?)?.toInt() ?? 0,
    goalCount: (json['goal_count'] as num?)?.toInt() ?? 0,
    activeRecommendations: (json['active_recommendations'] as num?)?.toInt() ?? 0,
    riskLevel: json['risk_level'] as String?,
  );
}

class ChatResponse {
  final String reply;
  final double? confidence;
  final String? provider;
  final String? sessionId;

  ChatResponse({
    required this.reply,
    this.confidence,
    this.provider,
    this.sessionId,
  });

  factory ChatResponse.fromJson(Map<String, dynamic> json) {
    final data = json['data'] as Map<String, dynamic>? ?? json;
    return ChatResponse(
      reply: data['reply'] as String? ?? '',
      confidence: (data['confidence'] as num?)?.toDouble(),
      provider: data['provider'] as String?,
      sessionId: data['session_id'] as String?,
    );
  }
}

class ProviderInfo {
  final String name;
  final List<String> capabilities;
  final bool isActive;

  ProviderInfo({required this.name, this.capabilities = const [], this.isActive = false});

  factory ProviderInfo.fromJson(Map<String, dynamic> json) => ProviderInfo(
    name: json['name'] as String? ?? '',
    capabilities: (json['capabilities'] as List?)?.map((e) => e.toString()).toList() ?? [],
    isActive: json['is_active'] as bool? ?? false,
  );
}
