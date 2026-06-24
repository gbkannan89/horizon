import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../providers/financial_provider.dart';
import 'main_shell.dart';

class OnboardingScreen extends StatefulWidget {
  const OnboardingScreen({super.key});

  @override
  State<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends State<OnboardingScreen>
    with SingleTickerProviderStateMixin {
  final PageController _pageController = PageController();
  int _currentIndex = 0;

  final _incomeCtrl = TextEditingController();
  final _rentCtrl = TextEditingController();
  final _cashCtrl = TextEditingController();

  double _runwayMonths = 6.0;
  bool _isSubmitting = false;

  late AnimationController _animController;
  late Animation<double> _fadeIn;

  static const _steps = ['Income', 'Housing', 'Savings', 'Review'];

  @override
  void initState() {
    super.initState();
    _animController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 600),
    );
    _fadeIn = CurvedAnimation(parent: _animController, curve: Curves.easeOut);
    _animController.forward();
  }

  @override
  void dispose() {
    _pageController.dispose();
    _incomeCtrl.dispose();
    _rentCtrl.dispose();
    _cashCtrl.dispose();
    _animController.dispose();
    super.dispose();
  }

  String _formatCurrency(String value) {
    final clean = value.replaceAll(',', '');
    if (clean.isEmpty) return '';
    final num = int.tryParse(clean);
    if (num == null) return value;
    final str = num.toString();
    if (str.length <= 3) return str;
    final last3 = str.substring(str.length - 3);
    var rest = str.substring(0, str.length - 3);
    final groups = <String>[last3];
    while (rest.isNotEmpty) {
      final take = rest.length > 2 ? 2 : rest.length;
      groups.insert(0, rest.substring(rest.length - take));
      rest = rest.substring(0, rest.length - take);
    }
    return groups.join(',');
  }

  void _onInputChanged(TextEditingController ctrl, String text) {
    final clean = text.replaceAll(',', '');
    if (clean.isEmpty || !RegExp(r'^\d+$').hasMatch(clean)) {
      ctrl.text = '';
      ctrl.selection = TextSelection.collapsed(offset: 0);
      return;
    }
    final formatted = _formatCurrency(clean);
    if (formatted != ctrl.text) {
      ctrl.text = formatted;
      ctrl.selection = TextSelection.collapsed(offset: formatted.length);
    }
  }

  Future<void> _completeOnboarding() async {
    setState(() => _isSubmitting = true);

    final provider = Provider.of<FinancialProvider>(context, listen: false);

    try {
      final income = double.tryParse(_incomeCtrl.text.replaceAll(',', ''));
      if (income != null && income > 0) {
        await provider.addIncome('Primary Salary', 'salary', income, 'monthly');
      }

      final rent = double.tryParse(_rentCtrl.text.replaceAll(',', ''));
      if (rent != null && rent > 0) {
        await provider.addExpense('Rent/Mortgage', rent, 'Housing', 'Needs');
      }

      final cash = double.tryParse(_cashCtrl.text.replaceAll(',', ''));
      if (cash != null && cash > 0) {
        await provider.addAsset('Primary Bank Account', 'bank', cash);
      }

      final prefs = await SharedPreferences.getInstance();
      await prefs.setDouble('emergency_runway_months', _runwayMonths);
      await prefs.setBool('has_onboarded', true);

      if (mounted) {
        Navigator.of(context).pushReplacement(
          MaterialPageRoute(builder: (_) => const MainShell()),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Error: $e'), backgroundColor: Colors.red),
        );
      }
    } finally {
      if (mounted) setState(() => _isSubmitting = false);
    }
  }

  void _goNext() {
    if (_currentIndex < 3) {
      _animController.reset();
      _pageController.nextPage(
        duration: const Duration(milliseconds: 400),
        curve: Curves.easeInOut,
      );
      _animController.forward();
    } else {
      _completeOnboarding();
    }
  }

  void _goBack() {
    if (_currentIndex > 0) {
      _animController.reset();
      _pageController.previousPage(
        duration: const Duration(milliseconds: 400),
        curve: Curves.easeInOut,
      );
      _animController.forward();
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            colors: [Color(0xFF0F172A), Color(0xFF1E1B4B), Color(0xFF312E81)],
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
          ),
        ),
        child: SafeArea(
          child: Column(
            children: [
              _buildTopBar(),
              _buildStepIndicator(),
              Expanded(
                child: PageView(
                  controller: _pageController,
                  physics: const NeverScrollableScrollPhysics(),
                  onPageChanged: (index) {
                    setState(() => _currentIndex = index);
                    _animController.reset();
                    _animController.forward();
                  },
                  children: [
                    _buildIncomeStep(),
                    _buildRentStep(),
                    _buildSavingsStep(),
                    _buildReviewStep(),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTopBar() {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 4),
      child: Row(
        children: [
          if (_currentIndex > 0)
            GestureDetector(
              onTap: _goBack,
              child: Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: const Icon(Icons.arrow_back_rounded, color: Colors.white, size: 22),
              ),
            )
          else
            const SizedBox(width: 38),
          const Spacer(),
          Text(
            'Step ${_currentIndex + 1} of 4',
            style: TextStyle(color: Colors.white.withValues(alpha: 0.5), fontSize: 13, fontWeight: FontWeight.w500),
          ),
          const Spacer(),
          const SizedBox(width: 38),
        ],
      ),
    );
  }

  Widget _buildStepIndicator() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 8),
      child: Row(
        children: List.generate(4, (index) {
          final isActive = _currentIndex >= index;
          final isCurrent = _currentIndex == index;
          return Expanded(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                Row(
                  children: [
                    if (index > 0)
                      Expanded(
                        child: Container(
                          height: 2,
                          decoration: BoxDecoration(
                            color: isActive
                                ? const Color(0xFF60A5FA)
                                : Colors.white.withValues(alpha: 0.15),
                          ),
                        ),
                      ),
                    Container(
                      width: 32,
                      height: 32,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        color: isActive
                            ? const Color(0xFF60A5FA)
                            : Colors.white.withValues(alpha: 0.08),
                        border: isCurrent && !isActive
                            ? Border.all(color: const Color(0xFF60A5FA), width: 2)
                            : null,
                      ),
                      child: Center(
                        child: isActive
                            ? const Icon(Icons.check, color: Colors.white, size: 16)
                            : Text(
                                '${index + 1}',
                                style: TextStyle(
                                  color: isCurrent
                                      ? const Color(0xFF60A5FA)
                                      : Colors.white.withValues(alpha: 0.4),
                                  fontSize: 13,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                      ),
                    ),
                    if (index < 3)
                      Expanded(
                        child: Container(
                          height: 2,
                          decoration: BoxDecoration(
                            color: _currentIndex > index
                                ? const Color(0xFF60A5FA)
                                : Colors.white.withValues(alpha: 0.15),
                          ),
                        ),
                      ),
                  ],
                ),
                const SizedBox(height: 6),
                Text(
                  _steps[index],
                  style: TextStyle(
                    color: isActive ? Colors.white : Colors.white.withValues(alpha: 0.3),
                    fontSize: 11,
                    fontWeight: isActive ? FontWeight.w600 : FontWeight.normal,
                  ),
                ),
              ],
            ),
          );
        }),
      ),
    );
  }

  Widget _buildInputField({
    required TextEditingController controller,
    required String hint,
  }) {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.08),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: Colors.white.withValues(alpha: 0.15)),
      ),
      child: TextField(
        controller: controller,
        keyboardType: TextInputType.number,
        style: const TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Colors.white),
        textAlign: TextAlign.center,
        onChanged: (v) => _onInputChanged(controller, v),
        decoration: InputDecoration(
          prefixText: '₹ ',
          prefixStyle: const TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Color(0xFF60A5FA)),
          hintText: hint,
          hintStyle: TextStyle(fontSize: 28, color: Colors.white.withValues(alpha: 0.15), fontWeight: FontWeight.normal),
          border: InputBorder.none,
          contentPadding: const EdgeInsets.symmetric(vertical: 24),
        ),
      ),
    );
  }

  Widget _buildStepCard({
    required String emoji,
    required String title,
    required String subtitle,
    required Widget input,
    required bool isLast,
  }) {
    return FadeTransition(
      opacity: _fadeIn,
      child: SingleChildScrollView(
        physics: const BouncingScrollPhysics(),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 8),
          child: Column(
            children: [
              const SizedBox(height: 24),
              Container(
                width: 88,
                height: 88,
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.08),
                  shape: BoxShape.circle,
                  border: Border.all(color: Colors.white.withValues(alpha: 0.15), width: 1),
                ),
                child: Center(child: Text(emoji, style: const TextStyle(fontSize: 44))),
              ),
              const SizedBox(height: 28),
              Text(
                title,
                textAlign: TextAlign.center,
                style: const TextStyle(fontSize: 26, fontWeight: FontWeight.w900, color: Colors.white, height: 1.2),
              ),
              const SizedBox(height: 12),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 8),
                child: Text(
                  subtitle,
                  textAlign: TextAlign.center,
                  style: TextStyle(fontSize: 15, color: Colors.white.withValues(alpha: 0.55), height: 1.5),
                ),
              ),
              const SizedBox(height: 36),
              input,
              const SizedBox(height: 32),
              SizedBox(
                width: double.infinity,
                height: 56,
                child: ElevatedButton(
                  onPressed: _goNext,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF3B82F6),
                    foregroundColor: Colors.white,
                    elevation: 4,
                    shadowColor: const Color(0xFF3B82F6).withValues(alpha: 0.4),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                  ),
                  child: const Text('Continue', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                ),
              ),
              if (!isLast) ...[
                const SizedBox(height: 12),
                TextButton(
                  onPressed: _goNext,
                  child: Text('Skip this step',
                      style: TextStyle(color: Colors.white.withValues(alpha: 0.45), fontWeight: FontWeight.w600, fontSize: 14)),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildIncomeStep() {
    return _buildStepCard(
      emoji: '💰',
      title: "What's your monthly salary?",
      subtitle: "We use this to set up your 50/30/20 budget baseline.",
      isLast: false,
      input: _buildInputField(controller: _incomeCtrl, hint: '50,000'),
    );
  }

  Widget _buildRentStep() {
    return _buildStepCard(
      emoji: '🏠',
      title: "Monthly rent or mortgage?",
      subtitle: "Your largest fixed cost helps us calculate your safety net target.",
      isLast: false,
      input: _buildInputField(controller: _rentCtrl, hint: '15,000'),
    );
  }

  Widget _buildSavingsStep() {
    return _buildStepCard(
      emoji: '🏦',
      title: "Current liquid savings?",
      subtitle: "Include checking, savings, and any instantly accessible cash.",
      isLast: false,
      input: _buildInputField(controller: _cashCtrl, hint: '1,00,000'),
    );
  }

  Widget _buildReviewStep() {
    final income = double.tryParse(_incomeCtrl.text.replaceAll(',', '')) ?? 0;
    final rent = double.tryParse(_rentCtrl.text.replaceAll(',', '')) ?? 0;
    final cash = double.tryParse(_cashCtrl.text.replaceAll(',', '')) ?? 0;
    final target = (rent * 2) * _runwayMonths;

    return FadeTransition(
      opacity: _fadeIn,
      child: SingleChildScrollView(
        physics: const BouncingScrollPhysics(),
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 8),
          child: Column(
            children: [
              const SizedBox(height: 16),
              Container(
                width: 88,
                height: 88,
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.08),
                  shape: BoxShape.circle,
                  border: Border.all(color: Colors.white.withValues(alpha: 0.15), width: 1),
                ),
                child: const Center(child: Text('✅', style: TextStyle(fontSize: 44))),
              ),
              const SizedBox(height: 24),
              const Text(
                "Ready to go!",
                textAlign: TextAlign.center,
                style: TextStyle(fontSize: 26, fontWeight: FontWeight.w900, color: Colors.white, height: 1.2),
              ),
              const SizedBox(height: 8),
              Text(
                'Review your details before we start.',
                textAlign: TextAlign.center,
                style: TextStyle(fontSize: 15, color: Colors.white.withValues(alpha: 0.55), height: 1.5),
              ),
              const SizedBox(height: 28),

              // Summary card
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.06),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(color: Colors.white.withValues(alpha: 0.1)),
                ),
                child: Column(
                  children: [
                    _summaryRow('💰 Salary',
                        income > 0 ? '₹${_formatCurrency(income.toStringAsFixed(0))}/mo' : 'Skipped', income > 0),
                    const Divider(color: Colors.white12, height: 20),
                    _summaryRow('🏠 Rent',
                        rent > 0 ? '₹${_formatCurrency(rent.toStringAsFixed(0))}/mo' : 'Skipped', rent > 0),
                    const Divider(color: Colors.white12, height: 20),
                    _summaryRow('🏦 Savings',
                        cash > 0 ? '₹${_formatCurrency(cash.toStringAsFixed(0))}' : 'Skipped', cash > 0),
                    const Divider(color: Colors.white12, height: 20),
                    _summaryRow('🛡️ Runway', '${_runwayMonths.toInt()} months', true),
                  ],
                ),
              ),

              const SizedBox(height: 20),

              // Runway slider card
              Container(
                width: double.infinity,
                padding: const EdgeInsets.all(20),
                decoration: BoxDecoration(
                  color: Colors.white.withValues(alpha: 0.06),
                  borderRadius: BorderRadius.circular(20),
                  border: Border.all(color: Colors.white.withValues(alpha: 0.1)),
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text('Emergency Runway',
                            style: TextStyle(color: Colors.white.withValues(alpha: 0.7), fontSize: 14)),
                        Text('${_runwayMonths.toInt()} months',
                            style: const TextStyle(
                                color: Color(0xFF10B981), fontSize: 18, fontWeight: FontWeight.bold)),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Slider(
                      value: _runwayMonths,
                      min: 1,
                      max: 12,
                      divisions: 11,
                      activeColor: const Color(0xFF10B981),
                      inactiveColor: Colors.white.withValues(alpha: 0.15),
                      onChanged: (val) => setState(() => _runwayMonths = val),
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text('1 mo', style: TextStyle(color: Colors.white.withValues(alpha: 0.3), fontSize: 11)),
                        if (rent > 0)
                          Text('Target: ₹${_formatCurrency(target.toStringAsFixed(0))}',
                              style: const TextStyle(
                                  color: Color(0xFF10B981), fontSize: 13, fontWeight: FontWeight.w600)),
                        Text('12 mo', style: TextStyle(color: Colors.white.withValues(alpha: 0.3), fontSize: 11)),
                      ],
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 28),

              SizedBox(
                width: double.infinity,
                height: 56,
                child: ElevatedButton(
                  onPressed: _isSubmitting ? null : _completeOnboarding,
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF10B981),
                    foregroundColor: Colors.white,
                    elevation: 4,
                    shadowColor: const Color(0xFF10B981).withValues(alpha: 0.4),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                  ),
                  child: _isSubmitting
                      ? const SizedBox(
                          height: 24,
                          width: 24,
                          child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2.5))
                      : const Text('Complete Setup',
                          style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                ),
              ),
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }

  Widget _summaryRow(String label, String value, bool filled) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label, style: TextStyle(color: Colors.white.withValues(alpha: 0.7), fontSize: 14)),
        Text(
          value,
          style: TextStyle(
            color: filled ? Colors.white : Colors.white.withValues(alpha: 0.35),
            fontSize: 14,
            fontWeight: FontWeight.bold,
          ),
        ),
      ],
    );
  }
}
