import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../providers/auth_provider.dart';
import '../utils/ui_utils.dart';
import '../screens/login_screen.dart';
import '../screens/edit_profile_screen.dart';
import '../screens/family_member_detail_screen.dart';

class ProfileScreen extends StatelessWidget {
  const ProfileScreen({super.key});

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
                            Navigator.of(context, rootNavigator: true).pushAndRemoveUntil(
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









  void _showAddFamilyMemberDetailsModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    String selectedRelation = 'child';
    String? selectedBloodGroup;
    DateTime? selectedDob;
    final dobCtrl = TextEditingController();

    final relationships = ['spouse', 'child', 'parent', 'sibling', 'other'];
    final bloodGroups = ['A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-'];
    final colors = ['#0D9488', '#059669', '#7C3AED', '#DB2777', '#EA580C', '#2563EB'];

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setDialogState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(context).viewInsets.bottom + 32,
            left: 24, right: 24, top: 24,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Add Family Member Profile',
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                const SizedBox(height: 16),
                const Text('Name', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                const SizedBox(height: 8),
                TextField(
                  controller: nameCtrl,
                  decoration: InputDecoration(
                    hintText: 'Enter name',
                    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                ),
                const SizedBox(height: 16),
                const Text('Relationship', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                const SizedBox(height: 8),
                DropdownButtonFormField<String>(
                  value: selectedRelation,
                  decoration: InputDecoration(border: OutlineInputBorder(borderRadius: BorderRadius.circular(12))),
                  items: relationships.map((r) => DropdownMenuItem(value: r, child: Text(r))).toList(),
                  onChanged: (v) { if (v != null) setDialogState(() => selectedRelation = v); },
                ),
                const SizedBox(height: 16),
                const Text('Blood Group (Optional)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                const SizedBox(height: 8),
                DropdownButtonFormField<String>(
                  value: selectedBloodGroup,
                  hint: const Text('Select Blood Group'),
                  decoration: InputDecoration(border: OutlineInputBorder(borderRadius: BorderRadius.circular(12))),
                  items: bloodGroups.map((bg) => DropdownMenuItem(value: bg, child: Text(bg))).toList(),
                  onChanged: (v) { setDialogState(() => selectedBloodGroup = v); },
                ),
                const SizedBox(height: 16),
                const Text('Date of Birth (Optional)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                const SizedBox(height: 8),
                TextField(
                  controller: dobCtrl,
                  readOnly: true,
                  decoration: InputDecoration(hintText: 'Select Date', border: OutlineInputBorder(borderRadius: BorderRadius.circular(12))),
                  onTap: () async {
                    final picked = await showDatePicker(
                      context: context,
                      initialDate: DateTime.now().subtract(const Duration(days: 3650)),
                      firstDate: DateTime(1900),
                      lastDate: DateTime.now(),
                    );
                    if (picked != null) {
                      setDialogState(() {
                        selectedDob = picked;
                        dobCtrl.text = picked.toIso8601String().split('T').first;
                      });
                    }
                  },
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity, height: 52,
                  child: ElevatedButton(
                    onPressed: () async {
                      if (nameCtrl.text.isEmpty) return;
                      final avatarCol = (colors..shuffle()).first;
                      try {
                        await Provider.of<FinancialProvider>(context, listen: false)
                            .addFamilyMember({
                          'name': nameCtrl.text.trim(),
                          'relationship': selectedRelation,
                          'blood_group': selectedBloodGroup,
                          'dob': selectedDob?.toIso8601String().split('T').first,
                          'avatar_color': avatarCol,
                        });
                        if (context.mounted) {
                          Navigator.pop(context);
                          UiUtils.showSnack(context, 'Family member profile added');
                        }
                      } catch (e) {
                        if (context.mounted) UiUtils.showSnack(context, 'Failed: $e', isError: true);
                      }
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF0D9488), foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                    ),
                    child: const Text('Add Profile', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _costRow(String emoji, String label, double amount, bool isLast) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 6),
      child: Row(children: [
        Text(emoji, style: const TextStyle(fontSize: 14)),
        const SizedBox(width: 8),
        Text(label, style: const TextStyle(fontSize: 13, color: Colors.black87)),
        const Spacer(),
        Text('₹${amount.toStringAsFixed(0)}', style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.black87)),
      ]),
    );
  }

  // ─── Build ─────────────────────────────────────────────────────────────────
  @override
  Widget build(BuildContext context) {
    final authProvider = Provider.of<AuthProvider>(context);
    final user = authProvider.user;
    final userName = user?['name'] ?? 'Guest';
    final userEmail = user?['email'] ?? 'guest@horizon.com';
    final userInitials = userName.isNotEmpty ? userName.substring(0, 1).toUpperCase() : 'G';
    final size = MediaQuery.of(context).size;

    return Scaffold(
      body: Stack(
        children: [
          Positioned.fill(
            child: const DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft, end: Alignment.bottomRight,
                  colors: [Color(0xFFF0FDFA), Color(0xFFF8FAFC), Color(0xFFF5F3FF)],
                ),
              ),
            ),
          ),
          Positioned(
            top: -size.height * 0.1, right: -size.width * 0.2,
            child: Container(width: size.width * 0.6, height: size.width * 0.6,
              decoration: BoxDecoration(shape: BoxShape.circle,
                gradient: RadialGradient(colors: [const Color(0xFF0D9488).withValues(alpha: 0.12), const Color(0xFF0D9488).withValues(alpha: 0.0)]))),
          ),
          Positioned(
            bottom: -size.height * 0.08, left: -size.width * 0.15,
            child: Container(width: size.width * 0.5, height: size.width * 0.5,
              decoration: BoxDecoration(shape: BoxShape.circle,
                gradient: RadialGradient(colors: [const Color(0xFF909AC6).withValues(alpha: 0.10), const Color(0xFF909AC6).withValues(alpha: 0.0)]))),
          ),
          Positioned.fill(child: CustomPaint(painter: _GridPainter())),
          Positioned.fill(
            child: ClipRRect(
              child: BackdropFilter(
                filter: ui.ImageFilter.blur(sigmaX: 12, sigmaY: 12),
                child: Container(
                  color: Colors.white.withValues(alpha: 0.45),
                  foregroundDecoration: BoxDecoration(border: Border.all(color: Colors.white.withValues(alpha: 0.25), width: 1)),
                ),
              ),
            ),
          ),
          SafeArea(
            child: SingleChildScrollView(
              child: Column(
                children: [
                  // ── Glass Header ────────────────────────────────────────
                  Container(
                    width: double.infinity,
                    padding: const EdgeInsets.fromLTRB(24, 16, 24, 24),
                    decoration: BoxDecoration(
                      color: Colors.white.withValues(alpha: 0.55),
                      border: Border(bottom: BorderSide(color: Colors.white.withValues(alpha: 0.5), width: 1)),
                    ),
                    child: Column(children: [
                      CircleAvatar(radius: 42, backgroundColor: const Color(0xFF0D9488).withValues(alpha: 0.15),
                        child: Text(userInitials, style: const TextStyle(color: Color(0xFF0D9488), fontSize: 30, fontWeight: FontWeight.bold))),
                      const SizedBox(height: 12),
                      Text(userName, style: const TextStyle(fontSize: 22, fontWeight: FontWeight.w900, color: Color(0xFF0F172A))),
                      const SizedBox(height: 2),
                      Text(userEmail, style: const TextStyle(color: Color(0xFF64748B), fontWeight: FontWeight.w500, fontSize: 14)),
                    ]),
                  ),

                  Padding(
                    padding: const EdgeInsets.all(20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                  // ── Recurring Family Costs ──────────────────────────────────────────
                  Consumer<FinancialProvider>(
                    builder: (context, provider, child) {
                      final costs = provider.recurringCostsBreakdown;
                      if (costs == null || provider.recurringFamilyCosts == 0) {
                        return const SizedBox.shrink();
                      }
                      final schooling = (costs['schooling'] ?? 0).toDouble();
                      final medicines = (costs['medicines'] ?? 0).toDouble();
                      final checkups = (costs['checkups'] ?? 0).toDouble();
                      final vaccinations = (costs['vaccinations'] ?? 0).toDouble();
                      final insurance = (costs['insurance'] ?? 0).toDouble();

                      return Container(
                        margin: const EdgeInsets.only(bottom: 20),
                        width: double.infinity,
                        padding: const EdgeInsets.all(16),
                        decoration: BoxDecoration(
                          color: const Color(0xFF0D9488).withValues(alpha: 0.08),
                          borderRadius: BorderRadius.circular(16),
                          border: Border.all(color: const Color(0xFF0D9488).withValues(alpha: 0.15)),
                        ),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(children: [
                              Container(
                                padding: const EdgeInsets.all(6),
                                decoration: BoxDecoration(color: const Color(0xFF0D9488).withValues(alpha: 0.12), shape: BoxShape.circle),
                                child: const Icon(Icons.people_alt_rounded, color: Color(0xFF0D9488), size: 16),
                              ),
                              const SizedBox(width: 10),
                              Text(
                                'Recurring Family Costs: ₹${provider.recurringFamilyCosts.toStringAsFixed(0)}/mo',
                                style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15, color: Color(0xFF0D9488)),
                              ),
                            ]),
                            const SizedBox(height: 12),
                            if (schooling > 0)
                              _costRow('📚', 'Schooling', schooling, false),
                            if (medicines > 0)
                              _costRow('💊', 'Medicines', medicines, false),
                            if (checkups > 0)
                              _costRow('🩺', 'Checkups', checkups, false),
                            if (vaccinations > 0)
                              _costRow('💉', 'Vaccinations', vaccinations, false),
                            if (insurance > 0)
                              _costRow('🛡️', 'Insurance', insurance, true),
                            const SizedBox(height: 8),
                            const Text('→ Auto-feed to Budget page', style: TextStyle(fontSize: 11, fontStyle: FontStyle.italic, color: const Color(0xFF64748B))),
                          ],
                        ),
                      );
                    },
                  ),

                  // ── Personal Information Navigation ────────────────────────────────
                  Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.white.withValues(alpha: 0.55),
                      borderRadius: BorderRadius.circular(16),
                      border: Border.all(color: Colors.white.withValues(alpha: 0.5), width: 1),
                    ),
                    child: InkWell(
                      onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const EditProfileScreen())),
                      borderRadius: BorderRadius.circular(12),
                      child: Row(children: [
                        Container(
                          padding: const EdgeInsets.all(10),
                          decoration: BoxDecoration(color: const Color(0xFF0D9488).withValues(alpha: 0.1), shape: BoxShape.circle),
                          child: const Icon(Icons.person_outline, color: Color(0xFF0D9488), size: 22),
                        ),
                        const SizedBox(width: 12),
                        Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                          const Text('Personal Information', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Color(0xFF0F172A))),
                          const SizedBox(height: 2),
                          const Text('Update your name, email, and phone', style: TextStyle(fontSize: 13, color: Color(0xFF64748B))),
                        ])),
                        const Icon(Icons.chevron_right_rounded, color: Color(0xFF64748B)),
                      ]),
                    ),
                  ),
                  const SizedBox(height: 20),

            

// ── Family Members Profile List ─────────────────────────────────────
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Family Profiles', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                      TextButton.icon(
                        onPressed: () => _showAddFamilyMemberDetailsModal(context),
                        icon: const Icon(Icons.add_rounded, color: Color(0xFF0D9488), size: 20),
                        label: const Text('Add Profile', style: TextStyle(color: Color(0xFF0D9488), fontWeight: FontWeight.bold)),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Consumer<FinancialProvider>(
                    builder: (context, provider, child) {
                      if (provider.familyMembers.isEmpty) {
                        return const Padding(
                          padding: EdgeInsets.symmetric(vertical: 8),
                          child: Text('No family member profiles found.', style: TextStyle(color: const Color(0xFF64748B))),
                        );
                      }
                      return Container(
                        height: 120,
                        margin: const EdgeInsets.only(bottom: 24),
                        child: ListView.builder(
                          scrollDirection: Axis.horizontal,
                          itemCount: provider.familyMembers.length,
                          itemBuilder: (context, idx) {
                            final m = provider.familyMembers[idx];
                            Color avatarColor;
                            try {
                              avatarColor = Color(int.parse(m.avatarColor!.replaceFirst('#', '0xFF')));
                            } catch (_) {
                              avatarColor = const Color(0xFF0D9488);
                            }
                            return GestureDetector(
                              onTap: () {
                                Navigator.push(
                                  context,
                                  MaterialPageRoute(builder: (_) => FamilyMemberDetailScreen(member: m)),
                                );
                              },
                              child: Container(
                                width: 110,
                                margin: const EdgeInsets.only(right: 12),
                                padding: const EdgeInsets.all(12),
                                decoration: BoxDecoration(
                                  color: Colors.white,
                                  borderRadius: BorderRadius.circular(16),
                                  boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 8, offset: const Offset(0, 2))],
                                ),
                                child: Column(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    CircleAvatar(
                                      radius: 24,
                                      backgroundColor: avatarColor,
                                      child: Text(
                                        m.name.isNotEmpty ? m.name.substring(0, 1).toUpperCase() : 'F',
                                        style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16),
                                      ),
                                    ),
                                    const SizedBox(height: 8),
                                    Text(
                                      m.name,
                                      maxLines: 1,
                                      overflow: TextOverflow.ellipsis,
                                      style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Color(0xFF1E293B)),
                                    ),
                                    Text(
                                      m.relationship,
                                      style: const TextStyle(fontSize: 11, color: const Color(0xFF64748B)),
                                    ),
                                  ],
                                ),
                              ),
                            );
                          },
                        ),
                      );
                    },
                  ),



            // ── Family Sharing ──────────────────────────────────────────────────
            Container(
              width: double.infinity,
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha: 0.55),
                borderRadius: BorderRadius.circular(16),
                border: Border.all(color: Colors.white.withValues(alpha: 0.5), width: 1),
              ),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Row(children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(color: const Color(0xFF0D9488).withValues(alpha: 0.1), shape: BoxShape.circle),
                    child: const Icon(Icons.people_alt_outlined, color: Color(0xFF0D9488), size: 20),
                  ),
                  const SizedBox(width: 12),
                  const Expanded(child: Text('Family Sharing', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Color(0xFF0F172A)))),
                ]),
                const SizedBox(height: 12),
                const Text('Share access with your family. Share your invite code or join an existing household.',
                    style: TextStyle(color: Color(0xFF64748B), fontSize: 13, height: 1.4)),
                const SizedBox(height: 16),
                Consumer<FinancialProvider>(
                  builder: (context, fp, _) {
                    return FutureBuilder<Map<String, dynamic>>(
                      future: fp.loadHouseholdSummary(),
                      builder: (context, snapshot) {
                        final code = snapshot.data?['inviteCode'] as String?;
                        return Row(children: [
                          Expanded(
                            child: Container(
                              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                              decoration: BoxDecoration(
                                color: Colors.white,
                                borderRadius: BorderRadius.circular(10),
                                border: Border.all(color: const Color(0xFF0D9488).withValues(alpha: 0.3)),
                              ),
                              child: Text(
                                code != null ? 'Your Code: $code' : 'Generating invite...',
                                style: TextStyle(
                                  fontWeight: FontWeight.bold, fontSize: 14, letterSpacing: 1,
                                  color: code != null ? const Color(0xFF0D9488) : const Color(0xFF94A3B8),
                                ),
                              ),
                            ),
                          ),
                          const SizedBox(width: 8),
                          if (code != null)
                            GestureDetector(
                              onTap: () => UiUtils.showSnack(context, 'Invite code copied!'),
                              child: Container(
                                padding: const EdgeInsets.all(10),
                                decoration: BoxDecoration(color: const Color(0xFF0D9488).withValues(alpha: 0.1), borderRadius: BorderRadius.circular(10)),
                                child: const Icon(Icons.copy_rounded, color: Color(0xFF0D9488), size: 20),
                              ),
                            ),
                        ]);
                      },
                    );
                  },
                ),
                const SizedBox(height: 16),
                const Divider(height: 1, color: Color(0xFFE2E8F0)),
                const SizedBox(height: 16),
                const Text('Join a Household', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Color(0xFF0F172A))),
                const SizedBox(height: 4),
                const Text('Enter the invite code shared by your family member.', style: TextStyle(color: Color(0xFF64748B), fontSize: 12)),
                const SizedBox(height: 12),
                StatefulBuilder(
                  builder: (context, setState) {
                    final joinCodeCtrl = TextEditingController();
                    return Row(children: [
                      Expanded(
                        child: TextField(
                          controller: joinCodeCtrl,
                          decoration: InputDecoration(
                            hintText: 'Enter invite code',
                            border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                            contentPadding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
                            isDense: true,
                          ),
                          textCapitalization: TextCapitalization.characters,
                        ),
                      ),
                      const SizedBox(width: 12),
                      ElevatedButton(
                        onPressed: () async {
                          final code = joinCodeCtrl.text.trim().toUpperCase();
                          if (code.isEmpty) return;
                          try {
                            final provider = Provider.of<FinancialProvider>(context, listen: false);
                            await provider.joinHousehold(code);
                            if (context.mounted) {
                              UiUtils.showSnack(context, 'Joined household successfully!');
                              joinCodeCtrl.clear();
                            }
                          } catch (e) {
                            if (context.mounted) {
                              UiUtils.showSnack(context, 'Failed: $e', isError: true);
                            }
                          }
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF0D9488),
                          foregroundColor: Colors.white,
                          padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                        ),
                        child: const Text('Join', style: TextStyle(fontWeight: FontWeight.bold)),
                      ),
                    ]);
                  },
                ),
              ]),
            ),
            const SizedBox(height: 24),

            // ── Export Report Button ─────────────────────────────────────────
            SizedBox(
              width: double.infinity,
              height: 56,
              child: ElevatedButton.icon(
                onPressed: () async {
                  try {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Downloading report...')),
                    );
                    final path = await context.read<FinancialProvider>().downloadReport();
                    if (context.mounted) {
                      UiUtils.showSnack(context, 'Report saved to: $path');
                    }
                  } catch (e) {
                    if (context.mounted) {
                      UiUtils.showSnack(context, 'Export failed: $e', isError: true);
                    }
                  }
                },
                icon: const Icon(Icons.description_outlined),
                label: const Text('Export Financial Report (PDF)', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF059669),
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                ),
              ),
            ),
            const SizedBox(height: 24),

            // ── Logout Button ────────────────────────────────────────────────
            SizedBox(
              width: double.infinity,
              height: 56,
              child: OutlinedButton(
                onPressed: () async {
                  await authProvider.logout();
                  if (context.mounted) {
                    Navigator.of(context, rootNavigator: true).pushAndRemoveUntil(
                      MaterialPageRoute(builder: (_) => const LoginScreen()),
                      (route) => false,
                    );
                  }
                },
                style: OutlinedButton.styleFrom(
                  side: const BorderSide(color: Colors.redAccent, width: 2),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                ),
                child: const Text('Log Out',
                    style: TextStyle(color: Colors.redAccent, fontSize: 16, fontWeight: FontWeight.bold)),
              ),
            ),
            const SizedBox(height: 12),
            // ── Delete Account Button ─────────────────────────────────────────
            SizedBox(
              width: double.infinity,
              height: 56,
              child: OutlinedButton(
                onPressed: () => _showDeleteConfirmationDialog(context),
                style: OutlinedButton.styleFrom(
                  side: const BorderSide(color: Colors.red, width: 1),
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                ),
                child: const Text('Delete Account',
                    style: TextStyle(color: Colors.red, fontSize: 16, fontWeight: FontWeight.bold)),
              ),
            ),
            const SizedBox(height: 40),
          ],
        ),
      ),
    ],
      ),
      ),
      ),
    ],
      ),
    );
  }
}

class _GridPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()..color = const Color(0xFF0D9488).withValues(alpha: 0.035)..strokeWidth = 0.5;
    const spacing = 40.0;
    for (double x = 0; x < size.width; x += spacing) canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
    for (double y = 0; y < size.height; y += spacing) canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
  }
  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
