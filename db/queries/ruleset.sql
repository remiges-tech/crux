-- name: GetApp :one
SELECT app
FROM ruleset
WHERE
    slice = $1
    AND app = $2
    AND class = $3
    AND realm = $4
    AND brwf = 'W';

-- name: GetClass :one
SELECT class
FROM ruleset
WHERE
    slice = $1
    AND app = $2
    AND class = $3
    AND realm = $4
    AND brwf = 'W';

-- name: GetWFActiveStatus :one
SELECT is_active
FROM ruleset
WHERE
    slice = $1
    AND app = $2
    AND class = $3
    AND realm = $4
    AND brwf = 'W'
    AND setname = $5;

-- name: GetWFInternalStatus :one
SELECT is_internal
FROM ruleset
WHERE
    slice = $1
    AND app = $2
    AND class = $3
    AND realm = $4
    AND brwf = 'W'
    AND setname = $5;

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
    and setname = $4
    and realm = @realm
    AND brwf = @brwf;

-- name: WorkFlowNew :one
INSERT INTO
    ruleset (
        realm, slice, app, brwf, class, setname, schemaid, is_active, is_internal, ruleset, createdat, createdby
    )
VALUES (
        @realm_name::varchar,
        (SELECT realmslice.id FROM realmslice WHERE realmslice.id= @slice AND realmslice.realm = @realm_name ),
        (SELECT app.shortnamelc FROM app WHERE app.shortnamelc= @app AND app.realm = @realm_name), 
        $1, $2, $3, $4, false, $5, $6, CURRENT_TIMESTAMP, $7
    )RETURNING id;

-- name: WorkFlowUpdate :execresult
UPDATE ruleset
SET
    ruleset = $4,
    editedat = CURRENT_TIMESTAMP,
    editedby = $5
WHERE
    realm = @realm_name::varchar
    AND slice = (SELECT realmslice.id FROM realmslice WHERE realmslice.id= @slice AND realmslice.realm = @realm_name )
    AND class = $1
    AND app = (SELECT app.shortnamelc FROM app WHERE app.shortnamelc= @app AND app.realm = @realm_name)
    AND brwf = $2
    AND setname = $3
    AND is_active = false;


-- name: WorkflowList :many
select
    id,
    slice,
    app,
    class,
    setname as name,
    is_active,
    is_internal,
    createdat,
    createdby,
    editedat,
    editedby
from ruleset
where
    realm = @realm
    AND (sqlc.narg('slice')::INTEGER is null OR slice = sqlc.narg('slice')::INTEGER)
    AND (sqlc.narg('app')::text[] is null OR app = any( sqlc.narg('app')::text[]))
    AND (sqlc.narg('class')::text is null OR class = sqlc.narg('class')::text)
    AND (sqlc.narg('setname')::text is null OR setname = sqlc.narg('setname')::text)
    AND (sqlc.narg('is_active')::BOOLEAN is null OR is_active = sqlc.narg('is_active')::BOOLEAN)
    AND (sqlc.narg('is_internal')::BOOLEAN is null OR is_internal = sqlc.narg('is_internal')::BOOLEAN)
    and brwf =  @brwf;

-- name: WorkflowDelete :execresult
DELETE from ruleset
where
    brwf = @brwf
    AND is_active = false
    and slice = $1
    and app = $2
    and class = $3
    and setname = $4
    AND realm = $5;

-- name: RulesetRowLock :one
SELECT * 
FROM ruleset 
WHERE
     realm = @realm_name::varchar
    AND slice = (SELECT realmslice.id FROM realmslice WHERE realmslice.id= @slice AND realmslice.realm = @realm_name )
    AND class = $1
    AND app = (SELECT app.shortnamelc FROM app WHERE app.shortnamelc= @app AND app.realm = @realm_name)
FOR UPDATE; 

-- name: AllRuleset :many
SELECT
    *
FROM
    public.ruleset;

-- name: ruleExists :one
select 1 
from ruleset 
where realm = $1
AND app= $2
AND slice= $3 
AND class = $4;

-- name: LoadRuleSet :one
SELECT * 
FROM ruleset 
WHERE
     realm = @realm_name::varchar
    AND slice = (SELECT realmslice.id FROM realmslice WHERE realmslice.id= @slice AND realmslice.realm = @realm_name )
    AND class = $1
    AND app = (SELECT app.shortnamelc FROM app WHERE app.shortnamelc= @app AND app.realm = @realm_name)
    AND setname = $2; 

-- name: IsWorkflowReferringSchema :one
select count(*)
From ruleset
Where realm = $1
AND slice = $2
AND app = $3
AND class = $4
AND is_active = true;

-- name: ActivateBRERuleSet :exec
UPDATE ruleset
SET is_active = true
WHERE realm = @realm
AND slice = @slice
AND app = @app
AND class = @class
AND setname = @setname
AND brwf = @brwf;

-- name: DeActivateBRERuleSet :exec
UPDATE ruleset
SET is_active = false
WHERE realm = @realm
AND slice = @slice
AND app = @app
AND class = @class
AND setname = @setname
AND brwf = @brwf;


-- name: GetBRERuleSetCount :one
SELECT count(*) FROM ruleset
WHERE realm = @realm
AND slice = @slice
AND app = @app
AND class = @class
AND setname = @setname
AND brwf = @brwf;

-- name: GetBRERuleSetActiveStatus :one
SELECT ruleset, is_active ,setname FROM ruleset
WHERE realm = @realm
AND slice = @slice
AND app = @app
AND class = @class
AND setname = @setname
AND brwf = @brwf;


