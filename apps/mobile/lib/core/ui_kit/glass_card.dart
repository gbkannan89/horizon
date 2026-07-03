import 'dart:ui';
import 'package:flutter/material.dart';
import '../theme/design_tokens.dart';

class GlassCard extends StatelessWidget {
  final Widget child;
  final double blur;
  final double opacity;
  final EdgeInsetsGeometry padding;
  final double borderRadius;
  final VoidCallback? onTap;

  const GlassCard({
    super.key,
    required this.child,
    this.blur = 10.0,
    this.opacity = 0.05, // very subtle fill
    this.padding = const EdgeInsets.all(AppSpacing.md),
    this.borderRadius = AppRadius.lg,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final isDark = Theme.of(context).brightness == Brightness.dark;
    final fillOpacity = isDark ? 0.2 : 0.4;
    final borderColor = isDark 
        ? Colors.white.withOpacity(0.1) 
        : Colors.white.withOpacity(0.5);
    final shadowColor = isDark 
        ? Colors.black.withOpacity(0.4) 
        : Colors.black.withOpacity(0.05);

    Widget content = Container(
      decoration: BoxDecoration(
        color: (isDark ? Colors.white : Colors.white).withOpacity(fillOpacity),
        borderRadius: BorderRadius.circular(borderRadius),
        border: Border.all(
          color: borderColor,
          width: 1.0,
        ),
        boxShadow: [
          BoxShadow(
            color: shadowColor,
            blurRadius: 20.0,
            offset: const Offset(0, 10),
          ),
        ],
      ),
      padding: padding,
      child: child,
    );

    if (onTap != null) {
      content = InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(borderRadius),
        child: content,
      );
    }

    return ClipRRect(
      borderRadius: BorderRadius.circular(borderRadius),
      child: BackdropFilter(
        filter: ImageFilter.blur(sigmaX: blur, sigmaY: blur),
        child: content,
      ),
    );
  }
}
