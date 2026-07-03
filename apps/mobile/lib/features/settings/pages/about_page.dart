import 'package:flutter/material.dart';

class AboutPage extends StatelessWidget {
  const AboutPage({super.key});

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('About')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Center(
            child: Column(children: [
              Container(
                width: 80, height: 80,
                decoration: BoxDecoration(color: theme.colorScheme.primaryContainer, borderRadius: BorderRadius.circular(20)),
                child: Icon(Icons.trending_up, size: 40, color: theme.colorScheme.primary),
              ),
              const SizedBox(height: 16),
              Text('Horizon', style: theme.textTheme.headlineSmall?.copyWith(fontWeight: FontWeight.bold)),
              const SizedBox(height: 4),
              Text('Goal-Centric Financial OS', style: theme.textTheme.bodyMedium?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              const SizedBox(height: 8),
              Text('Version 0.1.0-dev', style: theme.textTheme.bodySmall),
            ]),
          ),
          const SizedBox(height: 32),
          Card(child: Column(children: [
            ListTile(title: const Text('Privacy Policy'), trailing: const Icon(Icons.open_in_new, size: 16), onTap: () {}),
            ListTile(title: const Text('Terms of Service'), trailing: const Icon(Icons.open_in_new, size: 16), onTap: () {}),
            ListTile(title: const Text('Open Source Licenses'), trailing: const Icon(Icons.chevron_right, size: 18), onTap: () => showLicensePage(context: context)),
          ])),
          const SizedBox(height: 16),
          const Card(child: Column(children: [
            ListTile(title: Text('API Version'), subtitle: Text('v1'), trailing: Icon(Icons.chevron_right, size: 18)),
            ListTile(title: Text('Build Number'), subtitle: Text('1'), trailing: Icon(Icons.chevron_right, size: 18)),
          ])),
          const SizedBox(height: 24),
          Center(child: Text('© 2026 Horizon Financial Systems', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant))),
        ],
      ),
    );
  }
}
