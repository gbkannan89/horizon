import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/features/household/state/household_state.dart';

class HouseholdInvitePage extends ConsumerStatefulWidget {
  final String householdId;

  const HouseholdInvitePage({super.key, required this.householdId});

  @override
  ConsumerState<HouseholdInvitePage> createState() => _HouseholdInvitePageState();
}

class _HouseholdInvitePageState extends ConsumerState<HouseholdInvitePage> {
  final _formKey = GlobalKey<FormState>();
  final _userIdController = TextEditingController();
  String _selectedRole = 'Member';
  bool _isLoading = false;

  final List<String> _roles = ['Admin', 'Member', 'Viewer'];

  @override
  void dispose() {
    _userIdController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;
    
    setState(() { _isLoading = true; });
    try {
      final repo = ref.read(householdRepositoryProvider);
      await repo.inviteMember(widget.householdId, _userIdController.text.trim(), _selectedRole);
      
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(const SnackBar(content: Text('Invitation sent successfully')));
        ref.invalidate(householdDetailProvider(widget.householdId));
        context.pop();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    } finally {
      if (mounted) {
        setState(() { _isLoading = false; });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Invite Member')),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              TextFormField(
                controller: _userIdController,
                decoration: const InputDecoration(
                  labelText: 'User ID or Email',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.person),
                ),
                validator: (value) => value == null || value.isEmpty ? 'Please enter a user identifier' : null,
              ),
              const SizedBox(height: 16),
              DropdownButtonFormField<String>(
                value: _selectedRole,
                decoration: const InputDecoration(
                  labelText: 'Role',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.security),
                ),
                items: _roles.map((r) => DropdownMenuItem(value: r, child: Text(r))).toList(),
                onChanged: (val) {
                  if (val != null) setState(() { _selectedRole = val; });
                },
              ),
              const SizedBox(height: 32),
              ElevatedButton(
                onPressed: _isLoading ? null : _submit,
                style: ElevatedButton.styleFrom(padding: const EdgeInsets.symmetric(vertical: 16)),
                child: _isLoading 
                  ? const CircularProgressIndicator() 
                  : const Text('Send Invitation', style: TextStyle(fontSize: 16)),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
