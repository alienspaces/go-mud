import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/main.dart';

// Local Test Utilities
import './utility.dart';

void main() {
  testWidgets('Application', (WidgetTester tester) async {
    await tester.pumpWidget(MainApp(
      config: getConfig(),
      repositories: getRepositories(mockAPI: true),
    ));
  });
}
