import 'dart:async';
import 'package:flutter/foundation.dart';
import 'package:flutter/services.dart';
import 'sms_parser.dart';

/// Service for reading SMS messages and detecting transactions.
/// Uses platform channels since telephony package requires native setup.
class SmsService {
  final SmsParser _parser = SmsParser();

  /// Keywords that suggest a message is bank-transaction related
  static const _transactionKeywords = [
    'debited', 'credited', 'spent', 'paid', 'sent', 'received',
    'upi', 'nft', 'neft', 'imps', 'trf', 'transfer',
    'rs.', 'rs ', 'inr', 'a/c', 'account',
  ];

  /// Check if a message body looks like a bank transaction
  bool isTransactionSms(String body) {
    final lower = body.toLowerCase();
    return _transactionKeywords.any((kw) => lower.contains(kw));
  }

  /// Parse SMS messages into transactions.
  /// Only returns debit transactions (expenses) since credits go to collections.
  List<ParsedSmsTransaction> parseTransactions(List<String> messages) {
    final all = _parser.parseBatch(messages);
    return all.where((t) => t.isDebit).toList();
  }

  /// Parse credit SMS (money received — for collection matching)
  List<ParsedSmsTransaction> parseCreditTransactions(List<String> messages) {
    final all = _parser.parseBatch(messages);
    return all.where((t) => !t.isDebit).toList();
  }

  /// Filter messages to only transaction-related ones
  List<String> filterTransactionMessages(List<String> allMessages) {
    return allMessages.where(isTransactionSms).toList();
  }

  /// Format messages for API import (debits only)
  List<Map<String, dynamic>> formatForImport(List<ParsedSmsTransaction> txns) {
    return txns.map((t) => t.toJson()).toList();
  }
}

/// Android-specific SMS reader using MethodChannel.
/// This is a wrapper that communicates with the Android native SMS reader.
class AndroidSmsReader {
  static const _channel = MethodChannel('com.thaari.horizon/sms');

  /// Request READ_SMS permission on Android
  Future<bool> requestPermission() async {
    try {
      final result = await _channel.invokeMethod<bool>('requestSmsPermission');
      return result ?? false;
    } catch (e) {
      debugPrint('SMS permission error: $e');
      return false;
    }
  }

  /// Check if READ_SMS permission is granted
  Future<bool> hasPermission() async {
    try {
      final result = await _channel.invokeMethod<bool>('hasSmsPermission');
      return result ?? false;
    } catch (e) {
      debugPrint('SMS permission check error: $e');
      return false;
    }
  }

  /// Read SMS inbox messages from the last N days
  Future<List<String>> readInbox({int daysBack = 30}) async {
    try {
      final result = await _channel.invokeMethod<List<dynamic>>(
        'readSmsInbox',
        {'daysBack': daysBack},
      );
      if (result == null) return [];
      return result.cast<String>();
    } catch (e) {
      debugPrint('SMS read error: $e');
      return [];
    }
  }

  /// Start listening for new SMS messages
  Future<void> startListening({
    required Function(String message) onMessage,
  }) async {
    try {
      _channel.setMethodCallHandler((call) async {
        if (call.method == 'onSmsReceived') {
          final message = call.arguments as String?;
          if (message != null) {
            onMessage(message);
          }
        }
      });
      await _channel.invokeMethod('startSmsListener');
      debugPrint('SMS listener started');
    } catch (e) {
      debugPrint('SMS listener error: $e');
    }
  }

  /// Stop listening for new SMS messages
  Future<void> stopListening() async {
    try {
      await _channel.invokeMethod('stopSmsListener');
      _channel.setMethodCallHandler(null);
      debugPrint('SMS listener stopped');
    } catch (e) {
      debugPrint('SMS listener stop error: $e');
    }
  }
}
