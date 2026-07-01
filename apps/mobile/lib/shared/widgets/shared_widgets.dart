import 'dart:async';
import 'package:flutter/material.dart';
import '../../app/theme.dart';

class SharedLoadingView extends StatelessWidget {
  final String? message;
  const SharedLoadingView({super.key, this.message});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const CircularProgressIndicator(),
          if (message != null) ...[
            const SizedBox(height: AppSpacing.md),
            Text(message!, style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ],
        ],
      ),
    );
  }
}

class SharedSkeletonList extends StatelessWidget {
  final int itemCount;
  final double itemHeight;
  const SharedSkeletonList({super.key, this.itemCount = 8, this.itemHeight = 72});

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      padding: const EdgeInsets.all(AppSpacing.md),
      itemCount: itemCount,
      itemBuilder: (_, __) => Padding(
        padding: const EdgeInsets.only(bottom: AppSpacing.sm),
        child: Card(
          child: SizedBox(
            height: itemHeight,
            child: Center(
              child: CircularProgressIndicator(
                strokeWidth: 2,
                color: Colors.grey.withValues(alpha: 0.3),
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class SharedErrorView extends StatelessWidget {
  final String? message;
  final VoidCallback? onRetry;
  final VoidCallback? onDismiss;

  const SharedErrorView({super.key, this.message, this.onRetry, this.onDismiss});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.lg),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.cloud_off, size: AppIconSize.xl * 2, color: theme.colorScheme.error),
            const SizedBox(height: AppSpacing.md),
            Text('Could not load data', style: theme.textTheme.titleMedium),
            if (message != null) ...[
              const SizedBox(height: AppSpacing.sm),
              Text(message!, textAlign: TextAlign.center, style: theme.textTheme.bodySmall),
            ],
            const SizedBox(height: AppSpacing.lg),
            if (onRetry != null)
              FilledButton.icon(
                onPressed: onRetry,
                icon: const Icon(Icons.refresh),
                label: const Text('Try Again'),
              ),
            if (onDismiss != null) ...[
              const SizedBox(height: AppSpacing.sm),
              TextButton(onPressed: onDismiss, child: const Text('Dismiss')),
            ],
          ],
        ),
      ),
    );
  }
}

class SharedEmptyView extends StatelessWidget {
  final IconData icon;
  final String title;
  final String? subtitle;
  final String? actionLabel;
  final VoidCallback? onAction;

  const SharedEmptyView({
    super.key,
    required this.icon,
    required this.title,
    this.subtitle,
    this.actionLabel,
    this.onAction,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(AppSpacing.xl),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, size: AppIconSize.xl * 2, color: theme.colorScheme.onSurfaceVariant.withValues(alpha: 0.4)),
            const SizedBox(height: AppSpacing.md),
            Text(title, style: theme.textTheme.titleMedium, textAlign: TextAlign.center),
            if (subtitle != null) ...[
              const SizedBox(height: AppSpacing.sm),
              Text(subtitle!, textAlign: TextAlign.center, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            ],
            if (actionLabel != null && onAction != null) ...[
              const SizedBox(height: AppSpacing.md),
              TextButton.icon(onPressed: onAction, icon: const Icon(Icons.clear_all), label: Text(actionLabel!)),
            ],
          ],
        ),
      ),
    );
  }
}

class SharedOfflineView extends StatelessWidget {
  final String? message;
  const SharedOfflineView({super.key, this.message});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(Icons.wifi_off, size: AppIconSize.xl * 2, color: theme.colorScheme.onSurfaceVariant),
          const SizedBox(height: AppSpacing.md),
          Text("You're offline", style: theme.textTheme.titleMedium),
          if (message != null) ...[
            const SizedBox(height: AppSpacing.sm),
            Text(message!, textAlign: TextAlign.center, style: theme.textTheme.bodySmall),
          ],
        ],
      ),
    );
  }
}

class SharedSectionHeader extends StatelessWidget {
  final String title;
  final VoidCallback? onSeeAll;
  final Widget? trailing;

  const SharedSectionHeader({super.key, required this.title, this.onSeeAll, this.trailing});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Padding(
      padding: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(title, style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600, color: theme.colorScheme.primary)),
          if (trailing != null) trailing!,
          if (onSeeAll != null)
            TextButton(onPressed: onSeeAll, child: const Text('See all')),
        ],
      ),
    );
  }
}

class SharedCard extends StatelessWidget {
  final Widget child;
  final VoidCallback? onTap;
  final EdgeInsetsGeometry? padding;

  const SharedCard({super.key, required this.child, this.onTap, this.padding});

  @override
  Widget build(BuildContext context) {
    return Card(
      clipBehavior: Clip.antiAlias,
      child: onTap != null
          ? InkWell(onTap: onTap, child: Padding(padding: padding ?? const EdgeInsets.all(AppSpacing.md), child: child))
          : Padding(padding: padding ?? const EdgeInsets.all(AppSpacing.md), child: child),
    );
  }
}

class SharedStatusBadge extends StatelessWidget {
  final String label;
  final Color color;
  final IconData? icon;

  const SharedStatusBadge({super.key, required this.label, required this.color, this.icon});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: AppSpacing.sm, vertical: AppSpacing.xs),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(AppRadius.sm),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          if (icon != null) ...[Icon(icon, size: AppIconSize.sm, color: color), const SizedBox(width: AppSpacing.xs)],
          Text(label, style: theme.textTheme.labelSmall?.copyWith(color: color, fontWeight: FontWeight.w500)),
        ],
      ),
    );
  }
}

class SharedSearchBar extends StatefulWidget {
  final ValueChanged<String> onChanged;
  final VoidCallback? onCancel;
  final TextEditingController? controller;
  final String hintText;

  const SharedSearchBar({
    super.key,
    required this.onChanged,
    this.onCancel,
    this.controller,
    this.hintText = 'Search...',
  });

  @override
  State<SharedSearchBar> createState() => _SharedSearchBarState();
}

class _SharedSearchBarState extends State<SharedSearchBar> {
  late TextEditingController _controller;
  Timer? _debounce;

  @override
  void initState() {
    super.initState();
    _controller = widget.controller ?? TextEditingController();
  }

  @override
  void dispose() {
    if (widget.controller == null) _controller.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onChanged(String value) {
    _debounce?.cancel();
    _debounce = Timer(const Duration(milliseconds: 300), () => widget.onChanged(value));
  }

  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: _controller,
      autofocus: true,
      decoration: InputDecoration(
        hintText: widget.hintText,
        border: InputBorder.none,
        suffixIcon: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () {
            _controller.clear();
            widget.onCancel?.call();
          },
        ),
      ),
      onChanged: _onChanged,
    );
  }
}

class SharedFilterChipWidget extends StatelessWidget {
  final String label;
  final bool selected;
  final VoidCallback onTap;

  const SharedFilterChipWidget({super.key, required this.label, required this.selected, required this.onTap});

  @override
  Widget build(BuildContext context) {
    return FilterChip(
      label: Text(label),
      selected: selected,
      onSelected: (_) => onTap(),
      visualDensity: VisualDensity.compact,
    );
  }
}
