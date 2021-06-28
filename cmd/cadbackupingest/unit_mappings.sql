DROP TABLE IF EXISTS unit_mappings;

CREATE TABLE unit_mappings (
	fd_id	VARCHAR(10),
	station_id CHAR(3),
	station_name VARCHAR(50),
	UNIQUE( fd_id, station_id )
);

INSERT INTO unit_mappings VALUES
( '04020', '93', 'Canterbury' ),
( '04030', '06', 'Chaplin' ),
( '04040', '61', 'Danielson' ),
( '04041', '62', 'Attawaugan' ),
( '04042', '63', 'Dayville' ),
( '04043', '64', 'East Killingly' ),
( '04044', '65', 'South Killingly' ),
( '04045', '60', 'Williamsville' ),
( '04050', '92', 'East Brooklyn' ),
( '04051', '90', 'Mortlake' ),
( '04060', '71', 'Eastford' ),
( '04070', '12', 'Hampton' ),
( '04080', '95', 'Plainfield' ),
( '04081', '97', 'Central Village' ),
( '04082', '94', 'Moosup' ),
( '04083', '96', 'Atwood Hose' ),
( '04090', '70', 'Pomfret' ),
( '04101', '79', 'East Putnam' ),
( '04110', '16', 'Scotland' ),
( '04120', '67', 'Sterling' ),
( '04121', '68', 'Oneco' ),
( '04130', '84', 'Thompson Hill' ),
( '04131', '81', 'Community' ),
( '04132', '85', 'East Thompson' ),
( '04133', '83', 'Quinebaug' ),
( '04134', '82', 'West Thompson' ),
( '04150', '76', 'Woodstock' ),
( '04151', '77', 'Bungay' ),
( '04152', '75', 'Muddy Brook' ),
( '08010', '24', 'Baltic' ),
( '08020', '26', 'Bozrah' ),
( '08040', '56', 'Jewett City' ),
( '08060', '28', 'Colchester' ),
( '08090', '54', 'Lisbon' )
;

## Create derivative mapping

DROP TABLE IF EXISTS unit_station_mappings;

CREATE TABLE unit_station_mappings (
	fd_id VARCHAR(10),
	unit VARCHAR(32),
	UNIQUE( fd_id, unit )
);

INSERT INTO unit_station_mappings
SELECT m.fd_id, s.unit FROM unit_mappings m
LEFT OUTER JOIN unit_statuses s ON
  (s.unit LIKE CONCAT('%', m.station_id) AND NOT s.unit IN ( CONCAT('RES5', m.station_id), CONCAT('RES', m.station_id), CONCAT(m.station_id, 'FAST'), CONCAT(m.station_id, 'TECH'), CONCAT(m.station_id, 'PAID'), CONCAT('STA', m.station_id), CONCAT('STA5', m.station_id) ))
GROUP BY m.fd_id, s.unit
ORDER BY m.fd_id, s.unit;

