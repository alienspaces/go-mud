class CubitException implements Exception {
  final String message;
  const CubitException(this.message);

  @override
  String toString() {
    return message;
  }
}
