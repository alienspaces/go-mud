import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';
import 'package:go_mud_client/widgets/game/board/board_equipped.dart';
import 'package:go_mud_client/widgets/game/board/board_stashed.dart';
import 'package:go_mud_client/widgets/game/board/board_location.dart';

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
  _GameBoardWidgetState createState() => _GameBoardWidgetState();
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
    final log = getLogger('buildBoardButton');
    log.info('Building..');

    return LayoutBuilder(
      builder: (BuildContext context, BoxConstraints constraints) {
        log.info(
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
          margin: gameButtonMargin,
          child: ElevatedButton(
            onPressed: () {
              setState(() {
                panelIndex = boardButtonIndexes[boardButtonType]!;
              });
            },
            style: gameButtonStyle,
            child: Text(boardButtonLabels[boardButtonType]!),
          ),
        );
      },
    );
  }

  Widget buildBoardPanel(BuildContext context, {required Widget panel}) {
    final log = getLogger('buildBoardPanel');
    log.info('Building..');

    return LayoutBuilder(
      builder: (BuildContext context, BoxConstraints constraints) {
        log.info(
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
            color: Colors.purple[200],
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
    final log = getLogger('GameBoardWidget');
    log.info('Building..');

    return Row(
      children: <Widget>[
        // Board buttons
        Expanded(
          flex: 1,
          child: Container(
            padding: const EdgeInsets.fromLTRB(3, 0, 3, 0),
            decoration: BoxDecoration(
              color: Colors.purple[100],
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
