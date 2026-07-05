import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:provider/provider.dart';
import 'package:file_picker/file_picker.dart';
import 'package:fl_chart/fl_chart.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';
import '../main.dart';
import 'statement_insights_screen.dart';
import 'create_collection_screen.dart';
import 'collections_list_screen.dart';
import 'collection_detail_screen.dart';

class BudgetScreen extends StatefulWidget {
  const BudgetScreen({super.key});

  @override
  State<BudgetScreen> createState() => _BudgetScreenState();
}

class _BudgetScreenState extends State<BudgetScreen> {
  int _visibleCount = 5;
  int _touchedIndex = -1;
  bool _showActualSpent = true;
  String? _selectedBucket; // null = All, 'Needs', 'Wants', 'Savings'

  String _getMonthName(int month, int year) {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return '${months[month - 1]} $year';
  }

  void _showUploadProgressDialog(BuildContext context, String filePath) {
    final prov = Provider.of<FinancialProvider>(context, listen: false);

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Analyze Statement?'),
        content: const Text('Would you like to analyze this statement for recurring payments, spending patterns, and subscriptions?'),
        actions: [
          TextButton(
            onPressed: () {
              Navigator.of(ctx).pop();
              _doUpload(context, filePath, false);
            },
            child: const Text('Just Upload'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.of(ctx).pop();
              _doUpload(context, filePath, true);
            },
            child: const Text('Upload & Analyze'),
          ),
        ],
      ),
    );
  }

  void _doUpload(BuildContext context, String filePath, bool runAnalysis) {
    final prov = Provider.of<FinancialProvider>(context, listen: false);

    UiUtils.showSnack(
      context,
      runAnalysis ? 'Uploading and analyzing statement...' : 'Uploading statement...',
    );

    prov.uploadStatement(filePath, runAnalysis: runAnalysis).then((insertedCount) {
      if (navigatorKey.currentContext != null) {
        UiUtils.showSnack(
          navigatorKey.currentContext!,
          'Uploaded: $insertedCount new transactions',
        );
        if (runAnalysis && prov.lastAnalysisResult != null) {
          Navigator.of(navigatorKey.currentContext!).push(
            MaterialPageRoute(
              builder: (_) => const StatementInsightsScreen(),
            ),
          );
        }
      }
    }).catchError((e) {
      if (navigatorKey.currentContext != null) {
        UiUtils.showSnack(
          navigatorKey.currentContext!,
          'Upload failed: $e',
          isError: true,
        );
      }
    });
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
                      prefixStyle: const TextStyle(fontWeight: FontWeight.w700, color: Color(0xFF0D9488), fontSize: 16),
                      hintText: '500',
                      filled: true,
                      fillColor: Colors.white,
                      enabledBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(14),
                        borderSide: BorderSide(color: Colors.grey.shade200, width: 1.5),
                      ),
                      focusedBorder: OutlineInputBorder(
                        borderRadius: BorderRadius.circular(14),
                        borderSide: const BorderSide(color: Color(0xFF0D9488), width: 2),
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
                        borderSide: const BorderSide(color: Color(0xFF0D9488), width: 2),
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
                            color: selected ? const Color(0xFF0D9488) : Colors.white,
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(
                              color: selected ? const Color(0xFF0D9488) : Colors.grey.shade200,
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
                                Icon(Icons.autorenew_rounded, color: isRecurring ? const Color(0xFF0D9488) : Colors.grey, size: 20),
                                const SizedBox(width: 10),
                                Text('Make this a Recurring Bill', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 14, color: isRecurring ? const Color(0xFF1E293B) : Colors.grey.shade600)),
                              ],
                            ),
                            Switch(
                              value: isRecurring,
                              onChanged: (val) => setState(() => isRecurring = val),
                              activeThumbColor: const Color(0xFF0D9488),
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
                                activeThumbColor: const Color(0xFF0D9488),
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
                        backgroundColor: const Color(0xFF0D9488),
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

  List<dynamic> _filteredExpenses(FinancialProvider p) {
    if (_selectedBucket == null) return p.recentExpenses;
    return p.recentExpenses.where((e) => e['bucket'] == _selectedBucket).toList();
  }

  Widget _filterChip(String label, String? bucket, Color color) {
    final selected = _selectedBucket == bucket;
    return Expanded(
      child: GestureDetector(
        onTap: () {
          setState(() {
            _selectedBucket = bucket;
            _visibleCount = 5;
          });
        },
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          padding: const EdgeInsets.symmetric(vertical: 8),
          decoration: BoxDecoration(
            color: selected ? color : Colors.transparent,
            borderRadius: BorderRadius.circular(10),
          ),
          child: Text(
            label,
            textAlign: TextAlign.center,
            style: TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w700,
              color: selected ? Colors.white : color,
            ),
          ),
        ),
      ),
    );
  }

  Widget _formatBadge(String label) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
      decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: BorderRadius.circular(6)),
      child: Text(label, style: const TextStyle(color: Colors.white, fontSize: 10, fontWeight: FontWeight.w600)),
    );
  }

  // ── CATEGORY BREAKDOWN ──────────────────────────────────────────────────────
  Widget _buildCategoryBreakdown(Map<String, dynamic>? bd) {
    if (bd == null) return const SizedBox.shrink();
    final categories = bd['categoryBreakdown'] as List<dynamic>? ?? [];
    if (categories.isEmpty) return const SizedBox.shrink();

    final colors = [
      const Color(0xFF2DD4BF), const Color(0xFFF59E0B), const Color(0xFF10B981),
      const Color(0xFFEF4444), const Color(0xFF8B5CF6), const Color(0xFFEC4899),
      const Color(0xFF06B6D4), const Color(0xFFF97316),
    ];

    return _budgetCard(
      title: 'Spending by Category',
      icon: Icons.pie_chart_rounded,
      color: const Color(0xFF2DD4BF),
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
                color: (isUp ? (isWants ? Colors.red : Colors.green) : isDown ? (isWants ? Colors.green : Colors.red) : Colors.grey).withValues(alpha: 0.1),
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
        boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.03), blurRadius: 16, offset: const Offset(0, 6))],
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [
          Container(padding: const EdgeInsets.all(8), decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(10)), child: Icon(icon, color: color, size: 18)),
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
        color: isSelected ? const Color(0xFF0D9488) : Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: isSelected ? null : Border.all(color: Colors.grey.withValues(alpha: 0.2)),
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
            backgroundColor: Colors.grey.withValues(alpha: 0.2),
            barRadius: const Radius.circular(8),
            padding: EdgeInsets.zero,
          )
        ],
      ),
    );
  }

  List<PieChartSectionData> _getPieChartSections(FinancialProvider provider) {
    double needsValue = _showActualSpent ? provider.needsSpent : provider.needsBudget;
    double wantsValue = _showActualSpent ? provider.wantsSpent : provider.wantsBudget;
    double savingsValue = _showActualSpent ? provider.savingsSpent : provider.savingsBudget;
    
    double total = needsValue + wantsValue + savingsValue;
    if (total == 0) {
      return [
        PieChartSectionData(
          color: Colors.grey.shade300,
          value: 1,
          title: '',
          radius: 40,
        )
      ];
    }

    final needsPct = (needsValue / total * 100).round();
    final wantsPct = (wantsValue / total * 100).round();
    final savingsPct = (savingsValue / total * 100).round();

    return [
      PieChartSectionData(
        color: const Color(0xFFE88A1A),
        value: needsValue,
        title: needsPct > 5 ? '$needsPct%' : '',
        radius: _touchedIndex == 0 ? 46.0 : 40.0,
        titleStyle: TextStyle(
          fontSize: _touchedIndex == 0 ? 15.0 : 12.0,
          fontWeight: FontWeight.bold,
          color: Colors.white,
        ),
      ),
      PieChartSectionData(
        color: const Color(0xFF6B46C1),
        value: wantsValue,
        title: wantsPct > 5 ? '$wantsPct%' : '',
        radius: _touchedIndex == 1 ? 46.0 : 40.0,
        titleStyle: TextStyle(
          fontSize: _touchedIndex == 1 ? 15.0 : 12.0,
          fontWeight: FontWeight.bold,
          color: Colors.white,
        ),
      ),
      PieChartSectionData(
        color: const Color(0xFF059669),
        value: savingsValue,
        title: savingsPct > 5 ? '$savingsPct%' : '',
        radius: _touchedIndex == 2 ? 46.0 : 40.0,
        titleStyle: TextStyle(
          fontSize: _touchedIndex == 2 ? 15.0 : 12.0,
          fontWeight: FontWeight.bold,
          color: Colors.white,
        ),
      ),
    ];
  }

  Widget _buildPieChartCard(FinancialProvider provider) {
    double totalBudget = provider.needsBudget + provider.wantsBudget + provider.savingsBudget;
    double displayedTotal = _showActualSpent ? provider.totalSpent : totalBudget;

    return Container(
      margin: const EdgeInsets.only(bottom: 20),
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(color: Colors.grey.withValues(alpha: 0.1)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.03),
            blurRadius: 16,
            offset: const Offset(0, 6),
          )
        ],
      ),
      child: Column(
        children: [
          // Header / Toggle
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Text(
                _showActualSpent ? 'Spending Breakdown' : 'Budget Targets',
                style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Color(0xFF1E293B)),
              ),
              Container(
                decoration: BoxDecoration(
                  color: Colors.grey.shade100,
                  borderRadius: BorderRadius.circular(12),
                ),
                padding: const EdgeInsets.all(3),
                child: Row(
                  children: [
                    GestureDetector(
                      onTap: () => setState(() => _showActualSpent = true),
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                        decoration: BoxDecoration(
                          color: _showActualSpent ? Colors.white : Colors.transparent,
                          borderRadius: BorderRadius.circular(9),
                          boxShadow: _showActualSpent
                              ? [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 4, offset: const Offset(0, 2))]
                              : null,
                        ),
                        child: Text(
                          'Actual',
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.bold,
                            color: _showActualSpent ? const Color(0xFF0D9488) : Colors.grey.shade600,
                          ),
                        ),
                      ),
                    ),
                    GestureDetector(
                      onTap: () => setState(() => _showActualSpent = false),
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                        decoration: BoxDecoration(
                          color: !_showActualSpent ? Colors.white : Colors.transparent,
                          borderRadius: BorderRadius.circular(9),
                          boxShadow: !_showActualSpent
                              ? [BoxShadow(color: Colors.black.withValues(alpha: 0.05), blurRadius: 4, offset: const Offset(0, 2))]
                              : null,
                        ),
                        child: Text(
                          'Target',
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.bold,
                            color: !_showActualSpent ? const Color(0xFF0D9488) : Colors.grey.shade600,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
          const SizedBox(height: 24),
          
          // Pie Chart Stack
          SizedBox(
            height: 180,
            child: Stack(
              children: [
                PieChart(
                  PieChartData(
                    pieTouchData: PieTouchData(
                      touchCallback: (FlTouchEvent event, pieTouchResponse) {
                        setState(() {
                          if (!event.isInterestedForInteractions ||
                              pieTouchResponse == null ||
                              pieTouchResponse.touchedSection == null) {
                            _touchedIndex = -1;
                            return;
                          }
                          _touchedIndex = pieTouchResponse.touchedSection!.touchedSectionIndex;
                        });
                      },
                    ),
                    borderData: FlBorderData(show: false),
                    sectionsSpace: 4,
                    centerSpaceRadius: 55,
                    sections: _getPieChartSections(provider),
                  ),
                ),
                Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Text(
                        _showActualSpent ? 'Total Spent' : 'Total Target',
                        style: TextStyle(fontSize: 11, color: Colors.grey.shade500, fontWeight: FontWeight.w600),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '₹${displayedTotal.toStringAsFixed(0)}',
                        style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w900, color: Color(0xFF1E293B)),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),

          // Custom Legend / Info Rows
          Column(
            children: [
              _buildLegendRow(
                'Needs',
                _showActualSpent ? provider.needsSpent : provider.needsBudget,
                _showActualSpent ? provider.needsBudget : 0,
                const Color(0xFFE88A1A),
                50,
              ),
              const SizedBox(height: 10),
              _buildLegendRow(
                'Wants',
                _showActualSpent ? provider.wantsSpent : provider.wantsBudget,
                _showActualSpent ? provider.wantsBudget : 0,
                const Color(0xFF6B46C1),
                30,
              ),
              const SizedBox(height: 10),
              _buildLegendRow(
                'Savings',
                _showActualSpent ? provider.savingsSpent : provider.savingsBudget,
                _showActualSpent ? provider.savingsBudget : 0,
                const Color(0xFF059669),
                20,
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildLegendRow(String label, double amount, double limit, Color color, int targetPct) {
    double spentPctOfLimit = limit > 0 ? (amount / limit * 100) : 0.0;
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: Colors.grey.shade50.withValues(alpha: 0.8),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: Colors.grey.shade100),
      ),
      child: Row(
        children: [
          Container(
            width: 14,
            height: 14,
            decoration: BoxDecoration(
              color: color,
              borderRadius: BorderRadius.circular(4),
            ),
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      '$label ($targetPct% target)',
                      style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Color(0xFF1E293B)),
                    ),
                    Text(
                      '₹${amount.toStringAsFixed(0)}',
                      style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Color(0xFF1E293B)),
                    ),
                  ],
                ),
                if (_showActualSpent && limit > 0) ...[
                  const SizedBox(height: 4),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Budget: ₹${limit.toStringAsFixed(0)}',
                        style: TextStyle(color: Colors.grey.shade500, fontSize: 11),
                      ),
                      Text(
                        '${spentPctOfLimit.toInt()}% used',
                        style: TextStyle(
                          color: spentPctOfLimit > 100 ? Colors.red : Colors.grey.shade600,
                          fontWeight: FontWeight.bold,
                          fontSize: 11,
                        ),
                      ),
                    ],
                  ),
                ],
              ],
            ),
          ),
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

    return GestureDetector(
      onLongPress: () => _showExpenseActions(context, exp),
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(16),
          border: Border.all(color: Colors.grey.withValues(alpha: 0.08)),
          boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.03), blurRadius: 8, offset: const Offset(0, 2))],
        ),
        child: Row(
          children: [
            Container(
              width: 44, height: 44,
              decoration: BoxDecoration(color: color.withValues(alpha: 0.12), shape: BoxShape.circle),
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
                        decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(10)),
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
      ),
    );
  }

  void _showExpenseActions(BuildContext context, dynamic exp) {
    final name = exp['name'] ?? '';
    final amount = (exp['amount'] ?? 0).toDouble();
    final fp = Provider.of<FinancialProvider>(context, listen: false);

    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 40, height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey.shade300,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 20),
              Text('₹${amount.toStringAsFixed(0)} — $name',
                  style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w700),
                  textAlign: TextAlign.center),
              const SizedBox(height: 4),
              Text('What would you like to do with this transaction?',
                  style: TextStyle(color: Colors.grey.shade600, fontSize: 13)),
              const SizedBox(height: 24),

              // Create Collection
              ListTile(
                leading: Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: const Color(0xFF6B46C1).withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: const Icon(Icons.people_alt_outlined, color: Color(0xFF6B46C1)),
                ),
                title: const Text('Create Collection from this',
                    style: TextStyle(fontWeight: FontWeight.w600)),
                subtitle: const Text('Make this a member in a new collection'),
                onTap: () {
                  Navigator.of(ctx).pop();
                  Navigator.of(context).push(
                    MaterialPageRoute(
                      builder: (_) => CreateCollectionScreen(prefillAmount: amount, prefillName: name),
                    ),
                  );
                },
              ),
              const Divider(),
              ListTile(
                leading: Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: const Color(0xFF0D9488).withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: const Icon(Icons.link, color: Color(0xFF0D9488)),
                ),
                title: const Text('Link to existing Collection',
                    style: TextStyle(fontWeight: FontWeight.w600)),
                subtitle: const Text('Mark this as payment for a collection member'),
                onTap: () {
                  Navigator.of(ctx).pop();
                  _showCollectionPicker(context, exp);
                },
              ),
              const SizedBox(height: 12),
            ],
          ),
        ),
      ),
    );
  }

  void _showCollectionPicker(BuildContext context, dynamic exp) {
    final fp = Provider.of<FinancialProvider>(context, listen: false);
    final active = fp.collections.where((c) => c['status'] == 'active').toList();

    if (active.isEmpty) {
      UiUtils.showSnack(context, 'No active collections. Create one first.', isError: true);
      return;
    }

    showModalBottomSheet(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Container(
                width: 40, height: 4,
                decoration: BoxDecoration(
                  color: Colors.grey.shade300,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 20),
              const Text('Select Collection',
                  style: TextStyle(fontSize: 18, fontWeight: FontWeight.w700)),
              const SizedBox(height: 4),
              const Text('Which collection does this payment belong to?',
                  style: TextStyle(color: Colors.grey, fontSize: 13)),
              const SizedBox(height: 16),
              ...active.map((c) {
                final label = c['label'] ?? '';
                final collected = (c['total_collected'] ?? 0).toDouble();
                final expected = (c['total_expected'] ?? 0).toDouble();
                return ListTile(
                  leading: Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: const Color(0xFF0D9488).withValues(alpha: 0.08),
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: const Icon(Icons.people_outline, color: Color(0xFF0D9488), size: 20),
                  ),
                  title: Text(label, style: const TextStyle(fontWeight: FontWeight.w600)),
                  subtitle: Text('₹${collected.toStringAsFixed(0)} / ₹${expected.toStringAsFixed(0)} collected'),
                  onTap: () {
                    Navigator.of(ctx).pop();
                    _showMemberPicker(context, c, exp);
                  },
                );
              }),
            ],
          ),
        ),
      ),
    );
  }

  void _showMemberPicker(BuildContext context, Map<String, dynamic> collection, dynamic exp) {
    final amount = (exp['amount'] ?? 0).toDouble();
    final expName = exp['name'] ?? '';

    Navigator.of(context).push(
      MaterialPageRoute(
        builder: (_) => CollectionDetailScreen(
          collection: collection,
          highlightAmount: amount,
          highlightName: expName,
        ),
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
                  child: RefreshIndicator(
                    onRefresh: () async {
                      await provider.loadAllData();
                    },
                    child: SingleChildScrollView(
                      physics: const AlwaysScrollableScrollPhysics(),
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

                _buildPieChartCard(provider),
                const SizedBox(height: 20),

                // Committed Bills
                Container(
                  padding: const EdgeInsets.all(20),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(16),
                    border: Border.all(color: Colors.grey.withValues(alpha: 0.1)),
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
                    Text('${_filteredExpenses(provider).length} transactions',
                        style: TextStyle(color: Colors.grey.shade500, fontSize: 13)),
                  ],
                ),
                const SizedBox(height: 12),

                // Bucket Filter
                Container(
                  padding: const EdgeInsets.all(4),
                  decoration: BoxDecoration(
                    color: const Color(0xFFF1F5F9),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    children: [
                      _filterChip('All', null, const Color(0xFF64748B)),
                      const SizedBox(width: 4),
                      _filterChip('Needs', 'Needs', const Color(0xFFE88A1A)),
                      const SizedBox(width: 4),
                      _filterChip('Wants', 'Wants', const Color(0xFF6B46C1)),
                      const SizedBox(width: 4),
                      _filterChip('Savings', 'Savings', const Color(0xFF059669)),
                    ],
                  ),
                ),
                const SizedBox(height: 16),

                // Upload Statement Card
                Container(
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    gradient: const LinearGradient(colors: [Color(0xFF0D9488), Color(0xFF2DD4BF)]),
                    borderRadius: BorderRadius.circular(18),
                    boxShadow: [BoxShadow(color: const Color(0xFF0D9488).withValues(alpha: 0.25), blurRadius: 12, offset: const Offset(0, 4))],
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
                              decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.2), borderRadius: BorderRadius.circular(12)),
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

                // Analyze Now Card
                Container(
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(18),
                    border: Border.all(color: const Color(0xFFE2E8F0)),
                  ),
                  child: Material(
                    color: Colors.transparent,
                    child: InkWell(
                      borderRadius: BorderRadius.circular(18),
                      onTap: () async {
                        final fp = Provider.of<FinancialProvider>(context, listen: false);
                        UiUtils.showSnack(context, 'Analyzing your finances...');
                        await fp.runFullAnalysis();
                        if (context.mounted) {
                          Navigator.of(context).push(
                            MaterialPageRoute(builder: (_) => const StatementInsightsScreen()),
                          );
                        }
                      },
                      child: Padding(
                        padding: const EdgeInsets.all(18),
                        child: Row(
                          children: [
                            Container(
                              padding: const EdgeInsets.all(10),
                              decoration: BoxDecoration(
                                color: const Color(0xFF0D9488).withValues(alpha: 0.08),
                                borderRadius: BorderRadius.circular(12),
                              ),
                              child: const Icon(Icons.analytics_outlined, color: Color(0xFF0D9488), size: 26),
                            ),
                            const SizedBox(width: 16),
                            const Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text('Analyze Existing Data',
                                      style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15, color: Color(0xFF1E293B))),
                                  SizedBox(height: 4),
                                  Text('Detect recurring, patterns & subscriptions',
                                      style: TextStyle(fontSize: 12, color: Colors.grey)),
                                ],
                              ),
                            ),
                            const Icon(Icons.arrow_forward_ios, color: Colors.grey, size: 16),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),



                // Expense List
                ..._filteredExpenses(provider).take(_visibleCount).map((exp) => _buildExpenseItem(exp)),

                if (_filteredExpenses(provider).length > _visibleCount) ...[
                  const SizedBox(height: 12),
                  Center(
                    child: OutlinedButton.icon(
                      onPressed: () {
                        setState(() {
                          _visibleCount += 10;
                        });
                      },
                      style: OutlinedButton.styleFrom(
                        foregroundColor: const Color(0xFF0D9488),
                        side: const BorderSide(color: Color(0xFF0D9488), width: 1.5),
                        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                      ),
                      icon: const Icon(Icons.expand_more_rounded, size: 20),
                      label: Text(
                        'Show ${(_filteredExpenses(provider).length - _visibleCount).clamp(0, 10)} More',
                        style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
                      ),
                    ),
                  ),
                ],

                const SizedBox(height: 24),

                // Analytics
                _buildCategoryBreakdown(provider.budgetBreakdown),
                _buildBucketComparison(provider.budgetBreakdown),

                const SizedBox(height: 80),
              ],
            ),
          ),
          ),
          ),
        ],
      ),
    ),
    floatingActionButton: FloatingActionButton(
            onPressed: () => _showAddExpenseModal(context),
            backgroundColor: const Color(0xFF0D9488),
            child: const Icon(Icons.add, color: Colors.white),
          ),
        );
      }
    );
  }
}

