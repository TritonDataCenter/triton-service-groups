
INSERT INTO tsg_templates (id, "template_name", package, image_id, account_id, machine_name_prefix, firewall_enabled, networks, metadata, userdata, tags, archived) VALUES
	(319209784155176962, 'test-template-1', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 'joyent', 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
	(319209801539354626, 'test-template-2', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 'joyent', 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
	(319209812150321154, 'test-template-3', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 'joyent', 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
	(319209821014392834, 'test-template-4', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 'joyent', 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
	(319209831565656066, 'test-template-5', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 'joyent', 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false),
	(319209841670782978, 'test-template-6', 'test-package', '49b22aec-0c8a-11e6-8807-a3eb4db576ba', 'joyent', 'sample-', false, 'f7ed95d3-faaf-43ef-9346-15644403b963', NULL, 'bash script here', NULL, false);

/*INSERT INTO tsg_groups (id, "name", template_id, account_id, capacity, datacenter, health_check_interval, instance_tags, archived) VALUES
	(320376470673326082, 'test-group-1', 319209784155176962, 'joyent', 3, 'us-sw-1,us-east-1', 300, NULL, false),
	(320376528846356482, 'test-group-2', 319209801539354626, 'joyent', 3, 'us-west-1,us-east-2', 300, NULL, false),
	(320377354180919298, 'test-group-3', 319209812150321154, 'joyent', 6, 'us-east-1', 60, NULL, false),
	(320377452358238210, 'test-group-4', 319209821014392834, 'joyent', 1, 'us-east-1,eu-ams-1', 300, NULL, false),
	(320377513666019330, 'test-group-5', 319209831565656066, 'joyent', 3, 'eu-ams-1', 120, NULL, false),
	(320377641294168066, 'test-group-6', 319209841670782978, 'joyent', 12, 'us-east-3', 300, NULL, false);
*/