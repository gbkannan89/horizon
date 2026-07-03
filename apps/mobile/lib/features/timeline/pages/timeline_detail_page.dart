import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:flutter_animate/flutter_animate.dart';
import '../models/timeline_models.dart';
import '../widgets/timeline_widgets.dart';
import '../repository/timeline_repository.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';

class TimelineDetailPage extends ConsumerStatefulWidget {
  final String timelineId;
  const TimelineDetailPage({super.key, required this.timelineId});

  @override
  ConsumerState<TimelineDetailPage> createState() => _TimelineDetailPageState();
}

class _TimelineDetailPageState extends ConsumerState<TimelineDetailPage> {
  TimelineItem? _item;
  bool _loading = true;

  @override
  void initState() {
    super.initState();
    _load();
  }

  Future<void> _load() async {
    final arg = ModalRoute.of(context)?.settings.arguments;
    if (arg is TimelineItem) {
      setState(() { _item = arg; _loading = false; });
      return;
    }
    try {
      final item = await ref.read(timelineRepositoryProvider).getById(id: widget.timelineId);
      if (mounted) setState(() { _item = item; _loading = false; });
    } catch (_) {
      if (mounted) setState(() => _loading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        title: Text('Event Details', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: Container(
        decoration: BoxDecoration(
          gradient: AppColors.backgroundGradient(isDark),
        ),
        child: _buildBody(theme),
      ),
    );
  }

  Widget _buildBody(ThemeData theme) {
    if (_loading) {
      return Center(
        child: CircularProgressIndicator(color: AppColors.teal500).animate().fade(),
      );
    }
    if (_item == null) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.search_off_rounded, size: 64, color: AppColors.red500),
            const SizedBox(height: 16),
            Text('Event not found', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
          ],
        ),
      );
    }
    
    final children = <Widget>[
      TimelineDetailCard(item: _item!),
      const SizedBox(height: 24),
      OutlinedButton.icon(
        onPressed: () => context.pop(),
        icon: const Icon(Icons.arrow_back_rounded),
        label: const Text('Back to Timeline'),
        style: OutlinedButton.styleFrom(
          padding: const EdgeInsets.symmetric(vertical: AppSpacing.md),
        ),
      ),
      const SizedBox(height: 100),
    ];
    
    return ListView.builder(
      padding: const EdgeInsets.all(AppSpacing.md),
      itemCount: children.length,
      itemBuilder: (context, index) {
        return children[index]
          .animate()
          .fade(duration: 400.ms, delay: (20 * index).ms)
          .slideY(begin: 0.1, end: 0, duration: 400.ms, curve: Curves.easeOutQuad, delay: (20 * index).ms);
      },
    );
  }
}
