import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'dashboard_screen.dart';
import 'budget_screen.dart';
import 'profile_screen.dart';
import 'portfolio_screen.dart';
import 'discipline_dashboard_screen.dart';
import 'advisor_hub_screen.dart';
import '../services/discipline_service.dart';

class MainShell extends StatefulWidget {
  const MainShell({super.key});

  @override
  State<MainShell> createState() => _MainShellState();
}

class _MainShellState extends State<MainShell> with SingleTickerProviderStateMixin {
  int _currentIndex = 0;

  final List<Widget> _pages = [
    const DashboardScreen(),
    const PortfolioScreen(),
    const BudgetScreen(),
    const AdvisorHubScreen(),
    DisciplineDashboardScreen(disciplineService: DisciplineService()),
    const ProfileScreen(),
  ];

  final List<_NavItem> _navItems = const [
    _NavItem(icon: Icons.home_outlined,               activeIcon: Icons.home_rounded,             label: 'Home'),
    _NavItem(icon: Icons.account_balance_outlined,    activeIcon: Icons.account_balance_rounded,  label: 'Portfolio'),
    _NavItem(icon: Icons.pie_chart_outline_rounded,   activeIcon: Icons.pie_chart_rounded,        label: 'Budget'),
    _NavItem(icon: Icons.auto_awesome_outlined,       activeIcon: Icons.auto_awesome_rounded,     label: 'Advisor'),
    _NavItem(icon: Icons.shield_outlined,             activeIcon: Icons.shield_rounded,           label: 'Discipline'),
    _NavItem(icon: Icons.person_outline_rounded,      activeIcon: Icons.person_rounded,           label: 'Profile'),
  ];

  @override
  Widget build(BuildContext context) {
    SystemChrome.setSystemUIOverlayStyle(const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.light,
    ));

    return Scaffold(
      body: IndexedStack(index: _currentIndex, children: _pages),
      bottomNavigationBar: _buildBottomNav(),
    );
  }

  Widget _buildBottomNav() {
    return Container(
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: const Color(0xFF1E3A8A).withOpacity(0.08),
            blurRadius: 24,
            offset: const Offset(0, -4),
          ),
        ],
      ),
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 8),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            children: List.generate(_navItems.length, (i) => _buildNavTile(i)),
          ),
        ),
      ),
    );
  }

  Widget _buildNavTile(int index) {
    final item = _navItems[index];
    final selected = _currentIndex == index;

    return GestureDetector(
      onTap: () => setState(() => _currentIndex = index),
      behavior: HitTestBehavior.opaque,
      child: AnimatedContainer(
        duration: const Duration(milliseconds: 220),
        curve: Curves.easeOutCubic,
        padding: EdgeInsets.symmetric(horizontal: selected ? 18 : 12, vertical: 8),
        decoration: BoxDecoration(
          color: selected ? const Color(0xFF1E3A8A) : Colors.transparent,
          borderRadius: BorderRadius.circular(20),
        ),
        child: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              selected ? item.activeIcon : item.icon,
              color: selected ? Colors.white : const Color(0xFF94A3B8),
              size: 22,
            ),
            if (selected) ...[
              const SizedBox(width: 6),
              Text(
                item.label,
                style: const TextStyle(
                  color: Colors.white,
                  fontWeight: FontWeight.w700,
                  fontSize: 13,
                  letterSpacing: 0.2,
                ),
              ),
            ],
          ],
        ),
      ),
    );
  }
}

class _NavItem {
  final IconData icon;
  final IconData activeIcon;
  final String label;
  const _NavItem({required this.icon, required this.activeIcon, required this.label});
}
