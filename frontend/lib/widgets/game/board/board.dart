import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';

import 'package:go_mud_client/widgets/game/board/equipped_board.dart';
import 'package:go_mud_client/widgets/game/board/stashed_board.dart';
import 'package:go_mud_client/widgets/game/board/location_board.dart';

enum BoardButtonType { location, equipped, stashed }

Map<BoardButtonType, String> boardButtonLabels = {
  BoardButtonType.location: 'L',
  BoardButtonType.equipped: 'E',
  BoardButtonType.stashed: 'S',
};
Map<BoardButtonType, int> boardButtonIndexes = {
  BoardButtonType.location: 0,
  BoardButtonType.equipped: 1,
  BoardButtonType.stashed: 2,
};

class GameBoardWidget extends StatefulWidget {
  const GameBoardWidget({Key? key}) : super(key: key);

  @override
  State<GameBoardWidget> createState() => _GameBoardWidgetState();
}

class _GameBoardWidgetState extends State<GameBoardWidget> {
  int panelIndex = 0;
  double buttonWidth = 0;
  double buttonHeight = 0;

  double panelWidth = 0;
  double panelHeight = 0;

  Widget buildBoardButton(
    BuildContext context,
    BoardButtonType boardButtonType,
  ) {
    final log = getLogger('GameBoardWidget', 'buildBoardButton');
    log.fine('Building..');

    return LayoutBuilder(
      builder: (BuildContext context, BoxConstraints constraints) {
        log.fine(
          'Building width ${constraints.maxWidth} height ${constraints.maxHeight}',
        );

        // Set button dimensions
        buttonWidth = constraints.maxWidth - 2;
        buttonHeight = constraints.maxHeight - 2;
        if (buttonHeight > buttonWidth) {
          buttonHeight = buttonWidth;
        }
        if (buttonWidth > buttonHeight) {
          buttonWidth = buttonHeight;
        }

        return Container(
          width: buttonWidth,
          height: buttonHeight,
          margin: gameBoardButtonMargin,
          child: ElevatedButton(
            onPressed: () {
              setState(() {
                panelIndex = boardButtonIndexes[boardButtonType]!;
              });
            },
            style: gameBoardButtonStyle,
            child: Text(boardButtonLabels[boardButtonType]!),
          ),
        );
      },
    );
  }

  Widget buildBoardPanel(BuildContext context, {required Widget panel}) {
    final log = getLogger('GameBoardWidget', 'buildBoardPanel');
    log.fine('Building..');

    return LayoutBuilder(
      builder: (BuildContext context, BoxConstraints constraints) {
        log.fine(
          'Building width ${constraints.maxWidth} height ${constraints.maxHeight}',
        );

        // Set panel dimensions
        panelWidth = constraints.maxWidth - 2;
        panelHeight = constraints.maxHeight - 2;
        if (panelHeight > panelWidth) {
          panelHeight = panelWidth;
        }
        if (panelWidth > panelHeight) {
          panelWidth = panelHeight;
        }

        return Container(
          decoration: BoxDecoration(
            color: pageContainerBackgroundColor,
            border: null,
            borderRadius: const BorderRadius.all(Radius.zero),
          ),
          clipBehavior: Clip.antiAlias,
          width: panelWidth,
          height: panelHeight,
          margin: gamePanelMargin,
          child: panel,
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameBoardWidget', 'build');
    log.fine('Building..');

    return Row(
      children: <Widget>[
        // Board buttons
        Expanded(
          flex: 1,
          child: Container(
            padding: const EdgeInsets.fromLTRB(3, 0, 3, 0),
            decoration: BoxDecoration(
              color: pageContainerBackgroundColor,
              border: null,
              borderRadius: const BorderRadius.all(Radius.zero),
            ),
            alignment: Alignment.center,
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: <Widget>[
                buildBoardButton(context, BoardButtonType.location),
                buildBoardButton(context, BoardButtonType.equipped),
                buildBoardButton(context, BoardButtonType.stashed),
              ],
            ),
          ),
        ),
        // Board panels
        Expanded(
          flex: 5,
          child: Container(
            padding: const EdgeInsets.fromLTRB(3, 0, 3, 0),
            child: IndexedStack(
              alignment: Alignment.center,
              index: panelIndex,
              children: <Widget>[
                buildBoardPanel(context, panel: const BoardLocationWidget()),
                buildBoardPanel(context, panel: const BoardEquippedWidget()),
                buildBoardPanel(context, panel: const BoardStashedWidget()),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
