import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import '../repository/automation_repository.dart';
import '../state/automation_state.dart';

class AutomationFormPage extends ConsumerStatefulWidget {
  const AutomationFormPage({super.key});

  @override
  ConsumerState<AutomationFormPage> createState() => _AutomationFormPageState();
}

class _AutomationFormPageState extends ConsumerState<AutomationFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();
  final _descCtrl = TextEditingController();
  String _category = 'categorize';
  int _priority = 0;
  bool _isSaving = false;

  @override
  void dispose() {
    _nameCtrl.dispose();
    _descCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isSaving = true);
    try {
      await ref.read(automationRepositoryProvider).createRule({
        'name': _nameCtrl.text.trim(),
        'description': _descCtrl.text.trim(),
        'category': _category,
        'priority': _priority,
        'enabled': true,
        'conditions': [
          {'field': 'description', 'operator': 'contains', 'value': ''},
        ],
        'actions': [
          {'type': _category, 'params': {'category': 'auto_${_category}'}},
        ],
      });
      ref.invalidate(rulesProvider);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Rule created')));
        context.pop();
      }
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Create Rule')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(AppSpacing.md),
          children: [
            TextFormField(
              controller: _nameCtrl,
              decoration: const InputDecoration(labelText: 'Rule Name', border: OutlineInputBorder()),
              validator: (v) => v == null || v.trim().isEmpty ? 'Required' : null,
            ),
            const SizedBox(height: AppSpacing.md),
            TextFormField(
              controller: _descCtrl,
              decoration: const InputDecoration(labelText: 'Description (optional)', border: OutlineInputBorder()),
              maxLines: 2,
            ),
            const SizedBox(height: AppSpacing.md),
            DropdownButtonFormField<String>(
              decoration: const InputDecoration(labelText: 'Action Type', border: OutlineInputBorder()),
              initialValue: _category,
              items: ['categorize', 'tag', 'alert', 'notify', 'skip', 'split']
                  .map((e) => DropdownMenuItem(value: e, child: Text(e)))
                  .toList(),
              onChanged: (v) => setState(() => _category = v ?? 'categorize'),
            ),
            const SizedBox(height: AppSpacing.md),
            TextFormField(
              initialValue: _priority.toString(),
              decoration: const InputDecoration(labelText: 'Priority (lower = higher)', border: OutlineInputBorder()),
              keyboardType: TextInputType.number,
              onChanged: (v) => _priority = int.tryParse(v) ?? 0,
            ),
            const SizedBox(height: AppSpacing.xl),
            SizedBox(
              height: 48,
              child: ElevatedButton(
                onPressed: _isSaving ? null : _submit,
                child: _isSaving
                    ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                    : const Text('Create Rule'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
