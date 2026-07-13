import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:percent_indicator/percent_indicator.dart';
import '../providers/financial_provider.dart';
import 'create_collection_screen.dart';
import 'collection_detail_screen.dart';

class CollectionsListScreen extends StatefulWidget {
  const CollectionsListScreen({super.key});

  @override
  State<CollectionsListScreen> createState() => _CollectionsListScreenState();
}

class _CollectionsListScreenState extends State<CollectionsListScreen> {
  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<FinancialProvider>(context, listen: false).loadCollections();
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: const Color(0xFFF8FAFC),
      appBar: AppBar(
        title: const Text('Collections'),
        actions: [
          Consumer<FinancialProvider>(
            builder: (_, fp, _) => fp.collections.isNotEmpty
                ? PopupMenuButton<String>(
                    icon: const Icon(Icons.filter_list),
                    onSelected: (v) {
                      if (v == 'all') {
                        fp.setCollectionFilter(null);
                      } else {
                        fp.setCollectionFilter(v);
                      }
                    },
                    itemBuilder: (_) => [
                      CheckedPopupMenuItem(
                        value: 'all',
                        checked: fp.collectionFilter == null,
                        child: const Text('All'),
                      ),
                      CheckedPopupMenuItem(
                        value: 'active',
                        checked: fp.collectionFilter == 'active',
                        child: const Text('Active'),
                      ),
                      CheckedPopupMenuItem(
                        value: 'settled',
                        checked: fp.collectionFilter == 'settled',
                        child: const Text('Settled'),
                      ),
                    ],
                  )
                : const SizedBox(),
          ),
        ],
      ),
      body: Consumer<FinancialProvider>(
        builder: (context, fp, _) {
          if (fp.isLoadingCollections) {
            return const Center(child: CircularProgressIndicator());
          }

          final list = fp.filteredCollections;

          if (list.isEmpty) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.people_outline,
                    size: 80,
                    color: Colors.grey.shade300,
                  ),
                  const SizedBox(height: 20),
                  const Text(
                    'No collections yet',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                      color: Color(0xFF1E293B),
                    ),
                  ),
                  const SizedBox(height: 8),
                  const Text(
                    'Create a collection to track group payments.',
                    style: TextStyle(color: Colors.grey),
                  ),
                  const SizedBox(height: 24),
                  ElevatedButton.icon(
                    onPressed: () => _createCollection(context),
                    icon: const Icon(Icons.add),
                    label: const Text('Create Collection'),
                  ),
                ],
              ),
            );
          }

          return RefreshIndicator(
            onRefresh: fp.loadCollections,
            child: ListView.builder(
              padding: const EdgeInsets.fromLTRB(16, 16, 16, 100),
              itemCount: list.length,
              itemBuilder: (ctx, i) => _collectionCard(context, fp, list[i]),
            ),
          );
        },
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => _createCollection(context),
        backgroundColor: const Color(0xFF6366F1),
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  void _createCollection(BuildContext context) {
    Navigator.of(
      context,
    ).push(MaterialPageRoute(builder: (_) => const CreateCollectionScreen()));
  }

  Widget _collectionCard(
    BuildContext context,
    FinancialProvider fp,
    Map<String, dynamic> c,
  ) {
    final totalExpected = (c['total_expected'] ?? 0).toDouble();
    final totalCollected = (c['total_collected'] ?? 0).toDouble();
    final pct = totalExpected > 0
        ? (totalCollected / totalExpected).clamp(0.0, 1.0)
        : 0.0;
    final memberCount = c['member_count'] ?? 0;
    final paidCount = c['paid_count'] ?? 0;
    final status = c['status'] ?? 'active';
    final isActive = status == 'active';
    final label = c['label'] ?? '';

    return GestureDetector(
      onTap: () {
        Navigator.of(context).push(
          MaterialPageRoute(
            builder: (_) => CollectionDetailScreen(collection: c),
          ),
        );
      },
      child: Container(
        margin: const EdgeInsets.only(bottom: 14),
        padding: const EdgeInsets.all(18),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(20),
          border: Border.all(
            color: isActive
                ? const Color(0xFF6366F1).withValues(alpha: 0.1)
                : const Color(0xFF94A3B8).withValues(alpha: 0.2),
          ),
          boxShadow: [
            BoxShadow(
              color: const Color(0xFF6366F1).withValues(alpha: 0.05),
              blurRadius: 16,
              offset: const Offset(0, 6),
            ),
          ],
        ),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.all(10),
                  decoration: BoxDecoration(
                    color: isActive
                        ? const Color(0xFF6366F1).withValues(alpha: 0.1)
                        : Colors.grey.shade100,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Icon(
                    Icons.people_alt_outlined,
                    color: isActive ? const Color(0xFF6366F1) : Colors.grey,
                    size: 22,
                  ),
                ),
                const SizedBox(width: 14),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        label,
                        style: const TextStyle(
                          fontSize: 16,
                          fontWeight: FontWeight.w700,
                          color: Color(0xFF1E293B),
                        ),
                        overflow: TextOverflow.ellipsis,
                      ),
                      const SizedBox(height: 2),
                      Text(
                        '$paidCount / $memberCount paid',
                        style: TextStyle(
                          fontSize: 12,
                          color: isActive
                              ? const Color(0xFF64748B)
                              : Colors.grey,
                        ),
                      ),
                    ],
                  ),
                ),
                if (!isActive)
                  Container(
                    padding: const EdgeInsets.symmetric(
                      horizontal: 10,
                      vertical: 4,
                    ),
                    decoration: BoxDecoration(
                      color: status == 'settled'
                          ? Colors.green.shade50
                          : Colors.grey.shade100,
                      borderRadius: BorderRadius.circular(10),
                    ),
                    child: Text(
                      status.toUpperCase(),
                      style: TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.w700,
                        color: status == 'settled' ? Colors.green : Colors.grey,
                      ),
                    ),
                  ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        '₹${_fmt(totalCollected)}',
                        style: const TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.w900,
                          color: Color(0xFF6366F1),
                        ),
                      ),
                      Text(
                        'collected of ₹${_fmt(totalExpected)}',
                        style: const TextStyle(
                          fontSize: 11,
                          color: Colors.grey,
                        ),
                      ),
                    ],
                  ),
                ),
                Container(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                    vertical: 6,
                  ),
                  decoration: BoxDecoration(
                    color: pct >= 1.0
                        ? Colors.green.shade50
                        : const Color(0xFF6366F1).withValues(alpha: 0.05),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    '${(pct * 100).toInt()}%',
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w800,
                      color: pct >= 1.0
                          ? Colors.green
                          : const Color(0xFF6366F1),
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            ClipRRect(
              borderRadius: BorderRadius.circular(6),
              child: LinearPercentIndicator(
                lineHeight: 8,
                percent: pct,
                progressColor: pct >= 1.0
                    ? const Color(0xFF059669)
                    : const Color(0xFF6366F1),
                backgroundColor: const Color(
                  0xFF6366F1,
                ).withValues(alpha: 0.08),
                barRadius: const Radius.circular(6),
                padding: EdgeInsets.zero,
              ),
            ),
          ],
        ),
      ),
    );
  }

  String _fmt(dynamic val) {
    if (val == null) return '0';
    final n = (val is num) ? val : double.tryParse(val.toString()) ?? 0;
    if (n >= 100000) return '${(n / 100000).toStringAsFixed(1)}L';
    if (n >= 1000) return '${(n / 1000).toStringAsFixed(1)}K';
    return n.toStringAsFixed(0);
  }
}
