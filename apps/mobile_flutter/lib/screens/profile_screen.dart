import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';
import '../providers/financial_provider.dart';
import '../models/financial_models.dart';
import '../providers/auth_provider.dart';
import '../utils/ui_utils.dart';
import '../screens/login_screen.dart';
import '../screens/insurance_hub_screen.dart';

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
                            color: isSelected ? const Color(0xFF1E3A8A) : Colors.white,
                            borderRadius: BorderRadius.circular(14),
                            border: Border.all(
                              color: isSelected ? const Color(0xFF1E3A8A) : Colors.grey.shade200,
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
                    prefixStyle: const TextStyle(fontWeight: FontWeight.w700, color: Color(0xFF1E3A8A)),
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
                            color: isSelected ? const Color(0xFF1E3A8A) : Colors.white,
                            borderRadius: BorderRadius.circular(14),
                            border: Border.all(
                              color: isSelected ? const Color(0xFF1E3A8A) : Colors.grey.shade200,
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
                      backgroundColor: const Color(0xFF1E3A8A),
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

  // ─── Add Household Member Modal ────────────────────────────────────────────
  void _showAddMemberModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final incomeCtrl = TextEditingController();
    final contribCtrl = TextEditingController();
    final relationCtrl = TextEditingController();

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => Container(
        decoration: const BoxDecoration(
          color: Color(0xFFF8F9FA),
          borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
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
                  const Text('Add Contributing Member',
                      style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                  IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(context)),
                ],
              ),
              const SizedBox(height: 20),
              const Text('Name', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
              const SizedBox(height: 8),
              TextField(
                controller: nameCtrl,
                decoration: InputDecoration(
                  filled: true, fillColor: Colors.white,
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(14), borderSide: BorderSide.none),
                  contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                ),
              ),
              const SizedBox(height: 16),
              const Text('Relationship (e.g. Partner, Father)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
              const SizedBox(height: 8),
              TextField(
                controller: relationCtrl,
                decoration: InputDecoration(
                  filled: true, fillColor: Colors.white,
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(14), borderSide: BorderSide.none),
                  contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                ),
              ),
              const SizedBox(height: 16),
              const Text('Total Monthly Income', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
              const SizedBox(height: 8),
              TextField(
                controller: incomeCtrl,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(
                  prefixText: '₹ ',
                  prefixStyle: const TextStyle(fontWeight: FontWeight.w700, color: Color(0xFF1E3A8A)),
                  filled: true, fillColor: Colors.white,
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(14), borderSide: BorderSide.none),
                  contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                ),
              ),
              const SizedBox(height: 16),
              const Text('Contribution to Household', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
              const SizedBox(height: 8),
              TextField(
                controller: contribCtrl,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(
                  prefixText: '₹ ',
                  prefixStyle: const TextStyle(fontWeight: FontWeight.w700, color: Color(0xFF1E3A8A)),
                  filled: true, fillColor: Colors.white,
                  border: OutlineInputBorder(borderRadius: BorderRadius.circular(14), borderSide: BorderSide.none),
                  contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 14),
                ),
              ),
              const SizedBox(height: 32),
              SizedBox(
                width: double.infinity,
                height: 52,
                child: ElevatedButton(
                  onPressed: () async {
                    if (nameCtrl.text.isEmpty || incomeCtrl.text.isEmpty || contribCtrl.text.isEmpty) return;
                    final income = double.tryParse(incomeCtrl.text) ?? 0.0;
                    final contrib = double.tryParse(contribCtrl.text) ?? 0.0;
                    try {
                      final provider = Provider.of<FinancialProvider>(context, listen: false);
                      await provider.addContributingMember(nameCtrl.text, income, contrib, relationCtrl.text);
                      if (context.mounted) {
                        Navigator.pop(context);
                        UiUtils.showSnack(context, 'Member added successfully');
                      }
                    } catch (e) {
                      if (context.mounted) UiUtils.showSnack(context, 'Failed to add member: $e', isError: true);
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF1E3A8A),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)),
                    elevation: 0,
                  ),
                  child: const Text('Add Member',
                      style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // ─── Income type helpers ───────────────────────────────────────────────────
  Color _typeColor(String type) {
    switch (type.toLowerCase()) {
      case 'salary':   return const Color(0xFF1E3A8A);
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
                  colors: [Color(0xFF0F2057), Color(0xFF1E3A8A), Color(0xFF2563EB)],
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
                      backgroundColor: const Color(0xFF3B82F6),
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
                  // ── Personal Information Card ──────────────────────────────────────
                  Container(
                    padding: const EdgeInsets.all(24),
                    decoration: BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.circular(24),
                      boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.04), blurRadius: 10, offset: const Offset(0, 4))],
                    ),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                  const Text('Personal Information', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                  const SizedBox(height: 24),
                  const Text('Full Name', style: TextStyle(color: Colors.grey, fontWeight: FontWeight.w600, fontSize: 13)),
                  const SizedBox(height: 8),
                  TextField(
                    controller: TextEditingController(text: userName),
                    decoration: InputDecoration(
                      filled: true, fillColor: Colors.grey.withValues(alpha: 0.05),
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
                    ),
                  ),
                  const SizedBox(height: 20),
                  const Text('Email', style: TextStyle(color: Colors.grey, fontWeight: FontWeight.w600, fontSize: 13)),
                  const SizedBox(height: 8),
                  TextField(
                    controller: TextEditingController(text: userEmail),
                    decoration: InputDecoration(
                      filled: true, fillColor: Colors.grey.withValues(alpha: 0.05),
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
                    ),
                  ),
                  const SizedBox(height: 20),
                  const Text('Phone', style: TextStyle(color: Colors.grey, fontWeight: FontWeight.w600, fontSize: 13)),
                  const SizedBox(height: 8),
                  TextField(
                    controller: TextEditingController(text: ''),
                    decoration: InputDecoration(
                      filled: true, fillColor: Colors.grey.withValues(alpha: 0.05),
                      border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none),
                    ),
                  ),
                  const SizedBox(height: 24),
                  Align(
                    alignment: Alignment.centerRight,
                    child: ElevatedButton(
                      onPressed: () {},
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF1E3A8A),
                        padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                      ),
                      child: const Text('Save', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold)),
                    ),
                  ),
                ],
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
                  icon: const Icon(Icons.add_rounded, color: Color(0xFF1E3A8A), size: 20),
                  label: const Text('Add Source',
                      style: TextStyle(color: Color(0xFF1E3A8A), fontWeight: FontWeight.bold)),
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
                            color: const Color(0xFF1E3A8A).withValues(alpha: 0.08),
                            shape: BoxShape.circle,
                          ),
                          child: const Icon(Icons.account_balance_wallet_outlined,
                              color: Color(0xFF1E3A8A), size: 28),
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
                          colors: [Color(0xFF1E3A8A), Color(0xFF2563EB)],
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
                        child: ListTile(
                          contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                          leading: Container(
                            width: 44, height: 44,
                            decoration: BoxDecoration(
                              color: color.withValues(alpha: 0.1),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Icon(_typeIcon(income.type), color: color, size: 22),
                          ),
                          title: Text(income.label,
                              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15)),
                          subtitle: Row(
                            children: [
                              Container(
                                margin: const EdgeInsets.only(top: 4),
                                padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
                                decoration: BoxDecoration(
                                  color: color.withValues(alpha: 0.1),
                                  borderRadius: BorderRadius.circular(20),
                                ),
                                child: Text(
                                  income.type[0].toUpperCase() + income.type.substring(1),
                                  style: TextStyle(
                                    fontSize: 10,
                                    fontWeight: FontWeight.w700,
                                    color: color,
                                  ),
                                ),
                              ),
                              const SizedBox(width: 6),
                              Container(
                                margin: const EdgeInsets.only(top: 4),
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
                          trailing: Row(
                            mainAxisSize: MainAxisSize.min,
                            children: [
                              Text(
                                '₹${income.amount.toStringAsFixed(0)}',
                                style: TextStyle(
                                  fontWeight: FontWeight.w900,
                                  fontSize: 16,
                                  color: color,
                                ),
                              ),
                              const SizedBox(width: 4),
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
                                  margin: const EdgeInsets.only(left: 6),
                                  decoration: BoxDecoration(
                                    color: Colors.red.withValues(alpha: 0.08),
                                    borderRadius: BorderRadius.circular(10),
                                  ),
                                  child: const Icon(Icons.delete_outline_rounded,
                                      color: Colors.redAccent, size: 18),
                                ),
                              ),
                            ],
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
                  icon: const Icon(Icons.add_rounded, color: Color(0xFF1E3A8A), size: 20),
                  label: const Text('Add',
                      style: TextStyle(color: Color(0xFF1E3A8A), fontWeight: FontWeight.bold)),
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
            const SizedBox(height: 16),
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
                  backgroundColor: const Color(0xFF1E3A8A),
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
