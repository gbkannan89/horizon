import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/shared/widgets/index.dart';
import '../models/acct_models.dart';
import '../repository/acct_repository.dart';

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
    if (state.loading) return const LoadingView();
    if (state.error != null) return ErrorView(
      message: 'Could not load account details',
      onRetry: () => ref.read(acctDetailProvider(widget.accountId).notifier).load(),
    );
    if (state.detail == null) return const EmptyState(icon: Icons.search_off, title: 'No data');

    final d = state.detail!;
    return RefreshIndicator(
      onRefresh: () => ref.read(acctDetailProvider(widget.accountId).notifier).load(),
      child: ListView(
        padding: const EdgeInsets.all(AppTheme.spacingLg),
        children: [
          AppCard(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Row(children: [
              Icon(Icons.account_balance, size: 40, color: theme.colorScheme.primary),
              const SizedBox(width: AppTheme.spacingMd),
              Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(d.accountName, style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
                Text('${d.accountType}  •  ${d.currency}', style: theme.textTheme.bodySmall),
              ])),
            ]),
            const SizedBox(height: AppTheme.spacingLg),
            InfoRow(label: 'Institution', value: d.institutionName ?? '--'),
            InfoRow(label: 'Status', value: d.status),
            InfoRow(label: 'Health', value: d.accountHealth),
            InfoRow(label: 'Opened', value: d.openedDate),
            if (d.interestRate != null) InfoRow(label: 'Interest Rate', value: '${d.interestRate!.toStringAsFixed(2)}%'),
            if (d.creditLimit != null) InfoRow(label: 'Credit Limit', value: formatMoney(d.creditLimit!.toInt())),
          ])),
          const SizedBox(height: AppTheme.spacingMd),
          SectionHeader(title: 'Balances'),
          const SizedBox(height: AppTheme.spacingSm),
          ...d.balances.map((b) => Card(child: ListTile(
            dense: true,
            title: Text(b.label, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w500)),
            subtitle: Text(b.description, style: theme.textTheme.bodySmall),
            trailing: Text(formatMoney(b.value), style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
          ))),
          if (d.transactions.isNotEmpty) ...[
            const SizedBox(height: AppTheme.spacingMd),
            SectionHeader(title: 'Recent Transactions'),
            const SizedBox(height: AppTheme.spacingSm),
            ...d.transactions.take(10).map((t) => Card(child: ListTile(
              dense: true,
              leading: Icon(t.amount >= 0 ? Icons.arrow_upward : Icons.arrow_downward, size: AppTheme.iconSm, color: t.amount >= 0 ? Colors.green : Colors.red),
              title: Text(t.description, style: theme.textTheme.bodyMedium),
              subtitle: Text('${t.category}  •  ${t.eventDate}', style: theme.textTheme.bodySmall),
              trailing: Column(mainAxisAlignment: MainAxisAlignment.center, crossAxisAlignment: CrossAxisAlignment.end, children: [
                Text(formatMoney(t.amount), style: TextStyle(fontWeight: FontWeight.w600, color: t.amount >= 0 ? Colors.green : Colors.red, fontSize: 13)),
                Text(formatMoney(t.runningBalance), style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
              ]),
            ))),
          ],
        ],
      ),
    );
  }
}
