import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/dungeon/dungeon_list.dart';

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
    final log = getLogger('DungeonContainerWidget');
    log.fine('Initialising state..');

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonContainerWidget');
    log.fine('Building..');

    // ignore: avoid_unnecessary_containers
    return Container(
      child: DungeonListWidget(
        callbacks: widget.callbacks,
      ),
    );
  }
}
