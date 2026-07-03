import 'package:flutter/material.dart';

class FrequencySelector extends StatelessWidget {
  final String value;
  final ValueChanged<String> onChanged;

  const FrequencySelector({super.key, required this.value, required this.onChanged});

  @override
  Widget build(BuildContext context) {
    return DropdownButtonFormField<String>(
      decoration: const InputDecoration(labelText: 'Frequency', border: OutlineInputBorder()),
      initialValue: value,
      items: const [
        DropdownMenuItem(value: 'Daily', child: Text('Daily')),
        DropdownMenuItem(value: 'Weekly', child: Text('Weekly')),
        DropdownMenuItem(value: 'BiWeekly', child: Text('Bi-Weekly')),
        DropdownMenuItem(value: 'Monthly', child: Text('Monthly')),
        DropdownMenuItem(value: 'Quarterly', child: Text('Quarterly')),
        DropdownMenuItem(value: 'SemiAnnual', child: Text('Semi-Annual')),
        DropdownMenuItem(value: 'Annual', child: Text('Annual')),
      ],
      onChanged: (v) {
        if (v != null) onChanged(v);
      },
    );
  }
}
