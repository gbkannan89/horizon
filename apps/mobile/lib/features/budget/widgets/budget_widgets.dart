import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';
import 'package:horizon_mobile/app/theme.dart';
import '../models/budget_models.dart';

class BudgetProgressBar extends StatelessWidget {
  final double spent;
  final double budgeted;
  final double height;
  final Color? color;

  const BudgetProgressBar({
    super.key,
    required this.spent,
    required this.budgeted,
    this.height = 8,
    this.color,
  });

  @override
  Widget build(BuildContext context) {
    final value = budgeted > 0 ? (spent / budgeted).clamp(0.0, 1.0) : 0.0;
    final barColor = color ?? (spent > budgeted ? Colors.red : Colors.green);
    return LinearProgressIndicator(
      value: value,
      backgroundColor: Colors.grey.shade200,
      color: barColor,
      minHeight: height,
      borderRadius: BorderRadius.circular(height / 2),
    );
  }
}

class VsActualChart extends StatelessWidget {
  final List<BudgetCategoryModel> categories;

  const VsActualChart({super.key, required this.categories});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    if (categories.isEmpty) {
      return const SizedBox(
        height: 200,
        child: Center(child: Text('No categories to display')),
      );
    }
    return SizedBox(
      height: 200,
      child: BarChart(
        BarChartData(
          alignment: BarChartAlignment.spaceAround,
          maxY: categories.fold<double>(0, (max, c) => max > c.budgetedAmount ? max : c.budgetedAmount) * 1.2,
          barTouchData: BarTouchData(
            touchTooltipData: BarTouchTooltipData(
              getTooltipItem: (group, groupIndex, rod, rodIndex) {
                final cat = categories[groupIndex];
                final label = rodIndex == 0 ? 'Budget' : 'Spent';
                return BarTooltipItem('${cat.category}\n$label: ₹${rod.toY.toInt()}', const TextStyle(color: Colors.white, fontSize: 12));
              },
            ),
          ),
          titlesData: FlTitlesData(
            show: true,
            bottomTitles: AxisTitles(
              sideTitles: SideTitles(
                showTitles: true,
                reservedSize: 36,
                getTitlesWidget: (value, meta) {
                  final idx = value.toInt();
                  if (idx < 0 || idx >= categories.length) return const SizedBox();
                  return SideTitleWidget(
                    axisSide: meta.axisSide,
                    child: Text(
                      categories[idx].category.length > 6
                          ? '${categories[idx].category.substring(0, 6)}...'
                          : categories[idx].category,
                      style: theme.textTheme.labelSmall,
                    ),
                  );
                },
              ),
            ),
            topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
            rightTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
          ),
          borderData: FlBorderData(show: false),
          barGroups: categories.asMap().entries.map((entry) {
            final idx = entry.key;
            final cat = entry.value;
            return BarChartGroupData(
              x: idx,
              barRods: [
                BarChartRodData(
                  toY: cat.budgetedAmount,
                  color: Colors.blue.shade300,
                  width: 12,
                  borderRadius: const BorderRadius.vertical(top: Radius.circular(4)),
                ),
                BarChartRodData(
                  toY: cat.spentAmount,
                  color: cat.spentAmount > cat.budgetedAmount ? Colors.red : Colors.green.shade400,
                  width: 12,
                  borderRadius: const BorderRadius.vertical(top: Radius.circular(4)),
                ),
              ],
            );
          }).toList(),
        ),
      ),
    );
  }
}

class BudgetStatusBadge extends StatelessWidget {
  final String status;

  const BudgetStatusBadge({super.key, required this.status});

  @override
  Widget build(BuildContext context) {
    final (color, icon) = switch (status) {
      'Draft' => (Colors.grey, Icons.edit_note),
      'Active' => (Colors.green, Icons.check_circle),
      'Paused' => (Colors.orange, Icons.pause_circle),
      'Completed' => (Colors.blue, Icons.task_alt),
      'Archived' => (Colors.grey, Icons.archive),
      _ => (Colors.grey, Icons.help_outline),
    };
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.sm, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.15),
        borderRadius: BorderRadius.circular(AppRadius.sm),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 14, color: color),
          const SizedBox(width: 4),
          Text(status, style: TextStyle(fontSize: 12, color: color, fontWeight: FontWeight.w500)),
        ],
      ),
    );
  }
}
