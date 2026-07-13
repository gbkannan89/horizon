import 'package:flutter/material.dart';
import 'package:percent_indicator/percent_indicator.dart';
import 'package:provider/provider.dart';
import '../models/wishlist_item.dart';
import '../models/discipline_aggregates.dart';
import '../services/discipline_service.dart';
import '../providers/financial_provider.dart';
import '../utils/ui_utils.dart';
import 'collections_list_screen.dart';

class DisciplineDashboardScreen extends StatefulWidget {
  final DisciplineService disciplineService;

  const DisciplineDashboardScreen({super.key, required this.disciplineService});

  @override
  _DisciplineDashboardScreenState createState() =>
      _DisciplineDashboardScreenState();
}

class _DisciplineDashboardScreenState extends State<DisciplineDashboardScreen>
    with SingleTickerProviderStateMixin {
  bool _isLoading = true;
  PayYourselfFirst? _payYourselfFirst;
  ZeroBasedBudget? _zeroBasedBudget;
  List<WishlistItem> _wishlist = [];
  FinancialGuardrails? _guardrails;
  DebtRepaymentStrategy? _debtStrategy;

  late AnimationController _peekController;
  late Animation<double> _peekAnimation;

  @override
  void initState() {
    super.initState();

    _peekController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 1000),
    );
    _peekAnimation = TweenSequence([
      TweenSequenceItem(
        tween: Tween<double>(
          begin: 0,
          end: -60,
        ).chain(CurveTween(curve: Curves.easeOutCubic)),
        weight: 40,
      ),
      TweenSequenceItem(tween: ConstantTween<double>(-60), weight: 10),
      TweenSequenceItem(
        tween: Tween<double>(
          begin: -60,
          end: 0,
        ).chain(CurveTween(curve: Curves.easeInCubic)),
        weight: 50,
      ),
    ]).animate(_peekController);

    WidgetsBinding.instance.addPostFrameCallback((_) => _loadData());
  }

  @override
  void dispose() {
    _peekController.dispose();
    super.dispose();
  }

  Future<void> _loadData() async {
    setState(() => _isLoading = true);
    try {
      if (mounted) {
        Provider.of<FinancialProvider>(
          context,
          listen: false,
        ).loadCollections();
      }
      final futures = await Future.wait([
        widget.disciplineService.getPayYourselfFirst(),
        widget.disciplineService.getZeroBasedBudget(),
        widget.disciplineService.getWishlist(),
        widget.disciplineService.getGuardrails(),
        widget.disciplineService.getDebtStrategy(),
      ]);

      setState(() {
        _payYourselfFirst = futures[0] as PayYourselfFirst;
        _zeroBasedBudget = futures[1] as ZeroBasedBudget;
        _wishlist = futures[2] as List<WishlistItem>;
        _guardrails = futures[3] as FinancialGuardrails;
        _debtStrategy = futures[4] as DebtRepaymentStrategy;
      });
    } catch (e) {
      if (mounted)
        UiUtils.showSnack(
          context,
          'Error loading discipline data: $e',
          isError: true,
        );
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
        if (_wishlist.isNotEmpty) {
          Future.delayed(const Duration(milliseconds: 600), () {
            if (mounted) _peekController.forward();
          });
        }
      }
    }
  }

  Widget _buildGamifiedCard({
    required String title,
    required Widget child,
    required IconData icon,
    required List<Color> gradient,
  }) {
    return Container(
      margin: const EdgeInsets.only(bottom: 20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(24),
        border: Border.all(
          color: gradient[0].withValues(alpha: 0.2),
          width: 1.5,
        ),
        boxShadow: [
          BoxShadow(
            color: gradient[0].withValues(alpha: 0.08),
            blurRadius: 24,
            offset: const Offset(0, 12),
          ),
        ],
      ),
      child: ClipRRect(
        borderRadius: BorderRadius.circular(24),
        child: Container(
          decoration: BoxDecoration(
            border: Border(left: BorderSide(color: gradient[0], width: 6)),
          ),
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Container(
                    padding: const EdgeInsets.all(10),
                    decoration: BoxDecoration(
                      gradient: LinearGradient(colors: gradient),
                      shape: BoxShape.circle,
                      boxShadow: [
                        BoxShadow(
                          color: gradient[0].withValues(alpha: 0.4),
                          blurRadius: 12,
                          offset: const Offset(0, 4),
                        ),
                      ],
                    ),
                    child: Icon(icon, color: Colors.white, size: 20),
                  ),
                  const SizedBox(width: 14),
                  Text(
                    title,
                    style: const TextStyle(
                      fontSize: 19,
                      fontWeight: FontWeight.w900,
                      color: Color(0xFF1E293B),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 24),
              child,
            ],
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return Scaffold(
        body: Container(
          decoration: const BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [Color(0xFFF0FDFA), Color(0xFFF8FAFC), Color(0xFFF5F3FF)],
            ),
          ),
          child: const Center(
            child: CircularProgressIndicator(color: Color(0xFF0D9488)),
          ),
        ),
      );
    }

    return Scaffold(
      body: Container(
        decoration: const BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [Color(0xFFF0FDFA), Color(0xFFF8FAFC), Color(0xFFF5F3FF)],
          ),
        ),
        child: SafeArea(
          child: RefreshIndicator(
            onRefresh: _loadData,
            child: ListView(
              padding: const EdgeInsets.all(20),
              children: [
                _buildPayYourselfFirstCard(),
                _buildZeroBasedBudgetCard(),
                _buildDebtStrategyCard(),
                _buildGuardrailsCard(),
                const SizedBox(height: 16),
                _buildEmergencyFundCard(),
                const SizedBox(height: 16),
                _buildCollectionsCard(),
                const SizedBox(height: 16),
                _buildWishlistCard(),
                const SizedBox(height: 60),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildPayYourselfFirstCard() {
    if (_payYourselfFirst == null) return const SizedBox();
    final p = _payYourselfFirst!;
    final progress = p.targetAmount > 0
        ? (p.actualSavings / p.targetAmount)
        : 0.0;

    return _buildGamifiedCard(
      title: 'Pay Yourself First',
      icon: Icons.savings_rounded,
      gradient: const [Color(0xFF059669), Color(0xFF10B981)],
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Saved this month',
                    style: TextStyle(
                      color: Color(0xFF64748B),
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    '₹${p.actualSavings.toStringAsFixed(0)}',
                    style: const TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.w900,
                      color: Color(0xFF059669),
                    ),
                  ),
                ],
              ),
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 6,
                ),
                decoration: BoxDecoration(
                  color: p.status == 'On Track'
                      ? Colors.green.shade50
                      : Colors.red.shade50,
                  borderRadius: BorderRadius.circular(20),
                ),
                child: Text(
                  p.status,
                  style: TextStyle(
                    color: p.status == 'On Track'
                        ? Colors.green.shade700
                        : Colors.red.shade700,
                    fontWeight: FontWeight.bold,
                    fontSize: 12,
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 16),
          LinearPercentIndicator(
            lineHeight: 12.0,
            percent: progress.clamp(0.0, 1.0),
            barRadius: const Radius.circular(10),
            backgroundColor: Colors.grey.shade200,
            linearGradient: const LinearGradient(
              colors: [Color(0xFF059669), Color(0xFF34D399)],
            ),
            padding: EdgeInsets.zero,
          ),
          const SizedBox(height: 8),
          Text(
            'Target: 20% of Income (₹${p.targetAmount.toStringAsFixed(0)})',
            style: const TextStyle(
              color: Color(0xFF64748B),
              fontSize: 13,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildZeroBasedBudgetCard() {
    if (_zeroBasedBudget == null) return const SizedBox();
    final z = _zeroBasedBudget!;
    final diff = z.totalIncome - z.totalAllocated;

    return _buildGamifiedCard(
      title: 'Zero-Based Budget',
      icon: Icons.pie_chart_rounded,
      gradient: const [Color(0xFF6B46C1), Color(0xFF8B5CF6)],
      child: Column(
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: [
              CircularPercentIndicator(
                radius: 45.0,
                lineWidth: 10.0,
                percent: z.totalIncome > 0
                    ? (z.totalAllocated / z.totalIncome).clamp(0.0, 1.0)
                    : 0,
                circularStrokeCap: CircularStrokeCap.round,
                linearGradient: const LinearGradient(
                  colors: [Color(0xFF6B46C1), Color(0xFFA78BFA)],
                ),
                backgroundColor: Colors.grey.shade100,
                center: const Icon(
                  Icons.account_balance_wallet_rounded,
                  color: Color(0xFF6B46C1),
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  const Text(
                    'Unallocated',
                    style: TextStyle(
                      color: Color(0xFF64748B),
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                  Text(
                    '₹${diff.toStringAsFixed(0)}',
                    style: TextStyle(
                      fontSize: 24,
                      fontWeight: FontWeight.w900,
                      color: diff == 0 ? Colors.green : Colors.orange,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 10,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: z.status == 'Zero-Based'
                          ? Colors.green.shade50
                          : Colors.orange.shade50,
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      z.status,
                      style: TextStyle(
                        color: z.status == 'Zero-Based'
                            ? Colors.green.shade700
                            : Colors.orange.shade700,
                        fontSize: 11,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildGuardrailsCard() {
    if (_guardrails == null) return const SizedBox();
    final g = _guardrails!;
    return _buildGamifiedCard(
      title: 'Financial Guardrails',
      icon: Icons.shield_rounded,
      gradient: const [Color(0xFF14B8A6), Color(0xFF60A5FA)],
      child: Column(
        children: [
          _buildGuardrailRow(
            'Emergency Runway',
            '${g.emergencyFundRatio.toStringAsFixed(1)} months\nTarget: ₹${g.emergencyTargetAmount.toStringAsFixed(0)}',
            g.runwayStatus,
          ),
          const Padding(
            padding: EdgeInsets.symmetric(vertical: 12),
            child: Divider(height: 1),
          ),
          _buildGuardrailRow(
            'Housing Cost',
            '${g.housingCostRatio.toStringAsFixed(1)}%',
            g.housingStatus,
          ),
        ],
      ),
    );
  }

  Widget _buildGuardrailRow(String label, String value, String status) {
    bool isHealthy = status == 'Healthy';
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              label,
              style: const TextStyle(fontWeight: FontWeight.bold, fontSize: 15),
            ),
            const SizedBox(height: 4),
            Text(
              value,
              style: const TextStyle(
                color: Color(0xFF64748B),
                fontSize: 13,
                fontWeight: FontWeight.w600,
              ),
            ),
          ],
        ),
        Icon(
          isHealthy ? Icons.check_circle_rounded : Icons.warning_rounded,
          color: isHealthy ? Colors.green : Colors.orange,
          size: 28,
        ),
      ],
    );
  }

  Widget _buildDebtStrategyCard() {
    if (_debtStrategy == null) return const SizedBox();
    final d = _debtStrategy!;
    return _buildGamifiedCard(
      title: 'Debt Strategy',
      icon: Icons.trending_down_rounded,
      gradient: const [Color(0xFFDC2626), Color(0xFFF87171)],
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Container(
            width: double.infinity,
            padding: const EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.red.shade50,
              borderRadius: BorderRadius.circular(16),
            ),
            child: Row(
              children: [
                const Icon(
                  Icons.auto_awesome_rounded,
                  color: Color(0xFFDC2626),
                ),
                const SizedBox(width: 12),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Recommended Plan',
                      style: TextStyle(
                        color: Color(0xFF64748B),
                        fontSize: 12,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                    Text(
                      d.recommendedStrategy,
                      style: const TextStyle(
                        color: Color(0xFFDC2626),
                        fontWeight: FontWeight.bold,
                        fontSize: 16,
                      ),
                    ),
                  ],
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildEmergencyFundCard() {
    if (_guardrails == null) return const SizedBox.shrink();
    final targetAmount = _guardrails!.emergencyTargetAmount;
    // Just a placeholder calculation using liquid assets (bank/fd). We don't have direct access here,
    // but the guardrails object usually computes it or we can pass it from FinancialProvider.
    // For now, we'll prompt the user they can add Bank/FD assets in Portfolio to build this up.

    return _buildGamifiedCard(
      title: 'Emergency Fund',
      icon: Icons.health_and_safety_rounded,
      gradient: const [Color(0xFFEAB308), Color(0xFFFACC15)],
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                'Target Amount',
                style: TextStyle(fontWeight: FontWeight.bold, fontSize: 14),
              ),
              Text(
                '₹${targetAmount.toStringAsFixed(0)}',
                style: const TextStyle(
                  fontWeight: FontWeight.w900,
                  fontSize: 18,
                  color: Color(0xFFEAB308),
                ),
              ),
            ],
          ),
          const SizedBox(height: 8),
          const Text(
            'Your emergency fund is automatically calculated from your highly liquid assets (Bank Accounts and Fixed Deposits) added in your Portfolio.',
            style: TextStyle(color: Color(0xFF64748B), fontSize: 12),
          ),
          const SizedBox(height: 12),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: () {
                // Navigate to Portfolio or show Add Asset dialog
              },
              icon: const Icon(Icons.account_balance_rounded, size: 18),
              label: const Text('Manage Liquid Assets'),
              style: OutlinedButton.styleFrom(
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                foregroundColor: const Color(0xFFEAB308),
                side: const BorderSide(color: Color(0xFFFDE047)),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCollectionsCard() {
    return Consumer<FinancialProvider>(
      builder: (context, fp, _) {
        final active = fp.collections
            .where((c) => c['status'] == 'active')
            .toList();
        final totalExpected = active.fold<double>(
          0,
          (s, c) => s + ((c['total_expected'] ?? 0).toDouble()),
        );
        final totalCollected = active.fold<double>(
          0,
          (s, c) => s + ((c['total_collected'] ?? 0).toDouble()),
        );
        final memberCount = active.fold<int>(
          0,
          (s, c) => s + ((c['member_count'] ?? 0) as int),
        );
        final paidCount = active.fold<int>(
          0,
          (s, c) => s + ((c['paid_count'] ?? 0) as int),
        );

        return _buildGamifiedCard(
          title: 'Collections',
          icon: Icons.people_alt_rounded,
          gradient: const [Color(0xFF6B46C1), Color(0xFFA78BFA)],
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              if (active.isEmpty)
                const Padding(
                  padding: EdgeInsets.symmetric(vertical: 16),
                  child: Column(
                    children: [
                      Icon(
                        Icons.people_outline,
                        size: 36,
                        color: Color(0xFF64748B),
                      ),
                      SizedBox(height: 8),
                      Text(
                        'No active collections',
                        style: TextStyle(
                          color: Color(0xFF64748B),
                          fontWeight: FontWeight.w500,
                        ),
                      ),
                      SizedBox(height: 4),
                      Text(
                        'Track group payments from friends & family.',
                        style: TextStyle(
                          color: Color(0xFF64748B),
                          fontSize: 12,
                        ),
                      ),
                    ],
                  ),
                )
              else ...[
                Row(
                  children: [
                    _collectionStat(
                      'Active',
                      '${active.length}',
                      Icons.play_circle_outline,
                      const Color(0xFF6B46C1),
                    ),
                    const SizedBox(width: 20),
                    _collectionStat(
                      'Collected',
                      '₹${_fmt(totalCollected)}',
                      Icons.account_balance_wallet_outlined,
                      const Color(0xFF059669),
                    ),
                    const SizedBox(width: 20),
                    _collectionStat(
                      'Pending',
                      '${memberCount - paidCount}',
                      Icons.pending_actions,
                      const Color(0xFFF59E0B),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                ClipRRect(
                  borderRadius: BorderRadius.circular(6),
                  child: LinearPercentIndicator(
                    lineHeight: 10,
                    percent: totalExpected > 0
                        ? (totalCollected / totalExpected).clamp(0.0, 1.0)
                        : 0.0,
                    progressColor: const Color(0xFF6B46C1),
                    backgroundColor: const Color(
                      0xFF6B46C1,
                    ).withValues(alpha: 0.1),
                    barRadius: const Radius.circular(6),
                    padding: EdgeInsets.zero,
                  ),
                ),
                if (totalExpected > 0)
                  Padding(
                    padding: const EdgeInsets.only(top: 4),
                    child: Text(
                      '₹${_fmt(totalCollected)} of ₹${_fmt(totalExpected)} collected',
                      style: const TextStyle(
                        fontSize: 11,
                        color: Color(0xFF64748B),
                      ),
                    ),
                  ),
              ],
              const SizedBox(height: 16),
              SizedBox(
                width: double.infinity,
                child: OutlinedButton.icon(
                  onPressed: () {
                    Navigator.of(context).push(
                      MaterialPageRoute(
                        builder: (_) => const CollectionsListScreen(),
                      ),
                    );
                  },
                  icon: const Icon(Icons.open_in_new, size: 18),
                  label: Text(
                    active.isEmpty
                        ? 'Create Collection'
                        : 'View All Collections',
                  ),
                  style: OutlinedButton.styleFrom(
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    foregroundColor: const Color(0xFF6B46C1),
                    side: const BorderSide(color: Color(0xFFA78BFA)),
                  ),
                ),
              ),
            ],
          ),
        );
      },
    );
  }

  String _fmt(dynamic val) {
    if (val == null) return '0';
    final n = (val is num) ? val : double.tryParse(val.toString()) ?? 0;
    if (n >= 100000) return '${(n / 100000).toStringAsFixed(1)}L';
    if (n >= 1000) return '${(n / 1000).toStringAsFixed(1)}K';
    return n.toStringAsFixed(0);
  }

  Widget _collectionStat(
    String label,
    String value,
    IconData icon,
    Color color,
  ) {
    return Expanded(
      child: Column(
        children: [
          Icon(icon, color: color, size: 20),
          const SizedBox(height: 4),
          Text(
            value,
            style: TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.w900,
              color: color,
            ),
          ),
          Text(
            label,
            style: const TextStyle(fontSize: 10, color: Color(0xFF64748B)),
          ),
        ],
      ),
    );
  }

  Widget _buildWishlistCard() {
    return _buildGamifiedCard(
      title: 'Wishlist Lockbox',
      icon: Icons.lock_clock_rounded,
      gradient: const [Color(0xFF0891B2), Color(0xFF22D3EE)],
      child: Column(
        children: [
          if (_wishlist.isEmpty)
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 24),
              child: Column(
                children: [
                  Icon(
                    Icons.shopping_bag_outlined,
                    size: 40,
                    color: const Color(0xFF64748B).withValues(alpha: 0.5),
                  ),
                  const SizedBox(height: 8),
                  const Text(
                    'Your lockbox is empty',
                    style: TextStyle(
                      color: Color(0xFF64748B),
                      fontWeight: FontWeight.w500,
                      fontSize: 14,
                    ),
                  ),
                  const Text(
                    'Lock items here to delay impulse purchases.',
                    style: TextStyle(color: Color(0xFF64748B), fontSize: 12),
                  ),
                ],
              ),
            ),
          ..._wishlist.map((item) {
            bool unlocked = item.isUnlocked;
            bool bought = item.status == 'bought';

            BoxDecoration itemDecoration;
            if (bought) {
              itemDecoration = BoxDecoration(
                color: const Color(0xFFF1F5F9),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: const Color(0xFFE2E8F0)),
              );
            } else if (unlocked) {
              itemDecoration = BoxDecoration(
                color: const Color(0xFFECFDF5),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: const Color(0xFFA7F3D0), width: 1.5),
                boxShadow: [
                  BoxShadow(
                    color: const Color(0xFF10B981).withValues(alpha: 0.06),
                    blurRadius: 10,
                    offset: const Offset(0, 4),
                  ),
                ],
              );
            } else {
              itemDecoration = BoxDecoration(
                color: Colors.white,
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: const Color(0xFFE2E8F0)),
              );
            }

            return Dismissible(
              key: Key('wishlist_${item.id}'),
              direction: DismissDirection.endToStart,
              background: Container(
                margin: const EdgeInsets.only(bottom: 12),
                padding: const EdgeInsets.symmetric(horizontal: 20),
                decoration: BoxDecoration(
                  color: Colors.red.shade100,
                  borderRadius: BorderRadius.circular(20),
                ),
                alignment: Alignment.centerRight,
                child: Icon(
                  Icons.delete_outline_rounded,
                  color: Colors.red.shade700,
                  size: 28,
                ),
              ),
              confirmDismiss: (direction) async {
                final confirm = await showDialog<bool>(
                  context: context,
                  builder: (context) => AlertDialog(
                    title: const Text('Delete Item'),
                    content: const Text(
                      'Are you sure you want to delete this wishlist item?',
                    ),
                    actions: [
                      TextButton(
                        onPressed: () => Navigator.pop(context, false),
                        child: const Text('Cancel'),
                      ),
                      TextButton(
                        onPressed: () => Navigator.pop(context, true),
                        child: const Text(
                          'Delete',
                          style: TextStyle(color: Colors.red),
                        ),
                      ),
                    ],
                  ),
                );

                if (confirm == true) {
                  try {
                    await widget.disciplineService.deleteWishlistItem(item.id);
                    return true;
                  } catch (e) {
                    if (context.mounted)
                      UiUtils.showSnack(
                        context,
                        'Failed to delete: $e',
                        isError: true,
                      );
                    return false;
                  }
                }
                return false;
              },
              onDismissed: (direction) {
                setState(() {
                  _wishlist.removeWhere((w) => w.id == item.id);
                });
                UiUtils.showSnack(context, 'Wishlist item deleted');
              },
              child: AnimatedBuilder(
                animation: _peekAnimation,
                builder: (context, child) {
                  if (_peekAnimation.value == 0) return child!;
                  return Stack(
                    children: [
                      Positioned.fill(
                        child: Container(
                          margin: const EdgeInsets.only(bottom: 12),
                          padding: const EdgeInsets.symmetric(horizontal: 20),
                          decoration: BoxDecoration(
                            color: Colors.red.shade100,
                            borderRadius: BorderRadius.circular(20),
                          ),
                          alignment: Alignment.centerRight,
                          child: Icon(
                            Icons.delete_outline_rounded,
                            color: Colors.red.shade700,
                            size: 28,
                          ),
                        ),
                      ),
                      Transform.translate(
                        offset: Offset(_peekAnimation.value, 0),
                        child: child,
                      ),
                    ],
                  );
                },
                child: Container(
                  margin: const EdgeInsets.only(bottom: 12),
                  padding: const EdgeInsets.all(16),
                  decoration: itemDecoration,
                  child: Row(
                    children: [
                      Container(
                        padding: const EdgeInsets.all(10),
                        decoration: BoxDecoration(
                          color: bought
                              ? const Color(0xFFE2E8F0)
                              : unlocked
                              ? const Color(0xFFD1FAE5)
                              : const Color(0xFFFEF3C7),
                          shape: BoxShape.circle,
                        ),
                        child: Icon(
                          bought
                              ? Icons.check_circle_outline_rounded
                              : unlocked
                              ? Icons.lock_open_rounded
                              : Icons.lock_outline_rounded,
                          color: bought
                              ? const Color(0xFF64748B)
                              : unlocked
                              ? const Color(0xFF059669)
                              : const Color(0xFFD97706),
                          size: 20,
                        ),
                      ),
                      const SizedBox(width: 14),
                      Expanded(
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              item.name,
                              style: TextStyle(
                                fontWeight: FontWeight.w700,
                                fontSize: 15,
                                color: bought
                                    ? const Color(0xFF64748B)
                                    : unlocked
                                    ? const Color(0xFF064E3B)
                                    : const Color(0xFF1E293B),
                                decoration: bought
                                    ? TextDecoration.lineThrough
                                    : null,
                              ),
                            ),
                            const SizedBox(height: 2),
                            Text(
                              '₹${item.amount.toStringAsFixed(0)}',
                              style: TextStyle(
                                color: bought
                                    ? const Color(0xFF94A3B8)
                                    : unlocked
                                    ? const Color(0xFF047857)
                                    : const Color(0xFF0891B2),
                                fontWeight: FontWeight.w800,
                                fontSize: 14,
                              ),
                            ),
                          ],
                        ),
                      ),
                      if (bought)
                        Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 10,
                            vertical: 6,
                          ),
                          decoration: BoxDecoration(
                            color: const Color(0xFFE2E8F0),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: const Text(
                            'Purchased',
                            style: TextStyle(
                              color: Color(0xFF64748B),
                              fontWeight: FontWeight.bold,
                              fontSize: 11,
                            ),
                          ),
                        )
                      else if (unlocked)
                        FilledButton.icon(
                          onPressed: () =>
                              _updateWishlistStatus(item.id, 'bought'),
                          icon: const Icon(
                            Icons.shopping_cart_checkout_rounded,
                            size: 14,
                            color: Colors.white,
                          ),
                          label: const Text(
                            'Buy Now',
                            style: TextStyle(
                              fontWeight: FontWeight.bold,
                              fontSize: 12,
                              color: Colors.white,
                            ),
                          ),
                          style: FilledButton.styleFrom(
                            backgroundColor: const Color(0xFF059669),
                            foregroundColor: Colors.white,
                            padding: const EdgeInsets.symmetric(
                              horizontal: 14,
                              vertical: 10,
                            ),
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(12),
                            ),
                          ),
                        )
                      else
                        TweenAnimationBuilder<double>(
                          tween: Tween(
                            begin: 0,
                            end: item.unlockDate
                                .difference(DateTime.now())
                                .inDays
                                .toDouble(),
                          ),
                          duration: const Duration(milliseconds: 1500),
                          curve: Curves.easeOutCubic,
                          builder: (context, value, child) {
                            final days = value.toInt().clamp(0, 365);
                            return Container(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 10,
                                vertical: 6,
                              ),
                              decoration: BoxDecoration(
                                color: const Color(0xFFFFF7ED),
                                borderRadius: BorderRadius.circular(12),
                                border: Border.all(
                                  color: const Color(0xFFFED7AA),
                                ),
                              ),
                              child: Row(
                                mainAxisSize: MainAxisSize.min,
                                children: [
                                  const Icon(
                                    Icons.timer_outlined,
                                    color: Color(0xFFD97706),
                                    size: 14,
                                  ),
                                  const SizedBox(width: 4),
                                  Text(
                                    '${days}d left',
                                    style: const TextStyle(
                                      color: Color(0xFFD97706),
                                      fontSize: 11,
                                      fontWeight: FontWeight.w800,
                                    ),
                                  ),
                                ],
                              ),
                            );
                          },
                        ),
                    ],
                  ),
                ),
              ),
            );
          }),
          const SizedBox(height: 8),
          SizedBox(
            width: double.infinity,
            child: OutlinedButton.icon(
              onPressed: _addWishlistItemDialog,
              icon: const Icon(Icons.add_circle_outline_rounded, size: 18),
              label: const Text(
                'Add to Wishlist',
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              style: OutlinedButton.styleFrom(
                foregroundColor: const Color(0xFF0891B2),
                side: const BorderSide(color: Color(0xFF0891B2), width: 1.5),
                padding: const EdgeInsets.symmetric(vertical: 14),
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(16),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  void _addWishlistItemDialog() {
    final nameCtrl = TextEditingController();
    final amountCtrl = TextEditingController();
    int lockDays = 30;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => StatefulBuilder(
        builder: (context, setState) => Container(
          margin: EdgeInsets.only(
            bottom: MediaQuery.of(context).viewInsets.bottom + 24,
            left: 16,
            right: 16,
          ),
          padding: const EdgeInsets.all(24),
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(28),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withValues(alpha: 0.1),
                blurRadius: 20,
              ),
            ],
          ),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              // Handle
              Container(
                width: 40,
                height: 4,
                margin: const EdgeInsets.only(bottom: 24),
                decoration: BoxDecoration(
                  color: Colors.grey.shade300,
                  borderRadius: BorderRadius.circular(4),
                ),
              ),
              const Text(
                'Add Wishlist Item',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 24),
              TextField(
                controller: nameCtrl,
                textCapitalization: TextCapitalization.words,
                decoration: InputDecoration(
                  labelText: 'Item Name',
                  filled: true,
                  fillColor: Colors.grey.shade50,
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(14),
                    borderSide: BorderSide.none,
                  ),
                ),
              ),
              const SizedBox(height: 16),
              TextField(
                controller: amountCtrl,
                keyboardType: TextInputType.number,
                decoration: InputDecoration(
                  labelText: 'Amount (₹)',
                  filled: true,
                  fillColor: Colors.grey.shade50,
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(14),
                    borderSide: BorderSide.none,
                  ),
                ),
              ),
              const SizedBox(height: 20),
              const Align(
                alignment: Alignment.centerLeft,
                child: Text(
                  'Lock Duration',
                  style: TextStyle(
                    fontWeight: FontWeight.w600,
                    fontSize: 13,
                    color: Color(0xFF64748B),
                  ),
                ),
              ),
              const SizedBox(height: 12),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [7, 15, 30].map((days) {
                  final sel = lockDays == days;
                  return ChoiceChip(
                    label: Text('${days}d'),
                    selected: sel,
                    onSelected: (v) {
                      if (v) setState(() => lockDays = days);
                    },
                    selectedColor: const Color(0xFF0D9488),
                    labelStyle: TextStyle(
                      color: sel ? Colors.white : Colors.grey.shade700,
                      fontWeight: sel ? FontWeight.bold : FontWeight.normal,
                    ),
                    backgroundColor: Colors.grey.shade100,
                    side: BorderSide.none,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  );
                }).toList(),
              ),
              const SizedBox(height: 16),
              // Animated helper text
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: Colors.blue.shade50,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Row(
                  children: [
                    const Icon(
                      Icons.info_outline_rounded,
                      color: Colors.blue,
                      size: 20,
                    ),
                    const SizedBox(width: 12),
                    Expanded(
                      child: AnimatedSwitcher(
                        duration: const Duration(milliseconds: 300),
                        transitionBuilder:
                            (Widget child, Animation<double> animation) {
                              return FadeTransition(
                                opacity: animation,
                                child: SlideTransition(
                                  position: Tween<Offset>(
                                    begin: const Offset(0, 0.2),
                                    end: Offset.zero,
                                  ).animate(animation),
                                  child: child,
                                ),
                              );
                            },
                        child: Text(
                          'This item will be securely locked for $lockDays days to prevent impulse buying.',
                          key: ValueKey<int>(lockDays),
                          style: const TextStyle(
                            fontSize: 12,
                            color: Colors.blue,
                            height: 1.4,
                          ),
                        ),
                      ),
                    ),
                  ],
                ),
              ),
              const SizedBox(height: 24),
              Row(
                children: [
                  Expanded(
                    child: TextButton(
                      onPressed: () => Navigator.pop(context),
                      style: TextButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                      ),
                      child: const Text(
                        'Cancel',
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          color: Color(0xFF64748B),
                        ),
                      ),
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: ElevatedButton(
                      onPressed: () async {
                        try {
                          final item = await widget.disciplineService
                              .addWishlistItem(
                                nameCtrl.text,
                                double.parse(amountCtrl.text),
                                lockDays,
                              );
                          if (context.mounted) {
                            setState(() => _wishlist.add(item));
                            Navigator.pop(context);
                            UiUtils.showSnack(context, 'Wishlist item added!');
                          }
                        } catch (e) {
                          if (context.mounted)
                            UiUtils.showSnack(
                              context,
                              'Failed: $e',
                              isError: true,
                            );
                        }
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: const Color(0xFF0D9488),
                        padding: const EdgeInsets.symmetric(vertical: 14),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(14),
                        ),
                      ),
                      child: const Text(
                        'Lock Item',
                        style: TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _updateWishlistStatus(int id, String status) async {
    try {
      final updated = await widget.disciplineService.updateWishlistItemStatus(
        id,
        status,
      );
      if (mounted) {
        setState(() {
          final idx = _wishlist.indexWhere((w) => w.id == id);
          if (idx != -1) _wishlist[idx] = updated;
        });
        UiUtils.showSnack(context, 'Wishlist item purchased!');
      }
    } catch (e) {
      if (mounted) UiUtils.showSnack(context, 'Failed: $e', isError: true);
    }
  }
}
