import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:horizon_mobile/features/recurring/state/recurring_state.dart';

class RecurringListPage extends ConsumerWidget {
  const RecurringListPage({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final state = ref.watch(recurringListProvider);

    return Scaffold(
      appBar: AppBar(title: const Text('Recurring Transactions')),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/recurring/add'),
        child: const Icon(Icons.add),
      ),
      body: state.when(
        data: (items) {
          if (items.isEmpty) return const Center(child: Text('No recurring transactions setup.'));
          return ListView.builder(
            itemCount: items.length,
            itemBuilder: (context, index) {
              final item = items[index];
              return ListTile(
                title: Text(item.name),
                subtitle: Text('${item.frequency} - \$${item.amount.toStringAsFixed(2)}'),
                trailing: const Icon(Icons.chevron_right),
                onTap: () => context.push('/recurring/${item.id}'),
              );
            },
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (err, _) => Center(child: Text('Error: $err')),
      ),
    );
  }
}
