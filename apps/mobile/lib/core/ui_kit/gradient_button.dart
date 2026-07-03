import 'package:flutter/material.dart';
import '../theme/design_tokens.dart';

class GradientButton extends StatelessWidget {
  final VoidCallback onPressed;
  final Widget child;
  final Gradient? gradient;
  final double borderRadius;
  final EdgeInsetsGeometry padding;
  final bool isLoading;

  const GradientButton({
    super.key,
    required this.onPressed,
    required this.child,
    this.gradient,
    this.borderRadius = AppRadius.md,
    this.padding = const EdgeInsets.symmetric(vertical: AppSpacing.md, horizontal: AppSpacing.xl),
    this.isLoading = false,
  });

  @override
  Widget build(BuildContext context) {
    final defaultGradient = AppColors.primaryGradient;

    return Container(
      decoration: BoxDecoration(
        gradient: gradient ?? defaultGradient,
        borderRadius: BorderRadius.circular(borderRadius),
        boxShadow: AppShadows.soft,
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: isLoading ? null : onPressed,
          borderRadius: BorderRadius.circular(borderRadius),
          child: Padding(
            padding: padding,
            child: Center(
              child: isLoading 
                  ? const SizedBox(
                      width: 20, 
                      height: 20, 
                      child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2)
                    )
                  : DefaultTextStyle(
                      style: const TextStyle(
                        color: Colors.white, 
                        fontWeight: FontWeight.w600,
                        fontSize: 16,
                      ),
                      child: child,
                    ),
            ),
          ),
        ),
      ),
    );
  }
}
