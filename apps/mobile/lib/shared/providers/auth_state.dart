import 'package:flutter_riverpod/flutter_riverpod.dart';

final authStateProvider = StateNotifierProvider<AuthStateNotifier, AuthState>((ref) {
  return AuthStateNotifier();
});

enum AuthStatus { initial, authenticated, unauthenticated, loading }

class AuthState {
  final AuthStatus status;
  final String? userId;
  final String? email;
  final String? error;

  const AuthState({
    this.status = AuthStatus.initial,
    this.userId,
    this.email,
    this.error,
  });

  AuthState copyWith({AuthStatus? status, String? userId, String? email, String? error}) =>
      AuthState(status: status ?? this.status, userId: userId ?? this.userId, email: email ?? this.email, error: error);
}

class AuthStateNotifier extends StateNotifier<AuthState> {
  AuthStateNotifier() : super(const AuthState());

  void authenticated({String? userId, String? email}) {
    state = AuthState(status: AuthStatus.authenticated, userId: userId, email: email);
  }

  void unauthenticated({String? error}) {
    state = AuthState(status: AuthStatus.unauthenticated, error: error);
  }

  void loading() { state = const AuthState(status: AuthStatus.loading); }
  void initial() { state = const AuthState(); }
  void clearError() { state = state.copyWith(error: null); }
}
