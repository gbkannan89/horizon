import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/features/auth/providers/auth_form_state.dart';
import 'package:horizon_mobile/features/auth/repositories/auth_repository.dart';

class ResetPasswordPage extends ConsumerStatefulWidget {
  final String token;
  const ResetPasswordPage({super.key, required this.token});
  @override
  ConsumerState<ResetPasswordPage> createState() => _ResetPasswordPageState();
}

class _ResetPasswordPageState extends ConsumerState<ResetPasswordPage> {
  bool _success = false;

  Future<void> _handleSubmit() async {
    final notifier = ref.read(authFormStateProvider.notifier);
    if (!notifier.validateResetPassword()) return;
    notifier.setLoading(true);
    final state = ref.read(authFormStateProvider);
    final repo = ref.read(authRepositoryProvider);
    final result = await repo.resetPassword(widget.token, state.password);
    notifier.setLoading(false);
    if (result.isSuccess && mounted) {
      setState(() => _success = true);
    } else if (mounted) {
      notifier.setError(result.message ?? 'Failed to reset password');
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final formState = ref.watch(authFormStateProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Set New Password')),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.lg),
          child: _success ? _buildSuccess(theme) : _buildForm(theme, formState),
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
          Text('Password reset successful', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: 8),
          Text('Your password has been updated. You can now sign in with your new password.', textAlign: TextAlign.center, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          const SizedBox(height: 32),
          FilledButton(onPressed: () => context.go('/login'), child: const Text('Sign In')),
        ],
      ),
    );
  }

  Widget _buildForm(ThemeData theme, AuthFormState formState) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const SizedBox(height: AppSpacing.xl),
        Text('Choose a new password', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: AppSpacing.sm),
        Text('Must be at least 8 characters', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        const SizedBox(height: AppSpacing.xl),
        TextFormField(
          onChanged: (v) => ref.read(authFormStateProvider.notifier).setPassword(v),
          obscureText: formState.obscurePassword,
          decoration: InputDecoration(
            labelText: 'New Password',
            prefixIcon: const Icon(Icons.lock_outlined),
            suffixIcon: IconButton(
              icon: Icon(formState.obscurePassword ? Icons.visibility_outlined : Icons.visibility_off_outlined),
              onPressed: () => ref.read(authFormStateProvider.notifier).toggleObscurePassword(),
            ),
            errorText: formState.passwordError,
          ),
        ),
        const SizedBox(height: 16),
        TextFormField(
          onChanged: (v) => ref.read(authFormStateProvider.notifier).setConfirmPassword(v),
          obscureText: true,
          decoration: InputDecoration(labelText: 'Confirm Password', prefixIcon: const Icon(Icons.lock_outlined), errorText: formState.confirmPasswordError),
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
              : const Text('Reset Password', style: TextStyle(fontSize: 16)),
        ),
      ],
    );
  }
}
