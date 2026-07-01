import 'package:flutter/material.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/timeline_models.dart';

IconData categoryIcon(String category) {
  switch (category) {
    case 'FinancialEvent': return Icons.currency_rupee;
    case 'GoalEvent': return Icons.flag;
    case 'AccountEvent': return Icons.account_balance;
    case 'AssetEvent': return Icons.trending_up;
    case 'LiabilityEvent': return Icons.credit_card;
    case 'PortfolioEvent': return Icons.pie_chart;
    case 'HealthEvent': return Icons.favorite;
    case 'RiskEvent': return Icons.shield;
    case 'RecommendationEvent': return Icons.lightbulb;
    case 'SimulationEvent': return Icons.science;
    case 'OptimizationEvent': return Icons.auto_graph;
    case 'Achievement': return Icons.emoji_events;
    case 'Milestone': return Icons.flag;
    case 'UserEvent': return Icons.person;
    default: return Icons.circle;
  }
}

Color categoryColor(String category) {
  switch (category) {
    case 'FinancialEvent': return Colors.green;
    case 'GoalEvent': return Colors.blue;
    case 'AccountEvent': return Colors.indigo;
    case 'AssetEvent': return Colors.teal;
    case 'LiabilityEvent': return Colors.purple;
    case 'PortfolioEvent': return Colors.blue;
    case 'HealthEvent': return Colors.pink;
    case 'RiskEvent': return Colors.amber;
    case 'RecommendationEvent': return Colors.amber;
    case 'SimulationEvent': return Colors.cyan;
    case 'OptimizationEvent': return Colors.cyan;
    case 'Achievement': return Color(0xFFFFD700);
    case 'Milestone': return Colors.green;
    case 'UserEvent': return Colors.grey;
    default: return Colors.grey;
  }
}

Color severityColor(String severity) {
  switch (severity) {
    case 'critical': return Colors.red;
    case 'warning': return Colors.orange;
    case 'success': return Colors.green;
    case 'milestone': return Colors.blue;
    default: return Colors.grey;
  }
}

class SeverityBadge extends StatelessWidget {
  final String severity;
  const SeverityBadge({super.key, required this.severity});

  @override
  Widget build(BuildContext context) {
    return StatusChip(
      label: severity[0].toUpperCase() + severity.substring(1),
      color: severityColor(severity),
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
    return StatusChip(label: label, color: color, fontSize: 10);
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

    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 40, height: 40,
                decoration: BoxDecoration(color: color.withOpacity(0.12), borderRadius: BorderRadius.circular(10)),
                child: Icon(categoryIcon(item.category), color: color, size: 20),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(child: Text(item.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600), maxLines: 1, overflow: TextOverflow.ellipsis)),
                        const SizedBox(width: 8),
                        Text(item.formattedDate, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
                      ],
                    ),
                    if (item.summary.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Text(item.summary, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 2, overflow: TextOverflow.ellipsis),
                    ],
                    const SizedBox(height: 6),
                    Row(
                      children: [
                        CategoryBadge(category: item.category),
                        const SizedBox(width: 6),
                        SeverityBadge(severity: item.severity),
                        if (item.amount != 0) ...[
                          const Spacer(),
                          Text(formatMoney(item.amount), style: theme.textTheme.bodySmall?.copyWith(fontWeight: FontWeight.w600, color: item.amount > 0 ? Colors.green : Colors.red)),
                        ],
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class TimelineDetailCard extends StatelessWidget {
  final TimelineItem item;
  const TimelineDetailCard({super.key, required this.item});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final color = categoryColor(item.category);
    return AppCard(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(children: [
            Container(
              width: 48, height: 48,
              decoration: BoxDecoration(color: color.withOpacity(0.12), borderRadius: BorderRadius.circular(12)),
              child: Icon(categoryIcon(item.category), color: color, size: 24),
            ),
            const SizedBox(width: 12),
            Expanded(child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(item.title, style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
                Text(item.eventType, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ],
            )),
          ]),
          const SizedBox(height: AppTheme.spacingLg),
          if (item.description.isNotEmpty) ...[
            SectionHeader(title: 'Description'),
            const SizedBox(height: AppTheme.spacingXs),
            Text(item.description, style: theme.textTheme.bodyMedium),
            const SizedBox(height: AppTheme.spacingMd),
          ],
          InfoRow(label: 'Category', value: item.category),
          InfoRow(label: 'Severity', value: item.severity[0].toUpperCase() + item.severity.substring(1)),
          InfoRow(label: 'Event Type', value: item.eventType),
          InfoRow(label: 'Timestamp', value: item.timestamp),
          if (item.relatedEntity != null) InfoRow(label: 'Entity ID', value: item.relatedEntity!),
          if (item.relatedGoal != null) InfoRow(label: 'Related Goal', value: item.relatedGoal!),
          if (item.relatedAccount != null) InfoRow(label: 'Related Account', value: item.relatedAccount!),
          if (item.relatedAsset != null) InfoRow(label: 'Related Asset', value: item.relatedAsset!),
          if (item.relatedAgg != null) InfoRow(label: 'Related Aggregate', value: item.relatedAgg!),
          if (item.amount != 0) InfoRow(label: 'Amount', value: formatMoney(item.amount)),
          if (item.metadata != null && item.metadata!.isNotEmpty) ...[
            const SizedBox(height: AppTheme.spacingSm),
            SectionHeader(title: 'Metadata'),
            const SizedBox(height: AppTheme.spacingXs),
            ...item.metadata!.entries.map((e) => InfoRow(label: e.key, value: '${e.value}')),
          ],
        ],
      ),
    );
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
          padding: const EdgeInsets.all(AppTheme.spacingLg),
          child: ListView(
            controller: scrollCtrl,
            children: [
              Text('Filter Timeline', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: AppTheme.spacingLg),
              SectionHeader(title: 'Categories'),
              const SizedBox(height: AppTheme.spacingSm),
              Wrap(
                spacing: 8, runSpacing: 4,
                children: _categories.map((c) => FilterChip(
                  label: Text(_categoryLabel(c)), selected: _selectedCategories.contains(c),
                  onSelected: (v) => setState(() {
                    v ? _selectedCategories.add(c) : _selectedCategories.remove(c);
                  }),
                )).toList(),
              ),
              const SizedBox(height: AppTheme.spacingLg),
              SectionHeader(title: 'Severity'),
              const SizedBox(height: AppTheme.spacingSm),
              Wrap(
                spacing: 8,
                children: _severities.map((s) => FilterChip(
                  label: Text(s.isEmpty ? 'All' : '${s[0].toUpperCase()}${s.substring(1)}'),
                  selected: _selectedSeverity == s,
                  onSelected: (v) => setState(() => _selectedSeverity = v ? s : ''),
                )).toList(),
              ),
              const SizedBox(height: AppTheme.spacingLg),
              SectionHeader(title: 'Date Range'),
              const SizedBox(height: AppTheme.spacingSm),
              TextField(controller: _startCtrl, decoration: const InputDecoration(labelText: 'Start Date (YYYY-MM-DD)', prefixIcon: Icon(Icons.calendar_today))),
              const SizedBox(height: AppTheme.spacingSm),
              TextField(controller: _endCtrl, decoration: const InputDecoration(labelText: 'End Date (YYYY-MM-DD)', prefixIcon: Icon(Icons.calendar_today))),
              const SizedBox(height: AppTheme.spacingXl),
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
