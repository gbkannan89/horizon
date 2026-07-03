import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../core/theme/design_tokens.dart';

class AppSpacing {
  static const double xs = 4;
  static const double sm = 8;
  static const double md = 16;
  static const double lg = 24;
  static const double xl = 32;
  static const double xxl = 48;
}

class AppRadius {
  static const double sm = 8;
  static const double md = 12;
  static const double lg = 16;
  static const double xl = 24;
}

class AppElevation {
  static const double low = 1;
  static const double medium = 2;
  static const double high = 4;
}

class AppIconSize {
  static const double sm = 16;
  static const double md = 20;
  static const double lg = 24;
  static const double xl = 32;
}

final lightTheme = ThemeData(
  useMaterial3: true,
  colorSchemeSeed: AppColors.teal500,
  brightness: Brightness.light,
  scaffoldBackgroundColor: AppColors.lightBackground,
  textTheme: GoogleFonts.interTextTheme(ThemeData.light().textTheme),
  appBarTheme: const AppBarTheme(
    centerTitle: true, 
    elevation: 0,
    backgroundColor: AppColors.lightBackground,
    foregroundColor: AppColors.lightTextPrimary,
  ),
  cardTheme: const CardThemeData(clipBehavior: Clip.antiAlias),
  inputDecorationTheme: InputDecorationTheme(
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
    filled: true,
  ),
);

final darkTheme = ThemeData(
  useMaterial3: true,
  colorSchemeSeed: AppColors.teal500,
  brightness: Brightness.dark,
  scaffoldBackgroundColor: AppColors.darkBackground,
  textTheme: GoogleFonts.interTextTheme(ThemeData.dark().textTheme),
  appBarTheme: const AppBarTheme(
    centerTitle: true, 
    elevation: 0,
    backgroundColor: AppColors.darkBackground,
    foregroundColor: AppColors.darkTextPrimary,
  ),
  cardTheme: CardThemeData(
    clipBehavior: Clip.antiAlias,
    color: AppColors.darkSurface,
  ),
  inputDecorationTheme: InputDecorationTheme(
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
    filled: true,
  ),
);
