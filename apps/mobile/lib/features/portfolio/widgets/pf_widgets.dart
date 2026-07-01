import 'package:flutter/material.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/pf_models.dart';

Color _riskColor(int score) {
  if (score >= 80) return Colors.red;
  if (score >= 60) return Colors.orange;
  if (score >= 40) return Colors.amber;
  return Colors.green;
}

Color _cardColor(String type) {
  switch (type) {
    case 'Overview': return Colors.teal;
    case 'Allocation': return Colors.blue;
    case 'Performance': return Colors.green;
    case 'Holdings': return Colors.indigo;
    case 'Projection': return Colors.cyan;
    case 'Risk': return _riskColor(0);
    case 'Recommendation': return Colors.amber;
    case 'Optimization': return Colors.purple;
    case 'Simulation': return Colors.deepOrange;
    case 'Timeline': return Colors.brown;
    default: return Colors.grey;
  }
}

IconData _cardIcon(String type) {
  switch (type) {
    case 'Overview': return Icons.account_balance;
    case 'Allocation': return Icons.pie_chart;
    case 'Performance': return Icons.trending_up;
    case 'Holdings': return Icons.inventory_2;
    case 'Projection': return Icons.query_stats;
    case 'Risk': return Icons.shield;
    case 'Recommendation': return Icons.lightbulb;
    case 'Optimization': return Icons.auto_graph;
    case 'Simulation': return Icons.science;
    case 'Timeline': return Icons.history;
    default: return Icons.circle;
  }
}

class PortfolioCard extends StatelessWidget {
  final PfCard card;
  const PortfolioCard({super.key, required this.card});

  @override
  Widget build(BuildContext context) {
    return IconCard(
      title: card.title, subtitle: card.summary,
      icon: _cardIcon(card.cardType), color: _cardColor(card.cardType),
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
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(Icons.account_balance, color: theme.colorScheme.primary, size: AppTheme.iconLg),
            const SizedBox(width: 8),
            Text('Portfolio', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          ]),
          const SizedBox(height: AppTheme.spacingLg),
          Text(formatMoney(dash.portfolioValue), style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold)),
          const SizedBox(height: AppTheme.spacingXs),
          Row(children: [
            Icon(dash.totalReturn >= 0 ? Icons.trending_up : Icons.trending_down, size: AppTheme.iconSm, color: dash.totalReturn >= 0 ? Colors.green : Colors.red),
            const SizedBox(width: 4),
            Text('${dash.totalReturnPct.toStringAsFixed(1)}% (${formatMoney(dash.totalReturn.toInt())})', style: TextStyle(color: dash.totalReturn >= 0 ? Colors.green : Colors.red, fontWeight: FontWeight.w600)),
            const SizedBox(width: AppTheme.spacingMd),
            StatusChip(label: 'Risk: ${dash.riskScore} (${dash.riskLevel})', color: color, fontSize: 11),
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
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SectionHeader(title: 'Asset Allocation'),
          Text('Diversification: ${alloc.diversification.toStringAsFixed(1)}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          const SizedBox(height: AppTheme.spacingMd),
          ...alloc.allocations.map((a) => Padding(
            padding: const EdgeInsets.only(bottom: 8),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
                  Text(a.label, style: theme.textTheme.bodyMedium),
                  Text('${a.percent.toStringAsFixed(1)}%', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                ]),
                const SizedBox(height: AppTheme.spacingXs),
                ClipRRect(
                  borderRadius: BorderRadius.circular(AppTheme.radiusSm - 1),
                  child: LinearProgressIndicator(value: a.percent / 100.0, minHeight: 6, color: Colors.blue.shade300, backgroundColor: Colors.grey.withOpacity(0.15)),
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
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Icon(icon, color: color, size: 18),
            const SizedBox(width: 8),
            Text(title, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
          ]),
          const SizedBox(height: AppTheme.spacingMd),
          ...rows.map((r) => Padding(
            padding: const EdgeInsets.symmetric(vertical: 3),
            child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Text(r.key, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              Text(r.value, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
            ]),
          )),
        ],
      ),
    );
  }
}
