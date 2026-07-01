import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/search_models.dart';

class SearchResultCard extends StatelessWidget {
  final SearchResult result;
  final VoidCallback? onTap;

  const SearchResultCard({super.key, required this.result, this.onTap});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap ?? () => _navigate(context),
        child: Padding(
          padding: const EdgeInsets.all(12),
          child: Row(
            children: [
              Container(
                width: 40, height: 40,
                decoration: BoxDecoration(
                  color: _moduleColor(result.module).withOpacity(0.12),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: Icon(_moduleIcon(result.module), color: _moduleColor(result.module), size: 20),
              ),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(result.title, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600), maxLines: 1, overflow: TextOverflow.ellipsis),
                    const SizedBox(height: 2),
                    Row(
                      children: [
                        StatusChip(label: result.module.label, color: _moduleColor(result.module), fontSize: 10),
                        const SizedBox(width: 6),
                        Expanded(
                          child: Text(result.subtitle, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant), maxLines: 1, overflow: TextOverflow.ellipsis),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              const SizedBox(width: 4),
              Icon(Icons.chevron_right, size: 18, color: theme.colorScheme.onSurfaceVariant),
            ],
          ),
        ),
      ),
    );
  }

  void _navigate(BuildContext context) {
    final path = result.routePath;
    if (path.isNotEmpty) context.push(path);
  }

  Color _moduleColor(SearchModule m) {
    switch (m) {
      case SearchModule.accounts: return Colors.indigo;
      case SearchModule.transactions: return Colors.green;
      case SearchModule.goals: return Colors.blue;
      case SearchModule.portfolio: return Colors.teal;
      case SearchModule.timeline: return Colors.brown;
    }
  }

  IconData _moduleIcon(SearchModule m) {
    switch (m) {
      case SearchModule.accounts: return Icons.account_balance;
      case SearchModule.transactions: return Icons.currency_rupee;
      case SearchModule.goals: return Icons.flag;
      case SearchModule.portfolio: return Icons.pie_chart;
      case SearchModule.timeline: return Icons.timeline;
    }
  }
}

class RecentSearchTile extends StatelessWidget {
  final RecentSearch recent;
  final VoidCallback onTap;
  final VoidCallback? onDelete;

  const RecentSearchTile({super.key, required this.recent, required this.onTap, this.onDelete});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ListTile(
      dense: true,
      leading: Icon(Icons.history, size: 20, color: theme.colorScheme.onSurfaceVariant),
      title: Text(recent.query, style: theme.textTheme.bodyMedium),
      trailing: IconButton(
        icon: Icon(Icons.close, size: 16, color: theme.colorScheme.onSurfaceVariant),
        onPressed: onDelete,
      ),
      onTap: onTap,
    );
  }
}

class SearchFilterSheet extends StatefulWidget {
  final SearchFilters current;
  final ValueChanged<SearchFilters> onApply;

  const SearchFilterSheet({super.key, required this.current, required this.onApply});

  @override
  State<SearchFilterSheet> createState() => _SearchFilterSheetState();
}

class _SearchFilterSheetState extends State<SearchFilterSheet> {
  late Set<SearchModule> _selectedModules;
  String? _selectedDateRange;
  String? _selectedCategory;

  final _dateRanges = <String>['', 'Today', 'This Week', 'This Month', 'Last 3 Months', 'This Year'];
  final _categories = <String>['', 'Financial', 'Goal', 'Account', 'Portfolio', 'Health', 'Risk'];

  @override
  void initState() {
    super.initState();
    _selectedModules = Set.from(widget.current.modules);
    _selectedDateRange = widget.current.dateRange;
    _selectedCategory = widget.current.category;
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
      child: DraggableScrollableSheet(
        initialChildSize: 0.6,
        minChildSize: 0.4,
        maxChildSize: 0.8,
        expand: false,
        builder: (_, scrollCtrl) => Padding(
          padding: const EdgeInsets.all(16),
          child: ListView(
            controller: scrollCtrl,
            children: [
              Text('Search Filters', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 16),
              Text('Modules', style: theme.textTheme.titleSmall),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8, runSpacing: 4,
                children: SearchModule.values.map((m) => FilterChip(
                  label: Text(m.label),
                  selected: _selectedModules.contains(m),
                  onSelected: (v) => setState(() { v ? _selectedModules.add(m) : _selectedModules.remove(m); }),
                )).toList(),
              ),
              const SizedBox(height: 16),
              Text('Date Range', style: theme.textTheme.titleSmall),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8,
                children: _dateRanges.map((d) => FilterChip(
                  label: Text(d.isEmpty ? 'All Time' : d),
                  selected: _selectedDateRange == d,
                  onSelected: (_) => setState(() => _selectedDateRange = _selectedDateRange == d ? null : d),
                )).toList(),
              ),
              const SizedBox(height: 16),
              Text('Category', style: theme.textTheme.titleSmall),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8,
                children: _categories.map((c) => FilterChip(
                  label: Text(c.isEmpty ? 'All' : c),
                  selected: _selectedCategory == c,
                  onSelected: (_) => setState(() => _selectedCategory = _selectedCategory == c ? null : c),
                )).toList(),
              ),
              const SizedBox(height: 24),
              FilledButton(
                onPressed: () {
                  widget.onApply(SearchFilters(
                    modules: _selectedModules,
                    dateRange: _selectedDateRange,
                    category: _selectedCategory,
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
}
