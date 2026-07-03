import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';

class CreateCollectionScreen extends StatefulWidget {
  const CreateCollectionScreen({super.key});

  @override
  State<CreateCollectionScreen> createState() => _CreateCollectionScreenState();
}

class _CreateCollectionScreenState extends State<CreateCollectionScreen> {
  final _labelCtrl = TextEditingController();
  final _descCtrl = TextEditingController();
  final _bulkCountCtrl = TextEditingController(text: '1');
  final _bulkAmountCtrl = TextEditingController();
  bool _isSubmitting = false;

  List<_MemberEntry> _members = [];

  @override
  void dispose() {
    _labelCtrl.dispose();
    _descCtrl.dispose();
    _bulkCountCtrl.dispose();
    _bulkAmountCtrl.dispose();
    for (final m in _members) {
      m.nameCtrl.dispose();
      m.amountCtrl.dispose();
    }
    super.dispose();
  }

  void _addBulk() {
    final count = int.tryParse(_bulkCountCtrl.text) ?? 1;
    final amount = double.tryParse(_bulkAmountCtrl.text) ?? 0;
    if (count < 1 || amount <= 0) return;

    setState(() {
      for (int i = 0; i < count; i++) {
        _members.add(_MemberEntry(
          nameCtrl: TextEditingController(text: 'Person ${_members.length + 1}'),
          amountCtrl: TextEditingController(text: amount.toStringAsFixed(0)),
        ));
      }
      _bulkCountCtrl.text = '1';
      _bulkAmountCtrl.clear();
    });
  }

  void _addSingle() {
    setState(() {
      _members.add(_MemberEntry(
        nameCtrl: TextEditingController(text: 'Person ${_members.length + 1}'),
        amountCtrl: TextEditingController(),
      ));
    });
  }

  void _removeMember(int index) {
    setState(() {
      _members[index].nameCtrl.dispose();
      _members[index].amountCtrl.dispose();
      _members.removeAt(index);
    });
  }

  Future<void> _submit() async {
    if (_labelCtrl.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please enter a collection label'), backgroundColor: Colors.red),
      );
      return;
    }
    if (_members.isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Add at least one member'), backgroundColor: Colors.red),
      );
      return;
    }

    setState(() => _isSubmitting = true);
    try {
      final members = _members.map((m) => {
        'name': m.nameCtrl.text.trim(),
        'expected_amount': double.tryParse(m.amountCtrl.text) ?? 0,
      }).where((m) => (m['expected_amount'] as double) > 0).toList();

      if (members.isEmpty) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('All members need a valid amount'), backgroundColor: Colors.red),
        );
        setState(() => _isSubmitting = false);
        return;
      }

      final fp = Provider.of<FinancialProvider>(context, listen: false);
      await fp.createCollection({
        'label': _labelCtrl.text.trim(),
        'description': _descCtrl.text.trim().isEmpty ? null : _descCtrl.text.trim(),
        'members': members,
      });

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('Collection created!'), backgroundColor: Color(0xFF059669)),
        );
        Navigator.of(context).pop();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) setState(() => _isSubmitting = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        title: const Text('Create Collection'),
        actions: [
          TextButton(
            onPressed: _isSubmitting ? null : _submit,
            child: _isSubmitting
                ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2))
                : const Text('Create', style: TextStyle(fontWeight: FontWeight.bold)),
          ),
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          // Label
          TextField(
            controller: _labelCtrl,
            decoration: const InputDecoration(
              labelText: 'Collection Label',
              hintText: 'e.g. Badminton Court - July 2026',
              prefixIcon: Icon(Icons.label_outline),
            ),
          ),
          const SizedBox(height: 16),

          // Description
          TextField(
            controller: _descCtrl,
            maxLines: 2,
            decoration: const InputDecoration(
              labelText: 'Description (optional)',
              hintText: 'What is this collection for?',
              prefixIcon: Icon(Icons.description_outlined),
            ),
          ),
          const SizedBox(height: 28),

          // Bulk Add Section
          Container(
            padding: const EdgeInsets.all(16),
            decoration: BoxDecoration(
              color: const Color(0xFF1E3A8A).withValues(alpha: 0.04),
              borderRadius: BorderRadius.circular(16),
              border: Border.all(color: const Color(0xFF1E3A8A).withValues(alpha: 0.1)),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Icon(Icons.group_add, color: const Color(0xFF1E3A8A), size: 20),
                    const SizedBox(width: 8),
                    const Text('Bulk Add Members',
                        style: TextStyle(fontWeight: FontWeight.w700, fontSize: 15, color: Color(0xFF1E293B))),
                  ],
                ),
                const SizedBox(height: 14),
                Row(
                  children: [
                    Expanded(
                      child: TextField(
                        controller: _bulkCountCtrl,
                        keyboardType: TextInputType.number,
                        decoration: const InputDecoration(
                          labelText: 'Count',
                          isDense: true,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: TextField(
                        controller: _bulkAmountCtrl,
                        keyboardType: TextInputType.number,
                        decoration: const InputDecoration(
                          labelText: 'Amount each',
                          prefixText: '₹ ',
                          isDense: true,
                        ),
                      ),
                    ),
                    const SizedBox(width: 12),
                    ElevatedButton(
                      onPressed: _addBulk,
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
                        backgroundColor: const Color(0xFF1E3A8A),
                      ),
                      child: const Text('Add'),
                    ),
                  ],
                ),
              ],
            ),
          ),
          const SizedBox(height: 20),

          // Members List
          Row(
            children: [
              const Text('Members', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w700, color: Color(0xFF1E293B))),
              const Spacer(),
              Text('${_members.length} person${_members.length != 1 ? 's' : ''}',
                  style: const TextStyle(color: Colors.grey, fontSize: 13)),
            ],
          ),
          const SizedBox(height: 12),

          if (_members.isEmpty)
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: Colors.grey.shade200),
              ),
              child: const Center(
                child: Text('No members yet. Use bulk add above or add one below.',
                    style: TextStyle(color: Colors.grey)),
              ),
            )
          else
            ...List.generate(_members.length, (i) {
              final m = _members[i];
              return Container(
                margin: const EdgeInsets.only(bottom: 8),
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(12),
                  border: Border.all(color: Colors.grey.shade200),
                ),
                child: Row(
                  children: [
                    Expanded(
                      flex: 3,
                      child: TextField(
                        controller: m.nameCtrl,
                        decoration: const InputDecoration(
                          isDense: true,
                          border: InputBorder.none,
                          hintText: 'Name',
                        ),
                      ),
                    ),
                    const Text('₹ ', style: TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF64748B))),
                    SizedBox(
                      width: 80,
                      child: TextField(
                        controller: m.amountCtrl,
                        keyboardType: TextInputType.number,
                        decoration: const InputDecoration(
                          isDense: true,
                          border: InputBorder.none,
                          hintText: 'Amount',
                        ),
                      ),
                    ),
                    IconButton(
                      icon: const Icon(Icons.close, size: 18, color: Colors.red),
                      onPressed: () => _removeMember(i),
                      padding: EdgeInsets.zero,
                      constraints: const BoxConstraints(),
                    ),
                  ],
                ),
              );
            }),

          const SizedBox(height: 12),
          OutlinedButton.icon(
            onPressed: _addSingle,
            icon: const Icon(Icons.person_add_outlined, size: 18),
            label: const Text('Add Single Member'),
            style: OutlinedButton.styleFrom(
              foregroundColor: const Color(0xFF1E3A8A),
              side: const BorderSide(color: Color(0xFF1E3A8A)),
              padding: const EdgeInsets.symmetric(vertical: 14),
            ),
          ),
        ],
      ),
    );
  }
}

class _MemberEntry {
  final TextEditingController nameCtrl;
  final TextEditingController amountCtrl;
  _MemberEntry({required this.nameCtrl, required this.amountCtrl});
}
