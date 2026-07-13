import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class VehicleFuelForm extends StatefulWidget {
  final String vehicleId;

  const VehicleFuelForm({super.key, required this.vehicleId});

  @override
  State<VehicleFuelForm> createState() => _VehicleFuelFormState();
}

class _VehicleFuelFormState extends State<VehicleFuelForm> {
  final _formKey = GlobalKey<FormState>();

  late TextEditingController _dateCtrl;
  late TextEditingController _amountCtrl;
  late TextEditingController _litersCtrl;
  late TextEditingController _kmCtrl;
  late TextEditingController _priceCtrl;

  bool _isFullTank = true;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _dateCtrl = TextEditingController(
      text: DateTime.now().toIso8601String().substring(0, 10),
    );
    _amountCtrl = TextEditingController();
    _litersCtrl = TextEditingController();
    _kmCtrl = TextEditingController();
    _priceCtrl = TextEditingController();

    _amountCtrl.addListener(_recomputePricePerLiter);
    _litersCtrl.addListener(_recomputePricePerLiter);
  }

  @override
  void dispose() {
    _dateCtrl.dispose();
    _amountCtrl.dispose();
    _litersCtrl.dispose();
    _kmCtrl.dispose();
    _priceCtrl.dispose();
    super.dispose();
  }

  void _recomputePricePerLiter() {
    final amount = double.tryParse(_amountCtrl.text.trim()) ?? 0.0;
    final liters = double.tryParse(_litersCtrl.text.trim()) ?? 0.0;
    if (amount > 0.0 && liters > 0.0) {
      final price = amount / liters;
      _priceCtrl.text = price.toStringAsFixed(2);
    }
  }

  Future<void> _selectDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: DateTime.now(),
      firstDate: DateTime(2000),
      lastDate: DateTime(2100),
    );
    if (picked != null) {
      setState(() {
        _dateCtrl.text = picked.toIso8601String().substring(0, 10);
      });
    }
  }

  Future<void> _submit() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    final amount = double.tryParse(_amountCtrl.text.trim()) ?? 0.0;
    final liters = double.tryParse(_litersCtrl.text.trim()) ?? 0.0;
    final price = double.tryParse(_priceCtrl.text.trim()) ?? (amount / liters);

    final body = {
      'fill_date': _dateCtrl.text.trim(),
      'amount': amount,
      'liters': liters,
      'km_at_fill': double.tryParse(_kmCtrl.text.trim()),
      'price_per_liter': price,
      'is_full_tank': _isFullTank,
    };

    try {
      await Provider.of<FinancialProvider>(
        context,
        listen: false,
      ).addFuelRecord(widget.vehicleId, body);
      if (mounted) {
        UiUtils.showSnack(context, 'Fuel log added successfully');
        Navigator.pop(context);
      }
    } catch (e) {
      if (mounted)
        UiUtils.showSnack(context, 'Failed to add fuel log: $e', isError: true);
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Add Fuel Fill Log'),
        actions: [
          if (_isLoading)
            const Center(
              child: Padding(
                padding: EdgeInsets.symmetric(horizontal: 16),
                child: SizedBox(
                  width: 20,
                  height: 20,
                  child: CircularProgressIndicator(
                    color: Colors.white,
                    strokeWidth: 2,
                  ),
                ),
              ),
            )
          else
            IconButton(icon: const Icon(Icons.check), onPressed: _submit),
        ],
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(24),
        child: Form(
          key: _formKey,
          child: Column(
            children: [
              TextFormField(
                controller: _dateCtrl,
                readOnly: true,
                onTap: _selectDate,
                decoration: const InputDecoration(
                  labelText: 'Fill Date',
                  suffixIcon: Icon(Icons.calendar_today),
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _amountCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Total Amount Paid (₹)',
                  border: OutlineInputBorder(),
                ),
                validator: (val) =>
                    val == null || val.isEmpty ? 'Amount is required' : null,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _litersCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Liters Filled',
                  border: OutlineInputBorder(),
                ),
                validator: (val) => val == null || val.isEmpty
                    ? 'Liters filled is required'
                    : null,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _priceCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Price per Liter (₹/L - Optional)',
                  hintText: 'Auto-computed',
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _kmCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Odometer Reading (KM - Optional)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              SwitchListTile(
                title: const Text('Is Full Tank fill?'),
                activeThumbColor: const Color(0xFF6366F1),
                value: _isFullTank,
                onChanged: (val) {
                  setState(() => _isFullTank = val);
                },
              ),
              const SizedBox(height: 40),
            ],
          ),
        ),
      ),
    );
  }
}
