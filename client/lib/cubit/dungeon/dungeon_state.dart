part of 'dungeon_cubit.dart';

@immutable
abstract class DungeonState extends Equatable {
  const DungeonState();
}

@immutable
class DungeonStateInitial extends DungeonState {
  const DungeonStateInitial();

  @override
  List<Object> get props => [];
}

@immutable
class DungeonStateLoading extends DungeonState {
  const DungeonStateLoading();

  @override
  List<Object> get props => [];
}

@immutable
class DungeonStateLoaded extends DungeonState {
  final List<DungeonRecord>? dungeonRecords;
  final DungeonRecord? currentDungeonRecord;

  const DungeonStateLoaded({required this.dungeonRecords, this.currentDungeonRecord});

  @override
  List<Object?> get props => [dungeonRecords, currentDungeonRecord];
}
