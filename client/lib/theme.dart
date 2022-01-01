import 'package:flutter/material.dart';
import 'dart:ui';

// Application packages
import 'package:go_mud_client/logger.dart';

/// Returns the application theme data
ThemeData getTheme(BuildContext context) {
  final log = getLogger('Theme');

  // Font scaling
  double sizeFactor = 1;
  double sizeDelta = 0;

  var pixelRatio = window.devicePixelRatio;

  // Size in logical pixels
  var logicalScreenSize = window.physicalSize / pixelRatio;
  var logicalWidth = logicalScreenSize.width;
  var logicalHeight = logicalScreenSize.height;

  // Safe area paddings in logical pixels
  var paddingLeft = window.padding.left / window.devicePixelRatio;
  var paddingRight = window.padding.right / window.devicePixelRatio;
  var paddingTop = window.padding.top / window.devicePixelRatio;
  var paddingBottom = window.padding.bottom / window.devicePixelRatio;

  // Safe area in logical pixels
  var safeWidth = logicalWidth - paddingLeft - paddingRight;
  var safeHeight = logicalHeight - paddingTop - paddingBottom;

  log.info("Screen width >$safeWidth<");
  log.info("Screen height >$safeHeight<");

  if (safeWidth < 500) {
    sizeFactor = 0.85;
  } else if (safeWidth < 700) {
    sizeFactor = 0.9;
  } else if (safeWidth < 900) {
    sizeFactor = 0.95;
  }

  return ThemeData(
    textTheme: TextTheme(
      headline1: Theme.of(context).textTheme.headline1!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      headline2: Theme.of(context).textTheme.headline2!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      headline3: Theme.of(context).textTheme.headline3!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      headline4: Theme.of(context).textTheme.headline4!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      headline5: Theme.of(context).textTheme.headline5!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      headline6: Theme.of(context).textTheme.headline6!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      subtitle1: Theme.of(context).textTheme.subtitle1!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      subtitle2: Theme.of(context).textTheme.subtitle2!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      bodyText1: Theme.of(context).textTheme.bodyText1!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      bodyText2: Theme.of(context).textTheme.bodyText2!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      caption: Theme.of(context).textTheme.caption!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      button: Theme.of(context).textTheme.button!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
      overline: Theme.of(context).textTheme.overline!.apply(
            fontSizeFactor: sizeFactor,
            fontSizeDelta: sizeDelta,
          ),
    ),
    // NOTE: Setting the primaryColor provides some consistency
    // across screen colours as not all Flutter components correctly
    // use the defined colorTheme data provided further below.
    primaryColor: const Color(0xFF757575),

    colorScheme: const ColorScheme(
      brightness: Brightness.light,
      primary: Color(0xFF757575),
      primaryVariant: Color(0xFF757575),
      onPrimary: Colors.white,
      secondary: Color(0xFF66bb6a),
      secondaryVariant: Color(0xFF66bb6a),
      onSecondary: Colors.black,
      error: Color(0xFFc63f17),
      onError: Colors.white,
      background: Color(0xFFFAFAFA),
      onBackground: Colors.black,
      surface: Color(0xFFededed),
      onSurface: Colors.black,
    ),
  );
}
