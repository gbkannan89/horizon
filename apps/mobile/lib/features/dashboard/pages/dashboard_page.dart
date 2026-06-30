import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
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

    // Use Future.microtask to navigate after build — not during build
    // Navigation happens in async callbacks from card taps

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
        return _buildSkeleton(theme);
      case DashboardLoadStatus.error:
        return _buildError(context, theme, state.error);
      case DashboardLoadStatus.offline:
        return _buildOffline(context, theme);
      case DashboardLoadStatus.loaded:
        return _buildDashboard(context, theme, state);
    }
  }

  Widget _buildSkeleton(ThemeData theme) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: 6,
      itemBuilder: (_, __) => Padding(
        padding: const EdgeInsets.only(bottom: 12),
        child: Card(child: SizedBox(height: 100, child: Center(child: CircularProgressIndicator(strokeWidth: 2, color: theme.colorScheme.primary.withOpacity(0.3))))),
      ),
    );
  }

  Widget _buildError(BuildContext context, ThemeData theme, String? error) {
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.cloud_off, size: 64, color: theme.colorScheme.error),
            const SizedBox(height: 16),
            Text('Something went wrong', style: theme.textTheme.titleLarge),
            const SizedBox(height: 8),
            Text(error ?? 'Could not load dashboard', textAlign: TextAlign.center, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 24),
            FilledButton.icon(
              onPressed: () => ref.read(dashboardStateProvider.notifier).refresh(),
              icon: const Icon(Icons.refresh), label: const Text('Try Again'),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildOffline(BuildContext context, ThemeData theme) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.wifi_off, size: 64, color: theme.colorScheme.onSurfaceVariant),
          const SizedBox(height: 16),
          Text('You\'re offline', style: theme.textTheme.titleLarge),
          const SizedBox(height: 8),
          Text('Connect to the internet to see your dashboard', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
        ],
      ),
    );
  }

  Widget _buildDashboard(BuildContext context, ThemeData theme, DashboardState state) {
    final dash = state.dashboard;
    final summary = state.summary;
    if (dash == null) return _buildError(context, theme, 'No dashboard data');

    return RefreshIndicator(
      onRefresh: () => ref.read(dashboardStateProvider.notifier).refresh(),
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Critical Alert
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

          // Welcome + Summary row
          _buildWelcomeRow(context, theme, summary),

          const SizedBox(height: 12),

          // Tier 1 — Critical widgets
          ...dash.tier1.where((w) => w.visible).map((w) => Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: _buildTier1Widget(context, theme, w, summary),
          )),

          // Tier 2 — Important widgets (2-column grid)
          if (dash.tier2.where((w) => w.visible).isNotEmpty) ...[
            Padding(
              padding: const EdgeInsets.only(top: 8, bottom: 8),
              child: Text('Financial Overview', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            ),
            _buildTier2Grid(context, theme, dash.tier2.where((w) => w.visible).toList(), summary),
          ],

          // Tier 3 — Contextual
          if (dash.tier3.where((w) => w.visible).isNotEmpty) ...[
            Padding(
              padding: const EdgeInsets.only(top: 8, bottom: 8),
              child: Text('Activity', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            ),
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
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
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
                    Text('Net Worth: ${WidgetModel.formatMoney(summary.netWorth)}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ],
              ),
            ),
            IconButton(
              icon: const Icon(Icons.open_in_new),
              tooltip: 'View details',
              onPressed: () {},
            ),
          ],
        ),
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
        return Card(child: ListTile(title: Text(w.title), subtitle: Text(w.numericValue)));
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
        return Card(child: ListTile(leading: Icon(Icons.receipt, color: Colors.amber), title: Text(w.title), subtitle: Text(w.numericValue)));
      default:
        return Card(child: ListTile(title: Text(w.title), subtitle: Text(w.numericValue)));
    }
  }

  Widget _buildTier3Widget(BuildContext context, ThemeData theme, WidgetModel w, SummaryData? s) {
    switch (w.widgetType) {
      case 'recent_events':
        return Card(child: ListTile(leading: Icon(Icons.history, color: theme.colorScheme.primary), title: Text(w.title), trailing: Icon(Icons.chevron_right)));
      case 'achievements':
        return Card(child: ListTile(leading: Icon(Icons.emoji_events, color: Colors.amber), title: Text(w.title)));
      default:
        return Card(child: ListTile(title: Text(w.title)));
    }
  }

  Widget _buildTier2Grid(BuildContext context, ThemeData theme, List<WidgetModel> widgets, SummaryData? s) {
    final items = widgets.map((w) => _buildTier2Widget(context, theme, w, s)).toList();
    // Alternate between full-width and half-width
    return Column(
      children: items.asMap().entries.map((entry) {
        // Make financial cards half-width in pairs
        if (entry.value is NetWorthCard || entry.value is RiskScoreCard || entry.value is PortfolioCard) {
          final idx = entry.key;
          if (idx % 2 == 0 && idx + 1 < items.length) {
            return Padding(
              padding: const EdgeInsets.only(bottom: 12),
              child: Row(
                children: [
                  Expanded(child: items[idx]),
                  const SizedBox(width: 12),
                  Expanded(child: items[idx + 1]),
                ],
              ),
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
