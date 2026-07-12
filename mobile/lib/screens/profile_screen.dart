import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../providers/financial_provider.dart';
import '../models/financial_models.dart';
import '../providers/auth_provider.dart';
import '../utils/ui_utils.dart';
import '../screens/login_screen.dart';
import '../screens/insurance_hub_screen.dart';
import '../screens/edit_profile_screen.dart';
import '../screens/family_member_detail_screen.dart';
import '../screens/salary_history_screen.dart';

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

  // ─── Add Income Modal ──────────────────────────────────────────────────────
  void _showAddIncomeModal(BuildContext context) {
    final labelCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    String selectedType = 'salary';
    String selectedFrequency = 'monthly';

    final typeOptions = [
      {'value': 'salary',   'label': 'Salary',   'icon': Icons.work_outline_rounded},
      {'value': 'business', 'label': 'Business', 'icon': Icons.storefront_outlined},
      {'value': 'passive',  'label': 'Passive',  'icon': Icons.trending_up_rounded},
    ];

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Color(0xFFF8F9FA),
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24, right: 24, top: 8,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Handle
                Center(
                  child: Container(
                    width: 40, height: 4,
                    margin: const EdgeInsets.only(bottom: 20),
                    decoration: BoxDecoration(
                      color: Colors.grey.shade300,
                      borderRadius: BorderRadius.circular(4),
                    ),
                  ),
                ),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    const Text('Add Income Source',
                        style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                    IconButton(
                      icon: const Icon(Icons.close),
                      onPressed: () => Navigator.pop(ctx),
                    ),
                  ],
                ),
                const SizedBox(height: 20),

                // Label field
                const Text('Label', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                const SizedBox(height: 8),
                TextField(
                  controller: labelCtrl,
                  textCapitalization: TextCapitalization.words,
                  decoration: InputDecoration(
                    hintText: "e.g. My Salary, Partner Income, Father's Income",
                    hintStyle: TextStyle(color: Colors.grey.shade400, fontSize: 13),
                    filled: true,
                    fillColor: Colors.white,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(14),
                      borderSide: BorderSide.none,
                    ),
                    contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                  ),
                ),
                const SizedBox(height: 20),

                // Income type selector
                const Text('Income Type', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                const SizedBox(height: 10),
                Row(
                  children: typeOptions.map((t) {
                    final isSelected = selectedType == t['value'];
                    return Expanded(
                      child: GestureDetector(
                        onTap: () => setState(() => selectedType = t['value'] as String),
                        child: AnimatedContainer(
                          duration: const Duration(milliseconds: 200),
                          margin: const EdgeInsets.only(right: 8),
                          padding: const EdgeInsets.symmetric(vertical: 12),
                          decoration: BoxDecoration(
                            color: isSelected ? const Color(0xFF0D9488) : Colors.white,
                            borderRadius: BorderRadius.circular(14),
                            border: Border.all(
                              color: isSelected ? const Color(0xFF0D9488) : Colors.grey.shade200,
                            ),
                          ),
                          child: Column(
                            children: [
                              Icon(t['icon'] as IconData,
                                  color: isSelected ? Colors.white : Colors.grey,
                                  size: 22),
                              const SizedBox(height: 6),
                              Text(t['label'] as String,
                                  style: TextStyle(
                                    fontSize: 12,
                                    fontWeight: FontWeight.w600,
                                    color: isSelected ? Colors.white : Colors.grey.shade600,
                                  )),
                            ],
                          ),
                        ),
                      ),
                    );
                  }).toList(),
                ),
                const SizedBox(height: 20),

                // Amount
                const Text('Amount (₹)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                const SizedBox(height: 8),
                TextField(
                  controller: amountCtrl,
                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  decoration: InputDecoration(
                    prefixText: '₹ ',
                    prefixStyle: const TextStyle(fontWeight: FontWeight.w700, color: Color(0xFF0D9488)),
                    filled: true,
                    fillColor: Colors.white,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(14),
                      borderSide: BorderSide.none,
                    ),
                    contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                  ),
                ),
                const SizedBox(height: 20),

                // Frequency toggle
                const Text('Frequency', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                const SizedBox(height: 10),
                Row(
                  children: ['monthly', 'annual'].map((f) {
                    final isSelected = selectedFrequency == f;
                    return Expanded(
                      child: GestureDetector(
                        onTap: () => setState(() => selectedFrequency = f),
                        child: AnimatedContainer(
                          duration: const Duration(milliseconds: 200),
                          margin: EdgeInsets.only(right: f == 'monthly' ? 8 : 0),
                          padding: const EdgeInsets.symmetric(vertical: 14),
                          decoration: BoxDecoration(
                            color: isSelected ? const Color(0xFF0D9488) : Colors.white,
                            borderRadius: BorderRadius.circular(14),
                            border: Border.all(
                              color: isSelected ? const Color(0xFF0D9488) : Colors.grey.shade200,
                            ),
                          ),
                          child: Center(
                            child: Text(
                              f == 'monthly' ? 'Monthly' : 'Annual',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                color: isSelected ? Colors.white : Colors.grey.shade600,
                              ),
                            ),
                          ),
                        ),
                      ),
                    );
                  }).toList(),
                ),
                const SizedBox(height: 32),

                // Add button
                SizedBox(
                  width: double.infinity,
                  height: 52,
                  child: ElevatedButton(
                    onPressed: () async {
                      final label = labelCtrl.text.trim();
                      final amount = double.tryParse(amountCtrl.text.trim());
                      if (label.isEmpty || amount == null || amount <= 0) return;
                      final provider = Provider.of<FinancialProvider>(ctx, listen: false);
                      try {
                        await provider.addIncome(label, selectedType, amount, selectedFrequency);
                        if (ctx.mounted) {
                          Navigator.pop(ctx);
                          UiUtils.showSnack(ctx, 'Income added successfully');
                        }
                      } catch (e) {
                        if (ctx.mounted) UiUtils.showSnack(ctx, 'Failed to add income: $e', isError: true);
                      }
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF0D9488),
                      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                      elevation: 0,
                    ),
                    child: const Text('Add Income Source',
                        style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  // ─── Add Household Member Modal (Glassmorphism) ─────────────────────────
  void _showAddMemberModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final emailCtrl = TextEditingController();
    String selectedRelation = 'Partner';
    bool isEmailInvite = false;

    const relationships = ['Partner', 'Father', 'Mother', 'Brother', 'Sister', 'Son', 'Daughter', 'Friend'];

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setDialogState) => ClipRRect(
          borderRadius: const BorderRadius.vertical(top: Radius.circular(28)),
          child: BackdropFilter(
            filter: ui.ImageFilter.blur(sigmaX: 12, sigmaY: 12),
            child: Container(
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha: 0.9),
                borderRadius: const BorderRadius.vertical(top: Radius.circular(28)),
                border: Border(top: BorderSide(color: Colors.white.withValues(alpha: 0.3))),
              ),
              padding: EdgeInsets.only(
                bottom: MediaQuery.of(context).viewInsets.bottom + 32,
                left: 24, right: 24, top: 8,
              ),
              child: SingleChildScrollView(
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Center(
                      child: Container(
                        width: 40, height: 4,
                        margin: const EdgeInsets.only(bottom: 20),
                        decoration: BoxDecoration(
                          color: Colors.black.withValues(alpha: 0.15),
                          borderRadius: BorderRadius.circular(4),
                        ),
                      ),
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text('Add Family Member',
                            style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                        IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(context)),
                      ],
                    ),
                    const SizedBox(height: 20),
                    const Text('Name', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                    const SizedBox(height: 8),
                    TextField(
                      controller: nameCtrl,
                      style: const TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF1E293B)),
                      decoration: InputDecoration(
                        hintText: 'Enter full name',
                        hintStyle: const TextStyle(color: Color(0xFF94A3B8)),
                        filled: true,
                        fillColor: Colors.white.withValues(alpha: 0.6),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(14),
                          borderSide: BorderSide.none,
                        ),
                        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                      ),
                    ),
                    const SizedBox(height: 16),
                    const Text('Email (enter to send invite)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                    const SizedBox(height: 8),
                    TextField(
                      controller: emailCtrl,
                      keyboardType: TextInputType.emailAddress,
                      style: const TextStyle(fontWeight: FontWeight.w600, color: Color(0xFF1E293B)),
                      decoration: InputDecoration(
                        hintText: 'Optional — leave blank for manual entry',
                        hintStyle: const TextStyle(color: Color(0xFF94A3B8), fontSize: 13),
                        filled: true,
                        fillColor: Colors.white.withValues(alpha: 0.6),
                        border: OutlineInputBorder(
                          borderRadius: BorderRadius.circular(14),
                          borderSide: BorderSide.none,
                        ),
                        contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                      ),
                      onChanged: (_) => setDialogState(() => isEmailInvite = emailCtrl.text.trim().isNotEmpty),
                    ),
                    const SizedBox(height: 16),
                    const Text('Relationship', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                    const SizedBox(height: 8),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 16),
                      decoration: BoxDecoration(
                        color: Colors.white.withValues(alpha: 0.6),
                        borderRadius: BorderRadius.circular(14),
                      ),
                      child: DropdownButtonHideUnderline(
                        child: DropdownButton<String>(
                          value: selectedRelation,
                          isExpanded: true,
                          style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 15, color: Color(0xFF1E293B)),
                          items: relationships.map((r) => DropdownMenuItem(
                            value: r,
                            child: Row(children: [
                              Icon(_relationIcon(r), size: 18, color: const Color(0xFF0D9488)),
                              const SizedBox(width: 10),
                              Text(r),
                            ]),
                          )).toList(),
                          onChanged: (v) {
                            if (v != null) setDialogState(() => selectedRelation = v);
                          },
                        ),
                      ),
                    ),
                    const SizedBox(height: 32),
                    SizedBox(
                      width: double.infinity,
                      height: 52,
                      child: ElevatedButton(
                        onPressed: () async {
                          if (nameCtrl.text.isEmpty) return;
                          try {
                            final provider = Provider.of<FinancialProvider>(context, listen: false);
                            if (isEmailInvite) {
                              final result = await provider.inviteFamilyMember(
                                nameCtrl.text, emailCtrl.text.trim(), selectedRelation);
                              if (context.mounted) {
                                Navigator.pop(context);
                                UiUtils.showSnack(context, 'Invite sent! Code: ${result['inviteCode']}');
                              }
                            } else {
                              await provider.addContributingMember(nameCtrl.text, 0, 0, selectedRelation);
                              if (context.mounted) {
                                Navigator.pop(context);
                                UiUtils.showSnack(context, 'Member added successfully');
                              }
                            }
                          } catch (e) {
                            if (context.mounted) UiUtils.showSnack(context, 'Failed: $e', isError: true);
                          }
                        },
                        style: ElevatedButton.styleFrom(
                          backgroundColor: const Color(0xFF0D9488),
                          foregroundColor: Colors.white,
                          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                          elevation: 0,
                        ),
                        child: Text(isEmailInvite ? 'Send Invite' : 'Add Member',
                            style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  IconData _relationIcon(String relation) {
    switch (relation) {
      case 'Partner': return Icons.favorite_rounded;
      case 'Father': return Icons.man_3_rounded;
      case 'Mother': return Icons.woman_rounded;
      case 'Brother': return Icons.face_6_rounded;
      case 'Sister': return Icons.face_5_rounded;
      case 'Son': return Icons.child_care_rounded;
      case 'Daughter': return Icons.child_friendly_rounded;
      case 'Friend': return Icons.people_rounded;
      default: return Icons.person_rounded;
    }
  }

  IconData _familyRelationIcon(String relation) {
    switch (relation.toLowerCase()) {
      case 'spouse': return Icons.favorite_rounded;
      case 'parent': return Icons.man_3_rounded;
      case 'child': return Icons.child_care_rounded;
      case 'sibling': return Icons.people_rounded;
      default: return Icons.person_rounded;
    }
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
                  decoration: InputDecoration(
                    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                  items: relationships.map((r) => DropdownMenuItem(value: r, child: Text(r))).toList(),
                  onChanged: (v) {
                    if (v != null) setDialogState(() => selectedRelation = v);
                  },
                ),
                const SizedBox(height: 16),
                const Text('Blood Group (Optional)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                const SizedBox(height: 8),
                DropdownButtonFormField<String>(
                  value: selectedBloodGroup,
                  hint: const Text('Select Blood Group'),
                  decoration: InputDecoration(
                    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                  ),
                  items: bloodGroups.map((bg) => DropdownMenuItem(value: bg, child: Text(bg))).toList(),
                  onChanged: (v) {
                    setDialogState(() => selectedBloodGroup = v);
                  },
                ),
                const SizedBox(height: 16),
                const Text('Date of Birth (Optional)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Color(0xFF64748B))),
                const SizedBox(height: 8),
                TextField(
                  controller: dobCtrl,
                  readOnly: true,
                  decoration: InputDecoration(
                    hintText: 'Select Date',
                    border: OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                  ),
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
                  width: double.infinity,
                  height: 52,
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
                      backgroundColor: const Color(0xFF0D9488),
                      foregroundColor: Colors.white,
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

  // ─── Income type helpers ───────────────────────────────────────────────────
  Color _typeColor(String type) {
    switch (type.toLowerCase()) {
      case 'salary':   return const Color(0xFF0D9488);
      case 'business': return const Color(0xFF7C3AED);
      case 'passive':  return const Color(0xFF059669);
      default:         return const Color(0xFF475569);
    }
  }

  IconData _typeIcon(String type) {
    switch (type.toLowerCase()) {
      case 'salary':   return Icons.work_outline_rounded;
      case 'business': return Icons.storefront_outlined;
      case 'passive':  return Icons.trending_up_rounded;
      default:         return Icons.attach_money_rounded;
    }
  }

  String _frequencyLabel(String freq) {
    switch (freq.toLowerCase()) {
      case 'annual': return 'Annual';
      default:       return 'Monthly';
    }
  }

  /// Convert annual income to monthly equivalent for the total
  double _toMonthly(LocalIncome income) {
    if (income.frequency.toLowerCase() == 'annual') return income.amount / 12;
    return income.amount;
  }

  // ─── Build ─────────────────────────────────────────────────────────────────
  @override
  Widget build(BuildContext context) {
    final authProvider = Provider.of<AuthProvider>(context);
    final user = authProvider.user;
    final userName = user?['name'] ?? 'Guest';
    final userEmail = user?['email'] ?? 'guest@horizon.com';
    final userInitials = userName.isNotEmpty ? userName.substring(0, 1).toUpperCase() : 'G';
    final userType = user?['user_type'] ?? 'Standard';

    return Scaffold(
      backgroundColor: const Color(0xFFF1F5F9),
      body: SingleChildScrollView(
        child: Column(
          children: [
            // Background Gradient Header & Avatar
            Container(
              width: double.infinity,
              padding: EdgeInsets.fromLTRB(24, MediaQuery.of(context).padding.top + 24, 24, 40),
              decoration: const BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [Color(0xFF042F2E), Color(0xFF0D9488), Color(0xFF14B8A6)],
                ),
                borderRadius: BorderRadius.vertical(bottom: Radius.circular(40)),
              ),
              child: Column(
                children: [
                  Container(
                    padding: const EdgeInsets.all(4),
                    decoration: BoxDecoration(
                      shape: BoxShape.circle,
                      border: Border.all(color: Colors.white.withValues(alpha: 0.5), width: 3),
                    ),
                    child: CircleAvatar(
                      radius: 46,
                      backgroundColor: const Color(0xFF2DD4BF),
                      child: Text(userInitials, style: const TextStyle(color: Colors.white, fontSize: 32, fontWeight: FontWeight.bold)),
                    ),
                  ),
                  const SizedBox(height: 16),
                  Text(userName, style: const TextStyle(color: Colors.white, fontSize: 24, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 4),
                  Text(userEmail, style: TextStyle(color: Colors.white.withValues(alpha: 0.8), fontWeight: FontWeight.w500)),
                  const SizedBox(height: 12),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
                    decoration: BoxDecoration(
                      color: Colors.white.withValues(alpha: 0.2),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Text(userType[0].toUpperCase() + userType.substring(1), style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 13)),
                  ),
                ],
              ),
            ),
            
            // Rest of Body
            Padding(
              padding: const EdgeInsets.all(24.0),
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
                            Text(
                              'Recurring Family Costs: ₹${provider.recurringFamilyCosts.toStringAsFixed(0)}/mo',
                              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15, color: Color(0xFF0D9488)),
                            ),
                            const SizedBox(height: 8),
                            if (schooling > 0)
                              Text('  ├ Schooling:      ₹${schooling.toStringAsFixed(0)}', style: const TextStyle(fontSize: 13, color: Colors.black87)),
                            if (medicines > 0)
                              Text('  ├ Medicines:      ₹${medicines.toStringAsFixed(0)}', style: const TextStyle(fontSize: 13, color: Colors.black87)),
                            if (checkups > 0)
                              Text('  ├ Checkups:       ₹${checkups.toStringAsFixed(0)}', style: const TextStyle(fontSize: 13, color: Colors.black87)),
                            if (vaccinations > 0)
                              Text('  └ Vaccinations:   ₹${vaccinations.toStringAsFixed(0)}', style: const TextStyle(fontSize: 13, color: Colors.black87)),
                            const SizedBox(height: 8),
                            const Text('→ Auto-feed to Budget page', style: TextStyle(fontSize: 11, fontStyle: FontStyle.italic, color: Colors.grey)),
                          ],
                        ),
                      );
                    },
                  ),

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
                          child: Text('No family member profiles found.', style: TextStyle(color: Colors.grey)),
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
                                      style: const TextStyle(fontSize: 11, color: Colors.grey),
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

                  // ── Personal Information Navigation ────────────────────────────────
                  Container(
                    decoration: BoxDecoration(
                      borderRadius: BorderRadius.circular(16),
                      boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 10, offset: const Offset(0, 4))],
                    ),
                    child: Material(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(16),
                      clipBehavior: Clip.antiAlias,
                      child: ListTile(
                      contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                      leading: Container(
                        padding: const EdgeInsets.all(10),
                        decoration: BoxDecoration(
                          color: const Color(0xFF0D9488).withValues(alpha: 0.1),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: const Icon(Icons.person_outline, color: Color(0xFF0D9488)),
                      ),
                      title: const Text('Personal Information', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                      subtitle: const Text('Update your name, email, and phone', style: TextStyle(fontSize: 13, color: Colors.grey)),
                      trailing: const Icon(Icons.arrow_forward_ios, size: 16, color: Colors.grey),
                      onTap: () {
                        Navigator.push(context, MaterialPageRoute(builder: (_) => const EditProfileScreen()));
                      },
                    ),
                    ),
                  ),
            const SizedBox(height: 40),

            // ── MY INCOME SOURCES ────────────────────────────────────────────
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('My Income Sources',
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                TextButton.icon(
                  onPressed: () => _showAddIncomeModal(context),
                  icon: const Icon(Icons.add_rounded, color: Color(0xFF0D9488), size: 20),
                  label: const Text('Add Source',
                      style: TextStyle(color: Color(0xFF0D9488), fontWeight: FontWeight.bold)),
                ),
              ],
            ),
            const SizedBox(height: 4),
            const Align(
              alignment: Alignment.centerLeft,
              child: Text(
                'Add your salary and any additional income sources.',
                style: TextStyle(color: Colors.grey, fontSize: 13),
              ),
            ),
            const SizedBox(height: 16),
            Consumer<FinancialProvider>(
              builder: (context, provider, child) {
                if (provider.isLoading) {
                  return const Padding(
                    padding: EdgeInsets.symmetric(vertical: 24),
                    child: Center(child: CircularProgressIndicator()),
                  );
                }

                final incomes = provider.incomes;

                if (incomes.isEmpty) {
                  // Empty state
                  return Container(
                    width: double.infinity,
                    padding: const EdgeInsets.all(32),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(20),
                      border: Border.all(color: Colors.grey.shade100),
                    ),
                    child: Column(
                      children: [
                        Container(
                          width: 56, height: 56,
                          decoration: BoxDecoration(
                            color: const Color(0xFF0D9488).withValues(alpha: 0.08),
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(Icons.account_balance_wallet_outlined,
                              color: Color(0xFF0D9488), size: 28),
                        ),
                        const SizedBox(height: 16),
                        const Text('No income sources yet',
                            style: TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
                        const SizedBox(height: 6),
                        const Text('Tap "Add Source" to add your salary\nor any additional income.',
                            style: TextStyle(color: Colors.grey, fontSize: 13),
                            textAlign: TextAlign.center),
                      ],
                    ),
                  );
                }

                // Total monthly income
                final totalMonthly = incomes.fold(0.0, (sum, inc) => sum + _toMonthly(inc));

                return Column(
                  children: [
                    // Total banner
                    Container(
                      width: double.infinity,
                      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 16),
                      margin: const EdgeInsets.only(bottom: 12),
                      decoration: BoxDecoration(
                        gradient: const LinearGradient(
                          colors: [Color(0xFF0D9488), Color(0xFF14B8A6)],
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                        ),
                        borderRadius: BorderRadius.circular(18),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text('Total Monthly Income',
                                  style: TextStyle(color: Colors.white70, fontSize: 12, fontWeight: FontWeight.w600)),
                              SizedBox(height: 4),
                              Text('All sources combined',
                                  style: TextStyle(color: Colors.white54, fontSize: 11)),
                            ],
                          ),
                          Text(
                            '₹${totalMonthly.toStringAsFixed(0)}',
                            style: const TextStyle(
                              color: Colors.white,
                              fontSize: 24,
                              fontWeight: FontWeight.w900,
                              letterSpacing: -0.5,
                            ),
                          ),
                        ],
                      ),
                    ),

                    // Income cards
                    ...incomes.map((income) {
                      final color = _typeColor(income.type);
                      final isSalary = income.type.toLowerCase() == 'salary';
                      return Container(
                        margin: const EdgeInsets.only(bottom: 10),
                        decoration: BoxDecoration(
                          color: Colors.white,
                          borderRadius: BorderRadius.circular(18),
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withValues(alpha: 0.04),
                              blurRadius: 12,
                              offset: const Offset(0, 4),
                            ),
                          ],
                        ),
                        child: Material(
                          color: Colors.transparent,
                          child: InkWell(
                            borderRadius: BorderRadius.circular(18),
                            onTap: isSalary
                                ? () {
                                    Navigator.push(
                                      context,
                                      MaterialPageRoute(
                                        builder: (_) => SalaryHistoryScreen(income: income),
                                      ),
                                    );
                                  }
                                : null,
                            child: Padding(
                              padding: const EdgeInsets.all(12),
                              child: Row(
                                children: [
                                  Container(
                                    width: 44, height: 44,
                                    decoration: BoxDecoration(
                                      color: color.withValues(alpha: 0.1),
                                      borderRadius: BorderRadius.circular(12),
                                    ),
                                    child: Icon(_typeIcon(income.type), color: color, size: 22),
                                  ),
                                  const SizedBox(width: 16),
                                  Expanded(
                                    child: Column(
                                      crossAxisAlignment: CrossAxisAlignment.start,
                                      children: [
                                        Text(income.label, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
                                        const SizedBox(height: 4),
                                        Row(
                                          children: [
                                            Container(
                                              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                                              decoration: BoxDecoration(
                                                color: color.withValues(alpha: 0.1),
                                                borderRadius: BorderRadius.circular(20),
                                              ),
                                              child: Text(
                                                income.type[0].toUpperCase() + income.type.substring(1),
                                                style: TextStyle(fontSize: 10, fontWeight: FontWeight.w700, color: color),
                                              ),
                                            ),
                                            const SizedBox(width: 6),
                                            Container(
                                              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                                              decoration: BoxDecoration(
                                                color: Colors.grey.shade100,
                                                borderRadius: BorderRadius.circular(20),
                                              ),
                                              child: Text(
                                                _frequencyLabel(income.frequency),
                                                style: TextStyle(fontSize: 10, fontWeight: FontWeight.w600, color: Colors.grey.shade600),
                                              ),
                                            ),
                                          ],
                                        ),
                                        if (isSalary) ...[
                                          const SizedBox(height: 6),
                                          const Row(
                                            children: [
                                              Icon(Icons.trending_up_rounded, size: 14, color: Colors.grey),
                                              SizedBox(width: 4),
                                              Text(
                                                'Tap to view salary history & growth',
                                                style: TextStyle(fontSize: 11, color: Colors.grey, fontWeight: FontWeight.w500),
                                              ),
                                            ],
                                          ),
                                        ],
                                      ],
                                    ),
                                  ),
                                  Column(
                                    crossAxisAlignment: CrossAxisAlignment.end,
                                    children: [
                                      Text(
                                        '₹${income.amount.toStringAsFixed(0)}',
                                        style: TextStyle(
                                          fontWeight: FontWeight.w900,
                                          fontSize: 16,
                                          color: color,
                                        ),
                                      ),
                                      const SizedBox(height: 6),
                                      GestureDetector(
                                        onTap: () async {
                                          UiUtils.showDeleteBottomSheet(context, income.label, () async {
                                            await Provider.of<FinancialProvider>(context, listen: false)
                                                .deleteIncome(income.id);
                                            if (context.mounted) UiUtils.showSnack(context, 'Income removed');
                                          });
                                        },
                                        child: Container(
                                          width: 32, height: 32,
                                          decoration: BoxDecoration(
                                            color: Colors.red.withValues(alpha: 0.08),
                                            borderRadius: BorderRadius.circular(10),
                                          ),
                                          child: const Icon(Icons.delete_outline_rounded, color: Colors.redAccent, size: 18),
                                        ),
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      );
                    }),
                  ],
                );
              },
            ),
            const SizedBox(height: 40),

            // ── Household Members ────────────────────────────────────────────
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Household Members',
                    style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                TextButton.icon(
                  onPressed: () => _showAddMemberModal(context),
                  icon: const Icon(Icons.add_rounded, color: Color(0xFF0D9488), size: 20),
                  label: const Text('Add',
                      style: TextStyle(color: Color(0xFF0D9488), fontWeight: FontWeight.bold)),
                ),
              ],
            ),
            const SizedBox(height: 4),
            const Align(
              alignment: Alignment.centerLeft,
              child: Text(
                'Add contributing members to share your financial journey and goals.',
                style: TextStyle(color: Colors.grey, fontSize: 13),
              ),
            ),
            const SizedBox(height: 12),
            Consumer<FinancialProvider>(
              builder: (context, fp, child) {
                return FutureBuilder<Map<String, dynamic>>(
                  future: fp.loadHouseholdSummary(),
                  builder: (context, snapshot) {
                    if (!snapshot.hasData || snapshot.data?['inviteCode'] == null) {
                      return const SizedBox.shrink();
                    }
                    final code = snapshot.data!['inviteCode'];
                    return Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: const Color(0xFF0D9488).withValues(alpha: 0.06),
                        borderRadius: BorderRadius.circular(14),
                        border: Border.all(color: const Color(0xFF0D9488).withValues(alpha: 0.15)),
                      ),
                      child: Row(children: [
                        const Icon(Icons.vpn_key_rounded, color: Color(0xFF0D9488), size: 20),
                        const SizedBox(width: 10),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const Text('Family Invite Code', style: TextStyle(fontSize: 11, color: Colors.grey, fontWeight: FontWeight.w600)),
                              const SizedBox(height: 2),
                              Text(code, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Color(0xFF0D9488), letterSpacing: 1.5)),
                            ],
                          ),
                        ),
                        GestureDetector(
                          onTap: () {
                            UiUtils.showSnack(context, 'Code copied!');
                          },
                          child: Container(
                            padding: const EdgeInsets.all(8),
                            decoration: BoxDecoration(
                              color: const Color(0xFF0D9488).withValues(alpha: 0.1),
                              borderRadius: BorderRadius.circular(10),
                            ),
                            child: const Icon(Icons.copy_rounded, color: Color(0xFF0D9488), size: 18),
                          ),
                        ),
                      ]),
                    );
                  },
                );
              },
            ),
            const SizedBox(height: 12),
            Consumer<FinancialProvider>(
              builder: (context, provider, child) {
                if (provider.householdMembers.isEmpty) {
                  return const Padding(
                    padding: EdgeInsets.symmetric(vertical: 8),
                    child: Text('No members added yet.', style: TextStyle(color: Colors.grey)),
                  );
                }
                return Column(
                  children: provider.householdMembers.map((m) => Container(
                    margin: const EdgeInsets.only(bottom: 12),
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(18),
                      boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 10, offset: const Offset(0, 4))],
                    ),
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(m.name, style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                            Text(m.relationship, style: const TextStyle(color: Colors.grey, fontSize: 13)),
                          ],
                        ),
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.end,
                          children: [
                            Text('+₹${m.contribution.toStringAsFixed(0)}',
                                style: const TextStyle(fontWeight: FontWeight.w900, color: Colors.green, fontSize: 16)),
                            const Text('Contribution', style: TextStyle(color: Colors.grey, fontSize: 11)),
                          ],
                        )
                      ],
                    ),
                  )).toList(),
                );
              },
            ),
            const SizedBox(height: 40),


            // ── Insurance Hub Button ─────────────────────────────────────────
            SizedBox(
              width: double.infinity,
              height: 56,
              child: ElevatedButton.icon(
                onPressed: () {
                  Navigator.push(context, MaterialPageRoute(builder: (_) => const InsuranceHubScreen()));
                },
                icon: const Icon(Icons.shield_rounded),
                label: const Text('Insurance Hub', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF0D9488),
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                ),
              ),
            ),
            const SizedBox(height: 16),

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
    );
  }
}
