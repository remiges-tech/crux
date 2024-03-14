-- realm(realm table)
INSERT INTO
    realm
VALUES (
        11, 'BSE', 'BSE', 'Bombay Stock Exchange', 'colty', '2022-12-26T09:03:46Z', '378'
    );

INSERT INTO
    realm
VALUES (
        12, 'NSE', 'NSE', 'National Stock Exchange', 'umflw', '1993-06-05T12:56:23Z', '216'
    );

INSERT INTO
    realm
VALUES (
        13, 'MERCE', 'MERCE', 'MERCE pvt LTD', 'iuykj', '1978-10-06T12:04:41Z', '886'
    );

INSERT INTO
    realm
VALUES (
        14, 'REMIGES', 'REMIGES', 'REMIGES TECH PVT LTD', 'lgjix', '1970-11-24T02:08:37Z', '1009'
    );

-- realmslice(slice table)
INSERT INTO
    realmslice
VALUES (
        11, 'NSE', 'Stock Market', true, NULL, NULL
    );

INSERT INTO
    realmslice
VALUES (
        12, 'BSE', 'Stock Market', true, NULL, NULL
    );

INSERT INTO
    realmslice
VALUES (
        13, 'MERCE', 'Merce Pvt ltd', true, NULL, NULL
    );

INSERT INTO
    realmslice
VALUES (
        14, 'REMIGES', 'REMIGES Pvt ltd', true, NULL, NULL
    );

-- app
INSERT INTO
    "app"
VALUES (
        11, 'BSE', 'retailBANK', 'retailbank', 'retailbank pvt ltd', 'admin', '2024-01-29 00:00:00'
    );

INSERT INTO
    "app"
VALUES (
        12, 'NSE', 'retailbank1', 'retailbank1', 'retailbank pvt ltd', 'admin', '2024-01-29 00:00:00'
    );

INSERT INTO
    "app"
VALUES (
        13, 'MERCE', 'nedbank', 'nedbank', 'nedbank from canada', 'admin', '2024-01-29 00:00:00'
    );

INSERT INTO
    "app"
VALUES (
        14, 'BSE', 'nedBank1', 'nedbank1', 'netbank pvt ltd', 'admin', '2024-01-29 00:00:00'
    );

INSERT INTO
    "app"
VALUES (
        15, 'BSE', 'HDFCBank', 'hdfcbank', 'hdfcbank pvt ltd', 'admin', '2024-01-29 00:00:00'
    );
-- capgrant TABLE
INSERT INTO
    "capgrant"
VALUES (
        1, 11, 'john', 'hdfcbank', 'root', '2024-01-29 00:00:00', '2024-02-29 00:00:00', '2023-12-29 00:00:00', 'admin'
    );

INSERT INTO
    "capgrant"
VALUES (
        2, 11, 'Raj', 'hdfcbank', 'user', '2024-01-29 00:00:00', '2024-02-29 00:00:00', '2023-12-29 00:00:00', 'user'
    );

-- config TABLE
INSERT INTO
    config (
        realm, slice, name, descr, val, setby
    )
VALUES (
        11, 11, 'CONFIG_A', 'Description for CONFIG_A', 'Value for CONFIG_A', 'User1'
    ),
    (
        12, 12, 'CONFIG_B', 'Description for CONFIG_B', 'Value for CONFIG_B', 'User2'
    ),
    (
        13, 13, 'CONFIG_C', 'Description for CONFIG_C', 'Value for CONFIG_C', 'User3'
    );
-- schema(schema table)
INSERT INTO
    "schema"
VALUES (
        10, 11, 11, 'retailbank', 'B', 'custonboarding', '{"class":"custonboarding","attr":[{"name":"step","valtype":"enum","vals":["START","initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"]},{"name":"acctholdertype","valtype":"enum","vals":["individual","joint","corporate","hinduundivided","partnership"]},{"name":"branchtype","valtype":"enum","vals":["urban","semirural","rural"]},{"name":"branchcode","valtype":"str"},{"name":"refererquality","valtype":"int","valmin":0,"valmax":5},{"name":"districtcode","valtype":"int"},{"name":"accttype","valtype":"enum","vals":["savings","current","recurring","fixeddeposit","ppf"]}]}', '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}', '2022-12-26T09:03:46Z', 'Mal Houndsom', '2023-07-12T01:33:32Z', 'Clerc Careless'
    );

INSERT INTO
    "schema"
VALUES (
        11, 12, 12, 'nedbank', 'W', 'custonboarding', '{"class":"custonboarding","attr":[{"name":"step","valtype":"enum","vals":["START","initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"]},{"name":"acctholdertype","valtype":"enum","vals":["individual","joint","corporate","hinduundivided","partnership"]},{"name":"branchtype","valtype":"enum","vals":["urban","semirural","rural"]},{"name":"branchcode","valtype":"str"},{"name":"refererquality","valtype":"int","valmin":0,"valmax":5},{"name":"districtcode","valtype":"int"},{"name":"accttype","valtype":"enum","vals":["savings","current","recurring","fixeddeposit","ppf"]}]}', '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}', '2021-01-03T06:02:41Z', 'Marielle Strongitharm', '2021-06-07T02:28:17Z', 'Therese Roselli'
    );

insert into
    "schema"
VALUES (
        12, 11, 12, 'retailbank', 'W', 'inventoryitems', '{"attr": [{"name": "cat", "vals": ["textbook", "notebook", "stationery", "refbooks"], "valtype": "enum", "enumdesc": ["Text books", "Notebooks", "Stationery and miscellaneous items", "Reference books, library books"], "longdesc": "Each item can belong to one of the following categories: textbooks, notebooks, stationery, or reference books.", "shortdesc": "Category of item"}, {"name": "mrp", "valmax": 20000, "valmin": 0, "valtype": "float", "longdesc": "The maximum retail price of the item in INR as declared by the manufacturer.", "shortdesc": "Maximum retail price"}, {"name": "fullname", "lenmax": 40, "lenmin": 5, "valtype": "str", "longdesc": "The full human-readable name of the item. Not unique, therefore sometimes confusing.", "shortdesc": "Full name of item"}, {"name": "ageinstock", "valmax": 1000, "valmin": 1, "valtype": "int", "longdesc": "The age in days that the oldest sample of this item has been lying in stock", "shortdesc": "Age in stock, in days"}, {"name": "inventoryqty", "valmax": 10000, "valmin": 0, "valtype": "int", "longdesc": "How many of these items are currently present in the inventory", "shortdesc": "Number of items in inventory"}], "class": "inventoryitems"}', '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}', '2020-03-10T12:06:40Z', 'Marigold Sherwin', '2023-10-21T17:39:11Z', 'Brunhilde Bampkin'
    );

insert into
    "schema"
VALUES (
        13, 14, 12, 'retailbank', 'B', 'custonboarding', '{"class":"custonboarding","attr":[{"name":"step","valtype":"enum","vals":["START","initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"]},{"name":"acctholdertype","valtype":"enum","vals":["individual","joint","corporate","hinduundivided","partnership"]},{"name":"branchtype","valtype":"enum","vals":["urban","semirural","rural"]},{"name":"branchcode","valtype":"str"},{"name":"refererquality","valtype":"int","valmin":0,"valmax":5},{"name":"districtcode","valtype":"int"},{"name":"accttype","valtype":"enum","vals":["savings","current","recurring","fixeddeposit","ppf"]}]}', '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}', '2023-01-27T12:12:15Z', 'Adelaide Reape', '2023-01-04T22:00:12Z', 'Imogene Iaduccelli'
    );

insert into
    "schema"
VALUES (
        14, 11, 11, 'retailbank', 'B', 'temp', '{"class":"temp","attr":[{"name":"step","valtype":"enum","vals":["START","initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"]},{"name":"acctholdertype","valtype":"enum","vals":["individual","joint","corporate","hinduundivided","partnership"]},{"name":"branchtype","valtype":"enum","vals":["urban","semirural","rural"]},{"name":"branchcode","valtype":"str"},{"name":"refererquality","valtype":"int","valmin":0,"valmax":5},{"name":"districtcode","valtype":"int"},{"name":"accttype","valtype":"enum","vals":["savings","current","recurring","fixeddeposit","ppf"]}]}', '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}', '2022-12-24T19:38:52Z', 'Olly Gerrish', '2021-04-28T20:39:09Z', 'Ronni Matson'
    );

insert into
    "schema"
VALUES (
        15, 11, 12, 'retailbank', 'W', 'members', '{"attr":[{"name":"cat","vals":["textbook","notebook","stationery","refbooks"],"valtype":"enum","enumdesc":["Text books","Notebooks","Stationery and miscellaneous items","Reference books, library books"],"longdesc":"Each item can belong to one of the following categories: textbooks, notebooks, stationery, or reference books.","shortdesc":"Category of item"},{"name":"mrp","valmax":20000,"valmin":0,"valtype":"float","longdesc":"The maximum retail price of the item in INR as declared by the manufacturer.","shortdesc":"Maximum retail price"},{"name":"fullname","lenmax":40,"lenmin":5,"valtype":"str","longdesc":"The full human-readable name of the item. Not unique, therefore sometimes confusing.","shortdesc":"Full name of item"},{"name":"ageinstock","valmax":1000,"valmin":1,"valtype":"int","longdesc":"The age in days that the oldest sample of this item has been lying in stock","shortdesc":"Age in stock, in days"},{"name":"inventoryqty","valmax":10000,"valmin":0,"valtype":"int","longdesc":"How many of these items are currently present in the inventory","shortdesc":"Number of items in inventory"}],"class":"members"}', '{"tasks":["invitefordiwali","allowretailsale","assigntotrash"],"properties":["discount","shipby"]}', '2020-03-10T12:06:40Z', 'Marigold Sherwin', '2023-10-21T17:39:11Z', 'Brunhilde Bampkin'
    );

INSERT INTO
    "schema"
VALUES (
        16, 12, 12, 'retailbank', 'B', 'retailcustomer', '{"class":"retailcustomer","attr":[{"name":"step","valtype":"enum","vals":["START","initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"]},{"name":"acctholdertype","valtype":"enum","vals":["individual","joint","corporate","hinduundivided","partnership"]},{"name":"branchtype","valtype":"enum","vals":["urban","semirural","rural"]},{"name":"branchcode","valtype":"str"},{"name":"refererquality","valtype":"int","valmin":0,"valmax":5},{"name":"districtcode","valtype":"int"},{"name":"accttype","valtype":"enum","vals":["savings","current","recurring","fixeddeposit","ppf"]}]}', '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}', '2020-03-10T12:06:40Z', 'Marigold Sherwin', '2023-10-21T17:39:11Z', 'Brunhilde Bampkin'
    );

-- ruleset
INSERT INTO
    ruleset (
        id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby
    )
VALUES (
        5, 11, 11, 'retailbank', 'members', 'W', 'goldstatus', true, true, 10, '[{"rulepattern":[{"op":"eq","val":"initialdoc","attr":"step"},{"op":"eq","val":"rural","attr":"branchtype"},{"op":"eq","val":"savings","attr":"accttype"}],"ruleactions":{"tasks":["aadhaarcheck"],"properties":[{"val":"aadhaarcheck","name":"nextstep"}]}},{"rulepattern":[{"op":"eq","val":"initialdoc","attr":"step"},{"op":"eq","val":"semirural","attr":"branchtype"},{"op":"ne","val":"ppf","attr":"accttype"}],"ruleactions":{"tasks":["creditbureauchk","bankdetails","panchk"],"properties":[{"val":"creditbureauchk","name":"nextstep"}]}}]', '2024-01-28T00:00:00Z', 'admin', '2024-01-15T00:00:00Z', 'admin'
    );

INSERT INTO
    ruleset (
        id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby
    )
VALUES (
        6, 11, 12, 'retailbank', 'members', 'W', 'temp', true, false, 13, '[{"rulepattern":[{"op":"eq","val":"initialdoc","attr":"step"},{"op":"eq","val":"rural","attr":"branchtype"},{"op":"eq","val":"savings","attr":"accttype"}],"ruleactions":{"tasks":["aadhaarcheck"],"properties":[{"val":"aadhaarcheck","name":"nextstep"}]}},{"rulepattern":[{"op":"eq","val":"initialdoc","attr":"step"},{"op":"eq","val":"semirural","attr":"branchtype"},{"op":"ne","val":"ppf","attr":"accttype"}],"ruleactions":{"tasks":["creditbureauchk","bankdetails","panchk"],"properties":[{"val":"creditbureauchk","name":"nextstep"}]}}]', '2024-01-28T00:00:00Z', 'admin', '2024-01-15T00:00:00Z', 'admin'


);

INSERT INTO
    ruleset (
        id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby
    )
VALUES (
        7, 11, 13, 'nedbank', 'calls', 'W', 'vip', false, false, 13, '[{"rulepattern":[{"op":"eq","val":"initialdoc","attr":"step"},{"op":"eq","val":"rural","attr":"branchtype"},{"op":"eq","val":"savings","attr":"accttype"}],"ruleactions":{"tasks":["aadhaarcheck"],"properties":[{"val":"aadhaarcheck","name":"nextstep"}]}},{"rulepattern":[{"op":"eq","val":"initialdoc","attr":"step"},{"op":"eq","val":"semirural","attr":"branchtype"},{"op":"ne","val":"ppf","attr":"accttype"}],"ruleactions":{"tasks":["creditbureauchk","bankdetails","panchk"],"properties":[{"val":"creditbureauchk","name":"nextstep"}]}}]', '2024-01-28T00:00:00Z', 'aniket', '2024-01-15T00:00:00Z', 'tushar'
    );

INSERT INTO
    ruleset (
        id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby
    )
VALUES (
        8, 11, 12, 'retailbank', 'inventoryitems', 'W', 'discountcheck', true, false, 13, '[{"rulepattern":[{"op":"eq","val":"initialdoc","attr":"step"},{"op":"eq","val":"rural","attr":"branchtype"},{"op":"eq","val":"savings","attr":"accttype"}],"ruleactions":{"tasks":["aadhaarcheck"],"properties":[{"val":"aadhaarcheck","name":"nextstep"}]}},{"rulepattern":[{"op":"eq","val":"initialdoc","attr":"step"},{"op":"eq","val":"semirural","attr":"branchtype"},{"op":"ne","val":"ppf","attr":"accttype"}],"ruleactions":{"tasks":["creditbureauchk","bankdetails","panchk"],"properties":[{"val":"creditbureauchk","name":"nextstep"}]}}]', '2024-01-28T00:00:00Z', 'admin', '2024-01-15T00:00:00Z', 'admin'
    );
-- stepworkflow
INSERT INTO
    stepworkflow
VALUES (
        12, 'retailbank', 'yearendsale', 'doyearendsalechk'
    );

INSERT INTO
    stepworkflow
VALUES (
        12, 'retailbank', 'diwalisale', 'dodiscountcheck'
    );

--  for test case
INSERT INTO
    public.wfinstance (
        "id", "slice", "class", "step", "entityid", "app", "workflow", "loggedat", "nextstep"
    )
VALUES (
        777777, 12, 'inventoryitems', 'tempstep', 'tempentityid', 'retailbank', 'temp', '2024-02-05 00:00:00', 'temp'
    );

INSERT INTO
    public.wfinstance (
        "id", "slice", "class", "step", "entityid", "app", "workflow", "loggedat", "nextstep", "parent"
    )
VALUES (
        77, 12, 'inventoryitems', 'tempstep', 'tempentityid', 'retailbank', 'temp', '2024-02-05 00:00:00', 'temp', 78
    );

INSERT INTO
    public.wfinstance (
        "id", "slice", "class", "step", "entityid", "app", "workflow", "loggedat", "nextstep"
    )
VALUES (
        78, 12, 'inventoryitems', 'tempstep', 'tempentityid', 'retailbank', 'temp', '2024-02-05 00:00:00', 'temp'
    );

---- create above / drop below ----

-- wfinstance
DELETE FROM wfinstance;

-- stepworkflow
DELETE FROM stepworkflow;

-- ruleset
DELETE FROM ruleset;

-- capgrant
DELETE FROM "capgrant";
-- schema
DELETE FROM "schema";
-- config
DELETE FROM "config";
-- app
DELETE FROM "app";

-- realmslice
DELETE FROM realmslice;

-- realm
DELETE FROM realm;