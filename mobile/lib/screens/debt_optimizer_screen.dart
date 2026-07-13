import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';

class DebtOptimizerScreen extends StatefulWidget {
  const DebtOptimizerScreen({super.key});

  @override
  State<DebtOptimizerScreen> createState() => _DebtOptimizerScreenState();
}

class _DebtOptimizerScreenState extends State<DebtOptimizerScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) async {
      try {
        await context.read<FinancialProvider>().loadDebtOptimizer();
      } catch (_) {}
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(title: const Text('Debt Optimizer')),
      body: Consumer<FinancialProvider>(
        builder: (context, fp, _) {
          final data = fp.debtOptimizer;
          if (data == null) {
            return const Center(child: CircularProgressIndicator(color: Color(0xFF6366F1)));
          }

          final totalOutstanding = (data['total_outstanding'] ?? 0).toDouble();
          if (totalOutstanding == 0) {
            return const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.celebration_rounded, size: 80, color: Colors.green),
                  SizedBox(height: 20),
                  Text('No Debts!', style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                  SizedBox(height: 8),
                  Text('You are debt free. Keep it up!', style: TextStyle(color: Colors.grey)),
                ],
              ),
            );
          }

          final snowball = data['snowball'] ?? {};
          final avalanche = data['avalanche'] ?? {};
          final recommended = data['recommended_strategy'] ?? 'Avalanche';
          final freedomMonths = data['estimated_freedom_months'] ?? 0;
          final interestSavable = (data['total_interest_savable'] ?? 0).toDouble();

          return RefreshIndicator(
            onRefresh: () => fp.loadDebtOptimizer(),
            child: ListView(
              padding: const EdgeInsets.all(20),
              children: [
                _buildSummaryCard(totalOutstanding, (data['total_monthly_emis'] ?? 0).toDouble(), (data['total_interest_paid'] ?? 0).toDouble()),
                const SizedBox(height: 20),
                _buildComparisonCard('Snowball', snowball, recommended, const Color(0xFF0891B2)),
                const SizedBox(height: 16),
                _buildComparisonCard('Avalanche', avalanche, recommended, const Color(0xFFDC2626)),
                const SizedBox(height: 24),
                _buildRecommendationCard(recommended, freedomMonths, interestSavable),
                const SizedBox(height: 60),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildSummaryCard(double total, double emis, double interest) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        gradient: const LinearGradient(colors: [Color(0xFF6366F1), Color(0xFF312E81)]),
        borderRadius: BorderRadius.circular(24),
        boxShadow: [BoxShadow(color: const Color(0xFF6366F1).withValues(alpha: 0.3), blurRadius: 24, offset: const Offset(0, 12))],
      ),
      child: Column(
        children: [
          const Text('Total Debt Outstanding', style: TextStyle(color: Colors.white70, fontWeight: FontWeight.w600)),
          const SizedBox(height: 8),
          Text('₹${total.toStringAsFixed(0)}', style: const TextStyle(color: Colors.white, fontSize: 36, fontWeight: FontWeight.w900)),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _buildSummaryItem('Monthly EMIs', '₹${emis.toStringAsFixed(0)}', Colors.white70, Colors.white),
              _buildSummaryItem('Total Interest', '₹${interest.toStringAsFixed(0)}', Colors.white70, Colors.white),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildSummaryItem(String label, String value, Color labelColor, Color valueColor) {
    return Column(
      children: [
        Text(label, style: TextStyle(color: labelColor, fontSize: 12, fontWeight: FontWeight.w600)),
        const SizedBox(height: 4),
        Text(value, style: TextStyle(color: valueColor, fontSize: 16, fontWeight: FontWeight.bold)),
      ],
    );
  }

  Widget _buildComparisonCard(String name, Map<String, dynamic> data, String recommended, Color color) {
    final months = data['months_to_freedom'] ?? 0;
    final interest = (data['total_interest'] ?? 0).toDouble();
    final isRecommended = name.toLowerCase() == recommended.toLowerCase();

    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: isRecommended ? color : Colors.grey.shade200, width: isRecommended ? 2 : 1),
        boxShadow: [BoxShadow(color: color.withValues(alpha: isRecommended ? 0.1 : 0.04), blurRadius: 20, offset: const Offset(0, 8))],
      ),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(12)),
                  child: Icon(isRecommended ? Icons.star_rounded : Icons.trending_flat_rounded, color: color, size: 22),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(name, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 18, color: Color(0xFF1E293B))),
                ),
                if (isRecommended)
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                    decoration: BoxDecoration(color: color, borderRadius: BorderRadius.circular(12)),
                    child: const Text('BEST', style: TextStyle(color: Colors.white, fontSize: 11, fontWeight: FontWeight.bold)),
                  ),
              ],
            ),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                Column(
                  children: [
                    Text('$months', style: TextStyle(color: color, fontSize: 28, fontWeight: FontWeight.w900)),
                    const Text('Months', style: TextStyle(color: Colors.grey, fontSize: 12, fontWeight: FontWeight.w600)),
                  ],
                ),
                Column(
                  children: [
                    Text('₹${interest.toStringAsFixed(0)}', style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w900, color: Color(0xFF1E293B))),
                    const Text('Total Interest', style: TextStyle(color: Colors.grey, fontSize: 12, fontWeight: FontWeight.w600)),
                  ],
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRecommendationCard(String strategy, int months, double savings) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        gradient: LinearGradient(colors: strategy == 'Avalanche'
            ? [const Color(0xFFDC2626), const Color(0xFFF87171)]
            : [const Color(0xFF0891B2), const Color(0xFF22D3EE)]),
        borderRadius: BorderRadius.circular(24),
        boxShadow: [BoxShadow(color: const Color(0xFF6366F1).withValues(alpha: 0.1), blurRadius: 24, offset: const Offset(0, 12))],
      ),
      child: Column(
        children: [
          const Text('Recommended Strategy', style: TextStyle(color: Colors.white70, fontWeight: FontWeight.w600)),
          const SizedBox(height: 8),
          Text(strategy, style: const TextStyle(color: Colors.white, fontSize: 28, fontWeight: FontWeight.w900)),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              Column(
                children: [
                  Text('$months months', style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
                  const Text('To Freedom', style: TextStyle(color: Colors.white70, fontSize: 12)),
                ],
              ),
              Column(
                children: [
                  Text('₹${savings.toStringAsFixed(0)}', style: const TextStyle(color: Colors.white, fontSize: 20, fontWeight: FontWeight.bold)),
                  const Text('Interest Saved', style: TextStyle(color: Colors.white70, fontSize: 12)),
                ],
              ),
            ],
          ),
        ],
      ),
    );
  }
}
