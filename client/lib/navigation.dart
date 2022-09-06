import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/cubit/dungeon_action/dungeon_action_cubit.dart';

// Application page packages
import 'package:go_mud_client/pages/home.dart';
import 'package:go_mud_client/pages/character.dart';
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
  void openHomePage(BuildContext context) {
    final log = getLogger('Navigation');
    log.fine('(openHomePage) Opening home page..');
    setState(() {
      _pageList = [HomePage.pageName];
    });
  }

  void openCharacterPage(BuildContext context) {
    final log = getLogger('Navigation');
    log.fine('(openCharacterPage) Opening character page..');
    setState(() {
      _pageList = [CharacterPage.pageName];
    });
  }

  void openGamePage(BuildContext context) {
    final log = getLogger('Navigation');
    log.fine('(openGamePage) Opening game page..');

    // Clear all dungeon actions
    final dungeonActionCubit = BlocProvider.of<DungeonActionCubit>(context);
    log.fine(
        '(openGamePage) Dungeon action record count ${dungeonActionCubit.dungeonActionRecords.length}');
    dungeonActionCubit.clearActions();

    setState(() {
      _pageList = [GamePage.pageName];
    });
  }

  List<Page<dynamic>> _pages(BuildContext context) {
    final log = getLogger('Navigation - _pages');
    log.fine('Building pages..');

    List<Page<dynamic>> pages = [];

    NavigationCallbacks callbacks = NavigationCallbacks(
      openHomePage: openHomePage,
      openCharacterPage: openCharacterPage,
      openGamePage: openGamePage,
    );

    for (var pageName in _pageList) {
      switch (pageName) {
        case HomePage.pageName:
          log.fine('Adding ${HomePage.pageName}');
          pages.add(HomePage(callbacks: callbacks));
          break;
        case CharacterPage.pageName:
          log.fine('Adding ${CharacterPage.pageName}');
          pages.add(CharacterPage(callbacks: callbacks));
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
    final log = getLogger('Navigation - _onPopPage');
    log.fine('Page name ${route.settings.name}');

    if (!route.didPop(result)) {
      return false;
    }

    return true;
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('Navigation - build');
    log.fine('Building..');
    return Navigator(
      key: navigatorKey,
      pages: _pages(context),
      onPopPage: (route, result) => _onPopPage(route, result, context),
    );
  }
}
