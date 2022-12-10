import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/main.dart';

// Local Test Utilities
import './utility.dart';

// Warning: At least one test in this suite creates an HttpClient. When
// running a test suite that uses TestWidgetsFlutterBinding, all HTTP
// requests will return status code 400, and no network request will
// actually be made. Any test expecting a real network connection and
// status code will fail.
// To test code that needs an HttpClient, provide your own HttpClient
// implementation to the code under test, so that your test can
// consistently provide a testable response to the code under test.
// 00:01 +0 -1: Some tests failed.

void main() {
  testWidgets('Application displays expected', (WidgetTester tester) async {
    // Build our app and trigger a frame.
    await tester.pumpWidget(MainApp(
      config: getConfig(),
      repositories: getRepositories(mockAPI: true),
    ));

    // expect(find.text('Dungeon'), findsOneWidget);
  });
}
