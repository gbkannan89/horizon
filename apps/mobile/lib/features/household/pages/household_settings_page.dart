import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/features/household/state/household_state.dart';

class HouseholdSettingsPage extends ConsumerWidget {
  final String householdId;

  const HouseholdSettingsPage({super.key, required this.householdId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detailAsync = ref.watch(householdDetailProvider(householdId));

    return Scaffold(
      appBar: AppBar(title: const Text('Household Settings')),
      body: detailAsync.when(
        data: (detail) {
          return ListView(
            padding: const EdgeInsets.all(16),
            children: [
              const Text('Household', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              Card(
                child: Column(
                  children: [
                    ListTile(title: const Text('Name'), subtitle: Text(detail.name)),
                    const Divider(height: 1),
                    ListTile(title: const Text('Type'), subtitle: Text(detail.householdType)),
                    const Divider(height: 1),
                    ListTile(title: const Text('Currency'), subtitle: Text(detail.currency)),
                    const Divider(height: 1),
                    ListTile(title: const Text('Country'), subtitle: Text(detail.country)),
                  ],
                ),
              ),
              const SizedBox(height: 24),
              const Text('Members & Roles', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
              const SizedBox(height: 8),
              ...detail.members.map((member) => Card(
                child: _MemberRoleTile(householdId: householdId, member: member),
              )),
              const SizedBox(height: 24),
              ElevatedButton.icon(
                onPressed: () async {
                  final confirmed = await showDialog<bool>(
                    context: context,
                    builder: (ctx) => AlertDialog(
                      title: const Text('Dissolve Household?'),
                      content: const Text('This action cannot be undone. All members will be removed.'),
                      actions: [
                        TextButton(onPressed: () => ctx.pop(false), child: const Text('Cancel')),
                        TextButton(onPressed: () => ctx.pop(true), child: const Text('Dissolve', style: TextStyle(color: Colors.red))),
                      ],
                    ),
                  );
                  if (confirmed == true) {
                    try {
                      await ref.read(householdRepositoryProvider).dissolveHousehold(householdId);
                      ref.invalidate(householdsProvider);
                      if (context.mounted) context.go('/households');
                    } catch (e) {
                      if (context.mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
                    }
                  }
                },
                icon: const Icon(Icons.warning, color: Colors.white),
                label: const Text('Dissolve Household'),
                style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
              ),
            ],
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (e, st) => Center(child: Text('Error: $e')),
      ),
    );
  }
}

class _MemberRoleTile extends ConsumerStatefulWidget {
  final String householdId;
  final dynamic member;

  const _MemberRoleTile({required this.householdId, required this.member});

  @override
  ConsumerState<_MemberRoleTile> createState() => _MemberRoleTileState();
}

class _MemberRoleTileState extends ConsumerState<_MemberRoleTile> {
  late String _selectedRole;
  bool _isUpdating = false;

  @override
  void initState() {
    super.initState();
    _selectedRole = widget.member.role;
  }

  Future<void> _updateRole(String newRole) async {
    setState(() { _isUpdating = true; _selectedRole = newRole; });
    try {
      await ref.read(householdRepositoryProvider).updateMemberRole(widget.householdId, widget.member.userId, newRole);
      ref.invalidate(householdDetailProvider(widget.householdId));
    } catch (e) {
      setState(() { _selectedRole = widget.member.role; });
      if (mounted) ScaffoldMessenger.of(context).showSnackBar(SnackBar(content: Text('Error: $e')));
    } finally {
      if (mounted) setState(() { _isUpdating = false; });
    }
  }

  @override
  Widget build(BuildContext context) {
    final roles = ['Admin', 'Member', 'Viewer'];
    return ListTile(
      leading: CircleAvatar(child: Text(widget.member.userId.substring(0, 1).toUpperCase())),
      title: Text(widget.member.userId),
      subtitle: _isUpdating
          ? const Text('Updating...')
          : DropdownButton<String>(
              value: roles.contains(_selectedRole) ? _selectedRole : 'Member',
              items: roles.map((r) => DropdownMenuItem(value: r, child: Text(r))).toList(),
              onChanged: (val) {
                if (val != null && val != _selectedRole) _updateRole(val);
              },
            ),
    );
  }
}
