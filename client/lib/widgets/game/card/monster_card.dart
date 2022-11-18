import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';

import 'package:go_mud_client/widgets/common/bar.dart';
import 'package:go_mud_client/widgets/game/card/equipped.dart';

void displayMonsterCard(
  BuildContext context,
  DungeonActionRecord dungeonActionRecord,
) {
  final log = getLogger('displayMonsterCard');

  log.info(
    'Rendering look monster dialogue',
    dungeonActionRecord.targetMonster!,
  );

  MonsterDetailedData monster = dungeonActionRecord.targetMonster!;

  Widget content = Container(
    alignment: Alignment.center,
    color: Theme.of(context).colorScheme.background,
    padding: const EdgeInsets.all(5),
    child: Column(
      children: <Widget>[
        const Expanded(
          flex: 4,
          child: Text('IMAGE PLACEHOLDER'),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Strength",
            monster.strength,
            monster.currentStrength,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Dexterity",
            monster.dexterity,
            monster.currentDexterity,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Intelligence",
            monster.intelligence,
            monster.currentIntelligence,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Health",
            monster.health,
            monster.currentHealth,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Fatigue",
            monster.fatigue,
            monster.currentFatigue,
            null,
            null,
          ),
        ),
        const Expanded(
          flex: 3,
          child: Text('Description'),
        ),
        Expanded(
          flex: 3,
          child: GameCardEquippedWidget(
            objects: monster.equippedObjects,
          ),
        )
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
          title: Text(monster.name),
          content: content,
          actions: <Widget>[
            TextButton(
              onPressed: () {
                Navigator.of(context).pop();
              },
              child: const Text('Close'),
            ),
          ],
        ),
      );
    },
  );
}
