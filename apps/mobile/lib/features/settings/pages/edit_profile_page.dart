import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../repository/settings_repository.dart';

class EditProfilePage extends ConsumerStatefulWidget {
  const EditProfilePage({super.key});
  @override
  ConsumerState<EditProfilePage> createState() => _EditProfilePageState();
}

class _EditProfilePageState extends ConsumerState<EditProfilePage> {
  final _formKey = GlobalKey<FormState>();
  final _nameCtrl = TextEditingController();
  final _phoneCtrl = TextEditingController();
  bool _isLoading = false;
  bool _isSaving = false;
  String? _error;
  String? _dateOfBirth;

  @override
  void initState() {
    super.initState();
    _loadProfile();
  }

  Future<void> _loadProfile() async {
    setState(() { _isLoading = true; _error = null; });
    try {
      final repo = ref.read(settingsRepositoryProvider);
      final profile = await repo.getProfile();
      _nameCtrl.text = profile.name ?? '';
      _phoneCtrl.text = profile.phone ?? '';
    } catch (e) {
      _error = e.toString();
    }
    if (mounted) setState(() => _isLoading = false);
  }

  @override
  void dispose() { _nameCtrl.dispose(); _phoneCtrl.dispose(); super.dispose(); }

  Future<void> _save() async {
    if (!_formKey.currentState!.validate()) return;
    setState(() { _isSaving = true; _error = null; });
    try {
      final repo = ref.read(settingsRepositoryProvider);
      await repo.updateProfile(data: {
        'name': _nameCtrl.text.trim(),
        'phone': _phoneCtrl.text.trim(),
      });
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Profile updated'), behavior: SnackBarBehavior.floating));
        context.pop();
      }
    } catch (e) {
      if (mounted) setState(() => _error = e.toString());
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    if (_isLoading) {
      return Scaffold(
        appBar: AppBar(title: const Text('Edit Profile')),
        body: const Center(child: CircularProgressIndicator()),
      );
    }

    return Scaffold(
      appBar: AppBar(title: const Text('Edit Profile'), actions: [
        TextButton(
          onPressed: _isSaving ? null : _save,
          child: _isSaving
              ? const SizedBox(width: 18, height: 18, child: CircularProgressIndicator(strokeWidth: 2))
              : const Text('Save'),
        ),
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
            TextFormField(
              decoration: InputDecoration(labelText: 'Date of Birth', prefixIcon: const Icon(Icons.calendar_today), hintText: _dateOfBirth ?? 'Tap to select'),
              readOnly: true,
              onTap: () async {
                final picked = await showDatePicker(context: context, initialDate: DateTime.now().subtract(const Duration(days: 365 * 30)), firstDate: DateTime(1950), lastDate: DateTime.now());
                if (picked != null) setState(() => _dateOfBirth = '${picked.year}-${picked.month.toString().padLeft(2, '0')}-${picked.day.toString().padLeft(2, '0')}');
              },
            ),
            if (_error != null)
              Padding(
                padding: const EdgeInsets.only(top: 16),
                child: Container(
                  padding: const EdgeInsets.all(12),
                  decoration: BoxDecoration(color: theme.colorScheme.errorContainer, borderRadius: BorderRadius.circular(12)),
                  child: Row(children: [
                    Icon(Icons.error_outline, size: 20, color: theme.colorScheme.error),
                    const SizedBox(width: 8),
                    Expanded(child: Text(_error!, style: TextStyle(color: theme.colorScheme.error, fontSize: 14))),
                  ]),
                ),
              ),
          ],
        ),
      ),
    );
  }
}
