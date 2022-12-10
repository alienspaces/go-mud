import 'dart:convert';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/api/api.dart';
import 'package:go_mud_client/repository/repository.dart';

// Package
part 'dungeon_record.dart';

abstract class DungeonRepositoryInterface {
  Future<DungeonRecord?> getOne(String dungeonID);
  Future<List<DungeonRecord?>> getMany();
}

class DungeonRepository implements DungeonRepositoryInterface {
  final Map<String, String> config;
  final API api;

  DungeonRepository({required this.config, required this.api});

  @override
  Future<DungeonRecord?> getOne(String dungeonID) async {
    final log = getLogger('DungeonRepository');

    var response = await api.getDungeon(dungeonID);
    if (response.error != null) {
      log.warning('API responded with error ${response.error}');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    DungeonRecord? record;
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
        record = DungeonRecord.fromJson(data[0]);
      }
    }

    return record;
  }

  @override
  Future<List<DungeonRecord>> getMany() async {
    final log = getLogger('DungeonRepository');

    var response = await api.getDungeons();
    if (response.error != null) {
      log.warning('API responded with error ${response.error}');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    List<DungeonRecord> records = [];
    String? responseBody = response.body;
    if (responseBody != null && responseBody.isNotEmpty) {
      Map<String, dynamic> decoded = jsonDecode(responseBody);
      if (decoded['data'] != null) {
        List<dynamic> data = decoded['data'];
        log.fine('Decoded response $data');
        for (var element in data) {
          records.add(DungeonRecord.fromJson(element));
        }
      }
    }

    return records;
  }
}
