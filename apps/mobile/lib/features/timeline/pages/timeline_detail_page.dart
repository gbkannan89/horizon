import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
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
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: const Text('Event Details')),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : _item == null
              ? Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(Icons.search_off, size: 64, color: theme.colorScheme.error),
                      const SizedBox(height: 16),
                      Text('Event not found', style: theme.textTheme.titleMedium),
                    ],
                  ),
                )
              : ListView(
                  padding: const EdgeInsets.all(16),
                  children: [
                    TimelineDetailCard(item: _item!),
                    const SizedBox(height: 16),
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
