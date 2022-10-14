import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/common/header.dart';
import 'package:go_mud_client/widgets/dungeon/dungeon.dart';

class DungeonPage extends Page {
  static const String pageName = 'DungeonPage';
  final NavigationCallbacks callbacks;

  const DungeonPage({
    LocalKey key = const ValueKey(DungeonPage.pageName),
    name = DungeonPage.pageName,
    required this.callbacks,
  }) : super(key: key, name: name);

  @override
  Route createRoute(BuildContext context) {
    return PageRouteBuilder(
      settings: this,
      pageBuilder: (context, animation, secondaryAnimation) {
        return DungeonScreen(
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

class DungeonScreen extends StatefulWidget {
  final NavigationCallbacks callbacks;
  static String pageName = 'Dungeon';

  const DungeonScreen({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<DungeonScreen> createState() => _DungeonScreenState();
}

class _DungeonScreenState extends State<DungeonScreen> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonScreen');
    log.fine('Building..');

    return Scaffold(
      appBar: header(context, widget.callbacks),
      resizeToAvoidBottomInset: false,
      body: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        alignment: Alignment.center,
        child: DungeonContainerWidget(callbacks: widget.callbacks),
      ),
    );
  }
}
