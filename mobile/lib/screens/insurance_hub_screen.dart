import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class InsuranceHubScreen extends StatefulWidget {
  const InsuranceHubScreen({super.key});

  @override
  State<InsuranceHubScreen> createState() => _InsuranceHubScreenState();
}

class _InsuranceHubScreenState extends State<InsuranceHubScreen> {
  @override
  Widget build(BuildContext context) {
    final provider = context.watch<FinancialProvider>();
    final insurances = provider.insurances;

    return Scaffold(
      backgroundColor: const Color(0xFFF8F9FA),
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios_new_rounded, color: Color(0xFF1E293B)),
          onPressed: () => Navigator.pop(context),
        ),
        title: const Text('Insurance Hub', style: TextStyle(color: Color(0xFF1E293B), fontWeight: FontWeight.bold)),
        centerTitle: true,
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Educational Banner
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                gradient: const LinearGradient(colors: [Color(0xFF2DD4BF), Color(0xFF14B8A6)]),
                borderRadius: BorderRadius.circular(20),
              ),
              child: Row(
                children: [
                  const Icon(Icons.shield_rounded, color: Colors.white, size: 40),
                  const SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: const [
                        Text('Protect Your Future', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
                        SizedBox(height: 4),
                        Text('Always prefer pure Term Insurance over endowment (LIC) plans. Mix Health Insurance to cover medical emergencies.',
                          style: TextStyle(color: Colors.white70, fontSize: 13, height: 1.4)),
                      ],
                    ),
                  )
                ],
              ),
            ),
            const SizedBox(height: 24),
            
            const Text('Your Policies', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
            const SizedBox(height: 16),
            if (insurances.isEmpty)
              UiUtils.buildEmptyState(
                'No policies',
                'Tap Add Policy to track your health and term insurances.',
                Icons.health_and_safety_outlined,
                Colors.blue,
              )
            else
              ...insurances.map((ins) => _buildInsuranceCard(ins, provider)),
              
            const SizedBox(height: 24),
            SizedBox(
              width: double.infinity,
              height: 56,
              child: ElevatedButton.icon(
                onPressed: () => _showAddInsuranceModal(context, provider),
                icon: const Icon(Icons.add_rounded),
                label: const Text('Add Policy', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF6366F1),
                  foregroundColor: Colors.white,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInsuranceCard(dynamic ins, FinancialProvider p) {
    IconData icon;
    Color color;
    switch (ins['type']) {
      case 'health': icon = Icons.medical_services_rounded; color = const Color(0xFF10B981); break;
      case 'term': icon = Icons.shield_rounded; color = const Color(0xFF2DD4BF); break;
      case 'vehicle': icon = Icons.directions_car_rounded; color = const Color(0xFFF59E0B); break;
      default: icon = Icons.description_rounded; color = const Color(0xFF6366F1);
    }
    
    return Container(
      margin: const EdgeInsets.only(bottom: 16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: Colors.grey.shade100),
        boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.02), blurRadius: 10, offset: const Offset(0, 4))],
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(color: color.withValues(alpha: 0.1), shape: BoxShape.circle),
            child: Icon(icon, color: color, size: 24),
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(ins['provider'] ?? 'Provider', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 16, color: Color(0xFF1E293B))),
                    IconButton(
                      icon: const Icon(Icons.delete_outline_rounded, color: Colors.red, size: 20),
                      padding: EdgeInsets.zero,
                      constraints: const BoxConstraints(),
                      onPressed: () {
                        UiUtils.showDeleteBottomSheet(context, ins['provider'] ?? 'Policy', () async {
                          await p.deleteInsurance(ins['id'].toString());
                          if (mounted) UiUtils.showSnack(context, 'Policy removed.');
                        });
                      },
                    )
                  ],
                ),
                if (ins['policy_name'] != null && ins['policy_name'].isNotEmpty)
                  Text(ins['policy_name'], style: TextStyle(color: Colors.grey.shade600, fontSize: 13)),
                const SizedBox(height: 12),
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        const Text('Premium', style: TextStyle(color: Colors.grey, fontSize: 11)),
                        Text('₹${(ins['premium_amount'] ?? 0).toStringAsFixed(0)} / ${ins['premium_frequency']}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                      ],
                    ),
                    if (ins['coverage_amount'] != null)
                      Column(
                        crossAxisAlignment: CrossAxisAlignment.end,
                        children: [
                          const Text('Coverage', style: TextStyle(color: Colors.grey, fontSize: 11)),
                          Text('₹${(ins['coverage_amount'] ?? 0).toStringAsFixed(0)}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Color(0xFF10B981))),
                        ],
                      ),
                  ],
                )
              ],
            ),
          )
        ],
      ),
    );
  }

  void _showAddInsuranceModal(BuildContext context, FinancialProvider p) {
    String type = 'health';
    String frequency = 'yearly';
    final providerCtrl = TextEditingController();
    final policyCtrl = TextEditingController();
    final premiumCtrl = TextEditingController();
    final coverageCtrl = TextEditingController();

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
          child: Container(
            decoration: const BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
            ),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Center(
                    child: Container(
                      width: 40, height: 4,
                      margin: const EdgeInsets.only(top: 12, bottom: 20),
                      decoration: BoxDecoration(color: Colors.grey.shade300, borderRadius: BorderRadius.circular(4)),
                    ),
                  ),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Add Policy', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
                      IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(ctx)),
                    ],
                  ),
                  const SizedBox(height: 16),
                  DropdownButtonFormField<String>(
                    initialValue: type,
                    decoration: InputDecoration(labelText: 'Insurance Type', filled: true, border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none)),
                    items: const [
                      DropdownMenuItem(value: 'health', child: Text('Health Insurance')),
                      DropdownMenuItem(value: 'term', child: Text('Term Life Insurance')),
                      DropdownMenuItem(value: 'vehicle', child: Text('Vehicle Insurance')),
                      DropdownMenuItem(value: 'other', child: Text('Other')),
                    ],
                    onChanged: (v) => setState(() => type = v!),
                  ),
                  const SizedBox(height: 12),
                  TextField(controller: providerCtrl, decoration: InputDecoration(labelText: 'Provider (e.g. HDFC Ergo)', filled: true, border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none))),
                  const SizedBox(height: 12),
                  TextField(controller: policyCtrl, decoration: InputDecoration(labelText: 'Policy Name (Optional)', filled: true, border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none))),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(child: TextField(controller: premiumCtrl, decoration: InputDecoration(labelText: 'Premium (₹)', filled: true, border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none)), keyboardType: TextInputType.number)),
                      const SizedBox(width: 12),
                      Expanded(
                        child: DropdownButtonFormField<String>(
                          initialValue: frequency,
                          decoration: InputDecoration(labelText: 'Frequency', filled: true, border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none)),
                          items: ['monthly', 'quarterly', 'yearly'].map((e) => DropdownMenuItem(value: e, child: Text(e))).toList(),
                          onChanged: (v) => setState(() => frequency = v!),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  TextField(controller: coverageCtrl, decoration: InputDecoration(labelText: 'Coverage Amount (Optional)', filled: true, border: OutlineInputBorder(borderRadius: BorderRadius.circular(12), borderSide: BorderSide.none)), keyboardType: TextInputType.number),
                  const SizedBox(height: 24),
                  SizedBox(
                    width: double.infinity,
                    height: 56,
                    child: ElevatedButton(
                      onPressed: () async {
                        if (providerCtrl.text.isEmpty || premiumCtrl.text.isEmpty) {
                          UiUtils.showSnack(ctx, 'Please fill required fields.');
                          return;
                        }
                        try {
                          await p.addInsurance({
                            'type': type,
                            'provider': providerCtrl.text,
                            'policy_name': policyCtrl.text,
                            'premium_amount': double.tryParse(premiumCtrl.text) ?? 0,
                            'premium_frequency': frequency,
                            'coverage_amount': double.tryParse(coverageCtrl.text),
                          });
                          if (context.mounted) {
                            Navigator.pop(ctx);
                            UiUtils.showSnack(context, 'Policy added successfully!');
                          }
                        } catch (e) {
                          if (context.mounted) UiUtils.showSnack(context, 'Failed to add policy: $e', isError: true);
                        }
                      },
                      style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF6366F1), foregroundColor: Colors.white, shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16))),
                      child: const Text('Save Policy', style: TextStyle(fontWeight: FontWeight.bold, fontSize: 16)),
                    ),
                  )
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
