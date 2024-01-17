import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/style.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';
import 'package:go_mud_client/repository/character/character_repository.dart';

class CharacterListItemWidget extends StatelessWidget {
  final NavigationCallbacks callbacks;
  final CharacterRecord characterRecord;

  const CharacterListItemWidget({
    Key? key,
    required this.characterRecord,
    required this.callbacks,
  }) : super(key: key);

  void _selectCharacter(
    BuildContext context,
    CharacterRecord characterRecord,
  ) async {
    final log = getLogger('CharacterListItemWidget', '_selectCharacter');
    log.fine('Select character ID >${characterRecord.characterID}<');
    log.fine('Select character Name >${characterRecord.characterName}<');

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    characterCubit.select(characterRecord);

    // Open dungeon list page
    callbacks.openDungeonListPage(context);
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterListItemWidget', 'build');
    log.info('Display character ID ${characterRecord.characterID}');
    log.info('Display character Name ${characterRecord.characterName}');
    log.info('Display character Dungeon ID ${characterRecord.dungeonID}');
    log.info('Display character Dungeon Name ${characterRecord.dungeonName}');

    List<Widget> actionWidgets = <Widget>[
      Container(
        margin: const EdgeInsets.all(5),
        child: ElevatedButton(
          onPressed: () => {null},
          style: gameButtonStyle,
          child: Text(
            'Delete',
            style: gameButtonTextStyle(context),
          ),
        ),
      ),
      Container(
        margin: const EdgeInsets.all(5),
        child: ElevatedButton(
          onPressed: () => _selectCharacter(context, characterRecord),
          style: gameButtonStyle,
          child: Text(
            'Play',
            style: gameButtonTextStyle(context),
          ),
        ),
      ),
    ];

    return Container(
      margin: const EdgeInsets.all(5),
      child: Column(
        children: [
          Container(
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Text(
              characterRecord.characterName,
              style: Theme.of(context).textTheme.titleMedium,
            ),
          ),
          Container(
            margin: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: actionWidgets,
            ),
          ),
        ],
      ),
    );
  }
}
