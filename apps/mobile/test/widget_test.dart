import 'package:flutter_test/flutter_test.dart';
import 'package:horizon_mobile/app/app.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

void main() {
  testWidgets('Horizon app renders', (WidgetTester tester) async {
    await tester.pumpWidget(const ProviderScope(child: HorizonApp()));
    expect(find.text('Horizon'), findsNothing);
  });
}
