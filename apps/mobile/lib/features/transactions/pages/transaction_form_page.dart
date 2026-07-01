import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../models/transaction_models.dart';
import '../repository/transaction_repository.dart';

class TransactionFormPage extends ConsumerStatefulWidget {
  const TransactionFormPage({super.key});
  @override
  ConsumerState<TransactionFormPage> createState() => _TransactionFormPageState();
}

class _TransactionFormPageState extends ConsumerState<TransactionFormPage> {
  final _formKey = GlobalKey<FormState>();
  final _amountCtrl = TextEditingController();
  final _descriptionCtrl = TextEditingController();
  String _type = 'expense';
  String _category = 'purchase';
  bool _isSaving = false;
  bool _isLoading = false;
  String? _editId;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) => _loadIfEditing());
  }

  void _loadIfEditing() async {
    final extra = GoRouterState.of(context).extra;
    if (extra is String) {
      _editId = extra;
      setState(() => _isLoading = true);
      try {
        final repo = ref.read(transactionRepositoryProvider);
        final resp = await repo.getTransaction(id: extra);
        final tx = resp.data;
        if (tx != null) {
          _type = tx.isExpense ? 'expense' : 'income';
          _category = tx.category ?? tx.eventType;
          _amountCtrl.text = tx.amount.abs().toStringAsFixed(0);
          if (tx.description != null) _descriptionCtrl.text = tx.description!;
        }
      } catch (_) {}
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  void dispose() {
    _amountCtrl.dispose();
    _descriptionCtrl.dispose();
    super.dispose();
  }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() => _isSaving = true);
    try {
      final repo = ref.read(transactionRepositoryProvider);
      final amount = double.tryParse(_amountCtrl.text) ?? 0;
      final data = {
        'type': _type,
        'amount': _type == 'expense' ? -amount.abs() : amount.abs(),
        'currency': 'INR',
        if (_descriptionCtrl.text.isNotEmpty) 'description': _descriptionCtrl.text.trim(),
      };
      if (_editId != null) {
        await repo.updateTransaction(id: _editId!, data: data);
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Transaction updated'), behavior: SnackBarBehavior.floating));
          context.pop(true);
        }
      } else {
        await repo.createTransaction(data: data);
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Transaction created'), behavior: SnackBarBehavior.floating));
          context.pop(true);
        }
      }
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Failed: $e'), behavior: SnackBarBehavior.floating));
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    if (_isLoading) return Scaffold(appBar: AppBar(title: const Text('Edit Transaction')), body: const Center(child: CircularProgressIndicator()));
    return Scaffold(
      appBar: AppBar(
        title: Text(_editId != null ? 'Edit Transaction' : 'Add Transaction'),
        actions: [
          TextButton(
            onPressed: _isSaving ? null : _save,
            child: _isSaving
                ? const SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2))
                : const Text('Save'),
          ),
        ],
      ),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            Text('Type', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            SegmentedButton<String>(
              segments: const [
                ButtonSegment(value: 'expense', label: Text('Expense'), icon: Icon(Icons.arrow_upward)),
                ButtonSegment(value: 'income', label: Text('Income'), icon: Icon(Icons.arrow_downward)),
              ],
              selected: {_type},
              onSelectionChanged: (v) => setState(() {
                _type = v.first;
                _category = _type == 'income' ? 'salary' : 'purchase';
              }),
            ),
            const SizedBox(height: 16),
            Text('Amount', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            TextFormField(
              controller: _amountCtrl,
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: InputDecoration(
                prefixText: '₹ ',
                hintText: '0.00',
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
              ),
              validator: (v) {
                if (v == null || v.isEmpty) return 'Amount is required';
                final amount = double.tryParse(v);
                if (amount == null || amount <= 0) return 'Enter a valid amount';
                return null;
              },
            ),
            const SizedBox(height: 16),
            Text('Description', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            TextFormField(
              controller: _descriptionCtrl,
              maxLines: 2,
              decoration: InputDecoration(
                hintText: 'What was this transaction for?',
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 16),
            Text('Category', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            DropdownButtonFormField<String>(
              initialValue: _category,
              items: (_type == 'income'
                  ? ['salary', 'freelance', 'investment', 'refund', 'gift', 'other']
                  : ['purchase', 'bill', 'food', 'transport', 'entertainment', 'utilities', 'rent', 'healthcare', 'other']
              ).map((c) => DropdownMenuItem(value: c, child: Text(c[0].toUpperCase() + c.substring(1)))).toList(),
              onChanged: (v) => setState(() => _category = v ?? 'other'),
              decoration: InputDecoration(
                border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
              ),
            ),
            const SizedBox(height: 24),
            if (_editId == null)
              Text(
                'Note: The backend currently supports limited fields for new transactions.',
                style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant, fontStyle: FontStyle.italic),
              ),
          ],
        ),
      ),
    );
  }
}
