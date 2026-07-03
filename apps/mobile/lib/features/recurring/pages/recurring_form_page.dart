import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import '../repository/recurring_repository.dart';
import '../state/recurring_state.dart';
import '../widgets/recurring_widgets.dart';

class RecurringFormPage extends ConsumerStatefulWidget {
  const RecurringFormPage({super.key});

  @override
  ConsumerState<RecurringFormPage> createState() => _RecurringFormPageState();
}

class _RecurringFormPageState extends ConsumerState<RecurringFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();
  final _descCtrl = TextEditingController();
  final _amountCtrl = TextEditingController();
  String _frequency = 'Monthly';
  final int _interval = 1;
  String _startDate = '';
  String? _endDate;
  String _eventType = 'expense';
  final String _category = 'general';
  bool _skipWeekends = false;
  bool _isSaving = false;

  @override
  void initState() {
    super.initState();
    final now = DateTime.now();
    _startDate = _fmtDate(now);
  }

  @override
  void dispose() {
    _nameCtrl.dispose();
    _descCtrl.dispose();
    _amountCtrl.dispose();
    super.dispose();
  }

  String _fmtDate(DateTime d) => '${d.year}-${d.month.toString().padLeft(2, '0')}-${d.day.toString().padLeft(2, '0')}';

  Future<void> _pickDate(bool isStart) async {
    final initial = isStart
        ? (DateTime.tryParse(_startDate) ?? DateTime.now())
        : (_endDate != null ? DateTime.tryParse(_endDate!) ?? DateTime.now() : DateTime.now());
    final picked = await showDatePicker(
      context: context,
      initialDate: initial,
      firstDate: DateTime(2020),
      lastDate: DateTime(2035),
    );
    if (picked != null) {
      setState(() {
        if (isStart) {
          _startDate = _fmtDate(picked);
        } else {
          _endDate = _fmtDate(picked);
        }
      });
    }
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isSaving = true);
    try {
      final data = {
        'name': _nameCtrl.text.trim(),
        'description': _descCtrl.text.trim(),
        'amount': (double.tryParse(_amountCtrl.text) ?? 0).toInt(),
        'currency': 'INR',
        'frequency': _frequency,
        'interval': _interval,
        'start_date': _startDate,
        if (_endDate != null) 'end_date': _endDate,
        'event_type': _eventType,
        'category': _category,
        'skip_weekends': _skipWeekends,
      };
      await ref.read(recurringRepositoryProvider).createRecurring(data);
      ref.invalidate(recurringListProvider);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Recurring transaction created')));
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
    return Scaffold(
      appBar: AppBar(title: const Text('Add Recurring')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(AppSpacing.md),
          children: [
            TextFormField(
              controller: _nameCtrl,
              decoration: const InputDecoration(labelText: 'Name', border: OutlineInputBorder()),
              validator: (v) => v == null || v.trim().isEmpty ? 'Required' : null,
            ),
            const SizedBox(height: AppSpacing.md),
            TextFormField(
              controller: _descCtrl,
              decoration: const InputDecoration(labelText: 'Description (optional)', border: OutlineInputBorder()),
              maxLines: 2,
            ),
            const SizedBox(height: AppSpacing.md),
            TextFormField(
              controller: _amountCtrl,
              decoration: const InputDecoration(labelText: 'Amount', border: OutlineInputBorder(), prefixText: '₹ '),
              keyboardType: TextInputType.number,
              validator: (v) => v == null || v.trim().isEmpty ? 'Required' : null,
            ),
            const SizedBox(height: AppSpacing.md),
            FrequencySelector(value: _frequency, onChanged: (v) => setState(() => _frequency = v)),
            const SizedBox(height: AppSpacing.md),
            Row(
              children: [
                Expanded(
                  child: TextFormField(
                    readOnly: true,
                    decoration: const InputDecoration(labelText: 'Start Date', border: OutlineInputBorder(), suffixIcon: Icon(Icons.calendar_today)),
                    controller: TextEditingController(text: _startDate),
                    onTap: () => _pickDate(true),
                    validator: (v) => _startDate.isEmpty ? 'Required' : null,
                  ),
                ),
                const SizedBox(width: AppSpacing.md),
                Expanded(
                  child: TextFormField(
                    readOnly: true,
                    decoration: const InputDecoration(labelText: 'End Date (opt)', border: OutlineInputBorder(), suffixIcon: Icon(Icons.calendar_today)),
                    controller: TextEditingController(text: _endDate ?? ''),
                    onTap: () => _pickDate(false),
                  ),
                ),
              ],
            ),
            const SizedBox(height: AppSpacing.md),
            DropdownButtonFormField<String>(
              decoration: const InputDecoration(labelText: 'Event Type', border: OutlineInputBorder()),
              initialValue: _eventType,
              items: ['income', 'expense', 'transfer'].map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
              onChanged: (v) => setState(() => _eventType = v ?? 'expense'),
            ),
            const SizedBox(height: AppSpacing.md),
            SwitchListTile(
              title: const Text('Skip Weekends'),
              value: _skipWeekends,
              onChanged: (v) => setState(() => _skipWeekends = v),
            ),
            const SizedBox(height: AppSpacing.xl),
            SizedBox(
              height: 48,
              child: ElevatedButton(
                onPressed: _isSaving ? null : _submit,
                child: _isSaving ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2)) : const Text('Save Recurring'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
