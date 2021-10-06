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
class DungeonActionStateCreating extends DungeonActionState {
  const DungeonActionStateCreating();

  @override
  List<Object> get props => [];
}

@immutable
class DungeonActionStateCreated extends DungeonActionState {
  final DungeonActionRecord? dungeonActionRecord;

  const DungeonActionStateCreated({this.dungeonActionRecord});

  @override
  List<Object?> get props => [dungeonActionRecord];
}
