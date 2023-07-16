String normaliseName(String name) {
  return name.replaceFirst(RegExp(r'\(.*?\)'), '');
}
