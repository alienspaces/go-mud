import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';
import 'package:go_mud_client/utility.dart';
import 'package:go_mud_client/cubit/target.dart';
import 'package:go_mud_client/widgets/common/bar.dart';

class CharacterButtonWidget extends StatefulWidget {
  final String name;
  final int health;
  final int currentHealth;
  final int fatigue;
  final int currentFatigue;
  const CharacterButtonWidget({
    Key? key,
    required this.name,
    required this.health,
    required this.currentHealth,
    required this.fatigue,
    required this.currentFatigue,
  }) : super(key: key);

  @override
  State<CharacterButtonWidget> createState() => _CharacterButtonWidgetState();
}

class _CharacterButtonWidgetState extends State<CharacterButtonWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterButtonWidget', 'build');
    log.fine('Building..');

    log.info(
      'Character ${widget.name} '
      'health ${widget.health} '
      'currentHealth ${widget.currentHealth} '
      'fatigue ${widget.fatigue} '
      'currentFatigue ${widget.currentFatigue}',
    );

    return Stack(
      fit: StackFit.expand,
      children: [
        Container(
          margin: gameButtonMargin,
          child: ElevatedButton(
            onPressed: () {
              selectTarget(context, widget.name);
            },
            style: gameButtonStyle,
            child: Text(normaliseName(widget.name)),
          ),
        ),
        Container(
          margin: gameButtonMargin,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.end,
            children: [
              bar(
                "",
                widget.health,
                widget.currentHealth,
                null,
                null,
                Colors.green,
                0.5,
              ),
              bar(
                "",
                widget.fatigue,
                widget.currentFatigue,
                null,
                null,
                Colors.yellow,
                0.5,
              ),
            ],
          ),
        )
      ],
    );
  }
}
