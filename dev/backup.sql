SET sql_safe_updates = false;

DELETE FROM tsg_groups;
DELETE FROM tsg_templates;
DELETE FROM tsg_keys;
DELETE FROM tsg_users;
DELETE FROM tsg_accounts;

INSERT INTO tsg_keys (id, name, fingerprint, material, created_at, updated_at)
VALUES ('1d32f239-81e2-4e35-a258-a5649dc4e6f3', 'TSG_Management', '5a:ce:1e:1d:b0:96:78:c6:7a:f2:f8:26:e1:b3:55:79', 'just a test ssh private key yo', NOW(), NOW());

INSERT INTO tsg_accounts (id, account_name, triton_uuid, key_id, created_at, updated_at)
VALUES ('6f873d02-172c-418f-8416-4da2b50d5c53', 'joyent', '87307a00-ab96-4fec-8df7-1a256e49fbcc', '1d32f239-81e2-4e35-a258-a5649dc4e6f3', NOW(), NOW());

INSERT INTO tsg_templates (id, template_name, package, image_id, account_id, firewall_enabled, networks, metadata, userdata, tags, archived) VALUES
    ('ad74301e-ad62-404a-be44-3b2f24d082ac', 'test-template-1', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', '6f873d02-172c-418f-8416-4da2b50d5c53', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    ('f1ead2a9-92fc-4435-9eb8-9e520bc3e4f9', 'test-template-2', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', '6f873d02-172c-418f-8416-4da2b50d5c53', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    ('93a4a267-498c-4911-a463-196eca9a5d99', 'test-template-3', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', '6f873d02-172c-418f-8416-4da2b50d5c53', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    ('deee2b55-11ef-4ffa-b34d-d1035ae1943b', 'test-template-4', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', '6f873d02-172c-418f-8416-4da2b50d5c53', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    ('8b5a6001-8a59-4d85-bc72-1af83015b2c2', 'test-template-5', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', '6f873d02-172c-418f-8416-4da2b50d5c53', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    ('437c560d-b1a9-4dae-b3b3-6dbabb7d23a7', 'test-template-6', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', '6f873d02-172c-418f-8416-4da2b50d5c53', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false);

INSERT INTO tsg_groups (id, "name", template_id, account_id, capacity, health_check_interval, archived) VALUES
    ('9e075e5d-60d5-4cff-968e-b70db0badc12', 'test-group-1', 'ad74301e-ad62-404a-be44-3b2f24d082ac', '6f873d02-172c-418f-8416-4da2b50d5c53', 3, 300, false),
    ('77135218-9e49-4ef7-81da-09de9ec580ff', 'test-group-2', 'f1ead2a9-92fc-4435-9eb8-9e520bc3e4f9', '6f873d02-172c-418f-8416-4da2b50d5c53', 3, 300, false),
    ('95fb339f-8f8d-4184-ac5f-2c57e838136e', 'test-group-3', '93a4a267-498c-4911-a463-196eca9a5d99', '6f873d02-172c-418f-8416-4da2b50d5c53', 6, 60, false),
    ('9dd64b6a-3d2c-4e92-b280-ef0f6cf49d04', 'test-group-4', 'deee2b55-11ef-4ffa-b34d-d1035ae1943b', '6f873d02-172c-418f-8416-4da2b50d5c53', 1, 300, false),
    ('e6e5bc41-204f-4729-8189-ae3ec6385e95', 'test-group-5', '8b5a6001-8a59-4d85-bc72-1af83015b2c2', '6f873d02-172c-418f-8416-4da2b50d5c53', 3, 120, false),
    ('1398d5b9-5750-4ed8-af0d-fe328f5d8cd0', 'test-group-6', '437c560d-b1a9-4dae-b3b3-6dbabb7d23a7', '6f873d02-172c-418f-8416-4da2b50d5c53', 12, 300, false);
