package model

import "gitlab.com/alienspaces/go-mud/server/service/game/internal/record"

func (m *Model) performDungeonCharacterAction(
	dungeonActionRecord *record.DungeonAction,
	dungeonLocationRecordSet *DungeonLocationRecordSet,
) (*record.DungeonAction, error) {
	// const logger = this.loggerService.logger({
	// 	function: 'performDungeonCharacterAction',
	// });
	// const actionFuncs = {
	// 	move: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionMove(dungeonActionRecord, records),
	// 	look: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionLook(dungeonActionRecord, records),
	// 	equip: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionEquip(dungeonActionRecord, records),
	// 	stash: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionStash(dungeonActionRecord, records),
	// 	drop: (dungeonActionRecord: DungeonActionRepositoryRecord, records: DungeonLocationRecordSet) =>
	// 		this.performDungeonActionDrop(dungeonActionRecord, records),
	// };

	// const actionFunc = actionFuncs[dungeonActionRecord.resolved_command];
	// if (!actionFunc) {
	// 	throw new DomainError(`Action function for ${dungeonActionRecord.resolved_command} not supported`);
	// }

	// dungeonActionRecord = await actionFunc(dungeonActionRecord, records);

	// logger.info(`Have updated dungeon action record ${dungeonActionRecord}`);

	// return dungeonActionRecord;

	return nil, nil
}

// func (m *Model) performDungeonActionMove(
// 	dungeonActionRecord: DungeonActionRepositoryRecord,
// 	records: DungeonLocationRecordSet,
// ): Promise<DungeonActionRepositoryRecord> {
// 	if (dungeonActionRecord.dungeon_character_id) {
// 		// Move character
// 		let characterRecord = records.character;
// 		characterRecord.dungeon_location_id = dungeonActionRecord.resolved_target_dungeon_location_id;
// 		characterRecord = await this.dungeonCharacterRepository.updateOne({ record: characterRecord });

// 		// Update dungeon action entity
// 		dungeonActionRecord.dungeon_location_id = dungeonActionRecord.resolved_target_dungeon_location_id;
// 	} else if (dungeonActionRecord.dungeon_monster_id) {
// 		// Move monster
// 		// ...
// 	}
// 	return dungeonActionRecord;
// }
