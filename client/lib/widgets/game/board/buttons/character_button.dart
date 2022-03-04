import 'package:flutter/material.dart';
// import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';

// import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';
// import 'package:go_mud_client/cubit/dungeon_command/dungeon_command_cubit.dart';
// import 'package:go_mud_client/cubit/character/character_cubit.dart';

import 'package:go_mud_client/widgets/game/board/common.dart';

class CharacterButtonWidget extends StatefulWidget {
  final String characterName;
  const CharacterButtonWidget({Key? key, required this.characterName})
      : super(key: key);

  @override
  _CharacterButtonWidgetState createState() => _CharacterButtonWidgetState();
}

class _CharacterButtonWidgetState extends State<CharacterButtonWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterButtonWidget');
    log.info('Building..');

    return Container(
      margin: gameButtonMargin,
      child: ElevatedButton(
        onPressed: () {
          selectTarget(context, widget.characterName);
        },
        style: gameButtonStyle,
        child: Text(widget.characterName),
      ),
    );
  }
}
