SELECT um.fd_id AS 'fdid', um.station_id AS 'sta', i.incident_number, cs.priority AS 'pri', cs.call_type, GROUP_CONCAT(us.unit) AS units_arrived,
  IFNULL(ss.dispatch_time, MIN(us.dispatch_time)) AS 'dispatched', MIN(us.en_route_time) AS 'enroute', MIN(us.arrived_time) AS 'arrived'
FROM unit_mappings um
LEFT OUTER JOIN incidents i ON i.fd_id = um.fd_id
LEFT OUTER JOIN call_statuses cs ON cs.id = i.call_status_id
LEFT OUTER JOIN unit_statuses us ON cs.id = us.call_status_id AND (us.unit LIKE CONCAT('%', um.station_id) AND NOT us.unit IN ( CONCAT('RES5', um.station_id), CONCAT('RES', um.station_id), CONCAT(um.station_id, 'FAST'), CONCAT(um.station_id, 'TECH'), CONCAT(um.station_id, 'PAID'), CONCAT('STA', um.station_id), CONCAT('STA5', um.station_id) )) AND us.arrived_time <> ''
LEFT OUTER JOIN unit_statuses ss ON cs.id = ss.call_status_id AND ss.unit IN ( CONCAT('STA', um.station_id), CONCAT('RES', um.station_id), CONCAT(um.station_id, 'FAST'), CONCAT(um.station_id, 'TECH'), CONCAT(um.station_id, 'PAID') )
GROUP BY um.fd_id, i.incident_number;

