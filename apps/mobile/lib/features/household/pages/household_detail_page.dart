import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/features/goals/models/goal_models.dart';
import 'package:horizon_mobile/features/budget/models/budget_models.dart';
import 'package:horizon_mobile/features/household/state/household_state.dart';

class HouseholdDetailPage extends ConsumerWidget {
  final String householdId;

  const HouseholdDetailPage({super.key, required this.householdId});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final detailAsync = ref.watch(householdDetailProvider(householdId));
    final goalsAsync = ref.watch(householdGoalsProvider(householdId));
    final budgetsAsync = ref.watch(householdBudgetsProvider(householdId));

    return Scaffold(
      appBar: AppBar(
        title: const Text('Household Details'),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {
              context.push('/households/$householdId/settings');
            },
          )
        ],
      ),
      body: detailAsync.when(
        data: (detail) {
          return RefreshIndicator(
            onRefresh: () async {
              ref.invalidate(householdDetailProvider(householdId));
              ref.invalidate(householdGoalsProvider(householdId));
              ref.invalidate(householdBudgetsProvider(householdId));
            },
            child: ListView(
              padding: const EdgeInsets.all(16),
              children: [
                _buildHeader(detail),
                const SizedBox(height: 24),
                _buildFinancialSummary(detail),
                const SizedBox(height: 24),
                _buildMembersSection(context, detail),
                const SizedBox(height: 24),
                _buildSharedGoalsSection(context, ref, goalsAsync, detail),
                const SizedBox(height: 24),
                _buildSharedBudgetsSection(context, ref, budgetsAsync, detail),
                const SizedBox(height: 24),
                _buildLinkedAccountsSection(context, detail),
              ],
            ),
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(child: Text('Error: $error')),
      ),
    );
  }

  Widget _buildHeader(detail) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          detail.name,
          style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
        ),
        Text(
          '${detail.householdType} • ${detail.status}',
          style: const TextStyle(color: Colors.grey),
        ),
      ],
    );
  }

  Widget _buildFinancialSummary(detail) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text('Financial Summary', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            const SizedBox(height: 16),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Net Worth'),
                Text('${detail.currency} ${(detail.totalNetWorth / 100).toStringAsFixed(2)}', style: const TextStyle(fontWeight: FontWeight.bold)),
              ],
            ),
            const Divider(),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Assets'),
                Text('${detail.currency} ${(detail.totalAssets / 100).toStringAsFixed(2)}', style: const TextStyle(color: Colors.green)),
              ],
            ),
            const SizedBox(height: 8),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                const Text('Liabilities'),
                Text('${detail.currency} ${(detail.totalLiabilities / 100).toStringAsFixed(2)}', style: const TextStyle(color: Colors.red)),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMembersSection(BuildContext context, detail) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text('Members', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            TextButton.icon(
              onPressed: () {
                context.push('/households/$householdId/invite');
              },
              icon: const Icon(Icons.person_add),
              label: const Text('Invite'),
            ),
          ],
        ),
        Card(
          child: ListView.separated(
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            itemCount: detail.members.length,
            separatorBuilder: (context, index) => const Divider(height: 1),
            itemBuilder: (context, index) {
              final member = detail.members[index];
              return ListTile(
                leading: CircleAvatar(child: Text(member.userId.substring(0, 1).toUpperCase())),
                title: Text(member.userId),
                subtitle: Text('Status: ${member.inviteStatus}'),
                trailing: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: _roleColor(member.role).withValues(alpha: 0.15),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(member.role, style: TextStyle(fontSize: 12, color: _roleColor(member.role), fontWeight: FontWeight.w500)),
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildSharedGoalsSection(BuildContext context, WidgetRef ref, AsyncValue<List<GoalSummary>> goalsAsync, detail) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text('Shared Goals', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            Text('${detail.linkedGoals.length} linked'),
          ],
        ),
        const SizedBox(height: 8),
        goalsAsync.when(
          data: (goals) {
            if (goals.isEmpty) {
              return Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Text('No shared goals yet.', style: TextStyle(color: Colors.grey.shade600)),
                ),
              );
            }
            return Card(
              child: Column(
                children: goals.map((goal) => ListTile(
                  leading: CircleAvatar(
                    backgroundColor: _goalColor(goal.status),
                    child: Icon(Icons.flag, color: Colors.white, size: 18),
                  ),
                  title: Text(goal.name),
                  subtitle: Text('${goal.importance} • ${goal.progress}%'),
                  trailing: Text(goal.status, style: TextStyle(color: _goalColor(goal.status), fontSize: 12)),
                  onTap: () => context.push('/goals/${goal.goalId}'),
                )).toList(),
              ),
            );
          },
          loading: () => const Card(child: Padding(padding: EdgeInsets.all(16), child: Text('Loading...'))),
          error: (e, _) => Card(child: Padding(padding: EdgeInsets.all(16), child: Text('Error: $e'))),
        ),
      ],
    );
  }

  Widget _buildSharedBudgetsSection(BuildContext context, WidgetRef ref, AsyncValue<List<BudgetModel>> budgetsAsync, detail) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            const Text('Shared Budgets', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
            Text('${detail.linkedBudgets.length} linked'),
          ],
        ),
        const SizedBox(height: 8),
        budgetsAsync.when(
          data: (budgets) {
            if (budgets.isEmpty) {
              return Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Text('No shared budgets yet.', style: TextStyle(color: Colors.grey.shade600)),
                ),
              );
            }
            return Card(
              child: Column(
                children: budgets.map((budget) => ListTile(
                  leading: CircleAvatar(
                    backgroundColor: _budgetColor(budget.status),
                    child: Icon(Icons.pie_chart, color: Colors.white, size: 18),
                  ),
                  title: Text(budget.name),
                  subtitle: Text('${budget.period} • ₹${budget.totalBudgeted.toStringAsFixed(0)}'),
                  trailing: Text(budget.status, style: TextStyle(color: _budgetColor(budget.status), fontSize: 12)),
                  onTap: () => context.push('/planning/budget/${budget.budgetId}'),
                )).toList(),
              ),
            );
          },
          loading: () => const Card(child: Padding(padding: EdgeInsets.all(16), child: Text('Loading...'))),
          error: (e, _) => Card(child: Padding(padding: EdgeInsets.all(16), child: Text('Error: $e'))),
        ),
      ],
    );
  }

  Widget _buildLinkedAccountsSection(BuildContext context, detail) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const Text('Linked Accounts', style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold)),
        const SizedBox(height: 8),
        Card(
          child: detail.linkedAccounts.isEmpty
              ? const Padding(padding: EdgeInsets.all(16), child: Text('No linked accounts.'))
              : Column(
                  children: detail.linkedAccounts.map((acct) => ListTile(
                    leading: const Icon(Icons.account_balance),
                    title: Text(acct.accountId),
                    subtitle: Text('Added by: ${acct.addedBy}'),
                    onTap: () => context.push('/accounts/${acct.accountId}'),
                  )).toList(),
                ),
        ),
      ],
    );
  }

  Color _goalColor(String status) {
    switch (status) {
      case 'Active': return Colors.green;
      case 'Completed': return Colors.blue;
      case 'AtRisk': return Colors.red;
      case 'Paused': return Colors.orange;
      default: return Colors.grey;
    }
  }

  Color _budgetColor(String status) {
    switch (status) {
      case 'Active': return Colors.green;
      case 'Draft': return Colors.orange;
      case 'Completed': return Colors.blue;
      case 'Paused': return Colors.grey;
      case 'Archived': return Colors.grey;
      default: return Colors.grey;
    }
  }

  Color _roleColor(String role) {
    switch (role) {
      case 'Head': return Colors.purple;
      case 'Admin': return Colors.blue;
      case 'Member': return Colors.green;
      case 'Viewer': return Colors.grey;
      default: return Colors.grey;
    }
  }
}
