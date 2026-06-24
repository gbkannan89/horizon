import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:fl_chart/fl_chart.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class LifecycleScreen extends StatefulWidget {
  const LifecycleScreen({super.key});

  @override
  State<LifecycleScreen> createState() => _LifecycleScreenState();
}

class _LifecycleScreenState extends State<LifecycleScreen> {
  double _currentAge = 30;
  double _retirementAge = 60;
  double _lifeExpectancy = 100;
  double _inflationRate = 6.0;
  double _expectedReturn = 12.0;
  final _expensesCtrl = TextEditingController();
  Map<String, dynamic>? _result;
  bool _isLoading = false;

  @override
  void dispose() {
    _expensesCtrl.dispose();
    super.dispose();
  }

  Future<void> _runSimulation() async {
    setState(() => _isLoading = true);
    try {
      final fp = context.read<FinancialProvider>();
      final result = await fp.lifecycleSimulate({
        'current_age': _currentAge.toInt(),
        'retirement_age': _retirementAge.toInt(),
        'life_expectancy': _lifeExpectancy.toInt(),
        'inflation_rate': _inflationRate,
        'expected_return': _expectedReturn,
        'monthly_expenses': double.tryParse(_expensesCtrl.text) ?? 0,
      });
      setState(() => _result = result);
    } catch (e) {
      UiUtils.showSnack(context, 'Simulation failed: $e', isError: true);
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(title: const Text('Lifecycle Simulator')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            _buildParamsCard(),
            const SizedBox(height: 20),
            if (_result != null) _buildResults(),
            const SizedBox(height: 60),
          ],
        ),
      ),
    );
  }

  Widget _buildParamsCard() {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(24), boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 20)]),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [
          Container(padding: const EdgeInsets.all(10), decoration: BoxDecoration(color: const Color(0xFF6B46C1).withValues(alpha: 0.1), borderRadius: BorderRadius.circular(12)), child: const Icon(Icons.tune_rounded, color: Color(0xFF6B46C1), size: 20)),
          const SizedBox(width: 12),
          const Text('Parameters', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: Color(0xFF1E293B))),
        ]),
        const SizedBox(height: 24),
        _slider('Current Age', _currentAge, 18, 80, Icons.person_outlined, (v) => setState(() => _currentAge = v)),
        _slider('Retirement Age', _retirementAge, 40, 80, Icons.beach_access_outlined, (v) => setState(() => _retirementAge = v)),
        _slider('Life Expectancy', _lifeExpectancy, 70, 110, Icons.favorite_outline_rounded, (v) => setState(() => _lifeExpectancy = v)),
        _slider('Inflation Rate', _inflationRate, 0, 20, Icons.trending_up_rounded, (v) => setState(() => _inflationRate = v), suffix: '%'),
        _slider('Expected Return', _expectedReturn, 0, 30, Icons.monetization_on_outlined, (v) => setState(() => _expectedReturn = v), suffix: '%'),
        const SizedBox(height: 16),
        TextField(
          controller: _expensesCtrl,
          keyboardType: TextInputType.number,
          decoration: InputDecoration(
            labelText: 'Monthly Expenses (₹)',
            hintText: 'Leave empty to use your bills total',
            prefixText: '₹ ',
            prefixStyle: const TextStyle(fontWeight: FontWeight.w700, color: Color(0xFF6B46C1)),
            filled: true, fillColor: Colors.grey.shade50,
            border: OutlineInputBorder(borderRadius: BorderRadius.circular(14), borderSide: BorderSide.none),
          ),
        ),
        const SizedBox(height: 24),
        SizedBox(width: double.infinity, height: 52, child: ElevatedButton(
          onPressed: _isLoading ? null : _runSimulation,
          style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF6B46C1), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)), elevation: 0),
          child: _isLoading
            ? const SizedBox(height: 22, width: 22, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
            : const Text('Run Simulation', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
        )),
      ]),
    );
  }

  Widget _slider(String label, double value, double min, double max, IconData icon, ValueChanged<double> onChanged, {String suffix = ''}) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 16),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
          Row(children: [Icon(icon, size: 16, color: Colors.grey), const SizedBox(width: 6), Text(label, style: const TextStyle(fontSize: 13, fontWeight: FontWeight.w600, color: Colors.grey))]),
          Text('${value.toInt()}$suffix', style: const TextStyle(fontWeight: FontWeight.w900, fontSize: 16, color: Color(0xFF6B46C1))),
        ]),
        SliderTheme(
          data: SliderTheme.of(context).copyWith(
            activeTrackColor: const Color(0xFF6B46C1), thumbColor: const Color(0xFF6B46C1),
            inactiveTrackColor: const Color(0xFF6B46C1).withValues(alpha: 0.12),
            overlayColor: const Color(0xFF6B46C1).withValues(alpha: 0.12),
          ),
          child: Slider(value: value, min: min, max: max, onChanged: onChanged),
        ),
      ]),
    );
  }

  Widget _buildResults() {
    final proj = (_result!['projection'] as List<dynamic>?) ?? [];
    final milestones = (_result!['milestones'] as List<dynamic>?) ?? [];

    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(24), boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 20)]),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [
          Container(padding: const EdgeInsets.all(10), decoration: BoxDecoration(color: const Color(0xFF10B981).withValues(alpha: 0.1), borderRadius: BorderRadius.circular(12)), child: const Icon(Icons.auto_awesome_rounded, color: Color(0xFF10B981), size: 20)),
          const SizedBox(width: 12),
          const Text('Projection to Age 100', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: Color(0xFF1E293B))),
        ]),
        const SizedBox(height: 20),
        if (proj.length > 1) ...[
          SizedBox(
            height: 220,
            child: LineChart(LineChartData(
              gridData: FlGridData(show: true, drawVerticalLine: false, getDrawingHorizontalLine: (v) => FlLine(color: Colors.grey.shade100, strokeWidth: 1)),
              titlesData: FlTitlesData(
                leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, reservedSize: 50, getTitlesWidget: (v, _) => Text('₹${(v / 100000).toInt()}L', style: const TextStyle(fontSize: 10, color: Colors.grey)))),
                bottomTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, interval: 10, getTitlesWidget: (v, _) => Text('${v.toInt()}', style: const TextStyle(fontSize: 10, color: Colors.grey)))),
                rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
              ),
              borderData: FlBorderData(show: false),
              lineBarsData: [
                LineChartBarData(
                  spots: proj.asMap().entries.map((e) {
                    final age = (e.value['age'] ?? 0).toDouble();
                    final nw = (e.value['netWorth'] ?? 0).toDouble();
                    return FlSpot(age, nw);
                  }).toList(),
                  isCurved: true,
                  gradient: const LinearGradient(colors: [Color(0xFF6B46C1), Color(0xFF8B5CF6)]),
                  barWidth: 3, isStrokeCapRound: true, dotData: FlDotData(show: false),
                  belowBarData: BarAreaData(show: true, gradient: LinearGradient(colors: [Color(0xFF6B46C1).withValues(alpha: 0.15), Color(0xFF8B5CF6).withValues(alpha: 0.05)])),
                ),
              ],
            )),
          ),
          const SizedBox(height: 16),
        ],
        if (milestones.isNotEmpty) ...[
          const Divider(),
          const SizedBox(height: 12),
          const Text('Key Milestones', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
          const SizedBox(height: 12),
          ...milestones.map((m) {
            final type = m['type'] ?? '';
            final desc = m['description'] ?? '';
            final isValley = type == 'liquidity_valley';
            final isRetirement = type == 'retirement';
            return Padding(
              padding: const EdgeInsets.only(bottom: 10),
              child: Row(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Container(
                  padding: const EdgeInsets.all(6),
                  decoration: BoxDecoration(
                    color: isValley ? Colors.red.withValues(alpha: 0.1) : isRetirement ? const Color(0xFF059669).withValues(alpha: 0.1) : Colors.blue.withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Icon(
                    isValley ? Icons.warning_amber_rounded : isRetirement ? Icons.beach_access_rounded : Icons.check_circle_outline,
                    color: isValley ? Colors.red : isRetirement ? const Color(0xFF059669) : Colors.blue,
                    size: 18,
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(child: Text(desc, style: TextStyle(fontSize: 13, color: isValley ? Colors.red.shade700 : Colors.grey.shade800, fontWeight: FontWeight.w600))),
              ]),
            );
          }),
        ],
      ]),
    );
  }
}
