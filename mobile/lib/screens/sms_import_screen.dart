import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../services/sms_service.dart';
import '../services/sms_parser.dart';
import '../utils/ui_utils.dart';

class SmsImportScreen extends StatefulWidget {
  const SmsImportScreen({super.key});

  @override
  State<SmsImportScreen> createState() => _SmsImportScreenState();
}

class _SmsImportScreenState extends State<SmsImportScreen> {
  final AndroidSmsReader _smsReader = AndroidSmsReader();
  final SmsService _smsService = SmsService();

  bool _isLoading = false;
  bool _hasPermission = false;
  List<ParsedSmsTransaction> _detected = [];
  Set<int> _selected = {};

  @override
  void initState() {
    super.initState();
    _checkPermission();
  }

  Future<void> _checkPermission() async {
    final granted = await _smsReader.hasPermission();
    if (mounted) setState(() => _hasPermission = granted);
  }

  Future<void> _requestPermissionAndScan() async {
    setState(() => _isLoading = true);

    if (!_hasPermission) {
      final granted = await _smsReader.requestPermission();
      if (!granted) {
        if (mounted) {
          UiUtils.showSnack(
            context,
            'SMS permission is required to read bank transaction messages.',
            isError: true,
          );
          setState(() => _isLoading = false);
        }
        return;
      }
      setState(() => _hasPermission = true);
    }

    await _scanMessages();
  }

  Future<void> _scanMessages() async {
    setState(() => _isLoading = true);
    try {
      // Read last 30 days of SMS
      final messages = await _smsReader.readInbox(daysBack: 30);

      // Filter to transaction-related messages
      final txMessages = _smsService.filterTransactionMessages(messages);

      // Parse debit transactions (expenses)
      _detected = _smsService.parseTransactions(txMessages);

      // Select all by default
      _selected = Set.from(List.generate(_detected.length, (i) => i));

      if (mounted) {
        setState(() => _isLoading = false);
        if (_detected.isEmpty) {
          UiUtils.showSnack(context, 'No transaction SMS found in the last 30 days.');
        }
      }
    } catch (e) {
      if (mounted) {
        setState(() => _isLoading = false);
        UiUtils.showSnack(context, 'Error reading SMS: $e', isError: true);
      }
    }
  }

  Future<void> _importSelected() async {
    final selected = _selected.map((i) => _detected[i].toJson()).toList();
    if (selected.isEmpty) return;

    setState(() => _isLoading = true);
    try {
      final fp = Provider.of<FinancialProvider>(context, listen: false);
      final result = await fp.importSmsTransactions(selected);
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Imported ${result['inserted'] ?? 0} transactions. '
                '${result['duplicates'] ?? 0} duplicates skipped.'),
            backgroundColor: const Color(0xFF059669),
          ),
        );
        Navigator.of(context).pop();
      }
    } catch (e) {
      if (mounted) {
        UiUtils.showSnack(context, 'Import failed: $e', isError: true);
      }
    } finally {
      if (mounted) setState(() => _isLoading = false);
    }
  }

  String _fmtDate(DateTime d) {
    const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
      'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
    return '${d.day} ${months[d.month - 1]} ${d.year}';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        title: const Text('Import from SMS'),
        actions: [
          if (_detected.isNotEmpty)
            TextButton(
              onPressed: _isLoading ? null : _importSelected,
              child: _isLoading
                  ? const SizedBox(width: 20, height: 20,
                      child: CircularProgressIndicator(strokeWidth: 2))
                  : Text('Import ${_selected.length}',
                      style: const TextStyle(fontWeight: FontWeight.bold)),
            ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_isLoading && _detected.isEmpty) {
      return const Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircularProgressIndicator(color: Color(0xFF1E3A8A)),
            SizedBox(height: 20),
            Text('Reading SMS messages...',
                style: TextStyle(fontSize: 16, color: Colors.grey)),
          ],
        ),
      );
    }

    if (!_hasPermission && _detected.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.sms_outlined, size: 80, color: Colors.grey.shade300),
            const SizedBox(height: 20),
            const Text('Read Bank SMS Automatically',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            const Padding(
              padding: EdgeInsets.symmetric(horizontal: 40),
              child: Text(
                'Horizon can scan your SMS for bank transactions and add them automatically.',
                textAlign: TextAlign.center,
                style: TextStyle(color: Colors.grey),
              ),
            ),
            const SizedBox(height: 32),
            ElevatedButton.icon(
              onPressed: _requestPermissionAndScan,
              icon: const Icon(Icons.settings),
              label: const Text('Grant SMS Permission & Scan'),
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 14),
              ),
            ),
          ],
        ),
      );
    }

    if (_detected.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(Icons.check_circle_outline, size: 80, color: Colors.green.shade300),
            const SizedBox(height: 20),
            const Text('No Transactions Found',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold)),
            const SizedBox(height: 8),
            const Text('No bank transaction SMS found in the last 30 days.',
                style: TextStyle(color: Colors.grey)),
            const SizedBox(height: 24),
            OutlinedButton.icon(
              onPressed: _requestPermissionAndScan,
              icon: const Icon(Icons.refresh),
              label: const Text('Scan Again'),
            ),
          ],
        ),
      );
    }

    return Column(
      children: [
        // Summary bar
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
          color: Colors.white,
          child: Row(
            children: [
              Text('${_detected.length} transactions found',
                  style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 15)),
              const Spacer(),
              Text('${_selected.length} selected',
                  style: TextStyle(color: Colors.grey.shade600, fontSize: 13)),
            ],
          ),
        ),
        const Divider(height: 1),

        // List
        Expanded(
          child: ListView.builder(
            padding: const EdgeInsets.all(16),
            itemCount: _detected.length,
            itemBuilder: (ctx, i) {
              final t = _detected[i];
              final sel = _selected.contains(i);
              return GestureDetector(
                onTap: () {
                  setState(() {
                    if (sel) { _selected.remove(i); }
                    else { _selected.add(i); }
                  });
                },
                child: Container(
                  margin: const EdgeInsets.only(bottom: 8),
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(14),
                    border: Border.all(
                      color: sel ? const Color(0xFF1E3A8A) : Colors.grey.shade200,
                      width: sel ? 2 : 1,
                    ),
                  ),
                  child: Row(
                    children: [
                      Container(
                        width: 24, height: 24,
                        decoration: BoxDecoration(
                          color: sel
                              ? const Color(0xFF1E3A8A)
                              : Colors.transparent,
                          shape: BoxShape.circle,
                          border: Border.all(
                            color: sel ? const Color(0xFF1E3A8A) : Colors.grey.shade400,
                          ),
                        ),
                        child: sel
                            ? const Icon(Icons.check, color: Colors.white, size: 16)
                            : null,
                      ),
                      const SizedBox(width: 14),
                      Container(
                        width: 40, height: 40,
                        decoration: BoxDecoration(
                          color: t.isDebit
                              ? const Color(0xFFE88A1A).withValues(alpha: 0.1)
                              : const Color(0xFF059669).withValues(alpha: 0.1),
                          borderRadius: BorderRadius.circular(12),
                        ),
                        child: Icon(
                          t.isDebit ? Icons.arrow_upward : Icons.arrow_downward,
                          color: t.isDebit ? const Color(0xFFE88A1A) : const Color(0xFF059669),
                          size: 20,
                        ),
                      ),
                      const SizedBox(width: 12),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(t.name,
                                style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 14),
                                maxLines: 1, overflow: TextOverflow.ellipsis),
                            const SizedBox(height: 2),
                            Text(_fmtDate(t.date),
                                style: TextStyle(fontSize: 12, color: Colors.grey.shade500)),
                          ],
                        ),
                      ),
                      Text('₹${t.amount.toStringAsFixed(0)}',
                          style: TextStyle(
                            fontWeight: FontWeight.bold, fontSize: 15,
                            color: t.isDebit ? const Color(0xFFE88A1A) : const Color(0xFF059669),
                          )),
                      const SizedBox(width: 4),
                      Container(
                        padding: const EdgeInsets.symmetric(horizontal: 5, vertical: 2),
                        decoration: BoxDecoration(
                          color: Colors.grey.shade100,
                          borderRadius: BorderRadius.circular(4),
                        ),
                        child: Text('${(t.confidence * 100).round()}%',
                            style: TextStyle(fontSize: 9, color: Colors.grey.shade600)),
                      ),
                    ],
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }
}
