import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:fl_chart/fl_chart.dart';
import '../providers/financial_provider.dart';

class FinancialScoreScreen extends StatelessWidget {
  const FinancialScoreScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<FinancialProvider>();
    final score = provider.finScoreVal;
    final loading = provider.isLoading;
    
    final months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    
    List<Map<String, dynamic>> history = [];
    try {
      history = provider.scoreHistory.map((e) {
        final dt = DateTime.parse(e['calculatedAt']?.toString() ?? '');
        return {'month': months[dt.month - 1], 'score': e['score'] ?? 0};
      }).toList();
    } catch (_) {
      history = [];
    }
    
    if (history.isEmpty) {
      history = [{'month': months[DateTime.now().month - 1], 'score': score}];
    } else if (history.length > 6) {
      history = history.sublist(history.length - 6);
    }

    return Scaffold(
      backgroundColor: const Color(0xFFF1F5F9),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios_new_rounded, color: Color(0xFF1E293B)),
          onPressed: () => Navigator.pop(context),
        ),
        title: const Text('Score Breakdown', style: TextStyle(color: Color(0xFF1E293B), fontWeight: FontWeight.bold)),
        centerTitle: true,
      ),
      body: loading
        ? const Center(child: CircularProgressIndicator(color: Color(0xFF6366F1)))
        : SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            // ── SCORE HERO ──────────────────────────────────────────────────
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(32),
              decoration: BoxDecoration(
                gradient: const LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [Color(0xFF042F2E), Color(0xFF6366F1)]
                ),
                borderRadius: BorderRadius.circular(28),
                boxShadow: [BoxShadow(color: const Color(0xFF6366F1).withValues(alpha: 0.3), blurRadius: 20, offset: const Offset(0, 8))],
              ),
              child: Column(
                children: [
                  const Text('Your Financial Score', style: TextStyle(color: Color(0xFF93C5FD), fontSize: 14, fontWeight: FontWeight.w600)),
                  const SizedBox(height: 12),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    crossAxisAlignment: CrossAxisAlignment.baseline,
                    textBaseline: TextBaseline.alphabetic,
                    children: [
                      Text('$score', style: const TextStyle(color: Colors.white, fontSize: 64, fontWeight: FontWeight.w900, height: 1.0)),
                      const Text('/100', style: TextStyle(color: Color(0xFF93C5FD), fontSize: 20, fontWeight: FontWeight.bold)),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Text(
                    score >= 80 ? '🌟 Excellent! You are financially healthy.' : score >= 60 ? '👍 Good, but room for improvement.' : '⚡ Keep working on your fundamentals.',
                    style: const TextStyle(color: Colors.white, fontSize: 14, fontWeight: FontWeight.w500),
                    textAlign: TextAlign.center,
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // ── SCORE TREND CHART ───────────────────────────────────────────
            Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(24),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Score Trend', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                  const SizedBox(height: 24),
                  SizedBox(
                    height: 200,
                    child: BarChart(
                      BarChartData(
                        alignment: BarChartAlignment.spaceAround,
                        maxY: 100,
                        barTouchData: BarTouchData(enabled: false),
                        titlesData: FlTitlesData(
                          show: true,
                          bottomTitles: AxisTitles(
                            sideTitles: SideTitles(
                              showTitles: true,
                              getTitlesWidget: (value, meta) {
                                if (value.toInt() >= 0 && value.toInt() < history.length) {
                                  return Padding(
                                    padding: const EdgeInsets.only(top: 8.0),
                                    child: Text(history[value.toInt()]['month'] as String, style: const TextStyle(color: Colors.grey, fontSize: 12, fontWeight: FontWeight.bold)),
                                  );
                                }
                                return const SizedBox.shrink();
                              },
                            ),
                          ),
                          leftTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
                          rightTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
                          topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
                        ),
                        borderData: FlBorderData(show: false),
                        gridData: FlGridData(
                          show: true,
                          drawVerticalLine: false,
                          horizontalInterval: 25,
                          getDrawingHorizontalLine: (value) => FlLine(color: Colors.grey.shade200, strokeWidth: 1),
                        ),
                        barGroups: List.generate(history.length, (i) {
                          final val = (history[i]['score'] as int).toDouble();
                          return BarChartGroupData(
                            x: i,
                            barRods: [
                              BarChartRodData(
                                toY: val,
                                width: 24,
                                color: val == 0 ? Colors.grey.shade200 : const Color(0xFF2DD4BF),
                                borderRadius: const BorderRadius.vertical(top: Radius.circular(6)),
                                backDrawRodData: BackgroundBarChartRodData(
                                  show: true,
                                  toY: 100,
                                  color: Colors.grey.shade50,
                                ),
                              )
                            ],
                          );
                        }),
                      ),
                    ),
                  ),
                ],
              ),
            ),
            const SizedBox(height: 24),

            // ── PILLARS BREAKDOWN ───────────────────────────────────────────
            const Align(
              alignment: Alignment.centerLeft,
              child: Text('Score Pillars', style: TextStyle(fontSize: 20, fontWeight: FontWeight.w900, color: Color(0xFF1E293B))),
            ),
            const SizedBox(height: 16),
            _buildPillarCard(
              title: 'Savings Rate (30 Points)',
              icon: Icons.savings_rounded,
              color: const Color(0xFF10B981),
              description: 'Measures how much of your income is saved. Target is 20%+.',
            ),
            const SizedBox(height: 12),
            _buildPillarCard(
              title: 'Debt-to-Income (30 Points)',
              icon: Icons.credit_card_rounded,
              color: const Color(0xFFF59E0B),
              description: 'Evaluates your monthly EMIs against income. Target is 0%.',
            ),
            const SizedBox(height: 12),
            _buildPillarCard(
              title: 'Emergency Fund (40 Points)',
              icon: Icons.health_and_safety_rounded,
              color: const Color(0xFF2DD4BF),
              description: 'Looks at liquid assets (Bank/FDs) vs 3 months of basic expenses.',
            ),
            const SizedBox(height: 40),
          ],
        ),
      ),
    );
  }

  Widget _buildPillarCard({required String title, required IconData icon, required Color color, required String description}) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: Colors.grey.shade100),
        boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.02), blurRadius: 10, offset: const Offset(0, 4))],
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(color: color.withValues(alpha: 0.1), shape: BoxShape.circle),
            child: Icon(icon, color: color, size: 24),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(title, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Color(0xFF1E293B))),
                const SizedBox(height: 4),
                Text(description, style: TextStyle(color: Colors.grey.shade600, fontSize: 13, height: 1.4)),
              ],
            ),
          )
        ],
      ),
    );
  }
}
