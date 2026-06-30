import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'dart:async';

final connectivityMonitorProvider = Provider<ConnectivityMonitor>((ref) {
  return ConnectivityMonitor();
});

class ConnectivityMonitor {
  bool _isOnline = true;

  bool get isOnline => _isOnline;

  Stream<bool> get onConnectivityChanged => _connectivityController.stream;
  final _connectivityController = StreamController<bool>.broadcast();

  void setOnline(bool value) {
    if (_isOnline != value) {
      _isOnline = value;
      _connectivityController.add(value);
    }
  }

  void dispose() {
    _connectivityController.close();
  }
}
