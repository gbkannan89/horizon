import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/auth_provider.dart';
import 'login_screen.dart';

class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});

  void _showDeleteConfirmationDialog(BuildContext context) {
    final passwordController = TextEditingController();
    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setDialogState) {
          String? errorText;
          bool isLoading = false;

          return AlertDialog(
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
            title: const Row(
              children: [
                Icon(Icons.warning_amber_rounded, color: Colors.red, size: 28),
                SizedBox(width: 12),
                Text('Delete Account', style: TextStyle(fontWeight: FontWeight.bold)),
              ],
            ),
            content: SizedBox(
              width: double.maxFinite,
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: Colors.red.withValues(alpha: 0.08),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: const Text(
                      'This will permanently delete your account and all associated data including income, expenses, goals, investments, insurance, bills, and household data. This action cannot be undone.',
                      style: TextStyle(color: Colors.red, fontSize: 13, fontWeight: FontWeight.w500),
                    ),
                  ),
                  const SizedBox(height: 20),
                  TextField(
                    controller: passwordController,
                    obscureText: true,
                    decoration: InputDecoration(
                      labelText: 'Enter your password to confirm',
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                      errorText: errorText,
                    ),
                  ),
                ],
              ),
            ),
            actions: [
              TextButton(
                onPressed: () => Navigator.of(ctx).pop(),
                child: const Text('Cancel'),
              ),
              ElevatedButton(
                onPressed: isLoading
                    ? null
                    : () async {
                        setDialogState(() {
                          errorText = null;
                          isLoading = true;
                        });
                        try {
                          await Provider.of<AuthProvider>(context, listen: false)
                              .deleteAccount(passwordController.text.trim());
                          if (ctx.mounted) Navigator.of(ctx).pop();
                          if (context.mounted) {
                            Navigator.of(context).pushAndRemoveUntil(
                              MaterialPageRoute(builder: (_) => const LoginScreen()),
                              (route) => false,
                            );
                          }
                        } on Exception catch (e) {
                          setDialogState(() {
                            errorText = e.toString().replaceFirst('Exception: ', '');
                            isLoading = false;
                          });
                        }
                      },
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.red,
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                ),
                child: isLoading
                    ? const SizedBox(width: 20, height: 20, child: CircularProgressIndicator(strokeWidth: 2, color: Colors.white))
                    : const Text('Delete My Account'),
              ),
            ],
          );
        },
      ),
    );
  }

  Widget _buildSettingsTile(String title, IconData icon, ColorScheme colorScheme, {VoidCallback? onTap, bool isDestructive = false}) {
    return ListTile(
      contentPadding: const EdgeInsets.symmetric(horizontal: 24, vertical: 8),
      leading: Container(
        padding: const EdgeInsets.all(10),
        decoration: BoxDecoration(
          color: isDestructive ? const Color(0xFFFEE2E2) : colorScheme.primary.withValues(alpha: 0.1),
          shape: BoxShape.circle,
        ),
        child: Icon(
          icon, 
          color: isDestructive ? const Color(0xFFEF4444) : colorScheme.primary,
          size: 20,
        ),
      ),
      title: Text(
        title, 
        style: TextStyle(
          fontWeight: FontWeight.w600, 
          color: isDestructive ? const Color(0xFFEF4444) : colorScheme.onSurface,
        ),
      ),
      trailing: isDestructive ? null : Icon(Icons.chevron_right_rounded, color: Colors.grey[400]),
      onTap: onTap,
    );
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      backgroundColor: colorScheme.surface,
      appBar: AppBar(
        title: const Text('Settings'),
      ),
      body: SingleChildScrollView(
        child: Column(
          children: [
            // Profile Header
            Container(
              margin: const EdgeInsets.all(24),
              padding: const EdgeInsets.all(24),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(32),
                boxShadow: [
                  BoxShadow(
                    color: colorScheme.primary.withValues(alpha: 0.08),
                    blurRadius: 25,
                    offset: const Offset(0, 10),
                  ),
                ],
              ),
              child: Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      border: Border.all(color: colorScheme.primary.withValues(alpha: 0.2), width: 3),
                    ),
                    child: CircleAvatar(
                      radius: 36,
                      backgroundColor: colorScheme.primary.withValues(alpha: 0.1),
                      child: Icon(Icons.person_rounded, size: 40, color: colorScheme.primary),
                    ),
                  ),
                  const SizedBox(width: 20),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Demo User',
                        style: TextStyle(fontSize: 22, fontWeight: FontWeight.w800, letterSpacing: -0.5),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        'user@horizon.com',
                        style: TextStyle(color: Colors.grey[600], fontWeight: FontWeight.w500),
                      ),
                      const SizedBox(height: 8),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                        decoration: BoxDecoration(
                          color: const Color(0xFF81E6D9).withValues(alpha: 0.2),
                          borderRadius: BorderRadius.circular(10),
                        ),
                        child: const Text('Premium Member', style: TextStyle(color: Color(0xFF319795), fontSize: 11, fontWeight: FontWeight.bold)),
                      )
                    ],
                  ),
                ],
              ),
            ),
            
            const SizedBox(height: 8),
            
            // Settings List
            Container(
              decoration: const BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.only(topLeft: Radius.circular(32), topRight: Radius.circular(32)),
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Padding(
                    padding: EdgeInsets.only(left: 32, top: 32, bottom: 8),
                    child: Text('Account', style: TextStyle(fontWeight: FontWeight.w800, color: Colors.grey, fontSize: 13)),
                  ),
                  _buildSettingsTile('Personal Information', Icons.person_outline_rounded, colorScheme, onTap: () {}),
                  _buildSettingsTile('Bank Connections', Icons.account_balance_rounded, colorScheme, onTap: () {}),
                  _buildSettingsTile('Security', Icons.shield_outlined, colorScheme, onTap: () {}),
                  
                  const Padding(
                    padding: EdgeInsets.only(left: 32, top: 24, bottom: 8),
                    child: Text('Preferences', style: TextStyle(fontWeight: FontWeight.w800, color: Colors.grey, fontSize: 13)),
                  ),
                  _buildSettingsTile('Notifications', Icons.notifications_none_rounded, colorScheme, onTap: () {}),
                  _buildSettingsTile('Appearance', Icons.color_lens_outlined, colorScheme, onTap: () {}),
                  
                  const Padding(
                    padding: EdgeInsets.only(left: 32, top: 24, bottom: 8),
                    child: Text('More', style: TextStyle(fontWeight: FontWeight.w800, color: Colors.grey, fontSize: 13)),
                  ),
                  _buildSettingsTile('Help & Support', Icons.help_outline_rounded, colorScheme, onTap: () {}),
                  _buildSettingsTile(
                    'Log Out', 
                    Icons.logout_rounded, 
                    colorScheme, 
                    isDestructive: true,
                    onTap: () async {
                      await Provider.of<AuthProvider>(context, listen: false).logout();
                      if (context.mounted) {
                        Navigator.of(context).pushReplacement(
                          MaterialPageRoute(builder: (context) => const LoginScreen()),
                        );
                      }
                    },
                  ),
                  const SizedBox(height: 4),
                  _buildSettingsTile(
                    'Delete Account',
                    Icons.delete_forever_rounded,
                    colorScheme,
                    isDestructive: true,
                    onTap: () => _showDeleteConfirmationDialog(context),
                  ),
                  const SizedBox(height: 40),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
