import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/features/auth/providers/auth_form_state.dart';
import 'package:horizon_mobile/features/auth/repositories/auth_repository.dart';

class ForgotPasswordPage extends ConsumerStatefulWidget {
  const ForgotPasswordPage({super.key});
  @override
  ConsumerState<ForgotPasswordPage> createState() => _ForgotPasswordPageState();
}

class _ForgotPasswordPageState extends ConsumerState<ForgotPasswordPage> {
  bool _sent = false;

  Future<void> _handleSubmit() async {
    final notifier = ref.read(authFormStateProvider.notifier);
    if (!notifier.validateForgotPassword()) return;
    notifier.setLoading(true);
    final state = ref.read(authFormStateProvider);
    final repo = ref.read(authRepositoryProvider);
    final result = await repo.forgotPassword(state.email);
    notifier.setLoading(false);
    if (result.isSuccess && mounted) {
      setState(() => _sent = true);
    } else if (mounted) {
      notifier.setError(result.message ?? 'Failed to send reset email');
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final formState = ref.watch(authFormStateProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Reset Password')),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.lg),
          child: _sent ? _buildSuccess(theme) : _buildForm(theme, formState),
        ),
      ),
    );
  }

  Widget _buildSuccess(ThemeData theme) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
            Container(
              width: 80, height: 80,
              decoration: BoxDecoration(color: theme.colorScheme.primaryContainer, borderRadius: BorderRadius.circular(AppRadius.lg)),
              child: Icon(Icons.check_circle, size: 40, color: theme.colorScheme.primary),
            ),
            const SizedBox(height: AppSpacing.lg),
            Text('Check your email', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: AppSpacing.sm),
          Text('If an account exists with that email, we\'ve sent password reset instructions.', textAlign: TextAlign.center, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          const SizedBox(height: 32),
          FilledButton(
            onPressed: () {
              ref.read(authFormStateProvider.notifier).setStep(AuthFormStep.login);
              context.pop();
            },
            child: const Text('Back to Sign In'),
          ),
        ],
      ),
    );
  }

  Widget _buildForm(ThemeData theme, AuthFormState formState) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const SizedBox(height: AppSpacing.xl),
        Text('Enter your email', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: AppSpacing.sm),
        Text('We\'ll send you a link to reset your password', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        const SizedBox(height: AppSpacing.xl),
        TextFormField(
          initialValue: formState.email,
          onChanged: (v) => ref.read(authFormStateProvider.notifier).setEmail(v),
          decoration: InputDecoration(labelText: 'Email', prefixIcon: const Icon(Icons.email_outlined), errorText: formState.emailError),
          keyboardType: TextInputType.emailAddress,
        ),
        if (formState.generalError != null) Padding(
          padding: const EdgeInsets.only(top: 16),
          child: Text(formState.generalError!, style: TextStyle(color: theme.colorScheme.error)),
        ),
        const SizedBox(height: 24),
        FilledButton(
          onPressed: formState.isLoading ? null : _handleSubmit,
          style: FilledButton.styleFrom(minimumSize: const Size(double.infinity, 52), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppRadius.md))),
          child: formState.isLoading
              ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
              : const Text('Send Reset Link', style: TextStyle(fontSize: 16)),
        ),
      ],
    );
  }
}
