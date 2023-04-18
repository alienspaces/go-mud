--
-- action character (occupant)
--
select 
    ac.id, ac.action_id, ac.record_type,
    l.name,
    c.name
from action_character ac
join location_instance li on li.id = ac.location_instance_id
join location l on l.id = li.location_id
join character_instance ci on ci.id = ac.character_instance_id
join character c on c.id = ci.character_id
-- where ac.record_type = 'occupant'
;

