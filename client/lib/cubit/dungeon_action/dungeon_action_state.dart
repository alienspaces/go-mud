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
  const DungeonActionStateCreating({required this.sentence});

  @override
  List<Object> get props => [sentence];
}

@immutable
class DungeonActionStateCreated extends DungeonActionState {
  final DungeonActionRecord? dungeonActionRecord;

  const DungeonActionStateCreated({this.dungeonActionRecord});

  @override
  List<Object?> get props => [dungeonActionRecord];
}

@immutable
class DungeonActionStatePlaying extends DungeonActionState {
  final DungeonActionRecord previous;
  final DungeonActionRecord current;

  const DungeonActionStatePlaying({required this.previous, required this.current});

  @override
  List<Object?> get props => [previous, current];
}
