import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';

import 'package:go_mud_client/cubit/target.dart';

class ObjectButtonWidget extends StatefulWidget {
  final String objectName;
  const ObjectButtonWidget({Key? key, required this.objectName})
      : super(key: key);

  @override
  _ObjectButtonWidgetState createState() => _ObjectButtonWidgetState();
}

class _ObjectButtonWidgetState extends State<ObjectButtonWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('ObjectButtonWidget');
    log.info('Building..');

    return Container(
      margin: gameButtonMargin,
      child: ElevatedButton(
        onPressed: () {
          selectTarget(context, widget.objectName);
        },
        style: gameButtonStyle,
        child: Text(widget.objectName),
      ),
    );
  }
}
