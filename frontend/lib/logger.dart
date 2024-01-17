import 'package:logging/logging.dart';

void initLogger() {
  Logger.root.level = Level.INFO;
  Logger.root.onRecord.listen((record) {
    // ignore: avoid_print
    print(
        '${record.level.name}: ${record.time}: ${record.loggerName}: ${record.message}');
  });
}

Logger getClassLogger(String className) {
  return Logger(className);
}

Logger getLogger(String className, String? functionName) {
  Logger? logger;

  if (functionName == null) {
    logger = Logger('$className:');
  } else {
    logger = Logger('$className ($functionName):');
  }
  return logger;
}
