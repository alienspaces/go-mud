import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/utility.dart';
import 'package:go_mud_client/navigation.dart';

import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';

import 'package:go_mud_client/repository/character/character_repository.dart';
import 'package:go_mud_client/repository/dungeon/dungeon_repository.dart';

class DungeonListItemWidget extends StatelessWidget {
  final NavigationCallbacks callbacks;
  final CharacterRecord characterRecord;
  final DungeonRecord dungeonRecord;
  const DungeonListItemWidget({
    Key? key,
    required this.callbacks,
    required this.characterRecord,
    required this.dungeonRecord,
  }) : super(key: key);

  /// Sets the current dungeon state to the provided dungeon
  void _enterDungeon(
    BuildContext context,
    String dungeonID,
    String characterID,
  ) async {
    final log = getLogger('DungeonListItemWidget', '_enterDungeon');
    log.info('Enter dungeon $dungeonID with character $characterID');

    final dungeonCharacterCubit =
        BlocProvider.of<DungeonCharacterCubit>(context);

    final characterCubit = BlocProvider.of<CharacterCubit>(context);

    await dungeonCharacterCubit.enterDungeonCharacter(dungeonID, characterID);
    await characterCubit.refreshCharacter(characterID);
  }

  void _exitDungeon(
    BuildContext context,
    String dungeonID,
    String characterID,
  ) async {
    final log = getLogger('DungeonListItemWidget', '_exitDungeon');
    log.info('Exit dungeon $dungeonID with character $characterID');

    final dungeonCharacterCubit =
        BlocProvider.of<DungeonCharacterCubit>(context);

    final characterCubit = BlocProvider.of<CharacterCubit>(context);

    await dungeonCharacterCubit.exitDungeonCharacter(dungeonID, characterID);
    await characterCubit.refreshCharacter(characterID);
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonListItemWidget', 'build');
    log.fine('Dungeon ${dungeonRecord.dungeonID} ${dungeonRecord.dungeonName}');

    if (characterRecord.dungeonID != null &&
        characterRecord.dungeonID != dungeonRecord.dungeonID) {
      log.warning(
          "character is not in this dungeon, not displaying dungeon list item");
      return const SizedBox.shrink();
    }

    ButtonStyle buttonStyle = ElevatedButton.styleFrom(
      padding: const EdgeInsets.fromLTRB(30, 15, 30, 15),
      textStyle: Theme.of(context).textTheme.labelLarge!.copyWith(fontSize: 18),
    );

    List<Widget> children = [];

    // Character not in dungeon

    if (characterRecord.dungeonID == null) {
      children = <Widget>[
        Container(
          margin: const EdgeInsets.all(5),
          child: ElevatedButton(
            onPressed: () => _enterDungeon(
              context,
              dungeonRecord.dungeonID,
              characterRecord.characterID,
            ),
            style: buttonStyle,
            child: const Text('Enter'),
          ),
        ),
      ];
    }

    // Character in this dungeon

    else if (characterRecord.dungeonID != null &&
        characterRecord.dungeonID == dungeonRecord.dungeonID) {
      children = <Widget>[
        Container(
          margin: const EdgeInsets.all(5),
          child: ElevatedButton(
            onPressed: () => _enterDungeon(
              context,
              dungeonRecord.dungeonID,
              characterRecord.characterID,
            ),
            style: buttonStyle,
            child: const Text('Resume'),
          ),
        ),
        Container(
          margin: const EdgeInsets.all(5),
          child: ElevatedButton(
            onPressed: () => _exitDungeon(
              context,
              dungeonRecord.dungeonID,
              characterRecord.characterID,
            ),
            style: buttonStyle,
            child: const Text('Exit'),
          ),
        ),
      ];
    }

    // ignore: avoid_unnecessary_containers
    return Container(
      margin: const EdgeInsets.all(5),
      decoration: BoxDecoration(
        border: Border.all(width: 2),
      ),
      child: Column(
        children: [
          Container(
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Text(
              normaliseName(dungeonRecord.dungeonName),
              style: Theme.of(context).textTheme.titleMedium,
            ),
          ),
          Container(
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Text(
              dungeonRecord.dungeonDescription,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
          ),
          Container(
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: children,
            ),
          ),
        ],
      ),
    );
  }
}
