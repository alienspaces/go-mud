import 'dart:convert';

class RepositoryException implements Exception {
  String code;
  String message;
  RepositoryException(this.code, this.message);

  @override
  String toString() => "RepositoryException: $message";

  RepositoryException.fromJson(Map<String, dynamic> json)
      : code = json['code'],
        message = json['message'];
}

class DuplicateValueException extends RepositoryException {
  DuplicateValueException(String code, String message) : super(code, message);
  DuplicateValueException.fromJson(Map<String, dynamic> json)
      : super.fromJson(json);
}

class RecordCountException extends RepositoryException {
  RecordCountException(String recordName)
      : super('record.unexpected_count', recordName);
  RecordCountException.fromJson(Map<String, dynamic> json)
      : super.fromJson(json);
}

class RecordEmptyException extends RepositoryException {
  RecordEmptyException(String recordName) : super('record.empty', recordName);
  RecordEmptyException.fromJson(Map<String, dynamic> json)
      : super.fromJson(json);
}

// Analyses the API error message string and return a specific error class
RepositoryException resolveApiException(String message) {
  Map<String, dynamic> json = jsonDecode(message);

  // TODO: Specific error classes per error code we care about
  if (message.contains(
      'duplicate key value violates unique constraint \\"dungeon_character_name_key\\"')) {
    var e = DuplicateValueException.fromJson(json);
    e.message = 'Character name is taken, please try another.';
  }

  return RepositoryException.fromJson(json);
}
