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
  final String? action;
  final String? target;
  const DungeonActionStatePreparing({this.action, this.target});

  @override
  List<Object?> get props => [action, target];
}

@immutable
class DungeonActionStateCreating extends DungeonActionState {
  final String sentence;
  final DungeonActionRecord? current;
  const DungeonActionStateCreating({required this.sentence, this.current});

  @override
  List<Object?> get props => [sentence, current];
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
  final DungeonActionRecord current;
  final String action;
  final String? direction;

  const DungeonActionStateCreated(
      {required this.current, required this.action, this.direction});

  @override
  List<Object?> get props => [current, action, direction];
}

@immutable
class DungeonActionStatePlaying extends DungeonActionState {
  final DungeonActionRecord previous;
  final DungeonActionRecord current;
  final String action;
  final String? direction;

  const DungeonActionStatePlaying({
    required this.previous,
    required this.current,
    required this.action,
    this.direction,
  });

  @override
  List<Object?> get props => [
        previous,
        current,
        action,
        direction,
      ];
}
