import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';
import 'package:go_mud_client/utility.dart';
import 'package:go_mud_client/cubit/target.dart';

class ObjectButtonWidget extends StatefulWidget {
  final String name;
  const ObjectButtonWidget({
    Key? key,
    required this.name,
  }) : super(key: key);

  @override
  State<ObjectButtonWidget> createState() => _ObjectButtonWidgetState();
}

class _ObjectButtonWidgetState extends State<ObjectButtonWidget> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('ObjectButtonWidget', 'build');
    log.fine('Building button object name ${widget.name}');

    return Container(
      margin: gameButtonMargin,
      child: ElevatedButton(
        onPressed: () {
          selectTarget(context, widget.name);
        },
        style: gameButtonStyle,
        child: Text(normaliseName(widget.name)),
      ),
    );
  }
}
