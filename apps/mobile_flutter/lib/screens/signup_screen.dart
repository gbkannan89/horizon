import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'dart:ui';
import '../providers/auth_provider.dart';
import 'onboarding_screen.dart';
import '../utils/ui_utils.dart';

class SignupScreen extends StatefulWidget {
  const SignupScreen({super.key});

  @override
  State<SignupScreen> createState() => _SignupScreenState();
}

class _SignupScreenState extends State<SignupScreen> with SingleTickerProviderStateMixin {
  final _nameCtrl = TextEditingController();
  final _emailCtrl = TextEditingController();
  final _phoneCtrl = TextEditingController();
  final _passwordCtrl = TextEditingController();
  bool _obscurePassword = true;
  String _userType = 'salaried';
  String _riskProfile = 'moderate';
  double _passwordStrength = 0;

  late final AnimationController _animCtrl;
  late final Animation<double> _fadeAnim;
  late final Animation<Offset> _slideAnim;

  final _focusName = FocusNode();
  final _focusEmail = FocusNode();
  final _focusPhone = FocusNode();
  final _focusPassword = FocusNode();

  @override
  void initState() {
    super.initState();
    _animCtrl = AnimationController(vsync: this, duration: const Duration(milliseconds: 1000));
    _fadeAnim = CurvedAnimation(parent: _animCtrl, curve: Curves.easeOut);
    _slideAnim = Tween<Offset>(begin: const Offset(0, 0.12), end: Offset.zero)
        .animate(CurvedAnimation(parent: _animCtrl, curve: Curves.easeOutCubic));
    _animCtrl.forward();

    _passwordCtrl.addListener(_updateStrength);
  }

  void _updateStrength() {
    final p = _passwordCtrl.text;
    double s = 0;
    if (p.length >= 6) s += 0.25;
    if (p.length >= 10) s += 0.15;
    if (RegExp(r'[A-Z]').hasMatch(p)) s += 0.2;
    if (RegExp(r'[a-z]').hasMatch(p)) s += 0.1;
    if (RegExp(r'[0-9]').hasMatch(p)) s += 0.15;
    if (RegExp(r'[!@#$%^&*(),.?":{}|<>]').hasMatch(p)) s += 0.15;
    setState(() => _passwordStrength = s.clamp(0.0, 1.0));
  }

  @override
  void dispose() {
    _animCtrl.dispose();
    _nameCtrl.dispose();
    _emailCtrl.dispose();
    _phoneCtrl.dispose();
    _passwordCtrl.dispose();
    _focusName.dispose();
    _focusEmail.dispose();
    _focusPhone.dispose();
    _focusPassword.dispose();
    super.dispose();
  }

  Future<void> _handleSignup() async {
    if (_nameCtrl.text.trim().isEmpty || _emailCtrl.text.trim().isEmpty || _passwordCtrl.text.isEmpty) {
      UiUtils.showSnack(context, 'Please fill in name, email and password', isError: true);
      return;
    }
    try {
      final auth = Provider.of<AuthProvider>(context, listen: false);
      final ok = await auth.register(
        _nameCtrl.text.trim(), _emailCtrl.text.trim(), _passwordCtrl.text,
        userType: _userType, riskProfile: _riskProfile, phone: _phoneCtrl.text.trim(),
      );
      if (!mounted) return;
      if (ok) {
        Navigator.of(context).pushReplacement(
          MaterialPageRoute(builder: (_) => const OnboardingScreen()),
        );
      } else {
        UiUtils.showSnack(context, 'Registration failed or email already exists', isError: true);
      }
    } catch (e) {
      if (mounted) UiUtils.showSnack(context, 'Registration failed: ${e.toString()}', isError: true);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isLoading = context.watch<AuthProvider>().isLoading;
    final size = MediaQuery.of(context).size;

    return Scaffold(
      body: Container(
        width: double.infinity,
        height: double.infinity,
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFF0F2057), Color(0xFF1E3A8A), Color(0xFF2563EB)],
          ),
        ),
        child: SafeArea(
          child: SingleChildScrollView(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: FadeTransition(
              opacity: _fadeAnim,
              child: SlideTransition(
                position: _slideAnim,
                child: Column(
                  children: [
                    SizedBox(height: size.height * 0.05),

                    // ── Logo ─────────────────────────────────────────────────
                    Container(
                      width: 90,
                      height: 90,
                      decoration: BoxDecoration(
                        color: Colors.white,
                        shape: BoxShape.circle,
                        boxShadow: [
                          BoxShadow(
                            color: Colors.black.withOpacity(0.15),
                            blurRadius: 20,
                            offset: const Offset(0, 8),
                          ),
                        ],
                      ),
                      child: ClipOval(
                        child: Padding(
                          padding: const EdgeInsets.all(4.0),
                          child: Image.asset(
                            'assets/logo.png',
                            fit: BoxFit.cover,
                          ),
                        ),
                      ),
                    ),
                    const SizedBox(height: 16),
                    const Text('Horizon', style: TextStyle(
                      color: Colors.white, fontSize: 28,
                      fontWeight: FontWeight.w900, letterSpacing: -1,
                    )),
                    const SizedBox(height: 4),
                    Text('Start your wealth journey', style: TextStyle(
                      color: Colors.white.withOpacity(0.6), fontSize: 14,
                      fontWeight: FontWeight.w500,
                    )),
                    SizedBox(height: size.height * 0.04),

                    // ── Glass form card ──────────────────────────────────────
                    ClipRRect(
                      borderRadius: BorderRadius.circular(28),
                      child: BackdropFilter(
                        filter: ImageFilter.blur(sigmaX: 12, sigmaY: 12),
                        child: Container(
                          padding: const EdgeInsets.all(24),
                          decoration: BoxDecoration(
                            color: Colors.white.withOpacity(0.08),
                            border: Border.all(color: Colors.white.withOpacity(0.12)),
                            borderRadius: BorderRadius.circular(28),
                          ),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              // ── User Type Selector ──────────────────────────
                              Text("I am a...", style: TextStyle(
                                color: Colors.white.withOpacity(0.7),
                                fontSize: 12, fontWeight: FontWeight.w600,
                              )),
                              const SizedBox(height: 10),
                              Row(children: [
                                _typeChip('salaried', '💼', 'Salaried'),
                                const SizedBox(width: 8),
                                _typeChip('self-employed', '🦄', 'Freelancer'),
                                const SizedBox(width: 8),
                                _typeChip('business', '🏢', 'Business'),
                                const SizedBox(width: 8),
                                _typeChip('student', '🎓', 'Student'),
                              ]),
                              const SizedBox(height: 24),

                              // ── Name ────────────────────────────────────────
                              _glassField(
                                controller: _nameCtrl,
                                focusNode: _focusName,
                                icon: Icons.person_outline_rounded,
                                hint: 'Full name',
                                nextFocus: _focusEmail,
                              ),
                              const SizedBox(height: 14),

                              // ── Email ───────────────────────────────────────
                              _glassField(
                                controller: _emailCtrl,
                                focusNode: _focusEmail,
                                icon: Icons.email_outlined,
                                hint: 'Email address',
                                keyboardType: TextInputType.emailAddress,
                                nextFocus: _focusPhone,
                              ),
                              const SizedBox(height: 14),

                              // ── Phone (optional) ────────────────────────────
                              _glassField(
                                controller: _phoneCtrl,
                                focusNode: _focusPhone,
                                icon: Icons.phone_outlined,
                                hint: 'Phone (optional)',
                                keyboardType: TextInputType.phone,
                                nextFocus: _focusPassword,
                              ),
                              const SizedBox(height: 14),

                              // ── Password ────────────────────────────────────
                              _glassField(
                                controller: _passwordCtrl,
                                focusNode: _focusPassword,
                                icon: Icons.lock_outline_rounded,
                                hint: 'Create a password',
                                obscure: _obscurePassword,
                                onToggleObscure: () => setState(() => _obscurePassword = !_obscurePassword),
                                onSubmit: isLoading ? null : _handleSignup,
                              ),
                              if (_passwordCtrl.text.isNotEmpty) ...[
                                const SizedBox(height: 8),
                                _passwordStrengthBar(),
                              ],

                              const SizedBox(height: 20),

                              // ── Risk Profile ────────────────────────────────
                              Text("Investment style", style: TextStyle(
                                color: Colors.white.withOpacity(0.7),
                                fontSize: 12, fontWeight: FontWeight.w600,
                              )),
                              const SizedBox(height: 10),
                              Row(children: [
                                _riskChip('conservative', '🛡️', 'Safe'),
                                const SizedBox(width: 8),
                                _riskChip('moderate', '⚖️', 'Balanced'),
                                const SizedBox(width: 8),
                                _riskChip('aggressive', '🚀', 'Growth'),
                              ]),
                              const SizedBox(height: 24),

                              // ── Terms hint ──────────────────────────────────
                              Text(
                                'By creating an account, you agree to our Terms and Privacy Policy.',
                                style: TextStyle(
                                  color: Colors.white.withOpacity(0.4),
                                  fontSize: 11, fontWeight: FontWeight.w400,
                                ),
                              ),
                              const SizedBox(height: 16),

                              // ── Submit ──────────────────────────────────────
                              SizedBox(
                                width: double.infinity,
                                height: 54,
                                child: ElevatedButton(
                                  onPressed: isLoading ? null : _handleSignup,
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: Colors.white,
                                    foregroundColor: const Color(0xFF1E3A8A),
                                    elevation: 0,
                                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                                  ),
                                  child: isLoading
                                      ? const SizedBox(width: 22, height: 22,
                                          child: CircularProgressIndicator(strokeWidth: 2.5, color: Color(0xFF1E3A8A)))
                                      : const Text('Create Account', style: TextStyle(
                                          fontSize: 16, fontWeight: FontWeight.w800, letterSpacing: 0.3)),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),

                    const SizedBox(height: 24),

                    // ── Login link ────────────────────────────────────────────
                    Row(mainAxisAlignment: MainAxisAlignment.center, children: [
                      Text("Already have an account? ", style: TextStyle(
                        color: Colors.white.withOpacity(0.6), fontSize: 14)),
                      GestureDetector(
                        onTap: () => Navigator.of(context).pop(),
                        child: const Text('Sign In', style: TextStyle(
                          color: Colors.white, fontSize: 14, fontWeight: FontWeight.w700)),
                      ),
                    ]),

                    SizedBox(height: size.height * 0.05),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  // ── User Type Chip ─────────────────────────────────────────────────────────
  Widget _typeChip(String value, String emoji, String label) {
    final selected = _userType == value;
    return Expanded(
      child: GestureDetector(
        onTap: () => setState(() => _userType = value),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          padding: const EdgeInsets.symmetric(vertical: 12),
          decoration: BoxDecoration(
            color: selected ? Colors.white.withOpacity(0.2) : Colors.white.withOpacity(0.08),
            borderRadius: BorderRadius.circular(14),
            border: Border.all(
              color: selected ? Colors.white.withOpacity(0.5) : Colors.white.withOpacity(0.12),
              width: 1.5,
            ),
          ),
          child: Column(children: [
            Text(emoji, style: const TextStyle(fontSize: 20)),
            const SizedBox(height: 4),
            Text(label, style: TextStyle(
              color: selected ? Colors.white : Colors.white.withOpacity(0.6),
              fontSize: 11, fontWeight: FontWeight.w700,
            )),
          ]),
        ),
      ),
    );
  }

  // ── Risk Profile Chip ──────────────────────────────────────────────────────
  Widget _riskChip(String value, String emoji, String label) {
    final selected = _riskProfile == value;
    return Expanded(
      child: GestureDetector(
        onTap: () => setState(() => _riskProfile = value),
        child: AnimatedContainer(
          duration: const Duration(milliseconds: 200),
          padding: const EdgeInsets.symmetric(vertical: 14),
          decoration: BoxDecoration(
            color: selected ? Colors.white.withOpacity(0.2) : Colors.white.withOpacity(0.08),
            borderRadius: BorderRadius.circular(14),
            border: Border.all(
              color: selected ? Colors.white.withOpacity(0.5) : Colors.white.withOpacity(0.12),
              width: 1.5,
            ),
          ),
          child: Row(mainAxisAlignment: MainAxisAlignment.center, children: [
            Text(emoji, style: const TextStyle(fontSize: 16)),
            const SizedBox(width: 6),
            Text(label, style: TextStyle(
              color: selected ? Colors.white : Colors.white.withOpacity(0.6),
              fontSize: 13, fontWeight: FontWeight.w700,
            )),
          ]),
        ),
      ),
    );
  }

  // ── Glass Field ────────────────────────────────────────────────────────────
  Widget _glassField({
    required TextEditingController controller,
    required FocusNode focusNode,
    required IconData icon,
    required String hint,
    TextInputType? keyboardType,
    FocusNode? nextFocus,
    bool obscure = false,
    VoidCallback? onToggleObscure,
    VoidCallback? onSubmit,
  }) {
    return TextField(
      controller: controller,
      focusNode: focusNode,
      obscureText: obscure,
      keyboardType: keyboardType,
      textInputAction: nextFocus != null ? TextInputAction.next : TextInputAction.done,
      onSubmitted: (_) {
        if (nextFocus != null) {
          nextFocus.requestFocus();
        } else if (onSubmit != null) {
          onSubmit();
        }
      },
      style: const TextStyle(color: Color(0xFF1E293B), fontWeight: FontWeight.w600, fontSize: 15),
      decoration: InputDecoration(
        prefixIcon: Icon(icon, color: const Color(0xFF64748B), size: 20),
        suffixIcon: onToggleObscure != null
            ? IconButton(
                icon: Icon(
                  obscure ? Icons.visibility_outlined : Icons.visibility_off_outlined,
                  color: const Color(0xFF64748B), size: 20,
                ),
                onPressed: onToggleObscure,
              )
            : null,
        hintText: hint,
        hintStyle: const TextStyle(color: Color(0xFF94A3B8), fontWeight: FontWeight.w400, fontSize: 15),
        filled: true,
        fillColor: const Color(0xFFF8FAFC),
        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: BorderSide.none,
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: const BorderSide(color: Color(0xFF1E3A8A), width: 1.5),
        ),
      ),
    );
  }

  // ── Password Strength ──────────────────────────────────────────────────────
  Widget _passwordStrengthBar() {
    Color barColor;
    String label;
    if (_passwordStrength >= 0.8) {
      barColor = const Color(0xFF34D399);
      label = 'Strong';
    } else if (_passwordStrength >= 0.5) {
      barColor = const Color(0xFFFBBF24);
      label = 'Medium';
    } else if (_passwordStrength > 0) {
      barColor = const Color(0xFFF87171);
      label = 'Weak';
    } else {
      return const SizedBox();
    }

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        ClipRRect(
          borderRadius: BorderRadius.circular(4),
          child: LinearProgressIndicator(
            value: _passwordStrength,
            backgroundColor: Colors.white.withOpacity(0.1),
            valueColor: AlwaysStoppedAnimation(barColor),
            minHeight: 4,
          ),
        ),
        const SizedBox(height: 4),
        Text(label, style: TextStyle(color: barColor, fontSize: 11, fontWeight: FontWeight.w600)),
      ],
    );
  }
}
