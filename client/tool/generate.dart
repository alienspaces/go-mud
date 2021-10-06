import 'dart:convert';
import 'dart:io';
import 'package:dotenv/dotenv.dart' as dotenv;

// Application
import 'package:go_mud_client/logger.dart';

// Generates lib/env.dart from current environment
// USAGE: dart tool/generate_config.dart
Future<void> main() async {
  initLogger();
  final log = getLogger('main');

  dotenv.load();

  log.warning('APP_CLIENT_API_HOST = ${dotenv.env["APP_CLIENT_API_HOST"]}');
  log.warning('APP_CLIENT_API_PORT = ${dotenv.env["APP_CLIENT_API_PORT"]}');

  final config = {
    'serverHost': dotenv.env['APP_CLIENT_API_HOST'],
    'serverPort': dotenv.env['APP_CLIENT_API_PORT'],
  };

  const filename = 'lib/config.dart';
  await File(filename).writeAsString('final config = ${json.encode(config)};');
}
