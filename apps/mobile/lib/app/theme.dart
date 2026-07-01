import 'package:flutter/material.dart';

const _seedColor = Color(0xFF00897B);

class AppTheme {
  AppTheme._();

  // Spacing
  static const double spacingXs = 4;
  static const double spacingSm = 8;
  static const double spacingMd = 12;
  static const double spacingLg = 16;
  static const double spacingXl = 24;

  // Border radius
  static const double radiusSm = 4;
  static const double radiusMd = 8;
  static const double radiusLg = 12;
  static const double radiusXl = 16;

  // Icon sizing
  static const double iconSm = 16;
  static const double iconMd = 20;
  static const double iconLg = 24;

  // Card padding
  static const EdgeInsets cardPadding = EdgeInsets.all(16);
  static const EdgeInsets cardPaddingSm = EdgeInsets.all(12);
  static const EdgeInsets cardPaddingMd = EdgeInsets.all(14);

  // Section padding
  static const EdgeInsets sectionPadding = EdgeInsets.only(top: 16, bottom: 8);

  static const Color statusGreen = Colors.green;
  static const Color statusOrange = Colors.orange;
  static const Color statusRed = Colors.red;
  static const Color statusBlue = Colors.blue;
  static const Color statusGrey = Colors.grey;
  static const Color statusPurple = Colors.purple;
  static const Color statusTeal = Colors.teal;

  static Color healthColor(int score) {
    if (score >= 60) return statusGreen;
    if (score >= 40) return statusOrange;
    return statusRed;
  }

  static Color severityColor(String severity) {
    switch (severity) {
      case 'critical': return statusRed;
      case 'warning': return statusOrange;
      case 'success': return statusGreen;
      default: return statusGrey;
    }
  }
}

final lightTheme = ThemeData(
  useMaterial3: true,
  colorSchemeSeed: _seedColor,
  brightness: Brightness.light,
  appBarTheme: const AppBarTheme(centerTitle: true, elevation: 0),
  cardTheme: const CardThemeData(clipBehavior: Clip.antiAlias),
  inputDecorationTheme: InputDecorationTheme(
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(AppTheme.radiusLg)),
    filled: true,
  ),
  elevatedButtonTheme: ElevatedButtonThemeData(
    style: ElevatedButton.styleFrom(
      minimumSize: const Size(double.infinity, 48),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppTheme.radiusLg)),
    ),
  ),
  filledButtonTheme: FilledButtonThemeData(
    style: FilledButton.styleFrom(
      minimumSize: const Size(double.infinity, 48),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppTheme.radiusLg)),
    ),
  ),
  chipTheme: ChipThemeData(
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppTheme.radiusMd)),
  ),
);

final darkTheme = ThemeData(
  useMaterial3: true,
  colorSchemeSeed: _seedColor,
  brightness: Brightness.dark,
  appBarTheme: const AppBarTheme(centerTitle: true, elevation: 0),
  cardTheme: const CardThemeData(clipBehavior: Clip.antiAlias),
  inputDecorationTheme: InputDecorationTheme(
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(AppTheme.radiusLg)),
    filled: true,
  ),
  elevatedButtonTheme: ElevatedButtonThemeData(
    style: ElevatedButton.styleFrom(
      minimumSize: const Size(double.infinity, 48),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppTheme.radiusLg)),
    ),
  ),
  filledButtonTheme: FilledButtonThemeData(
    style: FilledButton.styleFrom(
      minimumSize: const Size(double.infinity, 48),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppTheme.radiusLg)),
    ),
  ),
  chipTheme: ChipThemeData(
    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(AppTheme.radiusMd)),
  ),
);
