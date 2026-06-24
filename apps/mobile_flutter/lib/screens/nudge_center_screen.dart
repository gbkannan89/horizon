import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';

class NudgeCenterScreen extends StatefulWidget {
  const NudgeCenterScreen({super.key});

  @override
  State<NudgeCenterScreen> createState() => _NudgeCenterScreenState();
}

class _NudgeCenterScreenState extends State<NudgeCenterScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(title: const Text('Nudge Center')),
      body: Consumer<FinancialProvider>(
        builder: (context, fp, _) {
          final active = fp.nudges.where((n) => n['is_dismissed'] == false).toList();

          if (active.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.check_circle_outline_rounded, size: 80, color: Colors.green.shade300),
                  const SizedBox(height: 20),
                  const Text('All caught up!', style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                  const SizedBox(height: 8),
                  const Text('No nudges or alerts right now.', style: TextStyle(color: Colors.grey)),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: () => fp.loadNudges(),
            child: ListView.builder(
              padding: const EdgeInsets.all(20),
              itemCount: active.length,
              itemBuilder: (context, index) {
                final nudge = active[index];
                return _buildNudgeCard(context, fp, nudge, index);
              },
            ),
          );
        },
      ),
    );
  }

  Widget _buildNudgeCard(BuildContext context, FinancialProvider fp, Map<String, dynamic> nudge, int index) {
    final severity = nudge['severity'] ?? 'info';
    final Color severityColor;
    final IconData severityIcon;

    switch (severity) {
      case 'critical':
        severityColor = Colors.red;
        severityIcon = Icons.error_rounded;
        break;
      case 'warning':
        severityColor = Colors.orange;
        severityIcon = Icons.warning_amber_rounded;
        break;
      default:
        severityColor = Colors.blue;
        severityIcon = Icons.info_rounded;
    }

    final category = nudge['category'] ?? '';
    String categoryLabel;
    IconData categoryIcon;
    switch (category) {
      case 'budget':
        categoryLabel = 'Budget';
        categoryIcon = Icons.pie_chart_rounded;
        break;
      case 'emergency_fund':
        categoryLabel = 'Emergency Fund';
        categoryIcon = Icons.shield_rounded;
        break;
      case 'savings':
        categoryLabel = 'Savings';
        categoryIcon = Icons.savings_rounded;
        break;
      case 'debt':
        categoryLabel = 'Debt';
        categoryIcon = Icons.trending_down_rounded;
        break;
      case 'subscription':
        categoryLabel = 'Subscription';
        categoryIcon = Icons.subscriptions_rounded;
        break;
      default:
        categoryLabel = category;
        categoryIcon = Icons.lightbulb_rounded;
    }

    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: severityColor.withValues(alpha: 0.2)),
        boxShadow: [BoxShadow(color: severityColor.withValues(alpha: 0.06), blurRadius: 16, offset: const Offset(0, 6))],
      ),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(8),
                  decoration: BoxDecoration(
                    color: severityColor.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Icon(severityIcon, color: severityColor, size: 20),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(nudge['title'] ?? '', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Color(0xFF1E293B))),
                      const SizedBox(height: 2),
                      Row(
                        children: [
                          Icon(categoryIcon, size: 14, color: Colors.grey),
                          const SizedBox(width: 4),
                          Text(categoryLabel, style: const TextStyle(color: Colors.grey, fontSize: 12, fontWeight: FontWeight.w600)),
                        ],
                      ),
                    ],
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                  decoration: BoxDecoration(
                    color: severityColor.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(severity.toUpperCase(), style: TextStyle(color: severityColor, fontSize: 10, fontWeight: FontWeight.bold)),
                ),
              ],
            ),
            const SizedBox(height: 14),
            Text(nudge['message'] ?? '', style: const TextStyle(color: Color(0xFF475569), height: 1.5, fontSize: 14)),
            const SizedBox(height: 16),
            Row(
              children: [
                if (nudge['action_label'] != null)
                  TextButton.icon(
                    onPressed: () {
                      final link = nudge['action_link'] ?? '';
                      if (link == 'budget') {
                        Navigator.pop(context);
                      } else if (link == 'portfolio') {
                        Navigator.pop(context);
                      }
                    },
                    icon: Icon(Icons.arrow_forward_rounded, size: 16, color: severityColor),
                    label: Text(nudge['action_label'], style: TextStyle(color: severityColor, fontWeight: FontWeight.bold)),
                    style: TextButton.styleFrom(
                      backgroundColor: severityColor.withValues(alpha: 0.08),
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                    ),
                  ),
                const Spacer(),
                IconButton(
                  onPressed: () => fp.dismissNudge(nudge['id']),
                  icon: const Icon(Icons.close_rounded, size: 20),
                  style: IconButton.styleFrom(backgroundColor: Colors.grey.shade100),
                  constraints: const BoxConstraints(minWidth: 36, minHeight: 36),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}
