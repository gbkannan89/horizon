import 'package:flutter/material.dart';

class StatusChip extends StatelessWidget {
  final String label;
  final Color color;
  final double fontSize;

  const StatusChip({
    super.key,
    required this.label,
    required this.color,
    this.fontSize = 11,
  });

  factory StatusChip.fromStatus(String status, {Map<String, Color>? colorMap, double fontSize = 11}) {
    final colors = colorMap ?? {
      'Active': Colors.green, 'Completed': Colors.blue, 'Paused': Colors.orange,
      'AtRisk': Colors.red, 'healthy': Colors.green, 'good': Colors.green,
      'warning': Colors.orange, 'critical': Colors.red,
    };
    return StatusChip(label: status, color: colors[status] ?? Colors.grey, fontSize: fontSize);
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(color: color.withOpacity(0.15), borderRadius: BorderRadius.circular(8)),
      child: Text(label, style: TextStyle(color: color, fontSize: fontSize, fontWeight: FontWeight.w600)),
    );
  }
}
