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
  final log = getLogger('Card', 'displayMonsterCard');

  log.fine(
    'Rendering look monster dialogue',
    dungeonActionRecord.actionTargetMonster!,
  );

  MonsterDetailedData monster = dungeonActionRecord.actionTargetMonster!;

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
            monster.monsterStrength,
            monster.monsterCurrentStrength,
            null,
            null,
            null,
            1,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Dexterity",
            monster.monsterDexterity,
            monster.monsterCurrentDexterity,
            null,
            null,
            null,
            1,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Intelligence",
            monster.monsterIntelligence,
            monster.monsterCurrentIntelligence,
            null,
            null,
            null,
            1,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Health",
            monster.monsterHealth,
            monster.monsterCurrentHealth,
            null,
            null,
            Colors.green,
            1,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Fatigue",
            monster.monsterFatigue,
            monster.monsterCurrentFatigue,
            null,
            null,
            Colors.yellow,
            1,
          ),
        ),
        const Expanded(
          flex: 3,
          child: Text('Description'),
        ),
        Expanded(
          flex: 3,
          child: GameCardEquippedWidget(
            objects: monster.monsterEquippedObjects,
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
          title: Text(monster.monsterName),
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
