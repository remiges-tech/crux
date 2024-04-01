-- realm(realm table)

INSERT INTO realm
VALUES (11,
        'BSE',
        'BSE',
        'Bombay Stock Exchange',
        'colty',
        '2022-12-26T09:03:46Z',
        '378');


INSERT INTO realm
VALUES (12,
        'NSE',
        'NSE',
        'National Stock Exchange',
        'umflw',
        '1993-06-05T12:56:23Z',
        '216');


INSERT INTO realm
VALUES (13,
        'MERCE',
        'MERCE',
        'MERCE pvt LTD',
        'iuykj',
        '1978-10-06T12:04:41Z',
        '886');


INSERT INTO realm
VALUES (14,
        'REMIGES',
        'REMIGES',
        'REMIGES TECH PVT LTD',
        'lgjix',
        '1970-11-24T02:08:37Z',
        '1009');

-- realmslice(slice table)

INSERT INTO realmslice
VALUES (11,
        'NSE',
        'Stock Market',
        true,
        NULL,
        NULL,
        '2024-03-01 00:00:00',
        'aniket',
        NULL,
        NULL);


INSERT INTO realmslice
VALUES (12,
        'BSE',
        'Stock Market',
        true,
        NULL,
        NULL,
        '2021-12-01 14:30:15',
        'aniket',
        NULL,
        NULL);


INSERT INTO realmslice
VALUES (13,
        'BSE',
        'Merce Pvt ltd',
        true,
        NULL,
        NULL,
        '2024-03-01 00:00:00',
        'aniket',
        NULL,
        NULL);


INSERT INTO realmslice
VALUES (14,
        'REMIGES',
        'REMIGES Pvt ltd',
        true,
        NULL,
        NULL,
        '2024-03-01 00:00:00',
        'aniket',
        NULL,
        NULL);

-- app

INSERT INTO "app"
VALUES (11,
        'BSE',
        'retailBANK',
        'retailbank',
        'retailbank pvt ltd',
        'admin',
        '2024-01-29 00:00:00');


INSERT INTO "app"
VALUES (12,
        'NSE',
        'retailbank1',
        'retailbank1',
        'retailbank pvt ltd',
        'admin',
        '2024-01-29 00:00:00');


INSERT INTO "app"
VALUES (13,
        'MERCE',
        'nedbank',
        'nedbank',
        'nedbank from canada',
        'admin',
        '2024-01-29 00:00:00');


INSERT INTO "app"
VALUES (14,
        'BSE',
        'nedBank1',
        'nedbank1',
        'netbank pvt ltd',
        'admin',
        '2024-01-29 00:00:00');


INSERT INTO "app"
VALUES (15,
        'BSE',
        'HDFCBank',
        'hdfcbank',
        'hdfcbank pvt ltd',
        'admin',
        '2024-01-29 00:00:00');


INSERT INTO "app"
VALUES (16,
        'BSE',
        'starMF',
        'starmf',
        'mutual Fund',
        'admin',
        '2024-01-29 00:00:00');

INSERT INTO "app"
VALUES (17,
        'BSE',
        'uccapp',
        'uccapp',
        'mutual Fund',
        'admin',
        '2024-01-29 00:00:00');
-- capgrant TABLE

INSERT INTO "capgrant"
VALUES (10,
        'BSE',
        'john',
        'hdfcbank',
        'root',
        '2024-01-29 00:00:00',
        '2024-02-29 00:00:00',
        '2023-12-29 00:00:00',
        'admin');


INSERT INTO "capgrant"
VALUES (12,
        'BSE',
        'Raj',
        'hdfcbank',
        'user',
        '2024-01-29 00:00:00',
        '2024-02-29 00:00:00',
        '2023-12-29 00:00:00',
        'user');


INSERT INTO "capgrant"
VALUES (13,
        'BSE',
        'Raj',
        'nedbank',
        'user',
        '2024-01-29 00:00:00',
        NULL,
        '2023-12-29 00:00:00',
        'user');

INSERT INTO "capgrant"
VALUES (14,
        'BSE',
        'kanchan@gmail.com',
        'nedbank',
        'auth',
        '2024-01-29 00:00:00',
        NULL,
        '2023-12-29 00:00:00',
        'user');

-- config TABLE

INSERT INTO config (realm, slice, name, descr, val, setby)
VALUES ('BSE',
        11,
        'CONFIG_A',
        'Description for CONFIG_A',
        'Value for CONFIG_A',
        'User1'), ('NSE',
                   12,
                   'CONFIG_B',
                   'Description for CONFIG_B',
                   'Value for CONFIG_B',
                   'User2'), ('MERCE',
                              13,
                              'CONFIG_C',
                              'Description for CONFIG_C',
                              'Value for CONFIG_C',
                              'User3');

-- schema(schema table)

INSERT INTO "schema"
VALUES (10,
        'BSE',
        11,
        'retailbank',
        'B',
        'custonboarding',
        '[{"attr": "cat", "valtype": "str"}, {"attr": "mrp", "valtype": "float"}, {"attr": "fullname", "valtype": "str"}, {"attr": "ageinstock", "valtype": "int"}, {"attr": "inventoryqty", "valtype": "int"}]',
        '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}',
        '2022-12-26T09:03:46Z',
        'Mal Houndsom',
        '2023-07-12T01:33:32Z',
        'Clerc Careless');


INSERT INTO "schema"
VALUES (11,
        'NSE',
        12,
        'nedbank',
        'W',
        'custonboarding',
        '[{"attr": "cat", "valtype": "str"}, {"attr": "mrp", "valtype": "float"}, {"attr": "fullname", "valtype": "str"}, {"attr": "ageinstock", "valtype": "int"}, {"attr": "inventoryqty", "valtype": "int"}]',
        '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}',
        '2021-01-03T06:02:41Z',
        'Marielle Strongitharm',
        '2021-06-07T02:28:17Z',
        'Therese Roselli');


insert into "schema"
VALUES (12,
        'BSE',
        12,
        'retailbank',
        'W',
        'inventoryitems',
        '[{"attr": "cat", "valtype": "str"}, {"attr": "mrp", "valtype": "float"}, {"attr": "fullname", "valtype": "str"}, {"attr": "ageinstock", "valtype": "int"}, {"attr": "inventoryqty", "valtype": "int"}]',
        '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');


insert into "schema"
VALUES (13,
        'REMIGES',
        12,
        'retailbank',
        'B',
        'custonboarding',
        '[{"attr": "cat", "valtype": "str"}, {"attr": "mrp", "valtype": "float"}, {"attr": "fullname", "valtype": "str"}, {"attr": "ageinstock", "valtype": "int"}, {"attr": "inventoryqty", "valtype": "int"}]',
        '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}',
        '2023-01-27T12:12:15Z',
        'Adelaide Reape',
        '2023-01-04T22:00:12Z',
        'Imogene Iaduccelli');


insert into "schema"
VALUES (14,
        'BSE',
        11,
        'retailbank',
        'B',
        'temp',
        '[{"attr": "cat", "valtype": "str"}, {"attr": "mrp", "valtype": "float"}, {"attr": "fullname", "valtype": "str"}, {"attr": "ageinstock", "valtype": "int"}, {"attr": "inventoryqty", "valtype": "int"}]',
        '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}',
        '2022-12-24T19:38:52Z',
        'Olly Gerrish',
        '2021-04-28T20:39:09Z',
        'Ronni Matson');


insert into "schema"
VALUES (15,
        'BSE',
        13,
        'retailbank',
        'W',
        'members',
        '[{"attr": "cat", "valtype": "str"}, {"attr": "mrp", "valtype": "float"}, {"attr": "fullname", "valtype": "str"}, {"attr": "ageinstock", "valtype": "int"}, {"attr": "inventoryqty", "valtype": "int"}]',
        '{"tasks":["invitefordiwali","allowretailsale","assigntotrash"],"properties":["discount","shipby"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');


INSERT INTO "schema"
VALUES (16,
        'BSE',
        13,
        'uccapp',
        'W',
        'ucc',
        '[{"attr":"step","vals":{"aof":{},"nomauth":{},"kycvalid":{},"bankaccvalid":{},"getcustdetails":{},"dpandbankaccvalid":{},"sendauthlinktoclient":{}},"valtype":"enum","longdesc":"","shortdesc":""},{"attr":"stepfailed","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"mode","vals":{"demat":{},"physical":{}},"valtype":"enum","longdesc":"","shortdesc":""}]',
        '{"tasks":["getcustdetails","aof","dpandbankaccvalid","kycvalid","nomauth","bankaccvalid","sendauthlinktoclient"],"properties":["nextstep","done"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');


INSERT INTO "schema"
VALUES (17,
        'BSE',
        12,
        'starmf',
        'W',
        'ucc',
        '[{"attr":"step","vals":{"aof":{},"nomauth":{},"kycvalid":{},"bankaccvalid":{},"getcustdetails":{},"dpandbankaccvalid":{},"sendauthlinktoclient":{}},"valtype":"enum","longdesc":"","shortdesc":""},{"attr":"stepfailed","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"mode","vals":{"demat":{},"physical":{}},"valtype":"enum","longdesc":"","shortdesc":""}]',
        '{"tasks":["getcustdetails","aof","dpandbankaccvalid","kycvalid","nomauth","bankaccvalid","sendauthlinktoclient"],"properties":["nextstep","done"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');


INSERT INTO "schema"
VALUES (18,
        'BSE',
        12,
        'starmf',
        'W',
        'ucctest',
        '[{"attr":"step","vals":{"step1":{},"step2":{}},"valtype":"enum","longdesc":"","shortdesc":""},{"attr":"stepfailed","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"mode","vals":{"demat":{},"physical":{}},"valtype":"enum","longdesc":"","shortdesc":""}]',
        '{"tasks":["step1", "step2"],"properties":["nextstep","done"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');

-- ruleset

INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (5,
        'BSE',
        11,
        'retailbank',
        'members',
        'W',
        'goldstatus',
        true,
        true,
        10,
        '[{"ruleactions": {"tasks": ["clearancesale"], "properties": {"shipby": "ups"}}, "rulepattern": [{"op": "eq", "val": "2", "attr": "inventoryqty"}, {"op": "eq", "val": "200", "attr": "mrp"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');


INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (6,
        'BSE',
        14,
        'retailbank',
        'members',
        'W',
        'temp',
        true,
        false,
        13,
        '[{"ruleactions": {"tasks": ["clearancesale"], "properties": {"shipby": "ups"}}, "rulepattern": [{"op": "eq", "val": "2", "attr": "inventoryqty"}, {"op": "eq", "val": "200", "attr": "mrp"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');


INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (7,
        'BSE',
        13,
        'nedbank',
        'calls',
        'W',
        'vip',
        false,
        false,
        13,
        '[{"ruleactions": {"tasks": ["clearancesale"], "properties": {"shipby": "ups"}}, "rulepattern": [{"op": "eq", "val": "2", "attr": "inventoryqty"}, {"op": "eq", "val": "200", "attr": "mrp"}]}]',
        '2024-01-28T00:00:00Z',
        'aniket',
        '2024-01-15T00:00:00Z',
        'tushar');


INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (8,
       'BSE',
        13,
        'uccapp',
        'ucc',
        'W',
        'ucc_user_cr',
        true,
        false,
        17,
        '[{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["getcustdetails"],"properties":{"nextstep":"getcustdetails"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["aof","dpandbankaccvalid","kycvalid","nomauth"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["aof","kycvalid","nomauth","bankaccvalid"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"physical","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"physical","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"kycvalid","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kycvalid","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"dpandbankvalidation","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"dpandbankvalidation","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"bankaccvalid","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"bankaccvalid","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"nomauth","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"nomauth","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"sendauthlinktoclient","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"sendauthlinktoclient","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');


INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (9,
        'BSE',
        12,
        'starmf',
        'ucc',
        'W',
        'ucc_user_cr',
        true,
        false,
        17,
        '[{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["getcustdetails"],"properties":{"nextstep":"getcustdetails"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["aof","dpandbankaccvalid","kycvalid","nomauth"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["aof","kycvalid","nomauth","bankaccvalid"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"physical","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"physical","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"kycvalid","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kycvalid","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"dpandbankvalidation","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"dpandbankvalidation","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"bankaccvalid","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"bankaccvalid","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"nomauth","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"nomauth","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"sendauthlinktoclient","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"sendauthlinktoclient","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');


INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (10,
        'BSE',
        12,
        'starmf',
        'ucctest',
        'W',
        'ucctest',
        true,
        false,
        18,
        '[{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["step1"],"properties":{"nextstep":"step1"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["step2"],"properties":{"nextstep":"step2"}},"rulepattern":[{"op":"eq","val":"step1","attr":"step"},{"op":"eq","val":"demat","attr":"mode"}]},{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"step2","attr":"step"},{"op":"eq","val":"demat","attr":"mode"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');

-- stepworkflow

INSERT INTO stepworkflow
VALUES (12,
        'retailbank',
        'yearendsale',
        'doyearendsalechk');


INSERT INTO stepworkflow
VALUES (12,
        'retailbank',
        'diwalisale',
        'dodiscountcheck');

insert into stepworkflow values (13, 'starmf', 'aof', 'aofworkflow');
insert into stepworkflow values (13, 'starmf', 'dpandbankaccvalid', 'dpandbankaccvalidWorkflow');

--  for test case

INSERT INTO public.wfinstance ("id", "slice", "class", "step", "entityid", "app", "workflow", "loggedat", "nextstep")
VALUES (777777,
        12,
        'inventoryitems',
        'tempstep',
        'tempentityid',
        'retailbank',
        'temp',
        '2024-02-05 00:00:00',
        'temp');


INSERT INTO public.wfinstance ("id", "slice", "class", "step", "entityid", "app", "workflow", "loggedat", "nextstep", "parent")
VALUES (77,
        12,
        'inventoryitems',
        'tempstep',
        'tempentityid',
        'retailbank',
        'temp',
        '2024-02-05 00:00:00',
        'temp',
        78);


INSERT INTO public.wfinstance ("id", "slice", "class", "step", "entityid", "app", "workflow", "loggedat", "nextstep")
VALUES (78,
        12,
        'inventoryitems',
        'tempstep',
        'tempentityid',
        'retailbank',
        'temp',
        '2024-02-05 00:00:00',
        'temp');

---- create above / drop below ----
 -- wfinstance

DELETE
FROM wfinstance;

-- stepworkflow

DELETE
FROM stepworkflow;

-- ruleset

DELETE
FROM ruleset;

-- capgrant

DELETE
FROM "capgrant";

-- schema

DELETE
FROM "schema";

-- config

DELETE
FROM "config";

-- app

DELETE
FROM "app";

-- realmslice

DELETE
FROM realmslice;