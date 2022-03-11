import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

class GameActionNarrativeWidget extends StatefulWidget {
  const GameActionNarrativeWidget({Key? key}) : super(key: key);

  @override
  _GameActionNarrativeWidgetState createState() =>
      _GameActionNarrativeWidgetState();
}

class _GameActionNarrativeWidgetState extends State<GameActionNarrativeWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameActionNarrativeWidget');
    log.info('Building..');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        if (state is DungeonActionStatePlaying) {
          // ignore: avoid_unnecessary_containers
          return IgnorePointer(
            ignoring: true,
            child: Container(
              color: Colors.brown[200]!.withOpacity(0.0),
              alignment: Alignment.center,
              child: Text(': ${state.current.narrative}'.trimRight()),
            ),
          );
        }
        return Container();
      },
    );
  }
}
