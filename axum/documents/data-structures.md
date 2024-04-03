# 数据结构

## events

```sql
-- crate.events definition

-- Drop table

-- DROP TABLE crate.events;
CREATE TABLE crate.events (
	id int8 NOT NULL,
	relation_id int8 NOT NULL,
	reference_id int8 NOT NULL,
	event_time timestamp NOT NULL,
	tags jsonb NOT NULL,
	details jsonb NOT NULL,
	CONSTRAINT actions_pk PRIMARY KEY (id)
);
CREATE INDEX actions_event_time_idx ON crate.events USING btree (event_time);
CREATE INDEX actions_reference_id_idx ON crate.events USING btree (reference_id);
CREATE INDEX actions_relation_id_idx ON crate.events USING btree (relation_id);
```

## settings

```sql
-- crate.settings definition

-- Drop table

-- DROP TABLE crate.settings;

CREATE TABLE crate.settings (
	id int8 NOT NULL,
	root_id int8 NOT NULL,
	parent_id int8 NOT NULL,
	tags jsonb NOT NULL,
	details jsonb NOT NULL,
	CONSTRAINT settings_pk PRIMARY KEY (id)
);
CREATE INDEX settings_parent_id_idx ON crate.settings USING btree (parent_id);
CREATE INDEX settings_root_id_idx ON crate.settings USING btree (root_id);
```