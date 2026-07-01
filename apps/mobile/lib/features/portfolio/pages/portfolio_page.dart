import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/pf_models.dart';
import '../repository/pf_repository.dart';
import '../widgets/pf_widgets.dart';

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

    return Scaffold(
      appBar: AppBar(title: const Text('Portfolio')),
      body: _buildBody(theme, state),
    );
  }

  Widget _buildBody(ThemeData theme, PfState state) {
    if (state.loading) return const SharedLoadingView();
    if (state.error != null) return SharedErrorView(
      message: state.error,
      onRetry: () => ref.read(pfStateProvider.notifier).load(),
    );

    return RefreshIndicator(
      onRefresh: () => ref.read(pfStateProvider.notifier).load(),
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          if (state.dash != null) PfSummaryCard(dash: state.dash!),
          const SizedBox(height: 16),
          if (state.dash != null) ...state.dash!.cards.map((c) => Padding(padding: const EdgeInsets.only(bottom: 8), child: PortfolioCard(card: c))),
          if (state.alloc != null) ...[const SizedBox(height: 8), AllocationCard(alloc: state.alloc!)],
          if (state.perf != null) ...[
            const SizedBox(height: 8),
            PfSectionCard(title: 'Performance', icon: Icons.trending_up, color: Colors.green, rows: [
              MapEntry('Period Return', '${state.perf!.periodReturnPct.toStringAsFixed(1)}%'),
              MapEntry('Benchmark', '${state.perf!.benchmarkReturn.toStringAsFixed(1)}%'),
              MapEntry('Unrealized G/L', _fmt(state.perf!.unrealizedGL)),
              MapEntry('Realized G/L', _fmt(state.perf!.realizedGL)),
            ]),
          ],
          if (state.risk != null) ...[
            const SizedBox(height: 8),
            PfSectionCard(title: 'Risk Assessment', icon: Icons.shield, color: Colors.amber, rows: [
              MapEntry('Score', '${state.risk!.riskScore} (${state.risk!.riskLevel})'),
              MapEntry('VaR', '${state.risk!.var_.toStringAsFixed(1)}%'),
              MapEntry('Sharpe Ratio', state.risk!.sharpeRatio.toStringAsFixed(2)),
              MapEntry('Volatility', '${state.risk!.volatility.toStringAsFixed(1)}%'),
              MapEntry('Max Drawdown', '${state.risk!.maxDrawdown.toStringAsFixed(1)}%'),
            ]),
          ],
          if (state.proj != null) ...[
            const SizedBox(height: 8),
            PfSectionCard(title: 'Projection', icon: Icons.query_stats, color: Colors.cyan, rows: [
              MapEntry('Projected Value', _fmt(state.proj!.projectedValue.toInt())),
              MapEntry('Confidence', state.proj!.confidence),
              MapEntry('Horizon', '${state.proj!.horizonYears} years'),
              MapEntry('Annual Return', '${state.proj!.annualReturn.toStringAsFixed(1)}%'),
            ]),
          ],
          if (state.recs != null && state.recs!.items.isNotEmpty) ...[
            const SizedBox(height: 8),
            PfSectionCard(title: 'Recommendations', icon: Icons.lightbulb, color: Colors.amber, rows: state.recs!.items.map((i) => MapEntry(i.title, i.value)).toList()),
          ],
          if (state.opts != null && state.opts!.items.isNotEmpty) ...[
            const SizedBox(height: 8),
            PfSectionCard(title: 'Optimizations', icon: Icons.auto_graph, color: Colors.purple, rows: state.opts!.items.map((i) => MapEntry(i.title, i.value)).toList()),
          ],
          if (state.sims != null && state.sims!.simulations.isNotEmpty) ...[
            const SizedBox(height: 8),
            PfSectionCard(title: 'Simulations', icon: Icons.science, color: Colors.deepOrange, rows: state.sims!.simulations.map((s) => MapEntry(s.scenario, s.outcome)).toList()),
          ],
          if (state.tl != null && state.tl!.items.isNotEmpty) ...[
            const SizedBox(height: 8),
            PfSectionCard(title: 'Timeline', icon: Icons.history, color: Colors.brown, rows: state.tl!.items.map((i) => MapEntry(i.title, i.value)).toList()),
          ],
        ],
      ),
    );
  }
}

String _fmt(int v) {
  if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr';
  if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L';
  return '₹$v';
}
