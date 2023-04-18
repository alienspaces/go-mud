part of 'dungeon_action_cubit.dart';

@immutable
abstract class DungeonActionState extends Equatable {
  const DungeonActionState();
}

@immutable
class DungeonActionStateInitial extends DungeonActionState {
  const DungeonActionStateInitial();

  @override
  List<Object> get props => [];
}

@immutable
class DungeonActionStatePreparing extends DungeonActionState {
  final String? actionCommand;
  final String? target;
  const DungeonActionStatePreparing({this.actionCommand, this.target});

  @override
  List<Object?> get props => [actionCommand, target];
}

@immutable
class DungeonActionStateCreating extends DungeonActionState {
  final String sentence;
  const DungeonActionStateCreating({required this.sentence});

  @override
  List<Object?> get props => [sentence];
}

@immutable
class DungeonActionStateError extends DungeonActionState {
  final String message;

  const DungeonActionStateError({required this.message});

  @override
  List<Object?> get props => [message];
}

@immutable
class DungeonActionStateCreated extends DungeonActionState {
  final DungeonActionRecord action;
  final DungeonActionRecord? previousAction;
  final String actionCommand;
  final String? direction;

  const DungeonActionStateCreated({
    required this.action,
    this.previousAction,
    required this.actionCommand,
    this.direction,
  });

  @override
  List<Object?> get props => [action, previousAction, actionCommand, direction];
}

@immutable
class DungeonActionStatePlaying extends DungeonActionState {
  final DungeonActionRecord currentActionRec;
  final DungeonActionRecord? previousActionRec;
  final String actionCommand;
  final String? actionDirection;

  const DungeonActionStatePlaying({
    required this.currentActionRec,
    this.previousActionRec,
    required this.actionCommand,
    this.actionDirection,
  });

  @override
  List<Object?> get props => [
        currentActionRec,
        previousActionRec,
        actionCommand,
        actionDirection,
      ];
}

@immutable
class DungeonActionStatePlayingOther extends DungeonActionState {
  final DungeonActionRecord actionRec;
  final String actionCommand;
  final String? actionDirection;

  const DungeonActionStatePlayingOther({
    required this.actionRec,
    required this.actionCommand,
    this.actionDirection,
  });

  @override
  List<Object?> get props => [
        actionRec,
        actionCommand,
        actionDirection,
      ];
}
