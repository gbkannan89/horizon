import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';

class BudgetScreen extends StatelessWidget {
  const BudgetScreen({super.key});

  void _showAddExpenseModal(BuildContext context) {
    final amountCtrl = TextEditingController();
    final nameCtrl = TextEditingController();
    String selectedCategory = 'Food';
    bool isLoading = false;

    String getBucketFor(String cat) {
      if (['Housing', 'Food', 'Bills', 'Transport'].contains(cat)) return 'Needs';
      if (['Shopping', 'Entertainment', 'Dining'].contains(cat)) return 'Wants';
      return 'Savings';
    }

    final categories = [
      {'label': 'Food',          'icon': Icons.restaurant_outlined},
      {'label': 'Housing',       'icon': Icons.home_outlined},
      {'label': 'Bills',         'icon': Icons.receipt_long_outlined},
      {'label': 'Transport',     'icon': Icons.directions_car_outlined},
      {'label': 'Shopping',      'icon': Icons.shopping_bag_outlined},
      {'label': 'Entertainment', 'icon': Icons.movie_outlined},
      {'label': 'Savings',       'icon': Icons.savings_outlined},
    ];

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
          child: Container(
            decoration: const BoxDecoration(
              color: Color(0xFFF8F9FA),
              borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
            ),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Handle bar
                  Center(
                    child: Container(
                      width: 40, height: 4,
                      margin: const EdgeInsets.only(top: 12, bottom: 20),
                      decoration: BoxDecoration(color: Colors.grey.shade300, borderRadius: BorderRadius.circular(4)),
                    ),
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Add Expense', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                      IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(ctx)),
                    ],
                  ),
                  const SizedBox(height: 20),

                  // Amount
                  const Text('Amount', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                  const SizedBox(height: 8),
                  TextField(
                    controller: amountCtrl,
                    keyboardType: const TextInputType.numberWithOptions(decimal: true),
                    autofocus: true,
                    decoration: InputDecoration(
                      prefixText: '₹ ',
                      prefixStyle: const TextStyle(fontWeight: FontWeight.w700, color: Color(0xFF1E3A8A), fontSize: 16),
                      hintText: '500',
                      filled: true,
                      fillColor: Colors.white,
                      enabledBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(14),
                        borderSide: BorderSide(color: Colors.grey.shade200, width: 1.5),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(14),
                        borderSide: const BorderSide(color: Color(0xFF1E3A8A), width: 2),
                      ),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                    ),
                  ),
                  const SizedBox(height: 20),

                  // Description
                  const Text('Description', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                  const SizedBox(height: 8),
                  TextField(
                    controller: nameCtrl,
                    textCapitalization: TextCapitalization.sentences,
                    decoration: InputDecoration(
                      hintText: 'e.g. Zomato dinner, Petrol',
                      filled: true,
                      fillColor: Colors.white,
                      enabledBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(14),
                        borderSide: BorderSide(color: Colors.grey.shade200, width: 1.5),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(14),
                        borderSide: const BorderSide(color: Color(0xFF1E3A8A), width: 2),
                      ),
                      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                    ),
                  ),
                  const SizedBox(height: 20),

                  // Category
                  const Text('Category', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                  const SizedBox(height: 12),
                  Wrap(
                    spacing: 10,
                    runSpacing: 10,
                    children: categories.map((c) {
                      final label = c['label'] as String;
                      final icon = c['icon'] as IconData;
                      final selected = selectedCategory == label;
                      return GestureDetector(
                        onTap: () => setState(() => selectedCategory = label),
                        child: AnimatedContainer(
                          duration: const Duration(milliseconds: 180),
                          padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 10),
                          decoration: BoxDecoration(
                            color: selected ? const Color(0xFF1E3A8A) : Colors.white,
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(
                              color: selected ? const Color(0xFF1E3A8A) : Colors.grey.shade200,
                              width: 1.5,
                            ),
                          ),
                          child: Row(mainAxisSize: MainAxisSize.min, children: [
                            Icon(icon, size: 16, color: selected ? Colors.white : Colors.grey.shade600),
                            const SizedBox(width: 6),
                            Text(label, style: TextStyle(
                              fontSize: 13, fontWeight: FontWeight.w600,
                              color: selected ? Colors.white : Colors.grey.shade700,
                            )),
                          ]),
                        ),
                      );
                    }).toList(),
                  ),
                  const SizedBox(height: 12),
                  // Bucket badge
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                    decoration: BoxDecoration(
                      color: _getBucketColor(getBucketFor(selectedCategory)).withOpacity(0.1),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(
                      'Will be counted as: ${getBucketFor(selectedCategory)}',
                      style: TextStyle(fontSize: 12, fontWeight: FontWeight.w600, color: _getBucketColor(getBucketFor(selectedCategory))),
                    ),
                  ),
                  const SizedBox(height: 32),

                  // Save button
                  SizedBox(
                    width: double.infinity, height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading ? null : () async {
                        final amt = double.tryParse(amountCtrl.text.trim());
                        final name = nameCtrl.text.trim();
                        if (amt == null || amt <= 0 || name.isEmpty) {
                          ScaffoldMessenger.of(ctx).showSnackBar(
                            const SnackBar(content: Text('Please fill in amount and description')),
                          );
                          return;
                        }
                        setState(() => isLoading = true);
                        try {
                          await Provider.of<FinancialProvider>(ctx, listen: false)
                            .addExpense(name, amt, selectedCategory, getBucketFor(selectedCategory));
                          if (ctx.mounted) Navigator.pop(ctx);
                        } catch (e) {
                          setState(() => isLoading = false);
                          if (ctx.mounted) ScaffoldMessenger.of(ctx).showSnackBar(SnackBar(content: Text('Error: $e')));
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1E3A8A),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                        elevation: 0,
                      ),
                      child: isLoading
                        ? const SizedBox(height: 22, width: 22, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                        : const Text('Save Expense', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildCategoryPill(String label, IconData icon, bool isSelected) {
    return Container(
      width: 70,
      padding: const EdgeInsets.symmetric(vertical: 12),
      decoration: BoxDecoration(
        color: isSelected ? const Color(0xFF1E3A8A) : Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: isSelected ? null : Border.all(color: Colors.grey.withOpacity(0.2)),
      ),
      child: Column(
        children: [
          Icon(icon, color: isSelected ? Colors.white : Colors.black87),
          const SizedBox(height: 8),
          Text(label, style: TextStyle(color: isSelected ? Colors.white : Colors.black87, fontSize: 11)),
        ],
      ),
    );
  }

  Widget _buildBudgetBar(String label, String percentText, String amount, double percent, Color color) {
    double safePercent = percent.clamp(0.0, 1.0);
    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  Container(width: 8, height: 8, decoration: BoxDecoration(shape: BoxShape.circle, color: color)),
                  const SizedBox(width: 8),
                  Text(label, style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 14)),
                ],
              ),
              Row(
                children: [
                  Text(percentText, style: TextStyle(color: color, fontWeight: FontWeight.bold, fontSize: 13)),
                  const Text(' — ', style: TextStyle(color: Colors.grey)),
                  Text(amount, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
                ],
              )
            ],
          ),
          const SizedBox(height: 8),
          LinearPercentIndicator(
            lineHeight: 8.0,
            percent: safePercent,
            progressColor: color,
            backgroundColor: Colors.grey.withOpacity(0.2),
            barRadius: const Radius.circular(8),
            padding: EdgeInsets.zero,
          )
        ],
      ),
    );
  }

  Widget _buildExpenseItem(String title, String subtitle, String amount, IconData icon, Color color) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: color.withOpacity(0.1),
                  shape: BoxShape.circle,
                ),
                child: Icon(icon, color: color, size: 20),
              ),
              const SizedBox(width: 16),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(title, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
                  Text(subtitle, style: const TextStyle(color: Colors.grey, fontSize: 12)),
                ],
              )
            ],
          ),
          Text(amount, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
        ],
      ),
    );
  }

  IconData _getIconData(String iconName) {
    switch(iconName) {
      case 'home_outlined': return Icons.home_outlined;
      case 'business_center_outlined': return Icons.business_center_outlined;
      case 'star_border': return Icons.star_border;
      default: return Icons.attach_money;
    }
  }

  Color _getColorForBucket(String bucket) {
    switch(bucket) {
      case 'Needs': return const Color(0xFFE88A1A);
      case 'Wants': return const Color(0xFF6B46C1);
      case 'Savings': return const Color(0xFF059669);
      default: return Colors.blue;
    }
  }

  Color _getBucketColor(String bucket) => _getColorForBucket(bucket);

  @override
  Widget build(BuildContext context) {
    return Consumer<FinancialProvider>(
      builder: (context, provider, child) {
        if (provider.isLoading) {
          return const Center(child: CircularProgressIndicator());
        }

        double totalBudget = provider.needsBudget + provider.wantsBudget + provider.savingsBudget;
        double overallPct = totalBudget > 0 ? provider.totalSpent / totalBudget : 0.0;
        
        return Scaffold(
          backgroundColor: const Color(0xFFF8F9FA),
          appBar: AppBar(
            title: const Text('Budget', style: TextStyle(fontWeight: FontWeight.bold)),
            backgroundColor: const Color(0xFFF8F9FA),
            elevation: 0,
          ),
          body: SingleChildScrollView(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              children: [
                // Month Selector
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Icon(Icons.chevron_left, color: Colors.grey.shade400),
                    const Text('June 2026', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                    Icon(Icons.chevron_right, color: Colors.grey.shade400),
                  ],
                ),
                const SizedBox(height: 24),

                // Budget Breakdown
                Container(
                  padding: const EdgeInsets.all(24),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(24),
                    border: Border.all(color: Colors.grey.withOpacity(0.1)),
                  ),
                  child: Column(
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text('Monthly Budget', style: TextStyle(color: Colors.grey, fontSize: 14)),
                          Text('₹${totalBudget.toStringAsFixed(0)}', style: const TextStyle(fontWeight: FontWeight.w900, fontSize: 18)),
                        ],
                      ),
                      const SizedBox(height: 20),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Row(
                            children: [
                              const Text('Spent ', style: TextStyle(fontWeight: FontWeight.bold)),
                              Text('₹${provider.totalSpent.toStringAsFixed(0)}', style: TextStyle(fontWeight: FontWeight.bold, color: Colors.blue.shade900)),
                            ],
                          ),
                          Text('${(overallPct * 100).toInt()}%', style: const TextStyle(color: Colors.grey, fontSize: 13)),
                        ],
                      ),
                      const SizedBox(height: 8),
                      LinearPercentIndicator(
                        lineHeight: 12.0,
                        percent: overallPct.clamp(0.0, 1.0),
                        linearGradient: const LinearGradient(colors: [Color(0xFFE88A1A), Color(0xFFDC2626)]),
                        backgroundColor: Colors.grey.withOpacity(0.2),
                        barRadius: const Radius.circular(8),
                        padding: EdgeInsets.zero,
                      ),
                      const Padding(
                        padding: EdgeInsets.symmetric(vertical: 24),
                        child: Divider(),
                      ),
                      _buildBudgetBar('Needs', '${(provider.needsSpent / provider.needsBudget * 100).toInt()}%', '₹${provider.needsSpent.toStringAsFixed(0)}', provider.needsSpent / provider.needsBudget, const Color(0xFFE88A1A)),
                      _buildBudgetBar('Wants', '${(provider.wantsSpent / provider.wantsBudget * 100).toInt()}%', '₹${provider.wantsSpent.toStringAsFixed(0)}', provider.wantsSpent / provider.wantsBudget, const Color(0xFF6B46C1)),
                      _buildBudgetBar('Savings', '${(provider.savingsSpent / provider.savingsBudget * 100).toInt()}%', '₹${provider.savingsSpent.toStringAsFixed(0)}', provider.savingsSpent / provider.savingsBudget, const Color(0xFF059669)),
                    ],
                  ),
                ),
                const SizedBox(height: 20),

                // Committed Bills
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: Colors.grey.withOpacity(0.1)),
                  ),
                  child: const Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text('Committed Bills', style: TextStyle(color: Colors.grey, fontSize: 13)),
                          SizedBox(height: 4),
                          Text('₹33,899/month (45% of income)', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                        ],
                      ),
                      Icon(Icons.chevron_right, color: Colors.grey),
                    ],
                  ),
                ),
                const SizedBox(height: 20),

                // Left to spend
                Row(
                  children: [
                    Expanded(
                      child: Container(
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(
                          color: const Color(0xFFFFF7ED),
                          borderRadius: BorderRadius.circular(16),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text('Needs Left', style: TextStyle(color: Colors.grey, fontSize: 12)),
                            const SizedBox(height: 4),
                            Text('₹${(provider.needsBudget - provider.needsSpent).clamp(0, double.infinity).toStringAsFixed(0)}', style: const TextStyle(color: Color(0xFFE88A1A), fontWeight: FontWeight.bold, fontSize: 16)),
                          ],
                        ),
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: Container(
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(
                          color: const Color(0xFFF5F3FF),
                          borderRadius: BorderRadius.circular(16),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text('Wants Left', style: TextStyle(color: Colors.grey, fontSize: 12)),
                            const SizedBox(height: 4),
                            Text('₹${(provider.wantsBudget - provider.wantsSpent).clamp(0, double.infinity).toStringAsFixed(0)}', style: const TextStyle(color: Color(0xFF6B46C1), fontWeight: FontWeight.bold, fontSize: 16)),
                          ],
                        ),
                      ),
                    )
                  ],
                ),
                const SizedBox(height: 32),

                // Recent Expenses
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Recent Expenses', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                    Text('See All', style: TextStyle(color: Colors.blue.shade700, fontSize: 13, fontWeight: FontWeight.bold)),
                  ],
                ),
                const SizedBox(height: 16),
                ...provider.recentExpenses.map((exp) {
                  return _buildExpenseItem(
                    exp['name'] ?? '',
                    exp['date'] ?? 'Today',
                    '₹${(exp['amount'] ?? 0).toStringAsFixed(0)}',
                    _getIconData(exp['icon'] ?? ''),
                    _getColorForBucket(exp['bucket'] ?? '')
                  );
                }).toList(),

                const SizedBox(height: 80),
              ],
            ),
          ),
          floatingActionButton: FloatingActionButton(
            onPressed: () => _showAddExpenseModal(context),
            backgroundColor: const Color(0xFF1E3A8A),
            child: const Icon(Icons.add, color: Colors.white),
          ),
        );
      }
    );
  }
}
