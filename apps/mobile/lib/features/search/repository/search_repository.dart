import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:horizon_mobile/core/network/api_client.dart';
import 'package:horizon_mobile/features/accounts/repository/acct_repository.dart';
import 'package:horizon_mobile/features/goals/repository/goal_repository.dart';
import 'package:horizon_mobile/features/portfolio/repository/pf_repository.dart';
import 'package:horizon_mobile/features/timeline/repository/timeline_repository.dart';
import '../models/search_models.dart';

final searchRepositoryProvider = Provider<SearchRepository>((ref) {
  return SearchRepository(
    apiClient: ref.read(apiClientProvider),
    acctRepo: ref.read(acctRepositoryProvider),
    goalRepo: ref.read(goalRepositoryProvider),
    pfRepo: ref.read(pfRepositoryProvider),
    tlRepo: ref.read(timelineRepositoryProvider),
  );
});

class SearchRepository {
  final ApiClient apiClient;
  final AcctRepository acctRepo;
  final GoalRepository goalRepo;
  final PfRepository pfRepo;
  final TimelineRepository tlRepo;

  SearchRepository({
    required this.apiClient,
    required this.acctRepo,
    required this.goalRepo,
    required this.pfRepo,
    required this.tlRepo,
  });

  Future<List<SearchResult>> search({
    required String query,
    Set<SearchModule>? modules,
  }) async {
    if (query.trim().isEmpty) return [];
    final results = <SearchResult>[];
    final q = query.toLowerCase();

    final searchAll = modules == null || modules.isEmpty;
    final activeMods = modules ?? <SearchModule>{};

    if (searchAll || activeMods.contains(SearchModule.accounts)) { await _searchAccounts(q, results); }
    if (searchAll || activeMods.contains(SearchModule.transactions)) { await _searchTransactions(q, results); }
    if (searchAll || activeMods.contains(SearchModule.goals)) { await _searchGoals(q, results); }
    if (searchAll || activeMods.contains(SearchModule.portfolio)) { await _searchPortfolio(q, results); }
    if (searchAll || activeMods.contains(SearchModule.timeline)) { await _searchTimeline(q, results); }

    return results;
  }

  Future<void> _searchAccounts(String q, List<SearchResult> results) async {
    try {
      final resp = await acctRepo.getAccounts();
      for (final a in resp.data?.accounts ?? []) {
        if (a.accountName.toLowerCase().contains(q) || a.accountType.toLowerCase().contains(q) || (a.institutionName?.toLowerCase().contains(q) == true)) {
          results.add(SearchResult(id: a.accountId, module: SearchModule.accounts, title: a.accountName, subtitle: '${a.accountType} · ${a.currency}', detail: a.status));
        }
      }
    } catch (_) {}
  }

  Future<void> _searchTransactions(String q, List<SearchResult> results) async {
    try {
      final resp = await acctRepo.getAccounts();
      for (final acct in resp.data?.accounts ?? []) {
        try {
          final detail = await acctRepo.getDetail(accountId: acct.accountId);
          for (final t in detail.data?.transactions ?? []) {
            if (t.description.toLowerCase().contains(q) || t.category.toLowerCase().contains(q)) {
              results.add(SearchResult(id: t.eventId, module: SearchModule.transactions, title: t.description, subtitle: '${t.category} · ${acct.accountName}', detail: t.eventDate));
            }
          }
        } catch (_) {}
      }
    } catch (_) {}
  }

  Future<void> _searchGoals(String q, List<SearchResult> results) async {
    try {
      final resp = await goalRepo.getGoals();
      for (final g in resp.data?.goals ?? []) {
        if (g.name.toLowerCase().contains(q) || g.status.toLowerCase().contains(q) || g.importance.toLowerCase().contains(q)) {
          results.add(SearchResult(id: g.goalId, module: SearchModule.goals, title: g.name, subtitle: '${g.status} · ${g.progress}% complete'));
        }
      }
    } catch (_) {}
  }

  Future<void> _searchPortfolio(String q, List<SearchResult> results) async {
    try {
      final resp = await pfRepo.getDashboard();
      final dash = resp.data;
      if (dash != null) {
        for (final c in dash.cards) {
          if (c.title.toLowerCase().contains(q) || c.summary.toLowerCase().contains(q)) {
            results.add(SearchResult(id: c.cardId, module: SearchModule.portfolio, title: c.title, subtitle: c.summary, route: '/portfolio'));
          }
        }
      }
    } catch (_) {}
  }

  Future<void> _searchTimeline(String q, List<SearchResult> results) async {
    try {
      final items = await tlRepo.search(query: q);
      results.addAll(items.map((t) => SearchResult(id: t.timelineId, module: SearchModule.timeline, title: t.title, subtitle: t.summary.isNotEmpty ? t.summary : t.category, detail: t.formattedDate)));
    } catch (_) {}
  }
}
