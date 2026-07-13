import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/financial_models.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class FamilyMemberDetailScreen extends StatefulWidget {
  final LocalFamilyMember member;
  const FamilyMemberDetailScreen({super.key, required this.member});

  @override
  State<FamilyMemberDetailScreen> createState() =>
      _FamilyMemberDetailScreenState();
}

class _FamilyMemberDetailScreenState extends State<FamilyMemberDetailScreen> {
  late LocalFamilyMember _member;
  bool _schoolingExpanded = false;
  bool _healthExpanded = false;
  bool _medicineExpanded = false;
  bool _vaccinationExpanded = false;
  bool _insuranceExpanded = false;
  bool _earningsExpanded = false;

  @override
  void initState() {
    super.initState();
    _member = widget.member;
    _schoolingExpanded = _member.relationship == 'child' || _member.schooling.isNotEmpty;
  }

  Color _avatarColor() {
    try {
      return Color(int.parse(_member.avatarColor!.replaceFirst('#', '0xFF')));
    } catch (_) {
      return const Color(0xFF6366F1);
    }
  }

  void _showDeleteConfirmation() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
        title: const Text('Delete Profile'),
        content: Text(
          'Remove ${_member.name} permanently? All records will be deleted.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              try {
                await Provider.of<FinancialProvider>(
                  context,
                  listen: false,
                ).deleteFamilyMember(_member.id);
                if (ctx.mounted) Navigator.pop(ctx);
                if (mounted) {
                  Navigator.pop(context);
                  UiUtils.showSnack(context, 'Member removed');
                }
              } catch (e) {
                if (ctx.mounted)
                  UiUtils.showSnack(ctx, 'Failed: $e', isError: true);
              }
            },
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Delete', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
    );
  }

  void _showEditMemberModal() {
    final nameCtrl = TextEditingController(text: _member.name);
    final phoneCtrl = TextEditingController(text: _member.phone ?? '');
    final emailCtrl = TextEditingController(text: _member.email ?? '');
    String rel = _member.relationship;
    String? bg = _member.bloodGroup;
    DateTime? dob;
    bool earningStatus = _member.earningStatus;
    double contribution = _member.contributionAmount;
    final dobCtrl = TextEditingController(text: _member.dob ?? '');

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24,
            right: 24,
            top: 24,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text(
                  'Edit Profile',
                  style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: nameCtrl,
                  decoration: const InputDecoration(
                    labelText: 'Name',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  initialValue: rel,
                  decoration: const InputDecoration(
                    labelText: 'Relationship',
                    border: OutlineInputBorder(),
                  ),
                  items: ['spouse', 'child', 'parent', 'sibling', 'other']
                      .map((r) => DropdownMenuItem(value: r, child: Text(r)))
                      .toList(),
                  onChanged: (v) {
                    if (v != null) setState(() => rel = v);
                  },
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String?>(
                  initialValue: bg,
                  decoration: const InputDecoration(
                    labelText: 'Blood Group',
                    border: OutlineInputBorder(),
                  ),
                  items:
                      [null, 'A+', 'A-', 'B+', 'B-', 'AB+', 'AB-', 'O+', 'O-']
                          .map(
                            (b) => DropdownMenuItem(
                              value: b,
                              child: Text(b ?? 'Select'),
                            ),
                          )
                          .toList(),
                  onChanged: (v) => setState(() => bg = v),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: dobCtrl,
                  readOnly: true,
                  decoration: const InputDecoration(
                    labelText: 'Date of Birth',
                    border: OutlineInputBorder(),
                  ),
                  onTap: () async {
                    final d = await showDatePicker(
                      context: ctx,
                      initialDate: DateTime.now().subtract(
                        const Duration(days: 3650),
                      ),
                      firstDate: DateTime(1900),
                      lastDate: DateTime.now(),
                    );
                    if (d != null) {
                      setState(() {
                        dob = d;
                        dobCtrl.text = d.toIso8601String().split('T').first;
                      });
                    }
                  },
                ),
                const SizedBox(height: 12),
                SwitchListTile(
                  title: const Text('Earning Member'),
                  value: earningStatus,
                  contentPadding: EdgeInsets.zero,
                  onChanged: (v) => setState(() => earningStatus = v),
                ),
                if (earningStatus)
                  TextField(
                    keyboardType: TextInputType.number,
                    decoration: const InputDecoration(
                      labelText: 'Monthly Contribution (₹)',
                      border: OutlineInputBorder(),
                    ),
                    controller: TextEditingController(
                      text: contribution > 0
                          ? contribution.toStringAsFixed(0)
                          : '',
                    ),
                    onChanged: (v) => contribution = double.tryParse(v) ?? 0,
                  ),
                const SizedBox(height: 12),
                TextField(
                  controller: phoneCtrl,
                  keyboardType: TextInputType.phone,
                  decoration: const InputDecoration(
                    labelText: 'Phone Number (Optional)',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: emailCtrl,
                  keyboardType: TextInputType.emailAddress,
                  decoration: const InputDecoration(
                    labelText: 'Email Address (Optional)',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 52,
                  child: ElevatedButton(
                    onPressed: () async {
                      if (nameCtrl.text.isEmpty) return;
                      await Provider.of<FinancialProvider>(
                        context,
                        listen: false,
                      ).updateFamilyMember(_member.id, {
                        'name': nameCtrl.text.trim(),
                        'relationship': rel,
                        'blood_group': bg,
                        'dob':
                            dob?.toIso8601String().split('T').first ??
                            _member.dob,
                        'earning_status': earningStatus,
                        'contribution_amount': contribution,
                        'phone': phoneCtrl.text.trim().isEmpty ? null : phoneCtrl.text.trim(),
                        'email': emailCtrl.text.trim().isEmpty ? null : emailCtrl.text.trim(),
                      });
                      if (ctx.mounted) Navigator.pop(ctx);
                      if (mounted)
                        UiUtils.showSnack(context, 'Profile updated');
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF6366F1),
                    ),
                    child: const Text(
                      'Save',
                      style: TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final size = MediaQuery.of(context).size;

    return Scaffold(
      body: Stack(
        children: [
          Positioned.fill(
            child: const DecoratedBox(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    Color(0xFFF0FDFA),
                    Color(0xFFF8FAFC),
                    Color(0xFFF5F3FF),
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            top: -size.height * 0.1,
            right: -size.width * 0.2,
            child: Container(
              width: size.width * 0.6,
              height: size.width * 0.6,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFF6366F1).withValues(alpha: 0.12),
                    const Color(0xFF6366F1).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),
          Positioned.fill(child: CustomPaint(painter: _GridPainter())),
          Positioned.fill(
            child: ClipRRect(
              child: BackdropFilter(
                filter: ui.ImageFilter.blur(sigmaX: 12, sigmaY: 12),
                child: Container(
                  color: Colors.white.withValues(alpha: 0.45),
                  foregroundDecoration: BoxDecoration(
                    border: Border.all(
                      color: Colors.white.withValues(alpha: 0.25),
                      width: 1,
                    ),
                  ),
                ),
              ),
            ),
          ),
          SafeArea(
            child: Column(
              children: [
                _buildHeader(),
                Expanded(child: _buildSections()),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHeader() {
    return Container(
      padding: const EdgeInsets.fromLTRB(20, 12, 20, 16),
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.55),
        border: Border(
          bottom: BorderSide(
            color: Colors.white.withValues(alpha: 0.5),
            width: 1,
          ),
        ),
      ),
      child: Row(
        children: [
          IconButton(
            icon: const Icon(
              Icons.arrow_back_rounded,
              color: Color(0xFF0F172A),
            ),
            onPressed: () => Navigator.pop(context),
          ),
          CircleAvatar(
            radius: 32,
            backgroundColor: _avatarColor(),
            child: Text(
              _member.name.isNotEmpty ? _member.name[0].toUpperCase() : '?',
              style: const TextStyle(
                color: Colors.white,
                fontSize: 28,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  _member.name,
                  style: const TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.w900,
                    color: Color(0xFF0F172A),
                  ),
                ),
                const SizedBox(height: 2),
                Row(
                  children: [
                    if (_member.relationship.isNotEmpty)
                      _infoChip(
                        _member.relationship[0].toUpperCase() +
                            _member.relationship.substring(1),
                        const Color(0xFF6366F1),
                      ),
                    if (_member.bloodGroup != null) ...[
                      const SizedBox(width: 6),
                      _infoChip(_member.bloodGroup!, const Color(0xFF909AC6)),
                    ],
                  ],
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.edit_outlined, color: Color(0xFF6366F1)),
            onPressed: _showEditMemberModal,
          ),
          IconButton(
            icon: const Icon(Icons.delete_outline, color: Colors.redAccent),
            onPressed: _showDeleteConfirmation,
          ),
        ],
      ),
    );
  }

  Widget _infoChip(String text, Color color) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(10),
      ),
      child: Text(
        text,
        style: TextStyle(
          fontSize: 10,
          fontWeight: FontWeight.w700,
          color: color,
        ),
      ),
    );
  }

  Widget _buildSections() {
    return ListView(
      padding: const EdgeInsets.fromLTRB(16, 12, 16, 32),
      children: [
        _buildSection(
          title: 'Schooling',
          icon: Icons.school_rounded,
          color: const Color(0xFF6366F1),
          expanded: _schoolingExpanded,
          onToggle: (v) => setState(() => _schoolingExpanded = v),
          child: _member.schooling.isEmpty
              ? _buildEmptySection(
                  'No schooling records',
                  Icons.school_outlined,
                  () => _showAddSchoolingModal(),
                )
              : Column(
                  children: [
                    ..._member.schooling.map((s) => _buildSchoolingCard(s)),
                    OutlinedButton.icon(
                      onPressed: _showFeeLedgerModal,
                      icon: const Icon(Icons.receipt_long_rounded, size: 14),
                      label: const Text(
                        'View Fee Ledger',
                        style: TextStyle(fontSize: 12),
                      ),
                      style: OutlinedButton.styleFrom(
                        foregroundColor: const Color(0xFF6366F1),
                      ),
                    ),
                  ],
                ),
        ),
        _buildSection(
          title: 'Health Checkups',
          icon: Icons.favorite_rounded,
          color: const Color(0xFFE88A1A),
          expanded: _healthExpanded,
          onToggle: (v) => setState(() => _healthExpanded = v),
          child: _member.checkups.isEmpty
              ? _buildEmptySection(
                  'No checkup records',
                  Icons.favorite_border,
                  () => _showAddCheckupModal(),
                )
              : Column(
                  children: _member.checkups
                      .map((c) => _buildCheckupCard(c))
                      .toList(),
                ),
        ),
        _buildSection(
          title: 'Regular Medicines',
          icon: Icons.medication_rounded,
          color: const Color(0xFF7C3AED),
          expanded: _medicineExpanded,
          onToggle: (v) => setState(() => _medicineExpanded = v),
          child: _member.medicines.isEmpty
              ? _buildEmptySection(
                  'No medicine records',
                  Icons.medication_outlined,
                  () => _showAddMedicineModal(),
                )
              : Column(
                  children: _member.medicines
                      .map((m) => _buildMedicineCard(m))
                      .toList(),
                ),
        ),
        _buildSection(
          title: 'Vaccinations',
          icon: Icons.vaccines_rounded,
          color: const Color(0xFF2563EB),
          expanded: _vaccinationExpanded,
          onToggle: (v) => setState(() => _vaccinationExpanded = v),
          child: _member.vaccinations.isEmpty
              ? _buildEmptySection(
                  'No vaccination records',
                  Icons.vaccines_outlined,
                  () => _showAddVaccinationModal(),
                )
              : Column(
                  children: _member.vaccinations
                      .map((v) => _buildVaccinationCard(v))
                      .toList(),
                ),
        ),
        _buildSection(
          title: 'Insurance',
          icon: Icons.shield_rounded,
          color: const Color(0xFF059669),
          expanded: _insuranceExpanded,
          onToggle: (v) => setState(() => _insuranceExpanded = v),
          child: _member.insurances.isEmpty
              ? _buildEmptySection(
                  'No insurance records',
                  Icons.shield_outlined,
                  () => _showAddInsuranceModal(),
                )
              : Column(
                  children: _member.insurances
                      .map((ins) => _buildInsuranceCard(ins))
                      .toList(),
                ),
        ),
        _buildSection(
          title: 'Earnings',
          icon: Icons.trending_up_rounded,
          color: const Color(0xFFDC2626),
          expanded: _earningsExpanded,
          onToggle: (v) => setState(() => _earningsExpanded = v),
          child: _member.earnings.isEmpty
              ? _buildEmptySection(
                  'No earning records',
                  Icons.trending_up_outlined,
                  () => _showAddEarningsModal(),
                )
              : _buildEarningsCard(),
        ),
      ],
    );
  }

  Widget _buildSection({
    required String title,
    required IconData icon,
    required Color color,
    required bool expanded,
    required ValueChanged<bool> onToggle,
    Widget? child,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: expanded ? 0.55 : 0.30),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: expanded
              ? Colors.white.withValues(alpha: 0.5)
              : Colors.white.withValues(alpha: 0.2),
          width: 1,
        ),
      ),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 6),
            child: Row(
              children: [
                Icon(
                  icon,
                  color: expanded ? color : const Color(0xFF94A3B8),
                  size: 20,
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Text(
                    title,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w800,
                      color: expanded
                          ? const Color(0xFF0F172A)
                          : const Color(0xFF94A3B8),
                    ),
                  ),
                ),
                Switch(
                  value: expanded,
                  onChanged: onToggle,
                  activeThumbColor: color,
                  inactiveThumbColor: const Color(0xFFCBD5E1),
                  inactiveTrackColor: const Color(0xFFE2E8F0),
                  materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                ),
              ],
            ),
          ),
          if (expanded && child != null)
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 0, 16, 16),
              child: child,
            ),
        ],
      ),
    );
  }

  Widget _buildEmptySection(String text, IconData icon, VoidCallback onAdd) {
    return Column(
      children: [
        const SizedBox(height: 8),
        Icon(icon, size: 40, color: const Color(0xFF94A3B8)),
        const SizedBox(height: 8),
        Text(
          text,
          style: const TextStyle(color: Color(0xFF94A3B8), fontSize: 13),
        ),
        const SizedBox(height: 12),
        ElevatedButton.icon(
          onPressed: onAdd,
          icon: const Icon(Icons.add, size: 16),
          label: const Text('Add', style: TextStyle(fontSize: 13)),
          style: ElevatedButton.styleFrom(
            backgroundColor: const Color(0xFF6366F1),
            foregroundColor: Colors.white,
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
          ),
        ),
      ],
    );
  }

  // ─── SCHOOLING ──────────────────────────────────────────────────────────────

  Widget _buildSchoolingCard(dynamic s) {
    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: const Color(0xFF6366F1).withValues(alpha: 0.06),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  s.institutionName,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 14,
                    color: Color(0xFF0F172A),
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  '₹${s.feeAmount.toStringAsFixed(0)} / ${s.feeFrequency}',
                  style: const TextStyle(
                    color: Color(0xFF64748B),
                    fontSize: 12,
                  ),
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(
              Icons.delete_outline,
              color: Colors.redAccent,
              size: 20,
            ),
            onPressed: () async {
              await Provider.of<FinancialProvider>(
                context,
                listen: false,
              ).deleteSchooling(_member.id, s.id);
              if (mounted) UiUtils.showSnack(context, 'Schooling removed');
            },
          ),
        ],
      ),
    );
  }

  void _showAddSchoolingModal() {
    final instCtrl = TextEditingController();
    final feeCtrl = TextEditingController();
    String freq = 'monthly';
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24,
            right: 24,
            top: 24,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Add Schooling',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: instCtrl,
                decoration: const InputDecoration(
                  labelText: 'Institution Name',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: feeCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Fee Amount (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                initialValue: freq,
                decoration: const InputDecoration(
                  labelText: 'Frequency',
                  border: OutlineInputBorder(),
                ),
                items: ['monthly', 'quarterly', 'yearly']
                    .map((f) => DropdownMenuItem(value: f, child: Text(f)))
                    .toList(),
                onChanged: (v) {
                  if (v != null) setState(() => freq = v);
                },
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: () async {
                    if (instCtrl.text.isEmpty || feeCtrl.text.isEmpty) return;
                    await Provider.of<FinancialProvider>(
                      context,
                      listen: false,
                    ).addSchooling(_member.id, {
                      'institution_name': instCtrl.text.trim(),
                      'fee_amount': double.tryParse(feeCtrl.text.trim()) ?? 0,
                      'fee_frequency': freq,
                    });
                    if (ctx.mounted) Navigator.pop(ctx);
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366F1),
                  ),
                  child: const Text(
                    'Save',
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showFeeLedgerModal() async {
    final s = _member.schooling.first;
    final provider = Provider.of<FinancialProvider>(context, listen: false);
    List<dynamic> payments = [];
    try {
      final data = await provider.getSchoolingPayments(_member.id, s.id);
      if (data is List) payments = data;
    } catch (_) {}

    if (!mounted) return;
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: const EdgeInsets.fromLTRB(24, 8, 24, 32),
          constraints: BoxConstraints(
            maxHeight: MediaQuery.of(context).size.height * 0.65,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Center(
                child: Container(
                  width: 40,
                  height: 4,
                  margin: const EdgeInsets.only(bottom: 12),
                  decoration: BoxDecoration(
                    color: Colors.grey.shade300,
                    borderRadius: BorderRadius.circular(4),
                  ),
                ),
              ),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Text(
                    'Fee Ledger — ${s.institutionName}',
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  IconButton(
                    icon: const Icon(Icons.add_rounded),
                    onPressed: () => _showAddPaymentModal(s.id, () async {
                      final data = await provider.getSchoolingPayments(
                        _member.id,
                        s.id,
                      );
                      if (data is List) setState(() => payments = data);
                    }),
                  ),
                ],
              ),
              const Divider(),
              Expanded(
                child: payments.isEmpty
                    ? const Center(
                        child: Text(
                          'No payments recorded.',
                          style: TextStyle(color: Colors.grey),
                        ),
                      )
                    : ListView.builder(
                        itemCount: payments.length,
                        itemBuilder: (ctx, i) {
                          final p = payments[i];
                          return ListTile(
                            dense: true,
                            leading: const Icon(
                              Icons.payment_rounded,
                              color: Color(0xFF6366F1),
                              size: 20,
                            ),
                            title: Text(
                              '₹${(p['amount'] ?? 0).toDouble().toStringAsFixed(0)}',
                              style: const TextStyle(
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            subtitle: Text(
                              p['paid_date'] ?? '',
                              style: const TextStyle(fontSize: 12),
                            ),
                            trailing: IconButton(
                              icon: const Icon(
                                Icons.delete_outline,
                                color: Colors.red,
                                size: 18,
                              ),
                              onPressed: () async {
                                await provider.deleteSchoolingPayment(
                                  _member.id,
                                  s.id,
                                  p['id'],
                                );
                                final data = await provider
                                    .getSchoolingPayments(_member.id, s.id);
                                if (data is List)
                                  setState(() => payments = data);
                              },
                            ),
                          );
                        },
                      ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _showAddPaymentModal(int schoolingId, VoidCallback onDone) {
    final amtCtrl = TextEditingController();
    DateTime selectedDate = DateTime.now();
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24,
            right: 24,
            top: 24,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Record Payment',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: amtCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Amount (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 12),
              Row(
                children: [
                  const Text(
                    'Date: ',
                    style: TextStyle(fontWeight: FontWeight.w600),
                  ),
                  TextButton(
                    onPressed: () async {
                      final d = await showDatePicker(
                        context: ctx,
                        initialDate: selectedDate,
                        firstDate: DateTime(2000),
                        lastDate: DateTime.now(),
                      );
                      if (d != null) setState(() => selectedDate = d);
                    },
                    child: Text(
                      selectedDate.toIso8601String().split('T').first,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: () async {
                    final amt = double.tryParse(amtCtrl.text.trim());
                    if (amt == null || amt <= 0) return;
                    await Provider.of<FinancialProvider>(
                      context,
                      listen: false,
                    ).recordSchoolingPayment(_member.id, schoolingId, {
                      'amount': amt,
                      'paid_date': selectedDate
                          .toIso8601String()
                          .split('T')
                          .first,
                    });
                    if (ctx.mounted) Navigator.pop(ctx);
                    UiUtils.showSnack(context, 'Payment recorded');
                    onDone();
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366F1),
                  ),
                  child: const Text(
                    'Record',
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // ─── CHECKUPS ───────────────────────────────────────────────────────────────

  Widget _buildCheckupCard(dynamic c) {
    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: const Color(0xFFE88A1A).withValues(alpha: 0.06),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '${c.frequency} checkup',
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 14,
                    color: Color(0xFF0F172A),
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  'Cost: ₹${c.recurringCost.toStringAsFixed(0)} / ${c.frequency}',
                  style: const TextStyle(
                    color: Color(0xFF64748B),
                    fontSize: 12,
                  ),
                ),
                if (c.lastCheckupDate != null)
                  Text(
                    'Last: ${c.lastCheckupDate}',
                    style: const TextStyle(
                      color: Color(0xFF64748B),
                      fontSize: 11,
                    ),
                  ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(
              Icons.delete_outline,
              color: Colors.redAccent,
              size: 20,
            ),
            onPressed: () async {
              await Provider.of<FinancialProvider>(
                context,
                listen: false,
              ).deleteCheckup(_member.id, c.id);
              if (mounted) UiUtils.showSnack(context, 'Checkup removed');
            },
          ),
        ],
      ),
    );
  }

  void _showAddCheckupModal() {
    String freq = 'yearly';
    final costCtrl = TextEditingController();
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24,
            right: 24,
            top: 24,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Add Checkup',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              DropdownButtonFormField<String>(
                initialValue: freq,
                decoration: const InputDecoration(
                  labelText: 'Frequency',
                  border: OutlineInputBorder(),
                ),
                items: ['monthly', 'quarterly', 'yearly']
                    .map((f) => DropdownMenuItem(value: f, child: Text(f)))
                    .toList(),
                onChanged: (v) {
                  if (v != null) setState(() => freq = v);
                },
              ),
              const SizedBox(height: 12),
              TextField(
                controller: costCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Cost (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: () async {
                    await Provider.of<FinancialProvider>(
                      context,
                      listen: false,
                    ).addCheckup(_member.id, {
                      'frequency': freq,
                      'recurring_cost':
                          double.tryParse(costCtrl.text.trim()) ?? 0,
                    });
                    if (ctx.mounted) Navigator.pop(ctx);
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366F1),
                  ),
                  child: const Text(
                    'Save',
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // ─── MEDICINES ──────────────────────────────────────────────────────────────

  Widget _buildMedicineCard(dynamic m) {
    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: const Color(0xFF7C3AED).withValues(alpha: 0.06),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  m.medicineName,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 14,
                    color: Color(0xFF0F172A),
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  '₹${m.monthlyCost.toStringAsFixed(0)} / month',
                  style: const TextStyle(
                    color: Color(0xFF64748B),
                    fontSize: 12,
                  ),
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(
              Icons.delete_outline,
              color: Colors.redAccent,
              size: 20,
            ),
            onPressed: () async {
              await Provider.of<FinancialProvider>(
                context,
                listen: false,
              ).deleteMedicine(_member.id, m.id);
              if (mounted) UiUtils.showSnack(context, 'Medicine removed');
            },
          ),
        ],
      ),
    );
  }

  void _showAddMedicineModal() {
    final nameCtrl = TextEditingController();
    final costCtrl = TextEditingController();
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24,
            right: 24,
            top: 24,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Add Medicine',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: nameCtrl,
                decoration: const InputDecoration(
                  labelText: 'Medicine Name',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: costCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Monthly Cost (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: () async {
                    if (nameCtrl.text.isEmpty) return;
                    await Provider.of<FinancialProvider>(
                      context,
                      listen: false,
                    ).addMedicine(_member.id, {
                      'medicine_name': nameCtrl.text.trim(),
                      'monthly_cost':
                          double.tryParse(costCtrl.text.trim()) ?? 0,
                    });
                    if (ctx.mounted) Navigator.pop(ctx);
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366F1),
                  ),
                  child: const Text(
                    'Save',
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // ─── VACCINATIONS ───────────────────────────────────────────────────────────

  Widget _buildVaccinationCard(dynamic v) {
    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: const Color(0xFF2563EB).withValues(alpha: 0.06),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  v.vaccineName,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 14,
                    color: Color(0xFF0F172A),
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  '₹${v.recurringCost.toStringAsFixed(0)} / ${v.frequency}',
                  style: const TextStyle(
                    color: Color(0xFF64748B),
                    fontSize: 12,
                  ),
                ),
                if (v.lastVaccinationDate != null)
                  Text(
                    'Last: ${v.lastVaccinationDate}',
                    style: const TextStyle(
                      color: Color(0xFF64748B),
                      fontSize: 11,
                    ),
                  ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(
              Icons.delete_outline,
              color: Colors.redAccent,
              size: 20,
            ),
            onPressed: () async {
              await Provider.of<FinancialProvider>(
                context,
                listen: false,
              ).deleteVaccination(_member.id, v.id);
              if (mounted) UiUtils.showSnack(context, 'Vaccination removed');
            },
          ),
        ],
      ),
    );
  }

  void _showAddVaccinationModal() {
    final nameCtrl = TextEditingController();
    final costCtrl = TextEditingController();
    String freq = 'yearly';
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24,
            right: 24,
            top: 24,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Add Vaccination',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: nameCtrl,
                decoration: const InputDecoration(
                  labelText: 'Vaccine Name',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                initialValue: freq,
                decoration: const InputDecoration(
                  labelText: 'Frequency',
                  border: OutlineInputBorder(),
                ),
                items: ['monthly', 'quarterly', 'yearly']
                    .map((f) => DropdownMenuItem(value: f, child: Text(f)))
                    .toList(),
                onChanged: (v) {
                  if (v != null) setState(() => freq = v);
                },
              ),
              const SizedBox(height: 12),
              TextField(
                controller: costCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Cost (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: () async {
                    if (nameCtrl.text.isEmpty) return;
                    await Provider.of<FinancialProvider>(
                      context,
                      listen: false,
                    ).addVaccination(_member.id, {
                      'vaccine_name': nameCtrl.text.trim(),
                      'frequency': freq,
                      'recurring_cost':
                          double.tryParse(costCtrl.text.trim()) ?? 0,
                    });
                    if (ctx.mounted) Navigator.pop(ctx);
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366F1),
                  ),
                  child: const Text(
                    'Save',
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // ─── INSURANCE ──────────────────────────────────────────────────────────────

  Widget _buildInsuranceCard(dynamic ins) {
    return Container(
      margin: const EdgeInsets.only(bottom: 10),
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: const Color(0xFF059669).withValues(alpha: 0.06),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  ins.providerName ?? 'Insurance',
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 14,
                    color: Color(0xFF0F172A),
                  ),
                ),
                const SizedBox(height: 2),
                Text(
                  'Premium: ₹${ins.premiumAmount.toStringAsFixed(0)} / ${ins.premiumFrequency}',
                  style: const TextStyle(
                    color: Color(0xFF64748B),
                    fontSize: 12,
                  ),
                ),
                if (ins.coverageAmount > 0)
                  Text(
                    'Coverage: ₹${ins.coverageAmount.toStringAsFixed(0)}',
                    style: const TextStyle(
                      color: Color(0xFF64748B),
                      fontSize: 11,
                    ),
                  ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(
              Icons.delete_outline,
              color: Colors.redAccent,
              size: 20,
            ),
            onPressed: () async {
              await Provider.of<FinancialProvider>(
                context,
                listen: false,
              ).unlinkInsurance(_member.id, ins.id);
              if (mounted) UiUtils.showSnack(context, 'Insurance removed');
            },
          ),
        ],
      ),
    );
  }

  void _showAddInsuranceModal() {
    final providerCtrl = TextEditingController();
    final premiumCtrl = TextEditingController();
    final coverageCtrl = TextEditingController();
    String freq = 'yearly';
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24,
            right: 24,
            top: 24,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Add Insurance',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: providerCtrl,
                decoration: const InputDecoration(
                  labelText: 'Provider / Policy Name',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 12),
              DropdownButtonFormField<String>(
                initialValue: freq,
                decoration: const InputDecoration(
                  labelText: 'Premium Frequency',
                  border: OutlineInputBorder(),
                ),
                items: ['monthly', 'quarterly', 'yearly']
                    .map((f) => DropdownMenuItem(value: f, child: Text(f)))
                    .toList(),
                onChanged: (v) {
                  if (v != null) setState(() => freq = v);
                },
              ),
              const SizedBox(height: 12),
              TextField(
                controller: premiumCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Premium Amount (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: coverageCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Coverage Amount (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: () async {
                    await Provider.of<FinancialProvider>(
                      context,
                      listen: false,
                    ).linkInsurance(_member.id, {
                      'provider_name': providerCtrl.text.trim(),
                      'premium_frequency': freq,
                      'premium_amount':
                          double.tryParse(premiumCtrl.text.trim()) ?? 0,
                      'coverage_amount':
                          double.tryParse(coverageCtrl.text.trim()) ?? 0,
                    });
                    if (ctx.mounted) Navigator.pop(ctx);
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366F1),
                  ),
                  child: const Text(
                    'Save',
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  // ─── EARNINGS ───────────────────────────────────────────────────────────────

  Widget _buildEarningsCard() {
    final e = _member.earnings.isNotEmpty ? _member.earnings.first : null;
    if (e == null) return const SizedBox.shrink();
    return Container(
      padding: const EdgeInsets.all(12),
      decoration: BoxDecoration(
        color: const Color(0xFFDC2626).withValues(alpha: 0.06),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      '₹${e.monthlyIncome.toStringAsFixed(0)} / month',
                      style: const TextStyle(
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                        color: Color(0xFF0F172A),
                      ),
                    ),
                    const SizedBox(height: 2),
                    Text(
                      'Contribution: ₹${e.contributionToHousehold.toStringAsFixed(0)}',
                      style: const TextStyle(
                        color: Color(0xFF64748B),
                        fontSize: 12,
                      ),
                    ),
                  ],
                ),
              ),
              IconButton(
                icon: const Icon(
                  Icons.delete_outline,
                  color: Colors.redAccent,
                  size: 20,
                ),
                onPressed: () async {
                  await Provider.of<FinancialProvider>(
                    context,
                    listen: false,
                  ).deleteEarnings(_member.id, e.id);
                  if (mounted) UiUtils.showSnack(context, 'Earnings removed');
                },
              ),
            ],
          ),
        ],
      ),
    );
  }

  void _showAddEarningsModal() {
    final incomeCtrl = TextEditingController();
    final contribCtrl = TextEditingController();
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
          ),
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom + 32,
            left: 24,
            right: 24,
            top: 24,
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const Text(
                'Add Earnings',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: incomeCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Monthly Income (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 12),
              TextField(
                controller: contribCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Contribution to Household (₹)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 24),
              SizedBox(
                width: double.infinity,
                height: 48,
                child: ElevatedButton(
                  onPressed: () async {
                    await Provider.of<FinancialProvider>(
                      context,
                      listen: false,
                    ).addEarnings(_member.id, {
                      'monthly_income':
                          double.tryParse(incomeCtrl.text.trim()) ?? 0,
                      'contribution_to_household':
                          double.tryParse(contribCtrl.text.trim()) ?? 0,
                    });
                    if (ctx.mounted) Navigator.pop(ctx);
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: const Color(0xFF6366F1),
                  ),
                  child: const Text(
                    'Save',
                    style: TextStyle(color: Colors.white),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class _GridPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = const Color(0xFF6366F1).withValues(alpha: 0.035)
      ..strokeWidth = 0.5;
    const spacing = 40.0;
    for (double x = 0; x < size.width; x += spacing) {
      canvas.drawLine(Offset(x, 0), Offset(x, size.height), paint);
    }
    for (double y = 0; y < size.height; y += spacing) {
      canvas.drawLine(Offset(0, y), Offset(size.width, y), paint);
    }
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) => false;
}
