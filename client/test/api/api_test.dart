import 'package:flutter_test/flutter_test.dart';

// Application
import 'package:go_mud_client/api/api.dart';

// Local Test Utilities
import '../utility.dart';

void main() {
  test('API should ', () async {
    final api = API(config: getConfig());
    expect(api, isNotNull);
    expect(await api.test(), isInstanceOf<APIResponse>());
  });
}
