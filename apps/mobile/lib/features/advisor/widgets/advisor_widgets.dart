import 'package:flutter/material.dart';
import 'package:flutter_markdown/flutter_markdown.dart';
import 'package:horizon_mobile/app/theme.dart';
import '../models/advisor_models.dart';

class ChatBubble extends StatelessWidget {
  final ChatMessage message;

  const ChatBubble({super.key, required this.message});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isUser = message.isUser;

    if (message.isLoading) {
      return _buildLoading(theme);
    }

    return Padding(
      padding: EdgeInsets.only(
        left: isUser ? AppSpacing.xl : 0,
        right: isUser ? 0 : AppSpacing.xl,
        bottom: AppSpacing.sm,
      ),
      child: Row(
        mainAxisAlignment: isUser ? MainAxisAlignment.end : MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (!isUser) ...[
            Container(
              width: 32, height: 32,
              decoration: BoxDecoration(
                color: theme.colorScheme.primaryContainer,
                borderRadius: BorderRadius.circular(AppRadius.sm),
              ),
              child: Icon(Icons.auto_awesome, size: 18, color: theme.colorScheme.primary),
            ),
            const SizedBox(width: AppSpacing.sm),
          ],
          Flexible(
            child: Container(
              padding: const EdgeInsets.all(AppSpacing.md),
              decoration: BoxDecoration(
                color: isUser ? theme.colorScheme.primary : theme.colorScheme.surfaceContainerHighest,
                borderRadius: BorderRadius.circular(AppRadius.md).copyWith(
                  bottomRight: isUser ? const Radius.circular(0) : null,
                  bottomLeft: isUser ? null : const Radius.circular(0),
                ),
              ),
              child: isUser
                  ? Text(message.content, style: TextStyle(color: theme.colorScheme.onPrimary))
                  : MarkdownBody(
                      data: message.content,
                      styleSheet: MarkdownStyleSheet(
                        p: TextStyle(color: theme.colorScheme.onSurface),
                        h1: TextStyle(color: theme.colorScheme.onSurface, fontWeight: FontWeight.bold),
                        h2: TextStyle(color: theme.colorScheme.onSurface, fontWeight: FontWeight.bold),
                        code: TextStyle(backgroundColor: theme.colorScheme.surfaceContainerLow),
                        codeblockDecoration: BoxDecoration(
                          color: theme.colorScheme.surfaceContainerLow,
                          borderRadius: BorderRadius.circular(AppRadius.sm),
                        ),
                      ),
                    ),
            ),
          ),
          if (isUser) ...[
            const SizedBox(width: AppSpacing.sm),
            Container(
              width: 32, height: 32,
              decoration: BoxDecoration(
                color: theme.colorScheme.primary,
                borderRadius: BorderRadius.circular(AppRadius.sm),
              ),
              child: Icon(Icons.person, size: 18, color: theme.colorScheme.onPrimary),
            ),
          ],
        ],
      ),
    );
  }

  Widget _buildLoading(ThemeData theme) {
    return Padding(
      padding: const EdgeInsets.only(right: AppSpacing.xl, bottom: AppSpacing.sm),
      child: Row(
        children: [
          Container(
            width: 32, height: 32,
            decoration: BoxDecoration(
              color: theme.colorScheme.primaryContainer,
              borderRadius: BorderRadius.circular(AppRadius.sm),
            ),
            child: Icon(Icons.auto_awesome, size: 18, color: theme.colorScheme.primary),
          ),
          const SizedBox(width: AppSpacing.sm),
          Container(
            padding: const EdgeInsets.all(AppSpacing.md),
            decoration: BoxDecoration(
              color: theme.colorScheme.surfaceContainerHighest,
              borderRadius: BorderRadius.circular(AppRadius.md).copyWith(bottomLeft: const Radius.circular(0)),
            ),
            child: const _AnimatedDots(),
          ),
        ],
      ),
    );
  }
}

class SuggestedPromptChip extends StatelessWidget {
  final SuggestedPrompt prompt;
  final VoidCallback onTap;

  const SuggestedPromptChip({super.key, required this.prompt, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return ActionChip(
      avatar: Icon(prompt.icon, size: AppIconSize.sm),
      label: Text(prompt.title),
      onPressed: onTap,
    );
  }
}

class ContextSummaryCard extends StatelessWidget {
  final AdvisorContext financialContext;

  const ContextSummaryCard({super.key, required this.financialContext});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      margin: const EdgeInsets.fromLTRB(AppSpacing.md, AppSpacing.sm, AppSpacing.md, 0),
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Row(
          children: [
            _stat(theme, 'Net Worth', '₹${_fmt(financialContext.netWorth)}', Colors.teal),
            _divider(theme),
            _stat(theme, 'Income', '₹${_fmt(financialContext.monthlyIncome)}', Colors.green),
            _divider(theme),
            _stat(theme, 'Health', '${financialContext.healthScore}', financialContext.healthScore >= 60 ? Colors.green : Colors.orange),
          ],
        ),
      ),
    );
  }

  Widget _stat(ThemeData theme, String label, String value, Color color) {
    return Expanded(
      child: Column(
        children: [
          Text(value, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold, color: color)),
          Text(label, style: theme.textTheme.labelSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        ],
      ),
    );
  }

  Widget _divider(ThemeData theme) {
    return Container(width: 1, height: 32, color: theme.colorScheme.outlineVariant);
  }

  String _fmt(double v) {
    if (v >= 10000000) return '${(v / 10000000).toStringAsFixed(1)}Cr';
    if (v >= 100000) return '${(v / 100000).toStringAsFixed(1)}L';
    if (v >= 1000) return '${(v / 1000).toStringAsFixed(1)}K';
    return v.toStringAsFixed(0);
  }
}

class _AnimatedDots extends StatefulWidget {
  const _AnimatedDots();
  @override
  State<_AnimatedDots> createState() => _AnimatedDotsState();
}

class _AnimatedDotsState extends State<_AnimatedDots> with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _anim1;
  late Animation<double> _anim2;
  late Animation<double> _anim3;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(vsync: this, duration: const Duration(milliseconds: 1200));
    _anim1 = Tween(begin: 0.3, end: 1.0).animate(CurvedAnimation(parent: _controller, curve: const Interval(0.0, 0.3, curve: Curves.easeInOut)));
    _anim2 = Tween(begin: 0.3, end: 1.0).animate(CurvedAnimation(parent: _controller, curve: const Interval(0.3, 0.6, curve: Curves.easeInOut)));
    _anim3 = Tween(begin: 0.3, end: 1.0).animate(CurvedAnimation(parent: _controller, curve: const Interval(0.6, 0.9, curve: Curves.easeInOut)));
    _controller.repeat(reverse: true);
  }

  @override
  void dispose() { _controller.dispose(); super.dispose(); }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return AnimatedBuilder(
      animation: _controller,
      builder: (context, _) => Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          _dot(theme, _anim1.value),
          const SizedBox(width: 4),
          _dot(theme, _anim2.value),
          const SizedBox(width: 4),
          _dot(theme, _anim3.value),
        ],
      ),
    );
  }

  Widget _dot(ThemeData theme, double opacity) {
    return Container(
      width: 8, height: 8,
      decoration: BoxDecoration(
        color: theme.colorScheme.onSurfaceVariant.withValues(alpha: opacity),
        shape: BoxShape.circle,
      ),
    );
  }
}
