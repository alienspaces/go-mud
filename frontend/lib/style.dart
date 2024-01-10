import 'package:flutter/material.dart';

// Decorated containers
EdgeInsetsGeometry decoratedContainerMargin = const EdgeInsets.all(2);

// Game buttons
EdgeInsetsGeometry gameButtonMargin = const EdgeInsets.all(2);
ButtonStyle gameButtonStyle = ElevatedButton.styleFrom(
  backgroundColor: Colors.brown,
  shape: const RoundedRectangleBorder(
    borderRadius: BorderRadius.all(Radius.circular(2)),
  ),
);

// Game panels
EdgeInsetsGeometry gamePanelMargin = const EdgeInsets.all(2);
