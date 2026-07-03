import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:fl_chart/fl_chart.dart';
import '../providers/financial_provider.dart';

class StatementInsightsScreen extends StatefulWidget {
  const StatementInsightsScreen({super.key});

  @override
  State<StatementInsightsScreen> createState() => _StatementInsightsScreenState();
}

class _StatementInsightsScreenState extends State<StatementInsightsScreen> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        title: const Text('Statement Insights'),
        actions: [
          Consumer<FinancialProvider>(
            builder: (ctx, fp, child) => fp.pendingRecurringSuggestions.isNotEmpty
                ? TextButton.icon(
                    onPressed: () => _showRecurringConfirmSheet(context, fp),
                    icon: const Icon(Icons.add_circle_outline, size: 18),
                    label: const Text('Review'),
                  )
                : const SizedBox(),
          ),
        ],
      ),
      body: Consumer<FinancialProvider>(
        builder: (context, fp, _) {
          final analysis = fp.lastAnalysisResult;
          if (analysis == null) {
            return const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(Icons.analytics_outlined, size: 80, color: Colors.grey),
                  SizedBox(height: 20),
                  Text('No analysis data available',
                      style: TextStyle(fontSize: 18, color: Colors.grey)),
                  SizedBox(height: 8),
                  Text('Upload a statement with analysis to see insights.',
                      style: TextStyle(color: Colors.grey)),
                ],
              ),
            );
          }

          final recurring = analysis['recurring_detections'] as List? ?? [];
          final subscriptions = analysis['subscription_candidates'] as List? ?? [];
          final lapsed = analysis['lapsed_subscriptions'] as List? ?? [];
          final patterns = analysis['spending_patterns'] as Map<String, dynamic>? ?? {};
          final nudgesCount = analysis['nudges_generated'] as int? ?? 0;

          return RefreshIndicator(
            onRefresh: () => fp.runFullAnalysis(),
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                _summaryHeader(context, recurring.length, subscriptions.length, nudgesCount),
                const SizedBox(height: 16),
                if (fp.pendingRecurringSuggestions.isNotEmpty) ...[
                  _sectionHeader('Recurring Payments Detected',
                      Icons.repeat, Colors.blue, fp.pendingRecurringSuggestions.length),
                  _recurringList(context, fp),
                  const SizedBox(height: 16),
                ],
                if (patterns.isNotEmpty) ...[
                  _sectionHeader('Spending Patterns', Icons.pie_chart, Colors.purple, null),
                  _spendingPatternsSection(patterns),
                  const SizedBox(height: 16),
                ],
                if (subscriptions.isNotEmpty) ...[
                  _sectionHeader('Subscription Candidates',
                      Icons.subscriptions, Colors.orange, subscriptions.length),
                  _subscriptionList(subscriptions),
                  const SizedBox(height: 16),
                ],
                if (lapsed.isNotEmpty) ...[
                  _sectionHeader('Lapsed Subscriptions',
                      Icons.cancel_schedule_send, Colors.red, lapsed.length),
                  _lapsedSubscriptionList(lapsed),
                  const SizedBox(height: 16),
                ],
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _summaryHeader(BuildContext context, int recurringCount, int subCount, int nudgesCount) {
    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Color(0xFF1E3A8A), Color(0xFF3B82F6)],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Analysis Complete',
              style: TextStyle(fontSize: 22, fontWeight: FontWeight.bold, color: Colors.white)),
          const SizedBox(height: 4),
          const Text('Key findings from your statement',
              style: TextStyle(color: Colors.white70, fontSize: 14)),
          const SizedBox(height: 20),
          Row(
            children: [
              _statBadge(Icons.repeat, '$recurringCount', 'Recurring', Colors.blue.shade200),
              const SizedBox(width: 12),
              _statBadge(Icons.subscriptions, '$subCount', 'Subscriptions', Colors.orange.shade200),
              const SizedBox(width: 12),
              _statBadge(Icons.lightbulb, '$nudgesCount', 'Nudges', Colors.green.shade200),
            ],
          ),
        ],
      ),
    );
  }

  Widget _statBadge(IconData icon, String value, String label, Color bgColor) {
    return Expanded(
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: bgColor.withValues(alpha: 0.2),
          borderRadius: BorderRadius.circular(12),
        ),
        child: Column(
          children: [
            Icon(icon, color: Colors.white, size: 22),
            const SizedBox(height: 6),
            Text(value,
                style: const TextStyle(
                    fontSize: 20, fontWeight: FontWeight.bold, color: Colors.white)),
            Text(label,
                style: const TextStyle(fontSize: 11, color: Colors.white70)),
          ],
        ),
      ),
    );
  }

  Widget _sectionHeader(String title, IconData icon, Color color, int? count) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        children: [
          Icon(icon, color: color, size: 22),
          const SizedBox(width: 8),
          Text(title,
              style:
                  const TextStyle(fontSize: 16, fontWeight: FontWeight.w600, color: Color(0xFF1E293B))),
          if (count != null) ...[
            const Spacer(),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
              decoration: BoxDecoration(
                color: color.withValues(alpha: 0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Text('$count found',
                  style: TextStyle(fontSize: 12, color: color, fontWeight: FontWeight.w600)),
            ),
          ],
        ],
      ),
    );
  }

  Widget _recurringList(BuildContext context, FinancialProvider fp) {
    final suggestions = fp.pendingRecurringSuggestions;
    return Column(
      children: suggestions.map((s) => _recurringCard(context, fp, s)).toList(),
    );
  }

  Widget _recurringCard(BuildContext context, FinancialProvider fp, Map<String, dynamic> s) {
    final confidence = (s['confidence'] as num?)?.toDouble() ?? 0;
    final pct = (confidence * 100).round();
    return Card(
      margin: const EdgeInsets.only(bottom: 10),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(14),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(s['name'] ?? '',
                      style: const TextStyle(
                          fontSize: 15, fontWeight: FontWeight.w600, color: Color(0xFF1E293B))),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                  decoration: BoxDecoration(
                    color: (confidence > 0.7 ? Colors.green : Colors.orange).withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(10),
                  ),
                  child: Text('$pct%',
                      style: TextStyle(
                          fontSize: 12,
                          fontWeight: FontWeight.bold,
                          color: confidence > 0.7 ? Colors.green : Colors.orange)),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              children: [
                _infoChip(Icons.currency_rupee, '₹${s['amount']}'),
                const SizedBox(width: 8),
                _infoChip(Icons.calendar_today, s['frequency'] ?? ''),
                const SizedBox(width: 8),
                _infoChip(Icons.repeat, '${s['occurrences']}x'),
              ],
            ),
            const SizedBox(height: 4),
            Row(
              children: [
                _infoChip(Icons.category, s['category'] ?? ''),
                const SizedBox(width: 8),
                _infoChip(Icons.inventory_2, s['bucket'] ?? ''),
              ],
            ),
            const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton.icon(
                  onPressed: () => fp.pendingRecurringSuggestions.removeWhere((x) => x['name'] == s['name']),
                  icon: const Icon(Icons.close, size: 16),
                  label: const Text('Ignore'),
                  style: TextButton.styleFrom(foregroundColor: Colors.grey),
                ),
                const SizedBox(width: 8),
                ElevatedButton.icon(
                  onPressed: () => _confirmRecurring(fp, s),
                  icon: const Icon(Icons.check, size: 16),
                  label: const Text('Add as Bill'),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF1E3A8A),
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    textStyle: const TextStyle(fontSize: 13),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  void _confirmRecurring(FinancialProvider fp, Map<String, dynamic> s) {
    showDialog(
      context: context,
      builder: (ctx) {
        final dueCtrl = TextEditingController(text: '1');
        return AlertDialog(
          title: const Text('Add Recurring Bill'),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Text('${s['name']} — ₹${s['amount']} / ${s['frequency']}'),
              const SizedBox(height: 12),
              TextField(
                controller: dueCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Due Day (1-31)',
                  border: OutlineInputBorder(),
                ),
              ),
            ],
          ),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(ctx).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(ctx).pop();
                final day = int.tryParse(dueCtrl.text) ?? 1;
                fp.confirmRecurringSuggestion(s, dueDay: day.clamp(1, 31));
              },
              child: const Text('Confirm'),
            ),
          ],
        );
      },
    );
  }

  Widget _infoChip(IconData icon, String label) {
    return Row(
      mainAxisSize: MainAxisSize.min,
      children: [
        Icon(icon, size: 14, color: Colors.grey),
        const SizedBox(width: 3),
        Text(label, style: const TextStyle(fontSize: 12, color: Colors.grey)),
      ],
    );
  }

  Widget _spendingPatternsSection(Map<String, dynamic> patterns) {
    final breakdown = patterns['category_breakdown'] as List? ?? [];
    final trends = patterns['category_trends'] as List? ?? [];
    final forecast = patterns['forecast'] as Map<String, dynamic>? ?? {};
    final weekendBoost = (patterns['weekend_boost_pct'] as num?)?.toDouble() ?? 0;
    final anomalies = patterns['anomalies'] as List? ?? [];
    final topMerchants = patterns['top_merchants'] as List? ?? [];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        if (breakdown.isNotEmpty) ...[
          _subSection('Category Breakdown', Icons.pie_chart_outline),
          _categoryBreakdownChart(breakdown),
          const SizedBox(height: 12),
        ],
        if (trends.isNotEmpty) ...[
          _subSection('Category Trends (MoM)', Icons.trending_up),
          ...trends.where((t) {
            final change = (t['change_pct'] as num?)?.toDouble() ?? 0;
            return change.abs() > 10;
          }).take(5).map((t) => _trendRow(t)),
          const SizedBox(height: 12),
        ],
        if (forecast.isNotEmpty) ...[
          _subSection('Monthly Forecast', Icons.query_stats),
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.white, borderRadius: BorderRadius.circular(12)),
            child: Row(
              children: [
                _forecastItem('Projected', '₹${_fmt(forecast['projected_total'])}', Colors.blue),
                _forecastItem('Spent', '₹${_fmt(forecast['current_spent'])}', Colors.orange),
                _forecastItem('Avg', '₹${_fmt(forecast['average_monthly'])}', Colors.grey),
              ],
            ),
          ),
          const SizedBox(height: 12),
        ],
        if (weekendBoost > 30) ...[
          _subSection('Weekend vs Weekday', Icons.calendar_view_week),
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.white, borderRadius: BorderRadius.circular(12)),
            child: Row(
              children: [
                Icon(weekendBoost > 0 ? Icons.trending_up : Icons.trending_down,
                    color: weekendBoost > 50 ? Colors.red : Colors.green),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    'Weekend spending is ${weekendBoost.abs().round()}% ${weekendBoost > 0 ? 'higher' : 'lower'} than weekdays.',
                    style: const TextStyle(fontSize: 14),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 12),
        ],
        if (anomalies.isNotEmpty) ...[
          _subSection('Unusual Transactions', Icons.warning_amber),
          ...anomalies.take(3).map((a) => _anomalyRow(a)),
          const SizedBox(height: 12),
        ],
        if (topMerchants.isNotEmpty) ...[
          _subSection('Top Merchants (This Month)', Icons.store),
          ...topMerchants.take(5).map((m) => _merchantRow(m)),
        ],
      ],
    );
  }

  Widget _categoryBreakdownChart(List breakdown) {
    if (breakdown.isEmpty) return const SizedBox();

    final colors = [
      Colors.blue, Colors.orange, Colors.purple, Colors.green,
      Colors.red, Colors.teal, Colors.indigo, Colors.amber,
      Colors.pink, Colors.cyan,
    ];

    double total = 0;
    for (var b in breakdown) {
      total += (b['amount'] as num?)?.toDouble() ?? 0;
    }

    if (total == 0) return const SizedBox();

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
            offset: const Offset(0, 4),
          )
        ],
      ),
      child: Column(
        children: [
          SizedBox(
            height: 200,
            child: PieChart(
              PieChartData(
                sectionsSpace: 2,
                centerSpaceRadius: 40,
                sections: breakdown.asMap().entries.map((e) {
                  final index = e.key;
                  final item = e.value;
                  final amount = (item['amount'] as num?)?.toDouble() ?? 0;
                  final pct = (amount / total * 100);
                  
                  return PieChartSectionData(
                    color: colors[index % colors.length],
                    value: amount,
                    title: '${pct.toStringAsFixed(1)}%',
                    radius: 50,
                    titleStyle: const TextStyle(
                      fontSize: 11,
                      fontWeight: FontWeight.bold,
                      color: Colors.white,
                    ),
                  );
                }).toList(),
              ),
            ),
          ),
          const SizedBox(height: 16),
          Wrap(
            spacing: 12,
            runSpacing: 8,
            children: breakdown.asMap().entries.map((e) {
              final index = e.key;
              final item = e.value;
              final amount = (item['amount'] as num?)?.toDouble() ?? 0;
              final catName = item['category']?.toString() ?? 'Misc';
              return Row(
                mainAxisSize: MainAxisSize.min,
                children: [
                  Container(
                    width: 12, height: 12,
                    decoration: BoxDecoration(
                      color: colors[index % colors.length],
                      shape: BoxShape.circle,
                    ),
                  ),
                  const SizedBox(width: 4),
                  Text('$catName (₹${_fmt(amount)})', style: const TextStyle(fontSize: 12)),
                ],
              );
            }).toList(),
          ),
        ],
      ),
    );
  }

  Widget _subSection(String title, IconData icon) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 6, top: 4),
      child: Row(
        children: [
          Icon(icon, size: 16, color: Colors.grey),
          const SizedBox(width: 6),
          Text(title,
              style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500, color: Color(0xFF475569))),
        ],
      ),
    );
  }

  Widget _trendRow(Map<String, dynamic> t) {
    final change = (t['change_pct'] as num?)?.toDouble() ?? 0;
    final isUp = change > 0;
    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      child: ListTile(
        dense: true,
        leading: Icon(
          isUp ? Icons.arrow_upward : Icons.arrow_downward,
          color: isUp ? Colors.red : Colors.green,
          size: 18,
        ),
        title: Text(t['category'] ?? '',
            style: const TextStyle(fontSize: 14, fontWeight: FontWeight.w500)),
        trailing: Text(
          '${isUp ? '+' : ''}${change.toStringAsFixed(1)}%',
          style: TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.bold,
            color: isUp ? Colors.red : Colors.green,
          ),
        ),
      ),
    );
  }

  Widget _forecastItem(String label, String value, Color color) {
    return Expanded(
      child: Column(
        children: [
          Text(value,
              style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: color)),
          Text(label, style: const TextStyle(fontSize: 12, color: Colors.grey)),
        ],
      ),
    );
  }

  Widget _anomalyRow(Map<String, dynamic> a) {
    return Card(
      margin: const EdgeInsets.only(bottom: 4),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
      child: ListTile(
        dense: true,
        leading: const Icon(Icons.error_outline, color: Colors.orange, size: 20),
        title: Text(a['name'] ?? '', style: const TextStyle(fontSize: 14)),
        subtitle: Text('₹${_fmt(a['amount'])} — ${a['category']}',
            style: const TextStyle(fontSize: 12)),
        trailing: Text(a['reason'] ?? '',
            style: const TextStyle(fontSize: 12, color: Colors.orange)),
      ),
    );
  }

  Widget _merchantRow(Map<String, dynamic> m) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 3, horizontal: 4),
      child: Row(
        children: [
          Icon(Icons.store, size: 16, color: Colors.grey.shade500),
          const SizedBox(width: 8),
          Expanded(
            child: Text(m['name'] ?? '',
                style: const TextStyle(fontSize: 13, color: Color(0xFF334155))),
          ),
          Text('₹${_fmt(m['amount'])}',
              style: const TextStyle(
                  fontSize: 13, fontWeight: FontWeight.w600, color: Color(0xFF1E293B))),
        ],
      ),
    );
  }

  Widget _subscriptionList(List subscriptions) {
    return Column(
      children: subscriptions.take(5).map((s) => _subscriptionCard(s)).toList(),
    );
  }

  Widget _subscriptionCard(Map<String, dynamic> s) {
    final isNew = s['is_new'] == true;
    return Card(
      margin: const EdgeInsets.only(bottom: 6),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
      child: ListTile(
        dense: true,
        leading: Icon(
          isNew ? Icons.fiber_new : Icons.subscriptions,
          color: isNew ? Colors.green : Colors.orange,
          size: 22,
        ),
        title: Text(s['name'] ?? '', style: const TextStyle(fontSize: 14)),
        subtitle: Text(
          '₹${_fmt(s['monthly_cost'])}/mo · ₹${_fmt(s['annual_cost'])}/yr',
          style: const TextStyle(fontSize: 12),
        ),
        trailing: s['savings_opportunity'] != null && (s['savings_opportunity'] as num) > 0
            ? Text('Save ₹${_fmt(s['savings_opportunity'])}',
                style: const TextStyle(fontSize: 12, color: Colors.green, fontWeight: FontWeight.w600))
            : null,
      ),
    );
  }

  Widget _lapsedSubscriptionList(List lapsed) {
    return Column(
      children: lapsed.map((s) => _lapsedCard(s)).toList(),
    );
  }

  Widget _lapsedCard(Map<String, dynamic> s) {
    return Card(
      margin: const EdgeInsets.only(bottom: 6),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(10)),
      child: ListTile(
        dense: true,
        leading: const Icon(Icons.cancel_schedule_send, color: Colors.red, size: 22),
        title: Text(s['name'] ?? '', style: const TextStyle(fontSize: 14)),
        subtitle: Text(
          '₹${_fmt(s['annual_cost'])}/yr · ${s['status']}',
          style: const TextStyle(fontSize: 12),
        ),
        trailing: Text('Save ₹${_fmt(s['savings_opportunity'])}',
            style: const TextStyle(fontSize: 12, color: Colors.green, fontWeight: FontWeight.w600)),
      ),
    );
  }

  void _showRecurringConfirmSheet(BuildContext context, FinancialProvider fp) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (ctx) => DraggableScrollableSheet(
        initialChildSize: 0.7,
        minChildSize: 0.3,
        maxChildSize: 0.9,
        expand: false,
        builder: (_, scrollCtrl) => Column(
          children: [
            Padding(
              padding: const EdgeInsets.all(16),
              child: Row(
                children: [
                  const Text('Review Recurring',
                      style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                  const Spacer(),
                  Text('${fp.pendingRecurringSuggestions.length} pending',
                      style: const TextStyle(color: Colors.grey)),
                ],
              ),
            ),
            const Divider(height: 1),
            Expanded(
              child: ListView(
                controller: scrollCtrl,
                padding: const EdgeInsets.all(16),
                children: fp.pendingRecurringSuggestions
                    .map((s) => _recurringCard(ctx, fp, s))
                    .toList(),
              ),
            ),
          ],
        ),
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
