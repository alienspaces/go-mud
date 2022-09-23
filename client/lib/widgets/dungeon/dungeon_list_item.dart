import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/repository/dungeon/dungeon_repository.dart';

class DungeonListItemWidget extends StatelessWidget {
  final NavigationCallbacks callbacks;
  final DungeonRecord dungeonRecord;
  const DungeonListItemWidget(
      {Key? key, required this.callbacks, required this.dungeonRecord})
      : super(key: key);

  /// Sets the current dungeon state to the provided dungeon
  void _selectDungeon(BuildContext context, DungeonRecord dungeonRecord) {
    final log = getLogger('HomeGameWidget');
    log.fine(
        'Select current dungeon ${dungeonRecord.id} ${dungeonRecord.name}');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    dungeonCubit.selectDungeon(dungeonRecord);

    callbacks.openCharacterPage(context);
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

    return BlocConsumer<DungeonCubit, DungeonState>(
      listener: (BuildContext context, DungeonState state) {
        //
      },
      builder: (BuildContext context, DungeonState state) {
        // ignore: avoid_unnecessary_containers
        return Container(
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
                child: ElevatedButton(
                  onPressed: () => _selectDungeon(context, dungeonRecord),
                  style: buttonStyle,
                  child: const Text('Play'),
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
