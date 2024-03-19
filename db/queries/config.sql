-- name: ConfigSet :exec
INSERT INTO
    config(
        realm, slice, name, descr, val, setby
    )
VALUES (
        $1, $2, $3, $4, $5, $6
    );