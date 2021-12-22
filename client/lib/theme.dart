import 'package:flutter/material.dart';

/// Returns the application theme data
ThemeData getTheme(BuildContext context) {
  return ThemeData(
    // // Default font family.
    // fontFamily: 'Lato',

    // NOTE: Setting the primaryColor provides some consistency
    // across screen colours as not all Flutter components correctly
    // use the defined colorTheme data provided further below.
    primaryColor: const Color(0xFF757575),

    // Default text theme.
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
