import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/common/header.dart';
import 'package:go_mud_client/widgets/home/home_container.dart';

class HomePage extends Page {
  static const String pageName = 'HomePage';
  final NavigationCallbacks callbacks;

  const HomePage({
    LocalKey key = const ValueKey(HomePage.pageName),
    name = HomePage.pageName,
    required this.callbacks,
  }) : super(key: key, name: name);

  @override
  Route createRoute(BuildContext context) {
    return PageRouteBuilder(
      settings: this,
      pageBuilder: (context, animation, secondaryAnimation) {
        return HomeScreen(
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

class HomeScreen extends StatefulWidget {
  final NavigationCallbacks callbacks;
  static String pageName = 'Home';

  const HomeScreen({
    Key? key,
    required this.callbacks,
  }) : super(key: key);

  @override
  _HomeScreenState createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  @override
  Widget build(BuildContext context) {
    final log = getLogger('HomeScreen');
    log.fine('Building..');

    return Scaffold(
      appBar: header(context, widget.callbacks),
      resizeToAvoidBottomInset: false,
      body: Container(
        padding: const EdgeInsets.symmetric(vertical: 16),
        alignment: Alignment.center,
        child: HomeContainerWidget(callbacks: widget.callbacks),
      ),
    );
  }
}
