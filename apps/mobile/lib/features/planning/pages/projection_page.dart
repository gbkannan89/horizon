import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/shared/widgets/shared_widgets.dart';
import '../models/planning_models.dart';
import '../repository/planning_repository.dart';

final projectionProvider = FutureProvider.autoDispose<ProjectionData>((ref) async {
  return ref.read(planningRepositoryProvider).getProjections();
});

class ProjectionPage extends ConsumerWidget {
  const ProjectionPage({super.key});
  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final async = ref.watch(projectionProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Projections')),
      body: async.when(
        loading: () => const SharedLoadingView(),
        error: (e, _) => SharedErrorView(message: e.toString()),
        data: (data) {
          if (data.points.isEmpty) return const SharedEmptyView(icon: Icons.query_stats, title: 'No projection data');
          return Padding(
            padding: const EdgeInsets.all(AppSpacing.md),
            child: Column(children: [
              Card(child: Padding(
                padding: const EdgeInsets.all(AppSpacing.md),
                child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Text('Income & Expenses', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                  const SizedBox(height: AppSpacing.md),
                  SizedBox(height: 250, child: LineChart(LineChartData(
                    gridData: const FlGridData(show: true, drawVerticalLine: false),
                    titlesData: FlTitlesData(leftTitles: const AxisTitles(sideTitles: SideTitles(showTitles: true, reservedSize: 40)), bottomTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, getTitlesWidget: (v, _) { final i = v.toInt(); if (i >= 0 && i < data.points.length) return Padding(padding: const EdgeInsets.only(top: 4), child: Text(data.points[i].period, style: theme.textTheme.labelSmall)); return const SizedBox(); })), topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)), rightTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false))),
                    borderData: FlBorderData(show: false),
                    lineBarsData: [
                      LineChartBarData(spots: data.points.asMap().entries.map((e) => FlSpot(e.key.toDouble(), e.value.income)).toList(), isCurved: true, color: Colors.green, barWidth: 2, dotData: const FlDotData(show: true)),
                      LineChartBarData(spots: data.points.asMap().entries.map((e) => FlSpot(e.key.toDouble(), e.value.expenses)).toList(), isCurved: true, color: Colors.red, barWidth: 2, dotData: const FlDotData(show: true)),
                      LineChartBarData(spots: data.points.asMap().entries.map((e) => FlSpot(e.key.toDouble(), e.value.netWorth)).toList(), isCurved: true, color: theme.colorScheme.primary, barWidth: 2, dotData: const FlDotData(show: true), belowBarData: BarAreaData(show: true, color: theme.colorScheme.primary.withValues(alpha: 0.05))),
                    ],
                  ))),
                ]),
              )),
            ]),
          );
        },
      ),
    );
  }
}
