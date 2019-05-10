-- Copyright 2019 Smart-Edge.com, Inc. All rights reserved.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

DROP DATABASE IF EXISTS controller_ce;

CREATE DATABASE controller_ce;

USE controller_ce

-- -------------
-- Entity tables
-- -------------

CREATE TABLE nodes (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE apps (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    type VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.type') STORED,
    entity JSON
);

CREATE TABLE vnfs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    type VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.type') STORED,
    entity JSON
);

CREATE TABLE traffic_policies (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE dns_configs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

CREATE TABLE credentials (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    entity JSON
);

-- -------------------
-- Primary join tables
-- -------------------

-- These tables join two entity tables.

-- dns_configs x apps
CREATE TABLE dns_configs_app_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    dns_config_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_config_id') STORED,
    app_id  VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.app_id') STORED,
    entity JSON,
    FOREIGN KEY (dns_config_id) REFERENCES dns_configs(id),
    FOREIGN KEY (app_id) REFERENCES apps(id),
    UNIQUE KEY (dns_config_id, app_id)
);

-- dns_configs x vnfs
CREATE TABLE dns_configs_vnf_aliases (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    dns_config_id  VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_config_id') STORED,
    vnf_id  VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.vnf_id') STORED,
    entity JSON,
    FOREIGN KEY (dns_config_id) REFERENCES dns_configs(id),
    FOREIGN KEY (vnf_id) REFERENCES vnfs(id),
    UNIQUE KEY (dns_config_id, vnf_id)
);

-- nodes x apps
CREATE TABLE nodes_apps (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    node_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.node_id') STORED,
    app_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.app_id') STORED,
    entity JSON,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (app_id) REFERENCES apps(id),
    UNIQUE KEY (node_id, app_id)
);

-- nodes x vnfs
CREATE TABLE nodes_vnfs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    node_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.node_id') STORED,
    vnf_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.vnf_id') STORED,
    entity JSON,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (vnf_id) REFERENCES vnfs(id),
    UNIQUE KEY (node_id, vnf_id)
);

-- nodes x dns_configs
CREATE TABLE nodes_dns_configs (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    node_id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.node_id') STORED,
    dns_config_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.dns_config_id') STORED,
    entity JSON,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (dns_config_id) REFERENCES dns_configs(id),
    UNIQUE KEY (node_id)
);

-- ---------------------
-- Secondary join tables
-- ---------------------

-- These tables join an entity table to a primary join table.

-- nodes_apps x traffic_policies
CREATE TABLE nodes_apps_traffic_policies (
    id VARCHAR(36) GENERATED ALWAYS AS (entity->>'$.id') STORED UNIQUE KEY,
    nodes_apps_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.nodes_apps_id') STORED,
    traffic_policy_id VARCHAR(36) GENERATED ALWAYS AS
        (entity->>'$.traffic_policy_id') STORED,
    entity JSON,
    FOREIGN KEY (nodes_apps_id) REFERENCES nodes_apps(id),
    FOREIGN KEY (traffic_policy_id) REFERENCES traffic_policies(id),
    UNIQUE KEY (nodes_apps_id, traffic_policy_id)
);
