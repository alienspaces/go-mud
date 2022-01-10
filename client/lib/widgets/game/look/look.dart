import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/look/look_character.dart';
import 'package:go_mud_client/widgets/game/look/look_object.dart';
import 'package:go_mud_client/widgets/game/look/look_monster.dart';

class GameLookWidget extends StatelessWidget {
  const GameLookWidget({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameLocationWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStatePlaying) {
          var dungeonActionRecord = state.current;
          if (dungeonActionRecord.command == 'look') {
            if (dungeonActionRecord.targetCharacter != null) {
              log.info('Registering look character dialogue');
              WidgetsBinding.instance?.addPostFrameCallback((_) {
                displayLookCharacterDialog(context, dungeonActionRecord);
              });
            } else if (dungeonActionRecord.targetMonster != null) {
              log.info('Registering look monster dialogue');
              WidgetsBinding.instance?.addPostFrameCallback((_) {
                displayLookMonsterDialog(context, dungeonActionRecord);
              });
            } else if (dungeonActionRecord.targetObject != null) {
              log.info('Registering look object dialogue');
              WidgetsBinding.instance?.addPostFrameCallback((_) {
                displayLookObjectDialog(context, dungeonActionRecord);
              });
            }
          }
        }
        return Container();
      },
    );
  }
}
