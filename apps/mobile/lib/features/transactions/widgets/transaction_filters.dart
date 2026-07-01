import 'package:flutter/material.dart';
import '../models/transaction_models.dart';

class TransactionFilterSheet extends StatefulWidget {
  final TransactionFilter current;
  final ValueChanged<TransactionFilter> onApply;

  const TransactionFilterSheet({
    super.key,
    required this.current,
    required this.onApply,
  });

  @override
  State<TransactionFilterSheet> createState() => _TransactionFilterSheetState();
}

class _TransactionFilterSheetState extends State<TransactionFilterSheet> {
  late TransactionFilter _filters;
  final _minAmountCtrl = TextEditingController();
  final _maxAmountCtrl = TextEditingController();

  @override
  void initState() {
    super.initState();
    _filters = widget.current;
    if (_filters.minAmount != null) _minAmountCtrl.text = _filters.minAmount!.toStringAsFixed(0);
    if (_filters.maxAmount != null) _maxAmountCtrl.text = _filters.maxAmount!.toStringAsFixed(0);
  }

  @override
  void dispose() {
    _minAmountCtrl.dispose();
    _maxAmountCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: EdgeInsets.only(bottom: MediaQuery.of(context).viewInsets.bottom),
      child: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(24, 16, 24, 24),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text('Filters', style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
                  TextButton(
                    onPressed: () {
                      setState(() {
                        _filters = const TransactionFilter();
                        _minAmountCtrl.clear();
                        _maxAmountCtrl.clear();
                      });
                    },
                    child: const Text('Clear All'),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              Text('Date Range', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
              const SizedBox(height: 8),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: () => _pickDate(context, isStart: true),
                      icon: const Icon(Icons.calendar_today, size: 16),
                      label: Text(_filters.startDate != null
                          ? '${_filters.startDate!.day}/${_filters.startDate!.month}/${_filters.startDate!.year}'
                          : 'Start date'),
                    ),
                  ),
                  const Padding(
                    padding: EdgeInsets.symmetric(horizontal: 8),
                    child: Text('to'),
                  ),
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: () => _pickDate(context, isStart: false),
                      icon: const Icon(Icons.calendar_today, size: 16),
                      label: Text(_filters.endDate != null
                          ? '${_filters.endDate!.day}/${_filters.endDate!.month}/${_filters.endDate!.year}'
                          : 'End date'),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              Text('Type', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
              const SizedBox(height: 8),
              Wrap(
                spacing: 8,
                children: [
                  _filterChip(theme, 'All', _filters.type == null, () => setState(() => _filters = _filters.copyWith(clearType: true, type: null))),
                  _filterChip(theme, 'Income', _filters.type == 'income', () => setState(() => _filters = _filters.copyWith(type: 'income'))),
                  _filterChip(theme, 'Expense', _filters.type == 'expense', () => setState(() => _filters = _filters.copyWith(type: 'expense'))),
                  _filterChip(theme, 'Transfer', _filters.type == 'transfer', () => setState(() => _filters = _filters.copyWith(type: 'transfer'))),
                ],
              ),
              const SizedBox(height: 16),
              Text('Amount Range', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
              const SizedBox(height: 8),
              Row(
                children: [
                  Expanded(
                    child: TextField(
                      controller: _minAmountCtrl,
                      keyboardType: TextInputType.number,
                      decoration: const InputDecoration(labelText: 'Min', prefixText: '₹ ', isDense: true),
                      onChanged: (v) => setState(() => _filters = _filters.copyWith(minAmount: double.tryParse(v))),
                    ),
                  ),
                  const Padding(padding: EdgeInsets.symmetric(horizontal: 8), child: Text('—')),
                  Expanded(
                    child: TextField(
                      controller: _maxAmountCtrl,
                      keyboardType: TextInputType.number,
                      decoration: const InputDecoration(labelText: 'Max', prefixText: '₹ ', isDense: true),
                      onChanged: (v) => setState(() => _filters = _filters.copyWith(maxAmount: double.tryParse(v))),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                child: FilledButton(
                  onPressed: () {
                    widget.onApply(_filters);
                    Navigator.pop(context);
                  },
                  child: const Text('Apply Filters'),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _filterChip(ThemeData theme, String label, bool selected, VoidCallback onTap) {
    return FilterChip(
      label: Text(label),
      selected: selected,
      onSelected: (_) => onTap(),
      visualDensity: VisualDensity.compact,
    );
  }

  Future<void> _pickDate(BuildContext context, {required bool isStart}) async {
    final now = DateTime.now();
    final picked = await showDatePicker(
      context: context,
      initialDate: isStart
          ? (_filters.startDate ?? now.subtract(const Duration(days: 30)))
          : (_filters.endDate ?? now),
      firstDate: DateTime(2020),
      lastDate: now.add(const Duration(days: 365)),
    );
    if (picked != null) {
      setState(() {
        if (isStart) {
          _filters = _filters.copyWith(startDate: picked);
        } else {
          _filters = _filters.copyWith(endDate: picked);
        }
      });
    }
  }
}
