-- realm(realm table)

INSERT INTO realm
VALUES (11,
        'Nova',
        'Nova',
        'Bombay Stock Exchange',
        'colty',
        '2022-12-26T09:03:46Z',
        '378');


INSERT INTO realm
VALUES (12,
        'Ecommerce',
        'ecommerce',
        'buying and selling goods and services online',
        'Tushar',
        CURRENT_TIMESTAMP,
        '{}');


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
        'Ecommerce',
        'buying and selling goods and services online',
        true,
        CURRENT_TIMESTAMP,
        NULL,
        CURRENT_TIMESTAMP,
        'aniket',
        NULL,
        NULL);


INSERT INTO realmslice
VALUES (12,
        'Nova',
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
        'Nova',
        'Stock Market',
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
        'Nova',
        'retailBANK',
        'retailbank',
        'retailbank pvt ltd',
        'admin',
        '2024-01-29 00:00:00');


INSERT INTO "app"
VALUES (12,
        'Ecommerce',
        'Amazon',
        'amazon',
        'American multinational technology company, engaged in e-commerce',
        'tushar',
        CURRENT_TIMESTAMP);



INSERT INTO "app"
VALUES (14,
        'Nova',
        'nedBank1',
        'nedbank1',
        'netbank pvt ltd',
        'admin',
        '2024-01-29 00:00:00');


INSERT INTO "app"
VALUES (15,
        'Nova',
        'HDFCBank',
        'hdfcbank',
        'hdfcbank pvt ltd',
        'admin',
        '2024-01-29 00:00:00');


INSERT INTO "app"
VALUES (16,
        'Nova',
        'fundify',
        'fundify',
        'mutual Fund',
        'admin',
        '2024-01-29 00:00:00');

INSERT INTO "app"
VALUES (17,
        'Nova',
        'uccapp',
        'uccapp',
        'mutual Fund',
        'admin',
        '2024-01-29 00:00:00');
INSERT INTO "app"
VALUES (18,
        'Ecommerce',
        'Myntra',
        'myntra',
        'American multinational technology company, engaged in e-commerce',
        'kanchan',
        CURRENT_TIMESTAMP);

-- capgrant TABLE

INSERT INTO "capgrant"
VALUES (10,
        'Nova',
        'Raj',
        NULL,
        'root',
        '2024-01-29 00:00:00',
        '2024-02-29 00:00:00',
        '2023-12-29 00:00:00',
        'admin');


INSERT INTO "capgrant"
VALUES (12,
        'Nova',
        'Raj',
        'hdfcbank',
        'rules',
        '2024-01-29 00:00:00',
        '2024-02-29 00:00:00',
        '2023-12-29 00:00:00',
        'user');


INSERT INTO "capgrant"
VALUES (14,
        'Nova',
        'Raj',
        NULL,
        'auth',
        '2024-01-29 00:00:00',
        NULL,
        '2023-12-29 00:00:00',
        'user');


INSERT INTO capgrant (id,realm, "user", app, cap, "from", "to", setby) VALUES
(15,'Nova', 'john_doe', 'fundify', 'read', '2023-01-01', '2023-12-31', 'admin'),
(16,'Nova', 'jane_smith', 'hdfcbank', 'write', '2023-02-15', NULL, 'manager'),
(18,'Nova', 'neha_gupta', 'uccapp', 'admin', '2023-01-01', '2024-01-01', 'admin');
INSERT INTO "capgrant"
VALUES (19,
        'Nova',
        'Raj',
        'fundify',
        'schema',
        '2024-01-29 00:00:00',
        NULL,
        '2023-12-29 00:00:00',
        'user');

INSERT INTO "capgrant"
VALUES (20,
        'Ecommerce',
        'Raj',
        'amazon',
        'schema',
        '2024-01-29 00:00:00',
        NULL,
        '2023-12-29 00:00:00',
        'kanchan');

        INSERT INTO "capgrant"
VALUES (21,
        'Ecommerce',
        'Raj',
        'amazon',
        'ruleset',
        '2024-01-29 00:00:00',
        NULL,
        '2023-12-29 00:00:00',
        'kanchan');


        INSERT INTO "capgrant"
VALUES (
        22,
        'Nova',
        'Raj',
        'amazon',
        'ruleset',
        '2024-01-29 00:00:00',
        NULL,
        '2023-12-29 00:00:00',
        'kanchan');


-- config TABLE

INSERT INTO config (realm, slice, name, descr, val, setby)
VALUES ('Nova',
        11,
        'CONFIG_A',
        'Description for CONFIG_A',
        'Value for CONFIG_A',
        'User1'),

         ('Ecommerce',
        12,
        'CONFIG_B',
        'Description for CONFIG_B',
        'Value for CONFIG_B',
        'User2');

-- schema(schema table)

INSERT INTO "schema"
VALUES (10,
        'Nova',
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
        'Ecommerce',
        11,
        'amazon',
        'B',
        'inventoryitems',
        '[{"attr":"cat","valtype":"enum","vals":{"textbook":{},"notebook":{},"stationery":{},"refbooks":{}}},{"attr":"mrp","shortdesc":"Maximum retail price","longdesc":"The maximum retail price of the item in INR as declared by the manufacturer.","valtype":"float"},{"attr":"fullname","valtype":"str"},{"attr":"ageinstock","valtype":"int"},{"attr":"inventoryqty","valtype":"int"}]',
        '{"tasks":["cat","mrp","fullname","ageinstock","inventoryqty"],"properties":["nextstep","done"]}',
        '2021-01-03T06:02:41Z',
        'Marielle Strongitharm',
        '2021-06-07T02:28:17Z',
        'Tushar');


insert into "schema"
VALUES (12,
        'Nova',
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
        'Nova',
        11,
        'retailbank',
        'B',
        'members',
        '[{"attr": "cat", "valtype": "str"}, {"attr": "mrp", "valtype": "float"}, {"attr": "fullname", "valtype": "str"}, {"attr": "ageinstock", "valtype": "int"}, {"attr": "inventoryqty", "valtype": "int"}]',
        '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}',
        '2022-12-24T19:38:52Z',
        'Olly Gerrish',
        '2021-04-28T20:39:09Z',
        'Ronni Matson');


insert into "schema"
VALUES (15,
        'Nova',
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
        'Nova',
        13,
        'uccapp',
        'W',
        'ucc',
        '[{"attr":"step","vals":{"aof":{},"nomauth":{},"kycvalid":{},"bankaccvalid":{},"getcustdetails":{},"dpandbankaccvalid":{},"auth_done":{},"sendauthlinktoclient":{}},"valtype":"enum","longdesc":"","shortdesc":""},{"attr":"stepfailed","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"mode","vals":{"demat":{},"physical":{}},"valtype":"enum","longdesc":"","shortdesc":""}]',
        '{"tasks":["getcustdetails","aof","dpandbankaccvalid","kycvalid","nomauth","bankaccvalid","auth_done","sendauthlinktoclient"],"properties":["nextstep","done"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');


INSERT INTO "schema"
VALUES (17,
        'Nova',
        12,
        'fundify',
        'W',
        'ucc',
        '[{"attr": "step", "vals": {"aof": {}, "start": {}, "kyc_done": {}, "dp_bank_done": {}, "dp_verification": {}, "kyc_verification": {}, "pan_verification": {}, "bank_verification": {}, "ucc_authentication": {}, "pan_aadhaar_linking": {}, "fataca_ubo_verification": {}, "nomination_authentication": {}}, "valtype": "enum", "longdesc": "", "shortdesc": ""}, {"attr": "stepfailed", "valtype": "bool", "longdesc": "", "shortdesc": ""}, {"attr": "member_type", "vals": {"broker": {}, "non-broker": {}}, "valtype": "enum", "longdesc": "", "shortdesc": ""}, {"attr": "ucc_type", "vals": {"demat": {}, "physical": {}}, "valtype": "enum", "longdesc": "", "shortdesc": ""}, {"attr": "tax_status_type", "vals": {"individual": {}, "non_individual": {}}, "valtype": "enum", "longdesc": "", "shortdesc": ""}]',
        '{"tasks": ["ucc_authentication", "pan_verification", "pan_aadhaar_linking", "kyc_verification", "dp_verification", "bank_verification", "nomination_authentication", "aof", "fataca_ubo_verification", "kyc_done", "dp_bank_done"], "properties": ["nextstep", "done"]}',
        CURRENT_TIMESTAMP,
        'tushar');

INSERT INTO "schema"
VALUES (18,
        'Nova',
        12,
        'fundify',
        'W',
        'ucctest',
        '[{"attr":"step","vals":{"step1":{},"step2":{}},"valtype":"enum","longdesc":"","shortdesc":""},{"attr":"stepfailed","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"mode","vals":{"demat":{},"physical":{}},"valtype":"enum","longdesc":"","shortdesc":""}]',
        '{"tasks":["step1", "step2"],"properties":["nextstep","done"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');

INSERT INTO "schema"
VALUES (19,
        'Nova',
        13,
        'uccapp',
        'W',
        'ucc_aof',
        '[{"attr":"step","vals":{"sendtorta":{},"getsigneddocument":{}},"valtype":"enum","longdesc":"","shortdesc":""},{"attr":"stepfailed","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"aofexists","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"mode","vals":{"demat":{},"physical":{}},"valtype":"enum","longdesc":"","shortdesc":""}]',
        '{"tasks":["getsigneddocument","sendtorta"],"properties":["nextstep","done"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');

        INSERT INTO "schema"
VALUES (20,
        'Nova',
        12,
        'fundify',
        'B',
        'custonboarding',
        '[{"attr": "cat", "valtype": "str"}, {"attr": "mrp", "valtype": "float"}, {"attr": "fullname", "valtype": "str"}, {"attr": "ageinstock", "valtype": "int"}, {"attr": "inventoryqty", "valtype": "int"}]',
        '{"class":"retailcustomer","tasks":["initialdoc","aadhaarcheck","creditbureauchk","panchk","bankdetails","referencechk","stage2done","complete"],"properties":["nextstep","done"]}',
        '2021-01-03T06:02:41Z',
        'Marielle Strongitharm',
        '2021-06-07T02:28:17Z',
        'Therese Roselli');

INSERT INTO "schema"
VALUES (21,
        'Ecommerce',
        11,
        'myntra',
        'B',
        'inventoryitems',
        '[{"attr":"cat","valtype":"enum","vals":{"textbook":{},"notebook":{},"stationery":{},"refbooks":{}}},{"attr":"mrp","shortdesc":"Maximum retail price","longdesc":"The maximum retail price of the item in INR as declared by the manufacturer.","valtype":"float"},{"attr":"fullname","valtype":"str"},{"attr":"ageinstock","valtype":"int"},{"attr":"inventoryqty","valtype":"int"}]',
        '{"tasks":["cat","mrp","fullname","ageinstock","inventoryqty"],"properties":["nextstep","done"]}',
        '2021-01-03T06:02:41Z',
        'Marielle Strongitharm',
        '2021-06-07T02:28:17Z',
        'kanchan');

        INSERT INTO "schema"
VALUES (22,
        'Nova',
        13,
        'amazon',
        'B',
        'ucc_aof',
        '[{"attr":"step","vals":{"sendtorta":{},"getsigneddocument":{}},"valtype":"enum","longdesc":"","shortdesc":""},{"attr":"stepfailed","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"aofexists","valtype":"bool","longdesc":"","shortdesc":""},{"attr":"mode","vals":{"demat":{},"physical":{}},"valtype":"enum","longdesc":"","shortdesc":""}]',
        '{"tasks":["getsigneddocument","sendtorta"],"properties":["nextstep","done"]}',
        '2020-03-10T12:06:40Z',
        'Marigold Sherwin',
        '2023-10-21T17:39:11Z',
        'Brunhilde Bampkin');


-- ruleset

INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (5,
        'Ecommerce',
        11,
        'amazon',
        'inventoryitems',
        'B',
        'amazonruleset',
        false,
        true,
        11,
        '[{"NFailed":0,"NMatched":0,"ruleactions":{"tasks":["christmassale"],"properties":{"shipby":"fedex"}},"rulepattern":[{"op":"eq","val":"textbook","attr":"cat"},{"op":"ge","val":5000,"attr":"mrp"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');


INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (6,
        'Nova',
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
VALUES (8,
       'Nova',
        13,
        'uccapp',
        'ucc',
        'W',
        'ucc_user_cr',
        true,
        false,
        17,
        '[{"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":true,"attr":"stepfailed"}]},{"ruleactions":{"tasks":["getcustdetails"],"properties":{"nextstep":"getcustdetails"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"demat","attr":"mode"}]},{"ruleactions":{"tasks":["aof","dpandbankaccvalid","kycvalid","nomauth"],"properties":{"nextstep":"auth_done"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"ruleactions":{"tasks":["aof","kycvalid","nomauth","bankaccvalid"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"getcustdetails","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"physical","attr":"mode"}]},{"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kycvalid","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"dpandbankaccvalid","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"bankaccvalid","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"nomauth","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"ruleactions":{"tasks":["sendauthlinktoclient"],"properties":{"nextstep":"sendauthlinktoclient"}},"rulepattern":[{"op":"eq","val":"auth_done","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]},{"ruleactions":{"tasks":[],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"sendauthlinktoclient","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":"demat","attr":"mode"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');

INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (11,
       'Nova',
        13,
        'uccapp',
        'ucc_aof',
        'W',
        'aofworkflow',
        true,
        false,
        19,
        '[{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":["getsigneddocument"],"properties":{"nextstep":"getsigneddocument"}},"rulePattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":["sendtorta"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":true,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');


INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby)
VALUES (9,
        'Nova',
        12,
        'fundify',
        'ucc',
        'W',
        'ucc_user_cr',
        true,
        false,
        17,
        '[{"ruleactions":{"tasks":["pan_verification","pan_aadhaar_linking","kyc_verification"],"properties":{"nextstep":"kyc_done"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_aadhaar_linking","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kyc_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["bank_verification"],"properties":{"nextstep":"bank_verification"}},"rulepattern":[{"op":"eq","val":"kyc_done","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["nomination_authentication"],"properties":{"nextstep":"nomination_authentication"}},"rulepattern":[{"op":"eq","val":"bank_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["aof"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"nomination_authentication","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["fataca_ubo_verification"],"properties":{"nextstep":"fataca_ubo_verification"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"fataca_ubo_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["pan_verification","kyc_verification"],"properties":{"nextstep":"kyc_done"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kyc_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["bank_verification"],"properties":{"nextstep":"bank_verification"}},"rulepattern":[{"op":"eq","val":"kyc_done","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["aof"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"bank_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["fataca_ubo_verification"],"properties":{"nextstep":"fataca_ubo_verification"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"fataca_ubo_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["pan_verification","pan_aadhaar_linking","kyc_verification"],"properties":{"nextstep":"kyc_done"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_aadhaar_linking","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kyc_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["bank_verification"],"properties":{"nextstep":"bank_verification"}},"rulepattern":[{"op":"eq","val":"kyc_done","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["nomination_authentication"],"properties":{"nextstep":"nomination_authentication"}},"rulepattern":[{"op":"eq","val":"bank_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"nomination_authentication","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["pan_verification","kyc_verification"],"properties":{"nextstep":"kyc_done"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kyc_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["bank_verification"],"properties":{"nextstep":"bank_verification"}},"rulepattern":[{"op":"eq","val":"kyc_done","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"bank_verification","attr":"step"},{"op":"eq","val":"broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["ucc_authentication"],"properties":{"nextstep":"ucc_authentication"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["pan_verification","pan_aadhaar_linking","kyc_verification"],"properties":{"nextstep":"kyc_done"}},"rulepattern":[{"op":"eq","val":"ucc_authentication","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_aadhaar_linking","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kyc_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["bank_verification"],"properties":{"nextstep":"bank_verification"}},"rulepattern":[{"op":"eq","val":"kyc_done","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["nomination_authentication"],"properties":{"nextstep":"nomination_authentication"}},"rulepattern":[{"op":"eq","val":"bank_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["aof"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"nomination_authentication","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["fataca_ubo_verification"],"properties":{"nextstep":"fataca_ubo_verification"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"fataca_ubo_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["ucc_authentication"],"properties":{"nextstep":"ucc_authentication"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["pan_verification","kyc_verification"],"properties":{"nextstep":"kyc_done"}},"rulepattern":[{"op":"eq","val":"ucc_authentication","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kyc_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["bank_verification"],"properties":{"nextstep":"bank_verification"}},"rulepattern":[{"op":"eq","val":"kyc_done","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["aof"],"properties":{"nextstep":"aof"}},"rulepattern":[{"op":"eq","val":"bank_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["fataca_ubo_verification"],"properties":{"nextstep":"fataca_ubo_verification"}},"rulepattern":[{"op":"eq","val":"aof","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"fataca_ubo_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"physical","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["pan_verification","pan_aadhaar_linking","kyc_verification"],"properties":{"nextstep":"kyc_done"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_aadhaar_linking","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kyc_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["dp_verification","bank_verification"],"properties":{"nextstep":"dp_bank_done"}},"rulepattern":[{"op":"eq","val":"kyc_done","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"dp_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"bank_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["nomination_authentication"],"properties":{"nextstep":"nomination_authentication"}},"rulepattern":[{"op":"eq","val":"dp_bank_done","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"nomination_authentication","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["pan_verification","kyc_verification"],"properties":{"nextstep":"kyc_done"}},"rulepattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"pan_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"kyc_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":["dp_verification","bank_verification"],"properties":{"nextstep":"dp_bank_done"}},"rulepattern":[{"op":"eq","val":"kyc_done","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"dp_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"bank_verification","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]},{"ruleactions":{"tasks":[""],"properties":{"done":"true"}},"rulepattern":[{"op":"eq","val":"dp_bank_done","attr":"step"},{"op":"eq","val":"non-broker","attr":"member_type"},{"op":"eq","val":"demat","attr":"ucc_type"},{"op":"eq","val":"non-individual","attr":"tax_status_type"}]}]',
        CURRENT_TIMESTAMP,
        'tushar');


INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (10,
        'Nova',
        12,
        'fundify',
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

        INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (12,
       'Nova',
        13,
        'uccapp',
        'ucc_aof',
        'W',
        'uccworkflow',
        true,
        false,
        19,
        '[{"ruleActions":{"tasks":["getsigneddocument"],"properties":{"nextstep1":"getsigneddocument"}},"rulePattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":["sendtorta"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":true,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');
  INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (13,
       'Nova',
        13,
        'uccapp',
        'ucc_aof',
        'W',
        'uccdoneworkflow',
        true,
        false,
        19,
        '[{"ruleActions":{"tasks":["getsigneddocument"],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":["sendtorta"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":true,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');

          INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (14,
       'Nova',
        13,
        'uccapp',
        'ucc_aof',
        'W',
        'uccmultiplestepsworkflow',
        true,
        false,
        19,
        '[{"ruleActions":{"tasks":["getsigneddocument","getadhaardocuments"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":["sendtorta"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":true,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');

         INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
VALUES (15,
       'Nova',
        13,
        'amazon',
        'ucc_aof',
        'B',
        'step_one',
        false,
        false,
        22,
        '[{"ruleActions":{"tasks":["getsigneddocument","getadhaardocuments"],"properties":{"nextstep":"sendtorta"},"thencall":"step_two"},"rulePattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":["sendtorta"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":true,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');
         INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
        VALUES (
        16,
       'Nova',
        13,
        'amazon',
        'ucc_aof',
        'B',
        'step_two',
        false,
        false,
        22,
        '[{"ruleActions":{"tasks":["getsigneddocument","getadhaardocuments"],"properties":{"nextstep":"sendtorta"},"thencall":"step_three"},"rulePattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":["sendtorta"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":true,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');
        INSERT INTO ruleset (id, realm, slice, app, class, brwf, setname, is_active, is_internal, schemaid, ruleset, createdat, createdby, editedat, editedby)
        VALUES (
        17,
       'Nova',
        13,
        'amazon',
        'ucc_aof',
        'B',
        'step_three',
        false,
        false,
        22,
        '[{"ruleActions":{"tasks":["getsigneddocument","getadhaardocuments"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"start","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":["sendtorta"],"properties":{"nextstep":"sendtorta"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"},{"op":"eq","val":false,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"getsigneddocument","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"},{"op":"eq","val":true,"attr":"aofexists"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":false,"attr":"stepfailed"}]},{"NFailed":0,"NMatched":0,"ruleActions":{"tasks":[],"properties":{"done":"true"}},"rulePattern":[{"op":"eq","val":"sendtorta","attr":"step"},{"op":"eq","val":true,"attr":"stepfailed"}]}]',
        '2024-01-28T00:00:00Z',
        'admin',
        '2024-01-15T00:00:00Z',
        'admin');
        
-- stepworkflow
INSERT INTO stepworkflow
VALUES (13,
        'uccapp',
        'getsigneddocument',
        'aofworkflow');

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

insert into stepworkflow values (13, 'uccapp', 'aof', 'aofworkflow');
-- insert into stepworkflow values (13, 'fundify', 'dpandbankaccvalid', 'dpandbankaccvalidWorkflow');

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