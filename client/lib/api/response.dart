class APIResponse {
  String? body;
  String? error;

  APIResponse({this.body, this.error});

  bool isNotEmpty() {
    return body != null || error != null;
  }
}
