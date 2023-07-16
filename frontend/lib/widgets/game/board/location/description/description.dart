import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/repository/dungeon_action/dungeon_action_repository.dart';
import 'package:go_mud_client/utility.dart';

enum DescriptionOpacity { fadeIn, fadeOut }

class GameLocationDescriptionWidget extends StatefulWidget {
  final DescriptionOpacity fade;
  final DungeonActionRecord dungeonActionRecord;

  const GameLocationDescriptionWidget(
      {Key? key, required this.fade, required this.dungeonActionRecord})
      : super(key: key);

  @override
  State<GameLocationDescriptionWidget> createState() =>
      _GameLocationDescriptionWidgetState();
}

class _GameLocationDescriptionWidgetState
    extends State<GameLocationDescriptionWidget> with TickerProviderStateMixin {
  late final AnimationController _controller = AnimationController(
    duration: const Duration(milliseconds: 500),
    vsync: this,
    lowerBound: 0.0,
    upperBound: 1.0,
  );
  late final Animation<double> _animation = CurvedAnimation(
    parent: _controller,
    curve: Curves.easeIn,
  );

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameLocationDescriptionWidget', 'build');
    log.fine('Building.. ${widget.fade}');

    return BlocConsumer<DungeonActionCubit, DungeonActionState>(
      listener: (BuildContext context, DungeonActionState state) {
        log.fine('listener...');
      },
      builder: (BuildContext context, DungeonActionState state) {
        widget.fade == DescriptionOpacity.fadeIn
            ? _controller.forward(from: 0.0)
            : _controller.reverse(from: 1.0);

        // ignore: avoid_unnecessary_containers
        return FadeTransition(
          opacity: _animation,
          child: Column(
            children: [
              Container(
                alignment: Alignment.center,
                margin: const EdgeInsets.fromLTRB(5, 10, 5, 5),
                child: Text(
                    normaliseName(
                      widget.dungeonActionRecord.actionLocation.locationName,
                    ),
                    style: Theme.of(context).textTheme.headlineSmall),
              ),
              Container(
                alignment: Alignment.center,
                margin: const EdgeInsets.fromLTRB(5, 5, 5, 10),
                child: Text(widget
                    .dungeonActionRecord.actionLocation.locationDescription),
              ),
            ],
          ),
        );
      },
    );
  }
}
