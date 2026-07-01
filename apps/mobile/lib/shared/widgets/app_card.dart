import 'package:flutter/material.dart';

class AppCard extends StatelessWidget {
  final Widget child;
  final EdgeInsetsGeometry? padding;
  final VoidCallback? onTap;
  final Color? color;

  const AppCard({super.key, required this.child, this.padding, this.onTap, this.color});

  @override
  Widget build(BuildContext context) {
    return onTap != null
        ? Card(clipBehavior: Clip.antiAlias, color: color, child: InkWell(onTap: onTap, child: Padding(padding: padding ?? const EdgeInsets.all(16), child: child)))
        : Card(color: color, child: Padding(padding: padding ?? const EdgeInsets.all(16), child: child));
  }
}
