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
  ) {
    final log = getLogger('DungeonListItemWidget');
    log.info('Enter dungeon $dungeonID with character $characterID');

    final dungeonCharacterCubit =
        BlocProvider.of<DungeonCharacterCubit>(context);

    dungeonCharacterCubit.enterDungeonCharacter(dungeonID, characterID);

    callbacks.openGamePage(context);
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonListItemWidget');
    log.fine(
        'Select current dungeon ${dungeonRecord.id} ${dungeonRecord.name}');

    ButtonStyle buttonStyle = ElevatedButton.styleFrom(
      padding: const EdgeInsets.fromLTRB(30, 15, 30, 15),
      textStyle: Theme.of(context).textTheme.button!.copyWith(fontSize: 18),
    );

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    var characterRecord = characterCubit.characterRecord;
    if (characterRecord == null) {
      return Container();
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
            child: Text(dungeonRecord.name,
                style: Theme.of(context).textTheme.headline3),
          ),
          Container(
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Text(dungeonRecord.description),
          ),
          Container(
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: <Widget>[
                Container(
                  margin: const EdgeInsets.all(5),
                  child: ElevatedButton(
                    onPressed: () => _enterDungeon(
                      context,
                      dungeonRecord.id,
                      characterRecord.id,
                    ),
                    style: buttonStyle,
                    child: const Text('Enter'),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
