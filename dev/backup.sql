SET sql_safe_updates = false;

DELETE FROM tsg_groups;
DELETE FROM tsg_templates;
DELETE FROM tsg_users;
DELETE FROM tsg_accounts;

INSERT INTO tsg_accounts (id, account_name, triton_uuid, created_at, updated_at) VALUES (332378521158418433, 'joyent', '87307a00-ab96-4fec-8df7-1a256e49fbcc', NOW(), NOW());

INSERT INTO tsg_templates (id, template_name, package, image_id, account_id, instance_name_prefix, firewall_enabled, networks, metadata, userdata, tags, archived) VALUES
    (319209784155176962, 'test-template-1', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 332378521158418433, 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    (319209801539354626, 'test-template-2', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 332378521158418433, 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    (319209812150321154, 'test-template-3', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 332378521158418433, 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    (319209821014392834, 'test-template-4', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 332378521158418433, 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    (319209831565656066, 'test-template-5', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 332378521158418433, 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
    (319209841670782978, 'test-template-6', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 332378521158418433, 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false);

INSERT INTO tsg_groups (id, "name", template_id, account_id, capacity, health_check_interval, archived) VALUES
    (320376470673326082, 'test-group-1', 319209784155176962, 332378521158418433, 3, 300, false),
    (320376528846356482, 'test-group-2', 319209801539354626, 332378521158418433, 3, 300, false),
    (320377354180919298, 'test-group-3', 319209812150321154, 332378521158418433, 6, 60, false),
    (320377452358238210, 'test-group-4', 319209821014392834, 332378521158418433, 1, 300, false),
    (320377513666019330, 'test-group-5', 319209831565656066, 332378521158418433, 3, 120, false),
    (320377641294168066, 'test-group-6', 319209841670782978, 332378521158418433, 12, 300, false);
