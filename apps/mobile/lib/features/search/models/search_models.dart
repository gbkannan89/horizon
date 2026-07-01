import 'package:flutter/material.dart';

class SearchResultItem {
  final String id;
  final String title;
  final String subtitle;
  final String module;
  final String? amount;
  final String? date;
  final IconData icon;
  final Color color;
  final String route;
  final Map<String, dynamic>? routeParams;

  SearchResultItem({
    required this.id,
    required this.title,
    required this.subtitle,
    required this.module,
    this.amount,
    this.date,
    required this.icon,
    this.color = Colors.grey,
    required this.route,
    this.routeParams,
  });
}

class SearchGroup {
  final String module;
  final IconData icon;
  final Color color;
  final List<SearchResultItem> items;
  final bool isAvailable;

  SearchGroup({
    required this.module,
    required this.icon,
    this.color = Colors.grey,
    this.items = const [],
    this.isAvailable = true,
  });

  bool get isEmpty => items.isEmpty;
}
