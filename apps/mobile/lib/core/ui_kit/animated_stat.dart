import 'package:flutter/material.dart';

class AnimatedStatValue extends StatelessWidget {
  final double value;
  final String format;
  final TextStyle? style;
  final Duration duration;

  const AnimatedStatValue({
    super.key,
    required this.value,
    this.format = '#,##0.00',
    this.style,
    this.duration = const Duration(milliseconds: 500),
  });

  @override
  Widget build(BuildContext context) {
    return TweenAnimationBuilder<double>(
      tween: Tween<double>(begin: 0, end: value),
      duration: duration,
      curve: Curves.easeOutCubic,
      builder: (context, currentValue, child) {
        // Simple formatting for demonstration, assuming a currency format or similar
        // For production, use intl package NumberFormat
        final formattedValue = currentValue.toStringAsFixed(2).replaceAllMapped(
          RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
          (Match m) => '${m[1]},'
        );
        return Text(
          formattedValue,
          style: style,
        );
      },
    );
  }
}
