import 'package:flutter/material.dart';

// Application
import 'package:go_mud_client/logger.dart';

Widget bar(
  String label,
  int fullValue,
  int currentValue,
  double? height,
  double? width,
  Color? color,
  double? opacity,
) {
  final log = getLogger('bar', null);

  height ??= 14;
  color ??= Colors.blue;
  opacity ??= 1;

  if (fullValue == 0) {
    fullValue = 1;
  }
  double widthFactor = currentValue / fullValue;

  log.fine('label $label widthFactor $widthFactor width $width height $height');

  return Container(
    height: height,
    width: width,
    // color: color ?? Colors.blue,
    alignment: Alignment.centerLeft,
    child: Stack(
      children: <Widget>[
        FractionallySizedBox(
          heightFactor: 1,
          widthFactor: widthFactor,
          child: Container(
            color: color.withOpacity(opacity),
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
