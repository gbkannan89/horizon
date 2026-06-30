import 'package:flutter/material.dart';
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
    final theme = Theme.of(context);
    final color = _cardColor(card.cardType);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(14),
        child: Row(children: [
          Container(
            width: 40, height: 40,
            decoration: BoxDecoration(color: color.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(10)),
            child: Icon(_cardIcon(card.cardType), color: color, size: 20),
          ),
          const SizedBox(width: 12),
          Expanded(child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(card.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
              Text(card.summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 1, overflow: TextOverflow.ellipsis),
            ],
          )),
        ]),
      ),
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
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(Icons.account_balance, color: theme.colorScheme.primary, size: 24),
              const SizedBox(width: 8),
              Text('Portfolio', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
            ]),
            const SizedBox(height: 16),
            Text(_fmt(dash.portfolioValue), style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold)),
            const SizedBox(height: 4),
            Row(children: [
              Icon(dash.totalReturn >= 0 ? Icons.trending_up : Icons.trending_down, size: 16, color: dash.totalReturn >= 0 ? Colors.green : Colors.red),
              const SizedBox(width: 4),
              Text('${dash.totalReturnPct.toStringAsFixed(1)}% (${_fmt(dash.totalReturn.toInt())})', style: TextStyle(color: dash.totalReturn >= 0 ? Colors.green : Colors.red, fontWeight: FontWeight.w600)),
              const SizedBox(width: 12),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                decoration: BoxDecoration(color: color.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(8)),
                child: Text('Risk: ${dash.riskScore} (${dash.riskLevel})', style: TextStyle(color: color, fontSize: 11, fontWeight: FontWeight.w600)),
              ),
            ]),
          ],
        ),
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
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Asset Allocation', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 4),
            Text('Diversification: ${alloc.diversification.toStringAsFixed(1)}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const SizedBox(height: 12),
            ...alloc.allocations.map((a) => Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
                    Text(a.label, style: theme.textTheme.bodyMedium),
                    Text('${a.percent.toStringAsFixed(1)}%', style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                  ]),
                  const SizedBox(height: 4),
                  ClipRRect(
                    borderRadius: BorderRadius.circular(3),
                    child: LinearProgressIndicator(value: a.percent / 100.0, minHeight: 6, color: Colors.blue.shade300, backgroundColor: Colors.grey.withValues(alpha: 0.15)),
                  ),
                ],
              ),
            )),
          ],
        ),
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
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Icon(icon, color: color, size: 18),
              const SizedBox(width: 8),
              Text(title, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            ]),
            const SizedBox(height: 12),
            ...rows.map((r) => Padding(
              padding: const EdgeInsets.symmetric(vertical: 3),
              child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
                Text(r.key, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                Text(r.value, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
              ]),
            )),
          ],
        ),
      ),
    );
  }
}

String _fmt(int v) {
  if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr';
  if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L';
  return '₹$v';
}
