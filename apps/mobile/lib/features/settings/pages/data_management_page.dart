import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/app/theme.dart';
import 'package:horizon_mobile/core/network/api_client.dart';

class DataManagementPage extends ConsumerWidget {
  const DataManagementPage({super.key});

  Future<void> _export(WidgetRef ref, BuildContext context) async {
    try {
      final api = ref.read(apiClientProvider);
      final resp = await api.get('/data/export');
      final data = resp.data is Map ? (resp.data as Map)['data'] : resp.data;
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Data exported successfully'), behavior: SnackBarBehavior.floating));
      }
    } catch (e) {
      if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Export failed: $e'), behavior: SnackBarBehavior.floating));
    }
  }

  Future<void> _backup(WidgetRef ref, BuildContext context) async {
    try {
      final api = ref.read(apiClientProvider);
      await api.post('/data/backup');
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Backup created'), behavior: SnackBarBehavior.floating));
      }
    } catch (e) {
      if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Backup failed: $e'), behavior: SnackBarBehavior.floating));
    }
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Data Management')),
      body: ListView(
        padding: const EdgeInsets.all(AppSpacing.md),
        children: [
          _card(theme, Icons.download, 'Export My Data', 'Download your financial data', 'Available', () => _export(ref, context)),
          _card(theme, Icons.upload, 'Import Data', 'Import transactions from other platforms', 'MVP', null),
          _card(theme, Icons.backup, 'Backup', 'Create a backup of your data', 'Available', () => _backup(ref, context)),
          _card(theme, Icons.restore, 'Restore', 'Restore from a previous backup', 'MVP', null),
        ],
      ),
    );
  }

  Widget _card(ThemeData t, IconData icon, String title, String subtitle, String status, VoidCallback? onTap) {
    return Card(
      margin: const EdgeInsets.only(bottom: AppSpacing.sm),
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: onTap != null ? () => onTap() : null,
        child: Opacity(
          opacity: onTap != null ? 1.0 : 0.5,
          child: ListTile(
            leading: Icon(icon),
            title: Text(title),
            subtitle: Text(subtitle),
            trailing: Chip(label: Text(status, style: t.textTheme.labelSmall), visualDensity: VisualDensity.compact),
          ),
        ),
      ),
    );
  }
}
