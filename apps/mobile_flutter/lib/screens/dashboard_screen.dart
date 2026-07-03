import 'dart:math' as math;
import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:provider/provider.dart';
import 'package:fl_chart/fl_chart.dart';
import '../providers/financial_provider.dart';
import '../providers/auth_provider.dart';
import 'financial_score_screen.dart';
import 'nudge_center_screen.dart';

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
                  _buildAnimatedSection(_buildSpendingTrend(provider, loaded), 250),
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
          GestureDetector(
            onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const NudgeCenterScreen())),
            child: Stack(
              children: [
                Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(color: Colors.white.withValues(alpha: 0.12), borderRadius: BorderRadius.circular(14)),
                  child: const Icon(Icons.notifications_none_rounded, color: Colors.white, size: 22),
                ),
                if (loaded && p.nudges.isNotEmpty)
                  Positioned(
                    right: 4, top: 4,
                    child: Container(
                      padding: const EdgeInsets.all(4),
                      decoration: const BoxDecoration(color: Color(0xFFEF4444), shape: BoxShape.circle),
                      child: Text(
                        '${p.nudges.length > 9 ? '9+' : p.nudges.length}',
                        style: const TextStyle(color: Colors.white, fontSize: 9, fontWeight: FontWeight.w800, height: 1),
                      ),
                    ),
                  ),
              ],
            ),
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
                backgroundColor: Colors.white.withValues(alpha: 0.15),
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
        boxShadow: [BoxShadow(color: color.withValues(alpha: 0.08), blurRadius: 16, offset: const Offset(0, 6))],
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Container(
          padding: const EdgeInsets.all(7),
          decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(10)),
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
        backgroundColor: color.withValues(alpha: 0.1),
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

  // ── SPENDING TREND ─────────────────────────────────────────────────────────
  Widget _buildSpendingTrend(FinancialProvider p, bool loaded) {
    if (!loaded || p.spendingTrend.isEmpty) return const SizedBox.shrink();

    final trend = p.spendingTrend;
    double maxVal = 0;
    for (final t in trend) {
      maxVal = [maxVal, (t['needs'] ?? 0).toDouble(), (t['wants'] ?? 0).toDouble(), (t['savings'] ?? 0).toDouble()].reduce((a, b) => a > b ? a : b);
    }

    return _card(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      _cardHeader('Spending Trend', Icons.trending_up_rounded, const Color(0xFF6B46C1)),
      const SizedBox(height: 20),
      SizedBox(
        height: 180,
        child: LineChart(
          LineChartData(
            gridData: FlGridData(show: true, drawVerticalLine: false, horizontalInterval: maxVal > 0 ? maxVal / 4 : 1,
              getDrawingHorizontalLine: (v) => FlLine(color: Colors.grey.shade100, strokeWidth: 1),
            ),
            titlesData: FlTitlesData(
              leftTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
              bottomTitles: AxisTitles(sideTitles: SideTitles(
                showTitles: true, interval: 1,
                getTitlesWidget: (v, _) {
                  final i = v.toInt();
                  if (i < 0 || i >= trend.length) return const SizedBox();
                  return Padding(padding: const EdgeInsets.only(top: 8), child: Text(trend[i]['month'] ?? '', style: const TextStyle(fontSize: 10, color: Colors.grey, fontWeight: FontWeight.w600)));
                },
              )),
              rightTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
              topTitles: AxisTitles(sideTitles: SideTitles(showTitles: false)),
            ),
            borderData: FlBorderData(show: false),
            lineBarsData: [
              _trendLine(trend, 'needs', const Color(0xFFE88A1A)),
              _trendLine(trend, 'wants', const Color(0xFF6B46C1)),
              _trendLine(trend, 'savings', const Color(0xFF059669)),
            ],
          ),
        ),
      ),
      const SizedBox(height: 12),
      Row(mainAxisAlignment: MainAxisAlignment.center, children: [
        _legendDot('Needs', const Color(0xFFE88A1A)),
        const SizedBox(width: 16),
        _legendDot('Wants', const Color(0xFF6B46C1)),
        const SizedBox(width: 16),
        _legendDot('Savings', const Color(0xFF059669)),
      ]),
    ]));
  }

  LineChartBarData _trendLine(List<dynamic> data, String key, Color color) {
    return LineChartBarData(
      spots: data.asMap().entries.map((e) => FlSpot(e.key.toDouble(), (e.value[key] ?? 0).toDouble())).toList(),
      isCurved: true,
      color: color,
      barWidth: 2.5,
      isStrokeCapRound: true,
      dotData: FlDotData(show: false),
      belowBarData: BarAreaData(show: true, gradient: LinearGradient(
        colors: [color.withValues(alpha: 0.15), color.withValues(alpha: 0.01)],
        begin: Alignment.topCenter, end: Alignment.bottomCenter,
      )),
    );
  }

  Widget _legendDot(String label, Color color) {
    return Row(mainAxisSize: MainAxisSize.min, children: [
      Container(width: 8, height: 8, decoration: BoxDecoration(color: color, shape: BoxShape.circle)),
      const SizedBox(width: 4),
      Text(label, style: TextStyle(fontSize: 11, color: Colors.grey.shade600, fontWeight: FontWeight.w600)),
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
        boxShadow: [BoxShadow(color: const Color(0xFF1E3A8A).withValues(alpha: 0.06), blurRadius: 24, offset: const Offset(0, 10))],
      ),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(children: [
          Container(
            padding: const EdgeInsets.all(10),
            decoration: BoxDecoration(color: const Color(0xFF3B82F6).withValues(alpha: 0.1), borderRadius: BorderRadius.circular(12)),
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
        decoration: BoxDecoration(color: color.withValues(alpha: 0.1), shape: BoxShape.circle),
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
        Padding(
          padding: const EdgeInsets.only(left: 4, bottom: 16),
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(8),
                decoration: BoxDecoration(
                  color: const Color(0xFF0891B2).withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(10),
                ),
                child: const Icon(Icons.ac_unit_rounded, color: Color(0xFF0891B2), size: 18),
              ),
              const SizedBox(width: 10),
              const Text('Cold Purchase Countdown', style: TextStyle(fontWeight: FontWeight.w800, fontSize: 16, color: Color(0xFF1E293B))),
            ],
          ),
        ),
        SizedBox(
          height: 270,
          child: ListView.builder(
            scrollDirection: Axis.horizontal,
            itemCount: p.wishlist.length,
            itemBuilder: (ctx, i) {
              final w = p.wishlist[i];
              final now = DateTime.now();
              final diff = w.unlockDate.difference(now);
              final totalDays = w.unlockDate.difference(w.addedDate).inDays;
              final daysLeft = diff.inDays;
              final hoursLeft = diff.inHours;
              final unlocked = daysLeft <= 0;
              final progress = totalDays > 0 ? (daysLeft / totalDays).clamp(0.0, 1.0) : 0.0;
              final isCloseToUnlock = !unlocked && daysLeft <= 3;

              return TweenAnimationBuilder<double>(
                tween: Tween(begin: 0, end: 1),
                duration: const Duration(milliseconds: 800),
                curve: Curves.easeOutCubic,
                builder: (context, anim, _) => Transform.translate(
                  offset: Offset(40 * (1 - anim), 0),
                  child: Opacity(
                    opacity: anim,
                    child: Container(
                      width: 200,
                      margin: const EdgeInsets.only(right: 16),
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        gradient: unlocked
                          ? const LinearGradient(colors: [Color(0xFFD1FAE5), Color(0xFFECFDF5)])
                          : isCloseToUnlock
                            ? const LinearGradient(colors: [Color(0xFFFEF3C7), Color(0xFFFFF7ED)])
                            : const LinearGradient(colors: [Color(0xFFE0E7FF), Color(0xFFEDE9FE)]),
                        borderRadius: BorderRadius.circular(24),
                        boxShadow: [
                          BoxShadow(
                            color: unlocked
                              ? const Color(0xFF059669).withValues(alpha: 0.12)
                              : isCloseToUnlock
                                ? const Color(0xFFD97706).withValues(alpha: 0.12)
                                : const Color(0xFF6366F1).withValues(alpha: 0.10),
                            blurRadius: 16,
                            offset: const Offset(0, 8),
                          ),
                        ],
                      ),
                      child: Column(
                        children: [
                          Row(
                            children: [
                              Icon(
                                unlocked ? Icons.lock_open_rounded : Icons.lock_outline_rounded,
                                color: unlocked ? const Color(0xFF059669) : isCloseToUnlock ? const Color(0xFFB45309) : const Color(0xFF6366F1),
                                size: 18,
                              ),
                              const Spacer(),
                              if (isCloseToUnlock)
                                _pulsingDot(),
                            ],
                          ),
                          const Spacer(),
                          // Circular countdown
                          Center(
                            child: SizedBox(
                              width: 130, height: 130,
                              child: Stack(
                                alignment: Alignment.center,
                                children: [
                                  unlocked
                                    ? _celebrationRing(size: 130)
                                    : TweenAnimationBuilder<double>(
                                        tween: Tween(begin: 1, end: progress),
                                        duration: const Duration(milliseconds: 1500),
                                        curve: Curves.easeOutCubic,
                                        builder: (context, val, _) => SizedBox(
                                          width: 130,
                                          height: 130,
                                          child: CircularProgressIndicator(
                                            value: val.clamp(0.0, 1.0),
                                            strokeWidth: 6,
                                            strokeCap: StrokeCap.round,
                                            backgroundColor: unlocked
                                              ? const Color(0xFF059669).withValues(alpha: 0.12)
                                              : isCloseToUnlock
                                                ? const Color(0xFFD97706).withValues(alpha: 0.12)
                                                : const Color(0xFF6366F1).withValues(alpha: 0.12),
                                            valueColor: AlwaysStoppedAnimation<Color>(
                                              unlocked
                                                ? const Color(0xFF10B981)
                                                : isCloseToUnlock
                                                  ? const Color(0xFFF59E0B)
                                                  : const Color(0xFF818CF8),
                                            ),
                                          ),
                                        ),
                                      ),
                                  Column(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      TweenAnimationBuilder<int>(
                                        tween: IntTween(begin: totalDays, end: daysLeft),
                                        duration: const Duration(milliseconds: 1500),
                                        curve: Curves.easeOutCubic,
                                        builder: (context, val, _) => Text(
                                          '${unlocked ? 0 : val}',
                                          style: TextStyle(
                                            fontSize: 28,
                                            fontWeight: FontWeight.w900,
                                            color: unlocked
                                              ? const Color(0xFF059669)
                                              : isCloseToUnlock
                                                ? const Color(0xFFB45309)
                                                : const Color(0xFF4338CA),
                                          ),
                                        ),
                                      ),
                                      Text(
                                        unlocked ? '' : 'days',
                                        style: TextStyle(
                                          fontSize: 10,
                                          color: unlocked
                                            ? const Color(0xFF059669).withValues(alpha: 0.6)
                                            : isCloseToUnlock
                                              ? const Color(0xFFB45309).withValues(alpha: 0.6)
                                              : const Color(0xFF4338CA).withValues(alpha: 0.6),
                                        ),
                                      ),
                                    ],
                                  ),
                                ],
                              ),
                            ),
                          ),
                          const Spacer(),
                          Text(
                            w.name,
                            style: TextStyle(
                              color: unlocked
                                ? const Color(0xFF065F46)
                                : isCloseToUnlock
                                  ? const Color(0xFF92400E)
                                  : const Color(0xFF1E1B4B),
                              fontWeight: FontWeight.bold, fontSize: 13,
                            ),
                            maxLines: 1, overflow: TextOverflow.ellipsis,
                          ),
                          const SizedBox(height: 2),
                          AnimatedSwitcher(
                            duration: const Duration(milliseconds: 400),
                            child: Text(
                              unlocked
                                ? '🎉 Ready to buy!'
                                : isCloseToUnlock
                                  ? '$hoursLeft hours left'
                                  : '$daysLeft days left',
                              key: ValueKey('${w.id}_$daysLeft'),
                              style: TextStyle(
                                color: unlocked
                                  ? const Color(0xFF047857)
                                  : isCloseToUnlock
                                    ? const Color(0xFFB45309)
                                    : const Color(0xFF6366F1),
                                fontWeight: FontWeight.w600, fontSize: 11,
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _pulsingDot() {
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0.0, end: 1.0),
      duration: const Duration(milliseconds: 1000),
      builder: (context, value, _) {
        return Container(
          width: 8,
          height: 8,
          decoration: BoxDecoration(
            color: Color.lerp(const Color(0xFFFBBF24), const Color(0xFFF59E0B), value)!,
            shape: BoxShape.circle,
            boxShadow: [
              BoxShadow(
                color: const Color(0xFFFBBF24).withValues(alpha: (0.5 + value * 0.5)),
                blurRadius: 4 + value * 4,
              ),
            ],
          ),
        );
      },
    );
  }

  Widget _celebrationRing({double size = 130}) {
    return TweenAnimationBuilder<double>(
      tween: Tween(begin: 0.0, end: 1.0),
      duration: const Duration(milliseconds: 600),
      curve: Curves.elasticOut,
      builder: (context, value, _) {
        return CustomPaint(
          size: Size(size, size),
          painter: _CelebrationPainter(value),
        );
      },
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
        ...p.goals.map((g) => _buildGoalRow(g))
      else
        const Center(child: CircularProgressIndicator()),
    ]));
  }

  Widget _buildGoalRow(dynamic g) {
    final title = g['name'] ?? '';
    final current = (g['current_amount'] ?? 0).toDouble();
    final target = (g['target_amount'] ?? 0).toDouble();
    final colorHex = g['color'] ?? '#059669';
    final projectedDate = g['projected_completion_date'];
    final monthlyNeeded = (g['monthly_saving_needed'] ?? 0).toDouble();

    Color c = _parseColor(colorHex);
    double pct = target > 0 ? (current / target).clamp(0.0, 1.0) : 0.0;

    String projectionText = '';
    if (projectedDate != null && monthlyNeeded > 0) {
      try {
        final dt = DateTime.parse(projectedDate.toString());
        const months = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
        projectionText = '📅 ${dt.day} ${months[dt.month - 1]} ${dt.year}  ·  ₹${_formatNum(monthlyNeeded)}/mo';
      } catch (_) {
        projectionText = '';
      }
    }

    return Padding(
      padding: const EdgeInsets.only(bottom: 16.0),
      child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
        Row(mainAxisAlignment: MainAxisAlignment.spaceBetween, children: [
          Text(title, style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 13)),
          Text('₹${_formatNum(current)} / ₹${_formatNum(target)}', style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 13)),
        ]),
        const SizedBox(height: 8),
        LinearPercentIndicator(
          lineHeight: 8.0,
          percent: pct,
          progressColor: c,
          backgroundColor: c.withValues(alpha: 0.15),
          barRadius: const Radius.circular(8),
          padding: EdgeInsets.zero,
        ),
        if (projectionText.isNotEmpty) ...[
          const SizedBox(height: 6),
          Text(projectionText, style: TextStyle(color: c, fontSize: 11, fontWeight: FontWeight.w600)),
        ],
      ]),
    );
  }

  // ── UPCOMING BILLS ──────────────────────────────────────────────────────────
  Widget _buildUpcomingBills(FinancialProvider p, bool loaded) {
    if (!loaded) return const SizedBox.shrink();

    final now = DateTime.now();
    final daysInMonth = DateTime(now.year, now.month + 1, 0).day;

    final List<Map<String, dynamic>> sortedBills = p.upcomingBills.map((b) {
      int dueDay = b['due_day'] ?? 1;
      DateTime dueDate = DateTime(now.year, now.month, dueDay.clamp(1, daysInMonth));
      if (dueDate.isBefore(now)) {
        dueDate = DateTime(now.year, now.month + 1, dueDay.clamp(1, daysInMonth));
      }
      int daysLeft = dueDate.difference(now).inDays;
      return {
        'bill': b,
        'dueDate': dueDate,
        'daysLeft': daysLeft,
      };
    }).toList();
    sortedBills.sort((a, b) => (a['daysLeft'] as int).compareTo(b['daysLeft'] as int));

    return _card(child: Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      _cardHeader('Upcoming Bills', Icons.event_note_rounded, const Color(0xFFEF4444)),
      const SizedBox(height: 16),
      if (sortedBills.isEmpty)
        const Padding(
          padding: EdgeInsets.symmetric(vertical: 8),
          child: Row(children: [
            Icon(Icons.check_circle_outline, color: Color(0xFF059669), size: 18),
            SizedBox(width: 8),
            Text('No upcoming bills', style: TextStyle(color: Color(0xFF059669), fontWeight: FontWeight.w600, fontSize: 14)),
          ]),
        )
      else
        ...sortedBills.take(6).map((item) {
          final bill = item['bill'] as Map<String, dynamic>;
          final daysLeft = item['daysLeft'] as int;
          final bucket = bill['bucket'] ?? 'Needs';
          final category = bill['category'] ?? '';
          final name = bill['name'] ?? 'Bill';
          final amount = (bill['amount'] ?? 0).toDouble();

          final bucketColor = _billBucketColor(bucket);
          final billIcon = _billCategoryIcon(category, bucket);
          final dueText = daysLeft == 0 ? 'Due Today' : daysLeft == 1 ? 'Due Tomorrow' : '$daysLeft days';
          final urgencyColor = daysLeft == 0 ? const Color(0xFFEF4444) : daysLeft <= 3 ? const Color(0xFFF59E0B) : const Color(0xFF94A3B8);
          final progress = daysLeft <= 30 ? (30 - daysLeft) / 30.0 : 0.0;

          return Container(
            margin: const EdgeInsets.only(bottom: 10),
            child: Row(
              children: [
                Container(
                  width: 44, height: 44,
                  decoration: BoxDecoration(
                    color: bucketColor.withValues(alpha: 0.12),
                    borderRadius: BorderRadius.circular(14),
                  ),
                  child: Icon(billIcon, color: bucketColor, size: 20),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Expanded(
                            child: Text(name,
                                style: const TextStyle(fontWeight: FontWeight.w600, fontSize: 14, color: Color(0xFF1E293B)),
                                overflow: TextOverflow.ellipsis),
                          ),
                          Text('₹${_formatNum(amount)}',
                              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 14, color: Color(0xFF1E293B))),
                        ],
                      ),
                      const SizedBox(height: 4),
                      Row(
                        children: [
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                            decoration: BoxDecoration(
                              color: bucketColor.withValues(alpha: 0.1),
                              borderRadius: BorderRadius.circular(6),
                            ),
                            child: Text(bucket,
                                style: TextStyle(fontSize: 9, fontWeight: FontWeight.w600, color: bucketColor)),
                          ),
                          const SizedBox(width: 6),
                          if (category.isNotEmpty)
                            Text(category,
                                style: const TextStyle(fontSize: 10, color: Color(0xFF94A3B8))),
                          const Spacer(),
                          Icon(Icons.schedule, size: 11, color: urgencyColor),
                          const SizedBox(width: 3),
                          Text(dueText,
                              style: TextStyle(fontSize: 11, fontWeight: FontWeight.w600, color: urgencyColor)),
                        ],
                      ),
                      const SizedBox(height: 6),
                      ClipRRect(
                        borderRadius: BorderRadius.circular(4),
                        child: LinearProgressIndicator(
                          value: progress.clamp(0.0, 1.0),
                          backgroundColor: bucketColor.withValues(alpha: 0.08),
                          valueColor: AlwaysStoppedAnimation<Color>(
                            daysLeft == 0 ? const Color(0xFFEF4444) : bucketColor,
                          ),
                          minHeight: 3,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          );
        }),
    ]));
  }

  Color _billBucketColor(String bucket) {
    switch (bucket) {
      case 'Savings': return const Color(0xFF059669);
      case 'Wants': return const Color(0xFF6B46C1);
      default: return const Color(0xFFE88A1A);
    }
  }

  IconData _billCategoryIcon(String category, String bucket) {
    switch (category.toUpperCase()) {
      case 'RENT':
      case 'HOUSING': return Icons.home_outlined;
      case 'ELECTRICITY':
      case 'UTILITIES':
      case 'WATER': return Icons.bolt_outlined;
      case 'INTERNET': return Icons.wifi_outlined;
      case 'INSURANCE': return Icons.health_and_safety_outlined;
      case 'EDUCATION': return Icons.school_outlined;
      case 'ENTERTAINMENT':
      case 'SUBSCRIPTION': return Icons.tv_outlined;
      case 'DEBT':
      case 'LOAN': return Icons.account_balance_outlined;
      case 'GROCERIES':
      case 'FOOD': return Icons.restaurant_outlined;
      case 'TRANSPORT': return Icons.directions_car_outlined;
      case 'INVESTMENT':
      case 'SAVINGS': return Icons.savings_outlined;
      default:
        if (bucket == 'Savings') return Icons.savings_outlined;
        if (bucket == 'Wants') return Icons.shopping_bag_outlined;
        return Icons.receipt_long_outlined;
    }
  }

  // ── Shared helpers ─────────────────────────────────────────────────────────
  Widget _card({required Widget child}) => Container(
    padding: const EdgeInsets.all(20),
    decoration: BoxDecoration(
      color: Colors.white,
      borderRadius: BorderRadius.circular(24),
      boxShadow: [BoxShadow(color: const Color(0xFF1E3A8A).withValues(alpha: 0.05), blurRadius: 20, offset: const Offset(0, 8))],
    ),
    child: child,
  );

  Widget _cardHeader(String title, IconData icon, Color color) => Row(children: [
    Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(color: color.withValues(alpha: 0.1), borderRadius: BorderRadius.circular(12)),
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

class _CelebrationPainter extends CustomPainter {
  final double progress;
  _CelebrationPainter(this.progress);

  @override
  void paint(Canvas canvas, Size size) {
    final center = size.center(Offset.zero);
    final radius = size.width / 2;
    final paint = Paint()
      ..color = const Color(0xFF34D399).withValues(alpha: 0.3 * progress)
      ..style = PaintingStyle.fill;

    for (int i = 0; i < 12; i++) {
      final angle = (i / 12) * math.pi * 2;
      final dist = (radius + 15) * progress;
      final dx = math.cos(angle) * dist;
      final dy = math.sin(angle) * dist;
      canvas.drawCircle(
        center + Offset(dx, dy),
        3 + 3 * progress,
        paint..color = [const Color(0xFF34D399), const Color(0xFF6EE7B7), const Color(0xFFA7F3D0), const Color(0xFF10B981)][i % 4].withValues(alpha: 0.5 * progress),
      );
    }

    // Outer celebration ring matching the progress indicator diameter
    canvas.drawCircle(center, (radius - 3) * progress, Paint()..color = const Color(0xFF34D399));
    // Inner hollow center to prevent overlap with the text
    canvas.drawCircle(
      center,
      (radius - 9) * progress,
      Paint()..color = Colors.white.withValues(alpha: 0.8 * progress)..style = PaintingStyle.fill,
    );
  }

  @override
  bool shouldRepaint(_CelebrationPainter old) => old.progress != progress;
}
