--
-- actions ordered by serial_number
--
select 
    a.id, a.serial_number, a.turn_number,
    li.id  as "location_id",
    l.name as "location_name",
    ci.id  as "character_id",
    c.name as "character_name",
    mi.id  as "monster_id",
    m.name as "monster_name"
from action a
join location_instance li on li.id = a.location_instance_id
join location l on l.id = li.location_id
left join character_instance ci on ci.id = a.character_instance_id
left join character c on c.id = ci.character_id
left join monster_instance mi on mi.id = a.monster_instance_id
left join monster m on m.id = mi.monster_id
order by a.serial_number desc
;
