import 'package:flutter/material.dart';

const _seedColor = Color(0xFF00897B);

final lightTheme = ThemeData(
  useMaterial3: true,
  colorSchemeSeed: _seedColor,
  brightness: Brightness.light,
  appBarTheme: const AppBarTheme(centerTitle: true, elevation: 0),
  cardTheme: CardThemeData(clipBehavior: Clip.antiAlias),
  inputDecorationTheme: InputDecorationTheme(
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
    filled: true,
  ),
);

final darkTheme = ThemeData(
  useMaterial3: true,
  colorSchemeSeed: _seedColor,
  brightness: Brightness.dark,
  appBarTheme: const AppBarTheme(centerTitle: true, elevation: 0),
  cardTheme: CardThemeData(clipBehavior: Clip.antiAlias),
  inputDecorationTheme: InputDecorationTheme(
    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
    filled: true,
  ),
);
