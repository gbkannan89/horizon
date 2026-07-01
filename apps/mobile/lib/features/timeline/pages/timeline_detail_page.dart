import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/timeline_models.dart';
import '../widgets/timeline_widgets.dart';
import '../repository/timeline_repository.dart';

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
    return Scaffold(
      appBar: AppBar(title: const Text('Event Details')),
      body: _loading
          ? const LoadingView()
          : _item == null
              ? EmptyState(icon: Icons.search_off, title: 'Event not found')
              : ListView(
                  padding: const EdgeInsets.all(AppTheme.spacingLg),
                  children: [
                    TimelineDetailCard(item: _item!),
                    const SizedBox(height: AppTheme.spacingLg),
                    OutlinedButton.icon(
                      onPressed: () => context.pop(),
                      icon: const Icon(Icons.arrow_back),
                      label: const Text('Back to Timeline'),
                    ),
                  ],
                ),
    );
  }
}
