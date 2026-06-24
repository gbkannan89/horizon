import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:provider/provider.dart';
import 'package:file_picker/file_picker.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class BudgetScreen extends StatelessWidget {
  const BudgetScreen({super.key});

  String _getMonthName(int month, int year) {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return '${months[month - 1]} $year';
  }

  void _showUploadProgressDialog(BuildContext context, String filePath) async {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (ctx) => const UploadAnimationDialog(),
    );
    try {
      await Provider.of<FinancialProvider>(context, listen: false).uploadStatement(filePath);
      if (context.mounted) {
        Navigator.of(context).pop(); // close dialog
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Upload successful!')));
      }
    } catch (e) {
      if (context.mounted) {
        Navigator.of(context).pop();
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Upload failed: $e')));
      }
    }
  }

  void _showAddExpenseModal(BuildContext context) {
    final amountCtrl = TextEditingController();
    final nameCtrl = TextEditingController();
    String selectedCategory = 'Food';
    String? manuallySelectedBucket;
    bool isRecurring = false;
    bool isEmi = false;
    final emiMonthsCtrl = TextEditingController();
    String recurringFreq = 'monthly';
    DateTime startDate = DateTime.now();
    bool isLoading = false;

    String getBucketFor(String cat) {
      if (['Housing', 'Food', 'Bills', 'Transport'].contains(cat)) return 'Needs';
      if (['Shopping', 'Entertainment', 'Dining'].contains(cat)) return 'Wants';
      return 'Savings';
    }

    String getCurrentBucket() => manuallySelectedBucket ?? getBucketFor(selectedCategory);

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
                        onTap: () {
                          setState(() {
                            selectedCategory = label;
                            manuallySelectedBucket = null; // reset to default for this category
                          });
                        },
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
                  const SizedBox(height: 20),

                  // Bucket selection
                  const Text('Bucket (Needs / Wants / Savings)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                  const SizedBox(height: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(14),
                      border: Border.all(color: Colors.grey.shade200, width: 1.5),
                    ),
                    child: DropdownButtonHideUnderline(
                      child: DropdownButton<String>(
                        value: getCurrentBucket(),
                        isExpanded: true,
                        icon: const Icon(Icons.arrow_drop_down, color: Colors.grey),
                        items: ['Needs', 'Wants', 'Savings'].map((b) {
                          return DropdownMenuItem(
                            value: b,
                            child: Row(
                              children: [
                                Container(
                                  width: 12, height: 12,
                                  decoration: BoxDecoration(shape: BoxShape.circle, color: _getBucketColor(b)),
                                ),
                                const SizedBox(width: 10),
                                Text(b, style: const TextStyle(fontWeight: FontWeight.w600)),
                              ],
                            ),
                          );
                        }).toList(),
                        onChanged: (val) {
                          if (val != null) {
                            setState(() => manuallySelectedBucket = val);
                          }
                        },
                      ),
                    ),
                  ),
                  const SizedBox(height: 20),

                  // Recurring Switch
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(14),
                      border: Border.all(color: Colors.grey.shade200, width: 1.5),
                    ),
                    child: Column(
                      children: [
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          children: [
                            Row(
                              children: [
                                Icon(Icons.autorenew_rounded, color: isRecurring ? const Color(0xFF1E3A8A) : Colors.grey, size: 20),
                                const SizedBox(width: 10),
                                Text('Make this a Recurring Bill', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 14, color: isRecurring ? const Color(0xFF1E293B) : Colors.grey.shade600)),
                              ],
                            ),
                            Switch(
                              value: isRecurring,
                              onChanged: (val) => setState(() => isRecurring = val),
                              activeColor: const Color(0xFF1E3A8A),
                            ),
                          ],
                        ),
                        if (isRecurring) ...[
                          const Divider(),
                          // Frequency
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              const Text('Frequency', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                              DropdownButton<String>(
                                value: recurringFreq,
                                underline: const SizedBox(),
                                items: ['monthly', 'quarterly', 'yearly'].map((f) => DropdownMenuItem(value: f, child: Text(f.toUpperCase(), style: const TextStyle(fontSize: 13)))).toList(),
                                onChanged: (v) => setState(() => recurringFreq = v!),
                              ),
                            ],
                          ),
                          // Start Date
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              const Text('Start Date', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                              TextButton(
                                onPressed: () async {
                                  final d = await showDatePicker(
                                    context: ctx,
                                    initialDate: startDate,
                                    firstDate: DateTime(2020),
                                    lastDate: DateTime(2100),
                                  );
                                  if (d != null) setState(() => startDate = d);
                                },
                                child: Text('${startDate.day}/${startDate.month}/${startDate.year}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                              )
                            ],
                          ),
                          // EMI Toggle
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              const Text('Is this an EMI?', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                              Switch(
                                value: isEmi,
                                onChanged: (val) => setState(() => isEmi = val),
                                activeColor: const Color(0xFF1E3A8A),
                              ),
                            ],
                          ),
                          if (isEmi) ...[
                            const SizedBox(height: 8),
                            TextField(
                              controller: emiMonthsCtrl,
                              keyboardType: TextInputType.number,
                              decoration: InputDecoration(
                                hintText: 'Total EMI Months (e.g. 6)',
                                filled: true,
                                fillColor: Colors.grey.shade50,
                                contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                                border: OutlineInputBorder(
                                  borderRadius: BorderRadius.circular(10),
                                  borderSide: BorderSide.none,
                                ),
                              ),
                            ),
                          ]
                        ]
                      ],
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
                          if (isRecurring) {
                            int? totalEmis;
                            if (isEmi) {
                              totalEmis = int.tryParse(emiMonthsCtrl.text);
                              if (totalEmis == null || totalEmis <= 0) {
                                ScaffoldMessenger.of(ctx).showSnackBar(const SnackBar(content: Text('Please enter valid total EMI months')));
                                setState(() => isLoading = false);
                                return;
                              }
                            }
                            await Provider.of<FinancialProvider>(ctx, listen: false)
                              .addRecurringBill(name, amt, selectedCategory, getCurrentBucket(), frequency: recurringFreq, isEmi: isEmi, emiTotalMonths: totalEmis, startDate: startDate);
                          } else {
                            await Provider.of<FinancialProvider>(ctx, listen: false)
                              .addExpense(name, amt, selectedCategory, getCurrentBucket());
                          }
                          if (ctx.mounted) {
                            Navigator.pop(ctx);
                            UiUtils.showSnack(ctx, 'Expense saved successfully!');
                          }
                        } catch (e) {
                          setState(() => isLoading = false);
                          if (ctx.mounted) UiUtils.showSnack(ctx, 'Error: $e', isError: true);
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

  Widget _formatBadge(String label) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(color: Colors.white.withOpacity(0.2), borderRadius: BorderRadius.circular(6)),
      child: Text(label, style: const TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.w600)),
    );
  }

  // ── CATEGORY BREAKDOWN ──────────────────────────────────────────────────────
  Widget _buildCategoryBreakdown(Map<String, dynamic>? bd) {
    if (bd == null) return const SizedBox.shrink();
    final categories = bd['categoryBreakdown'] as List<dynamic>? ?? [];
    if (categories.isEmpty) return const SizedBox.shrink();

    final colors = [
      const Color(0xFF3B82F6), const Color(0xFFF59E0B), const Color(0xFF10B981),
      const Color(0xFFEF4444), const Color(0xFF8B5CF6), const Color(0xFFEC4899),
      const Color(0xFF06B6D4), const Color(0xFFF97316),
    ];

    return _budgetCard(
      title: 'Spending by Category',
      icon: Icons.pie_chart_rounded,
      color: const Color(0xFF3B82F6),
      child: Column(children: [
        ...categories.take(6).toList().asMap().entries.map((entry) {
          final i = entry.key;
          final cat = entry.value;
          final name = cat['category'] ?? '';
          final amount = (cat['amount'] ?? 0).toDouble();
          final color = colors[i % colors.length];
          return Padding(
            padding: const EdgeInsets.only(bottom: 10),
            child: Row(children: [
              Container(width: 10, height: 10, decoration: BoxDecoration(color: color, borderRadius: BorderRadius.circular(3))),
              const SizedBox(width: 10),
              Expanded(child: Text(name, style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w600))),
              Text('₹${_fmt(amount)}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Color(0xFF1E293B))),
            ]),
          );
        }),
      ]),
    );
  }

  // ── MONTH-OVER-MONTH COMPARISON ─────────────────────────────────────────────
  Widget _buildBucketComparison(Map<String, dynamic>? bd) {
    if (bd == null) return const SizedBox.shrink();
    final buckets = bd['bucketComparison'] as List<dynamic>? ?? [];
    if (buckets.isEmpty) return const SizedBox.shrink();

    final prevMonth = bd['previousMonth'] ?? 'Last';
    final colors = {
      'Needs': const Color(0xFFE88A1A),
      'Wants': const Color(0xFF6B46C1),
      'Savings': const Color(0xFF059669),
    };

    return _budgetCard(
      title: 'vs $prevMonth',
      icon: Icons.compare_arrows_rounded,
      color: const Color(0xFF6B46C1),
      child: Column(children: buckets.map<Widget>((b) {
        final bucket = b['bucket'] ?? '';
        final curr = (b['currentAmount'] ?? 0).toDouble();
        final change = (b['changePct'] ?? 0).toDouble();
        final color = colors[bucket] ?? Colors.grey;
        final isUp = change > 0;
        final isDown = change < 0;
        final isWants = bucket == 'Wants';

        return Padding(
          padding: const EdgeInsets.only(bottom: 12),
          child: Row(children: [
            Container(width: 10, height: 10, decoration: BoxDecoration(color: color, shape: BoxShape.circle)),
            const SizedBox(width: 10),
            Expanded(child: Text(bucket, style: TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: color))),
            Text('₹${_fmt(curr)}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
            const SizedBox(width: 8),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
              decoration: BoxDecoration(
                color: (isUp ? (isWants ? Colors.red : Colors.green) : isDown ? (isWants ? Colors.green : Colors.red) : Colors.grey).withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
              ),
              child: Row(mainAxisSize: MainAxisSize.min, children: [
                Icon(
                  change == 0 ? Icons.remove_rounded : (isUp ? Icons.arrow_upward_rounded : Icons.arrow_downward_rounded),
                  size: 12,
                  color: isUp ? (isWants ? Colors.red : Colors.green) : isDown ? (isWants ? Colors.green : Colors.red) : Colors.grey,
                ),
                const SizedBox(width: 2),
                Text('${change.abs().toStringAsFixed(1)}%', style: TextStyle(
                  fontSize: 11, fontWeight: FontWeight.bold,
                  color: isUp ? (isWants ? Colors.red : Colors.green) : isDown ? (isWants ? Colors.green : Colors.red) : Colors.grey,
                )),
              ]),
            ),
          ]),
        );
      }).toList()),
    );
  }

  Widget _budgetCard({required String title, required IconData icon, required Color color, required Widget child}) {
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.03), blurRadius: 16, offset: const Offset(0, 6))],
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [
          Container(padding: const EdgeInsets.all(8), decoration: BoxDecoration(color: color.withOpacity(0.1), borderRadius: BorderRadius.circular(10)), child: Icon(icon, color: color, size: 18)),
          const SizedBox(width: 10),
          Text(title, style: const TextStyle(fontSize: 15, fontWeight: FontWeight.w800, color: Color(0xFF1E293B))),
        ]),
        const SizedBox(height: 16),
        child,
      ]),
    );
  }

  String _fmt(double v) => v >= 1000 ? '${(v / 1000).toStringAsFixed(0)}K' : v.toStringAsFixed(0);

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

  String _formatDate(dynamic dateVal) {
    if (dateVal == null) return '';
    try {
      final dt = DateTime.parse(dateVal.toString());
      final now = DateTime.now();
      final today = DateTime(now.year, now.month, now.day);
      final diff = today.difference(DateTime(dt.year, dt.month, dt.day)).inDays;
      if (diff == 0) return 'Today';
      if (diff == 1) return 'Yesterday';
      const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
      if (dt.year == now.year) return '${dt.day} ${months[dt.month - 1]}';
      return '${dt.day} ${months[dt.month - 1]} ${dt.year}';
    } catch (_) {
      return dateVal.toString();
    }
  }

  Widget _buildExpenseItem(dynamic exp) {
    final name = exp['name'] ?? '';
    final date = _formatDate(exp['date']);
    final amount = (exp['amount'] ?? 0).toDouble();
    final bucket = exp['bucket'] ?? 'Wants';
    final color = _getColorForBucket(bucket);
    final icon = _getIconData(exp['icon'] ?? '');

    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: Colors.grey.withOpacity(0.08)),
        boxShadow: [BoxShadow(color: Colors.black.withOpacity(0.03), blurRadius: 8, offset: const Offset(0, 2))],
      ),
      child: Row(
        children: [
          Container(
            width: 44, height: 44,
            decoration: BoxDecoration(color: color.withOpacity(0.12), shape: BoxShape.circle),
            child: Icon(icon, color: color, size: 20),
          ),
          const SizedBox(width: 14),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(name, style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 14), maxLines: 1, overflow: TextOverflow.ellipsis),
                const SizedBox(height: 4),
                Row(
                  children: [
                    Text(date, style: TextStyle(color: Colors.grey.shade500, fontSize: 12)),
                    const SizedBox(width: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                      decoration: BoxDecoration(color: color.withOpacity(0.1), borderRadius: BorderRadius.circular(10)),
                      child: Text(bucket, style: TextStyle(color: color, fontSize: 10, fontWeight: FontWeight.w700)),
                    ),
                  ],
                ),
              ],
            ),
          ),
          Text('₹${amount.toStringAsFixed(0)}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15, color: Color(0xFF1E293B))),
        ],
      ),
    );
  }

  IconData _getIconData(String iconName) {
    switch(iconName) {
      case 'movie_outlined': return Icons.movie_outlined;
      case 'shopping_bag_outlined': return Icons.shopping_bag_outlined;
      case 'restaurant_outlined': return Icons.restaurant_outlined;
      case 'directions_car_outlined': return Icons.directions_car_outlined;
      case 'receipt_long_outlined': return Icons.receipt_long_outlined;
      case 'savings_outlined': return Icons.savings_outlined;
      case 'home_outlined': return Icons.home_outlined;
      case 'business_center_outlined': return Icons.business_center_outlined;
      case 'star_border': return Icons.star_border;
      default: return Icons.receipt;
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


        double totalBudget = provider.needsBudget + provider.wantsBudget + provider.savingsBudget;
        double overallPct = totalBudget > 0 ? provider.totalSpent / totalBudget : 0.0;
        
        double totalBills = provider.bills.fold(0, (s, b) => s + b.amount);
        int billPct = provider.totalIncomeAgg > 0 ? (totalBills / provider.totalIncomeAgg * 100).toInt() : 0;
        
        return Scaffold(
          backgroundColor: const Color(0xFFF8F9FA),
          body: SafeArea(
            child: Column(
              children: [
                if (provider.isLoading)
                  LinearProgressIndicator(
                    backgroundColor: Colors.transparent,
                    valueColor: AlwaysStoppedAnimation<Color>(Colors.blue.shade700),
                    minHeight: 4,
                  ),
                Expanded(
                  child: SingleChildScrollView(
            padding: const EdgeInsets.all(20.0),
            child: Column(
              children: [
                // Month Selector
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    IconButton(
                      icon: Icon(Icons.chevron_left, color: Colors.blue.shade700),
                      onPressed: () => provider.previousMonth(),
                    ),
                    DropdownButtonHideUnderline(
                      child: DropdownButton<String>(
                        value: '${provider.selectedMonth}-${provider.selectedYear}',
                        icon: const Icon(Icons.arrow_drop_down, color: Colors.blue),
                        style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Colors.black87),
                        items: (() {
                          List<DropdownMenuItem<String>> items = [];
                          bool foundSelected = false;
                          String selectedVal = '${provider.selectedMonth}-${provider.selectedYear}';

                          for (int index = -6; index < 24; index++) {
                            int y = DateTime.now().year - (index ~/ 12);
                            int m = DateTime.now().month - (index % 12);
                            
                            // Handling negative modulo edge cases safely
                            while (m <= 0) { m += 12; y -= 1; }
                            while (m > 12) { m -= 12; y += 1; }

                            String val = '$m-$y';
                            if (val == selectedVal) foundSelected = true;
                            
                            // Prevent duplicates
                            if (!items.any((item) => item.value == val)) {
                              items.add(DropdownMenuItem(value: val, child: Text(_getMonthName(m, y))));
                            }
                          }
                          
                          if (!foundSelected) {
                            items.insert(0, DropdownMenuItem(
                              value: selectedVal,
                              child: Text(_getMonthName(provider.selectedMonth, provider.selectedYear)),
                            ));
                          }
                          
                          return items;
                        })(),
                        onChanged: (val) {
                          if (val != null) {
                            final parts = val.split('-');
                            provider.selectedMonth = int.parse(parts[0]);
                            provider.selectedYear = int.parse(parts[1]);
                            provider.loadAllData();
                          }
                        },
                      ),
                    ),
                    IconButton(
                      icon: Icon(Icons.chevron_right, color: Colors.blue.shade700),
                      onPressed: () => provider.nextMonth(),
                    ),
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
                      _buildBudgetBar('Needs', '${provider.needsBudget > 0 ? (provider.needsSpent / provider.needsBudget * 100).toInt() : 0}%', '₹${provider.needsSpent.toStringAsFixed(0)}', provider.needsBudget > 0 ? provider.needsSpent / provider.needsBudget : 0, const Color(0xFFE88A1A)),
                      _buildBudgetBar('Wants', '${provider.wantsBudget > 0 ? (provider.wantsSpent / provider.wantsBudget * 100).toInt() : 0}%', '₹${provider.wantsSpent.toStringAsFixed(0)}', provider.wantsBudget > 0 ? provider.wantsSpent / provider.wantsBudget : 0, const Color(0xFF6B46C1)),
                      _buildBudgetBar('Savings', '${provider.savingsBudget > 0 ? (provider.savingsSpent / provider.savingsBudget * 100).toInt() : 0}%', '₹${provider.savingsSpent.toStringAsFixed(0)}', provider.savingsBudget > 0 ? provider.savingsSpent / provider.savingsBudget : 0, const Color(0xFF059669)),
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
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text('Committed Bills', style: TextStyle(color: Colors.grey, fontSize: 13)),
                          const SizedBox(height: 4),
                          Text('₹${totalBills.toStringAsFixed(0)}/month ($billPct% of income)', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                        ],
                      ),
                      const Icon(Icons.chevron_right, color: Colors.grey),
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

                // Recent Expenses Header
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Recent Expenses', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                    Text('${provider.recentExpenses.length} transactions', style: TextStyle(color: Colors.grey.shade500, fontSize: 13)),
                  ],
                ),
                const SizedBox(height: 16),

                // Upload Statement Card
                Container(
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    gradient: const LinearGradient(colors: [Color(0xFF1E3A8A), Color(0xFF3B82F6)]),
                    borderRadius: BorderRadius.circular(18),
                    boxShadow: [BoxShadow(color: const Color(0xFF1E3A8A).withOpacity(0.25), blurRadius: 12, offset: const Offset(0, 4))],
                  ),
                  child: Material(
                    color: Colors.transparent,
                    child: InkWell(
                      borderRadius: BorderRadius.circular(18),
                      onTap: () async {
                        try {
                          FilePickerResult? result = await FilePicker.platform.pickFiles(
                            type: FileType.custom,
                            allowedExtensions: ['csv', 'xlsx'],
                          );
                          if (result != null && result.files.single.path != null) {
                            if (context.mounted) _showUploadProgressDialog(context, result.files.single.path!);
                          }
                        } catch (e) {
                          if (context.mounted) {
                            UiUtils.showSnack(context, 'Upload failed: $e', isError: true);
                          }
                        }
                      },
                      child: Padding(
                        padding: const EdgeInsets.all(18),
                        child: Row(
                          children: [
                            Container(
                              padding: const EdgeInsets.all(10),
                              decoration: BoxDecoration(color: Colors.white.withOpacity(0.2), borderRadius: BorderRadius.circular(12)),
                              child: const Icon(Icons.cloud_upload_outlined, color: Colors.white, size: 26),
                            ),
                            const SizedBox(width: 16),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text('Upload Bank Statement', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 15)),
                                  const SizedBox(height: 4),
                                  Row(
                                    children: [
                                      _formatBadge('CSV'),
                                      const SizedBox(width: 6),
                                      _formatBadge('XLSX'),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                            const Icon(Icons.arrow_forward_ios, color: Colors.white, size: 16),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),

                // Expense List
                ...provider.recentExpenses.map((exp) => _buildExpenseItem(exp)).toList(),

                const SizedBox(height: 24),

                // Analytics
                _buildCategoryBreakdown(provider.budgetBreakdown),
                _buildBucketComparison(provider.budgetBreakdown),

                const SizedBox(height: 80),
              ],
            ),
          ),
          ),
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

class UploadAnimationDialog extends StatefulWidget {
  const UploadAnimationDialog({super.key});

  @override
  State<UploadAnimationDialog> createState() => _UploadAnimationDialogState();
}

class _UploadAnimationDialogState extends State<UploadAnimationDialog> with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  bool _isSuccess = false;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(vsync: this, duration: const Duration(seconds: 2))..repeat();
    
    // Simulate parsing delay before morphing into success
    Future.delayed(const Duration(milliseconds: 1500), () {
      if (mounted) {
        setState(() => _isSuccess = true);
        _controller.stop();
      }
    });
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Dialog(
      backgroundColor: Colors.white,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
      child: Padding(
        padding: const EdgeInsets.symmetric(vertical: 40, horizontal: 24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            AnimatedSwitcher(
              duration: const Duration(milliseconds: 500),
              child: _isSuccess
                  ? const Icon(Icons.check_circle, color: Colors.green, size: 80, key: ValueKey('success'))
                  : RotationTransition(
                      turns: _controller,
                      key: const ValueKey('loading'),
                      child: const Icon(Icons.sync, color: Color(0xFF1E3A8A), size: 80),
                    ),
            ),
            const SizedBox(height: 24),
            Text(
              _isSuccess ? "Upload Complete!" : "Parsing file...",
              style: TextStyle(
                fontSize: 18, 
                fontWeight: FontWeight.bold,
                color: _isSuccess ? Colors.green.shade700 : const Color(0xFF1E3A8A)
              ),
            ),
            if (!_isSuccess) ...[
              const SizedBox(height: 8),
              const Text("Reading and categorizing transactions...", style: TextStyle(color: Colors.grey, fontSize: 13)),
            ]
          ],
        ),
      ),
    );
  }
}
