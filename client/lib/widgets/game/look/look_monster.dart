import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/widgets/common/bar.dart';

void displayLookMonsterDialog(
    BuildContext context, DungeonActionRecord dungeonActionRecord) {
  final log = getLogger('displayLookMonsterDialog');

  log.info('Rendering look Monster dialogue');
  Widget content = Container(
    alignment: Alignment.center,
    color: Theme.of(context).colorScheme.background,
    padding: const EdgeInsets.all(5),
    child: Column(
      children: <Widget>[
        const Expanded(
          flex: 10,
          child: Text('IMAGE PLACEHOLDER'),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Strength",
            dungeonActionRecord.targetMonster!.strength,
            dungeonActionRecord.targetMonster!.currentStrength,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Dexterity",
            dungeonActionRecord.targetMonster!.dexterity,
            dungeonActionRecord.targetMonster!.currentDexterity,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Intelligence",
            dungeonActionRecord.targetMonster!.intelligence,
            dungeonActionRecord.targetMonster!.currentIntelligence,
            null,
            null,
          ),
        ),
        const Expanded(
          flex: 10,
          child: Text('Description'),
        ),
      ],
    ),
  );

  showDialog<void>(
    context: context,
    barrierDismissible: false,
    builder: (BuildContext context) {
      return AlertDialog(
        title: Text(dungeonActionRecord.targetMonster!.name),
        content: content,
        actions: <Widget>[
          TextButton(
            child: const Text('Close'),
            onPressed: () {
              Navigator.of(context).pop();
            },
          ),
        ],
      );
    },
  );
}
