import 'package:flutter/foundation.dart';

class ParsedSmsTransaction {
  final String name;
  final double amount;
  final DateTime date;
  final bool isDebit; // true = expense, false = credit/income
  final double confidence;

  ParsedSmsTransaction({
    required this.name,
    required this.amount,
    required this.date,
    required this.isDebit,
    this.confidence = 0.0,
  });

  Map<String, dynamic> toJson() => {
    'name': name,
    'amount': amount,
    'date': '${date.year.toString().padLeft(4, '0')}-${date.month.toString().padLeft(2, '0')}-${date.day.toString().padLeft(2, '0')}',
  };
}

class SmsParser {
  // Regex to extract amount from messages like "Rs.1,200.00", "₹500", "INR 2,500"
  static final _amountRegex = RegExp(
    r'(?:Rs\.?\s*|₹\s*|INR\s*|INR\.?\s*)([\d,]+\.?\d*)',
    caseSensitive: false,
  );

  // Regex to extract date in various formats
  static final _dateRegex = RegExp(
    r'(\d{1,2})[-/](\d{1,2})[-/](\d{2,4})',
    caseSensitive: false,
  );

  static final _dateTextRegex = RegExp(
    r'(\d{1,2})\s+(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)[a-z]*\s+(\d{2,4})',
    caseSensitive: false,
  );

  // Debit indicators
  static final _debitKeywords = RegExp(
    r'\b(debited|debit|spent|paid|used|withdrawn|purchase|swiped|charged)\b',
    caseSensitive: false,
  );

  // Skip messages that are not transaction-related
  static final _skipKeywords = RegExp(
    r'\b(OTP|one\s*time\s*password|login|log\s*in|verification\s*code|registered|KYC|offer|reward|cashback|loan\s*approved|statement\s*generated)\b',
    caseSensitive: false,
  );

  List<ParsedSmsTransaction> parseBatch(List<String> smsMessages) {
    final results = <ParsedSmsTransaction>[];
    for (final msg in smsMessages) {
      final parsed = parseSingle(msg);
      if (parsed != null) {
        results.add(parsed);
      }
    }
    return results;
  }

  ParsedSmsTransaction? parseSingle(String message) {
    try {
      final text = message.trim();
      if (text.isEmpty) return null;

      // Skip non-transaction messages
      if (_skipKeywords.hasMatch(text)) return null;

      // Must have an amount
      final amountMatch = _amountRegex.firstMatch(text);
      if (amountMatch == null) return null;

      final amountStr = amountMatch.group(1)!.replaceAll(',', '');
      final amount = double.tryParse(amountStr);
      if (amount == null || amount <= 0) return null;

      final isDebit = _debitKeywords.hasMatch(text);

      // Extract date
      DateTime? txnDate = _extractDate(text);
      if (txnDate == null) {
        txnDate = DateTime.now();
      }

      // Extract name/merchant
      final name = _extractName(text, amountMatch.start);

      // Calculate confidence
      double confidence = 0.7;
      if (amountMatch.start > 0) confidence += 0.1;
      if (txnDate != DateTime.now()) confidence += 0.1;
      if (name.length > 3) confidence += 0.1;

      return ParsedSmsTransaction(
        name: name,
        amount: amount,
        date: txnDate,
        isDebit: isDebit,
        confidence: confidence.clamp(0.0, 1.0),
      );
    } catch (e) {
      debugPrint('SMS parse error: $e');
      return null;
    }
  }

  DateTime? _extractDate(String text) {
    // Try "dd Mon YYYY" format first
    final textMatch = _dateTextRegex.firstMatch(text);
    if (textMatch != null) {
      final day = int.tryParse(textMatch.group(1)!) ?? 1;
      final month = _monthNumber(textMatch.group(2)!);
      var year = int.tryParse(textMatch.group(3)!) ?? DateTime.now().year;
      if (year < 100) year += 2000;
      return DateTime(year, month, day);
    }

    // Try dd/mm/yyyy or dd-mm-yyyy
    final numMatch = _dateRegex.firstMatch(text);
    if (numMatch != null) {
      final a = int.tryParse(numMatch.group(1)!) ?? 1;
      final b = int.tryParse(numMatch.group(2)!) ?? 1;
      var c = int.tryParse(numMatch.group(3)!) ?? DateTime.now().year;
      if (c < 100) c += 2000;

      // Determine if it's DD/MM or MM/DD based on values
      if (a > 12) {
        // Must be DD
        return DateTime(c, b, a);
      } else if (b > 12) {
        // Must be MM/DD
        return DateTime(c, a, b);
      } else {
        // Ambiguous — assume DD/MM (common in India)
        return DateTime(c, b, a);
      }
    }

    return null;
  }

  String _extractName(String text, int amountEnd) {
    // Look for merchant name after "at", "to", "for", "at "
    final afterAmount = text.substring(amountEnd);
    final nameMatch = RegExp(
      r"\b(?:at|to|for|via|info[:\s]*|merchant[:\s]*)\s+([A-Za-z0-9\s.&,'-]+?)(?:\s+(?:on|ref|avl|available|txn|bank|\d|$))",
      caseSensitive: false,
    ).firstMatch(afterAmount);

    if (nameMatch != null) {
      return nameMatch.group(1)!.trim();
    }

    // Fallback: use text before "debited" or "credited"
    final actionMatch = RegExp(
      r'(?:Rs\.?\s*|₹\s*|INR\s*)[\d,]+\s+(?:debited|credited|spent|paid)\s+[^.]*?(?:from|to|at)\s+([A-Za-z\s]+)',
      caseSensitive: false,
    ).firstMatch(text);

    if (actionMatch != null) {
      return actionMatch.group(1)!.trim();
    }

    // Last resort: first meaningful word after amount
    final firstWord = text.split(RegExp(r'\s+')).firstWhere(
      (w) => w.length > 3 && !w.contains(RegExp(r'[0-9₹Rs]')),
      orElse: () => 'Unknown',
    );
    return firstWord;
  }

  int _monthNumber(String abbr) {
    const months = {
      'jan': 1, 'feb': 2, 'mar': 3, 'apr': 4, 'may': 5, 'jun': 6,
      'jul': 7, 'aug': 8, 'sep': 9, 'oct': 10, 'nov': 11, 'dec': 12,
    };
    return months[abbr.substring(0, 3).toLowerCase()] ?? DateTime.now().month;
  }
}
