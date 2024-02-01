-- name: Workflowget :one
select
    id,
    slice,
    app,
    class,
    setname as name,
    is_active,
    is_internal,
    ruleset as flowrules,
    createdat,
    createdby,
    editedat,
    editedby
from ruleset
where
    slice = $1
    and app = $2
    and class = $3
    and setname = $4;

-- -- :one
-- INSERT INTO
--     schema (
--         realm, slice, app, brwf, class, patternschema, actionschema, createdby, editedby
--     )
-- VALUES (
--         1, $1, $2, W, $3, $4, $5, $6, $7
--     )
-- RETURNING
--     id;

-- -- name: SchemaUpdate :one
-- -- :one
-- UPDATE schema
-- SET
--     app = $2,
--     brwf = $3,
--     class = $4,
--     patternschema = $5,
--     actionschema = $6,
--     editedat = CURRENT_TIMESTAMP,
--     editedby = $7
-- WHERE
--     id = $1
-- RETURNING
--     id;

-- -- name: SchemaDelete :one
-- -- :one
-- DELETE FROM schema WHERE id = $1 RETURNING id;

-- -- name: SchemaGet :one
-- -- :one
-- SELECT * FROM schema WHERE id = $1;

-- -- name: SchemaList :many
-- SELECT *
-- FROM schema
-- WHERE
--     slice = $1
--     AND class = $2
--     AND app = $3;

-- -- name: Wfschemaget :one
-- SELECT a.slice, a.app, a.class, b.longname, a.patternschema, a.actionschema, a.createdat, a.createdby, a.editedat, a.editedby
-- FROM schema as a, realm as b, realmslice as c
-- WHERE
--     a.realm = b.id
--     and a.slice = c.id
--     and a.slice = $1
--     and c.realm = b.shortname
--     and a.class = $3
--     AND a.app = $2;

-- -- name: Wfschemadelete :exec
-- DELETE from schema
-- where
--     id in (
--         SELECT a.id
--         FROM schema as a, realm as b, realmslice as c
--         WHERE
--             a.realm = b.id
--             and a.slice = c.id
--             and a.slice = $1
--             and c.realm = b.shortname
--             and a.class = $3
--             AND a.app = $2
--     );