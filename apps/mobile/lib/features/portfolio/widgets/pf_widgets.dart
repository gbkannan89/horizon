import 'package:flutter/material.dart';
import '../models/pf_models.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';
import 'package:horizon_mobile/core/ui_kit/animated_stat.dart';

Color _riskColor(int score) {
  if (score >= 80) return AppColors.red500;
  if (score >= 60) return AppColors.amber500;
  if (score >= 40) return AppColors.amber500;
  return AppColors.teal500;
}

Color _cardColor(String type) {
  switch (type) {
    case 'Overview': return AppColors.teal500;
    case 'Allocation': return Colors.blue;
    case 'Performance': return AppColors.teal500;
    case 'Holdings': return Colors.indigo;
    case 'Projection': return Colors.cyan;
    case 'Risk': return _riskColor(0);
    case 'Recommendation': return AppColors.amber500;
    case 'Optimization': return Colors.purple;
    case 'Simulation': return Colors.deepOrange;
    case 'Timeline': return Colors.brown;
    default: return AppColors.slate500;
  }
}

IconData _cardIcon(String type) {
  switch (type) {
    case 'Overview': return Icons.account_balance_rounded;
    case 'Allocation': return Icons.pie_chart_rounded;
    case 'Performance': return Icons.trending_up_rounded;
    case 'Holdings': return Icons.inventory_2_rounded;
    case 'Projection': return Icons.query_stats_rounded;
    case 'Risk': return Icons.shield_rounded;
    case 'Recommendation': return Icons.lightbulb_rounded;
    case 'Optimization': return Icons.auto_graph_rounded;
    case 'Simulation': return Icons.science_rounded;
    case 'Timeline': return Icons.history_rounded;
    default: return Icons.circle;
  }
}

class PortfolioCard extends StatelessWidget {
  final PfCard card;
  const PortfolioCard({super.key, required this.card});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = _cardColor(card.cardType);
    return GlassCard(
      child: Row(children: [
        Container(
          width: 44, height: 44,
          decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.md)),
          child: Icon(_cardIcon(card.cardType), color: color, size: 24),
        ),
        const SizedBox(width: AppSpacing.md),
        Expanded(child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(card.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
            Text(card.summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 1, overflow: TextOverflow.ellipsis),
          ],
        )),
      ]),
    );
  }
}

class PfSummaryCard extends StatelessWidget {
  final PfDashboardData dash;
  const PfSummaryCard({super.key, required this.dash});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = _riskColor(dash.riskScore);
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.account_balance_rounded, color: AppColors.teal500, size: 24),
            const SizedBox(width: 8),
            Text('Portfolio', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          ]),
          const SizedBox(height: 16),
          AnimatedStatValue(
            value: dash.portfolioValue.toDouble(),
            style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.w900),
          ),
          const SizedBox(height: 8),
          Row(children: [
            Icon(dash.totalReturn >= 0 ? Icons.trending_up_rounded : Icons.trending_down_rounded, size: 20, color: dash.totalReturn >= 0 ? AppColors.teal500 : AppColors.red500),
            const SizedBox(width: 4),
            Text('${dash.totalReturnPct.toStringAsFixed(1)}% (${_fmt(dash.totalReturn.toInt())})', style: TextStyle(color: dash.totalReturn >= 0 ? AppColors.teal500 : AppColors.red500, fontWeight: FontWeight.bold)),
            const SizedBox(width: 12),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
              decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.sm)),
              child: Text('Risk: ${dash.riskScore} (${dash.riskLevel})', style: TextStyle(color: color, fontSize: 11, fontWeight: FontWeight.bold)),
            ),
          ]),
        ],
      ),
    );
  }
}

class AllocationCard extends StatelessWidget {
  final AllocationData alloc;
  const AllocationCard({super.key, required this.alloc});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text('Asset Allocation', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w800)),
          const SizedBox(height: 4),
          Text('Diversification: ${alloc.diversification.toStringAsFixed(1)}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          const SizedBox(height: 16),
          ...alloc.allocations.map((a) => Padding(
            padding: const EdgeInsets.only(bottom: 12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
                  Text(a.label, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                  Text('${a.percent.toStringAsFixed(1)}%', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
                ]),
                const SizedBox(height: 6),
                ClipRRect(
                  borderRadius: BorderRadius.circular(AppRadius.sm),
                  child: LinearProgressIndicator(
                    value: a.percent / 100.0, 
                    minHeight: 8, 
                    color: AppColors.teal400, 
                    backgroundColor: AppColors.slate200.withOpacity(0.3)
                  ),
                ),
              ],
            ),
          )),
        ],
      ),
    );
  }
}

class PfSectionCard extends StatelessWidget {
  final String title; final List<MapEntry<String, String>> rows; final IconData icon; final Color color;
  const PfSectionCard({super.key, required this.title, required this.rows, required this.icon, required this.color});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(icon, color: color, size: 20),
            const SizedBox(width: 8),
            Text(title, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.bold)),
          ]),
          const SizedBox(height: 16),
          ...rows.map((r) => Padding(
            padding: const EdgeInsets.symmetric(vertical: 4),
            child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Text(r.key, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant, fontWeight: FontWeight.w500)),
              Text(r.value, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
            ]),
          )),
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
