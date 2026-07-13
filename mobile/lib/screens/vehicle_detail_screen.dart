import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../models/financial_models.dart';
import '../utils/ui_utils.dart';
import 'vehicle_service_form.dart';
import 'vehicle_fuel_form.dart';

class VehicleDetailScreen extends StatefulWidget {
  final LocalVehicle vehicle;

  const VehicleDetailScreen({super.key, required this.vehicle});

  @override
  State<VehicleDetailScreen> createState() => _VehicleDetailScreenState();
}

class _VehicleDetailScreenState extends State<VehicleDetailScreen>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 4, vsync: this);
    WidgetsBinding.instance.addPostFrameCallback((_) {
      _fetchData();
    });
  }

  void _fetchData() {
    final provider = Provider.of<FinancialProvider>(context, listen: false);
    provider.reloadVehicleServiceRecords(widget.vehicle.id);
    provider.reloadVehicleFuelRecords(widget.vehicle.id);
    provider.reloadVehicleLoan(widget.vehicle.id);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  Widget _buildOverviewTab(LocalVehicle v, LocalVehicleLoan? loan) {
    final insuranceRemaining = v.insuranceRenewalDate
        ?.difference(DateTime.now())
        .inDays;

    return RefreshIndicator(
      onRefresh: () async => _fetchData(),
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Vehicle Spec Card
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withValues(alpha: 0.03),
                    blurRadius: 10,
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        v.makeModel,
                        style: const TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.w900,
                          color: Color(0xFF1E293B),
                        ),
                      ),
                      if (v.fuelType != null)
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 10,
                            vertical: 4,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.blue.withValues(alpha: 0.1),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Text(
                            v.fuelType!.toUpperCase(),
                            style: const TextStyle(
                              color: Colors.blue,
                              fontWeight: FontWeight.bold,
                              fontSize: 11,
                            ),
                          ),
                        ),
                    ],
                  ),
                  const SizedBox(height: 6),
                  if (v.registrationNumber != null &&
                      v.registrationNumber!.isNotEmpty)
                    Text(
                      v.registrationNumber!,
                      style: TextStyle(
                        color: Colors.grey.shade600,
                        fontWeight: FontWeight.w600,
                        fontSize: 13,
                      ),
                    )
                  else
                    const Text(
                      'No Reg. Number',
                      style: TextStyle(color: Colors.grey, fontSize: 13),
                    ),
                  const Divider(height: 30),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      _buildSpecItem(
                        'Model Year',
                        v.modelYear?.toString() ?? 'N/A',
                      ),
                      _buildSpecItem(
                        'Purchase Year',
                        v.purchaseYear?.toString() ?? 'N/A',
                      ),
                      _buildSpecItem(
                        'Mileage (Kmpl)',
                        v.mileageKmpl != null ? '${v.mileageKmpl} km/l' : 'N/A',
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),

            // Cost Stats Card
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withValues(alpha: 0.03),
                    blurRadius: 10,
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Running Cost & Usage',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 20),
                  Row(
                    children: [
                      Expanded(
                        child: _buildMetricTile(
                          'KM Driven',
                          '${v.kmDriven.toStringAsFixed(0)} km',
                          Icons.speed,
                          Colors.purple,
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildMetricTile(
                          'Cost per KM',
                          '₹${v.costPerKm.toStringAsFixed(2)}',
                          Icons.payments_outlined,
                          Colors.teal,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: _buildMetricTile(
                          'Fuel Cost',
                          '₹${v.fuelCostTotal.toStringAsFixed(0)}',
                          Icons.local_gas_station_outlined,
                          Colors.orange,
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: _buildMetricTile(
                          'Service Cost',
                          '₹${v.totalServiceCost.toStringAsFixed(0)}',
                          Icons.build_circle_outlined,
                          Colors.blue,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
            const SizedBox(height: 20),

            // Insurance suggestion Card
            Container(
              padding: const EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withValues(alpha: 0.03),
                    blurRadius: 10,
                  ),
                ],
              ),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Insurance Status & Advisory',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 16),
                  if (insuranceRemaining != null) ...[
                    Row(
                      children: [
                        Icon(
                          insuranceRemaining > 30
                              ? Icons.check_circle_outline
                              : Icons.warning_amber_rounded,
                          color: insuranceRemaining > 30
                              ? Colors.green
                              : Colors.amber,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          insuranceRemaining > 0
                              ? 'Expires in $insuranceRemaining days'
                              : 'Expired ${-insuranceRemaining} days ago',
                          style: TextStyle(
                            fontWeight: FontWeight.bold,
                            color: insuranceRemaining > 30
                                ? Colors.green
                                : Colors.redAccent,
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),
                  ],
                  _buildAdvisoryRow(
                    'Current IDV Value',
                    v.insuranceIdv != null
                        ? '₹${v.insuranceIdv!.toStringAsFixed(0)}'
                        : 'N/A',
                  ),
                  const SizedBox(height: 10),
                  _buildAdvisoryRow(
                    'Recommended IDV',
                    '₹${v.suggestedIdv.toStringAsFixed(0)}',
                  ),
                  const SizedBox(height: 6),
                  Text(
                    'Standard depreciation is applied at 15% for year 1, and 10% per year for years 2-5. Keep IDV close to recommended values for optimum cover.',
                    style: TextStyle(fontSize: 11, color: Colors.grey.shade500),
                  ),
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
        Text(
          value,
          style: const TextStyle(
            fontSize: 14,
            fontWeight: FontWeight.bold,
            color: Color(0xFF1E293B),
          ),
        ),
      ],
    );
  }

  Widget _buildMetricTile(
    String label,
    String value,
    IconData icon,
    Color color,
  ) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.05),
        borderRadius: BorderRadius.circular(16),
        border: Border.all(color: color.withValues(alpha: 0.1)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, color: color, size: 24),
          const SizedBox(height: 12),
          Text(
            label,
            style: TextStyle(
              fontSize: 12,
              color: Colors.grey.shade600,
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 4),
          Text(
            value,
            style: const TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w900,
              color: Color(0xFF1E293B),
            ),
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

  Widget _buildServiceTab(List<LocalVehicleService> services) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8F9FA),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (_) => VehicleServiceForm(vehicleId: widget.vehicle.id),
            ),
          ).then((_) => _fetchData());
        },
        backgroundColor: const Color(0xFF6366F1),
        child: const Icon(Icons.add, color: Colors.white),
      ),
      body: services.isEmpty
          ? const Center(child: Text('No service records logged.'))
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
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withValues(alpha: 0.02),
                        blurRadius: 8,
                      ),
                    ],
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            s.serviceType.toUpperCase(),
                            style: const TextStyle(
                              fontWeight: FontWeight.bold,
                              color: Color(0xFF6366F1),
                            ),
                          ),
                          Text(
                            '₹${s.cost.toStringAsFixed(0)}',
                            style: const TextStyle(
                              fontWeight: FontWeight.w900,
                              fontSize: 16,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 6),
                      Text(
                        s.serviceDate,
                        style: const TextStyle(
                          color: Colors.grey,
                          fontSize: 12,
                        ),
                      ),
                      if (s.kmAtService != null) ...[
                        const SizedBox(height: 8),
                        Text(
                          'Odometer: ${s.kmAtService!.toStringAsFixed(0)} km',
                          style: const TextStyle(
                            fontSize: 13,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                      if (s.description != null &&
                          s.description!.isNotEmpty) ...[
                        const SizedBox(height: 8),
                        Text(
                          s.description!,
                          style: TextStyle(
                            color: Colors.grey.shade700,
                            fontSize: 13,
                          ),
                        ),
                      ],
                      if (s.nextServiceKm != null) ...[
                        const SizedBox(height: 8),
                        Container(
                          padding: const EdgeInsets.all(8),
                          decoration: BoxDecoration(
                            color: Colors.orange.withValues(alpha: 0.1),
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Row(
                            children: [
                              const Icon(
                                Icons.info_outline,
                                size: 14,
                                color: Colors.orange,
                              ),
                              const SizedBox(width: 6),
                              Text(
                                'Next service due at: ${s.nextServiceKm!.toStringAsFixed(0)} km',
                                style: const TextStyle(
                                  color: Colors.orange,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 11,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                      const Divider(height: 24),
                      Align(
                        alignment: Alignment.centerRight,
                        child: GestureDetector(
                          onTap: () async {
                            UiUtils.showDeleteBottomSheet(
                              context,
                              'Service Record',
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteServiceRecord(widget.vehicle.id, s.id);
                                _fetchData();
                              },
                            );
                          },
                          child: const Text(
                            'Delete',
                            style: TextStyle(
                              color: Colors.redAccent,
                              fontWeight: FontWeight.bold,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                );
              },
            ),
    );
  }

  Widget _buildFuelTab(List<LocalVehicleFuel> fuels) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8F9FA),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (_) => VehicleFuelForm(vehicleId: widget.vehicle.id),
            ),
          ).then((_) => _fetchData());
        },
        backgroundColor: const Color(0xFF6366F1),
        child: const Icon(Icons.add, color: Colors.white),
      ),
      body: fuels.isEmpty
          ? const Center(child: Text('No fuel logs found.'))
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: fuels.length,
              itemBuilder: (context, idx) {
                final f = fuels[idx];
                return Container(
                  margin: const EdgeInsets.only(bottom: 12),
                  padding: const EdgeInsets.all(16),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(16),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withValues(alpha: 0.02),
                        blurRadius: 8,
                      ),
                    ],
                  ),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            '${f.liters} Liters filled',
                            style: const TextStyle(fontWeight: FontWeight.bold),
                          ),
                          Text(
                            '₹${f.amount.toStringAsFixed(0)}',
                            style: const TextStyle(
                              fontWeight: FontWeight.w900,
                              fontSize: 16,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 6),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            f.fillDate,
                            style: const TextStyle(
                              color: Colors.grey,
                              fontSize: 12,
                            ),
                          ),
                          if (f.pricePerLiter != null)
                            Text(
                              '₹${f.pricePerLiter!.toStringAsFixed(2)}/L',
                              style: const TextStyle(
                                color: Colors.grey,
                                fontSize: 12,
                              ),
                            ),
                        ],
                      ),
                      if (f.kmAtFill != null) ...[
                        const SizedBox(height: 8),
                        Text(
                          'Odometer: ${f.kmAtFill!.toStringAsFixed(0)} km',
                          style: const TextStyle(
                            fontSize: 13,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                      const Divider(height: 24),
                      Align(
                        alignment: Alignment.centerRight,
                        child: GestureDetector(
                          onTap: () async {
                            UiUtils.showDeleteBottomSheet(
                              context,
                              'Fuel Fill Record',
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteFuelRecord(widget.vehicle.id, f.id);
                                _fetchData();
                              },
                            );
                          },
                          child: const Text(
                            'Delete',
                            style: TextStyle(
                              color: Colors.redAccent,
                              fontWeight: FontWeight.bold,
                              fontSize: 12,
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                );
              },
            ),
    );
  }

  Widget _buildLoanTab(LocalVehicleLoan? loan) {
    if (loan == null) {
      return Padding(
        padding: const EdgeInsets.all(24),
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const Icon(
                Icons.account_balance_outlined,
                size: 56,
                color: Colors.grey,
              ),
              const SizedBox(height: 12),
              const Text(
                'No loan linked to this vehicle.',
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 16),
              ElevatedButton(
                onPressed: () => _showLoanDialog(null),
                style: ElevatedButton.styleFrom(
                  backgroundColor: const Color(0xFF6366F1),
                ),
                child: const Text(
                  'Add Vehicle Loan',
                  style: TextStyle(color: Colors.white),
                ),
              ),
            ],
          ),
        ),
      );
    }

    final totalPaid = loan.emi * loan.emiPaid;

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
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.03),
                  blurRadius: 10,
                ),
              ],
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      loan.bankName ?? 'Vehicle Loan',
                      style: const TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    IconButton(
                      icon: const Icon(Icons.edit_note, color: Colors.blue),
                      onPressed: () => _showLoanDialog(loan),
                    ),
                  ],
                ),
                const Divider(),
                _buildAdvisoryRow(
                  'Loan Amount',
                  '₹${loan.loanAmount.toStringAsFixed(0)}',
                ),
                const SizedBox(height: 12),
                _buildAdvisoryRow('Interest Rate', '${loan.interestRate}%'),
                const SizedBox(height: 12),
                _buildAdvisoryRow('Tenure', '${loan.tenureMonths} Months'),
                const SizedBox(height: 12),
                _buildAdvisoryRow(
                  'Monthly EMI',
                  '₹${loan.emi.toStringAsFixed(0)}',
                ),
                const SizedBox(height: 12),
                _buildAdvisoryRow(
                  'EMIs Paid',
                  '${loan.emiPaid} / ${loan.tenureMonths}',
                ),
                const SizedBox(height: 12),
                _buildAdvisoryRow(
                  'Total Repaid',
                  '₹${totalPaid.toStringAsFixed(0)}',
                ),
                const SizedBox(height: 12),
                _buildAdvisoryRow(
                  'Outstanding Principal',
                  '₹${loan.outstanding.toStringAsFixed(0)}',
                ),
              ],
            ),
          ),
          const SizedBox(height: 24),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: () {
                UiUtils.showDeleteBottomSheet(
                  context,
                  'Vehicle Loan',
                  () async {
                    await Provider.of<FinancialProvider>(
                      context,
                      listen: false,
                    ).deleteVehicleLoan(widget.vehicle.id);
                    _fetchData();
                  },
                );
              },
              icon: const Icon(Icons.delete_outline, color: Colors.red),
              label: const Text(
                'Remove Loan Link',
                style: TextStyle(color: Colors.red),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _showLoanDialog(LocalVehicleLoan? loan) {
    final bankCtrl = TextEditingController(text: loan?.bankName);
    final amtCtrl = TextEditingController(
      text: loan?.loanAmount.toStringAsFixed(0),
    );
    final rateCtrl = TextEditingController(text: loan?.interestRate.toString());
    final tenureCtrl = TextEditingController(
      text: loan?.tenureMonths.toString(),
    );
    final emiCtrl = TextEditingController(text: loan?.emi.toStringAsFixed(0));
    final paidCtrl = TextEditingController(text: loan?.emiPaid.toString());
    final outstandingCtrl = TextEditingController(
      text: loan?.outstanding.toStringAsFixed(0),
    );

    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: Text(loan == null ? 'Add Vehicle Loan' : 'Update Loan'),
        content: SingleChildScrollView(
          child: Column(
            children: [
              TextField(
                controller: bankCtrl,
                decoration: const InputDecoration(labelText: 'Bank Name'),
              ),
              TextField(
                controller: amtCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: 'Loan Amount'),
              ),
              TextField(
                controller: rateCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Interest Rate (%)',
                ),
              ),
              TextField(
                controller: tenureCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: 'Tenure (Months)'),
              ),
              TextField(
                controller: emiCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: 'EMI Amount'),
              ),
              TextField(
                controller: paidCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(labelText: 'EMIs Paid'),
              ),
              TextField(
                controller: outstandingCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Outstanding Balance',
                ),
              ),
            ],
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(ctx),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () async {
              final body = {
                'bank_name': bankCtrl.text.trim(),
                'loan_amount': double.tryParse(amtCtrl.text.trim()) ?? 0.0,
                'interest_rate': double.tryParse(rateCtrl.text.trim()) ?? 0.0,
                'tenure_months': int.tryParse(tenureCtrl.text.trim()) ?? 12,
                'emi': double.tryParse(emiCtrl.text.trim()) ?? 0.0,
                'start_date': DateTime.now().toIso8601String().substring(0, 10),
                'emi_paid': int.tryParse(paidCtrl.text.trim()) ?? 0,
                'outstanding':
                    double.tryParse(outstandingCtrl.text.trim()) ?? 0.0,
              };
              await Provider.of<FinancialProvider>(
                context,
                listen: false,
              ).addOrUpdateVehicleLoan(
                widget.vehicle.id,
                body,
                isUpdate: loan != null,
              );
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

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.vehicle.makeModel),
        bottom: TabBar(
          controller: _tabController,
          labelColor: const Color(0xFF6366F1),
          unselectedLabelColor: Colors.grey,
          indicatorColor: const Color(0xFF6366F1),
          tabs: const [
            Tab(text: 'Overview'),
            Tab(text: 'Service'),
            Tab(text: 'Fuel'),
            Tab(text: 'Loan'),
          ],
        ),
      ),
      body: Consumer<FinancialProvider>(
        builder: (context, provider, child) {
          // Find this specific vehicle in provider to get refreshed computes
          final currentVehicle = provider.vehicles.firstWhere(
            (x) => x.id == widget.vehicle.id,
            orElse: () => widget.vehicle,
          );
          return TabBarView(
            controller: _tabController,
            children: [
              _buildOverviewTab(currentVehicle, provider.vehicleLoan),
              _buildServiceTab(provider.vehicleServices),
              _buildFuelTab(provider.vehicleFuels),
              _buildLoanTab(provider.vehicleLoan),
            ],
          );
        },
      ),
    );
  }
}
