import 'package:flutter/material.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/goal_models.dart';

Color _importanceColor(String imp) {
  switch (imp) { case 'Mandatory': return Colors.red; case 'Essential': return Colors.orange; case 'Lifestyle': return Colors.blue; case 'Dream': return Colors.purple; default: return Colors.grey; }
}

Color _statusColor(String status) {
  switch (status) { case 'Active': return Colors.green; case 'Paused': return Colors.orange; case 'Completed': return Colors.blue; case 'AtRisk': return Colors.red; default: return Colors.grey; }
}

String _fmt(double v) {
  return formatMoney(v.toInt());
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
          padding: AppTheme.cardPadding,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(children: [
                Expanded(child: Text(goal.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold))),
                StatusChip(label: goal.status, color: color),
              ]),
              const SizedBox(height: AppTheme.spacingSm),
              Row(children: [
                StatusChip(label: goal.importance, color: _importanceColor(goal.importance), fontSize: 10),
                const SizedBox(width: 8),
                Text('Priority: ${goal.priority}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                const Spacer(),
                if (goal.hasRecommendation)
                  Icon(Icons.lightbulb, color: Colors.amber, size: AppTheme.iconSm),
              ]),
              const SizedBox(height: AppTheme.spacingMd),
              Row(
                children: [
                  Expanded(
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(AppTheme.radiusSm),
                      child: LinearProgressIndicator(value: goal.progress / 100.0, backgroundColor: color.withOpacity(0.1), color: color, minHeight: 8),
                    ),
                  ),
                  const SizedBox(width: AppTheme.spacingMd),
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
    return IconCard(title: title, subtitle: summary, trailing: value, icon: icon, color: color, onTap: onTap);
  }
}

class GoalHeaderCard extends StatelessWidget {
  final GoalDashboardData goal;
  const GoalHeaderCard({super.key, required this.goal});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = _statusColor(goal.status);
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Container(
              width: 48, height: 48,
              decoration: BoxDecoration(color: color.withOpacity(0.12), borderRadius: BorderRadius.circular(12)),
              child: Icon(Icons.flag, color: color, size: 24),
            ),
            const SizedBox(width: AppTheme.spacingMd),
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
          const SizedBox(height: AppTheme.spacingLg),
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
          const SizedBox(height: AppTheme.spacingMd),
          ClipRRect(
            borderRadius: BorderRadius.circular(AppTheme.radiusSm),
            child: LinearProgressIndicator(value: goal.progress / 100.0, backgroundColor: color.withOpacity(0.1), color: color, minHeight: 10),
          ),
          if (goal.fundingGap > 0) ...[
            const SizedBox(height: AppTheme.spacingSm),
            Row(children: [
              Icon(Icons.warning_amber, size: 14, color: Colors.orange),
              const SizedBox(width: AppTheme.spacingXs),
              Text('Gap: ${_fmt(goal.fundingGap)}', style: TextStyle(color: Colors.orange, fontSize: 12)),
            ]),
          ],
        ],
      ),
    );
  }
}
