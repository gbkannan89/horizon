import 'package:flutter/material.dart';

class FilterChipRow extends StatelessWidget {
  final List<FilterChipOption> options;
  final String selected;
  final ValueChanged<String> onSelected;

  const FilterChipRow({
    super.key,
    required this.options,
    required this.selected,
    required this.onSelected,
  });

  @override
  Widget build(BuildContext context) {
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        children: options.map((o) {
          final isSelected = selected == o.value;
          return Padding(
            padding: const EdgeInsets.only(right: 8),
            child: FilterChip(
              label: Text(o.label),
              selected: isSelected,
              onSelected: (_) => onSelected(o.value),
            ),
          );
        }).toList(),
      ),
    );
  }
}

class FilterChipOption {
  final String label;
  final String value;

  const FilterChipOption({required this.label, required this.value});
}
