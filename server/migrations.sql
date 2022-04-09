-- https://docs.ohmtech.io/uThingVOC/bsec/#iaq-clasification
CREATE TABLE IF NOT EXISTS `measurements` (
  `temperature` INTEGER NOT NULL, -- celsius, ambient temperature estimated outside the device
  `pressure` FLOAT NOT NULL, -- hPa, atmospheric (barometric) pressure (1000 hPA = 1 bar)
  `humidity` FLOAT NOT NULL, -- ambient Relative-Humidity (*) estimated outside the device
  `gas_resistance` INTEGER NOT NULL, -- ohm, hot-plate gas resistance
  `iaq` FLOAT NOT NULL, -- Air Quality Index (check BSEC classification)
  `accuracy` INTEGER NOT NULL, -- 0: not-fixed-yet (IAQ invalid), 3: maximum-stability
  `co2_equivalent` FLOAT NOT NULL, -- ppm, CO2 equivalent estimate
  `voc_estimate` FLOAT NOT NULL, -- ppm, breath-VOC concentration estimate
  `created_at` TEXT NOT NULL -- ISO8601 extended timestamp in UTC. e.g. 2015-03-14T09:26:53Z
);