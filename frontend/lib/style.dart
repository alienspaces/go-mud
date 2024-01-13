import 'package:flutter/material.dart';

// Decorated containers
EdgeInsetsGeometry decoratedContainerMargin = const EdgeInsets.all(2);

// General game buttons
ButtonStyle gameButtonStyle = ElevatedButton.styleFrom(
  backgroundColor: const Color.fromARGB(255, 119, 119, 119),
  shape: const RoundedRectangleBorder(
    borderRadius: BorderRadius.all(Radius.circular(2)),
  ),
);

ButtonStyle gameCancelButtonStyle = ElevatedButton.styleFrom(
  backgroundColor: const Color.fromARGB(255, 119, 119, 119),
  shape: const RoundedRectangleBorder(
    borderRadius: BorderRadius.all(Radius.circular(2)),
  ),
);

TextStyle gameButtonTextStyle(BuildContext context) {
  return Theme.of(context).textTheme.bodyMedium!.copyWith(
        color: Theme.of(context).colorScheme.onPrimary,
      );
}

// Game board buttons
EdgeInsetsGeometry gameBoardButtonMargin = const EdgeInsets.all(2);

ButtonStyle gameBoardButtonStyle = ElevatedButton.styleFrom(
  backgroundColor: const Color.fromARGB(255, 84, 83, 82),
  shape: const RoundedRectangleBorder(
    borderRadius: BorderRadius.all(Radius.circular(2)),
  ),
);

TextStyle gameBoardButtonTextStyle(BuildContext context) {
  return Theme.of(context).textTheme.bodySmall!.copyWith(
        color: Theme.of(context).colorScheme.onPrimary,
      );
}

// Page containers
Color pageContainerBackgroundColor = const Color(0xFFFAFAFA);

// Game panels
EdgeInsetsGeometry gamePanelMargin = const EdgeInsets.all(2);
