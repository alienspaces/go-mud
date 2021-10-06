import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

// Application page packages
import 'package:go_mud_client/pages/home.dart';
import 'package:go_mud_client/pages/character.dart';
import 'package:go_mud_client/pages/game.dart';

final GlobalKey<NavigatorState> navigatorKey = GlobalKey<NavigatorState>();

class Navigation extends StatefulWidget {
  const Navigation({Key? key}) : super(key: key);

  @override
  _NavigationState createState() => _NavigationState();
}

typedef NavigationCallback = void Function();

// Navigation callbacks are passed down through the widget
// tree to any widgets that need to perform navigation.
class NavigationCallbacks {
  // Add a callback for every page we need to navigate to
  final NavigationCallback openHomePage;
  final NavigationCallback openCharacterPage;
  final NavigationCallback openGamePage;

  NavigationCallbacks({
    required this.openHomePage,
    required this.openCharacterPage,
    required this.openGamePage,
  });
}

class _NavigationState extends State<Navigation> {
  // List of supported pages
  List<String> _pageList = [HomePage.pageName];

  // Callback functions set the desired page stack
  void openHomePage() {
    final log = getLogger('Navigation');
    log.info('Opening home page..');
    setState(() {
      _pageList = [HomePage.pageName];
    });
  }

  void openCharacterPage() {
    final log = getLogger('Navigation');
    log.info('Opening character page..');
    setState(() {
      _pageList = [CharacterPage.pageName];
    });
  }

  void openGamePage() {
    final log = getLogger('Navigation');
    log.info('Opening game page..');
    setState(() {
      _pageList = [GamePage.pageName];
    });
  }

  List<Page<dynamic>> _pages(BuildContext context) {
    final log = getLogger('Navigation - _pages');
    log.info('Building pages..');

    List<Page<dynamic>> pages = [];

    NavigationCallbacks callbacks = NavigationCallbacks(
      openHomePage: openHomePage,
      openCharacterPage: openCharacterPage,
      openGamePage: openGamePage,
    );

    for (var pageName in _pageList) {
      switch (pageName) {
        case HomePage.pageName:
          log.info('Adding ${HomePage.pageName}');
          pages.add(HomePage(callbacks: callbacks));
          break;
        case CharacterPage.pageName:
          log.info('Adding ${CharacterPage.pageName}');
          pages.add(CharacterPage(callbacks: callbacks));
          break;
        case GamePage.pageName:
          log.info('Adding ${GamePage.pageName}');
          pages.add(GamePage(callbacks: callbacks));
          break;
        default:
        //
      }
    }
    return pages;
  }

  bool _onPopPage(Route<dynamic> route, dynamic result, BuildContext context) {
    final log = getLogger('Navigation - _onPopPage');
    log.info('Page name ${route.settings.name}');

    if (!route.didPop(result)) {
      return false;
    }

    return true;
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('Navigation - build');
    log.info('Building..');
    return Navigator(
      key: navigatorKey,
      pages: _pages(context),
      onPopPage: (route, result) => _onPopPage(route, result, context),
    );
  }
}
