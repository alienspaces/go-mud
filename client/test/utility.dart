import 'dart:io' show Platform;

// Application
import 'package:go_mud_client/api/api.dart';
import 'package:go_mud_client/api/response.dart';
import 'package:go_mud_client/repository/repository.dart';

Map<String, String> getConfig() {
  Map<String, String> envVars = Platform.environment;

  String? serverScheme = envVars['APP_CLIENT_API_SCHEME'];
  String? serverHost = envVars['APP_CLIENT_API_HOST'];
  String? serverPort = envVars['APP_CLIENT_API_PORT'];

  if (serverScheme == null) {
    throw Exception('Test setup failure, environment missing APP_CLIENT_API_SCHEME');
  }

  if (serverHost == null) {
    throw Exception('Test setup failure, environment missing APP_CLIENT_API_HOST');
  }

  if (serverPort == null) {
    throw Exception('Test setup failure, environment missing APP_CLIENT_API_PORT');
  }

  return {
    "serverScheme": envVars["APP_CLIENT_API_SCHEME"] ?? '',
    "serverHost": envVars['APP_CLIENT_API_HOST'] ?? '',
    "serverPort": envVars['APP_CLIENT_API_PORT'] ?? '',
  };
}

class MockAPI implements API {
  @override
  final Map<String, String> config;
  late final String hostname;
  late final String port;

  MockAPI({required this.config});

  @override
  Future<APIResponse> test() async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> getDungeon(String dungeonID) async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> getDungeons() async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> createCharacter(
    String dungeonID, {
    required String name,
    required int strength,
    required int dexterity,
    required int intelligence,
  }) async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> loadCharacter(String dungeonID, String characterID) async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> createDungeonAction(
    String dungeonID,
    String characterID,
    String sentence,
  ) async {
    return Future.value(APIResponse(body: ''));
  }
}

API getMockAPI() {
  return MockAPI(config: getConfig());
}

API getAPI() {
  return API(config: getConfig());
}

RepositoryCollection getRepositories({bool mockAPI = false}) {
  final API api = mockAPI ? getMockAPI() : getAPI();
  return RepositoryCollection(config: getConfig(), api: api);
}
