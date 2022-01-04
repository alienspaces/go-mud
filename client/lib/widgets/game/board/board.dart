import 'package:flutter/material.dart';

// Application packages
import 'package:go_mud_client/logger.dart';
import 'package:go_mud_client/widgets/game/inventory/equipped/equipped.dart';
import 'package:go_mud_client/widgets/game/inventory/stashed/stashed.dart';
import 'package:go_mud_client/widgets/game/location/location.dart';

class GameBoardWidget extends StatefulWidget {
  const GameBoardWidget({Key? key}) : super(key: key);

  @override
  _GameBoardWidgetState createState() => _GameBoardWidgetState();
}

enum ShowPanel { location, stashed, equipped }

class _GameBoardWidgetState extends State<GameBoardWidget> {
  ShowPanel showPanel = ShowPanel.location;

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameBoardWidget');
    log.info('Building..');

    late Widget panelWidget;
    if (showPanel == ShowPanel.location) {
      panelWidget = const GameLocationWidget();
    } else if (showPanel == ShowPanel.equipped) {
      panelWidget = const GameInventoryEquippedWidget();
    } else if (showPanel == ShowPanel.stashed) {
      panelWidget = const GameInventoryStashedWidget();
    }

    return Row(
      children: <Widget>[
        // Panel buttons
        Expanded(
          flex: 1,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: <Widget>[
              Container(
                color: Colors.blue,
                margin: const EdgeInsets.fromLTRB(5, 5, 5, 5),
                child: ElevatedButton(
                  onPressed: () {
                    setState(() {
                      showPanel = ShowPanel.location;
                    });
                  },
                  child: const Text('Location'),
                ),
              ),
              Container(
                color: Colors.blue[100],
                margin: const EdgeInsets.fromLTRB(5, 5, 5, 5),
                child: ElevatedButton(
                  onPressed: () {
                    setState(() {
                      showPanel = ShowPanel.equipped;
                    });
                  },
                  child: const Text('Equipped'),
                ),
              ),
              Container(
                color: Colors.blue[200],
                margin: const EdgeInsets.fromLTRB(5, 5, 5, 5),
                child: ElevatedButton(
                  onPressed: () {
                    setState(() {
                      showPanel = ShowPanel.stashed;
                    });
                  },
                  child: const Text('Stashed'),
                ),
              ),
            ],
          ),
        ),
        // Panel
        Expanded(
          flex: 7,
          child: Container(
            decoration: BoxDecoration(color: Colors.orange[100]),
            clipBehavior: Clip.antiAlias,
            child: panelWidget,
          ),
        ),
      ],
    );
  }
}
