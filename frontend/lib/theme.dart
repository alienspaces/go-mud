import 'package:flutter/material.dart';
import 'dart:ui';

// Application packages
import 'package:go_mud_client/logger.dart';

class FontScale {
  double sizeFactor;
  double sizeDelta;
  FontScale({required this.sizeFactor, required this.sizeDelta});
}

FontScale getFontScale(BuildContext context) {
  final log = getLogger('FontScale', 'getFontScale');

  FlutterView view = View.of(context);

  // Font scaling
  double sizeFactor = 1;
  double sizeDelta = 0;

  var pixelRatio = view.devicePixelRatio;

  // Size in logical pixels
  var logicalScreenSize = view.physicalSize / pixelRatio;
  var logicalWidth = logicalScreenSize.width;
  var logicalHeight = logicalScreenSize.height;

  // Safe area paddings in logical pixels
  var paddingLeft = view.padding.left / view.devicePixelRatio;
  var paddingRight = view.padding.right / view.devicePixelRatio;
  var paddingTop = view.padding.top / view.devicePixelRatio;
  var paddingBottom = view.padding.bottom / view.devicePixelRatio;

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
  FontScale fontScale = getFontScale(context);

  ThemeData theme = Theme.of(context);

  return ThemeData(
    textTheme: TextTheme(
      displayLarge: theme.textTheme.displayLarge!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      displayMedium: theme.textTheme.displayMedium!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      displaySmall: theme.textTheme.displaySmall!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      headlineMedium: theme.textTheme.headlineMedium!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      headlineSmall: theme.textTheme.headlineSmall!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      titleLarge: theme.textTheme.titleLarge!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      titleMedium: theme.textTheme.titleMedium!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      titleSmall: theme.textTheme.titleSmall!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      bodyLarge: theme.textTheme.bodyLarge!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      bodyMedium: theme.textTheme.bodyMedium!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      bodySmall: theme.textTheme.bodySmall!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      labelLarge: theme.textTheme.labelLarge!.apply(
        fontSizeFactor: fontScale.sizeFactor,
        fontSizeDelta: fontScale.sizeDelta,
      ),
      labelSmall: theme.textTheme.labelSmall!.apply(
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
