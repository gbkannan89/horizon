import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/shared/providers/auth_state.dart';

class SecurityPage extends ConsumerStatefulWidget {
  const SecurityPage({super.key});
  @override
  ConsumerState<SecurityPage> createState() => _SecurityPageState();
}

class _SecurityPageState extends ConsumerState<SecurityPage> {
  final _currentPwdCtrl = TextEditingController();
  final _newPwdCtrl = TextEditingController();
  final _confirmPwdCtrl = TextEditingController();
  bool _showForm = false;

  @override
  void dispose() {
    _currentPwdCtrl.dispose(); _newPwdCtrl.dispose(); _confirmPwdCtrl.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Security')),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          Card(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Padding(
                  padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
                  child: Text('Password', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                ),
                if (!_showForm)
                  ListTile(
                    leading: const Icon(Icons.lock_outline),
                    title: const Text('Change Password'),
                    trailing: const Icon(Icons.chevron_right, size: 18),
                    onTap: () => setState(() => _showForm = true),
                  )
                else ...[
                  Padding(
                    padding: const EdgeInsets.fromLTRB(16, 0, 16, 12),
                    child: Column(children: [
                      TextField(controller: _currentPwdCtrl, obscureText: true, decoration: const InputDecoration(labelText: 'Current Password', prefixIcon: Icon(Icons.lock_outlined))),
                      const SizedBox(height: 12),
                      TextField(controller: _newPwdCtrl, obscureText: true, decoration: const InputDecoration(labelText: 'New Password', prefixIcon: Icon(Icons.lock))),
                      const SizedBox(height: 12),
                      TextField(controller: _confirmPwdCtrl, obscureText: true, decoration: const InputDecoration(labelText: 'Confirm Password', prefixIcon: Icon(Icons.lock))),
                      const SizedBox(height: 16),
                      Row(children: [
                        Expanded(child: OutlinedButton(onPressed: () => setState(() => _showForm = false), child: const Text('Cancel'))),
                        const SizedBox(width: 12),
                        Expanded(child: FilledButton(onPressed: () {
                          ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Password changed')));
                          setState(() => _showForm = false);
                        }, child: const Text('Update'))),
                      ]),
                    ]),
                  ),
                ],
              ],
            ),
          ),
          const SizedBox(height: 16),
          Card(
            child: Column(children: [
              ListTile(
                leading: const Icon(Icons.devices),
                title: const Text('Active Sessions'),
                subtitle: const Text('Manage your logged-in devices'),
                trailing: const Icon(Icons.chevron_right, size: 18),
              ),
              ListTile(
                leading: Icon(Icons.logout, color: theme.colorScheme.error),
                title: Text('Logout All Devices', style: TextStyle(color: theme.colorScheme.error)),
                onTap: () {
                  ref.read(authStateProvider.notifier).unauthenticated();
                  context.go('/login');
                },
              ),
            ]),
          ),
        ],
      ),
    );
  }
}
