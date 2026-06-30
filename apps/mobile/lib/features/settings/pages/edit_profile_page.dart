import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

class EditProfilePage extends ConsumerStatefulWidget {
  const EditProfilePage({super.key});
  @override
  ConsumerState<EditProfilePage> createState() => _EditProfilePageState();
}

class _EditProfilePageState extends ConsumerState<EditProfilePage> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();
  final _phoneCtrl = TextEditingController();

  @override
  void dispose() { _nameCtrl.dispose(); _phoneCtrl.dispose(); super.dispose(); }

  void _save() {
    if (!_formKey.currentState!.validate()) return;
    ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Profile updated'), behavior: SnackBarBehavior.floating));
    context.pop();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Edit Profile'), actions: [
        TextButton(onPressed: _save, child: const Text('Save')),
      ]),
      body: Form(
        key: _formKey,
        child: ListView(
          padding: const EdgeInsets.all(16),
          children: [
            Center(
              child: Stack(children: [
                CircleAvatar(radius: 48, backgroundColor: theme.colorScheme.primaryContainer,
                  child: Icon(Icons.person, size: 40, color: theme.colorScheme.primary)),
                Positioned(bottom: 0, right: 0, child: Container(
                  padding: const EdgeInsets.all(4),
                  decoration: BoxDecoration(color: theme.colorScheme.primary, shape: BoxShape.circle),
                  child: Icon(Icons.camera_alt, size: 16, color: theme.colorScheme.onPrimary),
                )),
              ]),
            ),
            const SizedBox(height: 24),
            TextFormField(controller: _nameCtrl, decoration: const InputDecoration(labelText: 'Name', prefixIcon: Icon(Icons.person_outline)),
              validator: (v) => v == null || v.isEmpty ? 'Name is required' : null),
            const SizedBox(height: 16),
            TextFormField(controller: _phoneCtrl, decoration: const InputDecoration(labelText: 'Phone', prefixIcon: Icon(Icons.phone_outlined)),
              keyboardType: TextInputType.phone),
            const SizedBox(height: 16),
            TextFormField(decoration: const InputDecoration(labelText: 'Date of Birth', prefixIcon: Icon(Icons.calendar_today)), readOnly: true),
          ],
        ),
      ),
    );
  }
}
