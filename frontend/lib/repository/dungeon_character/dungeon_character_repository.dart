import 'dart:convert';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/api/api.dart';
import 'package:go_mud_client/repository/repository.dart';

// Package
part 'dungeon_character_record.dart';

abstract class DungeonCharacterRepositoryInterface {
  Future<CharacterRecord?> enterDungeonCharacter(
    String dungeonID,
    String characterID,
  );
  Future<CharacterRecord?> getDungeonCharacter(
    String dungeonID,
    String characterID,
  );
  Future<void> exitDungeonCharacter(
    String dungeonID,
    String characterID,
  );
}

class DungeonCharacterRepository
    implements DungeonCharacterRepositoryInterface {
  final Map<String, String> config;
  final API api;

  DungeonCharacterRepository({required this.config, required this.api});

  @override
  Future<CharacterRecord?> enterDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
    final log =
        getLogger('DungeonCharacterRepository', 'enterDungeonCharacter');

    var response = await api.enterDungeonCharacter(
      dungeonID,
      characterID,
    );

    log.fine('APIResponse body ${response.body}');
    log.fine('APIResponse error ${response.error}');

    if (response.error != null) {
      log.warning('API responded with error ${response.error}');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    CharacterRecord? record;
    String? responseBody = response.body;
    if (responseBody != null && responseBody.isNotEmpty) {
      Map<String, dynamic> decoded = jsonDecode(responseBody);
      if (decoded['data'] != null) {
        List<dynamic> data = decoded['data'];
        log.fine('Decoded response $data');
        if (data.length > 1) {
          log.warning('Unexpected number of records returned');
          throw RecordCountException('Unexpected number of records returned');
        }
        record = CharacterRecord.fromJson(data[0]);
      }
    }

    return record;
  }

  @override
  Future<CharacterRecord?> getDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
    final log = getLogger('DungeonCharacterRepository', 'getDungeonCharacter');

    var response = await api.getDungeonCharacter(
      dungeonID,
      characterID,
    );
    if (response.error != null) {
      log.warning('API responded with error ${response.error}');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    CharacterRecord? record;
    String? responseBody = response.body;
    if (responseBody != null && responseBody.isNotEmpty) {
      Map<String, dynamic> decoded = jsonDecode(responseBody);
      if (decoded['data'] != null) {
        List<dynamic> data = decoded['data'];
        log.fine('Decoded response $data');
        if (data.length > 1) {
          log.warning('Unexpected number of records returned');
          throw RecordCountException('Unexpected number of records returned');
        }
        record = CharacterRecord.fromJson(data[0]);
      }
    }

    return record;
  }

  @override
  Future<void> exitDungeonCharacter(
    String dungeonID,
    String characterID,
  ) async {
    final log = getLogger('DungeonCharacterRepository', 'exitDungeonCharacter');

    var response = await api.exitDungeonCharacter(
      dungeonID,
      characterID,
    );

    log.fine('APIResponse body ${response.body}');
    log.fine('APIResponse error ${response.error}');

    if (response.error != null) {
      log.warning('API responded with error ${response.error}');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    return;
  }
}
