import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/advisor_models.dart';
import '../repository/advisor_repository.dart';

final advisorProvider = StateNotifierProvider<AdvisorNotifier, AdvisorState>((ref) {
  return AdvisorNotifier(ref.read(advisorRepositoryProvider));
});

enum AdvisorStatus { initial, loading, ready, thinking, error }

class AdvisorState {
  final AdvisorStatus status;
  final List<ChatMessage> messages;
  final String? error;
  final AdvisorContext? context;
  final List<SuggestedPrompt> suggestedPrompts;
  final String? sessionId;
  final List<ProviderInfo> providers;
  final String? activeProvider;

  const AdvisorState({
    this.status = AdvisorStatus.initial,
    this.messages = const [],
    this.error,
    this.context,
    this.suggestedPrompts = const [],
    this.sessionId,
    this.providers = const [],
    this.activeProvider,
  });

  AdvisorState copyWith({
    AdvisorStatus? status,
    List<ChatMessage>? messages,
    String? error,
    AdvisorContext? context,
    List<SuggestedPrompt>? suggestedPrompts,
    String? sessionId,
    List<ProviderInfo>? providers,
    String? activeProvider,
    bool clearError = false,
  }) => AdvisorState(
    status: status ?? this.status,
    messages: messages ?? this.messages,
    error: clearError ? null : error ?? this.error,
    context: context ?? this.context,
    suggestedPrompts: suggestedPrompts ?? this.suggestedPrompts,
    sessionId: sessionId ?? this.sessionId,
    providers: providers ?? this.providers,
    activeProvider: activeProvider ?? this.activeProvider,
  );
}

class AdvisorNotifier extends StateNotifier<AdvisorState> {
  final AdvisorRepository _repo;
  int _msgCounter = 0;

  AdvisorNotifier(this._repo) : super(const AdvisorState());

  Future<void> init() async {
    state = state.copyWith(status: AdvisorStatus.loading);
    try {
      final results = await Future.wait([
        _repo.getContext(),
        _repo.getPrompts(),
        _repo.getProviders(),
      ]);
      state = state.copyWith(
        status: AdvisorStatus.ready,
        context: results[0] as AdvisorContext,
        suggestedPrompts: results[1] as List<SuggestedPrompt>,
        providers: results[2] as List<ProviderInfo>,
        activeProvider: (results[2] as List<ProviderInfo>).where((p) => p.isActive).firstOrNull?.name,
      );
    } catch (e) {
      state = state.copyWith(
        status: AdvisorStatus.ready,
        context: state.context ?? AdvisorContext(),
        suggestedPrompts: state.suggestedPrompts.isNotEmpty ? state.suggestedPrompts : _defaultPrompts(),
        error: e.toString(),
      );
    }
  }

  Future<void> sendMessage(String message) async {
    if (message.trim().isEmpty) return;

    final userMsg = ChatMessage(
      id: 'msg-${++_msgCounter}',
      role: 'user',
      content: message.trim(),
      timestamp: DateTime.now(),
    );

    final loadingMsg = ChatMessage(
      id: 'msg-loading',
      role: 'assistant',
      content: '',
      timestamp: DateTime.now(),
      isLoading: true,
    );

    state = state.copyWith(
      status: AdvisorStatus.thinking,
      messages: [...state.messages, userMsg, loadingMsg],
    );

    try {
      final resp = await _repo.sendMessage(
        sessionId: state.sessionId,
        message: message.trim(),
      );

      final aiMsg = ChatMessage(
        id: 'msg-${++_msgCounter}',
        role: 'assistant',
        content: resp.reply,
        timestamp: DateTime.now(),
      );

      final messages = [...state.messages, userMsg, aiMsg];

      state = state.copyWith(
        status: AdvisorStatus.ready,
        messages: messages,
        sessionId: resp.sessionId ?? state.sessionId,
      );
    } catch (e) {
      final messages = state.messages.where((m) => m.id != 'msg-loading').toList();
      final errorMsg = ChatMessage(
        id: 'msg-error',
        role: 'assistant',
        content: 'I\'m sorry, I encountered an error: $e\n\nPlease try again.',
        timestamp: DateTime.now(),
      );
      state = state.copyWith(
        status: AdvisorStatus.ready,
        messages: [...messages, errorMsg],
        error: e.toString(),
      );
    }
  }

  void clearConversation() {
    state = state.copyWith(
      messages: [],
      sessionId: null,
      status: AdvisorStatus.ready,
    );
  }

  void switchProvider(String name) async {
    try {
      await _repo.switchProvider(name);
      state = state.copyWith(activeProvider: name);
      final providers = state.providers.map((p) =>
        ProviderInfo(name: p.name, capabilities: p.capabilities, isActive: p.name == name)
      ).toList();
      state = state.copyWith(providers: providers);
    } catch (e) {
      state = state.copyWith(error: e.toString());
    }
  }

  List<SuggestedPrompt> _defaultPrompts() {
    return [
      SuggestedPrompt(id: '1', title: 'Financial Health', prompt: 'How is my financial health looking right now?', icon: Icons.favorite_outline),
      SuggestedPrompt(id: '2', title: 'Saving Tips', prompt: 'What can I do to save more money this month?', icon: Icons.savings_outlined),
      SuggestedPrompt(id: '3', title: 'Goal Progress', prompt: 'How are my goals tracking?', icon: Icons.flag_outlined),
      SuggestedPrompt(id: '4', title: 'Investment Ideas', prompt: 'Any investment ideas for my risk profile?', icon: Icons.trending_up),
    ];
  }
}
