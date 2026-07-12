import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:fl_chart/fl_chart.dart';
import '../providers/financial_provider.dart';
import '../models/financial_models.dart';
import '../utils/ui_utils.dart';
import 'salary_breakup_screen.dart';

class SalaryHistoryScreen extends StatefulWidget {
  final LocalIncome income;

  const SalaryHistoryScreen({
    super.key,
    required this.income,
  });

  @override
  State<SalaryHistoryScreen> createState() => _SalaryHistoryScreenState();
}

class _SalaryHistoryScreenState extends State<SalaryHistoryScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final provider = Provider.of<FinancialProvider>(context, listen: false);
      provider.reloadSalaryDetails(widget.income.id);
      provider.loadSalaryGrowth();
    });
  }

  Widget _buildGrowthChart(LocalSalaryGrowth growth) {
    if (growth.growthData.isEmpty) {
      return const SizedBox(
        height: 180,
        child: Center(child: Text('Add salary data to view growth graph', style: TextStyle(color: Colors.grey))),
      );
    }

    final data = growth.growthData;
    List<FlSpot> spots = [];
    double minYear = data.first.year.toDouble();
    double maxYear = data.last.year.toDouble();
    double maxGross = 0.0;

    for (int i = 0; i < data.length; i++) {
      spots.add(FlSpot(data[i].year.toDouble(), data[i].grossAnnual));
      if (data[i].grossAnnual > maxGross) maxGross = data[i].grossAnnual;
    }

    // Set nice round division for Y labels (thousands/Lakhs)
    double interval = maxGross > 0 ? (maxGross / 4).clamp(100000.0, double.infinity) : 500000.0;

    return Container(
      height: 220,
      padding: const EdgeInsets.only(right: 16, top: 12, bottom: 8),
      child: LineChart(
        LineChartData(
          gridData: FlGridData(
            show: true,
            drawVerticalLine: false,
            getDrawingHorizontalLine: (value) => FlLine(color: Colors.grey.shade100, strokeWidth: 1),
          ),
          titlesData: FlTitlesData(
            show: true,
            rightTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
            topTitles: const AxisTitles(sideTitles: SideTitles(showTitles: false)),
            bottomTitles: AxisTitles(
              sideTitles: SideTitles(
                showTitles: true,
                reservedSize: 26,
                interval: 1,
                getTitlesWidget: (value, meta) {
                  return Padding(
                    padding: const EdgeInsets.only(top: 8.0),
                    child: Text(
                      value.toInt().toString(),
                      style: TextStyle(color: Colors.grey.shade600, fontWeight: FontWeight.bold, fontSize: 10),
                    ),
                  );
                },
              ),
            ),
            leftTitles: AxisTitles(
              sideTitles: SideTitles(
                showTitles: true,
                interval: interval,
                reservedSize: 45,
                getTitlesWidget: (value, meta) {
                  String valStr;
                  if (value >= 100000) {
                    valStr = '${(value / 100000.0).toStringAsFixed(1)}L';
                  } else {
                    valStr = '${(value / 1000.0).toStringAsFixed(0)}K';
                  }
                  return Text(
                    valStr,
                    style: TextStyle(color: Colors.grey.shade600, fontSize: 10, fontWeight: FontWeight.bold),
                  );
                },
              ),
            ),
          ),
          borderData: FlBorderData(show: false),
          minX: minYear,
          maxX: maxYear,
          minY: 0,
          maxY: maxGross * 1.15,
          lineBarsData: [
            LineChartBarData(
              spots: spots,
              isCurved: true,
              color: const Color(0xFF0D9488),
              barWidth: 4,
              isStrokeCapRound: true,
              dotData: const FlDotData(show: true),
              belowBarData: BarAreaData(
                show: true,
                color: const Color(0xFF0D9488).withValues(alpha: 0.1),
              ),
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8F9FA),
      appBar: AppBar(
        title: Text('${widget.income.label} History'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (_) => SalaryBreakupScreen(income: widget.income),
                ),
              );
            },
          ),
        ],
      ),
      body: Consumer<FinancialProvider>(
        builder: (context, provider, child) {
          final details = provider.salaryDetails;
          final growth = provider.salaryGrowth;

          return RefreshIndicator(
            onRefresh: () async {
              await provider.reloadSalaryDetails(widget.income.id);
              await provider.loadSalaryGrowth();
            },
            child: SingleChildScrollView(
              physics: const AlwaysScrollableScrollPhysics(),
              padding: const EdgeInsets.all(20),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Growth chart card
                  if (growth != null && growth.growthData.isNotEmpty) ...[
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(20),
                        boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 10)],
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              const Text('Salary Growth Trend', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                              if (growth.avgGrowthPct > 0)
                                Container(
                                  padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                                  decoration: BoxDecoration(
                                    color: const Color(0xFF10B981).withValues(alpha: 0.1),
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                  child: Text(
                                    'Avg YoY: +${growth.avgGrowthPct.toStringAsFixed(1)}%',
                                    style: const TextStyle(color: Color(0xFF10B981), fontWeight: FontWeight.bold, fontSize: 11),
                                  ),
                                ),
                            ],
                          ),
                          const SizedBox(height: 16),
                          _buildGrowthChart(growth),
                        ],
                      ),
                    ),
                    const SizedBox(height: 24),
                  ],

                  // History list header
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Salary History', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                      TextButton.icon(
                        onPressed: () {
                          Navigator.push(
                            context,
                            MaterialPageRoute(
                              builder: (_) => SalaryBreakupScreen(income: widget.income),
                            ),
                          );
                        },
                        icon: const Icon(Icons.add, size: 18, color: Color(0xFF0D9488)),
                        label: const Text('Add Entry', style: TextStyle(color: Color(0xFF0D9488), fontWeight: FontWeight.bold)),
                      ),
                    ],
                  ),
                  const SizedBox(height: 10),

                  if (details.isEmpty)
                    Container(
                      padding: const EdgeInsets.all(32),
                      width: double.infinity,
                      decoration: BoxDecoration(
                        color: Colors.white,
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: const Column(
                        children: [
                          Icon(Icons.payments_outlined, size: 48, color: Colors.grey),
                          SizedBox(height: 12),
                          Text('No history records found.', style: TextStyle(fontWeight: FontWeight.bold)),
                          Text('Tap Add Entry to track your packages.', style: TextStyle(color: Colors.grey, fontSize: 12)),
                        ],
                      ),
                    )
                  else
                    ...details.map((sd) {
                      final company = sd.companyName ?? 'Unknown Company';
                      final yearStr = sd.toYear != null ? '${sd.fromYear} - ${sd.toYear}' : '${sd.fromYear} - Present';
                      final annualStr = sd.grossAnnual != null ? '₹${sd.grossAnnual!.toStringAsFixed(0)} /yr' : '₹0 /yr';
                      final monthlyStr = sd.monthlyInHand != null ? '₹${sd.monthlyInHand!.toStringAsFixed(0)} /mo in-hand' : 'Estimate unavailable';

                      // Find growth pct for this specific item from the growth model
                      double? growthPct;
                      if (growth != null) {
                        try {
                          growthPct = growth.growthData.firstWhere((x) => x.year == sd.fromYear).growthPct;
                        } catch (_) {}
                      }

                      return Container(
                        margin: const EdgeInsets.only(bottom: 12),
                        decoration: BoxDecoration(
                          color: Colors.white,
                          borderRadius: BorderRadius.circular(16),
                          border: sd.isCurrent ? Border.all(color: const Color(0xFF0D9488), width: 1.5) : null,
                          boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 8)],
                        ),
                        child: Material(
                          color: Colors.transparent,
                          child: InkWell(
                            borderRadius: BorderRadius.circular(16),
                            onTap: () {
                              Navigator.push(
                                  context,
                                  MaterialPageRoute(
                                      builder: (_) => SalaryBreakupScreen(income: widget.income, salaryDetail: sd)));
                            },
                            child: Padding(
                              padding: const EdgeInsets.all(16),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                    children: [
                                      Expanded(
                                        child: Text(company, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15, color: Color(0xFF1E293B))),
                                      ),
                                      if (sd.isCurrent)
                                        Container(
                                          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                                          decoration: BoxDecoration(
                                            color: const Color(0xFF0D9488).withValues(alpha: 0.1),
                                            borderRadius: BorderRadius.circular(20),
                                          ),
                                          child: const Text('Current', style: TextStyle(color: Color(0xFF0D9488), fontWeight: FontWeight.bold, fontSize: 10)),
                                        )
                                      else
                                        GestureDetector(
                                          onTap: () async {
                                            UiUtils.showSnack(context, 'Setting current salary detail...');
                                            await provider.setCurrentSalaryDetail(widget.income.id, sd.id);
                                          },
                                          child: Text('Make Current', style: TextStyle(color: Colors.blue.shade700, fontWeight: FontWeight.bold, fontSize: 12)),
                                        )
                                    ],
                                  ),
                                  const SizedBox(height: 4),
                                  Text(yearStr, style: const TextStyle(color: Colors.grey, fontSize: 12)),
                                  const SizedBox(height: 12),
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                    children: [
                                      Column(
                                        crossAxisAlignment: CrossAxisAlignment.start,
                                        children: [
                                          Text(annualStr, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                                          Text(monthlyStr, style: const TextStyle(color: Colors.grey, fontSize: 11)),
                                        ],
                                      ),
                                      if (growthPct != null && growthPct > 0)
                                        Row(
                                          children: [
                                            const Icon(Icons.arrow_upward_rounded, color: Color(0xFF10B981), size: 16),
                                            const SizedBox(width: 2),
                                            Text(
                                              '+${growthPct.toStringAsFixed(1)}%',
                                              style: const TextStyle(color: Color(0xFF10B981), fontWeight: FontWeight.bold, fontSize: 14),
                                            ),
                                          ],
                                        ),
                                    ],
                                  ),
                                  const SizedBox(height: 8),
                                  const Divider(),
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.end,
                                    children: [
                                      GestureDetector(
                                        onTap: () {
                                          Navigator.push(
                                              context,
                                              MaterialPageRoute(
                                                  builder: (_) => SalaryBreakupScreen(income: widget.income, salaryDetail: sd)));
                                        },
                                        child: const Text('Edit Details', style: TextStyle(color: Colors.grey, fontSize: 12, fontWeight: FontWeight.bold)),
                                      ),
                                      const SizedBox(width: 16),
                                      GestureDetector(
                                        onTap: () async {
                                          UiUtils.showDeleteBottomSheet(context, company, () async {
                                            await provider.deleteSalaryDetail(widget.income.id, sd.id);
                                          });
                                        },
                                        child: const Text('Delete', style: TextStyle(color: Colors.redAccent, fontSize: 12, fontWeight: FontWeight.bold)),
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      );
                    }),
                ],
              ),
            ),
          );
        },
      ),
    );
  }
}
