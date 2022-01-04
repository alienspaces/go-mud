class RepositoryException implements Exception {
  final String message;
  RepositoryException(this.message);

  @override
  String toString() => "RepositoryException: $message";
}

class DuplicateValueException extends RepositoryException {
  DuplicateValueException(String message) : super(message);
}

class RecordCountException extends RepositoryException {
  RecordCountException(String message) : super(message);
}

class RecordEmptyException extends RepositoryException {
  RecordEmptyException(String message) : super(message);
}

// Analyses the API error message string and return a specific error class
RepositoryException resolveApiException(String message) {
  if (message.contains(
      'duplicate key value violates unique constraint \\"dungeon_character_name_key\\"')) {
    return DuplicateValueException(
        'Character name is taken, please try another.');
  }

  return RepositoryException(message);
}
