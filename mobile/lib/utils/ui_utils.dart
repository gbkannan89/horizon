import 'package:flutter/material.dart';

class UiUtils {
  static void showSnack(BuildContext context, String message, {bool isError = false}) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      content: Row(
        children: [
          Icon(isError ? Icons.error_outline : Icons.check_circle_outline, color: Colors.white),
          const SizedBox(width: 12),
          Expanded(child: Text(message)),
        ],
      ),
      backgroundColor: isError ? Colors.redAccent : const Color(0xFF059669),
      behavior: SnackBarBehavior.floating,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      margin: const EdgeInsets.all(16),
    ));
  }

  static void showDeleteBottomSheet(BuildContext context, String title, VoidCallback onConfirm) {
    showModalBottomSheet(
      context: context,
      backgroundColor: Colors.transparent,
      builder: (ctx) => TweenAnimationBuilder<double>(
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeOutBack,
        tween: Tween(begin: 0.0, end: 1.0),
        builder: (ctx, value, child) => Transform.scale(
          scale: value,
          alignment: Alignment.bottomCenter,
          child: Container(
            margin: const EdgeInsets.all(16),
            padding: const EdgeInsets.all(24),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(24),
              boxShadow: const [BoxShadow(color: Colors.black12, blurRadius: 20)],
            ),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Icon(Icons.warning_amber_rounded, color: Colors.redAccent, size: 48),
                const SizedBox(height: 16),
                Text('Delete $title?', style: const TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Colors.redAccent)),
                const SizedBox(height: 8),
                const Text('This action cannot be undone.', style: TextStyle(color: Colors.grey)),
                const SizedBox(height: 24),
                Row(
                  children: [
                    Expanded(
                      child: OutlinedButton(
                        onPressed: () => Navigator.pop(ctx),
                        style: OutlinedButton.styleFrom(
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                          side: BorderSide(color: Colors.grey.shade300),
                        ),
                        child: const Text('Cancel', style: TextStyle(color: Colors.black87, fontWeight: FontWeight.bold)),
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: ElevatedButton(
                        onPressed: () {
                          Navigator.pop(ctx);
                          onConfirm();
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: Colors.redAccent,
                          padding: const EdgeInsets.symmetric(vertical: 16),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                        ),
                        child: const Text('Remove', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  static Widget buildEmptyState(String title, String desc, IconData icon, Color color) {
    String? assetPath;
    final lowerTitle = title.toLowerCase();
    if (lowerTitle.contains('gold')) {
      assetPath = 'assets/leaf.png';
    } else if (lowerTitle.contains('stock') || lowerTitle.contains('portfolio')) {
      assetPath = 'assets/camel.png';
    } else if (lowerTitle.contains('lending') || lowerTitle.contains('liability') || lowerTitle.contains('bill') || lowerTitle.contains('debt')) {
      assetPath = 'assets/flower-pot-1.png';
    } else if (lowerTitle.contains('asset') || lowerTitle.contains('device') || lowerTitle.contains('electronic') || lowerTitle.contains('insurance')) {
      assetPath = 'assets/flower-pot.png';
    }

    return Container(
      margin: const EdgeInsets.symmetric(vertical: 12),
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.05),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: color.withValues(alpha: 0.1)),
      ),
      child: Row(children: [
        assetPath != null
            ? Container(
                width: 44,
                height: 44,
                decoration: BoxDecoration(
                  color: Colors.white,
                  borderRadius: BorderRadius.circular(12),
                  boxShadow: [
                    BoxShadow(
                      color: Colors.black.withValues(alpha: 0.03),
                      blurRadius: 10,
                      offset: const Offset(0, 4),
                    ),
                  ],
                ),
                padding: const EdgeInsets.all(6),
                child: Image.asset(assetPath, fit: BoxFit.contain),
              )
            : Icon(icon, size: 36, color: color.withValues(alpha: 0.6)),
        const SizedBox(width: 16),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(title, style: TextStyle(fontWeight: FontWeight.bold, color: color.withValues(alpha: 0.8))),
              const SizedBox(height: 4),
              Text(desc, style: const TextStyle(fontSize: 12, color: Colors.grey)),
            ],
          ),
        )
      ]),
    );
  }
}
