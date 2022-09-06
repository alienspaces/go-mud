import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/widgets/game/card/character_card.dart';
import 'package:go_mud_client/widgets/game/card/object_card.dart';
import 'package:go_mud_client/widgets/game/card/monster_card.dart';

class GameCardWidget extends StatelessWidget {
  const GameCardWidget({Key? key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameLocationWidget');
    log.fine('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStatePlaying) {
          var dungeonActionRecord = state.current;
          if (dungeonActionRecord.command == 'look') {
            if (dungeonActionRecord.targetCharacter != null) {
              log.fine('Registering look character dialogue');
              WidgetsBinding.instance.addPostFrameCallback((_) {
                displayCharacterCard(context, dungeonActionRecord);
              });
            } else if (dungeonActionRecord.targetMonster != null) {
              log.fine('Registering look monster dialogue');
              WidgetsBinding.instance.addPostFrameCallback((_) {
                displayMonsterCard(context, dungeonActionRecord);
              });
            } else if (dungeonActionRecord.targetObject != null) {
              log.fine('Registering look object dialogue');
              WidgetsBinding.instance.addPostFrameCallback((_) {
                displayObjectCard(context, dungeonActionRecord);
              });
            }
          }
        }
        return Container();
      },
    );
  }
}
