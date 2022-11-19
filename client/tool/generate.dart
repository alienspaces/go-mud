import 'dart:convert';
import 'dart:io';
import 'package:dotenv/dotenv.dart';

// Application
import 'package:go_mud_client/logger.dart';

// Generates lib/env.dart from current environment
// USAGE: dart tool/generate_config.dart
Future<void> main() async {
  initLogger();
  final log = getLogger('main');

  var env = DotEnv();
  env.load();

  log.warning('APP_CLIENT_API_HOST = ${env["APP_CLIENT_API_HOST"]}');
  log.warning('APP_CLIENT_API_PORT = ${env["APP_CLIENT_API_PORT"]}');
  log.warning('APP_CLIENT_API_HOST = ${env["APP_CLIENT_API_SCHEME"]}');

  final config = {
    'serverHost': env['APP_CLIENT_API_HOST'],
    'serverScheme': env['APP_CLIENT_API_SCHEME'],
    'serverPort': env['APP_CLIENT_API_PORT'],
  };

  const filename = 'lib/config.dart';
  await File(filename).writeAsString('final config = ${json.encode(config)};');
}
