import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';

import 'package:go_mud_client/cubit/target.dart';

class MonsterButtonWidget extends StatefulWidget {
  final String monsterName;
  const MonsterButtonWidget({Key? key, required this.monsterName})
      : super(key: key);

  @override
  _MonsterButtonWidgetState createState() => _MonsterButtonWidgetState();
}

class _MonsterButtonWidgetState extends State<MonsterButtonWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('MonsterButtonWidget');
    log.info('Building..');

    return Container(
      margin: gameButtonMargin,
      child: ElevatedButton(
        onPressed: () {
          selectTarget(context, widget.monsterName);
        },
        style: gameButtonStyle,
        child: Text(widget.monsterName),
      ),
    );
  }
}
