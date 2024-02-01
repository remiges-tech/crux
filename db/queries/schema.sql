
-- name: Wfschemaget :one
SELECT a.slice, a.app, a.class, b.longname, a.patternschema, a.actionschema, a.createdat, a.createdby, a.editedat, a.editedby
FROM schema as a, realm as b, realmslice as c
WHERE
    a.realm = b.id
    and a.slice = c.id
    and a.slice = $1
    and c.realm = b.shortname
    and a.class = $3
    AND a.app = $2;

-- name: Wfschemadelete :exec
DELETE from schema
where
    id in (
        SELECT a.id
        FROM schema as a, realm as b, realmslice as c
        WHERE
            a.realm = b.id
            and a.slice = c.id
            and a.slice = $1
            and c.realm = b.shortname
            and a.class = $3
            AND a.app = $2
    );