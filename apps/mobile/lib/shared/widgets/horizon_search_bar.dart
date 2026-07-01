import 'dart:async';
import 'package:flutter/material.dart';

class HorizonSearchBar extends StatefulWidget {
  final String hintText;
  final ValueChanged<String> onSearch;
  final VoidCallback? onClose;
  final bool autofocus;
  final Duration debounceDuration;

  const HorizonSearchBar({
    super.key,
    required this.hintText,
    required this.onSearch,
    this.onClose,
    this.autofocus = true,
    this.debounceDuration = const Duration(milliseconds: 300),
  });

  @override
  State<HorizonSearchBar> createState() => _HorizonSearchBarState();
}

class _HorizonSearchBarState extends State<HorizonSearchBar> {
  final _ctrl = TextEditingController();
  Timer? _debounce;

  @override
  void initState() {
    super.initState();
    _ctrl.addListener(_onChanged);
  }

  @override
  void dispose() {
    _ctrl.removeListener(_onChanged);
    _ctrl.dispose();
    _debounce?.cancel();
    super.dispose();
  }

  void _onChanged() {
    _debounce?.cancel();
    _debounce = Timer(widget.debounceDuration, () => widget.onSearch(_ctrl.text));
  }

  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: _ctrl,
      autofocus: widget.autofocus,
      decoration: InputDecoration(
        hintText: widget.hintText,
        border: InputBorder.none,
        suffixIcon: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () {
            _ctrl.clear();
            widget.onSearch('');
            widget.onClose?.call();
          },
        ),
      ),
    );
  }
}
