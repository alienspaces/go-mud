import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/dungeon/list/list_item.dart';
import 'package:go_mud_client/repository/character/character_repository.dart';
import 'package:go_mud_client/repository/dungeon/dungeon_repository.dart';

class DungeonListWidget extends StatelessWidget {
  final NavigationCallbacks callbacks;
  final CharacterRecord characterRecord;
  final List<DungeonRecord>? dungeonRecords;

  const DungeonListWidget({
    Key? key,
    required this.callbacks,
    required this.characterRecord,
    required this.dungeonRecords,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonListWidget', 'build');
    log.fine('Building..');

    List<Widget> widgets = [];

    dungeonRecords?.forEach((dungeonRecord) {
      log.fine('Displaying dungeon widget');
      widgets.add(
        DungeonListItemWidget(
          callbacks: callbacks,
          characterRecord: characterRecord,
          dungeonRecord: dungeonRecord,
        ),
      );
    });

    return Column(
      children: widgets,
    );
  }
}
