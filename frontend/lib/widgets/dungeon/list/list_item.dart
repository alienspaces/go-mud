import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';
import 'package:go_mud_client/repository/dungeon/dungeon_repository.dart';

class DungeonListItemWidget extends StatelessWidget {
  final NavigationCallbacks callbacks;
  final DungeonRecord dungeonRecord;
  const DungeonListItemWidget(
      {Key? key, required this.callbacks, required this.dungeonRecord})
      : super(key: key);

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

    await dungeonCharacterCubit.enterDungeonCharacter(dungeonID, characterID);

    // ignore: use_build_context_synchronously
    callbacks.openGamePage(context);
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

    await dungeonCharacterCubit.exitDungeonCharacter(dungeonID, characterID);

    // ignore: use_build_context_synchronously
    callbacks.openCharacterPage(context);
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonListItemWidget', 'build');
    log.fine(
        'Select current dungeon ${dungeonRecord.dungeonID} ${dungeonRecord.dungeonName}');

    ButtonStyle buttonStyle = ElevatedButton.styleFrom(
      padding: const EdgeInsets.fromLTRB(30, 15, 30, 15),
      textStyle: Theme.of(context).textTheme.labelLarge!.copyWith(fontSize: 18),
    );

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    var characterRecord = characterCubit.characterRecord;
    if (characterRecord == null) {
      return Container();
    }
    if (characterRecord.dungeonID != null &&
        characterRecord.dungeonID != dungeonRecord.dungeonID) {
      // TODO: (client) Could choose to not display dungeon list item at all under this scenario
    }

    List<Widget> children = [];
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
    } else if (characterRecord.dungeonID != null &&
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
            child: Text(dungeonRecord.dungeonName,
                style: Theme.of(context).textTheme.displaySmall),
          ),
          Container(
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Text(dungeonRecord.dungeonDescription),
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
