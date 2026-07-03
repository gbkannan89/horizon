import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_animate/flutter_animate.dart';
import '../state/dashboard_state.dart';
import '../widgets/dashboard_cards.dart';
import '../models/dashboard_models.dart';
import '../widgets/dashboard_signature_widgets.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';
import 'package:horizon_mobile/core/ui_kit/animated_stat.dart';
import 'package:horizon_mobile/shared/providers/auth_state.dart';

class DashboardPage extends ConsumerStatefulWidget {
  const DashboardPage({super.key});
  @override
  ConsumerState<DashboardPage> createState() => _DashboardPageState();
}

class _DashboardPageState extends ConsumerState<DashboardPage> {
  bool _dismissedWindfall = false;
  bool _hasInsurance = false;

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
      backgroundColor: isDark ? const Color(0xFF0F172A) : const Color(0xFFF8FAFC),
      appBar: AppBar(
        toolbarHeight: 0, // Hide AppBar completely since we have top padding in Hero
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: _buildBody(context, theme, state),
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
            Text('Something went wrong', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
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
          Text('You\'re offline', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
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

    final authState = ref.watch(authStateProvider);
    final email = authState.email ?? '';
    final name = email.isNotEmpty ? email.split('@')[0] : 'Kannan';
    final userName = name.isNotEmpty ? name[0].toUpperCase() + name.substring(1) : 'Guest';

    return RefreshIndicator(
      color: AppColors.teal500,
      onRefresh: () => ref.read(dashboardStateProvider.notifier).refresh(),
      child: CustomScrollView(
        slivers: [
          // ── Hero gradient header ──────────────────────────────────────────
          SliverToBoxAdapter(
            child: _buildHeroHeader(theme, summary, userName),
          ),

          // ── Body cards ───────────────────────────────────────────────────
          SliverPadding(
            padding: const EdgeInsets.fromLTRB(20, 20, 20, 100),
            sliver: SliverList(
              delegate: SliverChildListDelegate([
                // Critical Alert
                if (dash.criticalAlert != null) ...[
                  GlassCard(
                    child: ListTile(
                      contentPadding: EdgeInsets.zero,
                      leading: Icon(Icons.warning_amber_rounded, color: AppColors.red500),
                      title: Text(dash.criticalAlert!.title, style: TextStyle(color: AppColors.red500, fontWeight: FontWeight.bold)),
                    ),
                  ),
                  const SizedBox(height: 20),
                ],

                // 3 Mini Cards Row (Income / Spent / Left)
                _buildIncomeSpentRow(summary),
                const SizedBox(height: 20),

                // Gatekeeper Card (Defensive Nudge)
                if (!_hasInsurance) ...[
                  GatekeeperCard(
                    onFixCover: () {
                      setState(() => _hasInsurance = true);
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Insurance added! Wealth steering unlocked.'), behavior: SnackBarBehavior.floating),
                      );
                    },
                  ),
                  const SizedBox(height: 20),
                ],

                // Windfall Interceptor
                if (!_dismissedWindfall) ...[
                  WindfallInterceptorCard(
                    amount: 50000,
                    onSelectJob: (job) {
                      setState(() => _dismissedWindfall = true);
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(content: Text('Windfall allocated to $job!'), behavior: SnackBarBehavior.floating),
                      );
                    },
                  ),
                  const SizedBox(height: 20),
                ],

                // Calm Coach Card
                const CalmCoachCard(
                  overspentAmount: 8000,
                  recoveryAmount: 2700,
                  months: 3,
                ),
                const SizedBox(height: 20),

                // Tier 1 — Critical widgets (Health Score / Goal Progress / Top Recommendation)
                ...dash.tier1.where((w) => w.visible).map((w) => Column(
                  children: [
                    _buildTier1Widget(context, theme, w, summary),
                    const SizedBox(height: 20),
                  ],
                )),

                // Tier 2 — Important widgets
                if (dash.tier2.where((w) => w.visible).isNotEmpty) ...[
                  Padding(
                    padding: const EdgeInsets.only(top: 8, bottom: 12),
                    child: Text('Financial Overview', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w800)),
                  ),
                  _buildTier2Grid(context, theme, dash.tier2.where((w) => w.visible).toList(), summary),
                  const SizedBox(height: 20),
                ],

                // Tier 3 — Contextual
                if (dash.tier3.where((w) => w.visible).isNotEmpty) ...[
                  Padding(
                    padding: const EdgeInsets.only(top: 8, bottom: 12),
                    child: Text('Activity', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w800)),
                  ),
                  ...dash.tier3.where((w) => w.visible).map((w) => Column(
                    children: [
                      _buildTier3Widget(context, theme, w, summary),
                      const SizedBox(height: 12),
                    ],
                  )),
                ],
              ]),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHeroHeader(ThemeData theme, SummaryData? summary, String userName) {
    final loaded = summary != null;
    final score = summary?.healthScore ?? 0;
    
    return Container(
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [Color(0xFF0F2057), Color(0xFF1E3A8A), Color(0xFF2563EB)],
        ),
        borderRadius: BorderRadius.vertical(bottom: Radius.circular(36)),
      ),
      padding: const EdgeInsets.fromLTRB(24, 48, 24, 32),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Top row
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(_greeting(), style: const TextStyle(color: Color(0xFF93C5FD), fontSize: 14, fontWeight: FontWeight.w500)),
                  const SizedBox(height: 2),
                  Text('$userName 👋', style: const TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.w800)),
                ],
              ),
              IconButton(
                icon: const Icon(Icons.notifications_none_rounded, color: Colors.white, size: 24),
                onPressed: () => context.push('/dashboard/notifications'),
              ),
            ],
          ),

          const SizedBox(height: 28),

          // Score + Net Worth row
          Row(
            children: [
              // Financial Score ring
              Stack(
                alignment: Alignment.center,
                children: [
                  SizedBox(
                    width: 90,
                    height: 90,
                    child: CircularProgressIndicator(
                      value: loaded ? (score / 100.0).clamp(0.0, 1.0) : 0.0,
                      strokeWidth: 7,
                      valueColor: const AlwaysStoppedAnimation<Color>(Color(0xFF34D399)),
                      backgroundColor: Colors.white.withOpacity(0.15),
                    ),
                  ),
                  Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        loaded ? '$score' : '--',
                        style: const TextStyle(fontSize: 26, fontWeight: FontWeight.w900, color: Colors.white),
                      ),
                      const Text('/100', style: TextStyle(fontSize: 10, color: Color(0xFF93C5FD))),
                    ],
                  ),
                ],
              ),
              const SizedBox(width: 20),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('Financial Score', style: TextStyle(color: Color(0xFF93C5FD), fontSize: 12, fontWeight: FontWeight.w500)),
                    const SizedBox(height: 4),
                    Text(
                      loaded && score >= 80 ? '🌟 Excellent' : loaded && score >= 60 ? '👍 Good' : '⚡ Improving',
                      style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 15),
                    ),
                    const SizedBox(height: 16),
                    const Text('Net Worth', style: TextStyle(color: Color(0xFF93C5FD), fontSize: 12, fontWeight: FontWeight.w500)),
                    const SizedBox(height: 2),
                    Text(
                      loaded ? '₹${_formatNum(summary!.netWorth.toDouble())}' : '₹--',
                      style: const TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.w900, letterSpacing: -0.5),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildIncomeSpentRow(SummaryData? summary) {
    final loaded = summary != null;
    final income = summary?.cashBalance.toDouble() ?? 0.0;
    final spent = summary?.totalDebt.toDouble() ?? 0.0;
    final left = summary?.portfolioValue.toDouble() ?? 0.0;

    return Row(
      children: [
        _miniCard('Income', loaded ? '₹${_formatNum(income)}' : '--', const Color(0xFF059669), Icons.trending_up_rounded),
        const SizedBox(width: 10),
        _miniCard('Spent', loaded ? '₹${_formatNum(spent)}' : '--', const Color(0xFFE88A1A), Icons.shopping_cart_outlined),
        const SizedBox(width: 10),
        _miniCard('Left', loaded ? '₹${_formatNum(left)}' : '--', const Color(0xFF1E3A8A), Icons.savings_outlined),
      ],
    );
  }

  Widget _miniCard(String label, String value, Color color, IconData icon) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.all(14),
        decoration: BoxDecoration(
          color: Theme.of(context).brightness == Brightness.dark ? const Color(0xFF151D2A) : Colors.white,
          borderRadius: BorderRadius.circular(20),
          boxShadow: [
            BoxShadow(
              color: color.withOpacity(0.08),
              blurRadius: 16,
              offset: const Offset(0, 6),
            ),
          ],
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Container(
              padding: const EdgeInsets.all(7),
              decoration: BoxDecoration(
                color: color.withOpacity(0.1),
                borderRadius: BorderRadius.circular(10),
              ),
              child: Icon(icon, color: color, size: 16),
            ),
            const SizedBox(height: 10),
            Text(label, style: const TextStyle(color: Color(0xFF94A3B8), fontSize: 11, fontWeight: FontWeight.w600)),
            const SizedBox(height: 2),
            Text(value, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
          ],
        ),
      ),
    );
  }

  String _formatNum(double v) {
    if (v >= 10000000) return '${(v / 10000000).toStringAsFixed(1)}Cr';
    if (v >= 100000) return '${(v / 100000).toStringAsFixed(1)}L';
    if (v >= 1000) return '${(v / 1000).toStringAsFixed(1)}K';
    return v.toStringAsFixed(0);
  }

  String _greeting() {
    final h = DateTime.now().hour;
    if (h < 12) return 'Good morning';
    if (h < 17) return 'Good afternoon';
    return 'Good evening';
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
