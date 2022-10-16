import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:http/retry.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/api/response.dart';

class API {
  final Map<String, String> config;
  late final String _scheme;
  late final String _hostname;
  late final String _port;

  API({required this.config}) {
    _scheme = config['serverScheme'].toString();
    _hostname = config['serverHost'].toString();
    _port = config['serverPort'].toString();
  }

  Future<APIResponse> test() async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());
    log.warning('Testing _hostname $_hostname _port $_port');

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
      );
      response = await client.get(uri,
          headers: {'Content-Type': 'application/json; charset=utf-8'});
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getDungeon(String dungeonID) async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID',
      );

      log.fine('URI $uri');

      response = await client.get(uri,
          headers: {'Content-Type': 'application/json; charset=utf-8'});
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getDungeons() async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());
    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons',
      );

      log.fine('URI $uri');

      response = await client.get(uri,
          headers: {'Content-Type': 'application/json; charset=utf-8'});
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> createCharacter({
    required String name,
    required int strength,
    required int dexterity,
    required int intelligence,
  }) async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/characters',
      );

      log.info('URI $uri');

      String bodyData = jsonEncode({
        "data": {
          "name": name,
          "strength": strength,
          "dexterity": dexterity,
          "intelligence": intelligence,
        },
      });
      log.warning('bodyData $bodyData');

      response = await client.post(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
        body: bodyData,
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    if (response.statusCode != 201 && response.statusCode != 200) {
      log.warning('Failed: $responseBody');
      return APIResponse(error: responseBody);
    }

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getCharacter(String characterID) async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/characters/$characterID',
      );

      log.fine('URI $uri');

      response = await client.get(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getCharacters() async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());
    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/characters',
      );

      log.fine('URI $uri');

      response = await client.get(uri,
          headers: {'Content-Type': 'application/json; charset=utf-8'});
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> createDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID/characters/$characterID/enter',
      );

      log.fine('URI $uri');

      response = await client.post(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID/characters/$characterID',
      );

      log.fine('URI $uri');

      response = await client.get(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> deleteDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID/characters/$characterID/exit',
      );

      log.fine('URI $uri');

      response = await client.post(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> createDungeonAction(
    String dungeonID,
    String characterID,
    String sentence,
  ) async {
    final log = getLogger('API');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID/characters/$characterID/actions',
      );

      log.fine('URI $uri');

      String bodyData = jsonEncode({
        "data": {
          "sentence": sentence,
        },
      });
      log.warning('bodyData $bodyData');

      response = await client.post(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
        body: bodyData,
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');
      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }
}
