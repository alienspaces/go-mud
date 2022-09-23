import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/widgets/dungeon/dungeon_list_item.dart';

class DungeonListWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const DungeonListWidget({Key? key, required this.callbacks})
      : super(key: key);

  @override
  State<DungeonListWidget> createState() => _DungeonListWidgetState();
}

class _DungeonListWidgetState extends State<DungeonListWidget> {
  @override
  void initState() {
    final log = getLogger('DungeonListWidget');
    log.fine('Initialising state..');

    super.initState();

    _loadDungeons(context);
  }

  void _loadDungeons(BuildContext context) {
    final log = getLogger('DungeonListWidget');
    log.fine('Loading dungeons');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    dungeonCubit.loadDungeons();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonListWidget');
    log.fine('Building..');

    return BlocConsumer<DungeonCubit, DungeonState>(
      listener: (context, state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonState state) {
        log.fine('builder...');
        List<Widget> widgets = [];

        if (state is DungeonStateLoaded) {
          // Dungeon list
          state.dungeonRecords?.forEach((dungeonRecord) {
            log.fine('Displaying dungeon widget');
            widgets.add(
              // ignore: avoid_unnecessary_containers
              Container(
                child: DungeonListItemWidget(
                  callbacks: widget.callbacks,
                  dungeonRecord: dungeonRecord,
                ),
              ),
            );
          });
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

        // ignore: avoid_unnecessary_containers
        return Container(
          child: Column(
            children: widgets,
          ),
        );
      },
    );
  }
}
