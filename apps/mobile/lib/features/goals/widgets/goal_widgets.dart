import 'package:flutter/material.dart';
import '../models/goal_models.dart';

Color _importanceColor(String imp) {
  switch (imp) { case 'Mandatory': return Colors.red; case 'Essential': return Colors.orange; case 'Lifestyle': return Colors.blue; case 'Dream': return Colors.purple; default: return Colors.grey; }
}

Color _statusColor(String status) {
  switch (status) { case 'Active': return Colors.green; case 'Paused': return Colors.orange; case 'Completed': return Colors.blue; case 'AtRisk': return Colors.red; default: return Colors.grey; }
}

String _fmt(double v) {
  if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr';
  if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L';
  return '₹${v.toInt()}';
}

class GoalListCard extends StatelessWidget {
  final GoalSummary goal;
  final VoidCallback? onTap;
  const GoalListCard({super.key, required this.goal, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = _statusColor(goal.status);
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(children: [
                Expanded(child: Text(goal.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold))),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                  decoration: BoxDecoration(color: color.withValues(alpha: 0.15), borderRadius: BorderRadius.circular(8)),
                  child: Text(goal.status, style: TextStyle(color: color, fontSize: 11, fontWeight: FontWeight.w600)),
                ),
              ]),
              const SizedBox(height: 8),
              Row(children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                  decoration: BoxDecoration(color: _importanceColor(goal.importance).withValues(alpha: 0.1), borderRadius: BorderRadius.circular(4)),
                  child: Text(goal.importance, style: TextStyle(color: _importanceColor(goal.importance), fontSize: 10)),
                ),
                const SizedBox(width: 8),
                Text('Priority: ${goal.priority}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                const Spacer(),
                if (goal.hasRecommendation)
                  const Icon(Icons.lightbulb, color: Colors.amber, size: 16),
              ]),
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(4),
                      child: LinearProgressIndicator(value: goal.progress / 100.0, backgroundColor: color.withValues(alpha: 0.1), color: color, minHeight: 8),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Text('${goal.progress}%', style: TextStyle(fontWeight: FontWeight.bold, color: color, fontSize: 14)),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class GoalDashboardCard extends StatelessWidget {
  final String title;
  final String summary;
  final String value;
  final IconData icon;
  final Color color;
  final VoidCallback? onTap;
  const GoalDashboardCard({super.key, required this.title, required this.summary, required this.value, required this.icon, required this.color, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(14),
          child: Row(children: [
            Container(
              width: 40, height: 40,
              decoration: BoxDecoration(color: color.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(10)),
              child: Icon(icon, color: color, size: 20),
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
                  Text(summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 1, overflow: TextOverflow.ellipsis),
                ],
              ),
            ),
            if (value.isNotEmpty) Text(value, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
          ]),
        ),
      ),
    );
  }
}

class GoalHeaderCard extends StatelessWidget {
  final GoalDashboardData goal;
  const GoalHeaderCard({super.key, required this.goal});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = _statusColor(goal.status);
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(children: [
              Container(
                width: 48, height: 48,
                decoration: BoxDecoration(color: color.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(12)),
                child: Icon(Icons.flag, color: color, size: 24),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(goal.goalName, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                    Text('${goal.importance} · ${goal.status}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  ],
                ),
              ),
            ]),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(_fmt(goal.currentValue), style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
                      Text('of ${_fmt(goal.targetAmount)}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                    ],
                  ),
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    Text('${goal.progress}%', style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold, color: color)),
                    Text(goal.targetDate ?? 'No target date', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 12),
            ClipRRect(
              borderRadius: BorderRadius.circular(4),
              child: LinearProgressIndicator(value: goal.progress / 100.0, backgroundColor: color.withValues(alpha: 0.1), color: color, minHeight: 10),
            ),
            if (goal.fundingGap > 0) ...[
              const SizedBox(height: 8),
              Row(children: [
                const Icon(Icons.warning_amber, size: 14, color: Colors.orange),
                const SizedBox(width: 4),
                Text('Gap: ${_fmt(goal.fundingGap)}', style: const TextStyle(color: Colors.orange, fontSize: 12)),
              ]),
            ],
          ],
        ),
      ),
    );
  }
}
