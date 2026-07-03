import 'dart:async';
import 'dart:convert';
import 'package:http/http.dart' as http;

class StreamClient {
  final String baseUrl;
  final String Function() getToken;
  
  http.Client? _client;
  final StreamController<Map<String, dynamic>> _controller = StreamController.broadcast();

  StreamClient({required this.baseUrl, required this.getToken});

  Stream<Map<String, dynamic>> get events => _controller.stream;

  void connect() async {
    _client = http.Client();
    final request = http.Request('GET', Uri.parse('$baseUrl/api/v1/stream/events'));
    request.headers['Authorization'] = 'Bearer ${getToken()}';
    request.headers['Accept'] = 'text/event-stream';

    try {
      final response = await _client!.send(request);
      
      response.stream
          .transform(utf8.decoder)
          .transform(const LineSplitter())
          .listen((line) {
        if (line.startsWith('data: ')) {
          final data = line.substring(6);
          if (data.isNotEmpty) {
            try {
              final json = jsonDecode(data);
              _controller.add(json);
            } catch (e) {
              print('Failed to parse SSE JSON: $e');
            }
          }
        }
      }, onError: (error) {
        print('SSE Error: $error');
        _reconnect();
      }, onDone: () {
        print('SSE Done');
        _reconnect();
      });
    } catch (e) {
      print('SSE Connection Error: $e');
      _reconnect();
    }
  }

  void _reconnect() {
    _client?.close();
    Future.delayed(const Duration(seconds: 5), () {
      connect();
    });
  }

  void dispose() {
    _client?.close();
    _controller.close();
  }
}
