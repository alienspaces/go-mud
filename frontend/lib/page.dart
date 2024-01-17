import 'package:flutter/material.dart';

Widget pageTransitionsBuilder(BuildContext context, Animation<double> animation,
    Animation<double> secondaryAnimation, Widget child) {
  const begin = 0.0;
  const end = 1.0;
  final tween = Tween(begin: begin, end: end);
  final opacityAnimation = animation.drive(tween);
  return FadeTransition(
    opacity: opacityAnimation,
    child: child,
  );
}
