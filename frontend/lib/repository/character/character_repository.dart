import 'dart:convert';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/api/api.dart';
import 'package:go_mud_client/repository/repository.dart';

// Package
part 'character_record.dart';

abstract class CharacterRepositoryInterface {
  Future<CharacterRecord?> createOne(
    CreateCharacterRecord record,
  );
  Future<CharacterRecord?> getOne(
    String characterID,
  );
  Future<List<CharacterRecord>> getMany();
}

class CharacterRepository implements CharacterRepositoryInterface {
  final Map<String, String> config;
  final API api;

  CharacterRepository({required this.config, required this.api});

  @override
  Future<CharacterRecord?> createOne(CreateCharacterRecord createRecord) async {
    final log = getLogger('CharacterRepository', 'createOne');
    log.info('Creating character ${createRecord.characterName}');

    var response = await api.createCharacter(
      characterName: createRecord.characterName,
      characterStrength: createRecord.characterStrength,
      characterDexterity: createRecord.characterDexterity,
      characterIntelligence: createRecord.characterIntelligence,
    );

    log.info('APIResponse body ${response.body}');
    log.info('APIResponse error ${response.error}');

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
          throw RecordCountException('CharacterRecord');
        }
        record = CharacterRecord.fromJson(data[0]);
      }
    }

    return record;
  }

  @override
  Future<CharacterRecord?> getOne(String characterID) async {
    final log = getLogger('CharacterRepository', 'getOne');

    var response = await api.getCharacter(
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
  Future<List<CharacterRecord>> getMany() async {
    final log = getLogger('CharacterRepository', 'getMany');

    var response = await api.getCharacters();
    if (response.error != null) {
      log.warning('API responded with error ${response.error}');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    List<CharacterRecord> records = [];
    String? responseBody = response.body;
    if (responseBody != null && responseBody.isNotEmpty) {
      Map<String, dynamic> decoded = jsonDecode(responseBody);
      if (decoded['data'] != null) {
        List<dynamic> data = decoded['data'];
        log.fine('Decoded response $data');
        for (var element in data) {
          records.add(CharacterRecord.fromJson(element));
        }
      }
    }

    return records;
  }
}
