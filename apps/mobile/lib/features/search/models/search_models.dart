enum SearchModule { accounts, transactions, goals, portfolio, timeline }

extension SearchModuleExt on SearchModule {
  String get label {
    switch (this) {
      case SearchModule.accounts: return 'Accounts';
      case SearchModule.transactions: return 'Transactions';
      case SearchModule.goals: return 'Goals';
      case SearchModule.portfolio: return 'Portfolio';
      case SearchModule.timeline: return 'Timeline';
    }
  }

  String get iconName {
    switch (this) {
      case SearchModule.accounts: return 'account_balance';
      case SearchModule.transactions: return 'currency_rupee';
      case SearchModule.goals: return 'flag';
      case SearchModule.portfolio: return 'pie_chart';
      case SearchModule.timeline: return 'timeline';
    }
  }

  String get routePrefix {
    switch (this) {
      case SearchModule.accounts: return '/accounts';
      case SearchModule.transactions: return '';
      case SearchModule.goals: return '/goals';
      case SearchModule.portfolio: return '/portfolio';
      case SearchModule.timeline: return '/timeline';
    }
  }
}

class SearchResult {
  final String id;
  final SearchModule module;
  final String title;
  final String subtitle;
  final String? route;
  final String? detail;

  const SearchResult({
    required this.id,
    required this.module,
    required this.title,
    required this.subtitle,
    this.route,
    this.detail,
  });

  String get routePath {
    if (route != null) return route!;
    if (module == SearchModule.transactions) return '';
    return '${module.routePrefix}/$id';
  }
}

class RecentSearch {
  final String query;
  final DateTime timestamp;

  const RecentSearch({required this.query, required this.timestamp});
}

class SearchFilters {
  final Set<SearchModule> modules;
  final String? dateRange;
  final String? category;

  const SearchFilters({
    this.modules = const {},
    this.dateRange,
    this.category,
  });

  bool get hasActiveFilters => modules.isNotEmpty || dateRange != null || category != null;

  SearchFilters copyWith({Set<SearchModule>? modules, String? dateRange, String? category}) {
    return SearchFilters(
      modules: modules ?? this.modules,
      dateRange: dateRange ?? this.dateRange,
      category: category ?? this.category,
    );
  }
}
