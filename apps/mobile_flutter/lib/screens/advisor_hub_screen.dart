import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/financial_provider.dart';
import 'nudge_center_screen.dart';
import 'simulation_screen.dart';
import 'debt_optimizer_screen.dart';
import 'subscription_insights_screen.dart';
import 'lifecycle_screen.dart';

class AdvisorHubScreen extends StatefulWidget {
  const AdvisorHubScreen({super.key});

  @override
  State<AdvisorHubScreen> createState() => _AdvisorHubScreenState();
}

class _AdvisorHubScreenState extends State<AdvisorHubScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final provider = context.read<FinancialProvider>();
      provider.loadAllAdvisorData();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(title: const Text('Financial Advisor')),
      body: Consumer<FinancialProvider>(
        builder: (context, fp, _) {
          final activeNudges = fp.nudges.where((n) => n['is_dismissed'] == false).length;
          final criticalNudges = fp.nudges.where((n) => n['severity'] == 'critical' && n['is_dismissed'] == false).length;
          final loading = fp.isLoading;

          return RefreshIndicator(
            onRefresh: () => fp.loadAllAdvisorData(),
            child: ListView(
              padding: const EdgeInsets.all(20),
              children: [
                if (loading)
                  Padding(
                    padding: const EdgeInsets.only(bottom: 16),
                    child: LinearProgressIndicator(color: const Color(0xFF6B46C1), backgroundColor: const Color(0xFF6B46C1).withOpacity(0.1)),
                  ),
                if (criticalNudges > 0)
                  Container(
                    margin: const EdgeInsets.only(bottom: 20),
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: Colors.red.shade50,
                      borderRadius: BorderRadius.circular(20),
                      border: Border.all(color: Colors.red.shade200),
                    ),
                    child: Row(
                      children: [
                        Icon(Icons.warning_amber_rounded, color: Colors.red.shade700, size: 28),
                        const SizedBox(width: 12),
                        Expanded(
                          child: Text('$criticalNudges critical alerts need your attention',
                              style: TextStyle(color: Colors.red.shade800, fontWeight: FontWeight.bold)),
                        ),
                      ],
                    ),
                  ),
                _buildFeatureCard(
                  context,
                  icon: Icons.notifications_active_rounded,
                  title: 'Nudge Center',
                  subtitle: '$activeNudges active coaching nudges',
                  gradient: const [Color(0xFF0891B2), Color(0xFF22D3EE)],
                  onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const NudgeCenterScreen())),
                ),
                const SizedBox(height: 16),
                _buildFeatureCard(
                  context,
                  icon: Icons.query_stats_rounded,
                  title: 'Scenario Simulator',
                  subtitle: 'Model financial decisions before making them',
                  gradient: const [Color(0xFF6B46C1), Color(0xFF8B5CF6)],
                  onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const SimulationScreen())),
                ),
                const SizedBox(height: 16),
                _buildFeatureCard(
                  context,
                  icon: Icons.trending_down_rounded,
                  title: 'Debt Optimizer',
                  subtitle: fp.debtOptimizer != null
                      ? '${(fp.debtOptimizer!['estimated_freedom_months'] ?? 0)} months to debt freedom'
                      : 'Analyze your debt payoff strategy',
                  gradient: const [Color(0xFFDC2626), Color(0xFFF87171)],
                  onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const DebtOptimizerScreen())),
                ),
                const SizedBox(height: 16),
                _buildFeatureCard(
                  context,
                  icon: Icons.subscriptions_rounded,
                  title: 'Subscription Audit',
                  subtitle: '${fp.subscriptionInsights.length} subscriptions · Find savings opportunities',
                  gradient: const [Color(0xFF059669), Color(0xFF10B981)],
                  onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const SubscriptionInsightsScreen())),
                ),
                const SizedBox(height: 16),
                _buildFeatureCard(
                  context,
                  icon: Icons.public_rounded,
                  title: 'Lifecycle Simulator',
                  subtitle: 'Project your wealth to age 100 with inflation',
                  gradient: const [Color(0xFF6B46C1), Color(0xFF8B5CF6)],
                  onTap: () => Navigator.push(context, MaterialPageRoute(builder: (_) => const LifecycleScreen())),
                ),
                const SizedBox(height: 60),
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildFeatureCard(BuildContext context, {required IconData icon, required String title, required String subtitle, required List<Color> gradient, required VoidCallback onTap}) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(24),
          border: Border.all(color: gradient[0].withOpacity(0.2), width: 1.5),
          boxShadow: [BoxShadow(color: gradient[0].withOpacity(0.08), blurRadius: 24, offset: const Offset(0, 12))],
        ),
        child: ClipRRect(
          borderRadius: BorderRadius.circular(24),
          child: Container(
            decoration: BoxDecoration(border: Border(left: BorderSide(color: gradient[0], width: 6))),
            padding: const EdgeInsets.all(24),
            child: Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(14),
                  decoration: BoxDecoration(
                    gradient: LinearGradient(colors: gradient),
                    shape: BoxShape.circle,
                    boxShadow: [BoxShadow(color: gradient[0].withOpacity(0.4), blurRadius: 12, offset: const Offset(0, 4))],
                  ),
                  child: Icon(icon, color: Colors.white, size: 24),
                ),
                const SizedBox(width: 18),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(title, style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w900, color: Color(0xFF1E293B))),
                      const SizedBox(height: 4),
                      Text(subtitle, style: const TextStyle(color: Colors.grey, fontWeight: FontWeight.w600, fontSize: 13)),
                    ],
                  ),
                ),
                Icon(Icons.chevron_right_rounded, color: gradient[0]),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
