import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/navigation.dart';
import 'package:go_mud_client/logger.dart';

class HomeContainerWidget extends StatefulWidget {
  final NavigationCallbacks callbacks;

  const HomeContainerWidget({Key? key, required this.callbacks})
      : super(key: key);

  @override
  State<HomeContainerWidget> createState() => _HomeContainerWidgetState();
}

class _HomeContainerWidgetState extends State<HomeContainerWidget> {
  @override
  void initState() {
    final log = getLogger('HomeContainerWidget');
    log.fine('Initialising state..');

    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('HomeContainerWidget');
    log.fine('Building..');

    // ignore: avoid_unnecessary_containers
    return Container(
      child: ElevatedButton(
        onPressed: () {
          widget.callbacks.openCharacterPage(context);
        },
        child: const Text(
          'Play',
        ),
      ),
    );
  }
}
