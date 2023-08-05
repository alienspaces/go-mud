import 'package:flutter/material.dart';
// import 'package:flutter_bloc/flutter_bloc.dart';

// Application
import 'package:go_mud_client/navigation.dart';
// import 'package:go_mud_client/logger.dart';
// import 'package:go_mud_client/cubit/character/character_cubit.dart';
// import 'package:go_mud_client/cubit/dungeon/dungeon_cubit.dart';

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

// TODO:12-implement-death: Remove the following?

// void _navigateHome(BuildContext context, NavigationCallbacks callbacks) {
//   final log = getLogger('Header', '_navigateHome');
//   log.fine('Navigating to home page...');
//   callbacks.openHomePage(context);
// }

// void _navigateCharacter(BuildContext context, NavigationCallbacks callbacks) {
//   final log = getLogger('Header', '_navigateCharacter');
//   log.fine('Navigating to character page...');
//   callbacks.openCharacterPage(context);
// }

// void _navigateDungeon(BuildContext context, NavigationCallbacks callbacks) {
//   final log = getLogger('Header', '_navigateDungeon');
//   log.fine('Navigating to dungeon list page...');
//   callbacks.openDungeonPage(context);
// }

// Widget _buildLink(
//   BuildContext context,
//   String label,
//   void Function() navigateFunc,
// ) {
//   return Container(
//     padding: const EdgeInsets.fromLTRB(1, 10, 1, 0),
//     child: ElevatedButton(
//       onPressed: navigateFunc,
//       style: ElevatedButton.styleFrom(
//         backgroundColor: Theme.of(context).colorScheme.secondary,
//         foregroundColor: Theme.of(context).colorScheme.onSecondary,
//         disabledForegroundColor: Theme.of(context).colorScheme.onSecondary,
//       ),
//       child: Text(
//         label,
//         style: Theme.of(context).textTheme.labelLarge!.copyWith(
//               fontSize: 14,
//               color: Theme.of(context).colorScheme.onPrimary,
//             ),
//       ),
//     ),
//   );
// }

AppBar header(BuildContext context, NavigationCallbacks callbacks) {
  // final log = getLogger('common', 'header');

  // final characterCubit = BlocProvider.of<CharacterCubit>(context);
  // final dungeonCubit = BlocProvider.of<DungeonCubit>(context);

  // var characterName = characterCubit.characterRecord != null
  //     ? characterCubit.characterRecord!.characterName
  //     : null;
  // log.info("Character record $characterName");

  // var dungeonName = dungeonCubit.dungeonRecord != null
  //     ? dungeonCubit.dungeonRecord!.dungeonName
  //     : null;
  // log.info("Dungeon record $dungeonName");

  List<Widget> links = [];

  // TODO:12-implement-death: Remove the following?

  // links.add(
  //   _buildLink(
  //     context,
  //     'Home',
  //     () => _navigateHome(context, callbacks),
  //   ),
  // );

  // if (characterName != null) {
  //   links.add(
  //     _buildLink(
  //       context,
  //       'Play',
  //       () => _navigateDungeon(context, callbacks),
  //     ),
  //   );
  // } else {
  //   links.add(
  //     _buildLink(
  //       context,
  //       'Play',
  //       () => _navigateCharacter(context, callbacks),
  //     ),
  //   );
  // }

  return AppBar(
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
