-- name: ActivateRecord :exec
DELETE FROM deactivated WHERE realm = @realm and "user"= @userId;

-- name: DeactivateRecord :exec
INSERT INTO deactivated (realm , "user", deactby, deactat)
VALUES (@realm,@userId,@deactby,@deactat);