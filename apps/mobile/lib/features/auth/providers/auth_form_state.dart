import 'package:flutter_riverpod/flutter_riverpod.dart';

final authFormStateProvider = StateNotifierProvider<AuthFormNotifier, AuthFormState>((ref) {
  return AuthFormNotifier();
});

class AuthFormState {
  final String email;
  final String password;
  final String confirmPassword;
  final bool obscurePassword;
  final bool rememberMe;
  final bool isLoading;
  final String? emailError;
  final String? passwordError;
  final String? confirmPasswordError;
  final String? generalError;
  final AuthFormStep step;

  const AuthFormState({
    this.email = '',
    this.password = '',
    this.confirmPassword = '',
    this.obscurePassword = true,
    this.rememberMe = false,
    this.isLoading = false,
    this.emailError,
    this.passwordError,
    this.confirmPasswordError,
    this.generalError,
    this.step = AuthFormStep.login,
  });

  AuthFormState copyWith({
    String? email, String? password, String? confirmPassword,
    bool? obscurePassword, bool? rememberMe, bool? isLoading,
    String? emailError, String? passwordError, String? confirmPasswordError,
    String? generalError, AuthFormStep? step,
  }) => AuthFormState(
    email: email ?? this.email,
    password: password ?? this.password,
    confirmPassword: confirmPassword ?? this.confirmPassword,
    obscurePassword: obscurePassword ?? this.obscurePassword,
    rememberMe: rememberMe ?? this.rememberMe,
    isLoading: isLoading ?? this.isLoading,
    emailError: emailError ?? this.emailError,
    passwordError: passwordError ?? this.passwordError,
    confirmPasswordError: confirmPasswordError ?? this.confirmPasswordError,
    generalError: generalError ?? this.generalError,
    step: step ?? this.step,
  );
}

enum AuthFormStep { login, forgotPassword, resetPassword }

class AuthFormNotifier extends StateNotifier<AuthFormState> {
  AuthFormNotifier() : super(const AuthFormState());

  void setEmail(String v) {
    state = state.copyWith(email: v, emailError: _validateEmail(v));
  }

  void setPassword(String v) {
    state = state.copyWith(password: v, passwordError: _validatePassword(v));
  }

  void setConfirmPassword(String v) {
    state = state.copyWith(confirmPassword: v, confirmPasswordError: _validateConfirmPassword(v));
  }

  void toggleObscurePassword() {
    state = state.copyWith(obscurePassword: !state.obscurePassword);
  }

  void toggleRememberMe() {
    state = state.copyWith(rememberMe: !state.rememberMe);
  }

  void setStep(AuthFormStep step) {
    state = state.copyWith(step: step, generalError: null);
  }

  void setLoading(bool v) { state = state.copyWith(isLoading: v); }
  void setError(String e) { state = state.copyWith(generalError: e, isLoading: false); }
  void clearError() { state = state.copyWith(generalError: null); }

  bool validateLogin() {
    final emailErr = _validateEmail(state.email);
    final passErr = _validatePassword(state.password);
    state = state.copyWith(emailError: emailErr, passwordError: passErr);
    return emailErr == null && passErr == null;
  }

  bool validateForgotPassword() {
    final emailErr = _validateEmail(state.email);
    state = state.copyWith(emailError: emailErr);
    return emailErr == null;
  }

  bool validateResetPassword() {
    final passErr = _validatePassword(state.password);
    final confirmErr = _validateConfirmPassword(state.confirmPassword);
    state = state.copyWith(passwordError: passErr, confirmPasswordError: confirmErr);
    return passErr == null && confirmErr == null;
  }

  String? _validateEmail(String v) {
    if (v.isEmpty) return 'Email is required';
    final emailRegex = RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,}$');
    if (!emailRegex.hasMatch(v)) return 'Enter a valid email address';
    return null;
  }

  String? _validatePassword(String v) {
    if (v.isEmpty) return 'Password is required';
    if (v.length < 8) return 'Password must be at least 8 characters';
    return null;
  }

  String? _validateConfirmPassword(String v) {
    if (v.isEmpty) return 'Please confirm your password';
    if (v != state.password) return 'Passwords do not match';
    return null;
  }
}
