import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import '../models/wishlist_item.dart';
import '../models/discipline_aggregates.dart';
import '../services/discipline_service.dart';

class DisciplineDashboardScreen extends StatefulWidget {
  final DisciplineService disciplineService;

  const DisciplineDashboardScreen({Key? key, required this.disciplineService}) : super(key: key);

  @override
  _DisciplineDashboardScreenState createState() => _DisciplineDashboardScreenState();
}

class _DisciplineDashboardScreenState extends State<DisciplineDashboardScreen> {
  bool _isLoading = true;
  PayYourselfFirst? _payYourselfFirst;
  ZeroBasedBudget? _zeroBasedBudget;
  List<WishlistItem> _wishlist = [];
  FinancialGuardrails? _guardrails;
  DebtRepaymentStrategy? _debtStrategy;

  @override
  void initState() {
    super.initState();
    _loadData();
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);
    try {
      final futures = await Future.wait([
        widget.disciplineService.getPayYourselfFirst(),
        widget.disciplineService.getZeroBasedBudget(),
        widget.disciplineService.getWishlist(),
        widget.disciplineService.getGuardrails(),
        widget.disciplineService.getDebtStrategy(),
      ]);

      setState(() {
        _payYourselfFirst = futures[0] as PayYourselfFirst;
        _zeroBasedBudget = futures[1] as ZeroBasedBudget;
        _wishlist = futures[2] as List<WishlistItem>;
        _guardrails = futures[3] as FinancialGuardrails;
        _debtStrategy = futures[4] as DebtRepaymentStrategy;
      });
    } catch (e) {
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error loading discipline data: $e')));
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  Widget _buildGamifiedCard({required String title, required Widget child, required IconData icon, required List<Color> gradient}) {
    return Container(
      margin: const EdgeInsets.only(bottom: 20),
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(28),
        boxShadow: [
          BoxShadow(color: gradient[0].withOpacity(0.15), blurRadius: 20, offset: const Offset(0, 10))
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Container(
                padding: const EdgeInsets.all(10),
                decoration: BoxDecoration(
                  gradient: LinearGradient(colors: gradient),
                  shape: BoxShape.circle,
                ),
                child: Icon(icon, color: Colors.white, size: 20),
              ),
              const SizedBox(width: 12),
              Text(title, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: Color(0xFF1E293B))),
            ],
          ),
          const SizedBox(height: 20),
          child,
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(
        backgroundColor: Color(0xFFF8FAFC),
        body: Center(child: CircularProgressIndicator(color: Color(0xFF1E3A8A))),
      );
    }

    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        title: const Text('Discipline HUD', style: TextStyle(fontWeight: FontWeight.w900)),
        backgroundColor: Colors.transparent,
        elevation: 0,
        centerTitle: false,
      ),
      body: RefreshIndicator(
        onRefresh: _loadData,
        child: ListView(
          padding: const EdgeInsets.all(20),
          children: [
            _buildPayYourselfFirstCard(),
            _buildZeroBasedBudgetCard(),
            _buildGuardrailsCard(),
            _buildDebtStrategyCard(),
            _buildWishlistCard(),
            const SizedBox(height: 60),
          ],
        ),
      ),
    );
  }

  Widget _buildPayYourselfFirstCard() {
    if (_payYourselfFirst == null) return const SizedBox();
    final p = _payYourselfFirst!;
    final progress = p.targetAmount > 0 ? (p.actualSavings / p.targetAmount) : 0.0;
    
    return _buildGamifiedCard(
      title: 'Pay Yourself First',
      icon: Icons.savings_rounded,
      gradient: const [Color(0xFF059669), Color(0xFF10B981)],
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Saved this month', style: TextStyle(color: Colors.grey, fontWeight: FontWeight.w600)),
                  Text('₹${p.actualSavings.toStringAsFixed(0)}', style: const TextStyle(fontSize: 24, fontWeight: FontWeight.w900, color: Color(0xFF059669))),
                ],
              ),
              Container(
                padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                decoration: BoxDecoration(color: p.status == 'On Track' ? Colors.green.shade50 : Colors.red.shade50, borderRadius: BorderRadius.circular(20)),
                child: Text(p.status, style: TextStyle(color: p.status == 'On Track' ? Colors.green.shade700 : Colors.red.shade700, fontWeight: FontWeight.bold, fontSize: 12)),
              )
            ],
          ),
          const SizedBox(height: 16),
          LinearPercentIndicator(
            lineHeight: 12.0,
            percent: progress.clamp(0.0, 1.0),
            barRadius: const Radius.circular(10),
            backgroundColor: Colors.grey.shade200,
            linearGradient: const LinearGradient(colors: [Color(0xFF059669), Color(0xFF34D399)]),
            padding: EdgeInsets.zero,
          ),
          const SizedBox(height: 8),
          Text('Target: 20% of Income (₹${p.targetAmount.toStringAsFixed(0)})', style: const TextStyle(color: Colors.grey, fontSize: 13, fontWeight: FontWeight.w600)),
        ],
      )
    );
  }

  Widget _buildZeroBasedBudgetCard() {
    if (_zeroBasedBudget == null) return const SizedBox();
    final z = _zeroBasedBudget!;
    final diff = z.totalIncome - z.totalAllocated;
    
    return _buildGamifiedCard(
      title: 'Zero-Based Budget',
      icon: Icons.pie_chart_rounded,
      gradient: const [Color(0xFF6B46C1), Color(0xFF8B5CF6)],
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              CircularPercentIndicator(
                radius: 45.0,
                lineWidth: 10.0,
                percent: z.totalIncome > 0 ? (z.totalAllocated / z.totalIncome).clamp(0.0, 1.0) : 0,
                circularStrokeCap: CircularStrokeCap.round,
                linearGradient: const LinearGradient(colors: [Color(0xFF6B46C1), Color(0xFFA78BFA)]),
                backgroundColor: Colors.grey.shade100,
                center: const Icon(Icons.account_balance_wallet_rounded, color: Color(0xFF6B46C1)),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Unallocated', style: TextStyle(color: Colors.grey, fontWeight: FontWeight.w600)),
                  Text('₹${diff.toStringAsFixed(0)}', style: TextStyle(fontSize: 24, fontWeight: FontWeight.w900, color: diff == 0 ? Colors.green : Colors.orange)),
                  const SizedBox(height: 4),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                    decoration: BoxDecoration(color: z.status == 'Zero-Based' ? Colors.green.shade50 : Colors.orange.shade50, borderRadius: BorderRadius.circular(12)),
                    child: Text(z.status, style: TextStyle(color: z.status == 'Zero-Based' ? Colors.green.shade700 : Colors.orange.shade700, fontSize: 11, fontWeight: FontWeight.bold)),
                  )
                ],
              )
            ],
          ),
        ],
      )
    );
  }

  Widget _buildGuardrailsCard() {
    if (_guardrails == null) return const SizedBox();
    final g = _guardrails!;
    return _buildGamifiedCard(
      title: 'Financial Guardrails',
      icon: Icons.shield_rounded,
      gradient: const [Color(0xFF2563EB), Color(0xFF60A5FA)],
      child: Column(
        children: [
          _buildGuardrailRow('Emergency Runway', '${g.emergencyFundRatio.toStringAsFixed(1)} months', g.runwayStatus),
          const Padding(padding: EdgeInsets.symmetric(vertical: 12), child: Divider(height: 1)),
          _buildGuardrailRow('Housing Cost', '${g.housingCostRatio.toStringAsFixed(1)}%', g.housingStatus),
        ],
      )
    );
  }

  Widget _buildGuardrailRow(String label, String value, String status) {
    bool isHealthy = status == 'Healthy';
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(label, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
            const SizedBox(height: 4),
            Text(value, style: const TextStyle(color: Colors.grey, fontSize: 13, fontWeight: FontWeight.w600)),
          ],
        ),
        Icon(isHealthy ? Icons.check_circle_rounded : Icons.warning_rounded, color: isHealthy ? Colors.green : Colors.orange, size: 28),
      ],
    );
  }

  Widget _buildDebtStrategyCard() {
    if (_debtStrategy == null) return const SizedBox();
    final d = _debtStrategy!;
    return _buildGamifiedCard(
      title: 'Debt Strategy',
      icon: Icons.trending_down_rounded,
      gradient: const [Color(0xFFDC2626), Color(0xFFF87171)],
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(color: Colors.red.shade50, borderRadius: BorderRadius.circular(16)),
            child: Row(
              children: [
                const Icon(Icons.auto_awesome_rounded, color: Color(0xFFDC2626)),
                const SizedBox(width: 12),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text('Recommended Plan', style: TextStyle(color: Colors.grey, fontSize: 12, fontWeight: FontWeight.w600)),
                    Text(d.recommendedStrategy, style: const TextStyle(color: Color(0xFFDC2626), fontWeight: FontWeight.bold, fontSize: 16)),
                  ],
                ),
              ],
            ),
          )
        ],
      )
    );
  }

  Widget _buildWishlistCard() {
    return _buildGamifiedCard(
      title: 'Cooling-off Wishlist',
      icon: Icons.ac_unit_rounded,
      gradient: const [Color(0xFF0891B2), Color(0xFF22D3EE)],
      child: Column(
        children: [
          ..._wishlist.map((item) {
            bool unlocked = item.isUnlocked;
            bool bought = item.status == 'bought';
            return Container(
              margin: const EdgeInsets.only(bottom: 12),
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                border: Border.all(color: Colors.grey.shade200),
                borderRadius: BorderRadius.circular(16),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(item.name, style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15, decoration: bought ? TextDecoration.lineThrough : null)),
                        const SizedBox(height: 4),
                        Text('₹${item.amount.toStringAsFixed(0)}', style: const TextStyle(color: Colors.grey, fontWeight: FontWeight.w600)),
                      ],
                    ),
                  ),
                  if (bought)
                    const Icon(Icons.check_circle_rounded, color: Colors.green)
                  else if (unlocked)
                    ElevatedButton(
                      onPressed: () => _updateWishlistStatus(item.id, 'bought'),
                      style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0891B2), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12))),
                      child: const Text('Buy Now', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                    )
                  else
                    Column(
                      children: [
                        const Icon(Icons.lock_rounded, color: Colors.orange, size: 20),
                        const SizedBox(height: 4),
                        Text('${item.unlockDate.difference(DateTime.now()).inDays}d left', style: const TextStyle(color: Colors.orange, fontSize: 11, fontWeight: FontWeight.bold)),
                      ],
                    )
                ],
              ),
            );
          }).toList(),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: _addWishlistItemDialog,
              icon: const Icon(Icons.add_rounded),
              label: const Text('Add to Wishlist'),
              style: OutlinedButton.styleFrom(shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12))),
            ),
          )
        ],
      )
    );
  }

  void _addWishlistItemDialog() {
    final nameCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Add Wishlist Item', style: TextStyle(fontWeight: FontWeight.bold)),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(24)),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(controller: nameCtrl, decoration: InputDecoration(labelText: 'Item Name', filled: true, border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none))),
            const SizedBox(height: 12),
            TextField(controller: amountCtrl, decoration: InputDecoration(labelText: 'Amount (₹)', filled: true, border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none)), keyboardType: TextInputType.number),
          ],
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          ElevatedButton(
            onPressed: () async {
              try {
                await widget.disciplineService.addWishlistItem(nameCtrl.text, double.parse(amountCtrl.text), 30);
                if (context.mounted) Navigator.pop(context);
                _loadData();
              } catch(e) {
                print(e);
              }
            },
            style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF1E3A8A), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12))),
            child: const Text('Lock for 30 Days', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
    );
  }

  void _updateWishlistStatus(int id, String status) async {
    try {
      await widget.disciplineService.updateWishlistItemStatus(id, status);
      _loadData();
    } catch(e) {
      print(e);
    }
  }
}
