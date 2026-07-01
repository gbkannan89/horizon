String formatMoney(num v) {
  final amount = v is double ? v.toInt() : v as int;
  if (amount >= 10000000) return '₹${(amount / 10000000).toStringAsFixed(2)}Cr';
  if (amount >= 100000) return '₹${(amount / 100000).toStringAsFixed(2)}L';
  return '₹$amount';
}
