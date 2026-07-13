import 'dart:ui' as ui;
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';
import 'vehicle_detail_screen.dart';
import 'device_detail_screen.dart';

class PortfolioScreen extends StatelessWidget {
  const PortfolioScreen({super.key});

  // ── Shared input decoration ────────────────────────────────────────────────
  static InputDecoration _inputDec(
    String label, {
    String? prefix,
    String? hint,
  }) => InputDecoration(
    labelText: label,
    hintText: hint,
    prefixText: prefix,
    prefixStyle: const TextStyle(
      fontWeight: FontWeight.w700,
      color: Color(0xFF0D9488),
    ),
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
      width: 40,
      height: 4,
      margin: const EdgeInsets.only(top: 12, bottom: 20),
      decoration: BoxDecoration(
        color: Colors.grey.shade300,
        borderRadius: BorderRadius.circular(4),
      ),
    ),
  );

  static void _showSnack(
    BuildContext context,
    String message, {
    bool isError = false,
  }) {
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Row(
          children: [
            Icon(
              isError ? Icons.error_outline : Icons.check_circle_outline,
              color: Colors.white,
            ),
            const SizedBox(width: 12),
            Expanded(child: Text(message)),
          ],
        ),
        backgroundColor: isError ? Colors.redAccent : const Color(0xFF059669),
        behavior: SnackBarBehavior.floating,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
        margin: const EdgeInsets.all(16),
      ),
    );
  }

  // ── ADD ASSET MODAL ────────────────────────────────────────────────────────
  void _showAddAssetModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    final interestRateCtrl = TextEditingController();
    final purchasePriceCtrl = TextEditingController();
    String selectedType = 'bank';
    bool generatesIncome = false;
    String incomeFrequency = 'monthly';
    DateTime? purchaseDate;
    bool isLoading = false;

    bool isEmergency = false;
    int? yearsOfDeposit;

    final types = [
      {
        'value': 'bank',
        'label': 'Bank FD/PF',
        'icon': Icons.account_balance_outlined,
      },
      {'value': 'physical', 'label': 'Gold', 'icon': Icons.diamond_outlined},
      {
        'value': 'property',
        'label': 'Property',
        'icon': Icons.real_estate_agent_rounded,
      },
      {
        'value': 'equity',
        'label': 'Stocks/MF',
        'icon': Icons.trending_up_rounded,
      },
    ];

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Add Asset',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close),
                        onPressed: () => Navigator.pop(ctx),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Asset Name',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: nameCtrl,
                    textCapitalization: TextCapitalization.words,
                    decoration: _inputDec('e.g. SBI Savings, HDFC FD, Gold'),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Asset Type',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 10),
                  Row(
                    children: types.map((t) {
                      final sel = selectedType == t['value'];
                      return Expanded(
                        child: GestureDetector(
                          onTap: () => setState(
                            () => selectedType = t['value'] as String,
                          ),
                          child: AnimatedContainer(
                            duration: const Duration(milliseconds: 180),
                            margin: const EdgeInsets.only(right: 8),
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            decoration: BoxDecoration(
                              color: sel
                                  ? const Color(0xFF059669)
                                  : Colors.white,
                              borderRadius: BorderRadius.circular(14),
                              border: Border.all(
                                color: sel
                                    ? const Color(0xFF059669)
                                    : Colors.grey.shade200,
                              ),
                            ),
                            child: Column(
                              children: [
                                Icon(
                                  t['icon'] as IconData,
                                  color: sel ? Colors.white : Colors.grey,
                                  size: 22,
                                ),
                                const SizedBox(height: 6),
                                Text(
                                  t['label'] as String,
                                  style: TextStyle(
                                    fontSize: 11,
                                    fontWeight: FontWeight.w600,
                                    color: sel
                                        ? Colors.white
                                        : Colors.grey.shade600,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      );
                    }).toList(),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Current Value',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: amountCtrl,
                    keyboardType: const TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    decoration: _inputDec(
                      'Amount',
                      prefix: '₹ ',
                      hint: '100000',
                    ),
                  ),
                  const SizedBox(height: 20),

                  // Enhanced asset details
                  const Text(
                    'Interest Rate (%)',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: interestRateCtrl,
                    keyboardType: const TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    decoration: _inputDec(
                      'e.g. 7.5',
                      hint: 'Interest Rate (if applicable)',
                    ),
                  ),
                  const SizedBox(height: 16),

                  if (selectedType == 'bank' || selectedType == 'property') ...[
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Generates Income?',
                          style: TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 13,
                            color: Colors.grey,
                          ),
                        ),
                        Switch(
                          value: generatesIncome,
                          onChanged: (v) => setState(() => generatesIncome = v),
                          activeThumbColor: const Color(0xFF059669),
                        ),
                      ],
                    ),
                    if (generatesIncome) ...[
                      const SizedBox(height: 12),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text(
                            'Income Frequency',
                            style: TextStyle(
                              fontWeight: FontWeight.w600,
                              fontSize: 13,
                              color: Colors.grey,
                            ),
                          ),
                          DropdownButton<String>(
                            value: incomeFrequency,
                            underline: const SizedBox(),
                            items: ['monthly', 'quarterly', 'yearly']
                                .map(
                                  (f) => DropdownMenuItem(
                                    value: f,
                                    child: Text(
                                      f.toUpperCase(),
                                      style: const TextStyle(fontSize: 13),
                                    ),
                                  ),
                                )
                                .toList(),
                            onChanged: (v) =>
                                setState(() => incomeFrequency = v!),
                          ),
                        ],
                      ),
                    ],
                  ],
                  if (selectedType == 'bank') ...[
                    const Text(
                      'FD Date',
                      style: TextStyle(
                        fontWeight: FontWeight.w600,
                        fontSize: 13,
                        color: Colors.grey,
                      ),
                    ),
                    const SizedBox(height: 8),
                    Row(
                      children: [
                        Expanded(
                          child: TextField(
                            readOnly: true,
                            decoration: _inputDec('Start Date'),
                            controller: TextEditingController(
                              text: purchaseDate != null
                                  ? '${purchaseDate!.day}/${purchaseDate!.month}/${purchaseDate!.year}'
                                  : '',
                            ),
                            onTap: () async {
                              final d = await showDatePicker(
                                context: ctx,
                                initialDate: purchaseDate ?? DateTime.now(),
                                firstDate: DateTime(2000),
                                lastDate: DateTime.now(),
                              );
                              if (d != null) setState(() => purchaseDate = d);
                            },
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: TextField(
                            keyboardType: TextInputType.number,
                            decoration: _inputDec('Years of FD'),
                            onChanged: (v) => yearsOfDeposit = int.tryParse(v),
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 12),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        const Text(
                          'Mark as Emergency Fund',
                          style: TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 13,
                            color: Colors.grey,
                          ),
                        ),
                        Switch(
                          value: isEmergency,
                          onChanged: (v) => setState(() => isEmergency = v),
                          activeThumbColor: const Color(0xFF059669),
                        ),
                      ],
                    ),
                  ],
                  if (selectedType == 'physical' ||
                      selectedType == 'property') ...[
                    const SizedBox(height: 12),
                  ],
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () async {
                              final name = nameCtrl.text.trim();
                              final amount = double.tryParse(
                                amountCtrl.text.trim(),
                              );
                              if (name.isEmpty || amount == null || amount <= 0)
                                return;
                              final rate =
                                  double.tryParse(
                                    interestRateCtrl.text.trim(),
                                  ) ??
                                  0.0;
                              setState(() => isLoading = true);
                              try {
                                await Provider.of<FinancialProvider>(
                                  ctx,
                                  listen: false,
                                ).addAsset(
                                  name,
                                  selectedType,
                                  amount,
                                  interestRate: rate,
                                  isLiability: false,
                                  isEmergency: selectedType == 'bank'
                                      ? isEmergency
                                      : false,
                                  generatesIncome: generatesIncome,
                                  incomeFrequency: generatesIncome
                                      ? incomeFrequency
                                      : null,
                                  purchasePrice: double.tryParse(
                                    purchasePriceCtrl.text.trim(),
                                  ),
                                  purchaseDate: purchaseDate
                                      ?.toIso8601String()
                                      .split('T')
                                      .first,
                                  yearsOfDeposit: selectedType == 'bank'
                                      ? yearsOfDeposit
                                      : null,
                                );
                                if (ctx.mounted) {
                                  Navigator.pop(ctx);
                                  _showSnack(ctx, 'Asset added successfully!');
                                }
                              } catch (e) {
                                setState(() => isLoading = false);
                                if (ctx.mounted)
                                  _showSnack(
                                    ctx,
                                    'Failed to add asset: $e',
                                    isError: true,
                                  );
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF059669),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text(
                              'Add Asset',
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  // ── EDIT ASSET MODAL ───────────────────────────────────────────────────────
  void _showEditAssetModal(BuildContext context, dynamic asset) {
    final nameCtrl = TextEditingController(text: asset.name);
    final amountCtrl = TextEditingController(text: asset.amount.toString());
    final rateCtrl = TextEditingController(
      text: asset.interestRate > 0 ? asset.interestRate.toString() : '',
    );
    bool isLoading = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Edit Asset',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close),
                        onPressed: () => Navigator.pop(ctx),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Asset Name',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: nameCtrl,
                    decoration: _inputDec('e.g. SBI Savings'),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Current Value',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: amountCtrl,
                    keyboardType: TextInputType.number,
                    decoration: _inputDec('Amount', prefix: '₹ '),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Interest Rate (%)',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: rateCtrl,
                    keyboardType: TextInputType.number,
                    decoration: _inputDec('e.g. 7.5'),
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () async {
                              final name = nameCtrl.text.trim();
                              final amount = double.tryParse(
                                amountCtrl.text.trim(),
                              );
                              if (name.isEmpty || amount == null || amount <= 0)
                                return;
                              setState(() => isLoading = true);
                              try {
                                await Provider.of<FinancialProvider>(
                                  ctx,
                                  listen: false,
                                ).updateAsset(int.parse(asset.id), {
                                  'name': name,
                                  'type': 'bank',
                                  'amount': amount,
                                  'interest_rate':
                                      double.tryParse(rateCtrl.text.trim()) ??
                                      0,
                                  'is_liability': false,
                                  'is_emergency': false,
                                });
                                if (ctx.mounted) {
                                  Navigator.pop(ctx);
                                  _showSnack(ctx, 'Asset updated');
                                }
                              } catch (e) {
                                setState(() => isLoading = false);
                                if (ctx.mounted)
                                  _showSnack(ctx, 'Failed: $e', isError: true);
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF0D9488),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text(
                              'Update Asset',
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
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
      {
        'value': 'home_loan',
        'label': 'Home\nLoan',
        'icon': Icons.home_outlined,
      },
      {
        'value': 'car_loan',
        'label': 'Car\nLoan',
        'icon': Icons.directions_car_outlined,
      },
      {
        'value': 'personal_loan',
        'label': 'Personal\nLoan',
        'icon': Icons.person_outline,
      },
      {
        'value': 'credit_card',
        'label': 'Credit\nCard',
        'icon': Icons.credit_card_outlined,
      },
    ];

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Add Liability',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close),
                        onPressed: () => Navigator.pop(ctx),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Loan / Debt Name',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: nameCtrl,
                    textCapitalization: TextCapitalization.words,
                    decoration: _inputDec(
                      'e.g. SBI Home Loan, HDFC Credit Card',
                    ),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Type',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 10),
                  Row(
                    children: types.map((t) {
                      final sel = selectedType == t['value'];
                      return Expanded(
                        child: GestureDetector(
                          onTap: () => setState(
                            () => selectedType = t['value'] as String,
                          ),
                          child: AnimatedContainer(
                            duration: const Duration(milliseconds: 180),
                            margin: const EdgeInsets.only(right: 6),
                            padding: const EdgeInsets.symmetric(vertical: 10),
                            decoration: BoxDecoration(
                              color: sel
                                  ? const Color(0xFFDC2626)
                                  : Colors.white,
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: sel
                                    ? const Color(0xFFDC2626)
                                    : Colors.grey.shade200,
                              ),
                            ),
                            child: Column(
                              children: [
                                Icon(
                                  t['icon'] as IconData,
                                  color: sel ? Colors.white : Colors.grey,
                                  size: 20,
                                ),
                                const SizedBox(height: 4),
                                Text(
                                  t['label'] as String,
                                  textAlign: TextAlign.center,
                                  style: TextStyle(
                                    fontSize: 10,
                                    fontWeight: FontWeight.w600,
                                    color: sel
                                        ? Colors.white
                                        : Colors.grey.shade600,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ),
                      );
                    }).toList(),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Outstanding Amount',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: outstandingCtrl,
                    keyboardType: const TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    decoration: _inputDec(
                      'Total remaining balance',
                      prefix: '₹ ',
                    ),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Monthly EMI',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: emiCtrl,
                              keyboardType:
                                  const TextInputType.numberWithOptions(
                                    decimal: true,
                                  ),
                              decoration: _inputDec('EMI', prefix: '₹ '),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Interest Rate',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: rateCtrl,
                              keyboardType:
                                  const TextInputType.numberWithOptions(
                                    decimal: true,
                                  ),
                              decoration: _inputDec('Rate', hint: '8.5'),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () async {
                              final name = nameCtrl.text.trim();
                              final outstanding = double.tryParse(
                                outstandingCtrl.text.trim(),
                              );
                              if (name.isEmpty ||
                                  outstanding == null ||
                                  outstanding <= 0)
                                return;
                              final emi =
                                  double.tryParse(emiCtrl.text.trim()) ?? 0.0;
                              final rate =
                                  double.tryParse(rateCtrl.text.trim()) ?? 0.0;
                              setState(() => isLoading = true);
                              try {
                                await Provider.of<FinancialProvider>(
                                  ctx,
                                  listen: false,
                                ).addLiability(
                                  name,
                                  selectedType,
                                  outstanding,
                                  emi,
                                  rate,
                                );
                                if (ctx.mounted) {
                                  Navigator.pop(ctx);
                                  _showSnack(
                                    ctx,
                                    'Liability added successfully!',
                                  );
                                }
                              } catch (e) {
                                setState(() => isLoading = false);
                                if (ctx.mounted)
                                  _showSnack(
                                    ctx,
                                    'Failed to add liability: $e',
                                    isError: true,
                                  );
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFFDC2626),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text(
                              'Add Liability',
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  // ── ADD LENDING MODAL ──────────────────────────────────────────────────────
  void _showAddLendingModal(BuildContext context) {
    final personCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    String direction = 'lent';
    DateTime? dateGiven;
    DateTime? promisedReturnDate;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  const Text(
                    'Track Lending',
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 20),
                  Row(
                    children: [
                      Expanded(
                        child: GestureDetector(
                          onTap: () => setState(() => direction = 'lent'),
                          child: Container(
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            decoration: BoxDecoration(
                              color: direction == 'lent'
                                  ? const Color(0xFF0D9488)
                                  : Colors.white,
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: direction == 'lent'
                                    ? const Color(0xFF0D9488)
                                    : Colors.grey.shade300,
                              ),
                            ),
                            child: Text(
                              'I Lent',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                color: direction == 'lent'
                                    ? Colors.white
                                    : Colors.grey.shade700,
                              ),
                            ),
                          ),
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: GestureDetector(
                          onTap: () => setState(() => direction = 'borrowed'),
                          child: Container(
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            decoration: BoxDecoration(
                              color: direction == 'borrowed'
                                  ? const Color(0xFF0D9488)
                                  : Colors.white,
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: direction == 'borrowed'
                                    ? const Color(0xFF0D9488)
                                    : Colors.grey.shade300,
                              ),
                            ),
                            child: Text(
                              'I Borrowed',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                color: direction == 'borrowed'
                                    ? Colors.white
                                    : Colors.grey.shade700,
                              ),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  Text(
                    direction == 'lent' ? 'Lent to' : 'Borrowed from',
                    style: const TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: personCtrl,
                    decoration: _inputDec('Person Name'),
                  ),
                  const SizedBox(height: 20),
                  TextField(
                    controller: amountCtrl,
                    keyboardType: TextInputType.number,
                    decoration: _inputDec('Amount', prefix: '₹ '),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Date Given',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    readOnly: true,
                    decoration: _inputDec(
                      dateGiven != null
                          ? '${dateGiven!.day}/${dateGiven!.month}/${dateGiven!.year}'
                          : 'Select Date',
                    ),
                    onTap: () async {
                      final d = await showDatePicker(
                        context: ctx,
                        initialDate: DateTime.now(),
                        firstDate: DateTime(2000),
                        lastDate: DateTime.now(),
                      );
                      if (d != null) setState(() => dateGiven = d);
                    },
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Promised Return Date',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    readOnly: true,
                    decoration: _inputDec(
                      promisedReturnDate != null
                          ? '${promisedReturnDate!.day}/${promisedReturnDate!.month}/${promisedReturnDate!.year}'
                          : 'Select Date',
                    ),
                    onTap: () async {
                      final d = await showDatePicker(
                        context: ctx,
                        initialDate: DateTime.now().add(
                          const Duration(days: 30),
                        ),
                        firstDate: DateTime.now(),
                        lastDate: DateTime(2100),
                      );
                      if (d != null) setState(() => promisedReturnDate = d);
                    },
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: () async {
                        if (personCtrl.text.isEmpty ||
                            amountCtrl.text.isEmpty ||
                            dateGiven == null)
                          return;
                        await Provider.of<FinancialProvider>(
                          ctx,
                          listen: false,
                        ).addLendingRecord({
                          'direction': direction,
                          'person_name': personCtrl.text.trim(),
                          'amount':
                              double.tryParse(amountCtrl.text.trim()) ?? 0,
                          'date_given': dateGiven!
                              .toIso8601String()
                              .split('T')
                              .first,
                          'promised_return_date': promisedReturnDate
                              ?.toIso8601String()
                              .split('T')
                              .first,
                        });
                        if (ctx.mounted) {
                          Navigator.pop(ctx);
                          _showSnack(ctx, 'Record added');
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFFF59E0B),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: const Text(
                        'Add Record',
                        style: TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  void _showEditLendingModal(BuildContext context, dynamic record) {
    final personCtrl = TextEditingController(text: record.personName);
    final amountCtrl = TextEditingController(text: record.amount.toString());

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _handle(),
                const Text(
                  'Edit Record',
                  style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 20),
                TextField(
                  controller: personCtrl,
                  decoration: _inputDec('Person Name'),
                ),
                const SizedBox(height: 12),
                TextField(
                  controller: amountCtrl,
                  keyboardType: TextInputType.number,
                  decoration: _inputDec('Amount', prefix: '₹ '),
                ),
                const SizedBox(height: 32),
                SizedBox(
                  width: double.infinity,
                  height: 52,
                  child: ElevatedButton(
                    onPressed: () async {
                      await Provider.of<FinancialProvider>(
                        ctx,
                        listen: false,
                      ).updateLendingRecord(record.id, {
                        'person_name': personCtrl.text.trim(),
                        'amount': double.tryParse(amountCtrl.text.trim()) ?? 0,
                      });
                      if (ctx.mounted) {
                        Navigator.pop(ctx);
                        _showSnack(ctx, 'Record updated');
                      }
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: const Color(0xFF0D9488),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(14),
                      ),
                      elevation: 0,
                    ),
                    child: const Text(
                      'Update',
                      style: TextStyle(
                        color: Colors.white,
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
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

  // ── ADD INCOME MODAL ───────────────────────────────────────────────────────
  void _showAddIncomeModal(BuildContext context) {
    final nameCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    String selectedType = 'salary';
    bool isLoading = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Add Income',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close),
                        onPressed: () => Navigator.pop(ctx),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Label',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: nameCtrl,
                    textCapitalization: TextCapitalization.words,
                    decoration: _inputDec("e.g. My Salary, Partner Income"),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Monthly Amount',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: amountCtrl,
                    keyboardType: const TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    decoration: _inputDec(
                      'Amount',
                      prefix: '₹ ',
                      hint: '50000',
                    ),
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () async {
                              final name = nameCtrl.text.trim();
                              final amount = double.tryParse(
                                amountCtrl.text.trim(),
                              );
                              if (name.isEmpty || amount == null || amount <= 0)
                                return;
                              setState(() => isLoading = true);
                              try {
                                await Provider.of<FinancialProvider>(
                                  ctx,
                                  listen: false,
                                ).addIncome(
                                  name,
                                  selectedType,
                                  amount,
                                  'monthly',
                                );
                                if (ctx.mounted) {
                                  Navigator.pop(ctx);
                                  _showSnack(ctx, 'Income added successfully!');
                                }
                              } catch (e) {
                                setState(() => isLoading = false);
                                if (ctx.mounted)
                                  _showSnack(
                                    ctx,
                                    'Failed to add income: $e',
                                    isError: true,
                                  );
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF0D9488),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text(
                              'Add Income',
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  // ── SALARY BREAKUP MODAL ──────────────────────────────────────────────────
  void _showSalaryBreakupModal(BuildContext context, dynamic income) {
    final basicCtrl = TextEditingController();
    final hraCtrl = TextEditingController();
    final ltaCtrl = TextEditingController();
    final pfCtrl = TextEditingController();
    final specialCtrl = TextEditingController();
    final mealCtrl = TextEditingController();
    final variablePctCtrl = TextEditingController();
    final companyCtrl = TextEditingController();
    final fromYearCtrl = TextEditingController();

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'Salary Breakup — ${income.label}',
                        style: const TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close),
                        onPressed: () => Navigator.pop(ctx),
                      ),
                    ],
                  ),
                  const SizedBox(height: 8),
                  Container(
                    padding: const EdgeInsets.all(12),
                    decoration: BoxDecoration(
                      color: const Color(0xFF0D9488).withValues(alpha: 0.06),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Row(
                      children: [
                        const Icon(
                          Icons.info_outline,
                          color: Color(0xFF0D9488),
                          size: 16,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          'Total monthly: ₹${income.amount.toStringAsFixed(0)} · ${income.frequency}',
                          style: const TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 13,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Company Name',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 12,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 6),
                  TextField(
                    controller: companyCtrl,
                    decoration: _inputDec('e.g. ABC Corp'),
                  ),
                  const SizedBox(height: 12),
                  const Text(
                    'From Year',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 12,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 6),
                  TextField(
                    controller: fromYearCtrl,
                    keyboardType: TextInputType.number,
                    decoration: _inputDec('e.g. 2023'),
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Monthly Breakup',
                    style: TextStyle(
                      fontWeight: FontWeight.w700,
                      fontSize: 14,
                      color: Color(0xFF0F172A),
                    ),
                  ),
                  const SizedBox(height: 4),
                  const Text(
                    'Basic Pay is typically 50% of fixed pay',
                    style: TextStyle(fontSize: 11, color: Colors.grey),
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: basicCtrl,
                          keyboardType: TextInputType.number,
                          decoration: _inputDec('Basic Pay', prefix: '₹'),
                        ),
                      ),
                      const SizedBox(width: 10),
                      Expanded(
                        child: TextField(
                          controller: hraCtrl,
                          keyboardType: TextInputType.number,
                          decoration: _inputDec('HRA', prefix: '₹'),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: ltaCtrl,
                          keyboardType: TextInputType.number,
                          decoration: _inputDec('LTA', prefix: '₹'),
                        ),
                      ),
                      const SizedBox(width: 10),
                      Expanded(
                        child: TextField(
                          controller: pfCtrl,
                          keyboardType: TextInputType.number,
                          decoration: _inputDec('PF', prefix: '₹'),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: specialCtrl,
                          keyboardType: TextInputType.number,
                          decoration: _inputDec(
                            'Special Allowance',
                            prefix: '₹',
                          ),
                        ),
                      ),
                      const SizedBox(width: 10),
                      Expanded(
                        child: TextField(
                          controller: mealCtrl,
                          keyboardType: TextInputType.number,
                          decoration: _inputDec('Meal Card', prefix: '₹'),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: variablePctCtrl,
                          keyboardType: TextInputType.number,
                          decoration: _inputDec('Variable Pay %'),
                        ),
                      ),
                      const SizedBox(width: 10),
                      Expanded(child: Container()),
                    ],
                  ),
                  const SizedBox(height: 24),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: () async {
                        await Provider.of<FinancialProvider>(
                          ctx,
                          listen: false,
                        ).addSalaryDetail(income.id, {
                          'company_name': companyCtrl.text.trim(),
                          'from_year':
                              int.tryParse(fromYearCtrl.text.trim()) ??
                              DateTime.now().year,
                          'is_current': true,
                          'fixed_pay': income.amount * 12,
                          'basic_pay': double.tryParse(basicCtrl.text.trim()),
                          'hra': double.tryParse(hraCtrl.text.trim()),
                          'lta': double.tryParse(ltaCtrl.text.trim()),
                          'pf_employee': double.tryParse(pfCtrl.text.trim()),
                          'special_allowance': double.tryParse(
                            specialCtrl.text.trim(),
                          ),
                          'meal_card': double.tryParse(mealCtrl.text.trim()),
                          'variable_pay_percentage': double.tryParse(
                            variablePctCtrl.text.trim(),
                          ),
                        });
                        if (ctx.mounted) {
                          Navigator.pop(ctx);
                          _showSnack(ctx, 'Salary breakup saved');
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF0D9488),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: const Text(
                        'Save Breakup',
                        style: TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                          fontSize: 16,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
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
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Add Vehicle',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close),
                        onPressed: () => Navigator.pop(ctx),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Make & Model',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: makeModelCtrl,
                    textCapitalization: TextCapitalization.words,
                    decoration: _inputDec("e.g. Honda City, Royal Enfield"),
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Purchase Cost',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: costCtrl,
                    keyboardType: const TextInputType.numberWithOptions(
                      decimal: true,
                    ),
                    decoration: _inputDec(
                      'Amount',
                      prefix: '₹ ',
                      hint: '800000',
                    ),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Model Year',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: modelYearCtrl,
                              keyboardType: TextInputType.number,
                              decoration: _inputDec('e.g. 2022'),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Purchase Year',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: purchaseYearCtrl,
                              keyboardType: TextInputType.number,
                              decoration: _inputDec('e.g. 2023'),
                            ),
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
                            const Text(
                              'Fuel Type',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            DropdownButtonFormField<String>(
                              initialValue: selectedFuel,
                              decoration: _inputDec('Fuel'),
                              items: const [
                                DropdownMenuItem(
                                  value: 'petrol',
                                  child: Text('Petrol'),
                                ),
                                DropdownMenuItem(
                                  value: 'diesel',
                                  child: Text('Diesel'),
                                ),
                                DropdownMenuItem(
                                  value: 'electric',
                                  child: Text('Electric'),
                                ),
                                DropdownMenuItem(
                                  value: 'cng',
                                  child: Text('CNG'),
                                ),
                                DropdownMenuItem(
                                  value: 'hybrid',
                                  child: Text('Hybrid'),
                                ),
                              ],
                              onChanged: (val) {
                                if (val != null)
                                  setState(() => selectedFuel = val);
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
                            const Text(
                              'Mileage (Kmpl)',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: mileageCtrl,
                              keyboardType:
                                  const TextInputType.numberWithOptions(
                                    decimal: true,
                                  ),
                              decoration: _inputDec('e.g. 15.4'),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Registration Number',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: regCtrl,
                    textCapitalization: TextCapitalization.characters,
                    decoration: _inputDec('e.g. TN 01 AB 1234'),
                  ),
                  const SizedBox(height: 20),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Insurance Renewal',
                        style: TextStyle(
                          fontWeight: FontWeight.w600,
                          fontSize: 13,
                          color: Colors.grey,
                        ),
                      ),
                      TextButton(
                        onPressed: () async {
                          final d = await showDatePicker(
                            context: ctx,
                            initialDate: DateTime.now().add(
                              const Duration(days: 30),
                            ),
                            firstDate: DateTime.now(),
                            lastDate: DateTime.now().add(
                              const Duration(days: 3650),
                            ),
                          );
                          if (d != null) setState(() => insuranceDate = d);
                        },
                        child: Text(
                          insuranceDate == null
                              ? 'Select Date'
                              : '${insuranceDate!.day}/${insuranceDate!.month}/${insuranceDate!.year}',
                          style: const TextStyle(
                            fontWeight: FontWeight.bold,
                            fontSize: 14,
                          ),
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () async {
                              final makeModel = makeModelCtrl.text.trim();
                              final cost = double.tryParse(
                                costCtrl.text.trim(),
                              );
                              if (makeModel.isEmpty ||
                                  cost == null ||
                                  cost <= 0)
                                return;
                              setState(() => isLoading = true);
                              try {
                                await Provider.of<FinancialProvider>(
                                  ctx,
                                  listen: false,
                                ).addVehicle(
                                  makeModel,
                                  cost,
                                  insuranceRenewalDate: insuranceDate,
                                  modelYear: int.tryParse(
                                    modelYearCtrl.text.trim(),
                                  ),
                                  purchaseYear: int.tryParse(
                                    purchaseYearCtrl.text.trim(),
                                  ),
                                  fuelType: selectedFuel,
                                  mileageKmpl: double.tryParse(
                                    mileageCtrl.text.trim(),
                                  ),
                                  registrationNumber:
                                      regCtrl.text.trim().isEmpty
                                      ? null
                                      : regCtrl.text.trim(),
                                );
                                if (ctx.mounted) {
                                  Navigator.pop(ctx);
                                  _showSnack(
                                    ctx,
                                    'Vehicle added successfully!',
                                  );
                                }
                              } catch (e) {
                                setState(() => isLoading = false);
                                if (ctx.mounted)
                                  _showSnack(
                                    ctx,
                                    'Failed to add vehicle: $e',
                                    isError: true,
                                  );
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFFE88A1A),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text(
                              'Add Vehicle',
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
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
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      const Text(
                        'Add Electronic Device',
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      IconButton(
                        icon: const Icon(Icons.close),
                        onPressed: () => Navigator.pop(ctx),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Device Name',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: nameCtrl,
                    textCapitalization: TextCapitalization.words,
                    decoration: _inputDec("e.g. iPhone 15 Pro, MacBook Air"),
                  ),
                  const SizedBox(height: 16),
                  Row(
                    children: [
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Brand (Optional)',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: brandCtrl,
                              textCapitalization: TextCapitalization.words,
                              decoration: _inputDec('e.g. Apple'),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Model (Optional)',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: modelCtrl,
                              decoration: _inputDec('e.g. A3106'),
                            ),
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
                            const Text(
                              'Category',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            DropdownButtonFormField<String>(
                              initialValue: selectedCategory,
                              decoration: _inputDec('Category'),
                              items: const [
                                DropdownMenuItem(
                                  value: 'mobile',
                                  child: Text('Mobile'),
                                ),
                                DropdownMenuItem(
                                  value: 'laptop',
                                  child: Text('Laptop'),
                                ),
                                DropdownMenuItem(
                                  value: 'tablet',
                                  child: Text('Tablet'),
                                ),
                                DropdownMenuItem(
                                  value: 'tv',
                                  child: Text('TV / Display'),
                                ),
                                DropdownMenuItem(
                                  value: 'appliance',
                                  child: Text('Appliance'),
                                ),
                                DropdownMenuItem(
                                  value: 'other',
                                  child: Text('Other'),
                                ),
                              ],
                              onChanged: (val) {
                                if (val != null)
                                  setState(() => selectedCategory = val);
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
                            const Text(
                              'Purchase Cost',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: amountCtrl,
                              keyboardType:
                                  const TextInputType.numberWithOptions(
                                    decimal: true,
                                  ),
                              decoration: _inputDec('₹ amount'),
                            ),
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
                            const Text(
                              'Warranty (Years)',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: warrantyCtrl,
                              keyboardType: TextInputType.number,
                              decoration: _inputDec('e.g. 1'),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            const Text(
                              'Expected Life (Years)',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: expectedLifeCtrl,
                              keyboardType: TextInputType.number,
                              decoration: _inputDec('e.g. 3'),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Notes',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: notesCtrl,
                    maxLines: 2,
                    decoration: _inputDec('Additional details'),
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () async {
                              final name = nameCtrl.text.trim();
                              final amt = double.tryParse(
                                amountCtrl.text.trim(),
                              );
                              if (name.isEmpty || amt == null || amt <= 0)
                                return;
                              setState(() => isLoading = true);
                              try {
                                await Provider.of<FinancialProvider>(
                                  ctx,
                                  listen: false,
                                ).addElectronic({
                                  'name': name,
                                  'category': selectedCategory,
                                  'brand': brandCtrl.text.trim().isEmpty
                                      ? null
                                      : brandCtrl.text.trim(),
                                  'model': modelCtrl.text.trim().isEmpty
                                      ? null
                                      : modelCtrl.text.trim(),
                                  'purchase_date': DateTime.now()
                                      .toIso8601String()
                                      .substring(0, 10),
                                  'purchase_amount': amt,
                                  'warranty_years':
                                      int.tryParse(warrantyCtrl.text.trim()) ??
                                      1,
                                  'expected_life_years':
                                      int.tryParse(
                                        expectedLifeCtrl.text.trim(),
                                      ) ??
                                      3,
                                  'notes': notesCtrl.text.trim().isEmpty
                                      ? null
                                      : notesCtrl.text.trim(),
                                });
                                if (ctx.mounted) {
                                  Navigator.pop(ctx);
                                  _showSnack(ctx, 'Device added successfully!');
                                }
                              } catch (e) {
                                setState(() => isLoading = false);
                                if (ctx.mounted)
                                  _showSnack(
                                    ctx,
                                    'Failed to add device: $e',
                                    isError: true,
                                  );
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF0D9488),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text(
                              'Add Device',
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  // ── ADD GOLD MODAL ─────────────────────────────────────────────────────────
  void _showAddGoldModal(BuildContext context) {
    final gramsCtrl = TextEditingController();
    final priceCtrl = TextEditingController();
    int carat = 24;
    DateTime? purchaseDate;
    bool isLoading = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  const Text(
                    'Add Gold Holding',
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Purity (Carat)',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Row(
                    children: [24, 22, 18].map((c) {
                      final sel = carat == c;
                      return Expanded(
                        child: GestureDetector(
                          onTap: () => setState(() => carat = c),
                          child: Container(
                            margin: const EdgeInsets.only(right: 8),
                            padding: const EdgeInsets.symmetric(vertical: 12),
                            decoration: BoxDecoration(
                              color: sel
                                  ? const Color(0xFFF59E0B)
                                  : Colors.white,
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: sel
                                    ? const Color(0xFFF59E0B)
                                    : Colors.grey.shade300,
                              ),
                            ),
                            child: Text(
                              '${c}K',
                              textAlign: TextAlign.center,
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                                color: sel
                                    ? Colors.white
                                    : Colors.grey.shade700,
                              ),
                            ),
                          ),
                        ),
                      );
                    }).toList(),
                  ),
                  const SizedBox(height: 20),
                  TextField(
                    controller: gramsCtrl,
                    keyboardType: TextInputType.number,
                    decoration: _inputDec('Grams Purchased', hint: 'e.g. 10'),
                  ),
                  const SizedBox(height: 20),
                  TextField(
                    controller: priceCtrl,
                    keyboardType: TextInputType.number,
                    decoration: _inputDec(
                      'Purchase Price per Gram',
                      prefix: '₹ ',
                    ),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Purchase Date',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    readOnly: true,
                    decoration: _inputDec(
                      purchaseDate != null
                          ? '${purchaseDate!.day}/${purchaseDate!.month}/${purchaseDate!.year}'
                          : 'Select Date',
                    ),
                    onTap: () async {
                      final d = await showDatePicker(
                        context: ctx,
                        initialDate: DateTime.now(),
                        firstDate: DateTime(2000),
                        lastDate: DateTime.now(),
                      );
                      if (d != null) setState(() => purchaseDate = d);
                    },
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () async {
                              final grams = double.tryParse(
                                gramsCtrl.text.trim(),
                              );
                              final ppg = double.tryParse(
                                priceCtrl.text.trim(),
                              );
                              if (grams == null ||
                                  grams <= 0 ||
                                  ppg == null ||
                                  ppg <= 0)
                                return;
                              setState(() => isLoading = true);
                              try {
                                await Provider.of<FinancialProvider>(
                                  ctx,
                                  listen: false,
                                ).addGoldAsset({
                                  'carat': carat,
                                  'grams': grams,
                                  'purchase_price_per_gram': ppg,
                                  'purchase_date':
                                      purchaseDate
                                          ?.toIso8601String()
                                          .split('T')
                                          .first ??
                                      DateTime.now()
                                          .toIso8601String()
                                          .split('T')
                                          .first,
                                });
                                if (ctx.mounted) {
                                  Navigator.pop(ctx);
                                  _showSnack(
                                    ctx,
                                    'Gold added! Live value will be shown.',
                                  );
                                }
                              } catch (e) {
                                setState(() => isLoading = false);
                                if (ctx.mounted)
                                  _showSnack(ctx, 'Failed: $e', isError: true);
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFFF59E0B),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text(
                              'Add Gold',
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  // ── ADD STOCK MODAL ────────────────────────────────────────────────────────
  void _showAddStockModal(BuildContext context) {
    final tickerCtrl = TextEditingController();
    final qtyCtrl = TextEditingController();
    final priceCtrl = TextEditingController();
    DateTime? purchaseDate;
    bool isLoading = false;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (ctx) => StatefulBuilder(
        builder: (ctx, setState) => Padding(
          padding: EdgeInsets.only(
            bottom: MediaQuery.of(ctx).viewInsets.bottom,
          ),
          child: Container(
            decoration: _sheetDec(),
            padding: const EdgeInsets.fromLTRB(24, 0, 24, 32),
            child: SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _handle(),
                  const Text(
                    'Add Stock / MF',
                    style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Stock Ticker (e.g. RELIANCE.NS)',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    controller: tickerCtrl,
                    decoration: _inputDec('e.g. RELIANCE.NS, TCS.NS'),
                  ),
                  const SizedBox(height: 20),
                  Row(
                    children: [
                      Expanded(
                        child: Column(
                          children: [
                            const Text(
                              'Quantity',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: qtyCtrl,
                              keyboardType: TextInputType.number,
                              decoration: _inputDec('e.g. 10'),
                            ),
                          ],
                        ),
                      ),
                      const SizedBox(width: 16),
                      Expanded(
                        child: Column(
                          children: [
                            const Text(
                              'Avg Purchase Price',
                              style: TextStyle(
                                fontWeight: FontWeight.w600,
                                fontSize: 13,
                                color: Colors.grey,
                              ),
                            ),
                            const SizedBox(height: 8),
                            TextField(
                              controller: priceCtrl,
                              keyboardType: TextInputType.number,
                              decoration: _inputDec('Price', prefix: '₹ '),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'Purchase Date',
                    style: TextStyle(
                      fontWeight: FontWeight.w600,
                      fontSize: 13,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 8),
                  TextField(
                    readOnly: true,
                    decoration: _inputDec(
                      purchaseDate != null
                          ? '${purchaseDate!.day}/${purchaseDate!.month}/${purchaseDate!.year}'
                          : 'Select Date',
                    ),
                    onTap: () async {
                      final d = await showDatePicker(
                        context: ctx,
                        initialDate: DateTime.now().subtract(
                          const Duration(days: 365),
                        ),
                        firstDate: DateTime(2000),
                        lastDate: DateTime.now(),
                      );
                      if (d != null) setState(() => purchaseDate = d);
                    },
                  ),
                  const SizedBox(height: 32),
                  SizedBox(
                    width: double.infinity,
                    height: 52,
                    child: ElevatedButton(
                      onPressed: isLoading
                          ? null
                          : () async {
                              final ticker = tickerCtrl.text
                                  .trim()
                                  .toUpperCase();
                              final qty = int.tryParse(qtyCtrl.text.trim());
                              final price = double.tryParse(
                                priceCtrl.text.trim(),
                              );
                              if (ticker.isEmpty || qty == null || qty <= 0)
                                return;
                              setState(() => isLoading = true);
                              try {
                                await Provider.of<FinancialProvider>(
                                  ctx,
                                  listen: false,
                                ).addStockHolding({
                                  'ticker': ticker,
                                  'quantity': qty,
                                  'avg_purchase_price': price ?? 0,
                                  'total_invested': (price ?? 0) * qty,
                                  'purchase_date':
                                      purchaseDate
                                          ?.toIso8601String()
                                          .split('T')
                                          .first ??
                                      DateTime.now()
                                          .toIso8601String()
                                          .split('T')
                                          .first,
                                });
                                if (ctx.mounted) {
                                  Navigator.pop(ctx);
                                  _showSnack(
                                    ctx,
                                    'Stock added! Live price will be shown.',
                                  );
                                }
                              } catch (e) {
                                setState(() => isLoading = false);
                                if (ctx.mounted)
                                  _showSnack(ctx, 'Failed: $e', isError: true);
                              }
                            },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF6B46C1),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                        elevation: 0,
                      ),
                      child: isLoading
                          ? const SizedBox(
                              height: 22,
                              width: 22,
                              child: CircularProgressIndicator(
                                color: Colors.white,
                                strokeWidth: 2,
                              ),
                            )
                          : const Text(
                              'Add Stock',
                              style: TextStyle(
                                color: Colors.white,
                                fontWeight: FontWeight.bold,
                                fontSize: 16,
                              ),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  // ── Section Header ─────────────────────────────────────────────────────────
  Widget _buildSectionHeader(
    String title,
    List<Color> gradient,
    VoidCallback onAdd,
  ) {
    return Padding(
      padding: const EdgeInsets.only(top: 24.0, bottom: 16.0),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            title,
            style: const TextStyle(
              fontSize: 19,
              fontWeight: FontWeight.w900,
              letterSpacing: -0.5,
              color: Color(0xFF0F172A),
            ),
          ),
          GestureDetector(
            onTap: onAdd,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 7),
              decoration: BoxDecoration(
                color: const Color(0xFF0D9488).withValues(alpha: 0.08),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(
                  color: const Color(0xFF0D9488).withValues(alpha: 0.2),
                  width: 1.2,
                ),
              ),
              child: const Row(
                children: [
                  Icon(Icons.add_rounded, size: 16, color: Color(0xFF0D9488)),
                  SizedBox(width: 4),
                  Text(
                    'Add',
                    style: TextStyle(
                      color: Color(0xFF0D9488),
                      fontWeight: FontWeight.w800,
                      fontSize: 12,
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

  // ── Portfolio Summary ──────────────────────────────────────────────────────
  List<Widget> _buildPortfolioSummary(FinancialProvider p) {
    final ps = p.portfolioSummary;
    if (ps == null || (ps['totalInvested'] ?? 0) == 0) return [];

    final totalInvested = (ps['totalInvested'] ?? 0).toDouble();
    final totalCurrent = (ps['totalCurrent'] ?? 0).toDouble();
    final totalReturn = (ps['totalReturn'] ?? 0).toDouble();
    final estimatedCagr = (ps['estimatedCagr'] ?? 0).toDouble();
    final returnColor = totalReturn >= 0
        ? const Color(0xFF34D399)
        : const Color(0xFFF87171);
    final returnSign = totalReturn >= 0 ? '+' : '';
    final returnPct = totalInvested > 0
        ? (totalReturn / totalInvested * 100)
        : 0.0;

    return [
      const SizedBox(height: 8),
      ClipRRect(
        borderRadius: BorderRadius.circular(24),
        child: BackdropFilter(
          filter: ui.ImageFilter.blur(sigmaX: 25, sigmaY: 25),
          child: Container(
            padding: const EdgeInsets.all(20),
            decoration: BoxDecoration(
              color: Colors.white.withValues(alpha: 0.60),
              borderRadius: BorderRadius.circular(24),
              border: Border.all(
                color: Colors.white.withValues(alpha: 0.7),
                width: 1.2,
              ),
              boxShadow: [
                BoxShadow(
                  color: Colors.black.withValues(alpha: 0.04),
                  blurRadius: 25,
                  offset: const Offset(0, 6),
                ),
              ],
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Container(
                      padding: const EdgeInsets.all(8),
                      decoration: BoxDecoration(
                        color: const Color(0xFF0D9488).withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(10),
                      ),
                      child: const Icon(
                        Icons.trending_up_rounded,
                        color: Color(0xFF0D9488),
                        size: 18,
                      ),
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Text(
                        'Portfolio Summary',
                        style: const TextStyle(
                          fontSize: 15,
                          fontWeight: FontWeight.w800,
                          color: Color(0xFF0F172A),
                        ),
                      ),
                    ),
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 10,
                        vertical: 4,
                      ),
                      decoration: BoxDecoration(
                        color: returnColor.withValues(alpha: 0.12),
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: Text(
                        '$returnSign${returnPct.toStringAsFixed(1)}%',
                        style: TextStyle(
                          fontWeight: FontWeight.w800,
                          fontSize: 12,
                          color: returnColor,
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 20),
                Row(
                  children: [
                    _psColumn(
                      'Invested',
                      '₹${_fmt(totalInvested)}',
                      const Color(0xFF64748B),
                    ),
                    _psColumn(
                      'Current',
                      '₹${_fmt(totalCurrent)}',
                      const Color(0xFF0F172A),
                    ),
                    _psColumn(
                      'Return',
                      '$returnSign₹${_fmt(totalReturn.abs())}',
                      returnColor,
                    ),
                    _psColumn(
                      'CAGR',
                      '${estimatedCagr.toStringAsFixed(1)}%',
                      const Color(0xFF909AC6),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    ];
  }

  Widget _psColumn(String label, String value, Color color) {
    return Expanded(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: const TextStyle(
              fontSize: 10,
              color: Color(0xFF94A3B8),
              fontWeight: FontWeight.w600,
            ),
          ),
          const SizedBox(height: 2),
          Text(
            value,
            style: TextStyle(
              fontSize: 13,
              fontWeight: FontWeight.w900,
              color: color,
            ),
          ),
        ],
      ),
    );
  }

  // ── Item card with delete swipe ────────────────────────────────────────────
  Widget _buildListItem(
    String title,
    double amount,
    String subtitle,
    List<Color> gradient,
    IconData icon,
    VoidCallback onDelete, {
    VoidCallback? onTap,
    VoidCallback? onEdit,
  }) {
    return ClipRRect(
      borderRadius: BorderRadius.circular(18),
      child: BackdropFilter(
        filter: ui.ImageFilter.blur(sigmaX: 25, sigmaY: 25),
        child: Container(
          margin: const EdgeInsets.only(bottom: 12),
          decoration: BoxDecoration(
            color: Colors.white.withValues(alpha: 0.55),
            borderRadius: BorderRadius.circular(18),
            border: Border.all(
              color: Colors.white.withValues(alpha: 0.6),
              width: 1.2,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.03),
                blurRadius: 15,
                offset: const Offset(0, 4),
              ),
            ],
          ),
          child: Material(
            color: Colors.transparent,
            child: InkWell(
              borderRadius: BorderRadius.circular(18),
              onTap: onTap,
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Expanded(
                      child: Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.all(10),
                            decoration: BoxDecoration(
                              color: gradient[0].withValues(alpha: 0.1),
                              shape: BoxShape.circle,
                            ),
                            child: Icon(icon, color: gradient[0], size: 20),
                          ),
                          const SizedBox(width: 12),
                          Expanded(
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                Text(
                                  title,
                                  style: const TextStyle(
                                    fontWeight: FontWeight.w800,
                                    fontSize: 14,
                                    color: Color(0xFF0F172A),
                                  ),
                                  overflow: TextOverflow.ellipsis,
                                ),
                                const SizedBox(height: 3),
                                Text(
                                  subtitle,
                                  style: const TextStyle(
                                    color: Color(0xFF64748B),
                                    fontSize: 11,
                                    fontWeight: FontWeight.w500,
                                  ),
                                  overflow: TextOverflow.ellipsis,
                                ),
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                    const SizedBox(width: 8),
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          '₹${amount.toStringAsFixed(0)}',
                          style: const TextStyle(
                            fontWeight: FontWeight.w900,
                            fontSize: 15,
                            color: Color(0xFF0F172A),
                          ),
                        ),
                        if (onEdit != null) ...[
                          const SizedBox(width: 4),
                          GestureDetector(
                            onTap: onEdit,
                            child: Container(
                              padding: const EdgeInsets.all(6),
                              decoration: BoxDecoration(
                                color: const Color(
                                  0xFF0D9488,
                                ).withValues(alpha: 0.1),
                                borderRadius: BorderRadius.circular(8),
                              ),
                              child: const Icon(
                                Icons.edit_outlined,
                                color: Color(0xFF0D9488),
                                size: 16,
                              ),
                            ),
                          ),
                        ],
                        const SizedBox(width: 4),
                        GestureDetector(
                          onTap: onDelete,
                          child: Container(
                            padding: const EdgeInsets.all(6),
                            decoration: BoxDecoration(
                              color: Colors.red.withValues(alpha: 0.08),
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: const Icon(
                              Icons.delete_outline_rounded,
                              color: Color(0xFFF87171),
                              size: 16,
                            ),
                          ),
                        ),
                      ],
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

  // ── BUILD ──────────────────────────────────────────────────────────────────
  @override
  Widget build(BuildContext context) {
    final provider = context.watch<FinancialProvider>();
    final size = MediaQuery.of(context).size;

    return Scaffold(
      body: Stack(
        children: [
          // ── Light gradient background ────────────────────────────────────
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
          // ── Accent blobs ────────────────────────────────────────────────
          Positioned(
            top: -size.height * 0.12,
            right: -size.width * 0.2,
            child: Container(
              width: size.width * 0.7,
              height: size.width * 0.7,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFF0D9488).withValues(alpha: 0.15),
                    const Color(0xFF0D9488).withValues(alpha: 0.04),
                    const Color(0xFF0D9488).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),
          Positioned(
            bottom: -size.height * 0.08,
            left: -size.width * 0.15,
            child: Container(
              width: size.width * 0.55,
              height: size.width * 0.55,
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                gradient: RadialGradient(
                  colors: [
                    const Color(0xFF909AC6).withValues(alpha: 0.12),
                    const Color(0xFF909AC6).withValues(alpha: 0.03),
                    const Color(0xFF909AC6).withValues(alpha: 0.0),
                  ],
                ),
              ),
            ),
          ),
          Positioned.fill(child: CustomPaint(painter: _GridPainter())),

          // ── Content ─────────────────────────────────────────────────────────
          SafeArea(
            child: provider.isLoading
                ? const Center(
                    child: CircularProgressIndicator(color: Color(0xFF0D9488)),
                  )
                : RefreshIndicator(
                    onRefresh: provider.loadAllData,
                    color: const Color(0xFF0D9488),
                    child: ListView(
                      physics: const AlwaysScrollableScrollPhysics(),
                      padding: const EdgeInsets.symmetric(
                        horizontal: 20,
                        vertical: 8,
                      ),
                      children: [
                        // Portfolio Summary
                        ..._buildPortfolioSummary(provider),

                        // Assets
                        _buildSectionHeader('Assets', const [
                          Color(0xFF059669),
                          Color(0xFF34D399),
                        ], () => _showAddAssetModal(context)),
                        if (provider.assets.isEmpty)
                          UiUtils.buildEmptyState(
                            'No assets yet',
                            'Tap Add to add your savings,\nFDs, gold, property etc.',
                            Icons.account_balance_wallet_outlined,
                            Colors.green,
                          ),
                        ...provider.assets.map(
                          (a) => _buildListItem(
                            a.name,
                            a.amount,
                            a.purchasePrice != null
                                ? 'Invested: ₹${a.purchasePrice!.toStringAsFixed(0)}'
                                : 'Asset',
                            const [Color(0xFF059669), Color(0xFF34D399)],
                            Icons.account_balance_wallet_rounded,
                            () => UiUtils.showDeleteBottomSheet(
                              context,
                              a.name,
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteAsset(a.id);
                                if (context.mounted)
                                  _showSnack(context, 'Asset deleted');
                              },
                            ),
                            onEdit: () => _showEditAssetModal(context, a),
                          ),
                        ),

                        // Gold Assets
                        _buildSectionHeader('Gold', const [
                          Color(0xFFF59E0B),
                          Color(0xFFFBBF24),
                        ], () => _showAddGoldModal(context)),
                        if (provider.goldAssets.isEmpty)
                          UiUtils.buildEmptyState(
                            'No gold holdings',
                            'Tap Add to track gold\nwith live value.',
                            Icons.monetization_on_outlined,
                            Colors.amber,
                          ),
                        ...provider.goldAssets.map(
                          (g) => _buildListItem(
                            '${g.grams} g · ${g.caratLabel}',
                            g.currentValue ?? g.totalPurchaseCost ?? 0,
                            'Invested: ₹${g.totalPurchaseCost?.toStringAsFixed(0) ?? '0'} · ${g.returnPct != null ? '${g.returnPct! >= 0 ? '+' : ''}${g.returnPct!.toStringAsFixed(1)}%' : '--'}',
                            const [Color(0xFFF59E0B), Color(0xFFFBBF24)],
                            Icons.monetization_on_rounded,
                            () => UiUtils.showDeleteBottomSheet(
                              context,
                              '${g.grams}g Gold',
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteGoldAsset(g.id);
                                if (context.mounted)
                                  _showSnack(context, 'Gold holding removed');
                              },
                            ),
                          ),
                        ),
                        const SizedBox(height: 16),

                        // Stock Holdings
                        _buildSectionHeader('Stocks / MF', const [
                          Color(0xFF6B46C1),
                          Color(0xFF8B5CF6),
                        ], () => _showAddStockModal(context)),
                        if (provider.stockHoldings.isEmpty)
                          UiUtils.buildEmptyState(
                            'No stock holdings',
                            'Tap Add to track stocks\nwith live prices.',
                            Icons.trending_up_outlined,
                            Colors.purple,
                          ),
                        ...provider.stockHoldings.map(
                          (s) => _buildListItem(
                            s.ticker,
                            s.currentValue ?? s.totalInvested,
                            'Qty: ${s.quantity} · Avg: ₹${s.avgPurchasePrice.toStringAsFixed(0)} · ${s.returnPct != null ? '${s.returnPct! >= 0 ? '+' : ''}${s.returnPct!.toStringAsFixed(1)}%' : '--'}',
                            const [Color(0xFF6B46C1), Color(0xFF8B5CF6)],
                            Icons.trending_up_rounded,
                            () => UiUtils.showDeleteBottomSheet(
                              context,
                              s.ticker,
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteStockHolding(s.id);
                                if (context.mounted)
                                  _showSnack(context, 'Stock holding removed');
                              },
                            ),
                          ),
                        ),
                        const SizedBox(height: 16),

                        // Lending Records
                        _buildSectionHeader('Money Lending', const [
                          Color(0xFFF59E0B),
                          Color(0xFFFBBF24),
                        ], () => _showAddLendingModal(context)),
                        if (provider.lendingRecords.isEmpty)
                          UiUtils.buildEmptyState(
                            'No lending records',
                            'Tap Add to track money lent\nor borrowed.',
                            Icons.currency_rupee_outlined,
                            Colors.amber,
                          ),
                        ...provider.lendingRecords.map(
                          (l) => _buildListItem(
                            l.personName,
                            l.remaining,
                            '${l.direction == 'lent' ? 'Lent' : 'Borrowed'} · ${l.status} · Due: ${l.promisedReturnDate ?? 'N/A'}',
                            const [Color(0xFFF59E0B), Color(0xFFFBBF24)],
                            l.direction == 'lent'
                                ? Icons.arrow_upward_rounded
                                : Icons.arrow_downward_rounded,
                            () => UiUtils.showDeleteBottomSheet(
                              context,
                              l.personName,
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteLendingRecord(l.id);
                                if (context.mounted)
                                  _showSnack(context, 'Record deleted');
                              },
                            ),
                            onEdit: () => _showEditLendingModal(context, l),
                          ),
                        ),
                        const SizedBox(height: 16),

                        // Liabilities
                        _buildSectionHeader('Liabilities', const [
                          Color(0xFFDC2626),
                          Color(0xFFF87171),
                        ], () => _showAddLiabilityModal(context)),
                        if (provider.liabilities.isEmpty)
                          UiUtils.buildEmptyState(
                            'No liabilities',
                            'Tap Add to track loans,\ncredit cards, EMIs etc.',
                            Icons.credit_card_outlined,
                            Colors.red,
                          ),
                        ...provider.liabilities.map(
                          (l) => _buildListItem(
                            l.name,
                            l.amount,
                            '${l.interestRate.toStringAsFixed(1)}% interest',
                            const [Color(0xFFDC2626), Color(0xFFF87171)],
                            Icons.credit_card_rounded,
                            () => UiUtils.showDeleteBottomSheet(
                              context,
                              l.name,
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteLiability(l.id);
                                if (context.mounted)
                                  _showSnack(context, 'Liability deleted');
                              },
                            ),
                          ),
                        ),

                        // Incomes
                        _buildSectionHeader('Income Sources', const [
                          Color(0xFF0D9488),
                          Color(0xFF2DD4BF),
                        ], () => _showAddIncomeModal(context)),
                        if (provider.incomes.isEmpty)
                          UiUtils.buildEmptyState(
                            'No income sources',
                            'Tap Add to record your salary\nand other income.',
                            Icons.trending_up_outlined,
                            Colors.blue,
                          ),
                        ...provider.incomes.map((i) {
                          final isSalary = i.type == 'salary';
                          return _buildListItem(
                            i.label,
                            i.amount,
                            isSalary
                                ? 'Salary · Tap for breakup & growth'
                                : '${i.frequency} income',
                            const [Color(0xFF0D9488), Color(0xFF2DD4BF)],
                            isSalary
                                ? Icons.work_outline_rounded
                                : Icons.trending_up_rounded,
                            () => UiUtils.showDeleteBottomSheet(
                              context,
                              i.label,
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteIncome(i.id);
                                if (context.mounted)
                                  _showSnack(context, 'Income deleted');
                              },
                            ),
                            onEdit: isSalary
                                ? () => _showSalaryBreakupModal(context, i)
                                : null,
                            onTap: isSalary
                                ? () => _showSalaryBreakupModal(context, i)
                                : null,
                          );
                        }),

                        // Vehicles
                        _buildSectionHeader('Vehicles', const [
                          Color(0xFFE88A1A),
                          Color(0xFFFBBF24),
                        ], () => _showAddVehicleModal(context)),
                        if (provider.vehicles.isEmpty)
                          UiUtils.buildEmptyState(
                            'No vehicles',
                            'Tap Add to track your cars/bikes\nand insurance.',
                            Icons.directions_car_outlined,
                            Colors.orange,
                          ),
                        ...provider.vehicles.map(
                          (v) => _buildListItem(
                            v.makeModel,
                            v.purchaseCost,
                            v.registrationNumber != null &&
                                    v.registrationNumber!.isNotEmpty
                                ? '${v.registrationNumber} • Cost/km: ₹${v.costPerKm.toStringAsFixed(2)}'
                                : 'Cost/km: ₹${v.costPerKm.toStringAsFixed(2)}',
                            const [Color(0xFFE88A1A), Color(0xFFFBBF24)],
                            Icons.directions_car_rounded,
                            () => UiUtils.showDeleteBottomSheet(
                              context,
                              v.makeModel,
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteVehicle(v.id);
                                if (context.mounted)
                                  _showSnack(context, 'Vehicle deleted');
                              },
                            ),
                            onTap: () {
                              Navigator.push(
                                context,
                                MaterialPageRoute(
                                  builder: (_) =>
                                      VehicleDetailScreen(vehicle: v),
                                ),
                              );
                            },
                          ),
                        ),

                        // Electronics
                        _buildSectionHeader('Electronics', const [
                          Color(0xFF0D9488),
                          Color(0xFF2DD4BF),
                        ], () => _showAddElectronicModal(context)),
                        if (provider.electronics.isEmpty)
                          UiUtils.buildEmptyState(
                            'No devices',
                            'Tap Add to track phones, laptops\nand warranty.',
                            Icons.phone_android_outlined,
                            Colors.teal,
                          ),
                        ...provider.electronics.map((e) {
                          IconData icon = Icons.devices_other_rounded;
                          if (e.category.toLowerCase() == 'mobile')
                            icon = Icons.phone_android_rounded;
                          if (e.category.toLowerCase() == 'laptop')
                            icon = Icons.laptop_chromebook_rounded;
                          if (e.category.toLowerCase() == 'tablet')
                            icon = Icons.tablet_mac_rounded;
                          if (e.category.toLowerCase() == 'tv')
                            icon = Icons.tv_rounded;
                          if (e.category.toLowerCase() == 'appliance')
                            icon = Icons.kitchen_rounded;

                          String subtitle =
                              'Warranty: ${e.warrantyStatus.replaceAll('_', ' ').toUpperCase()}';
                          if (e.warrantyStatus != 'expired') {
                            subtitle +=
                                ' (${e.warrantyDaysRemaining}d remaining)';
                          }

                          return _buildListItem(
                            e.name,
                            e.currentValue,
                            subtitle,
                            const [Color(0xFF0D9488), Color(0xFF2DD4BF)],
                            icon,
                            () => UiUtils.showDeleteBottomSheet(
                              context,
                              e.name,
                              () async {
                                await Provider.of<FinancialProvider>(
                                  context,
                                  listen: false,
                                ).deleteElectronic(e.id);
                                if (context.mounted)
                                  _showSnack(context, 'Device deleted');
                              },
                            ),
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

  String _fmt(double v) =>
      v >= 1000 ? '${(v / 1000).toStringAsFixed(0)}K' : v.toStringAsFixed(0);
}

class _GridPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    final paint = Paint()
      ..color = const Color(0xFF0D9488).withValues(alpha: 0.035)
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
