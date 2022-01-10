import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/repository/repository.dart';

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
          flex: 1,
          child: Text('IMAGE PLACEHOLDER'),
        ),
        Expanded(
          flex: 1,
          child: Text("${dungeonActionRecord.targetCharacter!.health}"),
        ),
        Expanded(
          flex: 1,
          child: Text("${dungeonActionRecord.targetCharacter!.fatigue}"),
        ),
      ],
    ),
  );

  showDialog<void>(
    context: context,
    barrierDismissible: false,
    builder: (BuildContext context) {
      return AlertDialog(
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
      );
    },
  );
}
