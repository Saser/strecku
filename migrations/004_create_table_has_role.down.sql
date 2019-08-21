DROP TABLE has_role;

-- `user_role` has to be dropped after `has_role` due to `user_role` being used in `has_role`.
DROP TYPE user_role;
