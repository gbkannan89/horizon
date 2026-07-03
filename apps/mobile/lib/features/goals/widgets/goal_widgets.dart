import 'package:flutter/material.dart';
import '../models/goal_models.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';
import 'package:horizon_mobile/core/ui_kit/animated_stat.dart';

Color _importanceColor(String imp) {
  switch (imp) { 
    case 'Mandatory': return AppColors.red500; 
    case 'Essential': return AppColors.amber500; 
    case 'Lifestyle': return Colors.blue; 
    case 'Dream': return Colors.purple; 
    default: return AppColors.slate500; 
  }
}

Color _statusColor(String status) {
  switch (status) { 
    case 'Active': return AppColors.teal500; 
    case 'Paused': return AppColors.amber500; 
    case 'Completed': return Colors.blue; 
    case 'AtRisk': return AppColors.red500; 
    default: return AppColors.slate500; 
  }
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
    return GlassCard(
      onTap: onTap,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Expanded(child: Text(goal.name, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold))),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
              decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.sm)),
              child: Text(goal.status, style: TextStyle(color: color, fontSize: 11, fontWeight: FontWeight.bold)),
            ),
          ]),
          const SizedBox(height: 12),
          Row(children: [
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
              decoration: BoxDecoration(color: _importanceColor(goal.importance).withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.xs)),
              child: Text(goal.importance, style: TextStyle(color: _importanceColor(goal.importance), fontSize: 10, fontWeight: FontWeight.w600)),
            ),
            const SizedBox(width: 8),
            Text('Priority: ${goal.priority}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            const Spacer(),
            if (goal.hasRecommendation)
              const Icon(Icons.lightbulb_rounded, color: AppColors.amber500, size: 18),
          ]),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: ClipRRect(
                  borderRadius: BorderRadius.circular(AppRadius.sm),
                  child: LinearProgressIndicator(
                    value: goal.progress / 100.0, 
                    backgroundColor: color.withOpacity(0.15), 
                    color: color, 
                    minHeight: 8
                  ),
                ),
              ),
              const SizedBox(width: 12),
              AnimatedStatValue(
                value: goal.progress.toDouble(),
                format: '#',
                style: TextStyle(fontWeight: FontWeight.w900, color: color, fontSize: 14),
              ),
              Text('%', style: TextStyle(fontWeight: FontWeight.bold, color: color, fontSize: 12)),
            ],
          ),
        ],
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
    return GlassCard(
      onTap: onTap,
      child: Row(children: [
        Container(
          width: 44, height: 44,
          decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.md)),
          child: Icon(icon, color: color, size: 24),
        ),
        const SizedBox(width: AppSpacing.md),
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
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Container(
              width: 48, height: 48,
              decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.md)),
              child: Icon(Icons.flag_rounded, color: color, size: 24),
            ),
            const SizedBox(width: AppSpacing.md),
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
          const SizedBox(height: 24),
          Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    AnimatedStatValue(
                      value: goal.currentValue,
                      style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.w900),
                    ),
                    Text('of ${_fmt(goal.targetAmount)}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.end,
                children: [
                  Row(
                    children: [
                      AnimatedStatValue(
                        value: goal.progress.toDouble(),
                        format: '#',
                        style: theme.textTheme.headlineMedium?.copyWith(fontWeight: FontWeight.bold, color: color),
                      ),
                      Text('%', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold, color: color)),
                    ],
                  ),
                  Text(goal.targetDate ?? 'No target date', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                ],
              ),
            ],
          ),
          const SizedBox(height: 16),
          ClipRRect(
            borderRadius: BorderRadius.circular(AppRadius.sm),
            child: LinearProgressIndicator(
              value: goal.progress / 100.0, 
              backgroundColor: color.withOpacity(0.15), 
              color: color, 
              minHeight: 12
            ),
          ),
          if (goal.fundingGap > 0) ...[
            const SizedBox(height: 12),
            Row(children: [
              const Icon(Icons.warning_amber_rounded, size: 16, color: AppColors.amber500),
              const SizedBox(width: 4),
              Text('Gap: ${_fmt(goal.fundingGap)}', style: const TextStyle(color: AppColors.amber500, fontSize: 12, fontWeight: FontWeight.bold)),
            ]),
          ],
        ],
      ),
    );
  }
}

