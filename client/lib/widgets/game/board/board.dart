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

class _GameBoardWidgetState extends State<GameBoardWidget> {
  int panelIndex = 0;

  @override
  Widget build(BuildContext context) {
    final log = getLogger('GameBoardWidget');
    log.info('Building..');

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
                      panelIndex = 0;
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
                      panelIndex = 1;
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
                      panelIndex = 2;
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
