import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';

import 'package:go_mud_client/widgets/common/bar.dart';
import 'package:go_mud_client/widgets/game/card/equipped.dart';

void displayCharacterCard(
    BuildContext context, DungeonActionRecord dungeonActionRecord) {
  final log = getLogger('Card', 'displayCharacterCard');

  log.fine('Rendering look character dialogue');

  CharacterDetailedData character = dungeonActionRecord.actionTargetCharacter!;

  Widget content = Container(
    alignment: Alignment.center,
    color: Theme.of(context).colorScheme.background,
    padding: const EdgeInsets.all(5),
    child: Column(
      children: <Widget>[
        // Image
        const Expanded(
          flex: 7,
          child: Text('IMAGE PLACEHOLDER'),
        ),
        // Statistics
        Expanded(
          flex: 1,
          child: bar(
            "Strength",
            character.characterStrength,
            character.characterCurrentStrength,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Dexterity",
            character.characterDexterity,
            character.characterCurrentDexterity,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Intelligence",
            character.characterIntelligence,
            character.characterCurrentIntelligence,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Health",
            character.characterHealth,
            character.characterCurrentHealth,
            null,
            null,
          ),
        ),
        Expanded(
          flex: 1,
          child: bar(
            "Fatigue",
            character.characterFatigue,
            character.characterCurrentFatigue,
            null,
            null,
          ),
        ),
        // Description
        const Expanded(
          flex: 3,
          child: Text('Description'),
        ),
        Expanded(
          flex: 3,
          child: GameCardEquippedWidget(
            objects: character.characterEquippedObjects,
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
          title: Text(character.characterName),
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
