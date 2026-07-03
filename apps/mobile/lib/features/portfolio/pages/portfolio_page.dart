import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_animate/flutter_animate.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/pf_models.dart';
import '../repository/pf_repository.dart';
import '../widgets/pf_widgets.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';

final pfStateProvider = StateNotifierProvider<PfStateNotifier, PfState>((ref) {
  return PfStateNotifier(ref.read(pfRepositoryProvider));
});

class PfState {
  final bool loading; final PfDashboardData? dash; final AllocationData? alloc;
  final PerformanceData? perf; final RiskData? risk; final ProjectionData? proj;
  final CardViewData? recs; final CardViewData? opts; final SimData? sims;
  final CardViewData? tl; final String? error;
  PfState({this.loading = false, this.dash, this.alloc, this.perf, this.risk, this.proj, this.recs, this.opts, this.sims, this.tl, this.error});
}

class PfStateNotifier extends StateNotifier<PfState> {
  final PfRepository _repo;
  PfStateNotifier(this._repo) : super(PfState());

  Future<void> load({String? userId}) async {
    state = PfState(loading: true);
    try {
      final results = await Future.wait([
        _repo.getDashboard(userId: userId),
        _repo.getAllocation(userId: userId),
        _repo.getPerformance(userId: userId),
        _repo.getRisk(userId: userId),
        _repo.getProjection(userId: userId),
        _repo.getRecommendations(userId: userId),
        _repo.getOptimization(userId: userId),
        _repo.getSimulations(userId: userId),
        _repo.getTimeline(userId: userId),
      ]);
      state = PfState(
        dash: (results[0] as PfDashboardResponse).data,
        alloc: (results[1] as AllocationResponse).data,
        perf: (results[2] as PerformanceResponse).data,
        risk: (results[3] as RiskResponse).data,
        proj: (results[4] as ProjectionResponse).data,
        recs: (results[5] as CardViewResponse).data,
        opts: (results[6] as CardViewResponse).data,
        sims: (results[7] as SimResponse).data,
        tl: (results[8] as CardViewResponse).data,
      );
    } catch (e) {
      state = PfState(error: e.toString());
    }
  }
}

class PortfolioPage extends ConsumerStatefulWidget {
  const PortfolioPage({super.key});
  @override
  ConsumerState<PortfolioPage> createState() => _PortfolioPageState();
}

class _PortfolioPageState extends ConsumerState<PortfolioPage> {
  @override
  void initState() { super.initState(); Future.microtask(() => ref.read(pfStateProvider.notifier).load()); }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(pfStateProvider);
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: Container(
        decoration: BoxDecoration(
          gradient: AppColors.backgroundGradient(isDark),
        ),
        child: _buildBody(theme, state),
      ),
    );
  }

  Widget _buildBody(ThemeData theme, PfState state) {
    if (state.loading) {
      return Center(
        child: CircularProgressIndicator(color: AppColors.teal500).animate().fade(),
      );
    }
    if (state.error != null) {
      return SharedErrorView(
        message: state.error,
        onRetry: () => ref.read(pfStateProvider.notifier).load(),
      );
    }

    final children = <Widget>[
      if (state.dash != null) PfSummaryCard(dash: state.dash!),
      if (state.dash != null) ...state.dash!.cards.map((c) => Padding(padding: const EdgeInsets.only(top: AppSpacing.md), child: PortfolioCard(card: c))),
      if (state.alloc != null) Padding(padding: const EdgeInsets.only(top: AppSpacing.md), child: AllocationCard(alloc: state.alloc!)),
      if (state.perf != null) 
        Padding(
          padding: const EdgeInsets.only(top: AppSpacing.md),
          child: PfSectionCard(title: 'Performance', icon: Icons.trending_up_rounded, color: AppColors.teal500, rows: [
            MapEntry('Period Return', '${state.perf!.periodReturnPct.toStringAsFixed(1)}%'),
            MapEntry('Benchmark', '${state.perf!.benchmarkReturn.toStringAsFixed(1)}%'),
            MapEntry('Unrealized G/L', _fmt(state.perf!.unrealizedGL)),
            MapEntry('Realized G/L', _fmt(state.perf!.realizedGL)),
          ]),
        ),
      if (state.risk != null) 
        Padding(
          padding: const EdgeInsets.only(top: AppSpacing.md),
          child: PfSectionCard(title: 'Risk Assessment', icon: Icons.shield_rounded, color: AppColors.amber500, rows: [
            MapEntry('Score', '${state.risk!.riskScore} (${state.risk!.riskLevel})'),
            MapEntry('VaR', '${state.risk!.var_.toStringAsFixed(1)}%'),
            MapEntry('Sharpe Ratio', state.risk!.sharpeRatio.toStringAsFixed(2)),
            MapEntry('Volatility', '${state.risk!.volatility.toStringAsFixed(1)}%'),
            MapEntry('Max Drawdown', '${state.risk!.maxDrawdown.toStringAsFixed(1)}%'),
          ]),
        ),
      if (state.proj != null) 
        Padding(
          padding: const EdgeInsets.only(top: AppSpacing.md),
          child: PfSectionCard(title: 'Projection', icon: Icons.query_stats_rounded, color: Colors.cyan, rows: [
            MapEntry('Projected Value', _fmt(state.proj!.projectedValue.toInt())),
            MapEntry('Confidence', state.proj!.confidence),
            MapEntry('Horizon', '${state.proj!.horizonYears} years'),
            MapEntry('Annual Return', '${state.proj!.annualReturn.toStringAsFixed(1)}%'),
          ]),
        ),
      if (state.recs != null && state.recs!.items.isNotEmpty) 
        Padding(
          padding: const EdgeInsets.only(top: AppSpacing.md),
          child: PfSectionCard(title: 'Recommendations', icon: Icons.lightbulb_rounded, color: AppColors.amber500, rows: state.recs!.items.map((i) => MapEntry(i.title, i.value)).toList()),
        ),
      if (state.opts != null && state.opts!.items.isNotEmpty) 
        Padding(
          padding: const EdgeInsets.only(top: AppSpacing.md),
          child: PfSectionCard(title: 'Optimizations', icon: Icons.auto_graph_rounded, color: Colors.purple, rows: state.opts!.items.map((i) => MapEntry(i.title, i.value)).toList()),
        ),
      if (state.sims != null && state.sims!.simulations.isNotEmpty) 
        Padding(
          padding: const EdgeInsets.only(top: AppSpacing.md),
          child: PfSectionCard(title: 'Simulations', icon: Icons.science_rounded, color: Colors.deepOrange, rows: state.sims!.simulations.map((s) => MapEntry(s.scenario, s.outcome)).toList()),
        ),
      if (state.tl != null && state.tl!.items.isNotEmpty) 
        Padding(
          padding: const EdgeInsets.only(top: AppSpacing.md),
          child: PfSectionCard(title: 'Timeline', icon: Icons.history_rounded, color: Colors.brown, rows: state.tl!.items.map((i) => MapEntry(i.title, i.value)).toList()),
        ),
      const SizedBox(height: 100), // padding for floating navbar
    ];

    return RefreshIndicator(
      color: AppColors.teal500,
      onRefresh: () => ref.read(pfStateProvider.notifier).load(),
      child: ListView.builder(
        padding: const EdgeInsets.all(AppSpacing.md),
        itemCount: children.length,
        itemBuilder: (context, index) {
          return children[index]
            .animate()
            .fade(duration: 400.ms)
            .slideY(begin: 0.1, end: 0, duration: 400.ms, curve: Curves.easeOutQuad);
        },
      ),
    );
  }
}

String _fmt(int v) {
  if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr';
  if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L';
  return '₹$v';
}
