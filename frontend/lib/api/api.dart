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
    final log = getLogger('API', 'test');
    final client = RetryClient(http.Client());
    log.warning('Testing _hostname $_hostname _port $_port');

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
      );

      log.info('URI $uri');

      // TODO: Remove
      print('URI $uri');

      response = await client.get(uri,
          headers: {'Content-Type': 'application/json; charset=utf-8'});
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getDungeon(String dungeonID) async {
    final log = getLogger('API', 'getDungeon');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID',
      );

      log.info('URI $uri');

      // TODO: Remove
      print('URI $uri');

      response = await client.get(uri,
          headers: {'Content-Type': 'application/json; charset=utf-8'});
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getDungeons() async {
    final log = getLogger('API', 'getDungeons');
    final client = RetryClient(http.Client());
    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons',
      );

      log.info('URI $uri');

      // TODO: Remove
      print('URI $uri');

      response = await client.get(uri,
          headers: {'Content-Type': 'application/json; charset=utf-8'});
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> createCharacter({
    required String characterName,
    required int characterStrength,
    required int characterDexterity,
    required int characterIntelligence,
  }) async {
    final log = getLogger('API', 'createCharacter');
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

      // TODO: Remove
      print('URI $uri');

      String bodyData = jsonEncode({
        "data": {
          "character_name": characterName,
          "character_strength": characterStrength,
          "character_dexterity": characterDexterity,
          "character_intelligence": characterIntelligence,
        },
      });
      log.info('bodyData $bodyData');

      response = await client.post(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
        body: bodyData,
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    if (response.statusCode != 201 && response.statusCode != 200) {
      log.warning('Failed: $responseBody');

      // TODO: Remove
      print('Failed with response: $responseBody');

      return APIResponse(error: responseBody);
    }

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getCharacter(String characterID) async {
    final log = getLogger('API', 'getCharacter');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/characters/$characterID',
      );

      log.info('URI $uri');

      // TODO: Remove
      print('URI $uri');

      response = await client.get(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    if (response.statusCode != 200) {
      log.warning('Failed: $responseBody');

      // TODO: Remove
      print('Failed with response: $responseBody');

      return APIResponse(error: responseBody);
    }

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getCharacters() async {
    final log = getLogger('API', 'getCharacters');
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

      // TODO: Remove
      print('URI $uri');

      response = await client.get(uri,
          headers: {'Content-Type': 'application/json; charset=utf-8'});
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    if (response.statusCode != 200) {
      log.warning('Failed: $responseBody');

      // TODO: Remove
      print('Failed with response: $responseBody');

      return APIResponse(error: responseBody);
    }

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> enterDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
    final log = getLogger('API', 'enterDungeonCharacter');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID/characters/$characterID/enter',
      );

      log.info('URI $uri');

      // TODO: Remove
      print('URI $uri');

      response = await client.post(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    log.info('Response: $responseBody');

    if (response.statusCode != 201 && response.statusCode != 200) {
      log.warning('Failed: $responseBody');

      // TODO: Remove
      print('Failed with response: $responseBody');

      return APIResponse(error: responseBody);
    }

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> getDungeonCharacter(
      String dungeonID, String characterID) async {
    final log = getLogger('API', 'getDungeonCharacter');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID/characters/$characterID',
      );

      log.info('URI $uri');

      // TODO: Remove
      print('URI $uri');

      response = await client.get(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    if (response.statusCode != 200) {
      log.warning('Failed: $responseBody');

      // TODO: Remove
      print('Failed with response: $responseBody');

      return APIResponse(error: responseBody);
    }

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }

  Future<APIResponse> exitDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
    final log = getLogger('API', 'exitDungeonCharacter');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID/characters/$characterID/exit',
      );

      log.info('URI $uri');

      // TODO: Remove
      print('URI $uri');

      response = await client.post(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

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
    final log = getLogger('API', 'createDungeonAction');
    final client = RetryClient(http.Client());

    http.Response? response;
    try {
      Uri uri = Uri(
        scheme: _scheme,
        host: _hostname,
        port: int.parse(_port),
        path: '/api/v1/dungeons/$dungeonID/characters/$characterID/actions',
      );

      log.info('URI $uri');

      // TODO: Remove
      print('URI $uri');

      String bodyData = jsonEncode({
        "data": {
          "sentence": sentence,
        },
      });
      log.info('bodyData $bodyData');

      response = await client.post(
        uri,
        headers: {
          'Content-Type': 'application/json; charset=utf-8',
        },
        body: bodyData,
      );
    } catch (err) {
      log.warning('Failed: ${err.toString()}');

      // TODO: Remove
      print('Failed: ${err.toString()}');

      return APIResponse(error: err.toString());
    } finally {
      client.close();
    }

    String responseBody = response.body;

    if (response.statusCode != 201 && response.statusCode != 200) {
      log.warning('Failed: $responseBody');

      // TODO: Remove
      print('Failed with response: $responseBody');

      return APIResponse(error: responseBody);
    }

    log.fine('Response: $responseBody');

    return APIResponse(body: responseBody);
  }
}
