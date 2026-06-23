import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../providers/financial_provider.dart';
import 'main_shell.dart';
import 'dart:ui';

class OnboardingScreen extends StatefulWidget {
  const OnboardingScreen({super.key});

  @override
  State<OnboardingScreen> createState() => _OnboardingScreenState();
}

class _OnboardingScreenState extends State<OnboardingScreen> {
  final PageController _pageController = PageController();
  int _currentIndex = 0;

  final _incomeCtrl = TextEditingController();
  final _rentCtrl = TextEditingController();
  final _cashCtrl = TextEditingController();
  
  double _runwayMonths = 6.0;
  bool _isLoading = false;

  Future<void> _completeOnboarding() async {
    setState(() => _isLoading = true);

    final provider = Provider.of<FinancialProvider>(context, listen: false);
    
    try {
      // Step 1: Add Income
      final income = double.tryParse(_incomeCtrl.text.replaceAll(',', ''));
      if (income != null && income > 0) {
        await provider.addIncome('Primary Salary', 'salary', income, 'monthly');
      }

      // Step 2: Add Rent/Housing Bill
      final rent = double.tryParse(_rentCtrl.text.replaceAll(',', ''));
      if (rent != null && rent > 0) {
        await provider.addExpense('Rent/Mortgage', rent, 'Housing', 'Needs');
      }

      // Step 3: Add Bank Asset (Liquid Cash)
      final cash = double.tryParse(_cashCtrl.text.replaceAll(',', ''));
      if (cash != null && cash > 0) {
        await provider.addAsset('Primary Bank Account', 'bank', cash);
      }

      // Save custom runway locally
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
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  void _nextPage() {
    if (_currentIndex < 3) {
      _pageController.nextPage(duration: const Duration(milliseconds: 500), curve: Curves.easeOutCubic);
    } else {
      _completeOnboarding();
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
              // Progress Indicator
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 32, vertical: 24),
                child: Row(
                  children: List.generate(4, (index) {
                    return Expanded(
                      child: AnimatedContainer(
                        duration: const Duration(milliseconds: 400),
                        height: 4,
                        margin: EdgeInsets.only(right: index < 3 ? 8 : 0),
                        decoration: BoxDecoration(
                          color: _currentIndex >= index ? const Color(0xFF60A5FA) : Colors.white.withOpacity(0.2),
                          borderRadius: BorderRadius.circular(2),
                        ),
                      ),
                    );
                  }),
                ),
              ),
              
              Expanded(
                child: PageView(
                  controller: _pageController,
                  physics: const NeverScrollableScrollPhysics(),
                  onPageChanged: (index) => setState(() => _currentIndex = index),
                  children: [
                    _buildStep(
                      title: "What's your monthly salary?",
                      subtitle: "We use this to establish your 50/30/20 baseline.",
                      icon: Icons.work_outline_rounded,
                      controller: _incomeCtrl,
                      hint: "50,000",
                    ),
                    _buildStep(
                      title: "What's your monthly rent?",
                      subtitle: "Your largest fixed cost helps us determine your required safety net.",
                      icon: Icons.home_work_outlined,
                      controller: _rentCtrl,
                      hint: "15,000",
                    ),
                    _buildStep(
                      title: "Current liquid cash?",
                      subtitle: "Include savings and checking accounts that are immediately accessible.",
                      icon: Icons.account_balance_wallet_outlined,
                      controller: _cashCtrl,
                      hint: "1,00,000",
                    ),
                    _buildRunwayStep(),
                  ],
                ),
              ),

              Padding(
                padding: const EdgeInsets.fromLTRB(32, 16, 32, 32),
                child: Column(
                  children: [
                    SizedBox(
                      width: double.infinity,
                      height: 56,
                      child: ElevatedButton(
                        onPressed: _isLoading ? null : _nextPage,
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF3B82F6),
                          foregroundColor: Colors.white,
                          elevation: 8,
                          shadowColor: const Color(0xFF3B82F6).withOpacity(0.5),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                        ),
                        child: _isLoading 
                          ? const SizedBox(height: 24, width: 24, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                          : Text(
                              _currentIndex == 3 ? 'Complete Setup' : 'Continue',
                              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                            ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    if (_currentIndex < 3)
                      TextButton(
                        onPressed: () {
                          _pageController.animateToPage(3, duration: const Duration(milliseconds: 500), curve: Curves.easeOutCubic);
                        },
                        child: const Text('Skip for now', style: TextStyle(color: Colors.white70, fontWeight: FontWeight.w600)),
                      ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildStep({
    required String title,
    required String subtitle,
    required IconData icon,
    required TextEditingController controller,
    required String hint,
  }) {
    return SingleChildScrollView(
      physics: const BouncingScrollPhysics(),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 32.0, vertical: 24.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const SizedBox(height: 20),
            TweenAnimationBuilder<double>(
            tween: Tween(begin: 0.0, end: 1.0),
            duration: const Duration(milliseconds: 800),
            curve: Curves.easeOutBack,
            builder: (context, value, child) {
              return Transform.scale(
                scale: value,
                child: child,
              );
            },
            child: Container(
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Colors.white.withOpacity(0.1),
                shape: BoxShape.circle,
                border: Border.all(color: Colors.white.withOpacity(0.2), width: 1),
              ),
              child: Icon(icon, size: 64, color: const Color(0xFF60A5FA)),
            ),
          ),
          const SizedBox(height: 40),
          Text(
            title,
            textAlign: TextAlign.center,
            style: const TextStyle(fontSize: 32, fontWeight: FontWeight.w900, color: Colors.white, height: 1.2),
          ),
          const SizedBox(height: 16),
          Text(
            subtitle,
            textAlign: TextAlign.center,
            style: TextStyle(fontSize: 15, color: Colors.white.withOpacity(0.7), height: 1.5),
          ),
          const SizedBox(height: 48),
          
          // Glassmorphism Input Field
          ClipRRect(
            borderRadius: BorderRadius.circular(20),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
              child: Container(
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.05),
                  border: Border.all(color: Colors.white.withOpacity(0.2)),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: TextField(
                  controller: controller,
                  keyboardType: TextInputType.number,
                  style: const TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Colors.white),
                  textAlign: TextAlign.center,
                  decoration: InputDecoration(
                    prefixText: '₹ ',
                    prefixStyle: const TextStyle(fontSize: 28, fontWeight: FontWeight.bold, color: Color(0xFF60A5FA)),
                    hintText: hint,
                    hintStyle: TextStyle(fontSize: 28, color: Colors.white.withOpacity(0.2), fontWeight: FontWeight.normal),
                    border: InputBorder.none,
                    contentPadding: const EdgeInsets.symmetric(vertical: 24),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
      ),
    );
  }

  Widget _buildRunwayStep() {
    final rentAmt = double.tryParse(_rentCtrl.text.replaceAll(',', '')) ?? 15000;
    // Estimate basic needs as double the rent for a quick onboarding heuristic
    final estimatedNeeds = rentAmt * 2; 
    final target = estimatedNeeds * _runwayMonths;

    return SingleChildScrollView(
      physics: const BouncingScrollPhysics(),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 32.0, vertical: 24.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const SizedBox(height: 20),
            Container(
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              color: Colors.white.withOpacity(0.1),
              shape: BoxShape.circle,
              border: Border.all(color: Colors.white.withOpacity(0.2), width: 1),
            ),
            child: const Icon(Icons.shield_outlined, size: 64, color: Color(0xFF10B981)),
          ),
          const SizedBox(height: 40),
          const Text(
            "Your Safety Net",
            textAlign: TextAlign.center,
            style: TextStyle(fontSize: 32, fontWeight: FontWeight.w900, color: Colors.white, height: 1.2),
          ),
          const SizedBox(height: 16),
          Text(
            "How many months of runway do you want in your Emergency Fund?",
            textAlign: TextAlign.center,
            style: TextStyle(fontSize: 15, color: Colors.white.withOpacity(0.7), height: 1.5),
          ),
          const SizedBox(height: 40),
          
          Text(
            "${_runwayMonths.toInt()} Months",
            style: const TextStyle(fontSize: 36, fontWeight: FontWeight.bold, color: Color(0xFF10B981)),
          ),
          Slider(
            value: _runwayMonths,
            min: 1,
            max: 12,
            divisions: 11,
            activeColor: const Color(0xFF10B981),
            inactiveColor: Colors.white.withOpacity(0.2),
            onChanged: (val) {
              setState(() {
                _runwayMonths = val;
              });
            },
          ),
          const SizedBox(height: 24),
          
          ClipRRect(
            borderRadius: BorderRadius.circular(20),
            child: BackdropFilter(
              filter: ImageFilter.blur(sigmaX: 10, sigmaY: 10),
              child: Container(
                padding: const EdgeInsets.all(24),
                width: double.infinity,
                decoration: BoxDecoration(
                  color: Colors.white.withOpacity(0.05),
                  border: Border.all(color: Colors.white.withOpacity(0.2)),
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Column(
                  children: [
                    Text("Target Goal", style: TextStyle(color: Colors.white.withOpacity(0.7), fontSize: 14)),
                    const SizedBox(height: 8),
                    Text("₹${target.toStringAsFixed(0)}", style: const TextStyle(color: Colors.white, fontSize: 28, fontWeight: FontWeight.bold)),
                  ],
                ),
              ),
            ),
          )
        ],
      ),
      ),
    );
  }
}
