import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/features/household/state/household_state.dart';

class HouseholdListPage extends ConsumerWidget {
  const HouseholdListPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final householdsAsync = ref.watch(householdsProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('My Households'),
        actions: [
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              // Navigate to create household page or show dialog
              context.push('/households/create');
            },
          ),
        ],
      ),
      body: householdsAsync.when(
        data: (households) {
          if (households.isEmpty) {
            return const Center(
              child: Text('You are not part of any households yet.'),
            );
          }
          return RefreshIndicator(
            onRefresh: () async {
              ref.invalidate(householdsProvider);
            },
            child: ListView.builder(
              itemCount: households.length,
              itemBuilder: (context, index) {
                final hh = households[index];
                return ListTile(
                  leading: const CircleAvatar(
                    child: Icon(Icons.home),
                  ),
                  title: Text(hh.name),
                  subtitle: Text('${hh.memberCount} members • ${hh.status}'),
                  trailing: const Icon(Icons.chevron_right),
                  onTap: () {
                    context.push('/households/${hh.householdId}');
                  },
                );
              },
            ),
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, stack) => Center(child: Text('Error: $error')),
      ),
    );
  }
}
