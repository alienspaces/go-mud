import 'dart:io' show Platform;

// Application
import 'package:go_mud_client/api/api.dart';
import 'package:go_mud_client/api/response.dart';
import 'package:go_mud_client/repository/repository.dart';

// Test dungeon - Cabin
const String testDungeonID = '34c5b913-3079-42a6-8228-3f1fb8f20dbe';

// Test character - Bolster
const String testCharacterID = '9d7b1765-5c76-461a-aa87-0e2c99803c0f';

Map<String, String> getConfig() {
  Map<String, String> envVars = Platform.environment;

  String? serverScheme = envVars['APP_CLIENT_API_SCHEME'];
  String? serverHost = envVars['APP_CLIENT_API_HOST'];
  String? serverPort = envVars['APP_CLIENT_API_PORT'];

  if (serverScheme == null) {
    throw Exception(
        'Test setup failure, environment missing APP_CLIENT_API_SCHEME');
  }

  if (serverHost == null) {
    throw Exception(
        'Test setup failure, environment missing APP_CLIENT_API_HOST');
  }

  if (serverPort == null) {
    throw Exception(
        'Test setup failure, environment missing APP_CLIENT_API_PORT');
  }

  Map<String, String> config = {
    "serverScheme": envVars["APP_CLIENT_API_SCHEME"] ?? '',
    "serverHost": envVars['APP_CLIENT_API_HOST'] ?? '',
    "serverPort": envVars['APP_CLIENT_API_PORT'] ?? '',
  };

  return config;
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
  Future<APIResponse> createCharacter({
    required String characterName,
    required int characterStrength,
    required int characterDexterity,
    required int characterIntelligence,
  }) async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> getCharacter(String characterID) async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> getCharacters() async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> enterDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> getDungeonCharacter(
      String dungeonID, String characterID) async {
    return Future.value(APIResponse(body: ''));
  }

  @override
  Future<APIResponse> exitDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
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
