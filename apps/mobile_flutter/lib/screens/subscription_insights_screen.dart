import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class SubscriptionInsightsScreen extends StatefulWidget {
  const SubscriptionInsightsScreen({super.key});

  @override
  State<SubscriptionInsightsScreen> createState() => _SubscriptionInsightsScreenState();
}

class _SubscriptionInsightsScreenState extends State<SubscriptionInsightsScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) async {
      try {
        await context.read<FinancialProvider>().loadSubscriptionInsights();
      } catch (_) {}
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(title: const Text('Subscription Audit')),
      body: Consumer<FinancialProvider>(
        builder: (context, fp, _) {
          final subs = fp.subscriptionInsights;
          if (fp.isLoading && subs.isEmpty) {
            return const Center(child: CircularProgressIndicator(color: Color(0xFF059669)));
          }

          double totalMonthly = 0;
          double totalAnnual = 0;
          double totalSavings = 0;
          for (final s in subs) {
            totalMonthly += (s['monthly_cost'] ?? 0).toDouble();
            totalAnnual += (s['annual_cost'] ?? 0).toDouble();
            totalSavings += (s['savings_opportunity'] ?? 0).toDouble();
          }

          if (subs.isEmpty) {
            return const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.subscriptions_rounded, size: 80, color: Colors.green),
                  SizedBox(height: 20),
                  Text('No Subscriptions', style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                  SizedBox(height: 8),
                  Text('Mark bills as subscriptions to track them here.', style: TextStyle(color: Colors.grey)),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () => fp.loadSubscriptionInsights(),
            child: ListView(
              padding: const EdgeInsets.all(20),
              children: [
                _buildCostSummary(totalMonthly, totalAnnual, totalSavings),
                const SizedBox(height: 24),
                ...subs.map((s) => _buildSubscriptionCard(context, fp, s)),
                const SizedBox(height: 60),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildCostSummary(double monthly, double annual, double savings) {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        gradient: const LinearGradient(colors: [Color(0xFF059669), Color(0xFF10B981)]),
        borderRadius: BorderRadius.circular(24),
        boxShadow: [BoxShadow(color: const Color(0xFF059669).withOpacity(0.3), blurRadius: 24, offset: const Offset(0, 12))],
      ),
      child: Column(
        children: [
          const Text('Total Subscription Cost', style: TextStyle(color: Colors.white70, fontWeight: FontWeight.w600)),
          const SizedBox(height: 8),
          Text('₹${annual.toStringAsFixed(0)}/year', style: const TextStyle(color: Colors.white, fontSize: 32, fontWeight: FontWeight.w900)),
          const SizedBox(height: 16),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              _buildCostItem('Monthly', '₹${monthly.toStringAsFixed(0)}'),
              _buildCostItem('Annual', '₹${annual.toStringAsFixed(0)}'),
              if (savings > 0) _buildCostItem('Savings', '₹${savings.toStringAsFixed(0)}'),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildCostItem(String label, String value) {
    return Column(
      children: [
        Text(value, style: const TextStyle(color: Colors.white, fontSize: 18, fontWeight: FontWeight.bold)),
        Text(label, style: const TextStyle(color: Colors.white70, fontSize: 12, fontWeight: FontWeight.w600)),
      ],
    );
  }

  Widget _buildSubscriptionCard(BuildContext context, FinancialProvider fp, Map<String, dynamic> sub) {
    final name = sub['name'] ?? 'Unknown';
    final monthly = (sub['monthly_cost'] ?? 0).toDouble();
    final annual = (sub['annual_cost'] ?? 0).toDouble();
    final status = sub['status'] ?? 'active';
    final savings = (sub['savings_opportunity'] ?? 0).toDouble();
    final billId = sub['bill_id'] ?? 0;

    final Color statusColor;
    final IconData statusIcon;
    switch (status) {
      case 'unused':
        statusColor = Colors.red;
        statusIcon = Icons.cancel_rounded;
        break;
      case 'flagged':
        statusColor = Colors.orange;
        statusIcon = Icons.flag_rounded;
        break;
      default:
        statusColor = Colors.green;
        statusIcon = Icons.check_circle_rounded;
    }

    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: statusColor.withOpacity(0.2)),
        boxShadow: [BoxShadow(color: statusColor.withOpacity(0.06), blurRadius: 16, offset: const Offset(0, 6))],
      ),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(color: statusColor.withOpacity(0.1), borderRadius: BorderRadius.circular(12)),
                  child: Icon(statusIcon, color: statusColor, size: 22),
                ),
                const SizedBox(width: 14),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(name, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Color(0xFF1E293B))),
                      const SizedBox(height: 2),
                      Text('₹${monthly.toStringAsFixed(0)}/mo · ₹${annual.toStringAsFixed(0)}/yr',
                          style: const TextStyle(color: Colors.grey, fontSize: 13, fontWeight: FontWeight.w600)),
                    ],
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(color: statusColor.withOpacity(0.1), borderRadius: BorderRadius.circular(12)),
                  child: Text(status.toUpperCase(), style: TextStyle(color: statusColor, fontSize: 10, fontWeight: FontWeight.bold)),
                ),
              ],
            ),
            if (savings > 0) ...[
              const SizedBox(height: 12),
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(color: Colors.green.shade50, borderRadius: BorderRadius.circular(12)),
                child: Row(
                  children: [
                    const Icon(Icons.savings_rounded, color: Colors.green, size: 18),
                    const SizedBox(width: 8),
                    Text('Potential savings: ₹${savings.toStringAsFixed(0)}/year',
                        style: const TextStyle(color: Colors.green, fontWeight: FontWeight.bold, fontSize: 13)),
                  ],
                ),
              ),
            ],
            if (status == 'active') ...[
              const SizedBox(height: 12),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: () {
                        fp.flagSubscription(billId, 'unused');
                        UiUtils.showSnack(context, 'Marked as unused');
                      },
                      icon: const Icon(Icons.remove_circle_outline, size: 16),
                      label: const Text('Mark Unused', style: TextStyle(fontSize: 12)),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: Colors.red, side: BorderSide(color: Colors.red.shade200),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                        padding: const EdgeInsets.symmetric(vertical: 10),
                      ),
                    ),
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: OutlinedButton.icon(
                      onPressed: () {
                        fp.flagSubscription(billId, 'flagged');
                        UiUtils.showSnack(context, 'Flagged for review');
                      },
                      icon: const Icon(Icons.flag_outlined, size: 16),
                      label: const Text('Flag', style: TextStyle(fontSize: 12)),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: Colors.orange, side: BorderSide(color: Colors.orange.shade200),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                        padding: const EdgeInsets.symmetric(vertical: 10),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }
}
