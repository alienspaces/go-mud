part of 'dungeon_character_cubit.dart';

@immutable
abstract class DungeonCharacterState extends Equatable {
  const DungeonCharacterState();
}

@immutable
class DungeonCharacterStateInitial extends DungeonCharacterState {
  const DungeonCharacterStateInitial();

  @override
  List<Object> get props => [];
}

@immutable
class DungeonCharacterStateLoading extends DungeonCharacterState {
  final String characterID;
  const DungeonCharacterStateLoading({required this.characterID});

  @override
  List<Object> get props => [characterID];
}

@immutable
class DungeonCharacterStateLoadError extends DungeonCharacterState {
  final String characterID;
  final String message;
  const DungeonCharacterStateLoadError(
      {required this.characterID, required this.message});

  @override
  List<Object> get props => [characterID];
}

@immutable
class DungeonCharacterStateLoaded extends DungeonCharacterState {
  final DungeonCharacterRecord dungeonCharacterRecord;
  const DungeonCharacterStateLoaded({required this.dungeonCharacterRecord});

  @override
  List<Object> get props => [dungeonCharacterRecord];
}

@immutable
class DungeonCharacterStateCreate extends DungeonCharacterState {
  const DungeonCharacterStateCreate();

  @override
  List<Object> get props => [];
}

@immutable
class DungeonCharacterStateCreateError extends DungeonCharacterState {
  final String dungeonID;
  final String characterID;
  final String message;
  const DungeonCharacterStateCreateError(
      {required this.dungeonID,
      required this.characterID,
      required this.message});

  @override
  List<Object> get props => [dungeonID, characterID, message];
}

@immutable
class DungeonCharacterStateCreated extends DungeonCharacterState {
  final DungeonCharacterRecord dungeonCharacterRecord;
  const DungeonCharacterStateCreated({required this.dungeonCharacterRecord});

  @override
  List<Object> get props => [dungeonCharacterRecord];
}

@immutable
class DungeonCharacterStateDelete extends DungeonCharacterState {
  const DungeonCharacterStateDelete();

  @override
  List<Object> get props => [];
}

@immutable
class DungeonCharacterStateDeleteError extends DungeonCharacterState {
  final String dungeonID;
  final String characterID;
  final String message;
  const DungeonCharacterStateDeleteError(
      {required this.dungeonID,
      required this.characterID,
      required this.message});

  @override
  List<Object> get props => [dungeonID, characterID, message];
}

@immutable
class DungeonCharacterStateDeleted extends DungeonCharacterState {
  final DungeonCharacterRecord dungeonCharacterRecord;
  const DungeonCharacterStateDeleted({required this.dungeonCharacterRecord});

  @override
  List<Object> get props => [dungeonCharacterRecord];
}
