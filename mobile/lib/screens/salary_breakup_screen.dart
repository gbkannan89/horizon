import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../models/financial_models.dart';
import '../utils/ui_utils.dart';

class SalaryBreakupScreen extends StatefulWidget {
  final LocalIncome income;
  final LocalSalaryDetail? salaryDetail; // null if creating a new one

  const SalaryBreakupScreen({
    super.key,
    required this.income,
    this.salaryDetail,
  });

  @override
  State<SalaryBreakupScreen> createState() => _SalaryBreakupScreenState();
}

class _SalaryBreakupScreenState extends State<SalaryBreakupScreen> {
  final _formKey = GlobalKey<FormState>();

  late TextEditingController _companyCtrl;
  late TextEditingController _fromYearCtrl;
  late TextEditingController _toYearCtrl;
  late TextEditingController _fixedPayCtrl;
  late TextEditingController _basicPayCtrl;
  late TextEditingController _hraCtrl;
  late TextEditingController _ltaCtrl;
  late TextEditingController _pfEmployeeCtrl;
  late TextEditingController _pfEmployerCtrl;
  late TextEditingController _specialCtrl;
  late TextEditingController _mealCardCtrl;
  late TextEditingController _varPctCtrl;
  late TextEditingController _varAmtCtrl;

  bool _isCurrent = false;
  bool _isLoading = false;

  // Estimated Outputs
  double _grossAnnual = 0.0;
  double _monthlyInHand = 0.0;

  @override
  void initState() {
    super.initState();
    final sd = widget.salaryDetail;

    _companyCtrl = TextEditingController(
      text: sd?.companyName ?? widget.income.companyName,
    );
    _fromYearCtrl = TextEditingController(
      text: sd?.fromYear.toString() ?? DateTime.now().year.toString(),
    );
    _toYearCtrl = TextEditingController(text: sd?.toYear?.toString() ?? '');
    _fixedPayCtrl = TextEditingController(
      text: sd?.fixedPay?.toStringAsFixed(0) ?? '',
    );
    _basicPayCtrl = TextEditingController(
      text: sd?.basicPay?.toStringAsFixed(0) ?? '',
    );
    _hraCtrl = TextEditingController(text: sd?.hra?.toStringAsFixed(0) ?? '');
    _ltaCtrl = TextEditingController(text: sd?.lta?.toStringAsFixed(0) ?? '');
    _pfEmployeeCtrl = TextEditingController(
      text: sd?.pfEmployee?.toStringAsFixed(0) ?? '',
    );
    _pfEmployerCtrl = TextEditingController(
      text: sd?.pfEmployer?.toStringAsFixed(0) ?? '',
    );
    _specialCtrl = TextEditingController(
      text: sd?.specialAllowance?.toStringAsFixed(0) ?? '',
    );
    _mealCardCtrl = TextEditingController(
      text: sd?.mealCard?.toStringAsFixed(0) ?? '',
    );
    _varPctCtrl = TextEditingController(
      text: sd?.variablePayPercentage?.toStringAsFixed(1) ?? '',
    );
    _varAmtCtrl = TextEditingController(
      text: sd?.variablePayAmount?.toStringAsFixed(0) ?? '',
    );

    _isCurrent = sd?.isCurrent ?? (sd == null); // default true for new entries

    _fixedPayCtrl.addListener(_recomputeEstimates);
    _basicPayCtrl.addListener(_recomputeEstimates);
    _pfEmployeeCtrl.addListener(_recomputeEstimates);
    _mealCardCtrl.addListener(_recomputeEstimates);
    _varPctCtrl.addListener(_recomputeEstimates);
    _varAmtCtrl.addListener(_recomputeEstimates);

    _recomputeEstimates();
  }

  @override
  void dispose() {
    _companyCtrl.dispose();
    _fromYearCtrl.dispose();
    _toYearCtrl.dispose();
    _fixedPayCtrl.dispose();
    _basicPayCtrl.dispose();
    _hraCtrl.dispose();
    _ltaCtrl.dispose();
    _pfEmployeeCtrl.dispose();
    _pfEmployerCtrl.dispose();
    _specialCtrl.dispose();
    _mealCardCtrl.dispose();
    _varPctCtrl.dispose();
    _varAmtCtrl.dispose();
    super.dispose();
  }

  void _recomputeEstimates() {
    final fixedPay = double.tryParse(_fixedPayCtrl.text.trim()) ?? 0.0;
    var basicPay = double.tryParse(_basicPayCtrl.text.trim()) ?? 0.0;
    if (fixedPay > 0.0 && basicPay == 0.0) {
      basicPay = fixedPay * 0.5;
    }

    final varPct = double.tryParse(_varPctCtrl.text.trim()) ?? 0.0;
    var varAmt = double.tryParse(_varAmtCtrl.text.trim()) ?? 0.0;
    if (fixedPay > 0.0 && varPct > 0.0 && varAmt == 0.0) {
      varAmt = fixedPay * (varPct / 100.0);
    }

    final pfEmp = double.tryParse(_pfEmployeeCtrl.text.trim()) ?? 0.0;
    final meal = double.tryParse(_mealCardCtrl.text.trim()) ?? 0.0;

    setState(() {
      _grossAnnual = fixedPay + varAmt;
      final monthlyGross = fixedPay / 12.0;
      _monthlyInHand = (monthlyGross - pfEmp - meal).clamp(
        0.0,
        double.infinity,
      );
    });
  }

  Future<void> _saveForm() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    final fromYear =
        int.tryParse(_fromYearCtrl.text.trim()) ?? DateTime.now().year;
    final toYear = int.tryParse(_toYearCtrl.text.trim());
    final fixedPay = double.tryParse(_fixedPayCtrl.text.trim());
    final basicPay = double.tryParse(_basicPayCtrl.text.trim());
    final hra = double.tryParse(_hraCtrl.text.trim());
    final lta = double.tryParse(_ltaCtrl.text.trim());
    final pfEmployee = double.tryParse(_pfEmployeeCtrl.text.trim());
    final pfEmployer = double.tryParse(_pfEmployerCtrl.text.trim());
    final special = double.tryParse(_specialCtrl.text.trim());
    final meal = double.tryParse(_mealCardCtrl.text.trim());
    final varPct = double.tryParse(_varPctCtrl.text.trim());
    final varAmt = double.tryParse(_varAmtCtrl.text.trim());

    final body = {
      'company_name': _companyCtrl.text.trim(),
      'from_year': fromYear,
      'to_year': toYear,
      'is_current': _isCurrent,
      'fixed_pay': ?fixedPay,
      'basic_pay': ?basicPay,
      'hra': ?hra,
      'lta': ?lta,
      'pf_employee': ?pfEmployee,
      'pf_employer': ?pfEmployer,
      'special_allowance': ?special,
      'meal_card': ?meal,
      'variable_pay_percentage': ?varPct,
      'variable_pay_amount': ?varAmt,
    };

    try {
      final provider = Provider.of<FinancialProvider>(context, listen: false);
      if (widget.salaryDetail != null) {
        await provider.updateSalaryDetail(
          widget.income.id,
          widget.salaryDetail!.id,
          body,
        );
      } else {
        await provider.addSalaryDetail(widget.income.id, body);
      }
      if (mounted) {
        UiUtils.showSnack(context, 'Salary breakup saved successfully');
        Navigator.pop(context);
      }
    } catch (e) {
      if (mounted) UiUtils.showSnack(context, 'Failed: $e', isError: true);
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isEdit = widget.salaryDetail != null;

    return Scaffold(
      appBar: AppBar(
        title: Text(isEdit ? 'Edit Salary Detail' : 'Add Salary Detail'),
        actions: [
          if (_isLoading)
            const Center(
              child: Padding(
                padding: EdgeInsets.symmetric(horizontal: 16),
                child: SizedBox(
                  width: 20,
                  height: 20,
                  child: CircularProgressIndicator(
                    color: Colors.white,
                    strokeWidth: 2,
                  ),
                ),
              ),
            )
          else
            IconButton(icon: const Icon(Icons.check), onPressed: _saveForm),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Summary computations card
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(18),
                decoration: BoxDecoration(
                  color: const Color(0xFF6366F1).withValues(alpha: 0.08),
                  borderRadius: BorderRadius.circular(16),
                  border: Border.all(
                    color: const Color(0xFF6366F1).withValues(alpha: 0.15),
                  ),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Computed Salary Summary',
                      style: TextStyle(
                        color: Color(0xFF6366F1),
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 12),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Gross Annual Package:',
                          style: TextStyle(color: Colors.black87),
                        ),
                        Text(
                          '₹${_grossAnnual.toStringAsFixed(0)} /yr',
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 15,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 8),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Est. Monthly In-Hand:',
                          style: TextStyle(color: Colors.black87),
                        ),
                        Text(
                          '₹${_monthlyInHand.toStringAsFixed(0)} /mo',
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            color: Color(0xFF6366F1),
                            fontSize: 16,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 24),

              const Text(
                'Company & Duration',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: Color(0xFF1E293B),
                ),
              ),
              const Divider(),
              const SizedBox(height: 8),
              TextFormField(
                controller: _companyCtrl,
                decoration: const InputDecoration(
                  labelText: 'Company Name',
                  border: OutlineInputBorder(),
                ),
                validator: (val) =>
                    val == null || val.isEmpty ? 'Company name required' : null,
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: TextFormField(
                      controller: _fromYearCtrl,
                      keyboardType: TextInputType.number,
                      decoration: const InputDecoration(
                        labelText: 'From Year (e.g. 2025)',
                        border: OutlineInputBorder(),
                      ),
                      validator: (val) => val == null || val.isEmpty
                          ? 'Start year required'
                          : null,
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: TextFormField(
                      controller: _toYearCtrl,
                      keyboardType: TextInputType.number,
                      decoration: const InputDecoration(
                        labelText: 'To Year (Optional)',
                        border: OutlineInputBorder(),
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              SwitchListTile(
                value: _isCurrent,
                title: const Text('Is Current Salary'),
                subtitle: const Text('Syncs with primary income list'),
                activeThumbColor: const Color(0xFF6366F1),
                onChanged: (val) {
                  setState(() => _isCurrent = val);
                },
              ),
              const SizedBox(height: 24),

              const Text(
                'Fixed Pay Components',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: Color(0xFF1E293B),
                ),
              ),
              const Divider(),
              const SizedBox(height: 8),
              TextFormField(
                controller: _fixedPayCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Fixed Pay (Annual, ₹)',
                  border: OutlineInputBorder(),
                ),
                validator: (val) =>
                    val == null || val.isEmpty ? 'Fixed pay required' : null,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _basicPayCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Basic Pay (Annual, ₹ - Optional)',
                  hintText: 'Defaults to 50% of Fixed Pay',
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _hraCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: 'HRA (Annual, ₹)'),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _ltaCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: 'LTA (Annual, ₹)'),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _specialCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Special Allowance (Annual, ₹)',
                ),
              ),
              const SizedBox(height: 24),

              const Text(
                'Deductions & Benefits (Monthly)',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: Color(0xFF1E293B),
                ),
              ),
              const Divider(),
              const SizedBox(height: 8),
              TextFormField(
                controller: _pfEmployeeCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'PF (Employee Contribution, Monthly ₹)',
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _pfEmployerCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'PF (Employer Contribution, Monthly ₹)',
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _mealCardCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Meal Card Deductions (Monthly ₹)',
                ),
              ),
              const SizedBox(height: 24),

              const Text(
                'Variable Pay',
                style: TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.bold,
                  color: Color(0xFF1E293B),
                ),
              ),
              const Divider(),
              const SizedBox(height: 8),
              TextFormField(
                controller: _varPctCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Variable Pay Percentage (%)',
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _varAmtCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Variable Pay Amount (Annual, ₹ - Optional)',
                  hintText: 'Auto-calculates if target % is set',
                ),
              ),
              const SizedBox(height: 40),
            ],
          ),
        ),
      ),
    );
  }
}
