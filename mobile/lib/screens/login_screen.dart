import 'dart:math' as math;
import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import 'main_shell.dart';
import 'signup_screen.dart';
import '../utils/ui_utils.dart';

class LoginScreen extends StatefulWidget {
  const LoginScreen({super.key});

  @override
  State<LoginScreen> createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> with SingleTickerProviderStateMixin {
  final _emailCtrl    = TextEditingController();
  final _passwordCtrl = TextEditingController();
  bool _obscurePassword = true;
  late final AnimationController _animCtrl;
  late final Animation<double> _fadeAnim;
  late final Animation<Offset> _slideAnim;

  @override
  void initState() {
    super.initState();
    _animCtrl = AnimationController(vsync: this, duration: const Duration(milliseconds: 800));
    _fadeAnim  = CurvedAnimation(parent: _animCtrl, curve: Curves.easeOut);
    _slideAnim = Tween<Offset>(begin: const Offset(0, 0.08), end: Offset.zero)
        .animate(CurvedAnimation(parent: _animCtrl, curve: Curves.easeOutCubic));
    _animCtrl.forward();
  }

  @override
  void dispose() {
    _animCtrl.dispose();
    _emailCtrl.dispose();
    _passwordCtrl.dispose();
    super.dispose();
  }

  Future<void> _handleLogin() async {
    try {
      final auth = Provider.of<AuthProvider>(context, listen: false);
      final success = await auth.login(_emailCtrl.text.trim(), _passwordCtrl.text);
      if (success) {
        if (mounted) {
          Navigator.of(context).pushReplacement(
            MaterialPageRoute(builder: (_) => const MainShell()),
          );
        }
      } else {
        UiUtils.showSnack(context, 'Invalid email or password', isError: true);
      }
    } catch (e) {
      if (mounted) UiUtils.showSnack(context, 'Login failed: ${e.toString()}', isError: true);
    }
  }

  @override
  Widget build(BuildContext context) {
    final isLoading = context.watch<AuthProvider>().isLoading;
    final size = MediaQuery.of(context).size;

    return Scaffold(
      body: Stack(
        children: [
          // ── Gradient Background ──────────────────────────────────────────
          Positioned.fill(
            child: DecoratedBox(
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    Color(0xFFF0FDFA), // teal-50
                    Color(0xFFF8FAFC), // slate-50
                    Color(0xFFF5F3FF), // violet-50
                  ],
                ),
              ),
            ),
          ),

          // ── Decorative Blobs ─────────────────────────────────────────────
          Positioned(
            top: -size.height * 0.12,
            right: -size.width * 0.2,
            child: Container(
              width: size.width * 0.7,
              height: size.width * 0.7,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFF0D9488).withValues(alpha: 0.20),
                    const Color(0xFF0D9488).withValues(alpha: 0.08),
                    const Color(0xFF0D9488).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            bottom: -size.height * 0.08,
            left: -size.width * 0.15,
            child: Container(
              width: size.width * 0.6,
              height: size.width * 0.6,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFF909AC6).withValues(alpha: 0.18),
                    const Color(0xFF909AC6).withValues(alpha: 0.06),
                    const Color(0xFF909AC6).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            top: size.height * 0.25,
            right: -size.width * 0.1,
            child: Container(
              width: size.width * 0.35,
              height: size.width * 0.35,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFFFBBF24).withValues(alpha: 0.12),
                    const Color(0xFFFBBF24).withValues(alpha: 0.04),
                    const Color(0xFFFBBF24).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),

          // ── Subtle Grid Pattern Overlay ─────────────────────────────────
          Positioned.fill(
            child: CustomPaint(
              painter: _GridPainter(),
            ),
          ),

          // ── Main Content ─────────────────────────────────────────────────
          SafeArea(
            child: Center(
              child: SingleChildScrollView(
                padding: const EdgeInsets.symmetric(horizontal: 28),
                child: FadeTransition(
                  opacity: _fadeAnim,
                  child: SlideTransition(
                    position: _slideAnim,
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        const SizedBox(height: 20),

                        // Logo
                        Container(
                          width: 80,
                          height: 80,
                          decoration: BoxDecoration(
                            shape: BoxShape.circle,
                            color: Colors.white,
                            boxShadow: [
                              BoxShadow(
                                color: Colors.black.withValues(alpha: 0.06),
                                blurRadius: 20,
                                offset: const Offset(0, 8),
                              ),
                            ],
                          ),
                          child: ClipOval(
                            child: Image.asset('assets/logo.png', fit: BoxFit.cover),
                          ),
                        ),
                        const SizedBox(height: 20),

                        const Text('Horizon', style: TextStyle(
                          color: Color(0xFF0F172A), fontSize: 32,
                          fontWeight: FontWeight.w900, letterSpacing: -1.5,
                        )),
                        const SizedBox(height: 6),
                        const Text('Your smart financial companion', style: TextStyle(
                          color: Color(0xFF64748B), fontSize: 13,
                          fontWeight: FontWeight.w500,
                        )),

                        const SizedBox(height: 28),

                        // ── Glassmorphism Card ─────────────────────────────
                        ClipRRect(
                          borderRadius: BorderRadius.circular(24),
                          child: BackdropFilter(
                            filter: ui.ImageFilter.blur(sigmaX: 25, sigmaY: 25),
                            child: Container(
                              padding: const EdgeInsets.all(28),
                              decoration: BoxDecoration(
                                color: Colors.white.withValues(alpha: 0.60),
                                borderRadius: BorderRadius.circular(24),
                                border: Border.all(
                                  color: Colors.white.withValues(alpha: 0.7),
                                  width: 1.5,
                                ),
                                boxShadow: [
                                  BoxShadow(
                                    color: Colors.black.withValues(alpha: 0.06),
                                    blurRadius: 40,
                                    offset: const Offset(0, 10),
                                  ),
                                ],
                              ),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text('Sign In', style: TextStyle(
                                    fontSize: 22, fontWeight: FontWeight.w900,
                                    color: Color(0xFF0F172A), letterSpacing: -0.5,
                                  )),
                                  const SizedBox(height: 4),
                                  const Text('Enter details to manage your wealth', style: TextStyle(
                                    fontSize: 13, color: Color(0xFF64748B),
                                  )),
                                  const SizedBox(height: 24),

                                  // Email
                                  TextField(
                                    controller: _emailCtrl,
                                    keyboardType: TextInputType.emailAddress,
                                    textInputAction: TextInputAction.next,
                                    style: const TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF0F172A)),
                                    decoration: InputDecoration(
                                      hintText: 'Email address',
                                      hintStyle: const TextStyle(color: Color(0xFF94A3B8), fontWeight: FontWeight.w400),
                                      prefixIcon: const Icon(Icons.email_outlined, color: Color(0xFF94A3B8)),
                                      filled: true,
                                      fillColor: Colors.white.withValues(alpha: 0.70),
                                      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                                      border: OutlineInputBorder(
                                        borderRadius: BorderRadius.circular(12),
                                        borderSide: BorderSide(color: Colors.grey.shade200),
                                      ),
                                      enabledBorder: OutlineInputBorder(
                                        borderRadius: BorderRadius.circular(12),
                                        borderSide: BorderSide(color: Colors.grey.shade200),
                                      ),
                                      focusedBorder: OutlineInputBorder(
                                        borderRadius: BorderRadius.circular(12),
                                        borderSide: const BorderSide(color: Color(0xFF0D9488), width: 1.5),
                                      ),
                                    ),
                                  ),
                                  const SizedBox(height: 16),

                                  // Password
                                  TextField(
                                    controller: _passwordCtrl,
                                    obscureText: _obscurePassword,
                                    textInputAction: TextInputAction.done,
                                    onSubmitted: (_) => isLoading ? null : _handleLogin(),
                                    style: const TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF0F172A)),
                                    decoration: InputDecoration(
                                      hintText: 'Password',
                                      hintStyle: const TextStyle(color: Color(0xFF94A3B8), fontWeight: FontWeight.w400),
                                      prefixIcon: const Icon(Icons.lock_outline_rounded, color: Color(0xFF94A3B8)),
                                      suffixIcon: IconButton(
                                        icon: Icon(
                                          _obscurePassword ? Icons.visibility_outlined : Icons.visibility_off_outlined,
                                          color: const Color(0xFF94A3B8),
                                        ),
                                        onPressed: () => setState(() => _obscurePassword = !_obscurePassword),
                                      ),
                                      filled: true,
                                      fillColor: Colors.white.withValues(alpha: 0.70),
                                      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
                                      border: OutlineInputBorder(
                                        borderRadius: BorderRadius.circular(12),
                                        borderSide: BorderSide(color: Colors.grey.shade200),
                                      ),
                                      enabledBorder: OutlineInputBorder(
                                        borderRadius: BorderRadius.circular(12),
                                        borderSide: BorderSide(color: Colors.grey.shade200),
                                      ),
                                      focusedBorder: OutlineInputBorder(
                                        borderRadius: BorderRadius.circular(12),
                                        borderSide: const BorderSide(color: Color(0xFF0D9488), width: 1.5),
                                      ),
                                    ),
                                  ),
                                  const SizedBox(height: 8),

                                  Align(
                                    alignment: Alignment.centerRight,
                                    child: TextButton(
                                      onPressed: () {},
                                      style: TextButton.styleFrom(
                                        foregroundColor: const Color(0xFF64748B),
                                        padding: EdgeInsets.zero,
                                      ),
                                      child: const Text('Forgot password?', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 12)),
                                    ),
                                  ),
                                  const SizedBox(height: 16),

                                  SizedBox(
                                    width: double.infinity, height: 52,
                                    child: ElevatedButton(
                                      onPressed: isLoading ? null : _handleLogin,
                                      style: ElevatedButton.styleFrom(
                                        backgroundColor: const Color(0xFF0D9488),
                                        foregroundColor: Colors.white,
                                        elevation: 0,
                                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                                      ),
                                      child: isLoading
                                          ? const SizedBox(width: 20, height: 20,
                                              child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2.5))
                                          : const Text('Sign In', style: TextStyle(fontSize: 15, fontWeight: FontWeight.w900)),
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),

                        const SizedBox(height: 24),

                        Row(mainAxisAlignment: MainAxisAlignment.center, children: [
                          const Text("Don't have an account? ", style: TextStyle(color: Color(0xFF64748B), fontSize: 13)),
                          GestureDetector(
                            onTap: () {
                              Navigator.of(context).push(
                                PageRouteBuilder(
                                  pageBuilder: (_, _, _) => const SignupScreen(),
                                  transitionsBuilder: (_, anim, _, child) => FadeTransition(opacity: anim, child: child),
                                  transitionDuration: const Duration(milliseconds: 300),
                                ),
                              );
                            },
                            child: const Text('Sign Up', style: TextStyle(
                              color: Color(0xFF0D9488), fontSize: 13, fontWeight: FontWeight.bold,
                            )),
                          ),
                        ]),

                        const SizedBox(height: 24),
                      ],
                    ),
                  ),
                ),
              ),
            ),
          ),
        ],
      ),
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
