import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:fl_chart/fl_chart.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class SimulationScreen extends StatefulWidget {
  const SimulationScreen({super.key});

  @override
  State<SimulationScreen> createState() => _SimulationScreenState();
}

class _SimulationScreenState extends State<SimulationScreen> {
  String _selectedType = 'purchase';
  bool _isLoading = false;
  Map<String, dynamic>? _result;

  final _nameCtrl = TextEditingController();
  final _amountCtrl = TextEditingController();
  final _rateCtrl = TextEditingController();
  final _tenureCtrl = TextEditingController();
  final _sipCtrl = TextEditingController();
  final _incomeCtrl = TextEditingController();

  @override
  void dispose() {
    _nameCtrl.dispose();
    _amountCtrl.dispose();
    _rateCtrl.dispose();
    _tenureCtrl.dispose();
    _sipCtrl.dispose();
    _incomeCtrl.dispose();
    super.dispose();
  }

  Future<void> _runSimulation() async {
    if (_nameCtrl.text.isEmpty) {
      UiUtils.showSnack(context, 'Give this scenario a name', isError: true);
      return;
    }

    setState(() => _isLoading = true);
    try {
      final fp = context.read<FinancialProvider>();
      final Map<String, dynamic> params = {'scenario_name': _nameCtrl.text, 'scenario_type': _selectedType};

      switch (_selectedType) {
        case 'purchase':
          params['target_amount'] = double.tryParse(_amountCtrl.text) ?? 0;
          params['interest_rate'] = double.tryParse(_rateCtrl.text) ?? 9.0;
          params['tenure_months'] = (double.tryParse(_tenureCtrl.text) ?? 60).toInt();
          break;
        case 'loan':
          params['loan_amount'] = double.tryParse(_amountCtrl.text) ?? 0;
          params['interest_rate'] = double.tryParse(_rateCtrl.text) ?? 10.0;
          params['tenure_months'] = (double.tryParse(_tenureCtrl.text) ?? 60).toInt();
          break;
        case 'sip_increase':
          params['sip_amount'] = double.tryParse(_sipCtrl.text) ?? 0;
          break;
        case 'income_change':
          params['income_change'] = double.tryParse(_incomeCtrl.text) ?? 0;
          break;
        case 'lumpsum':
          params['target_amount'] = double.tryParse(_amountCtrl.text) ?? 0;
          break;
      }

      final result = await fp.runSimulation(params);
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
      appBar: AppBar(title: const Text('Scenario Simulator')),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          children: [
            _buildScenarioTypeSelector(),
            const SizedBox(height: 24),
            _buildInputForm(),
            const SizedBox(height: 24),
            if (_result != null)
              _buildResults()
            else
              Container(
                padding: const EdgeInsets.all(32),
                decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(24), boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 20)]),
                child: Column(children: [
                  Icon(Icons.query_stats_outlined, size: 48, color: Colors.grey.shade300),
                  const SizedBox(height: 12),
                  Text('No simulations yet', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Colors.grey.shade500)),
                  const SizedBox(height: 4),
                  Text('Fill in the parameters above and tap "Run Simulation"', style: TextStyle(fontSize: 13, color: Colors.grey.shade400)),
                ]),
              ),
            const SizedBox(height: 60),
          ],
        ),
      ),
    );
  }

  Widget _buildScenarioTypeSelector() {
    final types = [
      ('purchase', Icons.shopping_cart_rounded, 'Purchase'),
      ('loan', Icons.account_balance_rounded, 'Loan'),
      ('sip_increase', Icons.trending_up_rounded, 'SIP Boost'),
      ('income_change', Icons.attach_money_rounded, 'Income'),
      ('lumpsum', Icons.monetization_on_rounded, 'Lumpsum'),
    ];

    return Container(
      padding: const EdgeInsets.all(4),
      decoration: BoxDecoration(color: Colors.grey.shade100, borderRadius: BorderRadius.circular(16)),
      child: Row(
        children: types.map((t) {
          final selected = _selectedType == t.$1;
          return Expanded(
            child: GestureDetector(
              onTap: () => setState(() {
                _selectedType = t.$1;
                _result = null;
              }),
              child: AnimatedContainer(
                duration: const Duration(milliseconds: 200),
                padding: const EdgeInsets.symmetric(vertical: 12),
                decoration: BoxDecoration(
                  color: selected ? const Color(0xFF1E3A8A) : Colors.transparent,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Column(
                  children: [
                    Icon(t.$2, color: selected ? Colors.white : Colors.grey, size: 20),
                    const SizedBox(height: 4),
                    Text(t.$3, style: TextStyle(
                      color: selected ? Colors.white : Colors.grey,
                      fontSize: 11, fontWeight: FontWeight.bold,
                    )),
                  ],
                ),
              ),
            ),
          );
        }).toList(),
      ),
    );
  }

  Widget _buildInputForm() {
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(24), boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 20)]),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text('Scenario Name', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
          const SizedBox(height: 8),
          TextField(
            controller: _nameCtrl,
            decoration: const InputDecoration(hintText: 'e.g., Buy a car in 2026'),
          ),
          const SizedBox(height: 20),
          if (_selectedType == 'purchase' || _selectedType == 'loan' || _selectedType == 'lumpsum') ...[
            _buildField('Amount (₹)', _amountCtrl, 'e.g., 1500000'),
            if (_selectedType != 'lumpsum') ...[
              const SizedBox(height: 16),
              _buildField('Interest Rate (%)', _rateCtrl, 'e.g., 9'),
              const SizedBox(height: 16),
              _buildField('Tenure (months)', _tenureCtrl, 'e.g., 60'),
            ],
          ],
          if (_selectedType == 'sip_increase') ...[
            _buildField('Additional SIP (₹/month)', _sipCtrl, 'e.g., 5000'),
          ],
          if (_selectedType == 'income_change') ...[
            _buildField('Income Change (₹/month)', _incomeCtrl, 'e.g., 25000 (positive or negative)'),
          ],
          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            child: ElevatedButton(
              onPressed: _isLoading ? null : _runSimulation,
              child: _isLoading
                  ? const SizedBox(height: 20, width: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                  : const Text('Run Simulation'),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildField(String label, TextEditingController ctrl, String hint) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
        const SizedBox(height: 8),
        TextField(controller: ctrl, keyboardType: TextInputType.number, decoration: InputDecoration(hintText: hint)),
      ],
    );
  }

  Widget _buildResults() {
    if (_result == null) return const SizedBox();
    final res = _result!;
    final projection = res['results'] is List ? res['results'] : (res['results']?['projection'] ?? []);
    final projData = projection is List ? projection : [];

    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(color: Colors.white, borderRadius: BorderRadius.circular(24), boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 20)]),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              const Icon(Icons.auto_awesome_rounded, color: Color(0xFF6B46C1)),
              const SizedBox(width: 10),
              const Text('Simulation Results', style: TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: Color(0xFF1E293B))),
            ],
          ),
          const SizedBox(height: 20),
          ..._buildResultRows(res),
          if (projData.isNotEmpty) ...[
            const SizedBox(height: 24),
            const Text('Net Worth Projection', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
            const SizedBox(height: 16),
            SizedBox(
              height: 220,
              child: projData.length > 1
                  ? LineChart(
                      LineChartData(
                        gridData: FlGridData(show: false),
                        titlesData: FlTitlesData(
                          leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, reservedSize: 50, getTitlesWidget: (v, _) => Text('₹${(v / 100000).toInt()}L', style: const TextStyle(fontSize: 10, color: Colors.grey)))),
                          bottomTitles: AxisTitles(sideTitles: SideTitles(showTitles: true, interval: 5, getTitlesWidget: (v, _) => Text('${v.toInt()}', style: const TextStyle(fontSize: 10, color: Colors.grey)))),
                          rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                          topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
                        ),
                        borderData: FlBorderData(show: false),
                        lineBarsData: [
                          LineChartBarData(
                            spots: projData.asMap().entries.map((e) => FlSpot(e.key.toDouble(), (e.value['net_worth'] ?? 0).toDouble())).toList(),
                            isCurved: true,
                            gradient: const LinearGradient(colors: [Color(0xFF6B46C1), Color(0xFF8B5CF6)]),
                            barWidth: 3,
                            isStrokeCapRound: true,
                            dotData: FlDotData(show: false),
                            belowBarData: BarAreaData(show: true, gradient: LinearGradient(colors: [Color(0xFF6B46C1).withValues(alpha: 0.15), Color(0xFF8B5CF6).withValues(alpha: 0.05)])),
                          ),
                        ],
                      ),
                    )
                  : const Center(child: Text('Simulate to see projection chart')),
            ),
          ],
        ],
      ),
    );
  }

  List<Widget> _buildResultRows(Map<String, dynamic> res) {
    final results = res['results'] ?? {};
    if (results is! Map) return [];
    final items = <Widget>[];

    for (final entry in (results as Map<String, dynamic>).entries) {
      if (entry.key == 'projection') continue;
      final value = entry.value;
      String display;
      if (value is double) {
        display = '₹${value.toStringAsFixed(0)}';
        if (entry.key.contains('rate') || entry.key.contains('impact') && entry.key != 'fin_score_impact') {
          display = '${value.toStringAsFixed(1)}%';
        }
      } else if (value is int) {
        display = value.toString();
        if (entry.key.contains('score')) {
          display = '$value pts';
        }
      } else {
        display = value.toString();
      }

      items.add(Padding(
        padding: const EdgeInsets.symmetric(vertical: 6),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(_formatLabel(entry.key), style: const TextStyle(color: Colors.grey, fontWeight: FontWeight.w600, fontSize: 13)),
            Text(display, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Color(0xFF1E293B))),
          ],
        ),
      ));
    }
    return items;
  }

  String _formatLabel(String key) {
    return key.replaceAll('_', ' ').split(' ').map((w) => w.isNotEmpty ? '${w[0].toUpperCase()}${w.substring(1)}' : '').join(' ');
  }
}
