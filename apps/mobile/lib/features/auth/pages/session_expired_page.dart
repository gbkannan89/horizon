import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';

class SessionExpiredPage extends StatelessWidget {
  const SessionExpiredPage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      body: Center(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
            Container(
              width: 80, height: 80,
              decoration: BoxDecoration(color: theme.colorScheme.errorContainer, borderRadius: BorderRadius.circular(AppRadius.lg)),
              child: Icon(Icons.timer_off, size: 40, color: theme.colorScheme.error),
            ),
            const SizedBox(height: AppSpacing.lg),
            Text('Session Expired', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: AppSpacing.sm),
              Text('Your session has timed out. Please sign in again to continue.', textAlign: TextAlign.center, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              const SizedBox(height: 32),
              FilledButton(
                onPressed: () => context.go('/login'),
                style: FilledButton.styleFrom(minimumSize: const Size(double.infinity, 52), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppRadius.md))),
                child: const Text('Sign In Again'),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
