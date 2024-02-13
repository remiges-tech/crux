-- realm(realm table)
INSERT INTO
    public.realm
VALUES (
        1, 'BSE', 'BSE', 'Bombay Stock Exchange', 'colty', '2022-12-26T09:03:46Z', '378'
    );

INSERT INTO
    public.realm
VALUES (
        2, 'NSE', 'NSE', 'National Stock Exchange', 'umflw', '1993-06-05T12:56:23Z', '216'
    );

INSERT INTO
    public.realm
VALUES (
        3, 'MERCE', 'MERCE', 'MERCE pvt LTD', 'iuykj', '1978-10-06T12:04:41Z', '886'
    );

INSERT INTO
    public.realm
VALUES (
        4, 'REMIGES', 'REMIGES', 'REMIGES TECH PVT LTD', 'lgjix', '1970-11-24T02:08:37Z', '1009'
    );

-- realmslice(slice table)
INSERT INTO
    public.realmslice
VALUES (
        1, 'BSE', 'Stock Market', true, NULL, NULL
    );

INSERT INTO
    public.realmslice
VALUES (
        2, 'NSE', 'Stock Market', true, NULL, NULL
    );

INSERT INTO
    public.realmslice
VALUES (
        3, 'MERCE', 'Merce Pvt ltd', true, NULL, NULL
    );

-- app
INSERT INTO
    public.app
VALUES (
        1, 'BSE', 'retailBANK', 'retailBANK', 'retailbank pvt ltd', 'admin', '2024-01-29 00:00:00'
    );

INSERT INTO
    public.app
VALUES (
        2, 'NSE', 'retailbank', 'retailbank', 'retailbank pvt ltd', 'admin', '2024-01-29 00:00:00'
    );

INSERT INTO
    public.app
VALUES (
        3, 'MERCE', 'nedbank', 'nedbank', 'nedbank from canada', 'admin', '2024-01-29 00:00:00'
    );

-- schema(schema table)
INSERT INTO
    public.schema
VALUES (
        10, 1, 1, 'retailbank', 'B', 'custonboarding', '{"class":"inventoryitems","attr":[{"name":"cat","valtype":"enum","vals":["textbook","notebook","stationery","refbooks"]},{"name":"mrp","valtype":"float"},{"name":"fullname","valtype":"str"},{"name":"ageinstock","valtype":"int"},{"name":"inventoryqty","valtype":"int"}]}', '{"tasks":["invitefordiwali","allowretailsale","assigntotrash"],"properties":["discount","shipby"]}', '2022-12-26T09:03:46Z', 'Mal Houndsom', '2023-07-12T01:33:32Z', 'Clerc Careless'
    );

INSERT INTO
    public.schema
VALUES (
        11, 2, 2, 'nedbank', 'W', 'custonboarding', '{"class":"inventoryitems","attr":[{"name":"cat","valtype":"enum","vals":["textbook","notebook","stationery","refbooks"]},{"name":"mrp","valtype":"float"},{"name":"fullname","valtype":"str"},{"name":"ageinstock","valtype":"int"},{"name":"inventoryqty","valtype":"int"}]}', '{"tasks":["invitefordiwali","allowretailsale","assigntotrash"],"properties":["discount","shipby"]}', '2021-01-03T06:02:41Z', 'Marielle Strongitharm', '2021-06-07T02:28:17Z', 'Therese Roselli'
    );

insert into
    public.schema
VALUES (
        12, 1, 2, 'retailBANK', 'W', 'inventoryitems', '{"attr":[{"name":"cat","vals":["textbook","notebook","stationery","refbooks"],"valtype":"enum","enumdesc":["Text books","Notebooks","Stationery and miscellaneous items","Reference books, library books"],"longdesc":"Each item can belong to one of the following categories: textbooks, notebooks, stationery, or reference books.","shortdesc":"Category of item"},{"name":"mrp","valmax":20000,"valmin":0,"valtype":"float","longdesc":"The maximum retail price of the item in INR as declared by the manufacturer.","shortdesc":"Maximum retail price"},{"name":"fullname","lenmax":40,"lenmin":5,"valtype":"str","longdesc":"The full human-readable name of the item. Not unique, therefore sometimes confusing.","shortdesc":"Full name of item"},{"name":"ageinstock","valmax":1000,"valmin":1,"valtype":"int","longdesc":"The age in days that the oldest sample of this item has been lying in stock","shortdesc":"Age in stock, in days"},{"name":"inventoryqty","valmax":10000,"valmin":0,"valtype":"int","longdesc":"How many of these items are currently present in the inventory","shortdesc":"Number of items in inventory"}],"class":"inventoryitems"}', '{"tasks":["invitefordiwali","allowretailsale","assigntotrash"],"properties":["discount","shipby"]}', '2022-12-26T09:03:46Z', 'Mal Houndsom', '2023-07-12T01:33:32Z', 'Clerc Careless'
    );

insert into
    public.schema
VALUES (
        13, 4, 2, 'retailBANK', 'B', 'custonboarding', '{"class":"inventoryitems","attr":[{"name":"cat","valtype":"enum","vals":["textbook","notebook","stationery","refbooks"]},{"name":"mrp","valtype":"float"},{"name":"fullname","valtype":"str"},{"name":"ageinstock","valtype":"int"},{"name":"inventoryqty","valtype":"int"}]}', '{"tasks":["invitefordiwali","allowretailsale","assigntotrash"],"properties":["discount","shipby"]}', '2023-01-27T12:12:15Z', 'Adelaide Reape', '2023-01-04T22:00:12Z', 'Imogene Iaduccelli'
    );

insert into
    public.schema
VALUES (
        14, 1, 1, 'retailBANK', 'B', 'tempclass', '{"class":"inventoryitems","attr":[{"name":"cat","valtype":"enum","vals":["textbook","notebook","stationery","refbooks"]},{"name":"mrp","valtype":"float"},{"name":"fullname","valtype":"str"},{"name":"ageinstock","valtype":"int"},{"name":"inventoryqty","valtype":"int"}]}', '{"tasks":["invitefordiwali","allowretailsale","assigntotrash"],"properties":["discount","shipby"]}', '2022-12-24T19:38:52Z', 'Olly Gerrish', '2021-04-28T20:39:09Z', 'Ronni Matson'
    );

insert into
    public.schema
VALUES (
        15, 2, 2, 'retailBANK', 'W', 'members', '{"class":"members","attr":[{"name":"cat","valtype":"enum","vals":["textbook","notebook","stationery","refbooks"]},{"name":"mrp","valtype":"float"},{"name":"fullname","valtype":"str"},{"name":"ageinstock","valtype":"int"},{"name":"inventoryqty","valtype":"int"}]}', '{"tasks":["invitefordiwali","allowretailsale","assigntotrash"],"properties":["discount","shipby"]}', '2020-03-10T12:06:40Z', 'Marigold Sherwin', '2023-10-21T17:39:11Z', 'Brunhilde Bampkin'
    );

INSERT INTO
    public.schema
VALUES (
        16, 2, 2, 'retailBANK', 'B', 'retailcustomer', '{"class":"retailcustomer","attr":[{"name":"step","valtype":"enum","vals":["START","initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"]},{"name":"acctholdertype","valtype":"enum","vals":["individual","joint","corporate","hinduundivided","partnership"]},{"name":"branchtype","valtype":"enum","vals":["urban","semirural","rural"]},{"name":"branchcode","valtype":"str"},{"name":"refererquality","valtype":"int","valmin":0,"valmax":5},{"name":"districtcode","valtype":"int"},{"name":"accttype","valtype":"enum","vals":["savings","current","recurring","fixeddeposit","ppf"]}]}', '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}', '2020-03-10T12:06:40Z', 'Marigold Sherwin', '2023-10-21T17:39:11Z', 'Brunhilde Bampkin'
    );

-- ruleset
INSERT INTO
    ruleset (
        id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby
    )
VALUES (
        5, 1, 2, 'retailbank', 'members', 'W', 'goldstatus', true, true, 10, '{
            "name": "step",
            "type": "enum",
            "vals": [ "START", "initialdoc", "aadhaarcheck", "creditbureauchk", "panchk", "bankdetails", "referencechk", "stage2done", "complete" ],
            "descr": "Current step completed"
        }', '2024-01-28T00:00:00Z', 'admin', '2024-01-15T00:00:00Z', 'admin'
    );

INSERT INTO
    ruleset (
        id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby
    )
VALUES (
        6, 1, 2, 'retailBANK', 'members', 'W', 'temp_set', true, false, 13, '{
            "name": "step",
            "type": "enum1",
            "vals": [ "START", "initialdoc", "aadhaarcheck", "creditbureauchk", "panchk", "bankdetails", "referencechk", "stage2done", "complete" ],
            "descr": "Current step completed"
        }', '2024-01-28T00:00:00Z', 'admin', '2024-01-15T00:00:00Z', 'admin'
    );

INSERT INTO
    ruleset (
        id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby
    )
VALUES (
        7, 1, 3, 'nedbank', 'calls', 'W', 'vip', true, true, 13, '{
            "name": "step",
            "type": "enum1",
            "vals": [ "START", "initialdoc", "aadhaarcheck", "creditbureauchk", "panchk", "bankdetails", "referencechk", "stage2done", "complete" ],
            "descr": "Current step completed"
        }', '2024-01-28T00:00:00Z', 'aniket', '2024-01-15T00:00:00Z', 'tushar'
    );

INSERT INTO
    ruleset (
        id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby
    )
VALUES (
        8, 1, 2, 'retailBANK', 'inventoryitems', 'W', 'discountcheck', true, false, 13, '{
            "name": "step",
            "type": "enum1",
            "vals": [ "START", "initialdoc", "aadhaarcheck", "creditbureauchk", "panchk", "bankdetails", "referencechk", "stage2done", "complete" ],
            "descr": "Current step completed"
        }', '2024-01-28T00:00:00Z', 'admin', '2024-01-15T00:00:00Z', 'admin'
    );

-- stepworkflow
INSERT INTO
    public.stepworkflow
VALUES (
        2, 'retailBANK', 'yearendsale', 'doyearendsalechk'
    );