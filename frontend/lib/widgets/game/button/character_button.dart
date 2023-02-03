import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';

import 'package:go_mud_client/cubit/target.dart';

class CharacterButtonWidget extends StatefulWidget {
  final String characterName;
  const CharacterButtonWidget({Key? key, required this.characterName})
      : super(key: key);

  @override
  State<CharacterButtonWidget> createState() => _CharacterButtonWidgetState();
}

class _CharacterButtonWidgetState extends State<CharacterButtonWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterButtonWidget', 'build');
    log.fine('Building..');

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
