extension MoneyFormat on int {
  String get toRupees {
    if (this >= 10000000) return '₹${(this / 10000000).toStringAsFixed(2)}Cr';
    if (this >= 100000) return '₹${(this / 100000).toStringAsFixed(2)}L';
    return '₹$this';
  }
}

extension StringCapitalize on String {
  String get capitalize => isNotEmpty ? '${this[0].toUpperCase()}${substring(1)}' : '';
}
