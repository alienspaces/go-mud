import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/common/header.dart';
import 'package:go_mud_client/widgets/character/character.dart';

class CharacterPage extends Page {
  static const String pageName = 'CharacterPage';
  final NavigationCallbacks callbacks;

  const CharacterPage({
    LocalKey key = const ValueKey(CharacterPage.pageName),
    name = CharacterPage.pageName,
    required this.callbacks,
  }) : super(key: key, name: name);

  @override
  Route createRoute(BuildContext context) {
    return PageRouteBuilder(
      settings: this,
      pageBuilder: (context, animation, secondaryAnimation) {
        return CharacterScreen(
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

class CharacterScreen extends StatefulWidget {
  final NavigationCallbacks callbacks;
  static String pageName = 'Character';

  const CharacterScreen({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<CharacterScreen> createState() => _CharacterScreenState();
}

class _CharacterScreenState extends State<CharacterScreen> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterScreen', 'build');
    log.fine('Building..');

    return Scaffold(
      appBar: header(context, widget.callbacks),
      resizeToAvoidBottomInset: false,
      body: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        alignment: Alignment.center,
        child: CharacterContainerWidget(callbacks: widget.callbacks),
      ),
    );
  }
}
