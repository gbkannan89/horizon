import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../models/financial_models.dart';
import '../utils/ui_utils.dart';

class DeviceDetailScreen extends StatefulWidget {
  final LocalElectronic device;

  const DeviceDetailScreen({
    super.key,
    required this.device,
  });

  @override
  State<DeviceDetailScreen> createState() => _DeviceDetailScreenState();
}

class _DeviceDetailScreenState extends State<DeviceDetailScreen> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _fetchData();
    });
  }

  void _fetchData() {
    final provider = Provider.of<FinancialProvider>(context, listen: false);
    provider.reloadElectronicServices(widget.device.id);
    provider.reloadElectronicEmi(widget.device.id);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  IconData _getCategoryIcon(String category) {
    switch (category.toLowerCase()) {
      case 'mobile':
        return Icons.phone_android;
      case 'laptop':
        return Icons.laptop;
      case 'tablet':
        return Icons.tablet_android;
      case 'tv':
        return Icons.tv;
      case 'appliance':
        return Icons.kitchen;
      default:
        return Icons.devices_other;
    }
  }

  Color _getWarrantyColor(String status) {
    switch (status.toLowerCase()) {
      case 'active':
        return Colors.green;
      case 'near_expiry':
        return Colors.amber;
      default:
        return Colors.redAccent;
    }
  }

  Widget _buildOverviewTab(LocalElectronic d, LocalElectronicEmi? emi) {
    final ageDays = DateTime.now().difference(d.purchaseDate).inDays;
    final ageYears = ageDays / 365.25;
    final depPct = (ageYears / d.expectedLifeYears * 100).clamp(0.0, 100.0);

    return RefreshIndicator(
      onRefresh: () async => _fetchData(),
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Specs
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.03), blurRadius: 10)],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.all(12),
                        decoration: BoxDecoration(color: const Color(0xFF0D9488).withValues(alpha: 0.08), shape: BoxShape.circle),
                        child: Icon(_getCategoryIcon(d.category), color: const Color(0xFF0D9488), size: 28),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(d.name, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: Color(0xFF1E293B))),
                            if (d.brand != null || d.model != null) ...[
                              const SizedBox(height: 4),
                              Text('${d.brand ?? ''} ${d.model ?? ''}'.trim(), style: TextStyle(color: Colors.grey.shade600, fontSize: 13, fontWeight: FontWeight.bold)),
                            ]
                          ],
                        ),
                      )
                    ],
                  ),
                  const Divider(height: 30),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      _buildSpecItem('Purchase Date', '${d.purchaseDate.day}/${d.purchaseDate.month}/${d.purchaseDate.year}'),
                      _buildSpecItem('Original Cost', '₹${d.purchaseAmount.toStringAsFixed(0)}'),
                      _buildSpecItem('Life Expectancy', '${d.expectedLifeYears} Years'),
                    ],
                  )
                ],
              ),
            ),
            const SizedBox(height: 20),

            // Depreciation & Value Card
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.03), blurRadius: 10)],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Current Value & Depreciation', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                  const SizedBox(height: 16),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Purchase Price', style: TextStyle(color: Colors.black87)),
                      Text('₹${d.purchaseAmount.toStringAsFixed(0)}', style: const TextStyle(fontWeight: FontWeight.bold)),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text('Depreciation (${depPct.toStringAsFixed(0)}%)', style: const TextStyle(color: Colors.redAccent)),
                      Text('-₹${(d.purchaseAmount - d.currentValue).toStringAsFixed(0)}', style: const TextStyle(color: Colors.redAccent, fontWeight: FontWeight.bold)),
                    ],
                  ),
                  const Divider(height: 24),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Estimated Current Value', style: TextStyle(color: Colors.black87, fontWeight: FontWeight.bold)),
                      Text('₹${d.currentValue.toStringAsFixed(0)}', style: const TextStyle(fontWeight: FontWeight.w900, color: Color(0xFF0D9488), fontSize: 18)),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),

            // Warranty Card
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.03), blurRadius: 10)],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text('Warranty Status', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
                  const SizedBox(height: 16),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Warranty Period', style: TextStyle(color: Colors.black87)),
                      Text('${d.warrantyYears} Year${d.warrantyYears > 1 ? 's' : ''}', style: const TextStyle(fontWeight: FontWeight.bold)),
                    ],
                  ),
                  if (d.warrantyExpiryDate != null) ...[
                    const SizedBox(height: 8),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text('Expiry Date', style: TextStyle(color: Colors.black87)),
                        Text('${d.warrantyExpiryDate!.day}/${d.warrantyExpiryDate!.month}/${d.warrantyExpiryDate!.year}', style: const TextStyle(fontWeight: FontWeight.bold)),
                      ],
                    ),
                  ],
                  const Divider(height: 24),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text('Status', style: TextStyle(color: Colors.black87)),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 4),
                        decoration: BoxDecoration(
                          color: _getWarrantyColor(d.warrantyStatus).withValues(alpha: 0.1),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Text(
                          d.warrantyStatus == 'expired'
                              ? 'EXPIRED'
                              : d.warrantyStatus == 'near_expiry'
                                  ? 'EXPIRING SOON'
                                  : 'ACTIVE',
                          style: TextStyle(
                            color: _getWarrantyColor(d.warrantyStatus),
                            fontWeight: FontWeight.bold,
                            fontSize: 11,
                          ),
                        ),
                      ),
                    ],
                  ),
                  if (d.warrantyStatus != 'expired') ...[
                    const SizedBox(height: 12),
                    Text(
                      '${d.warrantyDaysRemaining} days remaining in coverage.',
                      style: TextStyle(fontSize: 12, color: Colors.grey.shade600, fontWeight: FontWeight.bold),
                    ),
                  ],
                ],
              ),
            ),
            const SizedBox(height: 40),
          ],
        ),
      ),
    );
  }

  Widget _buildSpecItem(String label, String value) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(label, style: const TextStyle(fontSize: 11, color: Colors.grey)),
        const SizedBox(height: 4),
        Text(value, style: const TextStyle(fontSize: 14, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
      ],
    );
  }

  Widget _buildServiceTab(List<LocalElectronicService> services) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8F9FA),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _showAddServiceDialog(),
        backgroundColor: const Color(0xFF0D9488),
        child: const Icon(Icons.add, color: Colors.white),
      ),
      body: services.isEmpty
          ? const Center(child: Text('No service events recorded.'))
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: services.length,
              itemBuilder: (context, idx) {
                final s = services[idx];
                return Container(
                  margin: const EdgeInsets.only(bottom: 12),
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(16),
                    boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.02), blurRadius: 8)],
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            s.serviceType.toUpperCase().replaceAll('_', ' '),
                            style: const TextStyle(fontWeight: FontWeight.bold, color: Color(0xFF0D9488)),
                          ),
                          Text('₹${s.cost.toStringAsFixed(0)}', style: const TextStyle(fontWeight: FontWeight.w900, fontSize: 16)),
                        ],
                      ),
                      const SizedBox(height: 6),
                      Text(s.serviceDate, style: const TextStyle(color: Colors.grey, fontSize: 12)),
                      if (s.serviceCenter != null && s.serviceCenter!.isNotEmpty) ...[
                        const SizedBox(height: 8),
                        Text('Center: ${s.serviceCenter}', style: const TextStyle(fontSize: 12, fontWeight: FontWeight.bold)),
                      ],
                      if (s.description != null && s.description!.isNotEmpty) ...[
                        const SizedBox(height: 8),
                        Text(s.description!, style: TextStyle(color: Colors.grey.shade700, fontSize: 13)),
                      ],
                      const Divider(height: 24),
                      Align(
                        alignment: Alignment.centerRight,
                        child: GestureDetector(
                          onTap: () async {
                            UiUtils.showDeleteBottomSheet(context, 'Service Event', () async {
                              await Provider.of<FinancialProvider>(context, listen: false).deleteElectronicService(widget.device.id, s.id);
                              _fetchData();
                            });
                          },
                          child: const Text('Delete', style: TextStyle(color: Colors.redAccent, fontWeight: FontWeight.bold, fontSize: 12)),
                        ),
                      ),
                    ],
                  ),
                );
              },
            ),
    );
  }

  void _showAddServiceDialog() {
    final dateCtrl = TextEditingController(text: DateTime.now().toIso8601String().substring(0, 10));
    final costCtrl = TextEditingController();
    final centerCtrl = TextEditingController();
    final descCtrl = TextEditingController();
    String serviceType = 'repair';

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Add Service Record'),
        content: StatefulBuilder(
          builder: (ctx, setDialogState) => SingleChildScrollView(
            child: Column(
              children: [
                TextField(
                  controller: dateCtrl,
                  readOnly: true,
                  decoration: const InputDecoration(labelText: 'Service Date', suffixIcon: Icon(Icons.calendar_today)),
                  onTap: () async {
                    final d = await showDatePicker(context: ctx, initialDate: DateTime.now(), firstDate: DateTime(2000), lastDate: DateTime(2100));
                    if (d != null) {
                      setDialogState(() => dateCtrl.text = d.toIso8601String().substring(0, 10));
                    }
                  },
                ),
                DropdownButtonFormField<String>(
                  value: serviceType,
                  decoration: const InputDecoration(labelText: 'Type'),
                  items: const [
                    DropdownMenuItem(value: 'screen_replacement', child: Text('Screen Replacement')),
                    DropdownMenuItem(value: 'battery', child: Text('Battery Service')),
                    DropdownMenuItem(value: 'repair', child: Text('General Repair')),
                    DropdownMenuItem(value: 'annual_service', child: Text('Annual Service')),
                  ],
                  onChanged: (val) {
                    if (val != null) setDialogState(() => serviceType = val);
                  },
                ),
                TextField(controller: costCtrl, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'Cost (₹)')),
                TextField(controller: centerCtrl, decoration: const InputDecoration(labelText: 'Service Center')),
                TextField(controller: descCtrl, maxLines: 2, decoration: const InputDecoration(labelText: 'Description')),
              ],
            ),
          ),
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cancel')),
          TextButton(
            onPressed: () async {
              final body = {
                'service_date': dateCtrl.text,
                'service_type': serviceType,
                'cost': double.tryParse(costCtrl.text.trim()) ?? 0.0,
                'service_center': centerCtrl.text.trim().isEmpty ? null : centerCtrl.text.trim(),
                'description': descCtrl.text.trim().isEmpty ? null : descCtrl.text.trim(),
              };
              await Provider.of<FinancialProvider>(context, listen: false).addElectronicService(widget.device.id, body);
              if (mounted) {
                Navigator.pop(ctx);
                _fetchData();
              }
            },
            child: const Text('Save'),
          ),
        ],
      ),
    );
  }

  Widget _buildEmiTab(LocalElectronicEmi? emi) {
    if (emi == null) {
      return Padding(
        padding: const EdgeInsets.all(24),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(Icons.credit_card_outlined, size: 56, color: Colors.grey),
              const SizedBox(height: 12),
              const Text('No active EMI on this device.', style: TextStyle(fontWeight: FontWeight.bold)),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => _showEmiDialog(null),
                style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488)),
                child: const Text('Add EMI Schedule', style: TextStyle(color: Colors.white)),
              ),
            ],
          ),
        ),
      );
    }

    final remainingMonths = emi.totalMonths - emi.monthsPaid;
    final totalRemainingValue = emi.emiAmount * remainingMonths;

    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: Colors.white,
              borderRadius: BorderRadius.circular(20),
              boxShadow: [BoxShadow(color: Colors.black.withValues(alpha: 0.03), blurRadius: 10)],
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(emi.bankName, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
                    IconButton(icon: const Icon(Icons.edit_note, color: Colors.blue), onPressed: () => _showEmiDialog(emi)),
                  ],
                ),
                const Divider(),
                _buildAdvisoryRow('Monthly EMI', '₹${emi.emiAmount.toStringAsFixed(0)}'),
                const SizedBox(height: 12),
                _buildAdvisoryRow('EMIs Paid', '${emi.monthsPaid} / ${emi.totalMonths}'),
                const SizedBox(height: 12),
                _buildAdvisoryRow('Interest Rate', '${emi.interestRate}%'),
                const SizedBox(height: 12),
                _buildAdvisoryRow('Remaining Duration', '$remainingMonths Months'),
                const SizedBox(height: 12),
                _buildAdvisoryRow('Remaining Liability', '₹${totalRemainingValue.toStringAsFixed(0)}'),
              ],
            ),
          ),
          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: () {
                UiUtils.showDeleteBottomSheet(context, 'EMI Schedule', () async {
                  await Provider.of<FinancialProvider>(context, listen: false).deleteElectronicEmi(widget.device.id);
                  _fetchData();
                });
              },
              icon: const Icon(Icons.delete_outline, color: Colors.red),
              label: const Text('Delete EMI', style: TextStyle(color: Colors.red)),
            ),
          )
        ],
      ),
    );
  }

  void _showEmiDialog(LocalElectronicEmi? emi) {
    final bankCtrl = TextEditingController(text: emi?.bankName);
    final amtCtrl = TextEditingController(text: emi?.emiAmount.toStringAsFixed(0));
    final rateCtrl = TextEditingController(text: emi?.interestRate.toString());
    final monthsCtrl = TextEditingController(text: emi?.totalMonths.toString());
    final paidCtrl = TextEditingController(text: emi?.monthsPaid.toString());

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(emi == null ? 'Add EMI' : 'Update EMI'),
        content: SingleChildScrollView(
          child: Column(
            children: [
              TextField(controller: bankCtrl, decoration: const InputDecoration(labelText: 'Bank Name')),
              TextField(controller: amtCtrl, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'EMI Amount')),
              TextField(controller: rateCtrl, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'Interest Rate (%)')),
              TextField(controller: monthsCtrl, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'Total Months')),
              TextField(controller: paidCtrl, keyboardType: TextInputType.number, decoration: const InputDecoration(labelText: 'Months Paid')),
            ],
          ),
        ),
        actions: [
          TextButton(onPressed: () => Navigator.pop(ctx), child: const Text('Cancel')),
          TextButton(
            onPressed: () async {
              final body = {
                'bank_name': bankCtrl.text.trim(),
                'emi_amount': double.tryParse(amtCtrl.text.trim()) ?? 0.0,
                'interest_rate': double.tryParse(rateCtrl.text.trim()) ?? 0.0,
                'total_months': int.tryParse(monthsCtrl.text.trim()) ?? 12,
                'months_paid': int.tryParse(paidCtrl.text.trim()) ?? 0,
                'start_date': DateTime.now().toIso8601String().substring(0, 10),
              };
              await Provider.of<FinancialProvider>(context, listen: false)
                  .addOrUpdateElectronicEmi(widget.device.id, body, isUpdate: emi != null);
              if (mounted) {
                Navigator.pop(ctx);
                _fetchData();
              }
            },
            child: const Text('Save'),
          ),
        ],
      ),
    );
  }

  Widget _buildAdvisoryRow(String label, String value) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label, style: const TextStyle(color: Colors.black87)),
        Text(value, style: const TextStyle(fontWeight: FontWeight.bold)),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.device.name),
        bottom: TabBar(
          controller: _tabController,
          labelColor: const Color(0xFF0D9488),
          unselectedLabelColor: Colors.grey,
          indicatorColor: const Color(0xFF0D9488),
          tabs: const [
            Tab(text: 'Overview'),
            Tab(text: 'Service History'),
            Tab(text: 'EMI Details'),
          ],
        ),
      ),
      body: Consumer<FinancialProvider>(
        builder: (context, provider, child) {
          final currentDevice = provider.electronics.firstWhere((x) => x.id == widget.device.id, orElse: () => widget.device);
          return TabBarView(
            controller: _tabController,
            children: [
              _buildOverviewTab(currentDevice, provider.electronicEmi),
              _buildServiceTab(provider.electronicServices),
              _buildEmiTab(provider.electronicEmi),
            ],
          );
        },
      ),
    );
  }
}
