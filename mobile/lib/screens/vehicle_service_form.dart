import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';

class VehicleServiceForm extends StatefulWidget {
  final String vehicleId;

  const VehicleServiceForm({super.key, required this.vehicleId});

  @override
  State<VehicleServiceForm> createState() => _VehicleServiceFormState();
}

class _VehicleServiceFormState extends State<VehicleServiceForm> {
  final _formKey = GlobalKey<FormState>();

  late TextEditingController _dateCtrl;
  late TextEditingController _costCtrl;
  late TextEditingController _kmCtrl;
  late TextEditingController _centerCtrl;
  late TextEditingController _nextKmCtrl;
  late TextEditingController _descCtrl;

  String _selectedType = 'regular';
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _dateCtrl = TextEditingController(
      text: DateTime.now().toIso8601String().substring(0, 10),
    );
    _costCtrl = TextEditingController();
    _kmCtrl = TextEditingController();
    _centerCtrl = TextEditingController();
    _nextKmCtrl = TextEditingController();
    _descCtrl = TextEditingController();
  }

  @override
  void dispose() {
    _dateCtrl.dispose();
    _costCtrl.dispose();
    _kmCtrl.dispose();
    _centerCtrl.dispose();
    _nextKmCtrl.dispose();
    _descCtrl.dispose();
    super.dispose();
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

    final body = {
      'service_date': _dateCtrl.text.trim(),
      'service_type': _selectedType,
      'cost': double.tryParse(_costCtrl.text.trim()) ?? 0.0,
      'km_at_service': double.tryParse(_kmCtrl.text.trim()),
      'service_center': _centerCtrl.text.trim().isEmpty
          ? null
          : _centerCtrl.text.trim(),
      'next_service_km': double.tryParse(_nextKmCtrl.text.trim()),
      'description': _descCtrl.text.trim().isEmpty
          ? null
          : _descCtrl.text.trim(),
    };

    try {
      await Provider.of<FinancialProvider>(
        context,
        listen: false,
      ).addServiceRecord(widget.vehicleId, body);
      if (mounted) {
        UiUtils.showSnack(context, 'Service record added successfully');
        Navigator.pop(context);
      }
    } catch (e) {
      if (mounted)
        UiUtils.showSnack(
          context,
          'Failed to add service record: $e',
          isError: true,
        );
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Add Service Record'),
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
                  labelText: 'Service Date',
                  suffixIcon: Icon(Icons.calendar_today),
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              DropdownButtonFormField<String>(
                initialValue: _selectedType,
                decoration: const InputDecoration(
                  labelText: 'Service Type',
                  border: OutlineInputBorder(),
                ),
                items: const [
                  DropdownMenuItem(
                    value: 'regular',
                    child: Text('Regular Service'),
                  ),
                  DropdownMenuItem(
                    value: 'major',
                    child: Text('Major Service'),
                  ),
                  DropdownMenuItem(
                    value: 'repair',
                    child: Text('Repair / Parts Replacement'),
                  ),
                  DropdownMenuItem(
                    value: 'annual',
                    child: Text('Annual Maintenance'),
                  ),
                ],
                onChanged: (val) {
                  if (val != null) setState(() => _selectedType = val);
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _costCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Total Cost (₹)',
                  border: OutlineInputBorder(),
                ),
                validator: (val) =>
                    val == null || val.isEmpty ? 'Cost is required' : null,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _kmCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Odometer Reading (KM)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _centerCtrl,
                decoration: const InputDecoration(
                  labelText: 'Service Center Name',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _nextKmCtrl,
                keyboardType: TextInputType.number,
                decoration: const InputDecoration(
                  labelText: 'Next Service Odometer (KM - Optional)',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _descCtrl,
                maxLines: 3,
                decoration: const InputDecoration(
                  labelText: 'Description / Notes',
                  border: OutlineInputBorder(),
                ),
              ),
              const SizedBox(height: 40),
            ],
          ),
        ),
      ),
    );
  }
}
