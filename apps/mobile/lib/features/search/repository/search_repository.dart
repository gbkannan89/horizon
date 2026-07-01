import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/features/accounts/repository/acct_repository.dart';
import 'package:horizon_mobile/features/goals/repository/goal_repository.dart';
import '../../transactions/repository/transaction_repository.dart';
import '../../timeline/repository/timeline_repository.dart';
import '../models/search_models.dart';

final searchRepositoryProvider = Provider<SearchRepository>((ref) {
  return SearchRepository(
    apiClient: ref.read(apiClientProvider),
    transactionRepo: ref.read(transactionRepositoryProvider),
    timelineRepo: ref.read(timelineRepositoryProvider),
    accountsRepo: ref.read(acctRepositoryProvider),
    goalsRepo: ref.read(goalRepositoryProvider),
  );
});

class SearchRepository {
  final ApiClient apiClient;
  final TransactionRepository transactionRepo;
  final TimelineRepository timelineRepo;
  final AcctRepository accountsRepo;
  final GoalRepository goalsRepo;

  SearchRepository({
    required this.apiClient,
    required this.transactionRepo,
    required this.timelineRepo,
    required this.accountsRepo,
    required this.goalsRepo,
  });

  Future<List<SearchResultItem>> searchAll(String query) async {
    final results = <SearchResultItem>[];
    final futures = <Future<void>>[
      _searchTransactions(query, results),
      _searchTimeline(query, results),
      _searchAccounts(query, results),
      _searchGoals(query, results),
    ];
    await Future.wait(futures);
    results.sort((a, b) => a.module.compareTo(b.module));
    return results;
  }

  Future<void> _searchTransactions(String q, List<SearchResultItem> results) async {
    try {
      final resp = await transactionRepo.searchTransactions(query: q);
      final data = resp.data;
      if (data != null) {
        for (final tx in data.transactions) {
          results.add(SearchResultItem(
            id: tx.eventId,
            title: tx.description ?? tx.eventType,
            subtitle: '${tx.category ?? tx.eventType} · ${tx.formattedDate}',
            module: 'transaction',
            amount: tx.formattedAmount,
            icon: Icons.receipt_long_outlined,
            color: tx.isExpense ? Colors.red : Colors.green,
            route: '/transactions/${tx.eventId}',
          ));
        }
      }
    } catch (_) {}
  }

  Future<void> _searchTimeline(String q, List<SearchResultItem> results) async {
    try {
      final items = await timelineRepo.search(query: q);
      for (final item in items) {
        results.add(SearchResultItem(
          id: item.timelineId,
          title: item.title,
          subtitle: '${item.category} · ${item.formattedDate}',
          module: 'timeline',
          amount: item.amount > 0 ? '₹${item.amount}' : null,
          icon: Icons.timeline,
          color: Colors.indigo,
          route: '/timeline/${item.timelineId}',
        ));
      }
    } catch (_) {}
  }

  Future<void> _searchAccounts(String q, List<SearchResultItem> results) async {
    try {
      final resp = await accountsRepo.getAccounts();
      final data = resp.data;
      if (data != null) {
        final ql = q.toLowerCase();
        for (final acct in data.accounts) {
          if (acct.accountName.toLowerCase().contains(ql)) {
            results.add(SearchResultItem(
              id: acct.accountId,
              title: acct.accountName,
              subtitle: '${acct.accountType} · ${acct.institutionName ?? 'Account'}',
              module: 'account',
              icon: Icons.account_balance,
              color: Colors.teal,
              route: '/accounts/${acct.accountId}',
            ));
          }
        }
      }
    } catch (_) {}
  }

  Future<void> _searchGoals(String q, List<SearchResultItem> results) async {
    try {
      final resp = await goalsRepo.getGoals();
      final data = resp.data;
      if (data != null) {
        final ql = q.toLowerCase();
        for (final goal in data.goals) {
          if (goal.name.toLowerCase().contains(ql)) {
            results.add(SearchResultItem(
              id: goal.goalId,
              title: goal.name,
              subtitle: 'Goal · ${goal.status}',
              module: 'goal',
              icon: Icons.flag,
              color: Colors.blue,
              route: '/goals/${goal.goalId}',
            ));
          }
        }
      }
    } catch (_) {}
  }
}
