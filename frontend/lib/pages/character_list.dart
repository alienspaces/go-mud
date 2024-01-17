import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/page.dart';
import 'package:go_mud_client/widgets/common/header.dart';
import 'package:go_mud_client/widgets/character_list/container.dart';

class CharacterListPage extends Page {
  static const String pageName = 'CharacterListPage';
  final NavigationCallbacks callbacks;

  const CharacterListPage({
    LocalKey key = const ValueKey(CharacterListPage.pageName),
    name = CharacterListPage.pageName,
    required this.callbacks,
  }) : super(key: key, name: name);

  @override
  Route createRoute(BuildContext context) {
    return PageRouteBuilder(
      settings: this,
      pageBuilder: (context, animation, secondaryAnimation) {
        return CharacterListScreen(
          callbacks: callbacks,
        );
      },
      transitionDuration: const Duration(milliseconds: 300),
      transitionsBuilder: pageTransitionsBuilder,
    );
  }
}

class CharacterListScreen extends StatefulWidget {
  final NavigationCallbacks callbacks;
  static String pageName = 'Home';

  const CharacterListScreen({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<CharacterListScreen> createState() => _CharacterListScreenState();
}

class _CharacterListScreenState extends State<CharacterListScreen> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('CharacterListPage', 'build');
    log.fine('Building..');

    return Scaffold(
      appBar: header(context, widget.callbacks),
      resizeToAvoidBottomInset: false,
      body: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        alignment: Alignment.center,
        child: CharacterListContainerWidget(callbacks: widget.callbacks),
      ),
    );
  }
}
