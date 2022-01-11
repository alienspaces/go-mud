import 'package:flutter/material.dart';

// Application
import 'package:go_mud_client/logger.dart';

Widget bar(
  String label,
  int fullValue,
  int currentValue,
  double? height,
  double? width,
) {
  final log = getLogger('bar');

  height ??= 14;
  double widthFactor = currentValue / fullValue;

  log.info('label $label widthFactor $widthFactor width $width height $height');

  return Container(
    height: height,
    width: width,
    color: Colors.blue,
    alignment: Alignment.centerLeft,
    child: Stack(
      children: <Widget>[
        FractionallySizedBox(
          widthFactor: widthFactor,
          child: Container(
            color: Colors.blue[100],
            alignment: Alignment.centerLeft,
          ),
        ),
        Text(
          label,
          textAlign: TextAlign.center,
        ),
      ],
    ),
  );
}
