import 'package:logging/logging.dart';

Logger? logger;

void initLogger() {
  Logger.root.level = Level.INFO;
  Logger.root.onRecord.listen((record) {
    // ignore: avoid_print
    print('${record.level.name}: ${record.time}: ${record.loggerName}: ${record.message}');
  });
}

Logger getLogger(String name) {
  logger = Logger(name);
  return logger!;
}
