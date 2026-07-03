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
import 'package:horizon_mobile/features/planning/pages/projection_page.dart';
import 'package:horizon_mobile/features/budget/pages/budget_list_page.dart';
import 'package:horizon_mobile/features/budget/pages/budget_detail_page.dart';
import 'package:horizon_mobile/features/budget/pages/budget_form_page.dart';
import 'package:horizon_mobile/features/planning/pages/retirement_page.dart';
import 'package:horizon_mobile/features/planning/pages/emergency_fund_page.dart';
import 'package:horizon_mobile/features/planning/pages/debt_payoff_page.dart';
import 'package:horizon_mobile/features/planning/pages/investment_page.dart';
import 'package:horizon_mobile/features/accounts/pages/accounts_page.dart';
import 'package:horizon_mobile/features/timeline/pages/timeline_page.dart';
import 'package:horizon_mobile/features/timeline/pages/timeline_detail_page.dart';
import 'package:horizon_mobile/features/goals/pages/goal_detail_page.dart';
import 'package:horizon_mobile/features/accounts/pages/account_detail_page.dart';
import 'package:horizon_mobile/features/settings/pages/settings_page.dart';
import 'package:horizon_mobile/features/settings/pages/profile_page.dart';
import 'package:horizon_mobile/features/settings/pages/edit_profile_page.dart';
import 'package:horizon_mobile/features/settings/pages/preferences_page.dart';
import 'package:horizon_mobile/features/settings/pages/privacy_page.dart';
import 'package:horizon_mobile/features/settings/pages/security_page.dart';
import 'package:horizon_mobile/features/settings/pages/about_page.dart';
import 'package:horizon_mobile/features/settings/pages/ai_settings_page.dart';
import 'package:horizon_mobile/features/settings/pages/data_management_page.dart';
import 'package:horizon_mobile/features/settings/pages/notification_preferences_page.dart';
import 'package:horizon_mobile/features/insights/pages/insights_page.dart';
import 'package:horizon_mobile/features/insights/pages/insight_detail_page.dart';
import 'package:horizon_mobile/features/notifications/pages/notifications_page.dart';
import 'package:horizon_mobile/features/advisor/pages/advisor_page.dart';
import 'package:horizon_mobile/features/transactions/pages/transactions_page.dart';
import 'package:horizon_mobile/features/transactions/pages/transaction_detail_page.dart';
import 'package:horizon_mobile/features/transactions/pages/transaction_form_page.dart';
import 'package:horizon_mobile/features/search/pages/search_page.dart';
import 'package:horizon_mobile/features/recurring/pages/recurring_list_page.dart';
import 'package:horizon_mobile/features/recurring/pages/recurring_detail_page.dart';
import 'package:horizon_mobile/features/recurring/pages/recurring_form_page.dart';
import 'package:horizon_mobile/features/household/pages/household_list_page.dart';
import 'package:horizon_mobile/features/household/pages/household_detail_page.dart';
import 'package:horizon_mobile/features/household/pages/household_create_page.dart';
import 'package:horizon_mobile/features/household/pages/household_invite_page.dart';
import 'package:horizon_mobile/features/household/pages/household_settings_page.dart';
import 'package:horizon_mobile/features/automation/pages/automation_list_page.dart';
import 'package:horizon_mobile/features/automation/pages/automation_detail_page.dart';
import 'package:horizon_mobile/features/automation/pages/automation_form_page.dart';
import 'package:horizon_mobile/app/shell.dart';

Page<void> _slideTransition(Widget child) {
  return CustomTransitionPage<void>(
    child: child,
    transitionsBuilder: (context, animation, secondaryAnimation, child) {
      const begin = Offset(1.0, 0.0);
      const end = Offset.zero;
      const curve = Curves.easeInOutCubic;
      var tween = Tween(begin: begin, end: end).chain(CurveTween(curve: curve));
      return SlideTransition(position: animation.drive(tween), child: child);
    },
  );
}

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
      GoRoute(path: '/planning', builder: (_, __) => const PlanningPage(), routes: [
        GoRoute(path: 'projections', builder: (_, __) => const ProjectionPage()),
        GoRoute(path: 'budget', builder: (_, __) => const BudgetListPage(), routes: [
          GoRoute(path: 'add', builder: (_, __) => const BudgetFormPage()),
          GoRoute(path: ':id', builder: (_, state) => BudgetDetailPage(budgetId: state.pathParameters['id'] ?? '')),
        ]),
        GoRoute(path: 'retirement', builder: (_, __) => const RetirementPage()),
        GoRoute(path: 'emergency-fund', builder: (_, __) => const EmergencyFundPage()),
        GoRoute(path: 'debt-payoff', builder: (_, __) => const DebtPayoffPage()),
        GoRoute(path: 'investment', builder: (_, __) => const InvestmentPage()),
      ]),
      GoRoute(path: '/accounts', builder: (_, __) => const AccountsPage(), routes: [
        GoRoute(path: ':id', builder: (_, state) => AccountDetailPage(accountId: state.pathParameters['id'] ?? '')),
      ]),
      GoRoute(path: '/transactions', builder: (_, __) => const TransactionsPage(), routes: [
        GoRoute(path: ':id', builder: (_, state) => TransactionDetailPage(transactionId: state.pathParameters['id'] ?? '')),
      ]),
      GoRoute(path: '/transactions/add', builder: (_, __) => const TransactionFormPage()),
      GoRoute(path: '/recurring', builder: (_, __) => const RecurringListPage(), routes: [
        GoRoute(path: 'add', builder: (_, __) => const RecurringFormPage()),
        GoRoute(path: ':id', builder: (_, state) => RecurringDetailPage(id: state.pathParameters['id'] ?? '')),
      ]),
      GoRoute(path: '/automation', builder: (_, __) => const AutomationListPage(), routes: [
        GoRoute(path: 'add', builder: (_, __) => const AutomationFormPage()),
        GoRoute(path: ':id', builder: (_, state) => AutomationDetailPage(ruleId: state.pathParameters['id'] ?? '')),
      ]),
      GoRoute(path: '/search', pageBuilder: (_, __) => _slideTransition(const SearchPage())),
      GoRoute(path: '/insights', builder: (_, __) => const InsightsPage(), routes: [
        GoRoute(path: ':id', pageBuilder: (_, state) => _slideTransition(InsightDetailPage(insightId: state.pathParameters['id'] ?? ''))),
      ]),
      GoRoute(path: '/settings', pageBuilder: (_, __) => _slideTransition(const SettingsPage())),
      GoRoute(path: '/settings/profile', pageBuilder: (_, __) => _slideTransition(const ProfilePage())),
      GoRoute(path: '/settings/profile/edit', pageBuilder: (_, __) => _slideTransition(const EditProfilePage())),
      GoRoute(path: '/settings/preferences', pageBuilder: (_, __) => _slideTransition(const PreferencesPage())),
      GoRoute(path: '/settings/privacy', pageBuilder: (_, __) => _slideTransition(const PrivacyPage())),
      GoRoute(path: '/settings/security', pageBuilder: (_, __) => _slideTransition(const SecurityPage())),
      GoRoute(path: '/settings/about', pageBuilder: (_, __) => _slideTransition(const AboutPage())),
      GoRoute(path: '/settings/ai', pageBuilder: (_, __) => _slideTransition(const AiSettingsPage())),
      GoRoute(path: '/settings/data', pageBuilder: (_, __) => _slideTransition(const DataManagementPage())),
      GoRoute(path: '/settings/notifications', pageBuilder: (_, __) => _slideTransition(const NotificationPreferencesPage())),
      GoRoute(path: '/households', builder: (_, __) => const HouseholdListPage(), routes: [
        GoRoute(path: 'create', builder: (_, __) => const HouseholdCreatePage()),
        GoRoute(path: ':id', builder: (_, state) => HouseholdDetailPage(householdId: state.pathParameters['id'] ?? ''), routes: [
          GoRoute(path: 'settings', builder: (_, state) => HouseholdSettingsPage(householdId: state.pathParameters['id'] ?? '')),
          GoRoute(path: 'invite', builder: (_, state) => HouseholdInvitePage(householdId: state.pathParameters['id'] ?? '')),
        ]),
      ]),
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
        ListTile(leading: const Icon(Icons.receipt_long_outlined), title: const Text('Transactions'), onTap: () => context.push('/transactions')),
        ListTile(leading: const Icon(Icons.repeat), title: const Text('Recurring'), onTap: () => context.push('/recurring')),
        ListTile(leading: const Icon(Icons.auto_awesome), title: const Text('Automation'), onTap: () => context.push('/automation')),
        ListTile(leading: const Icon(Icons.account_balance), title: const Text('Accounts'), onTap: () => context.push('/accounts')),
        ListTile(leading: const Icon(Icons.lightbulb_outline), title: const Text('Insights'), onTap: () => context.push('/insights')),
        ListTile(leading: const Icon(Icons.auto_awesome), title: const Text('Advisor'), onTap: () => context.push('/advisor')),
        ListTile(leading: const Icon(Icons.home_work_outlined), title: const Text('Households'), onTap: () => context.push('/households')),
        ListTile(leading: const Icon(Icons.settings_outlined), title: const Text('Settings'), onTap: () => context.push('/settings')),
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
