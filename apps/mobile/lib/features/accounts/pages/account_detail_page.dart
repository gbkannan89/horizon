import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/acct_models.dart';
import '../repository/acct_repository.dart';
import '../widgets/acct_widgets.dart';

String _fmt(int v) { if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr'; if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L'; return '₹$v'; }

final acctDetailProvider = StateNotifierProvider.family<AcctDetailNotifier, AcctDetailState, String>((ref, id) {
  return AcctDetailNotifier(ref.read(acctRepositoryProvider), id);
});

class AcctDetailState {
  final bool loading; final AcctDetailData? detail; final String? error;
  AcctDetailState({this.loading = false, this.detail, this.error});
}

class AcctDetailNotifier extends StateNotifier<AcctDetailState> {
  final AcctRepository _repo; final String _id;
  AcctDetailNotifier(this._repo, this._id) : super(AcctDetailState());

  Future<void> load({String? userId}) async {
    state = AcctDetailState(loading: true);
    try {
      final r = await _repo.getDetail(userId: userId, accountId: _id);
      state = AcctDetailState(detail: r.data);
    } catch (e) { state = AcctDetailState(error: e.toString()); }
  }
}

class AccountDetailPage extends ConsumerStatefulWidget {
  final String accountId;
  const AccountDetailPage({super.key, required this.accountId});
  @override
  ConsumerState<AccountDetailPage> createState() => _AccountDetailPageState();
}

class _AccountDetailPageState extends ConsumerState<AccountDetailPage> {
  @override
  void initState() { super.initState(); Future.microtask(() => ref.read(acctDetailProvider(widget.accountId).notifier).load()); }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(acctDetailProvider(widget.accountId));
    final theme = Theme.of(context);

    return Scaffold(
      appBar: AppBar(title: Text(state.detail?.accountName ?? 'Account Details')),
      body: _buildBody(theme, state),
    );
  }

  Widget _buildBody(ThemeData theme, AcctDetailState state) {
    if (state.loading) return const Center(child: CircularProgressIndicator());
    if (state.error != null) {
      return Center(
      child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
        Icon(Icons.error_outline, size: 64, color: theme.colorScheme.error),
        const SizedBox(height: 16), const Text('Could not load account details'),
        FilledButton.icon(onPressed: () => ref.read(acctDetailProvider(widget.accountId).notifier).load(), icon: const Icon(Icons.refresh), label: const Text('Try Again')),
      ]),
    );
    }
    if (state.detail == null) return const Center(child: Text('No data'));

    final d = state.detail!;
    return RefreshIndicator(
      onRefresh: () => ref.read(acctDetailProvider(widget.accountId).notifier).load(),
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Header
          Card(child: Padding(
            padding: const EdgeInsets.all(16),
            child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
              Row(children: [
                Icon(Icons.account_balance, size: 40, color: theme.colorScheme.primary),
                const SizedBox(width: 12),
                Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                  Text(d.accountName, style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
                  Text('${d.accountType}  •  ${d.currency}', style: theme.textTheme.bodySmall),
                ])),
              ]),
              const SizedBox(height: 16),
              if (d.institutionName != null) _row('Institution', d.institutionName!, theme),
              _row('Status', d.status, theme),
              _row('Health', d.accountHealth, theme),
              _row('Opened', d.openedDate, theme),
              if (d.interestRate != null) _row('Interest Rate', '${d.interestRate!.toStringAsFixed(2)}%', theme),
              if (d.creditLimit != null) _row('Credit Limit', AcctFmt(d.creditLimit!.toInt()), theme),
            ]),
          )),
          const SizedBox(height: 12),
          // Balances
          Text('Balances', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
          const SizedBox(height: 8),
          ...d.balances.map((b) => Card(child: ListTile(
            dense: true,
            title: Text(b.label, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500)),
            subtitle: Text(b.description, style: theme.textTheme.bodySmall),
            trailing: Text(AcctFmt(b.value), style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
          ))),
          // Transactions
          if (d.transactions.isNotEmpty) ...[
            const SizedBox(height: 12),
            Text('Recent Transactions', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            ...d.transactions.take(10).map((t) => Card(child: ListTile(
              dense: true,
              leading: Icon(t.amount >= 0 ? Icons.arrow_upward : Icons.arrow_downward, size: 16, color: t.amount >= 0 ? Colors.green : Colors.red),
              title: Text(t.description, style: theme.textTheme.bodyMedium),
              subtitle: Text('${t.category}  •  ${t.eventDate}', style: theme.textTheme.bodySmall),
              trailing: Column(mainAxisAlignment: MainAxisAlignment.center, crossAxisAlignment: CrossAxisAlignment.end, children: [
                Text(AcctFmt(t.amount), style: TextStyle(fontWeight: FontWeight.w600, color: t.amount >= 0 ? Colors.green : Colors.red, fontSize: 13)),
                Text(AcctFmt(t.runningBalance), style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ]),
            ))),
          ],
        ],
      ),
    );
  }

  Widget _row(String l, String v, ThemeData t) => Padding(
    padding: const EdgeInsets.symmetric(vertical: 3),
    child: Row(children: [SizedBox(width: 100, child: Text(l, style: t.textTheme.bodySmall?.copyWith(color: t.colorScheme.onSurfaceVariant))), Text(v, style: t.textTheme.bodyMedium)]),
  );
}
