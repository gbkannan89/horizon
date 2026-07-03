import 'package:flutter/material.dart';

class AppColors {
  // Artha - Teal (Primary / Brand)
  static const Color teal50 = Color(0xFFF0FDF9);
  static const Color teal100 = Color(0xFFCCFBEF); // from image, standard is CCFBF1
  static const Color teal400 = Color(0xFF2DD4BF);
  static const Color teal500 = Color(0xFF14B8A6);
  static const Color teal600 = Color(0xFF0D9488);
  static const Color teal700 = Color(0xFF0F766E);
  static const Color teal900 = Color(0xFF134E4A);

  // Ink - Deep Navy (Base / Trust)
  static const Color navy700 = Color(0xFF1E293B);
  static const Color navy800 = Color(0xFF111A2E);
  static const Color navy900 = Color(0xFF0B1220);

  // Gold - Amber (Accent / Wealth)
  static const Color amber400 = Color(0xFFFBBF24);
  static const Color amber500 = Color(0xFFF59E0B);
  static const Color amber600 = Color(0xFFD97706);

  // Slate - Neutrals
  static const Color slate50 = Color(0xFFF8FAFC);
  static const Color slate100 = Color(0xFFF1F5F9);
  static const Color slate200 = Color(0xFFE2E8F0);
  static const Color slate500 = Color(0xFF64748B);
  static const Color slate600 = Color(0xFF475569);

  // Status Colors (Reds)
  static const Color red500 = Color(0xFFEF4444);
  static const Color red600 = Color(0xFFDC2626);
  static const Color red700 = Color(0xFFB91C1C);

  // Semantic Mappings - Light Mode
  static const Color lightBackground = slate50;
  static const Color lightSurface = Colors.white;
  static const Color lightSurfaceElevated = slate100;
  static const Color lightTextPrimary = navy900;
  static const Color lightTextSecondary = slate600;
  static const Color lightBorder = slate200;

  // Semantic Mappings - Dark Mode
  static const Color darkBackground = navy900;
  static const Color darkSurface = navy800;
  static const Color darkSurfaceElevated = navy700;
  static const Color darkTextPrimary = slate50;
  static const Color darkTextSecondary = slate200; 
  static const Color darkBorder = navy700;

  // Gradients
  static const LinearGradient primaryGradient = LinearGradient(
    colors: [teal500, teal700],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
  
  static const LinearGradient darkCardGradient = LinearGradient(
    colors: [navy800, navy900],
    begin: Alignment.topLeft,
    end: Alignment.bottomRight,
  );
}

class AppShadows {
  static final List<BoxShadow> soft = [
    BoxShadow(
      color: Colors.black.withOpacity(0.04),
      blurRadius: 10,
      offset: const Offset(0, 4),
    )
  ];
  
  static final List<BoxShadow> medium = [
    BoxShadow(
      color: Colors.black.withOpacity(0.08),
      blurRadius: 20,
      offset: const Offset(0, 8),
    )
  ];
}

class AppSpacing {
  static const double xs = 4.0;
  static const double sm = 8.0;
  static const double md = 16.0;
  static const double lg = 24.0;
  static const double xl = 32.0;
}

class AppRadius {
  static const double xs = 4.0;
  static const double sm = 8.0;
  static const double md = 12.0;
  static const double lg = 16.0;
  static const double xl = 24.0;
  static const double pill = 999.0;
}
