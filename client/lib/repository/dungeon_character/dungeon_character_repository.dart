import 'dart:convert';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/api/api.dart';
import 'package:go_mud_client/repository/repository.dart';

// Package
part 'dungeon_character_record.dart';

abstract class DungeonCharacterRepositoryInterface {
  Future<DungeonCharacterRecord?> getOne(
    String characterID,
  );
  Future<List<DungeonCharacterRecord>> getMany();
}

class DungeonCharacterRepository
    implements DungeonCharacterRepositoryInterface {
  final Map<String, String> config;
  final API api;

  DungeonCharacterRepository({required this.config, required this.api});

  @override
  Future<DungeonCharacterRecord?> getOne(String characterID) async {
    final log = getLogger('DungeonCharacterRepository');

    var response = await api.getCharacter(
      characterID,
    );
    if (response.error != null) {
      log.warning('API responded with error ${response.error}');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    DungeonCharacterRecord? record;
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
        record = DungeonCharacterRecord.fromJson(data[0]);
      }
    }

    return record;
  }

  @override
  Future<List<DungeonCharacterRecord>> getMany() async {
    final log = getLogger('DungeonCharacterRepository');

    var response = await api.getCharacters();
    if (response.error != null) {
      log.warning('API responded with error ${response.error}');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    List<DungeonCharacterRecord> records = [];
    String? responseBody = response.body;
    if (responseBody != null && responseBody.isNotEmpty) {
      Map<String, dynamic> decoded = jsonDecode(responseBody);
      if (decoded['data'] != null) {
        List<dynamic> data = decoded['data'];
        log.fine('Decoded response $data');
        for (var element in data) {
          records.add(DungeonCharacterRecord.fromJson(element));
        }
      }
    }

    return records;
  }
}
