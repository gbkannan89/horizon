import 'package:flutter/material.dart';
import 'package:fl_chart/fl_chart.dart';
import 'package:horizon_mobile/app/theme.dart';
import '../models/insight_models.dart';

class InsightCard extends StatelessWidget {
  final InsightItem insight;
  final VoidCallback? onTap;
  const InsightCard({super.key, required this.insight, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(AppSpacing.md),
          child: Row(
            children: [
              Container(
                width: 40, height: 40,
                decoration: BoxDecoration(
                  color: insight.priorityColor.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(AppRadius.sm),
                ),
                child: Icon(insight.categoryIcon, size: AppIconSize.md, color: insight.priorityColor),
              ),
              const SizedBox(width: AppSpacing.sm),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(insight.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500), maxLines: 1, overflow: TextOverflow.ellipsis),
                    Text(insight.summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 2, overflow: TextOverflow.ellipsis),
                  ],
                ),
              ),
              if (insight.priority.isNotEmpty)
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                  decoration: BoxDecoration(color: insight.priorityColor.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(4)),
                  child: Text(insight.priority, style: theme.textTheme.labelSmall?.copyWith(color: insight.priorityColor, fontWeight: FontWeight.w600)),
                ),
            ],
          ),
        ),
      ),
    );
  }
}

class InsightTrendChart extends StatefulWidget {
  final TrendData data;
  const InsightTrendChart({super.key, required this.data});
  @override
  State<InsightTrendChart> createState() => _InsightTrendChartState();
}

class _InsightTrendChartState extends State<InsightTrendChart> {
  String _period = '1M';

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final data = widget.data;
    if (data.values.isEmpty) {
      return const SizedBox(height: 200, child: Center(child: Text('No trend data available')));
    }
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.end,
          children: ['1M', '3M', '6M', '1Y'].map((p) => Padding(
            padding: const EdgeInsets.only(left: AppSpacing.xs),
            child: ChoiceChip(
              label: Text(p, style: theme.textTheme.labelSmall),
              selected: _period == p,
              onSelected: (_) => setState(() => _period = p),
              visualDensity: VisualDensity.compact,
            ),
          )).toList(),
        ),
        const SizedBox(height: AppSpacing.sm),
        SizedBox(height: 200,
      child: LineChart(
        LineChartData(
          gridData: FlGridData(show: true, drawVerticalLine: false, horizontalInterval: 1, getDrawingHorizontalLine: (value) => FlLine(color: theme.colorScheme.outlineVariant, strokeWidth: 1)),
          titlesData: FlTitlesData(
            leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, reservedSize: 40)),
            bottomTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, getTitlesWidget: (value, meta) {
              final idx = value.toInt();
              if (idx >= 0 && idx < data.labels.length) return Padding(padding: const EdgeInsets.only(top: 4), child: Text(data.labels[idx], style: theme.textTheme.labelSmall));
              return const SizedBox();
            })),
            topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
            rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
          ),
          borderData: FlBorderData(show: false),
          minX: 0, maxX: (data.values.length - 1).toDouble(),
          minY: data.values.reduce((a, b) => a < b ? a : b) * 0.9,
          maxY: data.values.reduce((a, b) => a > b ? a : b) * 1.1,
          lineBarsData: [LineChartBarData(
            spots: data.values.asMap().entries.map((e) => FlSpot(e.key.toDouble(), e.value)).toList(),
            isCurved: true, color: theme.colorScheme.primary, barWidth: 3,
            belowBarData: BarAreaData(show: true, color: theme.colorScheme.primary.withValues(alpha: 0.1)),
            dotData: FlDotData(show: true),
          )],
        ),
      ),
      ),
    ],
    );
  }
}

class InsightSummaryCard extends StatelessWidget {
  final InsightDashboard summary;
  const InsightSummaryCard({super.key, required this.summary});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      margin: const EdgeInsets.fromLTRB(AppSpacing.md, AppSpacing.md, AppSpacing.md, 0),
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.md),
        child: Row(
          children: [
            _stat(theme, 'Total', '${summary.total}', theme.colorScheme.primary),
            _divider(theme),
            _stat(theme, 'Critical', '${summary.criticalCount}', Colors.red),
            _divider(theme),
            _stat(theme, 'New', '${summary.newCount}', Colors.orange),
          ],
        ),
      ),
    );
  }

  Widget _stat(ThemeData theme, String label, String value, Color color) {
    return Expanded(child: Column(children: [
      Text(value, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold, color: color)),
      Text(label, style: theme.textTheme.labelSmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
    ]));
  }
  Widget _divider(ThemeData theme) => Container(width: 1, height: 32, color: theme.colorScheme.outlineVariant);
}
