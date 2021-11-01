part of 'dungeon_command_cubit.dart';

@immutable
abstract class DungeonCommandState extends Equatable {
  const DungeonCommandState();
}

@immutable
class DungeonCommandStateInitial extends DungeonCommandState {
  const DungeonCommandStateInitial();

  @override
  List<Object> get props => [];
}

@immutable
class DungeonCommandStatePreparing extends DungeonCommandState {
  final String? action;
  final String? target;
  const DungeonCommandStatePreparing({this.action, this.target});

  @override
  List<Object?> get props => [action, target];
}
