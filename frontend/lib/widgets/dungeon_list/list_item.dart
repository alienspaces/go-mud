import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/utility.dart';
import 'package:go_mud_client/navigation.dart';
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
    log.fine('Dungeon ${dungeonRecord.dungeonID} ${dungeonRecord.dungeonName}');

    if (characterRecord.dungeonID != null &&
        characterRecord.dungeonID != dungeonRecord.dungeonID) {
      log.warning("character is in dungeon ID ${characterRecord.dungeonID}");
      log.warning(
          "this dungeon ID ${dungeonRecord.dungeonID}, not displaying list item");
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
              characterRecord,
              dungeonRecord,
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
              characterRecord,
              dungeonRecord,
            ),
            style: buttonStyle,
            child: const Text('Resume'),
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
