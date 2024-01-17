import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/utility.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/style.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
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

  void _unselectCharacter(
    BuildContext context,
  ) async {
    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    characterCubit.unselect();

    // Open character list page
    callbacks.openCharacterListPage(context);
  }

  void _enterDungeon(
    BuildContext context,
    final CharacterRecord characterRecord,
    final DungeonRecord dungeonRecord,
  ) async {
    final log = getLogger('DungeonListItemWidget', '_enterDungeon');
    log.fine('Dungeon ID ${dungeonRecord.dungeonID}');
    log.fine('Character ID ${characterRecord.characterID}');

    final characterCubit = BlocProvider.of<CharacterCubit>(context);

    characterCubit
        .enter(dungeonRecord)
        .then((_) => callbacks.openGamePage(context));
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonListItemWidget', 'build');
    log.info('Display dungeon ID ${dungeonRecord.dungeonID}');
    log.info('Display dungeon Name ${dungeonRecord.dungeonName}');

    if (characterRecord.dungeonID != null &&
        characterRecord.dungeonID != dungeonRecord.dungeonID) {
      log.warning("Character dungeon ID ${characterRecord.dungeonID}");
      log.warning("Not displaying list item");
      return const SizedBox.shrink();
    }

    List<Widget> buttons = [
      Container(
        margin: const EdgeInsets.all(5),
        child: ElevatedButton(
          onPressed: () => _unselectCharacter(
            context,
          ),
          style: gameCancelButtonStyle,
          child: Text(
            'Cancel',
            style: gameButtonTextStyle(context),
          ),
        ),
      ),
    ];

    // Character not in dungeon
    if (characterRecord.dungeonID == null) {
      buttons.add(
        Container(
          margin: const EdgeInsets.all(5),
          child: ElevatedButton(
            onPressed: () => _enterDungeon(
              context,
              characterRecord,
              dungeonRecord,
            ),
            style: gameButtonStyle,
            child: Text(
              'Enter',
              style: gameButtonTextStyle(context),
            ),
          ),
        ),
      );
    }

    // Character in this dungeon
    else if (characterRecord.dungeonID != null &&
        characterRecord.dungeonID == dungeonRecord.dungeonID) {
      buttons.add(
        Container(
          margin: const EdgeInsets.all(5),
          child: ElevatedButton(
            onPressed: () => _enterDungeon(
              context,
              characterRecord,
              dungeonRecord,
            ),
            style: gameButtonStyle,
            child: Text(
              'Resume',
              style: gameButtonTextStyle(context),
            ),
          ),
        ),
      );
    }

    // ignore: avoid_unnecessary_containers
    return Container(
      margin: const EdgeInsets.all(5),
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
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: buttons,
            ),
          ),
        ],
      ),
    );
  }
}
