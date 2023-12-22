import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

// Application page packages
import 'package:go_mud_client/pages/character_list.dart';
import 'package:go_mud_client/pages/character_create.dart';
import 'package:go_mud_client/pages/dungeon_list.dart';
import 'package:go_mud_client/pages/game.dart';

final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();

class Navigation extends StatefulWidget {
  const Navigation({Key? key}) : super(key: key);

  @override
  State<Navigation> createState() => _NavigationState();
}

typedef NavigationCallback = void Function(BuildContext context);

// Navigation callbacks are passed down through the widget
// tree to any widgets that need to perform navigation.
class NavigationCallbacks {
  // Add a callback for every page we need to navigate to
  final NavigationCallback openCharacterListPage;
  final NavigationCallback closeCharacterListPage;
  final NavigationCallback openCharacterCreatePage;
  final NavigationCallback closeCharacterCreatePage;
  final NavigationCallback openDungeonListPage;
  final NavigationCallback closeDungeonListPage;
  final NavigationCallback openGamePage;
  final NavigationCallback closeGamePage;

  NavigationCallbacks({
    required this.openCharacterListPage,
    required this.closeCharacterListPage,
    required this.openCharacterCreatePage,
    required this.closeCharacterCreatePage,
    required this.openDungeonListPage,
    required this.closeDungeonListPage,
    required this.openGamePage,
    required this.closeGamePage,
  });
}

class _NavigationState extends State<Navigation> {
  // Initial page
  List<String> _pageList = [CharacterListPage.pageName];

  // Callback functions set the desired page stack
  void openCharacterListPage(BuildContext context) {
    final log = getLogger('Navigation', 'openCharacterListPage');
    log.fine('Opening character list page..');
    setState(() {
      _pageList = [CharacterListPage.pageName];
    });
  }

  void closeCharacterListPage(BuildContext context) {
    final log = getLogger('Navigation', 'closeCharacterListPage');
    log.warning('--- Closing character list page..');
    Navigator.pop(context);
    setState(() {
      _pageList.removeWhere(
        (pageName) => pageName == CharacterListPage.pageName,
      );
    });
  }

  void openCharacterCreatePage(BuildContext context) {
    final log = getLogger('Navigation', 'openCharacterCreatePage');
    log.fine('Opening character create page..');
    setState(() {
      _pageList = [CharacterListPage.pageName, CharacterCreatePage.pageName];
    });
  }

  void closeCharacterCreatePage(BuildContext context) {
    final log = getLogger('Navigation', 'closeCharacterCreatePage');
    log.warning('--- Closing character create page..');
    Navigator.pop(context);
    setState(() {
      _pageList.removeWhere(
        (pageName) => pageName == CharacterCreatePage.pageName,
      );
    });
  }

  void openDungeonListPage(BuildContext context) {
    final log = getLogger('Navigation', 'openDungeonListPage');
    log.fine('Opening dungeon list page..');
    setState(() {
      _pageList = [DungeonListPage.pageName];
    });
  }

  void closeDungeonListPage(BuildContext context) {
    final log = getLogger('Navigation', 'closeDungeonListPage');
    log.warning('--- Closing dungeon list page..');
    Navigator.pop(context);
    setState(() {
      _pageList.removeWhere(
        (pageName) => pageName == DungeonListPage.pageName,
      );
    });
  }

  void openGamePage(BuildContext context) {
    final log = getLogger('Navigation', 'openGamePage');
    log.fine('Opening game page..');
    setState(() {
      _pageList = [GamePage.pageName];
    });
  }

  void closeGamePage(BuildContext context) {
    final log = getLogger('Navigation', 'closeGamePage');
    log.warning('--- Closing game page..');
    Navigator.pop(context);
    setState(() {
      _pageList.removeWhere(
        (pageName) => pageName == DungeonListPage.pageName,
      );
    });
  }

  List<Page<dynamic>> _pages(BuildContext context) {
    final log = getLogger('Navigation', '_pages');
    log.fine('Building pages..');

    List<Page<dynamic>> pages = [];

    NavigationCallbacks callbacks = NavigationCallbacks(
      openCharacterListPage: openCharacterListPage,
      closeCharacterListPage: closeCharacterListPage,
      openCharacterCreatePage: openCharacterCreatePage,
      closeCharacterCreatePage: closeCharacterCreatePage,
      openDungeonListPage: openDungeonListPage,
      closeDungeonListPage: closeDungeonListPage,
      openGamePage: openGamePage,
      closeGamePage: closeGamePage,
    );

    for (var pageName in _pageList) {
      switch (pageName) {
        case CharacterListPage.pageName:
          log.fine('Adding ${CharacterListPage.pageName}');
          pages.add(CharacterListPage(callbacks: callbacks));
          break;
        case CharacterCreatePage.pageName:
          log.fine('Adding ${CharacterCreatePage.pageName}');
          pages.add(CharacterCreatePage(callbacks: callbacks));
          break;
        case DungeonListPage.pageName:
          log.fine('Adding ${DungeonListPage.pageName}');
          pages.add(DungeonListPage(callbacks: callbacks));
          break;
        case GamePage.pageName:
          log.fine('Adding ${GamePage.pageName}');
          pages.add(GamePage(callbacks: callbacks));
          break;
        default:
        //
      }
    }
    return pages;
  }

  bool _onPopPage(Route<dynamic> route, dynamic result, BuildContext context) {
    if (!route.didPop(result)) {
      return false;
    }

    setState(() {
      _pageList.removeWhere(
        (pageName) => pageName == route.settings.name,
      );
    });

    return true;
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('Navigation', 'build');
    log.fine('Building..');
    return Navigator(
      key: navigatorKey,
      pages: _pages(context),
      onPopPage: (route, result) => _onPopPage(route, result, context),
    );
  }
}
