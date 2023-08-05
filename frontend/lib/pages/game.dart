import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/widgets/common/header.dart';
import 'package:go_mud_client/widgets/game/game.dart';

class GamePage extends Page {
  static const String pageName = 'GamePage';
  final NavigationCallbacks callbacks;

  const GamePage({
    LocalKey key = const ValueKey(GamePage.pageName),
    name = GamePage.pageName,
    required this.callbacks,
  }) : super(key: key, name: name);

  @override
  Route createRoute(BuildContext context) {
    return PageRouteBuilder(
      settings: this,
      pageBuilder: (context, animation, secondaryAnimation) {
        return GameScreen(
          callbacks: callbacks,
        );
      },
      transitionDuration: const Duration(milliseconds: 300),
      transitionsBuilder: (context, animation, secondaryAnimation, child) {
        const begin = 0.0;
        const end = 1.0;
        final tween = Tween(begin: begin, end: end);
        final opacityAnimation = animation.drive(tween);
        return FadeTransition(
          opacity: opacityAnimation,
          child: child,
        );
      },
    );
  }
}

class GameScreen extends StatefulWidget {
  final NavigationCallbacks callbacks;
  static String pageName = 'Game';

  const GameScreen({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<GameScreen> createState() => _GameScreenState();
}

class _GameScreenState extends State<GameScreen> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameScreen', 'build');
    log.fine('Building..');

    return Scaffold(
      appBar: header(context, widget.callbacks),
      resizeToAvoidBottomInset: false,
      body: Container(
        alignment: Alignment.center,
        child: GameContainerWidget(callbacks: widget.callbacks),
      ),
    );
  }
}
