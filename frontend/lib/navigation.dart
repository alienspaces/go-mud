import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';

// Application page packages
import 'package:go_mud_client/pages/home.dart';

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

  NavigationCallbacks({
    required this.openHomePage,
  });
}

class _NavigationState extends State<Navigation> {
  // List of supported pages
  List<String> _pageList = [HomePage.pageName];

  // Callback functions set the desired page stack
  void openHomePage(BuildContext context) {
    final log = getLogger('Navigation', 'openHomePage');
    log.fine('Opening home page..');
    setState(() {
      _pageList = [HomePage.pageName];
    });
  }

  List<Page<dynamic>> _pages(BuildContext context) {
    final log = getLogger('Navigation', '_pages');
    log.fine('Building pages..');

    List<Page<dynamic>> pages = [];

    NavigationCallbacks callbacks = NavigationCallbacks(
      openHomePage: openHomePage,
    );

    for (var pageName in _pageList) {
      switch (pageName) {
        case HomePage.pageName:
          log.fine('Adding ${HomePage.pageName}');
          pages.add(HomePage(callbacks: callbacks));
          break;
        default:
        //
      }
    }
    return pages;
  }

  bool _onPopPage(Route<dynamic> route, dynamic result, BuildContext context) {
    final log = getLogger('Navigation', '_onPopPage');
    log.fine('Page name ${route.settings.name}');

    if (!route.didPop(result)) {
      return false;
    }

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
