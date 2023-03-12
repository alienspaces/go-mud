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
  final log = getLogger('FontScale', 'getFontScale');

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
      displayLarge: Theme.of(context).textTheme.displayLarge!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      displayMedium: Theme.of(context).textTheme.displayMedium!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      displaySmall: Theme.of(context).textTheme.displaySmall!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      headlineMedium: Theme.of(context).textTheme.headlineMedium!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      headlineSmall: Theme.of(context).textTheme.headlineSmall!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      titleLarge: Theme.of(context).textTheme.titleLarge!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      titleMedium: Theme.of(context).textTheme.titleMedium!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      titleSmall: Theme.of(context).textTheme.titleSmall!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      bodyLarge: Theme.of(context).textTheme.bodyLarge!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      bodyMedium: Theme.of(context).textTheme.bodyMedium!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      bodySmall: Theme.of(context).textTheme.bodySmall!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      labelLarge: Theme.of(context).textTheme.labelLarge!.apply(
            fontSizeFactor: fontScale.sizeFactor,
            fontSizeDelta: fontScale.sizeDelta,
          ),
      labelSmall: Theme.of(context).textTheme.labelSmall!.apply(
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
