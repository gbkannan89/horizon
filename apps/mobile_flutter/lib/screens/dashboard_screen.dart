import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import '../providers/auth_provider.dart';
import 'financial_score_screen.dart';

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
    final authProvider = context.watch<AuthProvider>();
    final userName = authProvider.user?['name'] ?? 'Guest';
    final userInitial = userName.isNotEmpty ? userName.substring(0, 1).toUpperCase() : 'G';
    final loaded = !provider.isLoading;

    return Scaffold(
      backgroundColor: const Color(0xFFF1F5F9),
      body: RefreshIndicator(
        onRefresh: provider.loadAllData,
        color: const Color(0xFF1E3A8A),
        child: CustomScrollView(
          slivers: [
            // ── Hero gradient header ──────────────────────────────────────────
            SliverToBoxAdapter(child: _buildHeroHeader(provider, userName, userInitial, loaded)),

            // ── Body cards ───────────────────────────────────────────────────
            SliverPadding(
              padding: const EdgeInsets.fromLTRB(20, 0, 20, 100),
              sliver: SliverList(
                delegate: SliverChildListDelegate([
                  const SizedBox(height: 20),
                  _buildAnimatedSection(_buildIncomeSpentRow(provider, loaded), 100),
                  const SizedBox(height: 20),
                  _buildAnimatedSection(_buildBudgetCard(provider, loaded), 200),
                  const SizedBox(height: 20),
                  _buildAnimatedSection(_buildFinancialInsights(provider, loaded), 300),
                  const SizedBox(height: 20),
                  _buildAnimatedSection(_buildColdPurchaseCountdown(provider, loaded), 400),
                  const SizedBox(height: 20),
                  _buildAnimatedSection(_buildGoalsCard(provider, loaded), 500),
                  const SizedBox(height: 20),
                  _buildAnimatedSection(_buildUpcomingBills(provider, loaded), 600),
                ]),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildAnimatedSection(Widget child, int delay) {
    return TweenAnimationBuilder<double>(
      duration: Duration(milliseconds: 600 + delay),
      curve: Curves.easeOutCubic,
      tween: Tween(begin: 0.0, end: 1.0),
      builder: (context, value, child) => Transform.translate(
        offset: Offset(0, 20 * (1 - value)),
        child: Opacity(opacity: value, child: child),
      ),
      child: child,
    );
  }

  // ── HERO HEADER ────────────────────────────────────────────────────────────
  Widget _buildHeroHeader(FinancialProvider p, String userName, String userInitial, bool loaded) {
    return TweenAnimationBuilder<double>(
      duration: const Duration(milliseconds: 600),
      curve: Curves.easeOutCubic,
      tween: Tween(begin: 0.0, end: 1.0),
      builder: (context, value, child) => Transform.translate(
        offset: Offset(0, -20 * (1 - value)),
        child: Opacity(
          opacity: value,
          child: Container(
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
            Text('$userName 👋', style: const TextStyle(color: Colors.white, fontSize: 22, fontWeight: FontWeight.w800)),
          ]),
          Container(
            padding: const EdgeInsets.all(10),
            decoration: BoxDecoration(color: Colors.white.withOpacity(0.12), borderRadius: BorderRadius.circular(14)),
            child: const Icon(Icons.notifications_none_rounded, color: Colors.white, size: 22),
          ),
        ]),

        const SizedBox(height: 28),

        // Score + Net Worth row
        Row(children: [
          // Financial Score ring
          GestureDetector(
            onTap: () {
              Navigator.push(context, MaterialPageRoute(builder: (_) => const FinancialScoreScreen()));
            },
            child: Stack(alignment: Alignment.center, children: [
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
          ),
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
            if (loaded && p.netWorthChange != 0) ...[
              const SizedBox(height: 4),
              Row(
                children: [
                  Icon(
                    p.netWorthChange > 0 ? Icons.arrow_upward_rounded : Icons.arrow_downward_rounded,
                    color: p.netWorthChange > 0 ? const Color(0xFF34D399) : const Color(0xFFF87171),
                    size: 14,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    '${p.netWorthChange > 0 ? '+' : '-'}₹${_formatNum(p.netWorthChange.abs())} ${p.netWorthChangePeriod}',
                    style: TextStyle(
                      color: p.netWorthChange > 0 ? const Color(0xFF34D399) : const Color(0xFFF87171),
                      fontSize: 12,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ],
              ),
            ],
          ])),
        ]),
      ]),
    ),
        ),
      ),
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

  // ── FINANCIAL INSIGHTS ─────────────────────────────────────────────────────
  Widget _buildFinancialInsights(FinancialProvider p, bool loaded) {
    if (!loaded) return const SizedBox.shrink();
    
    double savingsRate = p.totalIncomeAgg > 0 ? (p.savingsSpent / p.totalIncomeAgg * 100) : 0;
    
    return Container(
      padding: const EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24),
        boxShadow: [BoxShadow(color: const Color(0xFF1E3A8A).withOpacity(0.06), blurRadius: 24, offset: const Offset(0, 10))],
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [
          Container(
            padding: const EdgeInsets.all(10),
            decoration: BoxDecoration(color: const Color(0xFF3B82F6).withOpacity(0.1), borderRadius: BorderRadius.circular(12)),
            child: const Icon(Icons.insights_rounded, color: Color(0xFF3B82F6), size: 20),
          ),
          const SizedBox(width: 12),
          const Text('Financial Insights', style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold, color: Color(0xFF1E293B))),
        ]),
        const SizedBox(height: 20),
        _insightRow('Savings Rate', '${savingsRate.toStringAsFixed(1)}%', 'of your income goes to savings.', Icons.savings_outlined, const Color(0xFF059669)),
        const Padding(padding: EdgeInsets.symmetric(vertical: 12), child: Divider(height: 1, color: Color(0xFFF1F5F9))),
        _insightRow('Needs Consumption', '${p.totalIncomeAgg > 0 ? (p.needsSpent / p.totalIncomeAgg * 100).toStringAsFixed(1) : 0}%', 'of your income is spent on needs.', Icons.home_outlined, const Color(0xFFE88A1A)),
      ]),
    );
  }

  Widget _insightRow(String title, String value, String subtitle, IconData icon, Color color) {
    return Row(children: [
      Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(color: color.withOpacity(0.1), shape: BoxShape.circle),
        child: Icon(icon, color: color, size: 16),
      ),
      const SizedBox(width: 12),
      Expanded(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Text(title, style: const TextStyle(fontWeight: FontWeight.w700, fontSize: 13, color: Color(0xFF475569))),
        Row(children: [
          Text(value, style: TextStyle(fontWeight: FontWeight.w900, fontSize: 15, color: color)),
          const SizedBox(width: 6),
          Expanded(child: Text(subtitle, style: const TextStyle(fontSize: 12, color: Colors.grey), overflow: TextOverflow.ellipsis)),
        ]),
      ])),
    ]);
  }

  // ── COLD PURCHASE COUNTDOWN ────────────────────────────────────────────────
  Widget _buildColdPurchaseCountdown(FinancialProvider p, bool loaded) {
    if (!loaded || p.wishlist.isEmpty) return const SizedBox.shrink();

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Padding(
          padding: EdgeInsets.only(left: 4, bottom: 12),
          child: Text('Cold Purchase Countdown', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16, color: Color(0xFF1E293B))),
        ),
        SizedBox(
          height: 140,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: p.wishlist.length,
            itemBuilder: (ctx, i) {
              final w = p.wishlist[i];
              final daysLeft = w.unlockDate.difference(DateTime.now()).inDays;
              final unlocked = daysLeft <= 0;
              
              return Container(
                width: 160,
                margin: const EdgeInsets.only(right: 16),
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  gradient: unlocked 
                    ? const LinearGradient(colors: [Color(0xFF059669), Color(0xFF10B981)])
                    : const LinearGradient(colors: [Color(0xFF1E3A8A), Color(0xFF3B82F6)]),
                  borderRadius: BorderRadius.circular(20),
                  boxShadow: [BoxShadow(color: unlocked ? const Color(0xFF059669).withOpacity(0.3) : const Color(0xFF1E3A8A).withOpacity(0.3), blurRadius: 12, offset: const Offset(0, 6))],
                ),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Icon(unlocked ? Icons.lock_open_rounded : Icons.lock_outline_rounded, color: Colors.white70, size: 24),
                    const Spacer(),
                    Text(w.name, style: const TextStyle(color: Colors.white, fontWeight: FontWeight.bold, fontSize: 15), maxLines: 1, overflow: TextOverflow.ellipsis),
                    const SizedBox(height: 4),
                    Text(unlocked ? 'Ready to buy!' : '$daysLeft days left', style: const TextStyle(color: Colors.white70, fontWeight: FontWeight.w600, fontSize: 12)),
                  ],
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  // ── GOALS ──────────────────────────────────────────────────────────────────
  Widget _buildGoalsCard(FinancialProvider p, bool loaded) {
    return _card(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      _cardHeader('Financial Goals', Icons.flag_rounded, const Color(0xFF059669)),
      const SizedBox(height: 16),
      if (loaded && p.goals.isEmpty)
        const Text('No goals set yet.', style: TextStyle(color: Colors.grey, fontSize: 13))
      else if (loaded)
        ...p.goals.map((g) => _buildGoalRow(g['name'], g['current_amount'], g['target_amount'], g['color']))
      else
        const Center(child: CircularProgressIndicator()),
    ]));
  }

  Widget _buildGoalRow(String title, num current, num target, String colorHex) {
    Color c = _parseColor(colorHex);
    double pct = target > 0 ? (current / target).clamp(0.0, 1.0) : 0.0;
    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
          Text(title, style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13)),
          Text('₹${_formatNum(current.toDouble())} / ₹${_formatNum(target.toDouble())}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
        ]),
        const SizedBox(height: 8),
        LinearPercentIndicator(
          lineHeight: 8.0,
          percent: pct,
          progressColor: c,
          backgroundColor: c.withOpacity(0.15),
          barRadius: const Radius.circular(8),
          padding: EdgeInsets.zero,
        ),
      ]),
    );
  }

  // ── UPCOMING BILLS ──────────────────────────────────────────────────────────
  Widget _buildUpcomingBills(FinancialProvider p, bool loaded) {
    if (!loaded) return const SizedBox.shrink();
    if (p.upcomingBills.isEmpty) return const SizedBox.shrink();

    final now = DateTime.now();
    return _card(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      _cardHeader('Upcoming Bills', Icons.event_note_rounded, const Color(0xFFEF4444)),
      const SizedBox(height: 16),
      ...p.upcomingBills.map((bill) {
        int dueDay = bill['due_day'] ?? 1;
        int daysLeft = dueDay - now.day;
        if (daysLeft < 0) daysLeft += 30; // rough estimation for next month
        
        return Padding(
          padding: const EdgeInsets.only(bottom: 12.0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(color: const Color(0xFFEF4444).withOpacity(0.1), shape: BoxShape.circle),
                    child: const Icon(Icons.receipt_long_outlined, color: Color(0xFFEF4444), size: 16),
                  ),
                  const SizedBox(width: 12),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(bill['name'] ?? 'Bill', style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 14)),
                      Text(
                        daysLeft == 0 ? 'Due Today' : 'Due in $daysLeft day${daysLeft > 1 ? 's' : ''}',
                        style: TextStyle(color: daysLeft <= 3 ? const Color(0xFFEF4444) : Colors.grey, fontSize: 12, fontWeight: FontWeight.w500),
                      ),
                    ],
                  ),
                ],
              ),
              Text('₹${_formatNum((bill['amount'] ?? 0).toDouble())}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14)),
            ],
          ),
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
