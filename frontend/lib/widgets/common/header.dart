import 'package:flutter/material.dart';

// Application
import 'package:go_mud_client/navigation.dart';

// ignore: unused_element
void _showDialogue(
    BuildContext context, String content, void Function() continueFunc) {
  Widget cancelButton = ElevatedButton(
    child: const Text("Cancel"),
    onPressed: () {
      Navigator.of(context).pop();
    },
  );

  Widget continueButton = ElevatedButton(
    child: const Text("Continue"),
    onPressed: () {
      Navigator.of(context).pop();
      continueFunc();
    },
  );

  AlertDialog alert = AlertDialog(
    content: Text(content),
    actions: [
      cancelButton,
      continueButton,
    ],
  );

  showDialog(
    context: context,
    useRootNavigator: false,
    builder: (BuildContext context) {
      return alert;
    },
  );
}

AppBar header(BuildContext context, NavigationCallbacks callbacks) {
  List<Widget> links = [];

  return AppBar(
    shape: const Border(
      top: BorderSide(color: Colors.black),
      bottom: BorderSide(color: Colors.black),
    ),
    // ignore: avoid_unnecessary_containers
    title: Container(
      child: Column(
        children: <Widget>[
          Text(
            "Dungeon",
            style: Theme.of(context).textTheme.titleLarge!.copyWith(
                  color: Theme.of(context).colorScheme.onPrimary,
                  fontSize: 16,
                ),
          ),
          Text(
            "Doom",
            style: Theme.of(context).textTheme.titleLarge!.copyWith(
                  color: Theme.of(context).colorScheme.onPrimary,
                  fontSize: 16,
                ),
          ),
        ],
      ),
    ),
    actions: links,
  );
}
