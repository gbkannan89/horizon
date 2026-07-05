import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:percent_indicator/percent_indicator.dart';
import '../providers/financial_provider.dart';

class CollectionDetailScreen extends StatefulWidget {
  final Map<String, dynamic> collection;
  final double? highlightAmount;
  final String? highlightName;

  const CollectionDetailScreen({
    super.key,
    required this.collection,
    this.highlightAmount,
    this.highlightName,
  });

  @override
  State<CollectionDetailScreen> createState() => _CollectionDetailScreenState();
}

class _CollectionDetailScreenState extends State<CollectionDetailScreen> {
  late Map<String, dynamic> _collection;

  @override
  void initState() {
    super.initState();
    _collection = widget.collection;
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _refresh();
    });
  }

  Future<void> _refresh() async {
    final fp = Provider.of<FinancialProvider>(context, listen: false);
    final detail = await fp.getCollectionDetail(_collection['id']);
    if (detail != null && mounted) {
      setState(() => _collection = detail);
    }
  }

  Future<void> _recordPaymentFromSuggestion(Map<String, dynamic> member) async {
    final amount = widget.highlightAmount ?? (member['expected_amount'] ?? 0).toDouble();
    try {
      final fp = Provider.of<FinancialProvider>(context, listen: false);
      await fp.recordPayment(member['id'], amount);
      await _refresh();
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('₹${amount.toStringAsFixed(0)} linked to ${member['name']}'),
            backgroundColor: const Color(0xFF059669),
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
        );
      }
    }
  }

  Future<void> _recordPayment(Map<String, dynamic> member) async {
    final amtCtrl = TextEditingController(text: (member['expected_amount'] ?? 0).toString());

    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text('Record Payment - ${member['name']}'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Text('Expected: ₹${(member['expected_amount'] ?? 0).toDouble().toStringAsFixed(0)}'),
            const SizedBox(height: 12),
            TextField(
              controller: amtCtrl,
              keyboardType: TextInputType.number,
              decoration: const InputDecoration(
                labelText: 'Amount Received',
                prefixText: '₹ ',
                border: OutlineInputBorder(),
              ),
            ),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.of(ctx).pop(false), child: const Text('Cancel')),
          ElevatedButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            child: const Text('Confirm'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      final amount = double.tryParse(amtCtrl.text) ?? 0;
      if (amount > 0) {
        try {
          final fp = Provider.of<FinancialProvider>(context, listen: false);
          await fp.recordPayment(member['id'], amount);
          await _refresh();
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text('₹${amount.toStringAsFixed(0)} received from ${member['name']}'),
                backgroundColor: const Color(0xFF059669),
              ),
            );
          }
        } catch (e) {
          if (mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
            );
          }
        }
      }
    }
  }

  Future<void> _settleCollection() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Settle Collection?'),
        content: const Text('Mark this collection as settled?'),
        actions: [
          TextButton(onPressed: () => Navigator.of(ctx).pop(false), child: const Text('Cancel')),
          ElevatedButton(onPressed: () => Navigator.of(ctx).pop(true), child: const Text('Settle')),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      try {
        final fp = Provider.of<FinancialProvider>(context, listen: false);
        await fp.updateCollection(_collection['id'], status: 'settled');
        await _refresh();
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('Collection settled!'), backgroundColor: Color(0xFF059669)),
          );
        }
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
          );
        }
      }
    }
  }

  Future<void> _deleteCollection() async {
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Delete Collection?'),
        content: const Text('This will remove all members and payment records.'),
        actions: [
          TextButton(onPressed: () => Navigator.of(ctx).pop(false), child: const Text('Cancel')),
          ElevatedButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Delete'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      try {
        final fp = Provider.of<FinancialProvider>(context, listen: false);
        await fp.deleteCollection(_collection['id']);
        if (mounted) Navigator.of(context).pop();
      } catch (e) {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
          );
        }
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final totalExpected = (_collection['total_expected'] ?? 0).toDouble();
    final totalCollected = (_collection['total_collected'] ?? 0).toDouble();
    final pct = totalExpected > 0 ? (totalCollected / totalExpected).clamp(0.0, 1.0) : 0.0;
    final members = (_collection['members'] as List?) ?? [];
    final memberCount = (_collection['member_count'] ?? members.length) as int;
    final paidCount = (_collection['paid_count'] ?? 0) as int;
    final status = _collection['status'] ?? 'active';
    final label = _collection['label'] ?? '';

    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        title: Text(label, overflow: TextOverflow.ellipsis),
        actions: [
          if (status == 'active')
            PopupMenuButton<String>(
              onSelected: (v) {
                if (v == 'settle') _settleCollection();
                if (v == 'delete') _deleteCollection();
              },
              itemBuilder: (_) => [
                const PopupMenuItem(value: 'settle', child: ListTile(
                  leading: Icon(Icons.check_circle, color: Color(0xFF059669)),
                  title: Text('Settle'),
                  dense: true,
                )),
                const PopupMenuItem(value: 'delete', child: ListTile(
                  leading: Icon(Icons.delete, color: Colors.red),
                  title: Text('Delete'),
                  dense: true,
                )),
              ],
            ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _refresh,
        child: ListView(
          padding: const EdgeInsets.fromLTRB(16, 16, 16, 100),
          children: [
            // Progress card
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                gradient: const LinearGradient(
                  colors: [Color(0xFF0D9488), Color(0xFF2DD4BF)],
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                ),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Column(
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Progress',
                          style: TextStyle(fontSize: 14, color: Colors.white70, fontWeight: FontWeight.w500)),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                        decoration: BoxDecoration(
                          color: status == 'settled' ? Colors.green.shade200.withValues(alpha: 0.3) : Colors.white.withValues(alpha: 0.15),
                          borderRadius: BorderRadius.circular(10),
                        ),
                        child: Text(status.toUpperCase(),
                            style: TextStyle(
                                fontSize: 10, fontWeight: FontWeight.w700,
                                color: status == 'settled' ? Colors.green.shade200 : Colors.white70)),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Text('₹${_fmt(totalCollected)}',
                          style: const TextStyle(fontSize: 32, fontWeight: FontWeight.w900, color: Colors.white)),
                      const SizedBox(width: 8),
                      Padding(
                        padding: const EdgeInsets.only(bottom: 4),
                        child: Text('/ ₹${_fmt(totalExpected)}',
                            style: const TextStyle(fontSize: 16, color: Colors.white60)),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  ClipRRect(
                    borderRadius: BorderRadius.circular(6),
                    child: LinearPercentIndicator(
                      lineHeight: 10,
                      percent: pct,
                      progressColor: pct >= 1.0 ? const Color(0xFF34D399) : Colors.white,
                      backgroundColor: Colors.white.withValues(alpha: 0.2),
                      barRadius: const Radius.circular(6),
                      padding: EdgeInsets.zero,
                    ),
                  ),
                  const SizedBox(height: 12),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceAround,
                    children: [
                      _stat(Icons.people, '$paidCount / $memberCount', 'Paid'),
                      _stat(Icons.currency_rupee, '₹${_fmt(totalCollected)}', 'Collected'),
                      _stat(Icons.pending_actions, '${memberCount - paidCount}', 'Pending'),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),

            // Suggestion from transaction
            if (widget.highlightAmount != null) ...[
              Container(
                margin: const EdgeInsets.only(bottom: 16),
                padding: const EdgeInsets.all(14),
                decoration: BoxDecoration(
                  color: const Color(0xFF6B46C1).withValues(alpha: 0.06),
                  borderRadius: BorderRadius.circular(14),
                  border: Border.all(color: const Color(0xFF6B46C1).withValues(alpha: 0.15)),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        const Icon(Icons.lightbulb_outline, color: Color(0xFF6B46C1), size: 18),
                        const SizedBox(width: 8),
                        const Text('Link this transaction',
                            style: TextStyle(fontWeight: FontWeight.w700, fontSize: 14, color: Color(0xFF6B46C1))),
                      ],
                    ),
                    const SizedBox(height: 10),
                    Text(
                      '₹${widget.highlightAmount!.toStringAsFixed(0)} — ${widget.highlightName ?? ''}',
                      style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w600),
                    ),
                    const SizedBox(height: 4),
                    const Text('Which member does this payment belong to?',
                        style: TextStyle(fontSize: 12, color: Colors.grey)),
                    const SizedBox(height: 12),
                    ...members.where((m) => m['status'] != 'paid').take(3).map((m) {
                      final diff = (widget.highlightAmount! - (m['expected_amount'] ?? 0).toDouble()).abs();
                      final isMatch = diff < 1;
                      return Padding(
                        padding: const EdgeInsets.only(bottom: 6),
                        child: OutlinedButton.icon(
                          onPressed: () => _recordPaymentFromSuggestion(m),
                          icon: Icon(
                            isMatch ? Icons.recommend : Icons.person_outline,
                            size: 16,
                            color: isMatch ? const Color(0xFF0D9488) : null,
                          ),
                          label: Text(
                            '${m['name']} — ₹${(m['expected_amount'] ?? 0).toDouble().toStringAsFixed(0)}'
                            '${isMatch ? '  ✓ Best match' : ''}',
                            style: TextStyle(
                              fontSize: 13,
                              color: isMatch ? const Color(0xFF0D9488) : null,
                              fontWeight: isMatch ? FontWeight.w700 : null,
                            ),
                          ),
                          style: OutlinedButton.styleFrom(
                            foregroundColor: isMatch ? const Color(0xFF0D9488) : const Color(0xFF64748B),
                            side: BorderSide(
                              color: isMatch
                                  ? const Color(0xFF0D9488)
                                  : const Color(0xFFE2E8F0),
                              width: isMatch ? 2 : 1,
                            ),
                            backgroundColor: isMatch
                                ? const Color(0xFF0D9488).withValues(alpha: 0.05)
                                : null,
                            padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
                          ),
                        ),
                      );
                    }),
                    if (members.where((m) => m['status'] != 'paid').length > 3)
                      Text('+ ${members.where((m) => m['status'] != 'paid').length - 3} more pending members',
                          style: const TextStyle(fontSize: 11, color: Colors.grey)),
                  ],
                ),
              ),
            ],

            // Members
            Row(
              children: [
                const Text('Members', style: TextStyle(fontSize: 16, fontWeight: FontWeight.w700, color: Color(0xFF1E293B))),
                const Spacer(),
                Text('${members.length} person${members.length != 1 ? 's' : ''}',
                    style: const TextStyle(color: Colors.grey, fontSize: 13)),
              ],
            ),
            const SizedBox(height: 12),

            if (members.isEmpty)
              Container(
                padding: const EdgeInsets.all(24),
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(16),
                ),
                child: const Center(child: Text('No members', style: TextStyle(color: Colors.grey))),
              )
            else
              ...members.map((m) => _memberCard(context, m)),
          ],
        ),
      ),
    );
  }

  Widget _stat(IconData icon, String value, String label) {
    return Column(
      children: [
        Icon(icon, color: Colors.white70, size: 18),
        const SizedBox(height: 4),
        Text(value, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800, color: Colors.white)),
        Text(label, style: const TextStyle(fontSize: 10, color: Colors.white60)),
      ],
    );
  }

  Widget _memberCard(BuildContext context, Map<String, dynamic> m) {
    final name = m['name'] ?? '';
    final expected = (m['expected_amount'] ?? 0).toDouble();
    final paid = (m['paid_amount'] ?? 0).toDouble();
    final status = m['status'] ?? 'pending';
    final isPaid = status == 'paid';
    final isPartial = status == 'partial' || (paid > 0 && paid < expected);

    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(
          color: isPaid ? const Color(0xFF059669).withValues(alpha: 0.2) : Colors.grey.shade200,
        ),
      ),
      child: Row(
        children: [
          Container(
            width: 42, height: 42,
            decoration: BoxDecoration(
              color: isPaid
                  ? const Color(0xFF059669).withValues(alpha: 0.1)
                  : isPartial
                      ? const Color(0xFFF59E0B).withValues(alpha: 0.1)
                      : const Color(0xFFE2E8F0),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Icon(
              isPaid ? Icons.check_circle : isPartial ? Icons.hourglass_bottom : Icons.person_outline,
              color: isPaid ? const Color(0xFF059669) : isPartial ? const Color(0xFFF59E0B) : Colors.grey,
              size: 22,
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(name,
                          style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 14, color: Color(0xFF1E293B)),
                          overflow: TextOverflow.ellipsis),
                    ),
                    Text('₹${_fmt(paid)} / ₹${_fmt(expected)}',
                        style: TextStyle(
                            fontSize: 13, fontWeight: FontWeight.w700,
                            color: isPaid ? const Color(0xFF059669) : const Color(0xFF1E293B))),
                  ],
                ),
                const SizedBox(height: 4),
                Row(
                  children: [
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: isPaid
                            ? const Color(0xFF059669).withValues(alpha: 0.08)
                            : isPartial
                                ? const Color(0xFFF59E0B).withValues(alpha: 0.08)
                                : const Color(0xFFE2E8F0),
                        borderRadius: BorderRadius.circular(6),
                      ),
                      child: Text(
                        isPaid ? 'PAID' : isPartial ? 'PARTIAL' : 'PENDING',
                        style: TextStyle(
                          fontSize: 9, fontWeight: FontWeight.w700,
                          color: isPaid ? const Color(0xFF059669) : isPartial ? const Color(0xFFF59E0B) : Colors.grey,
                        ),
                      ),
                    ),
                    const Spacer(),
                    if (m['paid_date'] != null)
                      Text('${m['paid_date']}',
                          style: const TextStyle(fontSize: 11, color: Colors.grey)),
                  ],
                ),
                if (isPartial)
                  Padding(
                    padding: const EdgeInsets.only(top: 6),
                    child: ClipRRect(
                      borderRadius: BorderRadius.circular(3),
                      child: LinearProgressIndicator(
                        value: (paid / expected).clamp(0.0, 1.0),
                        backgroundColor: const Color(0xFFF59E0B).withValues(alpha: 0.1),
                        valueColor: const AlwaysStoppedAnimation<Color>(Color(0xFFF59E0B)),
                        minHeight: 3,
                      ),
                    ),
                  ),
              ],
            ),
          ),
          if (!isPaid)
            GestureDetector(
              onTap: () => _recordPayment(m),
              child: Container(
                margin: const EdgeInsets.only(left: 8),
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: const Color(0xFF0D9488).withValues(alpha: 0.08),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: const Icon(Icons.payments_outlined, color: Color(0xFF0D9488), size: 20),
              ),
            ),
        ],
      ),
    );
  }

  String _fmt(dynamic val) {
    if (val == null) return '0';
    final n = (val is num) ? val : double.tryParse(val.toString()) ?? 0;
    if (n >= 100000) return '${(n / 100000).toStringAsFixed(1)}L';
    if (n >= 1000) return '${(n / 1000).toStringAsFixed(1)}K';
    return n.toStringAsFixed(0);
  }
}
