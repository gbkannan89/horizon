import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_animate/flutter_animate.dart';
import '../state/dashboard_state.dart';
import '../widgets/dashboard_cards.dart';
import '../models/dashboard_models.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';
import 'package:horizon_mobile/core/ui_kit/animated_stat.dart';

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
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent, // Allow underlying app background
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Horizon', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            Text(_subtitle(state), style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ],
        ),
        backgroundColor: Colors.transparent,
        elevation: 0,
        actions: [
          IconButton(icon: const Icon(Icons.search_rounded), onPressed: () => context.push('/search')),
          IconButton(icon: const Icon(Icons.notifications_outlined), onPressed: () => context.push('/dashboard/notifications')),
          IconButton(
            icon: Icon(Icons.auto_awesome, color: AppColors.amber500), 
            onPressed: () => context.push('/advisor')
          ).animate(onPlay: (controller) => controller.repeat(reverse: true)).shimmer(duration: 2000.ms, color: AppColors.amber400),
        ],
      ),
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: isDark 
                ? [AppColors.darkBackground, AppColors.navy900]
                : [AppColors.lightBackground, AppColors.teal50.withOpacity(0.5)],
          ),
        ),
        child: _buildBody(context, theme, state),
      ),
    );
  }

  String _subtitle(DashboardState s) {
    if (s.status == DashboardLoadStatus.loading) return 'Loading...';
    if (s.dashboard != null) return s.dashboard!.state.replaceAllMapped(RegExp(r'([A-Z])'), (m) => ' ${m.group(1)}').trim();
    return 'Welcome back';
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
      padding: const EdgeInsets.all(AppSpacing.md),
      itemCount: 6,
      itemBuilder: (_, i) => Padding(
        padding: const EdgeInsets.only(bottom: AppSpacing.md),
        child: GlassCard(
          child: SizedBox(
            height: 120, 
            child: Center(
              child: CircularProgressIndicator(strokeWidth: 2, color: AppColors.teal500.withOpacity(0.5))
            )
          )
        ).animate().fade(duration: 500.ms, delay: (i * 100).ms),
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
            Icon(Icons.cloud_off_rounded, size: 64, color: AppColors.red500),
            const SizedBox(height: 16),
            Text('Something went wrong', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            Text(error ?? 'Could not load dashboard', textAlign: TextAlign.center, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 24),
            FilledButton.icon(
              style: FilledButton.styleFrom(backgroundColor: AppColors.teal500),
              onPressed: () => ref.read(dashboardStateProvider.notifier).refresh(),
              icon: const Icon(Icons.refresh_rounded), label: const Text('Try Again'),
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
          Icon(Icons.wifi_off_rounded, size: 64, color: theme.colorScheme.onSurfaceVariant),
          const SizedBox(height: 16),
          Text('You\'re offline', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
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

    final children = <Widget>[
      // Critical Alert
      if (dash.criticalAlert != null)
        Padding(
          padding: const EdgeInsets.only(bottom: AppSpacing.md),
          child: GlassCard(
            blur: 5,
            child: ListTile(
              contentPadding: EdgeInsets.zero,
              leading: Icon(Icons.warning_amber_rounded, color: AppColors.red500),
              title: Text(dash.criticalAlert!.title, style: TextStyle(color: AppColors.red500, fontWeight: FontWeight.bold)),
            ),
          ),
        ),

      // Welcome + Summary row
      _buildWelcomeRow(context, theme, summary),
      const SizedBox(height: AppSpacing.md),

      // Tier 1 — Critical widgets
      ...dash.tier1.where((w) => w.visible).map((w) => Padding(
        padding: const EdgeInsets.only(bottom: AppSpacing.md),
        child: _buildTier1Widget(context, theme, w, summary),
      )),

      // Tier 2 — Important widgets (2-column grid)
      if (dash.tier2.where((w) => w.visible).isNotEmpty) ...[
        Padding(
          padding: const EdgeInsets.only(top: 8, bottom: 12),
          child: Text('Financial Overview', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w800)),
        ),
        _buildTier2Grid(context, theme, dash.tier2.where((w) => w.visible).toList(), summary),
      ],

      // Tier 3 — Contextual
      if (dash.tier3.where((w) => w.visible).isNotEmpty) ...[
        Padding(
          padding: const EdgeInsets.only(top: 8, bottom: 12),
          child: Text('Activity', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w800)),
        ),
        ...dash.tier3.where((w) => w.visible).map((w) => Padding(
          padding: const EdgeInsets.only(bottom: AppSpacing.sm),
          child: _buildTier3Widget(context, theme, w, summary),
        )),
      ],
      
      const SizedBox(height: 100), // Bottom padding for floating nav
    ];

    return RefreshIndicator(
      color: AppColors.teal500,
      onRefresh: () => ref.read(dashboardStateProvider.notifier).refresh(),
      child: ListView.builder(
        padding: const EdgeInsets.all(AppSpacing.md),
        itemCount: children.length,
        itemBuilder: (context, index) {
          return children[index].animate().fade(duration: 400.ms).slideY(begin: 0.1, end: 0, duration: 400.ms, curve: Curves.easeOutQuad);
        },
      ),
    );
  }

  Widget _buildWelcomeRow(BuildContext context, ThemeData theme, SummaryData? summary) {
    return GlassCard(
      child: Row(
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              gradient: AppColors.primaryGradient,
              shape: BoxShape.circle,
            ),
            child: const Icon(Icons.person_rounded, color: Colors.white),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text('Total Wealth', style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600, color: theme.colorScheme.onSurfaceVariant)),
                if (summary != null)
                  AnimatedStatValue(
                    value: summary.netWorth.toDouble(),
                    style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.w900),
                  ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.analytics_rounded, color: AppColors.teal500),
            tooltip: 'View details',
            onPressed: () {},
          ),
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
        return GlassCard(child: ListTile(title: Text(w.title), subtitle: Text(w.numericValue)));
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
        return GlassCard(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [Icon(Icons.receipt_rounded, color: AppColors.amber500), SizedBox(height: 8), Text(w.title), Text(w.numericValue)]));
      default:
        return GlassCard(child: ListTile(title: Text(w.title), subtitle: Text(w.numericValue)));
    }
  }

  Widget _buildTier3Widget(BuildContext context, ThemeData theme, WidgetModel w, SummaryData? s) {
    switch (w.widgetType) {
      case 'recent_events':
        return GlassCard(child: ListTile(contentPadding: EdgeInsets.zero, leading: Icon(Icons.history_rounded, color: AppColors.teal500), title: Text(w.title, style: TextStyle(fontWeight: FontWeight.w600)), trailing: const Icon(Icons.chevron_right_rounded)));
      case 'achievements':
        return GlassCard(child: ListTile(contentPadding: EdgeInsets.zero, leading: Icon(Icons.emoji_events_rounded, color: AppColors.amber500), title: Text(w.title, style: TextStyle(fontWeight: FontWeight.w600))));
      default:
        return GlassCard(child: ListTile(title: Text(w.title)));
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
              padding: const EdgeInsets.only(bottom: AppSpacing.md),
              child: Row(
                children: [
                  Expanded(child: items[idx]),
                  const SizedBox(width: AppSpacing.md),
                  Expanded(child: items[idx + 1]),
                ],
              ),
            );
          }
          if (idx % 2 == 0) return Padding(padding: const EdgeInsets.only(bottom: AppSpacing.md), child: items[idx]);
          return const SizedBox.shrink();
        }
        return Padding(padding: const EdgeInsets.only(bottom: AppSpacing.md), child: entry.value);
      }).toList(),
    );
  }
}
