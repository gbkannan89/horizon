import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../models/financial_models.dart';
import '../utils/ui_utils.dart';
import 'vehicle_detail_screen.dart';
import 'device_detail_screen.dart';

class PortfolioScreen extends StatelessWidget {
  const PortfolioScreen({super.key});

  // ── Shared input decoration ────────────────────────────────────────────────
  static InputDecoration _inputDec(String label, {String? prefix, String? hint}) =>
    InputDecoration(
      labelText: label,
      hintText: hint,
      prefixText: prefix,
      prefixStyle: const TextStyle(fontWeight: FontWeight.w700, color: Color(0xFF0D9488)),
      labelStyle: const TextStyle(color: Colors.grey, fontSize: 14),
      filled: true,
      fillColor: Colors.white,
      enabledBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(14),
        borderSide: BorderSide(color: Colors.grey.shade200, width: 1.5),
      ),
      focusedBorder: OutlineInputBorder(
        borderRadius: BorderRadius.circular(14),
        borderSide: const BorderSide(color: Color(0xFF0D9488), width: 2),
      ),
      contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 16),
    );

  // ── Shared modal bottom sheet chrome ──────────────────────────────────────
  static BoxDecoration _sheetDec() => const BoxDecoration(
    color: Color(0xFFF8F9FA),
    borderRadius: BorderRadius.vertical(top: Radius.circular(28)),
  );

  static Widget _handle() => Center(
    child: Container(
      width: 40, height: 4, margin: const EdgeInsets.only(top: 12, bottom: 20),
      decoration: BoxDecoration(color: Colors.grey.shade300, borderRadius: BorderRadius.circular(4)),
    ),
  );

  static void _showSnack(BuildContext context, String message, {bool isError = false}) {
    ScaffoldMessenger.of(context).showSnackBar(SnackBar(
      content: Row(
        children: [
          Icon(isError ? Icons.error_outline : Icons.check_circle_outline, color: Colors.white),
          const SizedBox(width: 12),
          Expanded(child: Text(message)),
        ],
      ),
      backgroundColor: isError ? Colors.redAccent : const Color(0xFF059669),
      behavior: SnackBarBehavior.floating,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      margin: const EdgeInsets.all(16),
    ));
  }



  // ── ADD ASSET MODAL ────────────────────────────────────────────────────────
  void _showAddAssetModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    final interestRateCtrl = TextEditingController();
    final purchasePriceCtrl = TextEditingController();
    String selectedType = 'bank';
    bool isLiability = false;
    bool generatesIncome = false;
    String incomeFrequency = 'monthly';
    DateTime? purchaseDate;
    bool isLoading = false;

    final types = [
      {'value': 'bank',     'label': 'Bank/FD',   'icon': Icons.account_balance_outlined},
      {'value': 'physical', 'label': 'Gold/Property', 'icon': Icons.diamond_outlined},
      {'value': 'equity',   'label': 'Stocks/MF', 'icon': Icons.trending_up_rounded},
    ];

    showModalBottomSheet(
      context: context, isScrollControlled: true, backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(builder: (ctx, setState) => Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
        child: Container(
          decoration: _sheetDec(),
          padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
          child: SingleChildScrollView(child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
            _handle(),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              const Text('Add Asset', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
              IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(ctx)),
            ]),
            const SizedBox(height: 20),
            const Text('Asset Name', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: nameCtrl, textCapitalization: TextCapitalization.words,
              decoration: _inputDec('e.g. SBI Savings, HDFC FD, Gold')),
            const SizedBox(height: 20),
            const Text('Asset Type', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 10),
            Row(children: types.map((t) {
              final sel = selectedType == t['value'];
              return Expanded(child: GestureDetector(
                onTap: () => setState(() => selectedType = t['value'] as String),
                child: AnimatedContainer(
                  duration: const Duration(milliseconds: 180),
                  margin: const EdgeInsets.only(right: 8),
                  padding: const EdgeInsets.symmetric(vertical: 12),
                  decoration: BoxDecoration(
                    color: sel ? const Color(0xFF059669) : Colors.white,
                    borderRadius: BorderRadius.circular(14),
                    border: Border.all(color: sel ? const Color(0xFF059669) : Colors.grey.shade200),
                  ),
                  child: Column(children: [
                    Icon(t['icon'] as IconData, color: sel ? Colors.white : Colors.grey, size: 22),
                    const SizedBox(height: 6),
                    Text(t['label'] as String, style: TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: sel ? Colors.white : Colors.grey.shade600)),
                  ]),
                ),
              ));
            }).toList()),
            const SizedBox(height: 20),
            const Text('Current Value', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: amountCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: _inputDec('Amount', prefix: '₹ ', hint: '100000')),
            const SizedBox(height: 20),

            // Enhanced asset details
            const Text('Interest Rate (%)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: interestRateCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: _inputDec('e.g. 7.5', hint: 'Interest Rate (if applicable)')),
            const SizedBox(height: 16),

            const Text('Purchase Price (for ROI)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: purchasePriceCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: _inputDec('What you paid', prefix: '₹ ', hint: 'Optional')),
            const SizedBox(height: 16),

            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Purchase Date', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                TextButton(
                  onPressed: () async {
                    final d = await showDatePicker(context: ctx, initialDate: purchaseDate ?? DateTime.now(), firstDate: DateTime(2000), lastDate: DateTime.now());
                    if (d != null) setState(() => purchaseDate = d);
                  },
                  child: Text(purchaseDate == null ? 'Select (Optional)' : '${purchaseDate!.day}/${purchaseDate!.month}/${purchaseDate!.year}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                ),
              ],
            ),
            const SizedBox(height: 16),

            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Is this a Liability?', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                Switch(value: isLiability, onChanged: (v) => setState(() => isLiability = v), activeThumbColor: const Color(0xFF059669)),
              ],
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Generates Income?', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                Switch(value: generatesIncome, onChanged: (v) => setState(() => generatesIncome = v), activeThumbColor: const Color(0xFF059669)),
              ],
            ),
            if (generatesIncome) ...[
              const SizedBox(height: 12),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  const Text('Income Frequency', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                  DropdownButton<String>(
                    value: incomeFrequency,
                    underline: const SizedBox(),
                    items: ['monthly', 'quarterly', 'yearly'].map((f) => DropdownMenuItem(value: f, child: Text(f.toUpperCase(), style: const TextStyle(fontSize: 13)))).toList(),
                    onChanged: (v) => setState(() => incomeFrequency = v!),
                  ),
                ],
              ),
            ],
            const SizedBox(height: 32),
            SizedBox(width: double.infinity, height: 52, child: ElevatedButton(
              onPressed: isLoading ? null : () async {
                final name = nameCtrl.text.trim();
                final amount = double.tryParse(amountCtrl.text.trim());
                if (name.isEmpty || amount == null || amount <= 0) return;
                final rate = double.tryParse(interestRateCtrl.text.trim()) ?? 0.0;
                setState(() => isLoading = true);
                try {
                  final pp = double.tryParse(purchasePriceCtrl.text.trim());
                  final pd = purchaseDate?.toIso8601String().split('T').first;
                  await Provider.of<FinancialProvider>(ctx, listen: false).addAsset(name, selectedType, amount, interestRate: rate, isLiability: isLiability, generatesIncome: generatesIncome, incomeFrequency: generatesIncome ? incomeFrequency : null, purchasePrice: pp, purchaseDate: pd);
                  if (ctx.mounted) {
                    Navigator.pop(ctx);
                    _showSnack(ctx, 'Asset added successfully!');
                  }
                } catch (e) {
                  setState(() => isLoading = false);
                  if (ctx.mounted) _showSnack(ctx, 'Failed to add asset: $e', isError: true);
                }
              },
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF059669), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)), elevation: 0),
              child: isLoading ? const SizedBox(height: 22, width: 22, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                : const Text('Add Asset', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
            )),
          ])),
        ),
      )),
    );
  }

  // ── ADD LIABILITY MODAL ────────────────────────────────────────────────────
  void _showAddLiabilityModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final outstandingCtrl = TextEditingController();
    final emiCtrl = TextEditingController();
    final rateCtrl = TextEditingController();
    String selectedType = 'home_loan';
    bool isLoading = false;

    final types = [
      {'value': 'home_loan',     'label': 'Home\nLoan',     'icon': Icons.home_outlined},
      {'value': 'car_loan',      'label': 'Car\nLoan',      'icon': Icons.directions_car_outlined},
      {'value': 'personal_loan', 'label': 'Personal\nLoan', 'icon': Icons.person_outline},
      {'value': 'credit_card',   'label': 'Credit\nCard',   'icon': Icons.credit_card_outlined},
    ];

    showModalBottomSheet(
      context: context, isScrollControlled: true, backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(builder: (ctx, setState) => Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
        child: Container(
          decoration: _sheetDec(),
          padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
          child: SingleChildScrollView(child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
            _handle(),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              const Text('Add Liability', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
              IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(ctx)),
            ]),
            const SizedBox(height: 20),
            const Text('Loan / Debt Name', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: nameCtrl, textCapitalization: TextCapitalization.words,
              decoration: _inputDec('e.g. SBI Home Loan, HDFC Credit Card')),
            const SizedBox(height: 20),
            const Text('Type', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 10),
            Row(children: types.map((t) {
              final sel = selectedType == t['value'];
              return Expanded(child: GestureDetector(
                onTap: () => setState(() => selectedType = t['value'] as String),
                child: AnimatedContainer(
                  duration: const Duration(milliseconds: 180),
                  margin: const EdgeInsets.only(right: 6),
                  padding: const EdgeInsets.symmetric(vertical: 10),
                  decoration: BoxDecoration(
                    color: sel ? const Color(0xFFDC2626) : Colors.white,
                    borderRadius: BorderRadius.circular(12),
                    border: Border.all(color: sel ? const Color(0xFFDC2626) : Colors.grey.shade200),
                  ),
                  child: Column(children: [
                    Icon(t['icon'] as IconData, color: sel ? Colors.white : Colors.grey, size: 20),
                    const SizedBox(height: 4),
                    Text(t['label'] as String, textAlign: TextAlign.center, style: TextStyle(fontSize: 10, fontWeight: FontWeight.w600, color: sel ? Colors.white : Colors.grey.shade600)),
                  ]),
                ),
              ));
            }).toList()),
            const SizedBox(height: 20),
            const Text('Outstanding Amount', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: outstandingCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: _inputDec('Total remaining balance', prefix: '₹ ')),
            const SizedBox(height: 16),
            Row(children: [
              Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                const Text('Monthly EMI', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                const SizedBox(height: 8),
                TextField(controller: emiCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  decoration: _inputDec('EMI', prefix: '₹ ')),
              ])),
              const SizedBox(width: 12),
              Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                const Text('Interest Rate', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                const SizedBox(height: 8),
                TextField(controller: rateCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  decoration: _inputDec('Rate', hint: '8.5')),
              ])),
            ]),
            const SizedBox(height: 32),
            SizedBox(width: double.infinity, height: 52, child: ElevatedButton(
              onPressed: isLoading ? null : () async {
                final name = nameCtrl.text.trim();
                final outstanding = double.tryParse(outstandingCtrl.text.trim());
                if (name.isEmpty || outstanding == null || outstanding <= 0) return;
                final emi = double.tryParse(emiCtrl.text.trim()) ?? 0.0;
                final rate = double.tryParse(rateCtrl.text.trim()) ?? 0.0;
                setState(() => isLoading = true);
                try {
                  await Provider.of<FinancialProvider>(ctx, listen: false).addLiability(name, selectedType, outstanding, emi, rate);
                  if (ctx.mounted) {
                    Navigator.pop(ctx);
                    _showSnack(ctx, 'Liability added successfully!');
                  }
                } catch (e) {
                  setState(() => isLoading = false);
                  if (ctx.mounted) _showSnack(ctx, 'Failed to add liability: $e', isError: true);
                }
              },
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFFDC2626), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)), elevation: 0),
              child: isLoading ? const SizedBox(height: 22, width: 22, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                : const Text('Add Liability', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
            )),
          ])),
        ),
      )),
    );
  }

  // ── ADD INCOME MODAL ───────────────────────────────────────────────────────
  void _showAddIncomeModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    String selectedType = 'salary';
    bool isLoading = false;

    showModalBottomSheet(
      context: context, isScrollControlled: true, backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(builder: (ctx, setState) => Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
        child: Container(
          decoration: _sheetDec(),
          padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
          child: SingleChildScrollView(child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
            _handle(),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              const Text('Add Income', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
              IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(ctx)),
            ]),
            const SizedBox(height: 20),
            const Text('Label', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: nameCtrl, textCapitalization: TextCapitalization.words,
              decoration: _inputDec("e.g. My Salary, Partner Income")),
            const SizedBox(height: 20),
            const Text('Monthly Amount', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: amountCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: _inputDec('Amount', prefix: '₹ ', hint: '50000')),
            const SizedBox(height: 32),
            SizedBox(width: double.infinity, height: 52, child: ElevatedButton(
              onPressed: isLoading ? null : () async {
                final name = nameCtrl.text.trim();
                final amount = double.tryParse(amountCtrl.text.trim());
                if (name.isEmpty || amount == null || amount <= 0) return;
                setState(() => isLoading = true);
                try {
                  await Provider.of<FinancialProvider>(ctx, listen: false).addIncome(name, selectedType, amount, 'monthly');
                  if (ctx.mounted) {
                    Navigator.pop(ctx);
                    _showSnack(ctx, 'Income added successfully!');
                  }
                } catch (e) {
                  setState(() => isLoading = false);
                  if (ctx.mounted) _showSnack(ctx, 'Failed to add income: $e', isError: true);
                }
              },
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)), elevation: 0),
              child: isLoading ? const SizedBox(height: 22, width: 22, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                : const Text('Add Income', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
            )),
          ])),
        ),
      )),
    );
  }

  // ── ADD VEHICLE MODAL ──────────────────────────────────────────────────────
  void _showAddVehicleModal(BuildContext context) {
    final makeModelCtrl = TextEditingController();
    final costCtrl = TextEditingController();
    final modelYearCtrl = TextEditingController();
    final purchaseYearCtrl = TextEditingController();
    final mileageCtrl = TextEditingController();
    final regCtrl = TextEditingController();
    
    DateTime? insuranceDate;
    String selectedFuel = 'petrol';
    bool isLoading = false;

    showModalBottomSheet(
      context: context, isScrollControlled: true, backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(builder: (ctx, setState) => Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
        child: Container(
          decoration: _sheetDec(),
          padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
          child: SingleChildScrollView(child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
            _handle(),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              const Text('Add Vehicle', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
              IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(ctx)),
            ]),
            const SizedBox(height: 20),
            const Text('Make & Model', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: makeModelCtrl, textCapitalization: TextCapitalization.words,
              decoration: _inputDec("e.g. Honda City, Royal Enfield")),
            const SizedBox(height: 16),
            const Text('Purchase Cost', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: costCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true),
              decoration: _inputDec('Amount', prefix: '₹ ', hint: '800000')),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Model Year', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      TextField(controller: modelYearCtrl, keyboardType: TextInputType.number, decoration: _inputDec('e.g. 2022')),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Purchase Year', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      TextField(controller: purchaseYearCtrl, keyboardType: TextInputType.number, decoration: _inputDec('e.g. 2023')),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Fuel Type', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      DropdownButtonFormField<String>(
                        value: selectedFuel,
                        decoration: _inputDec('Fuel'),
                        items: const [
                          DropdownMenuItem(value: 'petrol', child: Text('Petrol')),
                          DropdownMenuItem(value: 'diesel', child: Text('Diesel')),
                          DropdownMenuItem(value: 'electric', child: Text('Electric')),
                          DropdownMenuItem(value: 'cng', child: Text('CNG')),
                          DropdownMenuItem(value: 'hybrid', child: Text('Hybrid')),
                        ],
                        onChanged: (val) {
                          if (val != null) setState(() => selectedFuel = val);
                        },
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Mileage (Kmpl)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      TextField(controller: mileageCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true), decoration: _inputDec('e.g. 15.4')),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            const Text('Registration Number', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: regCtrl, textCapitalization: TextCapitalization.characters, decoration: _inputDec('e.g. TN 01 AB 1234')),
            const SizedBox(height: 20),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Insurance Renewal', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                TextButton(
                  onPressed: () async {
                    final d = await showDatePicker(
                      context: ctx,
                      initialDate: DateTime.now().add(const Duration(days: 30)),
                      firstDate: DateTime.now(),
                      lastDate: DateTime.now().add(const Duration(days: 3650)),
                    );
                    if (d != null) setState(() => insuranceDate = d);
                  },
                  child: Text(insuranceDate == null ? 'Select Date' : '${insuranceDate!.day}/${insuranceDate!.month}/${insuranceDate!.year}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
                )
              ],
            ),
            const SizedBox(height: 32),
            SizedBox(width: double.infinity, height: 52, child: ElevatedButton(
              onPressed: isLoading ? null : () async {
                final makeModel = makeModelCtrl.text.trim();
                final cost = double.tryParse(costCtrl.text.trim());
                if (makeModel.isEmpty || cost == null || cost <= 0) return;
                setState(() => isLoading = true);
                try {
                  await Provider.of<FinancialProvider>(ctx, listen: false).addVehicle(
                    makeModel,
                    cost,
                    insuranceRenewalDate: insuranceDate,
                    modelYear: int.tryParse(modelYearCtrl.text.trim()),
                    purchaseYear: int.tryParse(purchaseYearCtrl.text.trim()),
                    fuelType: selectedFuel,
                    mileageKmpl: double.tryParse(mileageCtrl.text.trim()),
                    registrationNumber: regCtrl.text.trim().isEmpty ? null : regCtrl.text.trim(),
                  );
                  if (ctx.mounted) {
                    Navigator.pop(ctx);
                    _showSnack(ctx, 'Vehicle added successfully!');
                  }
                } catch (e) {
                  setState(() => isLoading = false);
                  if (ctx.mounted) _showSnack(ctx, 'Failed to add vehicle: $e', isError: true);
                }
              },
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFFE88A1A), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)), elevation: 0),
              child: isLoading ? const SizedBox(height: 22, width: 22, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                : const Text('Add Vehicle', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
            )),
          ])),
        ),
      )),
    );
  }

  void _showAddElectronicModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final brandCtrl = TextEditingController();
    final modelCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    final warrantyCtrl = TextEditingController(text: '1');
    final expectedLifeCtrl = TextEditingController(text: '3');
    final notesCtrl = TextEditingController();
    
    String selectedCategory = 'mobile';
    bool isLoading = false;

    showModalBottomSheet(
      context: context, isScrollControlled: true, backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(builder: (ctx, setState) => Padding(
        padding: EdgeInsets.only(bottom: MediaQuery.of(ctx).viewInsets.bottom),
        child: Container(
          decoration: _sheetDec(),
          padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
          child: SingleChildScrollView(child: Column(mainAxisSize: MainAxisSize.min, crossAxisAlignment: CrossAxisAlignment.start, children: [
            _handle(),
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              const Text('Add Electronic Device', style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
              IconButton(icon: const Icon(Icons.close), onPressed: () => Navigator.pop(ctx)),
            ]),
            const SizedBox(height: 20),
            const Text('Device Name', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: nameCtrl, textCapitalization: TextCapitalization.words,
              decoration: _inputDec("e.g. iPhone 15 Pro, MacBook Air")),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Brand (Optional)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      TextField(controller: brandCtrl, textCapitalization: TextCapitalization.words, decoration: _inputDec('e.g. Apple')),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Model (Optional)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      TextField(controller: modelCtrl, decoration: _inputDec('e.g. A3106')),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Category', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      DropdownButtonFormField<String>(
                        value: selectedCategory,
                        decoration: _inputDec('Category'),
                        items: const [
                          DropdownMenuItem(value: 'mobile', child: Text('Mobile')),
                          DropdownMenuItem(value: 'laptop', child: Text('Laptop')),
                          DropdownMenuItem(value: 'tablet', child: Text('Tablet')),
                          DropdownMenuItem(value: 'tv', child: Text('TV / Display')),
                          DropdownMenuItem(value: 'appliance', child: Text('Appliance')),
                          DropdownMenuItem(value: 'other', child: Text('Other')),
                        ],
                        onChanged: (val) {
                          if (val != null) setState(() => selectedCategory = val);
                        },
                      ),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Purchase Cost', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      TextField(controller: amountCtrl, keyboardType: const TextInputType.numberWithOptions(decimal: true), decoration: _inputDec('₹ amount')),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Warranty (Years)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      TextField(controller: warrantyCtrl, keyboardType: TextInputType.number, decoration: _inputDec('e.g. 1')),
                    ],
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text('Expected Life (Years)', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
                      const SizedBox(height: 8),
                      TextField(controller: expectedLifeCtrl, keyboardType: TextInputType.number, decoration: _inputDec('e.g. 3')),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            const Text('Notes', style: TextStyle(fontWeight: FontWeight.w600, fontSize: 13, color: Colors.grey)),
            const SizedBox(height: 8),
            TextField(controller: notesCtrl, maxLines: 2, decoration: _inputDec('Additional details')),
            const SizedBox(height: 32),
            SizedBox(width: double.infinity, height: 52, child: ElevatedButton(
              onPressed: isLoading ? null : () async {
                final name = nameCtrl.text.trim();
                final amt = double.tryParse(amountCtrl.text.trim());
                if (name.isEmpty || amt == null || amt <= 0) return;
                setState(() => isLoading = true);
                try {
                  await Provider.of<FinancialProvider>(ctx, listen: false).addElectronic({
                    'name': name,
                    'category': selectedCategory,
                    'brand': brandCtrl.text.trim().isEmpty ? null : brandCtrl.text.trim(),
                    'model': modelCtrl.text.trim().isEmpty ? null : modelCtrl.text.trim(),
                    'purchase_date': DateTime.now().toIso8601String().substring(0, 10),
                    'purchase_amount': amt,
                    'warranty_years': int.tryParse(warrantyCtrl.text.trim()) ?? 1,
                    'expected_life_years': int.tryParse(expectedLifeCtrl.text.trim()) ?? 3,
                    'notes': notesCtrl.text.trim().isEmpty ? null : notesCtrl.text.trim(),
                  });
                  if (ctx.mounted) {
                    Navigator.pop(ctx);
                    _showSnack(ctx, 'Device added successfully!');
                  }
                } catch (e) {
                  setState(() => isLoading = false);
                  if (ctx.mounted) _showSnack(ctx, 'Failed to add device: $e', isError: true);
                }
              },
              style: ElevatedButton.styleFrom(backgroundColor: const Color(0xFF0D9488), shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(14)), elevation: 0),
              child: isLoading ? const SizedBox(height: 22, width: 22, child: CircularProgressIndicator(color: Colors.white, strokeWidth: 2))
                : const Text('Add Device', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
            )),
          ])),
        ),
      )),
    );
  }

  // ── Section Header ─────────────────────────────────────────────────────────
  Widget _buildSectionHeader(String title, List<Color> gradient, VoidCallback onAdd) {
    return Padding(
      padding: const EdgeInsets.only(top: 24.0, bottom: 16.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(title, style: const TextStyle(fontSize: 19, fontWeight: FontWeight.w900, letterSpacing: -0.5, color: Colors.white)),
          GestureDetector(
            onTap: onAdd,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 7),
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha: 0.06),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: Colors.white.withValues(alpha: 0.12), width: 1.2),
              ),
              child: const Row(children: [
                Icon(Icons.add_rounded, size: 16, color: Colors.white),
                SizedBox(width: 4),
                Text('Add', style: TextStyle(color: Colors.white, fontWeight: FontWeight.w800, fontSize: 12)),
              ]),
            ),
          ),
        ],
      ),
    );
  }

  // ── Portfolio Summary ──────────────────────────────────────────────────────
  List<Widget> _buildPortfolioSummary(FinancialProvider p) {
    final ps = p.portfolioSummary;
    if (ps == null || (ps['totalInvested'] ?? 0) == 0) return [];

    final totalInvested = (ps['totalInvested'] ?? 0).toDouble();
    final totalCurrent = (ps['totalCurrent'] ?? 0).toDouble();
    final totalReturn = (ps['totalReturn'] ?? 0).toDouble();
    final estimatedCagr = (ps['estimatedCagr'] ?? 0).toDouble();
    final allocation = ps['allocation'] as List<dynamic>? ?? [];
    final returnColor = totalReturn >= 0 ? const Color(0xFF34D399) : const Color(0xFFF87171);
    final returnSign = totalReturn >= 0 ? '+' : '';

    return [
      const SizedBox(height: 20),
      Container(
        padding: const EdgeInsets.all(20),
        decoration: BoxDecoration(
          color: Colors.white.withValues(alpha: 0.04),
          borderRadius: BorderRadius.circular(24),
          border: Border.all(color: Colors.white.withValues(alpha: 0.08)),
        ),
        child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
          Row(children: [
            Container(
              padding: const EdgeInsets.all(8),
              decoration: BoxDecoration(
                color: Colors.white.withValues(alpha: 0.06),
                borderRadius: BorderRadius.circular(10),
              ),
              child: const Icon(Icons.trending_up_rounded, color: Colors.white, size: 18),
            ),
            const SizedBox(width: 12),
            const Text('Portfolio Summary', style: TextStyle(fontSize: 15, fontWeight: FontWeight.w800, color: Colors.white)),
          ]),
          const SizedBox(height: 16),
          Row(children: [
            _psColumn('Invested', '₹${_fmt(totalInvested)}', Colors.white.withValues(alpha: 0.7)),
            const SizedBox(width: 16),
            _psColumn('Current', '₹${_fmt(totalCurrent)}', Colors.white),
            const SizedBox(width: 16),
            _psColumn('Return', '$returnSign₹${_fmt(totalReturn.abs())}', returnColor),
            const SizedBox(width: 16),
            _psColumn('Est. CAGR', '${estimatedCagr.toStringAsFixed(1)}%', const Color(0xFFA78BFA)),
          ]),
          if (allocation.isNotEmpty) ...[
            const SizedBox(height: 16),
            Divider(height: 1, color: Colors.white.withValues(alpha: 0.08)),
            const SizedBox(height: 16),
            ...allocation.take(5).map((a) {
              final label = a['label'] ?? '';
              final pct = (a['percentage'] ?? 0).toDouble();
              final clr = _parseHex(a['color'] ?? '#6366F1');
              return Padding(
                padding: const EdgeInsets.only(bottom: 8),
                child: Row(children: [
                  Container(width: 10, height: 10, decoration: BoxDecoration(color: clr, borderRadius: BorderRadius.circular(3))),
                  const SizedBox(width: 10),
                  Expanded(child: Text(label, style: TextStyle(fontSize: 13, fontWeight: FontWeight.w500, color: Colors.white.withValues(alpha: 0.75)))),
                  Text('${pct.toStringAsFixed(1)}%', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13, color: Colors.white)),
                ]),
              );
            }),
          ],
        ]),
      ),
    ];
  }

  Widget _psColumn(String label, String value, Color color) {
    return Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Text(label, style: TextStyle(fontSize: 10, color: Colors.white.withValues(alpha: 0.4), fontWeight: FontWeight.w600)),
      const SizedBox(height: 2),
      Text(value, style: TextStyle(fontSize: 13, fontWeight: FontWeight.w900, color: color)),
    ]));
  }

  Color _parseHex(String hex) {
    hex = hex.replaceAll('#', '');
    if (hex.length == 6) hex = 'FF$hex';
    return Color(int.parse('0x$hex'));
  }

  // ── Item card with delete swipe ────────────────────────────────────────────
  Widget _buildListItem(String title, double amount, String subtitle, List<Color> gradient, IconData icon, VoidCallback onDelete, {VoidCallback? onTap}) {
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white.withValues(alpha: 0.04), // Glassmorphic background
        borderRadius: BorderRadius.circular(18),
        border: Border.all(color: Colors.white.withValues(alpha: 0.08)),
      ),
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          borderRadius: BorderRadius.circular(18),
          onTap: onTap,
          child: Padding(
            padding: const EdgeInsets.all(16),
            child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Expanded(
                child: Row(children: [
                  Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      color: Colors.white.withValues(alpha: 0.06),
                      shape: BoxShape.circle,
                    ),
                    child: Icon(icon, color: Colors.white, size: 20),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                      Text(title, style: const TextStyle(fontWeight: FontWeight.w800, fontSize: 14, color: Colors.white), overflow: TextOverflow.ellipsis),
                      const SizedBox(height: 3),
                      Text(subtitle, style: TextStyle(color: Colors.white.withValues(alpha: 0.4), fontSize: 11, fontWeight: FontWeight.w500), overflow: TextOverflow.ellipsis),
                    ]),
                  ),
                ]),
              ),
              const SizedBox(width: 8),
              Row(children: [
                Text('₹${amount.toStringAsFixed(0)}', style: const TextStyle(fontWeight: FontWeight.w900, fontSize: 15, color: Colors.white)),
                const SizedBox(width: 8),
                GestureDetector(
                  onTap: onDelete,
                  child: Container(
                    padding: const EdgeInsets.all(6),
                    decoration: BoxDecoration(color: Colors.red.withValues(alpha: 0.08), borderRadius: BorderRadius.circular(8)),
                    child: const Icon(Icons.delete_outline_rounded, color: Color(0xFFF87171), size: 16),
                  ),
                ),
              ]),
            ]),
          ),
        ),
      ),
    );
  }

  // ── BUILD ──────────────────────────────────────────────────────────────────
  @override
  Widget build(BuildContext context) {
    final provider = context.watch<FinancialProvider>();
    double totalAssets = provider.assets.fold(0.0, (s, a) => s + a.amount);
    double totalLiabilities = provider.liabilities.fold(0.0, (s, l) => s + l.amount);

    return Scaffold(
      backgroundColor: const Color(0xFF030712), // Premium dark theme
      body: Stack(
        children: [
          // ── Glowing Orbs (Glassmorphism BG effects) ─────────────────────────
          Positioned(
            top: -150,
            left: -150,
            child: Container(
              width: 350,
              height: 350,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFF0D9488).withValues(alpha: 0.25),
                    const Color(0xFF0D9488).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            bottom: 50,
            right: -150,
            child: Container(
              width: 350,
              height: 350,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFFA78BFA).withValues(alpha: 0.15),
                    const Color(0xFFA78BFA).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),

          // ── Content ─────────────────────────────────────────────────────────
          SafeArea(
            child: provider.isLoading
                ? const Center(child: CircularProgressIndicator(color: Colors.white))
                : RefreshIndicator(
                    onRefresh: provider.loadAllData,
                    color: Colors.white,
                    backgroundColor: const Color(0xFF0F172A),
                    child: ListView(
                      physics: const AlwaysScrollableScrollPhysics(),
                      padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 8),
                      children: [
                        // Net Worth card
                        Container(
                          padding: const EdgeInsets.all(24),
                          decoration: BoxDecoration(
                            color: Colors.black.withValues(alpha: 0.35), // Dark Glass
                            borderRadius: BorderRadius.circular(28),
                            border: Border.all(color: Colors.white.withValues(alpha: 0.1), width: 1.2),
                          ),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text('Total Net Worth', style: TextStyle(color: Colors.white.withValues(alpha: 0.5), fontWeight: FontWeight.w600)),
                                  const SizedBox(height: 8),
                                  Text('₹${provider.netWorth.toStringAsFixed(0)}',
                                      style: const TextStyle(fontSize: 26, fontWeight: FontWeight.w900, color: Colors.white)),
                                  const SizedBox(height: 12),
                                  Row(
                                    children: [
                                      _pill('Assets ₹${_fmt(totalAssets)}', const Color(0xFF34D399)),
                                      const SizedBox(width: 8),
                                      _pill('Debt ₹${_fmt(totalLiabilities)}', const Color(0xFFF87171)),
                                    ],
                                  ),
                                ],
                              ),
                              Container(
                                decoration: BoxDecoration(
                                  shape: BoxShape.circle,
                                  boxShadow: [BoxShadow(color: const Color(0xFF34D399).withValues(alpha: 0.2), blurRadius: 20)],
                                ),
                                child: CircularPercentIndicator(
                                  radius: 44,
                                  lineWidth: 8,
                                  percent: totalAssets > 0 ? (totalAssets / (totalAssets + totalLiabilities)).clamp(0.0, 1.0) : 0,
                                  progressColor: const Color(0xFF34D399),
                                  backgroundColor: Colors.white.withValues(alpha: 0.1),
                                  circularStrokeCap: CircularStrokeCap.round,
                                  center: const Text('NW', style: TextStyle(fontSize: 12, fontWeight: FontWeight.w900, color: Colors.white)),
                                ),
                              ),
                            ],
                          ),
                        ),

                        // Portfolio Summary
                        ..._buildPortfolioSummary(provider),

                        // Assets
                        _buildSectionHeader('Assets', const [Color(0xFF059669), Color(0xFF34D399)], () => _showAddAssetModal(context)),
                        if (provider.assets.isEmpty)
                          UiUtils.buildEmptyState('No assets yet', 'Tap Add to add your savings,\nFDs, gold, property etc.', Icons.account_balance_wallet_outlined, Colors.green),
                        ...provider.assets.map((a) => _buildListItem(
                              a.name,
                              a.amount,
                              a.purchasePrice != null ? 'Invested: ₹${a.purchasePrice!.toStringAsFixed(0)}' : 'Asset',
                              const [Color(0xFF059669), Color(0xFF34D399)],
                              Icons.account_balance_wallet_rounded,
                              () => UiUtils.showDeleteBottomSheet(context, a.name, () async {
                                await Provider.of<FinancialProvider>(context, listen: false).deleteAsset(a.id);
                                if (context.mounted) _showSnack(context, 'Asset deleted');
                              }),
                            )),

                        // Liabilities
                        _buildSectionHeader('Liabilities', const [Color(0xFFDC2626), Color(0xFFF87171)], () => _showAddLiabilityModal(context)),
                        if (provider.liabilities.isEmpty)
                          UiUtils.buildEmptyState('No liabilities', 'Tap Add to track loans,\ncredit cards, EMIs etc.', Icons.credit_card_outlined, Colors.red),
                        ...provider.liabilities.map((l) => _buildListItem(
                              l.name,
                              l.amount,
                              '${l.interestRate.toStringAsFixed(1)}% interest',
                              const [Color(0xFFDC2626), Color(0xFFF87171)],
                              Icons.credit_card_rounded,
                              () => UiUtils.showDeleteBottomSheet(context, l.name, () async {
                                await Provider.of<FinancialProvider>(context, listen: false).deleteLiability(l.id);
                                if (context.mounted) _showSnack(context, 'Liability deleted');
                              }),
                            )),

                        // Incomes
                        _buildSectionHeader('Income Sources', const [Color(0xFF0D9488), Color(0xFF2DD4BF)], () => _showAddIncomeModal(context)),
                        if (provider.incomes.isEmpty)
                          UiUtils.buildEmptyState('No income sources', 'Tap Add to record your salary\nand other income.', Icons.trending_up_outlined, Colors.blue),
                        ...provider.incomes.map((i) => _buildListItem(
                              i.label,
                              i.amount,
                              '${i.frequency} income',
                              const [Color(0xFF0D9488), Color(0xFF2DD4BF)],
                              Icons.trending_up_rounded,
                              () => UiUtils.showDeleteBottomSheet(context, i.label, () async {
                                await Provider.of<FinancialProvider>(context, listen: false).deleteIncome(i.id);
                                if (context.mounted) _showSnack(context, 'Income deleted');
                              }),
                            )),

                        // Vehicles
                        _buildSectionHeader('Vehicles', const [Color(0xFFE88A1A), Color(0xFFFBBF24)], () => _showAddVehicleModal(context)),
                        if (provider.vehicles.isEmpty)
                          UiUtils.buildEmptyState('No vehicles', 'Tap Add to track your cars/bikes\nand insurance.', Icons.directions_car_outlined, Colors.orange),
                        ...provider.vehicles.map((v) => _buildListItem(
                              v.makeModel,
                              v.purchaseCost,
                              v.registrationNumber != null && v.registrationNumber!.isNotEmpty
                                  ? '${v.registrationNumber} • Cost/km: ₹${v.costPerKm.toStringAsFixed(2)}'
                                  : 'Cost/km: ₹${v.costPerKm.toStringAsFixed(2)}',
                              const [Color(0xFFE88A1A), Color(0xFFFBBF24)],
                              Icons.directions_car_rounded,
                              () => UiUtils.showDeleteBottomSheet(context, v.makeModel, () async {
                                await Provider.of<FinancialProvider>(context, listen: false).deleteVehicle(v.id);
                                if (context.mounted) _showSnack(context, 'Vehicle deleted');
                              }),
                              onTap: () {
                                Navigator.push(
                                  context,
                                  MaterialPageRoute(
                                    builder: (_) => VehicleDetailScreen(vehicle: v),
                                  ),
                                );
                              },
                            )),

                        // Electronics
                        _buildSectionHeader('Electronics', const [Color(0xFF0D9488), Color(0xFF2DD4BF)], () => _showAddElectronicModal(context)),
                        if (provider.electronics.isEmpty)
                          UiUtils.buildEmptyState('No devices', 'Tap Add to track phones, laptops\nand warranty.', Icons.phone_android_outlined, Colors.teal),
                        ...provider.electronics.map((e) {
                          IconData icon = Icons.devices_other_rounded;
                          if (e.category.toLowerCase() == 'mobile') icon = Icons.phone_android_rounded;
                          if (e.category.toLowerCase() == 'laptop') icon = Icons.laptop_chromebook_rounded;
                          if (e.category.toLowerCase() == 'tablet') icon = Icons.tablet_mac_rounded;
                          if (e.category.toLowerCase() == 'tv') icon = Icons.tv_rounded;
                          if (e.category.toLowerCase() == 'appliance') icon = Icons.kitchen_rounded;

                          String subtitle = 'Warranty: ${e.warrantyStatus.replaceAll('_', ' ').toUpperCase()}';
                          if (e.warrantyStatus != 'expired') {
                            subtitle += ' (${e.warrantyDaysRemaining}d remaining)';
                          }

                          return _buildListItem(
                            e.name,
                            e.currentValue,
                            subtitle,
                            const [Color(0xFF0D9488), Color(0xFF2DD4BF)],
                            icon,
                            () => UiUtils.showDeleteBottomSheet(context, e.name, () async {
                              await Provider.of<FinancialProvider>(context, listen: false).deleteElectronic(e.id);
                              if (context.mounted) _showSnack(context, 'Device deleted');
                            }),
                            onTap: () {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                  builder: (_) => DeviceDetailScreen(device: e),
                                ),
                              );
                            },
                          );
                        }),

                        const SizedBox(height: 60),
                      ],
                    ),
                  ),
          ),
        ],
      ),
    );
  }

  String _fmt(double v) => v >= 1000 ? '${(v / 1000).toStringAsFixed(0)}K' : v.toStringAsFixed(0);

  Widget _pill(String text, Color color) => Container(
        padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
        decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(20)),
        child: Text(text, style: TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: color)),
      );
}
