import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/character_current/character_current_cubit.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/widgets/dungeon/list/list.dart';

class DungeonContainerWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const DungeonContainerWidget({Key? key, required this.callbacks})
      : super(key: key);

  @override
  State<DungeonContainerWidget> createState() => _DungeonContainerWidgetState();
}

class _DungeonContainerWidgetState extends State<DungeonContainerWidget> {
  @override
  void initState() {
    final log = getLogger('DungeonContainerWidget', 'initState');
    log.fine('Initialising state..');
    super.initState();
    _loadDungeons(context);
  }

  void _loadDungeons(BuildContext context) {
    final log = getLogger('DungeonContainerWidget', '_loadDungeons');
    log.fine('Loading dungeons');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    dungeonCubit.loadDungeons();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonContainerWidget', 'build');
    log.fine('Building..');

    final characterCubit = BlocProvider.of<CharacterCurrentCubit>(
      context,
      listen: true,
    );

    var characterRecord = characterCubit.characterRecord;
    if (characterRecord == null) {
      log.warning("character record is null, cannot display dungeons");
      return const SizedBox.shrink();
    }

    return BlocConsumer<DungeonCubit, DungeonState>(
      listener: (context, state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonState state) {
        log.fine('builder...');

        List<Widget> widgets = [];

        if (state is DungeonStateLoaded) {
          // Dungeon list
          widgets.add(
            DungeonListWidget(
              callbacks: widget.callbacks,
              characterRecord: characterRecord,
              dungeonRecords: state.dungeonRecords,
            ),
          );
        } else if (state is DungeonStateLoadError) {
          widgets.add(
            // ignore: avoid_unnecessary_containers
            Container(
              child: ElevatedButton(
                onPressed: () => _loadDungeons(context),
                child: const Text('Load Dungeons'),
              ),
            ),
          );
        } else {
          widgets.add(
            // ignore: avoid_unnecessary_containers
            Container(
              child: const Text('Loading dungeons...'),
            ),
          );
        }

        return Container(
          color: Colors.orange,
          child: Column(
            children: widgets,
          ),
        );
      },
    );
  }
}
