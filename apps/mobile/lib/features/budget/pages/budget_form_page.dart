import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import '../models/budget_models.dart';
import '../repository/budget_repository.dart';
import '../state/budget_state.dart';

class BudgetFormPage extends ConsumerStatefulWidget {
  const BudgetFormPage({super.key});

  @override
  ConsumerState<BudgetFormPage> createState() => _BudgetFormPageState();
}

class _BudgetFormPageState extends ConsumerState<BudgetFormPage> {
  final _formKey = GlobalKey<FormState>();
  String _name = '';
  String _period = 'Monthly';
  late TextEditingController _startDateCtrl;
  late TextEditingController _endDateCtrl;
  final List<_CategoryEntry> _categories = [];
  bool _isSaving = false;

  static const _periods = ['Weekly', 'Monthly', 'Quarterly', 'Yearly', 'Custom'];

  @override
  void initState() {
    super.initState();
    final now = DateTime.now();
    _startDateCtrl = TextEditingController(text: _fmtDate(now));
    _endDateCtrl = TextEditingController(text: _fmtDate(now.add(const Duration(days: 30))));
  }

  @override
  void dispose() {
    _startDateCtrl.dispose();
    _endDateCtrl.dispose();
    super.dispose();
  }

  String _fmtDate(DateTime d) => '${d.year}-${d.month.toString().padLeft(2, '0')}-${d.day.toString().padLeft(2, '0')}';

  Future<void> _pickDate(TextEditingController ctrl) async {
    final parsed = DateTime.tryParse(ctrl.text);
    final picked = await showDatePicker(
      context: context,
      initialDate: parsed ?? DateTime.now(),
      firstDate: DateTime(2020),
      lastDate: DateTime(2035),
    );
    if (picked != null) {
      ctrl.text = _fmtDate(picked);
    }
  }

  void _addCategory() {
    setState(() {
      _categories.add(_CategoryEntry(category: '', amount: 0));
    });
  }

  void _removeCategory(int idx) {
    setState(() => _categories.removeAt(idx));
  }

  void _updateCategory(int idx, String value) {
    _categories[idx].category = value;
  }

  void _updateAmount(int idx, double value) {
    _categories[idx].amount = value;
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    final validCategories = _categories.where((c) => c.category.isNotEmpty && c.amount > 0).toList();
    if (validCategories.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Add at least one category with a name and amount')));
      return;
    }

    setState(() => _isSaving = true);
    try {
      final req = CreateBudgetRequest(
        name: _name,
        period: _period,
        startDate: _startDateCtrl.text,
        endDate: _endDateCtrl.text,
        currency: 'INR',
        categories: validCategories.map((c) => CreateCategoryRequest(
          category: c.category,
          budgetedAmount: c.amount,
          rollover: false,
        )).toList(),
        tags: [],
      );
      await ref.read(budgetRepositoryProvider).createBudget(req);
      ref.invalidate(budgetsProvider);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Budget created')));
        context.pop();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Create Budget')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(AppSpacing.md),
          children: [
            TextFormField(
              decoration: const InputDecoration(labelText: 'Budget Name', border: OutlineInputBorder()),
              validator: (v) => v == null || v.trim().isEmpty ? 'Required' : null,
              onSaved: (v) => _name = v?.trim() ?? '',
            ),
            const SizedBox(height: AppSpacing.md),
            DropdownButtonFormField<String>(
              decoration: const InputDecoration(labelText: 'Period', border: OutlineInputBorder()),
              initialValue: _period,
              items: _periods.map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
              onChanged: (v) => setState(() => _period = v ?? 'Monthly'),
            ),
            const SizedBox(height: AppSpacing.md),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    controller: _startDateCtrl,
                    decoration: const InputDecoration(labelText: 'Start Date', border: OutlineInputBorder(), suffixIcon: Icon(Icons.calendar_today)),
                    readOnly: true,
                    onTap: () => _pickDate(_startDateCtrl),
                    validator: (v) => v == null || v.isEmpty ? 'Required' : null,
                  ),
                ),
                const SizedBox(width: AppSpacing.md),
                Expanded(
                  child: TextFormField(
                    controller: _endDateCtrl,
                    decoration: const InputDecoration(labelText: 'End Date', border: OutlineInputBorder(), suffixIcon: Icon(Icons.calendar_today)),
                    readOnly: true,
                    onTap: () => _pickDate(_endDateCtrl),
                    validator: (v) => v == null || v.isEmpty ? 'Required' : null,
                  ),
                ),
              ],
            ),
            const SizedBox(height: AppSpacing.lg),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text('Categories', style: theme.textTheme.titleMedium),
                TextButton.icon(
                  onPressed: _addCategory,
                  icon: const Icon(Icons.add),
                  label: const Text('Add Category'),
                ),
              ],
            ),
            if (_categories.isEmpty)
              Padding(
                padding: const EdgeInsets.symmetric(vertical: AppSpacing.md),
                child: Text('No categories added yet. Tap "Add Category" to start.', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ),
            ..._categories.asMap().entries.map((entry) {
              final idx = entry.key;
              final cat = entry.value;
              return Card(
                margin: const EdgeInsets.only(bottom: AppSpacing.sm),
                child: Padding(
                  padding: const EdgeInsets.all(AppSpacing.md),
                  child: Row(
                    children: [
                      Expanded(
                        flex: 3,
                        child: TextFormField(
                          initialValue: cat.category,
                          decoration: const InputDecoration(labelText: 'Category', isDense: true, border: OutlineInputBorder()),
                          onChanged: (v) => _updateCategory(idx, v),
                        ),
                      ),
                      const SizedBox(width: AppSpacing.sm),
                      Expanded(
                        flex: 2,
                        child: TextFormField(
                          initialValue: cat.amount > 0 ? cat.amount.toString() : '',
                          decoration: const InputDecoration(labelText: 'Amount', isDense: true, border: OutlineInputBorder()),
                          keyboardType: TextInputType.number,
                          onChanged: (v) => _updateAmount(idx, double.tryParse(v) ?? 0),
                        ),
                      ),
                      const SizedBox(width: 4),
                      IconButton(
                        icon: const Icon(Icons.delete, color: Colors.red, size: 20),
                        onPressed: () => _removeCategory(idx),
                      ),
                    ],
                  ),
                ),
              );
            }),
            const SizedBox(height: AppSpacing.xl),
            SizedBox(
              height: 48,
              child: ElevatedButton(
                onPressed: _isSaving ? null : _submit,
                child: _isSaving ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2)) : const Text('Save Budget'),
              ),
            ),
            const SizedBox(height: AppSpacing.xl),
          ],
        ),
      ),
    );
  }
}

class _CategoryEntry {
  String category;
  double amount;
  _CategoryEntry({required this.category, required this.amount});
}
