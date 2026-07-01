import 'package:flutter/material.dart';
import 'package:horizon_mobile/app/theme.dart';

class DataManagementPage extends StatelessWidget {
  const DataManagementPage({super.key});
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Data Management')),
      body: ListView(
        padding: const EdgeInsets.all(AppSpacing.md),
        children: [
          _card(theme, Icons.download, 'Export My Data', 'Download your financial data', 'Coming soon'),
          _card(theme, Icons.upload, 'Import Data', 'Import transactions from other platforms', 'Coming soon'),
          _card(theme, Icons.backup, 'Backup', 'Create a backup of your financial data', 'Coming soon'),
          _card(theme, Icons.restore, 'Restore', 'Restore from a previous backup', 'Coming soon'),
        ],
      ),
    );
  }

  Widget _card(ThemeData t, IconData icon, String title, String subtitle, String status) {
    return Card(
      margin: const EdgeInsets.only(bottom: AppSpacing.sm),
      child: Opacity(
        opacity: 0.5,
        child: ListTile(
          leading: Icon(icon),
          title: Text(title),
          subtitle: Text(subtitle),
          trailing: Chip(label: Text(status, style: t.textTheme.labelSmall), visualDensity: VisualDensity.compact),
        ),
      ),
    );
  }
}
