import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/providers/auth_state.dart';
import 'package:horizon_mobile/features/auth/repositories/auth_repository.dart';
import 'package:horizon_mobile/features/auth/pages/splash_page.dart';
import 'package:horizon_mobile/features/auth/pages/login_page.dart';
import 'package:horizon_mobile/features/auth/pages/forgot_password_page.dart';
import 'package:horizon_mobile/features/auth/pages/reset_password_page.dart';
import 'package:horizon_mobile/features/auth/pages/session_expired_page.dart';
import 'package:horizon_mobile/features/dashboard/pages/dashboard_page.dart';
import 'package:horizon_mobile/features/goals/pages/goals_page.dart';
import 'package:horizon_mobile/features/portfolio/pages/portfolio_page.dart';
import 'package:horizon_mobile/features/planning/pages/planning_page.dart';
import 'package:horizon_mobile/features/accounts/pages/accounts_page.dart';
import 'package:horizon_mobile/features/timeline/pages/timeline_page.dart';
import 'package:horizon_mobile/features/timeline/pages/timeline_detail_page.dart';
import 'package:horizon_mobile/features/goals/pages/goal_detail_page.dart';
import 'package:horizon_mobile/features/insights/pages/insights_page.dart';
import 'package:horizon_mobile/features/notifications/pages/notifications_page.dart';
import 'package:horizon_mobile/features/advisor/pages/advisor_page.dart';
import 'package:horizon_mobile/app/shell.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authStateProvider);
  return GoRouter(
    initialLocation: '/splash',
    debugLogDiagnostics: false,
    redirect: (context, state) {
      final location = state.matchedLocation;
      final isLoggedIn = authState.status == AuthStatus.authenticated;
      final isAuthRoute = location.startsWith('/login') || location == '/splash'
          || location.startsWith('/forgot-password') || location.startsWith('/reset-password');

      if (isLoggedIn && isAuthRoute) return '/dashboard';
      if (!isLoggedIn && !isAuthRoute && location != '/session-expired') return '/login';
      return null;
    },
    routes: [
      GoRoute(path: '/splash', builder: (_, __) => const SplashPage()),
      GoRoute(path: '/login', builder: (_, __) => const LoginPage()),
      GoRoute(path: '/forgot-password', builder: (_, __) => const ForgotPasswordPage()),
      GoRoute(path: '/reset-password', builder: (_, state) => ResetPasswordPage(token: state.uri.queryParameters['token'] ?? '')),
      GoRoute(path: '/session-expired', builder: (_, __) => const SessionExpiredPage()),
      GoRoute(path: '/advisor', builder: (_, __) => const AdvisorPage()),
      StatefulShellRoute.indexedStack(
        builder: (_, __, navigationShell) => AppShell(navigationShell: navigationShell),
        branches: [
          StatefulShellBranch(routes: [GoRoute(path: '/dashboard', builder: (_, __) => const DashboardPage(), routes: [
            GoRoute(path: 'notifications', builder: (_, __) => const NotificationsPage()),
          ])]),
          StatefulShellBranch(routes: [GoRoute(path: '/goals', builder: (_, __) => const GoalsPage(), routes: [
            GoRoute(path: ':id', builder: (_, state) => GoalDetailPage(goalId: state.pathParameters['id'] ?? '')),
          ])]),
          StatefulShellBranch(routes: [GoRoute(path: '/portfolio', builder: (_, __) => const PortfolioPage())]),
          StatefulShellBranch(routes: [GoRoute(path: '/timeline', builder: (_, __) => const TimelinePage(), routes: [
            GoRoute(path: ':id', builder: (_, state) => TimelineDetailPage(timelineId: state.pathParameters['id'] ?? '')),
          ])]),
          StatefulShellBranch(routes: [GoRoute(path: '/more', builder: (_, __) => const _MorePage())]),
        ],
      ),
      GoRoute(path: '/planning', builder: (_, __) => const PlanningPage()),
      GoRoute(path: '/accounts', builder: (_, __) => const AccountsPage()),
      GoRoute(path: '/insights', builder: (_, __) => const InsightsPage()),
    ],
  );
});

class _MorePage extends ConsumerWidget {
  const _MorePage();
  @override
  Widget build(BuildContext context, WidgetRef ref) => Scaffold(
    appBar: AppBar(title: const Text('More')),
    body: ListView(
      children: [
        ListTile(leading: const Icon(Icons.schema_outlined), title: const Text('Planning'), onTap: () => context.push('/planning')),
        ListTile(leading: const Icon(Icons.account_balance), title: const Text('Accounts'), onTap: () => context.push('/accounts')),
        ListTile(leading: const Icon(Icons.lightbulb_outline), title: const Text('Insights'), onTap: () => context.push('/insights')),
        ListTile(leading: const Icon(Icons.auto_awesome), title: const Text('Advisor'), onTap: () => context.push('/advisor')),
        const Divider(),
        ListTile(leading: const Icon(Icons.logout, color: Colors.red), title: const Text('Sign Out', style: TextStyle(color: Colors.red)), onTap: () async {
          final repo = ref.read(authRepositoryProvider);
          await repo.logout();
          ref.read(authStateProvider.notifier).unauthenticated();
          context.go('/login');
        }),
      ],
    ),
  );
}
