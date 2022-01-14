import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/widgets/common/bar.dart';

void displayLookCharacterDialog(
    BuildContext context, DungeonActionRecord dungeonActionRecord) {
  final log = getLogger('displayLookCharacterDialog');

  log.info('Rendering look character dialogue');
  Widget content = Container(
    alignment: Alignment.center,
    color: Theme.of(context).colorScheme.background,
    padding: const EdgeInsets.all(5),
    child: Column(
      children: <Widget>[
        const Expanded(
          flex: 7,
          child: Text('IMAGE PLACEHOLDER'),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Strength",
            dungeonActionRecord.targetCharacter!.strength,
            dungeonActionRecord.targetCharacter!.currentStrength,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Dexterity",
            dungeonActionRecord.targetCharacter!.dexterity,
            dungeonActionRecord.targetCharacter!.currentDexterity,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Intelligence",
            dungeonActionRecord.targetCharacter!.intelligence,
            dungeonActionRecord.targetCharacter!.currentIntelligence,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Health",
            dungeonActionRecord.targetCharacter!.health,
            dungeonActionRecord.targetCharacter!.currentHealth,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Fatigue",
            dungeonActionRecord.targetCharacter!.fatigue,
            dungeonActionRecord.targetCharacter!.currentFatigue,
            null,
            null,
          ),
        ),
        const Expanded(
          flex: 7,
          child: Text('Description'),
        ),
      ],
    ),
  );

  showDialog<void>(
    context: context,
    barrierDismissible: false,
    builder: (BuildContext context) {
      return FractionallySizedBox(
        heightFactor: .8,
        child: AlertDialog(
          title: Text(dungeonActionRecord.targetCharacter!.name),
          content: content,
          actions: <Widget>[
            TextButton(
              child: const Text('Close'),
              onPressed: () {
                Navigator.of(context).pop();
              },
            ),
          ],
        ),
      );
    },
  );
}
