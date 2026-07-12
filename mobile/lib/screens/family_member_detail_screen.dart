import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/financial_models.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class FamilyMemberDetailScreen extends StatefulWidget {
  final LocalFamilyMember member;

  const FamilyMemberDetailScreen({super.key, required this.member});

  @override
  State<FamilyMemberDetailScreen> createState() => _FamilyMemberDetailScreenState();
}

class _FamilyMemberDetailScreenState extends State<FamilyMemberDetailScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 6, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  void _showDeleteMemberConfirmation() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Delete Profile'),
        content: Text('Are you sure you want to delete ${widget.member.name}? This will remove all their records permanently.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () async {
              try {
                await Provider.of<FinancialProvider>(context, listen: false)
                    .deleteFamilyMember(widget.member.id);
                if (ctx.mounted) Navigator.pop(ctx);
                if (mounted) {
                  Navigator.pop(context);
                  UiUtils.showSnack(context, 'Member removed successfully');
                }
              } catch (e) {
                if (ctx.mounted) UiUtils.showSnack(ctx, 'Failed: $e', isError: true);
              }
            },
            style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
            child: const Text('Delete', style: TextStyle(color: Colors.white)),
          )
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<FinancialProvider>();
    // Find the refreshed member data from the provider
    final refreshedMember = provider.familyMembers.firstWhere(
      (m) => m.id == widget.member.id,
      orElse: () => widget.member,
    );

    return Scaffold(
      backgroundColor: const Color(0xFFF1F5F9),
      appBar: AppBar(
        title: Text(refreshedMember.name),
        backgroundColor: const Color(0xFF042F2E),
        foregroundColor: Colors.white,
        actions: [
          if (!refreshedMember.isSelf)
            IconButton(
              icon: const Icon(Icons.delete_outline),
              onPressed: _showDeleteMemberConfirmation,
            )
        ],
        bottom: TabBar(
          controller: _tabController,
          isScrollable: true,
          labelColor: Colors.white,
          unselectedLabelColor: Colors.white60,
          indicatorColor: const Color(0xFF14B8A6),
          tabs: const [
            Tab(text: 'Schooling', icon: Icon(Icons.school_outlined)),
            Tab(text: 'Health Checkups', icon: Icon(Icons.favorite_border_rounded)),
            Tab(text: 'Medicines', icon: Icon(Icons.medication_outlined)),
            Tab(text: 'Vaccinations', icon: Icon(Icons.vaccines_outlined)),
            Tab(text: 'Earnings', icon: Icon(Icons.payments_outlined)),
            Tab(text: 'Insurance', icon: Icon(Icons.shield_outlined)),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _SchoolingTab(member: refreshedMember),
          _CheckupTab(member: refreshedMember),
          _MedicineTab(member: refreshedMember),
          _VaccinationTab(member: refreshedMember),
          _EarningsTab(member: refreshedMember),
          _InsuranceTab(member: refreshedMember),
        ],
      ),
    );
  }
}

// ── SCHOOLING TAB ────────────────────────────────────────────────────────────
class _SchoolingTab extends StatelessWidget {
  final LocalFamilyMember member;

  const _SchoolingTab({required this.member});

  void _showAddSchoolingModal(BuildContext context) {
    final instCtrl = TextEditingController();
    final feeCtrl = TextEditingController();
    String freq = 'monthly';
    final notesCtrl = TextEditingController();

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
            left: 24, right: 24, top: 24,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Add Schooling Detail', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                const SizedBox(height: 16),
                TextField(
                  controller: instCtrl,
                  decoration: const InputDecoration(labelText: 'Institution Name'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: feeCtrl,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(labelText: 'Fee Amount (₹)'),
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  value: freq,
                  decoration: const InputDecoration(labelText: 'Frequency'),
                  items: const [
                    DropdownMenuItem(value: 'monthly', child: Text('Monthly')),
                    DropdownMenuItem(value: 'quarterly', child: Text('Quarterly')),
                    DropdownMenuItem(value: 'yearly', child: Text('Yearly')),
                  ],
                  onChanged: (v) {
                    if (v != null) setState(() => freq = v);
                  },
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: notesCtrl,
                  decoration: const InputDecoration(labelText: 'Notes'),
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: ElevatedButton(
                    onPressed: () async {
                      if (instCtrl.text.isEmpty || feeCtrl.text.isEmpty) return;
                      final fee = double.tryParse(feeCtrl.text.trim()) ?? 0.0;
                      try {
                        await Provider.of<FinancialProvider>(context, listen: false)
                            .addSchooling(member.id, {
                          'institution_name': instCtrl.text.trim(),
                          'fee_amount': fee,
                          'fee_frequency': freq,
                          'notes': notesCtrl.text.trim(),
                        });
                        if (ctx.mounted) Navigator.pop(ctx);
                        UiUtils.showSnack(context, 'Schooling detail added');
                      } catch (e) {
                        UiUtils.showSnack(context, 'Failed to add: $e', isError: true);
                      }
                    },
                    style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
                    child: const Text('Save', style: TextStyle(color: Colors.white)),
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
    if (member.schooling.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.school_outlined, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            const Text('No schooling details added.', style: TextStyle(color: Colors.grey)),
            const SizedBox(height: 12),
            ElevatedButton.icon(
              onPressed: () => _showAddSchoolingModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Schooling', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: member.schooling.length + 1,
      itemBuilder: (context, index) {
        if (index == member.schooling.length) {
          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: ElevatedButton.icon(
              onPressed: () => _showAddSchoolingModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Schooling Details', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          );
        }

        final sch = member.schooling[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Expanded(
                      child: Text(sch.institutionName, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                    ),
                    IconButton(
                      icon: const Icon(Icons.delete_outline, color: Colors.red),
                      onPressed: () async {
                        UiUtils.showDeleteBottomSheet(context, sch.institutionName, () async {
                          await Provider.of<FinancialProvider>(context, listen: false)
                              .deleteSchooling(member.id, sch.id);
                        });
                      },
                    )
                  ],
                ),
                const SizedBox(height: 8),
                Text('Fee: ₹${sch.feeAmount.toStringAsFixed(2)} / ${sch.feeFrequency}'),
                Text('Monthly equivalent: ₹${sch.monthlyEquivalent.toStringAsFixed(2)}'),
                if (sch.notes != null && sch.notes!.isNotEmpty) Text('Notes: ${sch.notes}'),
              ],
            ),
          ),
        );
      },
    );
  }
}

// ── HEALTH CHECKUP TAB ────────────────────────────────────────────────────────
class _CheckupTab extends StatelessWidget {
  final LocalFamilyMember member;

  const _CheckupTab({required this.member});

  void _showAddCheckupModal(BuildContext context) {
    final typeCtrl = TextEditingController();
    final costCtrl = TextEditingController();
    String freq = 'yearly';
    final notesCtrl = TextEditingController();

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
            left: 24, right: 24, top: 24,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Add Health Checkup Detail', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                const SizedBox(height: 16),
                TextField(
                  controller: typeCtrl,
                  decoration: const InputDecoration(labelText: 'Checkup Type (e.g. Eye, Dental)'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: costCtrl,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(labelText: 'Cost (₹)'),
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  value: freq,
                  decoration: const InputDecoration(labelText: 'Frequency'),
                  items: const [
                    DropdownMenuItem(value: 'monthly', child: Text('Monthly')),
                    DropdownMenuItem(value: 'quarterly', child: Text('Quarterly')),
                    DropdownMenuItem(value: 'yearly', child: Text('Yearly')),
                  ],
                  onChanged: (v) {
                    if (v != null) setState(() => freq = v);
                  },
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: notesCtrl,
                  decoration: const InputDecoration(labelText: 'Notes'),
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: ElevatedButton(
                    onPressed: () async {
                      if (typeCtrl.text.isEmpty || costCtrl.text.isEmpty) return;
                      final cost = double.tryParse(costCtrl.text.trim()) ?? 0.0;
                      try {
                        await Provider.of<FinancialProvider>(context, listen: false)
                            .addCheckup(member.id, {
                          'checkup_type': typeCtrl.text.trim(),
                          'recurring_cost': cost,
                          'frequency': freq,
                          'notes': notesCtrl.text.trim(),
                        });
                        if (ctx.mounted) Navigator.pop(ctx);
                        UiUtils.showSnack(context, 'Checkup detail added');
                      } catch (e) {
                        UiUtils.showSnack(context, 'Failed to add: $e', isError: true);
                      }
                    },
                    style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
                    child: const Text('Save', style: TextStyle(color: Colors.white)),
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
    if (member.checkups.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.favorite_border_rounded, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            const Text('No checkups added.', style: TextStyle(color: Colors.grey)),
            const SizedBox(height: 12),
            ElevatedButton.icon(
              onPressed: () => _showAddCheckupModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Checkup', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: member.checkups.length + 1,
      itemBuilder: (context, index) {
        if (index == member.checkups.length) {
          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: ElevatedButton.icon(
              onPressed: () => _showAddCheckupModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Checkup Details', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          );
        }

        final check = member.checkups[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Expanded(
                      child: Text(check.checkupType ?? 'Routine Checkup', style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                    ),
                    IconButton(
                      icon: const Icon(Icons.delete_outline, color: Colors.red),
                      onPressed: () async {
                        UiUtils.showDeleteBottomSheet(context, check.checkupType ?? 'Checkup', () async {
                          await Provider.of<FinancialProvider>(context, listen: false)
                              .deleteCheckup(member.id, check.id);
                        });
                      },
                    )
                  ],
                ),
                const SizedBox(height: 8),
                Text('Cost: ₹${check.recurringCost.toStringAsFixed(2)} / ${check.frequency}'),
                Text('Monthly equivalent: ₹${check.monthlyEquivalent.toStringAsFixed(2)}'),
                if (check.notes != null && check.notes!.isNotEmpty) Text('Notes: ${check.notes}'),
              ],
            ),
          ),
        );
      },
    );
  }
}

// ── MEDICINES TAB ────────────────────────────────────────────────────────────
class _MedicineTab extends StatelessWidget {
  final LocalFamilyMember member;

  const _MedicineTab({required this.member});

  void _showAddMedicineModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final costCtrl = TextEditingController();
    final purposeCtrl = TextEditingController();
    final docCtrl = TextEditingController();

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
            left: 24, right: 24, top: 24,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Add Medicine Detail', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                const SizedBox(height: 16),
                TextField(
                  controller: nameCtrl,
                  decoration: const InputDecoration(labelText: 'Medicine Name'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: costCtrl,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(labelText: 'Monthly Cost (₹)'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: purposeCtrl,
                  decoration: const InputDecoration(labelText: 'Purpose'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: docCtrl,
                  decoration: const InputDecoration(labelText: 'Prescribed By'),
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: ElevatedButton(
                    onPressed: () async {
                      if (nameCtrl.text.isEmpty || costCtrl.text.isEmpty) return;
                      final cost = double.tryParse(costCtrl.text.trim()) ?? 0.0;
                      try {
                        await Provider.of<FinancialProvider>(context, listen: false)
                            .addMedicine(member.id, {
                          'medicine_name': nameCtrl.text.trim(),
                          'monthly_cost': cost,
                          'purpose': purposeCtrl.text.trim(),
                          'prescribed_by': docCtrl.text.trim(),
                        });
                        if (ctx.mounted) Navigator.pop(ctx);
                        UiUtils.showSnack(context, 'Medicine added');
                      } catch (e) {
                        UiUtils.showSnack(context, 'Failed to add: $e', isError: true);
                      }
                    },
                    style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
                    child: const Text('Save', style: TextStyle(color: Colors.white)),
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
    if (member.medicines.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.medication_outlined, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            const Text('No regular medicines added.', style: TextStyle(color: Colors.grey)),
            const SizedBox(height: 12),
            ElevatedButton.icon(
              onPressed: () => _showAddMedicineModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Medicine', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: member.medicines.length + 1,
      itemBuilder: (context, index) {
        if (index == member.medicines.length) {
          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: ElevatedButton.icon(
              onPressed: () => _showAddMedicineModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Medicine Details', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          );
        }

        final med = member.medicines[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Expanded(
                      child: Text(med.medicineName, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                    ),
                    IconButton(
                      icon: const Icon(Icons.delete_outline, color: Colors.red),
                      onPressed: () async {
                        UiUtils.showDeleteBottomSheet(context, med.medicineName, () async {
                          await Provider.of<FinancialProvider>(context, listen: false)
                              .deleteMedicine(member.id, med.id);
                        });
                      },
                    )
                  ],
                ),
                const SizedBox(height: 8),
                Text('Monthly Cost: ₹${med.monthlyCost.toStringAsFixed(2)}'),
                if (med.purpose != null && med.purpose!.isNotEmpty) Text('Purpose: ${med.purpose}'),
                if (med.prescribedBy != null && med.prescribedBy!.isNotEmpty) Text('Prescribed by: ${med.prescribedBy}'),
              ],
            ),
          ),
        );
      },
    );
  }
}

// ── VACCINATIONS TAB ──────────────────────────────────────────────────────────
class _VaccinationTab extends StatelessWidget {
  final LocalFamilyMember member;

  const _VaccinationTab({required this.member});

  void _showAddVacModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final costCtrl = TextEditingController();
    String freq = 'one_time';
    final notesCtrl = TextEditingController();

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
            left: 24, right: 24, top: 24,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Add Vaccination Record', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                const SizedBox(height: 16),
                TextField(
                  controller: nameCtrl,
                  decoration: const InputDecoration(labelText: 'Vaccine Name'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: costCtrl,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(labelText: 'Cost (₹)'),
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  value: freq,
                  decoration: const InputDecoration(labelText: 'Schedule/Frequency'),
                  items: const [
                    DropdownMenuItem(value: 'one_time', child: Text('One Time')),
                    DropdownMenuItem(value: 'monthly', child: Text('Monthly')),
                    DropdownMenuItem(value: 'quarterly', child: Text('Quarterly')),
                    DropdownMenuItem(value: 'yearly', child: Text('Yearly')),
                  ],
                  onChanged: (v) {
                    if (v != null) setState(() => freq = v);
                  },
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: notesCtrl,
                  decoration: const InputDecoration(labelText: 'Notes'),
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: ElevatedButton(
                    onPressed: () async {
                      if (nameCtrl.text.isEmpty || costCtrl.text.isEmpty) return;
                      final cost = double.tryParse(costCtrl.text.trim()) ?? 0.0;
                      try {
                        await Provider.of<FinancialProvider>(context, listen: false)
                            .addVaccination(member.id, {
                          'vaccine_name': nameCtrl.text.trim(),
                          'recurring_cost': cost,
                          'frequency': freq,
                          'notes': notesCtrl.text.trim(),
                        });
                        if (ctx.mounted) Navigator.pop(ctx);
                        UiUtils.showSnack(context, 'Vaccination recorded');
                      } catch (e) {
                        UiUtils.showSnack(context, 'Failed to add: $e', isError: true);
                      }
                    },
                    style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
                    child: const Text('Save', style: TextStyle(color: Colors.white)),
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
    if (member.vaccinations.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.vaccines_outlined, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            const Text('No vaccination records added.', style: TextStyle(color: Colors.grey)),
            const SizedBox(height: 12),
            ElevatedButton.icon(
              onPressed: () => _showAddVacModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Vaccine', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: member.vaccinations.length + 1,
      itemBuilder: (context, index) {
        if (index == member.vaccinations.length) {
          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: ElevatedButton.icon(
              onPressed: () => _showAddVacModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Vaccine Details', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          );
        }

        final vac = member.vaccinations[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Expanded(
                      child: Text(vac.vaccineName, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                    ),
                    IconButton(
                      icon: const Icon(Icons.delete_outline, color: Colors.red),
                      onPressed: () async {
                        UiUtils.showDeleteBottomSheet(context, vac.vaccineName, () async {
                          await Provider.of<FinancialProvider>(context, listen: false)
                              .deleteVaccination(member.id, vac.id);
                        });
                      },
                    )
                  ],
                ),
                const SizedBox(height: 8),
                Text('Cost: ₹${vac.recurringCost.toStringAsFixed(2)} (${vac.frequency})'),
                if (vac.frequency != 'one_time') Text('Monthly equivalent: ₹${vac.monthlyEquivalent.toStringAsFixed(2)}'),
                if (vac.notes != null && vac.notes!.isNotEmpty) Text('Notes: ${vac.notes}'),
              ],
            ),
          ),
        );
      },
    );
  }
}

// ── EARNINGS TAB ─────────────────────────────────────────────────────────────
class _EarningsTab extends StatelessWidget {
  final LocalFamilyMember member;

  const _EarningsTab({required this.member});

  void _showAddEarningsModal(BuildContext context) {
    final incCtrl = TextEditingController();
    final contribCtrl = TextEditingController();
    final occCtrl = TextEditingController();

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
            left: 24, right: 24, top: 24,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Add Earnings / Income', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                const SizedBox(height: 16),
                TextField(
                  controller: incCtrl,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(labelText: 'Monthly Income (₹)'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: contribCtrl,
                  keyboardType: TextInputType.number,
                  decoration: const InputDecoration(labelText: 'Monthly Contribution to Household (₹)'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: occCtrl,
                  decoration: const InputDecoration(labelText: 'Occupation (e.g. Engineer, Business)'),
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: ElevatedButton(
                    onPressed: () async {
                      if (incCtrl.text.isEmpty) return;
                      final inc = double.tryParse(incCtrl.text.trim()) ?? 0.0;
                      final contrib = double.tryParse(contribCtrl.text.trim()) ?? 0.0;
                      try {
                        await Provider.of<FinancialProvider>(context, listen: false)
                            .addEarnings(member.id, {
                          'monthly_income': inc,
                          'contribution_to_household': contrib,
                          'occupation': occCtrl.text.trim(),
                        });
                        if (ctx.mounted) Navigator.pop(ctx);
                        UiUtils.showSnack(context, 'Earnings added');
                      } catch (e) {
                        UiUtils.showSnack(context, 'Failed to add: $e', isError: true);
                      }
                    },
                    style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
                    child: const Text('Save', style: TextStyle(color: Colors.white)),
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
    if (member.earnings.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.payments_outlined, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            const Text('No earnings details added.', style: TextStyle(color: Colors.grey)),
            const SizedBox(height: 12),
            ElevatedButton.icon(
              onPressed: () => _showAddEarningsModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Earnings', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: member.earnings.length + 1,
      itemBuilder: (context, index) {
        if (index == member.earnings.length) {
          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: ElevatedButton.icon(
              onPressed: () => _showAddEarningsModal(context),
              icon: const Icon(Icons.add, color: Colors.white),
              label: const Text('Add Earnings Details', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          );
        }

        final earn = member.earnings[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Expanded(
                      child: Text(earn.occupation ?? 'Salary/Income', style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                    ),
                    IconButton(
                      icon: const Icon(Icons.delete_outline, color: Colors.red),
                      onPressed: () async {
                        UiUtils.showDeleteBottomSheet(context, earn.occupation ?? 'Earnings', () async {
                          await Provider.of<FinancialProvider>(context, listen: false)
                              .deleteEarnings(member.id, earn.id);
                        });
                      },
                    )
                  ],
                ),
                const SizedBox(height: 8),
                Text('Monthly Income: ₹${earn.monthlyIncome.toStringAsFixed(2)}'),
                Text('Contribution to Household: ₹${earn.contributionToHousehold.toStringAsFixed(2)}'),
              ],
            ),
          ),
        );
      },
    );
  }
}

// ── INSURANCE TAB ────────────────────────────────────────────────────────────
class _InsuranceTab extends StatelessWidget {
  final LocalFamilyMember member;

  const _InsuranceTab({required this.member});

  void _showLinkInsuranceModal(BuildContext context) {
    final provider = Provider.of<FinancialProvider>(context, listen: false);
    // Find insurances that are NOT already linked to this member
    final unlinkedInsurances = provider.insurances.where((ins) {
      return !member.insuranceLinks.any((link) => link.insuranceId == ins['id']);
    }).toList();

    if (unlinkedInsurances.isEmpty) {
      showDialog(
        context: context,
        builder: (ctx) => AlertDialog(
          title: const Text('No Unlinked Insurance'),
          content: const Text('All existing insurance policies are already linked to this member, or you have not added policies in the Insurance Hub yet.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(ctx),
              child: const Text('OK'),
            )
          ],
        ),
      );
      return;
    }

    dynamic selectedIns = unlinkedInsurances.first;
    String rel = 'self';

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
            left: 24, right: 24, top: 24,
          ),
          child: SingleChildScrollView(
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                const Text('Link Insurance Policy', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                const SizedBox(height: 16),
                DropdownButtonFormField<dynamic>(
                  value: selectedIns,
                  decoration: const InputDecoration(labelText: 'Select Policy'),
                  items: unlinkedInsurances.map((ins) {
                    final label = '${ins['provider']} (${ins['type']})';
                    return DropdownMenuItem(value: ins, child: Text(label));
                  }).toList(),
                  onChanged: (v) {
                    if (v != null) setState(() => selectedIns = v);
                  },
                ),
                const SizedBox(height: 12),
                DropdownButtonFormField<String>(
                  value: rel,
                  decoration: const InputDecoration(labelText: 'Relationship'),
                  items: const [
                    DropdownMenuItem(value: 'self', child: Text('Self Covered')),
                    DropdownMenuItem(value: 'dependent', child: Text('Dependent')),
                    DropdownMenuItem(value: 'nominee', child: Text('Nominee')),
                  ],
                  onChanged: (v) {
                    if (v != null) setState(() => rel = v);
                  },
                ),
                const SizedBox(height: 24),
                SizedBox(
                  width: double.infinity,
                  height: 48,
                  child: ElevatedButton(
                    onPressed: () async {
                      try {
                        await Provider.of<FinancialProvider>(context, listen: false)
                            .linkInsurance(member.id, {
                          'insurance_id': selectedIns['id'],
                          'relationship': rel,
                        });
                        if (ctx.mounted) Navigator.pop(ctx);
                        UiUtils.showSnack(context, 'Policy linked successfully');
                      } catch (e) {
                        UiUtils.showSnack(context, 'Failed to link: $e', isError: true);
                      }
                    },
                    style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
                    child: const Text('Link', style: TextStyle(color: Colors.white)),
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
    if (member.insuranceLinks.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Icon(Icons.shield_outlined, size: 64, color: Colors.grey),
            const SizedBox(height: 16),
            const Text('No insurance policies linked.', style: TextStyle(color: Colors.grey)),
            const SizedBox(height: 12),
            ElevatedButton.icon(
              onPressed: () => _showLinkInsuranceModal(context),
              icon: const Icon(Icons.link, color: Colors.white),
              label: const Text('Link Policy', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: member.insuranceLinks.length + 1,
      itemBuilder: (context, index) {
        if (index == member.insuranceLinks.length) {
          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 16),
            child: ElevatedButton.icon(
              onPressed: () => _showLinkInsuranceModal(context),
              icon: const Icon(Icons.link, color: Colors.white),
              label: const Text('Link Insurance Policy', style: TextStyle(color: Colors.white)),
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
            ),
          );
        }

        final link = member.insuranceLinks[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Expanded(
                      child: Text(link.provider ?? 'Linked Policy', style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold)),
                    ),
                    IconButton(
                      icon: const Icon(Icons.link_off, color: Colors.red),
                      onPressed: () async {
                        UiUtils.showDeleteBottomSheet(context, 'Unlink ${link.provider ?? "policy"}?', () async {
                          await Provider.of<FinancialProvider>(context, listen: false)
                              .unlinkInsurance(member.id, link.insuranceId);
                        });
                      },
                    )
                  ],
                ),
                const SizedBox(height: 8),
                Text('Plan: ${link.policyName ?? "N/A"}'),
                Text('Type: ${link.type ?? "N/A"}'),
                Text('Role: ${link.relationship[0].toUpperCase()}${link.relationship.substring(1)}'),
              ],
            ),
          ),
        );
      },
    );
  }
}
