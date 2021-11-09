import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';
import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
import 'package:go_mud_client/cubit/character/character_cubit.dart';

class GameDungeonActionWidget extends StatefulWidget {
  const GameDungeonActionWidget({Key? key}) : super(key: key);

  @override
  _GameDungeonActionWidgetState createState() => _GameDungeonActionWidgetState();
}

double gridMemberWidth = 50;
double gridMemberHeight = 50;

class _GameDungeonActionWidgetState extends State<GameDungeonActionWidget> {
  List<Widget> _generateActions(BuildContext context) {
    return [
      _actionWidget(context, 'L', 'look'),
      _actionWidget(context, 'M', 'move'),
      _actionWidget(context, 'G', 'get'),
      _submitActionWidget(context),
    ];
  }

  Widget _actionWidget(BuildContext context, String label, String action) {
    return Container(
      margin: const EdgeInsets.all(2),
      width: gridMemberWidth,
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameDungeonActionWidget');
          log.info('Selecting action >$action<');
          _selectAction(context, action);
        },
        style: ElevatedButton.styleFrom(
          primary: Colors.green,
        ),
        child: Text(label),
      ),
    );
  }

  Widget _submitActionWidget(BuildContext context) {
    return Container(
      margin: const EdgeInsets.all(2),
      child: ElevatedButton(
        onPressed: () {
          final log = getLogger('GameDungeonActionWidget');
          log.info('Submitting action');
          _submitAction(context);
        },
        style: ElevatedButton.styleFrom(
          primary: Colors.green,
        ),
        child: const Text('Submit'),
      ),
    );
  }

  void _selectAction(BuildContext context, String action) {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Selecting action..');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    if (dungeonCubit.dungeonRecord == null) {
      log.warning('Dungeon cubit missing dungeon record, cannot initialise action');
      return;
    }

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    if (characterCubit.characterRecord == null) {
      log.warning('Character cubit missing character record, cannot initialise action');
      return;
    }

    final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
    if (dungeonCommandCubit.action == action) {
      log.info('++ Unselecting action $action');
      dungeonCommandCubit.unselectAction();
      return;
    }

    log.info('++ Selecting action $action');
    dungeonCommandCubit.selectAction(action);
  }

  void _submitAction(BuildContext context) async {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Submitting action..');

    final dungeonCubit = BlocProvider.of<DungeonCubit>(context);
    if (dungeonCubit.dungeonRecord == null) {
      log.warning('Dungeon cubit missing dungeon record, cannot initialise action');
      return;
    }

    final characterCubit = BlocProvider.of<CharacterCubit>(context);
    if (characterCubit.characterRecord == null) {
      log.warning('Character cubit missing character record, cannot initialise action');
      return;
    }

    log.info('++ Submitting action');
    final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
    final dungeonCommandCubit = BlocProvider.of<DungeonCommandCubit>(context);
    await dungeonActionCubit.createAction(
      dungeonCubit.dungeonRecord!.id,
      characterCubit.characterRecord!.id,
      dungeonCommandCubit.command(),
    );
    dungeonCommandCubit.unselectAll();

    // TODO: Loop this using a timer allowing animations to complete
    var moreActions = dungeonActionCubit.playAction();
    log.info('++ More actions >$moreActions<');
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameDungeonGridWidget');
    log.info('Building..');

    return BlocConsumer<DungeonCommandCubit, DungeonCommandState>(
      listener: (BuildContext context, DungeonCommandState state) {
        log.info('listener...');
      },
      builder: (BuildContext context, DungeonCommandState state) {
        return Container(
          decoration: BoxDecoration(
            color: const Color(0xFFDEDEDE),
            border: Border.all(
              color: const Color(0xFFDEDEDE),
            ),
            borderRadius: const BorderRadius.all(Radius.circular(5)),
          ),
          padding: const EdgeInsets.all(1),
          margin: const EdgeInsets.all(5),
          width: gridMemberWidth * 5,
          height: gridMemberHeight * 1,
          child: Row(
            children: _generateActions(context),
          ),
        );
      },
    );
  }
}
