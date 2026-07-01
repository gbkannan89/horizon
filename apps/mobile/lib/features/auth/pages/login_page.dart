import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/features/auth/providers/auth_form_state.dart';
import 'package:horizon_mobile/features/auth/repositories/auth_repository.dart';
import 'package:horizon_mobile/shared/providers/auth_state.dart';

class LoginPage extends ConsumerStatefulWidget {
  const LoginPage({super.key});
  @override
  ConsumerState<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends ConsumerState<LoginPage> with SingleTickerProviderStateMixin {
  final _formKey = GlobalKey<FormState>();
  late AnimationController _fadeIn;

  @override
  void initState() {
    super.initState();
    _fadeIn = AnimationController(vsync: this, duration: const Duration(milliseconds: 600))..forward();
  }

  @override
  void dispose() {
    _fadeIn.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    final formState = ref.read(authFormStateProvider.notifier);
    if (!formState.validateLogin()) return;

    formState.setLoading(true);
    final state = ref.read(authFormStateProvider);
    final repo = ref.read(authRepositoryProvider);
    final result = await repo.login(state.email, state.password, rememberMe: state.rememberMe);
    formState.setLoading(false);

    if (result.isSuccess && mounted) {
      ref.read(authStateProvider.notifier).authenticated(
        userId: result.user?.id, email: result.user?.email,
      );
      context.go('/dashboard');
    } else if (mounted) {
      formState.setError(result.message ?? 'Login failed');
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final formState = ref.watch(authFormStateProvider);

    return Scaffold(
      body: SafeArea(
        child: FadeTransition(
          opacity: _fadeIn,
          child: Center(
            child: SingleChildScrollView(
              padding: const EdgeInsets.all(AppSpacing.lg),
              child: Form(
                key: _formKey,
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    Center(
                      child: Container(
                        width: 80, height: 80,
                        decoration: BoxDecoration(
                          color: theme.colorScheme.primaryContainer,
                          borderRadius: BorderRadius.circular(AppRadius.lg),
                        ),
                        child: Icon(Icons.trending_up, size: 40, color: theme.colorScheme.primary),
                      ),
                    ),
                    const SizedBox(height: AppSpacing.lg),
                    Text('Welcome back', textAlign: TextAlign.center, style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
                    const SizedBox(height: AppSpacing.sm),
                    Text('Sign in to your financial command center', textAlign: TextAlign.center, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                    const SizedBox(height: AppSpacing.xl),

                    // Email
                    TextFormField(
                      initialValue: formState.email,
                      onChanged: (v) => ref.read(authFormStateProvider.notifier).setEmail(v),
                      decoration: InputDecoration(
                        labelText: 'Email',
                        prefixIcon: const Icon(Icons.email_outlined),
                        errorText: formState.emailError,
                      ),
                      keyboardType: TextInputType.emailAddress,
                      textInputAction: TextInputAction.next,
                    ),
                    const SizedBox(height: 16),

                    // Password
                    TextFormField(
                      initialValue: formState.password,
                      onChanged: (v) => ref.read(authFormStateProvider.notifier).setPassword(v),
                      obscureText: formState.obscurePassword,
                      decoration: InputDecoration(
                        labelText: 'Password',
                        prefixIcon: const Icon(Icons.lock_outlined),
                        suffixIcon: IconButton(
                          icon: Icon(formState.obscurePassword ? Icons.visibility_outlined : Icons.visibility_off_outlined),
                          onPressed: () => ref.read(authFormStateProvider.notifier).toggleObscurePassword(),
                        ),
                        errorText: formState.passwordError,
                      ),
                      textInputAction: TextInputAction.done,
                      onFieldSubmitted: (_) => _handleLogin(),
                    ),
                    const SizedBox(height: 8),

                    // Remember Me + Forgot Password
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Row(
                          children: [
                            Checkbox(
                              value: formState.rememberMe,
                              onChanged: (_) => ref.read(authFormStateProvider.notifier).toggleRememberMe(),
                            ),
                            const Text('Remember me'),
                          ],
                        ),
                        TextButton(
                          onPressed: () => context.push('/forgot-password'),
                          child: const Text('Forgot password?'),
                        ),
                      ],
                    ),

                    // Error
                    if (formState.generalError != null)
                      Padding(
                        padding: const EdgeInsets.only(bottom: 16),
                        child: Container(
                      padding: const EdgeInsets.all(AppSpacing.sm),
                      decoration: BoxDecoration(
                        color: theme.colorScheme.errorContainer,
                        borderRadius: BorderRadius.circular(AppRadius.md),
                      ),
                          child: Row(
                            children: [
                              Icon(Icons.error_outline, size: 20, color: theme.colorScheme.error),
                              const SizedBox(width: 8),
                              Expanded(child: Text(formState.generalError!, style: TextStyle(color: theme.colorScheme.error, fontSize: 14))),
                            ],
                          ),
                        ),
                      ),

                    // Login Button
                    FilledButton(
                      onPressed: formState.isLoading ? null : _handleLogin,
                      style: FilledButton.styleFrom(minimumSize: const Size(double.infinity, 52), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppRadius.md))),
                      child: formState.isLoading
                          ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                          : const Text('Sign In', style: TextStyle(fontSize: 16)),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
