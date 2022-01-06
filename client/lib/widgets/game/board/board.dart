import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/style.dart';
import 'package:go_mud_client/widgets/game/inventory/equipped/equipped.dart';
import 'package:go_mud_client/widgets/game/inventory/stashed/stashed.dart';
import 'package:go_mud_client/widgets/game/location/location.dart';

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
  double gridMemberWidth = 0;
  double gridMemberHeight = 0;

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

        // Set grid member dimensions
        gridMemberWidth = constraints.maxWidth - 2;
        gridMemberHeight = constraints.maxHeight - 2;
        if (gridMemberHeight > gridMemberWidth) {
          gridMemberHeight = gridMemberWidth;
        }
        if (gridMemberWidth > gridMemberHeight) {
          gridMemberWidth = gridMemberHeight;
        }

        log.info(
          '(A-**) Resulting button width $gridMemberWidth height $gridMemberHeight',
        );

        return Container(
          width: gridMemberWidth,
          height: gridMemberHeight,
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

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameBoardWidget');
    log.info('Building..');

    return Row(
      children: <Widget>[
        // Panel buttons
        Expanded(
          flex: 1,
          child: Container(
            padding: const EdgeInsets.all(1),
            margin: const EdgeInsets.all(1),
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
        // Panel
        Expanded(
          flex: 5,
          child: Container(
            decoration: const BoxDecoration(
              color: null,
              border: null,
              borderRadius: BorderRadius.all(Radius.zero),
            ),
            clipBehavior: Clip.antiAlias,
            child: IndexedStack(
              index: panelIndex,
              children: const <Widget>[
                GameLocationWidget(),
                GameInventoryEquippedWidget(),
                GameInventoryStashedWidget(),
              ],
            ),
          ),
        ),
      ],
    );
  }
}
