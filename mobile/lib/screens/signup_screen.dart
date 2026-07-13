import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import 'main_shell.dart';
import '../utils/ui_utils.dart';

class SignupScreen extends StatefulWidget {
  const SignupScreen({super.key});

  @override
  State<SignupScreen> createState() => _SignupScreenState();
}

class _SignupScreenState extends State<SignupScreen>
    with SingleTickerProviderStateMixin {
  final _nameCtrl = TextEditingController();
  final _emailCtrl = TextEditingController();
  final _phoneCtrl = TextEditingController();
  final _passwordCtrl = TextEditingController();
  final _inviteCtrl = TextEditingController();
  bool _obscurePassword = true;
  double _passwordStrength = 0;

  late final AnimationController _animCtrl;
  late final Animation<double> _fadeAnim;
  late final Animation<Offset> _slideAnim;

  final _focusName = FocusNode();
  final _focusEmail = FocusNode();
  final _focusPhone = FocusNode();
  final _focusPassword = FocusNode();
  final _focusInvite = FocusNode();

  static const _textDark = Color(0xFF0F172A);
  static const _textMuted = Color(0xFF64748B);
  static const _textHint = Color(0xFF94A3B8);
  static const _accent = Color(0xFF0D9488);
  static const _violet = Color(0xFF909AC6);

  @override
  void initState() {
    super.initState();
    _animCtrl = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1000),
    );
    _fadeAnim = CurvedAnimation(parent: _animCtrl, curve: Curves.easeOut);
    _slideAnim = Tween<Offset>(
      begin: const Offset(0, 0.12),
      end: Offset.zero,
    ).animate(CurvedAnimation(parent: _animCtrl, curve: Curves.easeOutCubic));
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
    if (RegExp(r'[!@#\$%^&*(),.?":{}|<>]').hasMatch(p)) s += 0.15;
    setState(() => _passwordStrength = s.clamp(0.0, 1.0));
  }

  @override
  void dispose() {
    _animCtrl.dispose();
    _nameCtrl.dispose();
    _emailCtrl.dispose();
    _phoneCtrl.dispose();
    _passwordCtrl.dispose();
    _inviteCtrl.dispose();
    _focusName.dispose();
    _focusEmail.dispose();
    _focusPhone.dispose();
    _focusPassword.dispose();
    _focusInvite.dispose();
    super.dispose();
  }

  Future<void> _handleSignup() async {
    if (_nameCtrl.text.trim().isEmpty ||
        _emailCtrl.text.trim().isEmpty ||
        _passwordCtrl.text.isEmpty) {
      UiUtils.showSnack(
        context,
        'Please fill in name, email and password',
        isError: true,
      );
      return;
    }
    try {
      final auth = Provider.of<AuthProvider>(context, listen: false);
      final inviteCode = _inviteCtrl.text.trim();
      final result = await auth.register(
        _nameCtrl.text.trim(),
        _emailCtrl.text.trim(),
        _passwordCtrl.text,
        phone: _phoneCtrl.text.trim(),
        inviteCode: inviteCode.isNotEmpty ? inviteCode : null,
      );
      if (!mounted) return;
      if (result['success'] == true) {
        Navigator.of(
          context,
        ).pushReplacement(MaterialPageRoute(builder: (_) => const MainShell()));
      } else if (result['detail'] is Map) {
        final detail = result['detail'] as Map;
        final confirmed = await _showPendingInviteDialog(
          invitedBy: detail['invitedBy'] ?? detail['householdName'] ?? '',
          householdName: detail['householdName'] ?? '',
          inviteCode: detail['inviteCode'] ?? '',
        );
        if (confirmed && mounted) {
          final retryResult = await auth.register(
            _nameCtrl.text.trim(),
            _emailCtrl.text.trim(),
            _passwordCtrl.text,
            phone: _phoneCtrl.text.trim(),
            inviteCode: detail['inviteCode'] ?? '',
          );
          if (mounted && retryResult['success'] == true) {
            Navigator.of(context).pushReplacement(
              MaterialPageRoute(builder: (_) => const MainShell()),
            );
          } else {
            UiUtils.showSnack(
              context,
              'Registration failed. Please try again.',
              isError: true,
            );
          }
        }
      } else {
        final msg = result['detail'] is String
            ? result['detail']
            : 'Registration failed or email already exists';
        UiUtils.showSnack(context, msg as String, isError: true);
      }
    } catch (e) {
      if (mounted)
        UiUtils.showSnack(
          context,
          'Registration failed: ${e.toString()}',
          isError: true,
        );
    }
  }

  Future<bool> _showPendingInviteDialog({
    required String invitedBy,
    required String householdName,
    required String inviteCode,
  }) async {
    return await showDialog<bool>(
          context: context,
          builder: (ctx) => AlertDialog(
            backgroundColor: Colors.white,
            shape: RoundedRectangleBorder(
              borderRadius: BorderRadius.circular(24),
            ),
            title: const Row(
              children: [
                Icon(Icons.family_restroom_rounded, color: _accent),
                SizedBox(width: 8),
                Text(
                  'Invitation Found',
                  style: TextStyle(
                    color: _textDark,
                    fontWeight: FontWeight.bold,
                  ),
                ),
              ],
            ),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '$invitedBy has invited you to join their household "$householdName".',
                  style: const TextStyle(color: _textMuted, height: 1.5),
                ),
                const SizedBox(height: 12),
                Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(
                    color: const Color(0xFFF0FDFA),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Row(
                    children: [
                      const Icon(
                        Icons.vpn_key_rounded,
                        color: _accent,
                        size: 20,
                      ),
                      const SizedBox(width: 8),
                      Text(
                        'Code: $inviteCode',
                        style: const TextStyle(
                          color: _textDark,
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 12),
                Text(
                  'By joining, you will be able to see shared finances and collaborate.',
                  style: TextStyle(
                    color: Colors.black.withValues(alpha: 0.5),
                    fontSize: 12,
                    height: 1.4,
                  ),
                ),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.pop(ctx, false),
                child: const Text(
                  'Create New Household',
                  style: TextStyle(color: _textMuted),
                ),
              ),
              ElevatedButton(
                onPressed: () => Navigator.pop(ctx, true),
                style: ElevatedButton.styleFrom(
                  backgroundColor: _accent,
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                ),
                child: const Text('Join Household'),
              ),
            ],
          ),
        ) ??
        false;
  }

  @override
  Widget build(BuildContext context) {
    final isLoading = context.watch<AuthProvider>().isLoading;
    final size = MediaQuery.of(context).size;

    return Scaffold(
      body: Stack(
        children: [
          Positioned.fill(
            child: const DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    Color(0xFFF0FDFA),
                    Color(0xFFF8FAFC),
                    Color(0xFFF5F3FF),
                  ],
                ),
              ),
            ),
          ),

          Positioned(
            top: -size.height * 0.1,
            right: -size.width * 0.18,
            child: Container(
              width: size.width * 0.65,
              height: size.width * 0.65,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFF0D9488).withValues(alpha: 0.18),
                    const Color(0xFF0D9488).withValues(alpha: 0.06),
                    const Color(0xFF0D9488).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            bottom: -size.height * 0.1,
            left: -size.width * 0.15,
            child: Container(
              width: size.width * 0.55,
              height: size.width * 0.55,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFF909AC6).withValues(alpha: 0.16),
                    const Color(0xFF909AC6).withValues(alpha: 0.05),
                    const Color(0xFF909AC6).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            top: size.height * 0.35,
            right: -size.width * 0.08,
            child: Container(
              width: size.width * 0.3,
              height: size.width * 0.3,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFFFBBF24).withValues(alpha: 0.10),
                    const Color(0xFFFBBF24).withValues(alpha: 0.03),
                    const Color(0xFFFBBF24).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),

          Positioned.fill(child: CustomPaint(painter: _GridPainter())),

          SafeArea(
            child: SingleChildScrollView(
              padding: const EdgeInsets.symmetric(horizontal: 24),
              child: FadeTransition(
                opacity: _fadeAnim,
                child: SlideTransition(
                  position: _slideAnim,
                  child: Column(
                    children: [
                      SizedBox(height: size.height * 0.05),
                      Container(
                        width: 90,
                        height: 90,
                        decoration: BoxDecoration(
                          color: Colors.white,
                          shape: BoxShape.circle,
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withValues(alpha: 0.06),
                              blurRadius: 20,
                              offset: const Offset(0, 8),
                            ),
                          ],
                        ),
                        child: ClipOval(
                          child: Padding(
                            padding: const EdgeInsets.all(4),
                            child: Image.asset(
                              'assets/logo.png',
                              fit: BoxFit.cover,
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(height: 16),
                      const Text(
                        'Horizon',
                        style: TextStyle(
                          color: _textDark,
                          fontSize: 28,
                          fontWeight: FontWeight.w900,
                          letterSpacing: -1,
                        ),
                      ),
                      const SizedBox(height: 4),
                      const Text(
                        'Start your wealth journey',
                        style: TextStyle(
                          color: _textMuted,
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      SizedBox(height: size.height * 0.035),

                      ClipRRect(
                        borderRadius: BorderRadius.circular(28),
                        child: BackdropFilter(
                          filter: ui.ImageFilter.blur(sigmaX: 25, sigmaY: 25),
                          child: Container(
                            padding: const EdgeInsets.all(24),
                            decoration: BoxDecoration(
                              color: Colors.white.withValues(alpha: 0.60),
                              border: Border.all(
                                color: Colors.white.withValues(alpha: 0.7),
                                width: 1.5,
                              ),
                              borderRadius: BorderRadius.circular(28),
                              boxShadow: [
                                BoxShadow(
                                  color: Colors.black.withValues(alpha: 0.05),
                                  blurRadius: 40,
                                  offset: const Offset(0, 10),
                                ),
                              ],
                            ),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                _glassField(
                                  controller: _nameCtrl,
                                  focusNode: _focusName,
                                  icon: Icons.person_outline_rounded,
                                  hint: 'Full name',
                                  nextFocus: _focusEmail,
                                ),
                                const SizedBox(height: 14),
                                _glassField(
                                  controller: _emailCtrl,
                                  focusNode: _focusEmail,
                                  icon: Icons.email_outlined,
                                  hint: 'Email address',
                                  keyboardType: TextInputType.emailAddress,
                                  nextFocus: _focusPhone,
                                ),
                                const SizedBox(height: 14),
                                _glassField(
                                  controller: _phoneCtrl,
                                  focusNode: _focusPhone,
                                  icon: Icons.phone_outlined,
                                  hint: 'Phone (optional)',
                                  keyboardType: TextInputType.phone,
                                  nextFocus: _focusPassword,
                                ),
                                const SizedBox(height: 14),
                                _glassField(
                                  controller: _passwordCtrl,
                                  focusNode: _focusPassword,
                                  icon: Icons.lock_outline_rounded,
                                  hint: 'Create a password',
                                  obscure: _obscurePassword,
                                  onToggleObscure: () => setState(
                                    () => _obscurePassword = !_obscurePassword,
                                  ),
                                  onSubmit: isLoading ? null : _handleSignup,
                                ),
                                if (_passwordCtrl.text.isNotEmpty) ...[
                                  const SizedBox(height: 8),
                                  _passwordStrengthBar(),
                                ],
                                const SizedBox(height: 14),
                                _glassField(
                                  controller: _inviteCtrl,
                                  focusNode: _focusInvite,
                                  icon: Icons.family_restroom_rounded,
                                  hint: 'Family invite code (optional)',
                                  textCapitalization:
                                      TextCapitalization.characters,
                                ),
                                const SizedBox(height: 8),
                                Text(
                                  'By creating an account, you agree to our Terms and Privacy Policy.',
                                  style: TextStyle(
                                    color: Colors.black.withValues(alpha: 0.4),
                                    fontSize: 11,
                                  ),
                                ),
                                const SizedBox(height: 16),
                                SizedBox(
                                  width: double.infinity,
                                  height: 54,
                                  child: ElevatedButton(
                                    onPressed: isLoading ? null : _handleSignup,
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor: _accent,
                                      foregroundColor: Colors.white,
                                      elevation: 0,
                                      shape: RoundedRectangleBorder(
                                        borderRadius: BorderRadius.circular(16),
                                      ),
                                    ),
                                    child: isLoading
                                        ? const SizedBox(
                                            width: 22,
                                            height: 22,
                                            child: CircularProgressIndicator(
                                              strokeWidth: 2.5,
                                              color: Colors.white,
                                            ),
                                          )
                                        : const Text(
                                            'Create Account',
                                            style: TextStyle(
                                              fontSize: 16,
                                              fontWeight: FontWeight.w800,
                                              letterSpacing: 0.3,
                                            ),
                                          ),
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      ),

                      const SizedBox(height: 24),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const Text(
                            "Already have an account? ",
                            style: TextStyle(color: _textMuted, fontSize: 14),
                          ),
                          GestureDetector(
                            onTap: () => Navigator.of(context).pop(),
                            child: const Text(
                              'Sign In',
                              style: TextStyle(
                                color: _accent,
                                fontSize: 14,
                                fontWeight: FontWeight.w700,
                              ),
                            ),
                          ),
                        ],
                      ),
                      SizedBox(height: size.height * 0.05),
                    ],
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

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
    TextCapitalization textCapitalization = TextCapitalization.none,
  }) {
    return TextField(
      controller: controller,
      focusNode: focusNode,
      obscureText: obscure,
      keyboardType: keyboardType,
      textCapitalization: textCapitalization,
      textInputAction: nextFocus != null
          ? TextInputAction.next
          : TextInputAction.done,
      onSubmitted: (_) {
        if (nextFocus != null) {
          nextFocus.requestFocus();
        } else if (onSubmit != null)
          onSubmit();
      },
      style: const TextStyle(
        color: _textDark,
        fontWeight: FontWeight.w600,
        fontSize: 15,
      ),
      decoration: InputDecoration(
        prefixIcon: Icon(icon, color: _textHint, size: 20),
        suffixIcon: onToggleObscure != null
            ? IconButton(
                icon: Icon(
                  obscure
                      ? Icons.visibility_outlined
                      : Icons.visibility_off_outlined,
                  color: _textHint,
                  size: 20,
                ),
                onPressed: onToggleObscure,
              )
            : null,
        hintText: hint,
        hintStyle: const TextStyle(
          color: _textHint,
          fontWeight: FontWeight.w400,
          fontSize: 15,
        ),
        filled: true,
        fillColor: Colors.white.withValues(alpha: 0.70),
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 16,
          vertical: 16,
        ),
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: BorderSide(color: Colors.grey.shade200),
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: BorderSide(color: Colors.grey.shade200),
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(16),
          borderSide: const BorderSide(color: _accent, width: 1.5),
        ),
      ),
    );
  }

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
            backgroundColor: Colors.grey.shade200,
            valueColor: AlwaysStoppedAnimation(barColor),
            minHeight: 4,
          ),
        ),
        const SizedBox(height: 4),
        Text(
          label,
          style: TextStyle(
            color: barColor,
            fontSize: 11,
            fontWeight: FontWeight.w600,
          ),
        ),
      ],
    );
  }
}

// ── Grid Pattern Painter ─────────────────────────────────────────────────────
class _GridPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = const Color(0xFF0D9488).withValues(alpha: 0.035)
      ..strokeWidth = 0.5;

    const spacing = 40.0;
    for (double x = 0; x < size.width; x += spacing) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
    }
    for (double y = 0; y < size.height; y += spacing) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
