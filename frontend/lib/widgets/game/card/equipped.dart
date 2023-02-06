import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

import 'package:go_mud_client/repository/repository.dart';
import 'package:go_mud_client/widgets/game/button/object_button.dart';

class GameCardEquippedWidget extends StatelessWidget {
  final List<ObjectDetailedData>? objects;

  const GameCardEquippedWidget({Key? key, this.objects}) : super(key: key);

  Widget buildEquippedPanel(BuildContext context, List<Widget> objectWidgets) {
    return Row(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: objectWidgets,
    );
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameCardEquippedWidget', 'build');

    log.fine('Building..');

    List<Widget> objectWidgets = [];
    if (objects == null || objects!.isEmpty) {
      return Container();
    }

    for (var i = 0; i < objects!.length; i++) {
      objectWidgets.add(
        ObjectButtonWidget(objectName: objects![i].objectName),
      );
    }

    return Container(
      child: buildEquippedPanel(context, objectWidgets),
    );
  }
}
