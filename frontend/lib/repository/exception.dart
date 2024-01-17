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

class ActionTooEarlyException extends RepositoryException {
  ActionTooEarlyException(String code, String message) : super(code, message);
  ActionTooEarlyException.fromJson(Map<String, dynamic> json)
      : super.fromJson(json);
}

class ActionInvalidCharacterException extends RepositoryException {
  ActionInvalidCharacterException(String code, String message)
      : super(code, message);
  ActionInvalidCharacterException.fromJson(Map<String, dynamic> json)
      : super.fromJson(json);
}

class ActionInvalidDungeonException extends RepositoryException {
  ActionInvalidDungeonException(String code, String message)
      : super(code, message);
  ActionInvalidDungeonException.fromJson(Map<String, dynamic> json)
      : super.fromJson(json);
}

class DuplicateCharacterNameException extends RepositoryException {
  DuplicateCharacterNameException(String code, String message)
      : super(code, message);
  DuplicateCharacterNameException.fromJson(Map<String, dynamic> json)
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

/// Examines the API error and return a specific error class
RepositoryException resolveApiException(String jsonString) {
  List<dynamic> jsonErrors = jsonDecode(jsonString);

  var json = jsonErrors[0];

  var code = json['code'];
  var message = json['message'];

  switch (code) {
    case "character.name_taken":
      {
        return DuplicateCharacterNameException(code, message);
      }
    case "action.too_early":
      {
        return ActionTooEarlyException(code, message);
      }
    case "action.invalid_character":
      {
        return ActionInvalidCharacterException(code, message);
      }
    case "action.invalid_dungeon":
      {
        return ActionInvalidDungeonException(code, message);
      }
  }

  return RepositoryException.fromJson(json);
}
