import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../models/acct_models.dart';
import '../repository/acct_repository.dart';
import '../widgets/acct_widgets.dart';

String _fmt(int v) { if (v >= 10000000) return '₹${(v / 10000000).toStringAsFixed(2)}Cr'; if (v >= 100000) return '₹${(v / 100000).toStringAsFixed(2)}L'; return '₹$v'; }

final acctStateProvider = StateNotifierProvider<AcctStateNotifier, AcctState>((ref) {
  return AcctStateNotifier(ref.read(acctRepositoryProvider));
});

class AcctState {
  final bool loading; final AcctDashboardData? dash; final AcctListData? list;
  final BalanceSummaryData? balances; final CashFlowData? cashFlow;
  final AcctHealthData? health; final String? error;
  final List<AcctCardData> filteredAccounts; final String filterType;

  AcctState({this.loading = false, this.dash, this.list, this.balances, this.cashFlow, this.health, this.error, this.filteredAccounts = const [], this.filterType = ''});
}

class AcctStateNotifier extends StateNotifier<AcctState> {
  final AcctRepository _repo;
  AcctStateNotifier(this._repo) : super(AcctState());

  Future<void> load({String? userId}) async {
    state = AcctState(loading: true);
    try {
      final r = await Future.wait([
        _repo.getDashboard(userId: userId),
        _repo.getAccounts(userId: userId),
        _repo.getBalances(userId: userId),
        _repo.getCashFlow(userId: userId),
        _repo.getHealth(userId: userId),
      ]);
      state = AcctState(
        dash: (r[0] as AcctDashboardResponse).data,
        list: (r[1] as AcctListResponse).data,
        balances: (r[2] as BalanceSummaryResponse).data,
        cashFlow: (r[3] as CashFlowResponse).data,
        health: (r[4] as AcctHealthResponse).data,
        filteredAccounts: (r[1] as AcctListResponse).data?.accounts ?? [],
      );
    } catch (e) {
      state = AcctState(error: e.toString());
    }
  }

  void filterByType(String type) {
    if (type.isEmpty || state.list == null) {
      state = AcctState(dash: state.dash, list: state.list, balances: state.balances, cashFlow: state.cashFlow, health: state.health, filteredAccounts: state.list?.accounts ?? []);
      return;
    }
    state = AcctState(dash: state.dash, list: state.list, balances: state.balances, cashFlow: state.cashFlow, health: state.health,
      filteredAccounts: state.list?.accounts.where((a) => a.accountType == type).toList() ?? [], filterType: type);
  }
}

class AccountsPage extends ConsumerStatefulWidget {
  const AccountsPage({super.key});
  @override
  ConsumerState<AccountsPage> createState() => _AccountsPageState();
}

class _AccountsPageState extends ConsumerState<AccountsPage> {
  @override
  void initState() { super.initState(); Future.microtask(() => ref.read(acctStateProvider.notifier).load()); }

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(acctStateProvider);
    final theme = Theme.of(context);
    return Scaffold(
      appBar: AppBar(title: const Text('Accounts')),
      body: _buildBody(theme, state),
    );
  }

  Widget _buildBody(ThemeData theme, AcctState state) {
    if (state.loading) return const Center(child: CircularProgressIndicator());
    if (state.error != null) return Center(
      child: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
        Icon(Icons.error_outline, size: 64, color: theme.colorScheme.error),
        const SizedBox(height: 16), Text('Could not load accounts'),
        FilledButton.icon(onPressed: () => ref.read(acctStateProvider.notifier).load(), icon: const Icon(Icons.refresh), label: const Text('Try Again')),
      ]),
    );

    final accts = state.filteredAccounts;
    return RefreshIndicator(
      onRefresh: () => ref.read(acctStateProvider.notifier).load(),
      child: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          if (state.balances != null) AcctBalanceCard(total: state.balances!.totalBalance, available: state.balances!.totalAvailable, spendable: state.balances!.totalSpendable),
          if (state.dash != null) ...[
            const SizedBox(height: 12),
            ...state.dash!.cards.map((c) => Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: AcctCardWidget(title: c.title, summary: c.summary, icon: _cardIcon(c.cardType), color: _cardColor(c.cardType)),
            )),
          ],
          if (state.cashFlow != null) ...[
            const SizedBox(height: 8),
            Card(child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text('Cash Flow', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 8),
                _cfRow('Inflow', AcctFmt(state.cashFlow!.periodInflow), Colors.green, theme),
                _cfRow('Outflow', AcctFmt(state.cashFlow!.periodOutflow), Colors.red, theme),
                const Divider(height: 16),
                _cfRow('Net', AcctFmt(state.cashFlow!.netFlow), state.cashFlow!.netFlow >= 0 ? Colors.green : Colors.red, theme),
                _cfRow('Projected', AcctFmt(state.cashFlow!.projectedFlow), Colors.blue, theme),
              ]),
            )),
          ],
          if (state.health != null) ...[
            const SizedBox(height: 8),
            Card(child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text('Account Health', style: theme.textTheme.titleSmall?.copyWith(fontWeight: FontWeight.w600)),
                const SizedBox(height: 8),
                Text('Healthy: ${state.health!.healthyPct.toStringAsFixed(0)}%', style: theme.textTheme.bodyMedium),
                ...state.health!.accountsByHealth.entries.map((e) => Padding(
                  padding: const EdgeInsets.symmetric(vertical: 2),
                  child: Text('${e.key}: ${e.value}', style: theme.textTheme.bodySmall),
                )),
              ]),
            )),
          ],
          if (accts.isNotEmpty) ...[
            const SizedBox(height: 12),
            Text('Accounts (${accts.length})', style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.w600)),
            const SizedBox(height: 8),
            // Quick filter chips
            SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(children: [
                _filterChip('All', '', state.filterType),
                _filterChip('Savings', 'Savings', state.filterType),
                _filterChip('Checking', 'Checking', state.filterType),
                _filterChip('Credit', 'Credit', state.filterType),
                _filterChip('Investment', 'Investment', state.filterType),
              ]),
            ),
            const SizedBox(height: 8),
            ...accts.map((a) => Padding(
              padding: const EdgeInsets.only(bottom: 8),
              child: AcctCard(acct: a, onTap: () => context.push('/accounts/${a.accountId}')),
            )),
          ],
        ],
      ),
    );
  }

  Widget _filterChip(String label, String value, String current) {
    final selected = current == value;
    return Padding(
      padding: const EdgeInsets.only(right: 8),
      child: FilterChip(
        label: Text(label), selected: selected,
        onSelected: (_) => ref.read(acctStateProvider.notifier).filterByType(value),
      ),
    );
  }

  Widget _cfRow(String label, String value, Color color, ThemeData t) => Padding(
    padding: const EdgeInsets.symmetric(vertical: 2),
    child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
      Text(label, style: t.textTheme.bodyMedium?.copyWith(color: t.colorScheme.onSurfaceVariant)),
      Text(value, style: TextStyle(fontWeight: FontWeight.w600, color: color)),
    ]),
  );

  IconData _cardIcon(String type) {
    switch (type) {
      case 'Overview': return Icons.dashboard;
      case 'CashBalance': return Icons.account_balance_wallet;
      case 'AccountHealth': return Icons.favorite;
      case 'CashFlow': return Icons.swap_horiz;
      case 'Risk': return Icons.shield;
      case 'Projection': return Icons.trending_up;
      case 'Recommendation': return Icons.lightbulb;
      case 'Timeline': return Icons.history;
      default: return Icons.circle;
    }
  }

  Color _cardColor(String type) {
    switch (type) {
      case 'Overview': return Colors.teal;
      case 'CashBalance': return Colors.green;
      case 'AccountHealth': return Colors.pink;
      case 'CashFlow': return Colors.blue;
      case 'Risk': return Colors.amber;
      case 'Projection': return Colors.cyan;
      case 'Recommendation': return Colors.amber;
      case 'Timeline': return Colors.brown;
      default: return Colors.grey;
    }
  }
}
