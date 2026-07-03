import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_animate/flutter_animate.dart';
import '../models/acct_models.dart';
import '../repository/acct_repository.dart';
import '../widgets/acct_widgets.dart';
import 'package:horizon_mobile/core/theme/design_tokens.dart';
import 'package:horizon_mobile/core/ui_kit/glass_card.dart';

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
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: Colors.transparent,
      appBar: AppBar(
        title: Text(state.detail?.accountName ?? 'Account Details', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        backgroundColor: Colors.transparent,
        elevation: 0,
      ),
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
            colors: isDark 
                ? [AppColors.darkBackground, AppColors.navy900]
                : [AppColors.lightBackground, AppColors.teal50.withOpacity(0.5)],
          ),
        ),
        child: _buildBody(theme, state),
      ),
    );
  }

  Widget _buildBody(ThemeData theme, AcctDetailState state) {
    if (state.loading) {
      return Center(
        child: CircularProgressIndicator(color: AppColors.teal500).animate().fade(),
      );
    }
    if (state.error != null) {
      return Center(
        child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
          Icon(Icons.error_outline_rounded, size: 64, color: AppColors.red500),
          const SizedBox(height: 16), const Text('Could not load account details'),
          const SizedBox(height: 16),
          FilledButton.icon(
            style: FilledButton.styleFrom(backgroundColor: AppColors.teal500),
            onPressed: () => ref.read(acctDetailProvider(widget.accountId).notifier).load(), 
            icon: const Icon(Icons.refresh_rounded), 
            label: const Text('Try Again')
          ),
        ]),
      );
    }
    if (state.detail == null) return const Center(child: Text('No data'));

    final d = state.detail!;
    final children = <Widget>[
      // Header
      GlassCard(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [
          Icon(Icons.account_balance_rounded, size: 40, color: AppColors.teal500),
          const SizedBox(width: AppSpacing.md),
          Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(d.accountName, style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold)),
            Text('${d.accountType}  •  ${d.currency}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          ])),
        ]),
        const SizedBox(height: 24),
        if (d.institutionName != null) _row('Institution', d.institutionName!, theme),
        _row('Status', d.status, theme),
        _row('Health', d.accountHealth, theme),
        _row('Opened', d.openedDate, theme),
        if (d.interestRate != null) _row('Interest Rate', '${d.interestRate!.toStringAsFixed(2)}%', theme),
        if (d.creditLimit != null) _row('Credit Limit', AcctFmt(d.creditLimit!.toInt()), theme),
      ])),
      const SizedBox(height: 16),
      // Balances
      Text('Balances', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
      const SizedBox(height: 12),
      ...d.balances.map((b) => Padding(
        padding: const EdgeInsets.only(bottom: AppSpacing.sm),
        child: GlassCard(child: ListTile(
          contentPadding: EdgeInsets.zero,
          dense: true,
          title: Text(b.label, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
          subtitle: Text(b.description, style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
          trailing: Text(AcctFmt(b.value), style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold)),
        )),
      )),
      // Transactions
      if (d.transactions.isNotEmpty) ...[
        const SizedBox(height: 16),
        Text('Recent Transactions', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold)),
        const SizedBox(height: 12),
        ...d.transactions.take(10).map((t) => Padding(
          padding: const EdgeInsets.only(bottom: AppSpacing.sm),
          child: GlassCard(child: ListTile(
            contentPadding: EdgeInsets.zero,
            dense: true,
            leading: Container(
              width: 36, height: 36,
              decoration: BoxDecoration(color: (t.amount >= 0 ? AppColors.teal500 : AppColors.red500).withOpacity(0.15), borderRadius: BorderRadius.circular(AppRadius.sm)),
              child: Icon(t.amount >= 0 ? Icons.arrow_upward_rounded : Icons.arrow_downward_rounded, size: 20, color: t.amount >= 0 ? AppColors.teal500 : AppColors.red500),
            ),
            title: Text(t.description, style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600)),
            subtitle: Text('${t.category}  •  ${t.eventDate}', style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            trailing: Column(mainAxisAlignment: MainAxisAlignment.center, crossAxisAlignment: CrossAxisAlignment.end, children: [
              Text(AcctFmt(t.amount), style: TextStyle(fontWeight: FontWeight.bold, color: t.amount >= 0 ? AppColors.teal500 : AppColors.red500, fontSize: 13)),
              Text(AcctFmt(t.runningBalance), style: theme.textTheme.bodySmall?.copyWith(color: theme.colorScheme.onSurfaceVariant)),
            ]),
          )),
        )),
      ],
      const SizedBox(height: 100), // padding for navbar
    ];

    return RefreshIndicator(
      color: AppColors.teal500,
      onRefresh: () => ref.read(acctDetailProvider(widget.accountId).notifier).load(),
      child: ListView.builder(
        padding: const EdgeInsets.all(AppSpacing.md),
        itemCount: children.length,
        itemBuilder: (context, index) {
          return children[index]
            .animate()
            .fade(duration: 400.ms, delay: (20 * index).ms)
            .slideY(begin: 0.1, end: 0, duration: 400.ms, curve: Curves.easeOutQuad, delay: (20 * index).ms);
        },
      ),
    );
  }

  Widget _row(String l, String v, ThemeData t) => Padding(
    padding: const EdgeInsets.symmetric(vertical: 4),
    child: Row(children: [
      SizedBox(width: 120, child: Text(l, style: t.textTheme.bodySmall?.copyWith(color: t.colorScheme.onSurfaceVariant, fontWeight: FontWeight.w500))), 
      Text(v, style: t.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.w600))
    ]),
  );
}
