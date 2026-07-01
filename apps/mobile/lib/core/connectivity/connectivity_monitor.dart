import 'dart:async';
import 'package:connectivity_plus/connectivity_plus.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

final connectivityMonitorProvider = Provider<ConnectivityMonitor>((ref) {
  final monitor = ConnectivityMonitor();
  ref.onDispose(() => monitor.dispose());
  return monitor;
});

class ConnectivityMonitor {
  final Connectivity _connectivity = Connectivity();
  bool _isOnline = true;
  StreamSubscription<List<ConnectivityResult>>? _subscription;
  final _connectivityController = StreamController<bool>.broadcast();

  bool get isOnline => _isOnline;
  Stream<bool> get onConnectivityChanged => _connectivityController.stream;

  ConnectivityMonitor() {
    _init();
  }

  Future<void> _init() async {
    try {
      final results = await _connectivity.checkConnectivity();
      _updateStatus(results);
    } catch (_) {
      _isOnline = true;
    }

    _subscription = _connectivity.onConnectivityChanged.listen(_updateStatus);
  }

  void _updateStatus(List<ConnectivityResult> results) {
    final online = results.any((r) => r != ConnectivityResult.none);
    if (_isOnline != online) {
      _isOnline = online;
      if (!_connectivityController.isClosed) {
        _connectivityController.add(online);
      }
    }
  }

  void dispose() {
    _subscription?.cancel();
    _connectivityController.close();
  }
}
