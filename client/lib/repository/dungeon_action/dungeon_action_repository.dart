import 'dart:convert';
import 'package:equatable/equatable.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/api/api.dart';
import 'package:go_mud_client/repository/repository.dart';

// Package
part 'dungeon_action_record.dart';

abstract class DungeonActionRepositoryInterface {
  Future<DungeonActionRecord?> create(
      String dungeonID, String characterID, String sentence);
}

class DungeonActionRepository implements DungeonActionRepositoryInterface {
  final Map<String, String> config;
  final API api;

  DungeonActionRepository({required this.config, required this.api});

  @override
  Future<DungeonActionRecord?> create(
      String dungeonID, String characterID, String sentence) async {
    final log = getLogger('DungeonActionRepository');

    var response = await api.createDungeonAction(
      dungeonID,
      characterID,
      sentence,
    );
    if (response.error != null) {
      log.warning('No records returned');
      RepositoryException exception = resolveApiException(response.error!);
      throw exception;
    }

    DungeonActionRecord? record;
    String? responseBody = response.body;
    if (responseBody != null && responseBody.isNotEmpty) {
      Map<String, dynamic> decoded = jsonDecode(responseBody);
      if (decoded['data'] != null) {
        List<dynamic> data = decoded['data'];
        log.fine('Decoded response $data');
        if (data.length > 1) {
          // TODO: (client) Add support for multiple dungeon actions in response
          log.warning('Unexpected number of records returned');
          throw RecordCountException('Unexpected number of records returned');
        }
        record = DungeonActionRecord.fromJson(data[0]);
      }
    }

    return record;
  }
}
