import 'package:flutter/material.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';
import 'package:flutter_animate/flutter_animate.dart';

class GatekeeperCard extends StatelessWidget {
  final VoidCallback onFixCover;
  
  const GatekeeperCard({
    super.key,
    required this.onFixCover,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: AppColors.red500.withOpacity(0.15),
                  shape: BoxShape.circle,
                ),
                child: Icon(Icons.shield_rounded, color: AppColors.red500, size: 20),
              ),
              const SizedBox(width: AppSpacing.sm),
              Text(
                'PROTECT THE EARNER',
                style: theme.textTheme.labelSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: AppColors.red500,
                  letterSpacing: 1.2,
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.md),
          Text(
            'Term Insurance Missing',
            style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: AppSpacing.xs),
          Text(
            'We refuse to cheer wealth growth strategies until your core family cover is secured. Fix this cover first.',
            style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
          ),
          const SizedBox(height: AppSpacing.md),
          OutlinedButton.icon(
            onPressed: onFixCover,
            style: OutlinedButton.styleFrom(
              foregroundColor: AppColors.red500,
              side: BorderSide(color: AppColors.red500.withOpacity(0.5)),
              shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppRadius.md)),
            ),
            icon: const Icon(Icons.add_moderator_rounded, size: 18),
            label: const Text('Add Insurance Details', style: TextStyle(fontWeight: FontWeight.bold)),
          ),
        ],
      ),
    ).animate().fade().slideY(begin: 0.1, end: 0, duration: 400.ms);
  }
}

class WindfallInterceptorCard extends StatelessWidget {
  final double amount;
  final Function(String job) onSelectJob;

  const WindfallInterceptorCard({
    super.key,
    required this.amount,
    required this.onSelectJob,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: AppColors.amber500.withOpacity(0.15),
                  shape: BoxShape.circle,
                ),
                child: Icon(Icons.bolt_rounded, color: AppColors.amber500, size: 20),
              ),
              const SizedBox(width: AppSpacing.sm),
              Text(
                'WINDFALL INTERCEPTOR',
                style: theme.textTheme.labelSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: AppColors.amber500,
                  letterSpacing: 1.2,
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.md),
          Text(
            'Unallocated Surplus Detected',
            style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: AppSpacing.xs),
          Text(
            'You received an extra ₹${amount.toStringAsFixed(0)}. Let\'s give it a job before it dissolves into consumption:',
            style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
          ),
          const SizedBox(height: AppSpacing.md),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              _jobButton(context, 'Top up Emergency Fund', () => onSelectJob('buffer')),
              _jobButton(context, 'Prepay Home Loan (Saves ₹18k)', () => onSelectJob('prepay')),
              _jobButton(context, 'Invest (Split 60/40)', () => onSelectJob('invest')),
            ],
          ),
        ],
      ),
    ).animate().fade().slideY(begin: 0.1, end: 0, duration: 450.ms);
  }

  Widget _jobButton(BuildContext context, String label, VoidCallback onTap) {
    return ActionChip(
      onPressed: onTap,
      backgroundColor: AppColors.teal500.withOpacity(0.15),
      label: Text(label, style: const TextStyle(fontWeight: FontWeight.bold, color: AppColors.teal500, fontSize: 12)),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppRadius.md)),
      side: BorderSide.none,
    );
  }
}

class CalmCoachCard extends StatelessWidget {
  final double overspentAmount;
  final double recoveryAmount;
  final int months;

  const CalmCoachCard({
    super.key,
    required this.overspentAmount,
    required this.recoveryAmount,
    required this.months,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: AppColors.teal500.withOpacity(0.15),
                  shape: BoxShape.circle,
                ),
                child: const Icon(Icons.spa_rounded, color: AppColors.teal500, size: 20),
              ),
              const SizedBox(width: AppSpacing.sm),
              Text(
                'THE CALM COACH',
                style: theme.textTheme.labelSmall?.copyWith(
                  fontWeight: FontWeight.bold,
                  color: AppColors.teal500,
                  letterSpacing: 1.2,
                ),
              ),
            ],
          ),
          const SizedBox(height: AppSpacing.md),
          Text(
            'No alarm, you are okay.',
            style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: AppSpacing.xs),
          Text(
            'You overspent ₹${overspentAmount.toStringAsFixed(0)} on discretionary things. That is what buffers are for. We have silently re-planned your recovery to ₹${recoveryAmount.toStringAsFixed(0)}/month for $months months. Your goals are still on track.',
            style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant),
          ),
        ],
      ),
    ).animate().fade().slideY(begin: 0.1, end: 0, duration: 500.ms);
  }
}
