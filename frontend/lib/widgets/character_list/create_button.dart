import 'package:flutter/material.dart';
import 'package:logging/logging.dart';

// Application
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';

class CharacterCreateButtonWidget extends StatelessWidget {
  final NavigationCallbacks callbacks;
  final Logger log = getClassLogger('CharacterCreateButtonWidget');

  CharacterCreateButtonWidget({Key? key, required this.callbacks})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterCreateButtonWidget', 'build');
    log.fine('Building..');

    ButtonStyle buttonStyle = ElevatedButton.styleFrom(
      padding: const EdgeInsets.fromLTRB(30, 15, 30, 15),
      textStyle: Theme.of(context).textTheme.labelLarge!.copyWith(fontSize: 18),
    );

    return Container(
      margin: const EdgeInsets.all(5),
      alignment: Alignment.centerRight,
      child: ElevatedButton(
        onPressed: () => callbacks.openCharacterCreatePage(context),
        style: buttonStyle,
        child: const Text('Create Character'),
      ),
    );
  }
}
