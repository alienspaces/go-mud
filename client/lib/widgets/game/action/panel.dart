import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_character/dungeon_character_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/style.dart';

class GameActionPanelWidget extends StatefulWidget {
  const GameActionPanelWidget({Key? key}) : super(key: key);

  @override
  State<GameActionPanelWidget> createState() => _GameActionPanelWidgetState();
}

double gridMemberWidth = 0;
double gridMemberHeight = 0;

class _GameActionPanelWidgetState extends State<GameActionPanelWidget> {
  List<Widget> _generateActions(BuildContext context) {
    return [
      _actionWidget(context, 'Move', 'move'),
      _actionWidget(context, 'Look', 'look'),
      _actionWidget(context, 'Equip', 'equip'),
      _actionWidget(context, 'Stash', 'stash'),
      _actionWidget(context, 'Drop', 'drop'),
      _actionWidget(context, 'Use', 'use'),
    ];
  }

  Widget _actionWidget(BuildContext context, String label, String action) {
    return Container(
      margin: gameButtonMargin,
      width: gridMemberWidth,
      height: gridMemberHeight,
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameActionPanelWidget');
          log.fine('Selecting action >$action<');
          _selectAction(context, action);
        },
        style: gameButtonStyle,
        child: Text(
          label,
        ),
      ),
    );
  }

  Widget _submitActionWidget(BuildContext context) {
    return Container(
      width: gridMemberWidth * 2,
      height: gridMemberHeight * 2,
      margin: const EdgeInsets.fromLTRB(5, 5, 5, 5),
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameActionPanelWidget');
          log.fine('Submitting action');
          _submitAction(context);
        },
        style: ElevatedButton.styleFrom(
          backgroundColor: Colors.green,
        ),
        child: const Text('Play'),
      ),
    );
  }

  void _selectAction(BuildContext context, String action) {
    final log = getLogger('GameActionPanelWidget');
    log.fine('Selecting action..');

    final dungeonCharacterCubit =
        BlocProvider.of<DungeonCharacterCubit>(context);
    if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
      log.warning(
          'Dungeon character cubit missing dungeon character record, cannot initialise action');
      return;
    }

    final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
    if (dungeonCommandCubit.action == action) {
      log.fine('++ Unselecting action $action');
      dungeonCommandCubit.unselectAction();
      return;
    }

    log.fine('++ Selecting action $action');
    dungeonCommandCubit.selectAction(action);
  }

  void _submitAction(BuildContext context) async {
    final log = getLogger('GameActionPanelWidget');
    log.fine('Submitting action..');

    final dungeonCharacterCubit =
        BlocProvider.of<DungeonCharacterCubit>(context);
    if (dungeonCharacterCubit.dungeonCharacterRecord == null) {
      log.warning(
          'Dungeon character cubit missing dungeon character record, cannot initialise action');
      return;
    }

    log.fine('++ Submitting action');
    final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
    final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);

    await dungeonActionCubit.createAction(
      dungeonCharacterCubit.dungeonCharacterRecord!.dungeonID,
      dungeonCharacterCubit.dungeonCharacterRecord!.characterID,
      dungeonCommandCubit.command(),
    );
    dungeonCommandCubit.unselectAll();

    // TODO: (client) Loop this using a timer allowing animations to complete
    var moreActions = dungeonActionCubit.playAction();
    log.fine('++ More actions >$moreActions<');
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameActionPanelWidget');
    log.fine('Building..');

    return LayoutBuilder(
      builder: (BuildContext context, BoxConstraints constraints) {
        log.fine(
            'Building width ${constraints.maxWidth} height ${constraints.maxHeight}');

        // Set grid member dimensions
        gridMemberWidth = (constraints.maxWidth / 5) - 2;
        gridMemberHeight = (constraints.maxHeight / 2) - 2;
        if (gridMemberHeight > gridMemberWidth) {
          gridMemberHeight = gridMemberWidth;
        }
        if (gridMemberWidth > gridMemberHeight) {
          gridMemberWidth = gridMemberHeight;
        }

        double gridWidth = gridMemberWidth * 5;
        double gridHeight = gridMemberHeight * 2;

        return Container(
          width: gridWidth,
          height: gridHeight,
          decoration: BoxDecoration(
            color: const Color(0xFFDEDEDE),
            border: Border.all(
              color: const Color(0xFFDEDEDE),
            ),
            borderRadius: const BorderRadius.all(Radius.circular(5)),
          ),
          padding: const EdgeInsets.all(1),
          // margin: const EdgeInsets.all(5),
          child: Row(
            children: <Widget>[
              Expanded(
                flex: 3,
                child: GridView.count(
                  crossAxisCount: 3,
                  children: _generateActions(context),
                ),
              ),
              Expanded(
                flex: 2,
                child: _submitActionWidget(
                  context,
                ),
              ),
            ],
          ),
        );
      },
    );
  }
}
