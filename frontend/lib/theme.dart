import 'package:flutter/material.dart';
import 'dart:ui';

// Application packages
import 'package:go_mud_client/logger.dart';

class FontScale {
  double sizeFactor;
  double sizeDelta;
  FontScale({required this.sizeFactor, required this.sizeDelta});
}

FontScale getFontScale() {
  final log = getLogger('getFontScale');

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

  log.fine("Screen width >$safeWidth< height >$safeHeight<");

  if (safeWidth < 370) {
    sizeFactor = 0.85;
  } else if (safeWidth < 390) {
    sizeFactor = 0.9;
  } else if (safeWidth < 410) {
    sizeFactor = 0.95;
  }

  log.fine("Font sizeDelta >$sizeDelta< sizeFactor >$sizeFactor<");

  return FontScale(sizeDelta: sizeDelta, sizeFactor: sizeFactor);
}

/// Returns the application theme data
ThemeData getTheme(BuildContext context) {
  FontScale fontScale = getFontScale();

  return ThemeData(
    textTheme: TextTheme(
      headline1: Theme.of(context).textTheme.headline1!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      headline2: Theme.of(context).textTheme.headline2!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      headline3: Theme.of(context).textTheme.headline3!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      headline4: Theme.of(context).textTheme.headline4!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      headline5: Theme.of(context).textTheme.headline5!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      headline6: Theme.of(context).textTheme.headline6!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      subtitle1: Theme.of(context).textTheme.subtitle1!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      subtitle2: Theme.of(context).textTheme.subtitle2!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      bodyText1: Theme.of(context).textTheme.bodyText1!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      bodyText2: Theme.of(context).textTheme.bodyText2!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      caption: Theme.of(context).textTheme.caption!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      button: Theme.of(context).textTheme.button!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      overline: Theme.of(context).textTheme.overline!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
    ),
    // NOTE: Setting the primaryColor provides some consistency
    // across screen colours as not all Flutter components correctly
    // use the defined colorTheme data provided further below.
    primaryColor: const Color(0xFF757575),

    colorScheme: const ColorScheme(
      brightness: Brightness.light,
      primary: Color(0xFF757575),
      onPrimary: Colors.white,
      secondary: Color(0xFF66bb6a),
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
