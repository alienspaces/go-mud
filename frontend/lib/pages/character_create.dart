import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/page.dart';
import 'package:go_mud_client/widgets/common/header.dart';
import 'package:go_mud_client/widgets/character_create/container.dart';

class CharacterCreatePage extends Page {
  static const String pageName = 'CharacterCreatePage';
  final NavigationCallbacks callbacks;

  const CharacterCreatePage({
    LocalKey key = const ValueKey(CharacterCreatePage.pageName),
    name = CharacterCreatePage.pageName,
    required this.callbacks,
  }) : super(key: key, name: name);

  @override
  Route createRoute(BuildContext context) {
    return PageRouteBuilder(
      settings: this,
      pageBuilder: (context, animation, secondaryAnimation) {
        return CharacterCreateScreen(
          callbacks: callbacks,
        );
      },
      transitionDuration: const Duration(milliseconds: 300),
      transitionsBuilder: pageTransitionsBuilder,
    );
  }
}

class CharacterCreateScreen extends StatefulWidget {
  final NavigationCallbacks callbacks;
  static String pageName = 'CharacterCreate';

  const CharacterCreateScreen({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<CharacterCreateScreen> createState() => _CharacterCreateScreenState();
}

class _CharacterCreateScreenState extends State<CharacterCreateScreen> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterCreatePage', 'build');
    log.fine('Building..');

    return Scaffold(
      appBar: header(context, widget.callbacks),
      resizeToAvoidBottomInset: false,
      body: CharacterCreateContainerWidget(
        callbacks: widget.callbacks,
      ),
    );
  }
}
