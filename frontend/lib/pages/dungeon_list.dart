import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/page.dart';
import 'package:go_mud_client/widgets/common/header.dart';
import 'package:go_mud_client/widgets/dungeon_list/container.dart';

class DungeonListPage extends Page {
  static const String pageName = 'DungeonListPage';
  final NavigationCallbacks callbacks;

  const DungeonListPage({
    LocalKey key = const ValueKey(DungeonListPage.pageName),
    name = DungeonListPage.pageName,
    required this.callbacks,
  }) : super(key: key, name: name);

  @override
  Route createRoute(BuildContext context) {
    return PageRouteBuilder(
      settings: this,
      pageBuilder: (context, animation, secondaryAnimation) {
        return DungeonListScreen(
          callbacks: callbacks,
        );
      },
      transitionDuration: const Duration(milliseconds: 300),
      transitionsBuilder: pageTransitionsBuilder,
    );
  }
}

class DungeonListScreen extends StatefulWidget {
  final NavigationCallbacks callbacks;
  static String pageName = 'Home';

  const DungeonListScreen({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  State<DungeonListScreen> createState() => _DungeonListScreenState();
}

class _DungeonListScreenState extends State<DungeonListScreen> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('DungeonListPage', 'build');
    log.fine('Building..');

    return Scaffold(
      appBar: header(context, widget.callbacks),
      resizeToAvoidBottomInset: false,
      body: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        alignment: Alignment.center,
        child: DungeonListContainerWidget(callbacks: widget.callbacks),
      ),
    );
  }
}
