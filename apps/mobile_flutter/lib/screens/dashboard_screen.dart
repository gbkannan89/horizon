import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';

class DashboardScreen extends StatefulWidget {
  const DashboardScreen({super.key});

  @override
  State<DashboardScreen> createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<FinancialProvider>(context, listen: false).loadAllData();
    });
  }

  Color _parseColor(String hex) {
    hex = hex.replaceAll('#', '');
    if (hex.length == 6) hex = 'FF$hex';
    return Color(int.parse('0x$hex'));
  }

  String _greeting() {
    final h = DateTime.now().hour;
    if (h < 12) return 'Good morning';
    if (h < 17) return 'Good afternoon';
    return 'Good evening';
  }

  @override
  Widget build(BuildContext context) {
    final provider = context.watch<FinancialProvider>();
    final loaded = !provider.isLoading;

    return Scaffold(
      backgroundColor: const Color(0xFFF1F5F9),
      body: RefreshIndicator(
        onRefresh: provider.loadAllData,
        color: const Color(0xFF1E3A8A),
        child: CustomScrollView(
          slivers: [
            // ── Hero gradient header ──────────────────────────────────────────
            SliverToBoxAdapter(child: _buildHeroHeader(provider, loaded)),

            // ── Body cards ───────────────────────────────────────────────────
            SliverPadding(
              padding: const EdgeInsets.fromLTRB(20, 0, 20, 100),
              sliver: SliverList(
                delegate: SliverChildListDelegate([
                  const SizedBox(height: 20),
                  _buildIncomeSpentRow(provider, loaded),
                  const SizedBox(height: 20),
                  _buildBudgetCard(provider, loaded),
                  const SizedBox(height: 20),
                  _buildRecentExpenses(provider, loaded),
                  const SizedBox(height: 20),
                  _buildGoalsCard(provider, loaded),
                ]),
              ),
            ),
          ],
        ),
      ),
    );
  }

  // ── HERO HEADER ────────────────────────────────────────────────────────────
  Widget _buildHeroHeader(FinancialProvider p, bool loaded) {
    return Container(
      decoration: const BoxDecoration(
        gradient: LinearGradient(
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
          colors: [Color(0xFF0F2057), Color(0xFF1E3A8A), Color(0xFF2563EB)],
        ),
        borderRadius: BorderRadius.vertical(bottom: Radius.circular(36)),
      ),
      padding: const EdgeInsets.fromLTRB(24, 56, 24, 32),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        // Top row
        Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
          Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Text(_greeting(), style: const TextStyle(color: Color(0xFF93C5FD), fontSize: 14, fontWeight: FontWeight.w500)),
            const SizedBox(height: 2),
            const Text('Kannan 👋', style: TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.w800)),
          ]),
          Row(children: [
            Container(
              padding: const EdgeInsets.all(10),
              decoration: BoxDecoration(color: Colors.white.withOpacity(0.12), borderRadius: BorderRadius.circular(14)),
              child: const Icon(Icons.notifications_none_rounded, color: Colors.white, size: 22),
            ),
            const SizedBox(width: 10),
            Container(
              padding: const EdgeInsets.all(2),
              decoration: BoxDecoration(
                shape: BoxShape.circle,
                border: Border.all(color: Colors.white.withOpacity(0.4), width: 2),
              ),
              child: const CircleAvatar(
                radius: 18, backgroundColor: Color(0xFF3B82F6),
                child: Text('K', style: TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 16)),
              ),
            ),
          ]),
        ]),

        const SizedBox(height: 28),

        // Score + Net Worth row
        Row(children: [
          // Financial Score ring
          Stack(alignment: Alignment.center, children: [
            CircularPercentIndicator(
              radius: 52,
              lineWidth: 7,
              percent: loaded ? (p.finScoreVal / 100.0).clamp(0.0, 1.0) : 0.0,
              progressColor: const Color(0xFF34D399),
              backgroundColor: Colors.white.withOpacity(0.15),
              circularStrokeCap: CircularStrokeCap.round,
              center: Column(mainAxisAlignment: MainAxisAlignment.center, children: [
                Text(loaded ? '${p.finScoreVal}' : '--',
                  style: const TextStyle(fontSize: 26, fontWeight: FontWeight.w900, color: Colors.white)),
                const Text('/100', style: TextStyle(fontSize: 10, color: Color(0xFF93C5FD))),
              ]),
            ),
          ]),
          const SizedBox(width: 20),
          Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            const Text('Financial Score', style: TextStyle(color: Color(0xFF93C5FD), fontSize: 12, fontWeight: FontWeight.w500)),
            const SizedBox(height: 4),
            Text(
              loaded && p.finScoreVal >= 80 ? '🌟 Excellent' : loaded && p.finScoreVal >= 60 ? '👍 Good' : '⚡ Improving',
              style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 15),
            ),
            const SizedBox(height: 16),
            const Text('Net Worth', style: TextStyle(color: Color(0xFF93C5FD), fontSize: 12, fontWeight: FontWeight.w500)),
            const SizedBox(height: 2),
            Text(
              loaded ? '₹${_formatNum(p.netWorth)}' : '₹--',
              style: const TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.w900, letterSpacing: -0.5),
            ),
          ])),
        ]),
      ]),
    );
  }

  // ── INCOME / SPENT / LEFT CARDS ────────────────────────────────────────────
  Widget _buildIncomeSpentRow(FinancialProvider p, bool loaded) {
    return Row(children: [
      _miniCard('Income',  loaded ? '₹${_formatNum(p.totalIncomeAgg)}' : '--', const Color(0xFF059669), Icons.trending_up_rounded),
      const SizedBox(width: 10),
      _miniCard('Spent',   loaded ? '₹${_formatNum(p.totalSpent)}'    : '--', const Color(0xFFE88A1A), Icons.shopping_cart_outlined),
      const SizedBox(width: 10),
      _miniCard('Left',    loaded ? '₹${_formatNum(p.totalLeft)}'     : '--', const Color(0xFF1E3A8A), Icons.savings_outlined),
    ]);
  }

  Widget _miniCard(String label, String value, Color color, IconData icon) {
    return Expanded(child: Container(
      padding: const EdgeInsets.all(14),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [BoxShadow(color: color.withOpacity(0.08), blurRadius: 16, offset: const Offset(0, 6))],
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Container(
          padding: const EdgeInsets.all(7),
          decoration: BoxDecoration(color: color.withOpacity(0.1), borderRadius: BorderRadius.circular(10)),
          child: Icon(icon, color: color, size: 16),
        ),
        const SizedBox(height: 10),
        Text(label, style: const TextStyle(color: Color(0xFF94A3B8), fontSize: 11, fontWeight: FontWeight.w600)),
        const SizedBox(height: 2),
        Text(value, style: TextStyle(color: color, fontSize: 15, fontWeight: FontWeight.w900)),
      ]),
    ));
  }

  // ── BUDGET OVERVIEW ────────────────────────────────────────────────────────
  Widget _buildBudgetCard(FinancialProvider p, bool loaded) {
    return _card(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      _cardHeader('Budget Overview', Icons.pie_chart_rounded, const Color(0xFF6B46C1)),
      const SizedBox(height: 20),
      Row(mainAxisAlignment: MainAxisAlignment.spaceAround, children: [
        _budgetCircle('Needs',    p.needsSpent,   p.needsBudget,   const Color(0xFFE88A1A), loaded),
        _budgetCircle('Wants',    p.wantsSpent,   p.wantsBudget,   const Color(0xFF6B46C1), loaded),
        _budgetCircle('Savings',  p.savingsSpent, p.savingsBudget, const Color(0xFF059669), loaded),
      ]),
    ]));
  }

  Widget _budgetCircle(String label, double spent, double budget, Color color, bool loaded) {
    final pct = loaded && budget > 0 ? (spent / budget).clamp(0.0, 1.0) : 0.0;
    return Column(children: [
      CircularPercentIndicator(
        radius: 40, lineWidth: 7,
        percent: pct,
        progressColor: color,
        backgroundColor: color.withOpacity(0.1),
        circularStrokeCap: CircularStrokeCap.round,
        center: Text('${(pct * 100).toInt()}%',
          style: TextStyle(fontSize: 13, fontWeight: FontWeight.w800, color: color)),
      ),
      const SizedBox(height: 10),
      Text(label, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 13, color: Color(0xFF1E293B))),
      const SizedBox(height: 2),
      Text(
        loaded ? '₹${_formatNum(spent)} / ₹${_formatNum(budget)}' : '--',
        style: const TextStyle(color: Color(0xFF94A3B8), fontSize: 10),
      ),
    ]);
  }

  // ── RECENT EXPENSES ────────────────────────────────────────────────────────
  Widget _buildRecentExpenses(FinancialProvider p, bool loaded) {
    final expenses = p.recentExpenses.take(5).toList();
    if (!loaded || expenses.isEmpty) return const SizedBox.shrink();

    final bucketColor = {'Needs': const Color(0xFFE88A1A), 'Wants': const Color(0xFF6B46C1), 'Savings': const Color(0xFF059669)};

    return _card(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      _cardHeader('Recent Expenses', Icons.receipt_long_rounded, const Color(0xFFE88A1A)),
      const SizedBox(height: 16),
      ...expenses.map((e) {
        final color = bucketColor[e['bucket']] ?? Colors.grey;
        return Padding(
          padding: const EdgeInsets.only(bottom: 14),
          child: Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
            Row(children: [
              Container(
                padding: const EdgeInsets.all(10),
                decoration: BoxDecoration(color: color.withOpacity(0.1), borderRadius: BorderRadius.circular(12)),
                child: Icon(Icons.receipt_outlined, color: color, size: 18),
              ),
              const SizedBox(width: 12),
              Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
                Text(e['name'] ?? '', style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 14, color: Color(0xFF1E293B))),
                Text(e['date'] ?? '', style: const TextStyle(color: Color(0xFF94A3B8), fontSize: 11)),
              ]),
            ]),
            Text('₹${(e['amount'] ?? 0).toStringAsFixed(0)}',
              style: TextStyle(fontWeight: FontWeight.w900, fontSize: 15, color: color)),
          ]),
        );
      }),
    ]));
  }

  // ── GOALS ──────────────────────────────────────────────────────────────────
  Widget _buildGoalsCard(FinancialProvider p, bool loaded) {
    if (!loaded || p.goals.isEmpty) return const SizedBox.shrink();

    return _card(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      _cardHeader('Goals Progress', Icons.flag_rounded, const Color(0xFF059669)),
      const SizedBox(height: 16),
      ...p.goals.map((g) {
        final target  = (g['target_amount'] ?? 1.0).toDouble();
        final current = (g['current_amount'] ?? 0.0).toDouble();
        final pct     = (current / (target > 0 ? target : 1)).clamp(0.0, 1.0);
        final color   = _parseColor(g['color'] ?? '#059669');
        return Padding(
          padding: const EdgeInsets.only(bottom: 16),
          child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
            Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
              Text(g['name'] ?? '', style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 14, color: Color(0xFF1E293B))),
              Text('${(pct * 100).toInt()}%', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 14, color: color)),
            ]),
            const SizedBox(height: 8),
            ClipRRect(
              borderRadius: BorderRadius.circular(8),
              child: LinearProgressIndicator(
                value: pct,
                minHeight: 8,
                color: color,
                backgroundColor: color.withOpacity(0.12),
              ),
            ),
            const SizedBox(height: 4),
            Text('₹${_formatNum(current)} of ₹${_formatNum(target)}',
              style: const TextStyle(color: Color(0xFF94A3B8), fontSize: 11)),
          ]),
        );
      }),
    ]));
  }

  // ── Shared helpers ─────────────────────────────────────────────────────────
  Widget _card({required Widget child}) => Container(
    padding: const EdgeInsets.all(20),
    decoration: BoxDecoration(
      color: Colors.white,
      borderRadius: BorderRadius.circular(24),
      boxShadow: [BoxShadow(color: const Color(0xFF1E3A8A).withOpacity(0.05), blurRadius: 20, offset: const Offset(0, 8))],
    ),
    child: child,
  );

  Widget _cardHeader(String title, IconData icon, Color color) => Row(children: [
    Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(color: color.withOpacity(0.1), borderRadius: BorderRadius.circular(12)),
      child: Icon(icon, color: color, size: 18),
    ),
    const SizedBox(width: 10),
    Text(title, style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w800, color: Color(0xFF1E293B))),
  ]);

  String _formatNum(double v) {
    if (v >= 10000000) return '${(v / 10000000).toStringAsFixed(1)}Cr';
    if (v >= 100000)  return '${(v / 100000).toStringAsFixed(1)}L';
    if (v >= 1000)    return '${(v / 1000).toStringAsFixed(1)}K';
    return v.toStringAsFixed(0);
  }
}
