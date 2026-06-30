import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'router.dart';
import 'theme.dart';

class HorizonApp extends ConsumerWidget {
  const HorizonApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return MaterialApp.router(
      title: 'Horizon',
      theme: horizonTheme,
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}
