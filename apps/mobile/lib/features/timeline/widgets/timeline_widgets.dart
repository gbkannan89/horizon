import 'package:flutter/material.dart';
import '../models/timeline_models.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';

IconData categoryIcon(String category) {
  switch (category) {
    case 'FinancialEvent': return Icons.currency_rupee_rounded;
    case 'GoalEvent': return Icons.flag_rounded;
    case 'AccountEvent': return Icons.account_balance_rounded;
    case 'AssetEvent': return Icons.trending_up_rounded;
    case 'LiabilityEvent': return Icons.credit_card_rounded;
    case 'PortfolioEvent': return Icons.pie_chart_rounded;
    case 'HealthEvent': return Icons.favorite_rounded;
    case 'RiskEvent': return Icons.shield_rounded;
    case 'RecommendationEvent': return Icons.lightbulb_rounded;
    case 'SimulationEvent': return Icons.science_rounded;
    case 'OptimizationEvent': return Icons.auto_graph_rounded;
    case 'Achievement': return Icons.emoji_events_rounded;
    case 'Milestone': return Icons.flag_rounded;
    case 'UserEvent': return Icons.person_rounded;
    default: return Icons.circle;
  }
}

Color categoryColor(String category) {
  switch (category) {
    case 'FinancialEvent': return AppColors.teal500;
    case 'GoalEvent': return Colors.blue;
    case 'AccountEvent': return Colors.indigo;
    case 'AssetEvent': return AppColors.teal500;
    case 'LiabilityEvent': return Colors.purple;
    case 'PortfolioEvent': return Colors.blue;
    case 'HealthEvent': return Colors.pink;
    case 'RiskEvent': return AppColors.amber500;
    case 'RecommendationEvent': return AppColors.amber500;
    case 'SimulationEvent': return Colors.cyan;
    case 'OptimizationEvent': return Colors.cyan;
    case 'Achievement': return const Color(0xFFFFD700);
    case 'Milestone': return AppColors.teal500;
    case 'UserEvent': return AppColors.slate500;
    default: return AppColors.slate500;
  }
}

Color severityColor(String severity) {
  switch (severity) {
    case 'critical': return AppColors.red500;
    case 'warning': return AppColors.amber500;
    case 'success': return AppColors.teal500;
    case 'milestone': return Colors.blue;
    default: return AppColors.slate500;
  }
}

class SeverityBadge extends StatelessWidget {
  final String severity;
  const SeverityBadge({super.key, required this.severity});

  @override
  Widget build(BuildContext context) {
    final color = severityColor(severity);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.xs)),
      child: Text(severity[0].toUpperCase() + severity.substring(1), style: TextStyle(color: color, fontSize: 10, fontWeight: FontWeight.bold)),
    );
  }
}

class CategoryBadge extends StatelessWidget {
  final String category;
  const CategoryBadge({super.key, required this.category});

  @override
  Widget build(BuildContext context) {
    final color = categoryColor(category);
    final label = category.replaceAllMapped(RegExp(r'([A-Z])'), (m) => ' ${m.group(1)}').trim();
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.xs)),
      child: Text(label, style: TextStyle(color: color, fontSize: 10, fontWeight: FontWeight.bold)),
    );
  }
}

class TimelineCard extends StatelessWidget {
  final TimelineItem item;
  final VoidCallback? onTap;

  const TimelineCard({super.key, required this.item, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = categoryColor(item.category);

    return GlassCard(
      onTap: onTap,
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: 44, height: 44,
            decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.md)),
            child: Icon(categoryIcon(item.category), color: color, size: 24),
          ),
          const SizedBox(width: AppSpacing.md),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(child: Text(item.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold), maxLines: 1, overflow: TextOverflow.ellipsis)),
                    const SizedBox(width: 8),
                    Text(item.formattedDate, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                  ],
                ),
                if (item.summary.isNotEmpty) ...[
                  const SizedBox(height: 6),
                  Text(item.summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 2, overflow: TextOverflow.ellipsis),
                ],
                const SizedBox(height: 10),
                Row(
                  children: [
                    CategoryBadge(category: item.category),
                    const SizedBox(width: 8),
                    SeverityBadge(severity: item.severity),
                    if (item.amount != 0) ...[
                      const Spacer(),
                      Text(_formatAmount(item.amount), style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.bold, color: item.amount > 0 ? AppColors.teal500 : AppColors.red500)),
                    ],
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  String _formatAmount(int amt) {
    if (amt >= 10000000) return '₹${(amt / 10000000).toStringAsFixed(2)}Cr';
    if (amt >= 100000) return '₹${(amt / 100000).toStringAsFixed(2)}L';
    return '₹$amt';
  }
}

class TimelineDetailCard extends StatelessWidget {
  final TimelineItem item;
  const TimelineDetailCard({super.key, required this.item});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = categoryColor(item.category);
    return GlassCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Container(
              width: 56, height: 56,
              decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.lg)),
              child: Icon(categoryIcon(item.category), color: color, size: 28),
            ),
            const SizedBox(width: AppSpacing.md),
            Expanded(child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(item.title, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                Text(item.eventType, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ],
            )),
          ]),
          const SizedBox(height: 24),
          if (item.description.isNotEmpty) ...[
            Text('Description', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            Text(item.description, style: theme.textTheme.bodyMedium),
            const SizedBox(height: 16),
          ],
          _detailRow(theme, 'Category', item.category),
          _detailRow(theme, 'Severity', item.severity[0].toUpperCase() + item.severity.substring(1)),
          _detailRow(theme, 'Event Type', item.eventType),
          _detailRow(theme, 'Timestamp', item.timestamp),
          if (item.relatedEntity != null) _detailRow(theme, 'Entity ID', item.relatedEntity!),
          if (item.relatedGoal != null) _detailRow(theme, 'Related Goal', item.relatedGoal!),
          if (item.relatedAccount != null) _detailRow(theme, 'Related Account', item.relatedAccount!),
          if (item.relatedAsset != null) _detailRow(theme, 'Related Asset', item.relatedAsset!),
          if (item.relatedAgg != null) _detailRow(theme, 'Related Aggregate', item.relatedAgg!),
          if (item.amount != 0) _detailRow(theme, 'Amount', _formatAmount(item.amount)),
          if (item.metadata != null && item.metadata!.isNotEmpty) ...[
            const SizedBox(height: 12),
            Text('Metadata', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            ...item.metadata!.entries.map((e) => _detailRow(theme, e.key, '${e.value}')),
          ],
        ],
      ),
    );
  }

  Widget _detailRow(ThemeData theme, String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(width: 140, child: Text(label, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant, fontWeight: FontWeight.w500))),
          Expanded(child: Text(value, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600))),
        ],
      ),
    );
  }

  String _formatAmount(int amt) {
    if (amt >= 10000000) return '₹${(amt / 10000000).toStringAsFixed(2)}Cr';
    if (amt >= 100000) return '₹${(amt / 100000).toStringAsFixed(2)}L';
    return '₹$amt';
  }
}

class FilterSheet extends StatefulWidget {
  final FilterParams current;
  final Function(FilterParams) onApply;

  const FilterSheet({super.key, required this.current, required this.onApply});

  @override
  State<FilterSheet> createState() => _FilterSheetState();
}

class _FilterSheetState extends State<FilterSheet> {
  late List<String> _selectedCategories;
  late String _selectedSeverity;
  final _startCtrl = TextEditingController();
  final _endCtrl = TextEditingController();

  final _categories = ['FinancialEvent', 'GoalEvent', 'AccountEvent', 'AssetEvent', 'PortfolioEvent', 'HealthEvent', 'RiskEvent', 'RecommendationEvent'];
  final _severities = ['', 'info', 'warning', 'critical', 'success'];

  @override
  void initState() {
    super.initState();
    _selectedCategories = List.from(widget.current.categories ?? []);
    _selectedSeverity = widget.current.severity ?? '';
  }

  @override
  void dispose() {
    _startCtrl.dispose();
    _endCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
      child: DraggableScrollableSheet(
        initialChildSize: 0.7,
        minChildSize: 0.5,
        maxChildSize: 0.9,
        expand: false,
        builder: (_, scrollCtrl) => Padding(
          padding: const EdgeInsets.all(16),
          child: ListView(
            controller: scrollCtrl,
            children: [
              Text('Filter Timeline', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 16),
              Text('Categories', style: theme.textTheme.titleSmall),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8, runSpacing: 4,
                children: _categories.map((c) => FilterChip(
                  label: Text(_categoryLabel(c)), selected: _selectedCategories.contains(c),
                  onSelected: (v) => setState(() {
                    v ? _selectedCategories.add(c) : _selectedCategories.remove(c);
                  }),
                )).toList(),
              ),
              const SizedBox(height: 16),
              Text('Severity', style: theme.textTheme.titleSmall),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8,
                children: _severities.map((s) => FilterChip(
                  label: Text(s.isEmpty ? 'All' : '${s[0].toUpperCase()}${s.substring(1)}'),
                  selected: _selectedSeverity == s,
                  onSelected: (v) => setState(() => _selectedSeverity = v ? s : ''),
                )).toList(),
              ),
              const SizedBox(height: 16),
              Text('Date Range', style: theme.textTheme.titleSmall),
              const SizedBox(height: 8),
              TextField(controller: _startCtrl, decoration: const InputDecoration(labelText: 'Start Date (YYYY-MM-DD)', prefixIcon: Icon(Icons.calendar_today))),
              const SizedBox(height: 8),
              TextField(controller: _endCtrl, decoration: const InputDecoration(labelText: 'End Date (YYYY-MM-DD)', prefixIcon: Icon(Icons.calendar_today))),
              const SizedBox(height: 24),
              FilledButton(
                onPressed: () {
                  widget.onApply(FilterParams(
                    categories: _selectedCategories.isNotEmpty ? _selectedCategories : null,
                    severity: _selectedSeverity.isNotEmpty ? _selectedSeverity : null,
                    startDate: _startCtrl.text.isNotEmpty ? _startCtrl.text : null,
                    endDate: _endCtrl.text.isNotEmpty ? _endCtrl.text : null,
                  ));
                  Navigator.pop(context);
                },
                style: FilledButton.styleFrom(minimumSize: const Size(double.infinity, 48)),
                child: const Text('Apply Filters'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  String _categoryLabel(String c) => c.replaceAllMapped(RegExp(r'([A-Z])'), (m) => ' ${m.group(1)}').trim();
}
