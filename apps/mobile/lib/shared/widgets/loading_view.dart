import 'package:flutter/material.dart';

class LoadingView extends StatelessWidget {
  final double? size;
  final double strokeWidth;

  const LoadingView({super.key, this.size, this.strokeWidth = 3});

  @override
  Widget build(BuildContext context) {
    return const Center(child: CircularProgressIndicator());
  }
}

class SkeletonList extends StatelessWidget {
  final int itemCount;
  final double itemHeight;

  const SkeletonList({super.key, this.itemCount = 6, this.itemHeight = 100});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: itemCount,
      itemBuilder: (_, __) => Padding(
        padding: const EdgeInsets.only(bottom: 12),
        child: Card(
          child: SizedBox(
            height: itemHeight,
            child: Center(
              child: CircularProgressIndicator(
                strokeWidth: 2,
                color: theme.colorScheme.primary.withOpacity(0.3),
              ),
            ),
          ),
        ),
      ),
    );
  }
}
