import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/insight_models.dart';
import '../state/insight_state.dart';
import '../widgets/insight_widgets.dart';

class InsightsPage extends ConsumerStatefulWidget {
  const InsightsPage({super.key});
  @override
  ConsumerState<InsightsPage> createState() => _InsightsPageState();
}

class _InsightsPageState extends ConsumerState<InsightsPage> {
  final _searchCtrl = TextEditingController();
  Timer? _debounce;
  bool _searching = false;

  @override
  void initState() { super.initState(); Future.microtask(() => ref.read(insightProvider.notifier).init()); }

  @override
  void dispose() { _searchCtrl.dispose(); _debounce?.cancel(); super.dispose(); }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(insightProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(
        title: _searching
            ? TextField(controller: _searchCtrl, autofocus: true,
                decoration: InputDecoration(hintText: 'Search insights...', border: InputBorder.none,
                  suffixIcon: IconButton(icon: const Icon(Icons.close), onPressed: () { setState(() { _searching = false; _searchCtrl.clear(); }); ref.read(insightProvider.notifier).refresh(); }),
                ),
                onChanged: (v) { _debounce?.cancel(); _debounce = Timer(const Duration(milliseconds: 300), () => ref.read(insightProvider.notifier).search(v.trim())); },
              )
            : const Text('Insights'),
        actions: [
          IconButton(icon: Icon(_searching ? Icons.search_off : Icons.search), onPressed: () => setState(() => _searching = !_searching)),
        ],
      ),
      body: _buildBody(theme, state),
    );
  }

  Widget _buildBody(ThemeData theme, InsightState state) {
    switch (state.status) {
      case InsightStatus.initial:
      case InsightStatus.loading:
        return const SharedSkeletonList();
      case InsightStatus.error:
        return SharedErrorView(message: state.error, onRetry: () => ref.read(insightProvider.notifier).refresh());
      case InsightStatus.empty:
        return const SharedEmptyView(icon: Icons.insights, title: 'No insights found', subtitle: 'Check back later for personalized insights');
      case InsightStatus.loaded:
        return RefreshIndicator(
          onRefresh: () => ref.read(insightProvider.notifier).refresh(),
          child: ListView(
            padding: const EdgeInsets.only(bottom: AppSpacing.md),
            children: [
              if (state.summary != null) InsightSummaryCard(summary: state.summary!),
              const SizedBox(height: AppSpacing.md),
              _buildTabs(theme, state),
              const SizedBox(height: AppSpacing.sm),
              ..._currentItems(state).map((item) => Padding(
                padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md, vertical: AppSpacing.xs),
                child: InsightCard(insight: item, onTap: () => context.push('/insights/${item.id}')),
              )),
              if (state.trends != null) ...[
                const SizedBox(height: AppSpacing.sm),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
                  child: Card(child: Padding(
                    padding: const EdgeInsets.all(AppSpacing.md),
                    child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                      const SharedSectionHeader(title: 'Trends'),
                      const SizedBox(height: AppSpacing.sm),
                      InsightTrendChart(data: state.trends!),
                    ]),
                  )),
                ),
              ],
            ],
          ),
        );
    }
  }

  Widget _buildTabs(ThemeData theme, InsightState state) {
    final tabs = ['All', 'Opportunities', 'Warnings', 'Achievements', 'Forecast', 'Trends'];
    return SizedBox(
      height: 40,
      child: ListView(
        scrollDirection: Axis.horizontal,
        padding: const EdgeInsets.symmetric(horizontal: AppSpacing.md),
        children: tabs.map((tab) => Padding(
          padding: const EdgeInsets.only(right: AppSpacing.sm),
          child: ChoiceChip(
            label: Text(tab),
            selected: state.activeTab == tab.toLowerCase(),
            onSelected: (_) => ref.read(insightProvider.notifier).loadTab(tab.toLowerCase()),
            visualDensity: VisualDensity.compact,
          ),
        )).toList(),
      ),
    );
  }

  List<InsightItem> _currentItems(InsightState state) {
    switch (state.activeTab) {
      case 'opportunities': return state.opportunities;
      case 'warnings': return state.warnings;
      case 'achievements': return state.achievements;
      case 'forecast': return state.achievements;
      default: return state.insights;
    }
  }
}
