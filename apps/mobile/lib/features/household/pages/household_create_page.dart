import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/features/household/state/household_state.dart';

class HouseholdCreatePage extends ConsumerStatefulWidget {
  const HouseholdCreatePage({super.key});

  @override
  ConsumerState<HouseholdCreatePage> createState() => _HouseholdCreatePageState();
}

class _HouseholdCreatePageState extends ConsumerState<HouseholdCreatePage> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();
  String _type = 'Family';
  String _currency = 'INR';
  String _country = 'IN';
  bool _isSaving = false;

  static const _types = ['Single', 'Couple', 'Family', 'Roommates', 'Custom'];
  static const _currencies = ['INR', 'USD', 'EUR', 'GBP', 'CAD', 'AUD'];
  static const _countries = ['IN', 'US', 'GB', 'CA', 'AU', 'DE', 'FR', 'SG', 'AE'];

  @override
  void dispose() {
    _nameCtrl.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isSaving = true);
    try {
      await ref.read(householdRepositoryProvider).createHousehold(
        _nameCtrl.text.trim(), _type, _currency, _country,
      );
      ref.invalidate(householdsProvider);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Household created')));
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
      appBar: AppBar(title: const Text('Create Household')),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(AppSpacing.md),
          children: [
            TextFormField(
              controller: _nameCtrl,
              decoration: const InputDecoration(labelText: 'Household Name', border: OutlineInputBorder()),
              validator: (v) => v == null || v.trim().isEmpty ? 'Required' : null,
            ),
            const SizedBox(height: AppSpacing.md),
            DropdownButtonFormField<String>(
              decoration: const InputDecoration(labelText: 'Type', border: OutlineInputBorder()),
              initialValue: _type,
              items: _types.map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
              onChanged: (v) => setState(() => _type = v ?? 'Family'),
            ),
            const SizedBox(height: AppSpacing.md),
            Row(
              children: [
                Expanded(
                  child: DropdownButtonFormField<String>(
                    decoration: const InputDecoration(labelText: 'Currency', border: OutlineInputBorder()),
                    initialValue: _currency,
                    items: _currencies.map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
                    onChanged: (v) => setState(() => _currency = v ?? 'INR'),
                  ),
                ),
                const SizedBox(width: AppSpacing.md),
                Expanded(
                  child: DropdownButtonFormField<String>(
                    decoration: const InputDecoration(labelText: 'Country', border: OutlineInputBorder()),
                    initialValue: _country,
                    items: _countries.map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
                    onChanged: (v) => setState(() => _country = v ?? 'IN'),
                  ),
                ),
              ],
            ),
            const SizedBox(height: AppSpacing.xl),
            SizedBox(
              height: 48,
              child: ElevatedButton(
                onPressed: _isSaving ? null : _submit,
                child: _isSaving
                    ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                    : const Text('Create Household'),
              ),
            ),
          ],
        ),
      ),
    );
  }
}