import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'screens/login_screen.dart';
import 'providers/financial_provider.dart';
import 'providers/auth_provider.dart';

void main() {
  runApp(
    MultiProvider(
      providers: [
        ChangeNotifierProvider(create: (_) => AuthProvider()),
        ChangeNotifierProvider(create: (_) => FinancialProvider()),
      ],
      child: const HorizonApp(),
    ),
  );
}

class HorizonApp extends StatelessWidget {
  const HorizonApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Horizon',
      debugShowCheckedModeBanner: false,
      theme: ThemeData(
        // Soft Pastel Palette
        colorScheme: ColorScheme.fromSeed(
          seedColor: const Color(0xFFA2A8D3), // Pastel Lavender/Indigo
          primary: const Color(0xFF909AC6),
          secondary: const Color(0xFFF1B7B0), // Pastel Peach/Pink
          background: const Color(0xFFF7F9FC), // Very soft cool white
          surface: Colors.white,
        ),
        scaffoldBackgroundColor: const Color(0xFFF7F9FC),
        useMaterial3: true,
        textTheme: GoogleFonts.interTextTheme(
          Theme.of(context).textTheme,
        ).apply(
          bodyColor: const Color(0xFF2D3748),
          displayColor: const Color(0xFF2D3748),
        ),
        appBarTheme: const AppBarTheme(
          backgroundColor: Colors.transparent,
          elevation: 0,
          centerTitle: true,
          iconTheme: IconThemeData(color: Color(0xFF2D3748)),
          titleTextStyle: TextStyle(color: Color(0xFF2D3748), fontSize: 20, fontWeight: FontWeight.w600),
        ),
        inputDecorationTheme: InputDecorationTheme(
          filled: true,
          fillColor: Colors.white,

          // Label floats above on focus so user always sees what field they're in
          floatingLabelBehavior: FloatingLabelBehavior.auto,
          labelStyle: const TextStyle(
            color: Color(0xFF64748B),
            fontSize: 15,
            fontWeight: FontWeight.w500,
          ),
          floatingLabelStyle: const TextStyle(
            color: Color(0xFF1E3A8A),
            fontSize: 13,
            fontWeight: FontWeight.w600,
          ),
          hintStyle: const TextStyle(
            color: Color(0xFFB0BCC8),
            fontSize: 15,
          ),

          // Rest state — visible light grey border (1.5px)
          enabledBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(14),
            borderSide: const BorderSide(color: Color(0xFFCBD5E1), width: 1.5),
          ),
          // Focus state — strong navy blue border (2px)
          focusedBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(14),
            borderSide: const BorderSide(color: Color(0xFF1E3A8A), width: 2.0),
          ),
          errorBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(14),
            borderSide: const BorderSide(color: Color(0xFFEF4444), width: 1.5),
          ),
          focusedErrorBorder: OutlineInputBorder(
            borderRadius: BorderRadius.circular(14),
            borderSide: const BorderSide(color: Color(0xFFEF4444), width: 2.0),
          ),

          // Tall comfortable tap target
          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 18),

          prefixIconColor: Color(0xFF1E3A8A),
          suffixIconColor: Color(0xFF64748B),
        ),
        elevatedButtonTheme: ElevatedButtonThemeData(
          style: ElevatedButton.styleFrom(
            backgroundColor: const Color(0xFF1E3A8A),
            foregroundColor: Colors.white,
            elevation: 0,
            minimumSize: const Size(double.infinity, 52),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
            textStyle: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w700,
              letterSpacing: 0.2,
            ),
          ),
        ),

      ),
      home: const LoginScreen(),
    );
  }
}
