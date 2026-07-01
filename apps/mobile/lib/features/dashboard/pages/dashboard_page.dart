import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../state/dashboard_state.dart';
import '../widgets/dashboard_cards.dart';
import '../models/dashboard_models.dart';

class DashboardPage extends ConsumerStatefulWidget {
  const DashboardPage({super.key});
  @override
  ConsumerState<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends ConsumerState<DashboardPage> {
  @override
  void initState() {
    super.initState();
    Future.microtask(() => ref.read(dashboardStateProvider.notifier).load());
  }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(dashboardStateProvider);
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Horizon', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            Text(_subtitle(state), style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ],
        ),
        actions: [
          IconButton(icon: const Icon(Icons.search), onPressed: () => context.push('/search')),
          IconButton(icon: const Icon(Icons.notifications_outlined), onPressed: () => context.push('/dashboard/notifications')),
          IconButton(icon: const Icon(Icons.auto_awesome), onPressed: () => context.push('/advisor')),
        ],
      ),
      body: _buildBody(context, theme, state),
    );
  }

  String _subtitle(DashboardState s) {
    if (s.status == DashboardLoadStatus.loading) return 'Loading...';
    if (s.dashboard != null) return s.dashboard!.state.replaceAllMapped(RegExp(r'([A-Z])'), (m) => ' ${m.group(1)}').trim();
    return 'Welcome';
  }

  Widget _buildBody(BuildContext context, ThemeData theme, DashboardState state) {
    switch (state.status) {
      case DashboardLoadStatus.initial:
      case DashboardLoadStatus.loading:
        return const SkeletonList();
      case DashboardLoadStatus.error:
        return ErrorView(
          message: 'Something went wrong',
          detail: state.error ?? 'Could not load dashboard',
          onRetry: () => ref.read(dashboardStateProvider.notifier).refresh(),
        );
      case DashboardLoadStatus.offline:
        return const OfflineView();
      case DashboardLoadStatus.loaded:
        return _buildDashboard(context, theme, state);
    }
  }

  Widget _buildDashboard(BuildContext context, ThemeData theme, DashboardState state) {
    final dash = state.dashboard;
    final summary = state.summary;
    if (dash == null) return ErrorView(message: 'No dashboard data', onRetry: () => ref.read(dashboardStateProvider.notifier).refresh());

    return RefreshIndicator(
      onRefresh: () => ref.read(dashboardStateProvider.notifier).refresh(),
      child: ListView(
        padding: const EdgeInsets.all(AppTheme.spacingLg),
        children: [
          if (dash.criticalAlert != null)
            Padding(
              padding: const EdgeInsets.only(bottom: 12),
              child: Card(
                color: theme.colorScheme.errorContainer,
                child: ListTile(
                  leading: Icon(Icons.warning_amber, color: theme.colorScheme.error),
                  title: Text(dash.criticalAlert!.title, style: TextStyle(color: theme.colorScheme.error, fontWeight: FontWeight.w600)),
                ),
              ),
            ),
          _buildWelcomeRow(context, theme, summary),
          const SizedBox(height: 12),
          ...dash.tier1.where((w) => w.visible).map((w) => Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: _buildTier1Widget(context, theme, w, summary),
          )),
          if (dash.tier2.where((w) => w.visible).isNotEmpty) ...[
            SectionHeader(title: 'Financial Overview'),
            _buildTier2Grid(context, theme, dash.tier2.where((w) => w.visible).toList(), summary),
          ],
          if (dash.tier3.where((w) => w.visible).isNotEmpty) ...[
            SectionHeader(title: 'Activity'),
            ...dash.tier3.where((w) => w.visible).map((w) => Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: _buildTier3Widget(context, theme, w, summary),
            )),
          ],
        ],
      ),
    );
  }

  Widget _buildWelcomeRow(BuildContext context, ThemeData theme, SummaryData? summary) {
    return AppCard(
      child: Row(
        children: [
          CircleAvatar(
            backgroundColor: theme.colorScheme.primaryContainer,
            child: Icon(Icons.person, color: theme.colorScheme.primary),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Financial Overview', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                if (summary != null)
                  Text('Net Worth: ${formatMoney(summary.netWorth)}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ],
            ),
          ),
          IconButton(icon: const Icon(Icons.open_in_new), tooltip: 'View details', onPressed: () {}),
        ],
      ),
    );
  }

  Widget _buildTier1Widget(BuildContext context, ThemeData theme, WidgetModel w, SummaryData? s) {
    switch (w.widgetType) {
      case 'health_score':
        return HealthScoreCard(widget: w, score: s?.healthScore ?? 0, grade: s?.healthGrade ?? '');
      case 'goal_progress':
        return GoalProgressCard(widget: w, onTrack: s?.goalsOnTrack ?? 0, total: s?.totalGoals ?? 0);
      case 'top_recommendation':
        return RecommendationCard(widget: w);
      default:
        return AppCard(child: ListTile(title: Text(w.title), subtitle: Text(w.numericValue)));
    }
  }

  Widget _buildTier2Widget(BuildContext context, ThemeData theme, WidgetModel w, SummaryData? s) {
    switch (w.widgetType) {
      case 'net_worth':
        return NetWorthCard(widget: w, netWorth: s?.netWorth ?? 0);
      case 'cash_position':
        return CashFlowCard(widget: w, income: s?.netWorth ?? 0, expenses: 0);
      case 'portfolio_snapshot':
        return PortfolioCard(widget: w, value: s?.portfolioValue ?? 0);
      case 'risk_summary':
        return RiskScoreCard(widget: w, score: s?.riskScore ?? 0, level: s?.riskLevel ?? '');
      case 'debt_summary':
        return DebtSummaryCard(widget: w, debt: s?.totalDebt ?? 0);
      case 'upcoming_bills':
        return AppCard(child: ListTile(leading: Icon(Icons.receipt, color: Colors.amber), title: Text(w.title), subtitle: Text(w.numericValue)));
      default:
        return AppCard(child: ListTile(title: Text(w.title), subtitle: Text(w.numericValue)));
    }
  }

  Widget _buildTier3Widget(BuildContext context, ThemeData theme, WidgetModel w, SummaryData? s) {
    switch (w.widgetType) {
      case 'recent_events':
        return AppCard(child: ListTile(leading: Icon(Icons.history, color: theme.colorScheme.primary), title: Text(w.title), trailing: const Icon(Icons.chevron_right)));
      case 'achievements':
        return AppCard(child: ListTile(leading: Icon(Icons.emoji_events, color: Colors.amber), title: Text(w.title)));
      default:
        return AppCard(child: ListTile(title: Text(w.title)));
    }
  }

  Widget _buildTier2Grid(BuildContext context, ThemeData theme, List<WidgetModel> widgets, SummaryData? s) {
    final items = widgets.map((w) => _buildTier2Widget(context, theme, w, s)).toList();
    return Column(
      children: items.asMap().entries.map((entry) {
        if (entry.value is NetWorthCard || entry.value is RiskScoreCard || entry.value is PortfolioCard) {
          final idx = entry.key;
          if (idx % 2 == 0 && idx + 1 < items.length) {
            return Padding(
              padding: const EdgeInsets.only(bottom: 12),
              child: Row(children: [Expanded(child: items[idx]), const SizedBox(width: 12), Expanded(child: items[idx + 1])]),
            );
          }
          if (idx % 2 == 0) return Padding(padding: const EdgeInsets.only(bottom: 12), child: items[idx]);
          return const SizedBox.shrink();
        }
        return Padding(padding: const EdgeInsets.only(bottom: 12), child: entry.value);
      }).toList(),
    );
  }
}
