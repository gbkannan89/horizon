import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../state/advisor_state.dart';
import '../widgets/advisor_widgets.dart';

class AdvisorPage extends ConsumerStatefulWidget {
  const AdvisorPage({super.key});
  @override
  ConsumerState<AdvisorPage> createState() => _AdvisorPageState();
}

class _AdvisorPageState extends ConsumerState<AdvisorPage> {
  final _inputCtrl = TextEditingController();
  final _scrollCtrl = ScrollController();

  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(advisorProvider.notifier).init());
  }

  @override
  void dispose() {
    _inputCtrl.dispose();
    _scrollCtrl.dispose();
    super.dispose();
  }

  void _sendMessage(String? text) {
    final message = text ?? _inputCtrl.text.trim();
    if (message.isEmpty) return;
    _inputCtrl.clear();
    ref.read(advisorProvider.notifier).sendMessage(message);
    _scrollToBottom();
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollCtrl.hasClients) {
        _scrollCtrl.animateTo(
          _scrollCtrl.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(advisorProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Advisor'),
        actions: [
          if (state.messages.isNotEmpty)
            IconButton(
              icon: const Icon(Icons.delete_outline),
              tooltip: 'Clear conversation',
              onPressed: () => ref.read(advisorProvider.notifier).clearConversation(),
            ),
        ],
      ),
      body: Column(
        children: [
          if (state.context != null)
            ContextSummaryCard(financialContext: state.context!),
          Expanded(child: _buildBody(theme, state)),
          _buildInputBar(theme, state),
        ],
      ),
    );
  }

  Widget _buildBody(ThemeData theme, AdvisorState state) {
    switch (state.status) {
      case AdvisorStatus.initial:
      case AdvisorStatus.loading:
        return const SharedLoadingView(message: 'Loading your financial context...');
      case AdvisorStatus.ready:
      case AdvisorStatus.thinking:
        if (state.messages.isEmpty) {
          return _buildWelcome(theme, state);
        }
        return _buildChat(theme, state);
      case AdvisorStatus.error:
        return SharedErrorView(
          message: state.error,
          onRetry: () => ref.read(advisorProvider.notifier).init(),
        );
    }
  }

  Widget _buildWelcome(ThemeData theme, AdvisorState state) {
    return ListView(
      padding: const EdgeInsets.all(AppSpacing.md),
      children: [
        const SizedBox(height: AppSpacing.xl),
        Container(
          width: 64, height: 64,
          decoration: BoxDecoration(
            color: theme.colorScheme.primaryContainer,
            borderRadius: BorderRadius.circular(AppRadius.lg),
          ),
          child: Icon(Icons.auto_awesome, size: 32, color: theme.colorScheme.primary),
        ),
        const SizedBox(height: AppSpacing.md),
        Text('Hello, I\'m your financial advisor', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: AppSpacing.sm),
        Text('Ask me anything about your finances, goals, or investments.', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        const SizedBox(height: AppSpacing.lg),
        Text('Suggested questions', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
        const SizedBox(height: AppSpacing.sm),
        Wrap(
          spacing: AppSpacing.sm,
          runSpacing: AppSpacing.sm,
          children: state.suggestedPrompts.map((p) =>
            SuggestedPromptChip(prompt: p, onTap: () => _sendMessage(p.prompt)),
          ).toList(),
        ),
      ],
    );
  }

  Widget _buildChat(ThemeData theme, AdvisorState state) {
    return ListView.builder(
      controller: _scrollCtrl,
      padding: const EdgeInsets.all(AppSpacing.md),
      itemCount: state.messages.length,
      itemBuilder: (context, index) {
        final msg = state.messages[index];
        return ChatBubble(message: msg);
      },
    );
  }

  Widget _buildInputBar(ThemeData theme, AdvisorState state) {
    final isThinking = state.status == AdvisorStatus.thinking;
    return Container(
      padding: const EdgeInsets.fromLTRB(AppSpacing.md, AppSpacing.sm, AppSpacing.sm, AppSpacing.md),
      decoration: BoxDecoration(
        color: theme.colorScheme.surface,
        border: Border(top: BorderSide(color: theme.colorScheme.outlineVariant)),
      ),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _inputCtrl,
              enabled: !isThinking,
              decoration: InputDecoration(
                hintText: isThinking ? 'Thinking...' : 'Ask your financial advisor...',
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(AppRadius.lg)),
                contentPadding: const EdgeInsets.symmetric(horizontal: AppSpacing.md, vertical: AppSpacing.sm),
              ),
              textInputAction: TextInputAction.send,
              onSubmitted: isThinking ? null : _sendMessage,
            ),
          ),
          const SizedBox(width: AppSpacing.sm),
          IconButton.filled(
            onPressed: isThinking ? null : () => _sendMessage(null),
            icon: isThinking
                ? const SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                : const Icon(Icons.send),
          ),
        ],
      ),
    );
  }
}
